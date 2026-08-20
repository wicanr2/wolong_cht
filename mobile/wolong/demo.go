//go:build !android

package wolongmobile

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/ui/phone"
)

// 推廣片的逐幀輸出。
//
// ⭐ **不錄螢幕，直接把每一幀寫出來**。X11 擷取要跟畫面搶時間，錄出來的
// 幀率不穩、還會抖；逐幀輸出是確定性的——同一個 seed 跑兩次得到同一批圖。
//
// 操作也不用 xdotool 送點擊：時間軸寫在下面，第幾幀做什麼一目了然，
// 而且**與畫面更新同一條時間線**，不會發生「點下去時畫面還沒到那一格」。

// demoStep 是時間軸上的一個動作。
type demoStep struct {
	frame int
	label string // 非空時寫進 marks.txt，給剪接用
	do    func(*phone.Session)
}

// demoTimeline 是推廣片的腳本。
//
// ⚠ **幀號是「畫出來的第幾張圖」，不是遊戲的 tick 數**。Ebiten 在慢速的
// 軟體算圖下會跳過 Draw（Update 照跑），拿 tick 當時間軸的話輸出的張數
// 對不上、動作也會落在不同的位置。一張圖 ＝ 一個影片幀 ＝ 1/30 秒。
func demoTimeline() []demoStep {
	tap := func(x, y float64) func(*phone.Session) {
		return func(s *phone.Session) { s.Tap(x, y) }
	}
	entry := func(c phone.Command) func(*phone.Session) {
		x, y, w, h := phone.CommandRect(int(c))
		return tap(float64(x+w/2), float64(y+h/2))
	}
	row := func(i int) func(*phone.Session) {
		_, my, _, _ := phone.MapRect()
		return tap(200, float64(my+44+i*30+12))
	}
	tab := func(i int) func(*phone.Session) {
		_, my, mw, _ := phone.MapRect()
		return func(s *phone.Session) {
			n := len(s.Tabs())
			if n == 0 {
				return
			}
			cell := mw / n
			s.Tap(float64(i*cell+cell/2), float64(my+20))
		}
	}
	return []demoStep{
		{frame: 1, label: "map", do: func(s *phone.Session) { s.SetZoom(2) }},
		{frame: 75, do: func(s *phone.Session) { s.Pan(6, 3) }},
		{frame: 105, do: func(s *phone.Session) { s.Pan(-6, -3) }},
		// 點自己的首都：浮出小卡。
		{frame: 135, label: "city", do: func(s *phone.Session) {
			s.Select(s.World().Factions[s.World().Player].Capital)
		}},
		{frame: 210, label: "list", do: entry(phone.CmdList)},
		{frame: 270, do: tab(1)},
		{frame: 330, do: tab(2)},
		{frame: 390, do: entry(phone.CmdList)}, // 收起來
		// 軍團 → 編成：挑一名武將，調兩個位置的兵種，送出。
		{frame: 405, label: "corps", do: entry(phone.CmdCorps)},
		{frame: 420, do: tab(1)},
		{frame: 450, do: row(0)},
		// ⚠ 兵種是「騎馬→弓兵→步兵→空」的循環，預設是步兵。
		// **點一下會變空**，推廣片上空槽看起來像瑕疵，所以點兩下換到騎馬。
		{frame: 480, do: row(2)},
		{frame: 490, do: row(2)},
		{frame: 510, do: row(3)},
		{frame: 520, do: row(3)},
		{frame: 545, do: row(7)},
		{frame: 570, do: entry(phone.CmdCorps)},
		// 進言：停戰提案 → 選對象 → 君主開口。
		{frame: 585, label: "advise", do: entry(phone.CmdAdvise)},
		{frame: 615, do: row(1)},
		{frame: 675, do: row(2)},
		{frame: 750, do: func(s *phone.Session) { s.CloseAdvise() }},
		// 戰場：開一場攻城，選一隊下突擊。
		{frame: 780, label: "battle", do: func(s *phone.Session) {
			if err := s.OpenDemoBattle(82); err != nil {
				fmt.Fprintln(os.Stderr, "demo battle:", err)
			}
		}},
		{frame: 840, do: func(s *phone.Session) {
			x, y, w, h := phone.SquadRect(3)
			s.Tap(float64(x+w/2), float64(y+h/2))
		}},
		{frame: 870, do: func(s *phone.Session) {
			x, y, w, h := phone.BattleCommandRect(0) // 第 0 個是突擊
			s.Tap(float64(x+w/2), float64(y+h/2))
		}},
		{frame: 1200, label: "end"},
	}
}

// demoRecorder 把每一幀寫成 PNG，並記下時間軸的標記。
type demoRecorder struct {
	dir   string
	total int
	steps []demoStep
	next  int
	marks *os.File
	n     int
}

func newDemoRecorder() *demoRecorder {
	dir := os.Getenv("WOLONG_FRAMES_DIR")
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	total, _ := strconv.Atoi(os.Getenv("WOLONG_FRAMES_N"))
	if total <= 0 {
		total = 1200 // 30 fps × 40 秒
	}
	f, err := os.Create(filepath.Join(dir, "marks.txt"))
	if err != nil {
		panic(err)
	}
	return &demoRecorder{dir: dir, total: total, steps: demoTimeline(), marks: f}
}

// step 跑到這一張圖該做的動作。
func (r *demoRecorder) step(s *phone.Session) {
	for r.next < len(r.steps) && r.steps[r.next].frame <= r.n {
		st := r.steps[r.next]
		if st.do != nil && s != nil {
			st.do(s)
		}
		if st.label != "" {
			// 標記寫的是**輸出的第幾張圖**，不是遊戲的幀號——
			// 剪接看的是圖的序號。
			fmt.Fprintf(r.marks, "%s %d\n", st.label, r.n)
		}
		r.next++
	}
}

// shot 寫出這一張。回傳 true 表示錄完了。
func (r *demoRecorder) shot(screen *ebiten.Image) bool {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	screen.ReadPixels(img.Pix)
	f, err := os.Create(filepath.Join(r.dir, fmt.Sprintf("f%05d.png", r.n)))
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	f.Close()
	r.n++
	return r.n >= r.total
}

func (r *demoRecorder) close() {
	fmt.Fprintf(r.marks, "total %d\n", r.n)
	r.marks.Close()
}
