// Package library 把原版素材目錄載成一組可檢視的項目。
//
// **這一層刻意不 import Ebiten。** Ebiten 在 init 期就要求顯示器，
// 混進來會讓無頭環境連截圖工具都跑不起來 —— 這個坑實際踩過一次：
// 第一版把截圖功能和 Ebiten 檢視器放同一個 package，結果
// `-shot` 模式在容器裡直接 panic 在 ebiten 的 init。
package library

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/palette"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
)

// Entry 是一種可檢視的素材。
type Entry struct {
	Label string
	Spec  gfx.Spec
	Data  []byte
	Count int
}

// Library 是一個原版素材目錄解出來的全部東西。
type Library struct {
	Palette *palette.Palette
	Talk    *text.Table
	Entries []Entry
	Warns   []string // 餘數不是 0 之類的警告，呼叫端要顯示出來

	// cursorPixels 是 DOS/V KI.EXE seg002:031B 的兩層游標 mask；
	// 載入時只保留已解出的 palette index，不把原始執行檔散布進引擎。
	cursorPixels []byte

	// Chrome 是 ICONGRF 段 3 的原始位元組——視窗外框的圖塊在裡面。
	// 用 gfx.RenderChrome ＋ gfx.ChromeEdge／ChromeCap／ChromeShaft 取。
	Chrome []byte

	// BattleUI 是 ICONGRF 段 1 的原始位元組。只由已定位的戰術指令
	// Spec／offset 解碼，不把混合尺寸的整段冒充固定圖集。
	BattleUI []byte

	// World 與 Tiles 是大地圖。原版是 384×256 格、每格 16×16 px，
	// 整張攤開是 6144×4096 —— 不預先畫成一張圖，畫面要多少畫多少。
	World *world.Map
	Tiles *world.TileSet
	// MCH 是 MMAP.MCH 的 strategic-map 物件圖塊與物件矩陣。
	// 事件 12 的火災／暴動呈現由它取用；未載入時畫面層會保留明確的
	// fallback marker，不把缺資產靜默當成 parity。
	MCH *world.MCH
}

// LoadOptions 控制呈現層的可選文字覆蓋。
// TalkJSON 是 tools/talkdat.py correct 產出的繁中訊息表；空字串表示只用
// 原始 TALK.DAT。原始素材仍只從 dir 讀取，翻譯覆蓋不會改寫它。
type LoadOptions struct {
	// TalkJSON 是預先產生、完整的呈現文字表。它適合研究與驗收，使用時必須由
	// 呼叫端明確提供；不能把它當成可隨發行包散布的原版文字來源。
	TalkJSON string
	// TalkCorrections 是可散布的最小校訂覆蓋表。它只記錄已定案的差異，會在
	// 本機讀到玩家自備的 raw TALK.DAT 後套用，避免發行包攜出完整原版文字表。
	TalkCorrections string
}

func read(dir, name string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return nil, fmt.Errorf("讀不到 %s：%w\n"+
			"原版素材不隨本專案散布，請用 -orig 指向自備的原版目錄", name, err)
	}
	return b, nil
}

// Load 載入一個原版素材目錄（dosv 或 pc98 都可以）。
func Load(dir string) (*Library, error) {
	return LoadWithOptions(dir, LoadOptions{})
}

