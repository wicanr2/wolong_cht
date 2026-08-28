package namechars

import (
	"os"
	"path/filepath"
	"testing"
)

// 解碼式的正對照：拿已知的兩個字（ㄅ＝A374、八＝A44B）照原版算式編回去。
func TestDecodeUndoesSwapAndOffset(t *testing.T) {
	enc := func(big5 uint16) []byte {
		sw := big5 + 0x1000 // 對調前先加（原版是對調後減，等價）
		w := (sw&0xFF)<<8 | sw>>8
		return []byte{byte(w), byte(w >> 8)}
	}
	data := append(enc(0xA374), enc(0xA44B)...)
	tb, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if string(tb.Runes) != "ㄅ八" || tb.Big5[1] != 0xA44B {
		t.Fatalf("解出 %q %04X", string(tb.Runes), tb.Big5)
	}
}

// 真檔（使用者自備，缺就跳過）：2,621 字、開頭照注音 ㄅ 起、全部是 Big5。
func TestRealTableIsBopomofoOrdered(t *testing.T) {
	path := filepath.Join("..", "..", "..", "workplace", "orig", "dosv", "END_S15.DAT")
	if _, err := os.Stat(path); err != nil {
		t.Skip("沒有原版素材")
	}
	tb, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(tb.Runes) != Count {
		t.Fatalf("%d 字，預期 %d", len(tb.Runes), Count)
	}
	if got := string(tb.Runes[:4]); got != "ㄅ八巴疤" {
		t.Fatalf("開頭是 %q", got)
	}
}
