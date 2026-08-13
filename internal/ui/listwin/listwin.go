// Package listwin 實作原版的「一覽表」。
//
// 規格出自日文原版說明書 3.8 節，那一節把操作講得很精確：
//
//	画面上部の項目名を左クリックする事でソートを行います。
//	データの選択は左クリックで選択部分がリバースし、
//	その状態でもう一度左クリックすると決定します。
//	リバース状態からキャンセルすると一覧表の選択に戻ります。
//	尚、ソート状態は武将一覧、拠点一覧など同種のウインドウ単位で
//	記録されていますので毎回ソートする必要はありません。
//
// 四條規則，每一條都影響手感，remake 不能省：
//
//  1. 點欄位名排序
//  2. **兩段式選取**：第一次點選取（反白），第二次點才決定
//  3. 右鍵取消；反白狀態下取消只退回選取層，不關視窗
//  4. **排序狀態以視窗種類為單位記住**，不是每次重來
//
// 這一層不含畫面——它只管狀態機與排序，這樣可以用測試把規則釘住，
// 畫面怎麼畫是 cmd/wlgame 的事。
package listwin

import "sort"

// Column 是一個可排序的欄位。
type Column struct {
	Title string
	// Less 比較兩筆資料。回傳 true 表示 a 應該排在 b 前面。
	Less func(a, b int) bool
}

// Kind 是視窗種類。排序狀態以它為單位記住（說明書 3.8）。
type Kind int

const (
	Generals Kind = iota // 武將一覽
	Cities               // 據點一覽
	Factions             // 勢力一覽
	Corps                // 軍團一覽
	NumKinds
)

// Phase 是兩段式選取的狀態。
type Phase int

const (
	// Browsing 是還沒選任何一列。
	Browsing Phase = iota
	// Selected 是有一列反白了，再點一次才決定。
	Selected
)

// sortState 是一種視窗記住的排序方式。
type sortState struct {
	col  int
	desc bool
	set  bool
}

// Memory 記住各種視窗的排序狀態，跨開關保留。
//
// 說明書：「ソート状態は…同種のウインドウ単位で記録されていますので
// 毎回ソートする必要はありません」。
type Memory [NumKinds]sortState

// List 是一個開著的一覽表。
type List struct {
	Kind    Kind
	Columns []Column
	// Rows 是資料列的索引（例如武將編號）。排序會重排這個切片。
	Rows []int

	// Cursor 是目前停在第幾列（Rows 的索引）。
	Cursor int
	// Top 是畫面最上面顯示第幾列，用於捲動。
	Top int
	// Height 是畫面一次顯示幾列。
	Height int

	phase Phase
	mem   *Memory
}

// New 開一個一覽表，並套用該種視窗記住的排序。
func New(kind Kind, cols []Column, rows []int, height int, mem *Memory) *List {
	l := &List{Kind: kind, Columns: cols, Rows: rows, Height: height, mem: mem}
	if mem != nil {
		if st := mem[kind]; st.set && st.col < len(cols) {
			l.applySort(st.col, st.desc)
		}
	}
	return l
}

// Phase 回報目前在瀏覽還是已選取。
func (l *List) Phase() Phase { return l.phase }

// Selection 回傳目前游標所指的資料索引。清單是空的時回 -1。
func (l *List) Selection() int {
	if l.Cursor < 0 || l.Cursor >= len(l.Rows) {
		return -1
	}
	return l.Rows[l.Cursor]
}

// SortBy 依第 col 欄排序。**再點同一欄會反向**——原版說明書沒寫這一點，
// 但「點欄位名排序」的介面若不能反向，使用者只能排出一半的順序，
// 所以這裡補上並標成 remake 差異。
func (l *List) SortBy(col int) {
	if col < 0 || col >= len(l.Columns) {
		return
	}
	desc := false
	if l.mem != nil {
		if st := l.mem[l.Kind]; st.set && st.col == col {
			desc = !st.desc
		}
	}
	l.applySort(col, desc)
}

