# 124 — 反白是「色號 XOR 12」，不是換一個底色

**狀態：CONFORMED。** 選單的反白列與指令列的反白格走的是同一支
VGA 常式：把那一塊的**色號 XOR 12**。黑底白字因此變成**黃底藍字**，
而 remake 先前畫的是綠底 ＋ 一圈米色外框，而且指令列**根本沒畫反白**。

- 日期：2026-09-03
- 實機證據：`workplace/promo-live/parity-tap5/menu.png`
  （松崗 DOS/V，196年4月20日，`docs/playtest/60`）
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_161CA`（`000161CA`）
  → `sub_10B46`（`00010B46`）；選單那一側是
  `sub_1036F` → `sub_10414` → `sub_1061F` → 同一支 `sub_10B46`
- 推論等級：**confirmed**（逐像素量到 ＋ 呼叫鏈在反組譯裡）
- 相關：[`54`](54-ui-colours-from-palette.md)（介面顏色查調色盤）、
  [`125`](125-menu-box-width-from-padding.md)（同一張截圖抓到的另一件事）

## 1. 原版做什麼

`sub_10B46` 是一支操作 VGA 繪圖控制器（`0x3CE`／`0x3CF`）的區塊常式，
把 `(dx, bx)` 起 `si × di` 的範圍**與 `0Ch` 做互斥或**。
呼叫它兩次就回到原狀——原版就是這樣「反白 → 做事 → 還原」的：

```asm
sub_161CA:                       ; 指令列的命中與反白
        mov ax, cx / sub ax, 18h ; 命中：x − 24
        cmp ax, 180h / jnb 不理   ; 八格 × 48 ＝ 384
        mov cl, 30h / div cl     ; 索引 ＝ (x − 24) ÷ 48
        …
        mov dx, ax               ; 格的左緣 ＝ 索引×48 + 24
        mov bx, 28h              ; ★ Y  ＝ 40
        mov si, 30h              ; ★ 寬 ＝ 48
        mov di, 10h              ; ★ 高 ＝ 16
        call sub_10B46           ; ← 反白
        call cs:funcs_161FE[bx]  ; ← 做那一格的事（整段期間都亮著）
        call sub_10B46           ; ← 再 XOR 一次 ＝ 還原
```

⚠ **`bx = 28h` 是 Y 不是寬。** 先前 `strategyCommandCellRect` 的註解記成
「原版另外用 40 px 寬畫高亮」，於是把「反白矩形」與「命中矩形」當成兩個東西。
它們是同一個：**(索引×48 + 24, 40, 48, 16)**。

## 2. 量到的顏色

`parity-tap5/menu.png` 的「軍團」格（216, 40, 48, 16）與 remake 同一格：

| | 底 | 字 |
|---|---|---|
| 原版 | **436 px** `(243,227,0)` ＝ 色 12 | **204 px** `(48,65,211)` ＝ 色 3 |
| remake（修前）| 436 px 色 0 黑 | 204 px 色 15 白 |

⭐ **兩邊的像素數一個不差**——字形與位置本來就對，差的只有那一次 XOR：

```
0（黑底）  XOR 12 = 12（黃）
15（白字） XOR 12 = 3（藍）
```

選單框裡的反白列同理：`(200, 72, 96, 16)` 全部是黃底藍字，
**沒有外框線**，而且**與那一列的文字對齊**（不是往上一格）。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 兩個顏色 | `internal/ui/chrome`：`HighlightXOR = 12`、`HighlightIndex = InkIndex ^ 12`、`HighlightInkIndex = PaperIndex ^ 12`，值一樣跟著 `GAMEPAL.BRG` 走（[`54`](54-ui-colours-from-palette.md)）|
| 選單列 | `cmd/wlgame/talkscene.go` 的 `drawLegacyChoiceBox`：填 `chrome.Highlight`、字用 `chrome.HighlightInk`，**拿掉外框線與 `y−1`** |
| 指令列 | `cmd/wlgame/strategyhud.go` 的 `drawNaturalStrategyHUD` ＋ `activeCommandCell()`，矩形直接用 `strategyCommandCellRect` |
| 差異 | 用「填黃 ＋ 字畫藍」代替真正的 XOR。框內底恆為黑、字恆為白，所以**逐點等價**；要在別種底色上反白才會分歧 |

⭐ **清單視窗的反白條走的是同一條規則**（§3.6）：`chrome.SelectIndex`
就是清單底色 9 XOR 12，不是另一個量出來的常數。

## 3.5 ⭐ 反白涵蓋八格，不是只有彈出選單那三格（2026-09-06）

`sub_161CA` 逐行讀完（74 B）之後，這一點沒有解釋空間：

```asm
000161E4  bx = 28h / si = 30h / di = 10h     ; 反白矩形 (24+48n, 40, 48, 16)
000161ED  push dx / push bx / push si / push di
000161F1  call sub_101B4 / call sub_10B46 / call sub_101DB   ; ★ 反白
000161FE  call cs:funcs_161FE[bx]            ; ← 那一格的 handler
00016203  pop di / pop si / pop bx / pop dx
00016207  call sub_101B4 / call sub_11C8D / call sub_10B46 / call sub_101DB
                                             ; ★ 再 XOR 一次 ＝ 還原
