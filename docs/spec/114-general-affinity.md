# 114 — 武將的心向勢力（`+0x19`）：在野出仕與俘虜歸降

**狀態：CONFORMED。** `state.General` 有 `Affinity`／`Sovereign`／
`VanishIfAffinityGone` 三個欄位，在野出仕接進月結，俘虜歸降加上心向的勢力這一閘。
隨機投靠那一條 2026-09-04 接上了（[`130`](130-freelance-random-join.md)）。

- 日期：2026-09-02
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_1585F`（月結逐武將）、
  `sub_15899`（在野出仕）、`sub_15940`（俘虜每月）
  ——[`../re/77`](../re/77-general-affinity-and-flags.md) §1、§2
- 推論等級：**confirmed（靜態 ＋ 四劇本 508 筆分布）**
- 相關：[`../formats/08`](../formats/08-sinario-save.md) §3、
  [`../mechanics/70`](../mechanics/70-ai.md) §3.9、
  [`../mechanics/60`](../mechanics/60-personnel.md) §4.1

## 1. 原版做什麼

`sub_1585F` 在月結（`sub_15358`，`sub_15695` 之後、`sub_155A6` 之前）
逐一掃 127 名武將：

```
[+0] bit 7 ＝ 0（不在場）            → 跳過
[+0x18] != 0（倒數計時器）           → 遞減，這個月不做別的
[+0x1C] == 0xFF（在野）              → sub_15899
[+0x1D] != 0xFF（俘虜）              → sub_15940
```

兩支都讀 `+0x19` ＝ **心向的勢力**（`0xFF` ＝ 沒有）。

### 1.1 在野（`sub_15899`）

| `+0x19` | 每月 |
|---|---|
| 有值 | 亂數 `< 0x40`（**25%**）才兌現：`+0x19` 清成 `0xFF`，那個勢力還在就 `+0x1C = 該勢力`；已滅而且旗標 bit 5 設著就 `[+0] = 0`（整筆歸零、退場）|
| `0xFF` | 走既有的隨機投靠（[`../mechanics/70`](../mechanics/70-ai.md) §3.9）|

出仕的對象是玩家勢力時跳訊息 `0x29`。

### 1.2 俘虜（`sub_15940`）

| 亂數 | 機率 | 結果 |
|---|---:|---|
| `< 0x20` | 12.5% | 逃走（訊息 `0x41`），`+0x1C = 0x18`（無所屬）|
| `0x20`–`0x3F` | 12.5% | **`+0x1C == +0x19` 才**歸降（訊息 `0x42`），清 `+0x1D`／`+0x17` |
| `≥ 0x40` | 75% | 不動 |

## 2. 欄位對應

| 原版 | remake |
|---|---|
| 武將 `+0x19` | `state.General.Affinity`（`0xFF` ＝ `noFaction`）|
| 旗標 `+0x00` bit 6 | `state.General.Sovereign` |
| 旗標 `+0x00` bit 5 | `state.General.VanishIfAffinityGone` |

旗標那個 byte 現在**逐位元蓋寫**（`setFlag`）：bit 7／6／5／4 由 remake 的欄位決定，
bit 0 還沒解，原樣保留——存檔寫回是「改寫不是重建」（`CLAUDE.md` §9）。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 三個欄位的載入與寫回 | `internal/state/state.go` |
| 在野出仕 | `internal/state/freelance.go` 的 `recruitFreelanceGenerals`，在 `compactEventQueue()`（＝ `sub_12BD9` 邊界）**之前**呼叫，對齊 `sub_15358` 的順序 |
| 出仕時的勢力武將數 | 同上，`raiseGeneralCount`——原版 `loc_1591A` 的最後一步是 `sub_12AD2` |
| 俘虜歸降的條件 | `internal/state/event10_approx.go`：`roll < 0x40` 那一支加 `g.Faction == g.Affinity` |
| 跨層傳遞 | `internal/scenario` 由 `state.General` 內嵌帶過去，不必改；戰鬥層不讀這三個欄位 |

⚠ **只做在野那一條，俘虜那一條沿用既有的近似 producer。**
兩者在原版是互斥的（`+0x1C == 0xFF` 對 `+0x1D != 0xFF`），
remake 也一樣——`recruitFreelanceGenerals` 只看 `Faction == noFaction`，
`produceApproximateEvent10` 只看 `Faction == Player && Captor != noFaction`。
不會有武將同時被兩邊處理，倒數計時器也不會被扣兩次。

⚠ **`sub_15899` 的隨機投靠那一條這一版不接。** 它要勢力的武將數
（`+0x18`）與據點數（`+0x23`）兩個執行期欄位——⭐ **2026-09-04 查證：
那兩個欄位早就在**（`state.Faction.Generals`／`Cities`），所以這個理由
在寫下的當時可能成立，現在已經不成立；
這一版只做「心向的勢力」那一條，因為開局 81 名在野武將全部有值，
隨機投靠要等他們兌現完才輪得到（[`../mechanics/70`](../mechanics/70-ai.md) §3.9）。
沒接的部分留在 §5。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestGeneralAffinityAndFlagsRoundTrip`：四劇本逐武將，`+0x19` 與旗標整個 byte 都 byte-for-byte 寫得回去 |
| 單元測試 | `TestFreelanceJoinsAffinityFaction`：亂數 `< 0x40` 出仕、`Affinity` 清成 `0xFF`、勢力武將數 +1 |
| 單元測試 | `TestFreelanceHoldsOnHighRoll`：`≥ 0x40` 連欄位都不清 |
| 單元測試 | `TestFreelanceVanishesOnlyWithFlag`：勢力已滅時，bit 5 設著才退場 |
| 單元測試 | `TestFreelanceTimerGate`：倒數沒歸零的月份只遞減 |
| 單元測試 | `TestCaptiveSurrenderNeedsAffinity`：關押方不是 `Affinity` 時不歸降，但逃走那一支照走 |
| 單元測試 | `TestEveryFreelanceGeneralHasAffinity`：四劇本的在野武將**全部**有值，附「一個在野武將都沒有」的正對照 |
| 既有測試 | `TestWorldStaysConsistent`：出仕沒有維護勢力武將數的話，818 天就會撞上不變量 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~隨機投靠那一條~~ | ✅ **2026-09-04 接上了**（[`130`](130-freelance-random-join.md)）。擋它的前提查證後本來就不成立——`state.Faction` 的 `Generals`／`Cities` 早就在且有維護 |
| ~~出仕的畫面通知~~ | ✅ **2026-09-04 接上了**：投靠玩家時排 `TalkNotice{Index: 0x29}`（TALK #41「{1}加入麾下了。」）。兩條路共用 `joinFaction`，通知條件與原版的 `loc_1591A` 相同（[`130`](130-freelance-random-join.md) §3）|
| 旗標 bit 5 之外的退場條件 | `sub_15899` 只在「心向的勢力已滅」時看 bit 5；bit 5 沒設的武將會留在原地等下一輪，這一點沒有實機驗證 |
