# 236：C 場景主入口與退出重畫

**狀態：CONFORMED。O0／O2 各 171 完整場景／裝置案例與 12 個錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[re/116](../re/116-c-scene-resume-restoration.md)

## 原始契約

保留原 far／near／尾跳、三句 TALK index 算法、對話等待、停／播曲順序、MMAP 恢復，
以及 world producer → cell renderer。鏡頭保留借位、dirty、記憶坐標與 `0xFFFF` 哨兵。
小地圖保留 192×128 圖庫、位元對齊、四 plane、latch、RCR carry 與 I/O 順序。
原版平台 API 用獨立成熟服務，C 不執行 guest CPU；不把受控 press／counter 外推實機時序。

## 驗證閘門

O0／O2 比完整 RAM、plane、FLAGS／暫存器／SS frame、呼叫入口、IN／OUT 與 DOS／mouse／sound API。
真實 ICONGRF／MMAP／IVENTGRF／TALK／SINARIO，覆蓋 byte 位移、舊／新相機、兩個尾跳分支、
八種 bit 對齊、far 座標／游標，以及完整場景主入口與三句／還原結果。
錯版須拒絕，來源由固定 IDA 乾淨重生；原始 C 機器碼匹配另驗。

## 實作與收據

`tools/c_recovery/resume.c` 與原 modal／input／resource／world 實作形成完整局部閉包。
[收據](../re/c-resume-verification.json)記錄真實素材、呼叫端 DS／TALK 基址、三句、
far／tail、原始 renderer 與 816 glyph，重跑 `tools/c_recovery_resume.sh`。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 正常玩家長程路徑 | 局部主入口不代證自然事件與通關 |
| 原作者 C compiler／machine code | 尚未確認 |
