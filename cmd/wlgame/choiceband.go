package main

// 原版文字選列器，見 spec/267、re/84 與已驗證的 sub_10414。
import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wicanr2/wolong_cht/internal/state"
)

const choiceBandTop, choiceBandBottom = 100, 131

type choiceBand struct {
	owner          *int
	count          int
	savedX, savedY int
}

// 一次輪詢只移動一列，端點夾住但仍回置滑鼠；不是依畫面列座標選取。
func choiceBandMove(row, count, y int) (next, pointerY int, warp bool) {
	if count <= 0 {
		return 0, choiceBandTop, false
	}
	row = max(0, min(count-1, row))
	switch {
	case y > choiceBandBottom:
		return min(row+1, count-1), choiceBandTop, true
	case y < choiceBandTop:
		return max(row-1, 0), choiceBandBottom, true
	default:
		return row, y, false
	}
}

// 與 Update 的輸入優先序一致，只列出走 sub_10414 的文字選單。
func (g *game) currentChoiceOwner() (*int, int) {
	if g.quitMenu.active {
		return &g.quitMenu.row, len(g.quitMenuRows())
	}
	if g.quitting || g.launcher != nil || g.messageActive() || g.endingActive() ||
		g.world != nil && (g.battleActive() || g.world.Outcome() != state.InProgress) {
		return nil, 0
	}
	if g.world != nil {
		if c := g.world.PendingDiplomacy(); c != nil {
			if !g.diplomacyEditingAmount {
				return &g.diplomacyRow, len(g.diplomacyChoiceLines(*c))
			}
			return nil, 0
		}
		if c := g.world.PendingFunding(); c != nil {
			if !g.fundingEditingAmount {
				return &g.fundingRow, len(g.fundingChoiceLines(*c))
			}
			return nil, 0
		}
	}
	if g.adviseActive() {
		switch g.advise {
		case advisePickCommand:
			return &g.adviseCmdRow, len(g.adviseCommandLabels())
		case advisePersuade:
			if !g.adviseTalking() && g.sess != nil {
				return &g.sessCur, len(g.adviseReasonLabels(g.adviseCmd))
			}
		}
		return nil, 0
	}
	if g.popupMenuActive() {
		return &g.cmdMenu.row, len(g.popupMenuLabels(g.cmdMenu.menu))
	}
	if g.marchMode.active {
		return &g.marchMode.row, len(g.marchModeRows())
	}
	if g.mapChoice.active {
		return &g.mapChoice.row, len(g.mapChoiceRows())
	}
	return nil, 0
}

func (g *game) choiceOwnsInput() bool {
	owner, count := g.currentChoiceOwner()
	return owner != nil && count > 0
}

func (g *game) beginChoiceBand(owner *int, count int) bool {
	if owner == nil || count <= 0 {
		return false
	}
	if g.choiceBand.owner == owner && g.choiceBand.count == count {
		return false
	}
	g.endChoiceBand()
	x, y := cursorPosition()
	g.choiceBand = choiceBand{owner: owner, count: count, savedX: x, savedY: y}
	*owner = 0 // sub_1041A：每次建立選列器從第一列開始。
	if desktopPointer.active {
		desktopPointer.x, desktopPointer.y = 0, choiceBandTop
	}
	return true
}

func (g *game) endChoiceBand() {
	if g.choiceBand.owner != nil && desktopPointer.active {
		desktopPointer.x, desktopPointer.y = g.choiceBand.savedX, g.choiceBand.savedY
	}
	g.choiceBand = choiceBand{}
}

// 在 Update 結束時同步，讓新選單的第一張畫面就隱藏並保存游標。
func (g *game) syncChoiceBand() {
	owner, count := g.currentChoiceOwner()
	if owner == nil || count <= 0 {
		g.endChoiceBand()
		return
	}
	g.beginChoiceBand(owner, count)
}

func (g *game) talkChoiceInput(selected *int, count int) (int, bool) {
	if selected == nil || count <= 0 {
		g.endChoiceBand()
		return 0, false
	}
	fresh := g.beginChoiceBand(selected, count)
	// 原版先處理左鍵，再問位置；父畫面按下那一下不確認新選單。
	if g.confirmed() || !fresh && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		row := *selected
		g.endChoiceBand()
		return row, true
	}
	if g.cancelled() {
		g.endChoiceBand()
		return 0, false
	}
	if fresh || !desktopPointer.active || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		return 0, false
	}
	row, y, warp := choiceBandMove(*selected, count, desktopPointer.y)
	*selected = row
	if warp {
		desktopPointer.x, desktopPointer.y = 0, y
	}
	return 0, false
}
