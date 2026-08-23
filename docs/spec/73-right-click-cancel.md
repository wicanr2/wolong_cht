# 73 — 右鍵取消是輸入層的語意，不是每個視窗各自的功能

**狀態：CONFORMED。七個面板改成問同一支 `cancelled()`，右鍵與 ESC 都退回上一層。**

- 日期：2026-08-23
- 出處：`sub_121E7`（等待按鍵，**CF=1 代表右鍵＝取消**，[`../re/22`](../re/22-strategy-command-tree.md) §3）；
  `sub_18F7C` 是取消時的擦除（[`../re/53`](../re/53-lord-select-window.md) §3）；
  [`../re/24`](../re/24-unread-function-catalogue.md) 列出 `sub_17E1F`／`sub_18B7C`／`sub_18DC8`／`sub_18E5A`
  四支等待常式全部標「CF=1 為右鍵取消」
- 推論等級：**confirmed**

## 1. 原版做什麼

原版的每一個模態畫面都停在同一組「等待按鍵」常式上，而那些常式
**回傳 CF=1 就是取消**。右鍵在哪裡按不影響——它是輸入層的結果，
不是某個熱區的功能。

所以原版的規則只有一條：

> **模態畫面開著時，任何位置的右鍵 ＝ 退回上一層。**

⚠ 這與**橫幅上那五個開關**不同。那五格是 `左鍵開、右鍵關`
（[`13`](13-main-window-toggles.md) §2.3），右鍵在那裡是「關掉這一格對應的常駐視窗」，
是熱區的功能。常駐視窗不是模態的，它不吃這一條。

## 2. remake 現況：七個面板漏掉

| 面板 | 右鍵 | ESC |
|---|---|---|
| 一覽表（`window.go`）| ✅ | ✅ |
| 軍團編成（`corps.go`）| ✅ | — |
| 行軍模式（`marchmode.go`）| ✅ | — |
| 存讀檔（`save.go`）| ✅ | — |
| 勢力一覽（`factionpicker.go`）| ✅ | ✅ |
| 據點情報（`cityinfo.go`）| ❌ | ✅ |
| 軍團情報（`corpsinfo.go`）| ❌ | ✅ |
| 財政（`finance.go`）| ❌ | ✅ |
| 進言（`advise.go`）| ❌ | ✅ |
| 外交（`diplomacy.go`）| ❌ | ✅ |
| 數值輸入（`amountpanel.go`）| ❌ | ❌ |
| 撥款（`funding.go`）| ❌ | ✅ |

⭐ **這是「一條規則散成十二份實作」的典型後果**（`CLAUDE.md` §7 第 6 條）。
原版只有一個地方決定「右鍵＝取消」，remake 卻讓每個面板自己記得要接——
於是漏掉七個，而且**漏掉的方式是安靜的**：面板打得開、ESC 也能關，
只有拿滑鼠玩的人會踩到。

## 3. 演算法

```
cancelPressed():
    右鍵剛按下 → true
    ESC 剛按下 → true        # remake 差異：鍵盤替代，原版沒有 ESC
    否則        → false
```

**每個模態面板的 update 一律先問 `cancelPressed()`**，成立就退回上一層。

⚠ **不要在全域攔截。** 面板有巢狀（進言 → 數值輸入），全域攔截會一次關掉兩層；
原版的取消是**一次退一層**，因為 CF=1 只回到當前那一層的呼叫端。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 共用判定 | `cmd/wlgame/cancel.go` 的 `cancelPressed()` ／ `game.cancelled()` |
| 接上的面板 | 上表七個 ❌ 全部改成先問 `cancelPressed()` |
| 差異 | ESC 是 remake 加的鍵盤替代（`CLAUDE.md` §9 的 ESC 規則），原版只有右鍵 |

## 5. 驗證

| 方式 | 結果 |
|---|---|
| 單元測試 | ✅ `TestCancelClosesEveryModalPanel`（`cmd/wlgame`）：據點情報／軍團情報／財政各開一次、送取消、斷言關掉 |
| 反向對照 | ✅ `TestPanelStaysOpenWithoutCancel`：沒按取消就不能關——否則「永遠關掉」也會讓上一條通過 |
| 巢狀不越級 | ✅ `TestCancelClosesOneLayerOnly` |

⚠ **取消必須是可注入的**（`game.cancelFn`）。`inpututil` 讀的是 Ebiten 的
全域輸入狀態，無頭測試裡永遠是 false——直接呼叫 `cancelPressed()` 的話，
「面板關不關」**測起來永遠通過**，而那正是這一輪要修的那種安靜失敗。

## 6. 未解

| 項目 | 現況 |
|---|---|
| 原版右鍵是否也關常駐視窗 | 沒量過。常駐視窗不走模態等待常式，推測不關，但**沒有實機證據** |
