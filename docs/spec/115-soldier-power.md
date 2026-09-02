# 115 — 兵的戰力來自統率力，不是士氣

**狀態：READY，卡在規格 116（退卻走不出城）。**
算式已經實作並有單測（`soldierPower`／`leaderPower`／`leaderHP`），
**接線刻意留著沒接**：接上之後兵活得久、同時退卻的人數跳一個量級，
會穩定撞上退卻走不出城的死鎖（§5）。

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
適性 = 主將的 +0x0E／+0x0F／+0x10，由戰場類別選（re/78 §2）
```

近戰有兩支，**remake 接的是 `sub_1B618`**（`internal/rules/tactical/damage.go`）：

```
sub_1B618：命中值 = rand(0..127) + 攻方戰力      ；≥ 0x46(70) 才命中
           傷害   = 攻方戰力（有利 +0x40、突擊 +0xC8，都飽和到 255）
sub_1B6BC：命中率 = 9.77% + (1 − 9.77%) × min(24, 攻 − 守) ÷ 128
           傷害   = max(1, 攻方戰力 ÷ 8)
```

⭐ **兩支的尺度差很多**，所以「戰力接錯」在 `sub_1B618` 那一支特別明顯：
戰力填成士氣（100–200）時命中值恆 ≥ 70、傷害 100–200 ⇒ **一擊必殺**。

⚠ **`sub_19B6D` 裡的騎兵戰場修正是死碼**（`cmp al, 1` 而 `al` 已經是
兵種 × 18，[`../re/78`](../re/78-soldier-power-from-command.md) §3.1）。
照抄原版就是不要那兩條。

## 2. remake 現在怎麼做

`internal/rules/tactical` 已經有逐槽的欄位（`Side.SquadPower [6]int`、
`Side.LeaderPower`、`Side.LeaderHP`），`Deploy` 也會用它們；
`Side.Power` 退成「沒填時的預設」。

**缺的只有 `internal/state/tactical.go` 那一段接線**——目前三個欄位寫 0，
於是 `Deploy` 退回 `Power`（＝軍團士氣）。理由見 §5。

## 3. remake 要怎麼改

| 項目 | 位置 |
|---|---|
| 三條算式 | ✅ `internal/state/soldierpower.go`（`soldierPower`／`leaderPower`／`leaderHP`），單測見 §4 |
| 戰場類別選適性欄 | ✅ 同檔的 `aptitudeIndex`：攻城恆取攻城適性，野戰看圖塊 `0xC0`／`0xD1` 兩道門檻（[`../re/78`](../re/78-soldier-power-from-command.md) §2.1）|
| 每槽戰力與大將 | ✅ `internal/rules/tactical`：`Side.SquadPower`／`LeaderPower`／`LeaderHP`，`Deploy` 用它們 |
| 主將能力進戰場 | ✅ `w.squadPowers`（軍團編號 ＝ 主將的武將編號）|
| **接線** | ⛔ `internal/state/tactical.go` 的 `deploy` 還沒接，卡在 [`116`](116-retreat-cannot-leave-the-city.md) |

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 ✅ | `TestSoldierPowerPerTroopType`（三個兵種係數 30／4／12）、`TestSoldierPowerRisesWithCommand`（統率單調，擋「又接回士氣」）、`TestLeaderPowerAndHP`、`TestLeaderHPIsNotMorale`、`TestAptitudeIndexByBattleClass` |
| 對原版 | ⛔ **戰術九區逐區對拍要重跑**（[`../playtest/40`](../playtest/40-tactical-parity.md)）——接線接上之後才做 |

## 5. ⚠ 為什麼還沒切：撞上 [`116`](116-retreat-cannot-leave-the-city.md)

實測接上去之後 `TestNormalScenarioTacticalBattleTerminates`
（濮陽攻城，固定種子 17）**穩定死鎖**：10 萬幀之內兩側剩餘兵數
一個都沒變。成因是城裡與城牆上的兵退卻時尋路找不到出城的路，
**那是既有缺陷**——先前被「一擊必殺」蓋住（[`116`](116-retreat-cannot-leave-the-city.md) §5）。

接線收在 `internal/state/tactical.go` 的 `deploy`：三個欄位目前寫 0，
`Deploy` 於是退回 `Power`（＝士氣）。修好 116 之後把那三行換成
`w.squadPowers(corps, siege, tile)` 的三個回傳值即可。

## 5.1 影響評估：這不是等價改寫

| | 現在 | 接上之後 |
|---|---:|---:|
| 一般兵的戰力 | ＝ 軍團士氣，**100–200** | 統率 10、適性 5、騎兵 → `(15 × 3 + 30) ÷ 4` ＝ **18** |
| `sub_1B618` 的命中值 | `rand(0..127) + 200` ⇒ **恆命中** | `rand(0..127) + 18` ⇒ 約 **59%** |
| `sub_1B618` 的每次傷害 | **100–200** ⇒ 一擊必殺 | **18** ⇒ 約 11 下才倒 |

⭐ **戰鬥會慢一個數量級**，而且勝負從「士氣誰高」變成「統率誰高」。
這是**朝原版靠**的改動，不是平衡調整——但它會讓現有的戰術對拍數字
（[`../playtest/40`](../playtest/40-tactical-parity.md)）全部需要重驗，
也會改變推廣片裡兩場戰鬥的長度（[`71`](71-promo-live-capture.md)）。

**所以這一份停在 READY**：算式沒有疑義，缺的是 [`116`](116-retreat-cannot-leave-the-city.md)，
以及切下去之後連同重跑對拍。

## 6. 未解

| 項目 | 現況 |
|---|---|
| 野戰／水戰的分界 | 攻城那一格 confirmed（據點編號恆 < `0xC0`），野戰與水戰的圖塊門檻是強證據（[`../re/78`](../re/78-soldier-power-from-command.md) §2.1）|
| `sub_1B618` 與 `sub_1B6BC` 的分工 | remake 只接前者。命中率公式兩支不同，要先確認玩家看到的是哪一支 |
