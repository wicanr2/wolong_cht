# 65 — 原版 oracle 換成 dosgolem：五格 × 五區全部 0 px

**狀態：通過。** 用 dosgolem 跑松崗 DOS/V 版 `KI.EXE`，照
`workplace/promo-live/parity-main/capture-metadata.txt` 的時間軸走，
開機 → 劇本選單 → 君主清單 → 選定君主 → 遊戲主畫面五格畫面，
對同一份 DOSBox-X 唯讀擷取**每一區都是 0 個不同像素**。
⭐ 整條鏈 **0.99 秒**；DOSBox-X 那邊同一條時間軸的 `wait` 加起來是 146 秒。

- 日期：2026-09-05
- 原版側：`workplace/promo-live/parity-main/` 的五張 `z*.png`（松崗 DOS/V 唯讀實機畫面，
  `machine=vgaonly`／`core=normal`／`cputype=486`／`cycles=fixed 20000`），
  用 `tools/parity_crop.py` 裁成 640×400
- dosgolem 側：`tools/dosgolem.sh workplace/dosgolem/golem-run "<時間軸>" dosbox`
- 工作副本：`~/cht/dosgolem-wolong`，HEAD `6689407`
- 比對：`tools/parity_diff.py`（分區座標出自 [`../re/46`](../re/46-strategy-chrome-cell-layer.md)，從機器碼算的）

## 1. 結果

| 畫面 | banner | command | map | minimap | faction |
|---|---:|---:|---:|---:|---:|
| `z1-newgame` NEW GAME 確認框 | 0 | 0 | 0 | 0 | 0 |
| `z2-scenario` 劇本選單 | 0 | 0 | 0 | 0 | 0 |
| `z3-leaders` 君主清單 | 0 | 0 | 0 | 0 | 0 |
| `z4-pick` 選定君主 | 0 | 0 | 0 | 0 | 0 |
| `z5-main` 遊戲主畫面 | 0 | 0 | 0 | 0 | 0 |

時間軸與原版擷取那一份**逐字相同**，只有 y 座標換算不同（§3）：

```
wait;shot:z1;click:320,215;wait;shot:z2;click:300,190;wait;shot:z3;
click:450,154;wait;shot:z4;click:360,336;wait;shot:z5
```

## 2. 這一輪改到的是 dosgolem，不是 remake

`KI.EXE` 餵進 dosgolem **不必改任何程式碼**就跑得動 330 萬道指令、
開得了 14 個資料檔。要補的是機器層的兩件事：

- **VGA 16 色平面模式**（mode 12h）——原本只有 mode 13h 的線性緩衝區
- **DOS/V 字型服務 `INT 15h AH=50h`**——照 [`../re/29`](../re/29-font-service-int15.md)
  的索引式接，字型檔用 `END_S13.DAT`／`END_S14.DAT`

以及遊戲時鐘：`int 61h AH=0Ch` 登記的回呼要真的發
（[`../re/61`](../re/61-timer-tick-source.md)）。

> ⭐ **第一次跑的錯誤訊息指向錯的地方。** 它跑了 330 萬道指令然後
> `int 21h AH=4Ch` 帶著離開碼 8 離開，印的是
> `ERROR: Memory not enougth ! ( at least 560KB )`——離開碼、訊息與
> 最後那道 `int 21h` **一致地**指向記憶體。真正的原因在 200 萬道指令之前：
> 沒有字型服務 → 兩條 `call far 0000:0000` 沒被 patch → 第一次畫字就是
> 一次遠呼叫到零位址，落點恰好是開機配置記憶體那一段。
> **能查出來的唯一理由是指令軌跡。**

## 3. ⚠ 視窗座標換算的分母是 479

DOSBox-X 的視窗是 640×480，`int 33h` 把**整個視窗**等比對映到遊戲的
640×400，所以

```
遊戲 y ＝ 視窗 y × 399 ÷ 479
```

不是 `× 400 ÷ 480`。大部分點兩種算法一樣，**只有少數差 1**：
視窗 336 是 279 不是 280。用 280 的話主畫面其餘四區全 0，
只有 `map` 差 62 點——而那 62 點是游標整塊偏一列。
**看起來像「只差一點點」，不像「換算式錯了」。**

[`../re/43`](../re/43-open-questions.md) 兩處寫的「送 y 要乘 1.2」**沒有錯**
——那是反方向（遊戲 → 視窗），而 `479÷399 ＝ 1.2005`，在 0–399 這個範圍內
四捨五入之後與 ×1.2 幾乎處處相同。**會咬人的是反過來除**：
`÷1.2` 與 `×399÷479` 在少數點上差 1，而視窗 336 正好是其中一個。

## 4. ⭐ 即時制的取樣點可以寫成遊戲日期

這一款是即時制，同一串操作在牆上時鐘的不同時刻停在不同的遊戲日期。
DOSBox 那邊只能用 `wait:3` 之類的秒數逼近——[`39`](39-system-window-parity.md)
的原版擷取停在 **4月9日**、remake 停在 4月1日，那 158 個不同像素
其實是日期，而當時記成「橫幅差三段」。

dosgolem 直接讀原版的時鐘欄位（`ds:0CF0` 日／`ds:0CF4` 月／`ds:0CF6` 年，
[`../re/06`](../re/06-game-clock.md)），所以取樣點寫成

```
until:196/4/9
```

實測時鐘推進速度：**約 500 萬道指令一個遊戲日**（第一章、預設戰略速度）。

## 5. 未解

| 項目 | 現況 | 下手點 |
|---|---|---|
| 載入既有存檔 | 沒試 | 把 `-save` 之類的旗標接進 `apps/wolong`，或直接讓遊戲走「載入」選單 |
| 戰術畫面 | 沒試 | 戰場走同一套繪製層，預期沒有新的機器層缺口——**但沒驗過就是沒驗過** |
| 事件對話會擋住時鐘 | 已觀察到：跑到 196年4月3日 2時會出現「就將呂布擊潰吧」的對話框，時鐘停住等點擊 | 那是原版行為不是缺口；腳本要把它點掉 |
| `STR.EXE` 的檔名懸案 | [`../re/29`](../re/29-font-service-int15.md) §6 仍未裁決 | dosgolem 的字型檔名是參數，把它改成 `END_S10/S11` 跑一次就知道 |
| PC-98 版 | dosgolem 沒有 PC-98 的機器層 | 那一版仍走 DOSBox-X |
