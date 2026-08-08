package main

// 財政、據點、勢力——命令視窗剩下的三格。
//
// **財政的設定值延遲到「次月末」才生效**（說明書：ここにセットされた値は
// 来月末より使用されます，`docs/mechanics/15-realtime.md` §4）。
// 所以這裡改的是 `NextTaxRate`／`NextRecruitCap` 而不是現行值，
// 而畫面要**同時顯示兩欄**——不然玩家看不出「我改的東西還沒生效」。
// 原版的財政視窗也是「今月末」與「來月」兩欄對照。

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// financeState 是財政畫面的狀態。四列：稅率、騎馬、弓兵、步兵
// （原版視窗裡四個綠色按鈕，由上而下就是這個順序）。
type financeState struct {
	active bool
	row    int
}

// TaxMax 是稅率的上限。
//
// ⚠ 這個值**還沒從原版讀出來**，先用 100 當上界。
// 平衡點在 33.75%（`internal/state` 的 TestTaxTippingPoint 有推導），
// 所以上限只要遠高於它就不影響玩法判斷。標成 remake 的暫定值。
const TaxMax = 100

// recruitStep 是募兵數的調整刻度。原版用小算盤直接輸入數字，
// remake 先用固定刻度——**這是操作方式的差異，不是規則的差異**。
const recruitStep = 100

func (g *game) beginFinance() { g.finance = financeState{active: true} }

func (g *game) updateFinance() {
	f := &g.finance
	switch {
	case pressed(ebiten.KeyEscape):
		f.active = false
		return
	case pressed(ebiten.KeyArrowUp):
		f.row = (f.row + 3) % 4
	case pressed(ebiten.KeyArrowDown):
		f.row = (f.row + 1) % 4
	}
	delta := 0
	if pressed(ebiten.KeyArrowRight) || pressed(ebiten.KeyEqual) {
		delta = 1
	}
	if pressed(ebiten.KeyArrowLeft) || pressed(ebiten.KeyMinus) {
		delta = -1
	}
	if delta == 0 {
		return
	}
	if f.row == 0 {
		g.world.NextTaxRate = clamp(g.world.NextTaxRate+delta, 0, TaxMax)
		return
	}
	t := f.row - 1
	g.world.NextRecruitCap[t] = clamp(
		g.world.NextRecruitCap[t]+delta*recruitStep, 0, 60000)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (g *game) drawFinance(screen *ebiten.Image) {
	if !g.finance.active {
		return
	}
	const x, y, w, h = 56, 64, 432, 216
	g.chrome.Window(screen, x, y, w, h, chrome.Menu)

	white := chrome.Paper
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{170, 170, 180, 255}
	sel := color.RGBA{140, 230, 140, 255}

	tx, ty := x+chrome.Tile+4, y+chrome.Tile+2
	g.td.Draw(screen, "財　政", tx, ty, amber)
	g.td.Draw(screen, "今月末", tx+180, ty, dim)
	g.td.Draw(screen, "來　月", tx+280, ty, dim)

	f := g.world.Factions[g.world.Player]
	rows := []struct {
		name     string
		now, nxt int
		unit     string
	}{
		{"稅　率", g.world.TaxRate, g.world.NextTaxRate, "%"},
		{"騎馬募兵", g.world.RecruitCap[economy.Cavalry], g.world.NextRecruitCap[economy.Cavalry], ""},
		{"弓兵募兵", g.world.RecruitCap[economy.Archer], g.world.NextRecruitCap[economy.Archer], ""},
		{"步兵募兵", g.world.RecruitCap[economy.Infantry], g.world.NextRecruitCap[economy.Infantry], ""},
	}
	ry := ty + textdraw.GlyphH + 8
	for i, r := range rows {
		col := white
		mark := "　"
		if i == g.finance.row {
			col, mark = sel, "●"
		}
		g.td.Draw(screen, mark+r.name, tx, ry, col)
		g.td.Draw(screen, fmt.Sprintf("%6d%s", r.now, r.unit), tx+180, ry, dim)
		g.td.Draw(screen, fmt.Sprintf("%6d%s", r.nxt, r.unit), tx+280, ry, col)
		ry += textdraw.GlyphH + 4
	}

	ry += 6
	g.td.Draw(screen, fmt.Sprintf("資金 %8d", f.Funds), tx, ry, white)
	ry += textdraw.GlyphH + 2
	// ⚠ 設定值**次月末**才生效，這一行不能省——
	// 少了它，玩家改完看不出為什麼數字沒動。
	g.td.Draw(screen, "設定值於次月末生效", tx, ry, amber)
	ry += textdraw.GlyphH + 2
	g.td.Draw(screen, "↑↓ 選欄　←→ 增減　ESC 關閉", tx, ry, dim)
}

// openCityList 是命令視窗的「據點」：自勢力據點一覽。
func (g *game) openCityList() {
	rows := g.playerCities()
	if len(rows) == 0 {
		g.lastEvent = "沒有據點"
		return
	}
	g.cityList(rows, "↑↓ 移動　Enter 選取／決定　1-5 排序　ESC 取消",
		func(city int) bool {
			// 說明書：「選了游標移過去」。這裡把鏡頭移到該據點。
			c := g.world.Cities[city]
			g.camX, g.camY = c.X-viewCols/2, c.Y-viewRows/2
			g.clampCam()
			g.lastEvent = "移動到 " + big5(c.Name)
			return true
		})
}

// openFactionList 是命令視窗的「勢力」：他勢力一覽。
func (g *game) openFactionList() {
	var rows []int
	for i := range g.world.Factions {
		if g.world.Factions[i].Alive {
			rows = append(rows, i)
		}
	}
	g.factionList(rows, "↑↓ 移動　Enter 選取／決定　1-4 排序　ESC 取消",
		func(f int) bool {
			g.lastEvent = big5(g.world.LordName(f)) + " 軍"
			return true
		})
}
