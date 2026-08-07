// Package clock 實作原版的遊戲時鐘。
//
// 規格來源是反組譯，不是猜的：`KI.EXE` 的 `sub_11D8E`（進位鏈）、
// `sub_11E17`（日期顯示）、`ds:98ABh`（每月天數表）。
// 完整反組譯筆記見 docs/re/06，機制說明見 docs/mechanics/15-realtime.md。
//
// 原版把時間切成五層，而說明書只提了其中三層：
//
//	子刻 (0–8，9 階) → 時 (1–24，24 階) → 日 → 月 → 年
//
// 一個遊戲日 = 24 × 9 = 216 個 tick。
//
// 「時」這一層在說明書裡從沒出現過，但它才是世界更新的實際節拍——
// 季節漸變、每時的世界更新、日期重繪都掛在它上面。remake 若照說明書
// 只做到「日」，這些事件就沒有地方掛。
package clock

// 每月天數表。原版在 `ds:98ABh`，索引 1–12（索引 0 是填充用的 0）。
//
// ⚠ 二月固定 28 天，**原版沒有閏年判斷**。這不是簡化，是照抄——
// 遊戲的年份跑到 999 就封頂，四年一閏在這裡沒有意義，
// 而加上閏年會讓月結的日期與原版錯開。
var daysInMonth = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// DaysInMonth 回傳該月的天數。月份超出 1–12 時回 0，與原版的表一致。
func DaysInMonth(month int) int {
	if month < 1 || month > 12 {
		return 0
	}
	return daysInMonth[month]
}

const (
	// SubticksPerHour 是「時」進位前的子刻數。
	// 原版：`cmp byte ptr ds:0CF2h, 8` / `jb` → 子刻 0..8 共 9 階。
	SubticksPerHour = 9

	// HoursPerDay 是「日」進位前的時數。
	// 原版：`cmp byte ptr ds:0CF3h, 17h` / `jb` → 時 1..24 共 24 階。
	HoursPerDay = 24

	// TicksPerDay 是一個遊戲日的 tick 數。
	TicksPerDay = HoursPerDay * SubticksPerHour // 216

	// MaxYear 是年份的上限。原版 `cmp word ptr ds:0CF6h, 3E8h`
	// 超過就設回 `3E6h` 再 `inc`，實際效果是停在 999。
	MaxYear = 999
)

// Clock 是遊戲時鐘的狀態。
//
// 欄位順序與原版劇本／存檔區塊最前面 8 byte 一一對應
// （`docs/formats/08` §1.1），這樣存讀檔可以直接對映：
//
//	+0x00 u8  Day
//	+0x01 u8  （該月天數，本結構用 DaysInMonth(Month) 算，不另存）
//	+0x02 u8  Subtick
//	+0x03 u8  Hour
//	+0x04 u16 Month
//	+0x06 u16 Year
type Clock struct {
	Year    int
	Month   int
	Day     int
	Hour    int // 1–24
	Subtick int // 0–8
}

// Event 是一次 Advance 造成的進位事件。
//
// 呼叫端靠這個決定要跑哪些系統，而不是自己比較前後狀態——
// 原版就是在進位鏈的各層直接呼叫對應的程式
// （換月 → 月結、每時 → 世界更新與季節漸變）。
type Event struct {
	Hour  bool // 進位到新的「時」
	Day   bool // 進位到新的一天
	Month bool // 進位到新的月 → 月結
	Year  bool // 進位到新的年

	// SeasonStep 表示這一 tick 要走一步季節調色盤內插。
	//
	// 原版 sub_19377 掛在「時」進位上，但開頭就 `cmp ds:0CF3h, 1`——
	// 只有時 == 1 才做事，所以實際頻率是**一天一次**。
	// 放進 Event 是為了讓呼叫端不必知道這個前提；
	// 直接每 tick 去問 InSeasonTransition() 會得到一天 9 次的錯誤答案。
	SeasonStep bool
}

