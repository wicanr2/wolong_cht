package state

import (
	"fmt"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func neutralDeclarationTestWorld(t *testing.T, scenario, power, alias int) *World {
	t.Helper()
	w := cooperationTestWorld(t, scenario)
	b := w.Bytes()
	off := factionBase + 24*factionSize
	for i := 0; i < 3; i++ {
		putU16(b, off+4+i*2, 0)
	}
	putU16(b, off+4, alias*4)
	putI24(b, off+0x20, 8192)
	b[off+0x23] = 8
	w, err := LoadBlock(b)
	if err != nil {
		t.Fatal(err)
	}
	w.Factions[1].Funds, w.Factions[1].Cities = 500_000, 8
	w.Factions[1].Reserves = [3]int{power * 4, 0, 0}
	w.Factions[1].InvasionTarget = diplomacy.NoTarget
	w.Factions[2].InvasionTarget = 24
	w.Friendship[1][2], w.Friendship[2][1] = diplomacy.Friendship(30), diplomacy.Friendship(40)
	return w
}

// 原版的等強分支保留舊目標。24會讀交友度表，不直接轉攻宣戰者。
func TestDeclarationNeutralTargetPowerBoundaries(t *testing.T) {
	for scenario := 0; scenario < 4; scenario++ {
		for _, power := range []int{99, 2000} {
			for _, alias := range []int{0, 99, 100, 2000} {
				t.Run(fmt.Sprintf("s%d-p%d-a%d", scenario, power, alias), func(t *testing.T) {
					w := neutralDeclarationTestWorld(t, scenario, power, alias)
					if !w.applyQueuedDeclaration(1, 2) {
						t.Fatal("合法宣戰被拒絕")
					}
					want := 24
					if power > alias {
						want = 1
					}
					if w.Factions[2].InvasionTarget != want || w.Factions[1].InvasionTarget != 2 ||
						w.Friendship[1][2].Raw() != 15 || w.Friendship[2][1].Raw() != 15 {
						t.Fatalf("宣戰結果錯誤：source=%d target=%d，want target=%d", w.Factions[1].InvasionTarget, w.Factions[2].InvasionTarget, want)
					}
					restored, err := LoadBlock(w.Bytes())
					if err != nil || restored.Factions[2].InvasionTarget != want {
						t.Fatal("宣戰反應未保存", err)
					}
				})
			}
		}
	}
}

func TestDeclarationNeutralTargetReadsCurrentRelations(t *testing.T) {
	w := neutralDeclarationTestWorld(t, 0, 99, 0)
	// 載入時的別名國力為零，之後交友度更新為8000預備兵的等價bytes。
	w.Friendship[0][4], w.Friendship[0][5] = diplomacy.Friendship(0x40), diplomacy.Friendship(0x1F)
	if !w.applyQueuedDeclaration(1, 2) || w.Factions[2].InvasionTarget != 24 {
		t.Fatal("使用舊raw快照，漏讀目前交友度")
	}
}

func TestNeutralDeclarationProducerThenIncomingWar(t *testing.T) {
	w := neutralDeclarationTestWorld(t, 0, 99, 2000)
	w.Factions[2].InvasionTarget = diplomacy.NoTarget
	w.Factions[2].Funds, w.Factions[2].Cities = 500_000, 1
	w.Cities[0].Owner, w.Cities[0].Adjacency = 2, 1
	w.Cities[0].Neighbours[0], w.Cities[1].Owner = 1, 24
	w.events = [eventQueueEntries]QueuedEvent{}
	var out []StrategyEvent
	if !w.queueNeutralDeclaration(rng.NewFixed(7), 2, &out) {
		t.Fatal("接壤中立且資金足夠時，沒有排入原版中立宣戰")
	}
	dispatch := func(code uint16) {
		t.Helper()
		for i, e := range w.events {
			if e.Code == code {
				w.eventCursor, w.eventDelay = i*eventQueueEntrySize, 1
				w.dispatchQueuedEvent(&Event{})
				return
			}
		}
		t.Fatalf("找不到事件%04X", code)
	}
	dispatch(0x0201)
	if w.Factions[2].InvasionTarget != 24 {
		t.Fatal("原版producer未經事件1設成中立目標")
	}
	if !w.queueEvent(rng.NewFixed(7), 1, 1, 0xFF02, 0xFF) {
		t.Fatal("沒有排入後續宣戰")
	}
	dispatch(0x0101)
	if w.Factions[2].InvasionTarget != 24 || w.Factions[1].InvasionTarget != 2 {
		t.Fatal("中立目標在後續宣戰被錯誤覆寫")
	}
}

func TestDeclarationAliasBoundsWithoutRawBlock(t *testing.T) {
	w := &World{Player: 0}
	w.Factions[1].Alive, w.Factions[2].Alive = true, true
	w.Factions[1].InvasionTarget, w.Factions[2].InvasionTarget = 255, 24
	if _, ok := w.declarationTargetPower(24); ok {
		t.Fatal("沒有原始區塊，卻嘗試讀別名")
	}
	if !w.applyQueuedDeclaration(1, 2) || w.Factions[2].InvasionTarget != 1 {
		t.Fatal("人工World的既有安全邊界改變")
	}
	for _, invalid := range []int{-1, 36, 255, 256} {
		if _, ok := w.declarationTargetPower(invalid); ok {
			t.Fatalf("越界索引%d仍讀取別名", invalid)
		}
	}
}
