# 107 — 啟動殼層的 UI 顏色也要查調色盤

**狀態：CONFORMED。啟動殼層沒有 `World`，`paletteInk` 因此整片退回硬寫的
近似色。改成用「殼層固定第 0 組」解析組別、並把勢力清單的底色、標題列、
數字欄與捲軸接回原版那一套之後，**清單本體（標題列 ＋ 九列資料）
逐像素 0 / 58,880 px**。**

- 日期：2026-09-01
- 出處：[`79`](79-new-game-faction-list.md) §1.1.1（`sub_11A6E` 的
  `mov al, 0 / call sub_10241` ＝ **調色盤組 0**）、
  [`54`](54-ui-colours-from-palette.md)（介面顏色一律查調色盤）、
  [`38`](38-list-windows.md) §1.1（清單本體是米色底）、
  [`../re/26`](../re/26-list-window-engine.md)
- 推論等級：**confirmed**——顏色是實機量出來的，而且與 `GAMEPAL.BRG`
  第 0 組的色號逐項對得上

## 1. 原版做什麼

新遊戲的四層畫在大地圖上，畫之前先切**調色盤組 0**。所以殼層與遊戲中
用的是同一組顏色，介面元件照同一套色號取色：

| 元件 | 色號 | 第 0 組的值（過 VGA 6-bit DAC 之後）|
|---|---:|---|
| 清單視窗的本體底 | 9 | `(243, 211, 146)` |
| 清單本體上的字 | 0 | `(0, 0, 0)` |
| 按鈕的面 | 7 | `(195, 130, 32)` |
| 按鈕的亮邊 | 9 | `(243, 211, 146)` |
| 按鈕的暗邊 | 6 | `(130, 65, 32)` |

## 2. 量到的落差

2026-09-01 拿呂布開局的君主卡逐像素比（[`../playtest/56`](../playtest/56-lubu-flow-parity.md) §4.7），
1,859 個不同像素**全部**落在兩顆鈕上，而且是**一對一的換色**、
形狀與像素數完全相同：

| 原版 | remake | 像素數 | 對應色號 |
|---|---|---:|---:|
| `(195, 130, 32)` | `(90, 90, 90)` | 1,267 | 7 |
| `(243, 211, 146)` | `(210, 210, 210)` | 131 | 9 |
| `(130, 65, 32)` | `(60, 60, 60)` | 83 | 6 |

那三個灰正是 `cmd/wlgame/displaylist.go` 的 `dlButtonFallback`／
`dlLightFallback`／`dlDarkFallback`。

勢力清單那一張是同一件事的另一半：本體底原版 `(243, 211, 146)`、
remake `(0, 0, 0)`；字原版 `(0, 0, 0)`、remake `(200, 200, 210)`。

### 2.1 底色蓋掉之後才露出來的三項

⚠ **黑底把三個真正的差異藏起來了**：原版的字是黑的，remake 的底也是黑的，
所以「原版有字、remake 沒畫」與「兩邊都畫對了」在差分圖上長得一模一樣
（[`../playtest/56`](../playtest/56-lubu-flow-parity.md) §4.1 第一版因此
把數字欄記成「逐像素相同」）。底色改對之後量到：

| 項目 | 原版 | remake（修前）|
|---|---|---|
| 標題列 | 黑底（5,301 px）＋ 白字（843 px）| 米色底 ＋ 同樣的白字 |
| 兩個數字欄的 X | 三位數欄位的**左緣**（絕對 304／352）；單一位數落在 321–327、369–375 | 當成右緣再減一個欄寬 ⇒ 297–303、345–351，**整欄左偏 24 px** |
| 數字的字模 | 原版 8×16 數字（[`38`](38-list-windows.md) §1.5）| 文字字型的 ASCII 數字，形狀不同。**君主卡的兩個數字也一樣**——那 59 px 是修完鈕的顏色之後才露出來的 |
| 沒有軍師的那一列 | 軍師名欄印「－－－」| 留白 |
| 捲軸 | 黑槽 ＋ 3D 綠鈕 ＋ 比例式滑塊 | 三個描邊的空框 |

