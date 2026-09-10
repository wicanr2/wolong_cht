# 188 — 玩家的第一鄰居：和平掉 8，交戰只掉 1

**狀態：CONFORMED。** `sub_12DF3` 對玩家勢力的排序第一鄰居先 −1，
接著 `cmp bh, 80h / jnb` —— **`bh` 的最高位元在「交戰」時是 1**，
所以那條 `jnb` 是**交戰就跳過**。額外的 −7 只發生在**和平**的對象身上。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_12C52`（IDA 線性 00012C52–00012CDE）、
  `sub_12DF3`（00012DF3–00012E32）
- 推論等級：confirmed（兩支逐條讀出，位元的設定端與消費端都找到；
  對拍歸零是第二個來源）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §46.22
- remake 實作：`internal/state/strategy.go` 的 `driftPlayerFriendship`
- 相關：[`../re/07`](../re/07-monthly-settlement.md) §22、
  [`../mechanics/70-ai.md`](../mechanics/70-ai.md)

## 1. 位元是誰設的：`sub_12C52`

排序緩衝區的每一筆是 **勢力偏移（勢力編號 × 0x40）**，
選擇排序每一輪找**交友度最小**的那個鄰居往前放：

```asm
mov  bx, si / shr bx,1 / shr bx,1     ; bx ＝ si/4
mov  ax, bx / shr bx,1 / add bx, ax   ; bx ＝ si×3/8 ＝ 勢力編號 × 24
add  bx, 600h                         ; ⇒ 交友度矩陣裡「自己」那一列
...
mov  ax, di / and ah, 7Fh             ; 清掉高位元再取編號
shl  ax,1 / shl ax,1 / mov al, ah     ; al ＝ (di & 7FFFh) >> 6 ＝ 勢力編號
xlat                                  ; al ← 交友度[自己][那個勢力]
cmp  dl, al / jbe loc_12CCA           ; 比目前最小值大 ⇒ 不換
mov  dl, al
cmp  al, 80h / jnb loc_12CBF
or   di, 8000h                        ; ★ 交友度 < 80h（交戰）⇒ 設 bit 15
xchg di, es:[bp+0]
```

⭐ 交友度的 **bit 7 是和平位元**（`sub_130F0` 用 `and cx, 807Fh`
把值與位元拆開，`docs/re/69`）。所以 `al < 0x80` ＝ **交戰中**，
而緩衝區那一筆的 **bit 15 ＝「這個鄰居正在跟我打」**。

## 2. 位元是誰看的：`sub_12DF3`

```asm
mov  di, es:[di] / cmp di, 0FFFFh / jz .rest
mov  bx, di                    ; ★ bx 留著原始值（含 bit 15）
and  di, 7FFFh                 ; di 才是乾淨的勢力偏移
mov  al, 1 / call sub_130F0    ; 先 −1
cmp  bh, 80h / jnb .rest       ; ★ bit 15 設了（交戰）⇒ 跳過
mov  al, 7 / call sub_130F0    ; 只有和平才再 −7
```

勢力偏移最大 `21 × 0x40 = 0x540`，`bh` 本身只到 `0x05`，
所以 `bh >= 0x80` 唯一的來源就是 §1 那個 `or di, 8000h`。

⇒ **和平的第一鄰居每月掉 8，交戰中的只掉 1。**

## 3. 這個極性才是自洽的

| | AI ↔ AI（`sub_12DB8`） | 玩家的第一鄰居（`sub_12DF3`） |
|---|---|---|
| 和平時 | −2（下限 20） | **−8** |
| 交戰時 | +1（上限 50，對玩家不加） | **−1** |

兩支的方向一致：**和平往下推、交戰不往下推**。
和平會自己惡化到開戰，開戰之後衰減幾乎停住——AI 那邊甚至回升（厭戰）。
「交戰再 −7」會讓交戰狀態自我強化，與 `sub_12DB8` 的厭戰設計互相矛盾。

## 4. remake 的差異與修法

`driftPlayerFriendship` 的條件反了：

```go
fr = fr.WithValue(fr.Value() - 1)
if !wasWar {                       // 原本是 if wasWar
    fr = fr.WithValue(fr.Value() - 7)
}
```

`wasWar` 取的是**排序緩衝區那一筆的快照**，不是當下的交友度：
`sub_12C52` 在月結一開始就把 22 個勢力的緩衝區全部建好，
之後每個勢力的漂移／合作／停戰都會改交友度，**後面的勢力仍然看自己那一份**。

## 5. 驗證

`tools/friendship_diff.py`（新工具，比區塊 `+0x680` 的 22 × 24 矩陣）
在修正前抓到同局面**勢力 2 → 13 差 7**（原版 `0x9C` ＝ 28、
remake `0xA3` ＝ 35），二分後定位到 5/1 月結（拍 2,967）。
修正後拍 3,400 與 5,080 兩個檢查點的交友度矩陣**逐 byte 相同**。

⚠ 這張表在快照裡一直存在，2026-09-10 之前沒有任何工具在比——
據點、勢力、軍團、全域、事件佇列五張表全綠了好幾天，
而**宣戰、停戰、協力的主要效果都落在這第六張表上**。
