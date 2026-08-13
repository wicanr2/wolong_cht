package main

// 玩家存檔流程。
//
// 原版的 SINARIO.DAT 同時承載四個劇本／槽位；remake 不在原始素材上就地
// 寫入，而是把目前 World 改寫到使用者明確指定的 overlay。這樣既能保留
// 尚未解出的位元組，也不會把自備的原版檔案變成工作產物。

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/savepath"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

type saveAction uint8

const (
	saveWrite saveAction = iota + 1
	saveRead
)

type saveUIState struct {
	active bool
	action saveAction
	slot   int
	slots  []launcherSlot
}

func (g *game) beginSaveUI(action saveAction) {
	g.saveUI = saveUIState{active: true, action: action}
	if action == saveRead {
		g.saveUI.slots = inspectLauncherSlots(g.saveFile)
	}
}

type saveUIActionKind uint8

const (
	saveActionSelect saveUIActionKind = iota + 1
	saveActionPrev
	saveActionNext
	saveActionConfirm
	saveActionCancel
)

type saveUIAction struct {
	kind saveUIActionKind
	slot int
}

// dispatchSaveUI 是存讀檔 modal 的單一輸入分派器。鍵盤數字鍵、方向鍵、
// 滑鼠與 Android 觸控映射都先轉成同一種 action，避免產生兩套選槽語意。
func (g *game) dispatchSaveUI(action saveUIAction) {
	if !g.saveUI.active {
		return
	}
	switch action.kind {
	case saveActionSelect:
		if action.slot >= 0 && action.slot < 4 {
			g.saveUI.slot = action.slot
		}
	case saveActionPrev:
		g.saveUI.slot = (g.saveUI.slot + 3) % 4
	case saveActionNext:
		g.saveUI.slot = (g.saveUI.slot + 1) % 4
	case saveActionCancel:
		g.saveUI = saveUIState{}
	case saveActionConfirm:
		var err error
		if g.saveUI.action == saveWrite {
			err = g.writeSave(g.saveUI.slot)
		} else {
			err = g.readSave(g.saveUI.slot)
		}
		if err != nil {
			g.lastEvent = err.Error()
			return
		}
		g.saveUI = saveUIState{}
	}
}

const (
	savePanelX = 112
	savePanelY = 72
	savePanelW = 416
	savePanelH = 248
	saveSlotX  = savePanelX + chrome.Tile + 4
	saveSlotY  = savePanelY + chrome.Tile + 2 + 2*(textdraw.GlyphH+4) + 8
	saveSlotW  = savePanelW - 2*chrome.Tile - 8
	saveSlotH  = textdraw.GlyphH + 4
)

func saveSlotRect(slot int) image.Rectangle {
	if slot < 0 || slot >= 4 {
		return image.Rectangle{}
	}
	return image.Rect(saveSlotX, saveSlotY+slot*(textdraw.GlyphH+4),
		saveSlotX+saveSlotW, saveSlotY+(slot+1)*(textdraw.GlyphH+4))
}

func saveFooterRect(confirm bool) image.Rectangle {
	y := savePanelY + savePanelH - chrome.Tile - textdraw.GlyphH
	if confirm {
		return image.Rect(savePanelX+savePanelW-144, y,
			savePanelX+savePanelW-72, y+textdraw.GlyphH+2)
	}
	return image.Rect(savePanelX+savePanelW-72, y,
		savePanelX+savePanelW, y+textdraw.GlyphH+2)
}

func (g *game) updateSaveUI() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.dispatchSaveUI(saveUIAction{kind: saveActionCancel})
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p := image.Point{X: x, Y: y}
		for slot := 0; slot < 4; slot++ {
			if p.In(saveSlotRect(slot)) {
				g.dispatchSaveUI(saveUIAction{kind: saveActionSelect, slot: slot})
				return
			}
		}
		if p.In(saveFooterRect(true)) {
			g.dispatchSaveUI(saveUIAction{kind: saveActionConfirm})
			return
		}
		if p.In(saveFooterRect(false)) {
			g.dispatchSaveUI(saveUIAction{kind: saveActionCancel})
			return
		}
		return
	}
	switch {
	case pressed(ebiten.KeyArrowUp):
		g.dispatchSaveUI(saveUIAction{kind: saveActionPrev})
	case pressed(ebiten.KeyArrowDown):
		g.dispatchSaveUI(saveUIAction{kind: saveActionNext})
	case pressed(ebiten.KeyEscape):
		g.dispatchSaveUI(saveUIAction{kind: saveActionCancel})
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		g.dispatchSaveUI(saveUIAction{kind: saveActionConfirm})
	default:
		for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
			if pressed(k) {
				g.dispatchSaveUI(saveUIAction{kind: saveActionSelect, slot: i})
				return
			}
		}
	}
}

