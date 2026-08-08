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
	isoCols   = 31
	isoRows   = 24
	isoColPx  = 8
	isoRowPx  = 8
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
	sprites  *battle.Sprites
	spCache  map[int]*ebiten.Image

	camCol, camRow int
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
	return &battleView{
		lib: g.battleLib, set: g.battleLib.TileSet(field),
		subs:  g.battleLib.SubTiles(field),
		cache: map[int]*ebiten.Image{}, pal: bank,
		sprites: g.battleSprites, spCache: map[int]*ebiten.Image{},
	}
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

// 相機的可動範圍。把 64 × 62 的格子投影過去之後，
// 欄落在 0–124、列落在 −5–62，超出去就只會看到空白。
const (
	maxCol = tactical.Width - 1 + tactical.Height - 1
	minRow = -(tactical.Width - 1) / 2
	maxRow = (tactical.Height-1)/2 + isoOriginY
)

// soldierImage 取一個兵的圖。
//
// ⭐ **圖號就是兵種的儲存值**（0／18／36／54）——兵種存成「× 18」正是
// 為了當索引用（docs/formats/07 §10）。姿勢那 18 張裡誰是誰還沒解，
// 所以這裡先用面向去挑，**那一段是 remake 的選擇**。
func (v *battleView) soldierImage(side int, kind tactical.Kind, facing int) *ebiten.Image {
	return v.frame(side, battle.SpriteFor(int(kind), facing))
}

// bannerImage 取軍旗。大將身邊插的那一支。
func (v *battleView) bannerImage(side, pose int) *ebiten.Image {
	return v.frame(side, battle.BannerSprite+pose%battle.PosesPerKind)
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

// centreOn 把相機對到某一格上，並夾在戰場的範圍內。
func (v *battleView) centreOn(x, y, z int) {
	col, row := isoProject(x, y, z)
	v.camCol = clampInt(col-isoCols/2, 0, maxCol-isoCols+1)
	v.camRow = clampInt(row-isoRows/2, minRow, maxRow-isoRows+1)
}

func clampInt(v, lo, hi int) int {
	if hi < lo {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// 畫面上的位置：原生 248 × 192，照 640 × 400 的視窗放大兩倍。
// 原版是 320 × 200，而這支的視窗正好是它的兩倍。
const (
	isoNativeW = isoCols * isoColPx // 248
	isoNativeH = isoRows * isoRowPx // 192
	isoScale   = 2
	isoScreenX = 0
	isoScreenY = 24
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
	v.centreOn(me.X, me.Y, me.Z)
	if v.buf == nil {
		v.buf = ebiten.NewImage(isoNativeW, isoNativeH)
	}
	v.buf.Fill(color.RGBA{16, 18, 20, 255})
	v.drawTerrain(v.buf, 0, 0)

	// 兵。畫成色塊疊在等角座標上——**人物圖形在 `BATTLE.SCH`，還沒解**，
	// 所以這一層仍然是暫代的（README 有標）。
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
			if img := v.soldierImage(i, s.Kind, s.Facing); img != nil {
				op := &ebiten.DrawImageOptions{}
				// 腳底對準格子：圖高 64，站的那一格在最下面那一列。
				op.GeoM.Translate(float64(px-battle.SpriteW/2),
					float64(py-battle.SpriteH+isoRowPx))
				v.buf.DrawImage(img, op)
				// 大將身邊插軍旗（圖 72–89 那一組）。
				if s.IsGeneral() {
					if b := v.bannerImage(i, s.Facing); b != nil {
						op.GeoM.Reset()
						op.GeoM.Translate(float64(px-battle.SpriteW/2+6),
							float64(py-battle.SpriteH+isoRowPx-4))
						v.buf.DrawImage(b, op)
					}
				}
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

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(isoScale, isoScale)
	op.GeoM.Translate(isoScreenX, isoScreenY)
	screen.DrawImage(v.buf, op)
}
