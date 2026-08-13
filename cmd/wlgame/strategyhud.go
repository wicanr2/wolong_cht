package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// DOS/V 自然策略畫面的固定骨架。影片 oracle 與說明書主畫面圖都對應：
// 橫幅下先有一列命令，再把右側 208 px 留給縮小地圖／自勢力情報；大地圖
// 因而是左側 27×21 格。這不是把模態視窗硬貼到地圖上，而是主畫面的
// resident HUD；命令／列表等暫存視窗仍由 window.go 疊在它上面。
const (
	strategyCommandY          = bannerH
	strategyCommandH          = 32
	strategyMapY              = strategyCommandY + strategyCommandH
	strategySidebarW          = 208 // 192 px minimap + 8 px 內縮與兩側 8 px 外框
	strategySidebarX          = screenW - strategySidebarW
	strategyMapW              = strategySidebarX
	strategyMapH              = screenH - strategyMapY
	strategySidebarInnerX     = strategySidebarX + chrome.Tile
	strategySidebarInnerRight = strategySidebarX + strategySidebarW - chrome.Tile
	strategyMinimapX          = strategySidebarInnerX
	strategyMinimapY          = bannerH + chrome.Tile
	strategyMinimapW          = strategySidebarW - 2*chrome.Tile
	// 上方右欄的內框內容是 128 px minimap，再留一列 8 px 高的勢力色標。
	// 下方情報框與上方框共用一條 8 px 分隔邊；不能把兩個 Window
	// 背靠背畫成 16 px 的雙邊，否則君主頭像會比原版低一格。
	strategyMinimapH       = 152
	strategyMinimapLegendY = bannerH + strategyMinimapH - 2*chrome.Tile
	strategyMinimapSwatchY = strategyMinimapLegendY + 4
	strategyFactionY       = bannerH + strategyMinimapH - chrome.Tile
	strategyFactionH       = screenH - strategyFactionY
	strategyFactionInnerY  = strategyFactionY + chrome.Tile
	strategyFactionInnerW  = strategySidebarW - 2*chrome.Tile

	// 以下是從 DOS/V 右欄畫面量得的文字安全矩形，不是玩法數值證據。
	// 8×15 半形字與 16×15 全形字必須共用同一個寬度契約。
	strategyCommandCellW = 52
	strategyCommandTextW = strategyCommandCellW - 4
	// 命令列內容區是 8 個 52 px cell；左右各留 2 px，讓 cell 之間的
	// 4 px gap 與外框不會誤觸。Y 只命中 16 px 內容區，不命中上下外框。
	strategyCommandHitInset    = 2
	strategyCommandHitY        = strategyCommandY + chrome.Tile
	strategyCommandHitH        = strategyCommandH - 2*chrome.Tile
	strategyInfoYOffset        = 24
	strategyInfoRowStep        = textdraw.GlyphH + textdraw.LineGap
	strategyInfoLabelXOffset   = 80
	strategyInfoLabelW         = 32
	strategyInfoValueXOffset   = 128
	strategyInfoValueW         = strategySidebarInnerRight - (strategySidebarX + strategyInfoValueXOffset)
	strategyInfoDividerXOffset = 120
	strategyInfoDividerH       = 48
	strategyTrustYOffset       = 88
	strategyResourceDividerY   = 120
	strategyResourceBoxY       = 128
	strategyResourceBoxH       = 88
	strategyResourceTextY      = 8
	strategyResourceReserveY   = 32
	strategyResourceRowStep    = 16
	strategyNumberSlots        = 5
	strategyNumberW            = strategyNumberSlots * textdraw.HalfW
	strategyNumberXOffset      = strategySidebarW - chrome.Tile - strategyNumberW
)

var naturalCommandLabels = [...]string{
	"進　言", "人　事", "財　政", "編　成",
	"軍　團", "據　點", "武　將", "勢　力",
}

