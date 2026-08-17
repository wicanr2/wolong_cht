# 15 — 事件 10 producer 深度逆向

**狀態：事件 10 dispatcher／consumer／queue writer 已證實；原版自然 producer
仍未知；remake 已增加明確標示的近似自然 producer，並保留 raw fixture 開關。**

- 日期：2026-08-11
- 原始輸入：DOS/V `workplace/orig/dosv/KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 資料庫：`workplace/ida/dosv/KI.EXE.i64`，SHA-256
  `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 工具：IDA Pro 9.4，容器 `ida-pro-9.4-ver2:uidfix-v1`；16-bit DOS/V 線性位址，
  segment base `0x10000`
- 可重生匯出器：[`tools/ida_event10_producer.idc`](../../tools/ida_event10_producer.idc)
- 證據契約：原始函式名、全域名、線性位址與運算元均保留；本文件的語意分級不會
  覆蓋 IDA 名稱。IDC 輸出是由 `.i64` 的函式邊界、code xref、data ref 與原始
  instruction 匯出，不把攤平 `.asm` 當成資料庫關係的替代品。

## 1. 先給結論

事件佇列的低 byte 是 dispatch code。`sub_131AE`（`000131AE`）從每筆 4-byte
record 取出 `Code`／`Param`，以 `Code.low−1` 索引 `CS:funcs_131E8`。表格中的
事件 10 entry 是 `sub_13496`（`00013496`）。

這個 consumer 位於無輸入自動 clock 的下游，而不是 clock driver：`sub_11D8E`
每小時進 `sub_13E11`，再呼叫受 `byte_131AD` 節流的 `sub_131AE`。`sub_12BD9`
初始節流值為 `7`，取到 record 後重設為 `0x0A`；完整的 idle call chain 與第 7／
之後每 10 個每時邊界的節拍見 [`docs/re/16-idle-clock-event10.md`](16-idle-clock-event10.md)。

`sub_13496` 不會產生事件，也不會寫回狀態；它只把 queue record 轉成 TALK 呼叫：

```asm
mov  al, ah       ; Code.high：來源 General／人物欄
mov  ah, 0FFh     ; TALK formatter 的原始 FFxx 形式
mov  cx, dx       ; Param：TALK index／參數
push ax
mov  di, sp
mov  al, 93h
call sub_18810
pop  ax
retn
```

因此「事件 10 的 producer」必須在某處建立：

```text
Code = (General << 8) | 0x0A
Param = 原始 TALK index／參數
```

對保存的 DOS/V `.i64` 做完 queue writer graph、直接 caller、data ref 與事件常數
追查後，沒有找到可證實的低 byte `0x0A` producer。這不是宣稱整個 binary 在所有
間接取址路徑都數學上排除；是已知 IDA 關係圖內的負證據，並以「原版 producer
未知」作限時結論。

## 2. dispatcher 與 consumer：已證實

### 2.1 `sub_131AE`（`000131AE`）

IDA raw instruction 的資料流如下：

1. `byte_131AD` 每次 dispatcher 呼叫遞減；歸零時設回 `0x0A`。這是 dispatcher
   cadence，不是事件 10 的事件碼。
2. `word_10D20` 是 queue segment 內的 byte offset；從 `word_10D56` 指向的
   segment 取 `[bx]`／`[bx+2]`。
3. 讀取後 `BX += 4` 並寫回 `word_10D20`，所以 record stride 是 4 bytes。
4. `AL=0` 直接返回；否則 `BL=AL`、`BX--`、`BX <<= 1`，呼叫
   `CS:funcs_131E8[BX]`。

這條鏈直接證實 `AL=0x0A` 對應 table entry 10；IDA 的 data ref 顯示
`funcs_131E8` 的唯一呼叫端是 `sub_131AE`（`000131E8` 的 table use），而事件
10 handler 的 entry 是 `sub_13496`。

### 2.2 `sub_13496`（`00013496`）

高 byte 只在 handler 內轉為 `FFxx` formatter word，低 byte 事件碼在 dispatch 後
不再傳給 TALK。`DX` 原樣進 `CX`，所以 remake 不能把未知 `Param` 改寫成泛用
文案或重新排序 TALK index；必須保存 raw queue payload。

## 3. queue layout 與 writer

### 3.1 `sub_12FBF`（`00012FBF`）

其保留的 raw 指令區段為 `00012FBF–0001300E`：

