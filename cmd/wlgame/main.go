// wlgame 是戰略主畫面的原型。
//
// 它把三層接起來跑：
//
//	internal/state    ← 從 SINARIO.DAT 載入真實的世界
//	internal/rules    ← 時鐘、月結、災害
//	internal/assets   ← 大地圖的圖塊與四季調色盤
//
//	tools/go.sh run ./cmd/wlgame -orig workplace/orig/dosv
//
// **這還不是遊戲。** 沒有指令、沒有軍團、沒有戰鬥、不能存檔。
// 它做的是把「時間在跑的世界」呈現出來，讓已定案的規格
// （docs/mechanics/15-realtime.md）能用眼睛驗收：
//
//   - 時間是連續的，不是回合制
//   - **開啟非常駐視窗會讓時間停止**（§2 的暫停規則）
//   - 季節在 3/6/9/12 月的前 16 天漸變，不是瞬間切換
//   - 月結時資金、預備兵、生產力會跳動
//
// 中文用倚天 16×15 點陣字（`-font` 指到自備的字型目錄）。
// **字型檔不隨本專案散布**，與原版資料同一個處理方式。
// 沒帶 `-font` 也跑得起來，只是中文會畫成空心方框——
// 缺字要看得出來，不能靜靜吃掉。
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/cjk"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/general"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 原版是 640×400。畫面分成三塊，比例照日文說明書 3.1 的截圖：
// 最上方 32 px 的橫幅、左邊的大地圖、右邊 160 px 的資訊欄。
const (
	screenW, screenH = 640, 400
	bannerH          = 32
	panelW           = 160
	mapW             = screenW - panelW
	viewCols         = mapW / 16                // 30
	viewRows         = (screenH - bannerH) / 16 // 23
)

// windowKind 是四個視窗開關（說明書 3.1）。
//
// **前三個開著時時間照跑，第四個會讓時間停止。**
// 這不是 UI 細節，是規則：docs/mechanics/15-realtime.md §2。
type windowKind int

const (
	winCommand windowKind = iota // 命令
	winFaction                   // 自勢力情報
	winMinimap                   // 縮小地圖
	winSystem                    // 系統
	numWindows
)

var windowNames = [numWindows]string{"命令", "自勢力情報", "縮小地圖", "系統"}

// residentWindows 是「開著也不會停時間」的那三個。
var residentWindows = [numWindows]bool{winCommand: true, winFaction: true, winMinimap: true}

type game struct {
	lib   *library.Library
	world *state.World
	rng   *rng.Rand
	td    *textdraw.Drawer

	open       [numWindows]bool
	camX, camY int

	// speed 是每個畫面更新要推進幾個遊戲 tick。
	// 原版的速度設定是「每 tick 之後等 N 個計時中斷」，沒有固定 tick rate
	// （docs/re/06 §4）；remake 改成固定 60 Hz 邏輯更新 ＋ 可調倍率，
	// 並把這件事標記為 remake 差異（15-realtime.md §7）。
	speed int

	// list 是開著的一覽表（武將／據點…）。它是**非常駐視窗**，
	// 所以開著的時候時間會停（15-realtime.md §2）。
	battleLib *battle.Library
	// view 是戰場的等角繪圖資源。沒有原版美術時是 nil，畫面退回高度圖。
	view *battleView

	list    *listwin.List
	sortMem listwin.Memory

	// listRow 把一列資料格式化成「名稱 ＋ 右邊那串數字」。
	// 一覽表本身只管狀態機（listwin），顯示什麼由開啟它的人決定。
	listRow func(id int) (name, cols string)
	// listPick 是決定一列之後要做的事。回傳 true 表示關掉一覽表。
	listPick func(id int) bool
	listHint string

	// form 是編成流程的狀態。
	form formState

	// 進言的狀態機：選指令 → 選對象 → 說服。
	advise    adviseStage
	adviseCmd persuasion.Command
	target    int
	sess      *persuasion.Session
	sessCur   int
	adviseLog []string

	lastEvent string
	quitting  bool

	// 截圖模式：跑 shotAt 幀之後把畫面存成 PNG 然後結束。
	// 這是這支程式**唯一**的自動驗收方式——Ebiten 要顯示器，
	// 所以 CI／容器裡要靠 Xvfb ＋ 這個旗標才驗得到畫面。
	shotPath string
	shotAt   int
	frame    int
}

