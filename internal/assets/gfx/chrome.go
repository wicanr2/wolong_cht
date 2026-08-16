package gfx

import (
	"fmt"
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

// 視窗外框的圖塊。**在 `ICONGRF` 段 3，不在段 1。**
//
// 段 3（`0x9700`，9,120 B）是混合內容；通用 8×8 外框圖塊仍只認已對上實機的部分，
// DOS/V 數值面板另已由 `sub_17D0D` 定位到 96×64 的 `0x14A0` 資源。
// 歷史上曾把整段標為「未解」，只知道它走 `sub_1F888`
// ——那支**位元對齊**的繪製常式，可以把圖畫在非 8 倍數的 x 上。
// 這正是介面元件需要的能力（視窗可以開在任意位置），
// 而地圖圖塊、頭像那些永遠對齊 8 的東西都走另一支。
// **「為什麼要用位元對齊的繪製常式」本身就指出了這一段是什麼。**
//
// 歷史找法曾拿 PC-98 實機截圖的邊框像素回掃段內圖塊；那只作資源交叉驗證，
// 本輪 DOS/V 數值內框則改以 `sub_17D0D` 的直接資源位址為準。比對時不看
// 實機調色盤，而看**「顏色的等價關係」**——把每個像素換成「它是這塊裡第幾種
// 出現的顏色」的簽章。詳見 docs/formats/03 §5.4。
//
// 尺寸是 **8×8**：畫面上量到的紅框 motif 週期正好 8 px，
// 邊柱也是 8 px 寬（黑 橙 黃 白 黃 橙 褐 黑）。
// 4 bpp × 8×8 ＝ 32 B／塊，而 9,120 ÷ 32 ＝ 285 整除。
//
// ⚠ **段 3 不是整段同尺寸的圖塊陣列**（與段 1 一樣是混合內容），
// 所以這裡只認已經對上實機畫面的那幾塊，不假裝整段都解得開。
const (
	// ChromeTile 是外框圖塊的邊長。
	ChromeTile = 8

	// ChromeEdge 是上下邊那條「紅框綠心」的 motif。
	ChromeEdge = 0x06C0
	// ChromeCap 是左右邊柱最上面那一塊（金色柱頭）。
	ChromeCap = 0x06E0
	// ChromeShaft 是邊柱的柱身，往下重複貼。
	// 段內 `0x0700`–`0x0720` 每 8 byte 都相符 ——
	// 那是因為柱身在垂直方向同紋，不是有五塊不同的圖。
	ChromeShaft = 0x0700

	// DOSVOrderIconOffset 是戰術底列那六張「目前命令」圖示的起點。
	//
	// `sub_1C673` 從 `word_10D48` 段取 `命令碼 × 0xC0`，而 `word_10D48`
	// 就是段 3 本身的載入段（`sub_1006B` 把檔案位移 0x9700 讀進去，
	// `docs/re/48` §6.1），所以段內位移就是 `碼 × 0xC0`，從 0 起算。
	//
	// 解出來的六張與命令碼逐張相符：0 陣形（陣地）／1 攻擊（長槍）／
	// 2 突擊（軍旗）／3 城壁（磚牆）／4 守陣（盾）／5 退卻（白旗），
	// 這是「換算對了」的內容檢查。後面第 7–9 張是兵種圖示的**橘色版**
	// （馬／弓／步），與 0x1BA0 的紅版、0x1EA0 的綠版同一批。
	DOSVOrderIconOffset = 0x0000
	DOSVOrderIconStride = 0xC0
	DOSVOrderIconCount  = 6

	// DOSVAmountPanelOffset 是 DOS/V sub_17D0D 的數值視窗圖形在
	// ICONGRF 第 3 段的 byte offset。
	//
	// 換算鏈（`docs/re/48` §6）：`sub_1006B` 把 `ICONGRF` 檔案位移 0x9700
	// 起（＝第 3 段）讀進 `word_10D48`，而 `sub_100DF` 讓
	// `word_10D50 = word_10D48 + 0x9A` 段 ＝ **+0x9A0 byte**。
	// `sub_17D0D` 讀 `word_10D50:0600h`，所以段內位移是 0x600 + 0x9A0。
	//
	// 解出來是完整的數字鍵盤（7 8 9 ◀ 取消／4 5 6 0 最大／1 2 3 00 決定），
	// 這是「換算對了」的內容檢查。
	DOSVAmountPanelOffset = 0x0FA0

	// ICONGRF 第 3 段尾段的四張 24×16 圖示：天秤（資金）、馬、弓、步。
	// 顯示清單 op 09 以 `word_10D50` 的 0x1200 起連號取用（`docs/re/48` §4），
	// 換算後落在段 3 的 0x1BA0；同樣四張的綠色版在 0x1EA0。
	// 紅色版用在自勢力情報與財政的「今月底」欄，綠色版用在財政的「次月」欄。
	// DOSVFactionLegendOffset 是縮小地圖視窗下方那條勢力色標
	// （左半紅、右半藍，各帶一個小色塊）。`sub_15A3A` 以
	// `word_10D50:0000h` 貼在 (440,168)，換算後是段 3 的 0x09A0。
	DOSVFactionLegendOffset = 0x09A0

	DOSVResourceIconOffset      = 0x1BA0
	DOSVResourceIconGreenOffset = 0x1EA0
	DOSVResourceIconStride      = 0xC0
	DOSVResourceIconCount       = 4

	// DOSVEmptySlotIconOffset 是編成畫面「兵種 4 ＝ 空槽」用的圖示。
	// `sub_16DA8` 以 `0x15C0 + (兵種−1)×0xC0` 取圖，兵種 1–3 落在綠色那
	// 一組的第 2–4 張，兵種 4 落在**綠組後面一張**（`docs/re/49` §3）。
	DOSVEmptySlotIconOffset = 0x21A0
)

// DOSVAmountPanel 是 sub_17D0D 以 AX=4006h 複製的 96×64 平面圖。
// 它不是可重複貼的 8×8 chrome tile，而是數值輸入器自己的完整內框。
var DOSVAmountPanel = Spec{Name: "ICONGRF/DOSV amount panel", Width: 96, Height: 64}

// DOSVResourceIcon 是資金／預備兵欄左邊那一直排圖示。
var DOSVResourceIcon = Spec{Name: "ICONGRF/DOSV resource icon", Width: 24, Height: 16}

// DOSVFactionLegend 是縮小地圖下方那條 192×16 的勢力色標。
var DOSVFactionLegend = Spec{Name: "ICONGRF/DOSV faction legend", Width: 192, Height: 16}

// chromeBytes 是一塊 8×8 4bpp 的大小。
const chromeBytes = ChromeTile * ChromeTile * Planes / 8

// DecodeChrome 解段 3 裡位於 off 的那一塊 8×8，回傳調色盤索引。
//
// seg 要傳**段 3 本身**（`IconRegions[3]` 那一段），不是整個 `ICONGRF.DAT`。
func DecodeChrome(seg []byte, off int) ([]byte, error) {
	if off < 0 || off+chromeBytes > len(seg) {
		return nil, fmt.Errorf("外框圖塊偏移 0x%X 超出段 3（%d byte）", off, len(seg))
	}
	out := make([]byte, ChromeTile*ChromeTile)
	for y := 0; y < ChromeTile; y++ {
		for x := 0; x < ChromeTile; x++ {
			var v byte
			for pl := 0; pl < Planes; pl++ {
				b := seg[off+pl*ChromeTile+y]
				v |= ((b >> (7 - uint(x))) & 1) << uint(pl)
			}
			out[y*ChromeTile+x] = v
		}
	}
	return out, nil
}

// RenderChrome 解一塊外框圖塊並上色。
func RenderChrome(seg []byte, off int, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	idx, err := DecodeChrome(seg, off)
	if err != nil {
		return nil, err
	}
	bank, err := p.Bank(bankIdx)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, ChromeTile, ChromeTile))
	for y := 0; y < ChromeTile; y++ {
		for x := 0; x < ChromeTile; x++ {
			img.SetRGBA(x, y, bank[idx[y*ChromeTile+x]])
		}
	}
	return img, nil
}

