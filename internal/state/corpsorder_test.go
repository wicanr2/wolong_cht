package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// formOne 編一支六槽滿編的軍團，回傳主將（＝軍團）編號。
func formOne(t *testing.T, w *World, faction int) int {
	t.Helper()
	w.Factions[faction].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	leader := w.Factions[faction].Lord
	var kinds [army.Positions]army.TroopType
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(leader, kinds, manned); err != nil {
		t.Fatal(err)
	}
	return leader
}

// 三個分支各自寫對 Delegated 與 Stage（docs/spec/39 §1.1）。
func TestSetMarchModeWritesRawValues(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	i := formOne(t, w, f)
	capital := w.Factions[f].Capital

	if err := w.March(i, capital); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		mode      MarchMode
		delegated bool
		stage     int
	}{
		{MarchDelegate, true, StageNormal},
		{MarchCommand, false, StageNormal},
		{MarchDisband, false, StageDisband},
	} {
		if err := w.SetMarchMode(i, tc.mode); err != nil {
			t.Fatalf("%v：%v", tc.mode, err)
		}
		c := w.Corps[i]
		if c.Delegated != tc.delegated || c.Stage != tc.stage {
			t.Errorf("%v → Delegated=%v Stage=%d，want %v／%d",
				tc.mode, c.Delegated, c.Stage, tc.delegated, tc.stage)
		}
		if c.Timer != 1 {
			t.Errorf("%v 之後計時器 = %d，原版寫 1", tc.mode, c.Timer)
		}
	}
}

// 「解體」只在目標據點就是自己的首都時才給（sub_17FDB 的 cmp [bx+3], [bp+0]）。
func TestDisbandOnlyOfferedAtCapital(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	i := formOne(t, w, f)
	capital := w.Factions[f].Capital

	if !w.DisbandAllowed(i) {
		t.Fatal("剛編成的軍團目標就是首都，應該可以解體")
	}
	other := -1
	for n := range w.Cities {
		if n != capital && w.Cities[n].Owner == f {
			other = n
			break
		}
	}
	if other < 0 {
		t.Skip("這個劇本的這個勢力只有一座城")
	}
	if err := w.March(i, other); err != nil {
		t.Fatal(err)
	}
	if w.DisbandAllowed(i) {
		t.Error("目標不是首都時不該給「解體」")
	}
	if err := w.SetMarchMode(i, MarchDisband); err == nil {
		t.Error("目標不是首都時 SetMarchMode(解體) 應該回錯誤")
	}
}

// 解散的五個動作（docs/re/64 §3）——remake 對得上的四個。
func TestDisbandReturnsMenAndFreesLeader(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	i := formOne(t, w, f)
	before := w.Factions[f].Reserves
	corpsCount := w.Factions[f].Corps
	men := w.Corps[i].Men
	if men == 0 {
		t.Fatal("軍團一個兵都沒有，測不到退回")
	}

	w.disbandCorps(i)

	if w.Corps[i].Alive {
		t.Error("解散之後軍團還在")
	}
	if w.Generals[i].Posted {
		t.Error("解散之後主將還掛著職務")
	}
	if got := w.Factions[f].Corps; got != corpsCount-1 {
		t.Errorf("軍團數 = %d，want %d", got, corpsCount-1)
	}
	total := 0
	for k := range before {
		total += w.Factions[f].Reserves[k] - before[k]
	}
	if total != men {
		t.Errorf("退回預備兵池 %d 點，軍團原本有 %d 點", total, men)
	}
}

// Stage 11 的軍團停在首都上，抵達處理就會解散它。
func TestArriveAtCapitalDisbands(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	i := formOne(t, w, f)
	capital := w.clampCity(w.Factions[f].Capital)
	c := &w.Corps[i]
	c.Node, c.TargetNode, c.Ordered = capital, capital, w.Factions[f].Capital
	c.X, c.Y = w.Cities[capital].X, w.Cities[capital].Y
	c.TargetX, c.TargetY = c.X, c.Y
	if err := w.SetMarchMode(i, MarchDisband); err != nil {
		t.Fatal(err)
	}
	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Alive {
		t.Error("在首都上的解體軍團沒有被解散")
	}
}

