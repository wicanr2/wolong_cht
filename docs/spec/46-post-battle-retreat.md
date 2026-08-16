# 46 — 戰後敗方退一站回家，退不了就壞滅

**狀態：CONFORMED。** `sub_1474A` 的後半段與 `sub_1487B` 都有機器碼出處，
remake 已實作並有單測。
⭐ **壞滅有第三個入口**：士氣 0、大將槽 0 之外，還有「找不到退路」。

- 日期：2026-08-17
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 分析時的 `KI.EXE.i64` SHA-256：`8a8fd7d528e0498000fd04282300588b637fbcb8aa48deb09242f1f41f532691`
- 出處：[`../re/09`](../re/09-combat.md) §5（`sub_1474A`）、
  [`../re/65`](../re/65-ai-march-decision-chain.md) §8.4（`sub_1487B`）
- 推論等級：**confirmed**（兩支都逐行讀過）

## 1. `sub_1474A` 的後半段

戰鬥結束後**兩邊各跑一次**（`sub_15130` 呼叫兩次，結果進 `ah` 的
bit 0／bit 1）。`cl` 是勝敗旗標：

```
士氣 0 或大將槽兵力 0            → 壞滅
cl == 0（勝方）                  → 原地不動，Stage 8
站在自家據點上                   → 原地不動，Stage 8
否則 sub_1487B 找退路
  找不到                         → ★ 壞滅
  找到 next                      → 目標改成 next、設「下一步要重算」
      兵力 ≤ 300 或 next ＝ 首都 → Stage 10（回首都補兵）
      否則                       → Stage 8（等士氣）
```

⭐ **敗方一次只退一站**，而且那一站必須是自己的地。
「退不了」與「被打光」在原版是同一個結局。

## 2. `sub_1487B`：回家路上的下一站

```
capital = 勢力記錄 +0x03；＝ 0xFF ⇒ 失敗（沒有首都就無處可退）
現在在野外（節點×8 ≥ 0x800）：
    先從這一格的兩個鄰接槽（+6／+8）裡挑一個**屬於自己**的當起點
    兩個都不是自己的 ⇒ 失敗
loc_1491B 廣度優先找路（往首都）
    回 carry ⇒ next ＝ 首都本身，直接成功（不再檢查歸屬）
    否則方向（±4）決定讀 +6 還是 +8 那一槽 ⇒ next
next 不是自己的據點 ⇒ 失敗
```

**目標是首都，不是攻擊目標。** 兩個歸屬檢查各擋一種情況：
出發那一格的鄰接要是自己的、下一站也要是自己的。

## 3. 這支不是 Stage 10 用的

`sub_1487B` 的呼叫者只有兩個：`sub_1474A`（戰後）與
`sub_14DA4`（據點失守時把受影響的軍團調頭）。

**Stage 10／11 走的是另一條路**：`sub_144A9`／`sub_144D6` 直接
`mov [si+20h], bh`（bh ＝ 首都）＋ `or byte ptr [si], 2`，
一次把目標設成首都，不逐站。remake 的 `headHomeResupply`／
`arriveDisband` 與這兩支一致。

## 4. remake 實作

| 項目 | 作法 |
|---|---|
| 找下一站 | `World.nextHopHome(i)`：`roads.Route(現在節點, 首都)` 的第 2 個節點；沒有首都、沒有下一站或下一站不是自己的 ⇒ −1 |
| 廣度優先回 carry | `Route` 回 nil 或長度 < 2 ⇒ 直接回首都（對應 `loc_1490C`）|
| 戰後處理 | `World.retreatOrPerish(i, won)`：回傳 true 表示「退不了 ⇒ 壞滅」，並在成功時 `March` 到下一站 ＋ 設 Stage |
| 接線 | `resolveCorpsBattle`／`fightGarrison`／戰術結算三處：`destroyed = 士氣判定 \|\| retreatOrPerish(...)`，再交給既有的 `afterBattle` |
| 沒有道路圖時 | `nextHopHome` 退回首都，等於一次走完。**規則層不讀檔案**，道路圖由呼叫端注入（`SetRoads`）|

⚠ **攻城時守方站在自己的城裡**，走「原地不動」那一支，
所以攻城的勝負與據點易主完全不受這條規則影響。

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestNextHopHomeStopsAtForeignGround`：下一站是別人的地 ⇒ −1 |
| 單元測試 | `TestNextHopHomeWithoutCapital`：沒有首都 ⇒ −1 |
| 單元測試 | `TestLoserRetreatsOneHop`：敗方目標變成路徑上的第一站，不是首都 |
| 單元測試 | `TestWinnerStandsStill`：勝方不動、Stage ＝ 8 |
| 單元測試 | `TestDefenderInOwnCityDoesNotRetreat`：守在自家城裡不退，攻城結果不變 |
| 單元測試 | `TestNoRetreatMeansDestroyed`：退不了的敗方走壞滅那條路 |
| 長跑 | `cmd/wlsim` 5 年 60 個月，不變量不違反 |

## 6. 未解

| 項目 | 現況 |
|---|---|
| `loc_1491B` 的方向回傳 | `±4` 決定讀哪一個鄰接槽，remake 用 `Route` 的第 2 個節點取代，沒有逐條對過兩者選的是不是同一站 |
| 野外那一格的鄰接槽 | remake 的 `Node` 在行軍途中停在上一個據點，所以走的是「從上一個據點找路」，不是原版的「從野外那一格的鄰接槽挑」|
| `sub_14DA4` | 據點失守時把受影響的軍團調頭，還沒接 |
