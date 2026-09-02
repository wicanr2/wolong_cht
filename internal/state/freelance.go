package state

import "github.com/wicanr2/wolong_cht/internal/rules/economy"

// 在野武將的每月處理。出處 `KI.EXE` 的 `sub_1585F` → `sub_15899`，
// 規格 docs/spec/114、機制 docs/mechanics/70 §3.9、欄位 docs/re/77 §2。

// affinityRollLimit 是兌現的亂數上限：`cmp al, 40h / jnb retn`，
// 也就是每月 0x40/0x100 ＝ **25%** 才會動。
const affinityRollLimit = 0x40

// recruitFreelanceGenerals 跑 `sub_15899` 的第一條路：**心向的勢力**。
//
// ⭐ 這一條是主規則，不是特例。開局在場的 81 名在野武將全部有 Affinity，
// 所以「投靠武將數最少的勢力」那條隨機投靠要等他們一一兌現之後才輪得到
// （docs/mechanics/70 §3.9）。隨機投靠那一條這一版還沒接，見 docs/spec/114 §5。
//
// 順序照原版：`sub_1585F` 在月結裡跑在 `sub_12BD9`（壓縮事件佇列）之前。
//
// 回傳這個月實際出仕的武將編號，給呼叫端做通知用；退場的不列入。
func (w *World) recruitFreelanceGenerals(rng economy.Rand) []int {
	if w == nil || rng == nil {
		return nil
	}
	var joined []int
	for id := range w.Generals {
		g := &w.Generals[id]
		// `cmp byte ptr [si], 80h / jb next`
		if !g.Alive {
			continue
		}
		// `cmp byte ptr [si+18h], 0 / jz + ; dec` —— 倒數沒歸零就只遞減。
		if g.Timer > 0 {
			g.Timer--
			continue
		}
		// 這一支只處理在野的。俘虜（Captor != noFaction）走
		// produceApproximateEvent10，兩者在原版是互斥的分支。
		if g.Faction != noFaction || g.Affinity == noFaction {
			continue
		}
		if int(rng.Next()&0xFF) >= affinityRollLimit {
			continue
		}
		want := g.Affinity
		// `mov byte ptr [si+19h], 0FFh` —— 不論後面走哪一條都先清掉。
		g.Affinity = noFaction
		if want >= 0 && want < numFactions && w.Factions[want].Alive {
			g.Faction = want
			// 原版 `loc_1591A` 的最後一步是 `sub_12AD2(al=勢力, ah=0FFh)`
			// ——勢力記錄 +0x18 的武將數要跟著加。漏掉的話那個數字會與
			// 實際人數脫節，而它是政略 AI 的輸入之一（aimarch.go 的註解）。
			w.raiseGeneralCount(want)
			joined = append(joined, id)
			continue
		}
		// 勢力已滅：旗標 bit 5 設著的整筆歸零，其餘留在原地。
		if g.VanishIfAffinityGone {
			g.Alive = false
			g.Posted = false
		}
	}
	return joined
}
