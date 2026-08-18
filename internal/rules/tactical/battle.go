package tactical

// 一場戰術戰鬥。每幀跑一次 Step，直到 Done。

// Rand 是規則層共用的亂數介面。
type Rand interface{ Next() int }

// Side 是交戰的一方。
type Side struct {
	// Soldiers 是場上的 48 個兵：第 k 隊佔 [k*8, k*8+8)，**每隊的第一個是隊長**
	// （`sub_1A754` 的 1 + 7 迴圈）。
	Soldiers [SoldiersOnFoot]Soldier

	// Power 是這一側的戰力（原版由士氣算進每個兵的 `+0x18`，§3.9）。
	Power int

	// Kinds[k] 是第 k 隊的兵種。
	//
	// ⚠ **不能拿隊長的兵種當一隊的兵種。** 第 0 隊的隊長是大將（兵種 0），
	// 照隊長補兵會讓整隊變成大將——而大將不攻擊也不會陣亡，一整隊
	// 站在那裡不動，戰鬥就永遠打不完。
	Kinds [Squads]Kind

	// Morale 是這一支軍團的士氣（原版 `word_1D30A:+0x06`，開打時
	// 由 `sub_19E70` 從軍團記錄抄進來）。
	//
	// ⭐ **它同時是每個兵的開場體力**（`sub_19B6D` 的 `mov es:[di+3], ah`，
	// docs/spec/61）：士氣高的軍團，每個兵開場血就厚。`MaxHP` 是
	// `sub_1B97E` 的**回復**上限，不是開場值——兩者可以不同。
	Morale int

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

	// Standing 是這一側目前的常令。補進場的兵接這個，**不是接隊長當下的命令**——
	// 隊長正在退卻時，新兵一出來就跟著往回走，整個待機池會在幾百幀內流光。
	Standing Command

	// Sortied 記錄守方是否已經下過突擊令（`sub_1B7CB` 那一發）。
	// 說明書 4.2 說「開了就關不回去」，所以這是不可逆的。
	Sortied bool

	// Selected 是底列六格的選取位元圖（原版 `byte_1D310` 的位元 0–5，
	// docs/spec/33）。**下完令就清空**，而且 0 表示「全隊」不是「沒有隊」。
	//
	// 原版只有玩家那一側有這個位元圖（`sub_1A8CC` 的 `cmp si, 600h`），
	// 這裡放在 Side 上，AI 那一側永遠是 0。
	Selected uint8
}

// AllSquadsMask 是六隊全選（原版 `xchg` 取到 0 時代入的 `0xFF` 的有效位）。
const AllSquadsMask = 1<<Squads - 1

// 攻守的側別是固定的：側 0 是攻方、側 1 是守方。
//
// `hitStructure` 的「守方碰不壞自己的城壁」與 `Order` 的「守方突擊才拆城壁」
// 都依賴這個約定。⚠ **原版的側別不是這樣**——原版側 0 恆為玩家、側 1 恆為
// 對方（docs/re/60 §5），攻守由 `sub_14ED7` 互換兩個軍團指標來表達。
// remake 選了「側別 ＝ 攻守」這一種，所以凡是原版寫「側 0／側 1」的地方，
// 都要先換算成攻守再對照。
const (
	AttackerSide = 0
	DefenderSide = 1
)

// SquadSelected 回傳第 k 隊有沒有被選中。
func (s *Side) SquadSelected(k int) bool {
	if k < 0 || k >= Squads {
		return false
	}
	return s.Selected&(1<<uint(k)) != 0
}

// ToggleSquadSelection 重現熱區 0x15–0x1A 的 `xor byte_1D310, 1<<k`。
func (b *Battle) ToggleSquadSelection(side, squad int) {
	if side < 0 || side >= len(b.Sides) || squad < 0 || squad >= Squads {
		return
	}
	b.Sides[side].Selected ^= 1 << uint(squad)
}

// TakeSquadSelection 重現 `0001C1CE` 的 `xchg`：取出選取位元圖並清空，
// **取到 0 就回全隊**。
func (b *Battle) TakeSquadSelection(side int) uint8 {
	if side < 0 || side >= len(b.Sides) {
		return AllSquadsMask
	}
	mask := b.Sides[side].Selected & AllSquadsMask
	b.Sides[side].Selected = 0
	if mask == 0 {
		return AllSquadsMask
	}
	return mask
}

