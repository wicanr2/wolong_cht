package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

// ⚠ 退兵回池要夾 `MaxReserve`（原版 `sub_155EC` 的 0xFFDC，
// 而那個上限正是在**退兵這條路徑**上驗到的，docs/spec/21 §5）。
// remake 先前只有月結加兵那一邊夾，退兵這份漏掉。
func TestPoolBackClampsToMaxReserve(t *testing.T) {
	w := &World{}
	w.Corps[0] = Corps{Alive: true, Faction: 0}
	w.Corps[0].Units[0] = combat.Unit{Kind: army.Cavalry, Men: 5000}
	w.Factions[0].Reserves[army.Cavalry] = economy.MaxReserve - 100

	w.poolBack(0)

	if got := w.Factions[0].Reserves[army.Cavalry]; got != economy.MaxReserve {
		t.Errorf("退兵後預備兵 = %d，應該夾在 %d", got, economy.MaxReserve)
	}
	// 槽位要清空、總兵力歸零——這是 poolBack 原本就有的行為，不可以弄壞。
	if w.Corps[0].Units[0].Men != 0 || w.Corps[0].Men != 0 {
		t.Error("退兵之後槽位或總兵力沒有歸零")
	}
}

// 沒到上限時照實加，不可以被夾成別的值。
func TestPoolBackAddsExactlyBelowCap(t *testing.T) {
	w := &World{}
	w.Corps[0] = Corps{Alive: true, Faction: 0}
	w.Corps[0].Units[0] = combat.Unit{Kind: army.Infantry, Men: 1234}
	w.Factions[0].Reserves[army.Infantry] = 1000

	w.poolBack(0)

	if got := w.Factions[0].Reserves[army.Infantry]; got != 2234 {
		t.Errorf("退兵後預備兵 = %d，預期 2234", got)
	}
}
