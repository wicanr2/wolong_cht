package main

import (

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/talkmenu"
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

// 第四項與第五項不走說服迴圈——君主只做一次驗收（docs/spec/49）。
// 用字照原版的指令列（`docs/re/22`）。
const (
	adviseRelocateRow = 3 // 遷都
	adviseSortieRow   = 4 // 請求君主出陣
)

// adviseFallbackNames 只在讀不到 `TALK.DAT` 時用。**正常路徑走 #77**
// （`sub_16224` 傳給 `sub_193E9` 的 `cx = 4Dh`，五列剛好對上）。
var adviseFallbackNames = []string{
	"敵對提案", "停戰提案", "請求協助", "遷　都", "請求君主出陣",
}

// adviseMenuTalk 是進言那五項的選單訊息。
const adviseMenuTalk = persuasion.TalkMenu

// adviseCommandLabels 是進言選單的五列，直接取自 `TALK.DAT`。
//
// ⚠ 原版寫「請求協助」，而 `persuasion.Cooperate.String()` 是「協力要請」
// ——**選單文字屬於松崗版的原文，不要拿內部術語頂替**。
func (g *game) adviseCommandLabels() []string {
	if g == nil || g.lib == nil {
		return append([]string(nil), adviseFallbackNames...)
	}
	return talkmenu.Labels(g.lib.Talk, adviseMenuTalk, nil, adviseFallbackNames)
}

// adviseVerdictBase 是那兩項的 TALK 起點——`sub_16909`／`sub_1699E`
// 傳給 `sub_13B08` 的 `cx`（docs/spec/49 §1）。
const (
	adviseRelocateTalkBase = persuasion.TalkRelocateBase
	adviseSortieTalkBase   = persuasion.TalkSortieBase
)

// adviseActive 回報進言流程是不是開著。開著時時間會停——
// 它是非常駐視窗（15-realtime.md §2）。
func (g *game) adviseActive() bool { return g.advise != adviseNone }

// adviseLordAwayEvent 是君主帶兵在外時，進言開不起來的那一行。
//
// **原版沒有這句**：原版的編成候選排除君主，君主帶兵只能靠進言的第五項，
// 所以「君主不在朝堂上而玩家要進言」這個狀態走不到（docs/spec/111）。
const adviseLordAwayEvent = "主公正在領軍，無法進言"

func (g *game) openAdvise() {
	// 進言是軍師對君主說話（docs/spec/45 §1：上框君主、下框軍師）。
	// 君主帶著軍團離開之後對話的另一方不在，五項沒有一項成立——
	// 與「請求君主出陣」用同一個判準，只是套用的層級往上提一格。
	if g.world != nil && g.world.LordLeadsCorps() {
		g.lastEvent = adviseLordAwayEvent
		return
	}
	g.advise = advisePickCommand
	g.adviseCmdRow = 0
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
	g.adviseQueue = nil
}

// adviseStep 是一句待演的話。**原版每一句演完都停下來等按鍵**
// （`sub_13C99`／`sub_13CDC` 結尾的 `sub_12216`，docs/spec/45 §1.1），
// 所以句子要排隊，不能一次全寫進框裡。
type adviseStep struct {
	who   adviseSpeaker
	index int
}

// adviseSay 把一句排進隊伍。查不到的索引直接不排——
// **不顯示半句、也不把索引當文字**（同 talkLines 的 fail-closed）。
func (g *game) adviseSay(who adviseSpeaker, index int) {
	if len(g.legacyTalkLines(index, g.adviseTalkVars(), talkTextWidth)) == 0 {
		return
	}
	g.adviseQueue = append(g.adviseQueue, adviseStep{who: who, index: index})
}

// adviseAdvance 演下一句，回傳還有沒有下一句可以演。
func (g *game) adviseAdvance() bool {
	if len(g.adviseQueue) == 0 {
		return false
	}
	step := g.adviseQueue[0]
	g.adviseQueue = g.adviseQueue[1:]
	lines := g.legacyTalkLines(step.index, g.adviseTalkVars(), talkTextWidth)
	if step.who == adviseLord {
		g.adviseLordSaid = lines
	} else {
		g.adviseAdvisorSaid = lines
	}
	return true
}

// adviseTalking 回報還有句子沒演完。**還在講話時不接受任何選擇**——
// 原版也是這樣，選單要等對話走完才出現。
func (g *game) adviseTalking() bool { return len(g.adviseQueue) > 0 }

// situation 把目前的世界狀態換成說服判定要的局勢。
// 組法在 `internal/state`，桌面版與手機版共用同一份。
func (g *game) situation(target int) persuasion.Situation {
	return g.world.PersuasionSituation(g.adviseCmd, target, g.ally)
}

// updateAdvise 處理進言流程的輸入。回傳 true 表示它吃掉了這一幀的輸入。
func (g *game) updateAdvise() bool {
	if !g.adviseActive() {
		return false
	}
	switch g.advise {
	case advisePickCommand:
		labels := g.adviseCommandLabels()
		if row, ok := g.talkChoiceClick(adviseMenuX, adviseMenuY, labels); ok {
			g.adviseCmdRow = row
			g.pickAdviseCommand(row)
			return true
		}
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.adviseCmdRow = (g.adviseCmdRow + len(labels) - 1) % len(labels)
		case pressed(ebiten.KeyArrowDown):
			g.adviseCmdRow = (g.adviseCmdRow + 1) % len(labels)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			g.pickAdviseCommand(g.adviseCmdRow)
		case g.cancelled():
			g.closeAdvise()
		}
		// 數字鍵是 remake 加的捷徑；原版只有游標選取。
		for i := range labels {
			if pressed(ebiten.Key1 + ebiten.Key(i)) {
				g.adviseCmdRow = i
				g.pickAdviseCommand(i)
				break
			}
		}

	case advisePickCapital:
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
			if node, ok := g.list.Confirm(); ok {
				g.list = nil
				g.beginRelocate(node)
			}
		case g.cancelled():
			if g.list.Cancel() {
				g.list = nil
				g.advise = advisePickCommand
			}
		}

	case adviseVerdict:
		if g.cancelled() || pressed(ebiten.KeyEnter) ||
			pressed(ebiten.KeySpace) {
			if !g.adviseAdvance() {
				g.closeAdvise()
			}
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
		case g.cancelled():
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
		if g.adviseTalking() {
			if pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) ||
				pressed(ebiten.KeyEscape) {
				g.adviseAdvance()
			}
			return true
		}
		opts := persuasion.Options(g.adviseCmd)
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.sessCur = (g.sessCur + len(opts) - 1) % len(opts)
		case pressed(ebiten.KeyArrowDown):
			g.sessCur = (g.sessCur + 1) % len(opts)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			g.offerReason(opts[g.sessCur])
		case g.cancelled():
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

// pickAdviseCommand 分派進言的五項（原版 `sub_16224` 的 `funcs_16255[選項×2]`）。
func (g *game) pickAdviseCommand(row int) {
	switch row {
	case adviseRelocateRow:
		g.openCapitalList()
	case adviseSortieRow:
		g.beginSortie()
	default:
		if row < 0 || row >= len(adviseCommands) {
			return
		}
		g.adviseCmd = adviseCommands[row]
		g.openTargetList()
		if g.adviseCmd == persuasion.Cooperate {
			g.advise = advisePickAlly
		} else {
			g.advise = advisePickTarget
		}
	}
}

// openCapitalList 列出自己的據點讓玩家挑遷都目標。
// 原版是 `sub_18853(cx=0Fh)` ＋ `sub_17400` 的地圖選點；
// remake 用一覽表挑，**這是操作方式的差異，不是規則的差異**。
func (g *game) openCapitalList() {
	var rows []int
	for i := range g.world.Cities {
		if g.world.Cities[i].Owner == g.world.Player {
			rows = append(rows, i)
		}
	}
	if len(rows) == 0 {
		return
	}
	g.openCityPicker(rows, "↑↓ 移動　Enter 選取／決定　ESC 取消", nil)
	g.advise = advisePickCapital
}

// beginRelocate 是進言第四項：君主看一眼就定案，沒有說服迴圈。
func (g *game) beginRelocate(node int) {
	g.adviseNode = node
	ok := g.world.AdviseRelocateAccepted(node)
	g.sayVerdict(adviseRelocateTalkBase, ok)
	if ok && g.world.AdviseRelocate(node) {
		// 搬成之後君主再講一句（`sub_133FD` 的玩家分支，docs/spec/64 §1.1）：
		// 組編號 0x1A4 配**君主自己的原始說話型**，主公型落在 0–2。
		g.sayCapitalMoved(node)
		g.lastEvent = "遷都 " + big5(g.world.Cities[node].Name)
	} else {
		g.lastEvent = "遷都：君主不同意"
	}
}

// sayCapitalMoved 是遷都定案之後君主講的那一句（docs/spec/64 §1.1）。
// 說話者與肖像都是**君主**（原版 `[si+1]` → `[bx+425Eh]`／`[bx+4241h]`）。
func (g *game) sayCapitalMoved(node int) {
	if g == nil || g.world == nil || g.world.Player < 0 ||
		g.world.Player >= len(g.world.Factions) {
		return
	}
	lord := g.world.Factions[g.world.Player].Lord
	if lord < 0 || lord >= len(g.world.Generals) {
		return
	}
	city := ""
	if node >= 0 && node < len(g.world.Cities) {
		city = big5(g.world.Cities[node].Name)
	}
	// #519 的 {4} 是軍師的**呼び名**（原版 marker \\4，`00010939` 取
	// 武將記錄 +0x08，docs/spec/119）。
	advisorName := ""
	if a := g.world.Factions[g.world.Player].Advisor; a >= 0 && a < len(g.world.Generals) {
		advisorName = big5(g.world.Generals[a].TalkName())
	}
	gen := g.world.Generals[lord]
	g.enqueueTalkWithPortrait(
		resolveBattleTalkIndex(state.CapitalMovedTalkBase, gen.TalkVariant),
		map[byte]string{'2': city, '6': "", '4': advisorName}, gen.Portrait)
}

// beginSortie 是進言第五項：兩道閘都過君主才親自出陣。
//
// ⚠ **看到君主說話不代表提議被接受**——原版無論通不通過都跳 #396，
// 差別只在第三句（`docs/mechanics/70-ai.md`）。
func (g *game) beginSortie() {
	ok := g.world.AdviseSortieAccepted()
	g.sayVerdict(adviseSortieTalkBase, ok)
	if ok && g.world.AdviseSortie() {
		g.lastEvent = "君主親自出陣"
	} else {
		g.lastEvent = "請求出陣：君主不同意"
	}
}

// sayVerdict 演 `sub_13B08` 的三句：上框君主開場、下框軍師、上框君主定案。
//
//	cx        君主開場（`sub_13C99` 自己加說話型變體）
//	cx + 3    軍師（`sub_13CDC`，不加變體）
//	cx + 4    接受；cx + 7 拒絕（原版 `and bp, bp / jz` 之後 `add cx, 3`）
func (g *game) sayVerdict(base int, accepted bool) {
	g.advise = adviseVerdict
	g.clearAdviseBoxes()
	g.adviseSay(adviseLord, base+g.playerTalkVariant())
	g.adviseSay(adviseAdvisor, base+3)
	reply := base + 4
	if !accepted {
		reply += 3
	}
	g.adviseSay(adviseLord, reply+g.playerTalkVariant())
	g.adviseAdvance()
}

func (g *game) beginPersuasion() {
	s := g.situation(g.target)
	g.sessCur = 0
	g.advise = advisePersuade
	// ① 君主開場（上框）、② 軍師的進言（下框）——原文與框的分工
	// 都照原版（docs/spec/44 §2、docs/spec/45 §1）。
	base := persuasion.TalkBase(g.adviseCmd)
	g.clearAdviseBoxes()
	g.adviseSay(adviseLord, base+g.playerTalkVariant())
	g.adviseSay(adviseAdvisor, base+3)
	queued := false
	if g.adviseCmd == persuasion.Hostility {
		queued = g.world.HasQueuedDeclaration(g.world.Player, g.target)
	}
	switch reaction := persuasion.FirstReaction(g.adviseCmd, s, queued); reaction {
	case persuasion.AskReason:
		// ③ 君主的回答也在上框——**四個反應碼共用同一個算式**，
		// AskReason 只是多了「接著開說服迴圈」。先前這一支沒演第三句，
		// 選單直接跳出來而上框停在開場那句（docs/spec/108）。
		g.adviseSay(adviseLord,
			persuasion.TalkReplyIndex(base, reaction, g.playerTalkVariant()))
		g.sess = persuasion.Begin(g.adviseCmd, s)
	default:
		g.sess = nil
		g.adjustTrust(persuasion.ReactionTrustDelta(reaction))
		base := persuasion.TalkBase(g.adviseCmd)
		// ③ 君主的回答又回到上框（`sub_13830` 的最後一步）。
		g.adviseSay(adviseLord, persuasion.TalkReplyIndex(base, reaction, g.playerTalkVariant()))
		if reaction == persuasion.Agree {
			g.commitAdvice()
		}
	}
	g.adviseAdvance() // 先演第一句，其餘等玩家按鍵
}

func (g *game) offerReason(r persuasion.Reason) {
	if g.sess == nil {
		return
	}
	repeat := g.sess.Offered(r)
	out, dt := g.sess.Offer(r)
	g.adjustTrust(dt)
	base := persuasion.TalkReasonBase(g.adviseCmd)
	slot := persuasion.TalkReasonSlot(g.adviseCmd, r)

	// 軍師說出那個理由 → 下框（`sub_13B5A` 的 `cx = base + 位置 + 1`
	// 之後接 `sub_13CDC`）；君主的反應 → 上框（`sub_13BA9` 結尾接
	// `sub_13C99`）。**兩個框各自只換掉最新的那一句。**
	g.adviseSay(adviseAdvisor, base+slot+1)
	g.adviseSay(adviseLord,
		persuasion.TalkReasonReply(base, slot, out, repeat, g.playerTalkVariant()))
	g.adviseAdvance()

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

// adviseReasonLabels 是說服選單的五列。原版把它們放在一則五行的
// TALK（#102／#166／#230，[`docs/re/66`](../../docs/re/66-message-box-geometry.md) §6
// 記的四則五行訊息之三）。讀不到就退回 Reason.String()，
// 那份用字也是照原版選單抄的。
func (g *game) adviseReasonLabels(c persuasion.Command) []string {
	opts := persuasion.Options(c)
	fallback := make([]string, len(opts))
	for i, r := range opts {
		fallback[i] = r.String()
	}
	if g == nil || g.lib == nil {
		return fallback
	}
	return talkmenu.Labels(g.lib.Talk, persuasion.TalkReasonBase(c),
		g.adviseTalkVars(), fallback)
}

// adjustTrust 對應原版 sub_13D91／sub_13DC9 的 byte 飽和行為。
func (g *game) adjustTrust(delta int) {
	if g.world != nil {
		g.world.AdjustTrust(delta)
	}
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

// drawAdviseBoxes 畫上下兩個講話框。**還沒講話的框不畫**——
// 原版的框是 `sub_13C99`／`sub_13CDC` 在講那一句時才畫出來的，
// 不是一直掛在那裡的空框。
func (g *game) drawAdviseBoxes(screen *ebiten.Image) {
	if len(g.adviseLordSaid) > 0 {
		g.drawLegacyTalkBox(screen, talkUpperBoxX, talkUpperBoxY,
			talkBoxW, talkBoxH, g.adviseLordSaid, g.playerLordPortrait())
	}
	if len(g.adviseAdvisorSaid) > 0 {
		g.drawLegacyTalkBox(screen, talkLowerBoxX, talkLowerBoxY,
			talkBoxW, talkBoxH, g.adviseAdvisorSaid, g.playerAdvisorPortrait())
	}
}

// drawAdvise 畫進言流程。
func (g *game) drawAdvise(screen *ebiten.Image) {
	switch g.advise {
	case advisePickCommand:
		// 原版的位置：`sub_16224` 的 `dx = 400h` ⇒ 粗格 (0, 4) ⇒ (0, 64)。
		// 大小由內容算（docs/spec/45 §2.2），所以不必再挑一個寬度。
		g.drawLegacyChoiceBox(screen, adviseMenuX, adviseMenuY,
			g.adviseCommandLabels(), g.adviseCmdRow)

	case adviseVerdict:
		// 第四、五項沒有說服迴圈，畫面就是 `sub_13B08` 的三句。
		g.drawIventScene(screen, 0)
		g.drawAdviseBoxes(screen)
		hint := "Enter／ESC 關閉"
		if g.adviseTalking() {
			hint = "Enter 繼續"
		}
		g.drawLegacyHint(screen, hint, talkUpperBoxY-4*chrome.Tile)

	case advisePersuade:
		// 原版的說服畫面（docs/spec/45 §1）：`IVENTGRF` 第 0 頁的插圖，
		// 上框君主、下框軍師，兩個框各只顯示最新的一句。
		// **說話者靠框的位置與肖像分辨**，句子裡不加標記。
		g.drawIventScene(screen, 0)
		g.drawAdviseBoxes(screen)
		// 提示放在橫幅與上框之間，不要壓到畫面底部的事件視窗。
		hintY := talkUpperBoxY - 4*chrome.Tile
		if g.adviseTalking() {
			g.drawLegacyHint(screen, "Enter 繼續", hintY)
			return
		}
		if g.sess == nil {
			g.drawLegacyHint(screen, "ESC 關閉", hintY)
			return
		}
		g.drawLegacyChoiceBox(screen, talkChoiceX, talkChoiceY,
			g.adviseReasonLabels(g.adviseCmd), g.sessCur)
		g.drawLegacyHint(screen, "↑↓ 選擇　Enter 提出　ESC 放棄", hintY)
	}
}
