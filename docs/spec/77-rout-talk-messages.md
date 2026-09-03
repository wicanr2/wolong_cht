# 77 — 敗走的兩段訊息：TALK #1F 與 #23 ＋ 八變體

**狀態：CONFORMED。** 敗走當下與 48 tick 後各一段，
索引與觸發條件都在機器碼裡，remake 已接上並有單測。

- 日期：2026-08-23
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_12977`（`00012977`）
  與 `sub_12A7E`（`00012A7E`），呼叫鏈見
  [`../re/65`](../re/65-ai-march-decision-chain.md) §8
- 推論等級：**confirmed**（兩支都是 IDA 直接反組譯出來的，沒有手工解碼）

## 1. 原版做什麼

[`43`](43-rout-on-blocked-return.md) 已經解出敗走的狀態機。這一份補的是它的**兩段訊息**。

### 1.1 敗走當下（`sub_12977`）

```asm
mov cx, 1Fh
mov al, cs:byte_10CFF        ; 玩家勢力
cmp al, [si+1]               ; 這支軍團的勢力
jz  short loc_129AE          ; → 是玩家的軍團，用 #1F
cmp al, cs:byte_12919        ; 呼叫 sub_1291A 時傳進來的那個勢力
jnz short loc_129BE          ; → 都不是，不出訊息
inc cx                       ; → 用 #20

loc_129AE:
mov al, [si+2]               ; 軍團編號 ＝ 主將的武將編號
mov ah, 0FFh                 ; 第二個變數留白
push ax / mov di, sp
mov al, 93h                  ; 一般通知的肖像
call sub_18810
```

| 索引 | 原文 |
|---|---|
| `#1F` | 「{1}大人的部隊遭殲滅！幸好已逃過敵軍之手。」|
| `#20` | 「很遺憾，沒能將{1}捉拿到手。」|

⭐ **敗走這條路徑走不到 `#20`。** `byte_12919` 是 `sub_1291A` 開頭
從 `al` 存下來的，而敗走的呼叫端（`sub_147BB` 的 `0x8000` 分支）
傳的 `al` 就是 `[si+1]`——**軍團自己的勢力**。
所以第二個 `cmp` 比的是同一個值，第一個沒中的話第二個也不會中。
`#20` 只在**戰鬥脫身**那條路徑出現，那時 `al` 是勝方
（[`../re/09`](../re/09-combat.md) §5）。

### 1.2 48 tick 之後（`sub_12A7E`）

```asm
mov al, cs:byte_10CFF
cmp al, [si+2241h]           ; 軍團的勢力
jnz short locret_12AD1       ; → 不是玩家的，不出訊息
add di, 4240h                ; di ＝ 主將的武將記錄
push di / mov di, sp / call sub_10CDE
mov cx, 23h  / mov al, 93h            / call sub_18810
pop bx
mov cx, 198h / mov ah, [bx+1Eh]       ; 主將的 +0x1E ＝ 變體
             / mov al, [bx+1]         ; 主將的勢力
             / call sub_18810
```

| 索引 | 原文 |
|---|---|
| `#23` | 「{1}大人平安歸來了。」|
| 組 `0x198` | 主將自己的檢討句，八格一組（422–429）|

**組 `0x198` 展開**（`talkVariantGroupBase = 0x196`，公式見
[`../re/25`](../re/25-message-variants-and-personnel.md) §1）：

| 變體 | 索引 | 原文 |
|---:|---:|---|
| 0 | 422 | 「．．．．」|
| 1 | 423 | 「我絕不會再受此等侮辱！！」|
| 2 | 424 | 「都是我的不德所致，請各位見諒。」|
| 3 | 425 | 「實，實在太丟臉了．．」|
| 4 | 426 | 「哦哦！！下次我一定會打回來的。」|
| 5 | 427 | 「厚顏的回來．．」|
| 6 | 428 | 「我把主公借給我的兵馬都損失了。」|
| 7 | 429 | 「若能再給我機會，下次一定．．」|

⚠ 展開用**原始的 `+0x1E`**，不是收斂成 0–2 的那個
（同 [`48`](48-governor-returns-on-city-fall.md) §2 的理由）。

⭐ 這一組與內政官歸來那一組（`0x1A6` → 534–541）是**兩組不同的句子**，
雖然形狀一模一樣。`0x1A6` 本身又剛好等於十進位 422，
**把組編號當索引用會落到 `0x198` 的第 0 格**——那是同一個坑。

## 2. 演算法

```
敗走成立時：
    if 軍團.勢力 == 玩家勢力:
        訊息 #1F，{1} ← 主將名，一般通知肖像

倒數歸零時（軍團記錄清掉、主將解職的那一刻）：
    if 軍團.勢力 == 玩家勢力:
        訊息 #23，{1} ← 主將名，一般通知肖像
        訊息 0x196 + ((0x198 − 0x196) << 3) + 主將.+0x1E，主將的肖像
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/corps.go`：`CorpsEvent.RoutEnded`；`tickRout` 改成回傳「這一 tick 剛結束」|
| 呈現層 | `cmd/wlgame/corps.go`：`reportRout`，掛在 `reportCorps` 的迴圈裡（與 `reportGovernorReturn` 同一個位置）|
| 常數 | `routTalk = 0x1F`、`routReturnTalk = 0x23`、`routRegretTalkBase = 0x198` |
| 差異 | 無 |

`#20` 在**敗走**這條路徑沒有實作——照 §1.1，它到不了。
戰鬥脫身那條路走 [`123`](123-captive-talk-messages.md)，那裡 `#20` 有接。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestRoutEndSignalsOnlyOnce`（`internal/state`）：倒數歸零那一 tick 才回 true |
| 單元測試 | `TestRoutTalkOnlyForPlayer`（`cmd/wlgame`）：別人的軍團敗走不出訊息 |
| 單元測試 | `TestRoutEndEnqueuesReturnAndRegret`（`cmd/wlgame`）：兩則、肖像分別是一般通知與主將 |
| 單元測試 | `TestRoutRegretTalkIndex`（`cmd/wlgame`）：`+0x1E` 選到 422–429 那一組 |
| 對原版 | **未做**——要讓原版跑出一支回不了家的軍團，得先有對應的存檔 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~`sub_10CDE` 做什麼~~ | **PC 喇叭發聲**（已解，[`../re/37`](../re/37-graphics-and-runtime-module-map.md) §1、[`../re/17`](../re/17-dosv-audio-tsr.md) §5）：整支只有 `mov ax, 101h` ＋ `call sub_1EB11`，而 `sub_1EB11` 直接操作 PPI port `61h`。**它跟訊息內容無關**，是訊息框跳出來時的提示音 |
| 對原版的實跑驗證 | §4 仍是**未做**：要讓原版跑出一支回不了家的軍團，得先有對應的存檔 |
| ~~戰鬥脫身的 `#1F`／`#20`~~ | **已接**，連被擒（`#21`／`#22`）與自刎（`#43`）一起——見 [`123`](123-captive-talk-messages.md) |
