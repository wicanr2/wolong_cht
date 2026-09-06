# 140 — 左下角的狀態列提示框

**狀態：CONFORMED。** 原版每個指令流程進行中都在畫面**左下角**掛一個
`(0, 320, 256, 80)` 的訊息框，肖像固定是通報者那張臉。remake 先前沒有這個框，
提示訊息不是併進別的視窗當標題，就是靠 remake 自己加的鍵盤提示條。

- 日期：2026-09-06
- 出處：[`../re/66`](../re/66-message-box-geometry.md) §1.2
  （`sub_18853` ＝ `00018853`、`sub_189A4` ＝ `000189A4`）、
  [`../re/22`](../re/22-strategy-command-tree.md) §4
- 推論等級：**confirmed**（機器碼的立即值 ＋ 原版擷取逐格量測兩條獨立證據）
- 相關：[`41`](41-message-box-geometry.md)（同一個框的另一個位置）、
  [`106`](106-message-box-reporter-portrait.md)（肖像 `0x93`）、
  [`39`](39-march-order-menu.md)（第一個用它的流程）

## 1. 原版做什麼

```asm
sub_18853(cx = TALK 索引)
    al = 93h / dx = 0 / bx = 12h
    push ax                       ; 肖像頁先存起來
    al = (cx == 0FFFFh) ? 0 : 1Eh ; 給 sub_189A4 的樣式：0 ＝ 擦除
    cx = 510h / call sub_189A4    ; 畫框
    pop ax                        ; al 還原成 93h
    cx == 0FFFFh → 收尾            ; 清除到此為止，不畫內容
    dx = 0 / bx = 14h / call sub_1075B  ; 肖像 ＋ 文字
```

**50 個呼叫點**全部落在指令 handler 內：**進入時設提示、離開時以 `0FFFFh`
清掉**，成對出現——23 處是清除，27 處帶索引，索引落在 `#0`–`#25`。

### 1.1 27 個帶索引的呼叫點

| TALK | 函式 | 什麼時候 |
|---:|---|---|
| #0 | `sub_16C5E` | 編成：選武將 |
| #1 | `sub_16C92` | 編成：下各部隊的編成指示 |
| #2 | `sub_1628F` | 軍團 → 行軍指示：選軍團 |
| #3 | `sub_17FDB` | 行軍：指示目標據點 |
| #4 | `sub_17F90` | 軍團情報面板——⭐ **別人的軍團**才掛（自己的直接進行軍指示，狀態列是 #3；[`149`](149-march-target-map-picker.md) §1.3）|
| #4 | `sub_15AD1` | 22 勢力的選擇視窗（縮小地圖圖例右半格，[`../re/31`](../re/31-faction-picker-screen.md)）|
| #5 | `sub_16405` | 進言 → 敵對提案 |
| #6 | `sub_164F1` | 進言 → 停戰提案 |
| #8 → #7 | `sub_16623` | 進言 → 請求協助：**先選協助勢力，再選協同進攻的對象** |
| #9 | `sub_16A9B`／`sub_16B71` | 人事：任命時選武將 |
| #11／#13 | `sub_16A9B`／`sub_16B08` | 內政官任命／解任 |
| #12／#14 | `sub_16B71`／`sub_16BE3` | 外交官任命／解任 |
| #15 | `sub_16909` | 進言 → 遷都：選目標據點 |
| #16 | `sub_1678D` | 財政主畫面 |
| #17／#18／#19／#20 | `sub_167CD`／`sub_167E6`／`sub_16806`／`sub_16826` | 稅率／騎兵／弓兵／步兵的募集人數 |
| #21 | `sub_17FDB` | 行軍三選一（帶 `\2` ＝ 目標據點名）|
| #22 | `sub_1628F` | 軍團 → 位置確認 |
| #23 | `sub_162FB` | 據點 → 據點一覽 |
| #24 | `sub_16366` | 武將一覽（[`145`](145-general-and-faction-cells.md)）|
| #25 | `sub_163BF` | 勢力一覽（同上）|

⭐ **`#10`「請選擇解任之武將。」一個呼叫點都沒有**——解任是選據點／勢力，
不選武將（[`142`](142-personnel-dismiss-flow.md)），所以那一則在原版是死文字。

## 2. 演算法

