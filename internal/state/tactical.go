package state

import (
	"errors"
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// 戰術戰鬥與戰略層的接點。
//
// 原版的分派規則（`sub_14E5C`／`sub_14ED7`，docs/re/09 §2）：
//
//   - **玩家的勢力捲進去 → 先問戰鬥指揮／委任**
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
	// Field 造這一場的戰場。**rotate 為真時整張轉 180 度**——
	// 攻城戰玩家守城時原版 `sub_14ED7` 會 `or byte_10D35, 0C0h`，
	// bit 6 就是這件事（docs/spec/56）。
	Field func(node int, siege, rotate bool) *tactical.Field

	// Script 回傳一側的 AI 腳本。tactic 是武將記錄 `+0x16`（0–7），
	// 段編號 ＝ tactic × 4 ＋ 戰場類別（`sub_1CBE5`，docs/re/11 §3.2）。
	// 沒設或回 nil 那一側就不由腳本驅動。
	Script func(node int, siege bool, tactic int) []byte

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

	// CityWall 是開打時守方據點的城兵數，攻城結算要用它算損害。
	CityWall int
}

// EncounterChoice 是一場等著玩家決定處理方式的遭遇。
//
// 這是原版行軍抵達敵軍／敵城時的中間狀態：
//
//   - 戰鬥指揮：進入戰術畫面，由玩家下令
//   - 委任：不開戰術畫面，直接用自動判定解決
//
// 它只存在於執行期，不是劇本／存檔欄位。
type EncounterChoice struct {
	Attacker int
	Defender int // -1 表示城兵；目前只有軍團對軍團會進這個選單
	Node     int
	Mode     combat.Mode
	Garrison int
}

// SetTactical 裝上戰場來源。傳 nil 就回到全自動判定。
func (w *World) SetTactical(t *TacticalSetup) { w.tactical = t }

// PendingBattle 回傳等著玩家打的那一場，沒有就回 nil。
func (w *World) PendingBattle() *Pending { return w.pending }

// PendingEncounter 回傳等著玩家選擇的遭遇，沒有就回 nil。
// 回傳副本，避免畫面層直接改動規則狀態。
func (w *World) PendingEncounter() *EncounterChoice {
	if w.encounter == nil {
		return nil
	}
	c := *w.encounter
	return &c
}

var errNoEncounter = errors.New("state: 沒有待處理的遭遇")

// ChooseBattleCommand 把待處理的遭遇轉成戰術戰鬥。
func (w *World) ChooseBattleCommand() error {
	if w.encounter == nil {
		return errNoEncounter
	}
	c := *w.encounter
	if !w.beginTactical(c.Attacker, c.Defender, c.Mode, c.Garrison) {
		return errNoTactical
	}
	w.encounter = nil
	return nil
}

// ChooseBattleDelegate 委任待處理的遭遇，回傳與正常自動戰鬥相同的事件。
func (w *World) ChooseBattleDelegate(rng combat.Rand) *CorpsEvent {
	if w.encounter == nil {
		return nil
	}
	c := *w.encounter
	w.encounter = nil
	w.rng = rng
	ev := &CorpsEvent{Corps: c.Attacker, Enemy: c.Defender, Mode: c.Mode,
		Captured: -1, Relocated: -1, GovernorReturned: noGovernor}
	w.resolveCorpsBattle(ev, c.Attacker, c.Defender, c.Mode, c.Garrison, rng)
	return ev
}

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
	// **委任中的軍團不跳選單**：`sub_14E5C` 先看玩家是攻方還是守方，
	// **兩條路各檢查各自那一方的委任位元**（`test [si], 4` 與
	// `test [di], 4`），設起來就退回自動判定（`docs/re/09` §2）。
	// 所以判準是「**玩家那一方**委任中」，與攻守無關。
	if w.Corps[att].Faction == w.Player {
		return !w.Corps[att].Delegated
	}
	if w.Corps[def].Faction == w.Player {
		return !w.Corps[def].Delegated
	}
	return false
}

