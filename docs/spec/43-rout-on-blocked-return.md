# 43 — 回不了家的軍團會敗走

**狀態：CONFORMED。** 觸發條件、立即效果與 48 tick 之後的收尾都有機器碼出處，
remake 已實作並有單測。

- 日期：2026-08-17
- 出處：[`../re/65`](../re/65-ai-march-decision-chain.md) §8
  （`sub_147BB` 的 `0x8000` 分支、`loc_1491B` 的成本、`sub_12977`、`sub_12A7E`）
- 推論等級：**強證據**（成本與高位元那六行是手工解碼原始 bytes）

## 1. 原版的規則

**觸發**：軍團要移動時重算下一步，而

1. 這條路要**穿過非己方的據點**（`loc_1491B` 每碰到一個
   `+0x01` 不等於自己的據點就 `add dx, 0A6h` ＋ `or dh, 80h`），**而且**
2. `Stage ≥ 10`（回首都補兵的路上，或解體的路上）。

野外節點（編號 ≥ 192）沒有歸屬，不算。

**立即**：佔用圖 −1、勢力的軍團數 −1、存在旗標寫 `8`、`+0x03` 寫 `48`、
大地圖上不再顯示；玩家或對手的軍團會跳 TALK #1F／#20。

**48 個 tick 之後**（`sub_12A7E`，每 tick 減 1）：軍團記錄歸零、
主將 `+0x17` 職務歸 0；原勢力已滅則主將改成在野。玩家的軍團跳 TALK #23 ＋ #198。

⭐ **兵員不回預備兵池。** 這是它與「解體」（[`../re/64`](../re/64-corps-arrival-state-machine.md) §3）
最大的差別——解體是回收，敗走是損失。

## 2. remake 怎麼做

| 項目 | 作法 |
|---|---|
| 狀態 | `Corps.Routing bool` ＋ `Corps.RoutTimer int`（軍團 `+0x00` 的旗標 8 與 `+0x03`）|
| 判定 | `World.returnBlocked(i)`：走一次 `w.routes[i]`，路上任何一個據點的 `Owner` 不是自己就成立 |
| 觸發點 | `headHomeResupply` 與 `arriveDisband` 重下行軍之後檢查一次 |
| 進入 | `routCorps(i)`：`Alive=false`、`Routing=true`、`RoutTimer=48`、勢力軍團數 −1。**兵不回池** |
| 收尾 | `tickCorps` 每 tick 對 `Routing` 的軍團減 1；歸零時清乾淨並把主將解職 |
| 存檔 | 旗標 8 與 `+0x03` 都寫回；載入時 `r[0]<0x80 && r[0]&8 != 0` ⇒ `Routing` |

**remake 差異**：原版是在**每一步移動**時判定，remake 在**重下行軍令時**判定。
兩者差在「路上情勢變了」的情況——remake 的路徑是一次算完的
（[`../re/64`](../re/64-corps-arrival-state-machine.md) §6），沒有逐步重算的地方可以掛。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestRoutedCorpsDisappearsAfterTimer`：48 tick 之後軍團歸零、主將解職 |
| 單元測試 | `TestRoutLosesTheMen`：預備兵池沒有增加（對照解體會增加）|
| 單元測試 | `TestRoutingSurvivesSaveRoundTrip`：旗標 8 與計時器 byte-for-byte 寫得回去 |
| 單元測試 | `TestReturnBlockedNeedsForeignCityOnTheRoute`：路上有別人的據點才成立 |

| 長跑 | `cmd/wlsim` 現在會自己掛道路圖（`-map`，預設 `MMAP.MAP`），5 年 60 個月跑完不變量不違反 |

> ⚠ **這一條的判定需要道路圖。** `returnBlocked` 走的是算好的格子路徑，
> 而 `wlsim` 先前沒有 `SetRoads`（規則層不讀檔案，道路圖由呼叫端注入），
> 路徑永遠是空的、判定永遠不成立——接上規則前後跑同一顆種子，
> 結果**逐字相同**。那證明的是「沒有回歸」，不是「規則有效」。
> 現在 `wlsim` 自己載 `MMAP.MAP`（253 條路），走的才是原版的道路。

## 4. 未解

| 項目 | 現況 |
|---|---|
| `loc_1491B` 的其他成本項 | 只解出「非己方據點 ＋0xA6」，廣度優先搜尋本身沒逐條讀 |
| TALK #1F／#20／#23／#198 | remake 還沒把這四則接上 |