// OrderSelected 是**玩家**下令那一條路（`0001C1B9`）：命令只送給被選中的隊，
// 一隊都沒選就送給全隊，送完把選取清空。
//
// ⚠ 回 false 表示**被拒絕**——城壁令在沒有城的戰場，原版是跳
// 「這哪裡有城壁啊！！」而不下令。這與腳本那一條（指令 3 自動降級成
// 指令 1，見 Order）是兩條不同的路，不要合併。
func (b *Battle) OrderSelected(side int, c Command) bool {
	if side < 0 || side >= len(b.Sides) {
		return false
	}
	if c == ScaleWal && !b.Field.IsSiege() {
		return false
	}
	mask := b.TakeSquadSelection(side)
	if mask == AllSquadsMask {
		b.Order(side, -1, c)
		return true
	}
	for k := 0; k < Squads; k++ {
		if mask&(1<<uint(k)) != 0 {
			b.Order(side, k, c)
		}
	}
	return true
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
	// PlayerSide 是玩家在哪一側（0 攻／1 守）。**原版的 side 0 永遠是玩家**，
	// remake 的 Sides[0] 固定是攻方，所以要另外記——見 SetPlayerSide。
	PlayerSide int

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

	// bar 是「門強度」條的狀態（docs/spec/32）。它是純呈現用的暫態，
	// 不進存檔。
	bar structureBar

	// Structures 是城壁與門，最多 16 段（docs/re/11 §5.9）。
	Structures []Structure

	rng Rand

	// projectiles 是飛在空中的箭。原版是一張 32 筆的表（docs/re/11 §5.1）。
	projectiles []projectile

	// sfx 是這一輪要送給音源 TSR 的效果碼。原版走 INT 61h `AH=0x05`、
	// `AL=`效果碼（docs/re/17 §3），碼本身就是 `SOUND.DAT` 的記錄編號
	// （docs/re/57 §6）。規則層只排隊、不播——**它不認識畫面也不認識喇叭**。
	sfx []uint8

	// siegeTick 是攻城方大將體力遞減的計時器。
	siegeTick int

	// scripts[i] 不是 nil 時，第 i 側由 `BATTLE.DAT` 的腳本驅動。
	// 玩家那一側留 nil，由畫面層下命令。
	scripts [2]*Script

	// Log 記錄值得回報的事件，供呼叫端呈現。
	Log []string
}

type projectile struct {
	side       int // 誰射的
	x, y, z    int
	tx, ty, tz int
	power      int
	special    bool
	// specialFrame 是原版特殊投射物 raw 0x214／0x215 的低位。它取自
	// 發射兵記錄 +0x02 bit 0，不能用方向或 side 代替。
	specialFrame uint8
	// direction 是原版飛道具記錄 +0x05：0 西、1 北、2 東、3 南；
	// CH=0x20 會再加上 0x80，sub_1BA2E 因而只更新高度。
	direction int
	// gridIndex／previousGridIndex 對應 +0x10／+0x12；畫面層尚未使用
	// previous*，但保留它們讓每幀狀態可直接對照 sub_1BAB7。
	gridIndex, previousGridIndex    int
	previousX, previousY, previousZ int
	heightFP                        int
	velocityFP                      int
}

// 原版已證實的三個戰術效果碼（docs/re/17 §3）。
const (
	SFXSpecialLaunch uint8 = 0x0A // sub_1AD7F 特殊投射物發射
	SFXProjectileHit uint8 = 0x0B // sub_1B97E 投射物命中
	SFXNormalLaunch  uint8 = 0x0C // sub_1AD2D 普通投射物發射
)

func (b *Battle) emitSFX(code uint8) {
	if b == nil {
		return
	}
	// 同一幀重複的碼只留一個：原版一次只送一個 code 給 TSR，
	// 而 TSR 只有三個 2-operator 通道（docs/re/57 §2）。
	for _, c := range b.sfx {
		if c == code {
			return
		}
	}
	b.sfx = append(b.sfx, code)
}

// TakeSoundEffects 取走並清空這一輪排隊的效果碼。
func (b *Battle) TakeSoundEffects() []uint8 {
	if b == nil || len(b.sfx) == 0 {
		return nil
	}
	out := b.sfx
	b.sfx = nil
	return out
}

