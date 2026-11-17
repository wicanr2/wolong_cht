package state

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func TestEnvoyReportConfirmationPreservesSynchronousResult(t *testing.T) {
	for _, friendship := range []int{0, 100} {
		t.Run(strconv.Itoa(friendship), func(t *testing.T) {
			w, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
			if err != nil {
				t.Fatal(err)
			}
			w.Player = 0
			w.Factions[0].Funds, w.Factions[3].Funds = 50000, 50000
			w.Factions[3].Diplomat, w.Factions[3].Aggression = 37, 0
			w.Factions[0].InvasionTarget, w.Factions[3].InvasionTarget = 3, diplomacy.NoTarget
			w.Generals[37].Duty, w.Generals[37].Politics = DutyDiplomat, 12
			w.Friendship[0][3] = diplomacy.Peace(friendship).WithWar(true)
			w.Friendship[3][0] = diplomacy.Peace(friendship).WithWar(true)
			w.events = [eventQueueEntries]QueuedEvent{}
			w.events[0] = QueuedEvent{Code: 0x0306}
			synchronous, err := LoadBlock(w.Bytes())
			if err != nil {
				t.Fatal(err)
			}
			w.rng, synchronous.rng = rng.NewFixed(1), rng.NewFixed(1)
			w.eventDelay, synchronous.eventDelay = 1, 1
			w.SetDiplomacyReportConfirmation(true)
			ev, syncEvent := Event{}, Event{}
			w.dispatchQueuedEvent(&ev)
			synchronous.dispatchQueuedEvent(&syncEvent)
			if ev.DiplomacyReport == nil || ev.DiplomacyReport.AfterNotice != 1 || syncEvent.DiplomacyReport != nil {
				t.Fatal("外交官結果未掛到主要結果確認，或同步入口也被延後")
			}
			if w.Factions[0].Funds != 50000 || w.Factions[3].Funds != 50000 || !w.Friendship[0][3].AtWar() {
				t.Fatal("主要結果未確認就付款或停戰")
			}
			settlement := ev.DiplomacyReport.Confirm()
			_, released := settlement.resolve()
			if !bytes.Equal(w.Bytes(), synchronous.Bytes()) || !reflect.DeepEqual(released, syncEvent.ReleasedGenerals) {
				t.Fatal("確認後與同步原規則結果不同")
			}
			before := append([]byte(nil), w.Bytes()...)
			if ev.DiplomacyReport.Confirm() != nil || !bytes.Equal(before, w.Bytes()) {
				t.Fatal("重複確認再次結算")
			}
			if _, err := json.Marshal(ev); err != nil {
				t.Fatal("runtime續行進入事件序列化", err)
			}
		})
	}
}