## 3. 成因

```go
func (g *game) paletteInk(index int, fallback color.RGBA) color.RGBA {
    if g.lib == nil || g.world == nil {   // ← 殼層的 g.world 是 nil
        return fallback
    }
    c, err := g.lib.PaletteColor(int(g.world.Clock.Season()), index)
    ...
}
```

`paletteInk` 用 `g.world.Clock.Season()` 取組別，而**啟動殼層還沒有世界**
——那一頁預覽用的世界存在另一個欄位 `launcherPreviewWorld`。
於是殼層畫的每一個顏色都拿到 fallback。

⭐ **反證很乾淨**：同一支 `dlButton` 進了遊戲就對了。編成面板的「確定」鈕
在同一輪對拍是逐像素相同的（[`../playtest/56`](../playtest/56-lubu-flow-parity.md) §3
第 13 列），因為那時候 `g.world` 已經有了。

## 4. 演算法

```
UI 顏色的調色盤組 =
    世界存在 → 世界時鐘的季節
    否則     → 0        ← 啟動殼層，照 sub_11A6E 的 mov al, 0
```

勢力清單的兩個顏色不走 `paletteInk`，直接用 `chrome` 已經載好的色票
（`chrome.Load(lib, 0)` 在殼層之前就跑過，所以 `chrome.Sheet`／`chrome.Ink`
在殼層也是對的）：

```
清單本體底 = chrome.Sheet   （色 9）
清單本體字 = chrome.Ink     （色 0）
```

## 5. remake 實作

| 項目 | 位置 |
|---|---|
| 組別解析 | `cmd/wlgame/strategyhud.go` 的 `uiPaletteBank()`，`paletteInk` 改用它 |
| 君主卡 | `cmd/wlgame/lordcard.go`：兩個數字改 `drawOriginalNumber` |
| 勢力清單 | `cmd/wlgame/factionlist.go` 的 `drawFactionList`：本體填 `chrome.Sheet`、標題列黑條、列字 `chrome.Ink`、數字欄改 `drawOriginalNumber` 且畫在欄位左緣、沒有軍師印「－－－」|
| 捲軸 | `cmd/wlgame/main.go` 的 `drawScrollbarAt`：矩形由呼叫端給，戰略層一覽表與勢力清單共用同一份畫法 |
| 差異 | **反白條是 remake 的鍵盤游標**（原版那一列沒有高亮，量到 5,420 px 都是米色）。原版是純滑鼠，沒有「目前這一列」的概念 |

## 6. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestUIPaletteBankFallsBackToLauncherSeason`、`TestPaletteInkWithoutWorldUsesRealPalette`（正對照：**沒有世界時也要拿到真的顏色**）、`TestPaletteInkWithoutLibraryKeepsFallback`（`cmd/wlgame`）|
| 對原版 | **君主卡 (160, 112) 240×192 → 0 / 46,080 px PASS**（修前 1,859 px）；**勢力清單本體 (152, 104) 368×160 → 0 / 58,880 px PASS**。整個清單視窗 (128, 96) 400×192 還剩 5,927 px，那是反白條與捲軸滑塊（見 §7）|

## 7. 未解

| 項目 | 現況 |
|---|---|
| 殼層其餘幾頁（ＹＥＳ／ＮＯ、劇本、四槽讀檔）的配色 | 這一輪只對過勢力清單與君主卡。其餘幾頁同樣走 `paletteInk`，修完應該一起好，但**沒有逐像素比過** |
| 捲軸滑塊的位置差 1 px | 量到：22 筆、`top` ＝ 4 時，原版的綠面在 y 161–216，remake 在 160–215（高度都是 56）。`⌊128×4/22⌋ ＝ 23` 給 159，原版對應的是 24。**只有這一個取樣點**，分不出是無條件進位、四捨五入還是槽的起點差 1；`38` §1.6 的實機量測只釘住高度沒釘位置。**不憑一個樣本改一支已經逐像素驗過的算式** |
| 反白條 | remake 的鍵盤游標，原版沒有。要不要照戰略層一覽表那樣「碰過才畫」（`g.listTouched`）沒有定案 |
