// Package economy 實作原版的月結。
//
// 每一條公式都是從 `KI.EXE` 的 `sub_15358` 及其子程式讀出來的，
// 完整反組譯見 docs/re/07，機制說明見 docs/mechanics/40-economy.md。
//
// 這一層刻意**不含畫面、不含存檔、不含 AI 的高階決策**，
// 只做「給定勢力與據點的狀態，跑一次月結」這件事。
// 這樣它可以被大量重複呼叫來驗證長期經濟行為，不必起遊戲。
package economy

// 資金的上下限。原版 `sub_15609`／`sub_1563B` 把資金鉗在 ±655,000
// （24 位有號值 0x09FE98 與 0xF60168）。
const (
	MaxFunds = 655000
	MinFunds = -655000

	// MaxReserve 是**每個兵種**的預備兵上限（原版 `sub_155EC` 的 0xFFDC）。
	MaxReserve = 65500
)

// TroopType 是三個兵種。順序與勢力記錄 +0x04／+0x06／+0x08 一致。
type TroopType int

const (
	Cavalry  TroopType = iota // 騎馬
	Archer                    // 弓兵
	Infantry                  // 步兵
	NumTroopTypes
)

// City 是一個據點在月結時會用到的部分。
// 欄位對應據點記錄的 +8／+0Ah／+0Eh（docs/re/07 §3）。
type City struct {
	X, Y       int // 大地圖格座標（0–383, 0–255）
	Production int // 生產力
	Owner      int // 所屬勢力編號
}

// Faction 是一個勢力在月結時會用到的部分。
type Faction struct {
	Funds    int                // 資金，可為負
	Reserves [NumTroopTypes]int // 預備兵
	Capital  City               // 首都（只需要座標）
	Cities   int                // 據點數，月結時重算後寫回

	// TaxRate 是「今月生效」的稅率（百分比）。
	// 玩家在財政視窗改的是下個月的值，月結時才搬過來，
	// 這裡拿到的已經是搬好的（docs/re/07 §7）。
	TaxRate int

	// RecruitCap 是三兵種各自的募兵數設定。
	// **這是上限不是目標**——實際募到多少由據點決定，取兩者較小者。
	RecruitCap [NumTroopTypes]int

	// Expense 是本月累計的支出。月結扣完之後歸零。
	Expense int

	// AI 為 true 時不使用 TaxRate，收入固定除以 2
	// （原版 `sub_15456` 對非玩家勢力就是 `shr`/`rcr` 一次）。
	AI bool
}

// distanceDivisor 回傳據點對首都的距離除數。
//
// 原版用**切比雪夫距離**（`max(|Δx|, |Δy|)`，不是歐氏也不是曼哈頓），
// 夾到 255 之後查 ds:5532h／ds:5535h 兩張表：
//
//	距離 ≤ 80  → 2
//	距離 ≤ 200 → 3
//	其餘        → 4
//
// **遠方據點的收入只有近處的一半。** 說明書完全沒提這件事，
// 但它讓「遷都到領土重心」變成一個真正的經濟決策。
func distanceDivisor(c, capital City) int {
	dx := c.X - capital.X
	if dx < 0 {
		dx = -dx
	}
	dy := c.Y - capital.Y
	if dy < 0 {
		dy = -dy
	}
	d := dx
	if dy > d {
		d = dy
	}
	if d > 255 {
		d = 255 // 原版 `and ah,ah / jnz / mov al,0FFh`
	}
	switch {
	case d <= 80:
		return 2
	case d <= 200:
		return 3
	default:
		return 4
	}
}

// splitRecruits 把一個據點的募兵基數 a 依地域拆成三兵種。
//
// 分區只看 **Y 座標**，與 X 無關；門檻是 80 與 150（大地圖高 256）。
// 概略比例（分母 32）：
//
//	北方 y < 80         騎馬 19（59%）／弓  1（3%）／步 12（38%）
//	中部 80 ≤ y < 150   騎馬  4（12%）／弓  4（12%）／步 24（75%）
//	南方 y ≥ 150        騎馬  1（3%）／弓 16（50%）／步 15（47%）
//
// 說明書只寫了「騎馬は北方でよく集まり、弓は南方でよく集まります」，
// **中部是步兵的天下（75%）這件事說明書沒說**。
//
// ⚠ 這裡照抄原版的移位序列，**不寫成 `a*19/32` 這種分數式**。
// 兩者在 a 不是 32 的倍數時會分岔：原版是「先算兩份、剩下的全給第三份」，
// 餘數被吸收；分數式則三份各自捨去，總和會少。
// 踩過一次：a=33 時原版給 (20,1,12)＝33，分數式給 (19,1,12)＝32。
func splitRecruits(a, y int) (cavalry, archer, infantry int) {
	switch {
	case y < 80: // 北方
		// dx=a>>2; cx=dx>>1; dx+=cx; cx>>=2; ax=a-dx-cx
		d := a >> 2
		c := d >> 1
		d += c
		c >>= 2
		return a - d - c, c, d

	case y < 150: // 中部
		// dx=a; ax=a>>3; cx=ax; dx-=ax; dx-=cx
		x := a >> 3
		return x, x, a - x - x

	default: // 南方
		// dx=a>>1; ax=a>>5; cx=dx; dx-=ax
		half := a >> 1
		x := a >> 5
		return x, half, half - x
	}
}

