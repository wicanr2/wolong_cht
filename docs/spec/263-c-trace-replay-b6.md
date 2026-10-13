# 263：C 軌跡回放 B-6（熱區／外框／填色函式吃真實互動輸入）

**狀態：CONFORMED。13 個目標共 3,723 次呼叫去重 552 案，O2／O0 雙跑
全同、原版／C 終態雜湊全同（[`c-rect-trace-verification.json`](../re/c-rect-trace-verification.json)）；
rect 矩陣 smoke 重跑不回歸。普查時「呼叫」計數 3,725 含 2 行教訓
文字的假命中，parser 實得 3,723。**

- 日期：2026-10-13。
- 出處：[`spec/262`](262-matching-interactive-longpath.md)（互動長程時間軸）、
  [`spec/261`](261-c-trace-replay-b5.md)（全暫存器回放法）、
  [`spec/222`](222-c-hotspot-map.md)／[`spec/225`](225-c-rectangle-bars.md)
  （熱區與矩形 C 還原）、`hotspot.c`／`rect.c` 靜讀（本片輸入完備性證明）。

## 普查（本輪實測，spec/262 時間軸加 WATCH 26 入口）

| 入口 | 呼叫 | 六暫存器去重 | 入口讀集（靜讀） |
|---|---:|---:|---|
| `1E453` 熱區 pixel query | 2,763 | 25 | CX／DX；`CS:E479/E47B`（harness 固定值）；計算位址上一 byte |
| `10CAC` VGA mode set | 775 | 354 | 無（AX／DX 只 push/pop）；port 寫出 |
| `1E3D7` 熱區矩形登記 | 52 | 50 | AX／BX／CX／DX；`CS:E479/E47B`；矩形內 map byte（重疊 CF）|
| `1F1A3` 矩形填色 | 35 | 31 | AX／BX／DX／SI／DI；VGA plane（DS/ES 函式內設 A0C8）|
| `10C60`／`10C77` 水平／垂直邊 | 20／20 | 19／18 | BX／CX／DX；內部呼叫 `1F9B0` blit（DS 由 caller 設定）|
| `1D5D4` tile flags | 14 | 13 | AX／BX／CX／DX；`CS:D84E` 段字；該段表 byte |
| `189DE` 16px→8px wrapper | 14 | 12 | AX／BX／DX；內部呼叫 `1E3D7` |
| `1895D` 視窗 wrapper | 14 | 14 | AX／BX／CX／DX；內部呼叫 `1D5D4`／`10C14`／`189DE` |
| `10C14` 外框幾何 | 10 | 10 | BX／CX／DX；`CS:0D4A`；內部呼叫 `10C60`／`10C77` |
| `1E3C0` 熱區初始化 | 2 | 2 | ES／DI（見下方佈置替代）；清 4,000 byte、寫 `CS:E479/E47B` |
| `10BCD` 視窗 caller | 2 | 2 | 全暫存器；內部呼叫 `10C14`／`1F1A3` |
| `10AD9` span 繪製 | 2 | 2 | 全暫存器；內部呼叫 `10CAC`；VGA plane |

零呼叫（本片不處理，僅矩陣覆蓋）：`1F020`、`1F140`、`1F17B`、`1E41B`、
`1C61F`、`1C6BF`、`1C6AE`、`1C6F6`、`1C74C`、`1C775`、`1C78E`。

## 輸入完備性與佈置

- FLAGS 一律用軌跡原值回放（不掃 CF/DF），段暫存器除 CS 外保留軌跡值
  兩側同值；runtime CS 映射到 harness CS（B-4/B-5 前例）。
- 函式讀的 CS 資料字（`E479/E47B`、`D84E`、`0D4A`）在 harness
  `initialRAM` 兩側同值固定（`6000:0`、`5800`、`226Ch`），不從軌跡佈置。
- map（`6000:0` 4,000 byte）、tile 表與 VGA plane 用 harness 兩側同構的
  模式內容；trace 案加 `realAsset`，`ICONGRF` 真實素材載入 `22000/30000`
  兩側，使 `1F9B0` blit 與 tile flags 讀到真實 byte。
- ⚠ `1E3C0` 的 ES/DI 例外：改佈置 harness 的 `6000:0` 兩側同值
  （B-3 的 far 指標替代前例）。軌跡 ES/DI 指向 caller 資料段，
  照原值清 4,000 byte 會蓋掉兩側同構但無意義的區段，且後續案各自
  重置 RAM，不影響其他目標的覆蓋。
- SP 用軌跡值，`SS:SP` 兩側同置返回哨兵 `0x100`（harness 既有機制）。

## 做法

1. 錄軌跡：spec/262 時間軸加 WATCH 上表 13 入口（全暫存器行，
   不需 watchmem：入口記憶體讀集全在 harness 固定區）。
2. 回放：`c_recovery_rect.go` 加 `--trace` 模式（比照 talk/glyph）：
   全暫存器＋FLAGS 去重後逐案 dual-run，比離開暫存器、完整 RAM、
   四 plane、device、port 與入口 trace。
3. 驗收：去重案 O2/O0 雙跑全同、原版/C 終態雜湊全同；
   rect 既有矩陣模式重跑不回歸。mismatch 即停，先查佈置再查 C 語意。

## 實作與驗收

- 產出：rect harness trace 模式、`tools/c_recovery_rect.sh` 的
  `trace` 模式、回放收據 `docs/re/c-rect-trace-verification.json`、
  本 spec 升 CONFORMED。
- 原版素材唯讀；軌跡輸出進 `workplace/matching-decompilation/c-rect/trace/`，
  不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 零呼叫的 11 個同族入口 | 本時間軸未觸發；僅 spec/222／225 矩陣覆蓋，另開時間軸才處理 |
| map／plane 內容 | harness 模式內容＋真實 ICONGRF，非原版該時刻的畫面內容；宣稱只到「真實暫存器輸入下原版/C 一致」 |
| 其他 slice | 下一片另開 spec |
