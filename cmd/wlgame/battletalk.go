package main

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// BattleTalkEntry 是戰術呈現層的最小 TALK payload。
// 文案不存入 state；Index 仍由目前載入的 TALK.DAT 解碼。
type BattleTalkEntry struct {
	Index    int
	Portrait int
	Side     int
	Duration int
}

const (
	battleTalkSides = 2 // 上下各一個框，對應原版的 word_1D322／word_1D324
	// battleTalkDuration 是對白框掛多久：原版 `sub_1C315` 設的到期時刻是
	// **目前節拍 ＋ 0x3C（60）**，由 `sub_1A12A` 每 tick 比對後擦掉
	// （docs/spec/60）。門強度條走同一個機制但常數是 0x14（20）。
	battleTalkDuration = 60
)

// resolveBattleTalkIndex 重現 sub_1075B 的 base／variant 索引公式。
func resolveBattleTalkIndex(base, variant int) int {
	if base < 0x196 {
		return base
	}
	return 0x196 + ((base - 0x196) << 3) + variant
}

// battleTalkSlots 是**每一側一個框**的對白狀態。
//
// ⭐ 原版就是這樣：`sub_1C315` 依側別（`dl`）把到期時刻寫進
// `word_1D322`（側 0）或 `word_1D324`（側 1），`sub_1A12A` 每 tick
// 各自比對後擦掉（docs/spec/60）。開場那一段你一句我一句，
// **兩個框是同時掛著的**——原版實錄截圖上兩個都在。
type battleTalkSlots struct {
	entry     [2]BattleTalkEntry
	remaining [2]int // 0 ＝ 這一側沒有框
}

// set 掛上（或換掉）某一側的框。新的一句會把到期時刻整個覆寫，不疊加。
func (q *battleTalkSlots) set(entry BattleTalkEntry) bool {
	side := entry.Side
	if side < 0 || side > 1 {
		return false
	}
	if entry.Duration <= 0 {
		entry.Duration = battleTalkDuration
	}
	q.entry[side], q.remaining[side] = entry, entry.Duration
	return true
}

// current 取某一側目前掛著的框。
func (q *battleTalkSlots) current(side int) (BattleTalkEntry, bool) {
	if side < 0 || side > 1 || q.remaining[side] <= 0 {
		return BattleTalkEntry{}, false
	}
	return q.entry[side], true
}

// active 回報有沒有任何一側掛著框。
func (q *battleTalkSlots) active() bool {
	return q.remaining[0] > 0 || q.remaining[1] > 0
}

// clear 收掉某一側的框。
func (q *battleTalkSlots) clear(side int) bool {
	if side < 0 || side > 1 || q.remaining[side] <= 0 {
		return false
	}
	q.remaining[side] = 0
	q.entry[side] = BattleTalkEntry{}
	return true
}

// clearAll 兩側一起收（玩家按鍵推進走這一條）。
func (q *battleTalkSlots) clearAll() bool {
	ok := q.clear(0)
	return q.clear(1) || ok
}

// tick 推進兩側各自的計時器。
func (q *battleTalkSlots) tick(frames int) {
	if frames <= 0 {
		return
	}
	for side := range q.remaining {
		if q.remaining[side] <= 0 {
			continue
		}
		if q.remaining[side] -= frames; q.remaining[side] <= 0 {
			q.clear(side)
		}
	}
}

type battleTalkSession struct {
	battle      *tactical.Battle
	queue       battleTalkSlots
	initialized bool
}

func clearBattleTalkSession(g *game, battle *tactical.Battle) {
	if g != nil && (battle == nil || g.battleTalkSession.battle == battle) {
		g.battleTalkSession = battleTalkSession{}
	}
}

func (g *game) bindBattleTalk(battle *tactical.Battle) {
	if g == nil || g.battleTalkSession.battle == battle {
		return
	}
	g.battleTalkSession = battleTalkSession{battle: battle}
}

func (s *battleTalkSession) initialize(battle *tactical.Battle, entries []BattleTalkEntry) {
	if s.battle != battle {
		*s = battleTalkSession{battle: battle}
	}
	if s.initialized {
		return
	}
	// 即使 entries 全部因 fail-closed 被丟棄，也必須封住本 battle 的初始化。
	s.initialized = true
	for _, entry := range entries {
		s.queue.set(entry)
	}
}

func (g *game) startBattleTalk(p *state.Pending) {
	if g == nil || g.world == nil || p == nil || p.Battle == nil {
		return
	}
	g.bindBattleTalk(p.Battle)
	s := &g.battleTalkSession
	if s.initialized {
		return
	}
	entries := make([]BattleTalkEntry, 0, battleTalkSides)

	// sub_1A3C3 的開戰 pair：第一句使用 0x1BA，第二句使用 0x1BB。
	// 上／下 slot 與攻／守的最終對應仍是強推論；採影片位置的最小接線。
	for _, spec := range []struct {
		base int
		side int
	}{
		{base: 0x1BA, side: 0}, // 攻方 → 上方（強推論）
		{base: 0x1BB, side: 1}, // 守方 → 下方（強推論）
	} {
		commander := g.battleCommander(p, spec.side)
		if commander < 0 || commander >= len(g.world.Generals) {
			continue
		}
		general := g.world.Generals[commander]
		entry := BattleTalkEntry{
			Index:    resolveBattleTalkIndex(spec.base, talkVariant(general.TalkVariant)),
			Portrait: general.Portrait,
			Side:     spec.side,
			Duration: battleTalkDuration,
		}
		// 不能安全代入的 marker 直接丟棄該 entry，不顯示 debug／半句文字。
		if _, ok := g.battleTalkText(entry); !ok {
			continue
		}
		entries = append(entries, entry)
	}
	s.initialize(p.Battle, entries)
}

func (g *game) battleTalkText(entry BattleTalkEntry) (string, bool) {
	if g == nil {
		return "", false
	}
	// 開戰 pair 的未知 payload 尚未有可證實的欄位映射；不猜 marker 值。
	lines, ok := g.talkLines(entry.Index, nil)
	if !ok || len(lines) == 0 {
		return "", false
	}
	text := strings.Join(lines, "\n")
	if text == "" {
		return "", false
	}
	return text, true
}

func (g *game) advanceBattleTalkInput() bool {
	if g == nil {
		return false
	}
	s := &g.battleTalkSession
	if !s.queue.active() {
		return false
	}
	if !pressed(ebiten.KeyEnter) && !pressed(ebiten.KeySpace) &&
		!inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}
	return s.queue.clearAll()
}

func (g *game) tickBattleTalk(frames int) {
	if g == nil {
		return
	}
	g.battleTalkSession.queue.tick(frames)
}
