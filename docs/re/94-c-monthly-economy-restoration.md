# 94：月結呼叫鏈與五個經濟 C 函式

**狀態：六個 C 函式局部對照通過。月結的資金、赤字與 RNG 已由 C 函式實際相接。**

- 日期：2026-10-08
- 輸入：松崗 DOS/V KI.EXE，67,099 bytes
- EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具：IDA Pro 9.4；位址空間為 database linear，檔案偏移另列
- 本輪 DB：`workplace/matching-decompilation/c-economy/ida/input.exe.i64`
- DB SHA-256：`7cacd5bfb39f92b26767c5f94b6dfb8fe3beef2c237d2c7c09a6fb137b66b149`
- 探針：[`ida_economy_probe.py`](../../tools/ida_economy_probe.py)，原始名稱、chunks、bytes、operands 與 xrefs 保留
- 契約：[`spec/205`](../spec/205-c-monthly-economy.md)、[`spec/206`](../spec/206-deficit-high-word-rounding.md)

## 1. 原始定位

| 原始函式 | IDA 線性區間，右界不含 | 檔案偏移 | bytes SHA-256 |
|---|---|---|---|
| `sub_15358` | `0x15358`–`0x153C6` | `0x5558`–`0x55C6` | `442201b9d1c4318beaa15c06fc08e647c6c7fc6158164d2f660da18a4b62cba7` |
| `sub_15609` | `0x15609`–`0x1562B` | `0x5809`–`0x582B` | `900be68bb0f82c658d1ca1a078d99944700aa5a2bd3f252763d3d2f007e7cad7` |
| `sub_1563B` | `0x1563B`–`0x15663` | `0x583B`–`0x5863` | `4a3ac32caa3ef009fba1c269a4a58cf55b511d8a1e23692fbbf41963b9a3760d` |
| `sub_154FC` | `0x154FC`–`0x15532` | `0x56FC`–`0x5732` | `6736ce700e02e63c8901a417028b4b34085b6ad735d77093aee428f5a16f9b3f` |
| `sub_155EC` | `0x155EC`–`0x155F9` | `0x57EC`–`0x57F9` | `5024845424b74eb69e6d7e8ab038ce0213d91718ab72fe001ef450dc132e48d7` |
| `sub_15828` | `0x15828`–`0x1585F` | `0x5A28`–`0x5A5F` | `442c1c7b26a8a7c7412220e2ab49a7204e5638155309085ccb1fa033950a36a5` |

原始指令共 306 bytes。probe 的每個 chunk 均與 EXE 對應原始 bytes 核對，
函式數 739 只驗 probe 已跑到，不用來宣稱 C 完成比例。

## 2. 已證實的控制流與 C 接線

[`economy.c`](../../tools/c_recovery/economy.c) 提供六個原始函式與 16-bit ABI 適配。
已證實的語意包括 22 個勢力的存在旗標、扣支出／結算／入帳／清支出／赤字順序，
最後九個 callee、四個設定 word 搬移、重畫呼叫及 DS=CS、AX 保存。

月結對資金加減與赤字函式的呼叫由新 C 函式直接完成；赤字再呼叫既有 C RNG。
原版側保留相同 callee 的真實指令，沒有用輸出常數代替這段呼叫鏈。
據點結算的收入由明示 MOV AX／DL fixture 提供，其餘十個 callee 使用 RET。
這些 fixture 限定未還原 callee 的輸入，不宣稱其完整功能已重建。

| 原始定位 | 已證實語意與適用邊界 |
|---|---|
| `sub_15609`／`sub_1563B` | 24-bit 加減後有號判斷；收入只鉗上界、扣款只鉗下界。正常 Go 比較限合法資金及非負額度 |
| `sub_154FC` | 切比雪夫距離，夾到 255；CS 門檻 80／200／255 對除數 2／3／4。矩陣限合法地圖座標差與各方向 |
| `sub_155EC` | 16-bit carry 或總和 >65500 時回傳 65500；carry 路徑的 FLAGS 與非 carry CMP 路徑不同 |
| `sub_15828` | 先取負資金高 16 位再 NEG、四次 SHL；各兵種以真實 RNG 取數，借位歸零 |

每個案例比較 14 個暫存器／段／FLAGS、勢力資料、設定、RNG 表與堆疊。
callee 入口快照另比 14 個暫存器、當時記錄、設定與 RNG c/s，避免呼叫順序差異被最後狀態掩蓋。
FLAGS 的 AND／XOR／SHL AF 按 dosgolem 模型為 0；不當作實機未定義旗標的證據。
RNG 表用 12:34:56 播種；每案例的 c/s 由案例索引在執行前固定，兩側寫同一份 258 bytes，
Go 使用同一份 `rng.FromRaw`。設定方式記在 verification.json，沒有重擲或改動正式遊戲種子。

