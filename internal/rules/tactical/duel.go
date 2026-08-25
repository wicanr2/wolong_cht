package tactical

// 開戰單挑（docs/spec/80、docs/re/74）。
//
// 原版 `sub_1A1C5` 在戰場主迴圈開頭跑一次：等 50 tick、依氣勢評估
// 決定挑戰／拒戰／應戰，應戰就進入「互嗆 → 對打段」的回合迴圈，
// 任一大將體力 < 0x46 決著。整段只做**編排**：命令 8（Duel）的處理
// 常式在原版是 nullsub（`funcs_1A7E1[8]`／`funcs_1A827[8]`），
// 大將隊全體靜止、位置由本狀態機傳送，體力變化來自戰場上其他部隊的
// 普通攻擊——單挑不需要專用傷害公式。

// Duel 是命令 8：單挑中，整隊靜止。
const Duel Command = 8

// duelMoraleFloor 是挑戰／應戰的氣勢門檻（`sub_1A2E8` 的 `cmp ax, 12C0h`）。
const duelMoraleFloor = 0x12C0

// duelEndHP 是決著門檻：任一大將體力低於它就分勝負（`cmp [±3], 46h`）。
const duelEndHP = 0x46

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
	duelRefuseB          // 強側回 0x1CC、等 20 後收尾
	duelBanter           // 互嗆 pair（`sub_1A3C3`）
	duelMelee            // 對打段（`sub_1A298`，計時 0x50）
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
	// loser 在決著段有意義。
	loser int
	// banterSecond 標記互嗆 pair 的第二句還沒說。
	banterSecond bool
	talks        []DuelTalk
}

// SetDuelInput 武裝單挑狀態機。編號不在 0xC0–0xD0 就不動作
// （攻城 < 0xC0 與其餘野戰 ≥ 0xD1 沒有單挑開場）。
func (b *Battle) SetDuelInput(in DuelInput) {
	if in.FieldNumber < 0xC0 || in.FieldNumber > 0xD0 {
		return
	}
	b.duel = duelState{input: in, armed: true, phase: duelWait, timer: 50}
}

// TakeDuelTalks 取走累積的喊話（呈現層轉成 TALK 索引畫對白框）。
func (b *Battle) TakeDuelTalks() []DuelTalk {
	out := b.duel.talks
	b.duel.talks = nil
	return out
}

// DuelActive 回報是否在單挑流程中（大將隊被命令 8 凍結的期間）。
func (b *Battle) DuelActive() bool {
	return b.duel.armed && b.duel.phase > duelWait
}

// duelOpeningFreeze 回報開場序（等待、挑戰、拒戰）是否凍結全場。
// 實機證據（playtest/43 的 b0–b3）：挑戰到拒戰交鋒的六秒內
// **兩軍都不動**，只有強側大將騎出去——原版把整段開場當成
// blocking sequence 跑。應戰進入回合迴圈後恢復正常更新
// （大將的體力變化來自周圍的普通戰鬥，docs/spec/80 §4）。
func (b *Battle) duelOpeningFreeze() bool {
	if !b.duel.armed {
		return false
	}
	switch b.duel.phase {
	case duelWait, duelChallenge, duelRefuseA, duelRefuseB:
		return true
	}
	return false
}

func (d *duelState) say(side, group int) {
	d.talks = append(d.talks, DuelTalk{Side: side, Group: group})
}

