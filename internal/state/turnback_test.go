package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 邊界掉頭（`docs/spec/132`）：軍團走進**和平**勢力的據點之前會折返。
//
// ⚠ 原版**不在下令時檢查**——指令下得成、訊息也跳，折返發生在路上。
// 所以這三條都是「下令成功之後跑一段時間，看它停在哪」。
// setupMarch 編一支軍團在 from，下令走到 to，跑 n 個 tick，回傳軍團編號。
func setupMarch(t *testing.T, w *World, from, to, ticks int) int {
	t.Helper()
	owner := w.Cities[from].Owner
	// ⚠ **要當玩家的軍團**，否則 AI 每隔一陣子就把它重新派走，
	// 量到的就不是掉頭而是 AI 的決策。
	w.Player = owner
	w.Factions[owner].Reserves = [economy.NumTroopTypes]int{9000, 9000, 9000}
	lord := w.Factions[owner].Lord
	kinds := [army.Positions]army.TroopType{}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err != nil {
		t.Fatal(err)
	}
	c := &w.Corps[lord]
	c.Node, c.X, c.Y = from, w.Cities[from].X, w.Cities[from].Y
	if err := w.March(lord, to); err != nil {
		t.Fatal(err)
	}
	r := rng.NewFixed(3)
	for i := 0; i < ticks; i++ {
		w.Tick(r)
	}
	return lord
}

// borderPair 找一對相鄰的據點：`from` 屬於某個活著的勢力，
// `to` 屬於**另一個非中立**的勢力。
//
// ⚠ **不要只看首都的鄰居。** 劇本 1 的曹操首都許昌四個鄰居全是自己的，
// 於是三條測試會一起 Skip——而 Skip 與 PASS 在輸出上一樣安靜。
func borderPair(w *World) (from, to int) {
	for i := range w.Cities {
		me := w.Cities[i].Owner
		if me == combat.NeutralFaction || me < 0 || me >= len(w.Friendship) {
			continue
		}
		for _, n := range w.Cities[i].Neighbours {
			if n < 0 || n >= len(w.Cities) {
				continue
			}
			o := w.Cities[n].Owner
			if o != me && o != combat.NeutralFaction &&
				o >= 0 && o < len(w.Friendship) {
				return i, n
			}
		}
	}
	return -1, -1
}

// TestMarchTurnsBackAtPeacefulBorder 是這條規則的正面樣本。
func TestMarchTurnsBackAtPeacefulBorder(t *testing.T) {
	w := load(t, 0)
	from, to := borderPair(w)
	if to < 0 {
		t.Fatal("整張圖找不到一對跨勢力的相鄰據點")
	}
	me, other := w.Cities[from].Owner, w.Cities[to].Owner
	w.Friendship[me][other] = w.Friendship[me][other].WithWar(false)

	lord := setupMarch(t, w, from, to, 4000)
	c := &w.Corps[lord]
	if c.Node == to {
		t.Fatalf("和平狀態下軍團仍然走進了據點 %d", to)
	}
	// 掉頭之後目標會被改寫成折返的據點——不是停在半路。
	if c.TargetNode == to {
		t.Errorf("目標仍是 %d，掉頭沒有改寫目標", to)
	}
	if c.Node != from || c.X != w.Cities[from].X || c.Y != w.Cities[from].Y {
		t.Errorf("軍團停在 node=%d (%d,%d)，應折返回 %d (%d,%d)",
			c.Node, c.X, c.Y, from, w.Cities[from].X, w.Cities[from].Y)
	}
}

// TestMarchPassesWhenAtWar 是反面樣本：只把和平位元清掉，其餘不變。
//
// **兩條要是同一個局面**，否則證明不了「差別來自交戰與否」。
func TestMarchPassesWhenAtWar(t *testing.T) {
	w := load(t, 0)
	from, to := borderPair(w)
	if to < 0 {
		t.Fatal("整張圖找不到一對跨勢力的相鄰據點")
	}
	me, other := w.Cities[from].Owner, w.Cities[to].Owner
	w.Friendship[me][other] = w.Friendship[me][other].WithWar(true)

	lord := setupMarch(t, w, from, to, 4000)
	if c := &w.Corps[lord]; c.Node == from &&
		c.X == w.Cities[from].X && c.Y == w.Cities[from].Y {
		t.Fatal("交戰中還是掉頭了")
	}
}

// TestMarchPassesIntoNeutral 釘住中立那一條例外（`sub_142AB` 的 `cmp al, 18h`）。
func TestMarchPassesIntoNeutral(t *testing.T) {
	w := load(t, 0)
	from, to := borderPair(w)
	if to < 0 {
		t.Fatal("整張圖找不到一對跨勢力的相鄰據點")
	}
	me := w.Cities[from].Owner
	// 把那個鄰居改成中立，並且與它「和平」——只有中立那一條能讓它通過。
	w.Cities[to].Owner = combat.NeutralFaction
	for i := range w.Friendship[me] {
		w.Friendship[me][i] = w.Friendship[me][i].WithWar(false)
	}

	lord := setupMarch(t, w, from, to, 4000)
	if c := &w.Corps[lord]; c.Node == from &&
		c.X == w.Cities[from].X && c.Y == w.Cities[from].Y {
		t.Fatal("中立據點也掉頭了")
	}
}
