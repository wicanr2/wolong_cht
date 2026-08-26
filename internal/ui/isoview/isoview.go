// Package isoview 是戰場的等角投影繪製。
//
// **桌面版與手機版共用這一份。** 投影公式、顯示格、深度關係全部照原版
// （`sub_1DAAA`／`sub_1DC9D`／`sub_1DDB4`），而那幾支的細節有好幾個
// 「差一列」的坑；各寫一份必然會有一邊踩到（`CLAUDE.md` §7 第 6 條）。
//
// 這一層**不認識版面**：它畫出一張原生解析度的畫布，
// 要放在畫面哪裡、放大幾倍由呼叫端決定。
package isoview

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// 戰場的**等角投影**。
//
// 公式出自 `sub_1DAAA`（docs/re/11 §5.12），那是原版把戰場座標換成畫面
// 位置的地方：
//
//	欄 ＝ X ＋ Y
//	列 ＝ (Y − X) ÷ 2 ＋ 32 − Z
//
// 兩個上限也在同一支常式裡：`cmp dx, 1Fh` 與 `cmp bx, 18h`，
// 所以**看得到的是 31 欄 × 24 列**。
//
// 一欄 8 px、一列 8 px，子圖塊 16 × 32——正好是 2 欄寬、4 列高，
// 也就是一個 16 × 8 的菱形頂加上底下的立面。**Z 每加一層就往上 8 px**，
// 所以疊起來的圖塊會自然堆成柱子。
const (
	isoCols    = 31
	isoRows    = 24
	isoColPx   = 16
	isoRowPx   = 16
	isoOriginY = 32 // `add bx, 20h`

	// originalRowBias 是原版的 Y 與我們的地形列號之間的差。
	//
	// 原版的繪圖端用**整塊 4,096 B** 的線性索引定址（`di = Y × 64 + X`，
	// 邊界是 `cmp di, 1000h` 不是 `0xF80`），而每張戰場的前 64 B 是表頭
	// ——所以原版的 `Y = 0` 是表頭那一列，第一列地形是 `Y = 1`
	// （docs/formats/07 §2.1）。`Tiles()` 只回 62 列地形、列號從 0 起算，
	// 兩個框差 1（docs/spec/57 §2）。
	//
	// **這一層只在投影的入口換一次**：規則層、小地圖、存檔全部用地形列號，
	// isoProject／cellOffset 收到之後才加上去。
	//
	// ⚠ 不能改成「鏡頭那一側減 1」——地形那條式子（走訪）減得掉，
	// 物件那條（各自投影再相減）減不掉：`floorDiv2` 不是線性的，
	// 兩邊各差一列會讓一半的物件比自己腳下的地形低一列。
	originalRowBias = 1

	// 戰場鏡頭的初值（`sub_199F3`：word_1D328=0x24、word_1D32A=0x0E）。
	// **鏡頭存的是原版的框**（含表頭那一列），換算集中在
	// isoProject／cellOffset 兩支的入口。
	battleCamInitX = 0x24
	battleCamInitY = 0x0e

	// 游標十字與鏡頭的兩個偏移（`0001C106` 的 −4、`0001C103` 的 +0x13）。
	cursorBiasX = -4
	cursorBiasY = 0x13

	// displayBandRows 是 `sub_1E0E1` 一次搬幾列：`mov cx, 8`。
	// 16×32 的編碼單位切成四帶，四帶依序落在輸出的第 0／8／16／24 列。
	displayBandRows = battle.SubTileH / 4
)

// isoProject 把**地形列號**（我們的框）換成「欄、列」。
// 入口先加 originalRowBias 換成原版的框——`floorDiv2` 不是線性的，
// 「兩邊各差一列」不會自己抵消，所以框一定要在進投影之前對齊。
func isoProject(x, y, z int) (col, row int) {
	return isoProjectOriginal(x, y+originalRowBias, z)
}

// isoProjectOriginal 吃**原版框**的 Y（含表頭那一列），給鏡頭自己用。
func isoProjectOriginal(x, y, z int) (col, row int) {
	return x + y, floorDiv2(y-x) + isoOriginY - z
}

// cellOffset 是「這一格在顯示格裡的第幾欄第幾列」，**相對於鏡頭算**。
//
// ⭐ 不能寫成 `isoProject(x,y,z) − isoProject(cam)`：
// `floorDiv2(a) − floorDiv2(b) ≠ floorDiv2(a − b)`，兩者在 **b 是奇數**時
// 對一半的格子差一列。原版根本不做這個減法——`sub_1DC9D` 從鏡頭那一格
// 開始走，交替 `di += 40h`（y+1）與 `di += 1`（x+1），
// 顯示格的位置是**走出來的**，所以鏡頭的奇偶不影響任何一格。
//
// 推導（`sub_1DC9D` 的兩層迴圈）：第 r 列從 (camX−r, camY+r) 起走，
// 第 s 格的半列位移 h ＝ (y−x) − (camY−camX) ＝ 2r ＋ (s mod 2)，
// 而 s 與 h 恆同奇偶，所以 r ＝ floorDiv2(h)。
func (v *View) cellOffset(x, y, z int) (dcol, drow int) {
	oy := y + originalRowBias // 換成原版的框（鏡頭本來就存那個框）
	return (x + oy) - (v.camWorldX + v.camWorldY),
		floorDiv2((oy-x)-(v.camWorldY-v.camWorldX)) - z
}

