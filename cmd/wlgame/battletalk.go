package main

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// BattleTalkEntry 是戰術呈現層的最小 TALK payload。
// 文案不存入 state；Index 仍由目前載入的 TALK.DAT 解碼。
type BattleTalkEntry struct {
	// Speaker 是說話武將的名字（Big5 已轉），供 \1 標記代入；
	// 空字串表示這一則沒有名字可代（fail-closed 會把帶 \1 的訊息丟棄）。
	Speaker string
	Index    int
	Portrait int
	Side     int
	Duration int
	// Other 是對手大將的呼び名。參數串流的第一格（docs/spec/136 §2）。
	Other string
}

const (
	battleTalkSides = 2 // 上下各一個框，對應原版的 word_1D322／word_1D324
	// battleTalkDuration 是對白框掛多久：原版 `sub_1C315` 設的到期時刻是
	// **目前節拍 ＋ 0x3C（60）**，由 `sub_1A12A` 每 tick 比對後擦掉
	// （docs/spec/60）。門強度條走同一個機制但常數是 0x14（20）。
	battleTalkDuration = 60
)

// talkVariantGroupBase 是 TALK.DAT 開始「八格一組」的門檻（`sub_1075B`
// 的 `cmp cx, 196h`）。之後的索引是**組編號**，不是訊息本身。
const talkVariantGroupBase = 0x196

// resolveBattleTalkIndex 重現 sub_1075B 的 base／variant 索引公式。
func resolveBattleTalkIndex(base, variant int) int {
	if base < talkVariantGroupBase {
		return base
	}
	if variant < 0 {
		variant = 0
	}
	if variant > 7 {
		variant = 7
	}
	return talkVariantGroupBase + ((base - talkVariantGroupBase) << 3) + variant
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
	if !s.initialized {
		// 開戰喊話由**單挑狀態機**產生（docs/spec/80）：這裡只負責武裝——
		// 閘是戰場編號 0xC0–0xD0（`sub_19A33` 的 `byte_1D34B == 1`），
		// 挑戰／拒戰／應戰、回合互嗆與決著全在 `tactical.stepDuel` 裡。
		in := tactical.DuelInput{
			FieldNumber: g.battle.FieldNumber(p.Node, p.Mode == combat.Siege),
		}
		for side := 0; side < 2; side++ {
			if c := g.battleCommander(p, side); c >= 0 && c < len(g.world.Generals) {
				in.Martial[side] = g.world.Generals[c].Martial
				in.CommandStat[side] = g.world.Generals[c].Command
			}
		}
		p.Battle.SetDuelInput(in)
		s.initialize(p.Battle, nil)
	}
	g.pumpDuelTalks(p)
}

// pumpDuelTalks 把狀態機累積的喊話換成 TALK 索引掛上對白框。
// 變體與肖像照說話大將（`sub_1075B` 的八變體公式）。
func (g *game) pumpDuelTalks(p *state.Pending) {
	s := &g.battleTalkSession
	for _, dt := range p.Battle.TakeDuelTalks() {
		commander := g.battleCommander(p, dt.Side)
		if commander < 0 || commander >= len(g.world.Generals) {
			continue
		}
		general := g.world.Generals[commander]
		// ⭐ **變體用 `+0x1E` 的原始值 0–7**（`sub_1C315` 的
		// `mov ah, [bx+1Eh]`，docs/spec/136 §3）。「模 3 收成 0–2」
		// 是進言那條路（`sub_13C99`）的作法，套到這裡會換成別人的台詞。
		entry := BattleTalkEntry{
			Speaker:  big5(general.TalkName()), // \1 ＝ 呼び名，docs/spec/119
			Index:    resolveBattleTalkIndex(dt.Group, general.TalkVariant),
			Portrait: general.Portrait,
			Side:     dt.Side,
			Duration: battleTalkDuration,
		}
		// 對手的呼び名：參數串流的第一格（docs/spec/136 §2）。
		if o := g.battleCommander(p, 1-dt.Side); o >= 0 && o < len(g.world.Generals) {
			entry.Other = big5(g.world.Generals[o].TalkName())
		}
		// 不能安全代入的 marker 直接丟棄該 entry，不顯示 debug／半句文字。
		if _, _, ok := g.battleTalkText(entry); !ok {
			continue
		}
		s.queue.set(entry)
	}
}

// battleTalkText 回傳整則文字，以及**句中要畫色 9 的那個名字**。
//
// ⚠ 那個名字不一定是說話者：`\1` 拿到誰由 `\6` 吃不吃參數決定
// （docs/spec/136）。renderer 靠它在句中找出那一段換色，
// 傳錯就會整句白字——與原版差在顏色，看起來卻像「字對了」。
func (g *game) battleTalkText(entry BattleTalkEntry) (string, string, bool) {
	if g == nil {
		return "", "", false
	}
	// ⭐ **一條共用的參數游標，不是一個標記一個值**（docs/spec/136）：
	// `sub_1C315` 推的是 `[對手, 說話者]`，而 `\6` **也吃一個**。
	// 於是 `{6}…{1}` 的 `\1` 是說話者、`{1}…` 的 `\1` 是對手——
	// 同一組八個變體裡兩種版面都有。
	lines, names, ok := g.talkStream(entry.Index, []string{entry.Other, entry.Speaker})
	if !ok || len(lines) == 0 {
		return "", "", false
	}
	text := strings.Join(lines, "\n")
	if text == "" {
		return "", "", false
	}
	name := ""
	if len(names) > 0 {
		name = names[0]
	}
	return text, name, true
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
