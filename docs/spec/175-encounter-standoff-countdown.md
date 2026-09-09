# 175 — 撞上敵人不是立刻開打：`+0x03` 先倒數 12 次（96 拍）

**狀態：READY（`+0x03` 的用途已分開，見 §3；對峙的觸發還沒接）。**
原版軍團要踏進「敵方軍團佔著的格子」或「別人的據點」時會**停住**，
設 `+0x00` 位元 5，`+0x03` 從 `0Ch` 開始倒數，每個軍團巡迴週期減 1，
**減到 1 的那一次才開打**。remake 一走到就結算，因此早 96 拍。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_12708`（IDA 線性 00012708）、
  `sub_12831`（00012831）、`sub_12880`（00012880）、`sub_1264A`（0001264A）、
  `sub_14A7B`（00014A7B）、`sub_14ADE`（00014ADE）、`sub_12977`（00012977）、
  `sub_12B3C`（00012B3C）、`sub_125A3`（000125A3）
- 推論等級：confirmed（四支常式逐條讀出，且實測的 96 拍與 12 × 8 吻合）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §32
- 相關：[`43`](43-rout-on-blocked-return.md)（`+0x03` 的另一個用途）、
  [`174`](174-relief-dispatch-does-not-move-marching-corps.md)

## 1. 原版做什麼

### 1.1 下一格有人就不走

`sub_12708` 寫座標之前先看佔用圖：

```asm
        cmp     byte ptr [di], 0      ; di ＝ 佔用圖上「要踏進去的那一格」
        jz      short loc_1273C       ; 空的 → 照常走
        mov     ds, cs:word_10D52
        mov     dx, di
        xor     ah, ah
        call    sub_12831
        jb      short loc_1273C       ; CF=1 → 還是走
        retn                          ; ⭐ CF=0 → **這一拍不動**
```

### 1.2 `sub_12831`：找出佔那一格的是誰

```asm
sub_12831:
        mov     bx, [si+0Eh]
        mov     di, 2240h             ; 軍團表
        mov     cx, 7Fh               ; 127 支
loc_1283D:
        cmp     byte ptr [di], 80h
        jb      short loc_1284C
        cmp     ax, [di+12h]          ; Y 相同？
        jnz     short loc_1284C
        cmp     dx, [di+10h]          ; X 相同？
        jz      short loc_12853
loc_1284C:
        add     di, 40h
        loop    loc_1283D
        jmp     short loc_1287B       ; 找不到 → stc（放行）
loc_12853:
        mov     cl, [si+1]
        cmp     cl, [di+1]
        jz      short loc_1287B       ; 自己人 → stc（放行）
        or      byte ptr [si], 20h    ; ⭐ 位元 5 ＝「前面卡著敵人」
        cmp     byte ptr [si+3], 1
        ja      short loc_1286C       ; > 1 → 還在倒數
        jz      short loc_12873       ; ＝ 1 → ⭐ 開打
        mov     byte ptr [si+3], 0Ch  ; ⭐ 第一次撞上：倒數設 12
        jmp     short loc_12876
loc_1286C:
        mov     al, 3
        call    sub_102F5             ; 倒數期間的音效
        jmp     short loc_12876
loc_12873:
        call    sub_14A7B             ; 開打
loc_12876:
        clc                           ; ⭐ 不放行
```

### 1.3 倒數在軍團巡迴裡減

`sub_125A3` 每次巡到一支軍團就先走 `sub_1264A`：

```asm
sub_1264A:
        test    byte ptr [si], 20h
        jnz     short loc_12658
        mov     byte ptr [si+3], 0    ; 沒卡著 → 歸零
        mov     byte ptr [si+21h], 0
        retn
loc_12658:
        dec     byte ptr [si+3]       ; ⭐ 卡著 → 每個巡迴週期減 1
        jnz     short locret_12661
        mov     byte ptr [si+3], 1    ; 減到 0 就停在 1，等下一次撞上時開打
```

⭐ **`sub_1264A` 由 `sub_125A3` 呼叫，而巡迴一圈是 8 拍**
（[`168`](168-corps-and-hour-cursors.md)），所以倒數 12 次 ＝ **96 拍**。
實測完全吻合：原版軍團 19 在拍 2,024 停住（旗標 `C5` → `F5`），
一直到拍 2,120 才開打（[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §32）。

### 1.35 據點那一條走 `sub_12880`，同一套倒數

`sub_12708` 還有第二個入口：位元 0 設著（已經走上路徑）而下一格的地形值
落在 `0CEh`–`0DDh`（**據點自己的圖塊**）時走 `sub_12880`：

```asm
sub_12880:
        mov     di, [si+0Eh]          ; 連結記錄
        cmp     byte ptr [si+0Ah], 4  ; 步進 ＝ 4（正向）？
        jnz     short loc_12892
        mov     di, es:[di+8]         ; 連結記錄 +08 ＝ B 端節點
        jmp     short loc_12896
loc_12892:
        mov     di, es:[di+6]         ; +06 ＝ A 端節點
loc_12896:
        shl     di, 1 / shl di, 1     ; → 據點記錄位址
        mov     al, [si+1]
        cmp     [di+841h], al         ; 那座據點是自己的嗎
        jz      short loc_128C7       ; 是 → stc（走進去）
        or      byte ptr [si], 20h    ; ⭐ 同一個位元 5
        cmp     byte ptr [si+3], 1
        ja      short loc_128B4       ; > 1 → 音效
        jz      short loc_128BB       ; ＝ 1 → ⭐ 攻城／接收
        mov     byte ptr [si+3], 0Ch  ; ⭐ 同一個 12
        jmp     short loc_128C2
