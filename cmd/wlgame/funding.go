package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// updateFunding 是事件 4／5 的玩家三選一接縫。選到指定金額後才進入
// sub_139E8 → sub_17C6E；滑鼠 3×6 格與跨平台鍵盤都呼叫同一個
// AmountEdit，完成格才會離開數值器並消費三選一。
func (g *game) updateFunding() {
	c := g.world.PendingFunding()
	if c == nil {
		g.fundingEditingAmount = false
		return
	}
	if g.fundingEditingAmount {
		if g.cancelled() {
			g.fundingEditingAmount = false
			g.amountCursorActive = false
			g.fundingRow = 0
			g.world.SetFundingAmount(c.RequestedAmount)
			return
		}
		if button, ok := g.amountPointerButton(); ok {
			if button.edit == state.AmountFinishInput {
				g.fundingEditingAmount = false
				g.resolveFunding(*c, state.FundingSetAmount)
				return
			}
			g.world.EditFundingAmount(button.edit, button.digit)
			return
		}
		if digit, ok := pressedAmountDigit(); ok {
			g.setAmountCursorAction(state.AmountAppendDigit, digit)
			g.world.EditFundingAmount(state.AmountAppendDigit, digit)
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
			g.fundingEditingAmount = false
			g.resolveFunding(*c, state.FundingSetAmount)
			return
		}
		if keyboardEdit != state.AmountEdit(255) {
			g.setAmountCursorAction(keyboardEdit, 0)
			g.world.EditFundingAmount(keyboardEdit, 0)
		}
		return
	}
	if row, ok := g.talkChoiceInput(&g.fundingRow, len(g.fundingChoiceLines(*c))); ok {
		g.fundingRow = row
		if row == int(state.FundingSetAmount) {
			g.beginFundingAmount()
			return
		}
		g.resolveFunding(*c, state.FundingOption(row))
		return
	}
	if pressed(ebiten.KeyArrowUp) {
		g.fundingRow = (g.fundingRow + 2) % 3
	}
	if pressed(ebiten.KeyArrowDown) {
		g.fundingRow = (g.fundingRow + 1) % 3
	}
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if pressed(k) {
			g.fundingRow = i
		}
	}
	if g.cancelled() {
		g.fundingRow = int(state.FundingReject)
	}
	if !g.confirmed() &&
		!g.cancelled() && !pressed(ebiten.Key1) &&
		!pressed(ebiten.Key2) && !pressed(ebiten.Key3) {
		return
	}

	option := state.FundingOption(g.fundingRow)
	if option == state.FundingSetAmount &&
		(g.confirmed()) {
		g.beginFundingAmount()
		return
	}
	g.resolveFunding(*c, option)
}

func (g *game) beginFundingAmount() {
	g.world.SetFundingAmount(0)
	g.fundingEditingAmount = true
	g.beginAmountEditor(amountCursorFunding, amountAnchorEventX, amountAnchorEventY)
}

func (g *game) resolveFunding(c state.FundingChoice, option state.FundingOption) {
	g.fundingRow = 0
	g.fundingEditingAmount = false
	g.amountCursorActive = false
	// sub_139E8 的選項 2（拒絕）與指定金額 0 都會回傳「未完成」類似
	// 的狀態結果，但兩者仍有各自已證實的 TALK 分支，不能落入泛用錯誤訊息。
	knownTalk := c.RequestedAmount == 0 || option == state.FundingReject ||
		(option == state.FundingSetAmount && c.OfferAmount == 0)
	g.enqueueFundingTalk(c, option)
	// 原版139AD9在官員最後答覆確認後才撥款；等待期間保留pending以暫停世界。
	g.afterTalk(func() {
		ok := g.world.ResolveFunding(option)
		if !ok && !knownTalk {
			g.setEvent(fundingTitle(c.Kind) + "：撥款未成立")
			return
		}
		g.lastEvent = ""
	})
}

// fundingTalkIndices 是 sub_139E8 的 TALK.DAT base+offset 映射。
// 事件 4 base=0x116（#278），事件 5 base=0x13F（#319）。指定金額
// 的 code 依序是：等於初值 0、低於初值 1、零 2、高於初值 3；
// 139A9F／139AA6／139A98同時改寫回應的兩個byte，收尾也使用正規化後的code。
func fundingTalkIndices(c state.FundingChoice, option state.FundingOption) (int, int) {
	base := 0x116
	if c.Kind == state.FundingDiplomat {
		base = 0x13F
	}
	if c.RequestedAmount == 0 {
		// 13A31–13A45：只確認官員通報，不進入三選一。
		return base + 35, base + 36
	}
	switch option {
	case state.FundingFullAmount:
		return base + 6, base + 10
	case state.FundingReject:
		return base + 8, base + 20
	case state.FundingSetAmount:
		amount := c.OfferAmount
		code := 0
		switch {
		case amount == 0:
			code = 2
		case amount < c.RequestedAmount:
			code = 1
		case amount > c.RequestedAmount:
			code = 3
		}
		return base + 6 + code, base + 10 + code*5
	default:
		return -1, -1
	}
}

