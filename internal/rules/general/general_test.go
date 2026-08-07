package general

import "testing"

// 劇本 1 的實際數值，用來釘住公式。
// 這些值是從 SINARIO.DAT 讀出來的（適性欄位已經 >>4）。
var scenario1 = []General{
	{Name: "呂布", Aptitude: [3]int{4, 10, 0}, Martial: 15, Command: 11, Politics: 1},
	{Name: "趙雲", Aptitude: [3]int{5, 6, 3}, Martial: 13, Command: 13, Politics: 6},
	{Name: "諸葛亮", Aptitude: [3]int{10, 8, 6}, Martial: 5, Command: 15, Politics: 15},
	{Name: "關羽", Aptitude: [3]int{3, 6, 2}, Martial: 14, Command: 11, Politics: 8},
	{Name: "孫乾", Aptitude: [3]int{2, 0, 0}, Martial: 1, Command: 1, Politics: 13},
	{Name: "滿寵", Aptitude: [3]int{3, 0, 0}, Martial: 1, Command: 1, Politics: 14},
}

func TestRatingKnownValues(t *testing.T) {
	want := map[string]int{
		"呂布": 66, "趙雲": 66, "諸葛亮": 64, "關羽": 61,
		"孫乾": 6, "滿寵": 7,
	}
	for _, g := range scenario1 {
		if got := g.Rating(); got != want[g.Name] {
			t.Errorf("%s 評價 = %d, want %d", g.Name, got, want[g.Name])
		}
	}
}

// ⭐ 武術與統率的權重相同 —— 說明書 10.5：
// 「武術と統率はその合計が同じであれば強さは同じ」。
func TestMartialAndCommandHaveEqualWeight(t *testing.T) {
	a := General{Martial: 15, Command: 5}
	b := General{Martial: 5, Command: 15}
	c := General{Martial: 10, Command: 10}
	if a.Rating() != b.Rating() || b.Rating() != c.Rating() {
		t.Errorf("和相同評價應相同：%d / %d / %d", a.Rating(), b.Rating(), c.Rating())
	}
	// 和不同就該不同。
	if (General{Martial: 10, Command: 11}).Rating() <= c.Rating() {
		t.Error("和較大的評價應該較高")
	}
}

// 政治完全不計入評價。
func TestPoliticsExcluded(t *testing.T) {
	lo := General{Martial: 5, Command: 5, Politics: 1}
	hi := General{Martial: 5, Command: 5, Politics: 15}
	if lo.Rating() != hi.Rating() {
		t.Errorf("政治不該影響評價：%d vs %d", lo.Rating(), hi.Rating())
	}
}

// 純文官應該排在最後 —— 這是「評價是純軍事價值」的可觀察後果。
func TestCivilOfficialsRankLast(t *testing.T) {
	var lowest General
	lowest.Aptitude = [3]int{99, 99, 99} // 先給一個不可能的高值
	lowestVal := 1 << 30
	for _, g := range scenario1 {
		if g.Rating() < lowestVal {
			lowestVal, lowest = g.Rating(), g
		}
	}
	if lowest.Name != "孫乾" {
		t.Errorf("評價最低的是 %s，預期是純文官孫乾", lowest.Name)
	}
	if lowest.Politics < 10 {
		t.Errorf("評價最低者的政治 = %d，應該是高政治的文官", lowest.Politics)
	}
}

// 評價相同時，玩家指揮選武術高的、委任選統率高的（說明書 10.5）。
func TestTieBreakByCommandStyle(t *testing.T) {
	warrior := General{Name: "武", Martial: 15, Command: 5}
	tactician := General{Name: "統", Martial: 5, Command: 15}
	if warrior.Rating() != tactician.Rating() {
		t.Fatal("前提錯了：兩人評價應該相同")
	}
	if !PreferForPlayerCommand(warrior, tactician) {
		t.Error("玩家親自指揮應該優先選武術高的")
	}
	if !PreferForDelegation(tactician, warrior) {
		t.Error("委任應該優先選統率高的")
	}
}
