package rle

import (
	"encoding/binary"
	"os"
	"testing"
)

// DecodeFile 要跳過 4 byte 長度頭，並且解出來剛好等於宣告值。
func TestDecodeFileSkipsHeaderAndChecksLength(t *testing.T) {
	// "AA" ＋ 重複 3 次 ＝ AAAAA，再接一個字面值 B。
	body := []byte{'A', 'A', 3, 'B'}
	src := make([]byte, 4, 4+len(body))
	binary.LittleEndian.PutUint32(src, 6)
	src = append(src, body...)

	out, err := DecodeFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "AAAAAB" {
		t.Fatalf("解出 %q，預期 %q", out, "AAAAAB")
	}
}

// 長度對不上就是解錯，不能默默放行——這是 docs/spec/113 的驗收條件。
func TestDecodeFileRejectsWrongLength(t *testing.T) {
	src := []byte{9, 0, 0, 0, 'A', 'A', 3, 'B'}
	if _, err := DecodeFile(src); err == nil {
		t.Fatal("宣告 9 B 但只解出 6 B，應該回錯誤")
	}
	if _, err := DecodeFile([]byte{1, 2}); err == nil {
		t.Fatal("連長度頭都不夠，應該回錯誤")
	}
}

// ⭐ `MMAP.MAP` 是唯一「從 0 解也對」的檔：頭 `00 80 01 00` 沒有相鄰重複，
// RLE 原樣吐出而且相位不變。這一條釘住那個例外，免得日後有人以為
// 「跳頭」和「切掉解出來的前四個 byte」是同一件事。
func TestMMapHeaderPassesThroughUnchanged(t *testing.T) {
	src, err := os.ReadFile("../../../workplace/orig/dosv/MMAP.MAP")
	if err != nil {
		t.Skip("找不到原版 MMAP.MAP，跳過")
	}
	want, err := DecodeFile(src)
	if err != nil {
		t.Fatal(err)
	}
	got := Decode(src)
	if len(got) != len(want)+Header {
		t.Fatalf("從 0 解出 %d B，預期 %d ＋ %d", len(got), len(want), Header)
	}
	for i := range want {
		if got[Header+i] != want[i] {
			t.Fatalf("第 %d 個 byte 起就不同了", i)
		}
	}
}