func fundingTalkBase(c state.FundingChoice) int {
	if c.Kind == state.FundingDiplomat {
		return 0x13F
	}
	return 0x116
}

// fundingTalkPromptIndex 對應 sub_139E8 的非零要求路徑：
// scene page 1 先畫 base TALK，再以 base+5 畫三選一。beginFunding
// 對一般事件要求已套用 500 下限，因此這是目前可進入 pending UI 的
// 主要 prompt；要求為零的原版自動分支另保留在 state／TALK 索引表。
func fundingTalkPromptIndex(c state.FundingChoice) int {
	if c.RequestedAmount == 0 {
		return fundingTalkBase(c) + 30
	}
	return fundingTalkBase(c)
}

func (g *game) fundingOfficerPortrait(c state.FundingChoice) int {
	if g.world != nil && c.Officer >= 0 && c.Officer < len(g.world.Generals) {
		return g.world.Generals[c.Officer].Portrait
	}
	return g.playerLordPortrait()
}

func (g *game) fundingOfficerTalkIndex(c state.FundingChoice, index int) int {
	if g.world == nil || c.Officer < 0 || c.Officer >= len(g.world.Generals) {
		return index
	}
	raw := g.world.Generals[c.Officer].TalkVariant
	if raw < 0 || raw > 7 {
		return index
	}
	if raw >= 3 {
		raw -= 3 // 13C99只減一次；合法raw7應讀base+4。
	}
	return index + raw
}

func (g *game) fundingTalkVars(c state.FundingChoice, amount int) map[byte]string {
	notice := state.TalkNotice{City: -1, Faction: -1, General: -1, Amount: amount}
	if c.Kind == state.FundingGovernor {
		notice.City = c.Subject
	} else {
		notice.Faction = c.Subject
	}
	vars, _ := g.world.TalkNoticeVars(notice, big5)
	return vars
}

func (g *game) enqueueFundingTalk(c state.FundingChoice, option state.FundingOption) {
	first, second := fundingTalkIndices(c, option)
	if first < 0 || second < 0 {
		return
	}
	amount := c.RequestedAmount
	if option == state.FundingSetAmount {
		amount = c.OfferAmount
	}
	vars := g.fundingTalkVars(c, amount)
	prompt := g.diplomacySceneMessage(g.fundingOfficerTalkIndex(c, fundingTalkPromptIndex(c)),
		g.fundingTalkVars(c, c.RequestedAmount), g.fundingOfficerPortrait(c))
	if prompt != nil {
		prompt.upper, prompt.scene = true, 1
	}
	before := len(g.messages)
	g.enqueueAdvisorTalk(first, vars, 1)
	var advisor *messageDialog
	if len(g.messages) > before {
		g.messages[len(g.messages)-1].underlay = prompt
		copy := g.messages[len(g.messages)-1]
		advisor = &copy
	}
	if reply := g.diplomacySceneMessage(g.fundingOfficerTalkIndex(c, second), vars, g.fundingOfficerPortrait(c)); reply != nil {
		reply.upper, reply.scene, reply.underlay = true, 1, advisor
		g.messages = append(g.messages, *reply)
	}
}

func fundingTitle(kind state.FundingKind) string {
	if kind == state.FundingDiplomat {
		return "外交官撥款"
	}
	return "內政官撥款"
}

// fundingChoiceLines 是三選一的那三列。**不折行**——選單的字寬
// 決定框寬（docs/spec/45 §2.2）。
func (g *game) fundingChoiceLines(c state.FundingChoice) []string {
	lines, _ := g.talkLines(fundingTalkBase(c)+5, g.fundingTalkVars(c, c.RequestedAmount))
	return lines
}

// drawFunding 畫出事件 4／5 的原版 composite：IVENTGRF page 1、
// 官員／TALK、三列選項；指定金額列確認後才顯示原版數值器。
func (g *game) drawFunding(screen *ebiten.Image, c *state.FundingChoice) {
	g.drawIventScene(screen, 1)
	vars, fields := padTalkVars(g.fundingTalkVars(*c, c.RequestedAmount))
	prompt := g.legacyTalkLines(g.fundingOfficerTalkIndex(*c, fundingTalkPromptIndex(*c)), vars, talkTextWidth)
	g.drawLegacyTalkBoxFields(screen, talkUpperBoxX, talkUpperBoxY, talkBoxW, talkBoxH,
		prompt, g.fundingOfficerPortrait(*c), fields)
	if g.fundingEditingAmount {
		g.drawAmountPanel(screen, c.OfferAmount, true)
		return
	}
	g.drawLegacyChoiceBox(screen, talkChoiceX, talkChoiceY,
		g.fundingChoiceLines(*c), g.fundingRow)
}
