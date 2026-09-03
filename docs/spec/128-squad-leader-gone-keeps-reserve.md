# 128 — 隊長倒下不會清掉那一隊的待機兵

**狀態：CONFORMED。** remake 的 `squadLeaderGone` 在隊長倒下時把
`Reserve[squad]` 設成 0，**原版沒有這一步**——碰結算緩衝區的 11 支函式
逐一看完，待機數只有 `dec`（開場取用、補兵），沒有任何地方寫 0
（[`../re/83`](../re/83-post-battle-troop-accounting.md) §1）。
方向是 **remake 少算兵力**：那些兵在原版會補進場、立刻退卻、走出畫面，
被算成生還（`../re/83` §4）。

- 日期：2026-09-03
- 出處：`sub_1A83F`（隊長不在時只改命令）、`sub_1B413`（補兵）、
  `sub_19F58`（戰後三項相加）、`sub_1A754`／`sub_1A785`（每幀施加）
- 推論等級：**confirmed（靜態）**——負證據那一半靠「候選集合封閉」
  （只有 11 支碰得到那塊緩衝區）
- 相關：[`65`](65-retreated-soldiers-survive.md)（退卻是保命不是損失）、
  [`../re/83`](../re/83-post-battle-troop-accounting.md)、
  [`../re/11`](../re/11-tactical-battle.md) §5.6

## 1. 原版與 remake 的差異

| | 原版 | remake（現況）|
|---|---|---|
| 隊長倒下時 | `sub_1A83F`：七名隊員的命令改成 5（退卻）。**待機數不動** | `squadLeaderGone`：七名排退卻 **＋ `Reserve[squad] = 0`** |
| 施加時機 | **每一幀**（`sub_1A754` 逐隊看隊長的 bit 7）| 一次性（`applyHit` 的死亡路徑）|
| 之後補兵 | 照補（`sub_1B413` 不看隊長在不在），補進來下一幀又被改成退卻 | 不補（`Reserve` 已經是 0）|
| 戰後那些兵 | 走「補進場 → 退卻 → 離場」，`ah = 0` ⇒ **算生還** | **從帳上消失** |

⭐ 兩個差異看起來會互相抵消——「不補」配「一次性」，行為都是那些兵不上場。
**但帳目不一樣**：原版最後把它們算進兵力，remake 沒有。

## 2. 要改什麼

| 項目 | 位置 | 改法 |
|---|---|---|
| 別清待機 | `internal/rules/tactical/soldier.go` `squadLeaderGone` | 刪掉 `s.Reserve[squad] = 0` |
| 每幀施加 | 同檔 ＋ `battle.go` 的 `Step()` | 逐隊檢查隊長的 `Alive`，不在就對在場的七名排退卻——與原版 `sub_1A754` 同一個位置（勝負判定那一段） |
| 補兵不受影響 | `battle.go` `reinforce` | 不動。⚠ 原版「補不補隊長格」還沒對（[`../re/83`](../re/83-post-battle-troop-accounting.md) §5），remake 維持不補 |

⚠ **`squadLeaderGoneFor` 那條指標回查的路徑要一起看**：改成每幀掃之後，
死亡當下不必再呼叫它，但**投射物命中**的路徑（`projectileHit` → `applyHit`）
也走同一支，所以留著不會錯——只是變成多餘。先留，等每幀版本驗過再拆。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestSquadLeaderGoneKeepsReserve`：隊長倒下之後 `Reserve[squad]` 不變 |
| 單元測試 ✅ | `TestLeaderLossRetreatsSquad`：七名排退卻、待機數不動 |
| 單元測試 ✅ | `TestSquadLeaderGoneReappliesEveryFrame`：補進場的兵也被排退卻 |
| 單元測試 ✅ | **`TestStepReappliesSquadLeaderGone`：接線本身**——前三支都直接呼叫 `applySquadLeaderGone`，把 `Step()` 裡那一行拔掉照樣綠（2026-09-03 的突變測試發現），所以這一支非有不可 |
| 突變測試 ✅ | 拔掉 `Step()` 裡那一行 ⇒ 只有 `TestStepReappliesSquadLeaderGone` 紅 |
| 對原版 ✅ | 野戰 ＋ 攻城四個取樣點重跑，**改前改後每個數字都相同**（[`../playtest/62`](../playtest/62-parity-retest-20260903.md)）|

⚠ **對拍證明的是「沒有回歸」，不是「改對了」**——四個取樣點裡沒有一個
走到「隊長倒下且那一隊還有待機兵」的狀態（[`../playtest/62`](../playtest/62-parity-retest-20260903.md) §3）。
改動本身的驗收靠上面那四支單元測試。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 戰後兵力的逐槽對拍 | 原版打完之後每槽兵數是三項相加（[`../re/83`](../re/83-post-battle-troop-accounting.md) §3），remake 的戰後回填**沒有逐槽比過原版** |
| 士氣按比例縮 | `sub_19F58` 最後三行：新士氣 ＝ 舊士氣 × 新總兵力 ÷ 舊總兵力。remake 的戰後士氣處理**還沒對照這一條** |
