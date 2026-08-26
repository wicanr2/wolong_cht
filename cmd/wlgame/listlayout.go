package main

import (
	"image"
	"image/color"
	"strings"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
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
	// listTextInset 是**資料列文字欄**距離欄起點的縮排：實機量測是
	// 半格 8 px（docs/spec/38 §1.5）。標題列、空列分隔線與數字欄不縮。
	listTextInset = 8
	// listScrollW 是左邊那條捲軸的寬度。**清單本體從 x+16 起，不是 x**
	// ——原版的清單熱區與反白列都從 `word_181AE + 0x10` 開始
	// （docs/re/26 §10）。
	listScrollW = 16
)

// listBodyX 是清單本體的左緣（捲軸右邊）。
func listBodyX() int { return listWinX + listScrollW }

// listBodyW 是清單本體的寬度。
func listBodyW() int { return listWinW - listScrollW }

// 捲軸三段（docs/re/26 §10 的熱區 0x3F–0x41）。三段在垂直方向接續：
// ▲ 佔 y+16–y+32、槽佔到 y+H−16、▼ 到 y+H。
func listScrollUpRect() image.Rectangle {
	return image.Rect(listWinX, listWinY+listRowH,
		listWinX+listScrollW, listWinY+2*listRowH)
}

func listScrollDownRect() image.Rectangle {
	return image.Rect(listWinX, listWinY+listWinH-listRowH,
		listWinX+listScrollW, listWinY+listWinH)
}

func listScrollTrackRect() image.Rectangle {
	return image.Rect(listWinX, listWinY+2*listRowH,
		listWinX+listScrollW, listWinY+listWinH-listRowH)
}

// listScrollThumbRect 是滑塊。**原版畫滑塊的常式沒讀**（docs/spec/38 §4），
// 這裡照 `top ÷ (筆數 − 一頁)` 取比例，是 remake 的畫法。
func listScrollThumbRect(top, total int) image.Rectangle {
	track := listScrollTrackRect()
	span := total - listRowsPerPage
	if span <= 0 {
		return track
	}
	h := track.Dy() * listRowsPerPage / total
	if h < listRowH {
		h = listRowH
	}
	y := track.Min.Y + (track.Dy()-h)*clamp(top, 0, span)/span
	return image.Rect(track.Min.X, y, track.Max.X, y+h)
}

// listFamily 是一組欄位定義。**標題與分隔線都是原版字串照抄**
// （docs/re/26 §4.1）——分隔線同時就是欄寬定義：
// 全形 `－` 佔 16 px、半形 `-` 佔 8 px，一段連續的 `－`／`-` 就是一欄。
type listFamily struct {
	Title string
	Sep   string
	// Extra 是分隔線之後那個**無標題**欄的全形字數。
	// 軍團表最右端有一格：委任中印「委任」，否則印兩個全形空白
	// （docs/re/27 §2）。它不在分隔線裡，所以要另外給。
	Extra int
}

// 四個家族（外加開局選勢力）。**「看」與「選」用同一組欄位**——
// 原版八個呼叫端只分五組 (bx, si, di)，兩種取法只差在建清單的 callback
// （docs/re/26 §4.2）。
var (
	listFamilyCorps = listFamily{
		Title: "武將名　總兵數　士氣值 現在位置 目標據點",
		// 數字欄 [80,112)／[144,168)（parity-menus7 的 m1：值右緣
		// 150／206、空列破折號左緣 120／184 ＝ body 40 ＋ 欄界）。
		Sep:   "－－－　　----　　---　 －－－　 －－－",
		Extra: 2,
	}
	listFamilyCities = listFamily{
		Title: "據點名　生產力　上昇率　防災　城兵　內政官",
		// 數字欄界照 parity-menus8 的 p0 實測（官渡列，值同基準）：
		// 生產力 [72,112)、上昇率 [144,176)、防災 [200,224)、
		// 城兵 [248,280)、內政官 X=288。
		Sep: "－－－　 -----　　----　 ---　 ---- －－－",
	}
	listFamilyGenerals = listFamily{
		Title: "武將名　武術 統率 政治　　勢力　　　身分",
		Sep:   "－－－　 --　 --　 --　 　－－－　　－－－",
	}
	listFamilyFactions = listFamily{
		Title: "勢力名　武將　據點　首都　　外交　　外交官",
		// 數字欄是**四格**破折號（右緣 136／184）——實機 orig-w3-target
		// 的「12」「1」右緣各在 134／181，比三格的右緣（128／176）多 8。
		Sep: "－－－　----　----　－－－　－－　　－－－",
	}
)

