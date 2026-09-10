# 181 — AI 的募兵節制；月結的顯示用收支在保存區塊裡

**狀態：CONFORMED。** 兩件事，都是靠**新開的兩張比對表**才看見的：

1. `sub_15456` 對 AI 勢力有一道**募兵節制閘**——軍團養太多相對於收入
   就整個月不募兵。remake 沒有這道閘，AI 的預備兵一路往上漂。
2. 財政視窗的月結快照（`cs:0D02h`／`0D05h`）**在保存區塊裡**
   （`+0x12`／`+0x15`，各 24 位），remake 標成「不進存檔」而沒寫。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_15456`（IDA 線性 00015456–0001548E）、
  `sub_1548F`（0001548F）
- 推論等級：confirmed（整支逐條讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §41
- 相關：[`../re/07`](../re/07-monthly-settlement.md) §8（原始反組譯）、
  [`14`](14-finance-window.md) §5（顯示用全域）

## 1. AI 的募兵節制

```asm
sub_15456:
        shr  byte ptr [bp+2], 1 / rcr word ptr [bp+0], 1   ; 收入 ÷ 2
        ax = si << 2 / al = 80h        ; ah ＝ 勢力編號、al ＝ 存在門檻
        xor dx, dx / cx = 7Fh / bx = 2240h
loop:   cmp [bx], al   / jb  skip      ; 軍團活著？
        cmp [bx+1], ah / jnz skip      ; 屬於這個勢力？
        add dx, [bx+4]                 ; 累加兵力（16 位，會繞回）
skip:   add bx, 20h / loop
        mov al, dh / xor ah, ah        ; ★ 只取和的**高 byte**
        add ax, [si+1Bh]               ; ★ ＋ 勢力記錄 +0x1B（本月支出 >> 8）
        shl ax, 1                      ; ★ × 2
        cmp ax, [bp+1] / jnb .no       ; ★ ≥ 收入 >> 8 ⇒ 這個月不募兵
```

四個「取高 byte」都是照抄。拿完整值去比會得到完全不同的答案——
**這道閘實際上幾乎總是成立**，AI 因此很少募兵。

同局面 196 年 5 月 1 日的月結，三個抽樣勢力全部被擋：

| 勢力 | 本月支出 >> 8 | 收入 >> 8 | 原版募到 |
|---:|---:|---:|---|
| 1 | 4 | 1 | 0 |
| 5 | 37 | 41 | 0 |
| 7 | 5 | 6 | 0 |

⛔ **原版掃錯了步進，remake 照抄。** `add bx, 20h` 走的是 32 B，而軍團
記錄是 64 B；127 次迴圈只走到軍團表的一半，而且**奇數次落在記錄中間**
——那時 `[bx]` 讀到的是 `+0x20`（意圖）、`[bx+1]` 是 `+0x21`、
`[bx+4]` 是 `+0x24`。所以被加總的不只是軍團兵力。remake 的
`aiCorpsWeight` 直接對當前狀態的區塊 bytes 掃，才拿得到那幾格
「不該被當成欄位」的值。

## 2. 顯示用的收支在保存區塊裡

`cs:0CF0h` 那 59 個 byte 是連在一起存的，其中 `+0x12`（＝ `cs:0D02h`）
是本月收入、`+0x15`（＝ `cs:0D05h`）是本月支出，各 24 位。
`sub_1548F` 在玩家的月結尾段寫它們。

remake 的 `IncomeSnap`／`ExpenseSnap` 原本註明「不進存檔」——
存檔語意上確實不需要，但**同局面對拍比的是那 59 個 byte**。
現在載入端讀、寫回端寫。

## 3. 改了哪幾支

| 檔案 | 改什麼 |
|---|---|
| `internal/rules/economy/economy.go` | `Faction.CorpsWeight`；`Settle` 在募兵前跑 `recruitBlocked` |
| `internal/state/state.go` | `aiCorpsWeight`（照抄錯步進）；`incomeSnapOffset`／`expenseSnapOffset` 的載入與寫回 |

## 4. 為什麼拖到現在才發現

**勢力表與全域欄位從來沒被比過。** `city_diff` 比據點、`corps_diff` 比軍團，
而預備兵、資金、本月支出、三個游標、稅率、募兵數全在另外兩張表裡。

於是 AI 的預備兵從 5 月 1 日的月結起一路偏高，而**症狀出現在完全不相干
的地方**——五百拍之後某支軍團補兵時每槽多分了 4 個兵。

現在 `tools/parity_ck.sh` 一次比四張表：據點、勢力、軍團、全域。
兩支新工具各有正反對照，接進 `check.sh`。

## 5. 未解

- 勢力 7 的本月支出在月結那一小時差一次累加（原版 46、remake 0），
  資金跟著差 46。像是「月結與每小時勢力更新」在同一小時內的順序。
- 拍 3,400 的每時勢力游標（`+0x2C`）差 1。
