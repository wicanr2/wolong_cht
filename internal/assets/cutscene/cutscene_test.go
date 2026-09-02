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

// 十九個過場檔全部要解到**檔頭宣告的長度**，一個 byte 都不差。
//
// ⭐ 這是「格式對不對」最強也最便宜的檢查（docs/spec/113）：
// 相位一掉，長度就對不上，而且畫面會整體位移。
func TestAllCutscenesDecodeToDeclaredLength(t *testing.T) {
	want := map[string]int{"END_S1.DAT": 91200, "OPEN_S1.DAT": 78016,
		"OPEN_S2.DAT": 384000, "OPEN_S3.DAT": 384000, "OPEN_S4.DAT": 384000,
		"OPEN_S5.DAT": 288000}
	names := []string{"GAMEOVER.DAT"}
	for i := 1; i <= 12; i++ {
		names = append(names, fmt.Sprintf("END_S%d.DAT", i))
	}
	for i := 1; i <= 6; i++ {
		names = append(names, fmt.Sprintf("OPEN_S%d.DAT", i))
	}
	for _, n := range names {
		out, err := rle.DecodeFile(read(t, n))
		if err != nil {
			t.Errorf("%s：%v", n, err)
			continue
		}
		exp := Size
		if v, ok := want[n]; ok {
			exp = v
		}
		if len(out) != exp {
			t.Errorf("%s 解出 %d B，預期 %d", n, len(out), exp)
		}
	}
}

// `OPEN_S2`–`S4` 是 12 幀 320 × 200 的動畫（docs/re/76 §7）：
// 第 0 幀與第 4 幀幾乎相同（四個相位 × 三個階段），
// 而第 0 幀與第 1 幀不同——相位掉了的話這兩件事會一起壞。
func TestOpenS2IsTwelveFrames(t *testing.T) {
	buf, err := rle.DecodeFile(read(t, "OPEN_S2.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	const frame = 32000
	if len(buf) != 12*frame {
		t.Fatalf("%d B，不是 12 × %d", len(buf), frame)
	}
	diff := func(a, b int) int {
		n := 0
		for i := 0; i < frame; i++ {
			if buf[a*frame+i] != buf[b*frame+i] {
				n++
			}
		}
		return n
	}
	if d := diff(0, 4); d > frame/100 {
		t.Errorf("第 0 幀與第 4 幀差 %d B，同一個相位不該差這麼多", d)
	}
	if d := diff(0, 1); d < frame/100 {
		t.Errorf("第 0 幀與第 1 幀只差 %d B，相鄰相位不該幾乎相同", d)
	}
}

// END_S2 是結局第二幕「三國志街道」。解出來的色號要真的用到多階，
// 而不是一整片同色——那是版面猜錯時最常見的樣子。
func TestEndS2LooksLikeAnImage(t *testing.T) {
	buf, err := Decode(read(t, "END_S2.DAT"))
	if err != nil {
		t.Fatal(err)
	}
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
