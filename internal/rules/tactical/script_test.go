package tactical

import (
	"os"
	"testing"
)

// runQuery 跑一個查詢指令再跑一個「等於 want 就下退卻」的分支，
// 回報分支有沒有成立——用得到的行為去檢查 cond，不用暴露內部欄位。
func runQuery(b *Battle, side int, op, arg, want int) bool {
	code := make([]byte, ScriptCodeSize)
	// 0: <op> <arg>
	// 2: 2a <want>  branch == want → 目標
	// 4: 00 04      目標 ＝ byte 8
	// 6: 00 ff      不成立就落到這裡等 255 幀（**不能與目標同一個位址**，
	//               否則成立與否都走到同一行，測不出東西）
	// 8: e3 05      order 全軍 退卻
	copy(code, []byte{byte(op), byte(arg), 0x2a, byte(want),
		0x00, 0x04, 0x00, 0xff, 0xe3, 0x05})
	s := NewScript(code, side)
	s.Step(b) // 查詢
	s.Step(b) // 分支
	s.Step(b) // 目標處
	return b.Sides[side].Soldiers[1].Next == Retreat
}

// 指令 4 查的是**敵方**選的陣形編號（`byte_1D346` 是玩家在選單上點的那一格，
// 而腳本永遠跑在對面那一側）。
func TestScriptQueryFoeFormation(t *testing.T) {
	b := newTestBattle(flatField())
	b.Sides[0].Formation = 11
	if !runQuery(b, 1, opQFoeForm, 0, 11) {
		t.Error("指令 4 沒回敵方的陣形編號")
	}
	b2 := newTestBattle(flatField())
	b2.Sides[0].Formation = 11
	if runQuery(b2, 1, opQFoeForm, 0, 3) {
		t.Error("指令 4 回了不該相等的值")
	}
}

// 指令 5 把敵方的陣形原點跟 28 比大小，回 0／1／2。
func TestScriptQueryFoeLine(t *testing.T) {
	cases := []struct {
		line, want int
		name       string
	}{
		{5, 0, "敵陣貼在自己那一側"},
		{28, 1, "敵陣在正中央"},
		{48, 2, "敵陣壓過來了"},
	}
	for _, c := range cases {
		b := newTestBattle(flatField())
		b.Sides[0].Line = c.line
		if !runQuery(b, 1, opQFoeLine, 0, c.want) {
			t.Errorf("%s（原點 %d）應回 %d", c.name, c.line, c.want)
		}
	}
}

// 指令 11／12 是同一個欄位的兩側：11 問自己、12 問對方，都夾在 255。
func TestScriptQueryMen(t *testing.T) {
	b := newTestBattle(flatField())
	// newTestBattle 兩側都是 6 × 100 ＝ 600，夾成 255。
	if !runQuery(b, 1, opQMyMen, 0, 255) {
		t.Error("指令 11 沒把自軍兵力夾成 255")
	}

	b2 := newTestBattle(flatField())
	// 讓對面只剩場上的 48 個。
	for k := range b2.Sides[0].Reserve {
		b2.Sides[0].Reserve[k] = 0
	}
	if !runQuery(b2, 1, opQFoeMen, 0, SoldiersOnFoot) {
		t.Errorf("指令 12 應回敵方剩餘 %d", SoldiersOnFoot)
	}
}

// 指令 14 問的是**場上**的人數（0–48），與待機的無關——
// 原版讀 `word_1D31C+1`，而那個計數器只掃 48 個場上的兵。
func TestScriptQueryOnField(t *testing.T) {
	b := newTestBattle(flatField())
	if !runQuery(b, 1, opQOnField, 0, SoldiersOnFoot) {
		t.Errorf("指令 14 應回場上的 %d，不是含待機的總數", SoldiersOnFoot)
	}
}

// 指令 15 在城壁還好好的時候回「最小耐久 ÷ 64」。
func TestScriptQueryWall(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 10)
	for k := 0; k < Squads; k++ {
		b.Deploy(0, k, Infantry, 10)
		b.Deploy(1, k, Infantry, 10)
	}
	b.Place()

	want := (SiegeWallDurability(10) * 4) >> 8
	if !runQuery(b, 1, opQWall, 0, want) {
		t.Errorf("指令 15 應回 %d", want)
	}
}

