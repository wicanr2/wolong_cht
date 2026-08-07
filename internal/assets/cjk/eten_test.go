package cjk

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// repoRoot 用 runtime.Caller 反推 repo 根目錄，不寫死絕對路徑。
func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	// file 是 .../internal/assets/cjk/eten_test.go
	return filepath.Join(filepath.Dir(file), "..", "..", "..")
}

// etenDir 回傳倚天字型的目錄。字型不隨本專案散布，找不到就跳過。
func etenDir(t *testing.T) string {
	t.Helper()
	for _, d := range []string{
		os.Getenv("ETEN_FONT_DIR"),
		filepath.Join(repoRoot(), "workplace", "eten"),
	} {
		if d == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(d, "STDFONT.15")); err == nil {
			return d
		}
	}
	t.Skip("cjk: 找不到倚天字型（設 ETEN_FONT_DIR 指向含 STDFONT.15/SPCFONT.15 的目錄）")
	return ""
}

func loadFont(t *testing.T) *Font {
	t.Helper()
	d := etenDir(t)
	f, err := Load(filepath.Join(d, "STDFONT.15"), filepath.Join(d, "SPCFONT.15"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return f
}

func TestLoad_GlyphCounts(t *testing.T) {
	f := loadFont(t)
	std, spc := f.GlyphCount()
	if std != 13094 {
		t.Errorf("STDFONT 字數 = %d，預期 13094", std)
	}
	if spc != 408 {
		t.Errorf("SPCFONT 字數 = %d，預期 408", spc)
	}
}

// render 把 glyph 畫成 ASCII art，方便比對。
func render(t *testing.T, f *Font, ch rune) []string {
	t.Helper()
	g, ok := f.Glyph(ch)
	if !ok {
		t.Fatalf("取不到 %q 的字模", ch)
	}
	var out []string
	for y := 0; y < GlyphHeight; y++ {
		var sb strings.Builder
		for x := 0; x < GlyphWidth; x++ {
			if g.AlphaAt(x, y).A != 0 {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		out = append(out, sb.String())
	}
	return out
}

// 索引公式的 oracle：「一」只有一列是連續橫線，其餘全空。
//
// 這一關沒過就別往下做 —— 整批字會整體偏移，看起來像「有字但都不對」。
func TestGlyph_YiIsASingleHorizontalLine(t *testing.T) {
	f := loadFont(t)

	dense := 0
	for _, row := range render(t, f, '一') {
		if strings.Count(row, "#") >= 12 {
			dense++
		}
	}
	if dense != 1 {
		t.Errorf("「一」應只有一列是連續橫線，實際有 %d 列", dense)
	}
}

// 幾個結構明確的字：驗證取到的不是空白也不是雜訊。
func TestGlyph_KnownCharactersAreDense(t *testing.T) {
	f := loadFont(t)

	for _, ch := range []rune{'中', '猴', '冬', '魔', '國', '一'} {
		rows := render(t, f, ch)
		ink := 0
		for _, r := range rows {
			ink += strings.Count(r, "#")
		}
		if ink < 10 {
			t.Errorf("%q 的筆劃只有 %d 個像素，疑似取到空白或索引錯誤", ch, ink)
		}
		if ink > GlyphWidth*GlyphHeight-10 {
			t.Errorf("%q 幾乎全黑（%d 像素），疑似取到雜訊", ch, ink)
		}
	}
}

// 把「一」的點陣逐位元釘死。這是索引公式最硬的一道驗收 ——
// 偏一格就會取到別的字，位元圖不可能還對得上。
//
// 刻意不用「左右對稱」之類的結構性判準：倚天的字模在 16 格內不是置中的
// （「中」的豎筆在 x=7、外框佔 x=1..14），對稱性測不出東西。
func TestGlyph_YiExactBitmap(t *testing.T) {
	f := loadFont(t)
	rows := render(t, f, '一')

	for y, row := range rows {
		want := strings.Repeat(".", GlyphWidth)
		switch y {
		case 6:
			want = "............." + "#" + ".."
		case 7:
			want = strings.Repeat("#", 15) + "."
		}
		if row != want {
			t.Errorf("「一」第 %d 列 = %q，預期 %q", y, row, want)
		}
	}
}

// 全形標點在 SPCFONT 不在 STDFONT。漏帶 SPCFONT 這些字會全部 fallback。
func TestGlyph_FullWidthPunctuationComesFromSPCFONT(t *testing.T) {
	f := loadFont(t)

	for _, ch := range []rune{'，', '。', '！', '？', '「', '」', '（', '）', '《', '》'} {
		hi, lo, ok := Big5(ch)
		if !ok {
			t.Errorf("%q 無法編成 Big5", ch)
			continue
		}
		src, _ := f.locate(hi, lo)
		if len(src) != len(f.spc) {
			t.Errorf("%q 應取自 SPCFONT，實際取自 STDFONT", ch)
		}
		if _, ok := f.Glyph(ch); !ok {
			t.Errorf("%q 取不到字模", ch)
		}
	}
}

// 手動補的字：Go 的 Big5 編碼器對 ～ 有歧義。
func TestBig5_ManualOverride(t *testing.T) {
	hi, lo, ok := Big5('～')
	if !ok {
		t.Fatal("～ 應由手動表補上")
	}
	if hi != 0xA1 || lo != 0xE3 {
		t.Errorf("～ 的 Big5 = %02X%02X，預期 A1E3", hi, lo)
	}
}

// Big5 沒有的字要回報 fallback，不能靜默給錯的字模。
func TestGlyph_NonBig5FallsBack(t *testing.T) {
	f := loadFont(t)

	// Big5 收錄的範圍比直覺廣 —— 日文假名（あ）與希臘字母（Ω）都在裡面，
	// 拿它們當「非 Big5」的例子會誤判。真正不在的是簡體字與非 BMP 字。
	//
	// 這批簡體字有一半會被 Go 的 Big5 編碼器映進造字區（`马`→`89C6`、
	// `着`→`FED3`），編碼「成功」但字型不涵蓋。列進來釘住那個破口。
	for _, ch := range []rune{'马', '仅', '试', '头', '顶', '着', '门', '东', '车', '国', '长', '买'} {
		if ch == '只' || ch == '儿' {
			continue
		}
		if _, ok := f.Glyph(ch); ok {
			t.Errorf("%q 是簡體字、不在 Big5 裡，應回報 fallback", ch)
		}
	}

	// 反過來確認：假名與希臘字母確實取得到，不該被誤判成缺字。
	for _, ch := range []rune{'あ', 'Ω'} {
		if _, ok := f.Glyph(ch); !ok {
			t.Errorf("%q 在 Big5 裡，應取得到字模", ch)
		}
	}
}

// 本專案 glossary 與手冊裡實際會用到的字，fallback 數量必須是 0。
//
// fallback 數量是品質指標：一大批字掉進 fallback 時，
// 先懷疑索引公式或漏帶 SPCFONT，不要無腦補字型。
func TestGlyph_ProjectVocabularyHasNoFallback(t *testing.T) {
	f := loadFont(t)

	// 這一串**不是手打的詞彙表**，是從原版 SINARIO.DAT 抽出來的：
	// 四個劇本全部武將名（含呼び名）＋ 全部 192 個據點名的相異字，
	// 共 377 個。只要有一個取不到字模，畫面上就會缺字。
	//
	// 用一手資料當測試樣本，而不是憑印象列詞——
	// 憑印象會漏掉「叡」「懿」「廮」這種只在特定劇本出現的字。
	const names = "丁丈上下且丕丘中九乾五京亭亮仁代仲任伍休余侯信倉傕備公六典冀冷凌函刑別劉功勳化北南原叡口史司合吳呂周嚴圃國地城埽壁壺壽夏天太夷奉始姜子孔孟孫安宕官定宛宮容富寧審寵封尚山岱州巫巴布帝師平度庶庸廖廣廩廬廮延建弁弋弘張彧彰徐德志忠恢恪息惇慈懷懿成房挺操攸散文新旋昂昌明易昕春昭昱晃晉普景暹曲曹會朗朱李杭東松柴桂桑桓梁梓植椎楊業榮槐樂樊橋權欽正武毗水氾汜汝江池沓沙沛沮河沾法泰洛洪海涇涪涼涿淄淮淵清渚渠渡渭湘溫滎滿漆漢漳潯潼濟濡濮牛牢狹獲玄王琦瑁瑜瑯瑾璋瓚甘田界留登白皖皮益盛盧真磐祖祝禪秭程稽竟章竹笣笮策系紀純紘索紹統維綿縣繇繡羽習聘肅肥胤臨興舒舞良艾苑苞英荀莞華萌葛葭蒙蓋蔡蔣蔪蘭虎融術街衛衡表袁褚襄西許詡談諸謖譙譚谷豐賈賢赤超越趙輿農遂道達遜遼邑邪邳郃郝郡郭郯都鄧鄱鄴配醜里野金銅鎮長閬關阪陜陰陳陵陶陸陽雒雙雛雲零雷霸靈鞏韋韓韶須顏風飛館馬騰高魏魯鳳黃黎龐"

	// 介面與訊息會用到的標點與數字。
	const ui = "，。、！？「」『』（）：；．０１２３４５６７８９年月日"

	var missing []rune
	for _, ch := range names + ui {
		if _, ok := f.Glyph(ch); !ok {
			missing = append(missing, ch)
		}
	}
	if len(missing) != 0 {
		t.Errorf("以下 %d 個原版用字取不到字模：%q", len(missing), string(missing))
	}
}

func TestLoad_RejectsMalformedFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.15")
	if err := os.WriteFile(bad, []byte{1, 2, 3}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(bad, bad); err == nil {
		t.Error("長度不是 30 的整數倍時應回傳錯誤")
	}
}

func TestEmboldenRow_GrowsRightAcrossByteBoundary(t *testing.T) {
	// x=7 是第一個 byte 的末位；加粗後 x=8 也必須亮，不能各 byte 分開處理。
	const x7 = uint16(1) << (15 - 7)
	got := emboldenRow(x7)
	want := x7 | uint16(1)<<(15-8)
	if got != want {
		t.Errorf("emboldenRow(%016b) = %016b，預期 %016b", x7, got, want)
	}
}

func TestEmboldenRow_DoesNotGrowLeftOrWrap(t *testing.T) {
	const (
		left  = uint16(1) << 15
		right = uint16(1)
	)
	if got := emboldenRow(left); got != left|left>>1 {
		t.Errorf("最左像素加粗 = %016b", got)
	}
	if got := emboldenRow(right); got != right {
		t.Errorf("最右像素不應溢位或回繞，得到 %016b", got)
	}
}

func TestGlyph_BoldContainsNormalAndAddsInk(t *testing.T) {
	d := etenDir(t)
	normal, err := LoadDir(d, Options{})
	if err != nil {
		t.Fatal(err)
	}
	bold, err := LoadDir(d, Options{Bold: true})
	if err != nil {
		t.Fatal(err)
	}
	n, _ := normal.Glyph('冬')
	b, _ := bold.Glyph('冬')
	added := 0
	for y := 0; y < GlyphHeight; y++ {
		for x := 0; x < GlyphWidth; x++ {
			na, ba := n.AlphaAt(x, y).A, b.AlphaAt(x, y).A
			if na != 0 && ba == 0 {
				t.Fatalf("粗體遺失原字像素 (%d,%d)", x, y)
			}
			if na == 0 && ba != 0 {
				added++
			}
		}
	}
	if added == 0 {
		t.Error("粗體沒有增加任何像素")
	}
}

// Go 的 Big5 編碼器會把一批簡體字映進造字區，編碼成功但沒有字模。
//
// Big() 必須驗碼位落在標準 Big5 區間（首 A1–F9、次 40–7E／A1–FE），
// 否則檢查工具會放行簡體字，畫面上才發現是一片空白。
func TestBig5_RejectsUserDefinedArea(t *testing.T) {
	// 注意：`只`（A575）與`儿`（A449）看起來像簡體，其實是標準 Big5 收錄的字，
	// 不能列進來 —— 「簡體字一定不在 Big5」這個直覺是錯的。
	for _, ch := range []rune{'头', '着', '门', '东', '车', '马', '国', '长'} {
		if _, _, ok := Big5(ch); ok {
			t.Errorf("%q 落在造字區，應判定為非 Big5", ch)
		}
	}
	// 反向確認沒有誤殺：標準區的字仍要通得過。
	for _, ch := range []rune{'冬', '魔', '一', '，', '「', 'あ', 'Ω'} {
		if _, _, ok := Big5(ch); !ok {
			t.Errorf("%q 在標準 Big5 裡，不該被擋掉", ch)
		}
	}
}
