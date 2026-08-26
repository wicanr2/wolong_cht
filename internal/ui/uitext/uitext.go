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
	"strings"
)

// Language 是語系代號。
type Language string

const (
	ZhHant Language = "zh-hant" // 母本，不轉換
	ZhHans Language = "zh-hans"
	En     Language = "en"
)

// Table 是載入好的語系轉換狀態。零值＝母本繁中（Convert 原樣回傳）。
type Table struct {
	lang     Language
	override map[string]string // 逐句覆寫：原文 → 譯文
	chars    map[rune]rune     // 字級表（zh-Hans）
}

// ParseLanguage 驗證語系代號。
func ParseLanguage(s string) (Language, error) {
	switch Language(strings.ToLower(s)) {
	case ZhHant, "":
		return ZhHant, nil
	case ZhHans:
		return ZhHans, nil
	case En:
		return En, nil
	}
	return "", fmt.Errorf("uitext: 不支援的語系 %q（可用：zh-hant／zh-hans／en）", s)
}

// Load 載入一個語系。
//
// overridePath 是逐句詞表（可空）；charsPath 是字級表（只有 zh-Hans 需要，
// 可空）。**兩個檔都缺時仍回傳可用的 Table**——缺檔時 Convert 退回原文，
// 遊戲照跑，缺譯直接以繁中顯示。
func Load(lang Language, overridePath, charsPath string) (*Table, error) {
	t := &Table{lang: lang}
	if overridePath != "" {
		raw, err := os.ReadFile(overridePath)
		if err != nil {
			return nil, fmt.Errorf("uitext: 讀不到詞表 %s：%w", overridePath, err)
		}
		if err := json.Unmarshal(raw, &t.override); err != nil {
			return nil, fmt.Errorf("uitext: %s 不是有效 JSON：%w", overridePath, err)
		}
	}
	if charsPath != "" {
		raw, err := os.ReadFile(charsPath)
		if err != nil {
			return nil, fmt.Errorf("uitext: 讀不到字級表 %s：%w", charsPath, err)
		}
		var m map[string]string
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("uitext: %s 不是有效 JSON：%w", charsPath, err)
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