// 兩側的陣形線常數要照抄，不能自己算對稱。
func TestFormationLineConstants(t *testing.T) {
	want := [2][3]int{{5, 28, 48}, {58, 36, 16}}
	for side := range want {
		for choice, v := range want[side] {
			if got := LineFor(side, choice); got != v {
				t.Errorf("side %d 的第 %d 個陣形線是 %d，應為 %d", side, choice, got, v)
			}
			if got := LineChoice(side, v); got != choice {
				t.Errorf("反查 side %d 的 %d 得到 %d，應為 %d", side, v, got, choice)
			}
		}
	}
	// 對稱算法（Width−1−l）在中央與敵軍側會算錯，留一條測試釘住。
	if LineFor(0, 1) == Width-1-LineFor(1, 1) {
		t.Error("中央那一組不該正好對稱——對稱是巧合，不是規則")
	}
}

// 原版的腳本在**有城壁的攻城戰場**上也跑得起來，而且城壁真的會被打。
func TestRealScriptsSiege(t *testing.T) {
	const dat = "../../../workplace/orig/dosv/BATTLE.DAT"
	raw, err := os.ReadFile(dat)
	if err != nil {
		t.Skip("找不到原版 BATTLE.DAT，跳過")
	}
	hit := 0
	for seg := 0; seg < 32; seg++ {
		f, _ := tiledField(32)
		b := NewBattle(f, SyntheticFormations(), &fixedRand{seq: []int{1, 7, 3}}, 10)
		for k := 0; k < Squads; k++ {
			b.Deploy(0, k, Infantry, 20)
			b.Deploy(1, k, Infantry, 20)
		}
		b.Place()
		b.SetScript(0, NewScript(raw[seg*256:(seg+1)*256], 0))
		b.SetScript(1, NewScript(raw[seg*256:(seg+1)*256], 1))
		start, _ := b.MinWallDurability()
		for i := 0; i < 5000 && !b.Done; i++ {
			b.Step()
		}
		if end, _ := b.MinWallDurability(); end < start {
			hit++
		}
		for _, s := range b.Structures {
			if s.Durability < 0 {
				t.Fatalf("段 %d 跑出負的耐久", seg)
			}
		}
	}
	if hit == 0 {
		t.Error("32 段腳本跑完沒有任何一段碰到城壁——碰撞沒接上")
	}
	t.Logf("32 段腳本裡有 %d 段打到城壁", hit)
}

// TestScriptLineUsesRoleNotSideIndex 釘住陣形線查表要用**角色**不是側編號。
//
// 原版的 side 0 永遠是玩家（docs/spec/56 §1），而 remake 的 side 0 是攻方。
// 玩家守城時兩者對不上：拿側編號去查 `lineX`，腳本會把自己的陣形線設成
// **玩家那一條（X=5）**，整支軍團往對面走——實測小地圖上敵方的部隊點
// 從 X=58 一路掉到 38（docs/spec/59 §6）。
func TestScriptLineUsesRoleNotSideIndex(t *testing.T) {
	b := &Battle{}
	b.SetPlayerSide(DefenderSide) // 玩家守城：side 1 是玩家、side 0 是腳本
	if b.Role(0) != 1 || b.Role(1) != 0 {
		t.Fatalf("角色對應錯了：side0=%d side1=%d", b.Role(0), b.Role(1))
	}
	// 腳本側（side 0）下 `line 0`：要拿到腳本那一列的 58，不是玩家的 5。
	s := &Script{side: 0}
	s.exec(b, opLine, 0, 0)
	if got, want := b.Sides[0].Line, LineFor(1, 0); got != want {
		t.Fatalf("腳本側的陣形線 = %d，預期 %d（腳本那一列）", got, want)
	}
	if b.Sides[0].Line == LineFor(0, 0) {
		t.Fatal("腳本側拿到了玩家那一條線")
	}
	// 玩家當攻方時側編號與角色一致，行為不能被這個修正改掉。
	b2 := &Battle{}
	s2 := &Script{side: 1}
	s2.exec(b2, opLine, 0, 0)
	if got, want := b2.Sides[1].Line, LineFor(1, 0); got != want {
		t.Fatalf("玩家攻方時腳本側的線 = %d，預期 %d", got, want)
	}
}
