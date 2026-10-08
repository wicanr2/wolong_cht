# 207：據點收入、募兵與 C 月結接線

**狀態：CONFORMED。五函式 O0/O2 各 1,779,300 組原版/C 相同，Go 1,385,882 組相同。**

- 日期：2026-10-08
- 輸入：松崗 DOS/V KI.EXE，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具位址空間：IDA Pro 9.4 database linear；檔案偏移另列
- 既有證據：[`re/07`](../re/07-monthly-settlement.md)、[`re/94`](../re/94-c-monthly-economy-restoration.md)
- 本輪 DB：`workplace/matching-decompilation/c-settlement/ida/input.exe.i64`
- DB SHA-256：`9f8ef4bf92d76f1cce93df6f61a86fbc6be7e41731a320322de724f20a0985d4`

## 1. 原始定位與契約

| 原始函式 | IDA 起點 | bytes | 行為 |
|---|---|---:|---|
| `sub_153C6` | `0x153C6` | 144 | SS 堆疊上建立 10-byte 累計，掃 192 個據點；依 owner 選擇，調距離／收入／募兵；分玩家與 AI，再加預備兵、回 AX/DX 收入與據點數 |
| `sub_15456` | `0x15456` | 57 | 累計收入除 2；以 0x20 步進掃 127 個軍團位置，16-bit 和後比較是否阻止募兵，CF 回傳 |
| `sub_1548F` | `0x1548F` | 109 | 原始兩段乘除稅率、word 加法與三兵種上限；更新玩家收入／支出顯示值 |
| `sub_15538` | `0x15538` | 15 | production÷BX 累加到 24-bit 收入；carry 進高 byte |
| `sub_15547` | `0x15547` | 95 | production÷BX÷32 按 Y 分區，以原版移位和餘數分配；三兵種各用 16-bit 累計 |

保留 DS 記錄與 SS:BP locals 的原始位址空間。函式回傳保留全部暫存器、FLAGS、堆疊 bytes。
稅率限 0–100、座標與首都索引合法；沒有為除零／DIV overflow 或非法記錄猜測行為。

## 2. 接線與驗證

原版側不修改這五個函式與其已還原經濟 callee 的 bytes。
C 據點結算將接入既有月結，替換前輪 MOV 收入 fixture；九個世界更新尾端 callee 及重畫仍為 RET fixture。

驗每個 callee 入口的完整暫存器、記錄、SS locals 與 CS 顯示欄位，另驗回傳與持續堆疊。
O0/O2 需通過原版矩陣與完整記憶體抽樣；刻意錯的收入 carry、募兵餘數、稅率、AI 步進及據點分支需拒絕。
Go 比較限同一合法 record 與等價輸入，不以普通整數公式替代尚未驗證的原始捨位或溢位。

原版/C/Go 矩陣、四劇本接線、435 次完整記憶體抽查及五個負對照已通過。
來源、工具身分與收據見 [`re/95`](../re/95-c-city-settlement-restoration.md)。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 月結尾端世界更新、UI 與玩家存檔流程 | 未驗證，保留明示 fixture 的範圍 |
| 非法索引、稅率及除法 fault | 未驗證，不擴張正式規則契約 |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知 |