// adviseStage 是進言的三個階段。
type adviseStage int

const (
	adviseNone        adviseStage = iota
	advisePickCommand             // 選敵對／停戰／協力
	advisePickTarget              // 選對象勢力
	advisePersuade                // 君主拒絕了，開始說服
)

// timeRuns 是暫停規則。
//
//	時間推進 ⟺ 開啟中的視窗集合 ⊆ {命令, 自勢力情報, 縮小地圖}
//
// 刻意寫成一個函式而不是散在各視窗的開關程式碼裡——
// 這樣「哪些視窗會停時間」只有一個地方可以改。
func (g *game) timeRuns() bool {
	// 一覽表、進言、編成都是非常駐視窗 —— 開著就停時間。
	if g.list != nil || g.adviseActive() || g.form.active {
		return false
	}
	for k := windowKind(0); k < numWindows; k++ {
		if g.open[k] && !residentWindows[k] {
			return false
		}
	}
	return true
}

// openGeneralList 開武將一覽。欄位與說明書 5.3 的一覽表一致。
func (g *game) openGeneralList() {
	var rows []int
	for i, gen := range g.world.Generals {
		if gen.Alive && gen.Faction == g.world.Player {
			rows = append(rows, i)
		}
	}
	gs := g.world.Generals
	rating := func(i int) int { return gs[i].Rules().Rating() }
	g.list = listwin.New(listwin.Generals, []listwin.Column{
		{Title: "武將名", Less: func(a, b int) bool { return gs[a].Name < gs[b].Name }},
		{Title: "武術", Less: func(a, b int) bool { return gs[a].Martial > gs[b].Martial }},
		{Title: "統率", Less: func(a, b int) bool { return gs[a].Command > gs[b].Command }},
		{Title: "政治", Less: func(a, b int) bool { return gs[a].Politics > gs[b].Politics }},
		{Title: "評價", Less: func(a, b int) bool { return rating(a) > rating(b) }},
	}, rows, 12, &g.sortMem)
	g.listRow = func(i int) (string, string) {
		gen := gs[i]
		rr := general.General{Aptitude: gen.Aptitude, Martial: gen.Martial,
			Command: gen.Command, Politics: gen.Politics}
		return big5(gen.Name), fmt.Sprintf("%4d%8d%8d%8d",
			gen.Martial, gen.Command, gen.Politics, rr.Rating())
	}
	g.listPick = func(i int) bool {
		g.lastEvent = "選擇了 " + big5(gs[i].Name)
		return true
	}
	g.listHint = "↑↓ 移動　Enter 選取／決定　1-5 排序　ESC 取消"
}

// drawList 畫一覽表。兩段式選取的「反白」用底色表示。
func (g *game) drawList(screen *ebiten.Image) {
	l := g.list
	const x, y, w = 40, 44, 400
	h := 26 + (l.Height+1)*(textdraw.GlyphH+2)
	vector.DrawFilledRect(screen, x, y, w, float32(h), color.RGBA{0, 0, 0, 225}, false)
	vector.StrokeRect(screen, x, y, w, float32(h), 1, color.RGBA{240, 200, 120, 255}, false)

	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}

	// 欄位名。數字鍵 1–5 對應排序，對應說明書「點欄位名排序」。
	cx := x + 8
	for i, c := range l.Columns {
		g.td.Draw(screen, fmt.Sprintf("%d%s", i+1, c.Title), cx, y+6, amber)
		cx += 80
	}

	rows, first := l.Visible()
	ry := y + 26
	for i, r := range rows {
		col := white
		if first+i == l.Cursor {
			hl := color.RGBA{70, 60, 30, 255}
			if l.Phase() == listwin.Selected {
				hl = color.RGBA{180, 140, 40, 255} // 反白
				col = color.RGBA{20, 20, 20, 255}
			}
			vector.DrawFilledRect(screen, x+4, float32(ry-1), w-8,
				float32(textdraw.GlyphH+2), hl, false)
		}
		name, cols := g.listRow(r)
		g.td.Draw(screen, name, x+8, ry, col)
		g.td.Draw(screen, cols, x+88, ry, col)
		ry += textdraw.GlyphH + 2
	}

	hint := g.listHint
	if l.Phase() == listwin.Selected {
		hint = "已選取　Enter 決定　ESC 退回"
	}
	g.td.Draw(screen, hint, x+8, y+h-textdraw.GlyphH-4, dim)
}

