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
| 截圖 | `tools/parity_shot.sh … -save-file SAVE-FIELD.DAT -load-slot 0 -shot-frames 400`：remake 同一份存檔照自然流程停在**遭遇訊息**上，與 `playtest/55` 的 `e4.png` 文字一字不差、肖像同一張、框的位置相同（差別只有鏡頭與日期，見 §5）|

## 4.5 遭遇訊息（2026-08-29 補上）

原版進戰術畫面之前先跳一則訊息（`sub_14EB9`／`sub_14F58`），`cx` 就是 TALK 索引：

| 情境 | `cx` | 訊息 |
|---|---|---|
| 野戰 | `1Dh` | #29「{1}大人的兵馬，遇上{1}的兵馬了！！」|
| 攻城・玩家是攻方 | `1Ch` | #28「{1}大人的兵馬，向{2}進攻了！！」|
| 攻城・玩家是守方 | `1Bh` | #27「{1}的兵馬，向{2}進攻過來了！！」|
| 攻城・空城（自動判定後）| `1Ah` | #26「{2}受到{1}兵馬的攻擊，被攻陷了！！」（§5 未接）|

⭐ **野戰那一則有兩個 `{1}`**，而且**第一個是玩家那一方的主將**：
原版玩家守方那條路在 `sub_14EB9` 之前先 `xchg si, di`（`0x14E9F`）。
攻城兩則沒有這一步，`{1}` 一律是攻方。變數依序取值的機制見
[`106`](106-message-box-reporter-portrait.md) §3。

remake：`state.encounterNotice` 產生 `TalkNotice`，掛在 `CorpsEvent.TalkNotices`，
由 `World.tick` 併進 `Event.TalkNotices`；桌面版**訊息還在時不畫戰場**
（`cmd/wlgame/main.go` 的 Draw 與 `updateMessageOnly`），按掉才換成戰術畫面。

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~遭遇訊息本身~~ | **做了**（§4.5）|
| 攻城打空城之後的 #26 | `sub_14ED7` 的 `cmp bx, 4200h` 那條路先自動判定，攻方贏（`al == 0`）才呼叫 `sub_14F71` 跳 #26「{2}受到{1}兵馬的攻擊，被攻陷了！！」。remake 的空城攻城走 `resolveCorpsBattle`，**那一則還沒接** |
| 遭遇當天的日期差一天 | 同一份存檔：原版在解凍約 4 秒後於 **196年4月17日** 跳訊息；remake 從載入（第 1 幀）到第 240 幀都停在 **4月16日**，遭遇也發生在 4月16日。存檔本身是 4月16日，所以差的是「原版那 4 秒推進了一天、remake 沒有」——載入後戰略時鐘的推進速率要另外量（`docs/spec/34` 的累加器 vs 原版 216 tick／日）|
