# 130 — 沒有心向的在野武將：隨機投靠，偏向武將最少的勢力

**狀態：CONFORMED。** `sub_15899` 的 `+0x19 == 0xFF` 分支——
[`114`](114-general-affinity.md) §5 掛著的那一條。規則有三層：
先找出**武將數最少**的勢力當起點，再擲一顆 1–64 的骰決定「往後數第幾個」，
而骰到 48 以上時走一條**玩家專屬的救濟**（玩家武將太少就送一個過去）。

⚠ **這一條不經過 25% 那道閘。** `cmp bh, 0FFh / jz loc_158C2` 在擲骰**之前**，
所以有心向的每月 25% 才動，沒有心向的**每月都跑**。
remake 現有的 `recruitFreelanceGenerals` 是先擲 25% 再分流，接的時候要拆開。

- 日期：2026-09-04
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_15899`
  （`000158C2`–`0001593F`），呼叫端 `sub_1585F`（月結逐武將）
- IDA database：`workplace/ida/dosv/KI.EXE.i64` SHA-256
  `65736f11b0b28a5b3a6db9e1a3d205cc24f0eaebc82b508ee0d7d283f6240572`（739 支函式）
- 推論等級：**confirmed（靜態，逐行讀出）**
- 相關：[`114`](114-general-affinity.md)（心向那一條）、
  [`../re/77`](../re/77-general-affinity-and-flags.md) §2、
  [`../mechanics/70`](../mechanics/70-ai.md) §3.9

## 1. 原版做什麼

```asm
loc_158C2:                        ; 心向 ＝ 0xFF 才走到這裡
        mov di, 0 / mov bx, di
        mov ax, 0FF16h            ; al ＝ 22（勢力數）、ah ＝ 0xFF（最小值的初值）
.scan:  cmp byte ptr [di], 80h    ; 勢力存在？（記錄 +0 的 bit 7）
        jb  .next
        cmp ah, [di+18h]          ; ★ +0x18 ＝ 武將數
        jb  .next
        mov ah, [di+18h]          ; 記下更小的
        mov bx, di                ; ★ bx ← 武將最少的那個勢力
.next:  add di, 40h / dec al / jnz .scan     ; 勢力記錄 64 B × 22

        call sub_1ECE0
        and al, 3Fh / inc al      ; ★ al ∈ [1, 64]
        cmp al, 30h / jnb .player ; ★ al ≥ 48 ⇒ 玩家救濟
        cmp al, 18h / jb .count
        mov al, 1                 ; ★ al ∈ [24, 47] ⇒ 當成 1
.count: cmp byte ptr [bx], 80h    ; 從 bx 起，數第 al 個「存在」的勢力
        jb  +
        dec al / jz .join         ; ⇒ 投靠 bx
+       add bx, 40h
        cmp bx, 580h / jb .count  ; 0x580 ＝ 22 × 64
        xor bx, bx / jmp .count   ; ⭐ 繞回勢力 0

.player:
        mov bx, cs:word_10CFD     ; ★ 玩家的勢力記錄
        mov al, [bx+23h]          ; ★ +0x23 ＝ 據點數
        shr al,1 / shr al,1 / inc al          ; 據點數 ÷ 4 ＋ 1
        cmp al, [bx+18h] / jbe .out           ; ≤ 武將數 ⇒ 什麼都不做
        ; 否則落到 .join，對象就是玩家

.join:  cmp bx, cs:word_10CFD / jnz +
          （對象是玩家 ⇒ 跳訊息 0x29）
+       shl bx,1 / shl bx,1
        mov [si+1Ch], bh          ; 武將的勢力 ← bx ÷ 64
        al = bh / ah = 0FFh / call sub_12AD2   ; 勢力的武將數 +1
