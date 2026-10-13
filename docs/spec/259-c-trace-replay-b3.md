# 259：C 軌跡回放 B-3（glyph 函式吃真實開局輸入）

**狀態：CONFORMED。468 次呼叫去重 451 組，O2／O0 雙跑全同、原版／C 終態雜湊全同（[`c-glyph-trace-verification.json`](../re/c-glyph-trace-verification.json)）。dosgolem 已補 `regs` 動詞與 WATCH 全暫存器行（BP／SP／DS／ES／SS／FLAGS）＋ `-watchmem` 傾印。**

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
   106FD×8、1F720×1（`re/29` 字型向量設定，只跑一次），終點指令數
   130,849,387 與 playtest/122 全同。
2. 逐去重呼叫：WATCH 擴充行一次給全暫存器（BP／SP／DS／ES／SS／FLAGS，
   `regs` 動詞同格式供定點查核），`-watchmem ES:BP:64,CS:#F780:48,DS:SI:64`
   同行傾印三個視窗。去重鍵是四者的完整內容。
   DS 由入口 AX 覆寫（`m->ds＝m->ax`），內容走 harness RAM，不必抓。
3. 回放：glyph harness `--trace` 模式（比照 talk）：全暫存器＋三視窗佈置進
   dual-run，device 初態兩側同構。FLAGS 用軌跡原值（不掃 CF／DF；
   B-1／B-2 掃 CF 只因舊 WATCH 行沒有 FLAGS）。
   ⚠ CS 視窗只佈置 F7A1..F7A3（模式字＋窄字旗標）：鄰接的字型 far 指標
   F78A–F793 是 harness patch 到 0080:0410/0414 stub 的呼叫操作元，
   佈置真值會把原版 CPU 送去真 TSR 位址（實測跑滿步數上限）。
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
| dosgolem `regs` 動詞 | 已落地（dosgolem `apps/wolong/cmd/shot`，含 `formatRegs`／`parseWatchMem` 單元測試）；本專案 `tools/dosgolem.sh` 以 `WOLONG_DOSGOLEM_WATCHMEM` 透傳。 |
| 字模窗長度 | 64B 已足：451 去重案全綠，未加長。 |
| glyph 全量 verify | `c_recovery_glyph_verify.py` 卡在無關的 rectangle 補充檔計數（實得 339／808，期望 338／805，HEAD 既有），矩陣 252 案本身通過；待另案修。 |
| 其他 slice | 下一片另開 spec。 |
