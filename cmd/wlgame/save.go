package main

// 玩家存檔流程。
//
// 原版的 SINARIO.DAT 同時承載四個劇本／槽位；remake 不在原始素材上就地
// 寫入，而是把目前 World 改寫到使用者明確指定的 overlay。這樣既能保留
// 尚未解出的位元組，也不會把自備的原版檔案變成工作產物。

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/savefile"
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

// 版面**全部出自原版**（docs/spec/25）：視窗矩形來自 `sub_1895D(cx=0F13h)`，
// 靜態層是顯示清單場景 6，數值座標由 `sub_18C20` 的 VRAM 位移換算。
const (
	savePanelX = 96
	savePanelY = 80
	savePanelW = 304
	savePanelH = 240

	saveTitleX, saveTitleY = 184, 91

	// 場景 6 的 op 05：一條橫跨內框的水平線，在標題與第一個槽之間。
	saveRuleX, saveRuleY, saveRuleW = 104, 111, 287

	// 名稱欄與日期欄。列距 48。
	saveNameBoxX, saveNameBoxY = 120, 118
	saveNameBoxW, saveNameBoxH = 256, 20
	saveNameX, saveNameY       = 120, 120
	saveSlotX, saveSlotY       = 256, 144
	saveSlotW, saveSlotH       = 120, 16
	saveSlotStep               = 48

	// 年、月、日各自嵌在「年　月　日」那三個字的前面。
	saveYearX, saveMonthX, saveDayX = 264, 304, 336
	saveDateLabelX                  = 288
	saveYearDigits                  = 3
	saveMonthDigits                 = 2

	// remake 差異：提示與兩個按鈕自己一個框，接在原版視窗下面。
	saveHintY = savePanelY + savePanelH
	saveHintH = 56
)

// saveSlotRect 是第 slot 個槽的可點矩形（原版熱區 0x20+slot），
// **與日期欄逐格重合**——左邊的名稱欄不可點。
func saveSlotRect(slot int) image.Rectangle {
	if slot < 0 || slot >= 4 {
		return image.Rectangle{}
	}
	y := saveSlotY + slot*saveSlotStep
	return image.Rect(saveSlotX, y, saveSlotX+saveSlotW, y+saveSlotH)
}

// saveFooterRect 是「確定」與「取消」——**原版沒有這兩顆**，
// 它按右鍵取消、點槽就決定。放在原版視窗外面的提示框裡。
func saveFooterRect(confirm bool) image.Rectangle {
	y := saveHintY + saveHintH - textdraw.GlyphH - 8
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
	// 原生檔（`docs/spec/20`）與上面那份原版格式**一起寫**。
	// 分成兩個動作就多一個「忘記匯出」的狀態，而原版格式正是這個專案的
	// 保存價值所在；一起寫沒有不一致的中間態。
	if err := g.writeNativeSave(slot); err != nil {
		return err
	}
	g.lastEvent = fmt.Sprintf("已儲存第 %d 槽", slot+1)
	return nil
}

// writeNativeSave 寫這一槽的原生存檔。
func (g *game) writeNativeSave(slot int) error {
	path, err := savepath.NativePath(g.saveFile, slot)
	if err != nil {
		return err
	}
	data, err := savefile.Encode(g.world, savefile.Origin{
		Source: filepath.Base(g.saveFile), Block: slot,
	})
	if err != nil {
		return fmt.Errorf("編碼原生存檔失敗：%w", err)
	}
	return writeFileAtomically(path, data)
}