// naturalCommandID 依畫面標籤命名，不依快捷鍵順序命名。快捷鍵只是另一個
// 輸入來源；滑鼠 cell 與鍵盤都必須先得到同一個 ID，再進入同一個 action。
type naturalCommandID uint8

const (
	naturalCommandAdvise naturalCommandID = iota
	naturalCommandPersonnel
	naturalCommandFinance
	naturalCommandFormation
	naturalCommandCorps
	naturalCommandCity
	naturalCommandGeneral
	naturalCommandFaction
)

var naturalCommandActions = [...]func(*game){
	(*game).openAdvise,
	(*game).openPersonnel,
	(*game).beginFinance,
	(*game).beginForm,
	(*game).openCorpsList,
	(*game).openCityList,
	(*game).openGeneralList,
	(*game).openFactionList,
}

// strategyCommandCellRect 是自然策略頂端命令列的純幾何契約。
// 不讀 game、不讀輸入狀態；外框、cell gap、地圖與右側 HUD 都不在矩形內。
func strategyCommandCellRect(index int) image.Rectangle {
	if index < 0 || index >= len(naturalCommandLabels) {
		return image.Rectangle{}
	}
	x := chrome.Tile + index*strategyCommandCellW + strategyCommandHitInset
	return image.Rect(x, strategyCommandHitY, x+strategyCommandCellW-2*strategyCommandHitInset,
		strategyCommandHitY+strategyCommandHitH)
}

func hitTestNaturalCommand(x, y int) (naturalCommandID, bool) {
	for i := range naturalCommandLabels {
		if image.Pt(x, y).In(strategyCommandCellRect(i)) {
			return naturalCommandID(i), true
		}
	}
	return 0, false
}

// strategyCommandForShortcut 僅負責把既有鍵盤快捷鍵映射到畫面標籤 ID。
// 特別保留畫面順序與快捷鍵順序的分離：M（行軍）不屬於頂端八格。
func strategyCommandForShortcut(key ebiten.Key) (naturalCommandID, bool) {
	switch key {
	case ebiten.KeyP:
		return naturalCommandAdvise, true
	case ebiten.KeyJ:
		return naturalCommandPersonnel, true
	case ebiten.KeyF:
		return naturalCommandFinance, true
	case ebiten.KeyA:
		return naturalCommandFormation, true
	case ebiten.KeyC:
		return naturalCommandCorps, true
	case ebiten.KeyT:
		return naturalCommandCity, true
	case ebiten.KeyG:
		return naturalCommandGeneral, true
	case ebiten.KeyK:
		return naturalCommandFaction, true
	default:
		return 0, false
	}
}

func (g *game) dispatchNaturalCommand(command naturalCommandID) bool {
	if int(command) < 0 || int(command) >= len(naturalCommandActions) {
		return false
	}
	naturalCommandActions[command](g)
	return true
}

// strategyHUDSingleLine 是常駐欄位的 fail-safe。正常 DOS/V 名稱不會進入截斷路徑；
// 若測試 fixture 或損壞資料帶入過長字串，不能讓字模穿過下一欄或下一列。
func strategyHUDSingleLine(s string, maxPixels int) string {
	if maxPixels <= 0 {
		return ""
	}
	if textdraw.StringWidth(s) <= maxPixels {
		return s
	}
	width := 0
	runes := make([]rune, 0, len([]rune(s)))
	for _, ch := range s {
		w := textdraw.RuneWidth(ch)
		if width+w > maxPixels {
			break
		}
		runes = append(runes, ch)
		width += w
	}
	if len(runes) == 0 {
		return "?"
	}
	return string(runes)
}

// strategyHUDNumber 維持原版右欄的五個半形數字槽位。原始資金在規則層是
// 有號 24 位，可能大於畫面可見欄寬；超出時以五個 9／負四位數飽和值顯示。
// 這不是原版數值規則，只是比 debug 用的井號更適合正式玩家畫面的可讀 fallback；
// 真實 state 不會被改寫。
func strategyHUDNumber(value int) string {
	if value > 99999 {
		value = 99999
	}
	if value < -9999 {
		value = -9999
	}
	return fmt.Sprintf("%5d", value)
}

