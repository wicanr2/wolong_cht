# 101 — 「遊戲結束」的兩項選單：框內逐像素 0 px

**狀態：通過。** 系統選單第 6 列跳出的「　終　　了　」／「　取　　消　」
兩項選單，**框內 5,376 px 全部相同**——連反白列的黃底藍字都對上。
remake 先前這一列走的是 F10 那條置中的 ＹＥＳ／ＮＯ 對話框。

- 日期：2026-09-06
- 規格：[`../spec/153`](../spec/153-quit-confirm-menu.md)
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh
  workplace/parity/sys5 "wait;click:320,200;click:300,151;runto:11CD0;
  sclick:448,15;steps:400000;stap:376,280;steps:900000;clock;shot:orig-sys5b"`
- remake 側：`tools/parity_shot.sh out.png -direct -scenario 0 -player 0 -seed 7
  -save-file workplace/dosgolem/root-noclouds/SAVE.DAT -load-slot 0
  -open-window 3 -quit-menu -cam 0,0`（`-quit-menu` 本輪新增）

## 1. 結果

| 比什麼 | 不同像素 |
|---|---:|
| `banner`／`command`／`minimap`／`faction` 四區 | **0 / 110848** |
| **選單框 (352,272,112,48)** | **0 / 5376** |
| 系統視窗的標題 ＋ 前六列 | 272 / 36608（＝音效值格，§2）|
| 視窗外的大地圖（0,64,208×336）| **0 / 69888** |
| 第 6 列以下（208,288,208×64）| 9908 / 13312（＝remake 多的兩列，§2）|

`map` 區那 10,180 px **全部**由上表最後兩列解釋完（272 ＋ 9,908），
沒有剩下的。

## 2. 兩個已知差異，都不是這一輪的

| 差異 | 出處 |
|---|---|
| 音效值格顯示「未接入」而不是「TYPE 1」 | 容器裡沒跑過 `tools/bgm2ogg.sh`，`Bank.Available()` 是 false（[`../spec/29`](../spec/29-audio.md) §5.1）。**驗收捷徑不要清空音檔目錄**——那會改到被驗收的畫面 |
| remake 的系統選單多「主君編成」「損害報告」兩列 | [`39`](39-system-window-parity.md) 記過的 remake 差異 |

## 3. ⭐ 選單開著時原版沒畫游標

框內 0 px 的意思是**連游標都沒有**——`stap:376,280` 之後游標就停在
(376, 280)，那個點落在框內第 0 列上，而原版一個像素都沒畫。

這與 [`../re/88`](../re/88-mouse-cursor-visibility.md) 解出來的規則一致——
畫不畫由 `byte_20100` 決定。四筆觀察現在都有解釋：

| 擷取 | 畫面 | 游標 |
|---|---|---|
| `orig-gen1` | 武將一覽（清單視窗）| ✗ |
| `orig-gov1` | 內政官任命（清單視窗）| ✅ |
| `cur-a` | 清單反白但沒決定 | ✗ |
| **`orig-sys5`** | **兩項選單 `sub_193E9`** | **✗** |

`orig-gen1` 與 `orig-gov1` 是同一種視窗卻一個有一個沒有——因為變數
不在畫面上，而是那個旗標。**remake 這一側還沒同步 hide／show**
（[`../re/88`](../re/88-mouse-cursor-visibility.md) §5），
目前只在地圖選點時自繪，這一張因此天然對齊。

## 4. 重跑一致

同一串 verbs 隔一輪重跑，與 9 月 6 日那一張 **0 / 256000**——
即時制的取樣點靠受控存檔（[`../spec/147`](../spec/147-controlled-parity-save.md)）
釘死之後是可重現的。

## 5. 未解

| 項目 | 現況 |
|---|---|
| remake 要把游標的 hide／show 接在哪 | 規則已解（[`../re/88`](../re/88-mouse-cursor-visibility.md)），但 54 個呼叫點還沒對應到 remake 的繪圖流程 |
| remake 多的兩列 | 「主君編成」「損害報告」是 remake 加的，[`39`](39-system-window-parity.md) 已裁定保留 |