// floorDiv2 是原版的 `sar bx, 1`——算術右移，負數往下取整。
// 用 Go 的 `/2` 會往零取整，**在 Y < X 的那半邊會差一列**。
func floorDiv2(v int) int {
	if v < 0 {
		return -((-v + 1) / 2)
	}
	return v / 2
}

// View 是戰場畫面的資源：解好的子圖塊圖與相機。
type View struct {
	lib   *battle.Library
	set   int
	field int        // 這一張戰場的編號，重新展開堆疊時要用
	subs  [][][]byte // subs[y][x] 是那一格由下往上的子圖塊
	// tileRev 是 subs 是照哪一個版本的圖塊展開的。城壁或門被打壞時
	// 規則層會改圖塊值並把版本 +1，這裡就要重新展開（docs/spec/66）。
	tileRev int
	cache map[int]*ebiten.Image
	pal   [16]color.RGBA
	// buf 是原生解析度的離屏畫布，畫完再整張放大。
	buf *ebiten.Image

	// sprites 是 `BATTLE.SCH` 的人物圖形，載不到就是 nil（兵畫成色點）。
	sprites *battle.Sprites
	spCache map[int]*ebiten.Image
	// sourceCache 是原版合併圖形表中 BATTLE.SCH 尾端物件的快取。
	sourceCache map[int]*ebiten.Image

	// banners 是場上插的旗（`sub_19E10`，docs/re/11 §5.14）。
	banners []battle.Banner

	// minimap 是戰術初始化時由 BATTLE.MAP raw tile 與 BATTLE.MDL
	// attribute table 產生的一次性 128×128 base image。
	//
	// ⚠ **原版會在城壁破壞後局部更新**（`sub_1B824` → `sub_1BB6D`，
	// docs/spec/66；戰場本體 remake 已經跟著換了）。**這張縮圖沒有**——
	// 它建一次就不再更新，是明知的 remake 差異，缺口記在
	// docs/playtest/19。不把兵混進這個快取則是刻意的：兵每幀都在動。
	minimap *ebiten.Image

	// camWorldX／camWorldY 保存原版 word_1D328／word_1D32A 的世界格原點。
	// sub_1DC9D 會在建立畫面串列後，才把它們換算成投影的 col／row
	// origin（word_1E160／word_1E162）；兩套座標不可混為一談。
	camWorldX, camWorldY int
	// cursorX, cursorY 是小地圖上那個十字的位置（原版 word_1D32C／
	// word_1D32E）。**與鏡頭是兩組變數**，只是被同一個點選一起改。
	cursorX, cursorY int
	camCol, camRow       int
}

// image 把一個子圖塊轉成 Ebiten 的圖，解過就快取起來。
func (v *View) image(n int) *ebiten.Image {
	if img, ok := v.cache[int(n)]; ok {
		return img
	}
	s := v.lib.SubTile(v.set, n)
	if s == nil {
		v.cache[n] = nil
		return nil
	}
	rgba := image.NewRGBA(image.Rect(0, 0, battle.SubTileW, battle.SubTileH))
	for y := 0; y < battle.SubTileH; y++ {
		for x := 0; x < battle.SubTileW; x++ {
			c := s.At(x, y)
			if c == battle.Transparent {
				continue
			}
			rgba.SetRGBA(x, y, v.pal[c&15])
		}
	}
	img := ebiten.NewImageFromImage(rgba)
	v.cache[n] = img
	return img
}

// soldierImage 取一個兵的圖。圖號的算法照 `sub_1B240` 的尾段
// （docs/re/11 §5.13）：兵種 ＋（面向 × 2 ｜ 狀態旗標）。
//
// flags 目前餵兩個位元：**走路的動畫幀**（bit 0，原版每次更新完
// `xor [si+2], 1`）與**受擊**（bit 4，原版 `sub_1B618` 設，
// 效果是面向一律當成正面）。bit 3 對應的狀態還沒解。
func (v *View) soldierImage(s *tactical.Soldier, side int) *ebiten.Image {
	flags := int(s.PoseStep) & battle.PoseFlagStep
	if s.HitGeneral {
		flags |= battle.PoseFlagHitGeneral // sub_1B6BC：+0x02 bit 3
	}
	if s.Hurt {
		flags |= battle.PoseFlagFront
	}
	return v.frame(side, battle.SpriteFor(int(s.Kind), s.Facing, flags))
}

func (v *View) frame(side, n int) *ebiten.Image {
	if v.sprites == nil {
		return nil
	}
	key := side*battle.SpritesPerSide + n
	if img, ok := v.spCache[key]; ok {
		return img
	}
	f := v.sprites.Sprite(side, n)
	if f == nil {
		v.spCache[key] = nil
		return nil
	}
	rgba := image.NewRGBA(image.Rect(0, 0, battle.SpriteW, battle.SpriteH))
	for y := 0; y < battle.SpriteH; y++ {
		for x := 0; x < battle.SpriteW; x++ {
			if c := f.At(x, y); c != battle.Transparent {
				rgba.SetRGBA(x, y, v.pal[c&15])
			}
		}
	}
	img := ebiten.NewImageFromImage(rgba)
	v.spCache[key] = img
	return img
}