// Result 是一次月結的明細，方便測試與畫面顯示。
type Result struct {
	Income    int                // 實際入帳的收入（已套稅率）
	GrossBase int                // 套稅率之前的 Σ(生產力 ÷ 距離除數)
	Recruited [NumTroopTypes]int // 實際募到的兵（已套上限）
	Deficit   [NumTroopTypes]int // 赤字懲罰扣掉的兵
	Cities    int                // 重算後的據點數
}

// Rand 是月結需要的亂數來源。原版赤字懲罰會呼叫 `sub_1ECE0`，
// 這裡抽出來讓測試可以給定序列。回傳值只取低 5 位（0–31）。
type Rand interface{ Next() int }

// Settle 跑一次月結，就地更新 f，並回傳明細。
//
// 順序照原版 `sub_15358`（docs/re/07 §1）：
//
//	① 先扣本月累計支出
//	② 據點結算（收入 ＋ 募兵 ＋ 據點數）
//	③ 收入入帳
//	④ 累計支出歸零
//	⑤ 資金為負 → 扣預備兵
//
// ⚠ 募兵在 ② 發生，**在赤字懲罰之前**——赤字時是「先募到再被扣掉」。
// 這個順序不是設計選擇，是照抄；改了會讓長期經濟行為與原版分岔。
func Settle(f *Faction, cities []City, owner int, rng Rand) Result {
	var res Result

	// ① 先扣支出。
	f.Funds = clampFunds(f.Funds - f.Expense)

	// ② 據點結算。
	var recruit [NumTroopTypes]int
	for _, c := range cities {
		if c.Owner != owner {
			continue
		}
		res.Cities++
		div := distanceDivisor(c, f.Capital)
		base := c.Production / div
		res.GrossBase += base

		// 募兵的基數再除以 32，然後依地域拆成三兵種。
		cav, arc, inf := splitRecruits(base/32, c.Y)
		recruit[Cavalry] += cav
		recruit[Archer] += arc
		recruit[Infantry] += inf
	}
	f.Cities = res.Cities

	if f.AI {
		// 原版對非玩家勢力固定除以 2，不看稅率欄位。
		res.Income = res.GrossBase / 2
	} else {
		res.Income = res.GrossBase * f.TaxRate / 100
	}

	// 募兵數設定是**上限**：取可募量與設定值的較小者。
	for t := TroopType(0); t < NumTroopTypes; t++ {
		got := recruit[t]
		if !f.AI && got > f.RecruitCap[t] {
			got = f.RecruitCap[t]
		}
		res.Recruited[t] = got
		f.Reserves[t] = clampReserve(f.Reserves[t] + got)
	}

	// ③ 收入入帳。
	f.Funds = clampFunds(f.Funds + res.Income)

	// ④ 累計支出歸零。
	f.Expense = 0

	// ⑤ 赤字懲罰。
	if f.Funds < 0 {
		// 原版：dx = (|資金| >> 8) × 16，等價於 |資金| / 16。
		// 這裡照原版先右移再左移，因為捨去的位元不會回來——
		// 直接寫 |資金|/16 在某些值上會差 1。
		penalty := ((-f.Funds) >> 8) * 16
		for t := TroopType(0); t < NumTroopTypes; t++ {
			cut := penalty + (rng.Next() & 0x1F)
			res.Deficit[t] = cut
			f.Reserves[t] -= cut
			if f.Reserves[t] < 0 {
				f.Reserves[t] = 0
			}
		}
	}
	return res
}

// ClampFunds 把值鉗在 ±655,000。原版收入（`sub_15609`）與支出
// （`sub_15673`）用的是同一組界限，所以累計支出也要走這裡——
// 每小時的預備兵維持費就是累加進 +0x1A 的（docs/re/08 §2）。
func ClampFunds(v int) int { return clampFunds(v) }

func clampFunds(v int) int {
	if v > MaxFunds {
		return MaxFunds
	}
	if v < MinFunds {
		return MinFunds
	}
	return v
}

