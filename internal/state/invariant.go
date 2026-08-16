package state

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
)

// 資料模型的不變量。
//
// **這一層驗的是「規則組合起來對不對」，不是「單條規則對不對」。**
// 單元測試把每一條公式釘在反組譯上，但那只保證單次呼叫；
// 佔領、壞滅、招降、陣亡、編成、解散這些事互相牽動的欄位跑久了會不會歪，
// 單元測試看不出來。
//
// 用的準繩是**原版檔案自己維護的冗餘**：同一件事在兩個地方各存一份，
// 只要一致就代表兩邊的更新路徑都對。這些冗餘不是我發明的——
// 是掃過四個劇本、192 個據點、127 名武將驗出來的（docs/formats/08 §1.4）：
//
//	據點表數出來的各勢力據點數  ==  勢力記錄 +0x23
//	武將表數出來的各勢力武將數  ==  勢力記錄 +0x18
//
// ⚠ **據點那一條在劇本 3、4 開局就對不上**，而那是**原版自己的資料瑕疵**：
// 武陵與南昌的 `+0x01`（執行期讀的）與 `+0x1A`（作者填的）互相矛盾，
// 而 `+0x23` 是照 `+0x1A` 算的（docs/formats/08 §1.6）。本專案照抄執行期
// 的 `+0x01`，所以那個差額一開始就在。
//
// 所以這裡驗的不是「相等」，而是**「差額不會變」**——
// 那反而更嚴格：任何一條規則改了據點歸屬卻忘了更新計數（或反過來），
// 差額就會動。基準線在載入時記下來（`bias`）。
//
// 另外幾條是資料本身的定義域（城兵 ≤ 上限、軍團兵力 == 六槽之和…），
// 違反就表示某條規則寫壞了欄位。

// Violation 是一條被打破的不變量。
type Violation struct {
	Kind   string // 哪一條
	Detail string
}

func (v Violation) String() string { return v.Kind + "：" + v.Detail }

// snapshotBias 在載入完成時記下「勢力記錄的據點數」與「實際數出來的」
// 之間的差額。原版劇本 3、4 開局就有差（見上），那不是我們的錯，
// 但**之後不准再變**。
func (w *World) snapshotBias() {
	for f := range w.cityBias {
		w.cityBias[f] = 0
	}
	cities := map[int]int{}
	for i := range w.Cities {
		if o := w.Cities[i].Owner; o >= 0 && o < numFactions {
			cities[o]++
		}
	}
	for f := 0; f < numFactions; f++ {
		w.cityBias[f] = w.Factions[f].Cities - cities[f]
	}
}

// CheckInvariants 掃一遍世界，回傳所有被打破的不變量。
//
// 回傳空切片表示這一刻的狀態是自洽的。
func (w *World) CheckInvariants() []Violation {
	var out []Violation
	add := func(kind, format string, a ...any) {
		out = append(out, Violation{Kind: kind, Detail: fmt.Sprintf(format, a...)})
	}

	// ① 據點數：勢力記錄的計數要等於實際數出來的。
	cities := map[int]int{}
	for i := range w.Cities {
		c := &w.Cities[i]
		if c.Owner >= 0 && c.Owner < numFactions {
			cities[c.Owner]++
		}
		if c.Garrison > c.GarrisonCap {
			add("城兵超過上限", "據點 %d 城兵 %d > 上限 %d", i, c.Garrison, c.GarrisonCap)
		}
		if c.Garrison < 0 || c.Prevention < 0 {
			add("據點欄位為負", "據點 %d 城兵 %d、防災 %d", i, c.Garrison, c.Prevention)
		}
	}

	// ② 武將數：同上。
	generals := map[int]int{}
	for i := range w.Generals {
		g := &w.Generals[i]
		if !g.Alive {
			continue
		}
		// **俘虜不算在任何一方的武將數裡。** 原版 `sub_12AD2` 在被俘時
		// 只 `dec` 舊勢力，沒有 `inc` 新勢力；要等釋放（`sub_150D7`）
		// 才重新入帳。被俘期間 +0x1C 雖然指向俘虜方，帳上卻是空的。
		if g.Captor != noFaction {
			continue
		}
		if g.Faction >= 0 && g.Faction < numFactions {
			generals[g.Faction]++
		}
	}

	// ③ 軍團數，以及「軍團活著 ⟺ 帶兵的武將出陣中」。
	corps := map[int]int{}
	for i := range w.Corps {
		c := &w.Corps[i]
		if !c.Alive {
			continue
		}
		if c.Faction >= 0 && c.Faction < numFactions {
			corps[c.Faction]++
		}
		total := 0
		for _, u := range c.Units {
			if u.Men < 0 {
				add("部隊兵力為負", "軍團 %d 有一槽是 %d", i, u.Men)
			}
			total += u.Men
		}
		if total != c.Men {
			add("軍團兵力與六槽不符", "軍團 %d 記著 %d，六槽加起來 %d", i, c.Men, total)
		}
		if g := &w.Generals[w.Leader(i)]; !g.Posted {
			add("帶兵的武將沒標出陣", "軍團 %d 的武將 %d", i, w.Leader(i))
		} else if !g.Alive {
			add("帶兵的武將不存在", "軍團 %d 的武將 %d", i, w.Leader(i))
		}
		if army.KindOf(c.Node) == army.CityNode &&
			w.Cities[c.Node].Owner != c.Faction &&
			w.Cities[c.Node].Owner != combat.NeutralFaction {
			// 停在敵方據點上是暫態（下一 tick 就會打起來），不算違規；
			// 這裡只記在旁邊，不列入違反。
			_ = c
		}
	}

	for f := 0; f < numFactions; f++ {
		fa := &w.Factions[f]
		if !fa.Alive {
			// 已滅的勢力不該還有據點或武將。
			if cities[f] != 0 {
				add("已滅勢力還有據點", "勢力 %d 還有 %d 個", f, cities[f])
			}
			if corps[f] != 0 {
				add("已滅勢力還有軍團", "勢力 %d 還有 %d 支", f, corps[f])
			}
			continue
		}
		if got, want := fa.Cities-cities[f], w.cityBias[f]; got != want {
			add("據點數的差額變了",
				"勢力 %d 記著 %d、實際 %d（差 %d，開局是 %d）",
				f, fa.Cities, cities[f], got, want)
		}
		if fa.Generals != generals[f] {
			add("武將數不符", "勢力 %d 記著 %d，實際 %d", f, fa.Generals, generals[f])
		}
		if fa.Corps != corps[f] {
			add("軍團數不符", "勢力 %d 記著 %d，實際 %d", f, fa.Corps, corps[f])
		}
		for t, n := range fa.Reserves {
			if n < 0 {
				add("預備兵為負", "勢力 %d 的兵種 %d 是 %d", f, t, n)
			}
		}
		if fa.Alive && fa.Lord >= 0 && fa.Lord < len(w.Generals) {
			if l := &w.Generals[fa.Lord]; !l.Alive {
				add("君主不存在", "勢力 %d 的君主 %d", f, fa.Lord)
			} else if l.Faction != f {
				add("君主不屬於自己的勢力", "勢力 %d 的君主 %d 屬於 %d", f, fa.Lord, l.Faction)
			}
		}
	}
	return out
}
