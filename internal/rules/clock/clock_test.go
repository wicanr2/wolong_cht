package clock

import "testing"

// 每月天數表照抄原版 ds:98ABh。二月固定 28，沒有閏年。
func TestDaysInMonth(t *testing.T) {
	want := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	for m := 0; m <= 12; m++ {
		if got := DaysInMonth(m); got != want[m] {
			t.Errorf("DaysInMonth(%d) = %d, want %d", m, got, want[m])
		}
	}
	// 原版的表只有 0–12，越界就是讀到別人的資料；這裡明確回 0。
	for _, m := range []int{-1, 13, 100} {
		if got := DaysInMonth(m); got != 0 {
			t.Errorf("DaysInMonth(%d) = %d, want 0", m, got)
		}
	}
}

// 四個劇本的起始日必須與 SINARIO.DAT 一致（docs/formats/08 §1.1）。
// 特別注意劇本 2 是 9 月、劇本 4 是 225 年 —— 與日文說明書的截圖不同，
// 以資料為準。
func TestScenarioStartDates(t *testing.T) {
	cases := []struct {
		name             string
		year, month, day int
		daysInStartMonth int
	}{
		{"呂布歸天", 196, 4, 1, 30},
		{"赤壁之戰", 208, 9, 1, 30},
		{"蜀地偏安", 212, 6, 1, 30},
		{"劉禪繼位", 225, 5, 1, 31},
	}
	for _, tc := range cases {
		c := New(tc.year, tc.month, tc.day)
		if c.Hour != 1 || c.Subtick != 0 {
			t.Errorf("%s: 起始時刻 = %d時%d子刻, want 1時0子刻", tc.name, c.Hour, c.Subtick)
		}
		if got := DaysInMonth(c.Month); got != tc.daysInStartMonth {
			t.Errorf("%s: %d月天數 = %d, want %d", tc.name, c.Month, got, tc.daysInStartMonth)
		}
	}
}

// 一個遊戲日 = 24 時 × 9 子刻 = 216 tick，而且整整 216 步之後
// 才會出現一次 Day 進位。
func TestTicksPerDay(t *testing.T) {
	if TicksPerDay != 216 {
		t.Fatalf("TicksPerDay = %d, want 216", TicksPerDay)
	}
	c := New(196, 4, 1)
	days, hours := 0, 0
	for i := 0; i < TicksPerDay; i++ {
		ev := c.Advance()
		if ev.Hour {
			hours++
		}
		if ev.Day {
			days++
		}
	}
	if hours != HoursPerDay {
		t.Errorf("216 tick 內的「時」進位 = %d, want %d", hours, HoursPerDay)
	}
	if days != 1 {
		t.Errorf("216 tick 內的「日」進位 = %d, want 1", days)
	}
	if c.Day != 2 || c.Hour != 1 || c.Subtick != 0 {
		t.Errorf("216 tick 後 = %d日 %d時 %d子刻, want 2日 1時 0子刻", c.Day, c.Hour, c.Subtick)
	}
}

// 進位鏈：跑滿一個月要正好 天數 × 216 個 tick，而且只觸發一次月進位。
func TestMonthRollover(t *testing.T) {
	c := New(196, 4, 1) // 4 月有 30 天
	months := 0
	for i := 0; i < 30*TicksPerDay; i++ {
		if c.Advance().Month {
			months++
		}
	}
	if months != 1 {
		t.Errorf("跑滿 4 月的月進位次數 = %d, want 1", months)
	}
	if c.Year != 196 || c.Month != 5 || c.Day != 1 {
		t.Errorf("= %d年%d月%d日, want 196年5月1日", c.Year, c.Month, c.Day)
	}
}

// 年進位：12 月跑完要跨年，而且月份重設成 1。
func TestYearRollover(t *testing.T) {
	c := New(196, 12, 1) // 12 月有 31 天
	years := 0
	for i := 0; i < 31*TicksPerDay; i++ {
		if c.Advance().Year {
			years++
		}
	}
	if years != 1 {
		t.Errorf("年進位次數 = %d, want 1", years)
	}
	if c.Year != 197 || c.Month != 1 || c.Day != 1 {
		t.Errorf("= %d年%d月%d日, want 197年1月1日", c.Year, c.Month, c.Day)
	}
}

