# 12 — 停戰說服訊息索引：`#190`–`#198`

**狀態：三變體與停戰說服索引已證實。2026-11-16補證事件2／3／6／7的返還對話；
13138固定CF0，原始caller不會觸發72／76，相關formatter不構成玩家阻塞。**

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

`sub_137D8`兩次呼叫13138後，只有CF1才會設DH的bit0／bit1，最後137F1寫AH=DH。
13138雖掃描Faction／Captor並計數，13167卻是 `cmp ax,ax`，接著固定跳1316D清CF。
所以DOS/V這條原始caller的AH固定0。舊稿把配對統計當成true回傳，已由本頁2026-11-16補證推翻。

`sub_1084A` 的 `\\2` handler（`000108DB` 附近）會從 `SS:[DI]` 取一個 word、
遞增 `DI` 兩次，再把它轉成城市繪圖索引。第一個 formatter 呼叫把 `DI=SP`，
所以 `\\2` 的參數來源可追查；第二個呼叫沿用恢復後的 caller `DI`，卻沒有建立
相同的堆疊 word。這構成「第二次呼叫缺少城市 formatter 參數」的已證實負證據，
而不是「AH 一定就是城市」的證據。

### 8.3 接線決定與 oracle 狀態

本輪不新增 #72／#76 的 `TalkNotice`，不把 `AH` 轉成城市、勢力或其他 state 欄位，
也不把 #73／#77 當成替代訊息。這是 **已證實 raw index／呼叫條件 + 已證實
formatter 缺參數 + 語意未知** 的明確邊界；這是當時尚未追查13138回傳契約的歷史邊界；目前原始caller的AH固定0，不形成玩家缺口。

曾以 Docker 內的暫時存檔嘗試直接進入 PC-98 事件 oracle，但啟動流程仍停在
`NEW GAME`／讀檔選擇畫面，沒有取得有效的事件畫面；這次嘗試是 inconclusive，
不作為原版行為證據，也沒有修改 `workplace/orig/`。

## 9. 2026-08-10 接線勘誤：事件 6／7 raw 次要 TALK

本節曾在Go按俘虜配對新增72／76的raw consumer，但誤認137D8的旗標會被設起。
2026-11-16已由13138原指令與真實事件6／7推翻，舊producer及專用配對helper已移除。
13C3D內的條件、CX+1Dh算術與formatter讀SS:DI仍保留作程式定位，不能當成正常玩家會看見的通知。
`TestQueuedDiplomacySecondaryTalkConditions`現在驗證有配對俘虜也不產生72／76，並核對實際返還ID。
`TestSecondaryTalkUsesCapturedRawFormatterWord`及 `TestSecondaryTalk76UsesNoPortraitMode`只驗證通用
呈現接縫，不作原始事件6／7的觸發證據。

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

只有人工強制AH非零而走進第二次#72時，才會讀 `SS:[還原後 DI]`。在非玩家回報路徑，`sub_13C3D`
先把 `DI=SI`，這個 runtime stack word 沒有可由持久 World state 重建的資料流；把 Go
零值 `0` 當成一個可用 formatter payload 是錯誤斷言。

- **已證實：** #72／#76 的 index、條件與第二次 handler 的 `SS:[DI]` 讀取。
- **強推論：** 這個 transient word 可能是有意義的 runtime payload；沒有動態 trace 前不能
  把它稱作垃圾或固定城市。
- **已證實不構成玩家阻塞：** 原版13138固定CF0，事件6／7不會自然走進該次要呼叫。

早期remake曾把事件6的raw payload標成invalid以略過通知；事件7的76沒有marker卻仍會誤顯示。
2026-11-16改為不建立原始caller不會產生的通知。通用raw formatter仍維持安全參數邊界。

`TestSecondaryTalkUsesCapturedRawFormatterWord` 只測呈現層對「真正已擷取的」通用 raw payload
是否能畫出，不可被引用為事件 6 已知 payload 的證據。

## 2026-11-16：事件3停戰的俘虜返還實測