.out:   retn
```

## 2. 三件值得單獨記的事

### 2.1 ⭐ 起點是「武將最少的勢力」，所以整條規則偏向弱勢方

`.count` 的迴圈**從 `bx` 開始**而不是從勢力 0——而 `bx` 是掃描階段選出來的
**武將數最少**的勢力。所以骰到 1（連同被壓成 1 的 24–47，合計 **37/64 ≈ 58%**）
就是投靠武將最少的那一個。

⚠ 平手時取**先掃到的**（`jb` 是嚴格小於，相等不更新），所以編號小的優先。

### 2.2 ⭐ `al ≥ 48` 是玩家專屬的救濟，而且門檻與據點數綁在一起

17/64 ≈ **26.6%** 的機率走這一條，它只看玩家：

```
玩家據點數 ÷ 4 ＋ 1 > 玩家武將數   ⇒ 送一個武將過來
否則                               ⇒ 這個月什麼都不做
```

⚠ **不是「投靠玩家」而是「玩家缺人才補」**——條件不成立時那一顆骰就浪費了，
不會轉去投靠別人。所以隨機投靠的實際發生率低於 100%。

### 2.3 骰面的分布

| `al` | 機率 | 行為 |
|---|---:|---|
| 1–23 | 23/64 ＝ 35.9% | 從武將最少的勢力起，往後數第 `al` 個存在的 |
| 24–47 | 24/64 ＝ 37.5% | **壓成 1** ⇒ 武將最少的那一個 |
| 48–64 | 17/64 ＝ 26.6% | 玩家救濟（可能什麼都不做）|

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 分流 | `internal/state/freelance.go`：`recruitFreelanceGenerals` 先看 `Affinity`，**有心向才擲 25%**；沒有心向的直接走 `randomJoin` |
| 掃最少 | 同檔 `fewestGeneralsFaction()`：只算 `Alive` 的勢力，平手取編號小的 |
| 骰與分流 | 同檔 `randomJoin()`，三段照 §2.3 |
| 玩家救濟 | 同檔，門檻 `Cities/4 + 1 > Generals` |
| 入帳 | 沿用 `raiseGeneralCount`（`sub_12AD2`）|
| 訊息 | ⚠ **沒做**——心向那一條（`recruitFreelanceGenerals` 的主路徑）同樣只改狀態不排事件（[`114`](114-general-affinity.md) §5）。兩條路一致，要補就一起補 |

⚠ **繞回的迴圈要有上限。** 原版 `.count` 在「一個存在的勢力都沒有」時
會無限繞（實務上不會發生——玩家自己一定在），remake 加一個
「掃滿 22 × 2 圈就放棄」的保險絲並標記為 remake 差異。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestRandomJoinPrefersFewestGenerals`：骰面 1 與 24–47 都投靠武將最少的勢力 |
| 單元測試 ✅ | `TestRandomJoinCountsFromFewest`：骰面 `n` 從最少那一個往後數第 `n` 個存在的勢力，會繞回 |
| 單元測試 ✅ | `TestRandomJoinPlayerReliefThreshold`：`al ≥ 48` 時，據點數 ÷ 4 ＋ 1 > 武將數才送人，否則不動 |
| 單元測試 ✅ | `TestRandomJoinSkipsTwentyFivePercentGate`：沒有心向的武將**每月都跑**，不受 25% 閘限制 |
| 突變測試 ✅ | 起點改成勢力 0（不找最少）⇒ 前兩支當場紅 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 平手時取編號小的 | 從 `jb`（嚴格小於）推出來的，**沒有實機驗過** |
| 訊息 `0x29` 的完整文字與觸發畫面 | 只知道索引；原版跳訊息時玩家看到什麼沒有對拍 |
| 訊息 `0x29` | 兩條路都沒排事件，見 §3 |
| 這一條的實際發生頻率 | 開局 81 名在野武將**全部有心向**，所以隨機投靠要等他們兌現完才輪得到（[`../mechanics/70`](../mechanics/70-ai.md) §3.9）。**長跑幾個月才會第一次觸發沒有量過** |
