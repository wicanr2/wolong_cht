// Package diplomacy 實作交友度與外交官。
//
// 公式出自 `KI.EXE` 的 `sub_13E8E`（外交官每小時的工作）與
// `sub_1578F`（外交官要錢），反組譯見 docs/re/08 §3、docs/re/07 §16。
// 機制說明見 docs/mechanics/50-diplomacy.md。
package diplomacy

// Rand 是規則層共用的亂數介面。
type Rand interface{ Next() int }

// 交友度的儲存格式：**一個 byte，低 7 位是 0–100 的值，
// 最高位元是「非交戰」旗標**（1 ＝ 和平，0 ＝ 交戰中）。
//
// ⚠ **第一版我把這個位元的極性寫反了**（寫成 1 ＝ 交戰）。
// 從 SINARIO.DAT 讀出交友度矩陣之後才發現：四個劇本裡
// **每一格的最高位元都是 1，而且沒有任何勢力有侵攻目標**——
// 若 1 代表交戰，開局就會是所有人互相宣戰。
//
// 極性確定之後，`sub_1578F` 的要價基準也跟著對上了：
// 位元為 1（和平）→ 基準 100；為 0（交戰）→ 基準 125。
// **交戰中要更多錢**，與說明書「交戦中でも外交官を派遣し、
// 懐柔策を講じておかなければ」的語氣一致。
const (
	MaxFriendship = 100
	peaceBit      = 0x80
)

// Friendship 是「某個勢力看另一個勢力」的交友度。**這是單向的**：
// A 對 B 與 B 對 A 是兩個獨立的值，畫面上只看得到自家君主那一側。
type Friendship uint8

// Value 回傳 0–100 的交友值（去掉交戰旗標）。
func (f Friendship) Value() int { return int(f) & 0x7F }

// AtWar 回報是不是交戰中。**最高位元是「和平」而不是「交戰」**——
// 位元清掉才代表開戰。
func (f Friendship) AtWar() bool { return f&peaceBit == 0 }

// Raw 是**含和平位元**的原始值。說服判定要的是這個——
// 原版的門檻常數（`0x80 + 好戰 × 4 + 60` 之類）本身就帶著那個位元，
// 傳 Value() 進去會讓每一道門檻都偏 128。
func (f Friendship) Raw() int { return int(f) }

// WithValue 換掉數值但保留交戰旗標。
func (f Friendship) WithValue(v int) Friendship {
	if v < 0 {
		v = 0
	}
	if v > MaxFriendship {
		v = MaxFriendship
	}
	return Friendship(byte(v) | (byte(f) & peaceBit))
}

// WithWar 設定或清除交戰狀態。開戰是**清掉**和平位元。
func (f Friendship) WithWar(at bool) Friendship {
	if at {
		return f &^ peaceBit
	}
	return f | peaceBit
}

// Peace 是一個和平狀態的交友值（最高位元設起）。
// 直接寫 Friendship(50) 會得到「交戰中」，那幾乎不會是呼叫端要的。
func Peace(v int) Friendship { return Friendship(peaceBit).WithValue(v) }

// Level 是說明書 5.1 的六階顯示。
type Level int

const (
	Intimate Level = iota // 親密
	Good                  // 良好
	Normal                // 普通
	Bad                   // 險惡
	Worst                 // 最惡
	AtWar                 // 交戰
)

func (l Level) String() string {
	return [...]string{"親密", "良好", "普通", "險惡", "最惡", "交戰"}[l]
}

// Level 把交友值換成六階。
//
// ⚠ **六階的實際門檻還沒反組譯出來。** 這裡先用等分，
// 並標成 remake 的暫定值——顯示層可以先跑，等解出真值再換。
// 唯一確定的是「交戰」壓過一切（說明書：交戰中固定顯示交戰）。
func (f Friendship) Level() Level {
	if f.AtWar() {
		return AtWar
	}
	switch v := f.Value(); {
	case v >= 80:
		return Intimate
	case v >= 60:
		return Good
	case v >= 40:
		return Normal
	case v >= 20:
		return Bad
	default:
		return Worst
	}
}

// Diplomat 是一名派駐中的外交官。
type Diplomat struct {
	Politics int // 武將 +0x13
	Budget   int // 武將 +0x1A，玩家撥的經費餘額
}

// 外交官每次輪到時的動作機率。原版 `cmp al, 20h` → 32/256。
const workChanceOutOf256 = 0x20

