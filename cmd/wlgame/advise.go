package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 進言（說明書 3.2、3.9）。
//
// **玩家扮演軍師，不是君主**——戰略指令不是直接執行，而是向君主提議。
// 君主可能拒絕，然後玩家挑理由說服他。
//
// 這一層只做畫面與流程；判定全在 internal/rules/persuasion，
// 那邊有 13 條測試把說明書 3.9 的規則釘住。

var adviseCommands = []persuasion.Command{
	persuasion.Hostility, persuasion.CeaseFire, persuasion.Cooperate,
}

// adviseActive 回報進言流程是不是開著。開著時時間會停——
// 它是非常駐視窗（15-realtime.md §2）。
func (g *game) adviseActive() bool { return g.advise != adviseNone }

func (g *game) openAdvise() {
	g.advise = advisePickCommand
	g.sessCur = 0
	g.ally = -1
	g.target = -1
	g.clearAdviseBoxes()
}

func (g *game) closeAdvise() {
	g.advise = adviseNone
	g.sess = nil
	g.clearAdviseBoxes()
}

// adviseSpeaker 是進言畫面上的兩個講話框。**原版靠框的位置與肖像
// 分辨誰在說話**，文字裡沒有說話者標記（docs/spec/45 §1）。
type adviseSpeaker int

const (
	adviseLord    adviseSpeaker = iota // 上框，`sub_13C99`
	adviseAdvisor                      // 下框，`sub_13CDC`，一定是玩家自己
)

func (g *game) clearAdviseBoxes() {
	g.adviseLordSaid = nil
	g.adviseAdvisorSaid = nil
}

// adviseSay 讓其中一個框說一則 TALK。查不到就維持原來那句，
// **不顯示半句、也不把索引當文字**（同 talkLines 的 fail-closed）。
func (g *game) adviseSay(who adviseSpeaker, index int) {
	lines := g.legacyTalkLines(index, g.adviseTalkVars(), talkTextWidth)
	if len(lines) == 0 {
		return
	}
	if who == adviseLord {
		g.adviseLordSaid = lines
		return
	}
	g.adviseAdvisorSaid = lines
}

// situation 把目前的世界狀態換成說服判定要的局勢。
//
// 每一欄都對應已解出的資料：國力用據點數、疲弊用資金 < 0、
// 侵攻狀態用勢力記錄 +0x19（docs/re/08 §1）。
func (g *game) situation(target int) persuasion.Situation {
	me := g.world.Factions[g.world.Player]
	them := g.world.Factions[target]
	// 「有別的勢力正在侵攻我方」——停戰提案的「我正在防禦戰」用這個，
	// 原版是 `sub_16577` 裡掃 22 個勢力的迴圈，**而且跳過交涉對象本身**。
	anyone := false
	for i := range g.world.Factions {
		f := &g.world.Factions[i]
		if f.Alive && i != target && f.InvasionTarget == g.world.Player {
			anyone = true
			break
		}
	}
	s := persuasion.Situation{
		Trust:       g.world.Trust,
		Aggression:  me.Aggression,
		OurCities:   me.Cities,
		TheirCities: them.Cities,
		OurFunds:    me.Funds,
		TheirFunds:  them.Funds,

		// 交友度是**有向的**，而且判定式要的是**含和平位元的原始值**
		// （`0x80` ＝ 和平）——門檻常數本身就帶著它。
		Friendship: g.world.Friendship[g.world.Player][target].Raw(),

		TheyInvadeThirdParty: them.InvasionTarget != diplomacy.NoTarget &&
			them.InvasionTarget != g.world.Player,
		TheyInvadeUs:    them.InvasionTarget == g.world.Player,
		AnyoneInvadesUs: anyone,
	}
	if g.adviseCmd == persuasion.Cooperate && g.ally >= 0 &&
		g.ally < len(g.world.Factions) {
		s.AllyCities = g.world.Factions[g.ally].Cities
		s.AllyFriendship = g.world.Friendship[g.world.Player][g.ally].Raw()
		s.SameFactionPicked = g.ally == target
	}
	return s
}

