# 168 — 另外兩個巡迴游標：軍團 `+0x28`、每「時」勢力 `+0x2C`

**狀態：READY**

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V）`sub_125A3`（`cs:word_10D18`）、
  `sub_13E11`（`ds:0D1Ch` ／ `cs:word_10D1C`）
- 相關：[`162`](162-city-cursor-from-save.md)（據點游標 `+0x2E`）、
  [`166`](166-save-writeback-keeps-unmodelled-bits.md)

## 1. 三個游標，remake 只還原了一個

原版的三條輪轉都把游標放在存檔區塊的前 59 byte 裡，單位都是**段內偏移**：

| 存檔 | 原版符號 | 用途 | 步進 | 一圈 | remake |
|---|---|---|---|---|---|
| `+0x28` | `cs:word_10D18` | 軍團更新（`sub_125A3`，每拍 16 支）| ×0x40 | 128 支 ÷ 16 ＝ 8 拍 | `corpsCursor`，**沒載入也沒寫回** |
| `+0x2C` | `cs:word_10D1C` | 每「時」的勢力更新（`sub_13E11`）| ×0x40 | 22 勢力 ＝ 22 小時 | `hourFaction`，**沒載入也沒寫回** |
| `+0x2E` | `cs:word_10D1E` | 據點整備（`sub_13EFD`）| ×0x20 | 192 拍 | 已還原（[`162`](162-city-cursor-from-save.md)）|

[`162`](162-city-cursor-from-save.md) 講過不還原游標的後果：**不是差一格，是整條時間軸相位錯開**。
另外兩條一樣。

- 軍團游標錯開 ⇒ 每支軍團的 `+0x0B`（移動計時器）在錯的拍被 `dec`，
  行軍抵達時刻整體平移。同局面對拍量到的就是這個：`+0x0B` 固定差 1
  （[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §16）。
- 勢力游標錯開 ⇒ 侵攻的財政檢查、預備兵維持費、外交官都在錯的小時跑，
  勢力記錄的 `+0x1A`（支出）與 `+0x20`（資金）跟著漂。

## 2. 兩支的巡迴寫法

```asm
sub_125A3:                              ; 每拍 16 支
        mov     si, cs:word_10D18
        mov     cx, 10h
loc_125B5:
        …處理 [si+2240h]…
        add     si, 40h
        loop    loc_125B5
        cmp     si, 1FC0h
        jb      short loc_125F6
        xor     si, si
loc_125F6:
        mov     cs:word_10D18, si
```

```asm
sub_13E11:                              ; 每「時」1 個勢力
        mov     si, ds:0D1Ch
        …處理 [si]…
        add     si, 40h
        cmp     si, 580h                ; 22 × 0x40
        jb      short loc_13E57
        xor     si, si
loc_13E57:
        mov     word_10D1C, si
```

⚠ 兩支都是**先處理再前進**，與 `sub_13EFD` 同形（[`163`](163-city-tick-order.md)）。
remake 兩處本來就是這個順序，缺的只有「從存檔還原／寫回」。

### 2.1 ⭐ 軍團的一圈是 **128** 格，不是 127

`cmp si, 1FC0h` 是在 16 次 `add si, 40h`（＝ `+0x400`）**之後**才比的，
所以 `si` 在檢查點上只會是 `0x400` 的倍數：`0x400 … 0x1C00`，下一個
`0x2000` ≥ `0x1FC0` 就歸零。一圈實際走過的是 `si = 0 … 0x1FC0`，
**128 格、8 拍**。

軍團表本身也是 128 格：`0x22C0`（軍團）到 `0x42C0`（武將）之間是
`0x2000` byte ÷ 64 ＝ 128。remake 的 `numCorps` 是 **127**（跟著武將數走），
游標便以 127 取模——**每 8 拍就比原版多轉一格**。

實測（同一份快照跑 200 拍）：原版游標回到起點 48（3,200 ÷ 128 整除），
remake 走到 `(48 + 3200) mod 127 = 73`。

⇒ 游標的模數要用 **128**（`corpsSlots`），而世界只建模前 127 格：
輪到第 128 格就跳過。

## 3. remake 要改什麼

`internal/state/state.go`：

```go
corpsCursorOffset = 0x28   // ÷ 0x40 ＝ 軍團編號
hourCursorOffset  = 0x2C   // ÷ 0x40 ＝ 勢力編號
corpsSlots        = 128    // 軍團游標的模數（§2.1）
```

載入時還原、`Bytes()` 時寫回，換算與 `cityCursorOffset` 同形。
越界一律退回 0——存檔可能來自別的版本。

## 4. 怎麼驗

1. 0 拍載入→寫回，全域 `+0x28`／`+0x2C` 與來源逐 byte 相同。
2. 同局面對拍：軍團 `+0x0B` 的差異數要下降。
3. `World.Fingerprint` 已經把兩個游標算進去（`fingerprint_test.go`），
   所以還原之後指紋才真的代表「同一個局面」。