func pressed(k ebiten.Key) bool { return inpututil.IsKeyJustPressed(k) }

func (g *game) Update() error {
	g.frame++
	if g.shotPath != "" && g.frame > g.shotAt {
		return ebiten.Termination
	}
	// [HARD] ESC 只取消／關視窗，F10 才離開（CLAUDE.md §10）。
	if g.quitting {
		switch {
		case pressed(ebiten.KeyY):
			return ebiten.Termination
		case pressed(ebiten.KeyN), pressed(ebiten.KeyEscape):
			g.quitting = false
		}
		return nil
	}
	if pressed(ebiten.KeyF10) {
		g.quitting = true
		return nil
	}
	// 戰場畫面獨佔一切——戰鬥中不能停時間，也不能開別的視窗。
	if g.battleActive() {
		g.updateBattle()
		return nil
	}
	// 進言流程是模態的，優先吃輸入。
	if g.adviseActive() {
		g.updateAdvise()
		return nil
	}
	if pressed(ebiten.KeyP) {
		g.openAdvise()
		return nil
	}
	// 一覽表開著時吃掉所有輸入 —— 它是模態的（說明書 3.8 的兩段式操作
	// 只有在獨佔輸入時才成立）。
	if g.list != nil {
		switch {
		case pressed(ebiten.KeyArrowUp):
			g.list.Move(-1)
		case pressed(ebiten.KeyArrowDown):
			g.list.Move(1)
		case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
			if id, ok := g.list.Confirm(); ok {
				if g.listPick(id) {
					g.list = nil
				}
			}
		case pressed(ebiten.KeyEscape):
			if g.list.Cancel() {
				g.list = nil
			}
		}
		for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5} {
			if pressed(k) {
				g.list.SortBy(i)
			}
		}
		return nil
	}
	// 編成畫面是模態的，優先吃輸入。
	if g.form.active {
		g.updateForm()
		return nil
	}
	switch {
	case pressed(ebiten.KeyG):
		g.openGeneralList()
		return nil
	case pressed(ebiten.KeyC):
		g.openCorpsList()
		return nil
	case pressed(ebiten.KeyA):
		g.beginForm()
		return nil
	case pressed(ebiten.KeyM):
		g.beginMarch()
		return nil
	}
	if pressed(ebiten.KeyEscape) {
		// 由上而下關掉最上面那個開著的視窗。
		for k := numWindows - 1; k >= 0; k-- {
			if g.open[k] {
				g.open[k] = false
				break
			}
		}
	}
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if pressed(k) {
			g.open[i] = !g.open[i]
		}
	}
	for i, k := range []ebiten.Key{ebiten.KeyMinus, ebiten.KeyEqual} {
		if pressed(k) {
			g.speed += []int{-1, 1}[i]
		}
	}
	if g.speed < 0 {
		g.speed = 0
	}
	if g.speed > 64 {
		g.speed = 64
	}

	step := 1
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		step = 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.camX += step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.camX -= step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.camY += step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.camY -= step
	}
	g.clampCam()

	if !g.timeRuns() {
		return nil
	}
	for i := 0; i < g.speed; i++ {
		ev := g.world.Tick(g.rng)
		if ev.Settled {
			g.lastEvent = "月結"
			if n := len(ev.Disaster); n > 0 {
				g.lastEvent += fmt.Sprintf("　災害%d", n)
			}
			if ev.Storm != nil {
				g.lastEvent += "　暴風雨"
			}
		}
		for _, i := range ev.Eliminated {
			g.lastEvent = big5(g.world.LordName(i)) + " 滅亡"
		}
		g.reportCorps(ev)
	}
	return nil
}

