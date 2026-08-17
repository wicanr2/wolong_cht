package main

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
)

// isoProject 把戰場座標換成「欄、列」。
func isoProject(x, y, z int) (col, row int) {
	return x + y, floorDiv2(y-x) + isoOriginY - z
}

// floorDiv2 是原版的 `sar bx, 1`——算術右移，負數往下取整。
// 用 Go 的 `/2` 會往零取整，**在 Y < X 的那半邊會差一列**。
func floorDiv2(v int) int {
	if v < 0 {
		return -((-v + 1) / 2)
	}
	return v / 2
}

// battleView 是戰場畫面的資源：解好的子圖塊圖與相機。
type battleView struct {
	lib   *battle.Library
	set   int
	subs  [][][]byte // subs[y][x] 是那一格由下往上的子圖塊
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
	// attribute table 產生的一次性 128×128 base image。原版是否在
	// 城壁破壞後做局部更新尚未知，因此不把兵或結構變更混入這個快取。
	minimap *ebiten.Image

	// camWorldX／camWorldY 保存原版 word_1D328／word_1D32A 的世界格原點。
	// sub_1DC9D 會在建立畫面串列後，才把它們換算成投影的 col／row
	// origin（word_1E160／word_1E162）；兩套座標不可混為一談。
	camWorldX, camWorldY int
	camCol, camRow       int
}

// newBattleView 準備一張戰場的繪圖資源。lib 為 nil 就回 nil，
// 呼叫端會退回沒有美術的畫法。
func (g *game) newBattleView(field int) *battleView {
	if g.battleLib == nil || g.lib == nil || g.lib.Palette == nil {
		return nil
	}
	bank, err := g.lib.Palette.Bank(0)
	if err != nil {
		return nil
	}
	var minimap *ebiten.Image
	var subs [][][]byte
	if field >= 0 && field < battle.NumFields {
		// 與 buildField 用同一個旗標：翻轉的戰場連小地圖一起翻
		// （docs/spec/56 §3）。
		tiles := g.battleLib.Tiles(field)
		if g.battleRotate {
			tiles = battle.Rotate180(tiles)
		}
		raw := battle.RenderTacticalMinimap(tiles, g.battleLib.TileAttributes(field))
		minimap = ebiten.NewImageFromImage(raw.RGBA(bank))
		subs = g.battleLib.SubTilesFor(field, tiles)
	}
	v := &battleView{
		lib: g.battleLib, set: g.battleLib.TileSet(field),
		subs:  subs,
		cache: map[int]*ebiten.Image{}, pal: bank,
		sprites: g.battleSprites, spCache: map[int]*ebiten.Image{},
		sourceCache: map[int]*ebiten.Image{},
		banners:     g.battleLib.Banners(field, g.rng.Next),
		minimap:     minimap,
		// sub_199F3：word_1D328=0x24、word_1D32A=0x0E，接著由
		// sub_1DC9D 換成投影 origin。原版只有 dirty flag 設定時才更新，
		// 不是每幀追著大將。
		camWorldX: 0x24,
		camWorldY: 0x0e,
	}
	if g.battleCamAt != nil {
		v.camWorldX, v.camWorldY = g.battleCamAt[0], g.battleCamAt[1]
	}
	v.applyCameraOrigin()
	return v
}

