# 12 — 停戰說服訊息索引：`#190`–`#198`

**狀態：三變體槽位與停戰說服這條索引路徑已證實；事件 6／7 的次要呼叫已定位，
但 formatter 參數契約與完整可見語意仍未知。**

- 日期：2026-08-09
- 輸入：`workplace/ida/dosv/KI.EXE.i64` 與同目錄 `KI.EXE`
- `KI.EXE` SHA-256：`FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`
- 工具：`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）
- 位址基準：以下 `000xxxxx` 是 IDA 線性位址；沒有把它和 EXE 檔案偏移混列

## 1. 入口與對話基底

`sub_13830`（`00013830`）保存呼叫端的 `CX`、`AX`、`DL`，並以 `CX` 作訊息
基底。停戰說服入口的呼叫端把 `CX` 設為 `0096h`；若玩家要求說明，
`00013882`–`00013889` 先顯示目前理由，再把 `[BP+0]` 加 `10h`，呼叫
`sub_13B5A`（`00013B5A`）。

`sub_13B5A` 的已證實資料流：

```asm
00013B5A  mov  byte ptr [bp+5], 0
00013B5E  mov  al, 5
00013B60  call sub_13B7E
00013B63  push ax
00013B64  xor  ah, ah
00013B66  add  ax, [bp+0]
00013B69  inc  ax
00013B6A  mov  cx, ax
00013B6C  call sub_13CDC
00013B6F  pop  ax
00013B70  call sub_13BA9
```

這表示使用者選到的理由／結果值由 `AX` 帶進 `sub_13BA9`，而不是由訊息表
本身隨機挑一則。

## 2. `sub_13BA9`：先選三組的起點

`sub_13BA9`（`00013BA9`）在一般分支計算：

```asm
bl = AL * 9 + AH * 3 + 6
CX = [BP+0] + BX
call sub_13C99
```

對「敵が他国侵攻中」這個理由（`AL=2`），`AH=0/1/2` 產生相對於說服
子流程基底 `CX=00A6h`（訊息索引 #166）的三個起點：

| `AH` | 相對偏移 | TALK.DAT 槽位組 | 語意（僅此索引路徑） |
|---:|---:|---|---|
| 0 | `0018h` | `#190`–`#192` | 沒有收到交戰報告 |
| 1 | `001Bh` | `#193`–`#195` | 允許停戰 |
| 2 | `001Eh` | `#196`–`#198` | 建議觀望 |

這裡的「語意」由日中原文並列確認；`AH` 的完整欄位名稱仍未知。

## 3. `sub_13C99`：每組按索引取三格

`sub_13C99`（`00013C99`）的關鍵指令是：

```asm
00013C9E  mov  dx, 0
00013CA1  mov  bx, 5
00013CA4  xor  ah, ah
00013CA6  mov  al, [si+1Eh]
00013CA9  cmp  al, 3
00013CAB  jb   .keep
00013CAD  sub  al, 3
.keep:
00013CAF  add  cx, ax
00013CB1  mov  al, [si+1]
00013CB4  call sub_1075B
```

在本段變體索引為 `0/1/2` 的情況，`[SI+1Eh]` 直接加到 `CX`，所以選取是
連續三格的**索引直取**，不是隨機抽句。值達 `3` 時另有 `sub al,3` 分支；
這份證據不把那個欄位命名成已知狀態，也不把所有值宣稱為只限 `0..2`。

## 4. 松崗槽位校訂

校訂前的繁中把停戰允許句放到 `#192`，使交戰情境顯示答非所問的回應。
`translations/corrections.json` 現在採「既有譯文重排」：

| 槽位 | 校訂後來源 |
|---|---|
| `#192` | 重用原 `#190` 的「正和他國交戰？我可沒接到那種報告啊！」 |
| `#193`、`#194` | 重用原 `#192` 的「現在停戰的話，就能賣個人情給{3}吧，好，准許停戰。」 |
| `#195` | 重用原 `#193` 的「嗯，現在的話，{3}也不至於拒絕。好，准許停戰！」 |

