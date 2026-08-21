# 49 — 進言的第四、五項：遷都與請求君主出陣

**狀態：CONFORMED。** 兩項的判定條件、TALK 起點與畫面流程都有機器碼出處，
remake 已接進畫面並有單測。
⭐ **這兩項沒有說服迴圈**——君主看一眼就定案，不問理由、不動信賴度。

- 日期：2026-08-17
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 分析時的 `KI.EXE.i64` SHA-256：`8a8fd7d528e0498000fd04282300588b637fbcb8aa48deb09242f1f41f532691`
- 出處：[`../mechanics/70-ai.md`](../mechanics/70-ai.md)（`sub_16909`／`sub_1699E` 的兩組閘）、
  `sub_13B08`（共用的對話）、`sub_16E8F`（出陣的編成）
- 推論等級：**confirmed**

## 1. 共用的對話：`sub_13B08`

兩項都不談判，所以走的不是 `sub_13830`（那一支後面接說服迴圈），
而是只有三句的 `sub_13B08(cx ＝ 起點, al ＝ 結果碼)`：

```asm
00013B21  sub_13D09(al=0)      ; IVENTGRF 第 0 頁
00013B29  sub_13C99            ; ① 上框：君主開場（cx，變體由 sub_13C99 自己加）
00013B2C  add cx, 3 / sub_13CDC ; ② 下框：軍師（cx+3，不加變體）
00013B32  inc cx               ; cx+4
00013B33  and bp, bp / jz →
00013B37  add cx, 3            ; 結果碼 ≠ 0 ⇒ cx+7
00013B3A  sub_13C99            ; ③ 上框：君主定案
```

| 項目 | 起點 | ① 開場 | ② 軍師 | ③ 接受 | ③ 拒絕 |
|---|---:|---|---|---|---|
| 遷都（`sub_16909`）| **386** | #386–388 | #389 | #390–392 | #393–395 |
| 請求出陣（`sub_1699E`）| **396** | #396–398 | #399 | #400–402 | #403–405 |

原版實錄影片 4 分 30 秒那一幕就是後者：
上框「軍師，是要我出陣嗎？」（#397）→ 下框「若請主公出陣，將士的士氣也會提高吧。」（#399）
→ 上框「好，來人啊，牽馬過來！我來大顯神勇！！」（#401）。

> ⚠ **看到君主說話不代表提議被接受。** ①② 兩句無論通不通過都會演，
> 差別只在第三句。

## 2. 遷都：兩個條件缺一不可

```
from.Kind >= to.Kind             類型編號越小城越大 ⇒ 不能往小的搬
to.Production > from.Production  嚴格大於，相等也不行
```

⚠ **只看生產力會誤判**——生產力更高的小都市仍然會被拒絕。

搬成之後要跑 `sub_14502`：目標還掛著舊首都的軍團一律改掛新首都。

## 3. 請求出陣：兩道閘，然後君主自己編一支軍團

| 閘 | 條件 |
|---|---|
| 國庫 | 資金 ≥ **(15 − 好戰等級) × 1,024**。好戰 15 的門檻是 0，好戰 0 的要 15,360 |
| 兵源 | 三種預備兵**總和** ≥ **600 點**（6,000 人）|

**兩道都要過**（原版用同一個 `dl` 旗標，任一條成立就寫 1 ＝ 擋下）。
兩道閘都不看敵情，只看自己的錢與兵。

過了就 `sub_16E8F`：用君主的武將記錄自動編一支軍團
（`sub_1461D`，與補兵同一條路），然後——

```asm
0001699E …
  call sub_16E8F
  and byte ptr [di], 0FBh    ; ★ 清掉軍團 +0x00 的位元 2
```

⭐ **`sub_16E8F` 一律 `or [si], 4`（委任），而出陣這條緊接著把它清掉。**
君主親自出陣的那一支**由玩家指揮**，不是委任。

## 4. remake 實作

| 項目 | 作法 |
|---|---|
| 遷都判定 | `capital.AcceptRelocation`（已存在，先前沒有任何呼叫端）|
| 出陣判定 | `persuasion.AcceptSortie` ＋ `SortieFundsGate`／`SortieReserveGate` |
| 入口 | `World.AdviseRelocateAccepted`／`AdviseRelocate`、`AdviseSortieAccepted`／`AdviseSortie`。**判定與執行分開**，畫面要先把君主的回答演完 |
| 編成 | `autoFormCorps(faction, leader, delegated)`——從 `formAICorpsTo` 抽出來的共用路徑，對應 `sub_16E8F`。出陣傳 `delegated = false` |
| 畫面 | 進言選單五列；第四項先用一覽表挑據點，第五項直接演。兩者都走 `sayVerdict`（`sub_13B08` 的三句）|

**remake 差異**：遷都的目標原版是在地圖上選點（`sub_17400`），
remake 用據點一覽表挑。**這是操作方式的差異，不是規則的差異**。

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestAcceptSortieNeedsBothGates`：兩道閘各自擋得住，都過才回 true |
| 單元測試 | `TestSortieFundsGateScalesWithAggression`：好戰 15 ⇒ 0、好戰 0 ⇒ 15,360 |
| 單元測試 | `TestAdviseSortieFormsUndelegatedCorps`：編出來的軍團**不是委任** |
| 單元測試 | `TestAdviseRelocateMovesCapitalAndCorps`：首都換了，掛舊首都的軍團跟著改掛 |
| 單元測試 | `TestAdviseRelocateRefusesSmallerCity`：只看生產力會過、加上類型就該拒絕 |
| 單元測試 | `TestVerdictTalkIndicesMatchOriginal`：兩項的六個位置 ＝ 386/389/390/393 與 396/399/400/403 |
| 截圖 | [`../playtest/35`](../playtest/35-advise-verdict-screens.md) |

## 6. 未解

| 項目 | 現況 |
|---|---|
| `sub_16E8F` 編成前的其餘檢查 | 只確認「君主還沒帶軍團」這一條。⚠ 它呼叫的 `sub_16EC9` 本身已解（六槽 × 三候選兵種表、每槽門檻 `0x32`、試算在堆疊副本上做，見 [`../re/30`](../re/30-corps-formation-ui.md) §7.3）|
| 遷都的地圖選點 | `sub_17400` 沒讀，remake 用一覽表代替 |
| 進言的指令列 | 五項在原版指令樹裡的排法（`docs/re/22`）沒有逐格對過，remake 用自己的小視窗 |
