# 108 — 進言問理由之前，君主要先講那一句

**狀態：CONFORMED。`AskReason` 這一支漏掉了君主的回答，
選單直接跳出來、上框還停在開場那句。**

- 日期：2026-09-01
- 出處：[`44`](44-advise-original-text.md)／[`45`](45-advise-scene-layout.md)
  的 `sub_13830`（`cx = base + al × 3 + 4`）；實機
  [`../playtest/56`](../playtest/56-lubu-flow-parity.md) §4.3
- 推論等級：**confirmed**——算式早就解出來且已用在其他三個反應上，
  這一輪只是補上漏掉的呼叫，實機也拍到了

## 1. 原版做什麼

進言（敵對提案）從按下指令列到出現理由選單，中間**演三句**：

| # | 說話者 | TALK 索引 | 呂布（說話類型 1）拿到的 |
|---|---|---|---|
| ① | 君主 | `base + 說話類型` | #87「是什麼事啊，{4}，說出來聽聽。」 |
| ② | 軍師 | `base + 3` | #89「想請主公答允對{3}的進兵。」 |
| ③ | 君主 | `base + 4 + 反應碼 × 3 + 說話類型` | #97「要是沒有勝算，那就不能答允！」 |

敵對提案的 `base` ＝ `0x56` ＝ 86。四個反應碼共用第三句的算式：

| 反應碼 | 意思 | 呂布拿到的 |
|---:|---|---|
| 0 | Refuse | #91 |
| 1 | Agree | #94 |
| **2** | **AskReason** | **#97** |
| 3 | AlreadyAtWar | #100 |

**第三句演完才出五項理由選單。**

## 2. 量到的落差

`cmd/wlgame/advise.go` 的 `beginPersuasion`：

```go
switch reaction := persuasion.FirstReaction(g.adviseCmd, s, queued); reaction {
case persuasion.AskReason:
    g.sess = persuasion.Begin(g.adviseCmd, s)   // ← 沒有第三句
default:
    …
    g.adviseSay(adviseLord, persuasion.TalkReplyIndex(base, reaction, g.playerTalkVariant()))
```

`TalkReplyIndex` 只在 `default` 分支叫得到。四個反應碼裡**只有 `AskReason`
少演一句**，而它正是唯一會進入說服迴圈的那一個——也就是玩家最常走的路。

說話類型的選法本身沒問題：同一顆執行檔走「已經交戰中」那條路顯示的是
#100，正好也是說話類型 1 的變體。

## 3. 演算法

```
第三句索引 = TalkReplyIndex(base, 反應碼, 說話類型)
           = base + 4 + 反應碼 × 3 + 說話類型

AskReason 也照這個算，唯一的差別是**算完之後還要開說服迴圈**。
```

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 呈現層 | `cmd/wlgame/advise.go` 的 `beginPersuasion`：`AskReason` 那一支補 `adviseSay` |
| 差異 | 無 |

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestAskReasonPlaysLordReplyBeforeMenu`（`cmd/wlgame`）|
| 對原版 | [`../playtest/56`](../playtest/56-lubu-flow-parity.md) §4.3 的原版影格（#97 ＋ 選單）|