```
框：X ＝ dx × 16       ＝  0 × 16 ＝   0
    Y ＝ (bx + 2) × 16 ＝ 20 × 16 ＝ 320      ; sub_189A4 自己加 2 格
    寬 ＝ cl × 16      ＝ 16 × 16 ＝ 256
    高 ＝ ch × 16      ＝  5 × 16 ＝  80
內容：肖像頁 ＝ 0x93（KAOGRF 第 147 張，通報者）
      版面與一般訊息框完全一樣（docs/spec/41），只有位置不同
清除：cx ＝ 0FFFFh
```

### 2.1 代入的名字：定寬三格、各有各的顏色

七支 marker handler 的表在 [`../formats/01`](../formats/01-talk-dat.md) §3
與 [`../re/79`](../re/79-talk-marker-handlers.md) §2。這一則（TALK #21）用到 `\2`：

| 標記 | 代入 | 顏色 |
|---|---|---:|
| `\1`／`\4` | 武將呼び名 | 9 |
| **`\2`** | **據點名** | **`0x0B`** |
| `\3`／`\5` | 君主姓名 | `0x0C` |

⭐ **五支都是 `al = 3` 再 `call loc_10701`**——也就是**照欄位寬度畫滿三個
全形字，補白也畫出來**。原版擷取上量到的正是這兩件事：
「許」「昌」兩格是色 `0x0B`（`(243,162,0)`），第三格是空白，
接下來的「移動下。」才回到色 15。

⚠ **補白要在語系轉換之後補**：`uitext.Convert` 的覆寫詞表是整串比對的，
帶著補白會查不到。

⭐ **這與一般訊息框是同一個框**——`sub_1895D`／`sub_189A4` 都吃 16 px 粗格、
都自己 `+2` 格，寬高也都是寫死的 `cx = 510h`。**變的只有位置**：
一般訊息框 `(160, 160)`、狀態列提示 `(0, 320)`。

