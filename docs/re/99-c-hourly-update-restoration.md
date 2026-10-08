# 99：C 每時世界更新與事件分派核心

**狀態：八函式 O0/O2 各 21,871 組原版/C 相同，七個負對照拒絕；限明示 handler／UI fixture。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-hourly/ida/input.exe.i64`
- DB SHA-256：`d15aa9875f3821827d3a5d76edf4019c34ae2f599322523e26541406637b2c7b`
- 原始函式／operand／chunk／xref 探針：[`ida_hourly_probe.py`](../../tools/ida_hourly_probe.py)
- 契約：[`spec/218`](../spec/218-c-hourly-update.md)

## 1. 原始定位與已證實語意

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 局部語意 |
|---|---|---|---:|---|
| `sub_13E11` | `0x13E11` | `0x4011` | 84 | event dispatch、財政／LowFunds、維持費、外交官、cursor、redraw 的原始次序 |
| `sub_13E65` | `0x13E65` | `0x4065` | 41 | 24-bit 兵池和、carry、五次 SHR／RCR、原始支出 helper |
| `sub_13E8E` | `0x13E8E` | `0x408E` | 111 | 兩個 RNG 閘、官員 budget byte、雙向 relation 更新 |
| `sub_15673` | `0x15673` | `0x5873` | 34 | 24-bit 支出累加與上界鉗制、AX／DX 保存 |
| `sub_131AE` | `0x131AE` | `0x33AE` | 68 | cadence、cursor、queue stride、code.low、原始十三碼間接表 |
| `sub_13496` | `0x13496` | `0x3696` | 16 | event 10 的 formatter、CX、SS stack 與 return |
| `sub_13507` | `0x13507` | `0x3707` | 19 | event 13 通知後的信賴度處理與原始參數 |
| `sub_13DC9` | `0x13DC9` | `0x3FC9` | 72 | trust byte borrow 歸零、通知／退出順序及保存暫存器 |

原始指令共 445 bytes，每函式的 bytes SHA-256 在 verification.json 核對 EXE 與 IDA chunk。
推論等級已證實，限合法 record／code.low 0–13、分離段與 IF/TF=0。

## 2. 接線與邊界

[`hourly.c`](../../tools/c_recovery/hourly.c) 接到既有 C 政治／世界／RNG。
原版與 C 均先執行事件分派，再取當時的 hourly faction cursor。
財政檢查只對存在勢力執行，原始兩支維持費／外交官 callee 在此 gate 之後仍被呼叫。
門檻以資金高 16 位做有號比較，byte status bit 6 先清後按原始分支設置。

維持費先將三個 word 兵池加為 AX/DL，兩次 carry 累入 DL，再五次 SHR DL／RCR AX。
支出 helper 保留原始 24-bit 加法、上界與最後 FLAGS，不用普通 word 和代替。

dispatcher 衰減 cadence，為零且 cursor<0x100 才取一筆 record，cursor 加 4，
低 byte=0 就返回，其餘按原始 CS:31F2 word table 呼叫。
event 10／13 與 trust helper 由 C 真實執行；其餘 11 個 handler 與 UI／音效／退出／重畫為 RET fixture。
本輪不宣稱其他政治／戰術事件 handler 或正常 modal 行為完成。

## 3. 完整矩陣

| 群組 | O0/O2 各原版/C |
|---|---:|
| 支出：五個初值 × 四個高位 × 五個 low-word 邊界 | 100 相同 |
| 維持費：三兵池各取五個 word 邊界 | 125 相同 |
| 派發：全部 256 cadence × 四個 cursor × 14 code.low 值 | 14,336 相同 |
| 外交官：256 個 c/s × 五個 budget × 五個 politics | 6400 相同 |
| trust 與 event 10／13：borrow、上界與原始通知入口 | 30 相同 |
| 每時：22 cursor × 四個 status × 七個資金 | 616 相同 |
| 四劇本 × 三玩家 × 22 cursor，每次同狀態一個每時入口 | 264 相同 |
| 合計 | **21,871 相同** |

每版 171 次完整 1 MB 抽查相同；所有案例驗 14 個暫存器／段／FLAGS、世界、globals、
queue、RNG、cadence 與 stack，每 callee 入口另驗 raw 參數與記錄。
首次兩個 trace 差異的暫存器／資料已一致，根因為重複監看及漏登錄 `sub_1310A`；
監看修正後用同一資料與判準完整重跑。負對照中的 cadence／globals 差異另明列於失敗收據。

O0/O2 JSON 相同，SHA-256：`9fd7647d0786005b2d62245012d80678a84cc296f911e1e55e14abd59c1cae51`。
兩側狀態與入口軌跡摘要：`43a185cd0d963205c40edeeb43dd569562fbafbab302da3d9016f64e6434b921`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 支出上界少 1 | expense 第 14 組 |
| 第一個兵池 carry 丟棄 | maintenance 第 46 組 |
| cadence 重設 9 | dispatch 第 57 組 |
| queue stride 改 2 | dispatch 第 57 組 |
| trust borrow 不歸零 | trust 第 2 組 |
| 維持費與外交官順序交換 | world 第 1 組 |
| 外交官 byte 經費係數改 24 | diplomat 第 231 組 |

## 4. 重跑與來源

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_hourly_probe.py workplace/matching-decompilation/c-hourly/ida KI.EXE
tools/c_recovery_hourly.sh
```

Go 1.26.7、GCC 12.2.0，原版／專案／dosgolem 唯讀、非 root、限資源 Docker。
dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import 來源前後雜湊相同。
RNG 表固定 12:34:56，每案例 raw c/s 於執行前同時寫兩側，沒有重擲或挑 seed。

| 來源 | SHA-256 |
|---|---|
| C | `dc115558ed470df49bd4438825b059d835a707e4f93262406ea8fb26c9536c04` |
| header | `1388fe3fe2c6bed876246e08fdfba0814ed948ec50e18ba236c0553626f3aff1` |
| fixture | `54864eeb3e0aeabcae321ffc537babe3e1abf8f53cb3f1d4e7cea1daebf36cb5` |
| Go 對照器 | `aa54aeed36a71497b0db18e5c13d91f6fbd7cfe11a2ffa6511caf1b870afb666` |

本機產物在 `workplace/matching-decompilation/c-hourly/`，收據核對入口是
[`c_recovery_hourly_verify.py`](../../tools/c_recovery_hourly_verify.py)，C 現況查
[`c-recovery-status.json`](c-recovery-status.json)。正式 Go 本輪沒有修改。

## 5. 未解範圍

後續 C handler 證據見 [`re/100`](100-c-event-handlers-restoration.md)。本輪舊收據的十一個
handler RET 保留；re/100 的新收據覆蓋真實 handler 與依賴，UI／戰術／正常玩家邊界仍未驗證。

| 項目 | 邊界 |
|---|---|
| 其他十一個 event handler 與完整 modal／UI／音效 | 未驗證，原始 table call 不替代 callee 完成 |
| 完整 Go 每時規則與月結前後三方 | 未驗證，本輪收據只量原版/C 核心 |
| 正常玩家日期流、硬體等待與長程 | 未驗證，one-hour 直接入口不取代玩家垂直鏈 |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知，整檔 match 仍是組語基準 |
