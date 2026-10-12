# 256：matching decomp 玩家長程驗證

**狀態：CONFORMED。Slice A 全綠：重建 EXE 從零跑到 196/4/8，8 截圖＋時鐘＋130,849,387 道指令數與原版全同（[`playtest/122`](../playtest/122-rebuilt-longpath-smoke.md)）。C 軌跡回放維持設計，不在本 spec 驗收。**

- 日期：2026-10-12。
- 出處：CONTEXT 各輪「完整玩家長程仍待」；[`re/90`](../re/90-assembly-reconstruction.md)
  §7（冷重建整檔 SHA 同）；[`spec/131`](131-dosgolem-oracle.md)（oracle 走 dosgolem）；
  CLAUDE.md §4.0（密碼頁留白按確定即過）。
- 輸入：原版 `KI.EXE`（`fffeba98…3868`）與冷重建 `KI.EXE`
  （`workplace/matching-decompilation/assembly/code-record/KI.EXE`，
  `cmp` 全同）；執行器 `tools/dosgolem.sh`。

## 為什麼要跑（byte 全同還跑什麼）

整檔 SHA 同證明 bytes 相等，不證明「重建產物在模擬器下跑得起來」：
載入路徑、檔名大小寫、工作目錄、dosgolem 的 EXE 解析任一環出錯，
SHA 都是綠的、跑起來是死的。Slice A 是 smoke，不是邏輯驗證——
邏輯相等由 SHA 背書，**跑得起來由本輪背書**。

## Slice A：重建 EXE 從零開局對拍（本輪）

1. 造 `workplace/longpath-rebuilt/`：複製 `orig/dosv` 69 檔，
   以冷重建 `KI.EXE` 蓋掉其中一支（`workplace/` 不入版控，原檔唯讀不动）。
2. 原版側先跑一次探索：`wait；shot:boot；hotspots`，從熱區圖找出密碼頁
   「確定」鈕座標，再把時間軸延伸到開場結束、大地圖出現
   （`clock` 確認遊戲日期在走）。
3. 同一時間軸跑兩次：`WOLONG_DOSGOLEM_GAMEDIR=orig/dosv`（對照）與
   `WOLONG_DOSGOLEM_GAMEDIR=longpath-rebuilt`（受測），取樣點寫成
   `clock`＋`shot`（即時制：日期＋畫面，不寫秒數）。
4. 驗收：同名截圖 SHA-256 全同、`clock` 輸出全同；任一不同即停，
   先懷疑載入／路徑，不懷疑邏輯（邏輯 SHA 已背書）。

## Slice B：還原 C 軌跡回放（設計，未來輪）

C 還原是逐函式 fixture，從未吃過真實遊玩序列。做法：
dosgolem `eventwatch`＋`WOLONG_DOSGOLEM_WATCH` 沿真實開局錄一串
（函式入口，暫存器，關鍵記憶體）軌跡，還原 C 依序吃同一串輸入，
在函式邊界比暫存器與記憶體。開局只到大地圖出現，不碰 AI 長程。
本輪不實作，實作前另開 spec 切片。

## 實作與驗收

- 產出：本 spec、playtest 收據（`docs/playtest/122-*.md`）、
  `workplace/longpath-rebuilt/`（不入版控）。
- 驗收：Slice A 四步全綠即 CONFORMED；C 軌跡回放不列入本 spec 驗收。
- 原版素材唯讀，不入 Git；`longpath-rebuilt/` 內含原版資料，只活在本地。

## 未解範圍

| 項目 | 限制 |
|---|---|
| C 軌跡回放實作 | 只留 Slice B 設計，實作另開 spec。 |
| AI 長程／90 日對拍 | remake 側 Issue #10 的範圍，不在 matching 線。 |
| PC-98 重建實跑 | matching 基線只含 DOS/V，不含 PC-98。 |