// image 把一個子圖塊轉成 Ebiten 的圖，解過就快取起來。
func (v *battleView) image(n int) *ebiten.Image {
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
func (v *battleView) soldierImage(s *tactical.Soldier, side int) *ebiten.Image {
	flags := int(s.PoseStep) & battle.PoseFlagStep
	if s.HitGeneral {
		flags |= battle.PoseFlagHitGeneral // sub_1B6BC：+0x02 bit 3
	}
	if s.Hurt {
		flags |= battle.PoseFlagFront
	}
	return v.frame(side, battle.SpriteFor(int(s.Kind), s.Facing, flags))
}

func (v *battleView) frame(side, n int) *ebiten.Image {
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
func (v *battleView) sourceImage(raw int) *ebiten.Image {
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
func (v *battleView) applyCameraOrigin() {
	v.camCol, v.camRow = isoProject(v.camWorldX, v.camWorldY, 0)
}

// setCameraFromMiniMap 重現 DOS/V `0001C0C6` 起的縮圖點選公式。
// 輸入是 640×400 邏輯畫布座標；原版直接使用 0x1F0（496）與
// 0xCF（207）作縮圖基準，再寫回 word_1D32A／word_1D328 並設 dirty flag。
func (v *battleView) setCameraFromMiniMap(screenX, screenY int) {
	v.camWorldY = (((screenX - 0x1f0) >> 1) | 1) - 0x13
	v.camWorldX = (((0xcf - screenY) >> 1) &^ 1) + 4
	v.applyCameraOrigin()
}

// setCameraWorld 直接指定鏡頭的世界格。**驗收用**：對拍時要讓 remake
// 跟原版停在同一個鏡頭（原版初值是 `sub_199F3` 的 36, 14）。
func (v *battleView) setCameraWorld(x, y int) {
	v.camWorldX, v.camWorldY = x, y
	v.applyCameraOrigin()
}

// 畫面上的位置：原版 `sub_1DAAA` 仍會測試 31 欄 × 24 列，但松崗 DOS/V
// 實際可見 viewport 是 480 × 368。最右一欄與最下一列可進入原版投影測試，
// 最後仍由硬體 viewport 裁掉；不能先建立 496 × 384 畫布，再靠 remake 的
// sidebar／命令列覆蓋，否則鏡頭與物件裁切的責任會落在錯誤圖層。
const (
	isoNativeW = dosvBattleNativeFieldW * 2 // 480
	isoNativeH = dosvBattleNativeFieldH * 2 // 368
	isoScale   = 1
)

// drawTerrain 畫地形。
//
// 疊法照 `sub_1BB6D`：一格的圖塊是一疊 1–7 個子圖塊，
// **第 z 個畫在 Z ＝ z 的位置**，所以往上長。
//
// 畫的順序用畫家演算法：**依畫面列由上往下**。子圖塊是往下長 32 px 的，
// 所以列數大（畫面上比較低）的後畫，正好蓋住後面的。
func (v *battleView) drawTerrain(dst *ebiten.Image, ox, oy int) {
	// 依畫面列分桶就排好了，不必真的排序（列數是有界的）。
	const above = 4 // 子圖塊高 4 列，上面那幾列的圖還會露出來
	buckets := make([][]int32, isoRows+above)
	for y := 0; y < len(v.subs); y++ {
		for x := 0; x < len(v.subs[y]); x++ {
			for z, n := range v.subs[y][x] {
				col, row := isoProject(x, y, z)
				if col < v.camCol || col >= v.camCol+isoCols {
					continue
				}
				r := row - v.camRow + above
				if r < 0 || r >= len(buckets) {
					continue
				}
				// 一筆 ＝ 欄（低 16 位）＋ 子圖塊編號（高 16 位）。
				buckets[r] = append(buckets[r], int32(col-v.camCol)|int32(n)<<16)
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

// ScreenPos 回傳一個戰場座標在畫面上的位置（左上角）。
func (v *battleView) ScreenPos(ox, oy, x, y, z int) (int, int, bool) {
	col, row := isoProject(x, y, z)
	if col < v.camCol || col >= v.camCol+isoCols ||
		row < v.camRow || row >= v.camRow+isoRows {
		return 0, 0, false
	}
	return ox + (col-v.camCol)*isoColPx, oy + (row-v.camRow)*isoRowPx, true
}

// drawBattleIso 用原版的子圖塊畫戰場。
func (g *game) drawBattleIso(screen *ebiten.Image, b *tactical.Battle, me *tactical.Soldier) {
	v := g.view
	l := dosvBattleLayoutFor(screenW, screenH)
	v.applyCameraOrigin()
	if v.buf == nil {
		v.buf = ebiten.NewImage(isoNativeW, isoNativeH)
	}
	v.buf.Fill(color.RGBA{16, 18, 20, 255})
	entries := v.buildDisplayList(b)
	v.drawDisplayGrid(v.buf, entries)

	// 人物圖形載不到時的明確 fallback；正常 raw 資產路徑已由同一份
	// display list 畫完，不再依 side／陣列順序覆蓋原版深度關係。
	if v.sprites == nil {
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

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(isoScale, isoScale)
	op.GeoM.Translate(float64(l.Field.X), float64(l.Field.Y))
	screen.DrawImage(v.buf, op)
}

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

func (v *battleView) buildDisplayList(b *tactical.Battle) []battleDisplayEntry {
	entries := make([]battleDisplayEntry, 0, 4096+2*tactical.SoldiersOnFoot)
	order := 0
	for y := 0; y < len(v.subs); y++ {
		for x := 0; x < len(v.subs[y]); x++ {
			for z, raw := range v.subs[y][x] {
				col, row := isoProject(x, y, z)
				entries = append(entries, battleDisplayEntry{kind: displayTerrain,
					col: col, row: row, cellCol: col - v.camCol, cellRow: row - v.camRow,
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
		phase := (banner.Phase + frame) % battle.BannerFrames
		raw := 0x150
		if banner.Side == 1 {
			raw = 0x204
		}
		if banner.Z == 6 {
			raw += 8 + phase
		} else {
			raw += phase * 2
		}
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
				side: p.Side, index: i, raw: projectileSourceIndex(p), order: order})
			order++
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
func (v *battleView) appendTallDisplayUnits(entries []battleDisplayEntry, order, col, row,
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

func (v *battleView) rawImage(raw int) *ebiten.Image {
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
func (v *battleView) drawDisplayGrid(dst *ebiten.Image, entries []battleDisplayEntry) {
	grid := makeDisplayGrid(entries)
	// sub_1E085／sub_1E0E1 先在 16×32 的暫存格式合成；sub_1DFE8／
	// sub_1E011 最後才把四段 16×8 重排成 VGA 上的 32×16。
	encoded := ebiten.NewImage(battle.SubTileW, battle.SubTileH)
	tile := ebiten.NewImage(battle.SubTileW*2, battle.SubTileH/2)
	op := &ebiten.DrawImageOptions{}
	for row := 0; row < battleDisplayScanRows; row++ {
		for ai := 0; ai < battleDisplayAnchors; ai++ {
			anchor := 1 + ai*2
			encoded.Clear()
			for layer := 0; layer < battle.MaxStack; layer++ {
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

func (v *battleView) drawDisplayFull(dst *ebiten.Image, slots [2]battleDisplaySlot) {
	for lane := 0; lane < 2; lane++ {
		if !slots[lane].set {
			continue
		}
		if img := v.rawImage(slots[lane].entry.raw); img != nil {
			dst.DrawImage(img, nil)
		}
	}
}

func (v *battleView) drawDisplaySlice(dst *ebiten.Image, slots [2]battleDisplaySlot, srcY, dstY int) {
	for lane := 0; lane < 2; lane++ {
		if !slots[lane].set {
			continue
		}
		img := v.rawImage(slots[lane].entry.raw)
		if img == nil {
			continue
		}
		slice := img.SubImage(image.Rect(0, srcY, battle.SubTileW, srcY+isoRowPx)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(dstY))
		dst.DrawImage(slice, op)
	}
}

func (v *battleView) drawDisplayEntry(dst *ebiten.Image, b *tactical.Battle, e battleDisplayEntry) {
	px, py, ok := v.ScreenPos(0, 0, e.x, e.y, e.z)
	if !ok {
		return
	}
	op := &ebiten.DrawImageOptions{}
	switch e.kind {
	case displayTerrain:
		img := v.image(e.raw)
		if img == nil {
			return
		}
		op.GeoM.Translate(float64(px-battle.SubTileW/2), float64(py))
		dst.DrawImage(img, op)
	case displayProjectile:
		img := v.sourceImage(e.raw)
		if img == nil {
			return
		}
		op.GeoM.Translate(float64(px-battle.SubTileW/2), float64(py))
		dst.DrawImage(img, op)
	case displayRawUnit:
		img := v.sourceImage(e.raw)
		if img == nil {
			return
		}
		op.GeoM.Translate(float64(px-battle.SubTileW/2+e.pxOff), float64(py+e.pyOff))
		dst.DrawImage(img, op)
	}
}

// drawProjectiles 將規則層的 raw 格座標畫成原版 BATTLE.SCH 投射物。
//
// `sub_1AD2D` 以 DX=0x210／0x211 選普通箭，方向只取奇偶；
// `sub_1AD7F` 以 DX=0x214／0x215 選特殊效果。特殊效果的低位由發射兵
// 記錄 +0x02 bit 0 帶出，規則層的 ProjectileView 以 SpecialFrame 保留，
// 因此兩張原版 raw frame 都能被取到。
func (v *battleView) drawProjectiles(dst *ebiten.Image, b *tactical.Battle) {
	for _, p := range b.Projectiles() {
		px, py, ok := v.ScreenPos(0, 0, p.X, p.Y, p.Z)
		if !ok {
			continue
		}
		raw := projectileSourceIndex(p)
		if img := v.sourceImage(raw); img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(px-battle.SubTileW/2), float64(py))
			dst.DrawImage(img, op)
			continue
		}
		c := color.RGBA{235, 90, 70, 255}
		if p.Side == 1 {
			c = color.RGBA{90, 150, 245, 255}
		}
		if p.Special {
			c = color.RGBA{245, 220, 110, 255}
			vector.DrawFilledRect(dst, float32(px-1), float32(py-7), 2, 6, c, false)
			continue
		}
		dx, dy := 0, 0
		switch p.Direction & 3 {
		case tactical.West:
			dx = -3
		case tactical.North:
			dy = -3
		case tactical.East:
			dx = 3
		case tactical.South:
			dy = 3
		}
		vector.DrawFilledRect(dst, float32(px-1+dx/2), float32(py-3+dy/2), 2, 2, c, false)
	}
}

// projectileSourceIndex 是原版 `sub_1AD2D`／`sub_1AD7F` 的 raw 圖號。
// 這些值是合併圖形表索引，不是 BATTLE.SCH 的直接單位編號。
func projectileSourceIndex(p tactical.ProjectileView) int {
	if p.Special {
		return 0x214 + (p.SpecialFrame & 1)
	}
	return 0x210 + (p.Direction & 1)
}

// drawBanners 畫場上插的旗（`sub_19E10`，docs/re/11 §5.14）。
//
// 位置開場就決定好，旗色由圖塊編號的最低位選、與交戰雙方無關；
// **但旗子會飄**——`sub_1BB10` 每幀把相位加一再 `and 3`，四張一循環。
// 建立時的亂數初值讓滿場的旗不會同步揮舞。
func (v *battleView) drawBanners(dst *ebiten.Image, frame int) {
	op := &ebiten.DrawImageOptions{}
	for _, b := range v.banners {
		px, py, ok := v.ScreenPos(0, 0, b.X, b.Y, b.Z)
		if !ok {
			continue
		}
		img := v.frame(b.Side, b.Sprite(frame))
		if img == nil {
			continue
		}
		op.GeoM.Reset()
		op.GeoM.Translate(float64(px-battle.SpriteW/2),
			float64(py-battle.SpriteH+isoRowPx))
		dst.DrawImage(img, op)
	}
}