// duelMorale 重現 `sub_1A34F`：氣勢 ＝ 大將隊兵數 × 大將體力，
// 武術門檻（`max(0, 武術×3−統率)/2`）沒過整個歸零，最後加亂數尾。
func (b *Battle) duelMorale(side int) int {
	s := &b.Sides[side]
	// 大將隊兵數＝場上＋待機（原版側記錄 +0x18 是整隊人數，
	// 不是場上那 8 個——只算場上永遠到不了 0x12C0 門檻）。
	men := s.Reserve[0]
	for i := 0; i < PerSquad; i++ {
		if s.Soldiers[i].Alive {
			men++
		}
	}
	core := men * s.Soldiers[0].HP
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

// duelSpot 是側 side 大將的對峙位（`sub_1A1C5` 迴圈裡寫死的
// (0x18,0x20)／(0x28,0x20)）。
func duelSpot(side int) (int, int) {
	if side == 0 {
		return 0x18, 0x20
	}
	return 0x28, 0x20
}

// duelPlace 把側 side 的大將直接放到 (x,y)（單挑的位置全由狀態機寫）。
func (b *Battle) duelPlace(side, x, y int) {
	g := &b.Sides[side].Soldiers[0]
	g.X, g.Y = x, y
	g.Z = b.standZ(g, x, y)
	g.GoalX, g.GoalY, g.GoalZ = g.X, g.Y, g.Z
	g.StepX, g.StepY, g.StepZ = g.X, g.Y, g.Z
}

// duelRideOut 讓側 side 的大將朝單挑位走一格（每 tick 一次）。
func (b *Battle) duelRideOut(side int) {
	g := &b.Sides[side].Soldiers[0]
	tx, ty := duelSpot(side)
	if g.X == tx && g.Y == ty {
		return
	}
	step := func(v, t int) int {
		if v < t {
			return v + 1
		}
		if v > t {
			return v - 1
		}
		return v
	}
	b.duelPlace(side, step(g.X, tx), step(g.Y, ty))
}

// duelFace 把兩側大將放回對峙位。
func (b *Battle) duelFace() {
	for i := range b.Sides {
		x, y := duelSpot(i)
		b.duelPlace(i, x, y)
	}
}

// duelTeleport 把兩側大將傳到同一個隨機格（`sub_1A298` 的
// x = rand&0xF + 0x18、y = rand&7 + 0x1C）。
func (b *Battle) duelTeleport() {
	y := b.rng.Next()&7 + 0x1C
	x := b.rng.Next()&0xF + 0x18
	for i := range b.Sides {
		b.duelPlace(i, x, y)
	}
}

// orderDuelSquad 設或清側 side 大將隊（第 0 隊）的命令 8。
func (b *Battle) orderDuelSquad(side int, on bool) {
	c := Duel
	if !on {
		c = Form
	}
	for k := 0; k < PerSquad; k++ {
		s := &b.Sides[side].Soldiers[k]
		if s.Alive {
			s.Cmd, s.Next = c, c
		}
	}
}

func (b *Battle) duelHP(side int) int { return b.Sides[side].Soldiers[0].HP }

// stepDuel 每戰場 tick 推進一次（掛在 Step() 開頭，腳本之前——
// 原版 `sub_1A1C5` 也在主迴圈最前面）。
func (b *Battle) stepDuel() {
	d := &b.duel
	if !d.armed || d.phase == duelIdle {
		return
	}
	// 對打／互嗆期間每 tick 檢查決著門檻（`sub_1A298` 的 `cmp [±3], 46h`）。
	if d.phase == duelBanter || d.phase == duelMelee {
		if b.duelHP(0) < duelEndHP || b.duelHP(1) < duelEndHP {
			d.loser = 0
			if b.duelHP(1) < b.duelHP(0) {
				d.loser = 1
			}
			d.phase = duelVerdictA
			d.timer = 0
		}
	}
	if d.timer > 0 {
		d.timer--
		// 挑戰／拒戰期間：強側大將**逐格騎向**單挑位（b3 的實機截圖裡
		// 它已在半路上；不是瞬移——b0 那一刻它還在原位）。
		if d.phase == duelChallenge || d.phase == duelRefuseA || d.phase == duelRefuseB {
			b.duelRideOut(d.strong)
		}
		// 對打段：計時 0x50 起算，前 0x20 tick 純等；
		// 之後每 tick 1/8 機率把兩大將傳到同一格。
		if d.phase == duelMelee && d.timer < 0x30 && b.rng.Next()&0xFF < 0x20 {
			b.duelTeleport()
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
		// `sub_1A398`：喊話＋強側大將隊命令 8、等 40。
		// **就定位在應戰才發生**——原版 b0（挑戰喊話當下）的實機截圖裡
		// 大將仍在原位（playtest/43）。
		d.say(d.strong, 0x1B7)
		b.orderDuelSquad(d.strong, true)
		d.phase = duelChallenge
		d.timer = 40
	case duelChallenge:
		if d.lo < duelMoraleFloor || d.lo < d.hi/2 {
			d.say(weak, 0x1B9) // 拒戰
			d.phase = duelRefuseA
			d.timer = 20
			return
		}
		// 應戰（loc_1A341）：兩大將對峙就位，回合迴圈開始。
		d.say(weak, 0x1B8)
		b.orderDuelSquad(weak, true)
		b.duelFace()
		d.round = 0
		d.phase = duelBanter
		d.banterSecond = false
		d.timer = 20
	case duelRefuseA:
		d.say(d.strong, 0x1CC)
		d.phase = duelRefuseB
		d.timer = 20
	case duelRefuseB:
		b.orderDuelSquad(d.strong, false)
		d.phase = duelIdle
	case duelBanter:
		// `sub_1A3C3`：round 0 用 0x1BA/0x1BB，之後每回合 4 組、
		// 體力差 < 0x14 再 +2 取「勢均」pair。
		pair := 0x1BA
		if d.round > 0 {
			pair = 0x1BC + (d.round-1)*4
			if diff := b.duelHP(0) - b.duelHP(1); diff < 0x14 && diff > -0x14 {
				pair += 2
			}
		}
		if !d.banterSecond {
			d.say(d.strong, pair)
			d.banterSecond = true
			d.timer = 40
			return
		}
		d.say(weak, pair+1)
		d.phase = duelMelee
		d.timer = 0x50
	case duelMelee:
		// 計時歸零、沒人倒下：兩大將回對峙位，情境碼 +1（上限 4）再互嗆。
		b.duelFace()
		if d.round < 4 {
			d.round++
		}
		d.phase = duelBanter
		d.banterSecond = false
		d.timer = 20
	case duelVerdictA:
		// 敗方已在退卻就不喊（loc_1A251 的 `cmp [di+1Ah], 5`）。
		if b.Sides[d.loser].Soldiers[0].Cmd != Retreat {
			d.say(d.loser, 0x1CC)
		}
		b.orderDuelSquad(d.loser, false)
		d.phase = duelVerdictB
		d.timer = 20
	case duelVerdictB:
		d.say(1-d.loser, 0x1CD)
		for i := range b.Sides {
			b.orderDuelSquad(i, false)
		}
		d.phase = duelIdle
	}
}
