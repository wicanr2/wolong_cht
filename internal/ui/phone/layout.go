// Package phone 是手機版的呈現層。
//
// ⚠ **這不是原版版面。** 手機版只共用規則層，畫面與操作重新設計
// （使用者裁定 2026-08-20，規格 docs/mobile/android-ux.md）。
// 桌面版仍是唯一的對拍基準——這裡的任何取捨都不影響那邊。
//
// 這一層認識 Ebiten。**縮放交給 Ebiten 的 Layout 契約**——
// mobile/wolong 的 game.Layout 回傳 LogicalW／LogicalH，之後所有座標
// 都在這個邏輯畫布上算，不必自己換算螢幕像素。
package phone

// 版面常數。單位是**邏輯像素**。
//
// 手機的邏輯畫布不是原版的 640×400：原版那個比例在現代手機上會留很寬的
// 黑邊，而且熱區只有 16×16。這裡用 16:9 的 960×540 當基準，
// Ebiten 再等比縮放到實際螢幕（docs/mobile/android-ux.md §2）。
const (
	LogicalW = 960
	LogicalH = 540

	// StatusH／CommandH 是上下兩條的高度。48 dp 的觸控下限換算到
	// 這個邏輯畫布大約是 56／64——**兩條都要比下限寬鬆**，
	// 因為手指按的是按鈕中心不是邊緣。
	StatusH  = 56
	CommandH = 64

	// CardW／CardH 是點到據點或軍團時浮出的小卡。
	CardW = 300
	CardH = 176
	// CardMargin 是小卡離畫面邊緣的距離。
	CardMargin = 16
)

// MapRect 回傳地圖區的邏輯矩形。
func MapRect() (x, y, w, h int) {
	return 0, StatusH, LogicalW, LogicalH - StatusH - CommandH
}

// Command 是底部指令列的四個入口（docs/mobile/android-ux.md §4）。
type Command int

const (
	CmdAdvise Command = iota // 進言
	CmdList                  // 一覽
	CmdCorps                 // 軍團
	CmdSystem                // 系統
	numCommands
)

func (c Command) Label() string {
	return [...]string{"進言", "一覽", "軍團", "系統"}[c]
}

// CommandRect 回傳第 i 個指令鈕的邏輯矩形。四個等分，各自留 8 px 間隙。
func CommandRect(i int) (x, y, w, h int) {
	const gap = 8
	cell := LogicalW / int(numCommands)
	return i*cell + gap, LogicalH - CommandH + gap,
		cell - gap*2, CommandH - gap*2
}

// Zoom 是地圖的縮放級距。**只允許整數倍**——原版是點陣圖，
// 非整數倍會把 16×16 的圖塊糊掉（docs/mobile/android-ux.md §3）。
const (
	MinZoom = 1
	MaxZoom = 3
)

// ClampZoom 把縮放夾回合法級距。
func ClampZoom(z int) int {
	if z < MinZoom {
		return MinZoom
	}
	if z > MaxZoom {
		return MaxZoom
	}
	return z
}

// TilePx 是大地圖一格的像素邊長（原版 16×16，`internal/assets/world`）。
const TilePx = 16
