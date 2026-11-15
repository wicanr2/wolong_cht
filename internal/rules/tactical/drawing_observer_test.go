package tactical

import "testing"

type drawingRecorder struct {
	erased [][3]int
	drawn  []SoldierDrawing
	deaths []Death
}

func (*drawingRecorder) BeginDrawingFrame(*Battle) {}
func (r *drawingRecorder) EraseSoldierDrawing(x, y, z int) {
	r.erased = append(r.erased, [3]int{x, y, z})
}
func (r *drawingRecorder) DrawSoldierDrawing(_ int, d SoldierDrawing) {
	r.drawn = append(r.drawn, d)
}
func (r *drawingRecorder) DrawDeathDrawing(_ int, d Death) {
	r.deaths = append(r.deaths, d)
}

func TestDrawingErasesSwappedPreviousCoordinates(t *testing.T) {
	b := &Battle{}
	b.Sides[0].Soldiers[0] = Soldier{Alive: true, Kind: Archer, X: 12, Y: 24}
	b.Sides[0].Soldiers[1] = Soldier{Alive: true, Kind: Archer, X: 30, Y: 28}
	b.initDrawingPrevious()
	b.Sides[0].Soldiers[0].X = 13
	b.Sides[0].Soldiers[1].X = 31
	r := &drawingRecorder{}
	b.SetDrawingObserver(r)
	if !b.swapWith(0, 0, 0, 1) {
		t.Fatal("換位應成功")
	}
	b.recordSoldierDrawing(0, 0, false)
	if len(r.erased) != 1 || r.erased[0] != [3]int{30, 28, 0} {
		t.Fatalf("原版交換word高byte後，應擦另一人的上一座標，得到%v", r.erased)
	}
	if len(r.drawn) != 1 || r.drawn[0].X != 31 {
		t.Fatalf("應於目前座標31重新登錄，得到%v", r.drawn)
	}
}

func TestReinforcementResetsPreviousDrawingBeforeErase(t *testing.T) {
	b := &Battle{Field: &Field{}}
	b.Sides[0].Reserve[0] = 1
	b.Sides[0].Soldiers[1] = Soldier{Kind: Archer, X: 42, Y: 30}
	b.initDrawingPrevious()
	r := &drawingRecorder{}
	b.SetDrawingObserver(r)
	b.reinforceSoldier(0, 1)
	if len(r.erased) != 1 || r.erased[0] != [3]int{1, 30, 0} {
		t.Fatalf("補兵應先重設上一繪圖座標，得到%v", r.erased)
	}
	if len(r.drawn) != 1 || r.drawn[0].X != 1 || !b.Sides[0].Soldiers[1].Alive {
		t.Fatal("補兵應在邊緣登錄")
	}
}

func TestDeathDrawingErasesThroughRemoval(t *testing.T) {
	b := &Battle{}
	b.Sides[0].Soldiers[1] = Soldier{Kind: Archer, X: 22, Y: 30}
	b.initUnitCollision()
	b.initDrawingPrevious()
	b.Sides[0].Soldiers[1].X = 23
	b.deaths = []Death{{Side: 0, Slot: 1, X: 22, Y: 30, Kind: Archer, FramesLeft: 4, Pending: true}}
	r := &drawingRecorder{}
	b.SetDrawingObserver(r)
	b.advanceDeath(0)
	if len(r.erased) != 1 || r.erased[0] != [3]int{22, 30, 0} || len(r.deaths) != 1 || r.deaths[0].X != 23 {
		t.Fatalf("倒地首次擦舊畫新不符：%+v", r)
	}
	for len(b.deaths) != 0 {
		b.advanceDeath(0)
	}
	if len(r.erased) != 4 || len(r.deaths) != 3 || r.erased[3] != [3]int{23, 30, 0} {
		t.Fatalf("最後一拍應只擦除，得到%+v", r)
	}
}
