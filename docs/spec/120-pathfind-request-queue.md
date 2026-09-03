# 120 — 尋路改成全域佇列：每幀兩個兵，FIFO

**狀態：CONFORMED。** 原版的重算節流是**一條 128 格的環狀佇列、每幀消化兩筆**，
不是每個兵各自的計時器。remake 的 `replanInterval = 30`（明示的 remake 差異）
換成同一個機制，`docs/re/43` 的「重算路徑的時機」跟著關掉。

- 日期：2026-09-03
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_1AED2`（消化端）、
  `sub_1C653`（入隊）、四個入隊點——[`../re/80`](../re/80-pathfind-request-queue.md)
- 推論等級：**confirmed（呼叫端窮舉 ＋ 資料區大小獨立佐證）**
- 相關：[`94`](94-retreat-path-not-cleared-every-frame.md)（繞路點不可以每幀清掉）、
  [`36`](36-ground-planes-and-climbing.md) §1.4（爬的閘）

## 1. 原版做什麼

```
每一幀（sub_1ADC8）：
  1. 大將體力／全軍退卻（sub_1AE56）
  2. ⭐ 消化尋路佇列（sub_1AED2）——最多兩筆
  3. 逐兵移動（sub_1AF69）

入隊（sub_1C653）：兵記錄 +0x00 bit 4 已設就不重複排；
                  否則設 bit 4，把兵記錄位址寫進尾端，尾游標 +2（環狀 128 格）
出隊（sub_1AED2）：清 bit 4、繞路點游標 +0x16 歸零、設起終點、
                  算一次波前擴散，成功就把點數寫 +0x17 並取第一個點
```

四個入隊點（[`../re/80`](../re/80-pathfind-request-queue.md) §3）：
**下退卻命令**、**三軸都走不動而且下一個繞路點也沒用**、
**目標格地形不通**、**碰撞處理的尾段**。

⭐ **原版預設不尋路**——平常直接朝目標三軸試走，走不動才要一條路。

## 2. remake 改之前是什麼

`replan()` 在「這一幀走不動或撞到地形」時被呼叫，內部用
`b.Frame - s.PathAt < replanInterval`（30 幀）擋掉太頻繁的重算。
形狀與原版不同：

| | 原版 | remake（改之前）|
|---|---|---|
| 預算 | **全域**每幀 2 次 | 每個兵每 30 幀 1 次 |
| 順序 | FIFO | 無序 |
| 最壞情況 | 96 個兵全卡住 ⇒ 每個兵約 48 幀輪到一次 | 同一幀最多 96 次波前擴散 |

## 3. remake 怎麼改

| 項目 | 位置 |
|---|---|
| 佇列本體 | ✅ `internal/rules/tactical/pathqueue.go`：`pathQueue`，128 格環狀、`Queued` 旗標去重 |
| 每幀消化兩筆 | ✅ `Battle.Step()` 在逐兵迴圈**之前**呼叫 `drainPathQueue()`（與 `sub_1ADC8` 同序）|
| 入隊 | ✅ `requestPath(side, k)`：退卻下令、走不動、撞地形、碰撞尾段四處 |
| 出隊做的事 | ✅ 繞路點游標歸零（`s.Path = nil` 之後重建）、算一次 `FindPath` |
| 旗標 | ✅ `Soldier.PathQueued`（對應兵記錄 `+0x00` bit 4）|

⚠ **`PathAt` 與 `replanInterval` 拿掉了。** `docs/spec/94` 的
「繞路點不可以每幀清掉」仍然成立，而且變得更自然：現在清不清由佇列決定，
不由呼叫端決定。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestPathQueueDedupes`（同一個兵排兩次只佔一格）、`TestPathQueueDrainsTwoPerFrame`（三筆要兩幀）、`TestPathQueueIsFIFO`、`TestPathQueueNeverOverflows`（96 個兵全排也塞不滿 128） |
| 迴歸 ✅ | 三場都仍然打得完，幀數跟著換節流而變（**那是預期的**，不是回歸）：`TestNormalScenarioTacticalBattleTerminates` 967 → **911**、`TestSiegeFixtureTerminates` 1,192 → **1,339**、`TestFieldBattleTerminates` **284** |
| 對原版 ⚠ | **對拍看不出來**（使用者裁定 2026-09-03）：尋路節流只改變兵走到同一格的**時刻**，單張截圖上分不出「這一幀誰在重算」。判準是機器碼，不是像素。實測也是這樣——攻城第 160 拍 `field` 0.85% → **0.99%**、`sb-minimap` 1.64% → **1.72%**，四個 PASS 區仍然 0 px、`bottom` 2 px、`sb-enemy` 12 px。**位置差，不是畫壞** |

## 5. 未解

| 項目 | 現況 |
|---|---|
| `sub_1ACA4` | 排隊之後緊接著呼叫，內容沒讀（[`../re/80`](../re/80-pathfind-request-queue.md) §5）|
| 碰撞尾段那一處的前提 | 哪幾條分支會走到 `loc_1B612` 沒有逐條追 |
| 佇列順序對戰局的影響 | FIFO 與無序在同一場攻城裡差多少，沒有量化 |
