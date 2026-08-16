package tactical

import (
	"os"
	"testing"
)

// tiledField 造一張帶原始圖塊值的戰場：X = 32 那一行從 top 到 bottom
// 是城壁圖塊，gate 那一格是門。
func tiledField(gateX int) (*Field, [][]byte) {
	tiles := make([][]byte, Height)
	for y := range tiles {
		tiles[y] = make([]byte, Width)
	}
	var heights [256]int
	heights[TileWallLo] = 4 // 城壁：四層，爬得上去但走不過去
	heights[TileGateLo] = 4
	heights[TileWallLo+brokenWallDelta] = 0 // 瓦礫：平的
	heights[TileGateLo+brokenGateDelta] = 0

	const top, bottom = 8, Height - 9
	gate := Height / 2
	for y := top; y <= bottom; y++ {
		tiles[y][32] = TileWallLo
	}
	tiles[gate][32] = TileGateLo
	return NewFieldFromTiles(tiles, &heights, gateX), tiles
}

// 一行連續的城壁只算一段，門另外算一筆——原版 `sub_19CE2` 用 dh 累計連續，
// 遇到非城壁的圖塊才把長度寫進 `+0x1A`。門把中間切開，所以是兩段。
func TestStructuresMergeRuns(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)

	var walls, gates int
	for _, s := range b.Structures {
		switch s.Kind {
		case KindWall:
			walls++
		case KindGate:
			gates++
		}
	}
	if walls != 2 {
		t.Errorf("城壁應為 2 段（被門切開），得到 %d", walls)
	}
	if gates != 1 {
		t.Errorf("門應為 1 道，得到 %d", gates)
	}

	const top, bottom = 8, Height - 9
	gate := Height / 2
	total := 0
	for _, s := range b.Structures {
		if s.Kind == KindWall {
			total += s.Run
		}
	}
	if want := (bottom - top + 1) - 1; total != want {
		t.Errorf("城壁總長 %d，應為 %d（%d 格扣掉門那一格）", total, want, want+1)
	}
	if b.Structures[0].Y != top {
		t.Errorf("第一段從 Y=%d 起，應為 %d", b.Structures[0].Y, top)
	}
	if b.Structures[1].Y != gate+1 {
		t.Errorf("第二段從 Y=%d 起，應為 %d", b.Structures[1].Y, gate+1)
	}
}

// 攻城戰的城壁耐久 ＝（城兵數 ＋ 50）× 10；野戰固定 300；門固定 80。
func TestStructureDurability(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 77)
	for _, s := range b.Structures {
		want := SiegeWallDurability(77)
		if s.Kind == KindGate {
			want = GateDurability
		}
		if s.Durability != want {
			t.Errorf("%s 耐久 %d，應為 %d", structureName(s.Kind), s.Durability, want)
		}
	}
	if got := SiegeWallDurability(77); got != 1270 {
		t.Errorf("(77+50)×10 應為 1270，得到 %d", got)
	}

	// gateX 為 0 就是野戰，城壁改成 300。
	f2, _ := tiledField(0)
	b2 := NewBattle(f2, SyntheticFormations(), &fixedRand{}, 77)
	for _, s := range b2.Structures {
		if s.Kind == KindWall && s.Durability != FieldWallDurability {
			t.Errorf("野戰城壁耐久 %d，應為 %d", s.Durability, FieldWallDurability)
		}
	}
}

// 指令 15 查到的是「最小耐久 ÷ 64」，但**只要有一段破了就回 0**。
func TestWallQuery(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 10)

	min, intact := b.MinWallDurability()
	if !intact {
		t.Fatal("開場不該有破掉的城壁")
	}
	if want := SiegeWallDurability(10); min != want {
		t.Fatalf("最小耐久 %d，應為 %d", min, want)
	}
	if got, want := b.WallQuery(), (min*4)>>8; got != want {
		t.Errorf("指令 15 回 %d，應為 %d", got, want)
	}

	// 打壞一段之後一律回 0。
	b.breakRow(b.Structures[0].Y)
	if _, intact := b.MinWallDurability(); intact {
		t.Error("打壞之後 intact 應為 false")
	}
	if got := b.WallQuery(); got != 0 {
		t.Errorf("破了之後指令 15 應回 0，得到 %d", got)
	}
}

