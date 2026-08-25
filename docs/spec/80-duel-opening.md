# 80 — 開戰單挑：挑戰、拒戰、應戰、回合互嗆、決著

**狀態：CONFORMED（2026-08-25）。狀態機在
`internal/rules/tactical/duel.go`，單元測試齊，b0（挑戰幀）對拍
回到 field 0.05%＝原版游標（[`../playtest/43`](../playtest/43-field-battle-parity.md)）。
機制全部有機器碼出處（[`../re/74`](../re/74-battle-opening-duel.md)），
傷害不需要新公式——單挑期間兩側大將隊靜止，體力變化來自
戰場上其他部隊的普通攻擊。**

- 日期：2026-08-25
- 出處：`sub_1A1C5`／`sub_1A2E8`／`sub_1A34F`／`sub_1A398`／`sub_1A3C3`／
  `sub_1A298`／`loc_1A23F`（`docs/re/74`）；`sub_19A33` 的 `byte_1D34B`；
  隊長命令跳表 `funcs_1A7E1[8] = nullsub`（`docs/re/11` §5.8）
- 推論等級：流程與門檻 **confirmed**（逐條反組譯）；
  `sub_1A34F` 的氣勢公式 confirmed（含亂數項）；
  拒戰分支有實機對照（`playtest/43` 的 b0–b3）

## 1. 觸發閘（`byte_1D34B`）

戰場初始化 `sub_19A33` 依**戰場編號**分類：

| 戰場編號 | `byte_1D34B` | 意義 |
|---|---:|---|
| `< 0xC0`（攻城＝據點編號）| 0 | 無單挑開場 |
| `0xC0`–`0xD0`（野戰，平原配對表前段）| **1** | **開場走單挑流程** |
| `≥ 0xD1`（其餘野戰＋水域）| 2 | 無單挑開場 |

主迴圈開頭 `sub_1A1C5` 等 50 tick，`byte_1D34B == 1` 才繼續。

## 2. 氣勢評估（`sub_1A34F`，帶亂數）

每側各算一次：

```
core = 大將隊兵數（單位記錄 +0x18）× 大將體力（+0x03）
threshold = max(0, 大將武術×3 − 大將統率) ÷ 2      ；武術、統率取自武將記錄
if rand(0..7)+8 > threshold: core = 0               ；武人氣質門檻
氣勢 = core + rand(0..7)×256
```

高側為「強側」（同分不換，攻方在前）。

## 3. 流程

```
強側氣勢 < 0x12C0（4800）           → 無事發生，正常開打
強側 ≥ 0x12C0：
  強側喊組 0x1B7、大將隊命令＝8、大將移到單挑位、等 40 tick
  弱側 < 0x12C0 或 < 強側一半：
      弱側喊 0x1B9（拒戰）、等 20、強側回 0x1CC → 兩隊命令歸 0，正常開打
  否則（應戰）：
      弱側喊 0x1B8 → 進入回合迴圈：
      round r ＝ 0..（情境碼上限 4）：
        互嗆 pair：r==0 → 0x1BA／0x1BB；否則 0x1BC+(r−1)×4
                   （兩大將體力差 < 20 再 +2 取「勢均」pair）
        對打段（sub_1A298）：計時 0x50；前 0x20 tick 純等；
          之後每 tick 1/8 機率把**兩大將傳到同一個隨機格**
          （x = rand&0xF + 0x18、y = rand&7 + 0x1C）
        任一大將體力 < 0x46（70）→ 決著；否則 r++ 回到互嗆
決著（loc_1A23F）：
  體力低的一方為敗方；敗方目前命令≠5（退卻）才喊 0x1CC、命令歸 0
  等 20 → 勝方喊 0x1CD → 等 20 → 兩隊命令歸 0，戰鬥照常繼續
```

## 3.1 開場序凍結全場、大將騎出（實機定案）

`parity-field14` 的 b0–b3（同一次錄影，跨約六秒）證明兩件
反組譯沒讀出來的行為：

1. **挑戰到拒戰交鋒期間，兩軍完全不動**（b0 對 b1／b2 的 field 差
   0／291 px）。原版把整段開場當 blocking sequence 跑；
   remake 以 `duelOpeningFreeze()` 對應——等待／挑戰／拒戰四個相位
   跳過所有 `updateSoldier`，應戰進入回合迴圈後恢復正常更新。