// writeSave 以暫存檔加同目錄改名完成一次儲存，避免中途停止留下半個
// SAVE.DAT。SaveInto 本身仍負責只覆寫已解欄位，其餘位元組照原檔保留。
func (g *game) writeSave(slot int) error {
	if g.saveFile == "" {
		return fmt.Errorf("未指定 -save-file，儲存已停用")
	}
	if savepath.SamePath(g.saveFile, g.sourceFile) {
		return fmt.Errorf("拒絕覆寫原始素材：%s", g.saveFile)
	}
	dir := filepath.Dir(g.saveFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("建立存檔目錄失敗：%w", err)
	}
	tmp, err := os.CreateTemp(dir, ".wlgame-save-*")
	if err != nil {
		return fmt.Errorf("建立存檔暫存檔失敗：%w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("關閉存檔暫存檔失敗：%w", err)
	}
	if err := g.world.SaveInto(g.saveBase, tmpName, slot); err != nil {
		return fmt.Errorf("儲存第 %d 槽失敗：%w", slot+1, err)
	}
	if err := os.Rename(tmpName, g.saveFile); err != nil {
		return fmt.Errorf("完成存檔改名失敗：%w", err)
	}
	g.saveBase = g.saveFile
	g.lastEvent = fmt.Sprintf("已儲存第 %d 槽", slot+1)
	return nil
}

// readSave 讀取後重新掛回所有執行期來源。
// 有效原版存檔會在全域區段 +0x0D／+0x0F 保存 Player；空槽／新劇本仍是
// 無玩家哨兵，此時才沿用目前執行期的勢力。Trust 位於 +0x10（IDA
// `byte_10D00`），由 LoadScenario 一起讀回。
func (g *game) readSave(slot int) error {
	if g.saveFile == "" {
		return fmt.Errorf("未指定 -save-file，讀取已停用")
	}
	if savepath.SamePath(g.saveFile, g.sourceFile) {
		return fmt.Errorf("拒絕把原始素材當成存檔讀取：%s", g.saveFile)
	}
	w, err := state.LoadScenario(g.saveFile, slot)
	if err != nil {
		return fmt.Errorf("讀取第 %d 槽失敗：%w", slot+1, err)
	}
	player := w.Player
	if player < 0 {
		player = g.world.Player
		w.Player = player
	}
	w.SetRoads(g.roads)
	w.SetTactical(g.tactical)
	g.world = w
	g.view = nil
	g.saveBase = g.saveFile
	if player >= 0 && player < len(w.Factions) {
		cap := w.Factions[player].Capital
		if cap >= 0 && cap < len(w.Cities) {
			g.camX, g.camY = w.Cities[cap].X-viewCols/2, w.Cities[cap].Y-viewRows/2
			g.clampCam()
		}
	}
	g.lastEvent = fmt.Sprintf("已讀取第 %d 槽；信賴度 %d", slot+1, w.Trust)
	return nil
}

func (g *game) drawSaveUI(screen *ebiten.Image) {
	if !g.saveUI.active {
		return
	}
	const x, y, w, h = savePanelX, savePanelY, savePanelW, savePanelH
	g.chrome.Window(screen, x, y, w, h, chrome.Menu)
	amber := color.RGBA{240, 200, 120, 255}
	white := chrome.Paper
	dim := color.RGBA{170, 170, 180, 255}
	red := color.RGBA{240, 140, 140, 255}
	tx := x + chrome.Tile + 4
	ty := y + chrome.Tile + 2
	title := "儲存資料"
	if g.saveUI.action == saveRead {
		title = "讀取資料"
	}
	g.td.Draw(screen, title, tx, ty, amber)
	fileText := "未指定 -save-file（目前停用）"
	fileCol := red
	if g.saveFile != "" {
		fileText = "檔案　" + filepath.Base(g.saveFile)
		fileCol = dim
	}
	g.td.Draw(screen, fileText, tx, ty+textdraw.GlyphH+4, fileCol)

	ry := ty + 2*(textdraw.GlyphH+4) + 8
	for i := 0; i < 4; i++ {
		col := white
		mark := "　"
		if i == g.saveUI.slot {
			col, mark = amber, "●"
		}
		label := fmt.Sprintf("%s%d　劇本／槽位 %d", mark, i+1, i+1)
		if g.saveUI.action == saveRead && (i >= len(g.saveUI.slots) || !g.saveUI.slots[i].Available) {
			col = dim
			label = fmt.Sprintf("%s%d　空白槽位", mark, i+1)
		}
		g.td.Draw(screen, label, tx, ry, col)
		ry += textdraw.GlyphH + 4
	}
	footerY := y + h - chrome.Tile - textdraw.GlyphH
	g.td.Draw(screen, "↑↓／1-4 選擇", tx, footerY, dim)
	g.td.Draw(screen, "確定", saveFooterRect(true).Min.X+16, footerY, white)
	g.td.Draw(screen, "取消", saveFooterRect(false).Min.X+16, footerY, white)
}
