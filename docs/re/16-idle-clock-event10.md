# 16 — DOS/V 無輸入自動時鐘與事件 10 關係

**狀態：無輸入時的自動時鐘／軍團行軍已由 IDA `.i64` 證實；事件 10 是該路徑
中的受節流 queue consumer，尚無證據顯示它是時鐘或行軍的 driver。**

- 日期：2026-08-11
- 研究問題：玩家不下命令、滑鼠不移動時，原版是否仍讓日期、軍團與 runtime 物件
  自動前進？事件 10 是否就是這條自動 clock？
- 原始輸入：DOS/V `workplace/orig/dosv/KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 資料庫：`workplace/ida/dosv/KI.EXE.i64`，SHA-256
  `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 工具：IDA Pro 9.4、容器 `ida-pro-9.4-ver2:uidfix-v1`；16-bit DOS/V 線性位址，
  segment base `0x10000`
- 可重生匯出器：[`tools/ida_idle_clock.idc`](../../tools/ida_idle_clock.idc)

## 1. 目前最可靠的結論

使用者描述的「滑鼠不動就自動跑」是正確的，但對應的原版鏈不是事件 10：

```text
sub_11BE0 主迴圈
  └─ sub_11F7F 讀取滑鼠按鍵／座標
       └─ 座標沒有變 → byte_198A3 設 bit 7
            └─ sub_11BE0 看到 byte_198A3 >= 80h
                 └─ sub_11CD0 無輸入更新路徑
                      ├─ sub_13EFD 據點 runtime 更新
                      ├─ sub_125A3 軍團計時／行軍
                      ├─ sub_12459 MCH 物件 timer／動畫
                      └─ sub_11D8E 子刻／時／日／月時鐘
                           └─ 每小時 sub_13E11
                                └─ sub_131AE 事件 queue dispatcher（受 byte_131AD 節流）
                                     └─ queue low byte 0x0A → sub_13496 TALK consumer
```

所以要分成兩句：

1. **無輸入自動 clock／軍團行軍：已證實，是 `sub_11CD0` 的主迴圈路徑。**
2. **事件 10：若 queue 裡已有 `Code.low=0x0A`，會在每小時更新的節流邊界被消費成
   TALK；它本身沒有呼叫時鐘或行軍。**

## 2. `sub_11F7F`：無滑鼠移動的條件

IDA raw instruction（`00011F7F–00012078`）顯示：

```asm
mov ax, 5 / xor bx, bx / call sub_20000   ; 讀按鍵／按鈕狀態 0
mov ax, 5 / mov bx, 1 / call sub_20000     ; 讀按鍵／按鈕狀態 1
mov ax, 3 / call sub_20000                 ; 讀滑鼠座標，回 CX／DX
cmp cx, word_1987E
jnz loc_11FD0
cmp dx, word_19880
jnz loc_11FD0
mov cx, word_19886
mov dx, word_19888
mov bx, bp
or  byte_198A3, 80h
retn
```

座標沒有變時，`byte_198A3` 的 bit 7 被設起；座標有變時，`loc_11FD0` 更新
`word_1987E`／`word_19880`、清除 bit 7，並更新可視地圖範圍。這不是「滑鼠事件
遺失」，而是原版用來判定是否可以讓世界繼續跑的 idle gate。

IDA data refs：

| 原始定位 | 內容 | 證據等級 |
|---|---|---|
| `00011FC8` | `sub_11F7F` 設 `byte_198A3` bit 7 | 已證實 |
| `00011C5B`、`00011C7A` | `sub_11BE0` 檢查 `byte_198A3 >= 0x80` | 已證實 |
| `00011FBE`、`00011FC2` | 無移動時取回 `word_19886`／`word_19888` | 已證實 |
| `00011D13`、`00011D16`、`00011D19`、`00011D1C` | idle path 依序呼叫據點、軍團、物件、時鐘 | 已證實 |

## 3. `sub_11CD0`：真正讓世界自動跑的地方

IDA function `00011CD0–00011D46` 在 `byte_10D2A != 1` 的一般遊戲路徑中直接呼叫：

```asm
call sub_13EFD       ; 據點 runtime／災害檢查
call sub_125A3       ; 每次掃 16 支軍團，計時器到零就 call sub_12662
call sub_12459       ; MCH 物件動畫 timer
call sub_11D8E       ; 推進子刻；進位時推進時／日／月
```

`sub_125A3` 內的 `dec [si+0Bh]`、重載 `[si+1Eh]`、`call sub_12662` 直接證實：
軍團不需要玩家再次移動滑鼠或下新命令，就會沿已存在的目標路徑繼續走。

## 4. 時鐘與事件 10 的真正相對位置

### 4.1 `sub_11D8E`（`00011D8E`）

它先處理子刻 `0–8`；進位到每小時時呼叫：

```asm
call sub_19377       ; 季節漸變
call sub_13E11       ; 每小時世界更新
call sub_11E17       ; 重畫日期
```

月進位才呼叫 `sub_15358`。因此日期流逝的 driver 是 `sub_11D8E`，不是
`sub_131AE` 或 `sub_13496`。

### 4.2 `sub_13E11`（`00013E11`）

IDA xref 顯示唯一 caller 是 `sub_11D8E` 的 `00011DEC`；函式開頭是：

