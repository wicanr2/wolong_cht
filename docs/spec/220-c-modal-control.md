# 220：C 外交／金額視窗與原始回應控制

**狀態：CONFORMED。十七函式 O0/O2 各 53,322 組原版/C 相同，十個負對照拒絕；底層 primitive fixture 邊界保留。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 出處：[`re/101`](../re/101-c-modal-restoration.md)、[`re/100`](../re/100-c-event-handlers-restoration.md)
- DB `workplace/matching-decompilation/c-modal/ida/input.exe.i64` SHA-256：`42181f593431ef87bb95cb5e66c6945fc7e92c86d0d8ef8f9e4df23209f338d2`

## 1. 原始契約

十七函式、1280 bytes 的定位由 re/101 列出。保留原始 AX／byte 回應、CX 訊息、DX 金額、
DS／ES、SS frame、register 保存及 FLAGS。0x13902 的 RNG／信賴度採納分支、0x139E8 的
非零金額下界／回應／扣款、原始選擇與金額 CF retry 均由 C 真實執行。
原始 SS 未初始化欄位不猜補。0x12216 的 code word patch 與恢復在 callee 入口觀測。

## 2. 受控 callee

IDA 0x20000 的遠呼叫以明示 fixture 提供座標、FLAGS clobber、RETF；不得混成 near RET。
兩個輸入 leaf 提供固定回應與 0–3 次 CF=1。選單忙碌與數值取消分開，完整原始迴圈仍在待驗函式。
數值服務 AX=30000 為上限；矩陣的 32768／65535 是 raw caller 邊界，不作合法玩家輸入證據。
文字／圖形／音效／其他底層服務為明示 primitive fixture。合法 pointer、分離 segment、
IF/TF=0；word 邊界與 raw byte 邊界另列，不當作合法玩家操作範圍。

## 3. 驗證閘門與垂直鏈

原始資料／record → C 視窗控制流 → 原始 helper／far ABI → world／globals／queue／RNG。
O0/O2 比全 register、FLAGS、callee 入口、stack、control fields、中途 patch 及完整 RAM 抽查。
矩陣須覆蓋十七入口、mouse word overflow、0–3 回應、byte variant、金額相鄰邊界、固定 raw RNG、
兩種 CF retry 與四劇本 handler 接線；十個錯誤版本均被拒絕。每版 417 次完整 1 MB 核對相同，十七函式都有實際入口收據，見 re/101。

研究入口不改 production 玩家路徑或存檔格式。來源與工具進 Git；原版、DB、記憶體收據
留在 ignored `workplace/matching-decompilation/c-modal/`。本規格不新增 remake 發行 gate。

## 4. 未解範圍

後續 C 數值輸入證據見 [`spec/221`](221-c-numeric-editor.md)。本規格的舊收據以固定值
取代 0x17C6E；新收據執行真正數值主迴圈、六鍵與財政 caller，圖形／實際 input 仍為 primitive。

| 項目 | 邊界 |
|---|---|
| 完整底層輸入／繪圖／音效與自然玩家視窗 | 本輪 primitive fixture 不作完成證據 |
| Go 每時／事件／視窗三方 | 尚未驗證 |
| 自然執行時序與硬體 | 尚未驗證 |
| C 機器碼匹配 | 尚未驗證 |
