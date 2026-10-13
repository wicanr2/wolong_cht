# 260：C 軌跡回放 B-4（1084A／1F4DF 吃真實開局輸入）

**狀態：CONFORMED。26 次呼叫去重 26 組，與 B-2 的 58 案同跑共 84 案，
O2／O0 雙跑全同、原版／C 終態雜湊全同（[`c-talk2-trace-verification.json`](../re/c-talk2-trace-verification.json)）；talk 全量矩陣 1277 案與 10 突變重跑通過。**

- 日期：2026-10-13。
- 出處：[`spec/258`](258-c-trace-replay-b2.md)（前例與未解表）、
  [`spec/259`](259-c-trace-replay-b3.md)（regs／watchmem 擷取法）、
  [`spec/230`](230-c-talk-rendering.md)（TALK C 還原）、
  [`playtest/122`](../playtest/122-rebuilt-longpath-smoke.md)（時間軸）。

## 輸入完備性證明（`talk_generated.inc` 靜讀）

| 函式 | 入口讀取 | 證明 |
|---|---|---|
| 1084A | AX（AH 字級）、BX／DX（座標）、CX、DS:SI 字串、CS 跳表 `cs:[si+8A4h]`（靜態） | BP／ES 進門即 push、出口 pop，全程不讀其值（BP 立即覆寫成 10h／8）；第一個旗標讀在 `and ah,ah` 之後，入口 FLAGS 不消費。字串以 0 結尾計長 |
| 1F4DF | AX／CX（lseek 位移）、BX（讀入段）、DS:DX 檔名、SI（讀入位移）、DI（長度）、DOS 檔案狀態 | BP 進門 push、建完自有 frame（`bp=sp`）才讀寫 `[bp+…]`，入口 BP 不消費；ES 全程不碰；第一個旗標讀是 int 21h 之後的 `jb`，由中斷結果設定 |

## 做法

1. 錄軌跡：slice A 時間軸加 WATCH 六函式（1084A／1F4DF 加上 B-2 四支做
   census 對照），`-watchmem DS:SI:256,DS:DX:64` 同行傾印。
   18853／189A4 預期維持零呼叫（本片不處理）。
2. 回放：talk harness `--trace` 擴充兩個 target：全暫存器佈置
   （執行期 CS 段對回 harness CS）、字串／檔名視窗佈置進 flat RAM、
   FLAGS 用軌跡原值（不掃 CF）。DOS 檔案由 harness 兩側同構的
   `dos.New(/orig)` 提供。
3. 驗收：去重呼叫全綠；mismatch 即停，先查佈置（視窗長度、段對照），
   再查 C 語意。

## 實作與驗收

- 產出：`c_recovery_talk.go`／`.sh` 的 trace 目標擴充、回放收據
  `docs/re/c-talk2-trace-verification.json`、本 spec 升 CONFORMED。
- 驗收：呼叫數、去重數、O2／O0 雙跑全同、原版／C 終態雜湊全同；
  B-2 既有 58 案不回歸（同一支 harness 的 trace 模式重跑）。
- 原版素材唯讀；軌跡輸出進 `workplace/`，不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 18853／189A4 | 本 timeline 零呼叫；需單挑等特殊路徑軌跡，另開 spec。 |
| 其他 slice | 下一片另開 spec。 |
