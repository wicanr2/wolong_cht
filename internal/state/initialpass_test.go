package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 新遊戲的開局政略評估（docs/spec/83）：sub_12BD9 的第二個呼叫點在
// sub_11AC3+66（選定君主當下），讀檔路徑跳過。孫策開局資金 word 82
// 只比財政閘 80 多 2，第一次月結後就掉下去——缺這一次評估，
// 它永遠不會對劉繇宣戰（playtest/45 的長程分歧根因）。
func TestInitialStrategyPassQueuesOpeningDeclarations(t *testing.T) {
	const sunCe, liuYao = 1, 10

	w := load(t, 0)
	w.Player = 0
	w.EnableStrategicAI()
	r := rng.NewFixed(1)
	w.RunInitialStrategyPass(r)

	// 宣戰走事件 1 佇列（Code 低 byte 1、來源孫策、Param 低 byte 目標）。
	found := false
	for _, e := range w.events {
		if e.Code&0xFF == 1 && int(e.Code>>8) == sunCe && int(e.Param&0xFF) == liuYao {
			found = true
		}
	}
	if !found {
		t.Fatal("開局評估後佇列裡沒有孫策對劉繇的事件 1")
	}

	// dispatch 之後交戰成立。跑到宣戰生效或超時。
	for i := 0; i < 50000 && !w.Friendship[sunCe][liuYao].AtWar(); i++ {
		if w.PendingEncounter() != nil {
			w.ChooseBattleDelegate(r)
			continue
		}
		if w.PendingDiplomacy() != nil {
			w.ResolveDiplomacy(DiplomacyReject)
			continue
		}
		if w.PendingFunding() != nil {
			w.ResolveFunding(FundingReject)
			continue
		}
		w.Tick(r)
	}
	if !w.Friendship[sunCe][liuYao].AtWar() {
		t.Fatal("事件 1 dispatch 後孫策與劉繇沒有進入交戰")
	}

	// 對照組：沒開政略 AI 時不做任何事。
	w2 := load(t, 0)
	w2.RunInitialStrategyPass(rng.NewFixed(1))
	for _, e := range w2.events {
		if e.Code != 0 {
			t.Fatal("未啟用政略 AI 卻排入了事件")
		}
	}
}
