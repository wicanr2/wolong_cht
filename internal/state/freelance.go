package state

import "github.com/wicanr2/wolong_cht/internal/rules/economy"

// 在野武將的每月處理。出處 `KI.EXE` 的 `sub_1585F` → `sub_15899`，
// 規格 docs/spec/114、機制 docs/mechanics/70 §3.9、欄位 docs/re/77 §2。

// affinityRollLimit 是兌現的亂數上限：`cmp al, 40h / jnb retn`，
// 也就是每月 0x40/0x100 ＝ **25%** 才會動。
const affinityRollLimit = 0x40

// freelanceJoinTalk 是出仕通知的 TALK 索引（`sub_15899` 的 `mov cx, 29h`）：
// #41「{1}加入麾下了。」⭐ **只有投靠玩家時才發**——原版的
// `cmp bx, cs:word_10CFD / jnz` 擋在前面，而那一段是心向與隨機投靠共用的。
const freelanceJoinTalk = 0x29

// recruitFreelanceGenerals 跑 `sub_15899` 的兩條路。
//
// ⭐ **有心向的每月 25% 才動，沒有心向的每月都跑**——原版的
// `cmp bh, 0FFh / jz loc_158C2` 在擲那顆 25% 的骰**之前**，
// 所以隨機投靠不受它限制（docs/spec/130）。
//
// 開局在場的 81 名在野武將全部有 Affinity，所以隨機投靠那一條要等
// 他們一一兌現之後才輪得到（docs/mechanics/70 §3.9）。
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
		if g.Faction != noFaction {
			continue
		}
		// ⭐ 沒有心向 ⇒ 隨機投靠那一條，**不經過 25% 的閘**
		// （`cmp bh, 0FFh / jz loc_158C2` 在擲骰之前，docs/spec/130）。
		if g.Affinity == noFaction {
			if to := w.randomJoin(id, rng); to >= 0 {
				joined = append(joined, id)
			}
			continue
		}
		if int(rng.Next()&0xFF) >= affinityRollLimit {
			continue
		}
		want := g.Affinity
		// `mov byte ptr [si+19h], 0FFh` —— 不論後面走哪一條都先清掉。
		g.Affinity = noFaction
		if want >= 0 && want < numFactions && w.Factions[want].Alive {
			// ⭐ 與隨機投靠共用 `loc_1591A`：寫勢力 ＋ `sub_12AD2` 入帳。
			// 勢力記錄 +0x18 的武將數要跟著加，漏掉的話那個數字會與實際
			// 人數脫節，而它是政略 AI 的輸入之一（aimarch.go 的註解）。
			w.joinFaction(id, want)
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

// fewestGeneralsFaction 是 `sub_15899` 的掃描段（`loc_158CA`）：
// 22 個勢力裡**存在且武將數最少**的那一個。
//
// ⚠ 平手取**編號小的**——原版是 `cmp ah, [di+18h] / jb` 嚴格小於，
// 相等時不更新。這一點是從指令推的，沒有實機驗過（docs/spec/130 §5）。
// 一個勢力都不存在時回 0，與原版的 `bx` 初值相同。
func (w *World) fewestGeneralsFaction() int {
	best, fewest := 0, 0xFF
	for i := 0; i < numFactions; i++ {
		if !w.Factions[i].Alive {
			continue
		}
		if w.Factions[i].Generals < fewest {
			fewest = w.Factions[i].Generals
			best = i
		}
	}
	return best
}

// randomJoinRelief 是骰面 ≥ 48 那一條（`loc_15907`）：**玩家專屬的救濟**。
// 據點數 ÷ 4 ＋ 1 > 武將數才送一個過去，否則這個月什麼都不做。
const randomJoinReliefRoll = 0x30

// randomJoinFlatRoll 是「壓成 1」的下界（`cmp al, 18h`）：
// 骰面 24–47 一律當成 1 ⇒ 投靠武將最少的那一個。
const randomJoinFlatRoll = 0x18

// randomJoin 跑 `sub_15899` 的 `+0x19 == 0xFF` 分支（docs/spec/130）。
// 回傳投靠的勢力編號，−1 ＝ 這個月沒有動。
func (w *World) randomJoin(id int, rng economy.Rand) int {
	start := w.fewestGeneralsFaction()

	// `call sub_1ECE0 / and al, 3Fh / inc al` ⇒ 1..64
	roll := int(rng.Next()&0x3F) + 1

	if roll >= randomJoinReliefRoll {
		p := w.Player
		if p < 0 || p >= numFactions || !w.Factions[p].Alive {
			return -1
		}
		// `mov al,[bx+23h] / shr ×2 / inc al / cmp al,[bx+18h] / jbe`
		if w.Factions[p].Cities/4+1 <= w.Factions[p].Generals {
			return -1
		}
		return w.joinFaction(id, p)
	}

	n := roll
	if n >= randomJoinFlatRoll {
		n = 1
	}
	// `loc_158F1`：從 start 起數第 n 個「存在」的勢力，走到底繞回 0。
	//
	// ⚠ **保險絲是 remake 差異**：原版在「一個勢力都不存在」時會無限繞
	// （實務上不會發生——玩家自己一定在）。
	at := start
	for step := 0; step < numFactions*2; step++ { //nolint:gocritic // 保險絲見上
		if w.Factions[at].Alive {
			n--
			if n == 0 {
				return w.joinFaction(id, at)
			}
		}
		at++
		if at >= numFactions {
			at = 0
		}
	}
	return -1
}

// joinFaction 是 `loc_1591A`：寫勢力、入帳。
//
// ⭐ **兩條路共用它**（心向兌現與隨機投靠），與原版相同——所以入帳與
// 「投靠玩家才通知」這兩件事只有一份實作。通知本身排在呼叫端
// （`World.Tick` 的月結段），因為 TalkNotice 要掛在那個 tick 的 Event 上。
func (w *World) joinFaction(id, to int) int {
	w.Generals[id].Faction = to
	w.raiseGeneralCount(to)
	return to
}