這只修復訊息索引與已存在的 `{3}` 參數，不宣稱中文行寬、語氣或畫面排版
已與 PC-98 原版完全相同。原始 `TALK.DAT` 不改寫；可重跑的輸出由
`tools/talkdat.py correct` 產生。

## 5. 證據檔

- `workplace/ida/dosv/func-sub_13830.txt`
- `workplace/ida/dosv/func-sub_13B5A.txt`
- `workplace/ida/dosv/func-sub_13BA9-current.txt`
- `workplace/ida/dosv/func-sub_13C99-current.txt`
- `workplace/ida/dosv/func-sub_1075B-current.txt`
- `workplace/ida/dosv/func-sub_1084A-current.txt`

`sub_1075B`／`sub_1084A` 的 TALK.DAT 格式化器細節另記於
`docs/formats/01-talk-dat.md`；本文件只記外交對話如何形成槽位索引。

## 6. 事件 2／3 玩家收到請求時的前置報告

事件 dispatch 對玩家成立的合作／停戰提案，會在三選一前先走下列兩個已證實的
TALK 基底：

| 事件 | IDA 線性位址 | `CX` | TALK | 原始句型 |
|---|---:|---:|---:|---|
| 3 停戰 | `sub_138C7` `000138C7` | `0168h` | #360 | `{3}前來請求停戰……意下如何？` |
| 2 協力 | `sub_138E6` `000138E6` | `0175h` | #373 | `{3}前來請求協助……意下如何？` |

`{3}` 經 `sub_13C3D`／formatter `00010904` 取請求方勢力的君主姓名；不能把事件字高
直接當成可顯示的文字。remake 的 `beginDiplomacy` 已將這兩個基底寫成
`TalkNotice`，通知關閉後才進入三選一；`TestQueuedDiplomacyChoiceTalkNotices` 與
`TestDiplomacyTalkExpansionUsesOriginalRequestMarkers` 已在 Docker／Xvfb 通過。

這一節只封閉前置報告。`sub_13902` 後續的接受／拒絕變體、數值視窗、PC-98 欄寬／游標、
逐頁動畫與完整玩家路徑仍待獨立 oracle，不能由 #360／#373 外推。

## 7. `sub_13902` 的選項與主要結果

`sub_13902` 在 `base+4+[choice]` 顯示玩家建言；事件 handler 隨後以
`sub_13C3D` 的 `CX=0x2B`／`0x2F` 顯示主要結果：

| 事件 | 建言 0／1／2 | 主要結果 response 0／1／2 |
|---|---|---|
| 停戰（base #360） | #364／#365／#366 | #43／#44／#45 |
| 協力（base #373） | #377／#378／#379 | #47／#48／#49 |

指定外交金額為 0 時回傳 response 0；正值且不超過初始要求時 response 1，結果句使用
`{7}`；拒絕回傳 response 2。超額輸入回傳 response 3，但 `sub_13C3D` 先以 response 2
顯示破裂句，再呼叫 `sub_13DC9(AL=1Eh)` 扣信賴度 30，且不執行外交收尾。remake 已把
選項句與主要結果依序排入 TALK modal，並以 `TestDiplomacyTalkIndicesMatchRaw13902Branches`
與 `TestDiplomacyAndFundingAmountOutcomeBounds` 固定索引、Trust 副作用與超額分支。
#367–#372／#380–#385 的 AH／信賴度次要回覆仍未解，不可視為已完成的原版對話流程。

## 8. 事件 6／7 次要呼叫：已證實索引，拒絕猜接 formatter

### 8.1 `sub_13C3D` 的第二次呼叫

本節使用唯讀 DOS/V 反組譯輸入：

- `workplace/ida/dosv/KI.EXE.asm` SHA-256：
  `FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`
- `workplace/orig/dosv/KI.EXE` SHA-256：
  `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`
- 工具：`ida-pro-9.4-ver2:uidfix-v1`（IDA Pro 9.4）；以下位址均為 IDA 線性位址。

