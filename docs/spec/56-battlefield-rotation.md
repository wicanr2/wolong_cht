# 56 — 戰場轉 180 度：什麼時候轉、轉的時候圖塊值要換

**狀態：CONFORMED。三段算式都接上了，並用原版的許昌攻防戰驗過：
`field` 區 87.8% → 46.1%、小地圖 41.5% → 13.4%
（[`../playtest/40`](../playtest/40-tactical-parity.md) §4）。**

- 日期：2026-08-17
- 出處：[`../formats/07`](../formats/07-battle.md) §2.2／§2.3（`sub_1CAEB`
  的旗標檢查、`sub_1CB9B` 的對調迴圈、`sub_1CBBC` 的值對映）、
  [`../re/58`](../re/58-bgm-scene-mapping.md) §4（`sub_14ED7` 誰設旗標）、
  [`../re/11`](../re/11-tactical-battle.md) §4.4（野戰的配對順序）、§4.5（`63 − v`）
- 推論等級：**confirmed（靜態）**；「玩家守城會轉」另有實機證據
  （[`../playtest/40`](../playtest/40-tactical-parity.md) §4：把 remake 的小地圖
  轉 180 度再比，相同像素 48.1% → 83.0%）

## 1. 原版做什麼

`byte_10D35` 的 **bit 6 ＝ 這一場的戰場要轉 180 度**。載入器 `sub_1CAEB`
在讀完 4,096 B 之後看這個位元，設著就轉。

誰設它，兩條路各一：

| 情形 | 誰設 | 條件 |
|---|---|---|
| **攻城** | `sub_14ED7` | **玩家守城**那一支 `or cs:byte_10D35, 0C0h`——bit 7 攻守對調 ＋ bit 6 翻轉。玩家攻城那一支不設（`sub_14ADE` 進來前先寫 0）|
| **野戰** | `sub_14BDD` | 兩格地形**換過順序才配上**那 21 筆表時 `xor dh, 40h` |

⭐ **攻城的判準是「玩家在哪一邊」，不是「誰攻誰守」。**
兩個位元一起設，所以玩家守城時側欄換邊與戰場翻轉是同一件事的兩面。

## 2. 演算法

三段，缺一不可：

```
① 地形格整個倒過來（sub_1CB9B）
   for i in 0 .. 1983:            ; cx = 7C0h
       swap cells[0x40 + i], cells[0xFBF − i]
   ＝ 3,968 格（64 × 62）首尾對調 ＝ 轉 180 度

② 每一格的值換成鏡射版（sub_1CBBC）
   v < 0x30            → 不變
   0x30 ≤ v ≤ 0xCF     → ((v − 0x30) ^ 0x10) + 0x30      ; 每 0x10 一組配對
   0xD0 ≤ v ≤ 0xEF     → (v & 3) ∈ {0,3} → v ^= 3；{1,2} → 不變
   v ≥ 0xF0            → v ^= 1

③ 城門的 X（BATTLE.MAP 索引第二欄）
   v ≠ 0 → v = 0x3F − v                                  ; 0x3F ＝ 寬 − 1
```

⚠ **只做 ① 會得到「地形對了但斜坡河岸都反向」的戰場**——
有方向性的圖塊要跟著鏡射，那正是 ② 那三段規則在做的事。
② 分三段代表圖塊表在這三個值域用了三種排列慣例，**照抄，不要合併**。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| ①② 純函式 | `internal/assets/battle`：`Rotate180(cells [][]byte)`、`RotateTile(v byte)` |
| ③ | `internal/assets/battle`：`RotateGateX(x int)` |
| 什麼時候轉 | `internal/state`：`beginTactical` 算出「玩家是不是守方」傳給 `TacticalSetup.Field` 回呼；野戰用 `battlefield.Select` 回的 `rotate` |
| 套用點 | `cmd/wlgame/battle.go` 的 `buildField`（規則層的地形）與 `newBattleView`（小地圖與繪圖） |

**兩個套用點要用同一份轉好的格子**，否則畫面與規則層會不一致
（兵走的路是規則層的、看起來的地形是繪圖層的）。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestRotateTileRanges`：三段的邊界值 ＋ 256 個值轉兩次都回到原樣 |
| 單元測試 ✅ | `TestRotate180IsInvolution`：整張轉兩次回原樣，而且**轉了跟沒轉不一樣**（正對照）|
| 單元測試 ✅ | `TestRotateKeepsWallAndGateCounts`：214 張逐張轉，城壁與門的格數不變——這是值對映表的正對照，換錯就對不上 |
| 對原版 ✅ | [`../playtest/40`](../playtest/40-tactical-parity.md) §4。另外攻方那一張改動前後**逐像素相同**，證明沒有動到不該轉的那一半 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 表頭與尾段那各 64 byte | 轉的時候原版**不動它們**（迴圈只掃 `0x40`–`0xFBF`）。內容仍未解 |
| 兵的初始位置 | 翻轉之後雙方的佈陣點怎麼算，還沒對過 |
| 鏡頭差一個等角格 | 翻轉之後戰場區還差 (−16, −8)（`../playtest/40` §4.1）。小地圖沒有位移，所以不是翻轉中心的問題 |