// Any 回報這次 Advance 有沒有任何進位發生。
func (e Event) Any() bool { return e.Hour || e.Day || e.Month || e.Year }

// New 建一個指向指定日期的時鐘，時與子刻設成原版的初值。
//
// 為什麼「時」的初值是 1 而不是 0：原版的進位鏈是層層 fall through，
// 進位之後一定會再跑一次最底層的 `inc`，所以任何新的一天都是從「時 = 1」
// 開始的。四個劇本的存檔在該欄位存的就是 1（docs/re/06 §2）。
func New(year, month, day int) Clock {
	return Clock{Year: year, Month: month, Day: day, Hour: 1, Subtick: 0}
}

// Advance 推進一個 tick，回報這一步發生了哪些進位。
//
// 這是 `sub_11D8E` 的直譯。原版用「條件不成立就往下落，落到底再一路
// fall through 回來」的寫法，等價於下面這種「由外而內找出最高進位層，
// 再由內而外重設」的結構。兩者的可觀察行為相同，而這種寫法在 Go 裡
// 才讀得懂。
func (c *Clock) Advance() Event {
	// 還沒到子刻上限：只加子刻，什麼都不觸發。
	if c.Subtick < SubticksPerHour-1 {
		c.Subtick++
		return Event{}
	}

	var ev Event

	switch {
	case c.Hour < HoursPerDay:
		// 時還沒滿，只進位到「時」。

	default:
		// 時滿了 → 至少要進位到「日」。
		ev.Day = true
		if c.Day < DaysInMonth(c.Month) {
			break
		}
		// 日也滿了 → 進位到「月」。
		ev.Month = true
		if c.Month < 12 {
			break
		}
		// 月也滿了 → 進位到「年」。
		ev.Year = true
	}

	if ev.Year {
		if c.Year >= MaxYear {
			// 原版：設 3E6h 後 inc → 停在 999。
			c.Year = MaxYear - 1
		}
		c.Year++
		c.Month = 0
	}
	if ev.Month {
		c.Month++
		c.Day = 0
	}
	if ev.Day {
		c.Day++
		c.Hour = 0
	}
	ev.Hour = true
	c.Hour++
	c.Subtick = 0
	ev.SeasonStep = c.InSeasonTransition()
	return ev
}

// InSeasonTransition 回報「此刻」是不是季節漸變的一步。
//
// 原版 `sub_19377`：只在 3、6、9、12 月的 1–16 日、且「時 == 1」
// 時走一步調色盤內插。**季節不是某天突然換的，是花 16 天褪過去的。**
//
// 條件裡的 Subtick == 0 是必要的：原版是在「時」進位的當下呼叫，
// 而一個「時」會持續 9 個子刻。少了這一項，每天會回報 9 次而不是 1 次。
// 一般情況用 Advance 回傳的 Event.SeasonStep 就好。
func (c Clock) InSeasonTransition() bool {
	if c.Hour != 1 || c.Subtick != 0 {
		return false
	}
	switch c.Month {
	case 3, 6, 9, 12:
		return c.Day <= 16
	}
	return false
}

// Season 是四季之一。
type Season int

const (
	Spring Season = iota
	Summer
	Autumn
	Winter
)

// Season 回傳目前的季節。
//
// 轉換月是 3／6／9／12（`sub_19377` 的四個分支），
// 而轉換是漸進的，所以「目前季節」指的是**已經轉入**的那一季：
// 3 月起是春，6 月起是夏，9 月起是秋，12 月起是冬。
// 1–2 月仍屬前一年的冬。
func (c Clock) Season() Season {
	switch {
	case c.Month >= 12 || c.Month < 3:
		return Winter
	case c.Month < 6:
		return Spring
	case c.Month < 9:
		return Summer
	default:
		return Autumn
	}
}

func (s Season) String() string {
	switch s {
	case Spring:
		return "春"
	case Summer:
		return "夏"
	case Autumn:
		return "秋"
	case Winter:
		return "冬"
	}
	return "?"
}
