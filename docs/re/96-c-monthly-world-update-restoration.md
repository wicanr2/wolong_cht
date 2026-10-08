# 96：C 月結世界更新與真實事件寫入

**狀態：十一函式 O0/O2 局部對照通過，七個月結尾端已由真實 C 函式執行。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear，檔案偏移另列
- DB：`workplace/matching-decompilation/c-world/ida/input.exe.i64`
- DB SHA-256：`c4d92484ac070341ffd090bf0a54ecbe9685ae8d8f203f287e5b14a18f851cf6`
- 原始函式／operand／chunk／xref 探針：[`ida_world_update_probe.py`](../../tools/ida_world_update_probe.py)
- 契約：[`spec/209`](../spec/209-c-monthly-world-update.md)、[`spec/210`](../spec/210-growth-word-parity.md)

## 1. 原始定位與已證實語意

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 已證實的局部語意 |
|---|---|---|---:|---|
| `sub_15695` | `0x15695` | `0x5895` | 128 | 192 據點生產力、上昇值、原始整數寬度與逐據點 RNG |
| `sub_155A6` | `0x155A6` | `0x57A6` | 70 | 127 武將存在旗標與五能力 byte 的移位評分 |
| `sub_12FBF` | `0x12FBF` | `0x31BF` | 79 | AX/DX raw 事件，指定／隨機起點、低 byte 空槽、向後搜尋與 CF |
| `sub_157FE` | `0x157FE` | `0x59FE` | 42 | 赤字高位門檻、欄位 +0x28 與事件 13 |
| `sub_15715` | `0x15715` | `0x5915` | 122 | 內政官狀態閘與請款金額、事件 4 |
| `sub_1578F` | `0x1578F` | `0x598F` | 111 | 外交官狀態閘、雙向關係較小值與事件 5 |
| `sub_130CB` | `0x130CB` | `0x32CB` | 8 | 原始關係 byte 查值及 BX 保存 |
| `sub_13119` | `0x13119` | `0x3319` | 31 | SI／DI 轉 0x600 起關係表的原始位址與 FLAGS |
| `sub_122DB` | `0x122DB` | `0x24DB` | 163 | 暴風雨四個座標 word、RNG 次序與 writer 成功閘 |
| `sub_12286` | `0x12286` | `0x2486` | 85 | 火災／暴動互斥路徑、失敗時 RNG 次序、事件 12 |
| `sub_1237E` | `0x1237E` | `0x257E` | 129 | 暴風雨標記衰減、20 格範圍與玩家通知參數 |

原始指令共 1068 bytes。推論等級已證實，限本輪合法索引／稅率／段與明示 fixture 的範圍。
原始位址、名稱與運算元不改寫，每個函式 bytes SHA-256 在 verification.json 保存並核對原版。

## 2. C 接線與副作用

[`world_update.c`](../../tools/c_recovery/world_update.c) 沿用既有經濟／據點／RNG 適配。
月結中的生產力、武將評分、內政官、外交官、暴風雨、火災／暴動、赤字通知已由 C 執行。
這些事件透過真實 C writer 寫到佇列，原版側也保留相同 writer 與 RNG 的真實指令。

原版 writer 以 queue code 的低 byte 判空。BL=FF 時取真實 RNG，`AL & 0x7C`
形成初始 byte 偏移，再加當前游標；其餘 BL 先乘 4。只向後掃到 0x100，沒有繞回。
佇列滿或起點超界仍保留已發生的 RNG 與 FLAGS，C 不將失敗事件當成沒有執行。

四劇本的 12 個接線案例比較完整世界、CS globals、事件佇列與 RNG。
`sub_1585F`、`sub_12BD9`、重畫、通知與音效目前仍是明示 RET fixture。
本輪不宣稱俘虜政治、完整佇列月度初始化或 UI 已完成。

## 3. 矩陣與收據

