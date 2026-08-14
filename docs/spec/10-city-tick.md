# 10 — 據點整備、威脅偵測與求援

**狀態：CONFORMED。整條鏈已實作，並在 PC-98 原版的執行期記憶體上取樣驗過
（`+0x18`／`+0x14` 各 0/192 不符，場上 21 支軍團）。**

- 日期：2026-08-14
- 出處：[`docs/re/44`](../re/44-threat-and-reinforcement-ai.md)（`sub_13EFD`／`sub_13FA9`／`sub_14575`）、
  [`docs/re/40`](../re/40-garrison-relief-request.md)（`sub_14028`／`sub_14057`／`sub_140C9`／`sub_14155`）
- 推論等級：**confirmed**（逐行讀過 ＋ 動態取樣）

## 1. 原版做什麼

`sub_11CD0` 每 tick 呼叫 `sub_13EFD` 一次，游標 `word_10D1E` 指到哪個據點就處理哪個，
處理完 `+= 0x20`，到 `0x1800`（192 × 32）繞回 0。

> **每 tick 只處理一個據點，192 個 tick 掃完一輪。**
> AI 的反應速度被這個掃描週期限制，不是被判斷邏輯限制。
> remake 照抄這個節奏——一次掃全圖會讓 AI 反應快 192 倍。

一個據點輪到時依序做五件事：

1. 求援冷卻 `+0x17` 減 1。
2. `+0x1A ≠ +0x01`（作者填的原主不是現在的所屬）→ 在原主的勢力記錄寫 `+0x17 = 據點編號`。
3. 從單位佔用圖抄 `+0x18`：`es = Y × 24 + cs:word_19872`、`al = es:[X] & 0x7Fh`。
4. **所屬不是中立（`0x18`）才**做威脅偵測與求援（`sub_13F74`）。
5. 無條件做內政（`sub_14194`）與災害 marker（`sub_14269`）。

## 2. 演算法

### 2.1 威脅偵測（`sub_13FA9`）

```
if 據點.+0x1B == 0: return          # 沒有敵方鄰居，連掃都不掃
威脅量 = 0
for 槽 in 據點.+0x1C..+0x1F:
    if 槽 == 0xFF: break
    鄰 = 據點表[槽]
    if 鄰.所屬 == 我方所屬: continue
    if 鄰.所屬 != 中立:
        if 交友度[我方][鄰.所屬] & 0x80: continue   # bit 7 ＝ 和平 → 不算威脅
        受威脅 = true
        威脅量 += 鄰.+0x18
    if 鄰.所屬 == 我方勢力.+0x19:                    # 侵攻目標
        受威脅 = true; 有具體目標 = true
        目標清單 += 槽
據點.+0x14 = 威脅量
據點.+0x00 的 bit 7 = 受威脅、bit 6 = 有具體目標
```

### 2.2 求援與派兵（`sub_14028`／`sub_14057`）

```
if 不受威脅: 據點.+0x17 = 0; return
if 有具體目標 and 據點.+0x18 == 0:  求援(1)          # 貼身威脅，立刻求援
else:
    目標 = 目標清單[亂數 & 3]
    要幾支 = max(0, 據點.+0x14 + 2 − 據點.+0x18)
    if 據點.+0x18 <= 1:
        if 要幾支 == 0 or 據點是玩家的: return       # ← 玩家不走這條
        求援(要幾支)
    else:
        派兵(目標, max(1, 要幾支), 跳過額度 = 據點.+0x18 − 要幾支)
        據點.+0x17 = 0
```

### 2.3 求援（`sub_140C9` ＋ `sub_140B3`）

```
if 據點.+0x17 != 0: return                          # 冷卻中
if 據點是玩家的:
    發訊息 #38「{2}前來請求援軍。」
    據點.+0x17 = 亂數(0..15) + 24
else:
    額度 = max(5, 勢力.資金 ÷ 8192) − 勢力.軍團數
    重複 min(額度, 要幾支) 次： 編一支軍團，目標 ＝ 這個據點
    據點.+0x17 = min(30, (|X − 首都X| + |Y − 首都X|) ÷ 8)
勢力.+0x16 = 據點編號
```

> ⚠ **距離算式的 Y 分量減的是首都的 X，不是 Y。** 兩版執行檔的指令 bytes
> 完全相同，出自 1994 年的原始碼——**照抄，不要修正**
> （[`docs/re/40`](../re/40-garrison-relief-request.md) §4.1）。

### 2.4 派兵（`sub_14155`）與編成（`sub_145C1`）

```
派兵：掃軍團表，只收「人已經在這個據點」的軍團
      還有跳過額度時，每支 25% 機率被跳過（亂數 < 0x40，額度只在跳過時才減）
      位元 2（委任）沒設、或 +0x23 >= 8（待解體）的不收
編成：同勢力、職務 == 0（無職）、+0x11（武力）最大的那一個
      掃 0x7F ＝ 127 次而武將表有 128 筆，**最後一筆掃不到**——照抄
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/threat`（`Scan`／`Requested`／`CorpsCap`／`Budget`／`EnemyMask`／`PlayerCooldown`／`AICooldown`／`Dispatch`）|
| 狀態層 | `internal/state/strategy.go`：`refreshCityThreat`／`relieve`／`requestRelief`／`dispatchGarrison`／`formAICorpsTo` |
| 掃描節奏 | `internal/state/state.go` 的 `tickCity`（游標 `cityCursor`）|
| 差異 | `+0x18` 直接數軍團而不是維護佔用圖計數器——**導出值，記帳會漂**。結果相同（§4 取樣 0/192 不符）|

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `internal/rules/threat/threat_test.go` 12 條（上限兩段在交界處連續、和平位元極性、`+0x1B ＝ 0` 直接跳過、中立不看交友度、bit 6 的分野、距離不對稱沒被改掉、只調在場軍團、跳過額度只在跳過時才減）|
| 單元測試 | `internal/state/ai_probe_test.go`：`TestCityThreatIsRecomputedOnTick`／`TestAICorpsCapFollowsFunds`／`TestPlayerCityAsksForRelief`／`TestReliefOnlyMovesDelegatedCorps` |
| **對原版** | [`docs/playtest/21`](../playtest/21-dosboxx-bridge-sampling.md) §4.2：PC-98 執行期記憶體，場上 21 支軍團時 `+0x18` 與 `+0x14` **各 0/192 不符**（非零 12–13 筆）|
| **對原版** | 同上 §4.1：`+0x00` 低 4 位 192/192，對照讀法只對 12/192 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| `sub_14194`／`sub_14269` | 內政與災害 marker 的細節在別的規格（`docs/mechanics/40`），本規格只保證呼叫順序 |
| 據點換手之後 `+0x00` 低 4 位會不會跟著變 | `sub_1890A` 靜態讀過，動態沒驗——要打下一座城才看得到 |
| 玩家據點求援的喇叭聲（`sub_10CDE`）| 呈現層未接 |
