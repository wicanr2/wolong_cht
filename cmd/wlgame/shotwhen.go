package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// `-shot-when` 的局面條件（docs/spec/118）。
//
// ⭐ **對拍的取樣點是規則層的函數，不是常數。** 寫死一個步數，規則層一改
// 就落在別的局面上，而**沒有任何測試會紅**——那些數字只寫在文件裡
// （docs/spec/91 §6）。這一組讓驗收腳本寫得出「等到門強度條亮著再截」。

// shotCondition 是一個具名的局面條件。
type shotCondition struct {
	name  string
	check func(*game) bool
}

// parseShotWhen 解析 `-shot-when` 的值。留白回 nil（照 `-shot-frames`）。
//
// ⛔ **不認得的值一律回錯誤，不要當成留白。** 打錯字時靜靜退回「照幀數截圖」
// 會拍到別的局面，而輸出看起來完全正常——那正是這一組旗標要消滅的失敗模式。
func parseShotWhen(s string) (*shotCondition, error) {
	s = strings.TrimSpace(s)
	switch {
	case s == "":
		return nil, nil
	case s == "battle":
		return &shotCondition{name: s, check: func(g *game) bool {
			return g.shotBattle() != nil
		}}, nil
	case s == "gate-bar":
		// docs/spec/91 §6 那張表的第二列：門強度條顯示中。
		return &shotCondition{name: s, check: func(g *game) bool {
			b := g.shotBattle()
			if b == nil {
				return false
			}
			_, shown := b.StructureBar()
			return shown
		}}, nil
	case s == "battle-settled":
		// ⭐ **開場布陣走完的那一刻**（docs/spec/91 §5.5）。原版把兵擺在
		// 戰場邊界再讓他們走進陣形（docs/spec/133），走完之後有一小段
		// 誰都不動的空窗，接著腳本才下第一道命令。
		// 實測兩邊都是**第 70 拍站定、第 75/76 拍重新開動**——
		// 那一段是戰術對拍唯一「兩邊都靜止」的取樣窗。
		//
		// ⚠ **判準是「這一拍沒有人動」，不是「座標範圍不再變」**：
		// 範圍在大部分人到位之後就不動了，而最後幾個還在走。
		// ⚠ **要跨「戰術拍」比，不能跨「呼叫次數」比。** 條件每個畫面幀
		// 都會被問一次，而戰鬥不是每一幀都走一拍——連問兩次看到同一組
		// 座標是常態，於是第 2 拍就「成立」了（實測就是這樣）。
		var prev map[int][2]int
		prevFrame := -1
		return &shotCondition{name: s, check: func(g *game) bool {
			b := g.shotBattle()
			if b == nil {
				return false
			}
			if b.Frame == prevFrame {
				return false
			}
			prevFrame = b.Frame
			now := map[int][2]int{}
			for side := range b.Sides {
				for k := range b.Sides[side].Soldiers {
					u := &b.Sides[side].Soldiers[k]
					if u.Alive {
						now[side*1000+k] = [2]int{u.X, u.Y}
					}
				}
			}
			settled := prev != nil && len(prev) == len(now)
			if settled {
				for k, v := range now {
					if prev[k] != v {
						settled = false
						break
					}
				}
			}
			prev = now
			// 開場的頭幾拍大家都還沒起步，也是「沒有人動」——
			// 所以要等真的走過一段才算數。
			return settled && b.Frame > 10
		}}, nil
	case strings.HasPrefix(s, "clock:"):
		// `clock:年/月/日[/時]`——**狀態層對拍的取樣點**（docs/spec/138）。
		// 兩邊比表要在同一個遊戲時刻，否則內政每小時都在動幾個據點，
		// 差異看起來像規則分歧，其實只是取樣點差了幾拍。
		f := strings.Split(strings.TrimPrefix(s, "clock:"), "/")
		if len(f) != 3 && len(f) != 4 {
			return nil, fmt.Errorf("-shot-when %q：`clock:` 後面要 `年/月/日` 或 `年/月/日/時`", s)
		}
		want := make([]int, len(f))
		for i, part := range f {
			n, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("-shot-when %q：`%s` 不是整數", s, part)
			}
			want[i] = n
		}
		return &shotCondition{name: s, check: func(g *game) bool {
			if g.world == nil {
				return false
			}
			c := g.world.Clock
			if c.Year != want[0] || c.Month != want[1] || c.Day != want[2] {
				return false
			}
			return len(want) == 3 || c.Hour == want[3]
		}}, nil
	case strings.HasPrefix(s, "battle-frame:"):
		n, err := strconv.Atoi(strings.TrimPrefix(s, "battle-frame:"))
		if err != nil || n < 0 {
			return nil, fmt.Errorf("-shot-when %q：`battle-frame:` 後面要一個非負整數", s)
		}
		return &shotCondition{name: s, check: func(g *game) bool {
			b := g.shotBattle()
			return b != nil && b.Frame >= n
		}}, nil
	}
	return nil, fmt.Errorf("-shot-when %q 不認得；可用的是 battle、battle-frame:N、battle-settled、gate-bar", s)
}

// shotBattle 取目前開著的戰術戰鬥，沒有就回 nil。
func (g *game) shotBattle() *tactical.Battle {
	if g == nil || g.world == nil {
		return nil
	}
	p := g.world.PendingBattle()
	if p == nil {
		return nil
	}
	return p.Battle
}

// shotReady 回報「現在可以截了沒」。沒帶 `-shot-when` 時恆真——
// 那條路完全照舊，由 `-shot-frames` 決定。
func (g *game) shotReady() bool {
	if g == nil || g.shotWhen == nil {
		return true
	}
	return g.shotWhen.check(g)
}
