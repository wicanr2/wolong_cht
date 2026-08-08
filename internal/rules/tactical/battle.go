package tactical

// 一場戰術戰鬥。每幀跑一次 Step，直到 Done。

// Rand 是規則層共用的亂數介面。
type Rand interface{ Next() int }

// Side 是交戰的一方。
type Side struct {
	// Soldiers 是場上的 48 個兵：第 k 隊佔 [k*8, k*8+8)，**每隊的第一個是隊長**
	// （`sub_1A754` 的 1 + 7 迴圈）。
	Soldiers [SoldiersOnFoot]Soldier

	// Reserve[k] 是第 k 隊還在畫面外待機的兵數。
	//
	// 說明書 4.1：「1つの戦場に全軍がはいることは出来ませんので、
	// 残りは**予備兵として画面外に待機**します」。一個編成位置 1,000 人
	// ＝ 100 個兵，場上只放得下 8 個。
	Reserve [Squads]int

	// Formation 是目前的陣形編號（0–15）。腳本用指令 1 切換。
	Formation int
	// Line 是陣形線：自軍側／中央／敵軍側三選一，值是戰場上的 X
	// （指令 2 寫的三個常數 58／36／16，docs/re/11 §5.8d）。
	Line int

	// Mirror 為真表示這一側的陣形要左右鏡射（原版對 0x600 側 `neg dl`）。
	Mirror bool

	// GateOpen 記錄守方是否已經開門。說明書 4.2：突擊會開門，
	// 而**開了的門這場戰鬥不能再關**。
	GateOpen bool
}

// Alive 回傳這一側還活著的兵數（不含待機的）。
func (s *Side) Alive() int {
	n := 0
	for i := range s.Soldiers {
		if s.Soldiers[i].Alive {
			n++
		}
	}
	return n
}

// Fresh 回傳「疲勞度還沒歸零」的兵數。原版每幀累加進 word_1D31A，
// 拿來算有利／不利（docs/re/11 §5.8g）。
func (s *Side) Fresh() int {
	n := 0
	for i := range s.Soldiers {
		if s.Soldiers[i].Alive && s.Soldiers[i].Stamina > 0 {
			n++
		}
	}
	return n
}

// Remaining 回傳這一側的總戰力：場上的 ＋ 待機的。
// **兩邊都歸零才算補不出兵**（說明書 4.1 的勝負條件）。
func (s *Side) Remaining() int {
	n := s.Alive()
	for _, r := range s.Reserve {
		n += r
	}
	return n
}

// Advantage 是有利／不利的三階（原版 byte_1D31E）。
type Advantage int

const (
	Disadvantaged Advantage = 0
	Even          Advantage = 1
	Advantaged    Advantage = 2
)

// advantageBias 是算有利／不利時給對方的讓分（原版 `add ah, 7`）。
const advantageBias = 7

// evenBand 是「差距在這個範圍內就算普通」的門檻（原版 `cmp al, 8 / ja`）。
//
// 說明書 6.1：「**敵も同じ状態の場合は通常と変わりません**」——
// 這個門檻就是那句話。
const evenBand = 8

// computeAdvantage 重現 `sub_1ADC8` 尾段。
func computeAdvantage(mine, theirs int) Advantage {
	a := Advantaged
	d := mine - (theirs + advantageBias)
	if d < 0 {
		a = Disadvantaged
		d = -d
	}
	if d <= evenBand {
		return Even
	}
	return a
}

// Battle 是一場進行中的戰術戰鬥。
type Battle struct {
	Field *Field
	// Sides[0] 是攻方、Sides[1] 是守方。
	Sides [2]Side

	// Forms 是 16 種陣形，每種 48 組相對座標（docs/re/11 §5.8d）。
	Forms *Formations

	Frame int
	Done  bool
	// Winner 在 Done 之後有意義：0 或 1。
	Winner int

	// Advantage[i] 是第 i 側目前的有利／不利。
	Advantage [2]Advantage

	rng Rand

	// projectiles 是飛在空中的箭。原版是一張 32 筆的表（docs/re/11 §5.1）。
	projectiles []projectile

	// siegeTick 是攻城方大將體力遞減的計時器。
	siegeTick int

	// Log 記錄值得回報的事件，供呼叫端呈現。
	Log []string
}

type projectile struct {
	side       int // 誰射的
	x, y, z    int
	tx, ty, tz int
	power      int
}

// NewBattle 開一場戰鬥。
func NewBattle(f *Field, forms *Formations, rng Rand) *Battle {
	b := &Battle{Field: f, Forms: forms, rng: rng, Winner: -1}
	b.Sides[1].Mirror = true
	return b
}

// Deploy 把一隊放上戰場。kind 是六個編成位置的兵種，men 是各位置的兵數
// （以「兵」為單位，一個位置滿編 100）。
//
// 場上一隊只放 8 個，其餘進待機（說明書 4.1）。
func (b *Battle) Deploy(side int, squad int, kind Kind, men int) {
	s := &b.Sides[side]
	on := men
	if on > PerSquad {
		on = PerSquad
	}
	s.Reserve[squad] = men - on
	for i := 0; i < on; i++ {
		k := kind
		if i == 0 && squad == 0 {
			k = General // 第 0 隊的隊長是大將
		}
		s.Soldiers[squad*PerSquad+i] = Soldier{
			Alive: true, Kind: k, HP: MaxHP, Stamina: StaminaFull,
			Target: -1, Cmd: Form, Next: Form,
		}
	}
}

