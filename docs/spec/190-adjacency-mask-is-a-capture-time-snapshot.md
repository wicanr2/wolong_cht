# 190 — 鄰接遮罩是易主時的快照，不是現況

**狀態：CONFORMED。** 據點記錄 `+0x00` 的低 4 位（四個鄰接槽裡哪幾個
屬於別的勢力）與 `+0x1B`（那幾位的個數）**只在據點易主時更新**
（`sub_14CF3` → `sub_188CC` → `sub_1890A`），而且**鄰接關係不對稱時
兩側都不更新**。remake 原本每個據點 tick 都由現況重算，兩者在同局面
196/5/17 的三座據點上分岔。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_188CC`（IDA 線性 000188CC–00018909）、
  `sub_1890A`（0001890A–0001895C）；唯一的呼叫鏈是 `sub_14CF3`（易主）
- 推論等級：confirmed（兩支逐條讀出，呼叫者只有一個；對拍樣本能分開兩種讀法）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §47.3
- remake 實作：`internal/state/strategy.go` 的 `refreshCityThreat`、
  `internal/state/corps.go` 的 `capture`
- 相關：[`../re/44`](../re/44-threat-and-reinforcement-ai.md) §5、
  [`187`](187-only-the-loser-is-judged.md)（同一支 `sub_14CF3`）

## 1. 原版怎麼維護

`sub_188CC`（`si` ＝ 剛易主的據點）掃自己的四個鄰接槽 `+0x1C`–`+0x1F`，
`dl` 從 1 開始每槽左移一位，槽值 `0xFF` 就跳過，其餘呼叫 `sub_1890A`：

```asm
sub_1890A:                       ; bx ＝ 鄰居記錄、si ＝ 我、ah ＝ 我的據點編號
        mov  di, bx / add di, 85Ch      ; 鄰居的 +0x1C
        mov  dh, 1 / mov cx, 4
.find:  cmp  ah, [di] / jz .hit
        shl  dh, 1 / inc di / loop .find
        retn                            ; ★ 鄰居的槽裡沒有我 ⇒ 兩側都不動
.hit:   cmp  al, [bx+841h] / jz .clear  ; 我的勢力 vs 鄰居的勢力
        ; ── 不同勢力：雙向設 ──
        test [bx+840h], dh / jnz .mine
        inc  byte ptr [bx+85Bh] / or [bx+840h], dh
.mine:  test [si], dl / jnz .done
        inc  byte ptr [si+1Bh] / or [si], dl
.done:  retn
.clear: ; ── 同勢力：雙向清（`not dh` / `and` / `not dh` 還原）──
```

三個性質：

1. **只在易主時跑。** `sub_188CC` 的唯一呼叫者是 `sub_14CF3`。
   平時的據點 tick（`sub_13EFD`）碰都不碰這兩欄。
2. **雙向。** 一次呼叫同時改鄰居那一側與自己這一側，計數各自 ±1。
3. ⭐ **鄰接不對稱時兩側都不更新。** A 的槽裡有 B，而 B 的槽裡沒有 A，
   `.find` 掃完四槽就 `retn`——連 A 自己那一位都不設。

⇒ 這兩欄是**易主那一刻的快照**，之後可以與現況不符，而原版就照著陳舊值
去判威脅、挑目標。

## 2. remake 錯在哪

`refreshCityThreat` 每個據點 tick 都跑
`threat.EnemyMask(c.Owner, ns)` 由現況重算。同局面 196/5/17 1 時
（拍 6,279）：

| 據點 | 原版 `+0x00` | remake | 原版 `+0x1B` | remake |
|---|---|---|---|---|
| 116 | `0B`（三個敵鄰）| `03` | 3 | 2 |
| 122 | `C2` | `C1` | 1 | 1 |
| 129 | `C6` | `C7` | 2 | 3 |

兩邊各自自洽（位元個數 ＝ 計數），差的是**該不該重算**。
據點 122 連個數都一樣，只有「是哪一槽」不同——**只比 `+0x1B` 看不出來**。

## 3. remake 的修法

- `refreshCityThreat` 不再寫 `c.Adjacency` 與 `c.EnemyNeighbours`，
  只讀它們餵 `threat.Scan`。
- `capture`（對應 `sub_14CF3`）在寫完新主之後呼叫 `updateAdjacencyOnCapture`，
  逐槽照 §1 的三個性質做增量維護。
- 快照載入本來就是從 `r[0x00]`／`r[0x1B]` 讀的，不必改。

⚠ **`threat.EnemyMask` 不刪**：開新遊戲與劇本載入時要有人算出初值
（原版的初值來自 `SINARIO.DAT`，remake 的新局面沒有那一份）。

## 4. 驗證

`tools/parity_ck.sh 196/5/17`：據點表由 5 個 byte／3 座變成 **0**。
