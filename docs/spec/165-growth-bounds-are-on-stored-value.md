# 165 — 上昇值的上下限套在**存值**上，不是實際值

**狀態：CONFORMED。** 已改（`internal/rules/governor`）並加了邊界測試；
同局面對拍的**據點表在五個檢查點逐 byte 相同**，`+0x10` 的差異歸零。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V）`sub_14194`（IDA 線性 00014194）尾段三處
- 相關：[`../formats/08`](../formats/08-sinario-save.md)（據點 `+0x10` 帶 +100 偏移）、
  [`../re/07`](../re/07-monthly-settlement.md)、[`164`](164-relief-request-gates.md)

## 1. 問題

據點記錄 `+0x10`（上昇值）在檔案裡是**帶 +100 偏移的存值**：存 100 ＝ 實際 0。
remake 的 `state.City.Growth` 存的是**實際值**（`int(r[0x10]) - 100`）。

而 `internal/rules/governor` 的 `MaxValue = 200` 被同時套在上昇值與防災值上：

```go
c.Growth = min(c.Growth+gain, MaxValue)      // ← 實際值 ≤ 200 ⇒ 存值 ≤ 300（byte 溢位）
c.Growth = max(c.Growth-draft, 0)            // ← 實際值 ≥ 0   ⇒ 存值 ≥ 100（提早觸底）
```

**防災值沒有偏移**，所以 200 對它是對的；上昇值差了整整 100。

## 2. 原版怎麼寫

```asm
; ① 上昇值：+ gain，夾 0C8h
        mov     al, [si+850h]        ; ← 存值
        add     al, ch
        cmp     al, 0C8h
        jbe     short loc_141EA
        mov     al, 0C8h             ; 上限 200（存值）＝ 實際 +100
loc_141EA:
        mov     [si+850h], al

; ② 防災值：+ (gain >> 1) + 1，夾 0C8h
        mov     al, [si+851h]        ; ← 沒有偏移
        shr     ch, 1
        inc     ch
        add     al, ch
        cmp     al, 0C8h
        jbe     short loc_14207
        mov     al, 0C8h             ; 上限 200
loc_14207:
        mov     [si+851h], al

; ③ 徵兵拿上昇值換，下限 0
        sub     [si+850h], dl
        jnb     short loc_14225
        mov     byte ptr [si+850h], 0 ; 下限 0（存值）＝ 實際 −100
```

`internal/rules/economy` 早就有正確的界（`MaxGrowth = 100`／`MinGrowth = -100`），
`governor` 沒有跟上。

## 3. 量到的差

同一份存檔、同一條亂數流、跑 200 拍（[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §15）：

| 據點 | 主 | 起點存值 | 原版 200 拍後 | remake 200 拍後 |
|---:|---:|---:|---:|---:|
| 57 | 24 | 102 | 99 | **100** |
| 103 | 3 | 101 | 97 | **100** |
| 129 | 1 | 96 | 92 | **100** |
| 152 | 18 | 101 | 98 | **100** |

⭐ **remake 全部停在 100**——因為 100 就是它的地板。原版一路往下掉
（據點 112 已經到 12）。徵兵每次扣 4，而回補每天期望 +9/16，
**淨值對已經在低檔的據點是負的**，所以這一格會一路探底；
remake 把它擋在 0 就等於免費保底，影響是每個據點每一拍。

## 4. remake 要改什麼

`internal/rules/governor/governor.go`：把上昇值的界分出來。

```go
// MaxValue 是防災值的上限（`cmp al, 0C8h`）。防災值沒有偏移。
MaxValue = 200
// 上昇值在檔案裡帶 +100 偏移，原版的夾值做在**存值**上：
// 上限 0C8h ＝ 實際 +100、下限 0 ＝ 實際 −100。
MaxGrowth, MinGrowth = 100, -100
```

```go
c.Growth = min(c.Growth+gain, MaxGrowth)
…
c.Growth = max(c.Growth-draft, MinGrowth)
```

## 5. 怎麼驗

1. `governor` 的單元測試加兩個邊界案例（實際 +100 不再長、實際 −100 不再掉）。
2. 200 拍狀態對拍：據點 `+0x10` 的差異數要歸零
   （改前 8 個據點不同，1,900 拍時 34 個）。
