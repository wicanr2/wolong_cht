# 90 — 編成第二層：4,046 → 78 px，抓到反白列的字色錯了

**狀態：通過。** 編成選完武將之後那一張（狀態列 #1「請下達各部隊編成之
指示。」）逐區比，`banner`／`command`／`minimap`／`faction` 全 0 px，
`map` **78 px** ——而那 78 px 就是原版自己的滑鼠游標（14×14 的紅箭頭）。
路上抓到兩個真的缺陷：**選完之後那一列沒留著反白**，以及
**反白列的字色是黑的（原版是黃的）**。

- 日期：2026-09-06
- 規格：[`../spec/38`](../spec/38-list-windows.md) §1.7、
  [`../spec/140`](../spec/140-status-message-box.md) §1.1（#1）
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh …
  "wait;click:320,200;click:300,151,1000;runto:11CD0;sclick:352,15;stap:192,47;
  steps:400000;stap:200,112;steps:400000;stap:200,112;steps:900000;shot:orig-form1"`
- remake 側：`-save-file workplace/dosgolem/root-noclouds/SAVE.DAT -load-slot 0
  -open-form -form-pick-row 0 -lord-corps=false -cam 0,0
  -fixture-when clock:196/4/17/6 -shot-when clock:196/4/17/6`
  ⚠ **受控存檔不能省**（[`../spec/147`](../spec/147-controlled-parity-save.md)）：原版側跑的是 `root-noclouds`，remake 只用 `-direct` 從劇本跑到同一時刻會有雲、也會有軌跡分歧。這一行本輪補上，補之前照著跑對不出文中的數字

## 1. 一路修下來的數字

| 改了什麼 | `map` |
|---|---:|
| 第一次拍 | 8,894 |
| `-lord-corps=false`（原版不讓君主編成，[`../spec/76`](../spec/76-lord-not-in-formation.md)）| 4,046 |
| **選完之後那一列留著反白**（`listwin.KeepSelected`）| 614 |
| **反白列的字改成黃的**（`chrome.Highlight`）| **78** |

剩下的 78 px 落在 x 200–213、y 112–125 ——**正好是最後一次點擊的位置**，
14×14，原版自己畫的滑鼠游標。

## 2. 兩個缺陷

**一、選完之後那一列沒留著反白。** 原版把編成視窗畫在武將一覽上面，
而**被選中的那一列一直反白著**。remake 的 `listwin.Confirm()` 依兩段式的
規則會退回 `Browsing`，於是那一列變回一般列。修法是在「`listPick` 回 false
（清單留著）」時 `KeepSelected()` 擺回去——**這一條對所有「選完不關清單」
的流程都成立**：編成、人事四條、武將自陳。

**二、反白列的字色是黑的。** 原版是 `f3e300`（色 12，`chrome.Highlight`），
remake 整份清單一律用 `chrome.Ink`（黑）。

⭐ **這兩個先前都對拍不到**，因為比過的每一張清單都是「剛開窗、沒有反白列」
（[`42`](42-window-parity.md) §4）。`chrome.Select` 的註解上還寫著
「⚠ 這一個還沒有實機證據」——現在有了，而且順帶發現字色是另一個顏色。
**沒有樣本的地方不會報錯，只會一直錯著。**

## 3. 順帶確認的

- 狀態列 #1 逐像素對上（`sub_16C92` 的 `sub_18853(cx = 1)`）。
- 原版的編成候選**不含夏侯惇**（已經是軍團長）與**曹操**（君主），
  與 [`../spec/143`](../spec/143-general-duty-field.md)、
  [`../spec/76`](../spec/76-lord-not-in-formation.md) 都對得上。
- 預備兵數三格（0／2000／6000）與六個部隊各 1000 兩邊相同。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 原版的滑鼠游標 | 選單上是 14×14 紅箭頭（白邊）、大地圖上是 15×15 白色空心框（＋1,+1 黑影）。**遊戲自己畫的**（dosgolem 的 INT 33h 不畫），所以它是 remake 的缺口；但繪製端還沒定位，而且 remake 用的是 OS 游標，要接得連「隱藏 OS 游標」一起決定 |
| ~~反白列的 `listCellInk` 覆寫~~ | **有樣本了**（[`98`](98-enemy-corps-panel.md) §3）：儲存格自己的顏色贏過反白列的字色，而且**換一個色號**——一般列色 10、反白列色 6 |
