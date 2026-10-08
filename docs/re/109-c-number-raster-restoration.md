# 109：C 數字 raster 與日期／參數入口

**狀態：四個原始 C 函式／一個 code 入口通過，O0/O2 各 4,881 組全 RAM／VGA／ABI 相同，八個錯版拒絕。**

- 日期：2026-10-09
- 輸入：DOS/V KI.EXE，67,099 bytes
- SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear，檔案偏移為線性位址減 `0x10000` 加 512
- DB：`workplace/matching-decompilation/c-numbers/ida/input.exe.i64`
- DB SHA-256：`07acecae5db3e134e1892c0c4b2618a4577fe214683101db5a933fb8ce657513`
- 入口：[`ida_numbers_probe.py`](../../tools/ida_numbers_probe.py)
- 前置證據：[`re/28`](28-text-number-rendering.md) 與 [`re/108`](108-c-glyph-raster-restoration.md)
- 契約：[`spec/229`](../spec/229-c-number-raster.md)

`sub_1062F` 保留 DX:AX、BL 寬度、BH 顏色與原始 signed 補數處理。
`sub_1069A` 讀取真正十個數字及負號字模，`sub_106DE` 填背景。
日期函式 `sub_11E17` 與無名稱參數區 `0x10984..0x109AF` 保留原始定位。
固定 DB 的函式台帳確認日期函式邊界；參數區仍另列 code 入口，不推測原作者 C 邊界。

| 原始定位 | IDA 線性起點 | 檔案起點 | bytes | 推論等級與證據 |
|---|---|---|---:|---|
| `sub_1062F` | `0x1062F` | `0x82F` | 107 | 已證實的 raw 數字控制，固定 probe／re/28 |
| `sub_1069A` | `0x1069A` | `0x89A` | 68 | 已證實的數字／負號前景與背景 raster，固定 probe |
| `sub_106DE` | `0x106DE` | `0x8DE` | 23 | 已證實的背景寫入，固定 probe |
| `sub_11E17` | `0x11E17` | `0x2017` | 47 | 已證實的三個原日期欄位呼叫，固定 probe／函式台帳 |
| 原始無名稱 code | `0x10984` | `0xB84` | 43 | 已證實的 SS word／CWD／座標計算，固定 probe；名稱未知 |

五個 C body 包含 143 條原始指令，沒有執行期 opcode 分派。日期直接讀原 DS:`0xCF6`／`0xCF4`／`0xCF0`，
參數入口讀 SS:`[di]` 後 CWD，保留輸入來源與原始寬度／顏色。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原作者 C 工具鏈與機器碼 | native C 語意測試不證明編譯匹配 |
| 自然玩家日期／訊息流程 | 局部入口不證明完整 UI 操作 |
| 超出 DIV quotient 範圍 | 原版可能產生除法例外，本輪不替原版猜補 |

## 同狀態驗證

兩側持有獨立 RAM／VGA，原版執行 guest 指令，C 用 native 運算與裝置 bus。
每組比較全 1 MB RAM、四個 65,536-byte plane、GC／seq／latch／port、十四 register／FLAGS、
原始入口與 32-byte stack 快照，以及 640×400 內容區像素。IF／TF 固定為 0；
未定義旗標以固定 dosgolem 模型比較，不宣稱實機時序。

| 群組 | 每種最佳化組數 |
|---|---:|
| 正負值／十進位邊界／六種欄寬／主入口 DF | 312 |
| 256 顏色與 plane mask／read map | 1,280 |
| 11 個實際數字／負號字模與 256 顏色 | 2,816 |
| 背景 primitive | 256 |
| word 邊界／直接 primitive 的 DF 與 CH | 192 |
| 原始日期 caller | 4 |
| SS word 數值參數 caller | 21 |
| 合計 | 4,881 |

C 除法不放寬 AX quotient，寬度 0 保留原 byte 回繞，寬度 1 的負數保留原保護分支。
數字字庫為本機 ICONGRF.DAT 第三段 `+0x840` 的 11×16 bytes；字庫形狀與雜湊在執行前核對，
沒有將原版輸出當成 C 輸入。兩種最佳化的完整收據逐 byte 相同。

| 刻意改錯 | 群組內首次拒絕 |
|---|---:|
| 刪除負數補數 carry | 25 |
| 刪除欄寬右緣調整 | 1 |
| 讀取錯誤字庫 bank | 1 |
| 前景色清成 0 | 2 |
| 刪除 latch read | 1 |
| STOSB 忽略 DF | 178 |
| 列距少一 byte | 1 |
| 背景色讀 AX 而非 BH | 42 |

[公開摘要](c-numbers-verification.json) 保存輸入、來源、原函式、DB 與收據雜湊。
來源：[numbers.c](../../tools/c_recovery/numbers.c)、[header](../../tools/c_recovery/numbers.h)、
[生成 C](../../tools/c_recovery/numbers_generated.inc)、[產生器](../../tools/c_recovery_numbers_generate.py)。
[重跑入口](../../tools/c_recovery_numbers.sh) 與 [驗證器](../../tools/c_recovery_numbers_verify.py)
檢查實際編譯 flags 的 source digest、錯版與舊規格 backlink；生成 C 在乾淨容器重生相同。

Go 1.26.7、GCC12.2.0、IDA9.4、非 root 限資源 Docker，dosgolem 固定來源 revision
`a9714ebdab2ad6b529f81225472680f2b11f2842`；原始 EXE／字庫與外部平台來源唯讀。
正式 Go remake 本輪保持既有實作；renderer 與原始兩個 caller 的綠色收據不代證正常玩家流程。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_numbers_probe.py workplace/matching-decompilation/c-numbers/ida KI.EXE
tools/c_recovery_numbers.sh
```