// sourceImage 取原版 render 合併表裡的 16×32 單位。
// BATTLE.SCH 的人物是兩個單位疊成 16×64；投射物則直接用單一單位。
func (v *View) sourceImage(raw int) *ebiten.Image {
	if v.sprites == nil {
		return nil
	}
	if v.sourceCache == nil {
		v.sourceCache = map[int]*ebiten.Image{}
	}
	if img, ok := v.sourceCache[raw]; ok {
		return img
	}
	s := v.sprites.SourceTile(raw)
	if s == nil {
		v.sourceCache[raw] = nil
		return nil
	}
	rgba := image.NewRGBA(image.Rect(0, 0, battle.SubTileW, battle.SubTileH))
	for y := 0; y < battle.SubTileH; y++ {
		for x := 0; x < battle.SubTileW; x++ {
			if c := s.At(x, y); c != battle.Transparent {
				rgba.SetRGBA(x, y, v.pal[c&15])
			}
		}
	}
	img := ebiten.NewImageFromImage(rgba)
	v.sourceCache[raw] = img
	return img
}

// applyCameraOrigin 重現 sub_1DC9D 尾端對 word_1E160／word_1E162 的換算：
// col = worldX+worldY；row = (worldY-worldX)>>1 + 32。
//
// 這裡刻意不把原點夾到 remake 自訂範圍。原版 producer 自己用 31×24
// 邊界裁切，負 row 也是有效相機狀態。
func (v *View) applyCameraOrigin() {
	v.camCol, v.camRow = isoProjectOriginal(v.camWorldX, v.camWorldY, 0)
}

// SetCameraFromMiniMap 重現 DOS/V `0001C0C6` 起的縮圖點選公式。
// 輸入是 640×400 邏輯畫布座標；原版直接使用 0x1F0（496）與
// 0xCF（207）作縮圖基準，再寫回 word_1D32A／word_1D328 並設 dirty flag。
func (v *View) SetCameraFromMiniMap(screenX, screenY int) {
	v.camWorldY = (((screenX - 0x1f0) >> 1) | 1) - 0x13
	v.camWorldX = (((0xcf - screenY) >> 1) &^ 1) + 4
	// 十字跟著點選走：原版 `0x1C0C6` 先用舊的 word_1D32C／word_1D32E
	// 還原十字，再寫新值（docs/re/60 §7）。
	//
	// **兩軸各有一個偏移**，不是直接抄鏡頭：`0001C106  add dx, 0FFFCh`
	// ＝ X − 4、`0001C103  add bx, 13h` ＝ 原版的 Y ＋ 19。
	// 用初值對一次就知道：鏡頭 (0x24, 0x0E) → 十字 (0x20, 0x21) ✓。
	v.cursorX = v.camWorldX + cursorBiasX
	v.cursorY = v.camWorldY + cursorBiasY
	v.applyCameraOrigin()
}

// SetCameraWorld 直接指定鏡頭的世界格。**驗收用**：對拍時要讓 remake
// 跟原版停在同一個鏡頭。吃的是**原版的框**（含表頭那一列），
// 所以 `sub_199F3` 的初值就是 `-battle-cam 36,14`。
func (v *View) SetCameraWorld(x, y int) {
	v.camWorldX, v.camWorldY = x, y
	v.applyCameraOrigin()
}

// 畫面上的位置：原版 `sub_1DAAA` 仍會測試 31 欄 × 24 列，但松崗 DOS/V
// 實際可見 viewport 是 480 × 368。最右一欄與最下一列可進入原版投影測試，
// 最後仍由硬體 viewport 裁掉；不能先建立 496 × 384 畫布，再靠 remake 的
// sidebar／命令列覆蓋，否則鏡頭與物件裁切的責任會落在錯誤圖層。
const (
	// nativeFieldW／H 是松崗 DOS/V 實際可見的 viewport 的一半
	// （版面常數在 `cmd/wlgame/battlelayout.go`，這裡不依賴那一層）。
	nativeFieldW = 240
	nativeFieldH = 184

	// NativeW／NativeH 是這張畫布的大小。
	NativeW = nativeFieldW * 2 // 480
	NativeH = nativeFieldH * 2 // 368
)

// drawTerrain 畫地形。
//
// 疊法照 `sub_1BB6D`：一格的圖塊是一疊 1–7 個子圖塊，
// **第 z 個畫在 Z ＝ z 的位置**，所以往上長。
//
// 畫的順序用畫家演算法：**依畫面列由上往下**。子圖塊是往下長 32 px 的，
// 所以列數大（畫面上比較低）的後畫，正好蓋住後面的。
func (v *View) drawTerrain(dst *ebiten.Image, ox, oy int) {
	// 依畫面列分桶就排好了，不必真的排序（列數是有界的）。
	const above = 4 // 子圖塊高 4 列，上面那幾列的圖還會露出來
	buckets := make([][]int32, isoRows+above)
	for y := 0; y < len(v.subs); y++ {
		for x := 0; x < len(v.subs[y]); x++ {
			for z, n := range v.subs[y][x] {
				dcol, drow := v.cellOffset(x, y, z)
				if dcol < 0 || dcol >= isoCols {
					continue
				}
				col := dcol
				r := drow + above
				if r < 0 || r >= len(buckets) {
					continue
				}
				// 一筆 ＝ 欄（低 16 位）＋ 子圖塊編號（高 16 位）。
				buckets[r] = append(buckets[r], int32(col)|int32(n)<<16)
			}
		}
	}
	op := &ebiten.DrawImageOptions{}
	for r, b := range buckets {
		for _, e := range b {
			img := v.image(int(e >> 16))
			if img == nil {
				continue
			}
			op.GeoM.Reset()
			op.GeoM.Translate(
				float64(ox+int(e&0xFFFF)*isoColPx-battle.SubTileW/2),
				float64(oy+(r-above)*isoRowPx),
			)
			dst.DrawImage(img, op)
		}
	}
}

