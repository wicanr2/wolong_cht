package tactical

import (
	"reflect"
	"testing"
)

func TestSoundRequestsKeepRepeatedCodesAndOrder(t *testing.T) {
	b := &Battle{}
	for _, code := range []uint8{12, 12, 10, 11, 12} {
		b.emitSFX(code)
	}
	if got := b.TakeSoundEffects(); !reflect.DeepEqual(got, []uint8{12, 12, 10, 11, 12}) {
		t.Fatalf("請求應保留次數與順序，得到%v", got)
	}
	if len(b.TakeSoundEffects()) != 0 {
		t.Fatal("取走後不能重播")
	}
}

func TestNextStepDiscardsUnconsumedSoundHistory(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
	b.emitSFX(12)
	b.emitSFX(12)
	b.Step()
	if got := b.TakeSoundEffects(); len(got) != 0 {
		t.Fatalf("未消費的上一拍聲音不能累積到下一拍：%v", got)
	}
}
