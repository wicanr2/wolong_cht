# 265：C 軌跡回放 B-7（世界顯示格函式吃真實攻城長程輸入）

**狀態：CONFORMED。38,055 次呼叫去重 12,525 案，O2／O0 雙跑全同、
原版／C 終態雜湊全同（[`c-mapcells-trace-verification.json`](../re/c-mapcells-trace-verification.json)）；
mapcells 矩陣 smoke 重跑不回歸。**

- 日期：2026-10-14。
- 出處：[`spec/264`](264-matching-siege-longpath.md)（攻城長程時間軸）、
  [`spec/263`](263-c-trace-replay-b6.md)（全暫存器回放法）、
  [`spec/234`](234-c-world-map-cells.md)（世界顯示格 C 還原）。

## 普查（本輪實測，spec/264 時間軸加 WATCH 10 入口）

| 入口 | 呼叫 | 全暫存器去重 |
|---|---:|---:|
| `1D4C7` | 9,586 | 1,009 |
| `1D796` | 7,848 | 5,090 |
| `1D782` | 5,036 | 4,949 |
| `1D615` | 3,571 | 484 |
| `1D66A` | 3,571 | 170 |
| `1D804` | 2,813 | 486 |
| `1D7E7` | 2,813 | 189 |
| `1D76B` | 2,813 | 144 |
| `1D46A` | 2 | 2 |
| `1D483` | 2 | 2 |

同軸另記：攻城鏈 `14ADE`／`14C72`／`11B5A` 各 1 次——那是
engagement 整場 harness 的範圍，本片不處理（另開 spec）。

## 輸入完備性與佈置（沿用 spec/234 的 harness 固定區）

- 十支的記憶體讀集是 CS 資料字 `D84A/D84C/D84E/D850/D854/D856`
  指到的四張表（`MMAP.MDL`／`MMAP.MCH`／旗標表／world raw）與
  VGA plane；harness `initial` 已把**真實素材**佈置在固定段
  （3000／3800／2500／5000）且 CS 資料字兩側同值，不從軌跡佈置。
- FLAGS 用軌跡原值回放；段暫存器除 CS（映射到 harness CS）外
  保留軌跡值兩側同值；SP 用軌跡值，`SS:SP` 兩側同置返回哨兵。
- 旗標表（2500 段）用 harness 既有的結構化佈置（同 matrix 的
  flags/render 組）：pattern 原始內容會讓 `1D66A` 的髒格迴圈
  不收斂（實測跑滿 400 萬步上限），結構化表是這支 harness
  已驗證的表形狀；宣稱因此只到「真實暫存器輸入＋harness 佈置
  的表」，與 B-6 的 map 內容界線相同。
- ⚠ `1D7E7` 的有效輸入要求 DF=0（spec/234：原 renderer 呼叫前
  先 `CLD`；DF=1 會反向寫進自己的後續指令）。普查實測 2,813 次
  呼叫的 FLAGS 全為 `7006`（DF=0），回放按軌跡原值即滿足前提；
  若未來軌跡出現 DF=1 的呼叫，該案要排除並記錄，不硬跑。

## 做法

1. 軌跡：spec/264 時間軸加 WATCH 上表 10 入口（已錄，存
   `workplace/matching-decompilation/c-mapcells/trace/watch.log`）。
2. 回放：`c_recovery_mapcells.go` 加 `--trace` 模式（比照 rect）：
   全暫存器＋FLAGS 去重後逐案 dual-run，比離開暫存器、完整 RAM、
   四 plane、device、port 與入口 trace。
3. 驗收：去重案 O2/O0 雙跑全同、原版/C 終態雜湊全同；
   mapcells 既有矩陣 smoke 重跑不回歸。mismatch 即停，
   先查佈置再查 C 語意。

## 實作與驗收

- 產出：mapcells harness trace 模式、`tools/c_recovery_mapcells.sh`
  的 `trace` 模式、回放收據 `docs/re/c-mapcells-trace-verification.json`、
  本 spec 升 CONFORMED。
- 原版素材唯讀；軌跡輸出進 `workplace/`，不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 攻城鏈（14ADE／14C72／11B5A）軌跡回放 | engagement 整場 harness 的範圍，另開 spec |
| 素材內容 | harness 固定段的真實 MMAP 素材；相機位置與髒格狀態由暫存器輸入帶入，非原版該時刻的完整顯示狀態 |
| 其他 slice | 下一片另開 spec |