// Tick 跑一次外交官的工作機會（原版是每「時」輪到該勢力時呼叫一次，
// 22 個勢力輪一圈，所以大約每天一次）。
//
// 回傳交友度是否提升了 1。
//
//	動作機率  12.5%
//	經費消耗  23 − 政治        ← 政治高的消耗慢
//	成功條件  rand(0..15) ≤ 政治 ← 政治高的成功率高
//	成功效果  交友度 +1，上限 100
//	經費歸零  完全停止工作
//
// **政治同時決定花費與效率**。說明書 5.3 只講了「要的錢比較少」，
// 機器碼顯示它也決定成功率。
func (d *Diplomat) Tick(f *Friendship, rng Rand) bool {
	if rng.Next()&0xFF >= workChanceOutOf256 {
		return false
	}
	if d.Budget <= 0 {
		return false // 經費用完 → 要再開口要錢才會動
	}
	cost := 23 - d.Politics
	if cost < 0 {
		cost = 0
	}
	d.Budget -= cost
	if d.Budget < 0 {
		d.Budget = 0
	}
	if rng.Next()&0x0F > d.Politics {
		return false // 這次沒有成果，但錢已經花掉了
	}
	if f.Value() >= MaxFriendship {
		return false
	}
	*f = f.WithValue(f.Value() + 1)
	return true
}

// Demand 回傳外交官這次會開口要多少錢（原版 `sub_1578F`）。
//
//	v    = min(交友度(我→對方), 交友度(對方→我))
//	base = 125 若 v 沒有交戰旗標，否則 100
//	金額 = (base − v 的數值) × 200
//
// **取兩個方向的較小者**——顯示是單向的，但外交官的工作量看比較差的那邊。
func Demand(mine, theirs Friendship) int {
	v := mine
	if theirs.Value() < mine.Value() {
		v = theirs
	}
	// 原版：`cmp al, 80h / jnb → 100`，否則 125。
	// 和平（位元為 1）→ 100；交戰（位元為 0）→ 125，**要更多錢**。
	base := 100
	if v.AtWar() {
		base = 125
	}
	n := base - v.Value()
	if n < 0 {
		n = 0
	}
	return n * 200
}

// CanRequestCooperation 回報協力外交的**局勢條件**是否成立。
//
// 說明書 7.4：
//
//	協力勢力がどこにも攻め込んでいない状態で、
//	協力勢力対自勢力の方が、協力勢力対侵攻勢力の方より交友関係がよい場合に成功します。
//
// ⭐ 本作的外交是**確定性**的（說明書 7.1）：條件滿足就成功，
// 與外交官能力無關；不滿足則重試永遠不會成功。**這裡不要擲骰子。**
//
// margin 是對方君主性格造成的額外門檻（好戰等級越高要求越多），
// 實際換算還沒反組譯出來 → 呼叫端先傳 0。
func CanRequestCooperation(allyInvadingSomeone bool, allyToUs, allyToTarget Friendship, margin int) bool {
	if allyInvadingSomeone {
		return false
	}
	return allyToUs.Value() > allyToTarget.Value()+margin
}

// CanCeaseFire 回報停戰的局勢條件（說明書 7.3）。
//
//	相手がこちらに侵攻してきている場合は停戦できません。
//
// 標準流程是先用協力外交讓對方轉向去打第三方，
// 對方就不再「正在侵攻我方」，停戰才談得成（說明書 10.1）。
func CanCeaseFire(theyAreInvadingUs bool) bool { return !theyAreInvadingUs }

// ---------------------------------------------------------------------------
// 侵攻狀態（勢力記錄 +0x19，docs/re/08 §1）
// ---------------------------------------------------------------------------

// NoTarget 表示這個勢力沒有侵攻對象。
const NoTarget = 0xFF

// CanSustainInvasion 回報一個勢力的財政撐不撐得起侵攻。
//
// 原版每小時輪到該勢力時檢查（`sub_13E11`）：
//
//	資金 >> 8  <  據點數 × 8 + 24   →  侵攻目標設回 0xFF
//
// **⭐ 財政撐不住就自動取消侵攻**，而且門檻與**據點數**掛鉤——
// 地盤越大，維持一場侵攻需要的最低資金越高。
//
// 這條規則是說明書幾個說服理由的底層機制：「敵が疲弊中」（資金 < 0）
// 之所以是提停戰的好時機，正是因為對方的侵攻**在機制上已經停了**
// （docs/mechanics/70-ai.md §3.11）。
//
// funds 傳有號的資金原值，程式內部照原版取 >> 8。
func CanSustainInvasion(funds, cities int) bool {
	return funds>>8 >= cities*8+24
}

// InvasionThreshold 回傳維持侵攻所需的最低資金（原值，不是 >> 8 之後的）。
// 只用於顯示與測試——判斷請用 CanSustainInvasion，那個才是照原版的比較法。
func InvasionThreshold(cities int) int { return (cities*8 + 24) << 8 }
