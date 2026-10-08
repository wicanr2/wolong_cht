# 95：據點收入、募兵與 C 經濟月結接線

**狀態：五個 C 函式 O0/O2 與 Go 局部對照通過；據點結算已接入 C 經濟月結。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列，區間右界不含
- DB：`workplace/matching-decompilation/c-settlement/ida/input.exe.i64`
- DB SHA-256：`9f8ef4bf92d76f1cce93df6f61a86fbc6be7e41731a320322de724f20a0985d4`
- 探針：[`ida_settlement_probe.py`](../../tools/ida_settlement_probe.py)
- 契約：[`spec/207`](../spec/207-c-city-settlement.md)、[`spec/208`](../spec/208-settlement-width-parity.md)

## 1. 原始定位與已證實語意

| 原始函式 | IDA 線性區間 | 檔案偏移 | bytes／SHA-256 |
|---|---|---|---|
| `sub_153C6` | `0x153C6`–`0x15456` | `0x55C6`–`0x5656` | 144；`8679095617271c468a25d4c272883f9ac31e407d14c016b52bea1dce6954c2ce` |
| `sub_15456` | `0x15456`–`0x1548F` | `0x5656`–`0x568F` | 57；`e344a10708f04207d6fe4131a77d312951aa00f7a9f1f1f47802ce89d533b41c` |
| `sub_1548F` | `0x1548F`–`0x154FC` | `0x568F`–`0x56FC` | 109；`82a7e3aaa9951ddb651bcba4990af63fea1c69224f34674d7d7e72adf9ac1b88` |
| `sub_15538` | `0x15538`–`0x15547` | `0x5738`–`0x5747` | 15；`d906aca6a2926a53d2e1ac9923a21190880a1f3e4a9d6644752667caf7762bb3` |
| `sub_15547` | `0x15547`–`0x155A6` | `0x5747`–`0x57A6` | 95；`b26946c1dfb9d43d8da7811e03c85705e03826014f7d3eadf8c58883fae158a1` |

原始指令共 420 bytes，IDA 的原始名稱、operand、chunk 與 xref 均保留。
以下語意等級為已證實，範圍限本輪合法輸入矩陣與明示的尾端 fixture。

`sub_153C6` 在 SS 堆疊上建 10-byte 累計，掃 DS 的 192 個據點。
只選 owner 相同的據點，先取距離除數，再累計收入與募兵。玩家分支套稅率與募兵上限，
AI 分支收入除 2 並跑原始募兵節制。最後加預備兵，AX 回收入低 word、DL 回收入高 byte、
DH 回據點數，其他保存暫存器及持續堆疊與原版相同。

[`settlement.c`](../../tools/c_recovery/settlement.c) 接到既有月結，原先 MOV 收入 fixture 已由真實 C 據點結算取代。
原版側保留五個函式及既有距離、資金、預備兵、赤字與 RNG 的真實 bytes。
九個尾端世界更新常式及最後重畫仍為 RET fixture，沒有把完整世界月結宣稱為完成。

## 2. 原始整數寬度

| 原始定位 | 已證實結果 |
|---|---|
| `0x1553F`／`0x15542` | 收入 low-word ADD 的 carry 進高 byte，24-bit 累計 |
| `0x1559A`／`0x1559D`／`0x155A0` | 三兵種累計各為 word；每個據點加完就以 16-bit 繞回 |
| `0x154B8`／`0x154BC` | 玩家稅率最後只做 low-word ADD，沒有把 carry 再加到高 byte |
| `sub_15456` | 0x20 錯步進、127 次掃描、16-bit 軍團和，見既有 [`spec/181`](../spec/181-ai-recruit-gate-and-monthly-globals.md) |

原版稅率反例 gross=66303、tax=99，收入為 103；普通乘除為 65639。
Go 已改為原始兩段乘除與 low-word ADD。

192 個北方據點各 production=65535、除數 2 時，分次移位分配再累計，
原版預備兵為 51584、5952、7808；舊 Go 無寬度累計後鉗制為 65500、5952、65500。
Go 已在每次募兵累計遮為 16-bit，再套兩種上限。
概略比例不作精確測試期望；本輪密集案例曾忽略分次捨位，已依直接原版回傳值訂正。

