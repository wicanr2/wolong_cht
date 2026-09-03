# 123 — 武將下場的五則訊息：脫身、被擒、自刎

**狀態：CONFORMED。** `sub_1291A` 擲完下場之後會跳訊息框，
而**敗方與勝方看到的是不同的一則**。remake 先前只有自己畫的一行狀態列
（「○○ 部隊壞滅（脫身）」），原文一則都沒接。

- 日期：2026-09-03
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_12977`（`00012977`）
  與 `sub_129C3`（`000129C3`），分派在 `sub_1291A`
  （[`../re/09`](../re/09-combat.md) §6）
- 推論等級：**confirmed**（IDA 直接反組譯）
- 相關：[`77`](77-rout-talk-messages.md)（敗走那兩則，同一支 `sub_12977`）

## 1. 原版做什麼

`sub_1291A` 的輸入 `al` ＝**勝方勢力**，開頭存進 `cs:byte_12919`。
擲完下場之後分兩支：脫身走 `sub_12977`、被擒走 `sub_129C3`。

**兩支的訊息都靠「玩家是哪一邊」選。** 玩家兩邊都不是就**不出訊息**——
二十二個勢力天天在打，全部都跳框的話畫面停不下來。

### 1.1 脫身（`sub_12977`）

```asm
mov cx, 1Fh
mov al, cs:byte_10CFF        ; 玩家勢力
cmp al, [si+1]  / jz  .出     ; ＝ 敗方 → #1F
cmp al, cs:byte_12919
jnz .不出                     ; ≠ 勝方 → 什麼都不出
inc cx                        ; ＝ 勝方 → #20
.出: mov al, [si+2] / mov ah, 0FFh / mov al, 93h / call sub_18810
```

### 1.2 被擒（`sub_129C3` 的 `loc_12A12` 起）

```asm
xor dx, dx
mov al, cs:byte_10CFF
cmp al, [bx+1Dh] / jz  .出    ; ★ 比的是**舊主**（+0x1D，剛剛才寫進去的）
cmp al, cs:byte_12919
jnz .不出
inc dx
.出: mov al, 93h / mov cx, 21h / add cx, dx / call sub_18810
     and dx, dx / jz .不出     ; ★ 只有勝方看得到第二則
     mov ah, [bx+1Eh] / mov al, [bx+1] / mov cx, 19Ah / call sub_18810
```

⭐ **被擒那一則比的是 `+0x1D`（舊主）不是 `+0x1C`（現主）。**
`xchg al,[bx+1Ch]` 已經把現主換成勝方了，拿現主去比會讓敗方看到勝方那一句。

### 1.3 自刎（`loc_12A57`）

```asm
mov byte ptr [bx], 0 / mov word ptr [bx+1Ch], 0FFFFh   ; 變成在野
mov al, cs:byte_10CFF
cmp al, cs:byte_12919 / jnz .不出   ; ★ 只有勝方
mov al, 93h / mov cx, 43h / call sub_18810
```

⭐ **自刎沒有敗方那一則**——自刎的前提就是舊主已滅
（[`../re/09`](../re/09-combat.md) §6），沒有人在敗方那一側看畫面。

### 1.4 五則原文

| 索引 | 誰看得到 | 原文（校訂後的母本）|
|---:|---|---|
| `#1F` 31 | 敗方 | 「{1}大人的部隊遭殲滅！幸好以逃過敵軍之手。」|
| `#20` 32 | 勝方 | 「很遺憾，沒能將{1}捉拿到手。」|
| `#21` 33 | 敗方 | 「{1}大人的部隊遭殲滅！很遺憾，似乎遭敵軍所擒了。」|
| `#22` 34 | 勝方 | 「抓到{1}了！！」|
| `#43` 67 | 勝方 | 「{1}在即將被我軍擒拿之前，自刎而死了。」|
| 組 `0x19A` 438–445 | 勝方（接在 `#22` 後）| 被擒武將自己的一句 |

