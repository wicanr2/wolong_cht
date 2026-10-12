# 257：C 軌跡回放 B-1（數字函式吃真實開局輸入）

**狀態：READY。第一個吃真實遊玩輸入的 C 回放：spec/229 數字函式，
輸入來自 slice A 開局時間軸的 dosgolem 軌跡，不再是生成 fixture。**

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
   spec/229 四函式入口，記每一次呼叫的暫存器＋步數。
2. 取入口記憶體：對每支函式另跑一次 `runto:<入口>` 停在首次到達，
   `ipeek` 取該 slice 讀集合（harness 的 read addrs 即清單）。
3. 回放：把（暫存器，記憶體）灌進數字 harness 的 dual-run
   （在 harness 內同時跑原版 bytes 模擬與還原 C），比離開狀態。
   不需要 dosgolem 錄離開狀態——harness 內雙跑即是比對。
4. 驗收：軌跡上所有呼叫全綠；任一 mismatch 即停，先查「輸入沒灌全」
   （讀集合漏了），再查 C 語意。

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