00016213  retn
```

⭐ **XOR 在 `call` 的前後各一次，而且不分索引**——八格走的是同一段程式碼，
所以**只要那一格的流程還沒回來，它就一直亮著**。

原版擷取（dosgolem，[`../playtest/81`](../playtest/81-command-cell-highlight.md)）：
進言的五項選單開著時「進言」亮、財政視窗開著時「財政」亮、
編成的武將一覽開著時「編成」亮——與這一段逐格對上。

### 3.5.1 remake 怎麼判斷「這一段還沒回來」

原版靠**呼叫堆疊**，remake 沒有那個結構（每個流程是一個狀態物件）。
所以改成**一格一個謂詞**：

| 格 | 謂詞 | 有原版擷取嗎 |
|---|---|---|
| 0 進言 | `adviseActive()` | ✅ [`../playtest/81`](../playtest/81-command-cell-highlight.md) |
| 1 人事 | 彈出選單開著，或它開出來的一覽還開著 | ✅ [`../playtest/61`](../playtest/61-city-personnel-menu-parity.md) |
| 2 財政 | `finance.active` | ✅ [`../playtest/81`](../playtest/81-command-cell-highlight.md) |
| 3 編成 | `form.active \|\| list != nil` | ✅ 同上 |
| 4 軍團 | 同 1 | ✅ [`../playtest/60`](../playtest/60-corps-menu-parity.md) |
| 5 據點 | 同 1 | ✅ [`../playtest/61`](../playtest/61-city-personnel-menu-parity.md) |
| 6 武將 | `list != nil` | ⚠ **沒有原版擷取** |
| 7 勢力 | `list != nil` | ⚠ **沒有原版擷取** |

⚠ **6／7 兩格是照公式接的，不是量出來的。** 原版那兩格走的是狀態列提示
＋ 地圖游標（[`../re/22`](../re/22-strategy-command-tree.md) §3.4），
而 remake 開的是一覽表——**流程本身就不一樣**，所以那兩格的反白時機
只能算強證據。要降級成 confirmed 得各拍一張。

## 3.6 ⭐ 清單視窗也是同一條規則，連換色的儲存格都是（2026-09-07）

[`../playtest/98`](../playtest/98-enemy-corps-panel.md) §3 量到軍團一覽
「委任」那一格在兩種狀態下的顏色，當時的結論是
「兩組都不是同一種變換，所以沒有硬推一條規則」。**那個結論是錯的**——
把三組都算一次 XOR 12 就對上了：

| 量到的 | 一般列 | 反白列 | XOR 12 |
|---|--:|--:|---|
| 清單底色 | 9 `(243,211,146)` | 5 `(81,146,65)` | 9 ^ 12 ＝ **5** ✅ |
| 一般文字 | 0（黑）| 12 `(243,227,0)` | 0 ^ 12 ＝ **12** ✅ |
| 「委任」那一格 | 10 `(211,0,0)` | 6 `(130,65,32)` | 10 ^ 12 ＝ **6** ✅ |

⇒ **反白不分「底色／文字／儲存格自己的顏色」，整塊一起 XOR。**
這正是 §1 那支 `sub_10B46` 的性質：它是一支**區塊**常式，
對矩形裡的每一個像素做同一件事，根本不知道哪一點是字、哪一點是底。

⭐ 之所以看起來像三條規則，是因為當時把「黑 → 黃」讀成「底色被換掉了」。
**一個色號一個色號地算 XOR，三組立刻收斂成一條。**

⚠ 這也讓 `chrome.SelectIndex = 5` 從「量出來的常數」變成**推導值**：
清單底色 9 XOR 12。`docs/re/27` §4 的 `ah = [si+6] | 0x90` 是同一件事的
另一半——顏色是屬性的**低半位元組**，所以 XOR `0Ch` 只動前景。

### 3.6.1 remake 要怎麼改

| 項目 | 作法 |
|---|---|
| `chrome.SelectIndex` | 改成 `SheetIndex ^ HighlightXOR`，不再是獨立常數 |
| `listInkWarnSelected` | 刪掉。反白列的儲存格色號改成 `idx ^ chrome.HighlightXOR` |
| 新的換色 | 只要給**一般列**的色號，反白自動正確——先前每加一種換色都要另外量一次 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 對原版 ✅ | [`../playtest/60`](../playtest/60-corps-menu-parity.md)：指令列 `command` 區與選單框 `(192,64,112,48)` **各 0 px** |
| 單元測試 | `TestHighlightIsXorOfInkAndPaper`（`internal/ui/chrome`）：兩個色號等於 XOR 12，且等於實機量到的 12／3 |
| 單元測試 | `TestCommandHighlightRectMatchesHitRect`（`cmd/wlgame`）：軍團格 ＝ (216,40,48,16) |
| 單元測試 | `TestActiveCommandCellFollowsCorpsMenu`（同上）：選單開著才亮 |
| 單元測試 | `TestActiveCommandCellCoversEveryFlow`（同上）：八格各自的謂詞，流程結束就熄 |
| 對原版 ✅ | [`../playtest/81`](../playtest/81-command-cell-highlight.md)：進言／財政／編成三格的 `command` 區各 0 px |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 武將／勢力兩格的反白時機 | **八格都接了**（§3.5），但這兩格**沒有原版擷取**：原版走狀態列提示 ＋ 地圖游標，remake 開的是一覽表，流程本身不同。其餘六格各有一張原版擷取對過 0 px |
| `sub_10B46` 的暫存器序列 | 只確認了它寫 `0Ch` 給繪圖控制器、而結果逐點等於 XOR 12。**中間那幾個 port 寫入沒有逐行讀** |
