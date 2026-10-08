# 92：C 播種函式：固定 RTC 的原版與 Go 對照

**狀態：局部播種對照通過。合法時分秒、明示 RTC AL 與堆疊邊界有原版/C/Go 收據。**

- 日期：2026-10-08
- EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 位址：IDA 9.4 線性 `0x1EC82`–`0x1ECE0`；檔案 `0xEE82`–`0xEEE0`
- 指令身分：94 bytes SHA-256 `420ac03341ea350b7ebb66a5e5d9210830a5dd18ccd4a5bc9dff123517271187`
- 固定 inventory：`workplace/matching-decompilation/assembly/dosv/ida-probe.json`，DB 與工具身分見 [`90`](90-assembly-reconstruction.md)
- 資料證據：[`10`](10-rng.md)；契約：[`spec/202`](../spec/202-c-rng-seed.md)

## 1. C 函式與 RTC 邊界

[`seed.c`](../../tools/c_recovery/seed.c) 先建立 `T[i]=i`，再依固定兩個 byte 索引交換 256 次。
最後寫入 `c=(RTC_AL + second + minute + 4×hour) mod 256`、`s=c XOR second`。
輸入仍是 BCD 原始 bytes，正常 Go 比較使用 RTC_AL=0。

原始函式在 INT 1Ah 之後直接讀 AL、CH、CL、DH，沒有檢查 CF 或把 AL 清零。
所以 C 介面必須保留 RTC_AL。它不能從「常見 BIOS 可能保留 AL」推定為恆等於零。

本機 dosgolem 的 AH=02h 只回報沒有 RTC，沒有成功時間來源。本輪在 IVT 裝入
明示的 RAM fixture：MOV CX、MOV DX、MOV AL、IRET。向量不再指向 dosgolem stub，
CPU 執行真實 INT／IRET，原始播種函式的 94 bytes 不變。
這是控制函式輸入的方法，不是原版硬體 RTC 的重建。

## 2. 完整矩陣與結果

| 條件 | 每個最佳化版本 |
|---|---:|
| 00:00:00 至 23:59:59，RTC_AL=0 | 86,400 組 |
| RTC_AL=`1`、`0x7F`、`0x80`、`0xFF`，23:59 的秒 `0`、`27`、`59` | 12 組 |
| 入口 SP=`0`、`2`、`0xFFFC`、`0xFFFE`，23:59 的秒 `0`、`59` | 8 組 |
| 原版/C 完整狀態與介面 | **86,420 組相同** |
| 原版/Go `rng.New`，限 RTC_AL=0 | **86,408 組相同** |
| 完整 1 MB 記憶體核對 | **85 次相同** |

O0/O2 的 JSON 收據相同，SHA-256：
`e89105ba729022f1d204533617731b45cc419680a1a5c17adca88392b4d3737a`。
原版/C 的 258-byte 狀態摘要相同：
`e3230de991da1ae62b69cc9a2d5b2c7633887872f8032ebfe6121667e3cf9863`。

每個原版案例固定執行 3,369 條指令，包含四條 RTC fixture 指令。
除資料狀態外，也比較完整 14 個 16-bit 暫存器／段／FLAGS 與持續堆疊 bytes。
寫入監看限制在 RNG 區域與堆疊，沒有未宣告的寫入。

## 3. 介面副作用與負對照

原版保存 DS、ES、AX、BX、DX，沒有保存 CX；CH 在計算時左移兩次，因此 CX 會改變。
INT 的返回 IP 留在堆疊較低的位置；FLAGS、CS 的槽稍後被播種暫存 AX、BX 覆蓋。
適配層保留這些持續 bytes，不只比較最終 c、s、表。

最後旗標來自 XOR c,second。CF／OF 清零，SF／ZF／PF 依結果；AF 在硬體上未定義，
本輪按實際 dosgolem 的 AF=0 比較，保留模型限定。

| C 突變 | 首次拒絕 |
|---|---:|
| 第一索引步長改 `0x4E` | 第 1 組 |
| 小時項不乘 4 | 第 3 組，RTC_AL 邊界 |
| s 改成 c XOR minute | 第 2 組 |

## 4. 來源身分與重跑

```sh
tools/c_recovery_seed.sh
```

對照器在 [`c_recovery_seed.go`](../../tools/c_recovery_seed.go)，驗證器在
[`c_recovery_seed_verify.py`](../../tools/c_recovery_seed_verify.py)。
全部 workload 在非 root、無網路、限資源 Docker 內執行；原版、dosgolem 與專案來源唯讀。
Go workspace 只在容器內建立。dosgolem 被 import 的來源在執行前後雜湊相同。

| 來源 | SHA-256 |
|---|---|
| C | `5310de9dcc0e40a374d168f73ffa8de2222de654378c7bfc6c50c4255e851cbb` |
| C header | `4f9e0c6e3026967450b6a29758bedd2e293010187a95e2fd98b0a1c438a0ffc9` |
| Go 對照器 | `0ef6fbffe2e23e967510b287fde68b47bcb5961f30b0eecb744bfdcc81a4d1c1` |

本機產物根為 `workplace/matching-decompilation/c-seed/`，results 下保存 O0/O2 與三份
突變收據，verification.json 保存輸入、來源、工具與結果的身分核對。

## 5. 未解範圍

| 項目 | 邊界 |
|---|---|
| 真實 RTC 與保留值 | 本輪只有受控回覆，沒有實機 RTC 證據 |
| 非法 BCD、其他 CS、堆疊重疊 | 未驗證，不能外推 |
| C 機器碼匹配與玩家流程 | 未驗證，不宣稱整個 matching decompilation 完成 |
| XOR 的 AF | 本輪按工具模型比對，不聲稱硬體定義 |

C 完成狀態查 [`c-recovery-status.json`](c-recovery-status.json)。原始組語基準與先前收據保留。
現行玩家工作仍查 [GitHub Issues](https://github.com/wicanr2/wolong_cht/issues)。