// ScreenPos 回傳一個**物件**（兵、投射物、旗）在畫面上的位置（左上角）。
//
// 走的是 `sub_1DAAA`：各自投影再減鏡頭原點，連 `jl`／`cmp 1Fh`／`cmp 18h`
// 的裁切也一樣。**不要換成 cellOffset**——那一支是地形走訪用的，
// 兩者在鏡頭是奇數時對一半的格子差一列，而原版的物件走的就是這一條
// （docs/spec/57 §2）。
func (v *View) ScreenPos(ox, oy, x, y, z int) (int, int, bool) {
	col, row := isoProject(x, y, z)
	if col < v.camCol || col >= v.camCol+isoCols ||
		row < v.camRow || row >= v.camRow+isoRows {
		return 0, 0, false
	}
	return ox + (col-v.camCol)*isoColPx, oy + (row-v.camRow)*isoRowPx, true
}

// drawBattleIso 用原版的子圖塊畫戰場。
type displayEntryKind uint8

const (
	displayTerrain displayEntryKind = iota
	displayProjectile
	displayRawUnit
)

// 原版 word_1E15C 的固定配置。sub_1D971 以 0x3c00 words 清除
// 0x7800 bytes；sub_1DC9D 每列推進 0x400，所以是 30×32 cells，
// 每 cell 0x20 bytes。sub_1DDB4 只消費其中 23 列、每列 15 個奇數欄錨點。
const (
	battleDisplayGridCols = 32
	battleDisplayGridRows = 30
	battleDisplayCellSize = 0x20
	battleDisplayRowSize  = 0x400
	battleDisplayScanRows = 23
	battleDisplayAnchors  = 15
)

// battleDisplaySlotOffset 重現 sub_1DA1C／sub_1DB34／sub_1DC03 的槽位算式：
// cell + 4 + z*4 + lane*2。lane 0 是 [bx]，lane 1 是 [bx+2]。
func battleDisplaySlotOffset(col, row, z, lane int) int {
	return row*battleDisplayRowSize + col*battleDisplayCellSize + 4 + z*4 + lane*2
}

// battleDisplayEntry 是原版 word_1E15C 的 typed remake 中介表示。
//
// 已證實：所有 producer 共用 sub_1DA1C／1DAAA／1DB34／1DB9B／1DC03
// 的投影格與 sub_1DDB4 consumer，而不是各 layer 直接畫到 VRAM。人物／旗的
// 16×64 圖也不是一次畫完：sub_1DA1C 把奇數 raw unit 放在目前 cell 的
// lane 1，再把偶數 raw unit 放到上一列、下一個 depth 的 lane 1。
//
// 尚未證實：Ebiten 端仍把每個 16×32 unit 畫一次，沒有逐 8-pixel strip
// 重現 sub_1DE95 對六個鄰格的遮罩合成；cell／slot 欄位保留原始定位，避免
// 再退回 side 或整張 16×64 sprite 排序。
type battleDisplayEntry struct {
	kind         displayEntryKind
	col, row     int // 實際繪圖基準
	cellCol      int // word_1E15C 的 producer cell
	cellRow      int
	layer, lane  int
	x, y, z      int
	side, index  int
	raw, order   int
	pxOff, pyOff int
}

// syncTiles 把打壞的城壁與門換到畫面上。
//
// 原版 `sub_1B824` 直接改戰場緩衝區、再由 `sub_1BB6D` 重新展開那一格的堆疊，
// 而繪圖端每一幀都從同一個緩衝區重建顯示格——所以牆一垮，下一幀畫的就是瓦礫。
// remake 的 `v.subs` 是進戰場那一幀展開的靜態副本，規則層的 `Field.Retile`
// 看不到它，於是「走得過去但牆還立著」（docs/spec/66 §2）。
//
// ⭐ 旗**不重掃**：`sub_19E10` 只在開場跑一次，重掃會重擲揮舞相位，
// 旗會在牆垮的那一幀跳一下。
func (v *View) syncTiles(b *tactical.Battle) {
	if v == nil || b == nil || v.lib == nil || v.subs == nil {
		return
	}
	rev := b.Field.Revision()
	if rev == v.tileRev {
		return
	}
	v.tileRev = rev
	if !b.Field.HasTiles() {
		return
	}
	cells := make([][]byte, len(v.subs))
	for y := range cells {
		cells[y] = make([]byte, len(v.subs[y]))
		for x := range cells[y] {
			cells[y][x] = b.Field.Tile(x, y)
		}
	}
	if next := v.lib.SubTilesFor(v.field, cells); next != nil {
		v.subs = next
	}
}