// 撞一次掉一點耐久；歸零那一下整段垮，而且**地形跟著變平**。
func TestHitStructureBreaksAndFlattens(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	s := &b.Structures[0]
	y := s.Y

	if f.StandLevel(32, y) != 4 {
		t.Fatalf("開場那一格應該是 4 層，得到 %d", f.StandLevel(32, y))
	}
	for i := s.Durability; i > 0; i-- {
		if !b.hitStructure(0, East, 32, y) {
			t.Fatal("撞在城壁上應該要有反應")
		}
	}
	if s.Durability != 0 {
		t.Fatalf("耐久應該歸零，得到 %d", s.Durability)
	}
	if s.Broken {
		t.Fatal("耐久歸零的那一幀還沒垮——原版是下一次撞上才垮")
	}
	b.hitStructure(0, East, 32, y)
	if !b.Structures[0].Broken {
		t.Fatal("耐久 0 再撞一次應該垮")
	}
	if got := f.StandLevel(32, y); got != 0 {
		t.Errorf("垮掉之後應該是平地，得到 %d 層", got)
	}
	if got := f.Tile(32, y); got != TileWallLo+brokenWallDelta {
		t.Errorf("圖塊值應換成 %#x，得到 %#x", TileWallLo+brokenWallDelta, got)
	}
}

// 突擊會把門全部打開，而且開了關不回去（說明書 4.2）。
// 守方突擊：`sub_1B7CB` 只挑**類型 1（城壁）**（`80 FC 01` ＝ `cmp ah, 1`），
// **門不在裡面**，而且**只有守方**那一側會觸發。
func TestDefenderChargeBreaksWallsNotGates(t *testing.T) {
	f, tiles := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	for k := 0; k < Squads; k++ {
		b.Deploy(DefenderSide, k, Infantry, 10)
	}
	b.Place()

	gate := Height / 2
	wall := gate - 3
	if f.StandLevel(32, wall) == 0 || f.StandLevel(32, gate) == 0 {
		t.Fatal("城壁與門那兩格開場都該是擋著的")
	}

	// 攻方突擊什麼都不拆。
	b.Order(AttackerSide, -1, Charge)
	if f.StandLevel(32, wall) == 0 {
		t.Error("攻方突擊不該拆城壁")
	}
	if b.Sides[AttackerSide].Sortied {
		t.Error("攻方不該被標成已出擊")
	}

	b.Order(DefenderSide, -1, Charge)
	if !b.Sides[DefenderSide].Sortied {
		t.Error("守方突擊之後要標成已出擊")
	}
	if f.StandLevel(32, wall) != 0 {
		t.Error("守方突擊之後城壁那一格應該走得過去")
	}
	if got := tiles[wall][32]; got != TileWallLo+brokenWallDelta {
		t.Errorf("城壁圖塊 = %#x，預期換成瓦礫 %#x", got, TileWallLo+brokenWallDelta)
	}
	for _, s := range b.Structures {
		if s.Kind == KindWall && !s.Broken {
			t.Error("守方突擊之後城壁應該全破")
		}
		if s.Kind == KindGate && s.Broken {
			t.Error("⚠ 門不該被突擊拆掉——sub_1B7CB 只挑類型 1")
		}
	}
	if f.StandLevel(32, gate) == 0 {
		t.Error("門那一格不該被突擊打通")
	}
}

// 城壁與門的圖塊區間拿原版的 `BATTLE.MAP` 對過。
//
// 三條檢查，每一條都是「解錯就會不成立」的性質：
//
//  1. **186 張攻城戰場全部有城壁圖塊，零例外**——區間抓錯就不會全中。
//  2. **一張圖的「城壁段數 ＋ 門數」從來不超過 16**，正好等於原版
//     0x0C00–0x0DFF 那 16 筆的額度。多算或少算都會有圖超出去。
//  3. 段數與門數都不是 0（每張攻城圖至少有一段城壁、一道門）。
func TestStructuresAgainstRealMap(t *testing.T) {
	const path = "../../../workplace/orig/dosv/BATTLE.MAP"
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Skip("找不到原版 BATTLE.MAP，跳過")
	}
	const idx, fieldSize, cellsOff, numFields = 512, 4096, 0x40, 214

	var heights [256]int
	for i := TileWallLo; i <= TileWallHi; i++ {
		heights[i] = 4
	}
	for i := TileGateLo; i <= TileGateHi; i++ {
		heights[i] = 4
	}

	siege, withWall, overflow := 0, 0, 0
	for n := 0; n < numFields; n++ {
		gateX := int(raw[n*2+1])
		base := idx + n*fieldSize + cellsOff
		tiles := make([][]byte, Height)
		for y := 0; y < Height; y++ {
			tiles[y] = append([]byte(nil), raw[base+y*Width:base+(y+1)*Width]...)
		}
		f := NewFieldFromTiles(tiles, &heights, gateX)
		got := buildStructures(tiles, f.IsSiege(), 0)

		if gateX == 0 {
			continue // 野戰
		}
		siege++
		walls, gates := 0, 0
		for _, s := range got {
			if s.Kind == KindWall {
				walls++
			} else {
				gates++
			}
		}
		if walls > 0 && gates > 0 {
			withWall++
		}
		if len(got) >= MaxStructures {
			overflow++
		}
	}
	if siege != 186 {
		t.Errorf("攻城戰場有 %d 張，應為 186", siege)
	}
	if withWall != siege {
		t.Errorf("只有 %d／%d 張攻城戰場同時有城壁與門——圖塊區間可能抓錯",
			withWall, siege)
	}
	if overflow != 0 {
		t.Errorf("有 %d 張圖的城壁段數＋門數塞滿或超過 16 筆的額度", overflow)
	}
}

