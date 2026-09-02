package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
)

// 一般部隊的戰力：((統率 + 適性) × 3 + 兵種係數) ÷ 4。
// 係數 30／4／12 出自 seg000:9C0F 的 bytes 1E 04 0C（docs/re/78 §3）。
func TestSoldierPowerPerTroopType(t *testing.T) {
	cases := []struct {
		command, aptitude int
		kind              army.TroopType
		want              int
	}{
		// 統率 10、適性 5 → (15 × 3) ＝ 45
		{10, 5, army.Cavalry, (45 + 30) / 4},  // 18
		{10, 5, army.Archer, (45 + 4) / 4},    // 12
		{10, 5, army.Infantry, (45 + 12) / 4}, // 14
		// 兩端
		{0, 0, army.Cavalry, 30 / 4},
		{15, 10, army.Infantry, (75 + 12) / 4},
	}
	for _, c := range cases {
		if got := soldierPower(c.command, c.aptitude, c.kind); got != c.want {
			t.Errorf("soldierPower(%d,%d,%v) = %d，want %d",
				c.command, c.aptitude, c.kind, got, c.want)
		}
	}
}

// ⭐ 統率是唯一會動的變數：兵種與適性固定時，戰力必須隨統率單調上升。
// 這一條擋的是「又把它接回士氣」。
func TestSoldierPowerRisesWithCommand(t *testing.T) {
	prev := -1
	for cmd := 0; cmd <= 15; cmd++ {
		p := soldierPower(cmd, 5, army.Infantry)
		if p < prev {
			t.Fatalf("統率 %d 的戰力 %d 比前一級 %d 還低", cmd, p, prev)
		}
		prev = p
	}
	if soldierPower(15, 5, army.Infantry) <= soldierPower(1, 5, army.Infantry) {
		t.Fatal("統率 15 與統率 1 的戰力一樣，這個欄位沒有接進來")
	}
}

// 大將那一格：戰力 (武力 × 2 + 適性) × 2、體力 max(70, (武力 × 4 + 50) × 士氣 ÷ 100)。
func TestLeaderPowerAndHP(t *testing.T) {
	if got, want := leaderPower(13, 4), (13*2+4)*2; got != want {
		t.Errorf("leaderPower = %d，want %d", got, want)
	}
	// 士氣 200、武力 13 → (52 + 50) × 2 ＝ 204
	if got, want := leaderHP(13, 200), 204; got != want {
		t.Errorf("leaderHP(13,200) = %d，want %d", got, want)
	}
	// 下限 0x46 ＝ 70：士氣 30、武力 1 → 54 × 30 ÷ 100 ＝ 16 → 夾成 70
	if got := leaderHP(1, 30); got != leaderHPFloor {
		t.Errorf("leaderHP(1,30) = %d，want 下限 %d", got, leaderHPFloor)
	}
}

// ⚠ 大將的開場體力**不是**軍團士氣（docs/spec/61 §2 的例外）。
func TestLeaderHPIsNotMorale(t *testing.T) {
	const morale = 200
	if leaderHP(1, morale) == morale {
		t.Fatal("武力 1 的大將體力等於士氣，sub_19B40 那一格沒有蓋掉")
	}
}

// 戰場類別選哪一個適性欄：攻城恆 0（原版寫進去的是據點編號 0–191），
// 野戰看圖塊（docs/re/78 §2.1）。
func TestAptitudeIndexByBattleClass(t *testing.T) {
	if got := aptitudeIndex(true, 0xFF); got != siegeAptitude {
		t.Errorf("攻城應該恆取攻城適性，得到 %d", got)
	}
	cases := []struct{ tile, want int }{
		{0x00, siegeAptitude},
		{0xBF, siegeAptitude},
		{0xC0, fieldAptitude},
		{0xD0, fieldAptitude},
		{0xD1, waterAptitude},
		{0xFF, waterAptitude},
	}
	for _, c := range cases {
		if got := aptitudeIndex(false, c.tile); got != c.want {
			t.Errorf("野戰圖塊 %#02x → %d，want %d", c.tile, got, c.want)
		}
	}
}
