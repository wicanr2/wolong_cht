# 118 — 截圖的時機用局面條件：`-shot-when` 與 `-auto-messages`

**狀態：CONFORMED。** 兩個驗收旗標補上兩個缺口：對拍的取樣點寫死步數
（改動規則層之後會安靜地爛掉），以及自然流程被遭遇訊息擋住走不到戰場。
**這是 remake 的驗收設施，不是遊戲規則**——兩個旗標都預設關閉，
不帶就是原本的行為。

- 日期：2026-09-03
- 出處：缺口本身有出處，設施沒有。
  取樣點那一條寫在 [`91`](91-tactical-parity.md) §6
  （「在有『跑到條件成立就截圖』的旗標之前，改動規則層之後要重新量」），
  訊息那一條寫在 [`105`](105-encounter-goes-straight-to-battle.md) §4
  與 [`../playtest/58`](../playtest/58-parity-retest-20260902.md) §3–§4
- 推論等級：**不適用**（沒有對應的原版行為；原版的判準是玩家的眼睛）
- 相關：[`117`](117-fixture-arms-duel-before-stepping.md)（同一族的驗收路徑缺陷）、
  [`112`](112-cursor-idle-resume-delay.md)（訊息收掉之後的暫停）

## 1. 缺口

### 1.1 取樣點寫死步數，會安靜地爛掉

`docs/playtest/40` 的攻城取樣點是「第 61 步」。那一步在 2026-08-18 是
**攻方正在攻門、兩個對白框剛過期**那一刻；規則層改過之後同一個步數落在
「攻方剛要出發」，而**沒有任何測試會紅**——那些數字只寫在文件裡
（[`91`](91-tactical-parity.md) §6）。2026-09-02 又踩了一次：
兵的戰力接回統率力之後三個取樣點全部要重找。

### 1.2 自然流程走不到戰場

2026-08-29 起遭遇會先跳一則訊息再進戰場
（[`105`](105-encounter-goes-straight-to-battle.md)），而截圖模式沒有人按掉它。
野戰對拍原本的 `-shot-frames 400` 因此停在訊息上——實測第 400 到
1,100 幀**五張 PNG 逐位元組相同**（[`../playtest/58`](../playtest/58-parity-retest-20260902.md) §3）。
現在靠 `-open-battle -siege-corps` 繞過去，但那條路跳過了整個戰略層，
**驗不到「遭遇怎麼進戰場」本身**。

## 2. 旗標

| 旗標 | 意思 |
|---|---|
| `-shot-when <條件>` | 截圖的時機改用局面條件。留白 ＝ 照 `-shot-frames`（原本的行為）|
| `-shot-deadline N` | 配 `-shot-when`：等到第 N 幀還不成立就**放棄並回非零**（預設 20000）|
| `-auto-messages` | 訊息框自動按掉（每幀推一頁），讓自然流程走得下去 |

`-shot-frames` 在有 `-shot-when` 時退成**下限**：先跑滿那麼多幀，再開始等條件。

### 2.1 條件

| 值 | 成立條件 |
|---|---|
| `battle` | 戰術戰鬥開著 |
| `battle-frame:N` | 戰術戰鬥開著，而且 `Battle.Frame ≥ N` |
| **`battle-settled`** | ⭐ **開場布陣走完的那一刻**：這一拍**沒有任何一個兵移動**（而且已經過了第 10 拍）。原版把兵擺在戰場邊界再讓他們走進陣形（[`133`](133-opening-deployment.md)），走完之後有一小段誰都不動的空窗，接著腳本才下第一道命令——**那是戰術對拍唯一「兩邊都靜止」的取樣窗**（[`../playtest/74`](../playtest/74-settled-tick-parity.md)）|
| `gate-bar` | 戰術戰鬥開著，而且 `Battle.StructureBar()` 的第二個回傳值為真（門強度條顯示中）|

`gate-bar` 就是 [`91`](91-tactical-parity.md) §6 那張表的第二列，
攻城取樣點的三個局面條件之一。

### 2.2 ⛔ 條件不成立就不截圖，而且要回非零

**沉默地截一張別的局面，比失敗更糟**——那張圖會被拿去比、
得到一個看起來像回歸的數字，而沒有人知道它拍的不是同一件事
（`~/diagnosis-notes/docs/03-silence-is-not-success` 的四個閘門）。
所以逾時的處置是**不寫檔 ＋ `log.Fatal`**，訊息裡帶條件名與已經跑過的幀數。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 條件解析 | `cmd/wlgame/shotwhen.go` 的 `parseShotWhen`（不認得的值在啟動時就失敗，不是跑到一半）|
| 截圖閘 | `cmd/wlgame/main.go` 的 `maybeSaveShot` |
| 逾時 | 同檔 `Update`：`-shot-when` 沒成立而且過了 `-shot-deadline` 就回錯誤（`RunGame` → `log.Fatal` → exit 1）|
| 訊息自動按掉 | `cmd/wlgame/messages.go` 的 `updateMessageOnly` |

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `cmd/wlgame/shotwhen_test.go`：各條件的成立／不成立、`battle-frame:` 的邊界、**不認得的值要回錯誤**（少了這一條，打錯字會靜靜退回「照幀數截圖」）、`TestBattleSettledWaitsForEveryoneToStop` |
| 對原版 | `battle-settled` 取到的畫面與寫死 `-battle-steps 70` **同一個狀態**（`sb-minimap` 都是 8 px，[`../playtest/74`](../playtest/74-settled-tick-parity.md)）|
| 對原版 | [`../playtest/59`](../playtest/59-shot-when-natural-flow.md)：野戰走**自然流程**（`-auto-messages -shot-when battle-frame:52`）與 `-open-battle` 那條捷徑截出同一張畫面；攻城用 `-shot-when gate-bar` 取樣 |

## 4.5 ⚠ `battle-settled` 的兩個坑，都是「成立在錯的地方」

| 坑 | 症狀 | 修法 |
|---|---|---|
| 拿**座標範圍**當判準 | 範圍在大部分人到位之後就不再變，而**最後幾個還在走** | 判準改成「這一拍動了幾個 ＝ 0」|
| 跨**呼叫次數**比而不是跨**戰術拍**比 | 條件每個畫面幀都被問一次，而戰鬥不是每一幀走一拍——連問兩次看到同一組座標是常態，於是**第 2 拍就成立** | 記住上一次的 `Battle.Frame`，同一拍直接回 false |

⭐ 第二個坑**回了一張圖、退出碼 0**，只是那張圖是第 2 幀的——
與這一份規格要消滅的失敗模式（打錯字靜靜退回照幀數截圖）是同一個形狀。

## 5. 未解

| 項目 | 現況 |
|---|---|
| 條件的組合 | 一次只吃一個條件。[`91`](91-tactical-parity.md) §6 的攻城取樣點其實是三個條件同時成立（城壁挨過打、條顯示中、對白框已收），現在只判得了第二個 |
| 對白框的收掉時刻 | 沒有條件可以判「兩側的對白框都到期」，那要規則層先把 `word_1D322`／`word_1D324` 的到期時刻露出來 |
