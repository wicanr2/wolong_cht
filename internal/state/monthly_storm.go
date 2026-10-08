package state

// saveMonthlyStormGlobals 寫回原版四個保存 word；未月結時 raw 錨點保留原值。
func (w *World) saveMonthlyStormGlobals() {
	values := [4]int{0xFFF0, 0xFFF0, 400, 400}
	if w.stormArea != nil {
		values = [4]int{w.stormArea.MinX, w.stormArea.MinY, w.stormArea.MaxX, w.stormArea.MaxY}
	}
	words := [4]uint16{}
	for i, value := range values {
		words[i] = uint16(value)
	}
	w.stormGlobals = &words
}
