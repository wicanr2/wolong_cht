// Package battle 解 `BATTLE.MAP`／`BATTLE.MDL`／`BATTLE.DAT`。
//
// 規格：docs/formats/07-battle.md
// 出處：docs/re/11（`sub_1CAEB` 載入、`sub_1BB6D` 展開堆疊、`sub_1CBE5` 挑腳本）
//
// 這一層是純解碼，回傳格子與腳本，不認識 Ebiten 也不認識規則層。
package battle

import "fmt"

// BATTLE.MAP 的佈局（confirmed，docs/formats/07 §2）。
const (
	IndexSize  = 512  // 檔頭索引，每筆 2 byte
	FieldSize  = 4096 // 一個戰場
	NumFields  = 214
	CellsOff   = 0x40 // 表頭 64 B 之後才是格子
	Width      = 64
	Height     = 62
	NumCells   = Width * Height // 3,968
	FieldsBase = IndexSize
)

// BATTLE.MDL 的佈局（confirmed，docs/formats/07 §5）。
const (
	MDLHeader   = 4096
	TileSetSize = 63488
	NumTileSets = 3
	// 每個圖塊組前 2,048 B 是 256 筆 × 8 B 的圖塊定義。
	TileDefs   = 256
	TileDefLen = 8
	// MaxStack 是一格最多疊幾層。原版把七層的通行性寫在相隔 0x1000 的
	// 七張圖上（`sub_1BB6D` 的 `mov ah, 7`）。
	MaxStack = 7
)

// BATTLE.DAT 是 32 段 × 256 B 的腳本（docs/formats/07 §6）。
const (
	ScriptSize = 256
	NumScripts = 32
)

// Library 是解好的戰場資料。
type Library struct {
	mapData []byte
	// stacks[組][圖塊] 是那個圖塊由下往上的子圖塊；長度就是堆疊高度。
	stacks [NumTileSets][TileDefs][]byte
	script []byte
}

// Parse 解三個檔。scripts 可以是 nil（沒有 `BATTLE.DAT` 就不驅動 AI）。
func Parse(mapData, mdl, scripts []byte) (*Library, error) {
	if n := IndexSize + NumFields*FieldSize; len(mapData) != n {
		return nil, fmt.Errorf("battle: BATTLE.MAP 是 %d B，預期 %d", len(mapData), n)
	}
	if n := MDLHeader + NumTileSets*TileSetSize; len(mdl) != n {
		return nil, fmt.Errorf("battle: BATTLE.MDL 是 %d B，預期 %d", len(mdl), n)
	}
	if scripts != nil && len(scripts) != ScriptSize*NumScripts {
		return nil, fmt.Errorf("battle: BATTLE.DAT 是 %d B，預期 %d",
			len(scripts), ScriptSize*NumScripts)
	}
	l := &Library{mapData: mapData, script: scripts}
	for t := 0; t < NumTileSets; t++ {
		base := MDLHeader + t*TileSetSize
		for i := 0; i < TileDefs; i++ {
			r := mdl[base+i*TileDefLen : base+(i+1)*TileDefLen]
			// [0] 是堆疊高度，[1..k] 是由下往上的子圖塊，[k+1..7] 全是 0。
			// 三個圖塊組共 768 筆驗過，零例外（docs/re/11 §4.2）。
			k := int(r[0])
			if k > MaxStack {
				k = MaxStack
			}
			l.stacks[t][i] = append([]byte(nil), r[1:1+k]...)
		}
	}
	return l, nil
}

// TileSet 回傳第 n 張戰場用哪一組圖塊。
func (l *Library) TileSet(n int) int { return int(l.mapData[n*2]) }

// GateX 回傳索引的第二欄：**命令 3（城壁移動）要走過去那一格的 X**
// （docs/re/11 §5.8i）。**0 表示這是野戰用的戰場**——
// 214 張裡零例外（§4.5）。
func (l *Library) GateX(n int) int { return int(l.mapData[n*2+1]) }

// IsSiege 回報第 n 張是不是攻城用的。
func (l *Library) IsSiege(n int) bool { return l.GateX(n) != 0 }

// Stacks 回傳第 n 張戰場每一格的堆疊高度，逐列排列（Height × Width）。
//
// **一格地圖存的是圖塊編號，而一個圖塊是一疊 1–7 層的子圖塊**——
// 堆疊高度就是那一格的地面高度，兵站在它上面（docs/re/11 §4.2、§5.2）。
func (l *Library) Stacks(n int) [][]int {
	t := l.TileSet(n)
	off := FieldsBase + n*FieldSize + CellsOff
	cells := l.mapData[off : off+NumCells]
	out := make([][]int, Height)
	for y := 0; y < Height; y++ {
		out[y] = make([]int, Width)
		for x := 0; x < Width; x++ {
			out[y][x] = len(l.stacks[t][cells[y*Width+x]])
		}
	}
	return out
}

// Tiles 回傳第 n 張戰場每一格的**原始圖塊值**，逐列排列（Height × Width）。
//
// 城壁與門是從這些值認出來的（0xD0–0xDF 與 0xF0–0xF7，docs/re/11 §5.9），
// 而且打壞時要把圖塊值換掉，所以規則層需要的是這個而不只是堆疊高度。
// 回傳的是複本，改它不會動到載入的檔案。
func (l *Library) Tiles(n int) [][]byte {
	off := FieldsBase + n*FieldSize + CellsOff
	cells := l.mapData[off : off+NumCells]
	out := make([][]byte, Height)
	for y := 0; y < Height; y++ {
		out[y] = append([]byte(nil), cells[y*Width:(y+1)*Width]...)
	}
	return out
}

// Heights 回傳第 n 張戰場那一組圖塊的堆疊高度表（圖塊值 → 層數）。
func (l *Library) Heights(n int) *[256]int {
	t := l.TileSet(n)
	var h [256]int
	for i := 0; i < TileDefs; i++ {
		h[i] = len(l.stacks[t][i])
	}
	return &h
}

// Script 回傳一段 AI 腳本。
//
// 段編號 ＝ **武將記錄 `+0x16` × 4 ＋ 戰場類別**（`sub_1CBE5`，docs/re/11 §3.2）。
// 沒載 `BATTLE.DAT` 時回 nil。
func (l *Library) Script(kind, category int) []byte {
	if l.script == nil {
		return nil
	}
	n := kind*4 + category
	if n < 0 || n >= NumScripts {
		return nil
	}
	return l.script[n*ScriptSize : (n+1)*ScriptSize]
}

// Category 回傳戰場編號對應的類別（`sub_19A33` 的三分）。
//
//	< 0xC0  →  0   攻城（戰場編號就是據點編號）
//	< 0xD1  →  1   野戰
//	否則    →  2   野戰（另一組）
func Category(field int) int {
	switch {
	case field < 0xC0:
		return 0
	case field < 0xD1:
		return 1
	}
	return 2
}
