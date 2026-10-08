# 202：sub_1EC82 的 C 播種契約

**狀態：CONFORMED。固定 RTC 回覆下，完整播種介面及合法時間矩陣的原版/C/Go 對照通過。**

- 日期：2026-10-08
- 來源：松崗 DOS/V `KI.EXE` SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 位址：IDA 9.4 線性 `0x1EC82`–`0x1ECE0`；檔案 `0xEE82`–`0xEEE0`
- 證據：[`re/10`](../re/10-rng.md)、[`re/90`](../re/90-assembly-reconstruction.md)
- 範圍：C 資料函式及局部播種，不驗 RTC 硬體時序、遊戲開機呼叫時機或正常玩家路徑

## 1. 審查後的資料契約

| 原始指令 | 行為 | 等級 |
|---|---|---|
| `0x1EC94`–`0x1EC99` | 初始化 `T[i]=i`，256 筆 | 已證實 |
| `0x1EC9D` | 呼叫 INT 1Ah AH=02h | 已證實 |
| `0x1EC9F`–`0x1ECAF` | `c=(RTC_AL + second + minute + (hour<<2)) mod 256` | 已證實 |
| `0x1ECB3`–`0x1ECB5` | i=second，j=second+1 | 已證實，限合法 BCD 秒 |
| `0x1ECB6`–`0x1ECCE` | 256 次交換；i+=0x4F、j+=0x89，低 byte 捲回 | 已證實 |
| `0x1ECD2`–`0x1ECD7` | 寫入 c，再寫 s=c XOR second | 已證實 |

輸入是 BIOS 回覆中的原始 BCD hour/minute/second 與 AL，不能把 AL 默默假設為零。
正常 Go 播種比較使用 AL=0；另以非零 AL 做 C 介面邊界對照。
本輪只使用合法 BCD 時分秒。非法秒 0xFF 會讓原版的 16-bit j 起點超出 256-byte 表，
不納入本輪資料層介面，保持原版邊界未知。

## 2. 驗證條件

dosgolem 尚未實作 RTC 成功回覆，本輪以明示的 RAM 中斷處理器提供固定 CX、DX、AL。
IVT 指向該處理器後，CPU 執行原始 INT／IRET；沒有改寫被驗證函式指令。
IF/TF=0、偶數 SP，堆疊與程式碼／RNG 狀態分離。

對照完整 258-byte 狀態與 Go `rng.New`，並檢查保存暫存器、CH 的低 byte 捲回與未宣告寫入。
適配層保留五個 PUSH、INT 的持續堆疊痕跡、播種暫存字與近返回。
最終旗標來自 XOR；其 AF 在硬體上未定義，本輪按 dosgolem 的 AF=0 比較，不外推實機。
先固定 RTC 回覆再執行，不能挑選碰巧通過的時間。
RTC 回覆、工具版本、程式雜湊、初始狀態與結果隨收據保存。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| RTC 的未指定暫存器 | 本輪固定回覆，不外推真實 BIOS 保留哪些值 |
| 非法 BCD 與堆疊重疊 | 未驗證，不能從合法資料層推定 |
| C 機器碼匹配 | 未驗證，這是 C 語意還原 |
| 玩家流程與 RTC 時序 | 未驗證，不作硬體時序考古 |

## 4. 驗證結果

O0/O2 各 86,420 組原版/C、86,408 組原版/Go 通過；每組原版執行 3,369 條指令，
其中四條為明示 RTC fixture ISR。每版 85 次完整 1 MB 記憶體核對相同。
三個 C 突變皆被拒絕，收據與範圍見 [`re/92`](../re/92-c-rng-seed-restoration.md)。
