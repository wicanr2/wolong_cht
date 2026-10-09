# 111 — 君主帶著軍團的時候，進言整個關掉

**狀態：CONFORMED。保留既有 remake 實作與測試紀錄；原版閘門依下列勘誤。
本輪 C 局部驗收範圍見 [spec/244](244-c-strategy-command.md)。**

- 日期：2026-09-01
- 出處：原版 `sub_16224`（IDA linear `0x16224`–`0x1623F`），
  固定輸入與資料庫身分見 [re/124](../re/124-c-strategy-command-restoration.md)。
- 推論等級：**已證實（原始指令）**；2026-09-01 的 remake 決策紀錄見 §2–§5。

## 1. 原版的進言閘門

**2026-10-09 勘誤（re/124，已證實）**：`sub_16224` 取得玩家君主後，
比較武將 `+0x17` 是否為 0。**任意非零職務都顯示 TALK #64，整個進言選單不開**。
原始門檻沒有查軍團 Alive，也沒有只允許職務 1 才阻擋。
證據見 [re/124](../re/124-c-strategy-command-restoration.md)。

原版的編成候選在建清單的 callback（線性 `0x176A0`）就排除君主
（`si ≠ 君主記錄`，confirmed），所以**君主帶兵只有一條路**：
進言的第五項「請求君主出陣」，由 `sub_16E8F` 讓君主本人編一支軍團。
`sub_16EC9` 是六槽兵種與預備兵的堆疊試算，不是君主職務閘門。
舊版將它指認為「君主已經帶兵」的檢查，依上述原始比較指令推翻。

remake 把「主君編成」做成系統選單的一列而且**預設放行**
（[`76`](76-lord-not-in-formation.md)，使用者裁定），
這增加了君主編成的入口；帶兵後無法進言已有原版 `sub_16224` 的依據。

## 2. 2026-09-01 的設計論證（歷史）

進言是軍師對君主說話（[`45`](45-advise-scene-layout.md)：上框君主、下框軍師）。
君主帶著軍團離開首都之後，這場對話的另一方不在——五項裡沒有一項成立：

| 進言的項目 | 君主不在時 |
|---|---|
| 敵對提案／停戰提案／請求協助 | 要君主點頭，而君主不在 |
| 遷都 | 同上 |
| 請求君主出陣 | 與另外四項共同被原版 `sub_16224` 的職務閘門擋下 |

這是當時的設計理由。舊版「把門檻往上提一層」的推論已被 §1 的原版入口證據推翻。

## 3. 既有 remake 演算法紀錄

```
君主在帶兵 = 玩家勢力的君主武將編號 n 落在軍團表範圍內，且 軍團[n].Alive

進言（指令列第 0 格 / 鍵盤 P）：
    君主在帶兵 → 不開視窗，事件列顯示「主公正在領軍，無法進言」
    否則        → 照舊

請求君主出陣（進言第五項）：
    維持原本的 AdviseSortieAccepted —— 判準本來就是同一個
```

⭐ **判準只有一份。** `World.LordLeadsCorps()` 是唯一入口，
`AdviseSortieAccepted` 與進言的開窗閘都呼叫它——
避免兩處各寫一份「君主在不在」而行為分岔（`CLAUDE.md` §7 第 6 條）。

上段保留歷史實作定位；其 Alive 判準與事件列呈現不等同 §1 的原始 duty 比較與 TALK #64。
本輪只修證據說明，未修改正式 Go。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/advise.go` 的 `LordLeadsCorps()`；`AdviseSortieAccepted` 改用它 |
| 呈現層 | `cmd/wlgame/advise.go` 的 `openAdvise` 開頭擋下 |
| 語系 | `translations/ui-ja.json`／`ui-en.json` 各補一條 |
| 差異 | 君主編成入口是 remake 差異；進言阻擋已有原版依據。Alive 判準與事件列呈現依 §3 留為實作差距 |

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestLordLeadsCorpsBlocksAdvise`、`TestLordLeadsCorpsSharedWithSortie`（`internal/state`）|
| 對原版 | 原始 duty 閘門可核對；固定指令見 [re/124](../re/124-c-strategy-command-restoration.md)。本頁舊測試只證明既有 remake 實作自洽，不代證本輪 C 驗收 |
