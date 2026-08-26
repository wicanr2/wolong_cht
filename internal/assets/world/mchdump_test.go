package world

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// 把 `MMAP.MCH` 的疊圖倒成一張總覽，肉眼看得到「軍團圖塊長什麼樣」。
//
// ⚠ 平常不跑：要看的時候給 `WOLONG_DUMP_DIR=<容器內目錄>`。
// 這不是驗收，是**眼睛的 oracle**——`docs/spec/74` 只寫了圖號怎麼算
//（勢力 × 5 ＋ 朝向），沒有一張圖說明那五張是什麼。
func TestDumpCorpsOverlayTiles(t *testing.T) {
	dir := os.Getenv("WOLONG_DUMP_DIR")
	if dir == "" {
		t.Skip("沒有 WOLONG_DUMP_DIR，跳過")
	}
	raw, err := os.ReadFile("../../../workplace/orig/dosv/MMAP.MCH")
	if err != nil {
		t.Skipf("找不到 MMAP.MCH：%v", err)
	}
	mch, err := ParseMCH(raw)
	if err != nil {
		t.Fatal(err)
	}
	const cols, cell = 20, 16
	rows := (256 + cols - 1) / cols
	img := image.NewRGBA(image.Rect(0, 0, cols*(cell+2), rows*(cell+2)))
	for n := 0; n < 256; n++ {
		tile := mch.Tile(byte(n))
		if tile == nil {
			continue
		}
		ox, oy := (n%cols)*(cell+2), (n/cols)*(cell+2)
		for y := 0; y < cell; y++ {
			for x := 0; x < cell; x++ {
				v := tile.Pix[y*cell+x]
				if v == MCHTransparent {
					img.Set(ox+x, oy+y, color.RGBA{24, 24, 32, 255})
					continue
				}
				g := uint8(16 * (int(v) % 16))
				img.Set(ox+x, oy+y, color.RGBA{g, uint8(255 - int(g)), 210, 255})
			}
		}
	}
	f, err := os.Create(dir + "/mch.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	t.Logf("寫出 %s/mch.png", dir)
}
