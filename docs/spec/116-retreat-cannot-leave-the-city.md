# 116 — 驗收用的戰場少了子圖塊表，打破的門反而把城封死

**狀態：CONFORMED。** `internal/state/state_test.go` 的攻城 fixture 用
`NewFieldFromTiles` 建戰場——那是「合成戰場」的退路，地面層直接取堆疊高度。
在那條路上**打破的門會變成 5 層高的方塊**，城反而封死。
改成與正式路徑一致的 `NewFieldFromTileLayers` 之後，同一場戰鬥
第 967 幀結束（0.49 秒）。

- 日期：2026-09-02
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_1BCA6`（地面層）——
  [`../re/63`](../re/63-ground-plane-map.md) §2
- 推論等級：**confirmed（決定性重現 ＋ 離線最小化）**
- 相關：[`115`](115-soldier-power.md)（把它照出來的那一步）、
  [`36`](36-ground-planes-and-climbing.md)（兩個平面的地面圖）

## 1. 原版怎麼算地面層

`sub_1BCA6` 讀的是**七層子圖塊表**（`word_1D302`，每個圖塊 8 byte：
層數 ＋ 由下往上七層的子圖塊編號）：

```
低平面 = 層 0–3 裡第一個「頂面」（子圖塊 ≥ 0x70）的層號
if 圖塊 ≥ 0xF0:  高平面 = 8（門的特殊標記），return
高平面 = 層 4–6 裡第一個「頂面」的層號
```

⭐ **地面層來自子圖塊，不是堆疊高度。** 一個「上面有拱、下面是空的」圖塊
（門就是這種）堆疊高度很高，但腳下那一層在 0——**人走得過去**。

## 2. remake 有兩條路，差別就在這裡

| 建構子 | `layers` | `cellGround` 走哪一支 |
|---|---|---|
| `NewFieldFromTileLayers` | 有 | `groundOf`：照 §1 的原版演算法 |
| `NewFieldFromTiles` | nil | 退路：**地面層 ＝ 堆疊高度**（`f.top`）|

退路是給沒有子圖塊表的合成戰場用的（單元測試那種）。
**用真的戰場資料走退路會出錯**，因為門的堆疊高度不等於它的地面層。

實測濮陽（節點 56，旋轉）門格 `(26, 8)`：

| | 圖塊 | 低平面 | 高平面 |
|---|---|---:|---:|
| 有子圖塊表・完好 | `0xF5` | **0** | 4 |
| 有子圖塊表・打破 | `0xFD` | **0** | 4 |
| 退路・完好 | `0xF5` | **1** | 5 |
| ⛔ 退路・打破 | `0xFD` | **5** | 5 |

⭐ **打破之後低平面從 1 跳到 5**——`heights[0xFD]` 是那個圖塊的堆疊高度。
一步最多差一層（[`../re/63`](../re/63-ground-plane-map.md) §3），
所以地面上的兵再也走不過那一格。**門打破反而封死。**

## 3. 症狀

`TestNormalScenarioTacticalBattleTerminates` 的 fixture 走退路，
於是攻方打進城之後八個門陸續破掉，城就封住了：

```
[probe] 10 萬幀：攻方 97 守方 26 done=false
[probe] 20 萬幀：攻方 97 守方 26 done=false
```

大將體力 49 < 50 讓全軍退卻，但七個退卻中的兵找不到出城的路
（在**空場**上也找不到，所以不是擠在一起），`Remaining()` 因此永遠不到 0。

⚠ **這個症狀先前看不到**，因為 remake 把兵的戰力接成軍團士氣，
`sub_1B618` 的命中值 `rand(0..127) + 200` 恆命中、傷害 100–200 ⇒ 一擊必殺，
兵還沒退卻就死了。接回原版的戰力（約 18，[`115`](115-soldier-power.md)）
才把這條路走到。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| fixture 改帶子圖塊表 | `internal/state/state_test.go` 的 `TacticalSetup.Field` |
| 正式路徑 | `internal/battlesetup` 本來就用 `NewFieldFromTileLayers`，不必改 |

⚠ **正式遊戲沒有這個問題**——`internal/battlesetup/battlesetup.go` 一直
帶著子圖塊表。壞的只有驗收 fixture，而那正是「驗收環境與正式環境不同」
會產生的那種假象：測試綠了半年，綠的是一張**與遊戲不同的地圖**。

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 迴歸 | `TestNormalScenarioTacticalBattleTerminates`：接上 [`115`](115-soldier-power.md) 的戰力之後**第 967 幀結束**（0.49 秒），先前 20 萬幀不結束 |
| 離線最小化 | 拿節點 56 旋轉後的戰場、照戰鬥當下打破八個門（`0xF5` → `0xFD`），比兩個建構子的 `GroundLevel` 與 `FindPath`——退路那一邊路徑歸零 |

## 6. 未解

| 項目 | 現況 |
|---|---|
| 還有誰在用 `NewFieldFromTiles` 配真戰場資料 | `internal/rules/tactical/tactical_test.go` 有一處。它驗的是圖塊解碼不是連通性，但同一個陷阱在那裡也成立 |
| 退路要不要留 | 合成戰場（`NewField`）確實只有堆疊高度。**能不能在圖塊 ≥ `0xF0` 時不用堆疊高度**，是另一個問題 |
