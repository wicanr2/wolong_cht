package economy

// 三種災害的觸發規則。反組譯見 docs/re/07 §17、§18，
// 機制說明見 docs/mechanics/40-economy.md §5。
//
// 說明書第 9 章把三種災害各給了一句描述，這裡是那三句話的實際公式。
// 兩處與說明書不同，都在下面標出來了。

// Disaster 是三種災害之一。
type Disaster int

const (
	NoDisaster Disaster = iota
	Fire                // 火災，比防災值
	Riot                // 暴動，比上昇值
	Storm               // 暴風雨，隨機挑據點，東部機率加倍
)

func (d Disaster) String() string {
	switch d {
	case Fire:
		return "火災"
	case Riot:
		return "暴動"
	case Storm:
		return "暴風雨"
	}
	return "無"
}

const (
	// disasterGate 是火災與暴動的第一道機率閘：亂數 < 24（滿值 256）。
	// 所以單一據點單月的最高災害機率是 24/256 ＝ 9.4%。
	disasterGate = 0x18

	// DisasterImmunity 是完全免疫的門檻。
	//
	// ⚠ 原版拿來比較的亂數是 `and al, 3Fh` ＝ **0–63**，
	// 所以防災值或上昇值存值只要 **≥ 64** 就完全不會觸發。
	// 說明書把安全區講得比實際小（「上昇値がプラスの場合は発生しません」），
	// 實際上昇值只要 ≥ −36 就免疫了。
	DisasterImmunity = 64
)

// RollCityDisaster 對一個據點擲一次月度災害。
//
// 原版 `sub_12286` 的順序是「先擲火災，沒中才擲暴動」，
// 所以**同一個據點同一個月不會同時發生兩種**，而且火災優先。
//
// disasterPrevention 是據點記錄 +11h（防災值），
// growthStored 是 +10h（上昇值的**存值**，＝實際值 ＋ 100）。
// 這裡刻意收存值而不是實際值，因為原版就是直接拿存值去比。
func RollCityDisaster(disasterPrevention, growthStored int, rng Rand) Disaster {
	if rng.Next()&0xFF < disasterGate {
		if rng.Next()&0x3F >= disasterPrevention {
			return Fire
		}
	}
	if rng.Next()&0xFF < disasterGate {
		if rng.Next()&0x3F >= growthStored {
			return Riot
		}
	}
	return NoDisaster
}

// StormArea 是暴風雨籠罩的範圍，以格為單位。
// 原版是以中心據點為準的 11 × 11 格（`[X−5, X+5] × [Y−5, Y+5]`）。
type StormArea struct {
	MinX, MinY, MaxX, MaxY int
}

// coastalThreshold 是「靠海」判定的門檻。
//
// ⚠ **原版比的是 X 座標的低位元組**（`cmp byte ptr [bx+848h], 0C0h`），
// 而 X 的值域是 4–370。結果是 `X mod 256 ≥ 192` 才算靠海：
//
//	程式判定為靠海            52 個據點
//	真正的東半部（X ≥ 192）  101 個據點
//	其中 X ≥ 256 被判成內陸    49 個據點  ← 最東邊、最靠海的那些
//
// 說明書寫「発生率は海側がわずかに高くなっています」，可見原意是東部加成，
// 是 `byte ptr` 讓它只對 X ∈ [192, 255] 那條帶狀區域生效。
//
// **這是原版的 bug，remake 刻意照抄**（CLAUDE.md §8：原版行為優先）。
// 要改成「真正的東半部」只需把下面那行換成 `x >= 192`。
const coastalThreshold = 0xC0

// IsCoastalForStorm 重現原版的「靠海」判定，含上面說的低位元組 bug。
func IsCoastalForStorm(x int) bool { return x&0xFF >= coastalThreshold }

// RollStorm 擲一次月度暴風雨。回傳 nil 表示這個月沒有暴風雨。
//
// cities 必須是完整的 192 筆據點表（原版用編號直接索引，不是只掃自己的）。
// 暴風雨**每月最多發生一次，而且不分勢力**——它挑的是全地圖的據點。
func RollStorm(cities []City, rng Rand) *StormArea {
	if rng.Next()&1 == 0 {
		return nil
	}
	idx := rng.Next() & 0xFF
	if idx >= 192 || idx >= len(cities) {
		return nil
	}
	c := cities[idx]
	if !IsCoastalForStorm(c.X) && rng.Next()&1 == 0 {
		return nil
	}
	return &StormArea{
		MinX: shiftBack(c.X), MinY: shiftBack(c.Y),
		MaxX: shiftBack(c.X) + 10, MaxY: shiftBack(c.Y) + 10,
	}
}

// shiftBack 是原版算範圍左上角的方式：座標 ≥ 10 時退 5，否則不動。
// 不是 `max(0, v-5)`——小於 10 的座標**完全不退**，所以範圍會偏右下。
func shiftBack(v int) int {
	if v < 10 {
		return v
	}
	return v - 5
}

// Contains 回報一個座標是否落在暴風雨範圍內。
func (s StormArea) Contains(x, y int) bool {
	return x >= s.MinX && x <= s.MaxX && y >= s.MinY && y <= s.MaxY
}
