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
			keyboardEdit = state.AmountRestoreInitial
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
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
	if row, ok := g.talkChoiceClick(talkChoiceX, talkChoiceY,
		g.diplomacyChoiceLines(*c)); ok {
		g.diplomacyRow = row
		if row == int(state.DiplomacyOfferFunds) {
			g.diplomacyEditingAmount = true
			g.beginAmountEditor(amountCursorDiplomacy)
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
	if !pressed(ebiten.KeyEnter) && !pressed(ebiten.KeySpace) &&
		!g.cancelled() &&
		!pressed(ebiten.Key1) && !pressed(ebiten.Key2) && !pressed(ebiten.Key3) {
		return
	}

	option := state.DiplomacyOption(g.diplomacyRow)
	if option == state.DiplomacyOfferFunds &&
		(pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace)) {
		g.diplomacyEditingAmount = true
		g.beginAmountEditor(amountCursorDiplomacy)
		return
	}
	g.resolveDiplomacy(*c, option)
}

func (g *game) resolveDiplomacy(c state.DiplomacyChoice, option state.DiplomacyOption) {
	kind := c.Kind
	ok := g.world.ResolveDiplomacy(option)
	g.enqueueDiplomacyTalk(c, option)
	g.diplomacyRow = 0
	g.diplomacyEditingAmount = false
	g.amountCursorActive = false
	if option == state.DiplomacyReject {
		g.setEvent(diplomacyTitle(kind) + "：拒絕")
		return
	}
	if !ok {
		g.setEvent(diplomacyTitle(kind) + "：交涉未成立")
		return
	}
	if option == state.DiplomacyOfferFunds {
		g.setEvent(diplomacyTitle(kind) + "：提供資金後成立")
	} else {
		g.setEvent(diplomacyTitle(kind) + "：無條件成立")
	}
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
// 玩家／合作方，所以真正提出協力請求的是 Invader；事件 3 則是 Source。
func (g *game) diplomacyTalkRequester(c state.DiplomacyChoice) string {
	id := c.Source
	if c.Kind == state.DiplomacyCooperation {
		id = c.Invader
	}
	return g.diplomacyFactionName(id)
}

// diplomacyTalkResponse 對應 sub_13902 回傳給 sub_13C3D 的 AL。外交指定金額
// 超過初始要求雖然 state 收尾 fail-closed，但原版仍顯示 response 2 的破裂句。
func diplomacyTalkResponse(c state.DiplomacyChoice, option state.DiplomacyOption) (int, int) {
	if option == state.DiplomacyReject ||
		(option == state.DiplomacyOfferFunds && c.OfferAmount > c.InitialAmount) {
		return 2, -1
	}
	if option == state.DiplomacyOfferFunds && c.OfferAmount > 0 {
		return 1, c.OfferAmount
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
func (g *game) enqueueDiplomacyTalk(c state.DiplomacyChoice, option state.DiplomacyOption) {
	choiceAmount := -1
	if option == state.DiplomacyOfferFunds {
		choiceAmount = c.OfferAmount
	}
	// 玩家挑的那一列由**軍師**在事件場景的下框說出來（原版 `sub_13CDC`，
	// docs/spec/42 §2）。插圖是事件 2／3 的第 0 頁，留在背後。
	g.enqueueAdvisorTalk(diplomacyTalkChoiceIndex(c, option),
		g.diplomacyTalkVars(c, choiceAmount), 0)
	response, resultAmount := diplomacyTalkResponse(c, option)
	resultBase := 0x2B // TALK #43; event 3
	if c.Kind == state.DiplomacyCooperation {
		resultBase = 0x2F // TALK #47; event 2
	}
	g.enqueueTalkWithPortrait(resultBase+response,
		g.diplomacyTalkVars(c, resultAmount), g.diplomacyResultPortrait(c))
}

// diplomacyResultPortrait 是結果句的講話者（原版 `sub_13C3D` 開頭）：
//
//	bl = 93h                       ; 預設
//	cmp si, cs:word_10CFD / jz →   ; si 就是玩家的勢力 ⇒ 用預設
//	bh = [si+2Ah]                  ; 那個勢力的「派駐外交官」
//	bl = [bx+4241h]                ; ★ 那名武將的 +0x01 頭像
//
// **回報結果的是派去的外交官**，不是一般通知的那張臉。
// 沒有外交官（原版 `+0x2A` 寫 `0xFF`）就退回預設，不畫錯人。
func (g *game) diplomacyResultPortrait(c state.DiplomacyChoice) int {
	if g == nil || g.world == nil {
		return defaultPortraitPage
	}
	id := c.Source
	if c.Kind == state.DiplomacyCooperation {
		id = c.Invader
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
	prompt := g.legacyTalkLines(g.diplomacyTalkPromptIndex(*c),
		g.diplomacyTalkVars(*c, -1), talkTextWidth)
	g.drawLegacyTalkBox(screen, talkUpperBoxX, talkUpperBoxY, talkBoxW, talkBoxH,
		prompt, g.playerLordPortrait())
	if g.diplomacyEditingAmount {
		g.drawAmountPanel(screen, c.OfferAmount, c.InitialAmount, true)
		return
	}
	g.drawLegacyChoiceBox(screen, talkChoiceX, talkChoiceY,
		g.diplomacyChoiceLines(*c), g.diplomacyRow)
}
