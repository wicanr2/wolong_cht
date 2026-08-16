// Package combat 實作戰略層的戰鬥自動判定。
//
// 公式全部出自 `KI.EXE` 的 `sub_15130`（勝負）、`sub_15285`／`sub_152D7`
// （戰力）、`sub_151B3`（傷亡與士氣）、`sub_1474A`（壞滅）與
// `sub_1291A`（武將的下場），完整反組譯見 docs/re/09，
// 機制說明見 docs/mechanics/30-combat.md。
//
// **這裡不含戰術戰鬥畫面。** 原版在玩家的勢力捲進去時先讓玩家選擇是否
// 親自指揮；選「委任」或玩家未捲入時都走這一套。兩條路的出口一致——勝方是誰、
// 哪一方壞滅——所以戰術層日後補上時，可以直接接在同一個介面後面。
package combat

import "github.com/wicanr2/wolong_cht/internal/rules/army"

// Rand 是規則層共用的亂數介面。原版是 `sub_1ECE0`，回傳 0–255。
type Rand interface{ Next() int }

// Mode 是戰鬥的場合。原版用它同時挑地形係數表的列與武將的適性欄位。
type Mode int

const (
	Siege Mode = iota // 攻城
	Field             // 野戰
)

// 地形係數表（`cs:5120h`，檔案偏移 0x5320）。四列 × 四欄，
// 欄位是兵種 − 1，第四欄沒有用到。
//
//	列 0  守城     騎 2  弓 3  步 3
//	列 1  野戰     騎 3  弓 2  步 1
//	列 2  沒有呼叫點會用到
//	列 3  攻城     騎 2  弓 1  步 2
//
// **列 2 沒有命名。** 戰略層傳得進去的只有 0、1、3；沒有證據之前
// 不要替它取名字（docs/re/09 §3.2）。
var terrainFactor = [4][4]int{
	{2, 3, 3, 0},
	{3, 2, 1, 0},
	{1, 3, 2, 0},
	{2, 1, 2, 0},
}

// rowFor 回傳某一方在某場合要用的地形係數列。
//
// 原版是這樣分出來的：`sub_15130` 先把攻城的 0 換成 3 再算攻方，
// 而算守方時 `al` 還是原始值 0 —— 因為 `sub_152D7` 出口會還原 ax。
func rowFor(m Mode, attacker bool) int {
	switch {
	case m == Field:
		return 1
	case attacker:
		return 3
	default:
		return 0
	}
}

// Leader 是帶兵的武將。欄位對應武將記錄（docs/formats/08 §3）。
type Leader struct {
	Martial int // +0x11 武力，1–15
	Command int // +0x12 統率，1–15

	// SiegeAptitude 與 FieldAptitude 是 +0x0E／+0x0F 的**高半位元組**（0–10）。
	//
	// ⚠ 這兩個欄位原本被記成「兵種適性」。索引它們的是場合（攻城 0／
	// 野戰 1），不是兵種 —— 見 docs/re/09 §3.3。
	SiegeAptitude, FieldAptitude int

	// Rating 是評價（＝適性和 ＋ 2×武力 ＋ 2×統率，見 internal/rules/general）。
	// 只在戰敗被擒的判定裡用到。
	Rating int
}

func (l Leader) aptitude(m Mode) int {
	if m == Field {
		return l.FieldAptitude
	}
	return l.SiegeAptitude
}

// Unit 是軍團六個編成位置之一。對應軍團記錄 +0x28 起的 4-byte 槽。
type Unit struct {
	Men  int            // +1，一格 10 人（滿編 100 ＝ 1000 人）
	Kind army.TroopType // +2 減 1
}

// Corps 是參戰的一方。只帶戰鬥要用的欄位，位置與行軍留在 army。
type Corps struct {
	Faction int
	Leader  Leader
	Units   [army.Positions]Unit
	Morale  int // +0x06
	Men     int // +0x04，總兵力
}

