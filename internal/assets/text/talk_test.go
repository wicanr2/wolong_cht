package text

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// origPath 回傳原版素材的路徑。素材不進版控，沒有就跳過測試
// ——CI 上不該因為缺原版檔而紅。
func origPath(t *testing.T, ver, name string) string {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", ver, name)
	if _, err := os.Stat(p); err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	return p
}

// TestRoundTrip 是本套件唯一有意義的驗收：解出來再組回去必須
// byte-for-byte 相同。寫不回去的中文化工具是不能用的。
func TestRoundTrip(t *testing.T) {
	for _, ver := range []string{"dosv", "pc98"} {
		t.Run(ver, func(t *testing.T) {
			raw, err := os.ReadFile(origPath(t, ver, "TALK.DAT"))
			if err != nil {
				t.Fatal(err)
			}
			tbl, err := Parse(raw)
			if err != nil {
				t.Fatal(err)
			}
			got := tbl.Bytes()
			if !bytes.Equal(got, raw) {
				for i := range got {
					if i >= len(raw) || got[i] != raw[i] {
						t.Fatalf("第一個差異 @ 0x%04x（原 %d B，重建 %d B）",
							i, len(raw), len(got))
					}
				}
				t.Fatalf("長度不同：原 %d B，重建 %d B", len(raw), len(got))
			}
		})
	}
}

// TestMarkerSet 釘住 docs/formats/01 §3 的統計：六種標記，沒有 \5。
func TestMarkerSet(t *testing.T) {
	raw, err := os.ReadFile(origPath(t, "dosv", "TALK.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	tbl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	count := map[byte]int{}
	for _, m := range tbl.Messages {
		for _, mk := range m.Markers() {
			count[mk]++
		}
	}
	want := map[byte]int{'1': 102, '2': 37, '3': 85, '4': 34, '6': 72, '7': 18}
	if len(count) != len(want) {
		t.Fatalf("標記種類 %d 種，預期 %d 種：%v", len(count), len(want), count)
	}
	for k, v := range want {
		if count[k] != v {
			t.Errorf("標記 \\%c 出現 %d 次，預期 %d", k, count[k], v)
		}
	}
	if _, ok := count['5']; ok {
		t.Error("出現了 \\5 —— 原版沒有這個標記")
	}
}

// TestTrailingBlankLinesArePreserved 釘住最容易踩的那個坑：
// 訊息結尾的空行是資料，正規化掉就寫不回去。
func TestTrailingBlankLinesArePreserved(t *testing.T) {
	raw, err := os.ReadFile(origPath(t, "dosv", "TALK.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	tbl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	var single, blank int
	for _, m := range tbl.Messages {
		if len(m.Lines) == 1 && len(m.Lines[0].Parts) == 0 {
			blank++ // 空訊息：只有一個 NUL
		}
		if n := len(m.Bytes()); n == 1 {
			single++
		}
	}
	if blank != 78 || single != 78 {
		t.Errorf("空訊息 %d 則（單 NUL %d 則），預期各 78 則", blank, single)
	}
}