// 玩家的軍團開回首都且未滿編 → 轉 Stage 9，下一次抵達處理補兵。
func TestArriveAtCapitalResupplies(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	w.Player = f
	i := formOne(t, w, f)
	capital := w.clampCity(w.Factions[f].Capital)
	c := &w.Corps[i]
	c.Node, c.TargetNode = capital, capital
	c.X, c.Y = w.Cities[capital].X, w.Cities[capital].Y
	c.TargetX, c.TargetY = c.X, c.Y
	// 打掉一半的兵，並把池子補回去（模擬打完仗回首都）。
	for k := range c.Units {
		if c.Units[k].Kind != EmptySlotKind {
			c.Men -= c.Units[k].Men / 2
			c.Units[k].Men -= c.Units[k].Men / 2
		}
	}
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{600, 600, 600}
	hurt := c.Men

	w.arriveCorps(i, &testRand{}) // Stage 0 → 9
	if w.Corps[i].Stage != StageResupply {
		t.Fatalf("未滿編的軍團回首都應該轉 Stage %d，得到 %d",
			StageResupply, w.Corps[i].Stage)
	}
	w.arriveCorps(i, &testRand{}) // Stage 9 → 補兵
	if w.Corps[i].Men <= hurt {
		t.Errorf("補兵之後兵力 %d，補之前 %d", w.Corps[i].Men, hurt)
	}
	if w.Corps[i].Stage != StageDone {
		t.Errorf("補完兵的 Stage = %d，原版寫 %d", w.Corps[i].Stage, StageDone)
	}
}

// 滿編的軍團回首都不補兵（原版 cmp word [si+4], 258h）。
func TestFullCorpsDoesNotResupply(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	w.Player = f
	i := formOne(t, w, f)
	capital := w.clampCity(w.Factions[f].Capital)
	c := &w.Corps[i]
	c.Node, c.TargetNode = capital, capital
	if c.Men < army60000Points {
		t.Skipf("這一支只有 %d 點，測不到滿編那條路", c.Men)
	}
	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Stage != StageNormal {
		t.Errorf("滿編的軍團回首都不該轉補兵，Stage = %d", w.Corps[i].Stage)
	}
}

// 委任看的是**玩家那一方**，不是攻方：原版兩條路各檢查各自那一方
// （`docs/re/09` §2）。
func TestDelegatedPlayerSideSkipsEncounter(t *testing.T) {
	w := load(t, 0)
	alive := w.AliveFactions()
	me, foe := alive[0], alive[1]
	w.Player = me
	mine, theirs := formOne(t, w, me), formOne(t, w, foe)
	// 假的戰場來源：只要不是 nil，wantsTactical 就會走到委任那一段。
	w.SetTactical(&TacticalSetup{
		Forms: &tactical.Formations{},
		Field: func(int, bool) *tactical.Field { return nil },
	})

	for _, tc := range []struct {
		name             string
		att, def         int
		delegate         int // 要把誰設成委任，−1 表示都不設
		wantsPlayerInput bool
	}{
		{"玩家是攻方、沒委任", mine, theirs, -1, true},
		{"玩家是攻方、委任", mine, theirs, mine, false},
		{"玩家是守方、沒委任", theirs, mine, -1, true},
		{"玩家是守方、委任", theirs, mine, mine, false},
		{"敵方委任不影響玩家", mine, theirs, theirs, true},
	} {
		w.Corps[mine].Delegated, w.Corps[theirs].Delegated = false, false
		if tc.delegate >= 0 {
			w.Corps[tc.delegate].Delegated = true
		}
		if got := w.wantsTactical(tc.att, tc.def); got != tc.wantsPlayerInput {
			t.Errorf("%s：wantsTactical = %v，want %v", tc.name, got, tc.wantsPlayerInput)
		}
	}
}