// leaderValue 是將領修正（原版 `sub_152D7` 前半段）。
//
//	武力 ≥ 統率  →  75% 用 武力 × 2、25% 用 武力 × 3/4 ＋ 統率
//	武力 < 統率  →  統率 × 3/4 ＋ 統率
//
// 那個 25% 的分支是原版擲 `rand & 3` 擲出來的，**不是筆誤**：
// 武力高的將領大多數時候發揮兩倍武力，偶爾退回混合值。
func leaderValue(l Leader, rng Rand) int {
	m, c := l.Martial, l.Command
	if m >= c {
		if rng.Next()&3 != 0 {
			return m * 2
		}
		c = m // 原版 mov dh, dl
	}
	return c - c>>2 + c
}

// Power 回傳一方的戰力。
//
//	基礎 ＝ (士氣 ÷ 8) × Σ(各槽兵力 × 地形係數[兵種])
//	戰力 ＝ 基礎 × 將領值 ÷ (64 × (16 − 適性))
//
// garrison 是守城方額外加上的城兵數（原版 `sub_15285` 的
// `add dl, [bx+13h]`，只有攻城的守方會有）。
//
// ⚠ 原版的基礎戰力算在 16 位裡。極端編成（六槽滿員、士氣 255）會溢位；
// 實際值域下不會發生，這裡不重現溢位。
func Power(c Corps, m Mode, attacker bool, garrison int, rng Rand) int {
	row := terrainFactor[rowFor(m, attacker)]
	sum := 0
	for _, u := range c.Units {
		sum += u.Men * row[int(u.Kind)]
	}
	if !attacker && m == Siege {
		sum += garrison
	}
	base := (c.Morale >> 3) * sum

	// 原版分兩步、中間截斷一次：先算 `將領值 × 16 ÷ (16 − 適性)`，
	// 再乘上基礎戰力並右移 10。**不要合併成一個除式**，
	// 中間那次整數截斷會改變結果。
	factor := leaderValue(c.Leader, rng) * 16 / (16 - c.Leader.aptitude(m))
	return base * factor >> 10
}

// Result 是一場自動判定的結果。
type Result struct {
	// DefenderWins 為真表示守方（di）贏了。原版用 al：0 ＝ 攻方、1 ＝ 守方。
	DefenderWins bool

	// Ratio 是戰力比值 ＝ 強的一方 × 8 ÷ 弱的一方，上限 100。
	// 它同時決定敗方的傷亡量與城市的損傷量。
	Ratio int

	// AttackerDestroyed／DefenderDestroyed 對應原版回傳的 ah bit 0／bit 1。
	AttackerDestroyed, DefenderDestroyed bool

	// CityDamage 是攻城戰扣掉的城兵／上昇值／防災值（三者同量）。
	CityDamage int
}

// Resolve 跑一場自動判定。attacker 與 defender 會被就地改寫
// （兵力、士氣），呼叫端要傳指標。
//
// **勝負不擲骰**：戰力大的一方直接贏。骰子只出現在傷亡量與將領修正。
func Resolve(attacker, defender *Corps, m Mode, garrison int, rng Rand) Result {
	pa := Power(*attacker, m, true, garrison, rng) + 8
	pd := Power(*defender, m, false, garrison, rng) + 8

	hi, lo := pa, pd
	defenderWins := pd > pa
	if defenderWins {
		hi, lo = pd, pa
	}
	ratio := hi * 8 / lo
	if ratio > 100 {
		ratio = 100
	}

	win, lose := attacker, defender
	if defenderWins {
		win, lose = defender, attacker
	}

	r := Result{DefenderWins: defenderWins, Ratio: ratio}
	if m == Siege {
		r.CityDamage = CityDamage(ratio)
	}
	applyCasualties(win, lose, ratio, rng)

	r.AttackerDestroyed = Destroyed(*attacker)
	r.DefenderDestroyed = Destroyed(*defender)
	return r
}

