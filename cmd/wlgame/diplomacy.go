package main

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// updateDiplomacy 是事件 2／3 的玩家三選一接縫。
// 選到資金列後才進入原版 sub_17C6E 的 3×6 數值器；滑鼠點格、數字鍵
// 與退格／Insert／Delete／Home 都共用同一個 AmountEdit，不再把數字
// 編輯誤當成現代 UI 的左右加減。
func (g *game) updateDiplomacy() {
	c := g.world.PendingDiplomacy()
	if c == nil {
		g.diplomacyEditingAmount = false
		return
	}
	if g.diplomacyEditingAmount {
		if g.cancelled() {
			g.diplomacyEditingAmount = false
			return
		}
		if button, ok := g.amountPointerButton(); ok {
			if button.edit == state.AmountFinishInput {
				g.diplomacyEditingAmount = false
				g.resolveDiplomacy(*c, state.DiplomacyOfferFunds)
				return
			}
			g.world.EditDiplomacyOfferAmount(button.edit, button.digit)
			return
		}
		if digit, ok := pressedAmountDigit(); ok {
			g.setAmountCursorAction(state.AmountAppendDigit, digit)
			g.world.EditDiplomacyOfferAmount(state.AmountAppendDigit, digit)
			return
		}
		keyboardEdit := state.AmountEdit(255)
		switch {
		case pressed(ebiten.KeyBackspace):
			keyboardEdit = state.AmountDeleteDigit
		case pressed(ebiten.KeyInsert):
			keyboardEdit = state.AmountAppendHundred
		case pressed(ebiten.KeyDelete):
			keyboardEdit = state.AmountClear
		case pressed(ebiten.KeyHome):
			keyboardEdit = state.AmountSetMax
		case g.confirmed():
			g.setAmountCursorAction(state.AmountFinishInput, 0)
			g.diplomacyEditingAmount = false
			g.resolveDiplomacy(*c, state.DiplomacyOfferFunds)
			return
		}
		if keyboardEdit != state.AmountEdit(255) {
			g.setAmountCursorAction(keyboardEdit, 0)
			g.world.EditDiplomacyOfferAmount(keyboardEdit, 0)
		}
		return
	}
	if row, ok := g.talkChoiceInput(&g.diplomacyRow, len(g.diplomacyChoiceLines(*c))); ok {
		g.diplomacyRow = row
		if row == int(state.DiplomacyOfferFunds) {
			g.beginDiplomacyAmount()
			return
		}
		g.resolveDiplomacy(*c, state.DiplomacyOption(row))
		return
	}
	if pressed(ebiten.KeyArrowUp) {
		g.diplomacyRow = (g.diplomacyRow + 2) % 3
	}
	if pressed(ebiten.KeyArrowDown) {
		g.diplomacyRow = (g.diplomacyRow + 1) % 3
	}
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if pressed(k) {
			g.diplomacyRow = i
		}
	}
	if g.cancelled() {
		g.diplomacyRow = int(state.DiplomacyReject)
	}
	if !g.confirmed() &&
		!g.cancelled() &&
		!pressed(ebiten.Key1) && !pressed(ebiten.Key2) && !pressed(ebiten.Key3) {
		return
	}

	option := state.DiplomacyOption(g.diplomacyRow)
	if option == state.DiplomacyOfferFunds &&
		(g.confirmed()) {
		g.beginDiplomacyAmount()
		return
	}
	g.resolveDiplomacy(*c, option)
}

func (g *game) beginDiplomacyAmount() {
	g.world.SetDiplomacyOfferAmount(0)
	g.diplomacyEditingAmount = true
	g.beginAmountEditor(amountCursorDiplomacy, amountAnchorEventX, amountAnchorEventY)
}

func (g *game) resolveDiplomacy(c state.DiplomacyChoice, option state.DiplomacyOption) {
	result, settle := g.world.PrepareDiplomacy(option)
	g.enqueueDiplomacyTalkWithSettlement(c, result, settle)
	g.diplomacyRow = 0
	g.diplomacyEditingAmount = false
	g.amountCursorActive = false
	g.lastEvent = ""
}

// diplomacyTalkBase 對應 sub_138C7／sub_138E6 傳給 sub_13902 的 CX。
func diplomacyTalkBase(kind state.DiplomacyKind) int {
	if kind == state.DiplomacyCeasefire {
		return 0x168 // TALK #360
	}
	return 0x175 // TALK #373
}

