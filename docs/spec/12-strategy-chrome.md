# 12 — 主畫面的視窗外框與指令列

**狀態：DRAFT。原版的外框是格子屬性層堆出來的，不是圖塊 blit
（[`docs/re/46`](../re/46-strategy-chrome-cell-layer.md)）。
底層的像素產生方式還沒讀出來，**還不能照著實作**。**

- 日期：2026-08-14
- 出處：[`docs/re/46`](../re/46-strategy-chrome-cell-layer.md)（`sub_1614A`／`sub_1895D`／`sub_1D5D4`／`sub_10C14`／`sub_189DE`）
- 推論等級：**呼叫鏈 confirmed**；像素產生方式**未知**

## 1. 原版做什麼

指令列 ＝ **一個框 ＋ 一串文字 ＋ 一個熱區**，沒有按鈕圖：

```
sub_1895D(樣式=0x0B, X=0, Y=0, 尺寸=0x021B)   框
sub_106F5(字串=cs:6181h, X=0x28, Y=8, 屬性=0x0F01)  八個指令名
sub_1E3D7(熱區=0x0C, X=0x28, Y=0x18, 尺寸=0x0230)
```

框由三層寫成，**三層都不碰圖庫**：格子屬性表（`cs:word_1D84E`，
每格 8 B、40 欄）、框線層（`cs:word_10D4A`，四條邊分開畫）、熱區。

座標單位：呼叫端傳 16 px 的粗格，內層轉成 8 px 格再轉像素。
指令列的 `0x021B` 換算後是 **432 × 32 pixel**，在畫面左上角。

## 2. 演算法

**還沒解。** `sub_10C60`／`sub_10C77` 怎麼把框線變成像素、
`cs:word_1D84E` 每格那 8 bytes 是什麼，都未讀
（[`docs/re/46`](../re/46-strategy-chrome-cell-layer.md) §4）。

**這一節填不出來之前，這份規格不能升 READY。**

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 呈現層 | `cmd/wlgame/strategyhud.go` 的 `drawNaturalStrategyHUD`、`internal/ui/chrome` |
| 差異 | **整段是自繪的近似**。外框用 `chrome.Window`、勢力色標用 `vector.DrawFilledRect`、顏色硬寫 RGBA |
| 證據來源 | 該函式自己的註解寫明是「影片 `af6xqcicXoI` 的 478×360 影像經黑邊還原」——**壓縮過的參考影片，不是執行檔** |

**指令列這一項倒是結構對的**：原版就是「框 ＋ 文字」，remake 也是。
不對的是框怎麼畫，以及字的來源（原版取 `cs:6181h`，remake 用自己的字串陣列）。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `cmd/wlgame/strategyhud_test.go`（驗的是**我們自己的**版面常數，不是原版）|
| 對原版 | **未做**。管線已就緒（[`docs/playtest/21`](../playtest/21-dosboxx-bridge-sampling.md)），可在原版執行期直接讀 VRAM 對拍 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| `sub_10C60`／`sub_10C77` | 框線的實際畫法——**主畫面 parity 的關鍵路徑** |
| `cs:word_1D84E` 每格 8 bytes | 內容意義未讀，消費端未找 |
| `cs:6181h` 的八個指令名 | 字串 bytes 未 dump |
| 樣式碼的值域 | 只確定 `0`＝擦除、`0x0B`＝指令列 |