## 3. 矩陣與已取得收據

| 條件 | O0/O2 各原版/C | Go |
|---|---:|---:|
| 收入：全部 production × 除數 2／3／4 × 累計 0／65535／0xFFFFFF | 589,824 | 196,608，限累計 0 |
| 募兵：全部 production × 三個除數 × Y 0／79／80／149／150／255 | 1,179,648 | 1,179,648 |
| 玩家：16 個收入邊界 × 稅率 0–100 × 兩組募兵設定 | 3,232 | 3,030，排除非法世界 gross |
| AI：8 個 gross × 6 個支出 × 128 個軍團位置樣本 | 6,144 | 6,144 |
| 據點：22 個勢力 × 五種據點數 × 四組玩家／AI／稅率輸入 | 440 | 440 |
| 四個原版劇本 × 三個玩家選擇，C 經濟月結接線 | 12 | 12 |
| 合計 | **1,779,300 相同** | **1,385,882 相同** |

每案例比較 14 個暫存器／段／FLAGS、完整世界資料、SS locals、CS 收支與設定、RNG 表及堆疊。
callee 入口快照另核對同時的記錄、locals、收支及 RNG c/s。
完整 1 MB 記憶體抽查 435 次相同。IF/TF=0，SS 與資料分離；未定義 FLAGS 按 dosgolem 模型限定。

O0/O2 JSON 相同，SHA-256：`6b56885bbdaded199a1d352586ce3921de61db0eaad362cb4914234af24fbb02`。
原版/C 狀態與入口軌跡摘要：`47a6567d46934f8cf96ab1669d0566666eb6074ec36689c5db366cc7fa8eca59`。
兩項 Go 修正的經濟、狀態層冷測與 Go vet 通過。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 收入忽略 low-word carry | income 第 65,539 組 |
| 北方募兵改概略比例 | recruit 第 65 組 |
| 玩家收入改普通乘除 | player-tax 第 2,421 組 |
| AI 步進改成 0x40 | ai-gate 第 1 組，含原始 BX／FLAGS 介面差異 |
| 據點 owner 選擇反轉 | city-settlement 第 1 組 |

## 4. 來源與重跑

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_settlement_probe.py workplace/matching-decompilation/c-settlement/ida KI.EXE
tools/c_recovery_settlement.sh
```

Go 1.26.7、GCC 12.2.0，原版與 dosgolem 唯讀；非 root、有界 CPU／記憶體／PID 的 Docker。
dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import 來源的前後雜湊核對。
RNG 表以 12:34:56 產生，各案例 c/s 在執行前由索引固定，原版、C、Go 使用相同 258-byte 狀態。
收據明示設定方法，沒有重擲，也沒有鎖定正式遊戲種子。

| 來源 | SHA-256 |
|---|---|
| C | `c6573bbfef750ca81c847efe88b0b3d2b51ac1623a321e8925983bbd59a750a3` |
| header | `5714b35cff296141a05bd25428728dc78cc4f22773c6d4d8b4b229f741f6c0f7` |
| 入口 fixture | `af8305a0f0610de1301bf81349fe59f97f22920e7f6d6c91660efa2de3f7f345` |
| Go 對照器 | `63837d611eaaf78039bcbfb381da7717b328f37dfd9de26d3bf7fa8b65480669` |

本機產物在 `workplace/matching-decompilation/c-settlement/`，驗證器為
[`c_recovery_settlement_verify.py`](../../tools/c_recovery_settlement_verify.py)。
先前經濟收據保留在 `workplace/matching-decompilation/c-economy/pre-settlement-proof/`，
現行六函式收據已按新 Go 來源重生，案例及結果保持相同。

## 5. 未解範圍

| 項目 | 邊界 |
|---|---|
| 尾端世界更新與重畫 | 未驗證完整內政、外交、天災及事件效果 |
| 正常玩家長程月結及完整 UI／存檔垂直鏈 | 未驗證，四劇本局部資料接線不替代正常玩家流程 |
| 非法資料、DIV fault 與硬體時序 | 未驗證，不從局部合法矩陣外推 |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知，整檔匹配仍由組語基準提供 |
