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
	battleTalkQueueLimit = 2
	battleTalkDuration   = 90
)

// resolveBattleTalkIndex 重現 sub_1075B 的 base／variant 索引公式。
func resolveBattleTalkIndex(base, variant int) int {
	if base < 0x196 {
		return base
	}
	return 0x196 + ((base - 0x196) << 3) + variant
}

type battleTalkQueue struct {
	entries   []BattleTalkEntry
	remaining int
}

func (q *battleTalkQueue) enqueue(entry BattleTalkEntry) bool {
	if len(q.entries) >= battleTalkQueueLimit {
		return false
	}
	if entry.Duration <= 0 {
		entry.Duration = battleTalkDuration
	}
	q.entries = append(q.entries, entry)
	if len(q.entries) == 1 {
		q.remaining = entry.Duration
	}
	return true
}

func (q *battleTalkQueue) current() (BattleTalkEntry, bool) {
	if len(q.entries) == 0 {
		return BattleTalkEntry{}, false
	}
	return q.entries[0], true
}

func (q *battleTalkQueue) advance() bool {
	if len(q.entries) == 0 {
		return false
	}
	q.entries = q.entries[1:]
	q.remaining = 0
	if len(q.entries) > 0 {
		q.remaining = q.entries[0].Duration
	}
	return true
}

func (q *battleTalkQueue) tick(frames int) {
	if frames <= 0 || len(q.entries) == 0 {
		return
	}
	q.remaining -= frames
	for q.remaining <= 0 && len(q.entries) > 0 {
		q.advance()
	}
}

type battleTalkSession struct {
	battle      *tactical.Battle
	queue       battleTalkQueue
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
		s.queue.enqueue(entry)
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
	entries := make([]BattleTalkEntry, 0, battleTalkQueueLimit)

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
	if len(s.queue.entries) == 0 {
		return false
	}
	if _, ok := s.queue.current(); !ok {
		return false
	}
	if !pressed(ebiten.KeyEnter) && !pressed(ebiten.KeySpace) &&
		!inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}
	s.queue.advance()
	return true
}

func (g *game) tickBattleTalk(frames int) {
	if g == nil {
		return
	}
	g.battleTalkSession.queue.tick(frames)
}
