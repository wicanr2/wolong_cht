# 116：場景主入口與鏡頭／小地圖退出重畫 C

**狀態：CONFORMED。9 支新函式與既有世界退出常式的完整局部 C 呼叫鏈通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 固定 IDA DB SHA-256：`4bf788aa04d543599828a061e4d177c8405950fc7d02c4690c1c2b5fdedf0204`
- 工具／位址：IDA Pro 9.4 database linear；file offset 為 linear 減 `0x10000` 加 512
- 入口：[ida_scene_resume_probe.py](../../tools/ida_scene_resume_probe.py)
- 前置：[re/113](113-c-resource-cue-restoration.md)、[re/114](114-c-world-map-cells-restoration.md)、[re/115](115-c-world-overlay-restoration.md)
- 契約：[spec/236](../spec/236-c-scene-resume.md)

`sub_13B08` 先停音、載入朝堂插圖與第 6 曲，再依原 TALK index 顯示三句對話，
還原 MMAP、世界畫面、滑鼠與場景音樂。`sub_11D46` 已有 C，本輪將其依賴接完整。
`sub_11F30` 的 Y 減 2 借位會直接尾跳 `sub_11F5A`，保留原堆疊，不改成 near call。
小地圖視野框依舊／新鏡頭判斷 dirty，舊區從 ICONGRF 段 2 還原，
新框從段 3 的 `0x08F0` 經 latch／bit shift 合成。

上述控制流由固定原指令證實；局部 C 行為由下列收據確認，未回填猜測性名稱。

## 來源與驗證

九個新函式共 578 bytes／255 原指令，既有 `sub_11D46` 沿用唯一原 C 實作。
[resume.c](../../tools/c_recovery/resume.c)、[生成內容](../../tools/c_recovery/resume_generated.inc)與
[收據](c-resume-verification.json)保存原名稱、位址、source hash 與實際編譯 flags。
重跑入口 [c_recovery_resume.sh](../../tools/c_recovery_resume.sh)。

| 矩陣 | O0／O2 每版組數 |
|---|---:|
| 四 plane 背景／位元對齊視野框、DF | 96 |
| 背景 X 邊界 clamp | 3 |
| 鏡頭 Y／借位尾跳／hidden | 10 |
| 鏡頭記憶／髒格 helper | 18 |
| 小地圖舊／新位置／哨兵 | 4 |
| 四劇本 × 小地圖 × 三相機的世界退出重畫 | 24 |
| 原兩個主入口分支 × 四劇本 × 小地圖 | 16 |

O0／O2 各 171 完整 RAM／plane／FLAGS／暫存器／SS frame／IN／OUT／API／Mouse／sound 狀態相同，
12 個 O0 錯版全拒絕，C 乾淨重生相同。16 個場景完整走三句、等待、停播曲、
IVENT→MMAP 恢復、原 world producer／renderer 與滑鼠，沒有以 no-op 代替相依函式。
每版 32 次君主上框／16 次軍師下框、816 全形 glyph、缺字 0；原文第一行字元給出
432 glyph 的保守下限，驗證沒有以 guessed workload 放寬為零。

輸入取自 `sub_16909` 的 AL=0／TALK `0x182` 與 `sub_1699E` 的 AL=1／TALK `0x18C`。
主入口沿用世界 DS，SI 是原人物記錄；最初誤給程式 DS 造成原版性格與 TALK index 讀錯，
原 trace／RAM 無程式碼覆寫，按原呼叫前提修正後同預算重跑。

每次執行前核對實際 MDL／MCH／MAP／ICON bytes，避免測試準備把素材段蓋成別種資料。
TALK／BGM／hotspot／world／字庫／stack 以明示獨立段布置。原 far loader 重定位與
跨函式 tail 的堆疊保留；受控 press／counter 不代表原實機 wall-clock 或自然事件。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然事件／長程玩家路徑 | 局部場景主入口不代證完整自然事件與通關 |
| 原作者 C 工具鏈／機器碼 | 尚未確認 |