// ── 半形語系的欄界（docs/spec/85）─────────────────────────────────
//
// 為什麼要另一套：羅馬化的人名最長 12 個半形格，而原版的姓名欄只有 6 格
// ——照原欄界畫，`夏侯淵` 與 `夏侯惇` 都會被裁成 `XIAHOU`，**兩個不同的
// 人在畫面上長得一樣**。四個家族的內容加起來吃不下本體的 368 px，
// 所以每一家都標明了誰讓步（spec/85 §3 有覆蓋率）。
//
// 中文與日文那一套（`fields()`，照原版量出來的）一個像素都不動。
//
// ⚠ **數值欄的間距是被標題決定的，不是被數字決定的**：數字只有兩位
// （16 px），但標題 `Mar`／`Cmd`／`Pol` 各要 24 px ＋ 一格間隔，
// 所以欄距拉到 32 px。照數字的寬度排會得到 `MarCmdPolFaction`。
var (
	latinFieldsGenerals = []listField{
		{X: 0, W: 96},                  // 武將名：12 字，149/149 放得下
		{X: 112, W: 16, Numeric: true}, // 武術
		{X: 144, W: 16, Numeric: true}, // 統率
		{X: 176, W: 16, Numeric: true}, // 政治
		{X: 208, W: 72},                // 勢力（君主名）：9 字
		{X: 296, W: 64},                // 身分
	}
	latinLabelsGenerals = []string{"Name", "Mar", "Cmd", "Pol", "Faction", "Rank"}

	latinFieldsCities = []listField{
		{X: 0, W: 96},                  // 據點名：12 字，194/194
		{X: 112, W: 40, Numeric: true}, // 生產力
		{X: 160, W: 32, Numeric: true}, // 上昇率
		{X: 200, W: 24, Numeric: true}, // 防災
		{X: 232, W: 32, Numeric: true}, // 城兵
		{X: 272, W: 88},                // 內政官：10 字，8 個會裁
	}
	latinLabelsCities = []string{"City", "Prod.", "Grow", "Fld", "Garr", "Gov."}

	latinFieldsCorps = []listField{
		{X: 0, W: 80},                  // 武將名：10 字
		{X: 96, W: 32, Numeric: true},  // 總兵數
		{X: 136, W: 24, Numeric: true}, // 士氣值
		{X: 168, W: 64},                // 現在位置：8 字
		{X: 248, W: 64},                // 目標據點：8 字
		{X: 328, W: 32, NoDash: true},  // 委任（無標題）
	}
	latinLabelsCorps = []string{"Name", "Men", "Mrl", "Position", "Target", ""}

	latinFieldsFactions = []listField{
		{X: 0, W: 72},                  // 勢力（君主名）：9 字
		{X: 88, W: 24, Numeric: true},  // 武將數
		{X: 128, W: 24, Numeric: true}, // 據點數
		{X: 160, W: 64},                // 首都（據點名）：8 字
		{X: 240, W: 32},                // 外交
		{X: 288, W: 64},                // 外交官（武將名）：8 字
	}
	latinLabelsFactions = []string{"Faction", "Gen", "Cty", "Capital", "Rel.", "Envoy"}
)

// latinTitle 依欄界生成標題列。
//
// 原版的標題是**一整條字串**（docs/re/26 §4.1），這裡沿用同一條路：
// 每個標籤落在自己欄位的文字起點（半形格對齊），生成的字串登記進語系
// 詞表，畫的時候照樣是一條字串。標籤太寬就往右擠開而不是疊字——
// 標題畫在黑底那一列，與資料列互不相干，推開只影響對齊。
func latinTitle(fields []listField, labels []string) string {
	var b strings.Builder
	cells := 0 // 已經填到第幾個半形格
	for i, label := range labels {
		if label == "" || i >= len(fields) {
			continue
		}
		want := (fields[i].X + listTextInset) / textdraw.HalfW
		for cells < want {
			b.WriteByte(' ')
			cells++
		}
		if cells > want {
			b.WriteByte(' ') // 被前一欄推開，至少留一格
			cells++
		}
		b.WriteString(label)
		cells += len(label)
	}
	return b.String()
}

// listField 是一欄：x 是相對視窗內緣的位移，W 是寬度（像素），
// Numeric 為真表示那一欄是半形的數字欄（靠右對齊）。
type listField struct {
	X, W    int
	Numeric bool
	// NoDash 標記「這一欄在空列不印破折號」——軍團表最右端那個
	// 無標題的委任格就是這樣（原版實錄影格上空列只有五組破折號）。
	NoDash bool
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
	if f.Extra > 0 {
		// 無標題欄接在分隔線右邊，中間空一個全形字。
		x += textdraw.GlyphW
		out = append(out, listField{X: x, W: f.Extra * textdraw.GlyphW, NoDash: true})
	}
	return out
}

// listRowY 是第 visible 列的 y（視窗內，標題佔第 0 列）。
func listRowY(visible int) int {
	return listWinY + listRowH + visible*listRowH
}