// ⭐ `seg000:B5B7` 那三條分支，行為照抄。
//
// 每一個運算元都查證過（見 hitStructure 的說明）：守方碰不壞城壁、
// 攻方背對城的方向撞上去耐久直接歸零、其餘只減 1。
func TestHitStructureSiegeBranches(t *testing.T) {
	newSiege := func() (*Battle, int) {
		f, _ := tiledField(32)
		b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 10)
		return b, b.Structures[0].Y
	}

	// ① 守方（side 1）碰不壞。
	b, y := newSiege()
	before := b.Structures[0].Durability
	if b.hitStructure(1, East, 32, y) {
		t.Error("守方不該碰得壞城壁")
	}
	if b.Structures[0].Durability != before {
		t.Errorf("守方撞了之後耐久從 %d 變成 %d", before, b.Structures[0].Durability)
	}

	// ② 攻方朝著城走（East）→ 只減 1。
	b, y = newSiege()
	before = b.Structures[0].Durability
	b.hitStructure(0, East, 32, y)
	if got := b.Structures[0].Durability; got != before-1 {
		t.Errorf("朝城撞一次耐久 %d，應為 %d", got, before-1)
	}

	// ③ 攻方背對城（West）→ 直接歸零。
	b, y = newSiege()
	b.hitStructure(0, West, 32, y)
	if got := b.Structures[0].Durability; got != 0 {
		t.Errorf("背對城撞上去耐久 %d，應為 0", got)
	}
	// 再撞一次就垮。
	b.hitStructure(0, West, 32, y)
	if !b.Structures[0].Broken {
		t.Error("耐久 0 再撞一次應該垮")
	}

	// ④ 野戰沒有這個分支——背對也只減 1。
	f, _ := tiledField(0) // gateX 0 ＝ 野戰
	fb := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	before = fb.Structures[0].Durability
	fb.hitStructure(0, West, 32, fb.Structures[0].Y)
	if got := fb.Structures[0].Durability; got != before-1 {
		t.Errorf("野戰背對撞一次耐久 %d，應為 %d（那個分支只在攻城戰）", got, before-1)
	}
}

// 面向只在**走成功**那一步才更新——被牆擋住的兵保持原本的面向。
//
// 原版把 `[si+5]` 寫在四個移動常式裡，而那些常式只有走得動時才被呼叫。
func TestFacingOnlyUpdatesOnSuccessfulMove(t *testing.T) {
	f := walledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	b.Deploy(0, 0, Infantry, 1)
	s := &b.Sides[0].Soldiers[0]
	s.X, s.Y, s.Z = 31, 10, 0 // 牆的正左邊（那一列沒有門）
	s.Facing = West
	s.GoalX, s.GoalY, s.GoalZ = 40, 10, 0 // 目標在牆的另一邊

	b.moveToward(0, 0)
	if s.X != 31 {
		t.Fatalf("兵越過了牆，X ＝ %d", s.X)
	}
	if s.Facing != West {
		t.Errorf("被牆擋住卻把面向改成 %d，應該保持 West(%d)", s.Facing, West)
	}
}

