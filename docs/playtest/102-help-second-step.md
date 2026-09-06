# 102 — 狀態列 #7：最後一條沒有原版擷取的，靠受控存檔派一個外交官走到

**狀態：通過。** 請求協助的**第二步**（「請選擇協同進攻之勢力。」）
五區逐像素比，只剩 **95 px 的原版游標**。
[`../spec/140`](../spec/140-status-message-box.md) §5 的
「只剩 #7 走不到」到此收掉——27 個帶索引的狀態列呼叫點**全部**有原版擷取。

- 日期：2026-09-07
- 規格：[`../spec/140`](../spec/140-status-message-box.md) §1.1（#8 → #7）、
  [`../spec/150`](../spec/150-diplomacy-preconditions.md)（兩道前置閘）、
  [`../spec/147`](../spec/147-controlled-parity-save.md) §3（`--diplomat`，本輪新增）
- 受控存檔：`tools/py.sh tools/parity_save.py
  workplace/dosgolem/root-noclouds/SAVE.DAT workplace/dosgolem/root-diplomat
  --diplomat 1:37`（武將 37 夏侯淵 → 勢力 1 孫策）
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-diplomat tools/dosgolem.sh …
  "wait;click:320,200;click:300,151,1000;runto:11CD0;sclick:352,15;stap:48,47;
  steps:300000;move:48,232;steps:100000;move:48,232;steps:100000;press;
  steps:900000;stap:200,112;steps:400000;stap:200,112;steps:900000;
  shot:orig-help7b"`
- remake 側：`-save-file workplace/dosgolem/root-diplomat/SAVE.DAT -load-slot 0
  -advise-target -advise-pick-row 2 -advise-list-row 0 -cam 0,0
  -fixture-when clock:196/4/17/6 -shot-when clock:196/4/17/6`

## 1. 結果

| 區 | 不同像素 |
|---|---:|
| `banner` | **0 / 20480** |
| `command` | **0 / 13824** |
| `map` | 95 / 145152 |
| `minimap` | **0 / 33280** |
| `faction` | **0 / 43264** |

那 95 px 落在 x 200–213、y 112–125 ＝ 最後一次點擊的位置，
14×14，原版自己畫的滑鼠游標——`ipeek:20100:1` 讀出來是 **1**
（[`../re/88`](../re/88-mouse-cursor-visibility.md)），與畫面一致。

**remake 一次就對**，沒有修任何東西。

## 2. ⭐ 走不到的狀態要用存檔做出來，不是等它發生

#7 掛了很久，理由是「要先派外交官到那個勢力，受控存檔還沒做那一版」。
`tools/parity_save.py` 加一個 `--diplomat 勢力:武將` 就走到了——
這是 [`../spec/147`](../spec/147-controlled-parity-save.md) 那句
「**存檔本身也可以是實驗器材**」的第二個例子（第一個是 `--no-clouds`）。

### 2.1 ⚠ 兩個欄位都要寫

同一件事在**兩張表各存一份**：勢力記錄 `+0x2A` ＝ 外交官的武將編號、
武將記錄 `+0x17` ＝ 職務 3（[`../spec/143`](../spec/143-general-duty-field.md) §2）。
只寫前者的話前置閘會過，但那個武將的「身分」欄還是「－－－」、
而且仍然出現在任命候選裡——**做出一份原版自己走不到的狀態**，
對拍就變成在比一個不存在的局面。

### 2.2 ⚠ 派的必須是自己的武將

第一次試跑派的是武將 1（袁尚，袁紹的人），閘照樣過、畫面也照畫，
**但原版的任命流程根本選不到他**（候選過濾第二條是「勢力 ＝ 玩家勢力」，
[`../spec/148`](../spec/148-shared-candidate-filter.md)）。
換成武將 37（夏侯淵）才是原版走得到的狀態。

⭐ **fixture 造得出來 ≠ 原版走得到。** 受控存檔繞過的應該是
「要花多久才會發生」，不是「規則允不允許」。

## 3. 順帶：兩段式在這裡也成立

第一次 `stap:200,112` 只把孫策那一列反白、狀態列仍是 #8；
第二次才決定並換成 #7。與 [`../spec/38`](../spec/38-list-windows.md) §1.7
的兩段式一致，remake 的 `-advise-list-row` 走的是同一條路。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 選完協同進攻對象之後 | 這一份停在 #7 的畫面，**再選下去**（成案／被拒）還沒拍 |