// updateAdvise 處理進言流程的輸入。回傳 true 表示它吃掉了這一幀的輸入。
func (g *game) updateAdvise() bool {
	if !g.adviseActive() {
		return false
	}
	switch g.advise {
	case advisePickCommand:
		for i := range adviseCommands {
			if pressed(ebiten.Key1 + ebiten.Key(i)) {
				g.adviseCmd = adviseCommands[i]
				g.openTargetList()
				if g.adviseCmd == persuasion.Cooperate {
					g.advise = advisePickAlly
				} else {
					g.advise = advisePickTarget
				}
			}
		}
		if pressed(ebiten.KeyEscape) {
			g.closeAdvise()
		}

	case advisePickAlly, advisePickTarget:
		// 對象用一覽表選，兩段式選取的規則由 listwin 負責。
		if g.list == nil {
			g.closeAdvise()
			return true
		}
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.list.Move(-1)
		case pressed(ebiten.KeyArrowDown):
			g.list.Move(1)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			if id, ok := g.list.Confirm(); ok {
				g.list = nil
				if g.advise == advisePickAlly {
					g.ally = id
					g.openTargetList()
					g.advise = advisePickTarget
				} else {
					g.target = id
					g.beginPersuasion()
				}
			}
		case pressed(ebiten.KeyEscape):
			if g.list.Cancel() {
				g.list = nil
				if g.advise == advisePickTarget && g.adviseCmd == persuasion.Cooperate {
					g.openTargetList()
					g.advise = advisePickAlly
				} else {
					g.advise = advisePickCommand
				}
			}
		}

	case advisePersuade:
		opts := persuasion.Options(g.adviseCmd)
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.sessCur = (g.sessCur + len(opts) - 1) % len(opts)
		case pressed(ebiten.KeyArrowDown):
			g.sessCur = (g.sessCur + 1) % len(opts)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			g.offerReason(opts[g.sessCur])
		case pressed(ebiten.KeyEscape):
			g.closeAdvise()
		}
	}
	return true
}

func (g *game) openTargetList() {
	var rows []int
	for _, i := range g.world.AliveFactions() {
		if i != g.world.Player {
			rows = append(rows, i)
		}
	}
	g.openFactionPicker(rows, "↑↓ 移動　Enter 選取／決定　1-6 排序　ESC 取消", nil)
}

func (g *game) beginPersuasion() {
	s := g.situation(g.target)
	g.sessCur = 0
	g.advise = advisePersuade
	// ① 君主開場（上框）、② 軍師的進言（下框）——原文與框的分工
	// 都照原版（docs/spec/44 §2、docs/spec/45 §1）。
	base := adviseTalkBase(g.adviseCmd)
	g.clearAdviseBoxes()
	g.adviseSay(adviseLord, base+g.playerTalkVariant())
	g.adviseSay(adviseAdvisor, base+3)
	queued := false
	if g.adviseCmd == persuasion.Hostility {
		queued = g.world.HasQueuedDeclaration(g.world.Player, g.target)
	}
	switch reaction := persuasion.FirstReaction(g.adviseCmd, s, queued); reaction {
	case persuasion.AskReason:
		g.sess = persuasion.Begin(g.adviseCmd, s)
	default:
		g.sess = nil
		g.adjustTrust(persuasion.ReactionTrustDelta(reaction))
		base := adviseTalkBase(g.adviseCmd)
		// ③ 君主的回答又回到上框（`sub_13830` 的最後一步）。
		g.adviseSay(adviseLord, adviseReplyIndex(base, reaction, g.playerTalkVariant()))
		if reaction == persuasion.Agree {
			g.commitAdvice()
		}
	}
}

func (g *game) offerReason(r persuasion.Reason) {
	if g.sess == nil {
		return
	}
	repeat := g.sess.Offered(r)
	out, dt := g.sess.Offer(r)
	g.adjustTrust(dt)
	base := adviseReasonBase(g.adviseCmd)
	slot := adviseReasonSlot(g.adviseCmd, r)

	// 軍師說出那個理由 → 下框（`sub_13B5A` 的 `cx = base + 位置 + 1`
	// 之後接 `sub_13CDC`）；君主的反應 → 上框（`sub_13BA9` 結尾接
	// `sub_13C99`）。**兩個框各自只換掉最新的那一句。**
	g.adviseSay(adviseAdvisor, base+slot+1)
	g.adviseSay(adviseLord,
		adviseReasonReply(base, slot, out, repeat, g.playerTalkVariant()))

	switch out {
	case persuasion.Agreed:
		g.commitAdvice()
	case persuasion.Failed:
		g.lastEvent = "說服失敗"
		g.sess = nil
	case persuasion.Withdrawn:
		g.lastEvent = "進言撤回"
		g.sess = nil
	}
}