本節新增範圍已CONFORMED。輸入KI.EXE雜湊與IDA位址空間沿本頁，執行器為dosgolem
`f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`。規格與完整輸入、fixed RNG及收據見
[spec220](../spec/220-c-modal-control.md#2026-11-16事件3停戰俘虜返還對話)。

七案原版13262→13C3D的AH皆0，含單向、雙向及多名俘虜，未觸發72；因此不能依俘虜存在
自行替事件3加暴動訊息。後續事件6／7及13138回傳補證見下一節，已證實相同固定AH0來源。
135ED依記錄順序呼叫150D7，返回玩家的39／45各先顯示37，再由本人raw4、肖像49／24
顯示199h群組的434。49返回勢力3時無玩家訊息，拒絕停戰不呼叫150D7。

Go原先已釋放但漏返還UI，現在以runtime ReleasedGenerals串到事件9共用對話入口，
不新增存檔欄位或猜補formatter。原版正常四槽讀檔多人案也實跑此鏈，與Go同入口完整區塊、
RNG及ID順序相同；直接6時與正常7時各自比較。Go仍在回覆入列前結算，原版在結果確認後，
不稱中間資金／人事寫回相位一致。校訂文字、鏡頭與部分場景邊緣差異仍保留。

## 2026-11-16：事件6／7返還與固定CF零

本節CONFORMED，規格與版本入口見[spec220](../spec/220-c-modal-control.md#2026-11-16外交官回報返還與固定cf零勘誤)。
[原始指令收據](../../workplace/scratch/closeout-envoy-captive-20261116/code-proof.json)核對五函式每條bytes與EXE檔案偏移，
保留13167 `cmp ax,ax`、13169跳轉、1316D `clc`，以及137DF／137E9的jnb。不是把原作者疑似錯字改成cmp ax,0。

原版事件6八案與事件7免費／付費雙向兩案、正常讀檔多人案，11份完整22208-byte區塊、258-byte RNG、
主通知索引、返還ID及存檔讀回相同。所有13C3D入口AH0，無72／76；成功返玩家先37再199h，
拒絕及無外交官不返還，返對方無玩家對話。Go先前既漏返還，還會在事件7多開戰鬥選單76，現已補正。

十條滑鼠104項／42圖、事件3與9共享回歸20項／10圖、冷測／vet與野戰62項／9圖通過。
原版同狀態數值比較固定完整RNG；GUI另用seed1與正常7時驗證畫面報價及實際付款方向，
不把該報價和直接原版6時的另一份RNG條件混比。原版直接/正常時鐘、鏡頭、校訂文字及入列前後寫回相位限制
均留在[整合收據](../../workplace/scratch/closeout-envoy-captive-20261116/integration.json)。

## 2026-11-16：事件2返還先於宣戰

本節CONFORMED，版本、原始指令與READY契約見[spec220](../spec/220-c-modal-control.md#2026-11-16事件2協力選擇的返還與宣戰順序)。
原始13251先13C3D，成功才13258→135ED逐筆150D7；全部返還對白完成後1325D才13526。
玩家原先與侵攻方和平才顯示1A0宣戰。Go既有bool入口丟掉返還ID，呈現層也先排宣戰notice，
本輪改接既有WithReturnees並先排所有返還，再排宣戰，不重寫規則。

原版九案與正常多人讀檔案共十份完整22208-byte區塊、258-byte RNG、返還ID及存檔讀回相同。
玩家39／45各先37通報，再用raw4、肖像49／24說199h；40返受攻方無玩家訊息。
拒絕與超額不返還；超額沿421扣30；原已交戰不重複宣戰。本片Source0、Invader12、Victim2，
名稱取原始資料，編號12是袁術，不憑題材印象映射名字。

九條實際滑鼠147項／77圖及兩錯版拒絕，並抽看原版／Go的37、434、487六框。
直接原版12時、正常原版與Go13時分開；GUI seed1檢查操作及付款，數值精確對拍另有完整raw RNG。
Go仍入列前結算、原版在結果確認後結算，不稱中間人事／資金相位或鏡頭逐像素一致。
來源及限制在[整合收據](../../workplace/scratch/closeout-cooperation-captive-20261116/integration.json)。

## 2026-11-16：中央結果確認後才付款及扣信賴

本節CONFORMED，沿本頁DOS/V KI.EXE及IDA9.4 linear，執行器與READY入口在
[spec220](../spec/220-c-modal-control.md#2026-11-16玩家外交結果確認後才付款及扣信賴)。
13251／13290先13C3D，13C6E的18810中央結果確認後才回13258／13297的135ED付款。
超額同樣等待中央破裂結果，再13C90→13DC9扣信賴。六案每個真實輸入等待前保存完整區塊與RNG，
中央結果期間的資金／Trust與來源相同；正確CS的131AE實跑，沒有視窗stub或量測中注入結果。

Go新增PrepareDiplomacy回傳預覽及一次性續行，桌面以既有afterTalk掛在中央結果之後。
原同步ResolveDiplomacyWithResult仍立即準備及結算，CLI／手機caller語意不變，總取數與結果未改。
原版確認前及結算後的十二份完整22208-byte區塊／258-byte RNG全等，同步等價與重複執行保護通過。
正式GUI六條114項／53圖、有俘虜共享兩條31項／17圖，以及舊binary兩錯版拒絕。

側欄資金位於中央框外，Go依目前Funds重畫；受控只開hudFaction位元的一案顯示，
結果等待時仍50000、右鍵後60000，已抽看兩個畫面與數字區像素。這是受控側欄可見性，
不冒稱正常操作打開側欄或原版同側欄像素對拍；正式binary另有六條正常標題／讀檔／滑鼠證據。

此修正只收斂付款／超額扣信賴的確認邊界。採納骰在軍師回覆之前的相位、逐名俘虜寫回、
1A0確認後的war位元及事件6／7仍不由此代證；較早各節的最終狀態證據保持有效，
不能將本片寫成全部外交中間相位一致。詳見[整合收據](../../workplace/scratch/closeout-diplomacy-confirm-20261116/integration.json)。

## 11. 未解

| 項目 | 現況 |
|---|---|
| `AH` 的完整欄位名稱 | 語意由日中原文並列確認，欄位名本身未定（§3）|
| #367–#372／#380–#385 的 AH／信賴度次要回覆 | 未解，不可當成完整的原版對話流程（§8）|
| #73／#77 | 未定位，不得拿來補接事件 6／7（§9）|
| 事件 6／7 次要 TALK 的 formatter 參數契約 | 僅人工強制分支的formatter來源未定；13138固定CF0使AH0，正常事件6／7不觸發72／76，不阻塞remake |
