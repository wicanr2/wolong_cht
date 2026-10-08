# 100：C 事件 handler 與政治／災害依賴

**狀態：三十函式 O0/O2 各 43,476 組原版/C 相同，八個負對照拒絕；限明示 modal／戰術／UI fixture。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-events/ida/input.exe.i64`
- DB SHA-256：`67e6b4f56f2faa716b9f18afd61260a20329374146d29a41a4204317eff6eb97`
- 原始 chunks／operand／xref 匯出：[`ida_events_probe.py`](../../tools/ida_events_probe.py)
- 契約：[`spec/219`](../spec/219-c-event-handlers.md)

## 1. 原始定位與證據等級

三十個原始函式共 1919 bytes。事件 1–9、11–12 的十一個 handler，加上 event 8
直接落入的 `sub_133FD`；十八個政治、任官、關係、災害及軍團／據點依賴均真實接線。
事件 10／13 沿用 [`re/99`](99-c-hourly-update-restoration.md)。

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 已證實的 raw 語意 |
|---|---|---|---:|---|
| `sub_1320C` | `0x1320C` | `0x340C` | 20 | 事件 1 勢力 gate、原始參數與目標更新 |
| `sub_13220` | `0x13220` | `0x3420` | 66 | 事件 2 三方 gate、外交回傳與目標更新 |
| `sub_13262` | `0x13262` | `0x3462` | 71 | 事件 3 外交回傳、清理、原始戰鬥 callee 與關係更新 |
| `sub_132A9` | `0x132A9` | `0x34A9` | 64 | 事件 4 內政官 byte sentinel 與原始金額視窗參數 |
| `sub_132E9` | `0x132E9` | `0x34E9` | 62 | 事件 5 外交官 byte sentinel 與原始金額視窗參數 |
| `sub_13327` | `0x13327` | `0x3527` | 97 | 事件 6 外交官、對玩家 gate、選擇與關係更新 |
| `sub_13388` | `0x13388` | `0x3588` | 98 | 事件 7 勢力／外交官 gate、選擇與目標更新 |
| `sub_133EA` | `0x133EA` | `0x35EA` | 19 | 事件 8 勢力 gate、據點選取及原始落下尾端 |
| `sub_133FD` | `0x133FD` | `0x35FD` | 136 | 遷都 raw 欄位、軍團 helper、兩種訊息的原始 formatter 參數 |
| `sub_13485` | `0x13485` | `0x3685` | 17 | 事件 9 官員記錄取址與原始清理 |
| `sub_134A6` | `0x134A6` | `0x36A6` | 11 | 事件 11 RNG、raw 24–39 storm 參數 |
| `sub_134B1` | `0x134B1` | `0x36B1` | 86 | 事件 12 災害建立／清理、RNG、延遲 queue writer |
| `sub_1351A` | `0x1351A` | `0x371A` | 12 | 以 AH 取勢力、byte 存在門檻與 CF |
| `sub_13526` | `0x13526` | `0x3726` | 133 | 原始目標 gate、通知參數與雙向關係更新 |
| `sub_135AB` | `0x135AB` | `0x37AB` | 66 | raw target 的 0x24 比較、間接 power 與 byte 寫入 |
| `sub_135ED` | `0x135ED` | `0x37ED` | 76 | 依回傳做 24-bit 金額移轉與原始官員清理 |
| `sub_13639` | `0x13639` | `0x3839` | 48 | 雙向關係 raw 最小值、清高位與 SHR |
| `sub_13669` | `0x13669` | `0x3869` | 46 | 雙向關係 raw 最小值與設高位 |
| `sub_13697` | `0x13697` | `0x3897` | 45 | 保留原始 byte carry 的雙向 relation 取址 |
| `sub_136C4` | `0x136C4` | `0x38C4` | 78 | 原始 byte 停戰條件、提案金額與保存暫存器 |
| `sub_13712` | `0x13712` | `0x3912` | 95 | 原始 byte 協力條件、提案金額與保存暫存器 |
| `sub_13771` | `0x13771` | `0x3971` | 103 | 外交官／官員取址、能力 gate、平手 RNG 與 byte score |
| `sub_137D8` | `0x137D8` | `0x39D8` | 29 | 雙向原始殘留記錄 helper 與 AH bit 組合 |
| `sub_137F5` | `0x137F5` | `0x39F5` | 59 | 127 官員 raw owner／狀態／能力最大值選取及 CF |
| `sub_13138` | `0x13138` | `0x3338` | 55 | 127 官員 raw pair 計數與原始 cmp ax,ax／CF |
| `sub_14502` | `0x14502` | `0x4702` | 70 | 127 軍團 raw +0x20／+0x14 欄位比較及更新 |
| `sub_16A3D` | `0x16A3D` | `0x6C3D` | 94 | 192 據點 raw owner、旗標／分類／生產力選取及 CF |
| `sub_123FF` | `0x123FF` | `0x25FF` | 57 | 十六災害槽 first-free、欄位寫入、保存與 CF |
| `sub_12438` | `0x12438` | `0x2638` | 33 | 同座標所有災害槽清理 |
| `sub_150D7` | `0x150D7` | `0x52D7` | 73 | 官員欄位、原始任官計數與通知參數 |

原始名稱、位址、operand、chunks、原版 bytes 與 source hashes 保留。
原始控制流與 byte 結果為已證實；高層玩法名只作導覽，不代替欄位 consumer 證據。
`sub_135AB` 的 `0x24` 及間接讀取照原版保留。raw target byte 0x23／0x24／0xFF
在 1 MB RAM 中受控比較，不據此推導合法勢力或中立玩家規則。GitHub Issue #27 維持 OPEN。
`sub_13138` 末端的 `cmp ax, ax` 照原版保留，沒有改成零值判斷。

## 2. 接線與明示外部邊界

[`events.c`](../../tools/c_recovery/events.c) 包含原始 C 名稱入口與可供既有 hooks 使用的 body。
原始 dispatcher 的十三碼 table 及每一個 handler 均執行真實 C。
event 8 從 `sub_133EA` 直接落入 `sub_133FD`，沒有新增 CALL／RET；收據驗這個入口的
原始 stack 與 IP。既有 C 經濟、任官、power、RNG、writer、storm helper 真實執行。

以下 callee 保留明示 fixture：

| IDA 線性位址 | fixture 契約 |
|---|---|
| 0x13C3D | `mov al, cs:[7000h]`／RET，回應 0–3，FLAGS 不變 |
| 0x138C7、0x138E6、0x139E8、0x12078、0x120D6 | 金額與 modal 視窗 RET，原始 caller／參數仍完整比較 |
| 0x145F8、0x14236 | 戰鬥準備 RET，caller 的 AX／SI／DI／stack 保留 |
| 0x15E60、0x15E80、0x10CE7、0x10CDE、0x18810 | 重畫／聲音／文字 RET，所有入口參數與堆疊比較 |
| 0x102F5、0x187FF、0x11CB1 | 既有 event 13 信賴度 callee 的 RET 邊界 |

正式 Go 本輪沒有修改，modal／戰術 fixture 不當作正常玩家驗收。
較早 [spec/218](../spec/218-c-hourly-update.md)、[re/99](99-c-hourly-update-restoration.md)
與 [spec/64](../spec/64-capital-relocation-report.md) 已加同版原始位址的範圍 backlink，
新 verifier 自動核對入口。它們的歷史收據與玩家功能範圍保留，不重新宣稱 UI 完成。

## 3. 完整矩陣

| 群組 | O0/O2 各原版/C |
|---|---:|
| presence：全部 byte 存在旗標 × 22 勢力 | 5,632 相同 |
| relation：雙向 raw 關係、原始 sentinel 與取址 | 4,620 相同 |
| target：raw 0x23／0x24 邊界、玩家／中立與 power | 240 相同 |
| official：全部 raw RNG × 能力 × 指定／選取官員 | 2,560 相同 |
| diplomacy：16 關係 × 6 policy × 2 官員 × 8 RNG × 2 gate | 3,072 相同 |
| disaster：兩 helper × 17 槽占用 × 四座標類型 | 136 相同 |
| disaster-event：17 槽占用 × 四碼 × 全部 raw RNG | 17,408 相同 |
| capital：三入口 × 22 勢力 × 32 據點／軍團邊界 | 2,112 相同 |
| cleanup：五 helper × 128 官員邊界 | 640 相同 |
| handlers：十二入口 × 三玩家 × 四 status × 四回應 × 三官員 × 三 RNG | 5,184 相同 |
| scenario-dispatch：四劇本 × 三玩家 × 13 碼 × 三 raw RNG × 四回應 | 1,872 相同 |
| 合計 | **43,476 相同** |

所有案例比較十四暫存器／segment／FLAGS、世界、queue、globals、RNG、cadence 與 stack。
每 callee 入口另比原始參數、record、SS stack、globals／RNG；每版 340 次完整 1 MB
抽查相同。三十個函式都有實際入口次數，保存在 `entries_seen`，不是由案例名稱推定覆蓋。
四劇本案例使用受控 queue 與回應，均經真實 dispatcher；不稱自然 producer 或正常玩家流程。

O0/O2 JSON 相同，SHA-256：`a33150b0c1297c3565431e559e12fc9692a35907f532c053b63247cefd3612c1`。
原版及 C 分別讀取各自記憶體計算摘要，兩側同為
`467c60e1141d88e52ec6e89f573507eee75f027204d295b7a32e6f68c9088cbc`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 勢力存在門檻改 0x81 | presence 第 1 組 |
| 關係取消不做 SHR | relation 第 11 組 |
| raw target 比較改 0x23 | target 第 79 組 |
| 災害槽容量減為 15 | disaster 第 61 組 |
| 軍團 raw 相符分支不更新 | capital 第 705 組 |
| 災害初值 timer 改 15 | disaster 第 1 組 |
| 殘留 count 改 cmp ax,0 | cleanup 第 257 組 |
| 協力上界改 59 | diplomacy 第 1537 組 |

## 4. 重跑與來源

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_events_probe.py workplace/matching-decompilation/c-events/ida KI.EXE
tools/c_recovery_events.sh
```

