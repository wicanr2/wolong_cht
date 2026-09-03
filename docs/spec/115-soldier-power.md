# 115 — 兵的戰力來自統率力，不是士氣

**狀態：CONFORMED。** 算式、戰場類別選欄與接線都上了，單測五支，
戰術九區對拍也重跑過（野戰七區 0 px，`../playtest/58`）。
接上之後那支攻城迴歸從「20 萬幀不結束」變成**第 967 幀結束**——
擋路的不是這條規則，是驗收 fixture 少帶子圖塊表（規格 116）。

- 日期：2026-09-02
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_19AF4`（布陣）、
  `sub_19B6D`（一般部隊）、`sub_19B40`（大將）、`sub_19C13`（取主將能力）、
  `sub_1B6BC`（近戰）——[`../re/78`](../re/78-soldier-power-from-command.md)
- 推論等級：**confirmed（靜態，原始 bytes 交叉解碼）**
- 相關：[`61`](61-soldier-initial-hp-from-morale.md)（開場體力）、
  [`80`](80-duel-opening.md)（單挑開場氣勢也讀戰力）

## 1. 原版做什麼

布陣時每個兵的戰力（記錄 `+0x18`）寫死一次，之後不再變：

```
一般兵：戰力 = ((統率 + 適性) × 3 + 兵種係數) ÷ 4
大將：  戰力 = (武力 × 2 + 適性) × 2
        體力 = max(70, (武力 × 4 + 50) × 軍團士氣 ÷ 100)

兵種係數（seg000:9C0F，bytes 1E 04 0C 00）：1 → 30、2 → 4、3 → 12
適性 = 主將的 +0x0E／+0x0F／+0x10，由**戰場類別**選（re/78 §2.1）
       類別 0 攻城／1 陸上／2 海上——原版自己的標題字串就是這三個
```

近戰有兩支，**分工在 `loc_1B5A1`**（[`../re/11`](../re/11-tactical-battle.md) §5.16–§5.17）：
擋路的對方不是大將走前者、是大將走後者。**remake 兩支都接了**
（`internal/rules/tactical/damage.go` 的 `attackCollision` → `meleeHit`／`hitGeneral`）：

```
sub_1B618：命中值 = rand(0..127) + 攻方戰力      ；≥ 0x46(70) 才命中
           傷害   = 攻方戰力（有利 +0x40、突擊 +0xC8，都飽和到 255）
sub_1B6BC：命中率 = 9.77% + (1 − 9.77%) × min(24, 攻 − 守) ÷ 128
           傷害   = max(1, 攻方戰力 ÷ 8)