`sub_13C3D`（`00013C3D`）的第一個 `sub_18810` 呼叫會先保存
`AX/CX/DX/DI`，再把 `DI` 改成 `SP`，所以主要結果 #43–#49 可以取得這次呼叫
建立的 formatter 堆疊。返回後它恢復 `DI/CX/AX`，再執行：

```asm
cmp  al, 3
jz   loc_13C8B          ; 先顯示 response 2，再扣信賴度 30
and  ah, ah
jz   loc_13C93
cmp  al, 2
jz   loc_13C93
add  cx, 1Dh
mov  al, 93h
call sub_18810
```

這個第二次呼叫沒有再次 `push` formatter 參數，也沒有再次把 `DI` 設成 `SP`。
因此可以確認「有條件的次要 TALK 呼叫」與 `CX+0x1D`，但不能把它當成第一次
呼叫的同一組 `{1}`／`{2}`／`{3}`／`{7}` 參數。

直接 caller 的索引如下：

| caller | 原始 `CX` | 第二次 `CX+0x1D` | 可直接證實的 TALK index |
|---|---:|---:|---:|
| `sub_13327`（事件 6） | `0x2B` | `0x48` | #72 |
| `sub_13388`（事件 7） | `0x2F` | `0x4C` | #76 |
| `sub_13262`（事件 3 停戰） | `0x2B` | `0x48` | #72 |
| `sub_13220`（事件 2 協力） | `0x2F` | `0x4C` | #76 |

原始 PC-98 `TALK.DAT` 的 #72 是「{2}で暴動が発生しました。」（暴動發生於
`{2}`），#76 則是「戰鬥指揮／委任／解體」的選單文字。兩者都不能在目前
`TalkNotice{City,Faction,General,Amount}` 沒有第二次 formatter 堆疊的情況下安全
展開成一般通知。鄰近的 #73／#77 不在上述四個直接 caller 的 `CX+0x1D` 算術結果中；
既有交接若把這段範圍寫成「#72／#73、#76／#77」，本節將其收窄為：#72／#76
是已證實的直接結果，#73／#77 仍是未定位的未知，不得拿來補接。

### 8.2 `AH` 不是信賴度，也沒有證明是城市索引

`sub_137D8`（`000137D8`）先呼叫 `sub_13138`（`00013138`）兩次，將兩個方向的
結果寫入 `AH` bit 0／bit 1。`sub_13138` 逐筆掃描勢力／武將表，將輸入勢力與
General 記錄的 `+0x1D`（俘虜／原勢力關係）及 `+0x1C`（目前擁有勢力）比較；
這是俘虜關係旗標的 raw 資料流，不是 `byte_10D00` 信賴度。

`sub_1084A` 的 `\\2` handler（`000108DB` 附近）會從 `SS:[DI]` 取一個 word、
遞增 `DI` 兩次，再把它轉成城市繪圖索引。第一個 formatter 呼叫把 `DI=SP`，
所以 `\\2` 的參數來源可追查；第二個呼叫沿用恢復後的 caller `DI`，卻沒有建立
相同的堆疊 word。這構成「第二次呼叫缺少城市 formatter 參數」的已證實負證據，
而不是「AH 一定就是城市」的證據。

### 8.3 接線決定與 oracle 狀態

本輪不新增 #72／#76 的 `TalkNotice`，不把 `AH` 轉成城市、勢力或其他 state 欄位，
也不把 #73／#77 當成替代訊息。這是 **已證實 raw index／呼叫條件 + 已證實
formatter 缺參數 + 語意未知** 的明確邊界；事件 6／7 次要 TALK 仍列未完成。

曾以 Docker 內的暫時存檔嘗試直接進入 PC-98 事件 oracle，但啟動流程仍停在
`NEW GAME`／讀檔選擇畫面，沒有取得有效的事件畫面；這次嘗試是 inconclusive，
不作為原版行為證據，也沒有修改 `workplace/orig/`。

## 9. 2026-08-10 接線勘誤：事件 6／7 raw 次要 TALK