loc_128BB:
        add     di, 840h
        call    sub_14ADE
loc_128C2:
        clc                           ; 不放行
```

⇒ **兩個入口共用位元 5 與 `+0x03`，倒數都是 12。**
差別只在倒數結束時呼叫誰：撞軍團走 `sub_14A7B`（野戰），
撞據點走 `sub_14ADE`（攻城／接收）。

⭐ 實測的那一場（軍團 19 → 據點 122）走的是**據點這一條**：
軍團停在最後一個路徑點 `(302,158)`，下一格 `(303,158)` 是 122 的城門格。

### 1.4 開打時把兩邊的旗標與倒數都清掉

```asm
sub_14A7B:
        …
        call    sub_14C72             ; 開戰前置；CF=1 就不打
        mov     di, bx
        and     byte ptr [si], 0DFh   ; 攻方清位元 5
        mov     byte ptr [si+3], 0
        and     byte ptr [di], 0DFh   ; 守方也清
        mov     byte ptr [di+3], 0
        call    sub_14E5C             ; 戰鬥判定，ah ＝ 結果
        …                             ; 贏的一方 sub_1291A（接收據點）
```

### 1.5 `+0x03` 是一個 byte、兩種用途

| 什麼時候 | 判準 | 用途 | 初值 |
|---|---|---|---|
| 敗走中 | `+0x00` ＝ `08h`（**已經不存在**，`sub_12977` 整個 byte 寫 8）| 敗走倒數 | `30h` ＝ 48 |
| 對峙中 | `+0x00` 位元 5，而且軍團**還活著**（≥ `80h`）| 對峙倒數 | `0Ch` ＝ 12 |
| 其他 | 兩者都不成立 | 恆為 **0**（`sub_1264A` 每個巡迴週期歸零）| — |

**兩種狀態互斥**：敗走的軍團旗標 < `80h`，對峙的軍團旗標 ≥ `C0h`。
原版檢查點的 44 筆活著的軍團記錄，`+0x03` **全部是 0**。

⚠ 對峙倒數還兼**動畫相位**：`sub_12B3C`（大地圖畫軍團）拿
`[si+3] & 3` ＋ `[si+21h] << 2` 當索引取圖。所以那 12 拍不只是等待，
畫面上是兩軍對砍的動畫。

## 2. remake 差在哪

`resolveContact` 一踏進敵人那一格就結算，**沒有對峙這一段**。
後果不是「早一點打完」，是**整條時間軸從那一刻起錯開 96 拍**——
同局面對拍的第一個分歧就是它（拍 2,023 對 2,119）。

⚠ 還有一個副作用：原版對峙期間軍團**停在原地**，所以那 96 拍裡
它佔的格子、據點的 `+0x18`、威脅量都跟著不同。

## 3. remake 要改什麼

分兩步做。**第一步（已完成）：把 `+0x03` 的兩種用途分開。**

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/corps.go`：`Corps.Standoff` ← `+0x00` 位元 5，進 `modelledCorpsBits` |
| 狀態層 | `internal/state/corps.go`：`RoutTimer` 改名 `Countdown`，一個欄位對一個 byte |
| 狀態層 | `internal/state/corps.go`：軍團巡迴每次先跑 `tickStandoff`（＝ `sub_1264A`）|
| 狀態層 | `internal/state/corps.go`：`saveCorps` 對**活著的**軍團也寫 `+0x03` |

⭐ **一個 byte 就用一個欄位。** 拆成 `RoutTimer` 與 `StandoffTimer` 兩個
Go 欄位會多出「兩個欄位、一個 byte」的失步風險，而原版本來就是靠
`+0x00` 分辨——照抄那個結構，載入、寫回、讀取三處都不可能對不上。
分開的是**用途**（哪個狀態讀它），不是儲存。

**第二步（未做）：接上觸發。** `step` 要在踏進「敵方軍團佔著的格子」
或「別人的據點」時不移動、設 `Standoff`、`Countdown ← 12`，
倒數到 1 才結算。⚠ 這一步會動到遭遇與攻城的入口，要另外驗。

## 4. 怎麼驗

第一步（`+0x03` 分開）：

1. 五個檢查點的軍團表維持**逐 byte 相同**——原版活著的軍團 `+0x03`
   全部是 0，所以「活著也寫 `+0x03`」不能讓任何一格變。
2. 單元測試：兩種倒數**不會同時非零**；載入→寫回 `+0x03` 與位元 5 無損；
   對峙中每個巡迴週期減 1、減到 0 停在 1；不在對峙就歸零。

第二步（接上觸發）：

3. 同局面對拍：第一個分歧從拍 2,023 往後推；原版與 remake 的第一場
   AI 自動判定戰鬥都落在拍 2,119。
4. 軍團表在拍 2,024–2,100 的旗標是 `F5`、`+0x03` 逐週期遞減。

## 5. 未解

| 項目 | 現況 |
|---|---|
| `+0x00` 位元 4 | `sub_12B3C` 設、`sub_12BA8` 清，而 `sub_12BA8` 接著呼叫 `sub_19656`／`sub_196ED`（繪圖）。**像是「這一格要重畫」的髒旗標**，不是規則狀態——待確認 |
| `+0x21` | `sub_1264A` 在沒卡住時一併歸零；`sub_12B3C` 拿 `<< 2` 與 `+0x03 & 3` 合成圖塊索引。**像是對峙動畫的第二個維度**，語意未讀 |
| `sub_102F5(al=3)` | 對峙期間每個週期呼叫一次，推測是音效 |
