# 128 — 隊長倒下不會清掉那一隊的待機兵

## 2026-11-15：應戰單挑與補兵白旗勘誤

輸入DOS/V KI.EXE SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`，dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`。位址為oracle.IDA DOS/V linear，以下為已證實的本案動態與IDA Pro 9.4原始指令證據。
已證實定位：`1A83F／1B437／1B445／1B44D`。隊員Next的下一批傳令與待機保留契約不變。白旗改由reinforceSoldier嘗試補兵時登錄，必須先通過unitDying；待機0或該側已全軍退卻即拒絕補兵，僅隊長更新圖示。退離沒有倒地動畫，該拍補兵入口即可畫白旗。
驗收及限制見[spec159](159-tactical-command-movement-order.md)與[收據](../../workplace/scratch/ui-windows-20261104/cao-duel-accept-integration.json)。

## 2026-11-15：野戰後段傳令、佇列與側欄勘誤

輸入DOS/V KI.EXE SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`，dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`。位址為oracle.IDA DOS/V linear；以下為已證實的正常野戰動態觀測。
已證實定位：`1A426／1A754／1A785／1A83F`。隊長在移動中退離邊界，下一個命令批次才處理隊員退卻，不能在同批立即回頭傳令。七名隊員包含離場記錄都寫Next=5，已在退卻者保留Next；腳本也仍可對離場隊長寫Next。待機數不清零。下方改前／現況對照表屬歷史快照，現行接線以此勘誤與spec159為準。
正式1545／400拍、38張完整圖及正常滑鼠範圍見[spec159](159-tactical-command-movement-order.md)與[收據](../../workplace/scratch/ui-windows-20261104/cao-field-queue-integration.json)。

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
| ~~士氣按比例縮~~ | **已接**（[`129`](129-post-battle-morale-scaling.md)＝CONFORMED）：`postBattleMorale` 就是 `sub_19EBD` ＋ `sub_19F58` 那兩段 |