func (l *List) applySort(col int, desc bool) {
	less := l.Columns[col].Less
	if less == nil {
		return
	}
	// 用穩定排序：同鍵值的列要保持原本的相對次序，
	// 否則連點兩次不同欄位會讓畫面跳來跳去。
	sort.SliceStable(l.Rows, func(i, j int) bool {
		if desc {
			return less(l.Rows[j], l.Rows[i])
		}
		return less(l.Rows[i], l.Rows[j])
	})
	if l.mem != nil {
		l.mem[l.Kind] = sortState{col: col, desc: desc, set: true}
	}
	// 排序之後游標回到最上面，否則會停在一個語意不同的位置。
	l.Cursor, l.Top = 0, 0
	l.phase = Browsing
}

// Move 移動游標，並跟著捲動。
func (l *List) Move(d int) {
	if len(l.Rows) == 0 {
		return
	}
	// 選取狀態下移動游標會取消選取——原版是兩段式，
	// 移動之後那個「已選取的那一列」就不再是游標所指的列了。
	if l.phase == Selected {
		l.phase = Browsing
	}
	l.Cursor += d
	if l.Cursor < 0 {
		l.Cursor = 0
	}
	if l.Cursor >= len(l.Rows) {
		l.Cursor = len(l.Rows) - 1
	}
	l.scrollIntoView()
}

// Select 把游標移到指定的 Rows 索引，供滑鼠／觸控列點擊使用。
// 點到不同列時，原本的兩段式反白必須先解除；點到同一列則保留
// Selected，讓呼叫端可用同一個 Confirm 完成第二次點擊確認。
func (l *List) Select(row int) bool {
	if row < 0 || row >= len(l.Rows) {
		return false
	}
	if l.phase == Selected && l.Cursor != row {
		l.phase = Browsing
	}
	l.Cursor = row
	l.scrollIntoView()
	return true
}

// Page 以目前可見列數翻頁。頁切換沿用 Move 的邊界與捲動規則，
// 並清除兩段式反白，避免翻頁後仍把不可見列當成待確認項目。
func (l *List) Page(delta int) {
	if len(l.Rows) == 0 || delta == 0 {
		return
	}
	step := l.Height
	if step <= 0 {
		step = 1
	}
	l.phase = Browsing
	l.Cursor += delta * step
	if l.Cursor < 0 {
		l.Cursor = 0
	}
	if l.Cursor >= len(l.Rows) {
		l.Cursor = len(l.Rows) - 1
	}
	l.scrollIntoView()
}

func (l *List) scrollIntoView() {
	if l.Height <= 0 {
		return
	}
	if l.Cursor < l.Top {
		l.Top = l.Cursor
	}
	if l.Cursor >= l.Top+l.Height {
		l.Top = l.Cursor - l.Height + 1
	}
	if max := len(l.Rows) - l.Height; l.Top > max {
		l.Top = max
	}
	if l.Top < 0 {
		l.Top = 0
	}
}

// Confirm 是「左鍵」。第一次選取（反白），第二次才決定。
//
// 回傳 (資料索引, true) 表示決定了；否則表示只是進入選取狀態。
func (l *List) Confirm() (int, bool) {
	if len(l.Rows) == 0 {
		return -1, false
	}
	if l.phase == Browsing {
		l.phase = Selected
		return -1, false
	}
	l.phase = Browsing
	return l.Selection(), true
}

// Cancel 是「右鍵」。
//
// 說明書：「リバース状態からキャンセルすると一覧表の選択に戻ります」——
// **反白狀態下取消只退回選取層，不關視窗**。回傳 true 表示應該關掉視窗。
func (l *List) Cancel() bool {
	if l.phase == Selected {
		l.phase = Browsing
		return false
	}
	return true
}

// Visible 回傳這一頁要畫的列（資料索引）與它們在 Rows 裡的位置。
func (l *List) Visible() (rows []int, first int) {
	if l.Height <= 0 || len(l.Rows) == 0 {
		return nil, 0
	}
	end := l.Top + l.Height
	if end > len(l.Rows) {
		end = len(l.Rows)
	}
	return l.Rows[l.Top:end], l.Top
}
