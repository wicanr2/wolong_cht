# 106：C 矩形、選取框與戰術計量條

**狀態：十六函式 O0/O2 各 134,958 組完整 RAM／四 plane／ABI 相同，原版／C／Go 131,072 個計量相同，十二個錯版拒絕。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V KI.EXE
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具：IDA Pro 9.4，database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-rect/ida/input.exe.i64`
- DB SHA-256：`3f3a95793967805c9df4f87de5db577b6eadb0aa3a39d4837be8c9d43a8dc9a5`
- 探針：[`ida_rect_probe.py`](../../tools/ida_rect_probe.py)
- 契約：[`spec/225`](../spec/225-c-rectangle-bars.md)
- 既有證據：[`re/18`](18-tactical-button-glyphs.md)、[`re/60`](60-tactical-sidebar.md)、[`re/105`](105-c-aligned-blit-restoration.md)

## 1. 原始定位與控制流

十六個新函式共 1,169 bytes，連同既有依賴與素材指標分配共核對 36 個函式。
四筆 MZ relocation 分別位於新的滑鼠查詢 caller 與既有裝置保存／恢復。
原始名稱、完整 chunks、operand、xref、file／IDA bytes 保留。以下 raw 控制流為已證實。

| 原始函式 | IDA 起點 | 檔案起點 | bytes | 控制流 |
|---|---|---|---:|---|
| `sub_1F020` | `0x1F020` | `0xF220` | 288 | 有號端點排序、640×400 裁切、四個 clip bits 與矩形描框 |
| `sub_1F140` | `0x1F140` | `0xF340` | 59 | 依左右 clip bits 寫兩條垂直邊，保留底端給水平 helper |
| `sub_1F17B` | `0x1F17B` | `0xF37B` | 39 | 首／中／尾 byte 的水平邊，REP STOSB 與 dummy latch read |
| `sub_1F1A3` | `0x1F1A3` | `0xF3A3` | 203 | 有號端點排序／裁切與包含端點的實心矩形 |
| `sub_10AAA` | `0x10AAA` | `0xCAA` | 47 | 彩色已填部分接黑色剩餘部分，兩列高 |
| `sub_10AD9` | `0x10AD9` | `0xCD9` | 109 | 像素 X 與 byte row offset 的列寫入／左右 mask |
| `sub_10CAC` | `0x10CAC` | `0xEAC` | 23 | write mode 3／Data Rotate／Bit Mask 的原始 I/O |
| `sub_10CC3` | `0x10CC3` | `0xEC3` | 27 | write mode 0 與 set/reset／mask 的原始恢復序列 |
| `sub_1C61F` | `0x1C61F` | `0xC81F` | 52 | 十六格選取的原始座標與 14×14 描框 caller |
| `sub_1C6BF` | `0x1C6BF` | `0xC8BF` | 55 | 位置表選隊，外／內兩層選取框 |
| `sub_1C6AE` | `0x1C6AE` | `0xC8AE` | 17 | 六隊選取框清除 caller |
| `sub_1C6F6` | `0x1C6F6` | `0xC8F6` | 86 | 固定滑鼠查詢、六隊待機條與四條側欄計量條 |
| `sub_1C74C` | `0x1C74C` | `0xC94C` | 41 | 六個待機值，各以 byte 上限 76 前進來源 |
| `sub_1C775` | `0x1C775` | `0xC975` | 25 | 兵力 word 右移兩次，取 AL 後以 byte 上限 124 |
| `sub_1C78E` | `0x1C78E` | `0xC98E` | 27 | 體力先右移一次，再加第二次右移，取 CL 後上限 124 |
| `sub_10BCD` | `0x10BCD` | `0xDCD` | 71 | 真實外框 caller，內部矩形填黑 |

描框只畫原始可見邊，裁切掉的邊不在 viewport 邊界重造。實心矩形則填滿包含端點的
裁切區域。兩者都使用原始 mask 與成熟 VGA bus，不能由平面 RAM 或一般線框代替。

計量條的取 byte 與位移順序是原始控制流的一部分。體力 3 經兩次截斷得到 1，
乘三除四會得到 2；原始 word 右移後取 byte 也不能換成不限位寬的飽和計算。
Go 顯示 helper 的修正另由 [spec/226](../spec/226-sidebar-bar-widths.md) 管理。

## 2. 素材配置與歷史收據

ICONGRF.DAT 的完整 SHA-256 與四段大小沿用 [formats/03](../formats/03-grf-images.md#5-icongrfdat四段組合檔部分解)。
`sub_100DF` 的分配以 paragraph 計算：`word_10D48` 先指向第 3 段起點，
加 0x6C paragraphs 後 `word_10D4A` 指向段內 0x6C0 的框圖塊。
本輪 C fixture 必須使用第 3 段命令圖示與框圖塊；六個位置 glyph 與背板仍取第 1 段。
較早 [re/105](105-c-aligned-blit-restoration.md) 的單一重畫／外框矩陣供應了第 1 段資料，
那份收據只證明同輸入的 raw caller 行為。本輪補真實 bank 的證據，保留歷史收據與來源雜湊。

## 3. 驗證契約

原版 guest CPU 與 C native 算法使用兩個獨立 RAM／VGA／latch 狀態。
沿用 [re/104](104-c-vga-blit-restoration.md#2-平台前提與獨立狀態) 的固定成熟平台契約，
IF/TF=0，底層滑鼠服務為固定明示 FAR fixture。畫面抽樣取 VGA 第 40 列起的 640×400。
O0/O2 需比較完整 RAM、四 plane、十四暫存器、FLAGS、stack、GC／seq／latch、port 與原始入口。
矩形反向端點／signed 邊界／單 byte／裁切旗標、兩層選取、真實素材、原始計量條與 Go 長度均須驗證。

## 4. 未解範圍

| 項目 | 邊界 |
|---|---|
| 正常玩家戰術／視窗操作 | 局部原始 caller 不代替自然流程 |
| 原始滑鼠時序與輸入裝置 | FAR fixture 只固定本輪回應 |
| 完整文字與其他繪圖路徑 | 本輪未涵蓋所有 renderer |
| C 機器碼匹配與原作者工具鏈 | 尚未驗證 |

## 5. 動態收據與分派 code 補充

後續顯示清單證據見 [re/107](107-c-display-interpreter-restoration.md)，沿用矩形／圖庫 C，
再補全部 handler、原始 DS near table 与未分類線段／底紋。較早二十 bytes 的來源快照保留。

O0/O2 各 134,958 組完整 1 MB RAM、四個 65536-byte plane、十四暫存器／FLAGS、
GC／seq／latch／index、port 與 90-byte 入口快照相同；2,441 次內容區像素相同。
兵力與體力各 0..65535，共 131,072 組原版／C／正式 Go 函式長度相同。
Go 函式與上限常數從正式來源逐字擷取，來源與 CGO flags 由 build info 固定。

| 群組 | 每種最佳化的組數 |
|---|---:|
| 反向／裁切／signed 矩形 | 3,072 |
| 首尾／垂直邊 helper | 108 |
| mode 設定／恢復 | 16 |
| 完整 word 計量 | 131,072 |
| 兩色計量 caller | 160 |
| byte 列／零高度 | 192 |
| 十六格／兩層／六隊選取 | 184 |
| 外框接實心內部 | 18 |
| 滑鼠區域／四側欄／六待機 | 75 |
| 原始待機來源 | 5 |
| 第 3 段九命令圖示重畫 | 54 |
| 第 1 段六位置 glyph／背板 | 2 |
| 合計 | **134,958** |

O0/O2 JSON 逐 byte 相同，SHA-256：`07927210a9bddcb437aafceef432bca42c861a70d0fff59d1844f1ef6f0477fa`。

| 刻意改錯 | 首次拒絕 |
|---|---:|
| 描框寬少算端點 | rectangle 第 1 組 |
| 不保留裁切邊消失旗標 | rectangle 第 289 組 |
| 實心高度少算端點 | rectangle 第 1537 組 |
| 實心右 mask 用 SHR 代替 SAR | rectangle 第 1633 組 |
| 刪 dummy latch read | rectangle 第 1 組 |
| write mode 改 0 | rectangle 第 1 組 |
| 體力改乘三除四 | gauge 第 8 組 |
| 上限改 125 | gauge 第 1 組 |
| 待機來源步幅改 3 | waiting 第 1 組 |
| 兩列 pitch 改 79 | bar 第 5 組 |
| 描框不清 DF | rectangle 第 2 組 |
| 少畫內層選取 | selection 第 129 組 |

原始 DS 間接 table 位於 IDA `0x1EA0D`，九個 word 的第 2／3 筆指向
`0x1EA26`／`0x1EA30`。這二十 bytes 在 IDA 被標成資料，通常匯出不含它們。
窄探針首次用了資料項的空 mnemonic 而無法停止，改取已解碼 insn 的 canonical mnemonic
後，保留原資料行、分類、bytes 與運算元，再確認兩個 near call 分別到 `0x1F020`／`0x1F1A3`。
共同轉換目標是 `0x1EAD0`；本輪未還原完整 interpreter。

補充探針：[ida_rect_handlers_probe.py](../../tools/ida_rect_handlers_probe.py)。

本次補充 DB 位於 `workplace/matching-decompilation/c-rect/handlers-ida/input.exe.i64`，
SHA-256 為 `694de3641d45608da76408a083e42f5036edc5136c60906b83aa93e5609916b7`。
八條指令／20 bytes 在 [rectangle-handler-code.json](rectangle-handler-code.json) 保留，
已合併至 [KI.code.S](../../tools/c_recovery/KI.code.S)。新版 24,384 條／55,412 bytes
全部組譯，完整 EXE 仍相同；已辨識的 handler 不再從原版當資料匯入。

C [rect.c](../../tools/c_recovery/rect.c)、[header](../../tools/c_recovery/rect.h)、
[fixture](../../tools/c_recovery/rect_fixture.h)、[driver](../../tools/c_recovery_rect.go)、
[重跑入口](../../tools/c_recovery_rect.sh) 與 [驗證器](../../tools/c_recovery_rect_verify.py) 進版控。
[驗證摘要](c-rect-verification.json) 固定來源、原始入口、完整矩陣與錯版收據。

Compiled source manifest 位於 `workplace/matching-decompilation/c-rect/results/c-source.sha256`，
SHA-256 為 `1cb7c2ee14ae9ad572b882a994eab6223c518ac8bea2167a02b20fddf2de0c5a`。
Go 1.26.7／GCC 12.2.0／IDA9.4；正式 UI 在既有 hr-go-ebiten SDK／Xvfb 冷測通過。
所有原版、素材與來源唯讀，非 root／限資源 Docker；測試裝置、自然玩家與 C 機器碼界線保留。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_rect_probe.py workplace/matching-decompilation/c-rect/ida KI.EXE
tools/c_recovery_rect.sh
```
