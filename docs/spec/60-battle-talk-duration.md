# 60 — 戰場對白顯示多久：60 個 tick，每側各一個

**狀態：CONFORMED。** 到期機制與常數都是逐行讀出來的，remake 照著改，
同一個局面的 `field` 從 19.88% 降到 **0.84%**。

- 日期：2026-08-18
- 出處：`sub_1C315`（顯示時設到期時刻）、`sub_1A12A`（每 tick 檢查三個到期）、
  `sub_1A69F`（腳本指令 16）
- 推論等級：**confirmed（靜態）**

## 1. 原版做什麼

戰術主節拍 `sub_1A12A` 每 tick 做四件事：

```asm
inc byte ptr cs:word_1D318          ; 節拍計數器（**byte**，0–255 循環）
ax = word_1D318
cmp cs:word_1D322, ax / jnz …       ; 側 0 的對白框到期 → sub_1C3B8(cl=0) 擦掉
cmp cs:word_1D324, ax / jnz …       ; 側 1 → sub_1C3B8(cl=1)
cmp cs:word_1D326, ax / jnz …       ; 門強度條 → sub_1C4A6（[`32`](32-gate-strength-bar.md)）
```

顯示端 `sub_1C315` 在畫完之後設到期時刻：

```asm
mov dx, cs:word_1D318
add dl, 3Ch                         ; ⭐ ＋60
mov cs:[bx-2CDEh], dx               ; bx ＝ 側 × 2 → word_1D322／word_1D324
```

| 東西 | 到期時刻 | 收掉的常式 |
|---|---|---|
| **戰場對白框** | 目前節拍 **＋ 0x3C（60）** | `sub_1C3B8`（每側一個）|
| 門強度條 | 目前節拍 ＋ 0x14（20）| `sub_1C4A6` |

⚠ **比較是 `==` 不是 `>=`**，而計數器是會繞回的 byte——所以「到期」是
**剛好那一格**。新的對白會把到期時刻整個覆寫，不會疊加。

⭐ **每一側各有一個框**：`word_1D322` 是側 0、`word_1D324` 是側 1，
兩邊可以同時掛著各自的台詞（開場那一段就是你一句我一句）。

## 2. remake 實作

| 項目 | 位置 |
|---|---|
| 壽命常數 | `cmd/wlgame/battletalk.go` 的 `battleTalkDuration` ＝ **60** |
| 每側一個框 | `battleTalkSlots`：`entry[2]` ＋ `remaining[2]`，`set`／`current(side)`／`tick`／`clear(side)`。同一側再來一句就覆寫，不排隊 |
| 上下對應 | `battleTalkState` 由側 0 填 `Top`、側 1 填 `Bottom`（`cmd/wlgame/battle.go`）|
| 推進 | 正常遊戲迴圈由 `updateBattle` 的 `g.tickBattleTalk(n)` 推；**驗收路徑**（`-battle-steps`）也要推同樣的格數，否則框永遠不消失 |

⚠ **不是佇列。** 一個一次只顯示一句的佇列會把兩句串成 120 格，
而原版兩個框各自從自己被設的那一刻算 60 格——開場那一段兩個框是疊著的。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestBattleTalkDurationMatchesOriginal`（常數 ＝ 60，且與門強度條的 20 是兩個不同的值）、`TestBattleTalkSlotsAreIndexedBySide`、`TestBattleTalkSlotsExpireIndependently` |
| 對原版 | 同一個局面的 `field` 在第 61 步 ＝ **0.84%**（[`../playtest/40`](../playtest/40-tactical-parity.md) §8）|

## 3.5 開戰那一對台詞的 TALK 索引

`sub_1A3C3` 在開戰時送**一對**台詞：第一句 `0x1BA`、第二句 `0x1BB`。
兩者都是八格一組的組編號（索引 ≥ `0x196`，[`../formats/01`](../formats/01-talk-dat.md) §3），
所以實際取到哪一則要看說話者的說話型。

⚠ **「上格是攻方、下格是守方」只是強推論**——remake 照原版實錄影格上的
位置接線（`cmd/wlgame/battletalk.go`），`sub_1A3C3` 的側別參數沒有逐行讀。
與 §1 的「每側各一個框」不衝突：那一對是分別掛到兩側，不是同一側兩句。

## 4. 未解

| 項目 | 現況 |
|---|---|
| 開戰 pair 的側別對應 | `0x1BA` → 上格、`0x1BB` → 下格是**強推論**（照影格位置接的）；`sub_1A3C3` 怎麼決定側別沒讀（§3.5）|
| `byte_1D349` 的三個值 | `sub_1A69F` 拿它當「這句要不要顯示」的閘（`al & 6` 那一段還沒逐位讀）。0／1／2 三種值由 `sub_1A6FA` 切換 |
| 玩家按鍵能不能提早關掉 | remake 可以按鍵推進；原版是否有這條路沒讀 |
