# 59 — 兩個驗收旗標的驗收：自然流程與捷徑截出**同一張**畫面

**狀態：四項全部通過。** 野戰走自然流程（遭遇訊息自動按掉 → 進戰場 →
等到戰場第 52 幀）截出來的畫面，與 `-open-battle -siege-corps` 那條捷徑
**九區逐像素相同（全 0 px）**；攻城的 `-shot-when gate-bar`
在戰場第 148 幀附近觸發，與 [`../spec/91`](../spec/91-tactical-parity.md) §6
記的「條第一次亮」同一刻；條件永遠不成立時**不寫檔、回非零**。

- 日期：2026-09-03
- 規格：[`../spec/118`](../spec/118-shot-when-condition.md)
- 原版側：沿用 [`58`](58-parity-retest-20260902.md) 的裁切
  （`parity-field14/b0.png`、`probe-march/e10.png`）
- 比對：`tools/py.sh tools/parity_diff.py … --regions tactical`

## 1. 野戰：自然流程終於走得到戰場

```
tools/parity_shot.sh out.png -direct -scenario 0 -player 0 -seed 1 \
  -save-file workplace/parity/SAVE-FIELD.DAT -load-slot 0 \
  -auto-messages -shot-when battle-frame:52 -shot-frames 1
```

第 **431** 幀截圖（遭遇訊息按掉 → 進戰場 → 戰場走到第 52 幀）。

| 比對對象 | 結果 |
|---|---|
| 原版 `b0` | `field` **95 px（0.05%）**、`sb-minimap` 32 px、其餘七區 **0 px** |
| `-open-battle -siege-corps 39,35 -battle-steps 52` 那一張 | ⭐ **九區全部 0 px** |

⭐ **兩條路徑逐像素相同，這是捷徑的正對照。** `-open-battle -siege-corps`
跳過整個戰略層直接擺兩支軍團，一直沒有東西證明它擺出來的局面
與自然遭遇一樣；現在有了。順帶也證明 `-auto-messages` 只是替玩家按了
Enter——它沒有改到任何被驗收的東西
（對比 [`49`](49-parity-retest-20260827.md) §2 那個把音效關掉的捷徑）。

## 2. 攻城：`gate-bar` 落在條第一次亮的那一刻

```
… -open-siege -siege-node 82 -siege-corps 81,39 -battle-steps 0 \
  -shot-when gate-bar -shot-frames 1
```

第 **478** 幀截圖，`field` 1.32%。換算得回戰場幀數：`stageEncounter`
收尾把戰術速度設成 1，於是每 **3.30** 個畫面推一步
（`speed.Steps`：`2913 / (1 × 16 × 600)`），478 ÷ 3.30 ≈ **145**——
與 [`../spec/91`](../spec/91-tactical-parity.md) §6 記的「第 148 幀第一次同時滿足」
是同一刻。

⚠ **條件標的是窗口的起點，不是最好的取樣點。** 那一刻攻方才剛碰到城壁，
`field` 1.32% 比第 160／240 拍的 0.85%／0.84% 差
（[`58`](58-parity-retest-20260902.md) §1.2）。

### 2.1 ⭐ 真正該用的寫法：步數 ＋ 條件一起帶

```
… -battle-steps 160 -shot-when gate-bar -shot-frames 1
```

| 比對對象 | 結果 |
|---|---|
| 原版 `e10` | `field` **0.85%**、`sb-enemy` 12 px、`bottom` 2 px、四區 0 px、`sb-minimap` 1.64% |
| 沒帶條件的那一張（`58` §1.2 的第 160 拍）| **九區全部 0 px** |

**帶條件不改變截到的畫面，改變的是「步數落空時會不會被發現」。**
第 160 拍今天落在條亮著的窗口裡；哪天規則層再改、它落到窗口外，
這條指令會**當場失敗**，而不是安靜地截一張別的局面回來讓人比。
寫死的步數就是這樣從「取樣點」爛成「一個數字」的
（[`../spec/91`](../spec/91-tactical-parity.md) §6）。

## 3. 條件不成立就要失敗

拿沒有城壁的野戰去等 `gate-bar`：

```
… -save-file SAVE-FIELD.DAT -auto-messages -shot-when gate-bar -shot-deadline 900
⚠ -shot-when gate-bar 到第 901 幀都沒有成立，沒有截圖（-shot-deadline 900）
parity_shot 回 1
ls: 無法存取 'workplace/parity/p118-never.png': 沒有此一檔案或目錄
```

三件事都要成立才算：**訊息帶條件名與幀數**、**exit 1**、**沒有寫出 PNG**。
第三件最重要——留下一個舊的或錯的 PNG，下一步的 diff 會照樣跑完並吐出數字。

## 4. 未解

| 項目 | 現況 | 下手點 |
|---|---|---|
| 一次只判一個條件 | `../spec/91` §6 的攻城取樣點是三個條件同時成立，現在只判得了「條顯示中」| 條件字串改成可以用 `+` 串接，或直接支援 `gate-bar+talk-clear` |
| 對白框到期沒有條件可判 | 規則層沒有把 `word_1D322`／`word_1D324` 的到期時刻露出來 | `internal/rules/tactical` 加一支唯讀的查詢 |
| 攻城取樣點的 0.84% 地板 | 局面不等價，與這兩個旗標無關（[`58`](58-parity-retest-20260902.md) §4）| 要「存檔與影格出自同一次擷取」的攻城素材 |
