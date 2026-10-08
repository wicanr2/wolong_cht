package state

import "testing"

func TestMonthlyGeneralScoreByteAndInactive(t *testing.T) {
	w := &World{}
	w.Generals[0] = General{Alive: true, Aptitude: [3]int{1, 2, 3}, Martial: 200, Command: 150}
	w.Generals[1] = General{Alive: false, MonthlyScore: 0xA5}
	w.refreshMonthlyGeneralScores()
	if w.Generals[0].MonthlyScore != 194 || w.Generals[1].MonthlyScore != 0xA5 {
		t.Fatalf("score %d inactive %d, want 194/165", w.Generals[0].MonthlyScore, w.Generals[1].MonthlyScore)
	}
}

func TestMonthlyGeneralScoreLoadSave(t *testing.T) {
	b := make([]byte, blockSize)
	b[generalBase+0x1F] = 0x93
	b[generalBase+garrisonGeneral*generalSize+0x1F] = 0xA7
	w, err := LoadBlock(b)
	if err != nil {
		t.Fatal(err)
	}
	if w.Generals[0].MonthlyScore != 0x93 {
		t.Fatal("score not loaded")
	}
	out := w.Bytes()
	if out[generalBase+0x1F] != 0x93 || out[generalBase+garrisonGeneral*generalSize+0x1F] != 0xA7 {
		t.Fatal("score or terminal slot not preserved")
	}
}
