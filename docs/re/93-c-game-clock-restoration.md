# 93：C 時鐘還原與年份邊界對照

**狀態：局部語意對照通過。第三個 C 函式可回查原版控制流，並修正 Go 年份邊界。**

- 日期：2026-10-08
- 輸入：松崗 DOS/V KI.EXE，67,099 bytes
- EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA Pro 9.4 線性位址：`sub_11D8E`，`0x11D8E`–`0x11E17`，右界不含
- 檔案偏移：`0x1F8E`–`0x2017`，137 bytes
- 指令 SHA-256：`4e2b6ecf618756da0f51e5dfb29c34b502fc70ffbf2d7050c5f5f51cb1a311dd`
- 固定 DB SHA-256：`e32713237d223145834531bffc8979146ed4bacaa98f73ec4bbf692b747c2c4f`
- 固定 IDA inventory、工具鏈與組語基準：[`90`](90-assembly-reconstruction.md)
- 契約：[`spec/203`](../spec/203-c-game-clock.md)、[`spec/204`](../spec/204-year-999-1000.md)

## 1. 還原範圍

[`clock.c`](../../tools/c_recovery/clock.c) 還原子刻、時、日、月、年的進位、四個近呼叫與速度等待。
原始名稱與位址保留，C 適配層明示暫存器、FLAGS、DS 相對記憶體與堆疊返回值。
它以結構化 C 表達控制流，沒有用原始指令 bytes 代替 C 實作。

| 已證實的局部語意 | 原始定位與證據 |
|---|---|
| 子刻小於 8 時只增子刻 | `0x11D8E` 起；固定 IDA 指令與合法日期矩陣 |
| 月結入口日為 0，月已換，小時與子刻未重設 | 呼叫 `sub_15358`；callee 入口快照 |
| 月結後日加 1，其後呼叫季節、每時世界更新、日期重畫 | `sub_19377`、`sub_13E11`、`sub_11E17`；順序及入口快照 |
| 年份 999 換年到 1000，入口年 >=1000 換年回到 999 | IDA `0x11DAA`–`0x11DBC`；全部 16-bit 年份矩陣 |
| 速度在 callee 返回後讀取，計數 >=速度才離開等待 | IDA `0x11DF2`–`0x11E17`；MOV 與延後計數 fixture |

推論等級「已證實」限上述原始函式與明示輸入。原版 callee 的完整語意仍未由本輪驗證。

## 2. 原版執行與輸入控制

dosgolem 載入原始 EXE，以 `IDAIn` 保留原始 CS 的段內入口。受測 137 bytes 與原版檔案相同。
四個 callee 在 RAM 換為 RET，另以 MOV fixture 測月結改日、季節改速度、世界更新改 AX、
重畫改 BX 的回傳影響。C callback 使用相同輸入，入口快照比較全部 14 個暫存器／段／FLAGS 與日期。

等待 fixture 在原始 CMP 前提供固定 ready／count 進展，兩側比較輪詢次數、FLAGS 及欄位清零。
IF/TF=0、偶數且分離的堆疊，沒有外部計時中斷。AND／XOR 的 AF 按 dosgolem 模型為 0，
不將硬體未定義旗標解讀為實機契約。

## 3. 驗證矩陣

| 條件 | O0、O2 各自結果 |
|---|---:|
| 年 196、999、1000；每個合法月日、時 1–23、子刻 0–8；速度 0 | 226,665 組相同 |
| 12/31 的換年；入口年 0–65535 | 65,536 組相同 |
| 速度 1–8；月份 1、2、12；即時／延後計數；RET／MOV callee | 96 組相同 |
| 原版/C 完整狀態、callee 快照與輪詢次數 | **292,297 組相同** |
| 原版/Go 日期及時／月事件，排除 48 個 MOV callee 案例 | **292,249 組相同** |
| 完整 1 MB 記憶體抽查 | **286 次相同** |

每個案例核對最終暫存器、日期、等待欄位與堆疊，完整記憶體定期抽查。
O0/O2 JSON 收據相同，SHA-256：
`a066044f86fe6a27158353f34a43cc7b0b8cb142a940b55175db8ca618d52c4a`。
兩側日期與 callee 入口軌跡摘要相同：
`32e97b0dae9153a9473c59902a7f240f27f3bef0d8ca5a4c676a65c2a3b7beca`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 日進位門檻從 23 改 24 | 第 207 組 |
| 換年比較從 1000 改 999 | 第 151,110 組 |
| 季節與世界更新呼叫順序交換 | 第 9 組，即使最終日期相同 |
| 等待從 >=速度改成 >速度 | 第 226,668 組，輪詢次數與 FLAGS 不同 |

## 4. Go 修正與機制文件

[`clock.go`](../../internal/rules/clock/clock.go) 的 `MaxYear` 現為可觀察最高年 1000。
換年時入口年 >=1000 才先設 998，再增為 999；999 可增為 1000。
冷測涵蓋五種年份入口及連續 999→1000→999，原版矩陣另驗全部 16-bit 入口。
日期機制現況見 [`mechanics/15`](../mechanics/15-realtime.md)，訂正紀錄見 `CONTEXT.md` §6。

## 5. 來源身分與重跑

```sh
tools/c_recovery_clock.sh
```

入口使用非 root、無網路、限資源 Docker，原版、專案與 dosgolem 唯讀。
對照器為 [`c_recovery_clock.go`](../../tools/c_recovery_clock.go)，收據驗證器為
[`c_recovery_clock_verify.py`](../../tools/c_recovery_clock_verify.py)。
原始 dosgolem 的被 import 來源在執行前後雜湊相同。

| 來源 | 身分 |
|---|---|
| Go / GCC | Go 1.26.7；GCC 12.2.0；C O0/O2 |
| 執行 image | `golang:1.26.7-bookworm`，ID `sha256:e8c859f5632dcfde7b32d2012b4351728f6437930887c2f6a91ea242459e5514` |
| dosgolem revision | `a9714ebdab2ad6b529f81225472680f2b11f2842`，來源樹身分另存收據 |
| C SHA-256 | `4a845847e050d0d7ff1ebd641fb827602a6ab10fc8fecd643584612e73158439` |
| header SHA-256 | `40d3d4b630dabed9b3f7f98d492fa60734b2a0a6d86d30e0a2961738d2df7284` |
| fixture SHA-256 | `0a7812df9b5331f463acc5fe8f03ae11c5f3313843b6d7119eaf29ac8433167b` |
| Go 對照器 SHA-256 | `7e1fb0a22eb91fd6404b31323e7e9b663ab920739b1878b7598509f547957869` |

本機產物根為 `workplace/matching-decompilation/c-clock/`；results 保存 O0/O2、四份負對照、
工具版本及來源雜湊，verification.json 驗證原版、來源與收據身分。
C 完成狀態見 [`c-recovery-status.json`](c-recovery-status.json)。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| 月結、季節、世界更新與繪圖 | 未驗證完整 callee，fixture 只驗呼叫者如何接收其作用 |
| 等待 wall-clock | 未驗證 PIT 與硬體時序，不外推速度的實際秒數 |
| 原版 1000 年 UI 與長期玩家流程 | 未驗證，不從日期 bytes 外推畫面或存檔垂直鏈完成 |
| C 機器碼匹配 | 未驗證，原版 C 工具鏈尚未知；整檔匹配成果仍是組語基準 |
