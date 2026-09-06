package main

// 指令列第 6 格「武將」與第 7 格「勢力」（`sub_16366`／`sub_163BF`，
// docs/spec/145）。兩格都有狀態列提示，而且**選完之後還有下一步**。

const (
	// 兩格各自的狀態列（`mov cx, 18h`／`19h` → `sub_18853`）。
	generalCellTalk = 0x18 // #24「確認指示之武將的能力。」
	factionCellTalk = 0x19 // #25「將游標移動至指示之勢力的首都據點。」

	// 被點到的武將自陳擅長的戰場：四組八格，由 +0x1E 選組內第幾個。
	captiveRefusesTalk = 0x1A8 // #550–557 俘虜的推託
	siegeBoastTalk     = 0x1A9 // #558–565 城塞戰
	fieldBoastTalk     = 0x1AA // #566–573 野戰
	navalBoastTalk     = 0x1AB // #574–581 海戰
)

// generalAptitudeTalk 挑那位武將要說哪一組（`loc_16380`）。
//
// 俘虜（`+0x1D ≠ 0xFF`）優先；否則取三個適性的最大值，再依序問
// 「等於攻城嗎 → 等於野戰嗎 → 否則水戰」。
// ⭐ **平手歸前面那一個**，所以三個一樣高時說城塞戰。
func (g *game) generalAptitudeTalk(who int) int {
	if g == nil || g.world == nil || who < 0 || who >= len(g.world.Generals) {
		return siegeBoastTalk
	}
	gen := &g.world.Generals[who]
	if gen.Captor != NoOfficial {
		return captiveRefusesTalk
	}
	siege, field, naval := gen.Aptitude[0], gen.Aptitude[1], gen.Aptitude[2]
	best := siege
	if field > best {
		best = field
	}
	if naval > best {
		best = naval
	}
	switch {
	case best == siege:
		return siegeBoastTalk
	case best == field:
		return fieldBoastTalk
	default:
		return navalBoastTalk
	}
}
