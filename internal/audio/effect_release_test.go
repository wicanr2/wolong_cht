package audio

import "testing"

// 合成的持續音用於區分「不再處理鏈結」與「真正關聲」。
// 這些是測試用音色參數，沒有複製原版音色記錄。
func sustainedTestEffect() Effect {
	raw := make([]byte, effectSize)
	raw[0], raw[1] = 0x21, 0x21
	raw[2] = 0x3f
	raw[4], raw[5] = 0xf0, 0xf0
	raw[6], raw[7] = 0x0f, 0x0f
	raw[10], raw[11], raw[12], raw[14] = 0x30, 0x98, 0x31, 6
	return Effect{Raw: raw}
}

func TestEffectContinuationRetriggersHeldNote(t *testing.T) {
	c := NewOPL3(0)
	effect := sustainedTestEffect()
	effect.program(c, 0, 6)
	for i := 0; i < 2000; i++ {
		c.Generate()
	}
	if c.ops[15].env != envSustain || c.ops[15].phase == 0 {
		t.Fatal("前置條件：測試音應已持續發聲")
	}
	effect.program(c, 0, 6)
	if c.ops[15].env != envAttack || c.ops[15].phase != 0 {
		t.Fatal("接續音必須重新起音，不能沿用上一段持續音")
	}
}

func TestEffectEndingStopsSustainedTone(t *testing.T) {
	list := []Effect{{Raw: make([]byte, effectSize)}, sustainedTestEffect()}
	pcm, err := RenderEffect(list, 1, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	segments := SegmentRMS(pcm, 6)
	if segments[0] < 1e-3 {
		t.Fatalf("前置條件：開頭應有持續音，RMS=%g", segments[0])
	}
	if segments[5] > segments[0]*0.001 {
		t.Fatalf("鏈結結束後仍發聲：開頭RMS=%g，末段RMS=%g", segments[0], segments[5])
	}
}
