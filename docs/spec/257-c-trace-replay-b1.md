# 257：C 軌跡回放 B-1（數字函式吃真實開局輸入）

**狀態：CONFORMED。522 次呼叫去重 233 組 × CF 兩態＝466 案，
O2／O0 雙跑全同、原版／C 終態雜湊全同（[`c-numbers-trace-verification.json`](../re/c-numbers-trace-verification.json)）。**

- 日期：2026-10-12。
- 出處：[`spec/256`](256-matching-longpath-verification.md) slice B 設計；
  [`spec/229`](229-c-number-raster.md)（數字函式 C 還原）；
  [`playtest/122`](../playtest/122-rebuilt-longpath-smoke.md)（開局時間軸）；
  dosgolem `WATCH`（入口暫存器）＋`runto`＋`ipeek`（入口記憶體）能力
  已讀碼確認（`apps/wolong/cmd/shot/main.go`：`installWatches`、`runto`）。
- 為什麼選數字函式：開局日期「196年4月1日」渲染必經它；
  輸入是暫存器＋RAM（時鐘記憶體），不碰滑鼠／IRQ；
  slice 內已有 VRAM／buffer 處理，回放沿用。

## 做法

1. 錄軌跡：slice A 時間軸（開機→4/8）加 `WOLONG_DOSGOLEM_WATCH` 盯
   三函式入口，記每一次呼叫的暫存器。實測：1062F 522 次、
   1069A 862 次（全來自 10669，內部呼叫）、106DE 372 次。
2. 只回放頂層 1062F（內部呼叫在 dual-run 內兩側同跑，不必單獨錄）。
   輸入完備性證明（`numbers_generated.inc` 靜讀）：入口 BP 只存不讀、
   CX 進門即覆寫、DS／ES 進門即覆寫、零 stack／日期 RAM 讀取；
   CF 用 seed 2／3 掃（既有矩陣已證其餘 flags-in 無關）。
   **不必 `runto`＋`ipeek`，不必動 dosgolem。**
3. 回放：WATCH 的六個 GPR 去重後灌進數字 harness dual-run
   （`--trace` 模式，`c_recovery_numbers.go`），CF 各掃 0／1，
   device 初態兩側同構。比離開暫存器＋RAM＋plane＋device＋port＋trace。
4. 驗收：去重呼叫 × CF 兩態全綠；任一 mismatch 即停，先查輸入完備性
   證明，再查 C 語意。

## 實作與驗收

- 產出：軌跡擷取腳本（`tools/`）、回放收據（`docs/re/c-number-trace-*.json`）、
  本 spec 升 CONFORMED。
- 驗收：呼叫數、輸入覆蓋（讀集合全灌）、雙跑全同，三者白紙黑字；
  不宣稱超出數字函式的覆蓋。
- 原版素材唯讀；軌跡輸出進 `/tmp` 或 `workplace/`，不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 其他 slice 的軌跡回放 | B-1 只做數字函式；下一 slice 另開 spec。 |
| WATCH 無離開狀態 | 靠 harness 內雙跑補，不改 dosgolem。 |
| 互動函式（選單／等待） | 輸入是硬體事件，不在 B-1。 |
