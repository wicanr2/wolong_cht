package tactical

// 開戰單挑（docs/spec/80、docs/re/74）。
//
// 原版 `sub_1A1C5` 是主迴圈**之前**的 blocking sequence：每場戰鬥先空跑
// 50 tick（實體照常更新，但腳本直譯器與玩家輸入在主迴圈裡，一步都不跑），
// `byte_1D34B == 1` 才走單挑——挑戰／拒戰／應戰、回合互嗆、對打段、決著。
// 整段只寫**目標座標**與命令 8：大將騎過去，移動與接觸命中走命令分派
// 之外的共通路徑，所以傷害就是一般碰撞命中，沒有單挑專用公式。

// Duel 是命令 8：命令分派是 nullsub，目標座標由單挑狀態機控制。
// **只掛大將本人**，隊員維持陣形。
const Duel Command = 8

// duelMoraleFloor 是挑戰／應戰的氣勢門檻（`sub_1A2E8` 的 `cmp ax, 12C0h`）。
const duelMoraleFloor = 0x12C0

// duelEndHP 是決著門檻：任一大將體力低於它就分勝負（`cmp [±3], 46h`）。
const duelEndHP = 0x46

// duelOpeningTicks 是每場戰鬥開頭的空跑（`sub_1A1C5` 的 `mov cx, 32h`）。
// 這段期間腳本與輸入都還沒開始，不分有沒有單挑。
const duelOpeningTicks = 50

// DuelInput 是單挑狀態機需要的外部資料（呈現層在建 Battle 後餵）。
type DuelInput struct {
	// FieldNumber 是戰場編號；只有 0xC0–0xD0（`sub_19A33` 分類出的
	// `byte_1D34B == 1`）觸發單挑開場。
	FieldNumber int
	// Martial／CommandStat 是兩側大將的武術與統率（武將記錄
	// +0x11／+0x12，`sub_19C13`），索引照 Battle 的側別（0 攻、1 守）。
	Martial, CommandStat [2]int
}

// DuelTalk 是一則單挑喊話：Side 是說話側、Group 是 TALK 組編號
// （0x1B7 挑戰、0x1B8 應戰、0x1B9 拒戰、0x1BA+ 互嗆、0x1CC/0x1CD 決著）。
type DuelTalk struct {
	Side  int
	Group int
}

// 狀態機的相位。
const (
	duelIdle      = iota // 未觸發或已結束
	duelWait             // 開場等 50 tick（`sub_1A1C5` 開頭）
	duelChallenge        // 強側喊完 0x1B7、等 40
	duelRefuseA          // 弱側喊 0x1B9、等 20
	duelBanterA          // 互嗆第一句喊完、等 10（`sub_1A3C3` 的 cx=0Ah）
	duelMelee            // 對打段（`sub_1A298`，計時 0x50）
	duelRegroup          // 歸位等 20（loc_1A202，每 tick 檢查體力）
	duelVerdictA         // 敗方 0x1CC、等 20
	duelVerdictB         // 勝方 0x1CD、等 20 後收尾
)

type duelState struct {
	input DuelInput
	armed bool
	phase int
	timer int
	// round 是互嗆情境碼（原版 `word_1D31F`，上限 4）。
	round int
	// strong 是氣勢高側；hi／lo 是兩側的氣勢值（同分不換、攻方在前）。
	strong int
	hi, lo int
	// first 是本回合互嗆先講的一側（r0＝強側；r≥1＝體力高側，同分不換）。
	first int
	// pair 是本回合的互嗆組（第一句），第二句＝pair+1。
	pair int
	// loser 在決著段有意義。
	loser int
	talks []DuelTalk
}

// SetDuelInput 武裝單挑狀態機。編號不在 0xC0–0xD0 就不動作
// （攻城 < 0xC0 與其餘野戰 ≥ 0xD1 沒有單挑開場，但 50 tick 的
// 開場空跑照樣有——那在 `OpeningActive` 裡，不看這裡）。
func (b *Battle) SetDuelInput(in DuelInput) {
	if in.FieldNumber < 0xC0 || in.FieldNumber > 0xD0 {
		return
	}
	b.duel = duelState{input: in, armed: true, phase: duelWait, timer: duelOpeningTicks}
}

// TakeDuelTalks 取走累積的喊話（呈現層轉成 TALK 索引畫對白框）。
//
// ⚠ 名字留著沒改，但**它不只有單挑**——腳本指令 16（攻城戰的開場
// 勸降，docs/spec/135）走同一條出口，原版兩者也都是 `sub_1C315`。
func (b *Battle) TakeDuelTalks() []DuelTalk {
	out := b.talks
	b.talks = nil
	return out
}