// Place 把所有兵放到陣形的起始位置。
func (b *Battle) Place() {
	for i := range b.Sides {
		for k := range b.Sides[i].Soldiers {
			s := &b.Sides[i].Soldiers[k]
			if !s.Alive {
				continue
			}
			x, y := b.formationSpot(i, k)
			s.X, s.Y = x, y
			s.Z = b.Field.StandLevel(x, y)
			s.GoalX, s.GoalY, s.GoalZ = x, y, s.Z
			s.StepX, s.StepY, s.StepZ = x, y, s.Z
		}
	}
}

// formationSpot 算出第 side 側第 k 個兵的陣形位置。
//
// 原版 `sub_1AA2C`：`陣形線 + 陣形表[陣形編號][兵編號]`，
// 另一側 X 鏡射，最後夾在 1–62。
func (b *Battle) formationSpot(side, k int) (int, int) {
	dx, dy := b.Forms.Offset(b.Sides[side].Formation, k)
	if b.Sides[side].Mirror {
		dx = -dx
	}
	return clamp(b.Sides[side].Line + dx), clamp(Height/2 + dy)
}

// Order 對一整側或單一隊下命令。squad 為 −1 表示全軍
// （原版指令 3 的參數 7）。
//
// **命令 3（城壁移動）在沒有城的戰場自動變成命令 1（攻擊）**——
// 原版指令 3 與 13 都有這一段（docs/re/11 §3.7）。
func (b *Battle) Order(side, squad int, c Command) {
	if c == ScaleWal && !b.Field.IsSiege() {
		c = Attack
	}
	if c == Charge && b.Field.IsSiege() {
		// 說明書 4.2：突擊時守方會開門，而**開了就關不回去**。
		b.Sides[side].GateOpen = true
	}
	lo, hi := 0, SoldiersOnFoot
	if squad >= 0 {
		lo, hi = squad*PerSquad, (squad+1)*PerSquad
	}
	for i := lo; i < hi; i++ {
		if b.Sides[side].Soldiers[i].Alive {
			b.Sides[side].Soldiers[i].Next = c
		}
	}
}

// Step 推進一幀。
func (b *Battle) Step() {
	if b.Done {
		return
	}
	b.Frame++

	for i := range b.Sides {
		for k := range b.Sides[i].Soldiers {
			if b.Sides[i].Soldiers[k].Alive {
				b.updateSoldier(i, k)
			}
		}
	}
	b.stepProjectiles()
	b.reinforce()

	b.Advantage[0] = computeAdvantage(b.Sides[0].Fresh(), b.Sides[1].Fresh())
	b.Advantage[1] = computeAdvantage(b.Sides[1].Fresh(), b.Sides[0].Fresh())

	b.drainSiegeGeneral()
	b.checkGeneralRetreat()
	b.checkVictory()
}

// checkVictory 重現 `sub_1A6FA`：任一側補不出兵就結束。
func (b *Battle) checkVictory() {
	for i := range b.Sides {
		if b.Sides[i].Remaining() == 0 {
			b.Done, b.Winner = true, 1-i
			b.Log = append(b.Log, sideName(1-i)+"勝：對方補不出兵了")
			return
		}
	}
}

// checkGeneralRetreat 重現 `sub_1AE56`：大將體力低於 50 → 全軍退卻。
func (b *Battle) checkGeneralRetreat() {
	for i := range b.Sides {
		g := &b.Sides[i].Soldiers[0]
		if !g.Alive || g.HP >= GeneralRetreatHP {
			continue
		}
		if g.Cmd == Retreat {
			continue
		}
		b.Order(i, -1, Retreat)
		b.Log = append(b.Log, sideName(i)+"的大將體力不支，全軍退卻")
	}
}

// drainSiegeGeneral 重現 `sub_1AE56` 的後半：**攻城方的大將體力持續下降**。
// 說明書 6.1：「攻城戦の攻め手は体力が下がり続けます」——內建的攻城計時器。
func (b *Battle) drainSiegeGeneral() {
	if !b.Field.IsSiege() {
		return
	}
	b.siegeTick++
	if b.siegeTick < SiegeDrainInterval {
		return
	}
	b.siegeTick = 0
	g := &b.Sides[0].Soldiers[0] // 0 側是攻方
	if g.Alive && g.HP > 0 {
		g.HP--
	}
}

// reinforce 讓待機的兵補進場（說明書 4.1 的輸送帶）。
func (b *Battle) reinforce() {
	for i := range b.Sides {
		s := &b.Sides[i]
		for k := 0; k < Squads; k++ {
			if s.Reserve[k] == 0 {
				continue
			}
			for j := k * PerSquad; j < (k+1)*PerSquad; j++ {
				if s.Soldiers[j].Alive {
					continue
				}
				// 隊長那一格不補——大將倒下是全隊退卻，不是換人。
				if j == k*PerSquad {
					continue
				}
				x, y := b.formationSpot(i, j)
				s.Soldiers[j] = Soldier{
					Alive: true, Kind: s.Soldiers[k*PerSquad].Kind,
					HP: MaxHP, Stamina: StaminaFull, Target: -1,
					Cmd: Form, Next: s.Soldiers[k*PerSquad].Cmd,
					X: x, Y: y, Z: b.Field.StandLevel(x, y),
				}
				s.Soldiers[j].GoalX, s.Soldiers[j].GoalY = x, y
				s.Soldiers[j].StepX, s.Soldiers[j].StepY = x, y
				s.Reserve[k]--
				break
			}
		}
	}
}

func sideName(i int) string {
	if i == 0 {
		return "攻方"
	}
	return "守方"
}
