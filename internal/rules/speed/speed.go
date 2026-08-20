// Package speed 是即時制的速度節流。
//
// ⭐ **戰略與戰術各有一個獨立的檔位**（原版 `ds:0CFAh`／`ds:0CFBh`），
// 這是這一款的特性：大地圖與戰場的時間流速分開設定。
//
// 桌面版與手機版共用這一份。**不要在呈現層各寫一份**——
// 兩份會慢慢長出行為差異，而差異出現在時間軸上最難查（CLAUDE.md §7 第 6 條）。
package speed

// 速度節流：把原版「推進一步之後等 N 個計時中斷」搬到 remake 的 60 Hz 畫面。
//
// 原版的中斷是音效驅動發的 291.3 Hz（`YNSOUND.COM` 把 PIT 設成 4660.9 Hz
// 再分頻 16，docs/re/61）。量子 3.43 ms 與 remake 的 16.7 ms 除不盡，
// 用整數「每畫面推進幾步」最好也只能差兩成，所以改成定點累加器——
// 規則與原版相同，只是把「等 N 個中斷」換成「攢滿 N 個中斷單位」。
//
// 規格：docs/spec/34-speed-steps.md。

const (
	// Scale 把 291.3 ÷ 60 ＝ 4.855 變成整數 2913，誤差 0.001%。
	// 不用浮點數是因為存檔與重播要位元一致。
	Scale = 600
	// UnitsPerFrame 是一個畫面更新累加的中斷單位數。
	UnitsPerFrame = 2913

	// TacticalMul 是戰術層的倍率。原版的第 5 列 handler
	// `sub_160A5` 就做這一件事：`ds:0CFCh = ds:0CFBh << 4`。
	TacticalMul = 16

	// Levels 是原版的檔位數（設定表 `ds:5FF4h` 寫著 5）。
	Levels = 5

	// 檔位 0（最高速）原版是「不等待」，速度由機器決定，沒有可抄的數字。
	// 這兩個上限是 remake 差異，取先前 remake 的最快值。
	HighSpeedStrategy = 16
	HighSpeedTactical = 4
)

// Throttle 是一層的累加器。戰略與戰術各一個，互不影響。
type Throttle struct{ acc int }

// Steps 回報這一個畫面更新該推進幾步。level 是原版檔位 0–4，
// mul 是該層的倍率（戰略 1、戰術 tacticalThrottleMul）。
func (t *Throttle) Steps(level, mul, highSpeed int) int {
	if level <= 0 {
		// 最高速：原版不等待，remake 用固定上限。累加器歸零，
		// 免得切回其他檔時一次補跑一大串。
		t.acc = 0
		return highSpeed
	}
	if level >= Levels {
		level = Levels - 1
	}
	cost := level * mul * Scale
	t.acc += UnitsPerFrame
	n := t.acc / cost
	t.acc -= n * cost
	return n
}

// Labels 是五個檔位在畫面上的字，取自原版系統選單（docs/re/55 §4）。
//
// ⚠ 那幾個半形空白是**原版版面的一部分**：系統選單的框寬由字數決定，
// 五列要一樣寬。呈現層要窄版面時自己 trim，不要改這裡。
var Labels = [Levels]string{"最高速", " 高速 ", " 普通 ", " 低速 ", "最低速"}
