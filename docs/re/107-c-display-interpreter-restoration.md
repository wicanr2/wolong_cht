# 107：C 顯示清單分派與未分類繪圖控制流

**狀態：九原始函式與十七 code 入口 O0/O2 各 382 組完整 RAM／plane／ABI 相同；十個錯版拒絕，glyph raster 保留明示 fixture。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具：IDA Pro 9.4，database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-display/ida/input.exe.i64`
- DB SHA-256：`135353e5c2ec69804fc6b7bb6a627a3008b0a0a55e78b7fa5feb331f6f6edc4c`
- 探針：[`ida_display_probe.py`](../../tools/ida_display_probe.py)
- 契約：[`spec/227`](../spec/227-c-display-interpreter.md)
- 既有證據：[`re/48`](48-window-display-list.md)、[`re/106`](106-c-rectangle-bars-restoration.md)

## 1. 原始入口與資料界線

原始五函式 `sub_1030F`、`sub_10337`、`sub_1E993`、`sub_1E9A7`、`sub_1E9C1`
負責登記、場景選擇與 12-byte 記錄分派。Opcode 取首 byte，near target 從當下 DS
的 `0xEA0D + 2*(opcode-1)` 讀取，不改成固定 C 順序。
九個原始 handler 於 IDA `0x1EA1F..0x1EAD0`，共同座標展開於 `0x1EAD0..0x1EAE9`。

直線、底紋與立體框位於 `0x1EDFE..0x1F020`、`0x1F26E..0x1F45B`、
`0x1F465..0x1F4A2`。解碼只在一次性 DB 建 code item，先保存原名稱、資料行、flags 與 bytes；
原版與既有專案 DB 不改寫。底紋區的段前綴始於 `0x1F311`，不能以 `0x1F312` 切入口。
直線的 `0x1EE62..0x1EE64` 與底紋的 `0x1F45B..0x1F465` 是有寫入／讀取端的原始資料槽，
它們維持資料；不把邊界解碼失敗逐 byte 猜成 code。

水平／垂直原始常式 `sub_1FB29`／`sub_1FBA7` 與字串 `sub_1F6DC`、碼寬分類
`sub_1F878` 接入同一個 C 呼叫閉包。真正字型 raster callback 仍採明示 fixture，
字串迴圈、窄／寬碼來源前進、陰影與呼叫順序使用真實原始 C。

## 2. C 還原與組語方式

以原始固定雜湊指令產生 native C 來源：每條指令是具體 C 暫存器／記憶體／旗標運算，
分支是原始位址 label，call 接已有 C 或明示裝置 callback。C 不讀 opcode bytes 執行，
也不呼叫原版 CPU。產生來源保留原始位址、運算元、分級與出處，編譯後以原版 guest CPU 驗證。
這是目前 16-bit ABI 研究層的控制流還原；不宣稱取得原作者的原始 C 命名或 C 機器碼匹配。

新增解碼 code 可合併到既有組語補充索引。每個指令逐 byte 組譯，資料槽與未解區段保留
本機原版輸入。輸出的完整 EXE 仍需逐 byte 與完整 SHA-256 相同，任何指令不匹配即停止。

## 3. 驗證與未解範圍

後續真實 glyph 證據見 [re/108](108-c-glyph-raster-restoration.md)：原字庫／far 服務与 raster
已接入真正 C，這份歷史 1F75E fixture 收據保留，原始分派程式未換成另一套規則。

O0/O2 比原始登記、九 opcode、十個實際場景、signed／word 邊界、DS 分派、完整
RAM／四平面／latch／port／register／FLAGS／call snapshot。原始圖庫與十個場景保留完整輸入。
文字 glyph 的裝置 fixture 不代替完整字型服務或正常玩家 UI；除以零／商溢位由原始資料契約隔離。

| 項目 | 邊界 |
|---|---|
| 完整字型 callback | 本輪維持明示 fixture，不稱完整字形對拍 |
| 正常玩家操作與自然時序 | 局部場景與 caller 矩陣不足以完成 |
| 其他未分類區段與完整 C | 本輪只增加本閉包 |
| C 機器碼匹配與原作者工具鏈 | 尚未驗證 |

## 4. 原始 C 與組語驗證

O0/O2 各 382 組原版／C 比全部十四 register、FLAGS、完整 1 MB RAM、四個
65536-byte plane、GC／seq／latch／index、port 與 90-byte 入口快照；每組都比內容區像素。
十個實際原場景各兩個原始入口、三個原點、兩種 DF 共 120 組，包含九個 handler。
54 組獨立 opcode 另驗當下 DS 的 table，交換 02／03 table 的反例由原版實際選 target。
純文字有限輸入使用真實 byte／word 讀取與碼寬分類，raster 僅在原始 1F75E 邊界固定。

| 群組 | 每種最佳化組數 |
|---|---:|
| 登記／初始化 | 12 |
| 十個原場景 | 120 |
| 水平／垂直線 | 144 |
| Bresenham／雙色框 | 36 |
| 原始底紋 | 5 |
| 寬窄碼分類 | 6 |
| 字串／陰影 | 4 |
| 九 opcode／替代 DS | 54 |
| DS table 交換 | 1 |
| 合計 | **382** |

兩版 JSON 逐 byte 相同，SHA-256：`bc6c702a5255b4f1cca70d7b3aa16b5431a05ee1b984907a90f53e1892ab254a`。

| 刻意改錯 | 首次拒絕 |
|---|---:|
| 第一筆記錄步進 24 | scene 第 1 組 |
| 原點 X 讀相對 Y | opcode 第 7 組 |
| DS table 改讀 CS | table-swap 第 1 組 |
| 圖庫來源改成字串段 | opcode 第 49 組 |
| 窄碼不回退來源 byte | text 第 1 組 |
| 略過陰影 bit | text 第 3 組 |
| 底紋來源每列多 1 | texture 第 1 組 |
| 垂直線 pitch 改 79 | axis 第 73 組 |
| 雙色框寬改成 1 | line-box 第 19 組 |
| 登記表偏移多 1 | register 第 1 組 |

原有對稱 X／Y 無法區分錯來源，負對照明確改為不對稱輸入後重生完整矩陣，不放寬判準。
IDA MUL／DIV 的第一個 operand 是隱含 AL，原本空顯示字串產生無效 C；改取唯一明示 operand。
內部像素入口 1EFEF 與舊矩形 helper 1F140／1F17B 補入兩側原始 trace 集合。
RAM／plane／register 先前已相同，這些是來源／驗證工具修正，失敗紀錄留在 WORKLOG。

native C 來源：[display.c](../../tools/c_recovery/display.c)、[header](../../tools/c_recovery/display.h)、
[fixture](../../tools/c_recovery/display_fixture.h)、[生成 C](../../tools/c_recovery/display_generated.inc)。
[產生器](../../tools/c_recovery_display_generate.py) 使用固定 IDA 832 條指令輸出具體 C 運算與 label，
無 runtime opcode 取值；不將整個位址區段硬當一個函式，原始內部 near 入口也保留。
[driver](../../tools/c_recovery_display.go)、[重跑入口](../../tools/c_recovery_display.sh)、
[驗證器](../../tools/c_recovery_display_verify.py) 與 [公開收據](c-display-verification.json) 綁定全部來源。

新增 317 條／746 bytes 的組語補充已逐條組譯，資料槽保持資料。整合後為
24,693 條／56,138 code bytes；全 67,099-byte EXE 仍相同，三個負對照拒絕。
[補充索引](rectangle-handler-code.json) 保留原名稱／資料行／分類與 decoder 身分；
[組譯工具](../../tools/display_code_supplement.py) 不允許指令 bytes fallback。

Compiled source manifest：`workplace/matching-decompilation/c-display/results/c-source.sha256`，
SHA-256 `f0f0684b9df8a06dc5a6a3681e0a04456d3450b0cb56d179fdddade3ecd267b7`。
Go 1.26.7／GCC 12.2.0／IDA9.4；IF/TF=0、來源與原版唯讀、獨立 mutable RAM／VGA、
非 root、限資源 Docker。已證實只套用原始 raw 控制流與本矩陣；沒有完整 glyph 或自然玩家聲明。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_display_probe.py workplace/matching-decompilation/c-display/ida KI.EXE
tools/c_recovery_display.sh
```