func (v *View) buildDisplayList(b *tactical.Battle) []battleDisplayEntry {
	entries := make([]battleDisplayEntry, 0, 4096+2*tactical.SoldiersOnFoot)
	order := 0
	for y := 0; y < len(v.subs); y++ {
		for x := 0; x < len(v.subs[y]); x++ {
			for z, raw := range v.subs[y][x] {
				col, row := isoProject(x, y, z)
				dcol, drow := v.cellOffset(x, y, z)
				entries = append(entries, battleDisplayEntry{kind: displayTerrain,
					col: col, row: row, cellCol: dcol, cellRow: drow,
					layer: z, lane: 0, x: x, y: y, z: z, raw: int(raw), order: order})
				order++
			}
		}
	}
	frame := 0
	if b != nil {
		frame = b.Frame
	}
	for _, banner := range v.banners {
		col, row := isoProject(banner.X, banner.Y, banner.Z)
		// sub_1BB10 → sub_1DC03：一般高度用 base+phase*2；Z=6
		// 的特殊分支用 base+8+phase，且兩者都只登錄一個 raw unit。
		raw := banner.SourceTile(frame)
		entries = append(entries, battleDisplayEntry{kind: displayRawUnit,
			col: col, row: row, cellCol: col - v.camCol, cellRow: row - v.camRow,
			layer: banner.Z, lane: 0, x: banner.X, y: banner.Y, z: banner.Z,
			raw: raw, order: order})
		order++
	}
	if b != nil {
		for i, p := range b.Projectiles() {
			col, row := isoProject(p.X, p.Y, p.Z)
			entries = append(entries, battleDisplayEntry{kind: displayProjectile,
				col: col, row: row, cellCol: col - v.camCol, cellRow: row - v.camRow,
				layer: p.Z, lane: 0, x: p.X, y: p.Y, z: p.Z,
				side: p.Side, index: i, raw: ProjectileSourceIndex(p), order: order})
			order++
		}
		// 倒地動畫：原版每幀由 `sub_1B360` 換一組圖重畫（docs/spec/68）。
		// 圖號已經是「這一側這一幀該用哪一張」，直接走人物那條路。
		for _, d := range b.Deaths() {
			col, row := isoProject(d.X, d.Y, d.Z)
			raw := battle.CombinedSourceTerrainTiles + d.Sprite()*2
			entries, order = v.appendTallDisplayUnits(entries, order, col, row,
				d.X, d.Y, d.Z, raw)
		}
		for side := range b.Sides {
			for i := range b.Sides[side].Soldiers {
				s := &b.Sides[side].Soldiers[i]
				if !s.Alive {
					continue
				}
				col, row := isoProject(s.X, s.Y, s.Z)
				flags := int(s.PoseStep) & battle.PoseFlagStep
				if s.HitGeneral {
					flags |= battle.PoseFlagHitGeneral
				}
				if s.Hurt {
					flags |= battle.PoseFlagFront
				}
				n := battle.SpriteFor(int(s.Kind), s.Facing, flags)
				raw := battle.CombinedSourceTerrainTiles + side*battle.UnitsPerSide + n*2
				entries, order = v.appendTallDisplayUnits(entries, order, col, row,
					s.X, s.Y, s.Z, raw)
			}
		}
	}
	return entries
}

// appendTallDisplayUnits 重現 sub_1DA1C 的兩次寫入。raw 必須是整張
// 16×64 圖的偶數（下半）unit；raw+1 是上半。實際像素位置仍維持
// Sprite() 已驗證的「奇數在上、偶數在下」。
func (v *View) appendTallDisplayUnits(entries []battleDisplayEntry, order, col, row,
	x, y, z, raw int) ([]battleDisplayEntry, int) {
	entries = append(entries,
		battleDisplayEntry{kind: displayRawUnit, col: col, row: row,
			cellCol: col - v.camCol, cellRow: row - v.camRow, layer: z, lane: 1,
			x: x, y: y, z: z, raw: raw + 1, pyOff: -battle.SpriteH + isoRowPx,
			order: order},
		battleDisplayEntry{kind: displayRawUnit, col: col, row: row,
			cellCol: col - v.camCol, cellRow: row - v.camRow - 1, layer: z + 1, lane: 1,
			x: x, y: y, z: z, raw: raw, pyOff: -battle.SpriteH + isoRowPx + battle.SubTileH,
			order: order + 1})
	return entries, order + 2
}

type battleDisplaySlot struct {
	entry battleDisplayEntry
	set   bool
}

type battleDisplayGrid [battleDisplayGridRows][battleDisplayGridCols][battle.MaxStack][2]battleDisplaySlot

// makeDisplayGrid 重現 producer 的「槽位已佔用就不覆蓋」。目前 typed IR
// 依原版初始化順序加入地形、場景旗、投射物、兵；同一槽只保留第一筆。
func makeDisplayGrid(entries []battleDisplayEntry) battleDisplayGrid {
	var grid battleDisplayGrid
	for _, e := range entries {
		if e.cellCol < 0 || e.cellCol >= battleDisplayGridCols ||
			e.cellRow < 0 || e.cellRow >= battleDisplayGridRows ||
			e.layer < 0 || e.layer >= battle.MaxStack || e.lane < 0 || e.lane > 1 {
			continue
		}
		s := &grid[e.cellRow][e.cellCol][e.layer][e.lane]
		if !s.set {
			s.entry, s.set = e, true
		}
	}
	return grid
}

