package world

import (
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

// 據點在大地圖上的樣子跟著歸屬換。規格 docs/spec/53、出處 docs/re/67。
//
// ⭐ **`MMAP.MAP` 只存「無所屬」那一張。** 中心格的三個值 205／208／211
// 各自是一組三張圖塊的最後一張，有主就往前挪：自勢力 −2、他勢力 −1。
// 全圖只有 192 格落在這個範圍，剛好等於據點數。

// CityCentreDX 是據點中心格相對於據點記錄座標的位移。
//
// ⭐ **是 0：記錄座標就是中心格。** 192 座據點的記錄座標逐座落在
// 值 205／208／211 的中心格上，零例外（`docs/spec/53` §5）。
// 這個常數留著是因為位移曾經被算成 4——成因是 `MMAP.MAP` 解壓後
// 開頭那 4 byte 是長度欄位（`world.MapHeader`），從 offset 0 讀會整張左移四格。
const CityCentreDX = 0

// NoOwner 是「無所屬」的勢力編號（原版 `0x18`，docs/re/62 §2）。
const NoOwner = 0x18

// 據點的裝飾格（角落的旗）在地圖檔裡是 222–229，**自勢力**時整批 +10。
//
// 帶是量出來的：把原版與 remake 逐格比，差異格的（檔案值 → 畫的值）
// 全部是 +10 且逐像素 0 差（`docs/spec/53` §4）。
//
// ⚠ **220、221、231 也在據點方塊裡但不換**，所以兩端是硬的；
// 230 更麻煩——它在關隘的上下會換，在大城方塊的左右邊**不換**
// （許昌的 (208,113)／(208,115) 是 230，原版照畫），
// 所以 230 由 `gateDecor` 單獨處理。
const (
	cityDecorLow  = 222
	cityDecorHigh = 229
	cityDecorStep = 10
	// 裝飾格離中心最遠 ±2（5×5 大城的四個角）。
	cityDecorRadius = 2
	// 關隘（中心 211）上下各一格的門柱。
	gateDecor       = 230
	gateCentreTile  = 211
	gateDecorOffset = 1
)

// Ownership 是據點與玩家的關係。三種狀態與縮小地圖的分類同一套。
type Ownership int

const (
	OwnedBySelf  Ownership = iota // 玩家所仕的勢力
	OwnedByOther                  // 別的勢力
	Unowned                       // 無所屬
)

// OwnershipOf 把據點的勢力編號換成三分類。
func OwnershipOf(owner, player int) Ownership {
	switch {
	case owner == NoOwner:
		return Unowned
	case owner == player:
		return OwnedBySelf
	default:
		return OwnedByOther
	}
}

// IsCityCentre 回報這個圖塊值是不是據點中心。
//
// 三個值兩兩相差 3：205 一般據點、208 另一種、211 關隘。
func IsCityCentre(tile byte) bool {
	return tile == 205 || tile == 208 || tile == 211
}

// CityCentreTile 回傳據點中心該畫的圖塊。
//
// base 是地圖檔裡的值；不是據點中心就原樣回傳，不要猜。
func CityCentreTile(base byte, own Ownership) byte {
	if !IsCityCentre(base) {
		return base
	}
	switch own {
	case OwnedBySelf:
		return base - 2
	case OwnedByOther:
		return base - 1
	default:
		return base
	}
}

// CapitalOverlayTile 是首都中心格要疊上去的 `MMAP.MCH` 圖塊。
//
// 112 個不透明點，疊在自勢力那一張（203）上與原版逐像素相同；
// 非首都疊上去反而差 75 點，所以**只有首都疊**。
const CapitalOverlayTile = 0xFF

// CorpsMark 是一支軍團要疊在大地圖上的圖塊（docs/spec/74）。
//
// Tile 是 **MMAP.MCH 的圖塊編號**，由呼叫端算好：原版 `sub_12B2A` 取
// 軍團記錄 `+0x09`（勢力編號 × 5）加 `+0x08`（朝向），每個勢力五張圖
// ——四個方向 ＋ 靜止。
type CorpsMark struct {
	// X、Y 是軍團所在的地圖格（原版軍團記錄 +0x10／+0x12）。
	X, Y int
	Tile byte
}

// CorpsHeadings 是朝向的值域：0／1 是 X 減增、2／3 是 Y 減增、4 是靜止。
const CorpsHeadings = 5

// CorpsTile 算一支軍團該畫哪一張 MCH 圖塊（docs/spec/74 §3）。
//
// 原版 `sub_12B2A`：`al = [si+9]`（勢力編號 × 5）`+ [si+8]`（朝向）。
// **每個勢力五張圖**，22 個勢力共 110 張。
//
// ⭐ 算式只寫在這裡一處。桌面與手機各自建疊圖清單，但**規則不重複**
// （CLAUDE.md §7 第 6 條）。
func CorpsTile(faction, heading int) byte {
	if heading < 0 || heading >= CorpsHeadings {
		// 朝向越界時退回「靜止」那一張，而不是畫出別的勢力的圖。
		heading = CorpsHeadings - 1
	}
	return byte(faction*CorpsHeadings + heading)
}

// CityMark 是一座據點要在大地圖上改的東西。
type CityMark struct {
	// X、Y 是**中心格**的地圖座標（＝據點記錄座標 + CityCentreDX）。
	X, Y    int
	Own     Ownership
	Capital bool
}

// applyDecor 把據點方塊裡的裝飾格換成自勢力的版本。
//
// ⚠ **只有自勢力換。** 別的勢力與無所屬的據點只換中心那一格——
// 拿嚴白虎的牛渚與無主的舒縣量過，它們的角落與地圖檔逐像素相同。
//
// 判準是**地圖檔本身寫了什麼**，不是據點的類型欄位：5×5 的大城角是
// 222–225、3×3 的是 226–229、關隘上下是 230，一條規則全部涵蓋。
// 全圖有 838 格落在這個帶裡，其中 43 格離任何據點中心超過兩格——
// 那些不是據點的裝飾，所以**半徑要限制住**，不能整張掃。
func (m *Map) applyDecor(mark CityMark, base byte, put func(x, y int, tile byte)) {
	if mark.Own != OwnedBySelf {
		return
	}
	for dy := -cityDecorRadius; dy <= cityDecorRadius; dy++ {
		for dx := -cityDecorRadius; dx <= cityDecorRadius; dx++ {
			x, y := mark.X+dx, mark.Y+dy
			t, err := m.Tile(x, y)
			if err != nil || t < cityDecorLow || t > cityDecorHigh {
				continue
			}
			put(x, y, t+cityDecorStep)
		}
	}
	if base != gateCentreTile {
		return
	}
	for _, dy := range [2]int{-gateDecorOffset, gateDecorOffset} {
		if t, err := m.Tile(mark.X, mark.Y+dy); err == nil && t == gateDecor {
			put(mark.X, mark.Y+dy, t+cityDecorStep)
		}
	}
}

// RenderMarked 與 Render 相同，但先照 marks 換掉據點中心（與大城的四個角）
// 的圖塊，再把首都的 MCH 圖塊疊上去。
//
// mch 可以是 nil（沒載到 `MMAP.MCH`）——那樣只少了首都那一張，
// 其餘照畫，不要整張失敗。
func (m *Map) RenderMarked(ts *TileSet, mch *MCH, pal *palette.Palette, bank,
	x0, y0, cols, rows int, marks []CityMark, corps []CorpsMark) (*image.RGBA, error) {
	swap := make(map[int]byte, len(marks)*5)
	overlay := make(map[int]byte, len(marks))
	put := func(x, y int, tile byte) { swap[y*Width+x] = tile }
	for _, mark := range marks {
		base, err := m.Tile(mark.X, mark.Y)
		if err != nil || !IsCityCentre(base) {
			// (X+4, Y) 不是據點中心 → 這一座跳過。呼叫端負責回報，
			// 這裡不要靜悄悄畫錯的圖塊。
			continue
		}
		put(mark.X, mark.Y, CityCentreTile(base, mark.Own))
		m.applyDecor(mark, base, put)
		if mark.Capital {
			overlay[mark.Y*Width+mark.X] = CapitalOverlayTile
		}
	}

	corpsAt := make(map[int]byte, len(corps))
	for _, c := range corps {
		corpsAt[c.Y*Width+c.X] = c.Tile
	}

	img := image.NewRGBA(image.Rect(0, 0, cols*TileSize, rows*TileSize))
	for ry := 0; ry < rows; ry++ {
		for rx := 0; rx < cols; rx++ {
			mx, my := x0+rx, y0+ry
			t, err := m.Tile(mx, my)
			if err != nil {
				return nil, err
			}
			if v, ok := swap[my*Width+mx]; ok {
				t = v
			}
			tile, err := ts.spec.RenderRGBA(ts.data, int(t), pal, bank)
			if err != nil {
				return nil, err
			}
			for py := 0; py < TileSize; py++ {
				for px := 0; px < TileSize; px++ {
					img.SetRGBA(rx*TileSize+px, ry*TileSize+py, tile.RGBAAt(px, py))
				}
			}
			if id, ok := overlay[my*Width+mx]; ok {
				blitMCH(img, mch, pal, bank, id, rx, ry)
			}
			// ⭐ 軍團**畫在首都疊圖之後**：同一格可能兩者都有
			// （軍團在自己的首都裡），而原版的顯示表是後推的層蓋在
			// 前面的層上（`sub_1D66A` 依序消費 si+3..si+6）。
			if id, ok := corpsAt[my*Width+mx]; ok {
				blitMCH(img, mch, pal, bank, id, rx, ry)
			}
		}
	}
	return img, nil
}

// blitMCH 把一張 MCH 圖塊疊到 img 的第 (rx, ry) 格。
//
// ⚠ MCH 的 0xFF 是**遮罩判定的透明像素**，不是色號——照畫會蓋掉地形。
func blitMCH(img *image.RGBA, mch *MCH, pal *palette.Palette, bank int,
	id byte, rx, ry int) {
	if mch == nil {
		return
	}
	over := mch.Tile(id)
	colours, err := pal.Bank(bank)
	if err != nil || over == nil {
		return
	}
	for py := 0; py < TileSize; py++ {
		for px := 0; px < TileSize; px++ {
			c := over.Pix[py*TileSize+px]
			if c == MCHTransparent || int(c) >= len(colours) {
				continue
			}
			img.SetRGBA(rx*TileSize+px, ry*TileSize+py, colours[c])
		}
	}
}