Go 1.26.7、GCC 12.2.0；IDA image `sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab`；
Go image `sha256:e8c859f5632dcfde7b32d2012b4351728f6437930887c2f6a91ea242459e5514`。
全部 workload 使用非 root、無網路、限資源 Docker，原版、專案及 dosgolem 唯讀，
只明確的輸出可寫。dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import 來源前後雜湊相同。
兩側每案例執行前寫相同 raw c/s，表固定 12:34:56，沒有重擲或挑 seed。
SHR／logic 的未定義 AF 沿用 dosgolem 模型，IF/TF=0，stack／資料段分離。

| 來源 | SHA-256 |
|---|---|
| events.c | `69c7a72c1168a3943824ced273702ace069ca3e7e42e3789baf0621341ef32b0` |
| events.h | `0845440ca2f3685f899772dc6ff78afef012aa3fdb4374f2611f1ca9175a6373` |
| events_fixture.h | `fdde12c4e3239aa79fe826c8def65edb059d54380352b03462cba427f597a676` |
| Go 對拍工具 | `0492e4bd0ecee39b0b5fc5bed4bc1fb4545d49a2d434b3a2dde371bde22f7279` |

本機產物在 `workplace/matching-decompilation/c-events/`，收據驗證入口是
[`c_recovery_events_verify.py`](../../tools/c_recovery_events_verify.py)。現行 C 分級台帳
[`c-recovery-status.json`](c-recovery-status.json) 共八十四個函式。

## 5. 未解範圍

| 項目 | 邊界 |
|---|---|
| 完整 UI／選擇／戰鬥／金額 callee | 本輪量 caller 與受控回傳，不能宣稱正常玩家 parity |
| 完整 Go 每時與事件三方 | 尚未驗證 |
| 中立 0x24 的玩家規則 | 原始 byte 行為與玩法語意分開，Issue #27 維持 OPEN |
| C 機器碼匹配與原作者工具鏈 | 尚未驗證，組語整檔基準仍是 binary match |