// CityDamage 是攻城戰對據點造成的損傷，城兵／上昇值／防災值各扣這麼多。
//
//	損傷 ＝ (63 − 比值) >> 2      （byte 運算）
//
// ⚠ **比值超過 63 會減出負數繞回去**，損傷反而跳到最大。
// 設計意圖是「苦戰傷城多、輾壓傷城少」，而 8 倍以上的懸殊戰力
// 會落進繞回區。這是原版行為，這裡照做（docs/re/09 §4.1）。
func CityDamage(ratio int) int { return int(byte(0x3F-ratio)) >> 2 }

// applyCasualties 是原版 `sub_151B3` 的後半段。
func applyCasualties(win, lose *Corps, ratio int, rng Rand) {
	oldWin, oldLose := win.Men, lose.Men
	newWin, newLose := 0, 0

	// 原版的除數在迴圈裡遞增（`inc ch`），所以第 i 槽除的是 比值 + i。
	// 照抄，不要化簡成同一個除數。
	div := ratio
	for i := 0; i < army.Positions; i++ {
		loss := rng.Next()&7 + 2 // 勝方：2–9，與比值無關
		win.Units[i].Men = shrink(win.Units[i].Men, loss, i == 0)
		newWin += win.Units[i].Men

		div++
		loss = rng.Next()%div + 8 // 敗方：8 ＋ rand mod (比值 + i)
		lose.Units[i].Men = shrink(lose.Units[i].Men, loss, i == 0)
		newLose += lose.Units[i].Men
	}
	win.Men, lose.Men = newWin, newLose

	// 勝方按戰前士氣衰減，敗方一律重設成 100 × 兵力比。
	// 戰前士氣不足 100 的一方直接歸零 —— 勝方也一樣。
	win.Morale = scaleMorale(win.Morale, win.Morale, newWin, oldWin)
	lose.Morale = scaleMorale(lose.Morale, 100, newLose, oldLose)
}

// shrink 扣掉一個槽的兵力。**第一槽（大將）永遠留 1**：
// 自動判定打不死大將的部隊（原版 `cmp cl, 6`）。
func shrink(men, loss int, isGeneralSlot bool) int {
	if men-loss > 0 {
		return men - loss
	}
	if isGeneralSlot {
		return 1
	}
	return 0
}

// scaleMorale 是戰後士氣。before 是戰前士氣（決定歸不歸零），
// base 是要被縮放的基準（勝方用自己的士氣，敗方固定 100）。
func scaleMorale(before, base, now, was int) int {
	if was == 0 || before < army.RoutMoraleGate {
		return 0
	}
	return base * now / was
}

// Destroyed 回報軍團戰後是否壞滅（原版 `sub_1474A` 的**前兩個**檢查）。
//
// 自動判定裡大將槽保底 1，所以這兩條裡實際上只有士氣 0 會成立，
// 也就是「戰前士氣不足 100」。大將槽那一條是留給戰術層的。
//
// ⚠ **壞滅還有第三個入口：敗方找不到退路**（`sub_1487B`）。
// 那一條需要道路圖與據點歸屬，所以在 `internal/state`
// （`retreatOrPerish`，docs/spec/46），不在這一層。
func Destroyed(c Corps) bool { return c.Morale == 0 || c.Units[0].Men == 0 }

// ---------------------------------------------------------------------------
// 敗方武將的下場（`sub_1291A`）
// ---------------------------------------------------------------------------

// NeutralFaction 是「無主」的勢力編號。勢力只有 22 個（0–21），
// 原版拿 0x18 當空城與中立方的哨兵值（`sub_14CF3`／`sub_1291A`）。
const NeutralFaction = 0x18

// Fate 是軍團壞滅後武將的去向。
type Fate int

const (
	Escaped  Fate = iota // 逃脫：軍團解散，武將保住
	Captured             // 被擒：改隸勝方，舊主記進 +0x1D
	Suicide              // 自刎：舊主已滅且武將不事二主
)