| 群組 | O0/O2 各原版/C |
|---|---:|
| 生產力：8 個 production × 7 個 growth 存值 × 5 個稅率 × 玩家／AI | 560 相同 |
| 評分：256 個能力 byte × 4 個 status × 128 槽含終止槽 | 131,072 相同 |
| writer：5 個游標 × 6 個起點 × 4 種空槽形狀 × 8 個種子狀態 | 960 相同 |
| 赤字通知：6 種資金 × 5 個門檻 × 16 個狀態 | 480 相同 |
| 內政官：192 據點位置 × 4 種欄位／計時狀態 | 768 相同 |
| 外交官：22×22 勢力對 × 7 個關係 byte | 3,388 相同 |
| 關係查值／取址：22×22×2 入口 | 968 相同 |
| 火災／暴動：256 個 c/s 狀態 × 6 個免疫值 | 1,536 相同 |
| 暴風雨：256 個狀態 × 新／舊區域 × 6 個座標邊界 | 3,072 相同 |
| 暴風雨標記：4 個強度 × 6 個距離 × 玩家／非玩家 | 48 相同 |
| 四劇本 × 三個玩家選擇的 C 月結接線 | 12 相同 |
| 合計 | **142,864 相同** |

Go `GrowCity` 的 **560 組**結果與最終 RNG 相同。其他原版事件 producer 的本輪證據限 C，
不從它推定完整 Go producer 或 GUI parity。
每版 **1117 次**完整 1 MB 記憶體抽查相同。
每案例另驗 14 個暫存器／段／FLAGS、完整世界、globals、佇列、RNG 及堆疊，
每次 callee 入口比較記錄、堆疊、globals 與 RNG c/s。

O0/O2 JSON 相同，SHA-256：`147f002b17e125c90a902a8afcad3baa7a53e41c2222c647520d9108b4c3772e`。
兩側狀態與入口軌跡摘要：`592a3bfda314832f41799f2a6595d12eae4be205c86b5890f51a043f65d4de9e`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 生產力改普通整數加法 | growth 第 429 組 |
| 評分首三能力改右移 3 | score 第 257 組 |
| writer 改用整個 word 判空 | writer 第 25 組 |
| 內政官費用改乘 51 | governor 第 1 組 |
| 外交官費用加 1 | diplomat 第 8 組 |
| 火災 RNG 門檻改 25 | disaster 第 1 組，含 callee 入口 FLAGS |
| 暴風雨範圍改 19 | marker 第 7 組 |

## 4. Go 生產力訂正

原版 IMUL 只取 AX 低 word，依其有號符號分支；正分支 ADD word 先繞回，再與上限比較。
65000／上昇值100／稅率30 的原版為 12114，舊 Go 鉗成 65535。
65535／100／稅率0 的 low-word 負乘積得到 33149，65535／-100／稅率100 得到 11092。
三個邊界冷測、560 個原版向量、經濟與狀態層冷測及 Go vet 均通過。
形成原因與訂正在 `CONTEXT.md` §6，機制現況在 [`mechanics/40`](../mechanics/40-economy.md)。

## 5. 來源與重跑

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_world_update_probe.py workplace/matching-decompilation/c-world/ida KI.EXE
tools/c_recovery_world.sh
```

Go 1.26.7、GCC 12.2.0，原版、專案與 dosgolem 唯讀，非 root、限資源 Docker。
dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import 來源前後雜湊相同。
RNG 表用 12:34:56 產生；每案例 c/s 在執行前固定並同時寫兩側，Go 使用相同 `rng.FromRaw`，沒有重擲。

| 來源 | SHA-256 |
|---|---|
| C | `90ea41a06b727fc12074c1347e831f6fbac7fb0535265bfbfa7eac7795437602` |
| header | `9782b21c1397d277680b0165b27d697bf7191ad099b399dfc0a875af91f5f088` |
| fixture | `0d2866252dc0339d08d45431beac7359019b29fd80717e07219396cf418ab853` |
| Go 對照器 | `81c0bcecba272dfe0b3a4c687f75f0be2e5b1d145981ea2d2519c03c2db7dfdf` |

本機產物在 `workplace/matching-decompilation/c-world/`，驗證器為
[`c_recovery_world_verify.py`](../../tools/c_recovery_world_verify.py)。
先前經濟與據點收據保留在各工作根的 `pre-world-proof/`，現行收據按新 Go 來源重生。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| 剩餘政治／俘虜及月度佇列初始化 | 未驗證，`sub_1585F`、`sub_12BD9` 保留明示 fixture |
| 完整 Go 事件 producer、UI、音效及正常長程玩家流程 | 未驗證，局部 C 接線不取代垂直鏈 |
| 非法指標與資料、硬體時序 | 未驗證，不從合法向量外推 |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知，整檔匹配仍由組語基準提供 |
