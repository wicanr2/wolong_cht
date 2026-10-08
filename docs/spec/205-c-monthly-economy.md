# 205：月結主流程與五個經濟 C 函式

**狀態：CONFORMED。六個函式 O0/O2 各 1,775,568 組原版/C 相同；尚未還原的 callee 使用明示 fixture。**

- 日期：2026-10-08
- 輸入：松崗 DOS/V KI.EXE，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址空間：IDA Pro 9.4 database linear；檔案偏移另列
- 既有證據：[`re/07`](../re/07-monthly-settlement.md)、[`re/90`](../re/90-assembly-reconstruction.md)
- 本輪 IDA 探針：`tools/ida_economy_probe.py`
- 本輪 DB：`workplace/matching-decompilation/c-economy/ida/input.exe.i64`
- DB SHA-256：`7cacd5bfb39f92b26767c5f94b6dfb8fe3beef2c237d2c7c09a6fb137b66b149`

## 1. 原始定位與審查契約

| 原始函式 | IDA 線性起點 | bytes | C 契約 |
|---|---|---:|---|
| `sub_15358` | `0x15358` | 110 | 22 個勢力；存在旗標 >=0x80 才扣款、結算、入帳、清支出、扣赤字兵；九個尾端 callee 後搬四個 CS word，再呼叫 0x15E80，最後 DS=CS、恢復 AX |
| `sub_15609` | `0x15609` | 34 | AX/DL 與勢力 +0x20 的 24-bit 資金相加；模 24-bit 後有號比較，只鉗上限 +655000；恢復 AX/DX |
| `sub_1563B` | `0x1563B` | 40 | 資金減 AX/DL；模 24-bit 後有號比較，只鉗下限 -655000；恢復 AX/DX |
| `sub_154FC` | `0x154FC` | 54 | 16-bit 無號座標差的絕對值，取最大並鉗距離 255；讀原始 CS 門檻／除數表，BX 回傳除數 |
| `sub_155EC` | `0x155EC` | 13 | AX+DX 有 carry 或結果大於 65500 時回傳 65500，保留原版最終 FLAGS |
| `sub_15828` | `0x15828` | 55 | 資金高 16 位有號時先 NEG、四次 SHL，再以三次原始 RNG 扣三個預備兵；借位歸零，保存 AX/BX/CX/DX |

所有區間右界不含；檔案起點為 IDA 起點減 0xFE00。保留原始 DS、CS、SI、DI 運算元。
C 函式由 KiMachine16 適配，RET 與 PUSH／POP 的持續記憶體副作用納入比較。

## 2. 呼叫鏈

月結中的資金加減、赤字扣兵及 RNG 使用已還原 C 函式實際相接，原版側保留其真實指令。
據點結算及尾端尚未還原的 callee 用 RET／明示 MOV fixture，沒有宣稱完整經濟 AI 或世界更新。
比較每個 callee 入口全部暫存器、日期／資金／支出狀態與順序，不能只比最後資金。

## 3. 驗證閘門

- O0/O2 與原版完整暫存器、FLAGS、資料及堆疊相同，定期核對完整 1 MB 記憶體。
- 資金測上下界相鄰值、24-bit 高位與 carry／borrow；Go 比較限合法資金與非負額度。
- 募兵上限枚舉每個 AX 並取數個 DX 邊界，包含 16-bit carry。
- 距離測合法座標差、80／200／255 門檻及各方向；直接讀原版查表資料。
- 赤字函式枚舉全部高 16 位，並取不同低 byte；固定 RNG 初值且同表對照。
- 月結測每個勢力位置與所有 status byte，並驗全部／交錯存在的多勢力流程。
- 刻意錯的 clamp、carry、距離門檻、赤字量及月結順序版本必須被拒絕。

六個突變版本皆被拒絕，Go 的 1,326,175 組局部規則對照相同，來源、DB 與收據核對通過。
完整矩陣與重跑入口在 [`re/94`](../re/94-c-monthly-economy-restoration.md)。

## 4. 未解範圍

| 項目 | 邊界 |
|---|---|
| 據點結算及九個尾端 callee | 未驗證完整原版語意，fixture 只量呼叫者契約 |
| 正常 UI、完整存檔與長期月結 | 未驗證，不從局部控制流宣稱玩家垂直鏈完成 |
| C 機器碼匹配 | 未驗證，原版 C 工具鏈仍未知 |
