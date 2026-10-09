# 113：C 完整資源載入、BGM cue 與分段設定

**狀態：CONFORMED。6 支函式的局部原版／C 行為通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 固定 IDA DB SHA-256：`da891546bc51a7e6864f2096a2349a65c178122967797b356cd17a387404e468`
- 工具／位址：IDA Pro 9.4 database linear；file offset 為 linear 減 `0x10000` 加 512
- 入口：[`ida_resource_probe.py`](../../tools/ida_resource_probe.py)
- 前置：[`re/110`](110-c-talk-rendering-restoration.md)、[`re/112`](112-c-choice-selector-restoration.md)
- 契約：[`spec/233`](../spec/233-c-resource-cue.md)

`sub_1F4A2` 用原 DOS open／read／close 讀完整檔案，單次 `0xF000` bytes，滿讀後
將目標 segment 前移 `0xF00` paragraphs。`sub_187AF` 依原 MMAP.MDL／MCH 入口載入。
`sub_10241` 保留 cue gate、index header、offset／length、停止／載入／播放與 mute 順序。
`sub_100DF` 保留兩筆 allocation 與原始資料／字庫分段關係，不用測試手寫地址取代。

## 已證實範圍與重跑

| 原始函式／IDA linear | bytes | 已證實的局部行為 |
|---|---:|---|
| `sub_187AF`／`0x187AF` | 29 | MMAP.MDL 後接 MMAP.MCH；第二個目標段加 `0x800` |
| `sub_1E378`／`0x1E378` | 20 | 完整讀檔成功的重試包裝；原 INT50 分支保留 |
| `sub_1F4A2`／`0x1F4A2` | 59 | 每讀 `0xF000` bytes 前移 `0xF00` paragraphs；短讀／EOF／close 與 `AX=0xFFFF` |
| `sub_10241`／`0x10241` | 129 | cue 相同即返回；停止、header／body、播放與 mute 順序 |
| `sub_102C2`／`0x102C2` | 14 | cue 設 `0xFF`；原 `INT61 AX=0x09F2`；AX 恢復 |
| `sub_100DF`／`0x100DF` | 118 | `0x1D5E`／`0x5356` allocation 成功後原 segment 加法與保存 |

共 369 bytes、165 條原始指令。原始函式名稱、位址與運算元保留；分級索引同時引用本文件與 spec/233。
[C 來源](../../tools/c_recovery/resource.c)、[生成內容](../../tools/c_recovery/resource_generated.inc)、
[驗證收據](c-resource-verification.json)與[函式台帳](c-recovery-status.json)均受版控。

`tools/c_recovery_resource.sh` 在固定 Go／GCC Docker 工具鏈執行。原始資料唯讀，
原版跑 dosgolem guest，C 直接執行 native 運算，兩側使用獨立 DOS／INT61 服務。
驗證器為 `tools/c_recovery_resource_verify.py`，C 由固定 IDA 輸出乾淨重生。

- O0／O2 各 93 組：分段邊界 16、完整真實檔 5、MMAP caller 2、cue 64、stop 5、allocation 1。
- 16 個 cue 各覆蓋重複／不重複與 mute 開／關。比對完整 RAM、四 VGA plane、暫存器、FLAGS、堆疊、函式入口、DOS／INT61 API 前後與 sound 狀態。
- 21 次完整 buffer 比對來源檔並確認 64-byte 尾端守衛。測試堆疊放在獨立段，避免大檔案覆蓋測試返回位址。
- 12 個錯版均拒絕；包括 read count、segment stride、exact-multiple EOF、cue gate／索引／mute、allocation 偏移與 stop／return。

cue 測試使用 `0x6100` 段。初版 `0x6000` 讓未清除的 BH 在三次左移後溢位消失，
第 6 個錯版因而漏過。修正測試段後重跑全部矩陣與錯版；沒有改動原始 C 的清除邏輯。

這份收據驗局部原始控制流；C 機器碼匹配、正常玩家流程與聲音硬體仍依未解範圍判讀。

## scroll helper 證據審查

原 `nullsub_5` 在 `0x105ED`，兩個直接 caller 為 `sub_1054D`／`sub_1059B`。
固定 IDA 的直接 xref 與逐段 operand 候選查詢沒有發現 `0x105ED..0x105F0` 的寫入／取址候選。
相鄰 `0x105EE` 段讀四個不同 display segment，不能拿它當作 DOS/V scroll helper 的安裝證據。
這份是「目前未定位」的負證據，xref 不涵蓋全部間接寫入，沒有宣稱原版不存在執行期替換。
查詢收據：[c-scroll-writer-verification.json](c-scroll-writer-verification.json)。
查詢：[ida_scroll_writer_probe.py](../../tools/ida_scroll_writer_probe.py)，原輸入與獨立 DB 保持唯讀／非破壞性。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原 INT50 媒體錯誤／完整退出流程 | 保留原 branch，不以成功讀檔代證全部錯誤介面 |
| 聲音 TSR 與實際音訊 | 本輪驗 cue／API／資料，不取代音訊解碼與人耳驗收 |
| scroll helper 間接改寫 | 直接查詢的負證據不構成完整不存在證明 |
| 原作者 C 工具鏈／機器碼 | 仍待確認 |
