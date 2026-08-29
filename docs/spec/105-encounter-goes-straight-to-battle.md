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
| 攻城・空城（自動判定後）| `1Ah` | #26「{2}受到{1}兵馬的攻擊，被攻陷了！！」|

⭐ **野戰那一則有兩個 `{1}`**，而且**第一個是玩家那一方的主將**：
原版玩家守方那條路在 `sub_14EB9` 之前先 `xchg si, di`（`0x14E9F`）。
攻城兩則沒有這一步，`{1}` 一律是攻方。變數依序取值的機制見
[`106`](106-message-box-reporter-portrait.md) §3。

⭐ **空城那一則只在玩家是守方時發**。`sub_14ED7` 分兩條路，兩條都問
`cmp bx, 4200h`（`ds:4200h` 是城兵用的臨時軍團 ＝ 城裡沒有駐守軍團），
但處置相反：

| 路徑 | 玩家 | 空城時 |
|---|---|---|
| `loc_14EF1`（`[di+1] == 玩家`）| 守方 | `sub_15130` 自動判定 → `and al, al` → `al == 0`（攻方贏）才 `sub_14F71` 發 #26 |
| `loc_14F2B`（`[si+1] == 玩家`）| 攻方 | `jz loc_14EE7` ——判定完就回，**一則訊息都不發** |

#26 的兩個變數照 `sub_14F71` 推堆疊的次序取：先 `di`（據點）後
`ax = [si+2]`（攻方主將），對應字串裡 `{2}` 先出現、`{1}` 後出現。
`sub_14F58`（#27／#28）推的順序相反，也與那兩則的字串順序一致——
**取值跟著標記在字串裡出現的次序走，不是跟著標記編號**（[`106`](106-message-box-reporter-portrait.md) §3）。

remake：`state.encounterNotice` 產生 `TalkNotice`，掛在 `CorpsEvent.TalkNotices`，
由 `World.tick` 併進 `Event.TalkNotices`；桌面版**訊息還在時不畫戰場**
（`cmd/wlgame/main.go` 的 Draw 與 `updateMessageOnly`），按掉才換成戰術畫面。
空城那一則在 `state.fightGarrison`——城裡沒有守軍軍團時走的是它，
不是 `resolveCorpsBattle`（後者的守方索引恆 ≥ 0）。
測試 `TestEmptyCityFallReportsTalk26` 兩個方向都驗（玩家守方有訊息、玩家攻方沒有）。

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~遭遇訊息本身~~ | **做了**（§4.5）|
| ~~攻城打空城之後的 #26~~ | **做了**（§4.5）：`state.fightGarrison` 在攻方攻下城、而且那是玩家的城時發 #26 |
| 遭遇當天的日期差一天 | **量到剩 2 個子刻**（§6）。原本記的「時鐘推進速率」不是成因——**接觸在第幾個子刻與速度檔無關**，節流只改牆鐘秒數不改 tick 數。缺的是原版接觸 tick 的一手數字 |

## 6. 日期差一天：差的是兩個子刻，不是時鐘速率

`workplace/parity/SAVE-FIELD.DAT` 的前 8 byte 是
`10 1e 03 0a 04 00 c4 00` ＝ **196年4月16日 10時 子刻 3**。

照 `clock.Advance` 逐格推（子刻 0–8 共 9 階、時 1–24）：
子刻 3 走到 8 要 5 個 tick，第 6 個 tick 進位到 11時；
之後每個時整 9 個 tick，11時 → 24時 是 13 個時 ＝ 117 個 tick，
第 123 個 tick 到 24時 子刻 0；再 9 個 tick，**第 132 個 tick 才換日**。

remake 這一側的接觸時間可以直接量——`workplace/parity/fieldtrace`
用行軍層逐 tick 印位置：

```
tick    11 軍團  35 (253,108)      ← 呂布軍每 24 個 tick 走一格
tick    34 軍團  35 (252,108)
…
tick   130 軍團  35 (248,108)      ← 踩進軍團 39（夏侯惇）那一格
```

**接觸在第 130 個 tick，換日在第 132 個**——remake 停在 4月16日是自洽的，
差的只有兩個子刻。原版顯示 4月17日，代表它的接觸落在第 132 個 tick 之後。

⭐ **原本把成因記成「載入後戰略時鐘的推進速率」是錯的。**
接觸發生在第幾個子刻由行軍節奏決定，而
[`34`](34-speed-steps.md) 的節流器只改「一個子刻花多少牆鐘時間」，
不改 tick 數——**任何速度檔跑到接觸都是同一個 tick**，
所以速率不可能造成日期差。

剩下的問題窄成一句：**原版從讀檔到接觸多花的那兩個以上的 tick 花在哪。**
最可能的位置是改造存檔配方（[`../playtest/43`](../playtest/43-field-battle-parity.md) §2）
特意設的 `+0x00 |= 2`「下一步要重算」——原版讀檔後自己重查道路表，
那段重算要不要佔 tick 沒有量過。**還沒改程式**：兩個 tick 的補償要有
原版的一手數字才動，照現況反推只會把偏差固定下來（`CLAUDE.md` §10）。
