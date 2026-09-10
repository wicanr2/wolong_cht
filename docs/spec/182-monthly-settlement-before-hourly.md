# 182 — 月結跑在每「時」更新之前

**狀態：CONFORMED。** 原版的時鐘常式是一層層 fall through：
`loc_11DC1` 的 `call sub_15358`（月結）落在 `loc_11DE0` 的
`call sub_13E11`（每「時」的世界更新）**之前**。remake 的順序相反。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_11D8E` 的 `loc_11DC1`／`loc_11DD7`／
  `loc_11DE0`（[`../re/06`](../re/06-game-clock.md) §1）
- 推論等級：confirmed
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §41.5
- 相關：[`../re/07`](../re/07-monthly-settlement.md)（`sub_15358`）、
  [`../re/08`](../re/08-hourly-update.md)（`sub_13E11`）

## 1. 原版做什麼

```asm
loc_11DC1:
        inc  byte ptr ds:0CF4h        ; 月++
        …查每月天數表、日 = 0…
        call sub_15358                ; ★ 月結
loc_11DD7:
        inc  byte ptr ds:0CF0h        ; 日++
        mov  byte ptr ds:0CF3h, 0     ; 時 = 0
loc_11DE0:
        inc  byte ptr ds:0CF3h        ; 時++
        mov  byte ptr ds:0CF2h, 0     ; 子刻 = 0
        call sub_19377                ; 季節漸變
        call sub_13E11                ; ★ 每「時」的世界更新
```

換月會 fall through 到換日、再到換時，所以**同一個小時裡月結先跑**。

## 2. 順序反過來的後果

不是「差一拍」：`sub_13E11` 的第二件事是**預備兵維持費累加**
（勢力記錄 `+0x1A`，[`../re/08`](../re/08-hourly-update.md) §2），
而月結的第四步是**把它歸零**。

順序反過來的話，月初那一小時輪到的勢力會先累加一次、再被歸零——
那一次累加**憑空消失**，而且資金在月結的第一步多扣了同一筆。

同局面 196 年 5 月 1 日 15 時的勢力 7：原版 `+0x1A` ＝ 46、remake ＝ 0，
資金跟著差 46。

## 3. 改了哪一支

`internal/state/state.go` 的 `TickMap`：`ev.Clock.Month` 為真時先跑月結，
月結的最後才呼叫 `hourly`。換月一定也換時，所以那裡不必再問一次
`ev.Clock.Hour`。

## 4. 驗證

`tools/parity_ck.sh` 的拍 3,100（5/1 15 時）與 3,250（5/2 9 時）：
**據點、勢力、軍團、全域四張表全部 0 個 byte**。

## 5. 派發開出玩家視窗時不要提早返回

同一支的第二個順序問題。`sub_13E11` 的第一個呼叫是 `sub_131AE`
（事件派發），而原版的玩家視窗是 **modal**：開完視窗、玩家回應之後，
`sub_13E11` 的其餘三步（侵攻財政檢查、預備兵維持費、外交官）
**照樣跑完**，每時勢力游標因此照樣推進一格。

remake 的 `hourly` 原本在 `w.diplomacy != nil || w.funding != nil` 時
提早返回——而世界暫停已經由 `tick` 開頭那道閘負責了
（`w.pending`／`w.diplomacy`／`w.funding` 任一非 nil 就整個不前進）。
再擋一次的後果是**那一小時的游標憑空少推一格**。

同局面拍 3,291 的外交三選一就是這樣差出來的：拍 3,400（5/3 2 時）
原版游標 11、remake 10。接上之後那一拍的四張表全部 0。

## 6. 未解

- 拍 3,420 的據點 129 上昇值與防災值各差 1——那是拍 3,415 的調兵分歧
  之後的下游（亂數流岔開）。根因已經量出來：據點 `+0x18` 在原版是
  **佔用圖抄來的**（增量維護），remake 是掃座標重算，兩者在
  「在那一格編成、從沒走進來過」的軍團上分岔
  （worklist `tick3415-relief-dispatch-occupancy`）。