func (g *game) clampCam() {
	if g.camX < 0 {
		g.camX = 0
	}
	if g.camY < 0 {
		g.camY = 0
	}
	if m := 384 - viewCols; g.camX > m {
		g.camX = m
	}
	if m := 256 - viewRows; g.camY > m {
		g.camY = m
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.battleActive() {
		g.drawBattle(screen)
		if g.shotPath != "" && g.frame == g.shotAt {
			g.saveShot(screen)
		}
		return
	}
	season := int(g.world.Clock.Season())

	// 大地圖。四季調色盤直接吃時鐘算出來的季節——
	// 所以畫面會隨遊戲時間換季，不需要另外驅動。
	if img, err := g.lib.RenderWorld(g.camX, g.camY, viewCols, viewRows, season); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, bannerH)
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}

	// 右側資訊欄。
	vector.DrawFilledRect(screen, mapW, 0, panelW, screenH, color.RGBA{16, 16, 32, 255}, false)
	vector.DrawFilledRect(screen, 0, 0, mapW, bannerH, color.RGBA{32, 24, 16, 255}, false)

	c := g.world.Clock
	p := g.world.Player
	f := g.world.Factions[p]

	white := color.RGBA{240, 240, 230, 255}
	amber := color.RGBA{240, 200, 120, 255}
	dim := color.RGBA{150, 150, 160, 255}

	// 橫幅：遊戲名與日期。原版右上角就是「196年 4月 8日」。
	g.td.Draw(screen, "臥龍傳", 10, 8, amber)
	g.td.Draw(screen, fmt.Sprintf("%d年%2d月%2d日", c.Year, c.Month, c.Day),
		mapW-150, 8, white)
	g.td.Draw(screen, seasonName(c.Season()), mapW-32, 8, amber)

	// 右側資訊欄。
	x := mapW + 6
	y := bannerH + 4
	line := func(s string, col color.RGBA) {
		g.td.Draw(screen, s, x, y, col)
		y += textdraw.GlyphH + textdraw.LineGap
	}
	line(big5(g.world.LordName(p))+" 軍", amber)
	line("軍師 "+big5(g.advisorName()), dim)
	line("", white)
	line(fmt.Sprintf("信賴度 %3d", g.world.Trust), white)
	line(fmt.Sprintf("據點   %3d", f.Cities), white)
	line(fmt.Sprintf("武將   %3d", f.Generals), white)
	line(fmt.Sprintf("軍團   %3d", f.Corps), white)
	line(fmt.Sprintf("資金 %6d", f.Funds), white)
	line("預備兵", white)
	line(fmt.Sprintf(" 騎馬 %5d", f.Reserves[economy.Cavalry]), white)
	line(fmt.Sprintf(" 弓兵 %5d", f.Reserves[economy.Archer]), white)
	line(fmt.Sprintf(" 步兵 %5d", f.Reserves[economy.Infantry]), white)
	line(fmt.Sprintf("稅率   %2d%%", g.world.TaxRate), white)
	line("", white)
	if g.timeRuns() {
		line("時間 進行中", color.RGBA{140, 230, 140, 255})
	} else {
		line("時間 停止", color.RGBA{240, 140, 140, 255})
	}
	line(fmt.Sprintf("速度 %d", g.speed), dim)
	line("", white)
	for k := windowKind(0); k < numWindows; k++ {
		// ⚠ 用 Big5 有的字。「▶」不在 Big5，會被畫成缺字方框——
		// 那個方框是字型層的設計（缺字要看得見），不該由介面自己觸發。
		mark := "　"
		col := dim
		if g.open[k] {
			mark = "●"
			col = amber
		}
		note := ""
		if !residentWindows[k] {
			note = "（停時間）"
		}
		line(fmt.Sprintf("%d%s%s%s", k+1, mark, windowNames[k], note), col)
	}

	// 開著的視窗畫成疊在地圖上的框。內容還沒做——
	// 現階段的重點是「開了會不會停時間」這條規則，不是視窗長什麼樣。
	wy := bannerH + 16
	for k := windowKind(0); k < numWindows; k++ {
		if !g.open[k] {
			continue
		}
		vector.DrawFilledRect(screen, 16, float32(wy), 300, 60, color.RGBA{0, 0, 0, 210}, false)
		vector.StrokeRect(screen, 16, float32(wy), 300, 60, 1, amber, false)
		g.td.Draw(screen, windowNames[k]+"視窗（尚未實作）", 24, wy+8, white)
		if !residentWindows[k] {
			g.td.Draw(screen, "此視窗開啟時，時間停止", 24, wy+30,
				color.RGBA{240, 140, 140, 255})
		}
		wy += 68
	}

	// 底部狀態列自己鋪底，否則字會壓在地圖上看不清楚。
	// **事件列與按鍵列分開兩行**——戰報比按鍵提示長得多，
	// 擠在同一行會被畫到資訊欄上或直接跑出畫面。
	if g.lastEvent != "" {
		vector.DrawFilledRect(screen, 0, screenH-38, mapW, 19, color.RGBA{0, 0, 0, 210}, false)
		g.td.Draw(screen, g.lastEvent, 4, screenH-36, amber)
	}
	vector.DrawFilledRect(screen, 0, screenH-19, mapW, 19, color.RGBA{0, 0, 0, 210}, false)
	g.td.Draw(screen, "1-4 視窗　P 進言　A 編成　M 行軍　C 軍團　G 武將　F10 離開",
		4, screenH-17, dim)
	if !g.td.Available() {
		g.td.Draw(screen, "（未載入字型）", mapW+6, screenH-17,
			color.RGBA{240, 140, 140, 255})
	}

	if g.list != nil {
		g.drawList(screen)
	}
	g.drawForm(screen)
	g.drawAdvise(screen)

	if g.quitting {
		vector.DrawFilledRect(screen, float32(screenW/2-90), float32(screenH/2-14),
			180, 28, color.RGBA{0, 0, 0, 230}, false)
		g.td.Draw(screen, "確定離開？（Y／N）", screenW/2-80, screenH/2-8, white)
	}
	if g.shotPath != "" && g.frame == g.shotAt {
		g.saveShot(screen)
	}
}