```asm
push ds / push ax,bx,cx,dx
mov  cx, ax
...
add  bx, cs:word_10D20
mov  ds, cs:word_10D56
cmp  byte ptr [bx], 0
...
mov  [bx], cx
mov  [bx+2], dx
```

此 helper 以 `AX` 成為 Code、`DX` 成為 Param；`BL=0xFF` 時先走原版亂數選槽
路徑，否則以 `BL` 選 record。它是已知 producer 的共用 queue writer，不能由
helper 名稱反推出呼叫端一定是哪個事件。

### 3.2 `sub_1301C`（`0001301C`）與 `sub_1300E`（`0001300E`）

`sub_1301C` 對 `BX×4` 加上 `word_10D20`，在 `word_10D56` segment 中檢查空槽，
寫入 `[bx]=AX`、`[bx+2]=DX`，並以 `0x400` queue 範圍做完整 256-slot 掃描。
`sub_1300E` 是將 `SI×4` 與 `AH=CH` 組合後呼叫它的 wrapper。這兩支仍保留
原始 offset／segment 運算元；語意只標為 queue writer，沒有臆測它們就是事件 10。

### 3.3 queue 的 data refs

由 IDC `DfirstB`／`DnextB` 輸出：

| 原始定位 | 觀察 | 等級 |
|---|---|---|
| `word_10D20` | `sub_12BD9` 初始化；`sub_12FBF`、`sub_1301C`、`sub_1304E` 讀寫 queue offset | 已證實 |
| `word_10D56` | `sub_100DF` 建立／保存 queue segment；兩個 writer 與 dispatcher 取用 | 已證實 |
| `byte_131AD` | `sub_12BD9` 寫 `7`；`sub_131AE` 寫 `0x0A` | 已證實為節拍 byte；不是事件 10 producer |
| `funcs_131E8` | `sub_131AE` 間接呼叫 | 已證實為 handler table |

`sub_130CB`（原始位址保留）只從 `[bx]` 讀取並回傳關係／狀態資料，沒有呼叫
`sub_12FBF` 或 `sub_1301C`；不能把它誤列為 queue writer。

## 4. caller graph 與事件常數

IDC 以 `RfirstB`／`RnextB` 列出下列直接 caller。這些是 IDA code xref 的原始定位，
不是由函式名猜測：

| writer／helper | 直接 caller（IDA 線性位址） | 已看到的 Code low byte | 等級 |
|---|---|---|---|
| `sub_12FBF` `00012FBF` | `sub_12286` `000122AF`／`000122CC`、`sub_122DB` `00012346`、`sub_12FB1` `00012FBA`、`sub_15715` `0001577E`、`sub_1578F` `000157ED`、`sub_157FE` `00015824` | `0x010C`、`0x020C`、`0x0B`、`0x04`、`0x05`、`0x0D` 等；無已證實 `0x0A` | 已證實的負證據 |
| `sub_12FB1` `00012FB1` | `sub_12D3A` `00012D51`、`sub_12E33` `00012E81`、`sub_12E89` `00012EE7`、`sub_12EFB` `00012F68`、`sub_12F71` `00012FA8` | 呼叫端各自組合 event code；已檢查值不含 `0x0A` | 已證實的負證據 |
| `sub_1301C` `0001301C` | `sub_1300E` `00013017`、`sub_134B1` `00013503`、`sub_15940` `00015981`、`sub_16623` `000166A9` | `0x0C`、`0x09`、`0x07` 等；無已證實 `0x0A` | 已證實的負證據 |

交叉核對幾個有語意的 producer：`sub_15715` 是事件 4、`sub_1578F` 是事件 5、
`sub_15940` 是事件 9、`sub_16623` 是事件 7 玩家路徑；`sub_134B1` 的事件 12
handler 寫入 `AL=0x0C`。這些路徑都沒有把 `AL` 組成 `0x0A`。

## 4.9 分派表整張攤開了：13 個事件碼

`funcs_131E8` 在 `000131F2`，**13 個 word**（IDA 把整段當一個陣列，
所以 `next_head` 從 `131F2` 一跳到 `1320C`）。逐格讀出來：

