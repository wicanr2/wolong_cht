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

48 個呼叫點全部落在指令 handler 內：**進入時設提示、離開時以 `0FFFFh` 清掉**，
成對出現。

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
| 差異 | 原版的 48 個呼叫點目前只接了行軍那一條（[`39`](39-march-order-menu.md)）。**其餘流程仍是 remake 自己的鍵盤提示條**，那是既有差異，不是這一份帶進來的 |
| 差異 | 名字的換色與定寬補白目前只在這個框成立。**一般訊息框（`drawMessage`）還是整段畫成白的、名字也還是裁掉補白**——同樣是既有差異 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestStatusBoxRectMatchesOriginal`（`cmd/wlgame`）：框 ＝ `(0, 320, 256, 80)`，四個數字各自對回 `dx`／`bx`／`cl`／`ch` 的立即值 |
| 單元測試 | `TestStatusTalkClearsLikeFFFF`：`clearStatusTalk()` 之後不畫 |
| 單元測試 | `TestPadTalkFieldKeepsFieldWidth`、`TestStatusTalkPadsAfterSubstitution` |
| 對原版 ✅ | [`../playtest/79`](../playtest/79-march-menu-original-layout.md)：整個框 **256×80 逐像素 0 px**（框、肖像、四列字、名字的色 `0x0B` 與補白全部對上）|

## 5. 未解

| 項目 | 現況 |
|---|---|
| 另外 47 個呼叫點 | 只接了行軍那一條。其餘 handler 的 TALK 索引在 [`../re/22`](../re/22-strategy-command-tree.md) §3 都有，但要一條一條接 |
| 樣式 `1Eh` 是什麼 | `sub_189A4` 把它傳給 `sub_189DE` 當 `ah`。**只知道 0 ＝ 擦除、非 0 ＝ 畫**，`1Eh` 這個值本身沒解（一般訊息框傳的也是 `1Eh`）|
| 一般訊息框要不要一起改 | `\1`–`\5` 的定寬補白與逐標記換色是**全域規則**（[`../re/79`](../re/79-talk-marker-handlers.md) §2），但一次改到所有訊息會動到四個語系的排版（[`87`](87-latin-screen-layout.md)、[`../playtest/32`](../playtest/32-talk-layout-fit.md)）。要另外開一份規格，先量再改 |

<!-- 缺口：兩項，見上表 -->
