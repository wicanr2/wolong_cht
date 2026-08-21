package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// 逐幀錄製（docs/spec/71）。推廣片的動態素材靠它產生。
//
// ⭐ **不錄螢幕，讓程式自己把每一幀寫出來**。X11 擷取要跟畫面搶時間，
// 幀率不穩還會抖；逐幀輸出是確定性的——同一個 seed 跑兩次得到同一批圖。

// recorder 把每一張畫出來的圖寫成 PNG。
type recorder struct {
	dir   string
	total int
	n     int
	buf   *image.RGBA
}

func newRecorder(dir string, total int) *recorder {
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("⚠ 建不出錄製目錄 %s：%v", dir, err)
	}
	if total <= 0 {
		total = 300
	}
	return &recorder{dir: dir, total: total}
}

// shot 寫出這一張。回傳 true 表示錄滿了。
//
// ⚠ 取像素用 `ReadPixels`，不要學 `saveShot` 逐點 `At`——那支一次只截一張
// 所以無所謂，錄三百張會變成瓶頸。
func (r *recorder) shot(screen *ebiten.Image) bool {
	if r.n >= r.total {
		return true
	}
	b := screen.Bounds()
	if r.buf == nil {
		r.buf = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	}
	screen.ReadPixels(r.buf.Pix)
	// 檔名的號碼是**畫出來的第幾張**，不是 g.frame——Ebiten 在慢速的
	// 軟體算圖下會跳過 Draw，拿 g.frame 當檔名會留下洞，
	// ffmpeg 的 -start_number 就接不下去。
	path := filepath.Join(r.dir, fmt.Sprintf("f%05d.png", r.n))
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("⚠ 寫不出 %s：%v", path, err)
	}
	defer f.Close()
	if err := png.Encode(f, r.buf); err != nil {
		log.Fatalf("⚠ PNG 編碼失敗：%v", err)
	}
	r.n++
	return r.n >= r.total
}