// saveShot 把目前畫面寫成 PNG。只在截圖模式用。
func (g *game) saveShot(screen *ebiten.Image) {
	b := screen.Bounds()
	img := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, screen.At(x, y))
		}
	}
	f, err := os.Create(g.shotPath)
	if err != nil {
		log.Printf("⚠ 截圖失敗：%v", err)
		return
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Printf("⚠ 截圖編碼失敗：%v", err)
		return
	}
	log.Printf("截圖 → %s（第 %d 幀，%d年%d月%d日）",
		g.shotPath, g.frame, g.world.Clock.Year, g.world.Clock.Month, g.world.Clock.Day)
}

func (g *game) Layout(int, int) (int, int) { return screenW, screenH }

// seasonName 用中文的季節名。clock.Season 的 String() 已經是中文，
// 這裡只是把它包成一個明確的名字，讓畫面程式讀起來清楚。
func seasonName(s clock.Season) string { return s.String() }

// big5 把 internal/state 保留的原始位元組轉成 UTF-8。
// 解析層刻意保留原始位元組才能 round-trip，所以轉換發生在最外層。
func big5(s string) string {
	if s == "" {
		return "－"
	}
	return text.Decode([]byte(s), text.Big5)
}

// advisorName 回傳玩家所仕勢力的軍師名。0xFF 表示沒有軍師——
// 開新遊戲時玩家本人就是要填進那一格的人。
func (g *game) advisorName() string {
	f := g.world.Factions[g.world.Player]
	if f.Advisor < 0 || f.Advisor >= len(g.world.Generals) {
		return ""
	}
	return g.world.Generals[f.Advisor].Name
}

