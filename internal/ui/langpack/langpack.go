// Package langpack 把「換一個語言要動哪些東西」收在一處（docs/spec/86）。
//
// 桌面與手機各有一套啟動流程，但語系的內容完全一樣——分成兩份寫的話
// 會長出不同的行為（CLAUDE.md §7 第 6 條：一條規則只留一份實作）。
//
// 語系檔**內嵌進執行檔**（共約 370 KB），檔案系統上的同名檔優先：
// 開發時改 `translations/*.json` 立即生效，Android 端則零設定。
// 原版資產仍然一個都不進執行檔——那條路是 `docs/spec/72` 的 `gamedata/`。
package langpack

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/assets/cjk"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
	"github.com/wicanr2/wolong_cht/internal/ui/uitext"
	"github.com/wicanr2/wolong_cht/translations"
)

// SearchPaths 是找語系檔的目錄，依序試；找不到才用內嵌的那一份。
// 呼叫端可以加上「執行檔旁邊」之類的路徑。
var SearchPaths = []string{"translations"}

// Pack 是一個語系的全部呈現資料。零值（Lang 為母本）代表不轉換。
type Pack struct {
	Lang uitext.Language
	// Text 是 UI 詞、人名地名與字級表；母本語系是 nil。
	Text *uitext.Table
	// Talk 是這個語系的 1,022 則訊息；母本語系是 nil（用原版解出來的那份）。
	Talk *text.Table
	// Font 是全形字型鏈；載不到就是 nil（呈現層畫缺字框）。
	Font textdraw.GlyphSource
}

// read 依「檔案系統優先、內嵌墊底」的順序取一份語系檔。
// 兩邊都沒有回 (nil, false)——**缺語系檔不該讓遊戲開不了**，
// 呼叫端退回母本繁中就好。
func read(name string) ([]byte, bool) {
	for _, dir := range SearchPaths {
		if dir == "" {
			continue
		}
		if b, err := os.ReadFile(filepath.Join(dir, name)); err == nil {
			return b, true
		}
	}
	if b, err := fs.ReadFile(translations.Files, name); err == nil {
		return b, true
	}
	return nil, false
}

// Available 回報有沒有這個語系的訊息包（母本永遠有）。
func Available(lang uitext.Language) bool {
	if lang == uitext.ZhHant {
		return true
	}
	_, ok := read(talkFile(lang))
	return ok
}

func talkFile(lang uitext.Language) string {
	switch lang {
	case uitext.ZhHans:
		return "talk-zh-hans.json"
	case uitext.Ja:
		return "talk-ja.json"
	case uitext.En:
		return "talk-en.json"
	}
	return ""
}

// Load 組出一個語系的全部呈現資料。
//
// fontDir 是使用者自備的點陣字目錄；載不到字型只是缺字，不是錯誤。
// **缺任何一塊都不擋**——缺的部分退回母本，缺譯要看得見。
func Load(lang uitext.Language, fontDir string) (*Pack, error) {
	p := &Pack{Lang: lang}
	if lang == uitext.ZhHant {
		p.Font = fontChain(lang, fontDir)
		return p, nil
	}

	if name := talkFile(lang); name != "" {
		if raw, ok := read(name); ok {
			t, err := text.ParseJSON(raw, text.UTF8)
			if err != nil {
				return nil, fmt.Errorf("langpack: %s：%w", name, err)
			}
			p.Talk = t
		}
	}

	var tables [][]byte
	for _, name := range []string{
		fmt.Sprintf("ui-%s.json", lang),
		fmt.Sprintf("names-%s.json", lang),
	} {
		if raw, ok := read(name); ok {
			tables = append(tables, raw)
		}
	}
	var chars []byte
	if lang == uitext.ZhHans {
		if raw, ok := read("t2s-chars.json"); ok {
			chars = raw
		}
	}
	t, err := uitext.Parse(lang, chars, tables...)
	if err != nil {
		return nil, fmt.Errorf("langpack: %s 的詞表：%w", lang, err)
	}
	p.Text = t
	p.Font = fontChain(lang, fontDir)
	return p, nil
}

// builtin 是原版資料裡的內建字型（`END_S13.DAT`，docs/spec/137）。
// 設了就當**主要**字型——那是遊戲自己用的那一份，全形標點與倚天不同。
var builtin textdraw.GlyphSource

// SetBuiltinFont 掛上原版內建的字型。傳 nil 表示沒有，退回倚天鏈。
func SetBuiltinFont(f textdraw.GlyphSource) { builtin = f }

