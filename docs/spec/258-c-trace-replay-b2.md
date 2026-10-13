# 258：C 軌跡回放 B-2（TALK 函式吃真實開局輸入）

**狀態：CONFORMED。35 次呼叫去重 29 組 × CF 兩態＝58 案，
O2／O0 雙跑全同、原版／C 終態雜湊全同（[`c-talk-trace-verification.json`](../re/c-talk-trace-verification.json)）。106F9 兩呼叫另佈置 DS:SI 各 256 bytes 真實內容（實機 DS＝CS，故是 code 段 bytes，雙側同）。**

- 日期：2026-10-12。
- 出處：[`spec/257`](257-c-trace-replay-b1.md)（方法）、[`spec/230`](230-c-talk-rendering.md)
  （TALK C 還原）、[`playtest/122`](../playtest/122-rebuilt-longpath-smoke.md)（時間軸）。
- 為什麼跳過 glyph 當 B-2：`1F7A4` 首行即讀 ZF、體內讀 BP／ES，
  WATCH 給不出來，需 dosgolem `regs` 動詞（跨 repo，待使用者決定）；
  TALK 四支全是 GPR 輸入，零工具鏈改動可做。

## 輸入完備性證明（`talk_generated.inc` 靜讀）

| 函式 | 呼叫數 | 證明 |
|---|---|---|
| 106F9 | 2 | 5 行，零狀態讀取，轉呼 1084A（內部狀態雙側同構） |
| 1075B | 2 | DS 存即覆寫（CS 靜態）；CF 讀前必有 `cmp`／`math`（L10／L14 先於 L12／L32）；零 BP／SP／SS／ES |
| 107D2 | 7 | DS 由入口 BX 覆寫；ZF／CF 讀前必有 `logic`／`cmp`（L22／L30 先於 L26／L32）；零 BP／SP／SS／ES |
| 1E38C | 24 | DS 只存取不讀；CF 讀（L32）前無 flagset → **入口 CF 有效，seed 2／3 掃**；零 BP／SP／SS／ES／DF |

未覆蓋：1084A（讀 BP／ES）、1F4DF（BP frame 讀 SS）、18853／189A4
（本 timeline 零呼叫）。Glyph `1F7A4` 同列待 `regs` 動詞。

## 做法

1. 錄軌跡：slice A 時間軸加 WATCH 八函式（上表四支＋1084A／1F4DF／18853／189A4
   做 census）。實測：106F9×2、1075B×2、107D2×7、1084A×2、1E38C×24、
   1F4DF×24、18853／189A4 零呼叫。
2. 回放：WATCH 六 GPR 去重後灌進 TALK harness dual-run（`--trace` 模式，
   比照 `c_recovery_numbers.go`），CF 各掃 0／1，device 初態兩側同構。
3. 驗收：去重呼叫 × CF 兩態全綠；mismatch 即停，先查完備性證明再查 C 語意。

## 實作與驗收

- 產出：`c_recovery_talk.go`／`.sh` 的 trace 模式、回放收據
  `docs/re/c-talk-trace-verification.json`、本 spec 升 CONFORMED。
- 驗收：呼叫數、去重數、CF 掃、O2／O0 雙跑全同、原版／C 終態雜湊全同；
  不宣稱 1084A／1F4DF／glyph。
- 原版素材唯讀；軌跡輸出進 `workplace/`，不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 1084A／1F4DF／glyph 1F7A4 回放 | 需入口 BP／ES／段暫存器，待 dosgolem `regs` 動詞（跨 repo）。 |
| 18853／189A4 | 本 timeline 零呼叫；需單挑等特殊路徑軌跡。 |
| 其他 slice | 下一片另開 spec。 |
