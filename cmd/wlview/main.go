// wlview 是原版素材的互動檢視器。
//
// 這不是遊戲本體 —— 遊戲的規則層要等 docs/mechanics/15-realtime.md
// （即時制的時間模型）標成 READY 才能動手，見 CLAUDE.md §9。
// 這支程式只做一件事：把已經標 READY 的三份規格（調色盤、圖庫、訊息表）
// 接成一條「解碼 → 呈現」的管道，讓格式的正確性能用眼睛驗收，
// 而不是只有測試綠。
//
//	tools/go.sh run ./cmd/wlview -orig workplace/orig/dosv
//
// 要在無頭環境輸出 PNG 請用 cmd/wlshot —— 那支不 import Ebiten。
// **Ebiten 在 init 期就要求顯示器**，跟它放同一個 binary 就跑不起來
// （這個坑實際踩過，見 internal/assets/library 的套件說明）。
//
// 操作：← → 換頁、↑ ↓ 換素材、1–4 換季節、ESC 退回第一張、F10 離開。
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

// 原版是 640×400（VGA 16 色 planar，見 docs/re/02 §5）。
// 檢視器沿用同一個畫布，素材才不會被縮放 —— 縮放過的畫面沒有驗收價值。
const (
	screenW = 640
	screenH = 400
)

var seasonNames = []string{"spring", "summer", "autumn", "winter"}

type app struct {
	lib      *library.Library
	cur      int  // 目前素材
	page     int  // 目前張數
	season   int  // 調色盤組（0–3 ＝ 春夏秋冬）
	quitting bool // F10 之後的 Y/N 確認
}

func (a *app) Update() error {
	// [HARD] ESC 只取消／退回上一層，F10 才離開（CLAUDE.md §10）。
	if a.quitting {
		switch {
		case pressed(ebiten.KeyY):
			return ebiten.Termination
		case pressed(ebiten.KeyN), pressed(ebiten.KeyEscape):
			a.quitting = false
		}
		return nil
	}
	switch {
	case pressed(ebiten.KeyF10):
		a.quitting = true
	case pressed(ebiten.KeyEscape):
		a.page = 0
	case pressed(ebiten.KeyArrowRight):
		a.page++
	case pressed(ebiten.KeyArrowLeft):
		a.page--
	case pressed(ebiten.KeyArrowDown):
		a.cur, a.page = (a.cur+1)%len(a.lib.Entries), 0
	case pressed(ebiten.KeyArrowUp):
		a.cur, a.page = (a.cur+len(a.lib.Entries)-1)%len(a.lib.Entries), 0
	}
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if pressed(k) {
			a.season = i
		}
	}
	if n := a.lib.Entries[a.cur].Count; n > 0 {
		a.page = ((a.page % n) + n) % n
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	e := a.lib.Entries[a.cur]
	img, err := a.lib.Render(a.cur, a.page, a.season)
	if err != nil {
		ebitenutil.DebugPrint(screen, err.Error())
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64((screenW-e.Spec.Width)/2),
		float64((screenH-e.Spec.Height)/2)+12,
	)
	screen.DrawImage(ebiten.NewImageFromImage(img), op)

	// ⚠ 狀態列刻意全用 ASCII。`ebitenutil.DebugPrint` 用的是內建的
	// ASCII 點陣字，中文字會被**靜靜吃掉**——畫面上看起來像少了東西，
	// 很容易被誤判成排版 bug。中文顯示要等倚天 16×15 點陣字那一項
	// （原版的字型來源見 CLAUDE.md §3.6，還沒結案）。
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"%s  %d/%d  %dx%d  season=%s  [<-/->]page [up/dn]asset [1-4]season [F10]quit",
		e.Label, a.page+1, e.Count, e.Spec.Width, e.Spec.Height,
		seasonNames[a.season]))
	if a.quitting {
		ebitenutil.DebugPrintAt(screen, "Quit? (Y/N)", screenW/2-40, screenH/2)
	}
}

func (a *app) Layout(int, int) (int, int) { return screenW, screenH }

func main() {
	dir := flag.String("orig", "workplace/orig/dosv",
		"原版素材目錄（不隨本專案散布，請自備）")
	page := flag.Int("page", 0, "起始張數")
	asset := flag.Int("asset", 0, "起始素材編號")
	season := flag.Int("season", 0, "季節 0-3（春夏秋冬）")
	flag.Parse()

	lib, err := library.Load(*dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range lib.Warns {
		log.Printf("⚠ %s", w)
	}
	if *asset < 0 || *asset >= len(lib.Entries) {
		log.Fatalf("素材編號 %d 超出範圍（共 %d 種）", *asset, len(lib.Entries))
	}

	ebiten.SetWindowSize(screenW*2, screenH*2)
	ebiten.SetWindowTitle("臥龍傳 素材檢視器")
	a := &app{lib: lib, cur: *asset, page: *page, season: *season}
	if err := ebiten.RunGame(a); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
