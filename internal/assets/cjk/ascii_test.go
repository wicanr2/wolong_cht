package cjk

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// 原版的 `END_S14.DAT` 與倚天 `ASCFONT.15` byte-for-byte 相同
// （`docs/re/29` §5、`docs/spec/137` §3）——所以半形字可以直接從
// 原版資料取，使用者不必自備字型。
//
// ⚠ 判準是**逐 byte 相同**，不是「兩邊都載得起來」：後者對
// 「兩份不同但大小一樣」的字型也會通過，而那正是全形標點踩過的坑。
func TestBuiltinASCIIMatchesEten(t *testing.T) {
	root := repoRoot()
	orig := filepath.Join(root, "workplace", "orig", "dosv", "END_S14.DAT")
	eten := filepath.Join(root, "workplace", "eten", "ASCFONT.15")
	a, err := os.ReadFile(orig)
	if err != nil {
		t.Skip("沒有原版素材，跳過")
	}
	b, err := os.ReadFile(eten)
	if err != nil {
		t.Skip("沒有倚天字型，跳過")
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("END_S14.DAT（%d B）與 ASCFONT.15（%d B）不同——"+
			"半形字不能直接從原版資料取", len(a), len(b))
	}
	f, err := LoadASCIIBuiltin(orig)
	if err != nil {
		t.Fatalf("LoadASCIIBuiltin：%v", err)
	}
	if _, ok := f.Glyph('A'); !ok {
		t.Error("取不到半形 A")
	}
}
