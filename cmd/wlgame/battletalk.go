package main

import (
	"strings"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// BattleTalkEntry 是戰術呈現層的最小 TALK payload。
// 文案不存入 state；Index 仍由目前載入的 TALK.DAT 解碼。
type BattleTalkEntry struct {
	// Speaker 是說話武將的名字（Big5 已轉），供 \1 標記代入；
	// 空字串表示這一則沒有名字可代（fail-closed 會把帶 \1 的訊息丟棄）。
	Speaker  string
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
			FieldNumber: p.FieldNumber,
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
	for _, dt := range p.Battle.TakeDuelTalks() {
		g.queueBattleTalk(p, dt.Side, dt.Group, dt.Frame)
	}
}

// queueBattleTalk 共用原版人物變體與參數串流，Side 保留實際攻守側。
func (g *game) queueBattleTalk(p *state.Pending, side, group, atFrame int) bool {
	if g == nil || g.world == nil || p == nil || p.Battle == nil || side < 0 || side > 1 {
		return false
	}
	age := p.Battle.Frame - atFrame
	if age < 0 || age >= battleTalkDuration {
		return false
	}
	commander := g.battleCommander(p, side)
	if commander < 0 || commander >= len(g.world.Generals) {
		return false
	}
	general := g.world.Generals[commander]
	// sub_1C315 以武將+1Eh原始說話型0–7選變體，不能套進言的模3收斂。
	// 名稱仍使用呼び名；參數串流順序是[對手,說話者]，沿spec119／136。
	entry := BattleTalkEntry{
		Speaker:  big5(general.TalkName()),
		Index:    resolveBattleTalkIndex(group, general.TalkVariant),
		Portrait: general.Portrait,
		Side:     side,
		Duration: battleTalkDuration - age,
	}
	if other := g.battleCommander(p, 1-side); other >= 0 && other < len(g.world.Generals) {
		entry.Other = big5(g.world.Generals[other].TalkName())
	}
	if _, _, ok := g.battleTalkText(entry); !ok {
		return false
	}
	g.bindBattleTalk(p.Battle)
	return g.battleTalkSession.queue.set(entry)
}

// playerBattleTalkGroup 對應熱區1C1B9的喊話，拒絕城壁令與退卻另有索引。
func playerBattleTalkGroup(command int, playerDefends, accepted bool) int {
	if command < 0 || command > 5 {
		return -1
	}
	if !accepted {
		if command == 3 {
			return 0x1AC
		}
		return -1
	}
	if command == 5 {
		return 0x1AF
	}
	if command == 3 && playerDefends {
		return 0x1B6
	}
	return 0x1B1 + command
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
	params := []string{entry.Other, entry.Speaker}
	for i, name := range params {
		if name != "" {
			params[i] = padTalkField(name)
		}
	}
	lines, names, ok := g.talkStream(entry.Index, params)
	if !ok || len(lines) == 0 {
		return "", "", false
	}
	text := strings.Join(lines, "\n")
	if text == "" {
		return "", "", false
	}
	name := ""
	if len(names) > 0 {
		name = strings.TrimRight(names[0], " \u3000")
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
	// 戰術喊話不阻擋左鍵選隊／下令（docs/spec/155）；原版喊話仍在時
	// 同一下點擊就執行熱區，不能先把它吃成額外的訊息確認。
	if !g.confirmed() {
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