// ---------------------------------------------------------------------------
// 視窗內部的龍紋
// ---------------------------------------------------------------------------

const (
	// WindowTextureOffset 是龍紋點陣在 `ICONGRF` 段 3 的位移。
	// 檔案位移 `0xBA20`（段 3 從 `0x9700` 起）**正好是檔尾**，
	// 而段 3 由 `sub_1006B` 以 `di = 0FFFFh` ＝「讀到檔尾」載入。
	WindowTextureOffset = 0x2320
	// WindowTextureSize 是磚塊的邊長：32×32、1 bpp、每列 4 byte ＝ 128 B。
	WindowTextureSize = 32
	// windowTextureBytes 是那 128 byte。
	windowTextureBytes = WindowTextureSize * WindowTextureSize / 8

	// 只有兩色（docs/formats/03 §5.5）。
	windowTextureInk   = 8 // 深藍
	windowTexturePaper = 0 // 黑
)

// DecodeWindowTexture 解出 32×32 的調色盤索引（只會是 0 或 8）。
func DecodeWindowTexture(seg []byte, off int) ([]byte, error) {
	if off < 0 || off+windowTextureBytes > len(seg) {
		return nil, fmt.Errorf("gfx: 視窗底紋位移 0x%X 超出段長 %d", off, len(seg))
	}
	out := make([]byte, WindowTextureSize*WindowTextureSize)
	for y := 0; y < WindowTextureSize; y++ {
		row := seg[off+y*4 : off+y*4+4]
		for x := 0; x < WindowTextureSize; x++ {
			if row[x>>3]>>(7-(x&7))&1 == 1 {
				out[y*WindowTextureSize+x] = windowTextureInk
			} else {
				out[y*WindowTextureSize+x] = windowTexturePaper
			}
		}
	}
	return out, nil
}

// RenderWindowTexture 解出龍紋並上色。
func RenderWindowTexture(seg []byte, off int, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	idx, err := DecodeWindowTexture(seg, off)
	if err != nil {
		return nil, err
	}
	bank, err := p.Bank(bankIdx)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, WindowTextureSize, WindowTextureSize))
	for y := 0; y < WindowTextureSize; y++ {
		for x := 0; x < WindowTextureSize; x++ {
			img.SetRGBA(x, y, bank[idx[y*WindowTextureSize+x]])
		}
	}
	return img, nil
}
