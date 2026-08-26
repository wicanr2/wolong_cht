# 88 — 三處顯示與原版對不上

**狀態：CONFORMED（2026-08-26 實作並實跑驗過）。**

- 日期：2026-08-26
- 出處：使用者看推廣片回報三件事，逐項比對實機錄影與既有 oracle 之後
  全部成立。兩版都比過（DOS/V 實機錄影 ＋ PC-98 oracle 截圖），
  結論一致——這三處不是版本差異。
- 推論等級：§1 **confirmed（實機像素）**、§2 **remake 自創的東西壞了**、
  §3 **confirmed（靜態，`docs/spec/58` 自己寫著答案）**

## 1. TALK 對話框的底是黑的，不是藍底龍紋

| 來源 | 框內主要顏色 |
|---|---|
| 松崗 DOS/V 實機（`o5-battle.mp4` 第 101 秒）| `(0,0,0)` 6,298 px（其餘是影片壓縮雜訊）|
| PC-98 oracle（`docs/images/pc98-oracle-event3-choice.png`）| 黑 |
| **remake 現況** | 黑 7,459 px ＋ **`(0,32,97)` 5,801 px** ← 龍紋藍 |

`drawLegacyTalkBox` 用的是 `chrome.Window(..., chrome.Menu)`，而 `Menu`
會鋪龍紋（`internal/ui/chrome.fillInterior`：`fill == Menu` 才鋪）。
TALK 框要的是 `chrome.Blank`——同一支 `Window` 已經支援，只是傳錯了值。

⭐ **不是所有視窗都要改。** 原版的系統選單確實是藍底龍紋
（[`../playtest/39`](../playtest/39-system-window-parity.md) 逐像素 PASS），
據點情報卡則是黑底（`pc98-oracle-city-panel.png`，remake 也已經是黑的）。
**改的只有 TALK 對話框那一種。**

### 1.1 文字沒有壓到肖像

使用者看到的「字壓在人像上」出自推廣片裡的 `wlgame-event3-choice.png`
——那是 **2026-08-10 的截圖**，現行版本的幾何已經沒有重疊：
肖像 `(168,168)` 64×64 → 右緣 232，文字起點 240，差 8 px。

> **這一條要記成規則**：**推廣片裡的靜態圖會凍結當時的 UI。**
> 程式改了、圖沒重拍，片子就同時放著三個世代的畫面——
> 而觀眾看到的是「同一款遊戲的對話框有三種樣子」。
> 每次重剪推廣片要連靜態圖一起重拍，不能只換動態段。

## 2. 事件列（「月結」那個框）不會消失

`g.lastEvent` 從頭到尾沒有人清掉，所以第一次月結之後那個框就**永遠留在畫面上**。
量過：大地圖 330 幀裡第 148 幀出現，之後 182 幀一直在。

原版的大地圖**沒有這個框**（`o2-strategy.mp4`、`o4-march.mp4` 全程沒有），
事件是用肖像對話框跳出來、玩家按鍵關掉。這個事件列是 **remake 自己加的**
（原版的月結只是數字變了），所以處置不是「照原版拿掉」而是「讓它像個提示」：

- 設定時記下當時的幀數，**六秒後自己消失**——與手機版同一個行為
  （`internal/ui/phone` 的事件通知，`TestEventNoticesAppearAndExpire`）。
- 六秒是手機版已經在用的值，兩端一致比另外挑一個數字好。

## 3. ⭐ 小兵只畫得出一半

`BATTLE.SCH` 的一個兵是 **16×64**：上半（奇數 unit）是頭與上身、
下半（偶數 unit）是身體與腿（`workplace/promo/dump/sprites.png` 看得到）。
`appendTallDisplayUnits` 把兩半推進顯示表——**上半在 (row, z)、
下半在 (row−1, z+1)**，這一段有測試釘著，是對的。

問題在**畫的時候只掃到 z**：

```go
// makeDisplayInfo（現況）
s := grid[row][col][z][0]   // ← 只看 lane 0
info[row][col].height = 2 * z
```

`displayDepthRange` 拿這個 `height` 決定畫哪幾層。平地上地形只有第 0 層，
於是 `z1 = 0`，**兵的下半（第 1 層）永遠畫不出來**——只剩頭與上身。
牆邊、屋簷下的兵看起來是完整的，因為鄰格的地形夠高把範圍撐開了。

而 [`58`](58-display-slot-depth-range.md) §1 早就寫著答案：

```
+1   高度（**含物件**）
+3   高度（**只有地形**）——物件擦掉時用它還原 +1
```

⭐ `+1` 與 `+3` 不是同一份的兩個複本：**物件 producer 只抬高 `+1`**
（`sub_1DB34`）。`sub_1DDB4` 取鄰格高度用的是 `+1`，也就是**含物件**那一份。
remake 把 `height` 算成只看 lane 0，等於用了 `+3`。

**改法**：`makeDisplayInfo` 的 `height` 兩個 lane 都算，
地形記 `2z`、物件記 `2z+1`（`58` §1 的編碼），`start` 維持只看 lane 0 的小 unit。

> **教訓寫成規則**：**規格寫對了不等於實作照著做。**
> `58` 是 CONFORMED，欄位語意一字不差，而實作只取了兩個欄位裡的一個；
> 對拍數字（`field` 3.4%）沒有把它抓出來，因為當時那一格的兵不多。
> **CONFORMED 的意思是「量過的那個局面對得上」，不是「每一條都實作了」。**

## 4. 改動

| 檔 | 內容 |
|---|---|
| `cmd/wlgame/talkscene.go` | `drawLegacyTalkBox` 的底改 `chrome.Blank` |
| `cmd/wlgame/main.go` | `lastEvent` 加時戳，六秒後不畫 |
| `internal/ui/isoview/isoview.go` | `makeDisplayInfo` 的高度含 lane 1（物件 `2z+1`）|

## 5. 驗證（2026-08-26）

| 項目 | 結果 |
|---|---|
| `TestDisplayHeightIncludesObjectsSoSoldiersDrawWhole` | 兵那一格高度 1（`2z+1`）、下半那一格 3；深度範圍涵蓋第 1 層；`start` 不被物件推高 |
| `TestEventLineExpiresAndRestartsOnNewText` | 六秒內看得到、之後不畫；換一句會重新計時 |
| TALK 框實測 | 框內 13,243 px 全黑，**沒有一個 `(0,32,97)`** |
| 系統選單回歸 | 仍有 `(0,32,97)` 3,436 px——藍底龍紋沒被改到 |
| 戰場實測 | 開闊地的兵有腿了（`workplace/promo/dump/rm-open-{before,after}.png`，差 **96 px**＝三個兵的下半）|
| 對拍回歸 | `SAVE-E` 那個 fixture 修正前後**逐位元組相同**——那一格的兵靠著城壁，鄰格高度本來就把範圍撐開了 |

⚠ **對拍數字沒有變，不代表沒有問題。** 已量過的那個局面剛好都是牆邊的兵，
所以 `playtest/40` 的 3.4% 一路都沒把這件事抓出來。**開闊地才是會出事的地方，
而對拍沒有選在那裡量。**

## 6. 未解

| 缺口 | 下手點 |
|---|---|
| 事件列本身是 remake 自創 | 原版怎麼提示月結（如果有）沒查過；目前只是讓它不擋畫面 |
| `playtest/40` 沒有涵蓋開闊地的兵 | 那一份量的兩個局面都在城壁邊。要擋住這一類回歸，得再加一個**開闊地**的對拍 fixture |
