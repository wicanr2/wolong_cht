# 135 — 腳本指令 16 是「掛一個對白框」：攻城戰的開場勸降

**狀態：CONFORMED。** `opMessage` 在 remake 是空的 stub（只寫進 `Battle.Log`），
於是**攻城戰開場那兩個對白框一個都不會出現**——同一場的 `field` 對拍
22.30% 全部來自它們。原版是腳本指令 16（`sub_1A69F`），
TALK 組 ＝ `0x1CE + 運算元`。

- 日期：2026-09-05
- 出處：`KI.EXE`（松崗 DOS/V）`sub_1A69F`（`0001A69F`，分派表
  `funcs_1A457` 的第 16 項）、`sub_1C315`（掛框，`docs/spec/60`）、
  `sub_1A426`（指令拆解）、`sub_1A8F6`（全軍退卻，寫 `byte_1D349`）
- 推論等級：confirmed（機器碼 ＋ 原版執行時的畫面）

## 1. 原版做什麼

指令拆解在 `sub_1A426`：

```asm
mov ax, [si]          ; ax ＝ 兩個 byte
mov bl, al / and bl, 1Fh      ; 低 5 位 ＝ 指令碼
and al, 0E0h / rol al,1 ×3    ; ★ al ＝ 參數（高 3 位）
                              ; ★ ah ＝ 第二個 byte（運算元）
```

指令 16（`sub_1A69F`）：

```asm
mov cl, ah / xor ch, ch
add cx, 1CEh            ; ★ TALK 組 ＝ 0x1CE + 運算元
mov dl, al / and dl, 1  ; ★ 側別 ＝ 參數 bit 0
mov ah, cs:byte_1D349
and al, 6
jnz  .rel
and  ah, ah / jz .show  ; 參數 bit1–2 都是 0 → 只在「戰鬥進行中」掛
jmp  .out
.rel:
and dl, dl / jz .cmp
xor al, 6               ; ★ 參數 bit 0 為 1 → 把「哪一側」的意思翻過來
.cmp:
shr al, 1
cmp ah, al / jnz .out
.show:
call sub_1C315          ; 掛框：cx ＝ TALK 組、dl ＝ 側別，活 60 拍
```

`byte_1D349` 是**戰鬥的收尾階段**（`sub_19A33` 開場歸零）：

| 值 | 意思 | 誰寫的 |
|---:|---|---|
| 0 | 戰鬥進行中 | 開場 |
| 1 | **側 0 全軍退卻中** | `sub_1A8F6`（指令 3 下命令 5、或 `sub_1AE56` 大將體力不支）|
| 2 | **側 1 全軍退卻中** | 同上 |

⇒ 指令 16 的參數其實是一道**「現在是哪個階段才講這句」**的閘：
開場的勸降走參數 bit1–2 ＝ 0（只在進行中講），
退卻時的台詞走另外兩個值。

## 2. 演算法

```
組 = 0x1CE + 運算元
側 = 參數 & 1
若 (參數 & 6) == 0：
    要 = 0
否則：
    v = 參數 & 6
    若 參數 & 1：v ^= 6        ← 只有這一支會翻，(參數 & 6)==0 那一支不翻
    要 = v >> 1
若 收尾階段 == 要 → 掛框（側, 組），活 60 拍
```

⭐ **32 段腳本全部以「訊息／等 15／訊息」三個 2-byte 指令開頭**
（`BATTLE.DAT` 實測，已記在 `Script.SkipBytes` 的註解）。
所以攻城戰開場會**先後掛兩個框**——第 50 拍一個、第 65 拍一個，
各活 60 拍，第 70 拍兩個都還在。

⚠ **野戰不會看到它們**：單挑挑戰成立時原版對腳本 PC `add 6`，
正好跳過這三個指令（[`80`](80-duel-opening.md) §3.2）。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/tactical/script.go` 的 `opMessage`；喊話佇列從 `duelState` 移到 `Battle`（腳本訊息與單挑喊話走同一條出口）|
| 收尾階段 | `internal/rules/tactical/battle.go` 新增 `Battle.endPhase`，在 `Order(side, −1, Retreat)` 設 1／2 |
| 呈現層 | 不動——`cmd/wlgame/battletalk.go` 的 `pumpDuelTalks` 已經在讀 `TakeDuelTalks()` |
| 差異 | 無 |

⚠ **不動勝負判定。** `sub_1A8F6` 還會起一個 120 拍的倒數
（`byte_1D34A = 0x78`）決定戰鬥何時結束，那是另一件事，
這一份只接「階段值」這個輸入。倒數本身列在 §5。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `internal/rules/tactical/script_test.go` 的 `TestScriptMessageEmitsATalk`（組 ＝ `0x1CE+運算元`、側 ＝ 參數 bit 0）與 `TestScriptMessageGateSkipsWhenPhaseDiffers` |
| 突變測試 | 把 `0x1CE` 改成別的值、把閘拿掉，兩支各自要紅 |
| 逐區對拍 | 同一場攻城戰第 70 拍，`field` **22.30% → 2.31%**（框畫出來了，內容還要 [`136`](136-battle-talk-parameters.md)），最後收在 0.08%（[`../playtest/76`](../playtest/76-battle-talk-parity.md)）|

## 5. 未解

| 項目 | 現況 |
|---|---|
| 退卻的 120 拍倒數（`byte_1D34A`）| `sub_1A8F6` 起、`sub_1A6FA` 遞減到 0 才 `sub_19FDC` 收尾。remake 的退卻沒有這段倒數 |
| 參數值 3 | `byte_1D349` 只會是 0／1／2，閘算得出 3 但沒有值對得上——沒有腳本用到，或是原版的死分支 |
