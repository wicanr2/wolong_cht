package state

import (
	"bytes"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
)

func TestDiplomacySettlementReleasesOneGeneralAtATime(t *testing.T) {
	w, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.Factions[0].Funds, w.Factions[3].Funds = 50000, 50000
	w.Friendship[0][3], w.Friendship[3][0] = diplomacy.Peace(40).WithWar(true), diplomacy.Peace(50).WithWar(true)
	for _, id := range []int{39, 45} {
		w.Generals[id].Faction, w.Generals[id].Captor, w.Generals[id].Duty = 3, 0, DutyCaptive
		w.Factions[0].Generals--
	}
	w.Generals[49].Faction, w.Generals[49].Captor, w.Generals[49].Duty = 0, 3, DutyCaptive
	w.Factions[3].Generals--
	synchronous, err := LoadBlock(w.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	synchronous.finishQueuedCeasefire(0, 3, 1, 10000)
	settlement := w.prepareQueuedCeasefire(0, 3, 1, 10000)
	if w.Factions[0].Funds != 40000 || w.Generals[39].Duty != DutyCaptive {
		t.Fatal("付款與第一名返回者的界線混在一起")
	}
	for _, id := range []int{39, 45, 49} {
		got, ok := settlement.NextReturnee()
		if !ok || got != id || w.Generals[id].Duty != DutyNone || w.Generals[id].Captor != noFaction {
			t.Fatal("返回記錄順序或狀態錯誤", got, id)
		}
		for _, later := range []int{39, 45, 49} {
			if later > id && w.Generals[later].Duty != DutyCaptive {
				t.Fatal("前一人對話未結束就釋放後一人", later)
			}
		}
		if settlement.Complete() || !w.Friendship[0][3].AtWar() {
			t.Fatal("返回者尚未耗盡就完成停戰")
		}
	}
	if _, ok := settlement.NextReturnee(); ok || !settlement.Complete() || settlement.Complete() {
		t.Fatal("最後收尾必須只執行一次")
	}
	if !bytes.Equal(w.Bytes(), synchronous.Bytes()) {
		t.Fatal("逐人完成後與同步結果不同")
	}
}

func TestDiplomacyDeclarationWaitsForConfirmation(t *testing.T) {
	w, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.Friendship[0][12], w.Friendship[12][0] = diplomacy.Peace(40), diplomacy.Peace(50)
	settlement := w.prepareQueuedCooperation(0, 12, 2, 0, 0)
	for {
		if _, ok := settlement.NextReturnee(); !ok {
			break
		}
	}
	if notice := settlement.CompletionNotice(); notice == nil || notice.Index != 0x1A0 || w.Friendship[0][12].AtWar() {
		t.Fatal("宣戰對話顯示前已交戰，或缺少原版宣戰句")
	}
	if !settlement.Complete() || !w.Friendship[0][12].AtWar() || !w.Friendship[12][0].AtWar() {
		t.Fatal("確認宣戰後未寫雙向狀態")
	}
}

func TestDiplomacyDeclarationNoticeUsesExistingTargetGate(t *testing.T) {
	w, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.Factions[0].InvasionTarget = 3
	w.Friendship[0][12], w.Friendship[12][0] = diplomacy.Peace(40), diplomacy.Peace(50)
	settlement := w.prepareQueuedCooperation(0, 12, 2, 0, 0)
	for {
		if _, ok := settlement.NextReturnee(); !ok {
			break
		}
	}
	if settlement.CompletionNotice() != nil || settlement.Complete() || w.Friendship[0][12].AtWar() {
		t.Fatal("已有侵攻目標時仍顯示宣戰或改寫交戰狀態")
	}
}