// ProjectileView 是畫面層可讀的飛道具快照；不暴露規則層的可變記錄。
type ProjectileView struct {
	Side, X, Y, Z                   int
	PreviousX, PreviousY, PreviousZ int
	Power                           int
	Direction                       int
	Special                         bool
	SpecialFrame                    int
}

// Projectiles 回傳目前仍在場上的飛道具，供戰術畫面以原生格座標繪製。
func (b *Battle) Projectiles() []ProjectileView {
	if b == nil || len(b.projectiles) == 0 {
		return nil
	}
	out := make([]ProjectileView, len(b.projectiles))
	for i, p := range b.projectiles {
		out[i] = ProjectileView{
			Side: p.side, X: p.x, Y: p.y, Z: p.z,
			PreviousX: p.previousX, PreviousY: p.previousY, PreviousZ: p.previousZ,
			Power: p.power, Direction: p.direction, Special: p.special,
			SpecialFrame: int(p.specialFrame & 1),
		}
	}
	return out
}

// NewBattle 開一場戰鬥。
//
// cityWall 是攻城時守方據點的城壁值（據點記錄 `+0x13`），決定城壁耐久；
// 野戰用不到，傳 0 即可。
func NewBattle(f *Field, forms *Formations, rng Rand, cityWall int) *Battle {
	b := &Battle{Field: f, Forms: forms, rng: rng, Winner: -1}
	// 原版把兩側的陣形原點分開存（side 0 → word_1D33C，side 1 → word_1D33E），
	// 而 side 1 是**把陣形表的 dx 取負**來鏡射（`sub_1AA2C` 的 `neg dl`），
	// 不是把原點對稱過去。
	b.SetPlayerSide(AttackerSide)
	if f != nil {
		b.Structures = buildStructures(f.tiles, f.IsSiege(), cityWall)
	}
	return b
}

// SetPlayerSide 指定哪一側是玩家。**要在 Place() 之前呼叫。**
//
// ⭐ 原版的 side 0 **永遠是玩家**：`sub_14E5C`（野戰）與 `sub_14ED7`（攻城）
// 在玩家是守方那一支互換 `word_10D2E`／`word_10D30` 並設 `byte_10D35` bit 7
// （docs/re/11 §5.11、docs/re/58 §4）。而陣形原點與鏡射綁的是
// **玩家／腳本**不是攻方／守方——`word_1D33C`（玩家）＝ X 5、
// `word_1D33E`（腳本）＝ X 58。
//
// remake 的 `Sides[0]` 固定是攻方（城壁、突擊、優勢度的判定都靠這個），
// 所以玩家守城時要把陣形原點與鏡射這兩樣換過來，其餘不動。
//
// 攻城時戰場另外會轉 180 度（docs/spec/56），兩件事合起來才對得上原版：
// 玩家永遠從 X=5 那一端出發，而地形被轉到讓 X=5 落在該落的地方。
func (b *Battle) SetPlayerSide(side int) {
	if side != DefenderSide {
		side = AttackerSide
	}
	b.PlayerSide = side
	b.Sides[side].Line, b.Sides[side].Mirror = LineFor(0, 0), false
	b.Sides[1-side].Line, b.Sides[1-side].Mirror = LineFor(1, 0), true
}

// Role 回傳這一側在**原版的框**裡是哪一個角色：0 ＝ 玩家、1 ＝ 腳本。
//
// ⭐ 陣形原點那兩組常數（`lineX`）綁的是**玩家／腳本**，不是攻方／守方，
// 也不是 remake 的側編號——原版的 side 0 永遠是玩家（docs/spec/56 §1）。
// remake 的 side 0 是攻方，玩家守城時就對不上，拿側編號去查表會讓
// **腳本把自己的陣形線設成玩家那一條（X=5）**，整支軍團往對面走。
func (b *Battle) Role(side int) int {
	if side == b.PlayerSide {
		return 0
	}
	return 1
}

// startHP 是這一側新上場的兵的體力：**軍團士氣**（docs/spec/61）。
//
// `MaxHP` 是 `sub_1B97E` 的**回復**上限，不是開場值——士氣 200 的軍團
// 開場每個兵 200 點，掉到 100 以下才開始回，回也只回到 100。
// 呼叫端沒給士氣時（測試、無頭模擬）退回 DefaultPower，與 Power 同一套預設。
func (s *Side) startHP() int {
	if s.Morale > 0 {
		return s.Morale
	}
	return DefaultPower
}