// LoadWithOptions 載入原版素材，必要時以版控的呈現文字覆蓋 TALK 表。
func LoadWithOptions(dir string, opts LoadOptions) (*Library, error) {
	palRaw, err := read(dir, "GAMEPAL.BRG")
	if err != nil {
		return nil, err
	}
	pal, err := palette.Parse(palRaw)
	if err != nil {
		return nil, err
	}
	talkRaw, err := read(dir, "TALK.DAT")
	if err != nil {
		return nil, err
	}
	talk, err := text.Parse(talkRaw)
	if err != nil {
		return nil, err
	}
	if opts.TalkJSON != "" && opts.TalkCorrections != "" {
		return nil, fmt.Errorf("TALK 覆蓋設定衝突：不可同時指定完整 JSON 與校訂覆蓋")
	}
	if opts.TalkJSON != "" {
		talk, err = text.LoadJSON(opts.TalkJSON, text.Big5)
		if err != nil {
			return nil, fmt.Errorf("載入校訂 TALK 表失敗：%w", err)
		}
	} else if opts.TalkCorrections != "" {
		talk, err = text.ApplyCorrections(talk, opts.TalkCorrections, text.Big5)
		if err != nil {
			return nil, fmt.Errorf("套用 TALK 校訂覆蓋失敗：%w", err)
		}
	}

	lib := &Library{Palette: pal, Talk: talk}
	add := func(label string, spec gfx.Spec, data []byte) {
		n, rem := spec.Count(data)
		if rem != 0 {
			// 餘數不是 0 就代表尺寸錯了。不靜悄悄吞掉——
			// 那正是「看起來能跑但其實解錯」的來源。
			lib.Warns = append(lib.Warns, fmt.Sprintf(
				"%s 餘 %d byte，尺寸 %dx%d 可能是錯的",
				label, rem, spec.Width, spec.Height))
		}
		lib.Entries = append(lib.Entries, Entry{label, spec, data, n})
	}

	for _, c := range []struct {
		file string
		spec gfx.Spec
		name string
	}{
		{"KAOGRF.DAT", gfx.Kao, "KAOGRF portraits"},
		{"KYOGRF.DAT", gfx.Kyo, "KYOGRF locations"},
		{"IVENTGRF.DAT", gfx.Ivent, "IVENTGRF cutscenes"},
	} {
		data, err := read(dir, c.file)
		if err != nil {
			return nil, err
		}
		add(c.name, c.spec, data)
	}

	icon, err := read(dir, "ICONGRF.DAT")
	if err != nil {
		return nil, err
	}
	for _, r := range gfx.IconRegions {
		if r.Spec.Width == 0 {
			continue // 段 1、段 3 的尺寸還沒解，不假裝畫得出來
		}
		add("ICONGRF/"+r.Name, r.Spec, icon[r.Offset:r.Offset+r.Length])
	}

	// 段 3 是視窗外框（見 gfx/chrome.go）。整段留著，畫面層要哪一塊自己取。
	lib.Chrome = icon[gfx.IconRegions[3].Offset : gfx.IconRegions[3].Offset+gfx.IconRegions[3].Length]
	lib.BattleUI = icon[gfx.IconRegions[1].Offset : gfx.IconRegions[1].Offset+gfx.IconRegions[1].Length]

	// DOS/V 游標是 KI.EXE 內建的兩層 16×16 mask，不是 MOUSE.MCH。
	// 這是可選呈現資產：缺少或不是目標 DOS/V 版時保留明確警告，
	// 不讓其它素材與規則載入一起失敗。
	if ki, err := read(dir, "KI.EXE"); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	} else if cursor, err := gfx.DecodeDOSVCursor(ki); err != nil {
		lib.Warns = append(lib.Warns, fmt.Sprintf("DOS/V 硬體游標未載入：%v", err))
	} else {
		lib.cursorPixels = cursor
	}

	// 大地圖。缺檔或解不開就記進 Warns 而不是整個失敗 ——
	// 檢視器的其他功能不該被它拖著一起壞。
	if mapRaw, err := read(dir, "MMAP.MAP"); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	} else if lib.World, err = world.ParseMap(mapRaw); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	}
	if mdl, err := read(dir, "MMAP.MDL"); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	} else if lib.Tiles, err = world.ParseTileSet(mdl); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	}
	if mch, err := read(dir, "MMAP.MCH"); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	} else if lib.MCH, err = world.ParseMCH(mch); err != nil {
		lib.Warns = append(lib.Warns, err.Error())
	}
	return lib, nil
}

