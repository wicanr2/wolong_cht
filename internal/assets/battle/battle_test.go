package battle

import (
	"os"
	"testing"
)

const dir = "../../../workplace/orig/dosv/"

func load(t *testing.T) *Library {
	t.Helper()
	read := func(n string) []byte {
		b, err := os.ReadFile(dir + n)
		if err != nil {
			t.Skip("找不到原版 " + n + "，跳過")
		}
		return b
	}
	l, err := Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), read("BATTLE.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// 索引第二欄為 0 ⟺ 野戰用的戰場。214 張零例外（docs/re/11 §4.5）。
func TestSiegeFlagMatchesGateColumn(t *testing.T) {
	l := load(t)
	fieldMaps := 0
	for n := 0; n < NumFields; n++ {
		if n >= 192 && l.IsSiege(n) {
			t.Errorf("戰場 %d 是野戰用的，第二欄卻是 %d", n, l.GateX(n))
		}
		if !l.IsSiege(n) {
			fieldMaps++
		}
	}
	// 22 張野戰圖（192–213）＋ 6 張第二欄為 0 的攻城圖。
	if fieldMaps != 28 {
		t.Errorf("第二欄為 0 的有 %d 張，應為 28", fieldMaps)
	}
}

// 堆疊高度：攻城圖長出城牆，平原的野戰圖幾乎是平的。
// 這是「一格是一疊 1–7 層圖塊」那條解讀最硬的驗證。
func TestStackHeightsShapeTheMap(t *testing.T) {
	l := load(t)
	tall := func(n int) int {
		c := 0
		for _, row := range l.Stacks(n) {
			for _, h := range row {
				if h >= 4 {
					c++
				}
			}
		}
		return c
	}
	// 戰場 198 是「平原 ＋ 平原」，只有 10 格高處。
	if got := tall(198); got != 10 {
		t.Errorf("戰場 198 的高處有 %d 格，應為 10", got)
	}
	// 戰場 192 是「山 ＋ 山」，最多。
	if got := tall(192); got != 596 {
		t.Errorf("戰場 192 的高處有 %d 格，應為 596", got)
	}
	// 攻城用的戰場 5 有一圈城牆。
	if got := tall(5); got != 320 {
		t.Errorf("戰場 5 的高處有 %d 格，應為 320", got)
	}
	// 每一格的堆疊都在 0–7。
	for _, row := range l.Stacks(5) {
		for x, h := range row {
			if h < 0 || h > MaxStack {
				t.Fatalf("第 %d 格的堆疊是 %d，應在 0–%d", x, h, MaxStack)
			}
		}
	}
}

// 腳本段編號 ＝ 武將 +0x16 × 4 ＋ 戰場類別。
func TestScriptSelection(t *testing.T) {
	l := load(t)
	for _, tc := range []struct{ field, want int }{
		{0, 0}, {191, 0}, // 攻城
		{192, 1}, {208, 1}, {0xD0, 1},
		{0xD1, 2}, {213, 2}, // 另一組野戰
	} {
		if got := Category(tc.field); got != tc.want {
			t.Errorf("Category(%d) ＝ %d，應為 %d", tc.field, got, tc.want)
		}
	}
	// 呂布那一型（+0x16 ＝ 0）在攻城戰用第 0 段。
	if got := l.Script(0, 0); len(got) != ScriptSize {
		t.Errorf("段 0 長 %d，應為 %d", len(got), ScriptSize)
	}
	// 兩段不該一樣。
	a, b := l.Script(0, 0), l.Script(7, 3)
	same := true
	for i := range a {
		if a[i] != b[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("第 0 段與第 31 段完全相同，段編號算錯了")
	}
}
