# 192 — 選路的成本不是格數：每個節點 +4，敵城 +0xA6

**狀態：CONFORMED。** 原版每次要走下一步都跑一次 `loc_1491B`——一個
**從目標往回**的 uniform-cost 搜尋。成本不是路徑格數：

```
成本 ＝ Σ(邊長 [bx+4])
     ＋ 4      × 經過的節點數
     ＋ 0xA6   × 經過的**非己方據點**數（同時設 bit 15）
```

remake 的 `march.Graph.Route` 用 Dijkstra，權重只有 `Edge.Steps`（格數）。
兩種度量在同局面 196/5/18 的軍團 37 上選出不同的路。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`loc_1491B`（IDA 線性 0001491B–00014A0E，
  **IDA 認不出函式，要用 `tools/ida_range.py` ＋ `ida_dump_bytes.idc`
  交叉解碼**）、`sub_14A0F`（00014A0F–00014A7A）、
  `sub_147BB`（000147BB，唯一的呼叫端）
- 推論等級：confirmed（逐條解碼 ＋ 對拍歸零）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §47.5
- 相關：[`43`](43-rout-on-blocked-return.md)（`0xA6` 與 bit 15 的消費端）、
  [`../re/08`](../re/08-hourly-update.md) §7

## 1. 這一段是自我修改碼

`0x1495C`–`0x149E0` 在 IDA 裡是一片 `db`，因為進入時有三行把值
**patch 進後面的立即值**：

```asm
mov cs:word_149B8, bx      ; → 149B6 的 `cmp si, imm16`
mov cs:word_149BE, cx      ; → 149BC 的 `cmp si, imm16`
mov cs:byte_149D2, dl      ; → 149CD 的 `cmp byte ptr es:[si+841h], imm8`
```

⇒ 交叉解碼才讀得出來（`CLAUDE.md` §7 第 28 條）。
被 patch 的三個值分別是**兩個終止條件**（當前所在邊的兩端）與
**自己的勢力編號**。

## 2. 搜尋方向是反的

```asm
mov si, ax                 ; ★ ax ＝ 目標節點 ⇒ 搜尋從**目標**開始
mov di, 8000h / mov cx, 400h / rep stosw    ; 清 0x400 word 的 visited
…
cmp si, word_149B8 / jz .hit                ; 碰到當前邊的 A 端
cmp si, word_149BE / jz .hit                ; 碰到當前邊的 B 端
```

⇒ **從目標往回搜，先碰到自己這條邊的哪一端，就往那一端走。**
回傳的 `ax`／`bx` 是那一筆佇列項目的 `+4`／`+6`
（`sub_14A0F` 寫的 `bp` ＝ ±4 的步進與鄰接項），`cx` ＝ 累計成本。

⚠ **不要把 `word_149B8`／`word_149BE` 讀成起點**——2026-09-10 就是
這樣讀錯的（`docs/lessons.json` 的 `self-modifying-code-reads-as-data`）。

## 3. 成本三項

| 位置 | 指令 | 加什麼 |
|---|---|---|
| 展開一個節點時 | `add dx, 4`（`149DC`）| **每個節點 4** |
| 那個節點是**據點**（`si < 600h`）而且 `es:[si+841h] ≠ 自己` | `add dx, 0A6h` ＋ `or dh, 80h`（`149D5`）| **敵城 166，並設 bit 15** |
| 推進一個鄰居時 | `add dl, [bx+4] / adc dh, 0`（`14A3D`）| **那條邊的長度** |

bit 15 的消費端在 `sub_147BB`：`cmp cx, 8000h / jnb` ⇒ 成本帶著那個位元
（路上有敵城）而軍團 `+0x23 ≥ 0x0A`（退卻中）時直接 `sub_1291A` 潰散
（[`43`](43-rout-on-blocked-return.md)）。

⇒ `0xA6` **不是「繞開敵城」的軟性偏好**，它大到足以蓋過任何合理的
邊長差；真正的效果是「有敵城的路只在沒有別條路時才走」。

## 4. 兩邊的圖是同構的

一度以為原版多了一層「野外節點」——軍團記錄的 `+0x0E` 平常確實停在
`≥ 0x800` 的值上。但那不是節點編號，是**連結記錄的位址**
（[`172`](172-corps-march-fields.md) §1）：節點只有 192 個據點，
`si < 600h` 就是 `據點 × 8`。佇列裡的節點來自
`sub_14A0F` 的 `di = [bx+6]`／`[bx+8]`，也就是連結記錄的**兩端節點**。

⇒ remake 的 192 節點 × 254 邊與原版同構，成本模型可以直接對齊：

```go
邊 (u → v) 的權 ＝ l.steps ＋ 4 ＋ penalty(v)
```

⚠ **節點成本要掛在 `l.to` 上。** 原版從**目標**往回搜，而且終止檢查
（`cmp si, word_149B8`）排在 `add dx, 4` **之前**——所以「軍團現在
站的那個節點」不算、目標算。反過來搜時掛在 `l.to` 才等價。

## 5. 驗證

軍團 37（據點 82 → 64）原本第一步往東西、原版往南北。改成原版的成本
模型之後，`tools/parity_ck.sh 196/5/16 196/5/17 196/5/18 196/5/20`
的軍團表**只剩 `corps-move-timer-phase` 那 1 個 byte**，六張表其餘全 0。

## 6. 未解

同成本時的 tie-break：原版是環形佇列 ＋「掃一遍取第一個等於最小值的
項目」，remake 是二元堆。**同成本的路一多就可能分岔**，目前還沒有
樣本能分開。