// Deploy 把一隊放上戰場。kind 是六個編成位置的兵種，men 是各位置的兵數
// （以「兵」為單位，一個位置滿編 100）。
//
// 場上一隊只放 8 個，其餘進待機（說明書 4.1）。
func (b *Battle) Deploy(side int, squad int, kind Kind, men int) {
	s := &b.Sides[side]
	if s.Power == 0 {
		// 沒有指定就用滿值。原版的 `+0x18` 是由軍團士氣算出來的（§3.9）；
		// 呼叫端沒給士氣時（測試、無頭模擬）用這個。
		s.Power = DefaultPower
	}
	s.Kinds[squad] = kind
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
			Alive: true, Kind: k, HP: s.startHP(), Stamina: StaminaFull,
			Power: s.Power, Target: -1, Cmd: Form, Next: Form,
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
			s.syncTerrain(b.Field, x, y, s.Z)
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
	return clamp(b.Sides[side].Line + dx), clampY(OriginY + dy)
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
	// ⚠ 只有**守方**的突擊會拆城壁（`sub_1B7CB` 的側別閘，
	// 見 SortieBreaksWalls）。攻方突擊只是加快疲勞、火力全開。
	if c == Charge && b.Field.IsSiege() && side == DefenderSide {
		b.Sides[side].Sortied = true
		b.SortieBreaksWalls()
	}
	lo, hi := 0, SoldiersOnFoot
	if squad >= 0 {
		lo, hi = squad*PerSquad, (squad+1)*PerSquad
	} else if c != Retreat {
		// 退卻是被動觸發的（受傷、大將不支），不該變成常令。
		b.Sides[side].Standing = c
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
	// 原版在同一個地方（0001A12A）遞增計時器並檢查三個 UI 的到期。
	b.expireStructureBar()

	// 腳本先跑：原版的主迴圈是「執行一個腳本指令 → 更新實體」。
	for i := range b.scripts {
		if b.scripts[i] != nil {
			b.scripts[i].Step(b)
		}
	}

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
					Alive: true, Kind: s.Kinds[k], Power: s.Power,
					HP: s.startHP(), Stamina: StaminaFull, Target: -1,
					Cmd: Form, Next: s.Standing,
					X: x, Y: y, Z: b.Field.StandLevel(x, y),
				}
				s.Soldiers[j].syncTerrain(b.Field, x, y, s.Soldiers[j].Z)
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

// ---------------------------------------------------------------------------
// 與戰略層的換算
// ---------------------------------------------------------------------------

// DefaultPower 是沒有指定士氣時每個兵的戰力。
const DefaultPower = 100

// MenPerSoldier 是戰場上一個兵等於戰略上幾個人。
//
// 說明書 4.1：「戦術では、**1兵士が戦略の兵数10人分に相当**します」。
const MenPerSoldier = 10

// Outcome 是一場戰術戰鬥的結果，換算回戰略層的單位。
type Outcome struct {
	AttackerWins bool
	// Men[i] 是第 i 側剩下的人數（已經 × 10 換回戰略單位）。
	Men [2]int
	// GeneralHP[i] 是第 i 側大將的體力，供「被擒／逃脫」之類的後續判定參考。
	GeneralHP [2]int
	Frames    int
}

// Result 把戰鬥結果換算成戰略層看得懂的數字。
func (b *Battle) Result() Outcome {
	o := Outcome{AttackerWins: b.Winner == 0, Frames: b.Frame}
	for i := range b.Sides {
		o.Men[i] = b.Sides[i].Remaining() * MenPerSoldier
		o.GeneralHP[i] = b.Sides[i].Soldiers[0].HP
	}
	return o
}

// Run 把戰鬥一路跑到結束，回傳結果。
//
// 給無頭環境（測試、cmd/wlsim）用；有畫面的呼叫端應該自己每幀呼叫 Step，
// 因為**戰鬥中不能停時間**（說明書 4.1），畫面要跟著跑。
//
// maxFrames 是保險絲：真的跑不完就判平手（回傳 false）。
func (b *Battle) Run(maxFrames int) bool {
	for i := 0; i < maxFrames && !b.Done; i++ {
		b.Step()
	}
	return b.Done
}