// fontChain 依語系決定字型的取用順序。
//
// **語系的字集不等於一份字型的字集**：日文人名裡的 PC-98 外字不在
// JIS X 0208 但在倚天 Big5 裡有，所以要接成鏈（docs/spec/84 §2.1）。
func fontChain(lang uitext.Language, dir string) textdraw.GlyphSource {
	if dir == "" {
		return nil
	}
	var primary textdraw.GlyphSource
	switch lang {
	case uitext.ZhHans:
		if f, err := cjk.LoadKuTen16Dir(dir, cjk.GB2312, cjk.Options{}); err == nil {
			primary = f
		}
	case uitext.Ja:
		if f, err := cjk.LoadKuTen16Dir(dir, cjk.JISX0208, cjk.Options{}); err == nil {
			primary = f
		}
	}
	var eten textdraw.GlyphSource
	if f, err := cjk.LoadDir(dir, cjk.Options{}); err == nil {
		eten = f
	}
	// ⭐ **每個語系都把另外兩套字集接在後面當墊底。**
	// 語言選單要用各語言自己的寫法寫（「简体中文」「日本語」），
	// 而「简」不在 Big5、平假名不在 GB2312——只掛自己那一套的話，
	// 選單上會有幾格方框，偏偏那幾格正是玩家要靠它認語言的地方。
	// 原版文字用不到這些碼位，所以墊底不影響既有的逐像素對拍。
	rest := make([]textdraw.GlyphSource, 0, 2)
	for _, cs := range []cjk.Charset{cjk.GB2312, cjk.JISX0208} {
		if (cs == cjk.GB2312 && lang == uitext.ZhHans) || (cs == cjk.JISX0208 && lang == uitext.Ja) {
			continue // 已經是 primary
		}
		if f, err := cjk.LoadKuTen16Dir(dir, cs, cjk.Options{}); err == nil {
			rest = append(rest, f)
		}
	}
	// ⭐ **內建字型排在倚天前面**：全形標點那 408 格兩份不一樣
	// （docs/spec/137），要逐像素對上原版就得用原版自己那一份。
	// 只對繁中母本成立——其他語系的主要字集另有來源。
	sources := make([]textdraw.GlyphSource, 0, 4+len(rest))
	if lang == uitext.ZhHant && builtin != nil {
		sources = append(sources, builtin)
	}
	sources = append(sources, primary, eten)
	sources = append(sources, rest...)
	return textdraw.Chain(sources...)
}

// Apply 把這個語系掛到 Drawer 上。
//
// ⚠ **切回母本時要把上一個語系清掉**：`SetRuneMap(nil)`／
// `SetTranslator(nil)` 都要送出去，否則簡體的字級表會留在那裡，
// 換成日文之後每個漢字仍然被換成簡體字形。
func (p *Pack) Apply(d *textdraw.Drawer) {
	if d == nil {
		return
	}
	if p == nil {
		d.SetRuneMap(nil)
		d.SetTranslator(nil)
		return
	}
	d.SetRuneMap(p.Text.RuneMap())
	if p.Text == nil {
		d.SetTranslator(nil)
	} else {
		d.SetTranslator(p.Text.Convert)
	}
	if p.Font != nil {
		d.SetFont(p.Font)
	}
}

// Convert 是這個語系的字串轉換（人名地名、UI 詞）。nil-safe。
func (p *Pack) Convert(s string) string {
	if p == nil {
		return s
	}
	return p.Text.Convert(s)
}

// Choice 是選單上的一個語言。
type Choice struct {
	Lang uitext.Language
	// Name **用該語言自己的寫法**。換過去之後畫面全是那個語言，
	// 用母本的寫法列出來，玩家反而認不出自己選了什麼。
	Name string
}

// Choices 是桌面啟動殼層與手機系統面板共用的語言清單，
// 也是 F9 的循環順序（docs/spec/86 §4）。
var Choices = []Choice{
	{uitext.ZhHant, "繁體中文"},
	{uitext.ZhHans, "简体中文"},
	{uitext.Ja, "日本語"},
	{uitext.En, "English"},
}

// Next 回傳循環順序的下一個語系（F9 用，docs/spec/86 §4）。
// 跳過沒有語系檔的那些，繞一圈都沒有就回母本。
func Next(cur uitext.Language) uitext.Language {
	order := make([]uitext.Language, len(Choices))
	for i, c := range Choices {
		order[i] = c.Lang
	}
	at := 0
	for i, l := range order {
		if l == cur {
			at = i
			break
		}
	}
	for i := 1; i <= len(order); i++ {
		next := order[(at+i)%len(order)]
		if Available(next) {
			return next
		}
	}
	return uitext.ZhHant
}
