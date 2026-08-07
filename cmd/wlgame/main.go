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
// ⚠ 中文還畫不出來。原版用倚天 16×15 點陣字，那一項還沒解
// （CLAUDE.md §3.6），而 Ebiten 內建的 debug 字型只有 ASCII，
// 中文會被**靜靜吃掉**。所以畫面上的字一律是 ASCII，
// 並在需要顯示人名／地名的地方標出編號。
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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/state"
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

var windowNames = [numWindows]string{"COMMAND", "FACTION", "MINIMAP", "SYSTEM"}

// residentWindows 是「開著也不會停時間」的那三個。
var residentWindows = [numWindows]bool{winCommand: true, winFaction: true, winMinimap: true}

type game struct {
	lib   *library.Library
	world *state.World
	rng   *lcg

	open       [numWindows]bool
	camX, camY int

	// speed 是每個畫面更新要推進幾個遊戲 tick。
	// 原版的速度設定是「每 tick 之後等 N 個計時中斷」，沒有固定 tick rate
	// （docs/re/06 §4）；remake 改成固定 60 Hz 邏輯更新 ＋ 可調倍率，
	// 並把這件事標記為 remake 差異（15-realtime.md §7）。
	speed int

	lastEvent string
	quitting  bool

	// 截圖模式：跑 shotAt 幀之後把畫面存成 PNG 然後結束。
	// 這是這支程式**唯一**的自動驗收方式——Ebiten 要顯示器，
	// 所以 CI／容器裡要靠 Xvfb ＋ 這個旗標才驗得到畫面。
	shotPath string
	shotAt   int
	frame    int
}

type lcg struct{ s uint32 }

func (r *lcg) Next() int {
	r.s = r.s*1664525 + 1013904223
	return int(r.s >> 16)
}

// timeRuns 是暫停規則。
//
//	時間推進 ⟺ 開啟中的視窗集合 ⊆ {命令, 自勢力情報, 縮小地圖}
//
// 刻意寫成一個函式而不是散在各視窗的開關程式碼裡——
// 這樣「哪些視窗會停時間」只有一個地方可以改。
func (g *game) timeRuns() bool {
	for k := windowKind(0); k < numWindows; k++ {
		if g.open[k] && !residentWindows[k] {
			return false
		}
	}
	return true
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
			g.lastEvent = fmt.Sprintf("month end  fires/riots=%d  storm=%v",
				len(ev.Disaster), ev.Storm != nil)
		}
		for _, f := range ev.Eliminated {
			g.lastEvent = fmt.Sprintf("faction %d eliminated", f)
		}
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

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf(
		"WOLONG   %3d/%02d/%02d  %02dh  %s",
		c.Year, c.Month, c.Day, c.Hour, seasonASCII(c.Season())), 8, 8)

	stat := fmt.Sprintf(""+
		"FACTION %d\n"+
		"lord    #%d\n"+
		"advisor #%d\n"+
		"trust   %d\n"+
		"cities  %d\n"+
		"funds   %d\n"+
		"reserve\n"+
		" cav %d\n"+
		" arc %d\n"+
		" inf %d\n"+
		"tax     %d%%\n"+
		"\n"+
		"TIME  %s\n"+
		"speed %d",
		p, f.Lord, f.Advisor, f.Trust, f.Cities, f.Funds,
		f.Reserves[economy.Cavalry], f.Reserves[economy.Archer],
		f.Reserves[economy.Infantry], g.world.TaxRate,
		map[bool]string{true: "RUNNING", false: "PAUSED"}[g.timeRuns()], g.speed)
	ebitenutil.DebugPrintAt(screen, stat, mapW+6, bannerH+4)

	const winListY = bannerH + 4 + 16*15
	ebitenutil.DebugPrintAt(screen, "WINDOWS", mapW+6, winListY)
	for k := windowKind(0); k < numWindows; k++ {
		mark := " "
		if g.open[k] {
			mark = "*"
		}
		note := ""
		if !residentWindows[k] {
			note = " (stops time)"
		}
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("%d%s %s%s", k+1, mark, windowNames[k], note),
			mapW+6, winListY+16+int(k)*14)
	}

	// 開著的視窗畫成疊在地圖上的框。內容還沒做——
	// 現階段的重點是「開了會不會停時間」這條規則，不是視窗長什麼樣。
	y := bannerH + 16
	for k := windowKind(0); k < numWindows; k++ {
		if !g.open[k] {
			continue
		}
		vector.DrawFilledRect(screen, 16, float32(y), 260, 72, color.RGBA{0, 0, 0, 200}, false)
		vector.StrokeRect(screen, 16, float32(y), 260, 72, 1, color.RGBA{200, 180, 120, 255}, false)
		ebitenutil.DebugPrintAt(screen, windowNames[k]+" window (not implemented)", 24, y+8)
		if !residentWindows[k] {
			ebitenutil.DebugPrintAt(screen, "time is STOPPED while this is open", 24, y+24)
		}
		y += 80
	}

	// 底部狀態列自己鋪底，否則字會壓在地圖上看不清楚。
	vector.DrawFilledRect(screen, 0, screenH-16, mapW, 16, color.RGBA{0, 0, 0, 200}, false)
	ebitenutil.DebugPrintAt(screen,
		"[1-4]win [arrows]scroll [-/=]speed [ESC]close [F10]quit", 4, screenH-15)
	if g.lastEvent != "" {
		ebitenutil.DebugPrintAt(screen, g.lastEvent, mapW+6, screenH-15)
	}

	if g.quitting {
		ebitenutil.DebugPrintAt(screen, "Quit? (Y/N)", screenW/2-40, screenH/2)
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

func seasonASCII(s clock.Season) string {
	return [...]string{"spring", "summer", "autumn", "winter"}[s]
}

func main() {
	dir := flag.String("orig", "workplace/orig/dosv", "原版素材目錄（請自備）")
	scenPath := flag.String("scenario-file", "", "劇本檔路徑（預設 <orig>/SINARIO.DAT）")
	scenario := flag.Int("scenario", 0, "劇本編號 0–3")
	player := flag.Int("player", 0, "玩家所仕的勢力編號")
	speed := flag.Int("speed", 4, "每個畫面更新推進幾個遊戲 tick")
	shot := flag.String("shot", "", "跑 N 幀之後截圖到這個路徑就結束（驗收用）")
	shotFrames := flag.Int("shot-frames", 120, "截圖前先跑幾幀")
	openWin := flag.Int("open-window", -1, "截圖前先打開第幾個視窗（0–3，驗收暫停規則用）")
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

	g := &game{lib: lib, world: w, rng: &lcg{s: 1}, speed: *speed,
		shotPath: *shot, shotAt: *shotFrames}
	if *openWin >= 0 && *openWin < int(numWindows) {
		g.open[*openWin] = true
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