// ClampReserve 把預備兵夾進 0–MaxReserve。
//
// ⚠ **退兵回池那條也要用它**（`internal/state` 的 `poolBack`）。
// 原版 `sub_155EC` 的 `0xFFDC` 是在退兵路徑上驗到的（docs/spec/21 §5），
// 而 remake 先前只有月結加兵那一邊夾——**同一條規則兩份實作，其中一份漏了**。
func ClampReserve(v int) int { return clampReserve(v) }

func clampReserve(v int) int {
	if v > MaxReserve {
		return MaxReserve
	}
	if v < 0 {
		return 0
	}
	return v
}

// Exhausted 回報這個勢力是不是「疲弊狀態」。
//
// 說明書把它定義成「資金がマイナスの状態」，而機器碼是用
// `cmp ah, 80h` 檢 24 位值的符號位（docs/formats/08 §1.5）——
// 兩者是同一件事。
//
// 這個判定在兩個地方被用到：說服理由「敵が疲弊中」／「我が国疲弊」
// 的可用條件（docs/mechanics/70-ai.md §1.3），以及赤字懲罰。
func (f Faction) Exhausted() bool { return f.Funds < 0 }

// ---------------------------------------------------------------------------
// 生產力與上昇值（原版 `sub_15695`，docs/re/07 §11）
// ---------------------------------------------------------------------------

// 上昇值的範圍。原版存成一個 byte，實際值 ＝ 存值 − 100。
const (
	MaxGrowth = 100
	MinGrowth = -100

	// TaxNeutral 是稅率對上昇值不加不減的那一點。
	// 原版 `sub ax, 1Eh`——稅率低於它就讓據點繁榮，高於就讓據點荒廢。
	TaxNeutral = 30
)

// CityState 是一個據點在生產力結算時會變動的部分。
// 欄位對應據點記錄的 +0Ch／+0Eh／+10h（docs/re/07 §11）。
type CityState struct {
	Production    int // +0Eh
	ProductionCap int // +0Ch，每個據點各自不同
	Growth        int // +10h 的實際值（−100…+100），是**成長率**不是增量
	Owner         int
}

// GrowCity 跑一個據點的月度生產力結算，就地更新 c。
//
// taxRate 只在這個據點屬於玩家時才有作用（原版用 `cmp ah, [si+1]` 判斷）；
// AI 的據點完全不受稅率影響，傳 applyTax=false。
//
// 公式（docs/re/07 §11）：
//
//	r = 上昇值 − (稅率 − 30)
//	d = (生產力 >> 8) × r        生產力 >> 8 為 0 時當 1
//	r ≥ 0 → 生產力 = min(生產力 + d/2, 上限)
//	r < 0 → 生產力 = max(生產力 − |d|, 0)
//	上昇值 = clamp(r − rand(0..15), −100, +100)
//
// **變化量與生產力本身成正比**——這是說明書「大きい数値の方が変化が
// 大きくなります」的來源，也代表這是複利模型：大據點長得快、崩得也快。
func GrowCity(c *CityState, taxRate int, applyTax bool, rng Rand) {
	r := c.Growth
	if applyTax {
		r -= taxRate - TaxNeutral
	}

	// 原版取的是生產力的**高位元組**，不是除以 256 之後四捨五入。
	scale := c.Production >> 8
	if scale == 0 {
		scale = 1
	}
	d := scale * r

	switch {
	case r >= 0:
		c.Production += d / 2
		if c.Production > c.ProductionCap {
			c.Production = c.ProductionCap
		}
	default:
		c.Production += d // d 已經是負的
		if c.Production < 0 {
			c.Production = 0
		}
	}

	// 上昇值每月自然衰減 rand(0..15)。期望值 7.5，
	// 所以稅率的實際平衡點大約在 22.5%（30 − 7.5），
	// 這就是攻略章「通常は税率を下げるだけで、内政の必要はありません」的機制。
	r -= rng.Next() & 0x0F
	if r > MaxGrowth {
		r = MaxGrowth
	}
	if r < MinGrowth {
		r = MinGrowth
	}
	c.Growth = r
}

// RiotRisk 回報這個據點有沒有暴動的可能。
//
// ⚠ 說明書 9.3 寫「上昇値がマイナスになると発生率が高くなり、
// プラスの場合は発生しません」，但機器碼的門檻比這寬鬆：
// 比較用的亂數是 0–63，而上昇值存值 ＝ 實際值 ＋ 100，
// 所以**實際上昇值 ≥ −36 就完全免疫**（docs/re/07 §17）。
//
// 說明書的說法是正確但保守的建議。實際擲骰用 RollCityDisaster。
func (c CityState) RiotRisk() bool {
	return c.Growth+100 < DisasterImmunity
}
