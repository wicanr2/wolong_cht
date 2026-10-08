# 111：C 訊息等待、滑鼠分派與游標

**狀態：19 個 C 函式／兩 raw 入口，O0/O2 各 210 全裝置／ABI 相同，三個背景還原與十二錯版通過。**

- 日期：2026-10-09
- 輸入：DOS/V KI.EXE，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear；檔案偏移為線性位址減 `0x10000` 加 512
- DB：`workplace/matching-decompilation/c-input/ida/input.exe.i64`
- DB SHA-256：`7e307f175321bfde14d8f7cfc77fa1351be676f1ee3ca8b978cdee1a66ed530c`
- 入口：[`ida_input_probe.py`](../../tools/ida_input_probe.py)
- 前置：[`re/110`](110-c-talk-rendering-restoration.md)、[`spec/45`](../spec/45-advise-scene-layout.md)
- 契約：[`spec/231`](../spec/231-c-input-mouse.md)

`sub_18810` 接完整繪圖／等待／擦除，`sub_121E7` 輪詢右鍵優先的 press count。
`sub_20000` 為另一個 segment 的 16 項 far 分派，包含實際游標保存、恢復與兩層繪圖。
原始 callback `0x20101` 與等待 `0x1222B` 保留 raw code 邊界，未知作者 C 邊界不改名。

原 `0x12239` 的 `EB 0F` 跳過左鍵詢問，仍問右鍵並等十個 counter 間隔。
patch 成 `90 90` 才同時問左／右鍵。兩者都可能由右鍵或 counter 結束，
較早「不 patch 就不等待」須以本輪完整原始分支與有界實驗勘誤。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原版 timer producer 與 wall-clock | 本輪注入明示 counter 輸入，驗等待消費端，不換算牆上秒數 |
| 自然玩家與自動 callback 派送 | direct callback／固定 press 序列不代證正常玩家操作 |
| 原作者 C 工具鏈與機器碼 | 仍需另驗 |

## 原始分段與可回查定位

main 的 `sub_18810`／`sub_121E7` 共 114 bytes；mouse segment 的 `sub_20000`、
`sub_2002E`、`nullsub_4`、`sub_20070`、`sub_2009A`、`sub_200BD`、`sub_200C0`、
`sub_20137`、`sub_2016B`、`sub_2019D`、`sub_201C6`、`sub_201E4`、`sub_2020C`、
`sub_20249`、`sub_202A0`、`sub_202BD`、`sub_202FE` 共 698 bytes。
原等待 `0x1222B..0x12286` 91 bytes，callback `0x20101..0x20137` 54 bytes。
全部 957 bytes／418 原始指令為 native C，沒有執行期 opcode 解碼器。
`sub_2002E` 尾端落到 `nullsub_4`，兩個原始入口仍分別記錄，C 以共用尾端保留 fall-through。

near call 的 offset 依原 segment 基準 `0x10000`／`0x20000` 解讀；far call 讀 runtime MZ
重定位 operand，返回同一份 SS frame。Mouse entry 的 CLI／STI 和 callback 的 PUSHF／POPF
照原 side effect 保留，不能只以 return value 比對。

## 同狀態矩陣與來源

兩側各 DOS／Mouse／Font／RAM／VGA；C 呼叫成熟 INT 33h／字庫 API，不跑 guest CPU。
每組比全 1 MB RAM、四平面、GC／seq／latch、IN／OUT、十四 register／FLAGS、
near／far／平台 API／counter 比較 checkpoint 與 SS 32-byte 快照；內容區逐像素比較。

| 群組 | O0／O2 各組數 |
|---|---:|
| 16 mouse service／三個 cursor 狀態 | 48 |
| 17 mouse／cursor／callback 入口 | 51 |
| 邊界 | 45 |
| 八種對齊 | 8 |
| patch／左／右／timeout／延遲 | 18 |
| 右鍵優先 press 與延遲 | 14 |
| 真實訊息繪圖／等待／擦除 | 15 |
| 原 `sub_12216` patch／還原 caller | 2 |
| show／callback move／hide 序列 | 9 |
| 合計 | 210 |

三個三步序列都完全恢復先前四平面背景，包含右／下邊界。
兩種最佳化完整收據逐 byte 相同，所有 19 個 named 與兩個 raw 入口實際進入。
原 `EB 0F` 模式不問左鍵，仍需 counter 或右鍵；NOP 加入左鍵。
18 組等待合計 140 個明示 counter checkpoint、220 次 mouse 查詢；
訊息組 110 counter checkpoint、130 次 mouse 查詢。這些是可重播控制輸入，不是原 timer producer
或自然玩家派送的驗收，也不把 `0xFF × 10` 換成未查證的秒數。

| 刻意改錯 | 群組 | 首次拒絕 |
|---|---|---:|
| 先問左鍵 | poll | 1 |
| 忽略原跳過左鍵 patch | wait | 1 |
| SI 改 9 | wait | 1 |
| 不清 pending 左鍵 | wait | 7 |
| 讀錯 cursor mask | cursor | 7 |
| 反轉下界比較 | cursor | 7 |
| 右界少一 pixel | edges | 31 |
| 保存範圍少一 byte | edges | 3 |
| 刪 latch read | cursor | 7 |
| 刪 STI | services | 1 |
| 刪 POPF | cursor | 49 |
| 不設訊息關閉標記 | message | 1 |

[來源](../../tools/c_recovery/input.c)、[header](../../tools/c_recovery/input.h)、
[生成 C](../../tools/c_recovery/input_generated.inc)、[產生器](../../tools/c_recovery_input_generate.py)、
[driver](../../tools/c_recovery_input.go) 與 [平台輸入](../../tools/c_recovery_input_platform.go)
進版控。真正原版 mouse masks／背景 buffer 仍由本機自備 KI.EXE 載入，不公開原資料。
[重跑入口](../../tools/c_recovery_input.sh)／[驗證器](../../tools/c_recovery_input_verify.py)
確認 source digest、原始 bytes、MZ／far、來源 flags、錯版、backlink 與乾淨 C 重生。
[公開摘要](c-input-verification.json) 保留完整來源與收據 hash，Go1.26.7／GCC12.2／IDA9.4。

等待區 21 條／71 bytes 補回公開組語，其中三個 far operand 分開記錄 file／IDA bytes，
GNU 真正組譯相同，不把 instruction bytes 當 fallback。總 24,735 指令／56,268 bytes，
完整 67,099-byte EXE SHA 相同；原 header／非指令資料仍只在本機匯入。
舊 glyph 補充檢查保留原 338 指令子集，新增資料不覆寫歷史收據。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_input_probe.py workplace/matching-decompilation/c-input/ida KI.EXE
tools/c_recovery_input.sh
```
