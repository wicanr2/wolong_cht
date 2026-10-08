# 105：C 位元對齊貼圖與戰術按鈕呼叫鏈

**狀態：五個新 C 函式 O0/O2 各 4,954 組完整 RAM／四平面／ABI 相同，十一個負對照拒絕；六按鈕與外框原始 caller 接線通過。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V KI.EXE，不外推 PC-98
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- ICONGRF.DAT SHA-256：`2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d`
- 工具：IDA Pro 9.4，database linear 位址；檔案偏移另列
- DB：`workplace/matching-decompilation/c-aligned/ida/input.exe.i64`
- DB SHA-256：`0d0bd6bded0fe104b1648c8fa3af852ef209b57ee7e1464b84b786dc0f2affa8`
- 探針：[`ida_aligned_probe.py`](../../tools/ida_aligned_probe.py)
- 契約：[`spec/224`](../spec/224-c-aligned-blit.md)
- 既有證據：[`re/18`](18-tactical-button-glyphs.md)、[`re/103`](103-c-hotspot-restoration.md)、[`re/104`](104-c-vga-blit-restoration.md)

## 1. 原始定位與證據審查

本輪五個新函式共 429 bytes。另核對十三個既有 C 依賴，十八個原始函式共有三筆 MZ
重定位，全部位於已還原的裝置保護鏈。原始名稱、完整 chunks、bytes、operand 與 xref 保留。
新函式沒有重定位；原始 bytes 與 IDA 匯出逐 byte 核對。以下控制流為已證實。

| 函式 | IDA 起點 | 檔案起點 | bytes | 原始控制流 |
|---|---|---|---:|---|
| `sub_1F888` | `0x1F888` | `0xFA88` | 176 | X 低三位、左右遮罩、四平面 blit；保存除 SI 外的原始暫存器 |
| `sub_1F938` | `0x1F938` | `0xFB38` | 97 | 首／中／尾 byte，word 交換與位移，dummy latch read；每列來源加 BL、目的加 80 |
| `sub_1F999` | `0x1F999` | `0xFB99` | 23 | 由 DL 設 Enable Set/Reset，再選 Bit Mask；保存 AX／DX |
| `sub_1C7F4` | `0x1C7F4` | `0xC9F4` | 74 | 六次原始 hit map 與背板，接六次連續 glyph 來源 |
| `sub_1C673` | `0x1C673` | `0xC873` | 59 | 原始位置表與圖塊索引，裝置保存／恢復包住單一 glyph 重畫 |

`sub_1F888` 的 CX 拆成 byte 寬與高。寬度與 X 低三位先以 byte 運算相加，BL 由
寬度加 7 後取整 byte，BH 由寬度加對齊量得到目的整 byte 數。row helper 的 `dec bh`
後使用 JL／JZ，不能把它簡化成無號尺寸迴圈。左右 mask 在首／尾寫入前各送 VGA Bit Mask。

來源 word 由兩個 byte 交換後右移；dummy ES read 在每次 VRAM 寫入前載入 latch。
SI 每列加原始 BL，四平面來源連續。`CX=0x1018` 的六個 glyph 每次前進 `0xC0`，
`sub_1C7F4` 不重設 SI，最後到 `0x3D80`。來源步幅與遮罩邊界採原始 byte／word wrap。

CS:`0xD2E4..0xD2E9` 的 hit 順序與 CS:`0xD2EA..0xD2F5` 的位置表逐 byte 核對；
兩張表各自由對應 caller 消費，不以 hit code 當作 glyph 索引。圖庫第 1 段的實際大小與
呼叫來源仍須在執行前驗證，不用版本參數代替資產檢查。

## 2. 平台與依賴界線

