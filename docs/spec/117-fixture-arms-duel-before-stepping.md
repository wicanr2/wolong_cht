# 117 — 驗收捷徑要先武裝開場喊話再推戰場

**狀態：CONFORMED。** `-battle-steps` 那條驗收路徑先推完 N 個戰場 tick
才呼叫 `startBattleTalk`，於是單挑狀態機在整段開場期間都沒有輸入——
挑戰喊話一次都不會產生。正常迴圈（`updateBattle`）是**先武裝再推**，
兩條路的順序對調過來就好。

- 日期：2026-09-02
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_1A1C5`（開場 50 tick
  ＋ 單挑）——[`80`](80-duel-opening.md)、[`../re/74`](../re/74-battle-opening-duel.md)
- 推論等級：**confirmed（決定性重現：16 個 seed 逐一量到同一個數字）**
- 相關：[`91`](91-tactical-parity.md) §6.1（另一個「旋鈕拿錯」的坑）、
  [`105`](105-encounter-goes-straight-to-battle.md)（遭遇直接進戰場）

## 1. 症狀

野戰的同狀態對拍（[`../playtest/43`](../playtest/43-field-battle-parity.md)）
原本九區裡七區 0 px，`field` 0.05%。用 `-open-battle -siege-corps 39,35
-battle-steps 52` 重開同一場，八區依舊 0 px 或 40 px，**而 `field` 是 11.24%**
——差的整塊是呂布的挑戰對白框，remake 那一張根本沒有框。

⭐ **判準是 seed 掃描**：挑戰側由氣勢比較決定，而氣勢平手時靠亂數尾
（[`80`](80-duel-opening.md) §6），所以換 seed 應該會換人挑戰、
數字必然跟著動。實測 seed 0–15 **十六個值一模一樣**：

```
seed 0   field 11.24%     seed 8   field 11.24%
…（中間十四個同值）…
seed 7   field 11.24%     seed 15  field 11.24%
```

一個「吃亂數」的機制對亂數完全免疫，只有一種解釋：**它沒有被啟動。**

## 2. 成因

`cmd/wlgame/battle.go` 的 `stageEncounter` 原本是：

```go
for i := 0; i < steps; i++ { p.Battle.Step() }   // ← 先推
g.startBattleTalk(p)                             // ← 才武裝
g.tickBattleTalk(steps)
```

`startBattleTalk` 裡才呼叫 `SetDuelInput`。單挑的開場只有 50 tick
（`duelOpeningTicks`），`-battle-steps 52` 早就跑完了，武裝時
`b.duel.armed` 才變 true、`timer` 才歸位——**該喊話的那一刻已經過去。**

正常迴圈沒有這個問題：`updateBattle` 第一行就是 `g.startBattleTalk(p)`，
每一幀推進之前都武裝一次（[`105`](105-encounter-goes-straight-to-battle.md) §3
那一列「第一幀武裝開戰喊話」）。

## 3. 改法

把驗收路徑改成與正常迴圈**同一個順序、同一個節拍**：

```go
g.startBattleTalk(p)              // 先武裝
for i := 0; i < steps; i++ {
    p.Battle.Step()
    g.pumpDuelTalks(p)            // 這一 tick 產生的喊話當下掛框
    g.tickBattleTalk(1)           // 對白的時鐘跟著走
}
```

⭐ **一次推完再補時鐘不等於逐拍推。** 對白的壽命是 60 tick
（[`60`](60-battle-talk-duration.md)），一次 `tickBattleTalk(steps)`
在總量上湊得出來，但「第 50 tick 產生、第 52 tick 還掛著」這種
中途狀態湊不出來——而對拍要比的正是中途狀態。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 順序與節拍 | `cmd/wlgame/battle.go` 的 `stageEncounter` |
| 迴歸 | `cmd/wlgame/stage_encounter_test.go` 的 `TestStageEncounterArmsDuelBeforeStepping` |

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestStageEncounterArmsDuelBeforeStepping`：同一條 fixture 推 0 拍時 `DuelActive()` 必須 false（開場 50 tick 還沒走完）、推 52 拍必須 true。**兩端都斷言**，少了前半就分不出「真的武裝了」與「永遠回 true」|
| 突變測試 | 把順序改回「先推再武裝」，那一支立刻紅：`第 52 拍 DuelActive() = false，want true` |
| 對原版 | [`../playtest/58`](../playtest/58-parity-retest-20260902.md)：野戰九區 |

## 6. ⭐ 教訓

**驗收捷徑與正常迴圈的差別，會以「被驗收的東西不見了」的形式出現。**
這是 [`../playtest/49`](../playtest/49-parity-retest-20260827.md) §2 那一條
（`-shot`把音效關掉，於是那一格印「未接入」）的同一族：
捷徑省掉的不是無關的雜項，而是**被驗收畫面的一部分**。

判準便宜且通用：**拿一個「應該會改變結果」的旋鈕轉一圈**。
數字紋風不動就表示那條路沒有被走到——不必先知道成因在哪。

## 7. 未解

| 項目 | 現況 |
|---|---|
| 野戰 `field` 的殘差 | 修好之後還剩多少，量在 [`../playtest/58`](../playtest/58-parity-retest-20260902.md) |
| 自然流程那條路 | 遭遇訊息會擋住 `-shot-frames`（[`105`](105-encounter-goes-straight-to-battle.md) §4 已寫明是預期行為）。要用自然流程做野戰對拍，得有一個「訊息自動按掉」的驗收旗標 |