func main() {
	dir := flag.String("orig", "workplace/orig/dosv", "原版素材目錄（請自備）")
	scenPath := flag.String("scenario-file", "", "劇本檔路徑（預設 <orig>/SINARIO.DAT）")
	scenario := flag.Int("scenario", 0, "劇本編號 0–3")
	player := flag.Int("player", 0, "玩家所仕的勢力編號")
	fontDir := flag.String("font", "workplace/eten",
		"倚天點陣字目錄（STDFONT.15／SPCFONT.15／ASCFONT.15，請自備）")
	speed := flag.Int("speed", 4, "每個畫面更新推進幾個遊戲 tick")
	shot := flag.String("shot", "", "跑 N 幀之後截圖到這個路徑就結束（驗收用）")
	shotFrames := flag.Int("shot-frames", 120, "截圖前先跑幾幀")
	openWin := flag.Int("open-window", -1, "截圖前先打開第幾個視窗（0–3，驗收暫停規則用）")
	openList := flag.Bool("open-list", false, "截圖前先開武將一覽（驗收用）")
	openAdvise := flag.Bool("open-advise", false, "截圖前先跑到說服畫面（驗收用）")
	openForm := flag.Bool("open-form", false, "截圖前先編一支軍團並開編成畫面（驗收用）")
	openCorps := flag.Bool("open-corps", false, "截圖前先編兩支軍團並開軍團一覽（驗收用）")
	openBattle := flag.Bool("open-battle", false, "截圖前先開一場野戰的戰術戰鬥（驗收用）")
	openSiege := flag.Bool("open-siege", false, "截圖前先開一場攻城的戰術戰鬥（驗收用）")
	flag.Parse()

	lib, err := library.Load(*dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range lib.Warns {
		log.Printf("⚠ %s", w)
	}
	path := *scenPath
	if path == "" {
		path = *dir + "/SINARIO.DAT"
	}
	w, err := state.LoadScenario(path, *scenario)
	if err != nil {
		log.Fatal(err)
	}
	w.Player = *player

	log.Printf("劇本 %d：%d年%d月%d日，勢力 %d 個，玩家所仕 %d（君主 %s）",
		*scenario+1, w.Clock.Year, w.Clock.Month, w.Clock.Day,
		len(w.AliveFactions()), *player,
		text.Decode([]byte(w.LordName(*player)), text.Big5))

	// 字型載不到不該讓遊戲開不了——只警告，然後把字畫成方框。
	var font *cjk.Font
	var ascii *cjk.ASCIIFont
	if f, err := cjk.LoadDir(*fontDir, cjk.Options{}); err != nil {
		log.Printf("⚠ 載不到倚天全形字型（%v）；中文會顯示成方框", err)
	} else {
		font = f
	}
	if a, err := cjk.LoadASCIIDir(*fontDir); err != nil {
		log.Printf("⚠ 載不到倚天半形字型（%v）", err)
	} else {
		ascii = a
	}

	g := &game{lib: lib, world: w, rng: rng.Now(), speed: *speed,
		td:       textdraw.New(font, ascii),
		shotPath: *shot, shotAt: *shotFrames}
	if *openWin >= 0 && *openWin < int(numWindows) {
		g.open[*openWin] = true
	}
	if *openList {
		g.openGeneralList()
		g.list.Move(2)
		g.list.Confirm() // 展示反白狀態
	}
	g.installTactical(*dir)
	if *openBattle || *openSiege {
		g.demoBattle(*openSiege)
	}
	if *openForm || *openCorps {
		// 驗收用：直接編幾支軍團出來，免得截圖前要按一長串鍵。
		g.demoCorps(*openCorps)
	}
	if *openAdvise {
		g.openAdvise()
		g.adviseCmd = persuasion.Hostility
		g.target = 13 // 劇本 1 的呂布
		g.beginPersuasion()
		// 曹操（14 據點）對呂布（7 據點）→「我國有利」成立。
		g.offerReason(persuasion.WeAreStronger)
	}
	// 開場把鏡頭移到首都附近。
	if cap := w.Factions[*player].Capital; cap >= 0 && cap < len(w.Cities) {
		g.camX = w.Cities[cap].X - viewCols/2
		g.camY = w.Cities[cap].Y - viewRows/2
	}
	g.clampCam()

	ebiten.SetWindowSize(screenW*2, screenH*2)
	ebiten.SetWindowTitle("臥龍傳 戰略畫面原型")
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