// battleDisplaySlotInfo 是顯示格表頭的兩個記帳欄位（`+1` 與 `+2`），
// **單位照原版是 2z**（docs/spec/58 §1）。
type battleDisplaySlotInfo struct {
	height int // +1：最後一個非零 unit 的 2z
	start  int // +2：最後一個**小 unit**（< 0x20）的 2z
}

// battleSmallUnit 是 `sub_1DD22` 的 `cmp al, 20h`——子圖塊編號 < 0x20 算「小」，
// 只有小的才會更新起始深度。
const battleSmallUnit = 0x20

// makeDisplayInfo 重算 `sub_1DD22` 寫在表頭的兩個欄位。
//
// 原版是邊走邊寫，而**每個（格, 深度）只有一個 producer**（docs/spec/58 §1），
// 各層又是依 z 遞增寫進來的，所以結果等於「取最大的那個 z」——
// 不必模擬走訪順序。地形在 lane 0，人物／旗在 lane 1。
//
// ⭐ **高度要含物件**（`+1`），不是只有地形那一份（`+3`）。
// `sub_1DDB4` 取鄰格高度用的是 `+1`；只算 lane 0 的話，平地上
// `z1 = 0`，兵的下半（第 1 層）就永遠畫不出來——**畫面上是半截人**
//（docs/spec/88 §3）。地形記 `2z`、物件記 `2z+1`（`58` §1 的編碼）。
//
// `start`（`+2`）維持只看 lane 0 的小 unit：那一格記的是地面，
// 物件不該把起始深度往上推。
func makeDisplayInfo(grid *battleDisplayGrid) [battleDisplayGridRows][battleDisplayGridCols]battleDisplaySlotInfo {
	var info [battleDisplayGridRows][battleDisplayGridCols]battleDisplaySlotInfo
	for row := range grid {
		for col := range grid[row] {
			for z := 0; z < battle.MaxStack; z++ {
				if s := grid[row][col][z][0]; s.set && s.entry.raw != 0 {
					info[row][col].height = 2 * z
					if s.entry.raw < battleSmallUnit {
						info[row][col].start = 2 * z
					}
				}
				if s := grid[row][col][z][1]; s.set && s.entry.raw != 0 {
					info[row][col].height = 2*z + 1
				}
			}
		}
	}
	return info
}

// displayDepthRange 回傳這個錨點要畫的深度範圍（含頭含尾），照 `sub_1DDB4`：
// 自己與**四個斜鄰格**的高度取最大，減掉自己的起始深度。
//
// 低於自己地面的深度被自己的地面擋住；要畫到最高的鄰居，
// 因為高的鄰居有半截探進這一格的圖裡（docs/spec/58 §3）。
func displayDepthRange(info *[battleDisplayGridRows][battleDisplayGridCols]battleDisplaySlotInfo,
	row, anchor int) (z0, z1 int) {
	high := info[row][anchor].height
	for _, n := range [...][2]int{{row, anchor - 1}, {row, anchor + 1},
		{row + 1, anchor - 1}, {row + 1, anchor + 1}} {
		r, c := n[0], n[1]
		if r < 0 || r >= battleDisplayGridRows || c < 0 || c >= battleDisplayGridCols {
			continue
		}
		if info[r][c].height > high {
			high = info[r][c].height
		}
	}
	z0 = info[row][anchor].start / 2
	z1 = high / 2
	if z1 >= battle.MaxStack {
		z1 = battle.MaxStack - 1
	}
	return z0, z1
}

func (v *View) rawImage(raw int) *ebiten.Image {
	if raw < battle.CombinedSourceTerrainTiles {
		return v.image(raw)
	}
	return v.sourceImage(raw)
}

// drawDisplayGrid 是 sub_1DDB4／sub_1DE95 的可追溯高階移植。
//
// 每列掃 15 個奇數欄 anchor（1,3,…,29），每個 anchor 合成一張
// 16×32 tile，畫在 ((anchor-1)*8,row*8)。sub_1E0E1 的 DX
// 0x30／0x20／0x10／0 分別選來源 unit 的第 24／16／8／0 列，
// 各取 8 px 放到輸出第 0／8／16／24 列；目前 cell 則由 sub_1E085
// 畫完整 32 px。這正是高物件跨相鄰菱形遮擋的來源。
func (v *View) drawDisplayGrid(dst *ebiten.Image, entries []battleDisplayEntry) {
	grid := makeDisplayGrid(entries)
	info := makeDisplayInfo(&grid)
	// sub_1E085／sub_1E0E1 先在 16×32 的暫存格式合成；sub_1DFE8／
	// sub_1E011 最後才把四段 16×8 重排成 VGA 上的 32×16。
	encoded := ebiten.NewImage(battle.SubTileW, battle.SubTileH)
	tile := ebiten.NewImage(battle.SubTileW*2, battle.SubTileH/2)
	op := &ebiten.DrawImageOptions{}
	for row := 0; row < battleDisplayScanRows; row++ {
		for ai := 0; ai < battleDisplayAnchors; ai++ {
			anchor := 1 + ai*2
			encoded.Clear()
			// 深度不是 0–6 全畫：從自己的地面畫到最高的鄰居
			// （`sub_1DDB4` ＋ `sub_1DE95`，docs/spec/58 §3）。
			// 多畫低於地面的那幾層，城壁的面上會多出一排亮邊。
			z0, z1 := displayDepthRange(&info, row, anchor)
			for layer := z0; layer <= z1; layer++ {
				v.drawDisplaySlice(encoded, grid[row][anchor-1][layer], 24, 0)
				v.drawDisplaySlice(encoded, grid[row][anchor+1][layer], 16, 8)
				v.drawDisplayFull(encoded, grid[row][anchor][layer])
				if row+1 < battleDisplayGridRows {
					v.drawDisplaySlice(encoded, grid[row+1][anchor-1][layer], 8, 16)
					v.drawDisplaySlice(encoded, grid[row+1][anchor+1][layer], 0, 24)
				}
			}
			unfoldDisplayTile(tile, encoded)
			op.GeoM.Reset()
			op.GeoM.Translate(float64(ai*battle.SubTileW*2), float64(row*isoRowPx))
			dst.DrawImage(tile, op)
		}
	}
}