## 3. 完整矩陣與結果

| 原版/C 條件 | O0、O2 各自案例 |
|---|---:|
| 收入：六種資金入口 × 全部 AX；另取極值與四種 DL | 393,440 組相同 |
| 扣款：同樣的高低位、carry／borrow 與上界相鄰值 | 393,440 組相同 |
| 預備兵：全部 AX × DX 0／1／35／36／65500／65535 | 393,216 組相同 |
| 距離：X 差 0–383、Y 差 0–255，四個方向 | 393,216 組相同 |
| 赤字：全部高 16 位 × 低 byte 0／1／255，固定 RNG 初值 | 196,608 組相同 |
| 月結：22 個位置 × 全部 status byte；另 16 個多勢力樣本 | 5,648 組相同 |
| 合計原版/C | **1,775,568 組相同** |
| Go 合法資金、募兵、距離與局部赤字規則 | **1,326,175 組相同** |
| 完整 1 MB 記憶體抽查 | **867 次相同** |

O0/O2 的 JSON 收據相同，SHA-256：
`df57663b933bfe0ec205d8672fc97000939a7ea68b6157a4f14fa08b590932fd`。
兩側狀態與入口軌跡摘要：`02c8c295ff7ac150199df82718a937f2728fa18b1ad905fee753cf34ec4e065a`。
Go 沒有完整月結 fixture 的完成聲明，非法資金與高位額度也未外推。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 資金上限少 1 | credit 第 262,145 組 |
| 資金下限多 1 | debit 第 1 組 |
| 募兵忽略 16-bit carry | reserve 第 131,072 組 |
| 距離第一門檻改為 79 | distance 第 321 組 |
| 赤字只 SHL 三次 | deficit 第 98,308 組 |
| 收入入帳與據點結算交換 | monthly 第 129 組，status=0x80 |

## 4. Go 赤字訂正

原版先讀 `[si+21h]`，再取負值；Go 曾先取欠款絕對值再右移。
欠款不是 256 的倍數時，Go 每個兵種少扣 16。原版/C probe 的合法負資金 -654849
已證明差異，固定同一 RNG 的弓兵剩餘量為原版 24547、舊 Go 24563。

Go 改為 `ceil(欠款/256)*16`。六個明確邊界與非零固定 RNG 冷測通過；
完整矩陣中的 15,351 組合法資金原版／Go 預備兵與最終 RNG 相同。
`internal/rules/economy` 與 `internal/state` 的冷測也通過。
機制現況在 [`mechanics/40`](../mechanics/40-economy.md)，形成原因與訂正在 `CONTEXT.md` §6。

## 5. 重跑與來源身分

首次建立本機 IDA 證據：

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_economy_probe.py workplace/matching-decompilation/c-economy/ida KI.EXE
tools/c_recovery_economy.sh
```

原版與外部 dosgolem 來源唯讀，workload 使用非 root、限資源 Docker。
Go 1.26.7、GCC 12.2.0，沿用 `golang:1.26.7-bookworm` 已驗證 image。
dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`；被 import 的來源在前後雜湊相同。

| 來源 | SHA-256 |
|---|---|
| C | `a8dfb96bc788326cc8641ffa78a64739d2544dfbf63ed7149dc334ef85f2ab8e` |
| header | `57958375c3f3cf604e0dffb8a56e2d0b464c727702d1e076b6f303a3806e89ec` |
| callee 入口 fixture | `8313a73e688631a8fe39db8ab12d66c979ed1f7e978ae9907e8910f24f91d0a8` |
| Go 對照器 | `cf91fc474696a32aa82243a2943b8d29f33819014e37a7e97722829f5db8b159` |

完整來源／DB／收據核對在 [`c_recovery_economy_verify.py`](../../tools/c_recovery_economy_verify.py)，
本機產物在 `workplace/matching-decompilation/c-economy/`。現行 C 完成台帳查
[`c-recovery-status.json`](c-recovery-status.json)。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| 據點完整結算與尾端 callee | 未驗證完整 AI、內政、外交與事件效果 |
| 月結正常玩家流程、UI 與完整存檔 | 未驗證，局部呼叫鏈不取代垂直鏈驗收 |
| C 機器碼匹配 | 未驗證，原版 C 工具鏈仍未知，整檔匹配仍由組語基準提供 |
