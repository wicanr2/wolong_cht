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
	return nil, fmt.Errorf("-shot-when %q 不認得；可用的是 battle、battle-frame:N、gate-bar", s)
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