2. **強側大將是逐格騎向單挑位，不是瞬移**——b0（挑戰喊話當下）
   它還在原位，b3（拒戰交鋒）已在半路。remake 以 `duelRideOut`
   每 tick 走一格對應；應戰時弱側就位用傳送（**推定**，
   b3 是拒戰路徑，沒有應戰側的實機參照）。

另一個實作面的定案：`sub_1A34F` 的「大將隊兵數」對應 remake 的
**場上＋待機**（`Reserve[0]` ＋ 活著的第 0 隊）——只算場上那 8 個
永遠到不了 0x12C0 門檻，單挑一次都不會觸發。

## 4. 命令 8 ＝「單挑中」：整隊靜止

隊長命令跳表 `funcs_1A7E1[8]` 與隊員表 `funcs_1A827[8]` 都是
**nullsub**——大將隊（隊長＋隊員）在命令 8 期間不移動、不攻擊，
位置完全由 §3 的傳送控制。**體力變化因此不需要任何單挑專用公式**：
戰場上其他部隊照常打（普通兵的鎖定包含大將），
單挑位在戰場中段（x 24–39、y 28–35），大將吃到的是一般攻擊。

## 5. remake 實作

| 項目 | 位置 |
|---|---|
| 單挑狀態機 | `internal/rules/tactical/duel.go`：`stepDuel` 掛在 `Battle.Step()` 開頭 |
| 輸入 | `DuelInput{FieldNumber, Martial, CommandStat}`——`startBattleTalk`（`cmd/wlgame/battletalk.go`）在綁定戰鬥時餵一次 |
| 台詞輸出 | `Battle.TakeDuelTalks()` 回 `[]DuelTalk{Side, Group}`；`pumpDuelTalks` 换成 TALK 索引（八變體照大將的 `TalkVariant`）掛上對白框 |
| 命令凍結 | 開場序 `duelOpeningFreeze()` 凍全場；之後命令 8（`Duel`）在 `updateSoldier` 是 no-op（nullsub 等效） |
| 位置 | 挑戰期 `duelRideOut` 逐格；應戰 `duelFace`、對打段 `duelTeleport` 傳送 |
| 對白版面 | 照 [`../re/74`](../re/74-battle-opening-duel.md) §4.1（既有 `drawBattleTalk`）|

先前 `startBattleTalk` 的挑戰段（只出 0x1B7／0x1B9 一輪、開場即出）
已由本狀態機取代——喊話時刻回到原版的 tick 50。

## 6. 驗證

| 方式 | 內容 | 結果 |
|---|---|---|
| 單元測試 | 閘（編號 0xC0–0xD0 才觸發）、氣勢公式（含歸零門檻）、拒戰／應戰分支、決著的退卻跳過、命令 8 凍結 | `internal/rules/tactical/duel_test.go` 全綠 |
| 對拍 | b0（挑戰幀，remake 於 tick ~52 截）field 0.05%＝原版游標、其餘八區七 PASS＋minimap 0.18% | [`../playtest/43`](../playtest/43-field-battle-parity.md) |

⚠ 氣勢帶亂數（§2），同一存檔每次跑挑戰側與拒戰／應戰路徑都可能不同
——對拍截圖要重試到 rng 落在與原版相同的路徑（挑戰側＝攻方）。

## 7. 未解

| 項目 | 現況 |
|---|---|
| `word_1D311 += 6` | 疑似喊話框位置位移，未驗；remake 未實作 |
| 組 `0x1B6`／`0x1C0`–`0x1CB` 的實際回合台詞 | 索引公式已定，逐組逐變體未抽驗 |
| 單挑期間玩家能不能對大將隊下令蓋掉命令 8 | 未讀；remake 暫不允許 |
| 應戰時弱側大將怎麼就位 | remake 用傳送（推定）；實機參照只有拒戰路徑（b3），應戰側的就位方式沒截到 |
| 對打段兩軍是否照常互打 | 依 `sub_1A298` 的讀法推定照常（remake 照做）；未在實機錄到應戰全程 |
