package main

import (
	"strings"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 一覽表的版面。出處 docs/re/26 §2／§4.1，規格 docs/spec/38。
//
// 原版的一覽視窗是 `sub_181C0(cx = 0x0B18)` ＝ 24 格 × 11 格（一格 16 px）
// ＝ **384 × 176**，左上角 (24, 88)。x 的範圍與命令列相同（24–408）。
const (
	listWinX = 24
	listWinY = 88
	listWinW = 384
	listWinH = 176
	// listRowH 是一列的高度：原版一格就是 16 px。
	listRowH = 16
	// listRowsPerPage 是一頁幾列：**原版的畫列 callback 就是十次迴圈**
	// （docs/re/26 §8.2），而 16（標題）＋ 10 × 16 ＝ 176 正好是視窗高。
	listRowsPerPage = 10
	// listTextInset 是文字距離視窗左緣的內縮。
	listTextInset = 2
)

// listFamily 是一組欄位定義。**標題與分隔線都是原版字串照抄**
// （docs/re/26 §4.1）——分隔線同時就是欄寬定義：
// 全形 `－` 佔 16 px、半形 `-` 佔 8 px，一段連續的 `－`／`-` 就是一欄。
type listFamily struct {
	Title string
	Sep   string
}

// 四個家族（外加開局選勢力）。**「看」與「選」用同一組欄位**——
// 原版八個呼叫端只分五組 (bx, si, di)，兩種取法只差在建清單的 callback
// （docs/re/26 §4.2）。
var (
	listFamilyCorps = listFamily{
		Title: "武將名　總兵數　士氣值 現在位置 目標據點",
		Sep:   "－－－　 ----　　---　　－－－　 －－－",
	}
	listFamilyCities = listFamily{
		Title: "據點名　生產力　上昇率　防災　城兵　內政官",
		Sep:   "－－－　-----　　----　 ---　 --- 　－－－",
	}
	listFamilyGenerals = listFamily{
		Title: "武將名　武術 統率 政治　　勢力　　　身分",
		Sep:   "－－－　 --　 --　 --　 　－－－　　－－－",
	}
	listFamilyFactions = listFamily{
		Title: "勢力名　武將　據點　首都　　外交　　外交官",
		Sep:   "－－－　--- 　--- 　－－－　－－　　－－－",
	}
)

// listField 是一欄：x 是相對視窗內緣的位移，W 是寬度（像素），
// Numeric 為真表示那一欄是半形的數字欄（靠右對齊）。
type listField struct {
	X, W    int
	Numeric bool
}

// fields 把分隔線切成欄。**欄數要與標題的欄位個數相同**，
// 這是 docs/re/26 §4.1 用來確認 `di` 語意的同一條性質。
func (f listFamily) fields() []listField {
	var out []listField
	x, run, numeric := 0, 0, false
	flush := func() {
		if run > 0 {
			out = append(out, listField{X: x - run, W: run, Numeric: numeric})
		}
		run = 0
	}
	for _, r := range f.Sep {
		w := textdraw.GlyphW
		if r < 0x80 {
			w = textdraw.HalfW
		}
		switch r {
		case '－':
			if run > 0 && numeric {
				flush()
			}
			numeric = false
			run += w
		case '-':
			if run > 0 && !numeric {
				flush()
			}
			numeric = true
			run += w
		default:
			flush()
		}
		x += w
	}
	flush()
	return out
}

// listRowY 是第 visible 列的 y（視窗內，標題佔第 0 列）。
func listRowY(visible int) int {
	return listWinY + listRowH + visible*listRowH
}

// listFooterStripY 是「上一頁／下一頁／確定／取消」那一條的 y。
//
// ⚠ **原版沒有這一條**：它靠捲軸翻頁、右鍵取消（說明書 3.8）。
// remake 把它畫在視窗**外面**，才不會吃掉第 10 列。
func listFooterStripY() int { return listWinY + listWinH + chrome.Tile }

// listFieldX 是第 col 欄的左緣。
func listFieldX(fields []listField, col int) int {
	if col < 0 || col >= len(fields) {
		return listWinX + listTextInset
	}
	return listWinX + listTextInset + fields[col].X
}

// listFieldRight 是第 col 欄的右緣（數字欄靠右對齊用）。
func listFieldRight(fields []listField, col int) int {
	if col < 0 || col >= len(fields) {
		return listWinX + listTextInset
	}
	return listWinX + listTextInset + fields[col].X + fields[col].W
}

// listRankNames 是身分名稱表（docs/re/26 §9，段內 0x75A4）。
// **松崗版用「俘虜」不是日文的「捕虜」。**
var listRankNames = [6]string{"－－－", "軍團長", "內政官", "外交官", "俘虜　", "君主　"}

// listDiplomacyNames 是外交關係六級 ＋ 自己那一格（docs/re/27 §4）。
// 索引就是原版表的筆序；交戰那一筆在畫面上會換色。
var listDiplomacyNames = [7]string{"交戰　", "最惡　", "險惡　", "普通　", "良好　", "親密　", "－－　"}

// listDiplomacyLevel 重現 `sub_17A7A`：交友度換成上表的索引。
//
//	最高位元未設      → 0 交戰
//	> 100            → 6 －－（自己那一格）
//	其餘 (v−1)/20 +1 → 1–5
func listDiplomacyLevel(raw int, atWar bool) int {
	if atWar {
		return 0
	}
	v := raw & 0x7F
	if v > 100 {
		return 6
	}
	if v == 100 {
		v = 99
	}
	return v/20 + 1
}

// listBlank 是空欄要印的東西（原版印同寬的全形空白）。
func listBlank(f listField) string {
	n := f.W / textdraw.GlyphW
	if n <= 0 {
		return ""
	}
	return strings.Repeat("　", n)
}
