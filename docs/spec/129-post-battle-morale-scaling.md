# 129 — 戰術戰鬥打完，兩側的士氣都要按兵力比縮

**狀態：CONFORMED。** 原版打完一場戰術戰鬥有兩段士氣處理：`sub_19EBD` 先把
**敗方**的士氣 clamp 成 99 或 0，`sub_19F58` 再對**兩側**做
`士氣 × 新兵力 ÷ 舊兵力`。remake 的 `ResolvePending` 兩行都把**分母寫成新值**，
所以勝方那一行是恆等式、敗方那一行恆等於約 100——**等於沒有縮**。

- 日期：2026-09-03
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_19EBD`（`00019EDD`–`00019EF3`
  與 `00019F0C`–`00019F22`）、`sub_19F58`（`00019F7D`–`00019F9B`）
- 推論等級：**confirmed（靜態）**
- 相關：[`../re/83`](../re/83-post-battle-troop-accounting.md) §3（同一支的前半＝兵力）、
  [`../re/09`](../re/09-combat.md) §4（**自動判定**那條線的同名規則，已實作）

## 1. 原版做什麼

### 1.1 `sub_19EBD`：只對敗方 clamp

```asm
; side 0 那一段（side 1 同形，只有 cmp 的值換成 1）
cmp cs:byte_1D349, 0 / jz .skip     ; byte_1D349 ＝ 勝方的 side ⇒ 自己贏就跳過
mov al, 63h                          ; 99
xchg al, es:[di+6]                   ; 士氣 ← 99；al ← 舊士氣
cmp al, 63h / ja .skip               ; 舊士氣 > 99 ⇒ 保持 99
mov byte ptr es:[di+6], 0            ; 否則歸零
.skip:
call sub_19F58
```

`es:[di+6]` ＝ 軍團記錄 `+0x06`（士氣，[`../re/08`](../re/08-hourly-update.md)）。

### 1.2 `sub_19F58` 尾段：兩側都乘兵力比

```asm
pop  di                        ; ★ di 回到軍團記錄基址（迴圈裡是 +0x28 起的六個槽）
mov  cx, dx                    ; dx ＝ 六槽新兵數的總和
xor  ax, ax
and  cx, cx / jz .out          ; 新兵力 0 ⇒ 士氣 0
xchg cx, es:[di+4]             ; ★ +0x04 ← 新兵力；cx ← 舊兵力
and  cx, cx / jz .out          ; 舊兵力 0 ⇒ 士氣 0
mov  al, es:[di+6]             ; 士氣
mul  word ptr es:[di+4]        ; × 新兵力
div  cx                        ; ÷ 舊兵力
.out:
mov  es:[di+6], al
```

⭐ **`pop di` 在迴圈之後**，所以這裡的 `[di+4]`／`[di+6]` 是軍團記錄的
`+0x04`（兵力）與 `+0x06`（士氣），不是六槽迴圈裡那個 `+0x28+4k` 的位址。
兩處都寫 `[di+6]` 而**指的是同一個欄位**——偏移基準相同。

### 1.3 合起來

| | 勝方 | 敗方 |
|---|---|---|
| 戰前士氣 ≥ 100 | `戰前 × 新 ÷ 舊` | `99 × 新 ÷ 舊` |
| 戰前士氣 < 100 | `戰前 × 新 ÷ 舊`（**不歸零**）| **0** |

⚠ **與自動判定那條線不一樣。** [`../re/09`](../re/09-combat.md) §4 的
`sub_1474A` 對**勝方也有** `cmp [si+6], 64h / jb ⇒ 歸零`，戰術這一條沒有。
兩條線的敗方 base 也差 1（自動判定 `64h` ＝ 100、戰術 `63h` ＝ 99）。
**同名規則在兩條路上不同，引用時要標明是哪一條。**

## 2. remake 現況：兩行的分母都是新值

```go
// internal/state/tactical.go ResolvePending 的 apply()
if (side == 0) == o.AttackerWins {
    c.Morale = c.Morale * men / maxInt(men, 1)      // ← men/men ＝ 1，恆等式
} else if c.Morale < army.RoutMoraleGate {
    c.Morale = 0
} else {
    c.Morale = 100 * total / maxInt(total+1, 1)     // ← total/(total+1) ≈ 1
}
```

⭐ 兩行是**同一個錯誤形狀**：分母該放舊兵力，放成了新兵力（或新兵力 +1）。
症狀是「看起來在算，實際什麼都沒做」——`go vet` 不會說話，
而且結果落在合理範圍內，所以測試也抓不到。

同一支函式裡的 `scale()` 反而是對的（`v * men / c.Men`），
**舊兵力就在手邊**，只是士氣那兩行沒有用它。

## 3. 要改什麼

| 項目 | 位置 | 改法 |
|---|---|---|
| 兩側都縮 | `internal/state/tactical.go` `apply()` | 先抓 `before := c.Men`，士氣算 `base * total / before` |
| 敗方的 base | 同上 | 戰前 ≥ 100 ⇒ `99`；< 100 ⇒ `0`（＝ 現有的 `RoutMoraleGate` 分支，只是 base 從 100 改成 99）|
| 勝方不歸零 | 同上 | 戰術這一條**沒有**「戰前 < 100 就歸零」，base 直接用戰前士氣 |
| 舊兵力為 0 | 同上 | 士氣 0（原版 `and cx,cx / jz`）|

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestPostBattleMoraleScalesWithMen`：勝方 `戰前 × 新 ÷ 舊`、敗方 `99 × 新 ÷ 舊` |
| 單元測試 ✅ | `TestPostBattleLoserBelowGateZeroes`：戰前 < 100 的敗方歸零 |
| 單元測試 ✅ | `TestPostBattleWinnerBelowGateSurvives`：戰前 < 100 的**勝方不歸零**（與自動判定的差別）|
| 單元測試 ✅ | `TestPostBattleMoraleZeroWhenNoMen`：新舊兵力任一為 0 ⇒ 士氣 0 |
| 單元測試 ✅ | **`TestResolvePendingScalesMorale`：接線本身**——上面四支直接呼叫 `postBattleMorale`，把 `apply()` 裡那一行拔掉照樣**全套綠**（突變測試發現）|
| 突變測試 ✅ | 拔掉 `apply()` 裡那一行 ⇒ 只有接線那一支紅（士氣停在 200，該是 100）|
| 對拍 ✅ | 野戰九區重跑，`field` 95 px／`sb-minimap` 40 px／其餘七區 0 px，與 [`../playtest/62`](../playtest/62-parity-retest-20260903.md) 同值。⚠ **這一條本來就不該動到對拍**——士氣在 `ResolvePending` 算，畫面已經結束了；跑它只是確認沒有意外的耦合 |

⭐ **抽成 `postBattleMorale` 是為了讓它可測**。原本那兩行在 `apply()` 這個
閉包裡，除了跑完整場戰鬥沒有別的辦法驗——而那正是它錯了很久沒被發現的原因。

## 5. 未解

| 項目 | 現況 |
|---|---|
| 為什麼敗方是 99 不是 100 | 自動判定用 `64h`、戰術用 `63h`。**兩個立即值都讀出來了**，但差 1 的理由沒有解釋——可能只是 `xchg` 那個寫法順手（先寫 99 再比 99）|
| 戰後士氣的實機對照 | 沒有。要打完一場戰術戰鬥再看軍團一覽的士氣欄 |
