// Package governor 是內政官：**每個據點每天一次**的整備。
//
// 出處是 `sub_14194`（docs/re/07 §19），掛在主迴圈的每 tick 據點更新
// `sub_13EFD` 上——一次一個據點，192 個輪一圈，而一天是 216 tick，
// 所以「每個據點大約每天被處理一次」。
//
// ⭐ **這一支是「120 個月 1872 次暴動」的正解。**
// 月結每月扣上昇值 `rand(0..15)`（期望 −7.5），而補回來的就是這裡：
// AI 的據點每天有 9/16 的機率 +1，月期望 +16.9 —— **淨值是正的**。
// 少了這一層，所有 AI 據點都會單調往下掉到暴動。
//
// ⚠ **玩家的據點基準比 AI 差。** `cl` 的預設是 8，但只要據點屬於玩家
// 就降到 5，要靠內政官的政治值補回來。這看起來不對稱，但兩個分支的
// 常數都是直接讀出來的——**照抄，不要「修正」成對稱的**。
// 它讓「人事」這個指令有實質意義：不派內政官的玩家據點比 AI 還糟。
package governor

// 兩組基準值。出自 `sub_14194` 開頭的 `mov cl,8 / mov dl,4` 與
// 玩家分支的 `mov cl,5 / mov dl,1`。
const (
	// BaseRate 是「成功率」那個計數（與 rand(0..15) 比大小）。
	BaseRate = 8
	// BaseDraft 是徵兵時一次的人數，同時也是徵兵對上昇值的代價。
	BaseDraft = 4

	// PlayerRate／PlayerDraft 是**玩家的據點**的基準——比上面更低。
	PlayerRate  = 5
	PlayerDraft = 1

	// gainOffset 是把 rate 換成增量的那個減數（`sub ch, 0Fh`）。
	gainOffset = 15
	// MaxValue 是上昇值與防災值的上限（`cmp al, 0C8h`）。
	MaxValue = 200
	// draftChance 是徵兵的門檻：`call rng / cmp al, 18h / jnb 不徵`。
	// 亂數是 0–255，所以機率是 24/256 ≈ 9.4%。
	draftChance = 24
	// rateMask 是成功判定用的亂數遮罩（`and al, 0Fh`）。
	rateMask = 0x0F
)

// City 是這一支會讀寫的據點欄位。名稱對應據點記錄的偏移。
type City struct {
	Growth      int // +10h 的實際值
	Prevention  int // +11h
	Garrison    int // +13h
	GarrisonCap int // +12h
}

// Official 是派駐在這個據點的內政官。Budget 是武將記錄 `+1Ah` 的經費餘額。
type Official struct {
	Politics int // 武將 +13h
	Martial  int // 武將 +11h
	Budget   int // 武將 +1Ah，0 ＝ 沒錢（這時內政官不作用）
}

// Tick 跑一個據點的一次整備。
//
// isPlayer 決定用哪一組基準。gov 為 nil 表示沒派內政官——
// **注意這不等於「跳過」**：沒有內政官照樣會跑，只是用基準值。
//
// gov 不是 nil 且 Budget > 0 時會**扣掉 1 點經費**（就地改 gov.Budget），
// 呼叫端要把它寫回武將記錄。
//
// rnd 要回 0–255 的值（原版的 `sub_1ECE0`）。
func Tick(c *City, gov *Official, isPlayer bool, rnd func() int) {
	rate, draft := BaseRate, BaseDraft
	if isPlayer {
		// ⚠ 玩家的據點基準更差，見 package 註解。
		rate, draft = PlayerRate, PlayerDraft
		if gov != nil && gov.Budget > 0 {
			gov.Budget--
			rate += gov.Politics
			draft = (draft + gov.Martial) >> 1
		}
	}

	// 增量：`mov ch, cl / sub ch, 0Fh / ja / mov ch, 1`。
	// `ja` 是**無號**比較，所以 rate ≤ 15 一律得到 1。
	gain := 1
	if rate > gainOffset {
		gain = rate - gainOffset
	}

	// ① 上昇值。成功條件是 `cmp cl, al / jb 跳過` ＝ rate ≥ rand(0..15)。
	if rate >= rnd()&rateMask {
		c.Growth = min(c.Growth+gain, MaxValue)
	}
	// ② 防災值。增量是 `shr ch,1 / inc ch`，**先除再加一**。
	if rate >= rnd()&rateMask {
		c.Prevention = min(c.Prevention+(gain>>1)+1, MaxValue)
	}
	// ③ 徵兵。三個細節都是從尾段那幾行讀出來的，不要簡化：
	//
	//   - **城兵已達或超過上限時，會被修剪回上限**（`jnb loc_1422F`
	//     直接跳到 `mov ch, cl`）——即使這一輪根本沒徵兵。
	//   - 徵到的兵**要拿上昇值去換**（`sub [si+850h], dl`，下限 0）。
	//   - 先夾 255（`add ch,dl / jnb / mov ch,0FFh`，byte 溢位）
	//     **再**夾上限（`cmp ch,cl / jbe`）。兩個夾都要，順序也要對。
	if c.Garrison >= c.GarrisonCap {
		c.Garrison = c.GarrisonCap
		return
	}
	if rnd() < draftChance {
		c.Growth = max(c.Growth-draft, 0)
		c.Garrison = min(min(c.Garrison+draft, 255), c.GarrisonCap)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
