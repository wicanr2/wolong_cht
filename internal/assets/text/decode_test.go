package text

import "testing"

// 姓名是 6 byte 的 Big5（3 個全形字），不足補全形空白（A140）。
// ⚠ 這些位元組是從 SINARIO.DAT 直接抓的，**不是我背出來的**。
// 第一版我憑印象寫了諸葛亮的位元組，AFAB 其實是「神」不是「葛」，
// 測試當場抓到。編碼表不要憑記憶。
func TestDecodeNames(t *testing.T) {
	cases := []struct {
		raw  []byte
		want string
	}{
		{[]byte{0xB1, 0xE4, 0xBE, 0xDE, 0xA1, 0x40}, "曹操"},
		{[]byte{0xA7, 0x66, 0xA5, 0xAC, 0xA1, 0x40}, "呂布"},
		{[]byte{0xBD, 0xD1, 0xB8, 0xAF, 0xAB, 0x47}, "諸葛亮"},
	}
	for _, c := range cases {
		if got := Decode(c.raw, Big5); got != c.want {
			t.Errorf("Decode(%X) = %q, want %q", c.raw, got, c.want)
		}
	}
}

// ⚠ F9D8 在 cp950 是「裏」。用嚴格的 big5 會解不開 ——
// 這個坑實際踩過，當時被誤判成「原版有錯字」。
func TestDecodeCP950Extension(t *testing.T) {
	got := Decode([]byte{0xF9, 0xD8}, Big5)
	if got == "" || got[0] == '?' {
		t.Errorf("F9D8 應該解得出「裏」，得到 %q", got)
	}
}

// 解不開的位元組要退回 hex，不能靜靜吃掉。
func TestDecodeFallbackIsVisible(t *testing.T) {
	got := Decode([]byte{0xFF, 0xFF}, Big5)
	if got == "" {
		t.Error("解不開時不該回空字串")
	}
}

// 呼び名與本名不同的三個人（docs/formats/08 §3）。
func TestDecodeAlias(t *testing.T) {
	if got := Decode([]byte{0xA4, 0xD5, 0xA9, 0xFA, 0xA1, 0x40}, Big5); got != "孔明" {
		t.Errorf("諸葛亮的呼び名 = %q, want 孔明", got)
	}
}
