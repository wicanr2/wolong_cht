package state

// refreshMonthlyGeneralScores 是 sub_155A6；不讀 RNG，僅更新存在的 127 槽。
func (w *World) refreshMonthlyGeneralScores() {
	for i := range w.Generals {
		g := &w.Generals[i]
		if !g.Alive {
			continue
		}
		var score uint8
		for _, aptitude := range g.Aptitude {
			score += uint8(aptitude)
		}
		score += uint8(g.Martial) << 1
		score += uint8(g.Command) << 1
		g.MonthlyScore = score
	}
}
