# 261：C 軌跡回放 B-5（18853／189A4 吃真實指令流程輸入）

**狀態：CONFORMED。指令流程（財政）錄得 18853×1、189A4×1，去重 2 案
O2／O0 雙跑全同、原版／C 終態雜湊全同（[`c-talk3-trace-verification.json`](../re/c-talk3-trace-verification.json)）；B-2／B-4 trace 重跑不回歸。**

- 日期：2026-10-13。
- 出處：[`spec/258`](258-c-trace-replay-b2.md)（未解表）、
  [`spec/260`](260-c-trace-replay-b4.md)（全暫存器回放法）、
  [`playtest/122`](../playtest/122-rebuilt-longpath-smoke.md)（開局時間軸）。

## 為什麼不在 slice A（實測）

WATCH 盯 18853／189A4 跑完 slice A：零呼叫。兩支的呼叫端都在清單／
指令流程（`formation`／`list`／`personnel`／`march` 的狀態訊息），
開局走到大地圖還沒開任何清單。改走：開局 → 關事件訊息（右鍵×4）→
指令列 → 財政，一條玩家真的會走的路，實測 18853×1、189A4×1
（CX=0010，走 CX≠FFFF 的第二訊息分支）。

## 輸入完備性證明（`talk_generated.inc` 靜讀）

| 函式 | 入口讀取 | 證明 |
|---|---|---|
| 18853 | CX（FFFF 分支）、DS:SI（1075B 的訊息參數表）；AL／BX／DX 進門即覆寫 | AX／BX／DX 只 push/pop；第一個旗標讀是自設 `xor al,al` 之後的 `cmp cx,0FFFFh`，入口 FLAGS 不消費 |
| 189A4 | AX（訊息號）、BX／DX／CX（窗座標，函式內 ×2） | 暫存器全 push/pop 或覆寫後使用；無入口旗標讀、無入口記憶體讀（窗原語 1D5D4／189DE 在 dual-run 內兩側同跑） |

## 做法

1. 錄軌跡：上述指令流程時間軸加 WATCH 18853／189A4，
   `-watchmem DS:SI:64,DS:DX:64,CS:#F780:16`。
2. 回放：talk harness `--trace` 的全暫存器 target 表加 18853／189A4
   （與 B-4 共用機制與計數欄）；FLAGS 用軌跡原值。
3. 驗收：去重呼叫全綠；mismatch 即停，先查佈置再查 C 語意。

## 實作與驗收

- 產出：回放收據 `docs/re/c-talk3-trace-verification.json`、
  本 spec 升 CONFORMED。
- 驗收：呼叫數、去重數、O2／O0 雙跑全同、原版／C 終態雜湊全同；
  B-2／B-4 的 trace 模式重跑不回歸。
- 原版素材唯讀；軌跡輸出進 `workplace/`，不入版控。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 18853 的 CX=FFFF 分支 | 本流程未觸發；矩陣案已覆蓋，軌跡回放未覆蓋。 |
| 其他 slice | 下一片另開 spec。 |
