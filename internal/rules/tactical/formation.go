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

// 陣形線的三個位置（原版指令 2 寫的三個常數）。
// 說明書 4.3：「自軍側／中央／敵軍側三選一」。
const (
	LineOwn    = 58
	LineCenter = 36
	LineEnemy  = 16
)

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

// Line 把「自軍側／中央／敵軍側」換成戰場上的 X。
//
// 三個常數是原版指令 2 寫進 `word_1D33E` 的（58／36／16），
// 而那個變數是**鏡射那一側**用的。另一側的原點要對稱過去，
// 兩軍才會面對面——用 LineFor 取，不要直接用這個。
func Line(choice int) int {
	switch choice {
	case 0:
		return LineOwn
	case 1:
		return LineCenter
	default:
		return LineEnemy
	}
}

// LineFor 回傳第 side 側在某個陣形線選擇下的原點 X。
//
// 鏡射那一側（1）直接用原版的常數；另一側對稱過去。
// **兩側都選「自軍側」時，各自貼著自己那一邊的邊緣。**
func LineFor(side, choice int) int {
	l := Line(choice)
	if side == 1 {
		return l
	}
	return Width - 1 - l
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
