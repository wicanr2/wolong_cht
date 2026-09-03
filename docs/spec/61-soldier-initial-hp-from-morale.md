# 61 — 兵的開場體力 ＝ 軍團士氣

**狀態：CONFORMED。** 出處是 `sub_19B6D` 的兩行，
`docs/re/11` §3.9 早就寫下來了。接上之後同一個局面的側欄
`sb-self` 逐像素相同、`sb-enemy` 1.48% → 0.86%。

⚠ 這一條要與 [`62`](62-swapped-unit-skips-its-turn.md) 一起看：
兵撐得比 100 點久之後才會擠在一起，而擠在一起會不會打得到人
由那一條決定。

- 日期：2026-08-18（§6 的大將體力補於 2026-09-03）
- 出處：`sub_19B6D`（`00019BBD`）、`sub_1B97E`（回復上限）、
  `sub_1AE56`（`0001AE56`，攻城計時器與退卻門檻，§6）、
  `sub_1C6F6`（把它畫成側欄那條體力條）
- 推論等級：**confirmed（靜態）**

## 1. 原版做什麼

編一隊上場時（`sub_19B6D`），八個兵的記錄一次寫齊：

```asm
00019B73  mov ah, [bx+6]            ; ★ 軍團摘要的 +0x06 ＝ 士氣
...
00019BBA  mov cx, 8                 ; 一隊八個兵
00019BBD  mov es:[di+3], ah         ; ★ 體力 ← 士氣
00019BC1  mov es:[di+4], al         ; 兵種 × 18
00019BC5  mov es:[di+18h], dl       ; 戰力（另一條算式）
00019BC9  mov byte ptr es:[di+19h], 80h
00019BCE  mov word ptr es:[di+1Ah], 1
00019BD4  add di, 20h / loop
```

`word_1D30A:+0x06` 開打時由 `sub_19E70` 從軍團記錄抄進來
（[`../re/11`](../re/11-tactical-battle.md) §3.9），也就是大地圖上看得到的
那個士氣值。**士氣高的軍團，每個兵開場血就厚。**

### 1.1 100 是回復上限，不是開場值

`sub_1B97E` 的 `cmp [bx+3], 64h` 把**回復**擋在 100
（[`../re/11`](../re/11-tactical-battle.md) §5）。兩件事不衝突：
士氣 200 的軍團開場每個兵 200 點，掉到 100 以下才開始回，
回也只回到 100。

### 1.2 側欄那條體力條讀的是大將

`sub_1C6F6` 用 `word_1D30E:+0x03`（我方）與 `+0x603`（對方），
也就是**每側第 0 個單位**的體力，畫成 `值 × 3 ÷ 4`、上限 124
（[`../re/60`](../re/60-tactical-sidebar.md) §5）。上限 124 對應到值 165，
所以士氣 200 的軍團開場那條是滿的。

## 2. 演算法

```
Deploy(側, 隊, 兵種, 兵數):
    每個兵.體力 ← 該側的軍團士氣          # 不是 100
    大將那一格再被蓋一次：max(70, (武力 × 4 + 50) × 士氣 ÷ 100)
```

⚠ **大將是例外。** `sub_19AF4` 先用 `sub_19B6D` 填滿 48 個兵，
再用 `sub_19B40` 覆蓋第 0 號——那一格的體力跟著武力走
（[`../re/78`](../re/78-soldier-power-from-command.md) §4）。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 士氣搬進戰場 | `internal/state/tactical.go` `beginTactical`：`b.Sides[side].Morale = c.Morale` |
| 開場體力 | `internal/rules/tactical/battle.go` 的 `Side.startHP()`，`Deploy` 與增援路徑共用 |
| 回復上限 | `MaxHP = 100` 不變——它是 `sub_1B97E` 的回復上限，不是開場值 |