上一節的「不新增 `TalkNotice`」是接線前的研究結論；本輪已在保留同一證據界線的
前提下完成 raw consumer，並以 state／呈現單測固定結果：

- 事件 6 `sub_13327` 在 `AH != 0` 且 response 不是 2／3 時，追加
  `Index=0x48`（TALK #72）、`Secondary=true`。
- 事件 7 `sub_13388` 在同一條件下，追加 `Index=0x4C`（TALK #76）、
  `Secondary=true`、`NoPortrait=true`；#76 是直接文字／選單樣式，不套玩家君主肖像。
- `sub_13C3D` 第二次呼叫仍沒有重建第一次的 `DI=SP` formatter stack。DOS/V TALK
  #72 的 `\\2` 不能由目前 state 安全推導城市，因此 `cmd/wlgame` 缺少該 marker 時
  整則 fail-closed；不會顯示半句，也不會把 `AH` 猜成城市／勢力／信賴度。
- `TestQueuedDiplomacySecondaryTalkConditions` 固定事件 6／7 的條件、索引、順序與
  portrait policy；`TestSecondaryTalkUsesCapturedRawFormatterWord`、
  `TestSecondaryTalk76UsesNoPortraitMode` 固定呈現層邊界。

因此本輪完成的是「事件 6／7 次要 TALK 的 raw 索引、條件、佇列順序與安全呈現接縫」；
#72 原版缺失 formatter payload 的不可重現記憶內容仍是明示差異，沒有被猜測填補。

## 10. 2026-08-12 勘誤：#72 不得合成 raw word `0`

本節使用松崗 DOS/V `KI.EXE`（SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`）、IDA Pro 9.4，
位址為 IDA 線性位址；`tools/ida_unresolved_research.idc` 是本輪唯讀匯出器。

`sub_13C3D`（`00013C3D`）的第一個 `sub_18810` 呼叫把 `DI=SP`，故 `\\2` handler
可讀第一組 formatter stack。它返回後先 `pop DI` 還原 caller 值，`add CX,1Dh`，便直接
第二次呼叫 `sub_18810`；第二次沒有新的 `push`，也沒有重設 `DI=SP`。

`\\2` handler（`000108DB`）的 raw 指令開頭是：

```asm
mov ds, cs:word_10D52
mov ax, ss:[di]
inc di
inc di
cmp ah, 0FFh
...
add ax, 4240h
```

因此事件 6 的第二次 #72 呼叫會讀 `SS:[還原後 DI]`。在非玩家回報路徑，`sub_13C3D`
先把 `DI=SI`，這個 runtime stack word 沒有可由持久 World state 重建的資料流；把 Go
零值 `0` 當成一個可用 formatter payload 是錯誤斷言。

- **已證實：** #72／#76 的 index、條件與第二次 handler 的 `SS:[DI]` 讀取。
- **強推論：** 這個 transient word 可能是有意義的 runtime payload；沒有動態 trace 前不能
  把它稱作垃圾或固定城市。
- **未知：** 原版自然執行時 `SS:[DI]` 的值與對應城市。

remake 因而把事件 6 的 `RawFormatterWord` 改為 `-1` 且
`RawFormatterWordValid=false`，與事件 7 一樣讓缺 payload 的訊息整則 fail-closed。
`TestQueuedDiplomacySecondaryTalkConditions` 固定這個邊界；
`TestSecondaryTalkUsesCapturedRawFormatterWord` 只測呈現層對「真正已擷取的」通用 raw payload
是否能畫出，不可被引用為事件 6 已知 payload 的證據。

## 11. 未解

| 項目 | 現況 |
|---|---|
| `AH` 的完整欄位名稱 | 語意由日中原文並列確認，欄位名本身未定（§3）|
| #367–#372／#380–#385 的 AH／信賴度次要回覆 | 未解，不可當成完整的原版對話流程（§8）|
| #73／#77 | 未定位，不得拿來補接事件 6／7（§9）|
| 事件 6／7 次要 TALK 的 formatter 參數契約 | 缺參數且語意未知，維持 fail-closed（§10）|