| 碼 | handler | 內容 |
|---:|---|---|
| 1 | `sub_1320C` | 宣戰（`sub_1351A` 的閘）|
| 2 | `sub_13220` | 協力要請 |
| 3 | `sub_13262` | 停戰交涉 |
| 4 | `sub_132A9` | 撥款要求（內政官）|
| 5 | `sub_132E9` | 撥款要求（外交官）|
| 6 | `sub_13327` | 外交官回報停戰結果 |
| 7 | `sub_13388` | 外交官回報協力結果 |
| 8 | `sub_133EA` | 遷都 |
| 9 | `sub_13485` | 釋放俘虜 |
| 10 | `sub_13496` | 訊息（本文件其餘各節）|
| 11 | `sub_134A6` | 暴風雨的動畫標記 |
| 12 | `sub_134B1` | 火災／暴動 |
| 13 | `sub_13507` | 信賴度變動 |

⭐ **13 個碼、13 個 handler，`internal/state/events.go` 的 `case 1` 到
`case 13` 一一對上**，沒有沒接到的碼。

## 5. 證據分級與未解範圍

### 已證實

- queue record 是 `Code`／`Param` 兩個 word，stride `4`。
- `sub_131AE` 以低 byte 事件碼索引 `funcs_131E8`；低碼 `0x0A` 對應事件 10。
- `sub_13496` 把高 byte／Param 轉成 TALK formatter 呼叫，沒有自然 producer 副作用。
- 已列出的 queue writer 與 direct caller、`word_10D20`／`word_10D56`／`byte_131AD`／
  `funcs_131E8` data refs 均可由同一 `.i64` 重新匯出。

### 強推論

在保存的 IDA function boundary 與直接 xref graph 中，所有已定位的 queue writer
caller 都使用別的事件碼；因此原版自然事件 10 producer 不在這個可見 producer
集合內，或不是靜態可直接追到的呼叫路徑。這足以停止無上限的 `0x0A` 猜測，但
不足以宣稱對所有 `register`／pointer／外部 loader 寫入做完形式化排除。

### 未知

以下來源沒有證據，不能補成事實：未被 IDA 建成函式的 far code、以暫存器或指標
間接寫入 queue segment 的路徑、外部載入器／密碼保護後才出現的劇本資料，以及
事件 10 的自然 TALK index／觸發時序。DOS/V 原版執行目前停在密碼頁，因此沒有
dynamic trace 可以越過這個限制。

## 6. remake 的明確邊界

`World.QueueEvent10` 只實作已證實的 raw 介面：檢查 `general` 與 `talkIndex` 邊界，
在完整 queue 中寫入 `(general<<8)|0x0A` 和原始 `Param`。它是 fixture／劇本注入口，
不是原版自然 producer 的名稱替代。`TestEvent10ProducerWritesRawTalkPayload`、
`TestQueuedEvent10TalkNotice` 固定 payload、滿槽搜尋、有效 General 與 fail-closed
邊界。

## 7. remake 近似自然 producer（substitute，不是原版證據）

使用者要求事件 10 也進入可遊玩的近似 remake，因此新增
`World.produceApproximateEvent10`。它不改寫已證實的原版結論，而是採一條可回溯的
替代規則：月結完成既有 queue producer 後，掃描玩家勢力目前收容的活武將；每月最多
選一名，沿 `General.Timer` 的月度倒數閘，依固定 RNG 邊界近似兩種已知鄰近 TALK：

| 近似條件 | 狀態變更 | raw event 10 |
|---|---|---|
| `rand & 0xFF < 0x20` | 武將轉在野、清除 `Captor`／`Posted` | `(general<<8)\|0x0A`, `Param=0x41`（逃走） |
| `0x20 <= rand & 0xFF < 0x40` | 保留玩家勢力、清除 `Captor`／`Posted` | `(general<<8)\|0x0A`, `Param=0x42`（歸降） |
| 其餘 | 不改狀態、不排事件 | — |

`0x41`／`0x42` 文字來自 DOS/V `sub_15940` 的已知月結俘虜處理鄰近路徑；
「這些文字可用於事件 10」是 **強推論／替代設計**，不是宣稱原版低碼 `0x0A`
writer 已找到。事件寫入後仍由原版對齊的每時 dispatcher 消費，GUI 沿既有
`TalkNotice`／TALK marker／分頁路徑顯示。

近似 producer 在 `LoadScenario` 預設開啟；`World.SetApproximateEvent10(false)` 可讓
raw fixture 只驗證原始 queue／consumer。`TestApproximateEvent10ProducerUsesKnownRawContract`、
`TestApproximateEvent10ProducerIsBoundedAndDisableable` 與
`TestApproximateEvent10ReentersIdleClockConsumer` 固定 payload、狀態邊界、關閉開關與
idle clock 重新進入 consumer 的行為。