士氣 0（測試、無頭模擬沒給）時 `startHP()` 退回 `DefaultPower`，
與 `Power` 同一套預設。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestDeployStartsHPAtMorale`（士氣 200 開場 200 點、士氣 0 退回預設、`MaxHP` 不跟著動）、`TestReinforcementUsesMoraleHP` |
| 對原版 | 同一個局面的 `sb-self` **逐像素相同**、`sb-enemy` 44 px（[`../playtest/40`](../playtest/40-tactical-parity.md) §10）|

## 5. 為什麼它與「被換位的兵不動」綁在一起

開場體力 100 時，第一次挨打就低於 `applyHit` 的退卻閘（`HP < MaxHP`），
第二次就排退卻——受傷的兵很快走光，前線永遠是稀疏的。
士氣 200 的兵要挨滿 100 點才排退卻，於是**會擠在一起**。

擠在一起之後打不打得到人，由 [`62`](62-swapped-unit-skips-its-turn.md) 決定：
攻擊是走進敵人格子撞出來的，一個兵一幀只走一步，所以「輪到自己時
就站在敵人旁邊」是必要條件。少了那一條，前排那一格會一直換人，
48 個兵圍著 2 個打卻一次也打不到。

## 6. ⭐ 大將的體力：兩個來源，斜率對過原版

**2026-09-03 收掉。** 大將那一格的體力有**兩個**扣減來源，remake 兩個都接了：

| 來源 | 原版 | remake | 對過原版了嗎 |
|---|---|---|---|
| **攻城計時器** | `sub_1AE56` 尾段：只在攻城戰、每 10 幀 −1，**扣攻方那一側** | `drainSiegeGeneral`（`SiegeDrainInterval = 10`）| ✅ [`../playtest/52`](../playtest/52-siege-timeseries-parity.md) §4：體力條 100 → 75 px ＝ 值 133 → 100，remake 花 338 幀 ⇒ **10.14 幀／點**，與 `0Ah` 吻合到 1.4%（誤差來自體力條的像素量化）|
| **挨打** | `sub_1B6BC`：傷害 ＝ 攻擊者戰力 ÷ 8，下限留 1（不會死）| `hitGeneral`（`generalDamageShift = 3`、`GeneralMinHP = 1`）| 逐次沒有對拍 |

⚠ **扣哪一側由 `cs:byte_10D35` 的 bit 7 決定**（`< 0x80` 扣 `di=0`、
`≥ 0x80` 扣 `di=0x600`），而那個 bit ＝ **玩家是守方**
（[`56`](56-battlefield-rotation.md)）。原版的 side 0 永遠是玩家，
所以「攻方」在玩家守城時是 side 1。**remake 的 `Sides[0]` 固定是攻方、
`PlayerSide` 另外記**，兩邊的索引語意不同而行為一致——
⭐ 讀這一段時很容易把原版的索引語意直接套到 remake 上，然後把
「寫死 `Sides[0]`」誤判成缺陷。

### 6.1 [`../playtest/40`](../playtest/40-tactical-parity.md) §10 那 60 點怎麼組成的

原版那一格打了 20 秒，大將體力 200 → 140。20 秒 ≈ 182 幀（戰術速度 2
＝ 9.1 fps，[`../re/61`](../re/61-timer-tick-source.md) §4）：

```
計時器  182 ÷ 10          ≈ 18 點
挨打    60 − 18 = 42 點   ≈ 21 次命中（每次 18 >> 3 = 2）
```

20 秒內被打 21 次，與「同一格死了約 115 個兵」的激烈程度相稱。
**兩個來源湊得出來**——所以那 60 點不需要第三個機制。
⚠ 這是量級估算，不是逐幀對拍；要逐幀得讓兩邊的時刻對齊。

## 7. 未解

| 項目 | 現況 |
|---|---|
| ~~`+0x18`（戰力）的算式~~ | **已解並接上**（confirmed，2026-09-02，[`../re/78`](../re/78-soldier-power-from-command.md)、[`115`](115-soldier-power.md)）|
| ~~大將的開場體力~~ | **已解並接上**：`sub_19B40` 把第 0 號兵蓋成 `max(70, (武力 × 4 + 50) × 士氣 ÷ 100)`（[`../re/78`](../re/78-soldier-power-from-command.md) §4）。remake 是 `internal/state/soldierpower.go` 的 `leaderHP`，`TestLeaderHPIsNotMorale` 釘住它不等於士氣 |
| 挨打那一半的逐次對拍 | 只對過量級（§6.1），沒有逐次比。要對得先讓兩邊的時刻對齊——同 [`../playtest/40`](../playtest/40-tactical-parity.md) §13 那一類 |
