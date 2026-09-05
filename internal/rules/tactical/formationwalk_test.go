package tactical

import "testing"

// ⭐ 一列縱隊的陣形裡，**大將站在正中間**，所以出生點在大將另一側的
// 隊員一定要越過大將那一格才到得了自己的位置——而大將是不能對調的
// （`docs/re/11` §5.16 的四道閘之一）。
//
// 唯一的出路是**繞出那一欄再回來**，也就是 `docs/spec/134` 的波數佇列。
// 原版實測就是這樣走的（`docs/playtest/75`）。
//
// ⚠ 判準是「每個人都站上自己的格子」，不是「大部分人到位」——
// 先前卡住的只有 96 個兵裡的一個。
func TestColumnFormationFillsEverySlot(t *testing.T) {
	stack := make([][]int, Height)
	for y := range stack {
		stack[y] = make([]int, Width)
	}
	// 原版陣形 0 的第一隊：dx 全是 −2，dy ＝ 0, −2, +2, −4, +4, −6, +6, −8
	// （`KI.EXE` 檔案偏移 0xCEE4，docs/re/11 §5.8d）。
	forms := &Formations{}
	col := [PerSquad][2]int8{
		{-2, 0}, {-2, -2}, {-2, 2}, {-2, -4},
		{-2, 4}, {-2, -6}, {-2, 6}, {-2, -8},
	}
	for form := 0; form < NumFormations; form++ {
		for k := 0; k < SoldiersOnFoot; k++ {
			forms.off[form][k] = col[k%PerSquad]
		}
	}
	var r seqRand
	b := NewBattle(NewField(stack, 0), forms, &r, 0)
	for sq := 0; sq < Squads; sq++ {
		b.Deploy(0, sq, Infantry, 100)
	}
	b.Deploy(1, 0, Infantry, 100)
	for k := range b.Sides[0].Soldiers {
		s := &b.Sides[0].Soldiers[k]
		s.X, s.Y = 1, 10+k%40
	}
	b.Sides[0].Line, b.Sides[0].Mirror = 3, true
	b.Sides[1].Line, b.Sides[1].Mirror = 62, false

	// 出生點刻意讓位 2／3／4 落在大將（目標 Y=32）的另一側。
	ys := [PerSquad]int{32, 27, 26, 36, 29, 21, 47, 31}
	for k := 0; k < PerSquad; k++ {
		s := &b.Sides[1].Soldiers[k]
		s.X, s.Y = 62, ys[k]
	}
	for f := 0; f < 200; f++ {
		b.Step()
	}
	for k := 0; k < PerSquad; k++ {
		s := &b.Sides[1].Soldiers[k]
		gx, gy := b.formationSpot(1, k)
		if s.X != gx || s.Y != gy {
			t.Errorf("位 %d 停在 (%d,%d)，陣形位置是 (%d,%d)——越不過大將那一格",
				k, s.X, s.Y, gx, gy)
		}
	}
}
