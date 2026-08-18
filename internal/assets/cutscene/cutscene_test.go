package cutscene

import (
	"fmt"
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/rle"
)

const dir = "../../../workplace/orig/dosv/"

func read(t *testing.T, n string) []byte {
	t.Helper()
	b, err := os.ReadFile(dir + n)
	if err != nil {
		t.Skip("找不到原版 " + n + "，跳過")
	}
	return b
}

// 十八個過場檔全部要能用同一支解壓走完，而且長度落在合理範圍。
//
// ⭐ 這是「格式對不對」最便宜的檢查：RLE 猜錯的話 run 長度會亂跑，
// 解出來的長度不會齊整地落在 128,000 附近。
func TestAllCutscenesDecode(t *testing.T) {
	names := []string{}
	for i := 1; i <= 12; i++ {
		names = append(names, fmt.Sprintf("END_S%d.DAT", i))
	}
	for i := 1; i <= 6; i++ {
		names = append(names, fmt.Sprintf("OPEN_S%d.DAT", i))
	}
	for _, n := range names {
		raw := len(rle.Decode(read(t, n)))
		if raw < Stride*Height {
			t.Errorf("%s 解出 %d B，比一張單平面畫面還小", n, raw)
		}
	}
}

// END_S2 是結局第二幕「三國志街道」。解出來的色號要真的用到多階，
// 而不是一整片同色——那是版面猜錯時最常見的樣子。
func TestEndS2LooksLikeAnImage(t *testing.T) {
	buf := Decode(read(t, "END_S2.DAT"))
	if len(buf) != Size {
		t.Fatalf("buffer = %d B，預期 %d", len(buf), Size)
	}
	px := Pixels(buf)
	var hist [16]int
	for _, v := range px {
		hist[v]++
	}
	used := 0
	for _, n := range hist {
		if n > Width*Height/1000 {
			used++
		}
	}
	if used < 8 {
		t.Fatalf("只用到 %d 個色號，版面多半猜錯了：%v", used, hist)
	}
	// 標題「三國志街道」在畫面中段；那一帶不能整列同色。
	row := 200
	first := px[row*Width]
	same := true
	for x := 0; x < Width; x++ {
		if px[row*Width+x] != first {
			same = false
			break
		}
	}
	if same {
		t.Fatal("第 200 列整列同色——兩半的切點多半錯了")
	}
}
