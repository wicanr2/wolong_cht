package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

var outcomePanel = image.Rect(112, 120, 528, 280)

func outcomeConfirmRect() image.Rectangle {
	return image.Rect(outcomePanel.Min.X+16, outcomePanel.Max.Y-48,
		outcomePanel.Max.X-16, outcomePanel.Max.Y-16)
}

func outcomeReason(kind state.OutcomeKind) string {
	switch kind {
	case state.DefeatTrustZero:
		return "信賴度歸零，已被逐出勢力。"
	case state.DefeatFactionEliminated:
		return "玩家勢力失去最後可替代據點。"
	default:
		return "遊戲已停止。"
	}
}

// outcomeLines 優先沿用已證實 selector 0x019E 對應的 DOS/V TALK.DAT 訊息。
// 只有 TALK 資產不存在、索引失效或 marker/payload 無法安全代入時，才退回
// 克制的 remake 原因句；不能在 GUI 顯示研究判定或反組譯備註。
func (g *game) outcomeLines() []string {
	if g != nil && g.world != nil && g.world.Outcome() == state.DefeatTrustZero {
		if lines, ok := g.talkLines(0x019E, map[byte]string{}); ok && len(lines) > 0 {
			return lines
		}
	}
	return []string{outcomeReason(g.world.Outcome())}
}

// updateOutcome 是 outcome latch 後唯一接受的遊戲內輸入：Enter／Space 或
// modal 內左鍵確認。它不重開同一局；回 launcher 是 remake 的呈現政策，
// 原版 sub_11CB1 之後是 unknown。
func (g *game) updateOutcome() error {
	if g == nil || g.world == nil || g.world.Outcome() == state.InProgress {
		return nil
	}
	if pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) {
		return g.returnToLauncher()
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if image.Pt(x, y).In(outcomeConfirmRect()) {
			return g.returnToLauncher()
		}
	}
	return nil
}

func (g *game) returnToLauncher() error {
	if g == nil {
		return nil
	}
	// 這裡只清理 runtime 指標；不改寫原始素材，也不把敗北局自動
	// 當成新局。launcher 的正式新局／讀檔流程仍由 launcher.go 負責。
	slots := inspectLauncherSlots(g.saveFile)
	g.world = nil
	g.roads = nil
	g.tactical = nil
	g.battleLib = nil
	g.battleSprites = nil
	g.view = nil
	g.open = [numWindows]bool{}
	g.list = nil
	g.form = formState{}
	g.finance = financeState{}
	g.advise = adviseNone
	g.sess = nil
	g.messages = nil
	g.saveUI = saveUIState{}
	g.lastEvent = ""
	g.quitting = false
	g.idleGate = idleClockGate{}
	g.launcher = newLauncher(hasAvailableLauncherSlot(slots), slots)
	return nil
}

func (g *game) drawOutcome(screen *ebiten.Image) {
	if g == nil || g.world == nil || g.world.Outcome() == state.InProgress {
		return
	}
	vector.DrawFilledRect(screen, 0, 0, screenW, screenH,
		color.RGBA{0, 0, 0, 150}, false)
	if g.chrome != nil {
		g.chrome.Window(screen, outcomePanel.Min.X, outcomePanel.Min.Y,
			outcomePanel.Dx(), outcomePanel.Dy(), chrome.Menu)
	} else {
		vector.DrawFilledRect(screen, float32(outcomePanel.Min.X), float32(outcomePanel.Min.Y),
			float32(outcomePanel.Dx()), float32(outcomePanel.Dy()), color.RGBA{32, 24, 16, 255}, false)
	}
	if g.td == nil {
		return
	}
	white := chrome.Paper
	amber := color.RGBA{240, 200, 120, 255}
	g.td.Draw(screen, "敗北", outcomePanel.Min.X+32, outcomePanel.Min.Y+24, amber)
	for i, line := range g.outcomeLines() {
		g.td.Draw(screen, line, outcomePanel.Min.X+32,
			outcomePanel.Min.Y+56+i*24, white)
	}
	vector.DrawFilledRect(screen, float32(outcomeConfirmRect().Min.X),
		float32(outcomeConfirmRect().Min.Y), float32(outcomeConfirmRect().Dx()),
		float32(outcomeConfirmRect().Dy()), chrome.Select, false)
	g.td.Draw(screen, "Enter／滑鼠左鍵：回到啟動畫面",
		outcomeConfirmRect().Min.X+16, outcomeConfirmRect().Min.Y+8, white)
}