⚠ **32 與 34 的松崗原文句首各有兩個殘字「配對」**，日文原句沒有對應成分；
`translations/corrections.json` 早就收了這兩筆（`text-error`），
上表引的是**校訂後**的母本，也就是 remake 實際會顯示的字。

**組 `0x19A` 展開**（`talkVariantGroupBase = 0x196`）：438 ＋ 變體。

| 變體 | 索引 | 原文 |
|---:|---:|---|
| 0 | 438 | 「嗚．．不甘心．」|
| 1 | 439 | 「快殺了我吧！！」|
| 2 | 440 | 「。事到如今，不能屈辱人下．．」|
| 3 | 441 | 「畢生的失策．．竟然被你這傢伙所擒．．」|
| 4 | 442 | 「放過我的話，日後你會知道厲害！！」|
| 5 | 443 | 「竟然會活著受辱．．」|
| 6 | 444 | 「嗚．．一切都完了嗎．．．」|
| 7 | 445 | 「竟然會活著受辱．．」|

⭐ 這一組與敗走檢討那一組（`0x198` → 422–429，[`77`](77-rout-talk-messages.md)）
形狀一樣、句子不同。**組編號差 2 就是 16 則**，弄錯會拿到別人的台詞。

## 2. 演算法

```
下場擲完之後（勝方 W、敗方 L、玩家 P）：
    脫身:  P == L → #1F｜P == W → #20｜否則不出
    被擒:  P == 舊主 → #21｜P == W → #22 ＋ 0x19A 組（被擒者自己的肖像）｜否則不出
    自刎:  P == W → #43｜否則不出
    三者的 {1} 都是那名武將的**呼び名**（docs/re/79），肖像 0x93（一般通知）
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/corps.go`：`CorpsEvent.FateSides`（`corps → {Winner, Loser}`），在 `corpsPerishes` 填 |
| 呈現層 | `cmd/wlgame/corps.go`：`reportCaptives`，掛在 `reportCorps` 的迴圈裡 |
| 常數 | `escapeTalk = 0x1F`、`escapeMissedTalk = 0x20`、`capturedTalk = 0x21`、`captureTakenTalk = 0x22`、`suicideTalk = 0x43`、`captiveRegretTalkBase = 0x19A` |
| 差異 | 狀態列那一行（`battleLine` 的「部隊壞滅（脫身）／被擒／自刎」）是 **remake 才有的摘要**，保留 |

⭐ **`FateSides` 是必要的，不能事後從武將記錄反推**：被擒之後
`Generals[i].Faction` 已經是勝方，而訊息要比的是**舊主**。
原版靠 `+0x1D` 保存舊主，remake 靠事件把當下的兩個勢力帶出來。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestCorpsPerishesRecordsFateSides`（`internal/state`）：勝敗兩方都記進事件，被擒之後仍拿得到舊主 |
| 單元測試 | `TestCaptiveTalkPicksViewerSide`（`cmd/wlgame`）：三種下場 × 敗方／勝方／局外人 ＝ 九種組合，索引與「不出訊息」都對 |
| 單元測試 | `TestCaptiveRegretTalkIndex`（`cmd/wlgame`）：`0x19A` 展開成 438–445，**不是** `0x198` 那一組 |
| 對原版 | **未做**——要跑到一支軍團壞滅且玩家在場，得先有對應的存檔 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~`sub_129C3` 的 `test [bx], 40h` 分支~~ | ✅ **2026-09-03 解掉並接上**：bit 6 ＝ **主公型**（`../re/77` §3），被俘就清掉它並把說話類型 `+3`，把 0／1／2 搬到 3／4／5 ＝ 臣下型。⚠ 順帶訂正：這一段**不以「舊主已滅」為條件**（[`127`](127-captured-sovereign-becomes-retainer.md) §1）|
| 城兵那一側 | `sub_14FCE` 也呼叫 `sub_129C3`（守城武將被擒）。remake 的城兵路徑有沒有走到同一則沒驗 |