// listFieldX 是第 col 欄**文字**的左緣：欄起點＋半格縮排
// （docs/spec/38 §1.5 的實機量測）。
func listFieldX(fields []listField, col int) int {
	if col < 0 || col >= len(fields) {
		return listBodyX()
	}
	return listBodyX() + fields[col].X + listTextInset
}

// listFieldRight 是第 col 欄的右緣（數字右靠在分隔線那段的右緣，
// 沒有縮排——武術 128／統率 168／政治 208，docs/spec/38 §1.5）。
func listFieldRight(fields []listField, col int) int {
	if col < 0 || col >= len(fields) {
		return listBodyX()
	}
	return listBodyX() + fields[col].X + fields[col].W
}

// listWarnInk 是換色用的前景色。原版把屬性的低 4 位由 0 改成 A
// （`bh = 0x9A`，docs/re/27 §5），這裡用一個接近的紅。
var listWarnInk = color.RGBA{200, 60, 40, 255}

// corpsHalfStrength 是「總兵數換色」的門檻：原版 `< 0x12C` ＝ 300 點
// ＝ 3,000 人（半編）。
const corpsHalfStrength = 300

// listRankNames 是身分名稱表（docs/re/26 §9，段內 0x75A4）。
// **松崗版用「俘虜」不是日文的「捕虜」。**
var listRankNames = [6]string{"－－－", "軍團長", "內政官", "外交官", "俘虜　", "君主　"}

// listDiplomacyNames 是外交關係六級 ＋ 自己那一格（docs/re/27 §4）。
// 索引就是原版表的筆序；交戰那一筆在畫面上會換色。
var listDiplomacyNames = [7]string{"交戰　", "最惡　", "險惡　", "普通　", "良好　", "親密　", "－－　"}

// listDiplomacyLevel 把交友度換成上表的索引。
//
// 算式在規則層（`diplomacy.DisplayIndex`，出處 `sub_17A7A`）——
// **一條規則只留一份實作**，這裡只負責把 UI 的兩個參數併回原始 byte。
func listDiplomacyLevel(raw int, atWar bool) int {
	if atWar {
		return 0
	}
	return diplomacy.DisplayIndex(raw&0x7F | 0x80)
}

// listDashes 是**沒有資料的那一列**在這一欄要印的東西：分隔線本身。
// 全形欄印 `－`、半形欄印 `-`，數量照欄寬（docs/spec/38 §1.4）。
func listDashes(f listField) string {
	if f.NoDash {
		return ""
	}
	if f.Numeric {
		return strings.Repeat("-", f.W/textdraw.HalfW)
	}
	n := f.W / textdraw.GlyphW
	if n <= 0 {
		return ""
	}
	return strings.Repeat("－", n)
}

// listBlank 是空欄要印的東西（原版印同寬的全形空白）。
func listBlank(f listField) string {
	n := f.W / textdraw.GlyphW
	if n <= 0 {
		return ""
	}
	return strings.Repeat("　", n)
}

// listCellRoom 是一個文字欄畫得下多少像素：到下一欄的起點為止
// （最後一欄到清單本體右緣）。數字欄靠右對齊、起點只會更右，
// 所以拿下一欄的 X 當界線是保守的。
func listCellRoom(fields []listField, col int) int {
	if col < 0 || col >= len(fields) {
		return 0
	}
	limit := listBodyW()
	if col+1 < len(fields) {
		limit = fields[col+1].X
	}
	// 留一格半形的間隙：緊貼著下一欄的數字（`XU-HUANG10`）
	// 看起來就是黏成一個詞。
	return limit - fields[col].X - listTextInset - textdraw.HalfW
}

// fitCell 把一格文字截到欄寬之內。
//
// **為什麼需要**：欄寬是照原版的三個全形字排的（人名欄 6 個半形格），
// 而羅馬化的人名可以到 12 個字母——不裁的話會畫進隔壁的數值欄，
// 兩欄疊在一起比截斷難讀得多（docs/spec/84 §6）。
//
// 裁法對羅馬化人名做了兩個調整，因為半吊子的截斷讀不出來：
// 切在連字號上就把連字號也去掉（`XIAHOU-` → `XIAHOU`）；
// 名字只剩一兩個字母時整段退成姓（`ZHUGE-L` → `ZHUGE`）。
func fitCell(s string, room int) string {
	if room <= 0 || textdraw.StringWidth(s) <= room {
		return s
	}
	out := make([]rune, 0, len(s))
	w := 0
	for _, ch := range s {
		cw := textdraw.RuneWidth(ch)
		if w+cw > room {
			break
		}
		out = append(out, ch)
		w += cw
	}
	cut := strings.TrimSuffix(string(out), "-")
	if i := strings.IndexByte(s, '-'); i > 0 && len(cut) > i && len(cut)-i-1 <= 2 {
		return s[:i] // 名只剩一兩個字母，不如只留姓
	}
	return cut
}
