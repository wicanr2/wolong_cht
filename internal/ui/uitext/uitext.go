// Package uitext 是 remake 的 UI 語系層（docs/spec/84）。
//
// 三層查找：逐句覆寫詞表（ui-<lang>.json）→ 字級轉換表（t2s-chars.json，
// 只有 zh-Hans 用）→ 原文。**fallback 一律回母本繁中**，缺譯要看得見，
// 不顯示空字串。
//
// 這一層刻意不依賴 Ebiten，也不認識畫面——它只做字串到字串的轉換，
// 排版仍由 textdraw 處理。人名地名（SINARIO 的 raw Big5）在呈現層
// 解碼成繁中之後，同樣經 Convert 出去。
package uitext

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Language 是語系代號。
type Language string

const (
	ZhHant Language = "zh-hant" // 母本，不轉換
	ZhHans Language = "zh-hans"
	Ja     Language = "ja" // 日文＝PC-98 原版的文字，不是翻譯
	En     Language = "en"
)

// Latin 回報這個語系是不是以半形拉丁字母為主。
//
// 版面差別很大：中日文一個字兩格、名字三個字就寫完；羅馬化的名字
// 要到 12 格。清單欄界因此另排一套（docs/spec/85）。
func (l Language) Latin() bool { return l == En }

// Table 是載入好的語系轉換狀態。零值＝母本繁中（Convert 原樣回傳）。
type Table struct {
	lang     Language
	override map[string]string // 逐句覆寫：原文 → 譯文
	chars    map[rune]rune     // 字級表（zh-Hans）
	// hasNumeric ＝ 詞表裡有帶 `%d` 的樣板。沒有就不必做樣板查找。
	hasNumeric bool
	// phrases 是片語替換，長的排前面（docs/spec/87 §4）。
	phrases []phrase
}

// phrase 是一條片語替換。
//
// 狀態列與遭遇視窗的句子是**拼出來的**（`名字 + " 對 " + 名字`），
// 名字先被換成羅馬拼音，接縫留在原地——逐句表永遠不會命中，
// 因為那一句每次都不同。片語層只補這個接縫。
type phrase struct{ from, to string }

// digits 用來把「已經填好數字的字串」還原成樣板。
var digits = regexp.MustCompile(`[0-9]+`)

// ParseLanguage 驗證語系代號。
func ParseLanguage(s string) (Language, error) {
	switch Language(strings.ToLower(s)) {
	case ZhHant, "":
		return ZhHant, nil
	case ZhHans:
		return ZhHans, nil
	case Ja:
		return Ja, nil
	case En:
		return En, nil
	}
	return "", fmt.Errorf("uitext: 不支援的語系 %q（可用：zh-hant／zh-hans／ja／en）", s)
}

// Load 載入一個語系。
//
// charsPath 是字級表（只有 zh-Hans 需要，可空）；overridePaths 是逐句詞表
// （UI 詞、人名地名…，可以多個，後面的蓋前面的）。**檔都缺時仍回傳可用的
// Table**——缺檔時 Convert 退回原文，遊戲照跑，缺譯直接以繁中顯示。
//
// 詞表的 JSON 允許兩種形狀：直接的 `{"原文": "譯文"}`，或包一層
// `{"note": ..., "names": {...}}`（名表用這種，好放產生方式的說明）。
func Load(lang Language, charsPath string, overridePaths ...string) (*Table, error) {
	var chars []byte
	if charsPath != "" {
		raw, err := os.ReadFile(charsPath)
		if err != nil {
			return nil, fmt.Errorf("uitext: 讀不到字級表 %s：%w", charsPath, err)
		}
		chars = raw
	}
	var tables [][]byte
	for _, path := range overridePaths {
		if path == "" {
			continue
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("uitext: 讀不到詞表 %s：%w", path, err)
		}
		tables = append(tables, raw)
	}
	return Parse(lang, chars, tables...)
}

