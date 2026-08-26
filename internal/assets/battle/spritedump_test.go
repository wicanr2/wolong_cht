package battle

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// 把人物圖形倒成一張 PNG，肉眼看得到「一個兵長什麼樣子」。
//
// ⚠ 平常不跑：要看的時候給 `WOLONG_DUMP_DIR=<容器內目錄>`（`tools/go.sh`
// 只把白名單裡的環境變數傳進容器，這一個在名單上）。
// 這不是驗收，是**眼睛的 oracle**——`docs/re/11` 只寫了圖號怎麼算，
// 沒有一張圖說明「完整的兵」該有多高。
func TestDumpSpriteSheet(t *testing.T) {
	dir := os.Getenv("WOLONG_DUMP_DIR")
	if dir == "" {
		t.Skip("沒有 WOLONG_DUMP_DIR，跳過")
	}
	out := dir + "/sprites.png"
	raw, err := os.ReadFile("../../../workplace/orig/dosv/BATTLE.SCH")
	if err != nil {
		t.Skipf("找不到 BATTLE.SCH：%v", err)
	}
	sp, err := ParseSprites(raw)
	if err != nil {
		t.Fatal(err)
	}
	const cols = 18
	rows := (SpritesPerSide + cols - 1) / cols
	img := image.NewRGBA(image.Rect(0, 0, cols*(SpriteW+2), rows*(SpriteH+2)))
	for i := range img.Pix {
		img.Pix[i] = 0
	}
	for n := 0; n < SpritesPerSide; n++ {
		f := sp.Sprite(0, n)
		if f == nil {
			continue
		}
		ox, oy := (n%cols)*(SpriteW+2), (n/cols)*(SpriteH+2)
		for y := 0; y < SpriteH; y++ {
			for x := 0; x < SpriteW; x++ {
				v := f.At(x, y)
				if v < 0 {
					img.Set(ox+x, oy+y, color.RGBA{20, 20, 30, 255}) // 透明畫成暗底
					continue
				}
				g := uint8(16 * (v % 16))
				img.Set(ox+x, oy+y, color.RGBA{g, uint8(255 - int(g)), 200, 255})
			}
		}
	}
	f, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	t.Logf("寫出 %s（%d 張，每張 %dx%d）", out, SpritesPerSide, SpriteW, SpriteH)
}