// drawNaturalStrategyHUD 畫 DOS/V 自然策略畫面的常駐 GUI 骨架。
//
// 證據等級：強推論。影片 `af6xqcicXoI` 的 478×360 影像經黑邊還原後，
// 可量得同一個 640×400 內框；說明書主畫面則直接列出命令列與右側縮小地圖。
// 影片是縮放後的參考，不把壓縮後像素當成原始 asset bytes。
func (g *game) drawNaturalStrategyHUD(screen *ebiten.Image) {
	// 命令列使用與原版視窗相同的深藍底／紅金外框；文字沿 8 px 邊界排列。
	g.chrome.Window(screen, 0, strategyCommandY, strategyMapW, strategyCommandH, chrome.Menu)
	for i, label := range naturalCommandLabels {
		x := chrome.Tile + i*strategyCommandCellW
		g.td.Draw(screen, strategyHUDSingleLine(label, strategyCommandTextW), x, strategyCommandY+8, chrome.Paper)
	}

	// 先畫下方情報框。它從 y=176 開始，讓上方 minimap 的底邊／勢力色標
	// 覆蓋共用分隔邊；這正是原版右欄在 y=168–184 的 16 px 色標列，
	// 避免兩個獨立 Window 疊出 16 px 厚的假分隔。
	g.chrome.Window(screen, strategySidebarX, strategyFactionY, strategySidebarW, strategyFactionH, chrome.Menu)
	g.drawNaturalFactionHUD(screen, strategySidebarX, strategyFactionY)

	// 右上縮小地圖。ICONGRF 段 2 本身就是 192×128，不能用大地圖降採樣替代。
	g.chrome.Window(screen, strategySidebarX, bannerH, strategySidebarW, strategyMinimapH, chrome.Menu)
	if img, err := g.lib.Render(g.minimapAsset(), 0, int(g.world.Clock.Season())); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(strategyMinimapX), float64(strategyMinimapY))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	// 原版縮小地圖下方是兩個半欄寬的勢力色標；圖像本身不含 runtime
	// 君主名，所以這一列由 state 填入。先前只有 8×8 色點，會讓右欄骨架
	// 看起來少掉一整列 16 px 的紅／藍帶。
	red := color.RGBA{210, 48, 40, 255}
	blue := color.RGBA{45, 105, 210, 255}
	legendX := strategySidebarX + chrome.Tile
	vector.DrawFilledRect(screen, float32(legendX), float32(strategyMinimapLegendY), 96, 16, red, false)
	vector.DrawFilledRect(screen, float32(legendX+96), float32(strategyMinimapLegendY), 96, 16, blue, false)
	vector.DrawFilledRect(screen, float32(legendX+8), float32(strategyMinimapSwatchY), 8, 8,
		color.RGBA{85, 154, 69, 255}, false)
	vector.DrawFilledRect(screen, float32(legendX+104), float32(strategyMinimapSwatchY), 8, 8,
		color.RGBA{85, 154, 69, 255}, false)
	g.td.Draw(screen, strategyHUDSingleLine(big5(g.world.LordName(g.world.Player)), 72), legendX+24, strategyMinimapLegendY, chrome.Paper)
	enemyName := "敵"
	for faction, f := range g.world.Factions {
		if faction != g.world.Player && f.Alive {
			enemyName = big5(g.world.LordName(faction))
			break
		}
	}
	g.td.Draw(screen, strategyHUDSingleLine(enemyName, 72), legendX+120, strategyMinimapLegendY, chrome.Paper)
}

