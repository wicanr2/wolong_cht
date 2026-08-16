# 47 — 據點易主之後，舊主留在那一格的軍團調頭回家

**狀態：CONFORMED。** 名單怎麼收、什麼時候調頭、退不了怎麼辦都有機器碼出處，
remake 已實作並有單測。
⭐ **調頭排在「遷都」之後、「滅亡判定」之前**——順序不是細節，
因為 `sub_1487B` 找的是**新首都**的方向。

- 日期：2026-08-17
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 分析時的 `KI.EXE.i64` SHA-256：`8a8fd7d528e0498000fd04282300588b637fbcb8aa48deb09242f1f41f532691`
- 工具：**IDAPython** `tools/ida_dump.py`
- 出處：`sub_14CF3`（易主）、`sub_14C72`（收名單）、`sub_14DA4`（調頭）；
  `sub_1487B` 的規則見 [`46`](46-post-battle-retreat.md) §2
- 推論等級：**confirmed**（四支都逐行讀過）

## 1. 名單是開打前收的

攻城的入口 `sub_14ADE` 先在堆疊上開 256 byte，再呼叫 `sub_14C72`：

```asm
00014C72  bx = 2240h / cl = 7Fh          ; 軍團表，127 筆，每筆 40h
00014C7B  cmp byte [bx], 80h   / jb →    ; 還活著
00014C80  cmp ax, [bx+12h]     / jnz →   ; 同一個 Y
00014C85  cmp dx, [bx+10h]     / jnz →   ; 同一個 X
00014C8A  cmp ch, [bx+1]       / jnz →   ; 同一個勢力（守方）
00014C8F  mov [bp+0], bx / inc si        ; 收進名單
…
00014C9D  mov [bp+0FEh], si              ; 筆數放在名單後面
00014CA2  筆數 ＝ 0 ⇒ STC                ; 沒有守軍 ⇒ 打城兵（sub_14F8A）
…                                        ; 否則挑最強的一支當代表出戰
```

⭐ **名單是「同一格、同一勢力、還活著」的全部軍團**，
出戰的只有最強的那一支。所以這一條處理的是
**疊在同一格上、沒被捲進那一場的守軍**。

## 2. 易主的順序

`sub_14CF3(al ＝ 新主, si ＝ 據點記錄)`：

```asm
00014CF5  xchg bh, [si+1]         ; 換旗，bh ＝ 舊主
00014CF8  mov [si+1Ah], bh        ; 據點記錄記下舊主
00014CFC  舊主 ＝ 18h（無主）⇒ 跳過下面整段
00014D01  call sub_14D63          ; 內政官歸來（TALK #68）
00014D0A  dec byte [bx+23h]       ; 舊主據點數 −1
00014D0D  call sub_14DF0          ; ★ 遷都（TALK #30）；回 carry ＝ 這一家完了
00014D11  名單筆數 ≠ 0 ⇒ call sub_14DA4   ; ★ 調頭
00014D1C  carry ⇒ call sub_14FCE          ; 滅亡：127 名武將逐一處置
00014D2A  inc byte [bx+23h]       ; 新主據點數 +1
```

⭐ **調頭夾在遷都與滅亡判定中間。** 排在遷都之後，是因為
`sub_1487B` 讀的是勢力記錄 `+0x03`——那時已經是**新首都**了。

### 2.1 這個順序寫反了看不出來

遷都會呼叫 `sub_14502`（`syncCorpsAfterCapitalChange`）：
把「目標還是舊首都」的軍團一律改掛新首都。所以先調頭再遷都，
軍團的目標會被那一步救回來，**兩條路在這裡收斂**。

remake 照原版的順序寫，但**不要以為有測試在擋它**——
跑過負對照：把調頭提前，`TestFallenCapitalRedirectsTowardTheNewCapital`
照樣綠。這一條要靠讀 `sub_14CF3` 維持，不是靠測試。

## 3. `sub_14DA4`：算一次，套給整份名單

```asm
00014DA7  mov al, [si+1]                 ; ★ 這時 si 還是據點記錄 ⇒ al ＝ 新主
00014DAA  mov si, [bp+0]                 ; 名單的**第一支**
00014DAE  call sub_1487B
00014DB2  jb → loc_14DDC                 ; ★ 找不到退路
00014DC0  cl = [bp+0FEh]                 ; 筆數
loc_14DC6:
00014DC9  mov [si+20h], bl               ; 目標 ＝ 回家的下一站
00014DCC  mov [si+14h], dx
00014DCF  mov byte [si+0Bh], 1           ; 計時器 ＝ 1，下一個 tick 就起步
00014DD3  or  byte [si], 2               ; 下一步要重算
00014DD6  inc bp / inc bp / loop         ; 名單裡的每一支
loc_14DDC:
00014DE5  call sub_1291A(al ＝ 新主)      ; ★ 主將擲下場，同樣走完整份名單
```

⭐ **`sub_1487B` 只算一次**（用名單的第一支），結果套給整份名單；
失敗也是整份一起失敗。這不影響正確性——名單裡的軍團同勢力、同一格，
而 `sub_1487B` 的輸入只有這兩樣，算出來必然相同。

## 4. remake 實作

| 項目 | 作法 |
|---|---|
| 名單 | `redirectFallenCityCorps` 現掃：`Alive && Faction == 舊主 && Node == 那一格`。**等價**——原版的名單是開打前收的，死掉的那幾支在 remake 已經 `Alive = false` |
| 逐支算 vs 算一次 | remake 對每一支各算一次 `nextHopHome`。輸入相同 ⇒ 結果相同（§3）|
| 順序 | `capture` 改成「換旗 → 舊主據點數 −1 → 遷都 → 調頭 → 滅亡 → 新主 +1」，與 `sub_14CF3` 逐行對齊 |
| 下一站 | `nextHopHome`（[`46`](46-post-battle-retreat.md) §2）|
| 退不了 | `corpsPerishes`：軍團消失、勢力軍團數 −1、主將擲一次下場。**與戰敗壞滅同一個出口**（原版也是同一支 `sub_1291A`）|
| 起步 | `March` ＋ `Timer = 1` |

## 5. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestFallenCityCorpsRetreatOneHop`：疊在那一格上、沒出戰的守軍目標變成回家的下一站 |
| 單元測試 | `TestFallenCityCorpsWithNoRetreatPerish`：退不了的走壞滅同一個出口，進 `ev.Destroyed` |
| 單元測試 | `TestFallenCapitalRedirectsTowardTheNewCapital`：首都被打下來時，那一格上的守軍最後朝新首都走。**驗的是結果不是順序**（§2.1）|
| 長跑 | `cmd/wlsim` 5 年 60 個月，不變量不違反 |

## 6. 未解

| 項目 | 現況 |
|---|---|
| `sub_14D63` | 「內政官因為據點被攻陷而歸來」（TALK #68）——remake 還沒把內政官送回去 |
| `[si+1Ah]` | 據點記錄記下舊主，remake 的 `OwnerRecorded` 是同一格但語意沒逐位元對過 |