// Captive 是判定武將下場需要的狀態。
type Captive struct {
	Rating       int  // 武將評價
	IsRuler      bool // 是不是敗方的君主
	HasCapital   bool // 敗方還有沒有首都
	LoyalToDeath bool // 武將旗標 bit 4（不事二主）
	LordSurvives bool // 舊主勢力是否還存在
}

// escapeThreshold 是逃脫判定的門檻：`rand(0..127) ≤ 評價 ÷ 2 + 40` 就逃掉。
func escapeThreshold(rating int) int { return rating>>1 + 40 }

// RollFate 判定敗方武將的下場。victor 是勝方的勢力編號。
//
// 三個必逃的例外（原版直接跳過擲骰）：
//
//   - **君主親征絕不被擒**
//   - 勝方與敗方同勢力
//   - 勝方是無主勢力
//
// 一個必被擒的例外：敗方**沒有首都**，無處可退。
//
// 逃脫機率 ＝ (評價 ÷ 2 + 41) ÷ 128。呂布、趙雲（評價 66）約 58%，
// 文官（評價 10）約 36%——**能力越高越容易脫身**。
func RollFate(c Captive, victor, loser int, rng Rand) Fate {
	if !c.HasCapital {
		return capturedOrSuicide(c)
	}
	if c.IsRuler || victor == loser || victor == NeutralFaction {
		return Escaped
	}
	if rng.Next()&0x7F > escapeThreshold(c.Rating) {
		return capturedOrSuicide(c)
	}
	return Escaped
}

// capturedOrSuicide 套用「不事二主」：舊主已滅 ＋ 武將旗標 bit 4 → 自刎。
// 原版訊息 0x43「{1}在即將被我軍擒拿之前，自刎而死了」。
func capturedOrSuicide(c Captive) Fate {
	if !c.LordSurvives && c.LoyalToDeath {
		return Suicide
	}
	return Captured
}

// ---------------------------------------------------------------------------
// 城兵（`sub_14F8A`）
// ---------------------------------------------------------------------------

// GarrisonLeader 是城兵的將領編號：127 號，也就是武將表最後那筆空記錄。
const GarrisonLeader = 0x7F

// GarrisonMorale 是城兵的士氣，原版直接寫 0xFF。
const GarrisonMorale = 0xFF

// Garrison 把據點的城兵數攤成一支臨時軍團：**六隊步兵、士氣 255、
// 將領是空記錄**。餘數散給前幾槽。
//
// 原版把它搭在 `ds:4200h`——軍團表與武將表中間那 64 byte，
// 戰後 `sub_14FC8` 清掉。
func Garrison(faction, men int) Corps {
	c := Corps{Faction: faction, Morale: GarrisonMorale, Men: men}
	each, extra := men/army.Positions, men%army.Positions
	for i := range c.Units {
		c.Units[i] = Unit{Men: each, Kind: army.Infantry}
		if i < extra {
			c.Units[i].Men++
		}
	}
	return c
}

// ---------------------------------------------------------------------------
// 補給與士氣回復（`sub_12600`）
// ---------------------------------------------------------------------------

// Upkeep 回傳軍團這一 tick 的軍費。
//
//	據點或道路上   兵力 ÷ 32 ＋ 1
//	野外           兵力 × 3/4
//
// **差距 24 倍。** 野外駐留貴到不可能長期維持，這是原版逼軍團回城的手段。
func Upkeep(men int, inField bool) int {
	if inField {
		return men>>1 + men>>2
	}
	return men>>5 + 1
}

// MoraleRegen 是每 tick 的士氣回復量。只有在據點或道路上才回。
const MoraleRegen = 10

// Recover 讓軍團回一次士氣，上限是勢力的士氣基準（勢力記錄 +0x1D，開局 200）。
// 在野外不回復。
func Recover(c *Corps, factionMorale int, inField bool) {
	if inField {
		return
	}
	c.Morale += MoraleRegen
	if c.Morale > factionMorale {
		c.Morale = factionMorale
	}
}