func (g *game) drawNaturalFactionHUD(dst *ebiten.Image, x, y int) {
	p := g.world.Player
	f := g.world.Factions[p]
	season := int(g.world.Clock.Season())

	if lord := f.Lord; lord >= 0 && lord < len(g.world.Generals) {
		if img, err := g.lib.Portrait(g.world.Generals[lord].Portrait, season); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(strategySidebarInnerX), float64(strategyFactionInnerY))
			dst.DrawImage(ebiten.NewImageFromImage(img), op)
		}
	}

	ink := chrome.Paper
	// 原版頭像右側是「君主／首都／軍師」三列，中央有一條垂直分隔線。
	// 這裡只重建已由影片量到的常駐骨架；完整勢力細節仍由按 2 的情報視窗呈現。
	labelX := x + strategyInfoLabelXOffset
	valueX := x + strategyInfoValueXOffset
	infoY := y + strategyInfoYOffset
	capital := "－"
	if f.Capital >= 0 && f.Capital < len(g.world.Cities) {
		capital = big5(g.world.Cities[f.Capital].Name)
	}
	info := [...]struct{ label, value string }{
		{"君主", big5(g.world.LordName(p))},
		{"首都", capital},
		{"軍師", big5(g.advisorName())},
	}
	for i, row := range info {
		py := infoY + i*strategyInfoRowStep
		g.td.Draw(dst, strategyHUDSingleLine(row.label, strategyInfoLabelW), labelX, py, ink)
		g.td.Draw(dst, strategyHUDSingleLine(row.value, strategyInfoValueW), valueX, py, ink)
	}
	vector.DrawFilledRect(dst, float32(x+strategyInfoDividerXOffset), float32(infoY), 2, strategyInfoDividerH,
		color.RGBA{190, 190, 190, 255}, false)

	// 頭像下方先放信賴度，再以原版紅色橫線分隔資金／預備兵區。
	g.td.Draw(dst, strategyHUDSingleLine(fmt.Sprintf("信賴度 %3d", g.world.Trust), strategyFactionInnerW-16), x+16, y+strategyTrustYOffset, ink)
	vector.DrawFilledRect(dst, float32(x+8), float32(y+strategyResourceDividerY), float32(strategySidebarW-16), 2,
		color.RGBA{192, 32, 32, 255}, false)

	// 原版資源區是黑底內框，不是把四列數字直接寫在藍色情報底上。
	// 中央三段紅色圖形的語意尚未取得獨立 raw asset，因此只重建已量到的
	// 欄位骨架；數值仍使用同一份 state，不以向量圖形冒充原版 glyph。
	resourceBoxY := y + strategyResourceBoxY
	vector.DrawFilledRect(dst, float32(strategySidebarInnerX), float32(resourceBoxY),
		float32(strategyFactionInnerW), strategyResourceBoxH, color.Black, false)
	resourceY := resourceBoxY + strategyResourceTextY
	resourceInk := color.RGBA{255, 223, 154, 255}
	g.td.Draw(dst, "資金", x+16, resourceY, resourceInk)
	g.td.Draw(dst, "預備兵", x+16, resourceY+textdraw.GlyphH+textdraw.LineGap, resourceInk)
	g.td.Draw(dst, strategyHUDNumber(f.Funds), x+strategyNumberXOffset, resourceY, ink)

	reserveValues := [...]int{
		f.Reserves[economy.Cavalry],
		f.Reserves[economy.Archer],
		f.Reserves[economy.Infantry],
	}
	for i, value := range reserveValues {
		iconY := resourceY + strategyResourceReserveY + i*strategyResourceRowStep
		vector.DrawFilledRect(dst, float32(x+96), float32(iconY), 24, 16,
			color.RGBA{192, 32, 32, 255}, false)
		vector.DrawFilledRect(dst, float32(x+104), float32(iconY+4), 8, 8,
			color.RGBA{64, 0, 0, 255}, false)
		g.td.Draw(dst, strategyHUDNumber(value), x+strategyNumberXOffset, iconY, ink)
	}
}