// writeFileAtomically 用同目錄暫存檔加改名，避免中途停止留下半個檔案。
func writeFileAtomically(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".wlgame-native-*")
	if err != nil {
		return fmt.Errorf("建立暫存檔失敗：%w", err)
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("寫入暫存檔失敗：%w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("關閉暫存檔失敗：%w", err)
	}
	return os.Rename(name, path)
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
	// 優先讀原生檔——原版格式裝不下 routes／游標那些欄位，
	// 從它讀回來的世界會少掉一部分執行期狀態（`docs/spec/20` §2.4）。
	w, native, err := g.readNativeSave(slot)
	if err != nil {
		return err
	}
	if w == nil {
		if w, err = state.LoadScenario(g.saveFile, slot); err != nil {
			return fmt.Errorf("讀取第 %d 槽失敗：%w", slot+1, err)
		}
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
	kind := "原版格式"
	if native {
		kind = "原生檔"
	}
	g.lastEvent = fmt.Sprintf("已讀取第 %d 槽（%s）；信賴度 %d", slot+1, kind, w.Trust)
	return nil
}

// readNativeSave 讀這一槽的原生存檔。檔案不存在時回 (nil, false, nil)，
// 讓呼叫端退回原版格式；**存在但壞掉一律回錯誤**——
// 靜靜退回去會讓玩家以為讀到了最新進度，實際上少了執行期狀態。
func (g *game) readNativeSave(slot int) (*state.World, bool, error) {
	path, err := savepath.NativePath(g.saveFile, slot)
	if err != nil {
		return nil, false, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("讀取原生存檔失敗：%w", err)
	}
	w, _, err := savefile.Decode(data)
	if err != nil {
		return nil, false, fmt.Errorf("第 %d 槽的原生存檔有問題：%w", slot+1, err)
	}
	return w, true, nil
}

func (g *game) drawSaveUI(screen *ebiten.Image) {
	if !g.saveUI.active {
		return
	}
	g.chrome.Window(screen, savePanelX, savePanelY, savePanelW, savePanelH, chrome.Menu)

	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})
	dateInk := g.paletteInk(0x05, color.RGBA{200, 200, 255, 255})
	warnInk := g.paletteInk(strategyInkGauge, color.RGBA{210, 48, 40, 255})
	amber := color.RGBA{240, 200, 120, 255}

	// 標題：原版是 ＮＥＷ　ＧＡＭＥ／ＬＯＡＤ　ＤＡＴＡ／ＳＡＶＥ　ＤＡＴＡ
	// 三選一（docs/re/52 §2）。remake 這個視窗只做讀取與儲存。
	title := "ＳＡＶＥ　ＤＡＴＡ"
	if g.saveUI.action == saveRead {
		title = "ＬＯＡＤ　ＤＡＴＡ"
	}
	g.td.Draw(screen, title, saveTitleX, saveTitleY, ink)
	vector.DrawFilledRect(screen, saveRuleX, saveRuleY, saveRuleW, 1, ink, false)

	for i := 0; i < 4; i++ {
		dy := i * saveSlotStep
		vector.DrawFilledRect(screen, saveNameBoxX, float32(saveNameBoxY+dy),
			saveNameBoxW, saveNameBoxH, color.Black, false)
		// 日期欄是**凹槽**：底色 5、外圈 2／0、內圈 D／4（docs/re/48 §2.1）。
		g.dlSunken(screen, saveSlotX, saveSlotY+dy, saveSlotW, saveSlotH)
		g.td.Draw(screen, "年　月　日", saveDateLabelX, saveSlotY+dy, ink)

		slot, ok := launcherSlot{}, false
		if i < len(g.saveUI.slots) {
			slot, ok = g.saveUI.slots[i], g.saveUI.slots[i].Available
		}
		nameInk := labelInk
		if i == g.saveUI.slot {
			nameInk = amber
		}
		if !ok {
			g.td.Draw(screen, "空白槽位", saveNameX, saveNameY+dy, warnInk)
			continue
		}
		g.td.Draw(screen, strategyHUDSingleLine(slot.Title, saveNameBoxW-16),
			saveNameX, saveNameY+dy, nameInk)
		g.td.Draw(screen, strategyHUDNumber(slot.Year, saveYearDigits),
			saveYearX, saveSlotY+dy, dateInk)
		g.td.Draw(screen, strategyHUDNumber(slot.Month, saveMonthDigits),
			saveMonthX, saveSlotY+dy, dateInk)
		g.td.Draw(screen, strategyHUDNumber(slot.Day, saveMonthDigits),
			saveDayX, saveSlotY+dy, dateInk)
	}

	// ↓ remake 差異：原版點槽就決定、右鍵取消，沒有這一框。
	g.chrome.Window(screen, savePanelX, saveHintY, savePanelW, saveHintH, chrome.Menu)
	hx, hy := savePanelX+8, saveHintY+8
	if g.saveFile == "" {
		g.td.Draw(screen, "未指定 -save-file（目前停用）", hx, hy, warnInk)
	} else {
		g.td.Draw(screen, "檔案　"+filepath.Base(g.saveFile), hx, hy, labelInk)
	}
	footerY := saveHintY + saveHintH - textdraw.GlyphH - 8
	g.td.Draw(screen, "↑↓／1-4 選擇", hx, footerY, labelInk)
	g.td.Draw(screen, "確定", saveFooterRect(true).Min.X+16, footerY, ink)
	g.td.Draw(screen, "取消", saveFooterRect(false).Min.X+16, footerY, ink)
}
