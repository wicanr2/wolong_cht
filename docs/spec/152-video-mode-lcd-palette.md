# 152 — 畫面模式：「１６色」／「液晶」是同一份調色盤的兩組 bank

**狀態：CONFORMED。** 系統選單第 1 列切的是 `GAMEPAL.BRG` 的
**bank 0–3 ↔ bank 4–7**——同樣四季、換一整組顏色。
remake 先前那一列**點下去沒有任何反應**，值格固定顯示「１６色」。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）
  - 字串表 `ds:6002h`：「１６色」／「 液晶 」（[`../re/55`](../re/55-system-menu-window.md) §4）
  - `GAMEPAL.BRG` 的 8 組（[`../formats/02`](../formats/02-brg-palette.md) §6）
  - 說明書：「アナログカラーと8階調液晶とを切り替えます」
    （[`../reference/01`](../reference/01-jp-manual.md)）
- 推論等級：**confirmed**（字串表 ＋ 調色盤檔 ＋ 原版擷取 `orig-sys1`）
- 驗證：[`../playtest/100`](../playtest/100-video-mode-parity.md)（地圖 0 px、系統選單逐列 0 px）
- 相關：[`13`](13-main-window-toggles.md) §「畫面模式」、
  [`54`](54-ui-colours-from-palette.md)（顏色一律查調色盤）

## 1. 原版做什麼

點系統選單第 1 列的值格（**右邊那個格子，不是左邊的標籤**）就切換，
`ＯＫ` 那種確認一步都沒有——**當場整個畫面換色**。

原版擷取 `orig-sys1`：切成「液晶」之後，地表由綠變黃、標題由紅變桃紅、
視窗底紋由藍變黑白——**每一個色號都變了**，版面一格沒動。

⭐ **bank 就是調色盤，不是美術**。`TileSet.RenderRGBA(data, index, pal, bank)`
的 `bank` 只餵給 `pal.Bank(bank)`；圖塊資料是同一份。所以「換畫面模式」
＝ **把所有取 bank 的地方 +4**，一行圖形程式都不必動。

## 2. remake 要怎麼改

| 項目 | 作法 |
|---|---|
| 狀態 | `game.videoLCD bool`——**零值 ＝ 16 色**，與原版開機預設一致 |
| bank | `g.paletteBank()`：`季節 + 4×(videoLCD)`。**所有傳 bank 的地方共用這一支**，不要各自 `int(g.world.Clock.Season())` |
| 切換 | `dispatchSystemRow` 補 `sysRowVideo`：左右鍵都 toggle（只有兩個值，同「主君編成」那一列的理由）|
| 值格 | `videoModeLabels[0|1]`，字串已經在 |
| 差異 | 無 |

### 2.1 ⭐ 外框與底紋也要跟著重畫

`chrome.Load(lib, bank)` 是**開局時載一次**的。換季只動色號 14
（[`../formats/02`](../formats/02-brg-palette.md) §4）所以看不出來，
換畫面模式一眼就看得出來：地圖整片變黃了，視窗還是藍底紅框。

記住「現在的外框是用哪一組畫的」（`game.chromeBank`），`Update` 裡比一次，
不同就 `loadChrome()`。⭐ **比的是 bank 不是旗標**，所以順帶把換季也涵蓋了。

⚠ **`Clock.Season()` 有兩種用途**：一種是調色盤 bank（要跟著換），
另一種是**素材頁**（肖像的四季版本、據點景觀圖）。這一份只動前者。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 對原版 ✅ | [`../playtest/100`](../playtest/100-video-mode-parity.md)：視窗外的整片大地圖 **0 px** |
| 單元測試 | `TestPaletteBankFollowsVideoMode`：四季 × 兩種模式 ＝ 0–7，且零值是 16 色 |
| 單元測試 | `TestSystemMenuVideoRowToggles`：點那一列會切、值格文字跟著換 |

## 4. 未解

| 項目 | 現況 |
|---|---|
| 素材頁的四季 | 肖像與據點景觀圖是**另一種**四季（換的是圖不是色），這一份沒動它們——切到液晶時那些圖的顏色會不會也跟著換，還沒對過 |
