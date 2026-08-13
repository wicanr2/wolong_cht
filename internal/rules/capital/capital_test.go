package capital

import "testing"

// 類型值照 docs/formats/08 §1.6：0 大城／1 中城／2 小城／3 關／4 戰場。
const (
	big = iota
	mid
	small
	pass
	field
)

// ⚠ **「類型優先」只在生產力也不落後時成立。**
//
// 兩個條件是同時檢查的，門檻又邊掃邊更新，所以一個生產力很高的小城
// 會把後面所有大城擋掉——即使大城才是「應該」的首都。
// 這兩個案例只差在順序，結果完全相反。
func TestPickKindWinsOnlyWhenProductionAlsoWins(t *testing.T) {
	// 大城排前面：先被記下，後面的小城過不了類型門檻。
	front := []Site{
		{Owner: 1, Kind: big, Production: 100},
		{Owner: 1, Kind: small, Production: 9000},
	}
	if got := Pick(front, 1); got != 0 {
		t.Fatalf("大城在前應選 #0，得到 %d", got)
	}
	// 小城排前面：它的 9000 變成門檻，後面的大城因為生產力太低被擋掉。
	back := []Site{
		{Owner: 1, Kind: small, Production: 9000},
		{Owner: 1, Kind: big, Production: 100},
	}
	if got := Pick(back, 1); got != 0 {
		t.Fatalf("小城在前應仍選 #0（原版就是這樣），得到 %d", got)
	}
}

func TestPickIgnoresOtherFactions(t *testing.T) {
	sites := []Site{
		{Owner: 2, Kind: big, Production: 40000},
		{Owner: 1, Kind: small, Production: 500},
	}
	if got := Pick(sites, 1); got != 1 {
		t.Fatalf("別人的大城不該被選，得到 %d", got)
	}
}

func TestPickNoCities(t *testing.T) {
	if got := Pick([]Site{{Owner: 2, Kind: big}}, 1); got != None {
		t.Fatalf("沒有據點應回 None，得到 %d", got)
	}
}

// 同類型時比生產力。這正是劇本 1 勢力 19 的形狀：
// 三個都是小城，演算法選生產力最高的「黃」而不是作者填的「北海」。
func TestPickBreaksTieByProduction(t *testing.T) {
	sites := []Site{
		{Owner: 1, Kind: small, Production: 2857}, // 黃
		{Owner: 1, Kind: small, Production: 1562}, // 挺
		{Owner: 1, Kind: small, Production: 2562}, // 北海
	}
	if got := Pick(sites, 1); got != 0 {
		t.Fatalf("應選生產力最高的 #0，得到 %d", got)
	}
}

// ⭐ 門檻邊掃邊更新，所以**兩個條件是「且」不是排序**。
//
// 這個案例分得出兩種實作：先按類型排序再取生產力最高的話會選 #2；
// 照抄原版的線性掃描會選 #0——因為掃到 #2 時 bestProd 已經是 9000，
// 而 #2 的 5000 沒過「生產力 ≥ 目前最佳」那道門檻。
//
// 釘住它是為了防止哪天有人把這段「整理」成排序。
func TestPickIsLinearScanNotSort(t *testing.T) {
	sites := []Site{
		{Owner: 1, Kind: mid, Production: 9000},
		{Owner: 1, Kind: small, Production: 1000},
		{Owner: 1, Kind: big, Production: 5000},
	}
	if got := Pick(sites, 1); got != 0 {
		t.Fatalf("線性掃描應停在 #0，得到 %d（被改成排序了？）", got)
	}
}

// 戰場（類型 4）城兵上限是 0，但演算法沒有排除它——
// 劇本 2 的勢力 2 首都就填成長阪。只有在沒有更好的選擇時才會輪到。
func TestPickAllowsBattlefieldWhenNothingElse(t *testing.T) {
	sites := []Site{{Owner: 1, Kind: field, Production: 1100}}
	if got := Pick(sites, 1); got != 0 {
		t.Fatalf("只有戰場時也要選得出來，得到 %d", got)
	}
}

func TestAcceptRelocation(t *testing.T) {
	// 原版 sub_16909：先比類型（編號小＝城大），再比生產力，兩個都要過。
	cases := []struct {
		name     string
		from, to Site
		want     bool
	}{
		{"同級但生產力更高 → 接受", Site{Kind: 1, Production: 100}, Site{Kind: 1, Production: 200}, true},
		{"目標是更大的城且生產力更高 → 接受", Site{Kind: 2, Production: 100}, Site{Kind: 0, Production: 150}, true},
		{"生產力更高但城更小 → 拒絕", Site{Kind: 0, Production: 100}, Site{Kind: 2, Production: 9999}, false},
		{"城更大但生產力較低 → 拒絕", Site{Kind: 2, Production: 300}, Site{Kind: 0, Production: 200}, false},
		{"生產力相等 → 拒絕（嚴格大於）", Site{Kind: 1, Production: 100}, Site{Kind: 1, Production: 100}, false},
		{"高 4 位不參與比較", Site{Kind: 0xF1, Production: 100}, Site{Kind: 0xA1, Production: 200}, true},
	}
	for _, c := range cases {
		if got := AcceptRelocation(c.from, c.to); got != c.want {
			t.Errorf("%s：AcceptRelocation(%+v, %+v) = %v，想要 %v", c.name, c.from, c.to, got, c.want)
		}
	}
}
