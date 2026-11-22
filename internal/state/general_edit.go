package state

// GeneralAbilities are the seven editable stats. Aptitude order is siege,
// land, and naval. Tactic selects one of eight delegated battle scripts.
type GeneralAbilities struct {
	Martial, Command, Politics, Tactic int
	Aptitude                           [3]int
}

func (g General) Abilities() GeneralAbilities {
	return GeneralAbilities{g.Martial, g.Command, g.Politics, g.Tactic, g.Aptitude}
}

// ApplyGeneralAbilities validates the complete draft before changing the world.
// Unchanged legacy values can survive an edit to an unrelated field.
func (w *World) ApplyGeneralAbilities(id int, next GeneralAbilities) bool {
	if id < 0 || id >= len(w.Generals) || w.HasIncompleteTick() || w.Outcome() != InProgress {
		return false
	}
	g := &w.Generals[id]
	valid := func(value, previous, low, high int) bool {
		return value == previous || value >= low && value <= high
	}
	if !valid(next.Martial, g.Martial, 1, 15) || !valid(next.Command, g.Command, 1, 15) ||
		!valid(next.Politics, g.Politics, 1, 15) || !valid(next.Tactic, g.Tactic, 0, 7) {
		return false
	}
	for i, v := range next.Aptitude {
		if !valid(v, g.Aptitude[i], 0, 10) {
			return false
		}
	}
	g.Martial, g.Command, g.Politics, g.Tactic = next.Martial, next.Command, next.Politics, next.Tactic
	g.Aptitude = next.Aptitude
	return true
}
