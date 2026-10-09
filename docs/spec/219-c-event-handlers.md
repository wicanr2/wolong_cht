# 219：C 十三碼事件 handler 與原始依賴

**狀態：CONFORMED。三十函式 O0/O2 各 43,476 組原版/C 相同，八個負對照拒絕；限明示 callee fixture。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear；完整 chunks 與檔案偏移見 probe
- 出處：[`re/100`](../re/100-c-event-handlers-restoration.md)、[`re/99`](../re/99-c-hourly-update-restoration.md)
- 證據審查：DB `workplace/matching-decompilation/c-events/ida/input.exe.i64` 的 SHA-256 為 `67e6b4f56f2faa716b9f18afd61260a20329374146d29a41a4204317eff6eb97`。原始名稱、完整 chunks、直接 xref 與 instruction／operand 匯出保留，UID/GID 1000:1000。

## 1. 契約

三十函式、1919 bytes 的原始定位由 re/100 列出。原始暫存器與 segment 為 typed
`KiMachine16` 輸入；保存、byte／word wrap、FLAGS、SS stack、memory 寫入及
呼叫次序均保留。事件 8 的落下尾端沒有新增 CALL 或 stack return。

政治 helper 真實接線到既有經濟、任官、power、relation、RNG 與 queue writer。
災害 allocator／clear、遷都候選與軍團 raw 欄位更新由 C 真實執行。
`sub_135AB` 保留 `0x24` 與 raw memory；`sub_13138` 保留 `cmp ax, ax`。
不得修正原版看似不合理的行為，也不得由跨版本或 remake 推定分支。

## 2. 受控外部 callee

UI／文字／聲音／重畫、金額視窗、外交選擇、戰鬥準備作明示 fixture。
選擇回傳覆蓋 0–3，其他 fixture 為 RET；兩側入口參數與 stack 逐 byte 相同。
合法勢力、據點、官員、record、queue code 0–13、分離 segment 與 IF/TF=0。
原始 pointer 算法在 1 MB RAM 保留；沒有猜補未知高層玩法。

## 3. 驗證閘門與垂直鏈

原始 EXE／SINARIO.DAT → 原始 record → C handler／真實依賴 → globals／world／queue
逐 byte 與完整原版 RAM 抽查。O0/O2 結果相同，回傳、callee trace、FLAGS、stack 不作 mask。
矩陣須含每個 handler、每個 helper、多個 raw RNG、關係／支出／災害／任官邊界、
四劇本經真實 dispatcher 接線及八個可拒絕的錯誤版本。
每版 340 次完整 1 MB 核對相同，三十函式均有實際入口收據；矩陣與來源見 re/100。

正式 Go、玩家 UI、存檔流程不受本輪研究入口變更。C 來源與工具進 Git，原版、DB、
完整組語與記憶體收據留在本機 ignored `workplace/matching-decompilation/c-events/`。

## 4. 未解範圍

後續 C 視窗證據見 [`spec/220`](220-c-modal-control.md)。本規格的舊收據仍使用 modal
fixture；新收據執行真正視窗控制流與 CF retry，底層輸入／繪圖及正常玩家範圍仍未驗證。

| 項目 | 邊界 |
|---|---|
| 正常玩家事件 modal 與戰術 | 未驗證，fixture 不作完成證據 |
| Go 每時與事件三方 | 未驗證 |
| 中立 0x24 的玩家語意 | 未驗證，本輪只還原原始控制流 |
| C 機器碼、原作 compiler | 未驗證 |

## 據點持久效果的後續原生 C 驗證

據點整備與災害已由原生 C 驗證。`sub_14269` 的原始 `0x14269` 已接入據點輪轉、內政官與marker扣減，保留byte及word回繞，見 [spec/248](248-c-city-tick.md)。本頁較早事件handler的測試範圍仍保留；動畫物件更新的C驗證不代證完整事件畫面與自然玩家流程。
