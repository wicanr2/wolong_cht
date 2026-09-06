# 82 — 戰略畫面的六張對拍：日期也對齊之後，四張整片 0 px

**狀態：通過。** 加了 `-fixture-when` 讓 remake 那一側「先跑到那一天再擺
fixture」，兩邊用**同一個遊戲時刻**取樣，於是每一張戰略畫面都剩著的那塊
`banner` 116 px（日期）**收掉了**。六張裡四張**五區全部 0 px**。

- 日期：2026-09-06
- 規格：[`../spec/118`](../spec/118-shot-when-condition.md) §2.3（`-fixture-when`）、
  [`../spec/124`](../spec/124-menu-highlight-xor.md)、
  [`../spec/140`](../spec/140-status-message-box.md)、
  [`../spec/38`](../spec/38-list-windows.md)（空列破折號）
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-saveb tools/dosgolem.sh
  workplace/dosgolem/advise "wait;click:320,200;click:300,151;until:196/4/20;
  sclick:352,15;save:cmd;stap:<x>,47;steps:400000;shot:<名>"`
  （`<x>`：進言 48、人事 96、財政 144、編成 192、軍團 240、據點 288）
- remake 側：`tools/parity_shot.sh out.png -direct -scenario 0 -player 0 -seed 7
  -save-file workplace/dosgolem/root-saveb/SAVE.DAT -load-slot 0 <fixture>
  -cam 0,0 -fixture-when clock:196/4/20 -shot-when clock:196/4/20`

## 1. 結果

| 畫面 | `banner` | `command` | `map` | `minimap` | `faction` |
|---|---:|---:|---:|---:|---:|
| 進言（五項選單）| **0** | **0** | **0** | **0** | **0** |
| 人事（四項選單）| **0** | **0** | **0** | **0** | **0** |
| 軍團（兩項選單）| **0** | **0** | **0** | **0** | **0** |
| 據點（兩項選單）| **0** | **0** | **0** | **0** | **0** |
| 編成（武將一覽）| **0** | 88 | **0** | **0** | **0** |
| 財政（財政視窗）| **0** | 88 | 141 | **0** | **0** |

那兩塊非零都有解釋：

| 塊 | px | 是什麼 |
|---|---:|---|
| `command` | 88 | **原版自己畫的滑鼠游標**（14×14，停在剛點過的那一格上）。remake 的截圖模式不畫——與 [`76`](76-battle-talk-parity.md) §4 的 95 px 同一類 |
| `map`（財政）| 141 | **M7 的校訂**：TALK #16「財政**予定**」→「財政**計畫**」（`translations/corrections.json` 的 `id: 16`）。⭐ 差異的位置正好落在改過的那兩個字上——**那是接對索引的證據**，不是缺陷 |

⭐ 軍團／據點／人事那三張把 [`60`](60-corps-menu-parity.md)／[`61`](61-city-personnel-menu-parity.md)
記著的「`banner` 各剩 110–116 px」也一起收掉了。

## 2. ⭐ `-fixture-when`：一個順序問題

`-shot-when clock:196/4/20` 早就有了，缺的是**讓 fixture 也等到那一刻**：
remake 的驗收 fixture 是**啟動時**擺的，那一刻時鐘還在存檔的日期上，
而窗一開時間就停（`timeRuns()`），追不上去。

第一版把檢查寫在 `Update` 的開頭，結果是 **「日期對了、視窗沒開」**：

> 時鐘是在 `Update` 中間走的，而截圖的判定在 `Draw`。
> 走到那一天的**那一幀**，`Update` 開頭看到的還是前一天 → fixture 晚一幀擺；
> 而 `Draw` 已經在那一幀拍掉了。**症狀看起來像 fixture 壞了。**

改成 `defer`（在這一幀的世界推進完之後才檢查）就對了。

## 3. 抓到的東西：空列的破折號差一個半格

編成那一張最後剩 42 px，全部在**一條掃描線**上（y=255，第十列＝空列）。
量出來：

| | 武術 | 統率 | 政治 |
|---|---|---|---|
| 數值右緣（右靠）| 128 | 168 | 208 |
| **空列破折號左緣** | 原版 **120**／remake 112 | 160／152 | 200／192 |

⇒ **原版的分隔線與欄界是兩份資料。** 數值右緣說欄是 `[112,128)`，
而破折號從 120 起 ＝ 欄起點 ＋ 8。軍團一覽卻是**貼欄起點**的
（parity-menus7 的 m1，[`../spec/38`](../spec/38-list-windows.md)）——
兩張的關係不同，同一條 `Sep` 推不出來。

remake 先前用同一條 `Sep` 推出欄界與破折號，於是武將一覽的空列破折號
整排左移 8 px。加 `listFamily.NumericDashInset` 之後 `map` 收到 **0 px**。

⚠ **十列裡只有沒有資料的那幾列看得到**——這一張剛好只有一列是空的，
所以它只值 42 px。清單短的時候會放大成十倍。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 其餘家族的破折號縮排 | 只有軍團（0）與武將（8）量過。據點／勢力那兩張**沒有原版擷取**，維持 0 |
| 財政的 141 px | 校訂造成的，**刻意的差異**。要 0 px 得拿未校訂的文本跑，那不是遊戲會出的畫面 |
| 原版擷取裡的滑鼠游標 | 88 px。要消掉得讓原版把游標移開再截圖——`move:` 之後游標會拖動鏡頭（[`../re/84`](../re/84-popup-row-band-and-world-cursor.md) §2），得先確認拖不動的位置 |
| 武將／勢力兩格的指令列反白 | 原版走狀態列提示 ＋ 地圖游標，remake 開一覽表，**流程不同**（[`../spec/124`](../spec/124-menu-highlight-xor.md) §5）|