// RenderWorld 畫出以 (x0, y0) 格為左上角、cols×rows 格的一塊大地圖。
func (l *Library) RenderWorld(x0, y0, cols, rows, bank int) (*image.RGBA, error) {
	if l.World == nil || l.Tiles == nil {
		return nil, fmt.Errorf("大地圖沒有載入成功，看 Warns")
	}
	return l.World.Render(l.Tiles, l.Palette, bank, x0, y0, cols, rows)
}

// RenderChrome 畫出一塊視窗外框圖塊（8×8）。
func (l *Library) RenderChrome(off, bank int) (*image.RGBA, error) {
	if l.Chrome == nil {
		return nil, fmt.Errorf("ICONGRF 段 3 沒有載入")
	}
	return gfx.RenderChrome(l.Chrome, off, l.Palette, bank)
}

// PaletteColor 回傳原版調色盤指定色，供反組譯已直接給出 palette index
// 的向量 overlay 使用；不把 RGB 常數重抄進呈現層。
func (l *Library) PaletteColor(bank, index int) (color.RGBA, error) {
	if l == nil || l.Palette == nil {
		return color.RGBA{}, fmt.Errorf("GAMEPAL.BRG 沒有載入")
	}
	colors, err := l.Palette.Bank(bank)
	if err != nil {
		return color.RGBA{}, err
	}
	if index < 0 || index >= len(colors) {
		return color.RGBA{}, fmt.Errorf("調色盤索引 %d 超出範圍", index)
	}
	return colors[index], nil
}

// DOSVAmountPanel 解出松崗 DOS/V sub_17D0D 使用的數值輸入器完整內框。
// 資源位於 ICONGRF 第 3 段，不是 chrome.Window 的 8×8 可貼圖邊框。
func (l *Library) DOSVAmountPanel(bank int) (*image.RGBA, error) {
	if l == nil || l.Chrome == nil {
		return nil, fmt.Errorf("ICONGRF 段 3 沒有載入")
	}
	return gfx.DOSVAmountPanel.RenderRGBAAt(l.Chrome, gfx.DOSVAmountPanelOffset, l.Palette, bank)
}

// DOSVFactionLegend 解出縮小地圖下方那條 192×16 勢力色標。
func (l *Library) DOSVFactionLegend(bank int) (*image.RGBA, error) {
	if l == nil || l.Chrome == nil {
		return nil, fmt.Errorf("ICONGRF 段 3 沒有載入")
	}
	return gfx.DOSVFactionLegend.RenderRGBAAt(l.Chrome,
		gfx.DOSVFactionLegendOffset, l.Palette, bank)
}

// DOSVResourceIcon 解出資金／預備兵欄左邊那一直排 24×16 圖示。
//
// index 0–3 依序是天秤（資金）、馬、弓、步；green 選綠色那一組
// （財政視窗「次月」欄用它，「今月底」用紅色）。位址換算見
// `docs/re/48` §6。
func (l *Library) DOSVResourceIcon(index int, green bool, bank int) (*image.RGBA, error) {
	if l == nil || l.Chrome == nil {
		return nil, fmt.Errorf("ICONGRF 段 3 沒有載入")
	}
	if index < 0 || index >= gfx.DOSVResourceIconCount {
		return nil, fmt.Errorf("資源圖示 %d 超出 0–%d", index, gfx.DOSVResourceIconCount-1)
	}
	base := gfx.DOSVResourceIconOffset
	if green {
		base = gfx.DOSVResourceIconGreenOffset
	}
	return gfx.DOSVResourceIcon.RenderRGBAAt(l.Chrome,
		base+index*gfx.DOSVResourceIconStride, l.Palette, bank)
}