// 年份封頂在 999（原版 cmp 3E8h → 設 3E6h → inc）。
func TestYearCap(t *testing.T) {
	c := New(999, 12, 31)
	for i := 0; i < TicksPerDay; i++ {
		c.Advance()
	}
	if c.Year != MaxYear {
		t.Errorf("999年12月31日跑完一天後 = %d年, want %d年", c.Year, MaxYear)
	}
	if c.Month != 1 || c.Day != 1 {
		t.Errorf("= %d月%d日, want 1月1日", c.Month, c.Day)
	}
}

// 跑滿一整年（非閏年）要正好 365 天。這同時驗了每月天數表。
func TestOneFullYear(t *testing.T) {
	c := New(200, 1, 1)
	days := 0
	for c.Year == 200 {
		if c.Advance().Day {
			days++
		}
	}
	if days != 365 {
		t.Errorf("200 年共 %d 天, want 365（二月 28 天，無閏年）", days)
	}
}

// 季節漸變只在 3/6/9/12 月的 1–16 日、時 == 1 時成立，
// 而且一個轉換月裡正好會成立 16 次。
func TestSeasonTransition(t *testing.T) {
	for _, month := range []int{3, 6, 9, 12} {
		c := New(200, month, 1)
		hits := 0
		if c.InSeasonTransition() {
			hits++ // 起始那一刻本身就符合條件
		}
		steps := 0
		for i := 0; i < DaysInMonth(month)*TicksPerDay; i++ {
			ev := c.Advance()
			if c.InSeasonTransition() {
				hits++
			}
			if ev.SeasonStep {
				steps++
			}
		}
		// Event.SeasonStep 與直接查詢必須一致（起始那一刻不經過 Advance，
		// 所以查詢會多算到它）。
		if steps != hits-1 {
			t.Errorf("%d 月：Event.SeasonStep = %d 次，查詢 = %d 次", month, steps, hits)
		}
		if hits != 16 {
			t.Errorf("%d 月的季節漸變步數 = %d, want 16", month, hits)
		}
	}
	// 非轉換月一次都不該成立。
	for _, month := range []int{1, 2, 4, 5, 7, 8, 10, 11} {
		c := New(200, month, 1)
		for i := 0; i < DaysInMonth(month)*TicksPerDay; i++ {
			if c.InSeasonTransition() {
				t.Fatalf("%d 月不該有季節漸變", month)
			}
			c.Advance()
		}
	}
}

func TestSeason(t *testing.T) {
	cases := []struct {
		month int
		want  Season
	}{
		{1, Winter}, {2, Winter},
		{3, Spring}, {4, Spring}, {5, Spring},
		{6, Summer}, {7, Summer}, {8, Summer},
		{9, Autumn}, {10, Autumn}, {11, Autumn},
		{12, Winter},
	}
	for _, tc := range cases {
		if got := New(200, tc.month, 1).Season(); got != tc.want {
			t.Errorf("%d 月 = %v, want %v", tc.month, got, tc.want)
		}
	}
}

// 進位是階梯式的：發生月進位時一定同時有日進位與時進位。
// 原版的 fall through 寫法保證了這件事，remake 不能破壞它——
// 月結（sub_15358）與每時的世界更新是在同一個 tick 裡依序跑的。
func TestEventCascade(t *testing.T) {
	c := New(196, 12, 31)
	for i := 0; i < TicksPerDay*2; i++ {
		ev := c.Advance()
		if ev.Year && !(ev.Month && ev.Day && ev.Hour) {
			t.Fatalf("年進位時缺了下層進位: %+v", ev)
		}
		if ev.Month && !(ev.Day && ev.Hour) {
			t.Fatalf("月進位時缺了下層進位: %+v", ev)
		}
		if ev.Day && !ev.Hour {
			t.Fatalf("日進位時缺了時進位: %+v", ev)
		}
	}
}

// 任何時刻的欄位都必須落在原版的值域內。
func TestInvariants(t *testing.T) {
	c := New(196, 4, 1)
	for i := 0; i < 400*TicksPerDay; i++ {
		c.Advance()
		if c.Month < 1 || c.Month > 12 {
			t.Fatalf("月 = %d 超出 1–12", c.Month)
		}
		if c.Day < 1 || c.Day > DaysInMonth(c.Month) {
			t.Fatalf("%d 月的日 = %d 超出範圍", c.Month, c.Day)
		}
		if c.Hour < 1 || c.Hour > HoursPerDay {
			t.Fatalf("時 = %d 超出 1–24", c.Hour)
		}
		if c.Subtick < 0 || c.Subtick >= SubticksPerHour {
			t.Fatalf("子刻 = %d 超出 0–8", c.Subtick)
		}
	}
}