⭐ 它**不是模態的**：框掛著的同時，指令流程繼續（選據點、跳選單）。
原版的模態訊息走 `sub_18810`（[`41`](41-message-box-geometry.md)），兩者的呼叫端不同。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 版面常數 ＋ 繪製 | `cmd/wlgame/statusbox.go`：`statusBoxX`／`statusBoxY`、`drawStatusBox` |
| 狀態 | `game.statusBox`（`statusBoxState`；`lines` 是空的就是沒掛——**零值安全**）|
| 設／清 | `setStatusTalk(index, vars)`／`clearStatusTalk()`，對應 `cx = 索引`／`cx = 0FFFFh` |
| 框本體 | 沿用 `drawLegacyTalkBox`（與一般訊息框同一支），肖像 `defaultPortraitPage` ＝ `0x93` |
| 名字換色 | `talkInkCity` ＝ `0x0B`，用既有的 `drawTalkLineWithName`（戰場對白的 `\1` 走同一支，色 9）|
| 定寬補白 | `padTalkField`／`talkFieldCells` ＝ 3：`setStatusTalk` 對 `\1`–`\5` 一律補到三個全形字 |
| 財政（`cx = 10h`）| `beginFinance`／`endFinance`：TALK **#16**「請指示下個月以後的財政予定。」——`sub_1678D` 進入時 `mov cx, 10h`、離開時 `mov cx, 0FFFFh`，**成對出現**（[`../re/22`](../re/22-strategy-command-tree.md) §3.4）|
| 編成（`cx = 0`）| `beginForm`／收尾：TALK **#0**「進行軍隊編組。請選擇武將。」——`sub_16288` 只是 `mov cx, 1 / call sub_16C5E`，狀態列在 `sub_16C5E` 裡設（[`../re/30`](../re/30-corps-formation-ui.md) §1）|
| 進言四條 | `setAdviseStatus()`（`cmd/wlgame/advise.go`）：敵對 #5、停戰 #6、**請求協助是兩步兩則**（先 #8 選協助勢力、再 #7 選協同進攻的對象，取消退回上一步就換回去）、遷都 #15 |
| 財政四格 | `financeAmountTalk`（`cmd/wlgame/finance.go`）：稅率 #17／騎兵 #18／弓兵 #19／步兵 #20；離開數值輸入器就換回 #16 |
| 編成第二層 | `formOrderTalk` ＝ #1，選完武將那一刻換上去（`sub_16C92`）|
| 行軍目標 | `marchTargetTalk` ＝ #3，`pickDestination` 掛（`sub_17FDB`）|
| 軍團情報 | `corpsInfoTalk` ＝ #4，`openCorpsInfo` 掛（`sub_17F90` 的 `loc_17FBC`）。**自己的軍團走 `showCorpsPanel` 不掛 #4**（[`149`](149-march-target-map-picker.md) §1.3）|
| 勢力選擇視窗 | `factionPickerTalk` ＝ #4「以滑鼠的右鍵回復。」，`openFactionPicker`／`closeFactionPicker` 成對（`sub_15AD1` 進 `sub_15AFC` 之前掛、離開後清；[`../re/31`](../re/31-faction-picker-screen.md) §1.3）|
| 差異 | **27 個帶索引的呼叫點全部接上了。**〔#10 除外——原版一個呼叫點都沒有，是死文字（§1.1）〕|
| 差異 | 名字的換色與定寬補白目前只在這個框成立。**一般訊息框（`drawMessage`）還是整段畫成白的、名字也還是裁掉補白**——同樣是既有差異 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestStatusBoxRectMatchesOriginal`（`cmd/wlgame`）：框 ＝ `(0, 320, 256, 80)`，四個數字各自對回 `dx`／`bx`／`cl`／`ch` 的立即值 |
| 單元測試 | `TestStatusTalkClearsLikeFFFF`：`clearStatusTalk()` 之後不畫 |
| 單元測試 | `TestPadTalkFieldKeepsFieldWidth`、`TestStatusTalkPadsAfterSubstitution` |
| 對原版 ✅ | [`../playtest/79`](../playtest/79-march-menu-original-layout.md)：整個框 **256×80 逐像素 0 px**（框、肖像、四列字、名字的色 `0x0B` 與補白全部對上）|
| 對原版 ✅ | [`../playtest/81`](../playtest/81-command-cell-highlight.md)：財政（#16）與編成（#0）兩條流程的狀態列框也各自對過 |
| 單元測試 | `TestStatusTalkIndexesMatchOriginal`：**25 個索引逐格釘住**（`cmd/wlgame`）|
| 單元測試 | `TestAdviseCooperateSwapsStatusBetweenSteps`：**真的走流程再看框裡的字**——請求協助 #8 → #7 → 退回 #8，敵對 #5、停戰 #6。⭐ 只驗常數擋不住「常數對、流程沒接」，先前那三條就是這樣 |
| 對原版 ✅ | 財政的稅率輸入（#17）：`banner`／`command`／`minimap`／`faction` 全 0 px，`map` 95 px（原版自己的滑鼠游標 88 px）|
| 對原版 ✅ | 行軍目標（#3）：原版擷取確認選完軍團之後掛的是「請指示行軍目標之據點。」；⚠ **remake 用一覽表選目的地、原版在地圖上點**，所以那一張不做逐像素比 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~另外 45 個呼叫點~~ | **接完了**：27 個帶索引的呼叫點全部有對應（§3）。⭐ `#10`「請選擇解任之武將。」在原版是**死文字**，一個呼叫點都沒有 |
| ~~逐條的原版擷取~~ | **27 個帶索引的呼叫點全部有原版擷取了**，逐條的結果收在 [`../playtest/93`](../playtest/93-status-bar-sweep.md)。⚠ 只剩 #7（請求協助的第二步）走不到——要先派外交官到那個勢力（[`150`](150-diplomacy-preconditions.md)），受控存檔還沒做那一版 |
| 樣式 `1Eh` 是什麼 | `sub_189A4` 把它傳給 `sub_189DE` 當 `ah`。**只知道 0 ＝ 擦除、非 0 ＝ 畫**，`1Eh` 這個值本身沒解（一般訊息框傳的也是 `1Eh`）|
| ~~一般訊息框要不要一起改~~ | **改了**（[`119`](119-talk-marker-fields.md) §3.1）：`\1`–`\5` 的定寬補白與逐標記換色收成一支，兩個框共用——原版本來就是同一支 `sub_1075B`。[`../playtest/92`](../playtest/92-personnel-assign-parity.md) §3 量到那 150 px |

<!-- 缺口：四項，見上表 -->