// unfoldDisplayTile 重現 sub_1E011 的四段搬移：encoded rows 0–7、8–15、
// 16–23、24–31 依序成為輸出左上、右上、左下、右下。
func unfoldDisplayTile(dst, encoded *ebiten.Image) {
	dst.Clear()
	for i := 0; i < 4; i++ {
		srcY := i * (battle.SubTileH / 4)
		part := encoded.SubImage(image.Rect(0, srcY, battle.SubTileW,
			srcY+battle.SubTileH/4)).(*ebiten.Image)
		dstX, dstY := displayUnfoldDestination(i)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(dstX), float64(dstY))
		dst.DrawImage(part, op)
	}
}

func displayUnfoldDestination(part int) (x, y int) {
	return (part & 1) * battle.SubTileW, (part >> 1) * (battle.SubTileH / 4)
}

// displaySlotEmpty 回報這一格是不是空的。
//
// ⭐ **圖號 0 ＝ 空**：原版每幀先把整個顯示格緩衝區 `rep stosw` 清成 0
// （`sub_1D971`），所以「沒有東西」與「圖號 0」是同一個狀態——
// 圖塊堆疊中間的 0 只是那一層沒東西，不是「畫第 0 張圖」。
//
// 照著把 0 畫出來會在城門那一帶多出一圈白色的菱形邊：第 0 張子圖塊
// 是有外框的，而堆疊裡的 0 又剛好出現在門洞那幾層
// （docs/playtest/40 §12）。`makeDisplayInfo` 早就跳過 `raw == 0`，
// 畫的那一邊漏了同一條。
func displaySlotEmpty(s battleDisplaySlot) bool {
	return !s.set || s.entry.raw == 0
}

func (v *View) drawDisplayFull(dst *ebiten.Image, slots [2]battleDisplaySlot) {
	for lane := 0; lane < 2; lane++ {
		if displaySlotEmpty(slots[lane]) {
			continue
		}
		if img := v.rawImage(slots[lane].entry.raw); img != nil {
			dst.DrawImage(img, nil)
		}
	}
}

func (v *View) drawDisplaySlice(dst *ebiten.Image, slots [2]battleDisplaySlot, srcY, dstY int) {
	for lane := 0; lane < 2; lane++ {
		if displaySlotEmpty(slots[lane]) {
			continue
		}
		img := v.rawImage(slots[lane].entry.raw)
		if img == nil {
			continue
		}
		// ⭐ 一帶是 **8 列**不是 16：`sub_1E0E1` 的遮罩與四個色平面各只搬
		// `mov cx, 8` 個 word。搬 16 列會讓下一帶的內容被上一帶蓋掉一半，
		// 症狀是城壁的面上多出一排亮邊（docs/spec/58 §5）。
		slice := img.SubImage(image.Rect(0, srcY, battle.SubTileW,
			srcY+displayBandRows)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(dstY))
		dst.DrawImage(slice, op)
	}
}


// ProjectileSourceIndex 是原版 `sub_1AD2D`／`sub_1AD7F` 的 raw 圖號。
// 這些值是合併圖形表索引，不是 BATTLE.SCH 的直接單位編號。
func ProjectileSourceIndex(p tactical.ProjectileView) int {
	if p.Special {
		return 0x214 + (p.SpecialFrame & 1)
	}
	return 0x210 + (p.Direction & 1)
}


// Options 是開一張戰場要的東西。
type Options struct {
	Lib     *battle.Library
	Palette [16]color.RGBA
	Sprites *battle.Sprites

	// Field 是戰場編號。超出範圍時只建相機，不建地形——
	// 呼叫端仍拿得到一張（空的）畫布，不必為此分兩條路。
	Field int

	// Rotate 是翻轉的戰場（docs/spec/56 §3）。**地形、小地圖與旗要用
	// 同一個旗標**，否則旗會插在翻轉前的位置。
	Rotate bool

	// Rand 是插旗用的亂數（`sub_19E10` 只在開場跑一次）。
	Rand func() int

	// CamAt 不是 nil 就覆寫鏡頭初值，給驗收路徑定位用。
	CamAt *[2]int
}