// say 把一則對白排進出口。
func (b *Battle) say(side, group int) {
	b.talks = append(b.talks, DuelTalk{Side: side, Group: group})
}

// DuelActive 回報單挑流程是否進行中（挑戰喊話起、決著收尾止）。
func (b *Battle) DuelActive() bool {
	return b.duel.armed && b.duel.phase > duelWait
}

// OpeningActive 回報「主迴圈還沒開始」：每場戰鬥開頭的 50 tick 空跑，
// 加上整段單挑（原版都在 `sub_1A1C5` 裡）。這段期間**腳本不跑、
// 玩家輸入不收**——兩軍靜止靠的是「還沒有人下過命令」，不是實體凍結。
func (b *Battle) OpeningActive() bool {
	if b.Frame <= duelOpeningTicks {
		return true
	}
	return b.DuelActive()
}



// duelMorale 重現 `sub_1A34F`：氣勢 ＝ 大將**戰力 × 體力**，
// 武術門檻（`max(0, 武術×3−統率)/2`）沒過整個歸零，最後加亂數尾。
//
// `[bx+0x18]` 是單位記錄的**戰力**（`sub_19B6D` 由軍團士氣寫入，
// remake 的 `Power`），不是兵數——第一版讀成兵數，兵多的一側永遠
// 先挑戰，與實機兩輪錄影（氣勢高的呂布先喊）相反。
// 體力開場＝士氣（docs/spec/61），所以核心實質是士氣的平方。
func (b *Battle) duelMorale(side int) int {
	g := &b.Sides[side].Soldiers[0]
	core := g.Power * g.HP
	threshold := b.duel.input.Martial[side]*3 - b.duel.input.CommandStat[side]
	if threshold < 0 {
		threshold = 0
	}
	threshold /= 2
	if (b.rng.Next()&7)+8 > threshold {
		core = 0
	}
	return core + (b.rng.Next()&7)<<8
}

// duelSpot 是側 side 大將的單挑位（攻 (0x18,0x20)／守 (0x28,0x20)）。
func duelSpot(side int) (int, int) {
	if side == 0 {
		return 0x18, 0x20
	}
	return 0x28, 0x20
}

// duelGoal 把側 side 大將的目標座標指到 (x,y)——**不搬人**，
// 大將自己騎過去（原版寫 `+0x10/+0x11/+0x14`，移動在共通路徑）。
func (b *Battle) duelGoal(side, x, y int) {
	g := &b.Sides[side].Soldiers[0]
	g.GoalX, g.GoalY = x, y
	g.GoalZ = b.standZ(g, x, y)
	g.Path = nil
}

// duelFace 把兩側大將的目標指回各自的單挑位（loc_1A202）。
func (b *Battle) duelFace() {
	for i := range b.Sides {
		x, y := duelSpot(i)
		b.duelGoal(i, x, y)
	}
}

// duelMeet 把兩側大將的目標指到同一格。
func (b *Battle) duelMeet(x, y int) {
	for i := range b.Sides {
		b.duelGoal(i, x, y)
	}
}

// orderDuelLeader 設或清側 side **大將本人**的命令 8（`[si+1Bh]`）。
func (b *Battle) orderDuelLeader(side int, on bool) {
	c := Duel
	if !on {
		c = Form
	}
	g := &b.Sides[side].Soldiers[0]
	if g.Alive {
		g.Cmd, g.Next = c, c
	}
}

func (b *Battle) duelHP(side int) int { return b.Sides[side].Soldiers[0].HP }

// duelKO 檢查決著門檻；成立就切到決著相位並回報 true。
func (b *Battle) duelKO() bool {
	d := &b.duel
	if b.duelHP(0) >= duelEndHP && b.duelHP(1) >= duelEndHP {
		return false
	}
	d.loser = 0
	if b.duelHP(1) < b.duelHP(0) {
		d.loser = 1
	}
	d.phase = duelVerdictA
	d.timer = 0
	return true
}

// duelBanterPair 算本回合的互嗆 pair 與先講側（`sub_1A3C3`）。
func (b *Battle) duelBanterPair() int {
	d := &b.duel
	if d.round == 0 {
		d.first = d.strong
		return 0x1BA
	}
	pair := 0x1BC + (d.round-1)*4
	diff := b.duelHP(0) - b.duelHP(1)
	// 先講的是體力高側；同分不換（沿用上一回合的先講側）。
	if diff > 0 {
		d.first = 0
	} else if diff < 0 {
		d.first = 1
	}
	if diff < 0x14 && diff > -0x14 {
		pair += 2
	}
	return pair
}