// adviseReasonBase 是說服迴圈的起點。原版在進迴圈前 `add [bp+0], 10h`
// （`sub_13830`），所以理由那一段全部相對於 base + 16 —— 而 base + 16
// 正好是那三則五選一的選單（#102／#166／#230，docs/spec/44 §1）。
func adviseReasonBase(c persuasion.Command) int { return adviseTalkBase(c) + 0x10 }

// adviseReasonSlot 是理由在這個指令的選單裡排第幾（0–4，4 是撤回）。
// **順序也是資料**——原版的索引算式直接吃這個位置。
func adviseReasonSlot(c persuasion.Command, r persuasion.Reason) int {
	for i, o := range persuasion.Options(c) {
		if o == r {
			return i
		}
	}
	return len(persuasion.Options(c)) - 1 // 找不到就當撤回，不要算出界
}

// adviseReasonReply 是君主對一個理由的反應（原版 `sub_13BA9` 的結尾）：
//
//	撤回（第 5 項）      base + 42
//	同一個理由講第二次   base + 45
//	否則                 base + 位置×9 + 結果×3 + 6
//
// 結果碼 0 ＝ 理由不成立、1 ＝ 湊夠了、2 ＝ 還要再一個。
// 每個位置佔三則是君主的**說話型**變體（`sub_13C99` 的 `add cx, ax`）。
func adviseReasonReply(base, slot int, out persuasion.Outcome, repeat bool, variant int) int {
	switch {
	case out == persuasion.Withdrawn:
		return base + 42 + variant
	case repeat:
		return base + 45 + variant
	}
	code := 2 // Continue：還要再一個
	switch out {
	case persuasion.Failed:
		code = 0
	case persuasion.Agreed:
		code = 1
	}
	return base + slot*9 + code*3 + 6 + variant
}

// adviseReasonLabels 是說服選單的五列。原版把它們放在一則五行的
// TALK（#102／#166／#230，[`docs/re/66`](../../docs/re/66-message-box-geometry.md) §6
// 記的四則五行訊息之三），**順序就是索引算式吃的位置**。
// 讀不到就退回 Reason.String()，那份用字也是照原版選單抄的。
func (g *game) adviseReasonLabels(c persuasion.Command) []string {
	opts := persuasion.Options(c)
	out := make([]string, len(opts))
	for i, r := range opts {
		out[i] = r.String()
	}
	lines, ok := g.talkLines(adviseReasonBase(c), g.adviseTalkVars())
	if !ok || len(lines) != len(opts) {
		return out
	}
	for i, l := range lines {
		if t := strings.TrimRight(l, "　 "); t != "" {
			out[i] = t
		}
	}
	return out
}

// adjustTrust 對應原版 sub_13D91／sub_13DC9 的 byte 飽和行為。
func (g *game) adjustTrust(delta int) {
	if g.world != nil {
		g.world.AdjustTrust(delta)
	}
}

// adviseTalkBase 是三個進言指令的 TALK 起點——`sub_16405`／`sub_164F1`／
// `sub_16623` 傳給 `sub_13830` 的 `cx`（docs/spec/44 §1）。
//
// ⚠ **三組措辭各不相同**，不能共用一組。
func adviseTalkBase(c persuasion.Command) int {
	switch c {
	case persuasion.CeaseFire:
		return 0x96 // 150
	case persuasion.Cooperate:
		return 0xD6 // 214
	default:
		return 0x56 // 86，敵對（進兵）
	}
}

// adviseReplyIndex 是君主回答的 TALK 索引：`base + 4 + 結果碼×3`
// （`sub_13830` 的 `cx = base + al×3 + 4`）。碼 ≥ 4 一律用 83。
//
// 每個位置佔三則，因為 `sub_13C99` 會把君主的**說話型**加進索引。
func adviseReplyIndex(base int, r persuasion.Reaction, variant int) int {
	if r >= persuasion.SameFaction {
		return 0x53 + variant // #83「我想軍師並不是來談笑的。」
	}
	return base + 4 + 3*int(r) + variant
}

