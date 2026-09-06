# 106 — 啟動殼層的前三頁：背景 0 px，三個版面都是 remake 自己設計的

**狀態：進行中（已歸因）。** `-shot` 本身會觸發直啟，所以啟動殼層的畫面
**一直拍不到**；本輪加了 `-open-launcher` 之後第一次比。
**背景地圖與橫幅 0 px**，三個視窗的**版面各自不同**——那不是幾個常數的偏差，
是 remake 的殼層從一開始就照自己的設計走。

- 日期：2026-09-07
- 規格：[`../spec/90`](../spec/90-same-state-parity.md) §5.1（`-open-launcher`，本輪新增）、
  [`../spec/79`](../spec/79-new-game-faction-list.md)（勢力清單）、
  [`../spec/25`](../spec/25-slot-select-window.md)（四槽視窗）
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh
  workplace/parity/naming "wait;shot:s0;click:320,200;steps:600000;shot:s2"`
  （`s0` ＝ NEW GAME ＹＥＳ／ＮＯ、`s2` ＝ LOAD DATA）
  ＋ `"wait;click:320,176;steps:600000;shot:n1"`（劇本選擇）
- remake 側：`tools/parity_shot.sh out.png
  -save-file workplace/dosgolem/root-noclouds/SAVE.DAT -open-launcher 階段`
  ⚠ 這一份走的是**啟動殼層，不載存檔**——`-save-file` 只是讓 LOAD DATA
  那一頁有槽位資料可讀。

## 1. 結果

| 頁 | `banner` | 視窗以外的地圖 |
|---|---:|---|
| `title` | **0 / 20480** | 背景相同（`map`／`minimap`／`faction` 的差異全在視窗矩形內）|
| `scenario` | **0 / 20480** | 同上 |
| `load` | **0 / 20480** | 同上 |

⭐ **背景那張地圖與橫幅逐像素相同**——殼層底下鋪的是同一張圖、同一個位置。

## 2. 三個版面都不一樣

| 頁 | 原版 | remake |
|---|---|---|
| `title` | **ＮＥＷ　ＧＡＭＥ ＹＥＳ／ＮＯ 兩項**，框 (208,128,224,88) | 三項選單 `NEW GAME`／`LOAD DATA`／`LANGUAGE`，框 (112,56,416,288) |
| `scenario` | 標題 `NEW GAME` ＋ **四章各兩列**（「第一章．「呂布歸天」之卷」／「196年 4月 1日」）| 標題「選擇劇本」＋ 一列摘要（「第一章勢力：曹操　軍師：荀彧」）＋ 三條虛線 |
| `load` | `sub_18B5D` 的四槽視窗：副標「第一章勢力：曹操　軍師：荀彧」＋ 四列日期 | 自己的一份：「第 1 槽　196年4月16日　曹操」／「第 2 槽　空白槽位」 |

⚠ **`load` 這一頁與遊戲中的四槽視窗在原版是同一支常式**
（`sub_18B5D` → `sub_18B7C`，只有標題不同），而
[`../playtest/99`](../playtest/99-slot-window-parity.md) 已經把遊戲中那一份
對到 0 px——包括「空槽照畫日期，不寫『空白槽位』」。
**remake 的殼層是第二份實作**，所以那一輪的修正沒有跟過來。
這違反 `CLAUDE.md` §7 第 6 條「一條規則只留一份實作」。

## 3. 這不是「差幾個常數」

三頁的框矩形、欄位數、每一列的內容都不同，
**逐像素的數字（`map` 七萬多）在這個階段沒有資訊量**。
要對得先有一份「原版殼層版面」的規格（框、列距、每一列印什麼），
再讓 remake 照著重畫——那是一份新的 `docs/spec/`，不是這一輪的修補。

⭐ 但 §1 已經有一個實質結論：**背景不必動**。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 原版殼層三頁的版面 | 框矩形與欄位都還沒從機器碼讀出來。入口：`sub_11AC3`（新遊戲流程）、`sub_18B5D`（四槽視窗）|
| `title` 的三項選單 | remake 多一項 `LANGUAGE`（[`../spec/86`](../spec/86-runtime-language-switch.md) §4 的 remake 差異）。原版只有 ＹＥＳ／ＮＯ——**要不要改成兩項＋另找語言入口還沒裁定** |
| `load` 的兩份實作 | 殼層那一份與 `drawSaveUI` 應該收成一支（§2）|
