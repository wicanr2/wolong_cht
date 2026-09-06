# 150 — 停戰與請求協助的兩道前置閘：先派外交官、同一件事不重複提

**狀態：CONFORMED。** 原版在**選完勢力的那一刻**就檢查兩件事，不通過就跳訊息
並回去重選：**那個勢力要有我方的外交官**，而且**同型的使者不能已經在路上**。
remake 的規則層有這兩個條件，但**外交官那一條是反的**（要求「沒有外交官」），
而且兩條都只在最後入列時無聲回 false，玩家看不到原因。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）
  - `sub_165EF`（`000165EF`）：外交官閘，**停戰與協助共用**
  - `sub_16605`（`00016605`）：停戰的重複閘
  - `sub_1676F`（`0001676F`）：協助的重複閘
  - `sub_1304E`（`0001304E`）：事件佇列查詢
  - 呼叫端 `sub_164F1`（停戰）、`sub_16623`（請求協助）
- 推論等級：**confirmed**（五支逐行讀完 ＋ 原版擷取 `orig-help7`／`orig-cease55`）
- 驗證：[`../playtest/95`](../playtest/95-diplomacy-preconditions.md)（兩張 `map` 只剩游標）
- 相關：[`140`](140-status-message-box.md) §1.1（#6／#8／#7）、
  [`../re/25`](../re/25-message-variants-and-personnel.md)（外交官 `+0x2A`）、
  [`143`](143-general-duty-field.md)（任命時同時寫勢力 `+0x2A`）

## 1. 原版做什麼

```asm
sub_164F1（停戰提案）:            sub_16623（請求協助）:
  cx = 6 / sub_18853               cx = 8 / sub_18853
loc_16506:                       loc_1663E:
  call sub_17906  ; 選勢力          call sub_17906  ; 選**協助**勢力
  jb  離開                          jb  離開
  call sub_165EF  ; ★ 閘一          call sub_165EF  ; ★ 閘一（同一支）
  jb  loc_16506   ; 回去重選        jb  loc_1663E
  call sub_16605  ; ★ 閘二          call sub_1676F  ; ★ 閘二
  jb  loc_16506                     jb  loc_1663E
  …成立…                            cx = 7 / sub_18853   ; 換成 #7
                                    call sub_17906  ; 再選協同進攻的對象
```

### 1.1 閘一 `sub_165EF`：那個勢力要有我方的外交官

```asm
sub_165EF:
    cmp byte ptr [bx+2Ah], 0FFh   ; bx ＝ 剛選的那個勢力的記錄
    jz  short loc_165F7           ; ★ 等於 0FFh（沒派人）→ 拒絕
    clc / retn                    ;   不等於 → 通過
loc_165F7:
    al = 93h / cx = 37h / call sub_18810   ; TALK #55
    stc / retn
```

TALK **#55**「`\3` 勢力仍未派遣任何人．．．」

⭐ **是「要有」不是「不能有」。** 勢力記錄 `+0x2A` 是**我方派駐在那裡的
外交官**（[`143`](143-general-duty-field.md) §2：任命時同時寫），
所以這一條的意思是「沒有管道就談不成」。

⚠ **敵對提案沒有這道閘**（`sub_16405` 走 `sub_16475` 的局勢評估），
只有停戰與協助有。

### 1.2 閘二：同型的使者不能已經在路上

```asm
sub_16605（停戰）:                sub_1676F（協助）:
    si = bx / al = 6                 si = bx / al = 7
    dx = 0FFFFh                      dx = 0FFFFh
    call sub_1304E                   call sub_1304E
    jnb loc_16615   ; ★ 找到 → 拒絕  jnb loc_1677F
    clc / retn                       clc / retn
loc_16615: cx = 49h → TALK #73    loc_1677F: cx = 4Ah → TALK #74
    stc / retn                       stc / retn
```

`sub_1304E` 掃事件佇列（`cs:word_10D20` 起，4 B 一筆、上限 `0x400`），
比對 `[bx] == ax`——`al` 是事件型別、`ah` 是勢力編號（由 `si` 的記錄位移
`>>6` 得到）。`dx = 0FFFFh` 表示兩個附加欄位都不比。

TALK **#73**「遵照命令，已派遣停戰使者前往 `\3`。」
TALK **#74**「遵照命令，已派遣使者前往 `\3` 請求協助。」

⭐ **訊息是肯定句，作用卻是拒絕**——它在說「已經派過了」，
所以這一次不再受理。看字面會以為成功了。

## 2. remake 要怎麼改

| 項目 | 現況 → 作法 |
|---|---|
| 外交官條件 | ⚠ `internal/state/events.go` 的 `QueuePlayerCeasefire`／`QueuePlayerCooperation` 寫成 `Diplomat != noFaction → return false`——**反了**。改成 `== noFaction → return false` |
| 檢查時機 | 原版在**選完勢力那一刻**就查，不通過跳訊息並回清單。remake 只在最後入列時回 false，玩家看不到原因。`cmd/wlgame/advise.go` 的清單確定分支加上兩道閘 |
| 協助那條查的是**協助勢力** | 兩道閘都在第一張清單（選協助勢力）之後，不是第二張 |
| 訊息 | #55／#73／#74，`\3` 代入君主姓名（[`119`](119-talk-marker-fields.md) §1 的色 `0x0C` ＋ 定寬三格）|
| 敵對提案 | **不加**——原版沒有 |

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 對原版 ✅ | [`../playtest/95`](../playtest/95-diplomacy-preconditions.md)：停戰與請求協助各一張，四區 0 px、`map` 只剩游標 95 px |
| 單元測試 | `TestPlayerDiplomacyProducers`（`internal/state`）：**沒外交官要回 false、派了才過**——這一支原本把錯的行為釘住了 |
| 單元測試 | `TestAdviseAllyGateShowsTalk`：UI 層選到沒派人的勢力跳 #55 且清單留著 |
| 突變測試 | 把條件改回 `!=`，`TestPlayerDiplomacyProducers` 要變紅 |

## 4. 未解

| 項目 | 現況 |
|---|---|
| `sub_1304E` 的 `dx` 附加欄位 | 這兩個呼叫點都傳 `0FFFFh`（不比），別的呼叫點傳什麼還沒逐一讀 |
