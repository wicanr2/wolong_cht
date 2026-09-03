package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
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

// 戰場類別選哪一個適性欄。⛔ **參數是戰場編號，不是大地圖圖塊**——
// `sub_19A33` 比的是 `byte_10D34`，而野戰的編號由 `sub_14B63` 算
// （docs/re/05、docs/re/58 §4）。
func TestAptitudeIndexByBattleClass(t *testing.T) {
	cases := []struct {
		field, want int
		note        string
	}{
		{0x00, siegeAptitude, "據點 0"},
		{0xBF, siegeAptitude, "據點 191（攻城的上界）"},
		{0xC0, fieldAptitude, "平原配對表的基底"},
		{0xD0, fieldAptitude, "地形類型 2"},
		{0xD1, waterAptitude, "地形類型 3"},
		{0xD5, waterAptitude, "地形類型 7"},
	}
	for _, c := range cases {
		if got := aptitudeIndex(c.field); got != c.want {
			t.Errorf("戰場編號 %#02x（%s）→ %d，want %d",
				c.field, c.note, got, c.want)
		}
	}
	// ⭐ 這一格是先前接錯的證據：橋的**圖塊**是 `0xCA`（< 0xD1，看起來像陸上），
	// 而它的**戰場編號**是 0xD1–0xD4 ＝ 海上。拿圖塊去比會取到錯的適性欄。
	if aptitudeIndex(0xCA) == waterAptitude {
		t.Error("0xCA 當成戰場編號應該是陸上——這一格在測 raw 圖塊就錯了")
	}
	if aptitudeIndex(0xD3) != waterAptitude {
		t.Error("橋的戰場編號 0xD3 應該取海戰適性")
	}
}

// ⭐ 接線本身：`squadPowers` 要用 `TacticalSetup.Category` 給的類別，
// 三個適性欄各給不同的值，看它取到哪一個。
//
// 少了這一支，`Category` 沒被設（先前 `Tile` 就是**從來沒有人設**）
// 也不會有任何測試變紅——野戰永遠取 `+0x0F`，海上那一欄一輩子用不到。
func TestSquadPowersUsesBattleCategory(t *testing.T) {
	for _, tc := range []struct {
		category, wantApt int
	}{
		{siegeAptitude, 2},
		{fieldAptitude, 5},
		{waterAptitude, 9},
	} {
		w := &World{}
		w.Corps[0] = Corps{Alive: true, Morale: 100}
		w.Generals[0] = General{Alive: true, Command: 10, Martial: 10}
		w.Generals[0].Aptitude = [3]int{2, 5, 9}
		w.Corps[0].Units[0] = combat.Unit{Men: 100, Kind: army.Cavalry}
		squads, _, _ := w.squadPowers(0, tc.category)
		want := soldierPower(10, tc.wantApt, army.Cavalry)
		if squads[0] != want {
			t.Errorf("類別 %d：戰力 %d，want %d（適性 %d）",
				tc.category, squads[0], want, tc.wantApt)
		}
	}
}
