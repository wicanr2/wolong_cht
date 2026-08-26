package text

import (
	"os"
	"path/filepath"
	"testing"
)

// UTF8 語系包（docs/spec/84）：Raw 直接存 UTF-8，Lines() 用表內編碼解，
// 簡體字不會被 Big5 編碼卡死。
func TestLoadJSONUTF8SimplifiedPack(t *testing.T) {
	p := filepath.Join("..", "..", "..", "translations", "talk-zh-hans.json")
	if _, err := os.Stat(p); err != nil {
		t.Skip("找不到 talk-zh-hans.json")
	}
	table, err := LoadJSON(p, UTF8)
	if err != nil {
		t.Fatal(err)
	}
	// #70「{2}发生了暴风雨。」：marker 展開＋簡體原樣出來。
	lines, ok := table.Lines(70, map[byte]string{'2': "许昌"})
	if !ok || len(lines) == 0 || lines[0] != "许昌发生了暴风雨。" {
		t.Fatalf("#70 = %#v ok=%v", lines, ok)
	}
	// marker 缺值仍要 fail-closed，與母本表同一套行為。
	if _, ok := table.Lines(70, nil); ok {
		t.Fatal("缺變數應 fail-closed")
	}
}

// 零值 Table 的編碼是 Big5——既有 Parse 呼叫端不受 enc 欄位影響。
func TestParseTableStaysBig5(t *testing.T) {
	// 第 0 則是「馬\0」（3 byte），其餘全部指到 body 結尾＝空則。
	// 0xB0 0xA8 是 Big5 的「馬」。
	body := []byte{0xB0, 0xA8, 0x00}
	raw := make([]byte, TableBytes)
	for i := 0; i < TableEntries; i++ {
		off := TableBytes + len(body)
		if i == 0 {
			off = TableBytes
		}
		raw[i*2] = byte(off & 0xFF)
		raw[i*2+1] = byte(off >> 8)
	}
	table, err := Parse(append(raw, body...))
	if err != nil {
		t.Fatal(err)
	}
	lines, ok := table.Lines(0, nil)
	if !ok || len(lines) == 0 || lines[0] != "馬" {
		t.Fatalf("Big5 預設解碼壞了：%#v ok=%v", lines, ok)
	}
}
