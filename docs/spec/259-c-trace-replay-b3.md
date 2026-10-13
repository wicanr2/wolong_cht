# 259：C 軌跡回放 B-3（glyph 函式吃真實開局輸入）

**狀態：READY。B-1／B-2 方法套到 glyph slice：1F7A4／106F5／106FD／1F720。
實作待 dosgolem `regs` 動詞（使用者套用中）；spec 先行。**

- 日期：2026-10-12。
- 出處：[`spec/257`](257-c-trace-replay-b1.md)（方法）、[`spec/258`](258-c-trace-replay-b2.md)
  （TALK 前例＋訊息佈置）、[`spec/228`](228-c-glyph-raster.md)（glyph C 還原）、
  [`playtest/122`](../playtest/122-rebuilt-longpath-smoke.md)（時間軸）。

## 為什麼 glyph 需要 regs 動詞（TALK 不需要）

`1F7A4` 首行即讀 ZF（`test cs:…；if ZF goto`），體內 `SI＝BP` 後讀 `ES:[SI]`：
入口 BP、ES、ZF 全是有效輸入，WATCH 只給六個 GPR。
TALK 四支經靜讀證明不讀 BP／SP／SS／ES／DF（B-2），glyph 證明失敗，
故 B-3 的 capture 必須有全暫存器。

## 做法

1. 錄軌跡：slice A 時間軸加 WATCH 四函式。實測：1F7A4×452、106F5×7、
   106FD×8、1F720×1（`re/29` 字型向量設定，只跑一次）。
2. 逐去重呼叫：`runto:<入口>,N`＋`regs` 取全暫存器（含 BP／ES／FLAGS），
   `peek:<ES>:<SI>:64` 取字模 bytes（`peek` 三欄位形吃指定段）。
   DS 由入口 AX 覆寫（`m->ds＝m->ax`），內容走 harness RAM，不必抓。
3. 回放：glyph harness `--trace` 模式（比照 talk）：全暫存器＋字模佈置進
   dual-run，device 初態兩側同構。1F720 另掃 DF（矩陣前例）。
4. 驗收：去重呼叫全綠；mismatch 即停，先查佈置（字模窗長度、ES 段），
   再查 C 語意。

## 實作與驗收

- 產出：`c_recovery_glyph.go`／`.sh` 的 trace 模式、回放收據
  `docs/re/c-glyph-trace-verification.json`、本 spec 升 CONFORMED。
- 驗收：呼叫數、去重數、O2／O0 雙跑全同、原版／C 終態雜湊全同；
  字模佈置的位址與長度白紙黑字。
- 原版素材唯讀；軌跡輸出進 `workplace/`，不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| dosgolem `regs` 動詞 | 跨 repo，使用者套用中；未落地前 B-3 不開工。 |
| 字模窗長度 | 先抓 64B；mismatch 再加長，不預設。 |
| 其他 slice | 下一片另開 spec。 |