本次結論是「深度靜態逆向完成，原版 producer 保持 unknown；remake 另有可關閉、可測試、
不冒充原版的近似自然 producer」；只有取得新的原版密碼答案、可觀測 SAVE／事件樣本或
另一份已知版本 binary，才值得重新開啟原版 producer 研究線。

## 8. 2026-08-11 remake 功能封口與驗證範圍

使用者指定以**松崗繁中版**作為唯一驗證 oracle。現有原始檔工作目錄仍叫
`workplace/orig/dosv`，那只是既有路徑／二進位標識；本輪不再以 PC-98 或跨版本
DOS/V 差異作為事件 10 的阻塞條件。

remake 的事件 10 功能現在由三個明確層次完成：

1. `idleClockGate` 對應 `sub_11F7F` 的游標穩定條件：首次觀測、游標移動、按鈕或
   命令 frame 都不推進世界；下一個游標穩定且無輸入的 frame 才允許 `World.TickMap`。
2. `World.TickMap` 依已證實順序更新據點、軍團、MCH 物件與時鐘；每時才進既有 queue
   dispatcher，因此已下達命令的軍團、日期與已有 queue event 會在玩家停手時自然前進。
3. raw queue consumer 與可關閉的月結 substitute 維持原有界線；近似 producer 不因 UI
   功能封口而被升格為原版 `0x0A` writer 證據。

`TestIdleClockGateRequiresStablePointerAndNoCommand` 固定 UI 輸入閘門，
`TestIdleClockDispatchesQueuedEvent10OnHourlyCadence` 固定第 7 個每時邊界的事件 10
consumer，`TestApproximateEvent10*` 固定 substitute 的 raw payload 與狀態原子性。
固定 `seed=17`、30 frame 的 Docker/Xvfb smoke 仍產生既有自然畫面 hash
`45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。

因此「玩家不下命令、游標不移動時自動跑」這個 remake 功能已完成；原版自然低碼
`0x0A` writer 則維持已記錄的 **unknown** 研究邊界，不能改寫成已證實。

## 9. 2026-08-12 再審：逐一展開所有可見 direct producer

本次以新的暫存 IDA 9.4 資料庫重新開啟同一份松崗 DOS/V `KI.EXE`，輸入雜湊、工具版本
與線性位址基準同文件開頭；沒有修改保存的 `.i64`。`tools/ida_event10_deep.idc` 與
`tools/ida_unresolved_research.idc` 保留可重生的唯讀匯出流程。

| 原始函式／位址 | 直接可見的 queue Code | 判讀等級 |
|---|---|---|
| `sub_12286` `00012286` | `0x010C`、`0x020C` | 已證實 |
| `sub_122DB` `000122DB` | `0x000B` | 已證實 |
| `sub_12D3A`／`sub_12E33`／`sub_12E89`／`sub_12EFB`／`sub_12F71` | 經 `sub_12FB1` 分別組成 `0x08`、`0x02`、`0x03`、`0x01`、`0x01` | 已證實 |
| `sub_15715`／`sub_1578F`／`sub_157FE` | `0x04`、`0x05`、`0x0D` | 已證實 |
| `sub_164F1` → `sub_1300E` | `0x06` | 已證實 |
| `sub_134B1`／`sub_15940`／`sub_16623` | `0x0C`、`0x09`、`0x07` | 已證實 |

`sub_12BD9` 只進行 queue 的月度初始化／搬移：它先設 `word_10D20=0`、
`byte_131AD=7`，複製既有範圍後清空後段，沒有寫入 `0x0A` record。`word_10D56`
的直接 data ref 也只落在 `sub_100DF` 初始化、`sub_12BD9`、兩個 writer
`sub_12FBF`／`sub_1301C` 與 scanner `sub_1304E`。

這使結論更精確而非更寬：**目前 IDA 可見的每條 direct queue-writer 路徑均未寫出
low byte `0x0A`。** 這是已證實的有限負證據；自然 producer 仍可能經由未辨識的間接
位址建構、外部程式或密碼頁後才可觀測的流程寫入，故仍必須標為 **unknown**。remake 的
可關閉 substitute 仍是產品近似，不是這個表格所推出的原版規則。