```

⭐ **兩支的尺度差很多**，所以「戰力接錯」在 `sub_1B618` 那一支特別明顯（大將那一支的傷害是 ÷8，藏得住）：
戰力填成士氣（100–200）時命中值恆 ≥ 70、傷害 100–200 ⇒ **一擊必殺**。

⚠ **`sub_19B6D` 裡的騎兵戰場修正是死碼**（`cmp al, 1` 而 `al` 已經是
兵種 × 18，[`../re/78`](../re/78-soldier-power-from-command.md) §3.1）。
照抄原版就是不要那兩條。

## 2. remake 現在怎麼做

`internal/rules/tactical` 有逐槽的欄位（`Side.SquadPower [6]int`、
`Side.LeaderPower`、`Side.LeaderHP`），`Deploy` 用它們布陣；
`Side.Power` 退成「沒填時的預設」。
`internal/state/tactical.go` 的 `deploy` 在開戰時把三個欄位填好
（`w.squadPowers`），所以正式路徑不會走到那個預設。

## 3. remake 要怎麼改

| 項目 | 位置 |
|---|---|
| 三條算式 | ✅ `internal/state/soldierpower.go`（`soldierPower`／`leaderPower`／`leaderHP`），單測見 §4 |
| 戰場類別選適性欄 | ✅ `aptitudeIndex(戰場編號)` ＋ `TacticalSetup.Category`（＝ `battle.Category(FieldNumber)`）。⛔ **不是大地圖圖塊**——見 §3.1 |
| 每槽戰力與大將 | ✅ `internal/rules/tactical`：`Side.SquadPower`／`LeaderPower`／`LeaderHP`，`Deploy` 用它們 |
| 主將能力進戰場 | ✅ `w.squadPowers`（軍團編號 ＝ 主將的武將編號）|
| 接線 | ✅ `internal/state/tactical.go` 的 `deploy` |

### 3.1 ⛔ 第一版接錯了：拿大地圖圖塊去比門檻

`aptitudeIndex` 原本吃的是 `TacticalSetup.Tile(x, y)` ＝ 軍團所在格的
**大地圖圖塊值**，拿它去比 `0xC0`／`0xD1`。原版比的是 `byte_10D34`
＝ **戰場編號**，而野戰的編號是 `sub_14B63` 從五格地形算出來的
（[`../re/78`](../re/78-soldier-power-from-command.md) §2.1）。兩邊都錯：

| | 圖塊 | 戰場編號 | 適性欄 |
|---|---|---|---|
| 橋 | `0xCA` ⇒ 看起來像陸上 | `0xD1`–`0xD4` | **海戰**（接錯時取到陸戰）|
| 一般水路 | `0xD8` ⇒ 看起來像海上 | 平原配對表 `0xC0+n` | **陸戰**（接錯時取到海戰）|

⚠ **而且 `Tile` 從來沒有人設**——`grep -rn "Tile:"` 在 `internal/battlesetup`
一筆都沒有，所以 `battleTile` 永遠回退成 `0xC0`：**野戰恆取陸戰適性，
海戰那一欄一輩子用不到**。這一版加了 `TestSquadPowersUsesBattleCategory`
把接線本身釘住——先前沒有任何測試會因為「沒有人設」而變紅。

⭐ **現成的答案就在隔壁**：`battle.Category(p.FieldNumber(node, siege))`
早就用在腳本段的選擇上（`internal/battlesetup`），照抄一行就好。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestSoldierPowerPerTroopType`（三個兵種係數 30／4／12）、`TestSoldierPowerRisesWithCommand`（統率單調，擋「又接回士氣」）、`TestLeaderPowerAndHP`、`TestLeaderHPIsNotMorale`、`TestAptitudeIndexByBattleClass`（含「`0xCA` 當戰場編號是陸上」那一格）、**`TestSquadPowersUsesBattleCategory`**（三個適性欄給三個不同的值，看接線取到哪一個）|
| 迴歸 ✅ | `TestNormalScenarioTacticalBattleTerminates`：第 967 幀結束（0.49 秒）|
| 對原版 ✅ | 戰術九區逐區對拍已重跑（[`../playtest/58`](../playtest/58-parity-retest-20260902.md)）：**野戰七區 0 px、`field` 95 px**，攻城前兩個取樣點與 08-27 同值。沒有回歸 |

## 5. ⭐ 它照出一個 fixture 缺陷

接上去的第一次跑，`TestNormalScenarioTacticalBattleTerminates` 死鎖：
10 萬幀之內兩側剩餘兵數一個都沒變。追下去成因**不在這條規則**——
驗收 fixture 用 `NewFieldFromTiles` 建戰場，那條退路把打破的門
算成 5 層高的方塊，城反而封死（[`116`](116-retreat-cannot-leave-the-city.md)）。

⭐ **一擊必殺把那個缺陷藏了半年**：兵還沒退卻就死了，退卻的尋路一次都沒被走到。
把戰力接對之後兵活得久、退卻的人數跳一個量級，缺陷才浮出來。

## 5.1 影響評估：這不是等價改寫

| | 接上之前 | 接上之後 |
|---|---:|---:|
| 一般兵的戰力 | ＝ 軍團士氣，**100–200** | 統率 10、適性 5、騎兵 → `(15 × 3 + 30) ÷ 4` ＝ **18** |
| `sub_1B618` 的命中值 | `rand(0..127) + 200` ⇒ **恆命中** | `rand(0..127) + 18` ⇒ 約 **59%** |
| `sub_1B618` 的每次傷害 | **100–200** ⇒ 一擊必殺 | **18** ⇒ 約 11 下才倒 |

⭐ **戰鬥慢了一個數量級**，而且勝負從「士氣誰高」變成「統率誰高」。
這是**朝原版靠**的改動，不是平衡調整。戰術對拍已經重驗過
（§4，[`../playtest/58`](../playtest/58-parity-retest-20260902.md)）；
推廣片裡兩場戰鬥的長度也會跟著變（[`71`](71-promo-live-capture.md)）。

同一場攻城的結果也跟著翻面：接對之前是**攻方勝、剩 566 點**（第 1,034 幀），
接對之後是**守方勝**（第 967 幀）。

## 6. 未解

| 項目 | 現況 |
|---|---|
| 海戰適性實際被取到過沒有 | 算式與接線都對了，但**沒有跑過一場橋上的野戰**。要驗得讓兩支軍團在圖塊 `0xCA` 那一格遭遇 |
| 地形類型 3–7 各是什麼地形 | `cs:982Fh` 的範圍已攤開（[`../re/05`](../re/05-battle-selection.md)），但「類型 3 ＝ 山地／丘陵」那一欄的標籤是從圖塊外觀推的，沒有機器碼出處 |
