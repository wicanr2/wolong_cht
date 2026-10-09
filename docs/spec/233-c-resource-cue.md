# 233：C 完整檔案讀取與原始 cue

**狀態：CONFORMED。O0／O2 各 93 組局部同狀態驗證與 12 個錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[`re/113`](../re/113-c-resource-cue-restoration.md)

## 原始契約

保留原 DOS file handle、`0xF000` 單讀、segment `+0xF00`、短讀／EOF、stack／CF 與
`0xFFFF` 返回。cue 依原 BGM header 的 index／offset／length載入，不拿音訊結果回填 C。
同 cue gate、先停、載入 header／曲段、播放與 mute API順序保留。allocation 所衍生的
TALK／各圖庫／字庫／世界區 segment 保留原加法、uint16 與兩筆原大小。

使用獨立成熟 DOS／sound API，C 不跑 guest CPU。input bytes與平台來源唯讀，
平台不把原 TSR 的逐週期聲音時序當成已證實 parity。成功 allocation／錯誤返回分開聲明。

## 驗證閘門

O0/O2 比完整 RAM／VGA／register／FLAGS／near／DOS API前後與SS frame。
完整真實檔、`0xF000` 邊界／exact-multiple／多段來源與尾端守衛、MMAP.MDL／MCH caller、
cue／重複 cue／mute／stop／allocation 分段必須覆蓋，錯版須拒絕，C從固定IDA乾淨重生。

## 實作與收據

`tools/c_recovery/resource.c` 與生成內容保留 6 支原始函式的 165 條指令語意。
重跑入口為 `tools/c_recovery_resource.sh`；[收據](../re/c-resource-verification.json)
記錄固定輸入、IDA DB、工具鏈、完整來源雜湊、21 次 buffer 守衛與錯版。
台帳為 [c-recovery-status.json](../re/c-recovery-status.json)。

## 未解範圍

| 項目 | 限制 |
|---|---|
| INT50／完整退出情境 | 原 branch 保留，平台成功路徑不代證錯誤 GUI |
| 原 TSR／硬體聲音 | 僅比 API／資料狀態；原 TSR 的輸出音色與硬體時序尚未驗證 |
| 原作者 C compiler／machine code | 尚未確認 |
