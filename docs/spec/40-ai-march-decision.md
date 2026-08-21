# 40 — 電腦勢力的行軍決策鏈（Stage 0–3／8／10）

**狀態：CONFORMED。** 四支 AI handler 的轉移條件全部有機器碼出處，
remake 已實作並有單測。

- 日期：2026-08-17
- 出處：[`../re/65`](../re/65-ai-march-decision-chain.md)（四支 handler 逐行）、
  [`../re/64`](../re/64-corps-arrival-state-machine.md)（分派表、Stage 8–11）、
  [`../re/44`](../re/44-threat-and-reinforcement-ai.md)（勢力 `+0x16`／`+0x17` 誰寫）、
  [`../formats/08`](../formats/08-sinario-save.md) §1.5–§1.7（欄位）
- 推論等級：**confirmed（靜態）**，除了 §3 標明的一處

## 1. 為什麼這件事非做不可

remake 一直有兩個空轉的欄位：勢力的 `ReliefSite`（求援）與 `LostSite`（失土）
**有人寫、沒人讀**。原版讀它們的就是 AI 軍團的 Stage 2——
軍團走到據點、發現沒事做，就去把這兩格待辦領走。

沒有這條鏈，AI 的軍團編出來、走到第一個目標，然後就停在那裡不動了。

## 2. 狀態轉移

`arriveCorps`（軍團停在目標據點上時，每個 tick 跑一次）依 `Stage` 分派。
**Stage ≥ 8 不分玩家**，0–3 才分。

| Stage | 名稱 | 條件 → 下一個 Stage |
|---:|---|---|
| 0 | 移動中 | 目標據點的「威脅有具體目標」旗標亮 → 留在 0（駐守）；否則 → 1 |
| 1 | 在據點決策 | 現在不是據點 → 0；兵力 ≤ 300 → 10；旗標亮 → 0；**不受威脅 或 這一格軍團數 > 2** → 2（延遲 1–8 tick）；否則兵力 < 600 且在首都 → 9 |
| 2 | 挑目標 | 旗標亮 → 1；受威脅且這一格只有一支 → 1；領到待辦 → 0；領不到且不受威脅 → 11 |
| 3 | 補完兵的體檢 | 六槽任一 < 30 點（300 人）→ 11；否則 → 8 |
| 8 | 等士氣 | 士氣 < 勢力基準 → 留著；達標 → 1 |
| 9 | 首都補兵 | 退回 ＋ 重分配 → 3 |
| 10 | 走回首都 | 目標校正成首都；到了 → 9 |
| 11 | 解體 | 目標校正成首都；到了 → 解散 |

### 2.1 挑目標的優先序

| 順序 | 來源 | 誰寫 |
|---:|---|---|
| 1 | 勢力 `LostSite`（`+0x17`，被打下來的據點）| 據點換手時 |
| 2 | 勢力 `ReliefSite`（`+0x16`，最近求援的據點）| 據點求援時 |

**兩格都是「取走就清空」的一格佇列**，所以一件待辦只會派出一支軍團。
**資金吃緊（`LowFunds`）的勢力跳過第一項**——沒錢就不主動反攻失土。

`LowFunds` 是勢力記錄 `+0x00` 的位元 6，remake 早就在每小時的財政檢查裡
算出來了（`internal/rules/diplomacy` 的 `CanSustainInvasion`），
缺的只是這個消費端。

## 3. ⚠ 一處與原版不同

`sub_143AF` 的「這一格軍團數 > 2」那個比較，原版的 `di` 沒設就用，
實際讀到的位址比據點記錄少 `0x840`
（[`../re/65`](../re/65-ai-march-decision-chain.md) §3.2）。

**remake 實作作者意圖的版本**（讀據點的 `Occupancy`）。
乾淨重寫裡沒有那個錯位的位址可以讀，所以「照抄」不成立。
**這是明確標記的 remake 差異。**

## 4. remake 要改哪裡

| 項目 | 作法 |
|---|---|
| `City.Threatened`／`City.Specific` | 據點記錄 `+0x00` 的位元 7／6。`tickCity` 掃完威脅就寫，存檔 round-trip 也帶上（先前只存低 4 位的 `Adjacency`）|
| `internal/state/aimarch.go` | 四支 handler ＋ Stage 8／10 |
| `arriveCorps` | 先分 Stage ≥ 8，再分玩家／非玩家 |
| `Stage` 常數 | 補 `StageWaitMorale = 8`、`StageHomeResupply = 10` |
| 延遲 | Stage 1 → 2 時 `Timer = 亂數(0–7) + 1`；remake 用既有的 `combat.Rand` |
| `Occupancy` 的即時扣減 | 原版 `dec [di+18h]`。remake 照做——`Occupancy` 每 tick 從位置重算，扣減只在同一個 tick 內生效，語意相同 |
| 勢力滅亡的收尾 | `disperseFaction`：俘虜釋放、有職務的變在野、君主與無職成為勝方的俘虜，**軍團全部消失**（`../re/59` §4.1）|
| 武將數的維護 | 被俘／戰死 `−1`、釋放 `+1`，**俘虜期間不算在任何一方**（`../re/59` §4.2）。不變量層跟著改成不數俘虜 |

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestAIStage0StaysWhileThreatIsSpecific`：旗標亮就不換檔 |
| 單元測試 | `TestAIStage1SendsWeakCorpsHome`：兵力 ≤ 300 → Stage 10 |
| 單元測試 | `TestAIStage2TakesLostSiteFirst`：先領失土再領求援，兩格都被清空 |
| 單元測試 | `TestAIStage2SkipsLostSiteWhenLowFunds`：資金吃緊時只領求援 |
| 單元測試 | `TestAIStage2DisbandsWithNothingToDo`：沒待辦又不受威脅 → Stage 11 |
| 單元測試 | `TestAIStage3DisbandsUnderstrengthCorps`：任一槽 < 30 → 11，否則 → 8 |
| 單元測試 | `TestAIStage8WaitsForMorale`：士氣達標才回 Stage 1 |
| 單元測試 | `TestAICorpsLifecycleReachesRelief`：整條鏈跑到領走求援並出發 |
| 單元測試 | `TestEliminatedFactionLeavesNoCorps`：滅亡的勢力不留軍團、武將數歸零 |

## 6. 未解

（無。原先掛在這裡的兩條都已收掉：`sub_147BB` 的 `0x8000` 分支由
[`43`](43-rout-on-blocked-return.md) 完整實作並驗過，`sub_1487B` 的
演算法在 [`46`](46-post-battle-retreat.md) §2、消費端在
[`47`](47-city-fall-corps-redirect.md) §3。）

<!-- 缺口：無 -->