沿用 [re/104](104-c-vga-blit-restoration.md#2-平台前提與獨立狀態) 的固定 VGA 契約，
原版 guest CPU 與 native C 各持有獨立 machine、RAM、四平面與 latch。
word 讀取按固定 dosgolem CPU 的分段契約：第二個 byte 的 offset 也以 16-bit 前進，
SI=FFFFh 時回到同段 0000h。位移計算仍為原始 16-bit。
未定義位移旗標只對齊該 CPU 模型，不宣稱實機逐週期相同。

`sub_1C673` 的底層裝置服務維持明示固定 FAR fixture，保存／恢復 caller 使用既有真實 C。
其他貼圖、hit map、外框與窗口 caller 接入實際 C 與 VGA，不能用平面 RAM 代替 VRAM。
畫面抽樣採 dosgolem 臥龍傳的 640×400 內容區，從 mode 12h 的第 40 列起取樣。

## 3. 驗證與產物

O0/O2 每版 4,954 組全部比較十四暫存器、FLAGS、stack、完整 1 MB RAM、四個
65536-byte plane、GC／seq／latch／index、port 順序與 90-byte 原始入口快照。
683 次 640×400 內容區像素相同，36 次原始 hit map 消費相同。原版與 C 都從各自狀態取結果。

| 群組 | 每種最佳化的組數 |
|---|---:|
| 位元對齊貼圖 | 2,496 |
| row helper | 1,344 |
| 原始 CH=0 的 256 列 | 56 |
| plane／Bit Mask 設定 | 512 |
| 實際六 glyph／八種對齊 | 96 |
| 單一重畫與三種 FAR 回應 | 216 |
| 六按鈕建立 | 6 |
| 六按鈕 hit 查詢 | 36 |
| 真實外框與 hit map | 192 |
| 合計 | **4,954** |

所有尺寸、來源、mask、plane pattern 與控制回應在執行前固定。首／中／尾 byte、八種
對齊、byte 寬相加回繞、有號 BH、SI／目的回繞、零高度、DF 與原始位置表均有收據。Read Map Select 在本輪固定 0。
ROW 主矩陣保持 write mode 0，避免 mask 與忽略 CPU 資料的 mode 1 關聯而掩蓋來源錯誤。
原始六按鈕程序六次連續消費 glyph，之後以真實 `sub_1E453` 查六個 slot 的 hit code。
既有 `sub_1895D` 外框透過真實 horizontal／vertical C 與 VGA 執行，沒有繪圖 no-op。

O0/O2 收據逐 byte 相同，SHA-256：`f3f26ed82f6f045896caf71c33125f28b54854ba74cf5e5e8ece090ed8346f92`。

| 刻意改錯 | 首次拒絕 |
|---|---:|
| 首 mask 不按 X 位移 | aligned 第 313 組 |
| 尾 mask 全設 FFh | aligned 第 1 組 |
| 刪 dummy latch read | row 第 1 組 |
| 每列來源多進一 byte | aligned 第 1 組 |
| 來源 word 不交換 byte | row 第 26 組 |
| 每平面重設 SI | aligned 第 1 組 |
| JL 只看 SF，忽略 OF | row 第 121 組 |
| 六按鈕每次重設 3900h | buttons 第 1 組 |
| 重畫位置表未乘二 | redraw 第 37 組 |
| 刪 CLD | aligned 第 2 組 |
| 外框列距改成 639 | window 第 17 組 |

C 來源在 [`aligned.c`](../../tools/c_recovery/aligned.c)、[header](../../tools/c_recovery/aligned.h)
與 [fixture](../../tools/c_recovery/aligned_fixture.h)，driver 與重跑入口為
[`c_recovery_aligned.go`](../../tools/c_recovery_aligned.go)、
[`c_recovery_aligned.sh`](../../tools/c_recovery_aligned.sh)。
驗證器 [`c_recovery_aligned_verify.py`](../../tools/c_recovery_aligned_verify.py) 核對原版、來源、
實際 binary 的 CGO flags、全部入口與負對照；版控摘要在
[`c-aligned-verification.json`](c-aligned-verification.json)。

Compiled source manifest SHA-256：`8392206341feb61f339c6efbb64adfc810704a0f6ef9c19caedf8ab96dabc7a9`。
全 C／header 與兩個 Go 檔雜湊進實際編譯參數，避免外部 include 的舊快取。
Go 1.26.7／GCC 12.2.0；dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`。
原版、工具與來源唯讀掛載，非 root、無網路、限資源 Docker，只本機輸出可寫。
C 機器碼匹配仍為 false，正式 Go／UI／存檔格式沒有修改。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_aligned_probe.py workplace/matching-decompilation/c-aligned/ida KI.EXE
tools/c_recovery_aligned.sh
```

原版圖庫只在容器內轉成輸入 bank，沒有新增圖像外送或公開原版資料。

初版 C 誤用了 Machine.Read16 的連續實體位址語意。SI=FFFFh、BH=1、CH=0 的
smoke 反例中，暫存器、RAM 與裝置暫存器相同，但 plane 不同。回查 CPU.read16
確認第二個 offset 會在同段回繞；修正只落在來源 word helper，沒有剔除跨界 case。
舊收據、C 來源與編譯摘要留在 `results/word-linear-counterexample-*`。
FAR stub 首次漏了 CS prefix，導致原版讀 DS 而 C 讀 CS；補原先宣告的段覆寫後，同矩陣全部重生。
此差異屬驗證工具，不改遊戲 caller，舊收據留在 `results/service-ds-counterexample-*`。

## 4. 未解範圍

後續正確素材 bank 證據見 [re/106](106-c-rectangle-bars-restoration.md)：`word_10D48` 命令圖示
改取 ICONGRF 第 3 段起點，`word_10D4A` 改取其段內 0x6C0。這份歷史矩陣使用第 1 段
輸入，保留為 raw caller 比較；正確素材的 54 個重畫與 18 個外框收據在新矩陣。

| 項目 | 邊界 |
|---|---|
| 正常玩家戰術流程 | 局部 caller 對照不代替自然戰鬥操作，Issue #22 保持原範圍 |
| 裝置服務的自然輸入與時序 | FAR fixture 只固定本輪輸入，不外推真實滑鼠時序 |
| 完整文字／字型服務 | 本輪是位元對齊圖塊，不完成其他字型 blitter |
| C 機器碼匹配與原作者工具鏈 | 尚未驗證 |
