package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// 戰術戰鬥與戰略層的接點。
//
// 原版的分派規則（`sub_14E5C`／`sub_14ED7`，docs/re/09 §2）：
//
//   - **玩家的勢力捲進去 → 開戰術畫面**
//   - 其餘一律自動判定
//   - 例外：**攻打空城**（城裡沒有軍團）也走自動判定
//
// 這裡照抄那條規則。`Tactical` 沒設好時全部退回自動判定——
// 戰場資產（`BATTLE.*`）不是必要的，無頭模擬照樣跑得起來。

// TacticalSetup 是呼叫端提供的戰場來源。
//
// node 是據點編號（攻城戰）或戰場編號（野戰，見 docs/re/05）；
// 回傳 nil 表示這一場開不了戰術畫面，退回自動判定。
type TacticalSetup struct {
	Field func(node int, siege bool) *tactical.Field
	Forms *tactical.Formations
}

// Pending 是一場等著被跑完的戰術戰鬥。
//
// 世界迴圈在它還在的時候**不會前進**——原版進戰術畫面時戰略時間就停了，
// 而戰術層本身是硬即時的（說明書 4.1：「戦闘中は絶対に時間を止められません」）。
type Pending struct {
	Battle *tactical.Battle

	Attacker int // 攻方的軍團編號
	Defender int // 守方的軍團編號，−1 表示對手是城兵
	Node     int
	Mode     combat.Mode
	Garrison int
}

// SetTactical 裝上戰場來源。傳 nil 就回到全自動判定。
func (w *World) SetTactical(t *TacticalSetup) { w.tactical = t }

// PendingBattle 回傳等著玩家打的那一場，沒有就回 nil。
func (w *World) PendingBattle() *Pending { return w.pending }

// wantsTactical 回報這一場該不該開戰術畫面。
func (w *World) wantsTactical(att, def int) bool {
	if w.tactical == nil || w.tactical.Forms == nil || w.tactical.Field == nil {
		return false
	}
	if w.Player < 0 {
		return false
	}
	// **打空城不進戰術畫面**（原版 `cmp bx, 4200h`，docs/re/09 §2）。
	if def < 0 {
		return false
	}
	return w.Corps[att].Faction == w.Player || w.Corps[def].Faction == w.Player
}

// beginTactical 準備一場戰術戰鬥。回傳 false 表示開不成，呼叫端該自動判定。
func (w *World) beginTactical(att, def int, m combat.Mode, garrison int) bool {
	node := w.Corps[att].Node
	f := w.tactical.Field(node, m == combat.Siege)
	if f == nil {
		return false
	}
	b := tactical.NewBattle(f, w.tactical.Forms, w.rng)
	// 陣形線：兩軍都擺在自己那一側（說明書 4.3 的三選一，預設自軍側）。
	b.Sides[0].Line = tactical.LineFor(0, 0)
	b.Sides[1].Line = tactical.LineFor(1, 0)

	deploy := func(side, corps int) {
		c := w.Corps[corps]
		for k, u := range c.Units {
			if u.Men == 0 {
				continue
			}
			// 戰略上一點兵力 ＝ 10 人，戰術上一個兵 ＝ 10 人，所以
			// 「一點」正好換一個兵（說明書 4.1 的 1 : 10）。
			b.Deploy(side, k, kindOf(u.Kind), u.Men)
		}
	}
	deploy(0, att)
	deploy(1, def)
	b.Place()

	// ⚠ 開場都下「攻擊」。
	//
	// 原版是由 `BATTLE.DAT` 的腳本驅動 AI（段編號 ＝ 武將 +0x16 × 4 ＋
	// 戰場類別，docs/re/11 §3.2），十九個指令都讀出來了，但**那台直譯器
	// 還沒實作**——它查的那些全域變數有一半還沒解。
	// 在那之前用這個代替，**不要當成原版行為**。
	b.Order(0, -1, tactical.Attack)
	b.Order(1, -1, tactical.Attack)

	w.pending = &Pending{Battle: b, Attacker: att, Defender: def,
		Node: node, Mode: m, Garrison: garrison}
	return true
}

func kindOf(t army.TroopType) tactical.Kind {
	switch t {
	case army.Cavalry:
		return tactical.Cavalry
	case army.Archer:
		return tactical.Archer
	default:
		return tactical.Infantry
	}
}

// ResolvePending 把打完的戰術戰鬥結算回戰略層。
//
// 換算照說明書 4.1 的 1 : 10：戰場上剩下的兵數 × 10 ＝ 戰略的人數，
// 再除以 10 換回「點」。壞滅與敗將的下場走與自動判定同一條路
// （`afterBattle`），所以兩種戰鬥的出口一致——原版也是這樣
// （docs/re/09 §10：戰術層回傳的 `al`／`ah` 與自動判定同格式）。
func (w *World) ResolvePending(rng combat.Rand) *CorpsEvent {
	p := w.pending
	if p == nil || !p.Battle.Done {
		return nil
	}
	w.pending = nil
	o := p.Battle.Result()

	ev := &CorpsEvent{Corps: p.Attacker, Enemy: p.Defender, Captured: -1}
	r := combat.Result{DefenderWins: !o.AttackerWins}
	ev.Battle = &r

	apply := func(corps, side int) {
		c := &w.Corps[corps]
		men := o.Men[side] / tactical.MenPerSoldier
		scale := func(v int) int {
			if c.Men == 0 {
				return 0
			}
			return v * men / c.Men
		}
		total := 0
		for k := range c.Units {
			c.Units[k].Men = scale(c.Units[k].Men)
			total += c.Units[k].Men
		}
		c.Men = total
		// 士氣照自動判定的規則走：敗方重設成 100 × 兵力比。
		if (side == 0) == o.AttackerWins {
			c.Morale = c.Morale * men / maxInt(men, 1)
		} else if c.Morale < army.RoutMoraleGate {
			c.Morale = 0
		} else {
			c.Morale = 100 * total / maxInt(total+1, 1)
		}
		if total == 0 {
			c.Morale = 0
		}
	}
	apply(p.Attacker, 0)
	apply(p.Defender, 1)

	r.AttackerDestroyed = combat.Destroyed(w.battle(p.Attacker))
	r.DefenderDestroyed = combat.Destroyed(w.battle(p.Defender))
	w.damageCity(p.Node, p.Mode, r)
	w.afterBattle(ev, p.Attacker, r.AttackerDestroyed, p.Defender, rng)
	w.afterBattle(ev, p.Defender, r.DefenderDestroyed, p.Attacker, rng)
	if o.AttackerWins && p.Mode == combat.Siege && !r.AttackerDestroyed {
		w.capture(p.Attacker, ev)
	}
	return ev
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
