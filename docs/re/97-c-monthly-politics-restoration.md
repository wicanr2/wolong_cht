# 97：C 月結政治、俘虜與完整規則接線

**狀態：二十一函式 O0/O2 局部對照通過，完整月結規則不再使用政治／俘虜替身。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-politics/ida/input.exe.i64`
- DB SHA-256：`74479720151e78166907794ded71267553e4d65e0f92251feb87c835029144c3`
- 原始函式／operand／chunk／xref 探針：[`ida_politics_probe.py`](../../tools/ida_politics_probe.py)
- 契約：[`spec/211`](../spec/211-c-monthly-politics.md)

## 1. 原始定位與已證實語意

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 已證實的局部角色 |
|---|---|---|---:|---|
| `sub_1585F` | `0x1585F` | `0x5A5F` | 58 | 127 武將存在、月度 timer 與入仕／俘虜分派 |
| `sub_15899` | `0x15899` | `0x5A99` | 167 | 求職冷卻、勢力選擇、玩家門檻與入仕 |
| `sub_15940` | `0x15940` | `0x5B40` | 80 | 俘虜逃脫／歸降、事件 9、通知與原始欄位寫回 |
| `sub_12AD2` | `0x12AD2` | `0x2CD2` | 34 | 舊／新勢力的武將數 byte 減／增，FF sentinel |
| `sub_15990` | `0x15990` | `0x5B90` | 22 | 原始玩家通知參數及 SS 堆疊 |
| `sub_1301C` | `0x1301C` | `0x321C` | 50 | 256-slot raw writer、低 byte 空槽與向後搜尋 |
| `sub_12BD9` | `0x12BD9` | `0x2DD9` | 121 | queue 搬移／清零、cursor／cadence、scratch 初始化與兩輪勢力更新 |
| `sub_12C52` | `0x12C52` | `0x2E52` | 141 | 據點邊境收集、relation 查表、原始排序與標記 bit |
| `sub_12CDF` | `0x12CDF` | `0x2EDF` | 91 | 四個鄰接 bit、重複排除與中立 word 標記 |
| `sub_12D3A` | `0x12D3A` | `0x2F3A` | 30 | 無目標勢力的 RNG 遷都事件 8 |
| `sub_12FB1` | `0x12FB1` | `0x31B1` | 14 | SI 轉 event code.high，接既有 64-slot writer |
| `sub_12D58` | `0x12D58` | `0x2F58` | 96 | 玩家／AI 政治呼叫順序、CF 分支與目標保留／清除 |
| `sub_12DB8` | `0x12DB8` | `0x2FB8` | 59 | AI 第一鄰居的原始交友度漂移 |
| `sub_12DF3` | `0x12DF3` | `0x2FF3` | 64 | 玩家第一鄰居、反向勢力的交友度下降 |
| `sub_130F0` | `0x130F0` | `0x32F0` | 26 | 低 7-bit 減值歸零、war bit 保存 |
| `sub_1310A` | `0x1310A` | `0x330A` | 15 | 原始 SI／DI 雙向位址、AX 保存、DX 回傳 |
| `sub_13091` | `0x13091` | `0x3291` | 58 | 三兵種右移累計、word 和、2000 上界、資金高位／據點數閘 |
| `sub_12E33` | `0x12E33` | `0x3033` | 86 | 玩家協力要請條件、雙向 relation 與事件 2 |
| `sub_12E89` | `0x12E89` | `0x3089` | 114 | 累減力量、停戰候選、事件 3 與 CF |
| `sub_12EFB` | `0x12EFB` | `0x30FB` | 118 | 資金／交友／力量門檻、事件 1 與 CF |
| `sub_12F71` | `0x12F71` | `0x3171` | 64 | 中立目標分支、word threshold 與原始目標清除 |

原始指令共 1508 bytes；每個函式的原始 bytes SHA-256 在 verification.json 驗 EXE 與 IDA chunk。
語意等級已證實，限合法勢力／武將／鄰接索引、分離段、至少一個存在勢力與明示 UI fixture。

## 2. 完整 C 月結規則

[`politics.c`](../../tools/c_recovery/politics.c) 接到既有 world／settlement／economy／RNG。
`sub_15358` 的九個尾端 callee 均有真實 C 實作，原版側同樣保留全部真實指令。
佇列搬移、邊境表與政治順序、入仕、俘虜、經濟、成長、災害及 writer 已實際串接。
僅通知、音效與重畫使用 RET fixture；這個完成範圍不包括 GUI 或正常玩家長程。

月度佇列先把 offset 0x100 起的 0x300 bytes 前移，再清後 0x100 bytes。
cursor 設 0、cadence 設 7，ES scratch 的 0x210 words 清為 FFFF，最後依序建立邊境及跑政治。
原版 CLD 的 DF 清除與 REP 後 SI／DI／CX 介面也納入比較，不能只驗 queue 內容。

原版 256-slot writer 沒有隨機起點，以 BL×4+cursor 開始讀；空槽看 code.low。
新／舊任官計數是 byte，有 FF sentinel 與繞回。relation 更新保留 war bit，
政治分派保留原始 CF，不用抽象 boolean 取代最後 FLAGS。

## 3. 矩陣與結果

| 群組 | O0/O2 各原版/C |
|---|---:|
| 任官：23×23 新／舊 byte 含 FF × 3 個計數邊界 | 1587 相同 |
| writer：合法 cursor／slot、四種空槽形狀 | 60 相同 |
| 力量：5 個兵池 × 5 個據點數 × 5 個資金高位 | 125 相同 |
| relation：22×22×2 入口 × 6 個 relation byte | 5808 相同 |
| 入仕／俘虜：256 RNG 初值 × 2 入口 × 4 種歸屬欄位 | 2048 相同 |
| 武將更新：128 槽含終止槽 × 4 個 status × 3 個 timer | 1536 相同 |
| 邊境：22 勢力 × 2 入口 × 4 relation × 普通／中立鄰接 | 352 相同 |
| 十個政治／通知入口 × 8 種條件 | 80 相同 |
| 滿足協力前提的相鄰關係門檻 | 8 相同 |
| 16 個 queue／世界樣本的完整月度初始化 | 16 相同 |
| 四劇本 × 3 個玩家選擇的完整 C 月結規則 | 12 相同 |
| 合計 | **11,632 相同** |

每版 91 次完整 1 MB 記憶體抽查相同。所有案例核對 14 個暫存器／段／FLAGS、世界、
queue、scratch、globals、cadence、RNG 與持續堆疊；每次 callee 入口另比記錄、SS stack 與 globals。
原版/C 全量一致不代替 Go 政治層或 GUI 的同狀態比較，本輪沒有新的完整 Go 月結完成聲明。

O0/O2 JSON 相同，SHA-256：`ad2b970892169f03a17ff5bc80babd1a8cb3ba62d948f88189cfb78b024fe2c5`。
兩側狀態與入口軌跡摘要：`8eab99c52669e98448e29ba0d6c2bd8ea9be9c81241b0ae93441d59d61ac6b3b`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| writer 用整個 word 判空 | writer 第 4 組 |
| relation 刪 war bit | relation 第 4 組 |
| queue 從 offset 0 前移 | initialize 第 1 組 |
| 任官多加 1 | assignment 第 1 組 |
| 力量上界改 2001 | power 第 3 組 |
| timer 減 2 | generals 第 8 組 |
| 邊境排序反轉 | frontier 第 1 組 |
| 協力門檻改 A4 | cooperation 第 1 組，含原始 FLAGS |

協力門檻突變曾通過原先 80 個政治樣本；那些輸入沒有進入協力深層分支。
補上滿足前提的八個相鄰門檻後，正常版本通過、錯誤版本被拒絕，沒有放寬閘門。

## 4. 來源與重跑

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_politics_probe.py workplace/matching-decompilation/c-politics/ida KI.EXE
tools/c_recovery_politics.sh
```