// DOSVTroopIcon 解出編成畫面六個槽用的兵種圖示。
//
// kind 用原版的兵種編碼：1 騎馬／2 弓兵／3 步兵／4 空槽。
// 前三個就是 `DOSVResourceIcon` 綠色組的第 2–4 張，空槽另有一張
// （`docs/re/49` §3）。
func (l *Library) DOSVTroopIcon(kind, bank int) (*image.RGBA, error) {
	if l == nil || l.Chrome == nil {
		return nil, fmt.Errorf("ICONGRF 段 3 沒有載入")
	}
	if kind < 1 || kind > 4 {
		return nil, fmt.Errorf("兵種 %d 超出 1–4", kind)
	}
	if kind == 4 {
		return gfx.DOSVResourceIcon.RenderRGBAAt(l.Chrome,
			gfx.DOSVEmptySlotIconOffset, l.Palette, bank)
	}
	return l.DOSVResourceIcon(kind, true, bank)
}

// DOSVOrderIcon 解出戰術底列每一格右半那張「目前命令」的圖示。
//
// code 是命令碼 0–5（陣形／攻擊／突擊／城壁／守陣／退卻），
// 與 `tactical.Command` 同一組編碼。原版走 `sub_1C673`，
// 位址換算與內容檢查見 `gfx.DOSVOrderIconOffset`。
func (l *Library) DOSVOrderIcon(code, bank int) (*image.RGBA, error) {
	if l == nil || l.Chrome == nil {
		return nil, fmt.Errorf("ICONGRF 段 3 沒有載入")
	}
	if code < 0 || code >= gfx.DOSVOrderIconCount {
		return nil, fmt.Errorf("命令碼 %d 超出 0–%d", code, gfx.DOSVOrderIconCount-1)
	}
	return gfx.DOSVResourceIcon.RenderRGBAAt(l.Chrome,
		gfx.DOSVOrderIconOffset+code*gfx.DOSVOrderIconStride, l.Palette, bank)
}

// DOSVBattleCommandBase 解出 sub_1C7F4 在底列重複六次的 80×32 底板。
func (l *Library) DOSVBattleCommandBase(bank int) (*image.RGBA, error) {
	if l == nil || l.BattleUI == nil {
		return nil, fmt.Errorf("ICONGRF 段 1 沒有載入")
	}
	return gfx.DOSVBattleCommandBase.RenderRGBAAt(l.BattleUI,
		gfx.DOSVBattleCommandBaseOffset, l.Palette, bank)
}

// DOSVBattleCommandGlyph 解出 sub_1F888 連續消費的六張 24×16 指令 glyph。
func (l *Library) DOSVBattleCommandGlyph(index, bank int) (*image.RGBA, error) {
	if l == nil || l.BattleUI == nil {
		return nil, fmt.Errorf("ICONGRF 段 1 沒有載入")
	}
	if index < 0 || index >= gfx.DOSVBattleCommandGlyphCount {
		return nil, fmt.Errorf("DOS/V 戰術指令 glyph %d 超出範圍", index)
	}
	off := gfx.DOSVBattleCommandGlyphOffset + index*gfx.DOSVBattleCommandGlyphStride
	return gfx.DOSVBattleCommandGlyph.RenderRGBAAt(l.BattleUI, off, l.Palette, bank)
}

// DOSVBattleSideCommands 解出 sub_1C863 直接貼到右欄的 128×96 六列面板。
func (l *Library) DOSVBattleSideCommands(bank int) (*image.RGBA, error) {
	if l == nil || l.BattleUI == nil {
		return nil, fmt.Errorf("ICONGRF 段 1 沒有載入")
	}
	return gfx.DOSVBattleSideCommands.RenderRGBAAt(l.BattleUI,
		gfx.DOSVBattleSideCommandsOffset, l.Palette, bank)
}