// diplomacyTalkPromptIndex 對應 sub_13902 → sub_13C99：base 不是固定
// 顯示頁，還要套用玩家君主 General +0x1E 的 0–2 變體。
func (g *game) diplomacyTalkPromptIndex(c state.DiplomacyChoice) int {
	return diplomacyTalkPromptIndex(c, g.playerTalkVariant())
}

// diplomacyTalkChoiceIndex 對應 sub_13902 的 base+4+choice：
// choice 0／1／2 分別是無條件、指定金額、拒絕建言。
func diplomacyTalkChoiceIndex(c state.DiplomacyChoice, option state.DiplomacyOption) int {
	return diplomacyTalkBase(c.Kind) + 4 + int(option)
}

// diplomacyTalkRequester 是 TALK.DAT 的 \\3 參數來源。事件 2 的事件字高是
// 玩家／合作方，請求方名稱來自DI指向的Target；事件3則是Source。
func (g *game) diplomacyTalkRequester(c state.DiplomacyChoice) string {
	id := c.Source
	if c.Kind == state.DiplomacyCooperation {
		id = c.Target
	}
	return g.diplomacyFactionName(id)
}

// diplomacyTalkResponse 對應 sub_13902 回傳給 sub_13C3D 的 AL。外交指定金額
// 超過初始要求雖然 state 收尾 fail-closed，但原版仍顯示 response 2 的破裂句。
func diplomacyTalkResponse(result state.DiplomacyResult) (int, int) {
	if result.Response >= 2 {
		return 2, -1
	}
	if result.Response == 1 {
		return 1, result.Amount
	}
	return 0, -1
}

func (g *game) diplomacyTalkVars(c state.DiplomacyChoice, amount int) map[byte]string {
	vars := map[byte]string{
		'3': g.diplomacyTalkRequester(c),
		'6': "",
	}
	if amount >= 0 {
		vars['7'] = strconv.Itoa(amount)
	}
	return vars
}

// enqueueDiplomacyTalk 接回 sub_13902 的已證實文字索引與
// sub_13C3D 的主要結果句；原版次要信賴度／AH 分支仍留在未完成邊界。
func (g *game) enqueueDiplomacyTalk(c state.DiplomacyChoice, result state.DiplomacyResult) {
	g.enqueueDiplomacyTalkWithSettlement(c, result, nil)
}

func (g *game) enqueueDiplomacyTalkWithSettlement(c state.DiplomacyChoice, result state.DiplomacyResult,
	settle func() (state.DiplomacyResult, bool)) {
	choiceAmount := -1
	if result.AdviceOption == state.DiplomacyOfferFunds {
		choiceAmount = c.OfferAmount
	}
	prompt := g.diplomacySceneMessage(g.diplomacyTalkPromptIndex(c), g.diplomacyTalkVars(c, -1), g.playerLordPortrait())
	if prompt != nil {
		prompt.upper, prompt.scene = true, 0
	}
	// 玩家挑的那一列由**軍師**在事件場景的下框說出來（原版 `sub_13CDC`，
	// docs/spec/42 §2）。插圖是事件 2／3 的第 0 頁，留在背後。
	before := len(g.messages)
	g.enqueueAdvisorTalk(diplomacyTalkChoiceIndex(c, result.AdviceOption),
		g.diplomacyTalkVars(c, choiceAmount), 0)
	var advisor *messageDialog
	if len(g.messages) > before {
		g.messages[len(g.messages)-1].underlay = prompt
		copy := g.messages[len(g.messages)-1]
		advisor = &copy
	}
	replyIndex := g.diplomacyTalkPromptIndex(c) + 7
	if result.AdviceAdopted {
		replyIndex += 3
	}
	if reply := g.diplomacySceneMessage(replyIndex, g.diplomacyTalkVars(c, -1), g.playerLordPortrait()); reply != nil {
		reply.upper, reply.scene, reply.underlay = true, 0, advisor
		g.messages = append(g.messages, *reply)
	}
	response, resultAmount := diplomacyTalkResponse(result)
	resultBase := 0x2B // TALK #43; event 3
	if c.Kind == state.DiplomacyCooperation {
		resultBase = 0x2F // TALK #47; event 2
	}
	g.enqueueTalkWithPortrait(resultBase+response,
		g.diplomacyTalkVars(c, resultAmount), g.diplomacyResultPortrait(c))
	if settle != nil {
		// 13C6E的中央結果確認後，才付款或13DC9扣信賴。
		g.afterTalk(func() {
			settled, _ := settle()
			g.enqueueDiplomacyConsequences(c, settled)
		})
		return
	}
	g.enqueueDiplomacyConsequences(c, result)
}

