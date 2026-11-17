package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

// DiplomacySettlement是135ED逐人150D7及後續和平／宣戰的runtime續行。
// 每人的狀態先寫回；呈現層確認該人全部對話後，才取下一人。
type DiplomacySettlement struct {
	next      func() (int, bool)
	finish    func() bool
	notice    func() *TalkNotice
	drained   bool
	completed bool
}

func (s *DiplomacySettlement) NextReturnee() (int, bool) {
	if s == nil || s.drained || s.completed {
		return -1, false
	}
	id, ok := s.next()
	if !ok {
		s.drained = true
	}
	return id, ok
}

func (s *DiplomacySettlement) CompletionNotice() *TalkNotice {
	if s == nil || !s.drained || s.completed || s.notice == nil {
		return nil
	}
	return s.notice()
}

// Complete只在全部返回者處理完後執行一次；1A0由呈現層先確認。
func (s *DiplomacySettlement) Complete() bool {
	if s == nil || !s.drained || s.completed {
		return false
	}
	s.completed = true
	return s.finish()
}

func (s *DiplomacySettlement) resolveResult() (bool, []int, *TalkNotice) {
	if s == nil {
		return false, nil, nil
	}
	var released []int
	for {
		id, ok := s.NextReturnee()
		if !ok {
			break
		}
		released = append(released, id)
	}
	notice := s.CompletionNotice()
	return s.Complete(), released, notice
}

func (s *DiplomacySettlement) resolve() (bool, []int) {
	applied, released, _ := s.resolveResult()
	return applied, released
}

func (w *World) prepareQueuedCooperation(source, invader, target, response, amount int) *DiplomacySettlement {
	if response < 0 || response >= 2 {
		return nil
	}
	if response == 1 {
		w.Factions[source].Funds = economy.ClampFunds(w.Factions[source].Funds + amount)
		w.Factions[target].Funds = economy.ClampFunds(w.Factions[target].Funds - amount)
	}
	return &DiplomacySettlement{
		next:   w.diplomaticCaptiveReleases(source, target),
		finish: func() bool { return w.applyQueuedDeclaration(source, invader) },
		notice: func() *TalkNotice {
			if source != w.Player || source < 0 || source >= numFactions || invader < 0 || invader >= numFactions ||
				w.Friendship[source][invader].AtWar() || !w.canApplyQueuedDeclaration(source, invader) {
				return nil
			}
			return &TalkNotice{Index: 0x1A0, City: -1, Faction: invader,
				General: w.Factions[source].Lord, Amount: -1, SpeakerPortrait: true}
		},
	}
}

func (w *World) prepareQueuedCeasefire(source, target, response, amount int) *DiplomacySettlement {
	if response < 0 || response >= 2 {
		return nil
	}
	if response == 1 {
		w.Factions[target].Funds = economy.ClampFunds(w.Factions[target].Funds + amount)
		w.Factions[source].Funds = economy.ClampFunds(w.Factions[source].Funds - amount)
	}
	return &DiplomacySettlement{
		next:   w.diplomaticCaptiveReleases(source, target),
		finish: func() bool { return w.finishDiplomaticPeace(source, target) },
	}
}

// 145F8／14236／13669在所有150D7對話完成後收尾；同步caller也走這一份。
func (w *World) finishDiplomaticPeace(source, target int) bool {
	if w.Factions[source].InvasionTarget == target {
		w.Factions[source].InvasionTarget = diplomacy.NoTarget
	}
	if w.Factions[target].InvasionTarget == source {
		w.Factions[target].InvasionTarget = diplomacy.NoTarget
	}
	for i := range w.Cities {
		c := &w.Cities[i]
		ownerInPair := c.Owner == source || c.Owner == target
		recordedInPair := c.OwnerRecorded == source || c.OwnerRecorded == target
		if ownerInPair && recordedInPair {
			c.OwnerRecorded = c.Owner
		}
	}
	value := w.Friendship[source][target].Value()
	if other := w.Friendship[target][source].Value(); other < value {
		value = other
	}
	w.Friendship[source][target] = diplomacy.Peace(value)
	w.Friendship[target][source] = diplomacy.Peace(value)
	return true
}
