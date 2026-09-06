# 95 — 停戰與請求協助的前置閘：抓到一個**反向**的條件

**狀態：通過。** 兩張逐區比對，`banner`／`command`／`minimap`／`faction`
全 0 px，`map` 只剩原版自己畫的滑鼠游標（95 px）。
路上抓到一個規則層的反向條件：**remake 要求「那個勢力沒有外交官」，
原版要求「有」**。

- 日期：2026-09-06
- 規格：[`../spec/150`](../spec/150-diplomacy-preconditions.md)
- 原版側：`…;sclick:352,15;stap:48,47;steps:200000;<move:48,200 × N>;press;
  steps:400000;stap:200,112;steps:400000;stap:200,112;steps:900000;shot:…`
  （N ＝ 1 停戰／2 請求協助）
- remake 側：`-advise-target -advise-pick-row 1|2 -advise-list-row 0`
  （`-advise-list-row` 本輪新增）

## 1. 結果

| 進言那一項 | 選到誰 | 原版跳什麼 | `map` |
|---|---|---|---:|
| 1 停戰提案 | 孫策（沒派外交官）| TALK #55 | **95** |
| 2 請求協助 | 孫策（沒派外交官）| TALK #55 | **95** |

兩張都是「訊息掛著、勢力一覽留著、那一列反白、狀態列還是 #6／#8」。

## 2. ⭐ 抓到的：條件是反的

`internal/state/events.go` 的兩支 producer 寫著

```go
w.Factions[target].Diplomat != noFaction || … { return false }
```

也就是「**那個勢力有外交官就不准提**」。原版 `sub_165EF` 是

```asm
cmp byte ptr [bx+2Ah], 0FFh
jz  → TALK #55 ＋ 拒絕      ; ★ 沒派人才拒絕
clc / retn                  ;   派了人就通過
```

勢力記錄 `+0x2A` 是**我方派駐在那裡的外交官**
（[`../spec/143`](../spec/143-general-duty-field.md) §2：任命時同時寫），
所以那一條的意思是「**沒有管道就談不成**」。

⚠ **這個 bug 有單元測試護著。** `TestPlayerDiplomacyProducers` 明確設
`Diplomat = noFaction` 再斷言成功——測試把錯的行為釘住了。
⭐ **只有拿原版跑同一步才分得出來**：兩種寫法在測試裡都自洽，
畫面上也都「有反應」，差別只在**哪一種情況會被擋**。

## 3. 第二道閘：訊息是肯定句，作用卻是拒絕

`sub_16605`（停戰）與 `sub_1676F`（協助）查事件佇列裡有沒有同型的事件，
**有就拒絕**，跳的訊息卻是

> #73「遵照命令，已派遣停戰使者前往`\3`。」
> #74「遵照命令，已派遣使者前往`\3`請求協助。」

看字面會以為成功了，實際是「已經派過一次了，這次不受理」。
remake 的規則層本來就有這個去重，只是同樣沒有訊息、也不在選完的那一刻查。

## 4. 順帶修的 fixture

`-advise-list-row` 一開始接的是 `pickListRow`，**一點反應都沒有**——
進言的清單自己處理確定，不走 `confirmListSelection`。
抽出 `confirmAdviseSelection()` 之後 fixture 與真實操作走同一條路。

⭐ 這與 [`92`](92-personnel-assign-parity.md) §2 是同一個教訓的兩面：
**fixture 要走真實流程**，否則擺出來的是一個真實操作走不到的狀態。

## 5. 未解

| 項目 | 現況 |
|---|---|
| 協同進攻的對象選完之後 | #7 那一張已經有了（[`102`](102-help-second-step.md)），**再選下去**（成案／被拒）還沒拍 |
| `sub_1304E` 的 `dx` 附加欄位 | 這兩個呼叫點都傳 `0FFFFh`（不比），別的呼叫點還沒逐一讀 |
| 原版擷取裡的滑鼠游標 | 兩張都是 95 px |