```asm
call sub_131AE
...
call sub_13E65       ; 預備兵維持費
call sub_13E8E       ; 外交官效果
...
mov al, 4
call sub_15E80       ; redraw／狀態畫面更新
```

這證實事件 queue dispatcher 是「每小時更新的子步驟」，但不是每小時都真的取一筆。
`sub_131AE` 先遞減 `byte_131AD`；`sub_12BD9` 初始化它為 `7`，取到一筆後重設為
`0x0A`。因此在佇列尚未到 `0x100` 結尾時，初次取事件是在第 7 次每小時呼叫，
之後每 10 次每小時呼叫取一筆。它不是無輸入主迴圈的入口，也不是軍團行軍的入口。

### 4.3 `sub_131AE`（`000131AE`）

它先遞減 `byte_131AD`；歸零才讀 queue，然後把 `byte_131AD` 設回 `0x0A`。這裡
有兩個容易混淆但不同的 `0x0A`：

| `0x0A` | 意義 |
|---|---|
| `byte_131AD = 0x0A` | queue dispatcher 的節流計數器；不是事件 10 |
| `Code.low = 0x0A` | 事件 10 handler index，經 `funcs_131E8` 進 `sub_13496` |

事件 10 handler 只把 `Code.high`／`Param` 轉成 TALK formatter 呼叫；沒有
`sub_11D8E`、`sub_125A3`、`sub_12662` 或 clock 欄位的寫入。

## 5. 對「事件 10 就是自動 clock」假說的判定

### 已證實的部分

- 玩家不下命令且滑鼠不移動時，原版會自動推進時鐘。
- 已下達目的地的軍團會在這條 idle path 中自動行軍。
- 每小時世界更新會呼叫 queue dispatcher；dispatcher 依 `byte_131AD` 節流，並非每小時
  都消費一筆事件 queue。
- 如果 queue 中恰有低碼 `0x0A`，事件 10 TALK 會在這個自動時鐘期間出現，並可能
  讓玩家看到訊息視窗；這是「事件 10 隨 clock 被處理」。

### 目前不成立的部分

- 沒有 IDA xref 顯示 `sub_13496` 反呼叫 clock／軍團更新。
- `sub_11CD0`、`sub_11D8E`、`sub_13E11`、`sub_125A3`、`sub_13EFD` 的函式鏈中，
  沒有直接 queue writer `sub_12FBF`／`sub_1301C` 呼叫。
- 已檢查的 queue writer caller 沒有可證實的 `Code.low=0x0A`。因此目前不能把
  「無輸入自動跑」實作成 `QueueEvent10` 的副作用。

### 最合理的新模型

```text
玩家停手／滑鼠停住
  → idle clock path 運行
  → 時鐘每小時進入世界更新
  → queue dispatcher 嘗試取事件
  → 若其他路徑曾寫入 0x0A，才顯示事件 10 TALK
```

原版 event 10 producer 仍可能是未建函式、間接寫入或特定劇本資料，但「它就是
時鐘」已被目前 `.i64` 的 caller graph 否定；應改查「哪個自動世界更新或 AI 路徑
在 idle 前／月結時寫入 `0x0A`」，而不是把時鐘本身接到 `QueueEvent10`。

## 6. 對 remake 的影響

remake 的對應層應是：

| 原版 | remake |
|---|---|
| `sub_11CD0` idle world loop | `game.timeRuns` + `idleClockGate` → `World.TickMap` |
| `sub_125A3` 軍團更新 | `World.tickCorps` |
| `sub_11D8E` 時鐘 | `World.Clock.Advance`（由 `World.TickMap`／`World.Tick` 尾端執行） |
| `sub_13E11` 每小時更新 | `World.hourly` |
| `sub_131AE` 受節流 queue dispatch | `World.dispatchQueuedEvent` |
| `sub_13496` event 10 TALK | `TalkNotice{Index, General}` |

目前 `World.QueueEvent10` 仍只應作受控 fixture／劇本注入口；不能用它代替
`World.Tick`。正常 UI 現在以 `idleClockGate` 補足原版的輸入條件：首次座標觀測、游標
位移、滑鼠按鈕或命令 frame 都不進世界更新；下一個游標穩定且無輸入的 frame 才跑
`World.TickMap`。這讓日期、軍團位置、物件 timer 與 queue 依原本順序一起前進，並以
`TestIdleClockGateRequiresStablePointerAndNoCommand` 與
`TestIdleClockDispatchesQueuedEvent10OnHourlyCadence` 固定 UI gate 加 state cadence。
同一畫面的額外 `g.speed` 規則 tick 使用不含物件的 `World.Tick`，維持物件每個可見
map-loop 一次的 cadence。事件 10 producer 仍是獨立的 unknown 邊界，不能與更新順序混淆。

## 7. 2026-08-11 松崗繁中驗證封口

本輪唯一 oracle 是松崗繁中版；`workplace/orig/dosv` 是沿用的資料夾名稱，不把 PC-98
或其他版本差異納入驗收。固定 `seed=17`、速度 1、30 frame 的 Docker/Xvfb `wlgame`
smoke 在加入 gate 後仍輸出
`45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`，證明正常自然
畫面沒有因首個非 idle frame 而漂移。這是 remake 行為與畫面 smoke，不是密碼頁後原版
動態 trace，也不把未知 `0x0A` writer 說成已逆向。