// stepDuel 每戰場 tick 推進一次（掛在 Step() 開頭——
// 原版 `sub_1A1C5` 在主迴圈之前，remake 以相位機拆進逐 tick 的 Step）。
func (b *Battle) stepDuel() {
	d := &b.duel
	if !d.armed || d.phase == duelIdle {
		return
	}
	// 對打段與歸位期間每 tick 檢查決著門檻
	// （`sub_1A298` 與 loc_1A219 的 `cmp [±3], 46h`）。
	if (d.phase == duelMelee || d.phase == duelRegroup) && b.duelKO() {
		return
	}
	if d.timer > 0 {
		d.timer--
		// 對打段：計時 0x50 起算，剩餘 ≥ 0x30 純等；之後每 tick
		// 機率 0x20/256 把兩大將的目標改成同一隨機格
		// （抽的順序照原版：先 y ＝ rand&7+0x1C、再 x ＝ rand&0xF+0x18）。
		if d.phase == duelMelee && d.timer < 0x30 && b.rng.Next()&0xFF < 0x20 {
			y := b.rng.Next()&7 + 0x1C
			x := b.rng.Next()&0xF + 0x18
			b.duelMeet(x, y)
		}
		if d.timer > 0 {
			return
		}
	}
	weak := 1 - d.strong
	switch d.phase {
	case duelWait:
		// 評估（`sub_1A2E8`）。同分不換：攻方（記錄 0）在前。
		m0, m1 := b.duelMorale(0), b.duelMorale(1)
		d.strong, d.hi, d.lo = 0, m0, m1
		if m1 > m0 {
			d.strong, d.hi, d.lo = 1, m1, m0
		}
		if d.hi < duelMoraleFloor {
			d.phase = duelIdle // 沒人夠氣勢，正常開打
			return
		}
		// `sub_1A398`：喊話＋目標指向己側單挑位＋大將命令 8＋等 40。
		b.say(d.strong, 0x1B7)
		x, y := duelSpot(d.strong)
		b.duelGoal(d.strong, x, y)
		b.orderDuelLeader(d.strong, true)
		// `add word_1D311, 6`：跳過腳本開頭的「message／wait 15／message」，
		// 抑制腳本自帶的開場喊話（docs/spec/80 §3.2）。
		for i := range b.scripts {
			b.scripts[i].SkipBytes(6)
		}
		d.phase = duelChallenge
		d.timer = 40
	case duelChallenge:
		if d.lo < duelMoraleFloor || d.lo < d.hi/2 {
			b.say(weak, 0x1B9) // 拒戰
			d.phase = duelRefuseA
			d.timer = 20
			return
		}
		// 應戰（loc_1A341 → sub_1A398）：弱側喊 0x1B8、騎向己側單挑位。
		b.say(weak, 0x1B8)
		x, y := duelSpot(weak)
		b.duelGoal(weak, x, y)
		b.orderDuelLeader(weak, true)
		d.round = 0
		// 應戰後接的是互嗆第一句（回合迴圈開頭），先等應戰的 40 tick。
		d.phase = duelRegroup
		d.timer = 40
	case duelRefuseA:
		// 拒戰回應之後**立即**清命令回到正常戰鬥（原版 0x1CC 後沒有等待）。
		b.say(d.strong, 0x1CC)
		b.orderDuelLeader(d.strong, false)
		d.phase = duelIdle
	case duelRegroup:
		// 回合迴圈開頭：互嗆第一句（pair 與先講側一回合只算一次）。
		d.pair = b.duelBanterPair()
		b.say(d.first, d.pair)
		d.phase = duelBanterA
		d.timer = 10
	case duelBanterA:
		// 互嗆第二句＋兩大將騎向會合點 (0x20,0x20)，進對打段。
		b.say(1-d.first, d.pair+1)
		b.duelMeet(0x20, 0x20)
		d.phase = duelMelee
		d.timer = 0x50
	case duelMelee:
		// 計時歸零、沒人倒下：歸位、等 20，情境碼 +1（上限 4）再互嗆。
		b.duelFace()
		if d.round < 4 {
			d.round++
		}
		d.phase = duelRegroup
		d.timer = 20
	case duelVerdictA:
		// 敗方已在退卻就不喊、也不清命令（loc_1A251 的 `cmp [di+1Ah], 5`——
		// 清命令那一行在 ≠5 的分支裡，退卻的維持退卻）。
		if b.Sides[d.loser].Soldiers[0].Cmd != Retreat {
			b.say(d.loser, 0x1CC)
			b.orderDuelLeader(d.loser, false)
		}
		d.phase = duelVerdictB
		d.timer = 20
	case duelVerdictB:
		b.say(1-d.loser, 0x1CD)
		// 收尾清命令（loc_1A27A）。退卻中的維持退卻不清。
		for i := range b.Sides {
			if b.Sides[i].Soldiers[0].Cmd != Retreat {
				b.orderDuelLeader(i, false)
			}
		}
		d.phase = duelIdle
	}
}