// adviseTalkVars 是進言那幾則的變數：`{3}` 交涉對象、`{4}` 軍師（玩家）、
// `{6}` 排版標記（空字串，原版 handler 只調 X 不輸出字元）。
func (g *game) adviseTalkVars() map[byte]string {
	vars := map[byte]string{'6': ""}
	if g.target >= 0 && g.target < len(g.world.Factions) {
		vars['3'] = big5(g.world.LordName(g.target))
	}
	if p := g.world.Player; p >= 0 && p < len(g.world.Factions) {
		if a := g.world.Factions[p].Advisor; a >= 0 && a < len(g.world.Generals) {
			vars['4'] = big5(g.world.Generals[a].Name)
		}
	}
	return vars
}

// commitAdvice 將第一反應直接同意或說服成功接到原版 producer：敵對
// 立即走 sub_13526，停戰／協力分別寫入事件 6／7，讓後續每小時 handler
// 依原版節拍處理外交結果。
func (g *game) commitAdvice() bool {
	var ok bool
	switch g.adviseCmd {
	case persuasion.Hostility:
		ok = g.world.ApplyPlayerHostility(g.target)
	case persuasion.CeaseFire:
		ok = g.world.QueuePlayerCeasefire(g.target)
	case persuasion.Cooperate:
		ok = g.world.QueuePlayerCooperation(g.ally, g.target)
	}
	if !ok {
		// remake 專屬的守門句：君主點頭之後規則層才發現條件變了。
		// 原版沒有這條路徑，所以沒有對應的原文（docs/spec/44 §6）。
		g.adviseLordSaid = []string{"局勢已變，這項進言", "沒有成立。"}
		g.lastEvent = "進言失效"
		g.sess = nil
		return false
	}
	g.lastEvent = g.adviseCmd.String() + " 成立"
	if g.adviseCmd == persuasion.CeaseFire || g.adviseCmd == persuasion.Cooperate {
		g.lastEvent += "（外交事件已排入）"
	}
	g.sess = nil
	return true
}

// drawAdvise 畫進言流程。
func (g *game) drawAdvise(screen *ebiten.Image) {
	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}

	// 進言的視窗也走原版外框（ICONGRF 段 3）。尺寸先進位到 8 的倍數，
	// 不然邊框會切在半塊上。
	box := func(x, y, w, h int) {
		up := func(v int) int { return (v + chrome.Tile - 1) / chrome.Tile * chrome.Tile }
		g.chrome.Window(screen, x, y, up(w), up(h), chrome.Menu)
	}
	lh := textdraw.GlyphH + 2

	switch g.advise {
	case advisePickCommand:
		box(40, 60, 260, 30+len(adviseCommands)*lh)
		g.td.Draw(screen, "進　言", 48, 66, amber)
		for i, c := range adviseCommands {
			g.td.Draw(screen, fmt.Sprintf("%d　%s", i+1, c.String()),
				48, 86+i*lh, white)
		}

	case advisePersuade:
		// 原版的說服畫面（docs/spec/45 §1）：`IVENTGRF` 第 0 頁的插圖，
		// 上框君主、下框軍師，兩個框各只顯示最新的一句。
		// **說話者靠框的位置與肖像分辨**，句子裡不加標記。
		g.drawIventScene(screen, 0)
		g.drawLegacyTalkBox(screen, talkUpperBoxX, talkUpperBoxY,
			talkBoxW, talkBoxH, g.adviseLordSaid, g.playerLordPortrait())
		g.drawLegacyTalkBox(screen, talkLowerBoxX, talkLowerBoxY,
			talkBoxW, talkBoxH, g.adviseAdvisorSaid, g.playerAdvisorPortrait())
		// 提示放在橫幅與上框之間，不要壓到畫面底部的事件視窗。
		hintY := talkUpperBoxY - 4*chrome.Tile
		if g.sess == nil {
			g.drawLegacyHint(screen, "ESC 關閉", hintY)
			return
		}
		g.drawLegacyChoiceBox(screen, g.adviseReasonLabels(g.adviseCmd), g.sessCur)
		g.drawLegacyHint(screen, "↑↓ 選擇　Enter 提出　ESC 放棄", hintY)
	}
}