// New 準備一張戰場的繪圖資源。缺素材就回 nil，
// 呼叫端會退回沒有美術的畫法。
func New(opt Options) *View {
	if opt.Lib == nil {
		return nil
	}
	var minimap *ebiten.Image
	var subs [][][]byte
	var banners []battle.Banner
	if opt.Field >= 0 && opt.Field < battle.NumFields {
		// 與規則層的戰場建構用同一個旗標：翻轉的戰場連小地圖一起翻
		// （docs/spec/56 §3）。
		tiles := opt.Lib.Tiles(opt.Field)
		if opt.Rotate {
			tiles = battle.Rotate180(tiles)
		}
		// 縮圖畫的是原版那個 64×64 緩衝區：表頭那一列也在裡面
		// （docs/formats/07 §2.1）。地形那 62 列換成可能已翻轉的版本。
		rows := opt.Lib.MinimapRows(opt.Field)
		copy(rows[1:1+len(tiles)], tiles)
		raw := battle.RenderTacticalMinimap(rows, opt.Lib.TileAttributes(opt.Field))
		minimap = ebiten.NewImageFromImage(raw.RGBA(opt.Palette))
		subs = opt.Lib.SubTilesFor(opt.Field, tiles)
		// 旗與地形要用**同一份**格子：翻轉的戰場連旗一起翻
		// （docs/playtest/40 §11）。
		rnd := opt.Rand
		if rnd == nil {
			rnd = func() int { return 0 }
		}
		banners = opt.Lib.BannersFor(opt.Field, tiles, rnd)
	}
	v := &View{
		lib: opt.Lib, set: opt.Lib.TileSet(opt.Field),
		field: opt.Field,
		subs:  subs,
		cache: map[int]*ebiten.Image{}, pal: opt.Palette,
		sprites: opt.Sprites, spCache: map[int]*ebiten.Image{},
		sourceCache: map[int]*ebiten.Image{},
		banners:     banners,
		minimap:     minimap,
		// sub_199F3：word_1D328=0x24、word_1D32A=0x0E，接著由
		// sub_1DC9D 換成投影 origin。原版只有 dirty flag 設定時才更新，
		// 不是每幀追著大將。
		camWorldX: battleCamInitX,
		camWorldY: battleCamInitY,
		// 游標十字的位置是**另一組變數**（`sub_199F3` 的 word_1D32C／
		// word_1D32E ＝ 0x20／0x21），不是鏡頭；縮圖點選時兩者一起更新。
		cursorX: battleCamInitX + cursorBiasX,
		cursorY: battleCamInitY + cursorBiasY,
	}
	if opt.CamAt != nil {
		v.camWorldX, v.camWorldY = opt.CamAt[0], opt.CamAt[1]
	}
	v.applyCameraOrigin()
	return v
}

// Render 畫出這一幀的戰場，回傳原生解析度（480×368）的畫布。
//
// ⚠ 回傳的是**內部重用的畫布**，不要保留它的參考跨幀使用。
func (v *View) Render(b *tactical.Battle) *ebiten.Image {
	v.applyCameraOrigin()
	if v.buf == nil {
		v.buf = ebiten.NewImage(NativeW, NativeH)
	}
	v.buf.Fill(color.RGBA{16, 18, 20, 255})
	v.syncTiles(b)
	v.drawDisplayGrid(v.buf, v.buildDisplayList(b))
	v.drawFallbackDots(b)
	return v.buf
}

// drawFallbackDots 是人物圖形載不到時的明確 fallback。
//
// ⚠ 正常路徑由同一份 display list 畫完，**不再依 side／陣列順序覆蓋**
// 原版的深度關係——那樣畫出來的遮擋是錯的。
func (v *View) drawFallbackDots(b *tactical.Battle) {
	if v.sprites != nil {
		return
	}
	for i := range b.Sides {
		base := color.RGBA{235, 90, 70, 255}
		if i == 1 {
			base = color.RGBA{90, 150, 245, 255}
		}
		for k := range b.Sides[i].Soldiers {
			s := &b.Sides[i].Soldiers[k]
			if !s.Alive {
				continue
			}
			px, py, ok := v.ScreenPos(0, 0, s.X, s.Y, s.Z)
			if !ok {
				continue
			}
			c := base
			size := 4
			if s.Cmd == tactical.Retreat {
				c = color.RGBA{130, 130, 130, 255}
			}
			if s.IsGeneral() {
				size = 6
				c = color.RGBA{250, 220, 130, 255}
				if i == 1 {
					c = color.RGBA{210, 230, 255, 255}
				}
			}
			vector.DrawFilledRect(v.buf, float32(px-size/2), float32(py-size),
				float32(size), float32(size), c, false)
		}
	}
}

// Minimap 是戰術初始化時產生的 128×128 底圖，沒有素材時是 nil。
func (v *View) Minimap() *ebiten.Image { return v.minimap }

// Cursor 是小地圖上那個十字的位置。**與鏡頭是兩組變數**，
// 只是被同一個點選一起改。
func (v *View) Cursor() (x, y int) { return v.cursorX, v.cursorY }

// Camera 是鏡頭的世界格原點（原版的框，含表頭那一列）。
func (v *View) Camera() (x, y int) { return v.camWorldX, v.camWorldY }

// Field 是這一張戰場的編號。
func (v *View) Field() int { return v.field }
