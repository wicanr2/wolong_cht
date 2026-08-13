// Package strategyai 提供政略 AI 的純規則判定。
//
// 這裡只放不依賴存檔或畫面的資料流：鄰接勢力按交友度排序、國力評分、
// 以及原版 sub_12EFB 的三道侵攻閘。World 另外負責把這些判定接到月結、
// 交戰狀態與軍團資料；如此可用固定輸入直接測試機器碼規則。
package strategyai

import (
	"sort"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
)

// NoTarget 是勢力記錄 +0x19 的無侵攻目標哨兵值。
const NoTarget = 0xFF

// Faction 是 sub_13091／sub_12EFB 會讀到的勢力欄位。
// Funds 保留原始有號 24 位資金的數值；Power 會按原版的 24 位元組讀法
// 處理負數的高兩 byte。
type Faction struct {
	Alive          bool
	Player         bool
	Cities         int
	Funds          int
	Reserves       [3]int
	Aggression     int
	InvasionTarget int
}

// Candidate 是一個已由據點鄰接槽去重的敵對勢力候選。
type Candidate struct {
	Faction    int
	Friendship diplomacy.Friendship
}

// SortCandidates 複製並按原版 sub_12C52 的交友度由低到高排序。
// 同值時以勢力編號作穩定 tie-break；原版的選擇排序遇同值會保留掃描順序，
// 而據點掃描本身就是編號順序，這裡把該結果明寫出來。
func SortCandidates(in []Candidate) []Candidate {
	out := append([]Candidate(nil), in...)
	sort.SliceStable(out, func(i, j int) bool {
		li, lj := out[i].Friendship.Raw(), out[j].Friendship.Raw()
		if li != lj {
			return li < lj
		}
		return out[i].Faction < out[j].Faction
	})
	return out
}

// Power 重現 sub_13091 的已讀資料流。
//
// 先把三種預備兵各右移 2 位後相加；若結果的高 byte 不小於據點數，
// 原版直接把分數拉到 0x7D0，之後再做 0x7D0 上限。資金記錄 +0x21
// 是原始 24 位元組的高兩 byte，以 unsigned word 與 0x13 比較；太低時
// 國力直接為零。這個細節不能寫成 Go 的算術右移，否則負資金會分岔。
func Power(f Faction) int {
	p := (f.Reserves[0] >> 2) + (f.Reserves[1] >> 2) + (f.Reserves[2] >> 2)
	if p>>8 >= f.Cities {
		p = 0x7D0
	}
	if p > 0x7D0 {
		p = 0x7D0
	}
	if fundsWord(f.Funds) <= 0x13 {
		return 0
	}
	return p
}

func fundsWord(funds int) int {
	u := uint32(funds) & 0xFFFFFF
	return int((u >> 8) & 0xFFFF)
}

// FundLimit 是 sub_12EFB 的侵攻資金門檻（比較用的原始 word）。
func FundLimit(cities int) int {
	v := cities*16 + 64
	if v > 0x61A {
		return 0x61A
	}
	return v
}

// FriendshipLimit 是 sub_12EFB 的交友度門檻，含和平位元。
func FriendshipLimit(aggression int) int {
	if aggression < 0 {
		aggression = 0
	}
	return 0x80 + 20 + aggression + aggression/2
}

// ShouldDeclareWar 重現 sub_12EFB。成功條件全部採原版的嚴格比較：
// 資金必須大於門檻、交友度不得大於門檻、己方國力至少是目標的四分之三。
func ShouldDeclareWar(self, target Faction, candidate Candidate) bool {
	if !self.Alive || !target.Alive || self.Player || self.InvasionTarget == candidate.Faction {
		return false
	}
	if fundsWord(self.Funds) <= FundLimit(self.Cities) {
		return false
	}
	if candidate.Friendship.Raw() > FriendshipLimit(self.Aggression) {
		return false
	}
	selfPower, targetPower := Power(self), Power(target)
	return selfPower >= targetPower-targetPower>>2
}

// AtWarValue 是 sub_13639 的雙向交戰友好度更新：取兩方原始值的較小
// 7-bit 數值，再右移一位，最高位保持清除。
func AtWarValue(a, b diplomacy.Friendship) diplomacy.Friendship {
	v := a.Value()
	if b.Value() < v {
		v = b.Value()
	}
	return diplomacy.Friendship(byte(v >> 1)).WithWar(true)
}
