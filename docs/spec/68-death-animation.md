# 68 — 倒地動畫：四幀，換一組圖

**狀態：CONFORMED。** 圖號、幀數與分組都是逐行讀出來的，remake 已經畫出來了。

- 日期：2026-08-18
- 出處：`sub_1B360`（`0001B360`–`0001B3B2`），呼叫者 `sub_1ADC8`（`0001AE09`）；
  進入條件 `sub_1B618` 的 `0001B697`
- 推論等級：**confirmed（靜態）**
- 輸入檔：`workplace/orig/dosv/KI.EXE`

## 1. 原版做什麼

血扣到 0 時（`sub_1B618`）：

```asm
0001B697  and byte ptr [di],   10h   ; ★ 只留 bit 4，bit 7（在場）被清掉
0001B69A  or  byte ptr [di],    1    ; bit 0 ＝ 陣亡
0001B69D  mov byte ptr [di+1],  4    ; ★ 計時 ← 4
```

`[si+1]` 與硬直用的是**同一個計時器**（[`63`](63-hit-stun.md)）。
bit 7 被清掉之後每幀的單位迴圈走另一條（`0001AE00`），
呼叫 `sub_1B360` 重畫，計時歸零時由 `sub_1B4B8(ah = 1)` 收掉——
**不算生還**（[`65`](65-retreated-soldiers-survive.md)）。

### 1.1 圖號

```asm
0001B38A  mov cx, 168h              ; 側 0
0001B38D  cmp si, 600h / jb .s0
0001B393  mov cx, 21Ch              ; ★ 側 1
.s0:
0001B396  cmp byte ptr [si+4], 24h  ; ★ 兵種（記錄裡已經是 × 18）
0001B39A  jb  .arm0                 ;   < 36 → 大將／騎馬
0001B39C  ja  .arm2
0001B39E  add cx, 4                 ;   ＝ 36 → 弓
0001B3A1  jmp .arm0
.arm2:
0001B3A3  add cx, 8                 ;   > 36 → 步
.arm0:
0001B3A6  cmp byte ptr [si+1], 2
0001B3AA  ja  .draw
0001B3AC  inc cx / inc cx           ; ★ 計時 ≤ 2（後兩幀）→ ＋2
.draw:
0001B3AE  call sub_1DA1C
```

`cx` 這裡已經是**合併圖形表的 raw unit**，不像 `sub_1B240` 還要
`shl cx,1 / add cx, 0C0h`。換算回圖號：

```
raw 0x168 ＝ 360 ＝ 192 + 84 × 2   → 圖號 84（側 0）
raw 0x21C ＝ 540 ＝ 192 + 174 × 2  → 圖號 174 ＝ 84 + 90（側 1）
```

90 正是一側的圖數（`sub_1B240` 的 `add cx, 5Ah`），所以兩側是同一張表的兩段。

| 兵種（`+0x04`）| raw ＋ | 圖號 ＋ |
|---|---:|---:|
| 0（大將）、18（騎馬）| 0 | 0 |
| 36（弓）| 4 | 2 |
| 54（步）| 8 | 4 |

| 計時 | raw ＋ | 圖號 ＋ |
|---:|---:|---:|
| 4、3 | 0 | 0 |
| 2、1 | 2 | 1 |

⭐ **大將與騎馬共用同一組倒地圖**——`cmp` 只分三段，不是四種兵種各一組。

### 1.2 ⚠ 倒地中的兵不擋路

`0001B697` 的 `and [di], 10h` 把 **bit 7（在場）清掉**，而占格與選目標
都看那個位元。所以倒地那四幀**只是畫面**：不擋路、不被打、不算場上人數。

## 2. 演算法

```
血歸零:  記一筆倒地（位置、側、兵種），計時 ← 4；該兵離場（不算生還）
每幀:    計時 −= 1；歸零就把那一筆刪掉
畫:      圖號 ＝ 84 ＋ 側 × 90 ＋ 兵種組 × 2 ＋ (計時 ≤ 2 ? 1 : 0)
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 倒地清單 | `internal/rules/tactical`：`Battle.deaths`，`Deaths()` 給呈現層讀 |
| 產生 | `damage.go` 的致命傷那一條，與既有的 `Alive = false` 同一處 |
| 每幀遞減 | `Battle.Step` 的同一輪，與投射物同級 |
| 畫 | `internal/ui/isoview/isoview.go` 的 `buildDisplayList`，走既有的 `appendTallDisplayUnits` |

⭐ **不擋路是照 §1.2 做的**：倒地的兵一離開 `Soldiers` 就不再參與規則，
只留一筆給畫面用。這也讓「四幀不動的屍體」不會把戰鬥卡住。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestDeathSpriteNumbers`：八組（側 × 兵種 × 計時）的圖號，外加把兩個 raw 常數 `0x168`／`0x21C` 對回來 |
| 單元測試 ✅ | `TestDeathAnimationLastsFourFrames`：打死一個兵 → 一筆倒地、四幀後消失，而且那四幀不進 `Remaining()` |
| 迴歸 ✅ | `TestTacticalBattleAlwaysResolves` 仍在 5,000 tick 內分出勝負 |
| 對拍迴歸 ✅ | 第 61 步九區的數字不變（`field` 307 px／0.17%）。⚠ 那個取樣點 2026-08-27 起不再等價（`../playtest/49`），這一列是當時的迴歸證據 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 大將陣亡 | 大將不會死（`sub_1B618` 的 `IsGeneral` 那一條），所以 `+0` 那一組實際只有騎馬會用到；大將的倒地圖是不是死碼還沒查 |