// 門強度條：只在更小的耐久出現時更新，20 幀後自己收掉（docs/spec/32）。
func TestStructureBarFollowsRawMinimumAndExpiry(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	if _, shown := b.StructureBar(); shown {
		t.Fatal("還沒挨打就不該有門強度條")
	}

	s := &b.Structures[0]
	y, start := s.Y, s.Durability
	b.hitStructure(0, East, 32, y)
	got, shown := b.StructureBar()
	if !shown || got != start-1 {
		t.Fatalf("第一次挨打後 = %d,%v，預期 %d,true", got, shown, start-1)
	}

	// 更大的耐久不覆蓋已畫過的最小值——原版 word_1C405 只往下走。
	b.noteStructureDamage(start + 100)
	if got, _ := b.StructureBar(); got != start-1 {
		t.Errorf("更大的耐久不該覆蓋，得到 %d", got)
	}

	// 到期：原版 word_1D318 每幀 +1，**等於** word_1D326 才收。
	for i := 0; i < StructureBarLifetime-1; i++ {
		b.Frame++
		b.expireStructureBar()
		if _, shown := b.StructureBar(); !shown {
			t.Fatalf("第 %d 幀就收掉了，預期撐滿 %d 幀", i+1, StructureBarLifetime)
		}
	}
	b.Frame++
	b.expireStructureBar()
	if _, shown := b.StructureBar(); shown {
		t.Fatalf("第 %d 幀應該收掉", StructureBarLifetime)
	}

	// 收掉之後重新挨打會重建，最小值也跟著重設。
	b.hitStructure(0, East, 32, y)
	if got, shown := b.StructureBar(); !shown || got != start-2 {
		t.Fatalf("重建後 = %d,%v，預期 %d,true", got, shown, start-2)
	}
}

// Step 有把到期檢查接進去（原版 0001A12A 是遞增計時器的同一支）。
func TestStructureBarExpiresThroughStep(t *testing.T) {
	b := newTestBattle(walledField(32))
	b.noteStructureDamage(500)
	for i := 0; i < StructureBarLifetime; i++ {
		if b.Done {
			t.Fatalf("第 %d 幀戰鬥就結束了，這個 fixture 撐不到到期", i)
		}
		b.Step()
	}
	if _, shown := b.StructureBar(); shown {
		t.Fatal("Step 應該在第 0x14 幀收掉門強度條")
	}
}

// 選取位元圖：切換、取出即清空、0 代表全隊（docs/spec/33 §1.3）。
func TestSquadSelectionMatchesRawBitfield(t *testing.T) {
	b := newTestBattle(walledField(32))
	if got := b.TakeSquadSelection(0); got != AllSquadsMask {
		t.Fatalf("一隊都沒選時應回全隊 %#x，得到 %#x", AllSquadsMask, got)
	}

	b.ToggleSquadSelection(0, 2)
	b.ToggleSquadSelection(0, 5)
	if !b.Sides[0].SquadSelected(2) || !b.Sides[0].SquadSelected(5) ||
		b.Sides[0].SquadSelected(0) {
		t.Fatalf("選取狀態錯誤：%#x", b.Sides[0].Selected)
	}
	b.ToggleSquadSelection(0, 2) // 再點一次取消
	if b.Sides[0].SquadSelected(2) {
		t.Error("再點一次應該取消")
	}

	if got, want := b.TakeSquadSelection(0), uint8(1<<5); got != want {
		t.Fatalf("取出 = %#x，預期 %#x", got, want)
	}
	if b.Sides[0].Selected != 0 {
		t.Error("取出之後位元圖必須清空（原版是 xchg）")
	}
	// 越界不得改到別人。
	b.ToggleSquadSelection(0, -1)
	b.ToggleSquadSelection(0, Squads)
	if b.Sides[0].Selected != 0 {
		t.Errorf("越界的隊編號不該改動位元圖：%#x", b.Sides[0].Selected)
	}
}

// 只有被選中的隊收到命令；選取在下令後清空。
func TestOrderSelectedOnlyReachesPickedSquads(t *testing.T) {
	b := newTestBattle(walledField(32))
	b.ToggleSquadSelection(0, 1)
	if !b.OrderSelected(0, Charge) {
		t.Fatal("突擊不該被拒絕")
	}
	for k := 0; k < Squads; k++ {
		got := b.Sides[0].Soldiers[k*PerSquad].Next
		want := Command(Charge)
		if k != 1 {
			want = b.Sides[0].Soldiers[k*PerSquad].Cmd
		}
		if k == 1 && got != want {
			t.Errorf("第 %d 隊應收到突擊，得到 %v", k, got)
		}
		if k != 1 && got == Charge {
			t.Errorf("第 %d 隊沒被選中，不該收到突擊", k)
		}
	}
	if b.Sides[0].Selected != 0 {
		t.Error("下完令選取要清空")
	}
}