func (g *game) enqueueDiplomacyConsequences(c state.DiplomacyChoice, result state.DiplomacyResult) {
	if result.Response == 3 {
		if g.world != nil && g.world.Player >= 0 && g.world.Player < len(g.world.Factions) {
			g.enqueueTalkNotice(state.TalkNotice{
				Index: 0x1A5, City: -1, Faction: c.Source,
				General: g.world.Factions[g.world.Player].Lord, Amount: -1,
			})
		}
	}
	// 13258的135ED先等所有150D7返還對白，1325D才13526宣戰。
	for _, id := range result.ReleasedGenerals {
		g.enqueueReleasedGeneral(id)
	}
	for _, notice := range result.TalkNotices {
		g.enqueueTalkNotice(notice)
	}
}

// 只取已由共用TALK排版器建立的視窗，供同一場景保留上下框。
func (g *game) diplomacySceneMessage(index int, vars map[byte]string, portrait int) *messageDialog {
	before := len(g.messages)
	g.enqueueTalkWithPortrait(index, vars, portrait)
	if len(g.messages) == before {
		return nil
	}
	d := g.messages[len(g.messages)-1]
	g.messages = g.messages[:before]
	return &d
}

// diplomacyResultPortrait沿sub_13C3D入口的SI取說話者：停戰是Target，合作是Source。
// SI等於玩家勢力時使用147；其他勢力才從該勢力+2Ah取外交官肖像。
func (g *game) diplomacyResultPortrait(c state.DiplomacyChoice) int {
	if g == nil || g.world == nil {
		return defaultPortraitPage
	}
	id := c.Source
	if c.Kind == state.DiplomacyCeasefire {
		// 13262在13C3D入口的SI是停戰接收方，玩家收到來使時用通報者。
		id = c.Target
	}
	if id < 0 || id >= len(g.world.Factions) || id == g.world.Player {
		return defaultPortraitPage
	}
	envoy := g.world.Factions[id].Diplomat
	if envoy < 0 || envoy >= len(g.world.Generals) || !g.world.Generals[envoy].Alive {
		return defaultPortraitPage
	}
	return g.world.Generals[envoy].Portrait
}

func diplomacyTitle(kind state.DiplomacyKind) string {
	if kind == state.DiplomacyCeasefire {
		return "停戰交涉"
	}
	return "協力交涉"
}

func (g *game) diplomacyFactionName(id int) string {
	if id < 0 || id >= len(g.world.Factions) || !g.world.Factions[id].Alive {
		return "－"
	}
	return big5(g.world.LordName(id))
}

// diplomacyChoiceLines 是三選一的那三列。**不折行**——選單的字寬
// 決定框寬（docs/spec/45 §2.2），折了框就跟著錯。
func (g *game) diplomacyChoiceLines(c state.DiplomacyChoice) []string {
	lines, _ := g.talkLines(diplomacyTalkBase(c.Kind)+3, g.diplomacyTalkVars(c, -1))
	return lines
}

// drawDiplomacy 畫出事件 2／3 的原版 composite：IVENTGRF page 0、
// 玩家君主肖像／TALK、三列 TALK 選項；只有選到資金列後才覆蓋原版
// `(80,176)` 數值面板。
func (g *game) drawDiplomacy(screen *ebiten.Image, c *state.DiplomacyChoice) {
	g.drawIventScene(screen, 0)
	vars, fields := padTalkVars(g.diplomacyTalkVars(*c, -1))
	prompt := g.legacyTalkLines(g.diplomacyTalkPromptIndex(*c), vars, talkTextWidth)
	g.drawLegacyTalkBoxFields(screen, talkUpperBoxX, talkUpperBoxY, talkBoxW, talkBoxH,
		prompt, g.playerLordPortrait(), fields)
	if g.diplomacyEditingAmount {
		g.drawAmountPanel(screen, c.OfferAmount, true)
		return
	}
	g.drawLegacyChoiceBox(screen, talkChoiceX, talkChoiceY,
		g.diplomacyChoiceLines(*c), g.diplomacyRow)
}
