package cjk

import (
	"os"
	"testing"
)

// 轉出來的區位字型要真的取得到字模（tools/pcf2raw.py → 載入器）。
//
// 字型檔不進版控，所以缺檔就跳過——與原版素材同一個處理方式。
// 這一支擋的是「轉換工具與載入器對索引的理解不一致」：兩邊各自看起來
// 都對，合起來卻整批偏移，症狀是「有字但都不是那個字」。
func TestConvertedKuTenFontsResolveGlyphs(t *testing.T) {
	const dir = "../../../workplace/eten"
	for _, tc := range []struct {
		name  string
		cs    Charset
		probe []rune
	}{
		{"JISKAN16", JISX0208, []rune{'あ', 'ア', '国', '発', '戦', '軍', '拠', '点'}},
		{"HZK16", GB2312, []rune{'啊', '国', '发', '战', '军', '据', '点'}},
	} {
		if _, err := os.Stat(dir + "/" + tc.name); err != nil {
			t.Logf("跳過 %s：%v", tc.name, err)
			continue
		}
		f, err := LoadKuTen16Dir(dir, tc.cs, Options{})
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		for _, ch := range tc.probe {
			a, ok := f.Glyph(ch)
			if !ok || countAlpha(a) < 8 {
				t.Errorf("%s 的「%c」缺字或幾乎全空（ok=%v）", tc.name, ch, ok)
			}
		}
	}
}