Go 1.26.7、GCC 12.2.0，原版、專案與 dosgolem 唯讀，非 root、限資源 Docker。
dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`；import 來源前後雜湊相同。
RNG 表用 12:34:56，各案例 c/s 執行前固定且兩側相同，沒有重擲或改正式種子。

| 來源 | SHA-256 |
|---|---|
| C | `3b60a81cd5cc0979ca99f9bb03c5ce341f36b539ad318e70cced5ac516aa0eef` |
| header | `291e2425b3056eac0c107b93da341f2ad5e4bd32b9d71275e0c007b4f95b9d8c` |
| fixture | `1241042349ce6d9ea5b08355a53b9fdec54171a44f0501ccfd11b17d4b1f8171` |
| Go 對照器 | `1af1b0ab9df0a217c548fe36913bc09ec50be41fb55283614d42cc2bd3379887` |

本機產物在 `workplace/matching-decompilation/c-politics/`，驗證器為
[`c_recovery_politics_verify.py`](../../tools/c_recovery_politics_verify.py)。現行完成台帳查
[`c-recovery-status.json`](c-recovery-status.json)。前輪的 fixture 與原版收據保留。

## 5. 未解範圍

| 項目 | 邊界 |
|---|---|
| UI、音效、重畫及正常玩家長程 | 未驗證，完整月結規則接線不等於正常玩家垂直鏈 |
| 完整 Go 政治／俘虜同狀態比較 | 未驗證，本輪原版/C 收據不外推 Go parity |
| 非法記錄、無勢力、索引／段重疊 | 未驗證，不用無限尋找的 prototype 猜補終止條件 |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知，整檔匹配仍由組語基準提供 |