// ⚠ 城壁令對**玩家**是拒絕，對腳本是降級成攻擊——兩條路不同。
func TestOrderSelectedRejectsScaleWallOffSiege(t *testing.T) {
	field := newTestBattle(NewField(make([][]int, Height), 0)) // gateX=0 ＝ 沒有城
	if field.Field.IsSiege() {
		t.Fatal("這個 fixture 應該是野戰")
	}
	if field.OrderSelected(0, ScaleWal) {
		t.Error("野戰的城壁令應該被拒絕")
	}
	if field.Sides[0].Soldiers[0].Next == ScaleWal {
		t.Error("被拒絕就不該下令")
	}
	// 腳本那一條仍然是降級，不是拒絕。
	field.Order(0, -1, ScaleWal)
	if got := field.Sides[0].Soldiers[0].Next; got != Attack {
		t.Errorf("腳本路徑應降級成攻擊，得到 %v", got)
	}

	siege := newTestBattle(walledField(32))
	if !siege.OrderSelected(0, ScaleWal) {
		t.Error("攻城戰的城壁令不該被拒絕")
	}
}

// 純地形走不通時，尋路要改成「可以拆的就穿過去」——兵才會走到那一格
// 撞上去，把耐久打掉。⚠ 這是 remake 的近似（原版用 0x1800 那張繞路點
// 清單，演算法未解）。
// sealedField 是一張被城壁完整封死的圖：X=32 整column 都是城壁，
// 中間一格是門。純地形繞不過去。
func sealedField() (*Field, [][]byte) {
	tiles := make([][]byte, Height)
	for y := range tiles {
		tiles[y] = make([]byte, Width)
		tiles[y][32] = TileWallLo
	}
	tiles[Height/2][32] = TileGateLo
	var heights [256]int
	heights[TileWallLo] = 4
	heights[TileGateLo] = 4
	heights[TileWallLo+brokenWallDelta] = 0
	heights[TileGateLo+brokenGateDelta] = 0
	return NewFieldFromTiles(tiles, &heights, 32), tiles
}

// 純地形走不通時，尋路要改成「可以拆的就穿過去」——兵才會走到那一格
// 撞上去，把耐久打掉。⚠ 這是 remake 的近似（原版用 0x1800 那張繞路點
// 清單，演算法未解）。
func TestPathFallsBackThroughBreakableCells(t *testing.T) {
	f, _ := sealedField()
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)

	from, to := Point{X: 20, Y: Height / 2}, Point{X: 45, Y: Height / 2}
	if pts := f.FindPath(from, to, true, nil); len(pts) != 0 {
		t.Fatalf("這張圖被城壁封死，純地形尋路應該回空，得到 %d 段", len(pts))
	}
	if len(f.FindPathForcing(from, to, true, b.breachCost, b.breakableAt)) == 0 {
		t.Fatal("允許撞穿之後應該找得到路")
	}

	// 門比城壁便宜（耐久 80 對上千），成本要反映這一點。
	if b.breachCost(32, Height/2) >= b.breachCost(32, 1) {
		t.Error("門的撞穿成本應該低於城壁")
	}
	if b.breachCost(20, 20) != 0 {
		t.Error("空地不該有撞穿成本")
	}

	// 打壞之後純地形就通了，不必再靠 forcing。
	b.breakRow(Height / 2)
	if len(f.FindPath(from, to, true, nil)) == 0 {
		t.Error("門破了之後純地形尋路應該通")
	}
}

// 被地形擋住的兵要重算繞路，**即使這一幀靠跟同伴對調而「動了」**。
// 沒有這一條的話，前排會靠著跟後排換位一直算「有移動」，永遠不重算路。
func TestTryMoveReportsTerrainBlock(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	b.Deploy(AttackerSide, 0, Infantry, 8)
	b.Place()

	s := &b.Sides[AttackerSide].Soldiers[0]
	y := b.Structures[0].Y
	s.X, s.Y, s.Z = 31, y, 0
	moved, walled := b.tryMove(AttackerSide, 0, 32, y, 0)
	if moved || !walled {
		t.Errorf("撞城壁應該回 (false, true)，得到 (%v, %v)", moved, walled)
	}

	// 走得動的那一步不算被地形擋住。
	s.X, s.Y, s.Z = 20, y, 0
	moved, walled = b.tryMove(AttackerSide, 0, 21, y, 0)
	if !moved || walled {
		t.Errorf("空地應該回 (true, false)，得到 (%v, %v)", moved, walled)
	}
}
