package tactical

import (
	"fmt"
	"os"
)

// 十六種陣形。表在 `KI.EXE` 的 seg000 偏移 `0xCCE4`（檔案偏移 `0xCEE4`），
// 一個陣形 48 組 (dx, dy)，每組兩個有號 byte（docs/re/11 §5.8d）。
//
// **表不隨本專案散布**——與原版資產同一個處理方式，載入時請自備 KI.EXE。
const (
	NumFormations = 16
	formTableOff  = 0xCEE4
	formStride    = SoldiersOnFoot * 2 // 96，正是原版指令 1 乘的那個數
)

// 陣形原點。說明書 4.3：陣形線「自軍側／中央／敵軍側」三選一。
//
// ⭐ **兩側的常數不一樣，而且不是互相對稱的。** 原版把原點分成兩個變數：
//
//	word_1D33C  side 0（玩家）　初值 0x2005 ＝ (X=5,  Y=32)
//	word_1D33E  side 1（腳本）　初值 0x2039 ＝ (X=58, Y=32)
//
// 玩家那三個值在選單的處理常式裡（`seg000:C168`／`C184`／`C1A0`）寫死成
// 48／28／5；腳本那三個在指令 2（`sub_1A4AB`）裡寫死成 58／36／16。
// 鏡射是靠**把陣形表的 dx 取負**做的（`sub_1AA2C` 的 `neg dl`），
// 所以原點不必對稱——照抄兩組就好，不要自己算 `Width-1-l`。
var lineX = [2][3]int{
	{5, 28, 48},  // side 0：自軍側／中央／敵軍側
	{58, 36, 16}, // side 1：同上
}

// OriginY 是陣形原點的 Y，兩側都一樣（兩個初值的高位元組都是 0x20）。
const OriginY = 32

// Formations 是十六種陣形的相對座標表。
type Formations struct {
	off [NumFormations][SoldiersOnFoot][2]int8
}

// LoadFormations 從原版的 KI.EXE 讀出陣形表。
//
// ⚠ 偏移只對松崗 DOS/V 版。PC-98 是另一次編譯，同一個偏移讀到的是別的東西
// （CLAUDE.md §8 第 9 條：位址不能跨版本外推）。
func LoadFormations(exePath string) (*Formations, error) {
	b, err := os.ReadFile(exePath)
	if err != nil {
		return nil, err
	}
	end := formTableOff + NumFormations*formStride
	if len(b) < end {
		return nil, fmt.Errorf("tactical: %s 只有 %d B，讀不到 0x%X 的陣形表",
			exePath, len(b), formTableOff)
	}
	f := &Formations{}
	for i := 0; i < NumFormations; i++ {
		for k := 0; k < SoldiersOnFoot; k++ {
			o := formTableOff + i*formStride + k*2
			f.off[i][k] = [2]int8{int8(b[o]), int8(b[o+1])}
		}
	}
	return f, nil
}

// Offset 回傳第 form 個陣形裡第 k 個兵的相對座標。
func (f *Formations) Offset(form, k int) (dx, dy int) {
	if form < 0 || form >= NumFormations || k < 0 || k >= SoldiersOnFoot {
		return 0, 0
	}
	o := f.off[form][k]
	return int(o[0]), int(o[1])
}

// Bounds 回傳一個陣形的外框，給畫面與測試用。
func (f *Formations) Bounds(form int) (minX, maxX, minY, maxY int) {
	minX, minY = 127, 127
	maxX, maxY = -128, -128
	for k := 0; k < SoldiersOnFoot; k++ {
		x, y := f.Offset(form, k)
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return
}

// LineFor 回傳第 side 側在某個陣形線選擇下的原點 X。
// choice 0／1／2 ＝ 自軍側／中央／敵軍側。
func LineFor(side, choice int) int {
	if side < 0 || side > 1 {
		side = 0
	}
	if choice < 0 || choice > 2 {
		choice = 0
	}
	return lineX[side][choice]
}

// LineChoice 是 LineFor 的反查：把原點 X 換回三選一的編號，找不到回 −1。
func LineChoice(side, x int) int {
	if side < 0 || side > 1 {
		return -1
	}
	for i, v := range lineX[side] {
		if v == x {
			return i
		}
	}
	return -1
}

// SyntheticFormations 造一組不需要原版檔的陣形，給測試與無資產環境用。
//
// **這不是原版資料**：它只保證「六隊各成一列、隊長在最後面」這個
// 已經驗證過的性質（docs/re/11 §5.8d 的 ASCII 圖），數值是自己排的。
func SyntheticFormations() *Formations {
	f := &Formations{}
	for form := 0; form < NumFormations; form++ {
		width := 1 + form%4 // 讓不同陣形有不同的疏密
		for k := 0; k < SoldiersOnFoot; k++ {
			squad, idx := k/PerSquad, k%PerSquad
			if idx == 0 {
				// 隊長全部站在最後面那一行。
				f.off[form][k] = [2]int8{-2, int8(squad*3 - 7)}
				continue
			}
			f.off[form][k] = [2]int8{
				int8(2 + squad*width),
				int8((idx-4)*2 + squad),
			}
		}
	}
	return f
}
