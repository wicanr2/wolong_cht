# 201：sub_1ECE0 的 C 還原契約

**狀態：CONFORMED。原版與 C 的局部介面對照、C 與 Go 的資料規則對照均通過。**

- 日期：2026-10-08
- 範圍：松崗 DOS/V `KI.EXE` 的亂數取數，限局部常式；不修改正式 Go 玩家路徑
- 輸入 SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 位址：IDA 9.4 線性 `0x1ECE0`–`0x1ECFC`；檔案 `0xEEE0`–`0xEEFC`
- 指令身分：28 bytes SHA-256 `f56e5a1d08956bc9c68476dd5245b85a7e0169d52757b5e476c32b28e1911a9d`
- 證據入口：[`re/10`](../re/10-rng.md)、[`re/90`](../re/90-assembly-reconstruction.md)

## 1. 證據審查

DRAFT 問題是 C 只回傳亂數值是否足夠。13 條原始指令還會改 AX、FLAGS、堆疊內容、
SP 與 IP，因此研究實作分成可讀的資料函式與 16-bit 呼叫介面適配層。
兩層都保持原始 `sub_1ECE0` 定位，不把新 C 介面當成原作者的介面。

| 原始定位 | 契約 | 等級與來源 |
|---|---|---|
| `byte_1ECFC`、`byte_1ECFD`、`0x1ECFE` 表 | `c`、`s` 與 256 byte 表 | 已證實，re/10 與固定 bytes |
| `0x1ECE0`、`0x1ECE1` | 先將 DS、BX 推入 SS:SP | 已證實，原始 PUSH 指令 |
| `0x1ECE2`–`0x1ECEC` | AX／DS 改成 CS；用表取值 | 已證實，原始 MOV／XLAT 指令 |
| `0x1ECED`–`0x1ECF6` | `value = T[s] + c`；先寫 `c += 0x89`，再寫 `s = value`，皆取低 8 bit | 已證實，原始 ADD／MOV 指令 |
| `0x1ECF9`–`0x1ECFB` | 從記憶體取回 BX、DS，再近返回 | 已證實，原始 POP／RET 指令 |

ADD 影響 CF、PF、AF、ZF、SF、OF；其餘指令不改這些運算旗標。離開時的運算旗標
來自 `c + 0x89`，不能用回傳亂數值的加法代替。CPU 契約來源：
[Intel SDM Vol. 2A ADD](https://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-software-developer-vol-2a-manual.pdf)，
3-31／3-32，查閱日期 2026-10-08。

## 2. typed input 與輸出

資料層收一份 258-byte 狀態：`uint8_t c`、`uint8_t s`、`uint8_t table[256]`。
回傳 `uint8_t`，且只更新 c、s。表不要求合法播種的置換，因為本契約驗的是取數本體。

適配層收 14 個 16-bit 暫存器／段／旗標欄位與 1 MB 記憶體。入口是近呼叫已推入
返回位址、IP 位於原始取數函式的時刻。完整返回後比較暫存器、旗標與記憶體。
本輪限定原始程式載入段、偶數 SP、表可連續讀取、入口 IF／TF 固定為 0，
堆疊不得覆蓋原始指令。DF 與六個運算旗標仍依案例變化。

一般堆疊下 BX、DS 取回原值。堆疊若與 c／s／表重疊，必須按實際記憶體讀寫處理，
不能無條件把本機變數中的舊 BX 還原。此為受控邊界實驗，不代表正常玩家的記憶體佈局。

## 3. 驗證閘門

| 驗證 | 判準 |
|---|---|
| 原版 oracle | dosgolem 載入原始 EXE，以真實 CS 的 `IDAIn` 指定入口；直接執行原始 13 條指令 |
| 初始狀態 | 執行前固定表、c、s、暫存器、旗標、SS:SP；沒有重擲或挑選結果 |
| C 適配層 | 14 個欄位與受影響記憶體逐項相同；原版寫入僅限狀態與堆疊 |
| C 資料層與 Go | 非重疊堆疊下，回傳值、258-byte 狀態與 `rng.Next` 相同 |
| 抽樣完整記憶體 | 比完整 1 MB，檢查未宣告的額外寫入 |
| 最佳化 | GCC `-O0`／`-O2` 都對同一組輸入通過 |
| 負對照 | 改計數器增量、AH、旗標來源或 BX 取回方式，必須被拒絕 |

本輪 C 為原生語意還原，沒有聲稱 C 產物機器碼與 DOS 原版相同。
局部呼叫不驗播種、取數呼叫順序、玩家流程或整款 remake 的行為。

## 4. 驗證收據

GCC `-O0`、`-O2` 各 263,680 組原版／C 對照與 263,168 組 Go 資料規則對照通過；
每個最佳化版本做 258 次完整 1 MB 記憶體核對。四種故意錯誤的 C 版本皆被拒絕。
具體輸入、工具身分、摘要與範圍見 [`re/91`](../re/91-c-rng-restoration.md)。
