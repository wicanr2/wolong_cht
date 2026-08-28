# 105 — 遭遇時沒有「戰鬥指揮／委任」選單：直接進戰場

**狀態：CONFORMED（2026-08-29 機器碼 ＋ 實機兩條證據，已實作與單測）。**

- 日期：2026-08-29
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_14E5C`（野戰遭遇）與
  `sub_14ED7`（攻城遭遇），`tools/ida_dump.py` 全文；實機 [`../playtest/55`](../playtest/55-encounter-menu-parity.md)。
  「戰鬥指揮／委任／解體」的決定點在**行軍指示**（[`39`](39-march-order-menu.md)）。
- 推論等級：**confirmed**

## 1. 原版做什麼

```
sub_14E5C:                       ; 野戰；攻城 sub_14ED7 同形
  al = 玩家勢力
  cmp al, [si+1] / jz 玩家是攻方
  cmp al, [di+1] / jz 玩家是守方
  al = 1 / call sub_15130        ; 都不是玩家 → 自動判定
玩家是攻方:  test byte ptr [si], 4 / jnz 自動判定   ; 玩家那一方委任中
             call sub_14EB9 / call sub_11B5A         ; ★ 直接進戰術畫面
玩家是守方:  test byte ptr [di], 4 / jnz 自動判定
             call sub_10CDE（喇叭）/ sub_14EB9 / sub_11B5A
```

**兩條路都沒有呼叫選單引擎 `sub_193E9`**。實機也一樣：遭遇訊息
「夏侯惇大人的兵馬，遇上呂布的兵馬了！！」按掉之後**下一個畫面就是戰場**
（`playtest/55`）。「委任」是行軍指示三選一（TALK #76）那一刻寫進軍團 `+0x00`
位元 2 的，遭遇當下只讀那個位元。

## 2. remake 先前做錯了什麼

`state.EncounterChoice` ＋ 桌面的「遭遇戰：戰鬥指揮／委任」視窗 ＋ 手機的遭遇 modal，
是 remake 自己加的一層，原版沒有。`re/09` §2 寫「只有玩家捲進去才會出現選單」
是把「進不進戰術畫面」的分支讀成了一個選單。

## 3. 改動

| 層 | 內容 |
|---|---|
| `internal/state` | `fight()`：`wantsTactical` 成立就 `beginTactical`，**沒有中間狀態**；刪 `EncounterChoice`／`PendingEncounter`／`ChooseBattleCommand`／`ChooseBattleDelegate`；指紋少一欄（`spec/69`）|
| `cmd/wlgame` | 刪遭遇視窗（`updateBattleChoice`／`drawBattleChoice`／`battleChoiceLayout`）；`updateBattle` 第一幀武裝開戰喊話；刪旗標 open-battle-choice、encounter-choose（`-open-battle`／`-open-siege` 直接開戰場，載存檔照自然流程就會進戰場）|
| `internal/ui/phone` | 刪 `modalEncounter`；戰場在 `tickBattle` 看到 `PendingBattle` 就自己建視圖 |
| `cmd/wlsim`、`internal/battlesetup`、`workplace/parity/*` | 同步拿掉遭遇分支 |
| 差異 | 無（這一條是把 remake 多出來的東西拿掉）|

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `internal/state`：玩家軍團沒委任 → `fight` 之後 `PendingBattle != nil`；委任 → 直接有戰果、不開戰場（`state_test.go` 的遭遇測試改寫）|
| 實機 | `playtest/55`：訊息 → 按掉 → 戰場，中間沒有任何選單 |
| 截圖 | `tools/parity_shot.sh … -save-file SAVE-FIELD.DAT -load-slot 0 -shot-frames 400`：remake 同一份存檔照自然流程停在戰場第一拍 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 遭遇訊息本身 | 原版先跳「{1}大人的兵馬，遇上{1}的兵馬了！！」（TALK #29 那一組）要按鍵才進戰場；remake 直接進戰場，**沒有這一則**。要補得接 `messages` 佇列並在清空後才開戰場 |