// DOSVBattleFlag 解出側欄上／下兩格的將旗底圖（各 128×32）。
// foe 為真時取上格（對方，段 1 `0x0800`），否則取下格（我方，`0x1000`）。
func (l *Library) DOSVBattleFlag(foe bool, bank int) (*image.RGBA, error) {
	if l == nil || l.BattleUI == nil {
		return nil, fmt.Errorf("ICONGRF 段 1 沒有載入")
	}
	off := gfx.DOSVBattleFlagAllyOffset
	if foe {
		off = gfx.DOSVBattleFlagFoeOffset
	}
	return gfx.DOSVBattleFlag.RenderRGBAAt(l.BattleUI, off, l.Palette, bank)
}

// DOSVBattleFormationStrip 解出十六個陣形那一格（128×32，8 欄 × 2 列）。
func (l *Library) DOSVBattleFormationStrip(bank int) (*image.RGBA, error) {
	if l == nil || l.BattleUI == nil {
		return nil, fmt.Errorf("ICONGRF 段 1 沒有載入")
	}
	return gfx.DOSVBattleFormation.RenderRGBAAt(l.BattleUI,
		gfx.DOSVBattleFormationOffset, l.Palette, bank)
}

// DOSVBattleSideFooter 解出側欄最底那一條（128×16）。
func (l *Library) DOSVBattleSideFooter(bank int) (*image.RGBA, error) {
	if l == nil || l.BattleUI == nil {
		return nil, fmt.Errorf("ICONGRF 段 1 沒有載入")
	}
	return gfx.DOSVBattleSideFooter.RenderRGBAAt(l.BattleUI,
		gfx.DOSVBattleSideFooterOffset, l.Palette, bank)
}

// DOSVCursor 畫出 DOS/V KI.EXE 內建的 16×16 白框／紅填游標。
func (l *Library) DOSVCursor(bank int) (*image.RGBA, error) {
	if l == nil || len(l.cursorPixels) == 0 {
		return nil, fmt.Errorf("DOS/V KI.EXE 內建游標沒有載入")
	}
	return gfx.RenderDOSVCursorPixels(l.cursorPixels, l.Palette, bank)
}

// Portrait 畫出 KAOGRF 的第 page 張頭像。
// page 要傳武將記錄的 `+0x01`（state.General.Portrait），**不是武將編號**。
func (l *Library) Portrait(page, bank int) (*image.RGBA, error) {
	for i, e := range l.Entries {
		if e.Spec.Name == gfx.Kao.Name {
			return l.Render(i, page, bank)
		}
	}
	return nil, fmt.Errorf("KAOGRF 沒有載入")
}

// Location 畫出 KYOGRF 的第 page 張據點景觀（96×96）。
//
// page 要傳據點記錄 `+0x16` 的**高 4 位**（state.City.KindHigh），
// 值域 0–14（`docs/re/50` §3）。
func (l *Library) Location(page, bank int) (*image.RGBA, error) {
	for i, e := range l.Entries {
		if e.Spec.Name == gfx.Kyo.Name {
			return l.Render(i, page, bank)
		}
	}
	return nil, fmt.Errorf("KYOGRF 沒有載入")
}

// Banner 畫出最上方那條 640×32 的標題橫幅（ICONGRF 段 0）。
func (l *Library) Banner(bank int) (*image.RGBA, error) {
	for i, e := range l.Entries {
		if e.Label == "ICONGRF/banner" {
			return l.Render(i, 0, bank)
		}
	}
	return nil, fmt.Errorf("ICONGRF 段 0 沒有載入")
}

// Render 畫出第 asset 種素材的第 page 張。
func (l *Library) Render(asset, page, bank int) (*image.RGBA, error) {
	if asset < 0 || asset >= len(l.Entries) {
		return nil, fmt.Errorf("素材編號 %d 超出範圍（共 %d 種）",
			asset, len(l.Entries))
	}
	e := l.Entries[asset]
	if e.Count > 0 {
		page = ((page % e.Count) + e.Count) % e.Count
	}
	return e.Spec.RenderRGBA(e.Data, page, l.Palette, bank)
}
