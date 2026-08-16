package main

import (
	"fmt"

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
		if pressed(ebiten.KeyEscape) {
			g.fundingEditingAmount = false
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
			keyboardEdit = state.AmountRestoreInitial
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
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
	if row, ok := g.talkChoiceClick(talkChoiceX, talkChoiceY,
		g.fundingChoiceLines(*c)); ok {
		g.fundingRow = row
		if row == int(state.FundingSetAmount) {
			g.fundingEditingAmount = true
			g.beginAmountEditor(amountCursorFunding)
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
	if pressed(ebiten.KeyEscape) {
		g.fundingRow = int(state.FundingReject)
	}
	if !pressed(ebiten.KeyEnter) && !pressed(ebiten.KeySpace) &&
		!pressed(ebiten.KeyEscape) && !pressed(ebiten.Key1) &&
		!pressed(ebiten.Key2) && !pressed(ebiten.Key3) {
		return
	}

	option := state.FundingOption(g.fundingRow)
	if option == state.FundingSetAmount &&
		(pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace)) {
		g.fundingEditingAmount = true
		g.beginAmountEditor(amountCursorFunding)
		return
	}
	g.resolveFunding(*c, option)
}

func (g *game) resolveFunding(c state.FundingChoice, option state.FundingOption) {
	kind := c.Kind
	ok := g.world.ResolveFunding(option)
	g.fundingRow = 0
	g.fundingEditingAmount = false
	g.amountCursorActive = false
	// sub_139E8 的選項 2（拒絕）與指定金額 0 都會回傳「未完成」類似
	// 的狀態結果，但兩者仍有各自已證實的 TALK 分支，不能落入泛用錯誤訊息。
	knownTalk := option == state.FundingReject ||
		(option == state.FundingSetAmount && c.OfferAmount == 0)
	if ok || knownTalk {
		g.enqueueFundingTalk(c, option)
	}
	if option == state.FundingReject {
		g.setEvent(fundingTitle(kind) + "：拒絕")
		return
	}
	if !ok {
		g.setEvent(fundingTitle(kind) + "：撥款未成立")
		return
	}
	if option == state.FundingSetAmount {
		g.setEvent(fundingTitle(kind) + "：指定金額撥款")
	} else {
		g.setEvent(fundingTitle(kind) + "：全額撥款")
	}
}

// fundingTalkIndices 是 sub_139E8 的 TALK.DAT base+offset 映射。
// 事件 4 base=0x116（#278），事件 5 base=0x13F（#319）。指定金額
// 的 code 依序是：等於初值 0、低於初值 1、零 2、高於初值 3；
// 原版的收尾則只看三選一的原始列號。
func fundingTalkIndices(c state.FundingChoice, option state.FundingOption) (int, int) {
	base := 0x116
	if c.Kind == state.FundingDiplomat {
		base = 0x13F
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
		return base + 6 + code, base + 15
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
	return fundingTalkBase(c)
}

func (g *game) fundingTalkVars(c state.FundingChoice, amount int) map[byte]string {
	vars := map[byte]string{'6': "", '7': fmt.Sprintf("%d", amount)}
	if c.Kind == state.FundingGovernor {
		if c.Subject >= 0 && c.Subject < len(g.world.Cities) {
			vars['2'] = big5(g.world.Cities[c.Subject].Name)
		}
		if g.world.Player >= 0 && g.world.Player < len(g.world.Factions) {
			advisor := g.world.Factions[g.world.Player].Advisor
			if advisor >= 0 && advisor < len(g.world.Generals) && g.world.Generals[advisor].Alive {
				vars['4'] = big5(g.world.Generals[advisor].Name)
			}
		}
	} else if c.Subject >= 0 && c.Subject < len(g.world.Factions) {
		vars['3'] = big5(g.world.LordName(c.Subject))
	}
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
	g.enqueueTalk(first, vars)
	g.enqueueTalk(second, vars)
}

func fundingTitle(kind state.FundingKind) string {
	if kind == state.FundingDiplomat {
		return "外交官撥款"
	}
	return "內政官撥款"
}

func (g *game) fundingOfficerName(id int) string {
	if id < 0 || id >= len(g.world.Generals) || !g.world.Generals[id].Alive {
		return "－"
	}
	return big5(g.world.Generals[id].Name)
}

func (g *game) fundingSubjectName(c *state.FundingChoice) string {
	if c == nil {
		return "－"
	}
	if c.Kind == state.FundingGovernor {
		if c.Subject >= 0 && c.Subject < len(g.world.Cities) {
			return big5(g.world.Cities[c.Subject].Name)
		}
		return "－"
	}
	if c.Subject >= 0 && c.Subject < len(g.world.Factions) {
		return g.diplomacyFactionName(c.Subject)
	}
	return "－"
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
	vars := g.fundingTalkVars(*c, c.RequestedAmount)
	prompt := g.legacyTalkLines(fundingTalkPromptIndex(*c), vars, talkTextWidth)
	portrait := -1
	if c.Officer >= 0 && c.Officer < len(g.world.Generals) {
		portrait = g.world.Generals[c.Officer].Portrait
	} else {
		portrait = g.playerLordPortrait()
	}
	g.drawLegacyTalkBox(screen, talkUpperBoxX, talkUpperBoxY, talkBoxW, talkBoxH,
		prompt, portrait)
	if g.fundingEditingAmount {
		g.drawAmountPanel(screen, c.OfferAmount, c.RequestedAmount, true)
		return
	}
	g.drawLegacyChoiceBox(screen, talkChoiceX, talkChoiceY,
		g.fundingChoiceLines(*c), g.fundingRow)
}