// beginTactical 準備一場戰術戰鬥。回傳 false 表示開不成，呼叫端該自動判定。
func (w *World) beginTactical(att, def int, m combat.Mode, garrison int) bool {
	node := w.Corps[att].Node
	// ⭐ **玩家守城 → 戰場轉 180 度**（docs/spec/56 §1）。判準是「玩家在哪一邊」，
	// 不是誰攻誰守：原版兩個位元一起設，側欄換邊與戰場翻轉是同一件事的兩面。
	rotate := m == combat.Siege && def >= 0 && def < len(w.Corps) &&
		w.Corps[def].Faction == w.Player
	f := w.tactical.Field(node, m == combat.Siege, rotate)
	if f == nil {
		return false
	}
	// 攻城戰的城壁耐久由**守方據點的城兵數**決定：（城兵數 ＋ 50）× 10
	// （`sub_19CE2` 讀據點記錄 `+0x13`）。野戰用不到。
	cityWall := 0
	if m == combat.Siege && node >= 0 && node < len(w.Cities) {
		cityWall = w.Cities[node].Garrison
	}
	b := tactical.NewBattle(f, w.tactical.Forms, w.rng, cityWall)
	// ⭐ 原版的 side 0 永遠是玩家（`sub_14E5C`／`sub_14ED7` 在玩家守方
	// 那一支互換兩個記錄）。陣形原點與鏡射跟著玩家走，不跟攻守走。
	if def >= 0 && def < len(w.Corps) && w.Corps[def].Faction == w.Player {
		b.SetPlayerSide(tactical.DefenderSide)
	}

	deploy := func(side, corps int) {
		c := w.Corps[corps]
		// 戰力由士氣來（原版 `sub_19B6D` 把軍團士氣寫進每個兵的 +0x18）。
		b.Sides[side].Power = c.Morale
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

	// AI 那一側交給 `BATTLE.DAT` 的腳本（段編號 ＝ 武將 +0x16 × 4 ＋
	// 戰場類別，docs/re/11 §3.2）。**玩家那一側不裝**，由畫面層下命令。
	//
	// 腳本第一個指令通常是查詢，所以掛上去之前先給一個「攻擊」當底，
	// 免得沒有腳本的那一場開局所有人都站著不動。
	corpsOf := [2]int{att, def}
	for side := range b.Sides {
		b.Order(side, -1, tactical.Attack)
		if w.Corps[corpsOf[side]].Faction == w.Player || w.tactical.Script == nil {
			continue
		}
		code := w.tactical.Script(node, m == combat.Siege,
			w.Generals[w.Leader(corpsOf[side])].Tactic)
		if code != nil {
			b.SetScript(side, tactical.NewScript(code, side))
		}
	}

	w.pending = &Pending{Battle: b, Attacker: att, Defender: def,
		Node: node, Mode: m, Garrison: garrison, CityWall: cityWall}
	return true
}

// StageBattle 直接擺一場戰術戰鬥，**繞過遭遇判定**。
//
// 給驗收與測試用：`resolveContact` 的觸發順序是「先野戰再攻城」
// （照 `sub_12708` 的 `cmp byte ptr [di], 0` 那個閘），而守軍就待在城裡時
// 兩者的座標相同，所以走正常流程擺不出攻城的戰術畫面。
// 那個順序對不對是**戰略層的獨立問題**（見 CONTEXT.md），不在這裡繞過去。
func (w *World) StageBattle(att, def int, m combat.Mode, rng combat.Rand) error {
	if w.tactical == nil || w.tactical.Forms == nil || w.tactical.Field == nil {
		return errNoTactical
	}
	// 正常路徑是 `Tick` 把亂數源記進 w.rng；這裡繞過了 Tick，要自己給。
	w.rng = rng
	for _, i := range []int{att, def} {
		if i < 0 || i >= numCorps || !w.Corps[i].Alive {
			return fmt.Errorf("state: 軍團 %d 不存在", i)
		}
	}
	if !w.beginTactical(att, def, m, w.Cities[w.clampCity(w.Corps[att].Node)].Garrison) {
		return fmt.Errorf("state: 擺不出戰場（據點 %d）", w.Corps[att].Node)
	}
	return nil
}

var errNoTactical = errors.New("state: 沒有裝戰場來源")

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
	if w == nil || w.outcome != InProgress {
		return nil
	}
	p := w.pending
	if p == nil || !p.Battle.Done {
		return nil
	}
	w.pending = nil
	o := p.Battle.Result()

	ev := &CorpsEvent{Corps: p.Attacker, Enemy: p.Defender, Captured: -1,
		Relocated: -1, GovernorReturned: noGovernor}
	ev.BattleBefore = [2]int{w.Corps[p.Attacker].Men, w.Corps[p.Defender].Men}
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
	// 攻城的損害走**戰術層自己的公式**（`sub_19FDC` 前半：由城壁被打掉
	// 多少算，不是自動判定的兵力比），三個欄位各扣同一個值。
	r.CityDamage = p.Battle.CityDamage(p.CityWall)
	ev.BattleAfter = [2]int{w.Corps[p.Attacker].Men, w.Corps[p.Defender].Men}
	ev.BattleCityDamage = r.CityDamage
	w.damageCity(p.Node, p.Mode, r)
	attDead := r.AttackerDestroyed || w.retreatOrPerish(p.Attacker, !r.DefenderWins)
	defDead := r.DefenderDestroyed || w.retreatOrPerish(p.Defender, r.DefenderWins)
	w.afterBattle(ev, p.Attacker, attDead, p.Defender, rng)
	w.afterBattle(ev, p.Defender, defDead, p.Attacker, rng)
	if o.AttackerWins && p.Mode == combat.Siege && !attDead {
		w.capture(p.Attacker, ev, rng)
	}
	return ev
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
