# 195 — AI 階段 1 問的是**意圖**，不是行軍目標

**狀態：CONFORMED。** `sub_143AF`（AI Stage 1）的三個判斷都讀
`bx = ax`——而 `ax` 是 `sub_14325` 從 **`[si+20h]`（意圖）** 算出來的
據點記錄位址，不是 `+0x14`（行軍目標）。

```asm
sub_14325:
        mov ah, [si+20h] / xor al, al       ; ★ ax ＝ 意圖 × 256
        shr ax,1 / shr ax,1 / shr ax,1
        add ax, 840h                        ; ⇒ 意圖那座據點的記錄位址
        …
        call cs:funcs_1434F[bx]             ; handler 收到的 ax 就是它

sub_143AF:
        cmp word ptr [si+0Eh], 800h / jnb → Stage 0     ; 還在路段上
        cmp word ptr [si+4], 12Ch  / jbe → Stage 10     ; 兵力 ≤ 300
        mov bx, ax                          ; ★ bx ＝ **意圖**據點
        test byte ptr [bx], 40h  / jnz → Stage 0        ; 威脅有具體目標
        cmp  byte ptr [bx], 80h  / jb  → 出擊            ; 不受威脅
        cmp  byte ptr [di+18h], 2 / ja → 出擊            ; 這一格軍團數 > 2
        …留守／補兵
```

remake 的 `aiStage1` 三個條件都用 `c.TargetNode`。
兩者在「意圖 ＝ 目標」時相同，**不同時就分岔**——而 AI 挑到新目標之後
要兩拍才落實（[`170`](170-ai-target-takes-two-ticks.md)），中間那一拍
正是兩者不同的時候。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_143AF`（IDA 線性 000143AF–0001440E）、
  `sub_14325`（00014325）
- 推論等級：confirmed（`ax` 的來源逐條讀出；`di` 的來源由原版的入口
  暫存器實測定案，見 §2）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §48.3
- remake 實作：`internal/state/aimarch.go` 的 `aiStage1`
- 相關：[`194`](194-ai-holds-a-threatened-city-on-the-way.md)（同一場對拍的另一半）、
  [`170`](170-ai-target-takes-two-ticks.md)、[`../re/65`](../re/65-ai-march-decision-chain.md) §3.2

## 1. 對拍上長什麼樣

同局面軍團 88（勢力 0）在拍 7,091 打下據點 74，之後意圖是 74
（[`194`](194-ai-holds-a-threatened-city-on-the-way.md) 讓它留守時寫的），
而行軍目標仍是 56。到了 Stage 1：

| 讀哪一個 | 判斷 | 結果 |
|---|---|---|
| 原版：意圖 74（剛打下、受威脅、只有自己）| `test 40h` 不成立、`cmp 80h` 成立、軍團數 ≤ 2 | **留守**，Stage 不變 |
| remake：目標 56 | 不受威脅 ⇒ 出擊 | 取一次亂數、Stage → 2 → 挑新待辦 |

於是 remake 的軍團 88 多取走勢力 0 的 `+0x17`（失守據點 88），
原版留著——勢力表、據點表、軍團表跟著全差。

## 2. ⭐ 第三個條件的 `di` 是**上一支軍團留下的勢力記錄**

`sub_14325` 只設 `ax` 與 `bx`（跳表索引），**`di` 是呼叫鏈上留下來的**。
攔 `sub_143AF` 的入口暫存器（`WOLONG_DOSGOLEM_WATCH=143AF`）拍 9,364：

```
#314445725 呼叫 143AF AX=0940 BX=000A CX=0002 DX=0004 SI=3500 DI=0100
```

`AX=0940` ⇒ 意圖據點 8、`SI=3500` ⇒ 軍團 75，而 **`DI=0100`**——
那不是據點記錄（據點表在 `0x840`），是**勢力表**的偏移
（勢力表基址 0，每筆 64 ⇒ 勢力 4）。

把 181 次呼叫的 `DI` 數一遍：

| DI | 次數 | 指到 |
|---|---:|---|
| `0B00` | 108 | 據點 22 |
| `1020` | 29 | 據點 63 |
| `09E0` | 14 | 據點 13 |
| `0E80` | 7 | 據點 50 |
| **`0100`** | **4** | **勢力 4** |
| 其餘 | 19 | 據點為主 |

⇒ **兩個來源**：`sub_1440F` 的第一行 `mov di, ax`（意圖據點的記錄）
與 `sub_143AF` **自己**的留守分支（勢力的記錄）：

```asm
.hold:  cmp  word ptr [si+4], 258h / jnb retn   ; 兵力 ≥ 600 ⇒ 直接回，di 不設
        mov  ah, [si+1] / xor al, al
        shr  ax,1 / shr ax,1 / mov di, ax       ; ★ di ← 勢力 × 64
        …
```

⇒ `cmp byte ptr [di+18h], 2` 讀的是**上一支軍團留下的那個位址** `+0x18`，
而同一個位移在兩張表是不同的東西：據點記錄是「那一格的軍團數」、
勢力記錄是「武將數」（[`../formats/08`](../formats/08-sinario-save.md)）。

這是原版的未初始化變數，但它**是決定性的**——同一個局面每次跑
都讀到同一個值。

⇒ remake 保存同一份狀態（`World.stage1DI`，存 raw 位址，
`≥ 0x840` 就是據點）。第一次呼叫之前 `di` 的值來自更早的呼叫鏈，
remake 用「未知 ⇒ 視為 > 2」。

⚠ **模型還不完整**：改成追兩個來源之後 5/25 與 5/31 維持全綠，
但 6/1 之前的拍 9,274 仍對不上——原版那時的 `DI=0100`（勢力 4），
而 remake 已經被某一次 Stage 2 覆蓋成據點。條件本身會回頭影響
「誰走到留守分支」，所以早期的一點分歧就會自我放大。

## 3. 驗證

`tools/parity_ck.sh 196/5/25`（拍 7,935）：
據點 26 B／18 座、勢力 1 B、軍團 19 B／2 支 → **六張表全 0**。

⭐ **對照實驗**：把 [`194`](194-ai-holds-a-threatened-city-on-the-way.md)
的 `holdThreatenedCity` 關掉再跑，5/25 變成 **180 個 byte／118 座**——
兩條規則都必要，不是互相抵消（`CLAUDE.md` §7 第 18 條）。