// Parse 與 Load 相同，但吃已經讀進來的位元組——語系包會內嵌進執行檔
// （docs/spec/86 §2），那條路沒有檔案可以開。
func Parse(lang Language, chars []byte, tables ...[]byte) (*Table, error) {
	t := &Table{lang: lang, override: map[string]string{}}
	for _, raw := range tables {
		if len(raw) == 0 {
			continue
		}
		const path = "詞表"
		var any map[string]json.RawMessage
		if err := json.Unmarshal(raw, &any); err != nil {
			return nil, fmt.Errorf("uitext: %s 不是有效 JSON：%w", path, err)
		}
		if nested, ok := any["phrases"]; ok {
			var m map[string]string
			if err := json.Unmarshal(nested, &m); err != nil {
				return nil, fmt.Errorf("uitext: %s 的 phrases 不是字串表：%w", path, err)
			}
			for k, v := range m {
				t.phrases = append(t.phrases, phrase{k, v})
			}
			delete(any, "phrases")
		}
		if nested, ok := any["names"]; ok {
			var m map[string]string
			if err := json.Unmarshal(nested, &m); err != nil {
				return nil, fmt.Errorf("uitext: %s 的 names 不是字串表：%w", path, err)
			}
			for k, v := range m {
				t.override[k] = v
			}
			continue
		}
		for k, v := range any {
			var str string
			if err := json.Unmarshal(v, &str); err != nil {
				continue // note 之類的非字串欄位略過
			}
			t.override[k] = str
		}
	}
	for k := range t.override {
		if strings.Contains(k, "%d") {
			t.hasNumeric = true
			break
		}
	}
	// ⭐ **長的排前面**：`" 對 "` 與 `" 對我方宣戰"` 同時存在時，
	// 先換短的會把長的那一條切掉一半。
	sort.SliceStable(t.phrases, func(i, j int) bool {
		return len(t.phrases[i].from) > len(t.phrases[j].from)
	})
	if len(chars) > 0 {
		var m map[string]string
		if err := json.Unmarshal(chars, &m); err != nil {
			return nil, fmt.Errorf("uitext: 字級表不是有效 JSON：%w", err)
		}
		t.chars = make(map[rune]rune, len(m))
		for k, v := range m {
			kr, vr := []rune(k), []rune(v)
			if len(kr) != 1 || len(vr) != 1 {
				return nil, fmt.Errorf("uitext: 字級表的 %q→%q 不是單字對單字", k, v)
			}
			t.chars[kr[0]] = vr[0]
		}
	}
	return t, nil
}

// Lang 回報這張表的語系。
func (t *Table) Lang() Language {
	if t == nil {
		return ZhHant
	}
	return t.lang
}

// Add 補一筆逐句覆寫。給執行期才生成的字串用（例如依欄界算出來的
// 清單標題，docs/spec/85 §4）；nil 表安全忽略。
func (t *Table) Add(from, to string) {
	if t == nil || from == "" {
		return
	}
	if t.override == nil {
		t.override = map[string]string{}
	}
	t.override[from] = to
	if strings.Contains(from, "%d") {
		t.hasNumeric = true
	}
}

// RuneMap 回傳字級表（給 textdraw.SetRuneMap 用）；沒有就回 nil。
func (t *Table) RuneMap() map[rune]rune {
	if t == nil {
		return nil
	}
	return t.chars
}

// Convert 把一句母本繁中換成目前語系的呈現文字。
//
// 覆寫詞表命中就用它（en 只有這一層）；zh-Hans 再走字級表逐字換；
// 都沒有就原樣回傳。
func (t *Table) Convert(s string) string {
	if t == nil || t.lang == ZhHant || s == "" {
		return s
	}
	if v, ok := t.override[s]; ok {
		return v
	}
	if v, ok := t.numeric(s); ok {
		return v
	}
	// ⚠ 片語要在字級表**之前**做：簡體字形換過之後，
	// 片語鍵（繁體）就對不上了。
	s = t.applyPhrases(s)
	if len(t.chars) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if v, ok := t.chars[ch]; ok {
			b.WriteRune(v)
			continue
		}
		b.WriteRune(ch)
	}
	return b.String()
}

// applyPhrases 逐條替換片語。**只動明列的那幾條**，不做任何自動切詞——
// 這一層的風險全部來自「換到了不該換的地方」，把清單握在手上是唯一的
// 控制方式。沒有片語命中時回傳的是同一個字串，一個 byte 都不動。
func (t *Table) applyPhrases(s string) string {
	for _, p := range t.phrases {
		if strings.Contains(s, p.from) {
			s = strings.ReplaceAll(s, p.from, p.to)
		}
	}
	return s
}

// numeric 處理「畫的時候已經填好數字」的標籤。
//
// UI 的數值欄位是 `fmt.Sprintf("兵 %d", n)` 先組好才畫的，所以畫到
// 這一層時看到的是 `兵 300`——**逐句詞表查不到**。把數字抽回 `%d` 去查
// 樣板，命中再把原來的數字依序填回去。`%d` 的個數對不上就放棄，
// 寧可顯示原文也不要把數字擺錯位置。
func (t *Table) numeric(s string) (string, bool) {
	if !t.hasNumeric {
		return "", false
	}
	nums := digits.FindAllString(s, -1)
	if len(nums) == 0 {
		return "", false
	}
	v, ok := t.override[digits.ReplaceAllString(s, "%d")]
	if !ok || strings.Count(v, "%d") != len(nums) {
		return "", false
	}
	var b strings.Builder
	i := 0
	for {
		k := strings.Index(v, "%d")
		if k < 0 {
			b.WriteString(v)
			break
		}
		b.WriteString(v[:k])
		b.WriteString(nums[i])
		v, i = v[k+2:], i+1
	}
	return b.String(), true
}

// ConvertLines 對多列逐列 Convert。
func (t *Table) ConvertLines(lines []string) []string {
	if t == nil || t.lang == ZhHant {
		return lines
	}
	out := make([]string, len(lines))
	for i, ln := range lines {
		out[i] = t.Convert(ln)
	}
	return out
}
