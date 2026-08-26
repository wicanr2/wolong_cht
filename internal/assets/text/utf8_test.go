package text

import (
	"os"
	"path/filepath"
	"strings"
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

// 語系包的空白是內容不是補位：切掉會把 "struck {2}." 變成 "struckXUCHANG."。
func TestUTF8DecodeKeepsTrailingSpace(t *testing.T) {
	if got := Decode([]byte("struck \x00"), UTF8); got != "struck " {
		t.Fatalf("UTF8 尾端空白被切掉了：%q", got)
	}
	// Big5 的全形補位空白仍然要切——那是定長欄位的填充。
	if got := Decode([]byte{0xA4, 0x40, 0xA1, 0x40}, Big5); got != "一" {
		t.Fatalf("Big5 補位空白沒切掉：%q", got)
	}
}

// 三個語系包都要能載、都要能代入變數（docs/spec/84 §5）。
//
// 這一支擋的是「語系檔與程式漂開」：檔案格式改了、marker 寫錯、
// 則數少一則，都會在這裡當場失敗，而不是等玩到那一則才發現。
func TestAllLanguagePacksLoad(t *testing.T) {
	// 變數值本身也是語系的（名表提供），所以逐語系給對應的城名。
	for _, tc := range []struct {
		file, city, want string
	}{
		{"talk-zh-hans.json", "许昌", "许昌发生了暴风雨。"},
		{"talk-ja.json", "許昌", "許昌で暴風雨が"},
		{"talk-en.json", "XUCHANG", "A storm has struck XUCHANG."},
	} {
		p := filepath.Join("..", "..", "..", "translations", tc.file)
		if _, err := os.Stat(p); err != nil {
			t.Skipf("找不到 %s", tc.file)
		}
		table, err := LoadJSON(p, UTF8)
		if err != nil {
			t.Fatalf("%s: %v", tc.file, err)
		}
		lines, ok := table.Lines(70, map[byte]string{'2': tc.city})
		if !ok || len(lines) == 0 {
			t.Fatalf("%s #70 取不到內容", tc.file)
		}
		if lines[0] != tc.want {
			t.Errorf("%s #70 = %q，預期 %q", tc.file, lines[0], tc.want)
		}
		// 全表掃一遍：marker 一律要被代入或不存在，不能把 `{2}` 印出去。
		for i := 0; i < MessageCount; i++ {
			ls, ok := table.Lines(i, map[byte]string{
				'1': "X", '2': "X", '3': "X", '4': "X", '6': "", '7': "1"})
			if !ok {
				t.Fatalf("%s #%d 代入失敗（有未知 marker）", tc.file, i)
			}
			for _, l := range ls {
				if strings.Contains(l, "{") || strings.Contains(l, "\\") {
					t.Fatalf("%s #%d 殘留未展開的標記：%q", tc.file, i, l)
				}
			}
		}
	}
}
