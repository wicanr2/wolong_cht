# 252：C 軍團輪轉與行軍更新

**狀態：READY。17函式閉包與窄驗證契約已審查，原生C驗收尚未完成。**

- 日期：2026-10-10。
- 證據：[re/132](../re/132-c-army-update-restoration.md)。
- 輸入：DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 固定IDA probe SHA-256：`937f8738225ae0f9e5efc133245316dd075d845ffded97a22f5f211b4804f2a7`；DB SHA-256：`4f01452cc2df29a9306aa4b444d1d4125e3f27f7730fcfdadbefd9f9da3a4667`。
- 工具與位址：IDA Pro 9.4，database linear段基址10000。

## 原始契約

- 已證實，17函式434指令／1068bytes，沒有間接CALL／JMP或新MZ重定位；外部依賴使用既有CONFORMED C。
- 根入口`sub_125A3`從各側CS:D52／CS:9874建立DS／ES，讀CS:D18後固定處理16格，最後才比較1FC0及回寫cursor。
- 由原始0初值產生的八個批首皆需驗證；末批包含DS:4200臨時城兵。原始C保留第128格，不把Go的127常駐軍團範圍套回原版。
- 其他入口使用其原caller的DS、ES、SI、BX、方向byte與有效資料指標。佔用格由各側自身原LDS欄位取得，不跨機器複製。
- 原`126FF→12708`及`12804→12808`落下邊、byte／word回繞、CBW、CF與callee回傳都保留。
- 軍費、士氣、第一個同座標active軍團等行為按[re/132](../re/132-c-army-update-restoration.md)逐位址契約，不猜補規則。

## 實作與驗收

- 證據審查：17函式434指令／1068bytes已逐條核對。初態只使用合法raw記錄、完整4-byte道路點與明示guest frame；沿用既有route／outcome初態後，須各自重填本輪道路點及佔用格。
- 固定矩陣涵蓋四劇本的八批cursor、timer 0／1／2／FF、費用的兩次截斷、士氣byte回繞、127筆搜尋與第128格輪轉、兩條落下邊、道路方向、佔用格及野戰／攻城接線。
- 原版先通過獨立數值與資料模型，再比較C。根入口不推定callee的堆疊深度；真戰術frame以原始呼叫軌跡與既有已驗frame契約檢查。
- CS:CF3由fixture明示設定；不宣稱它與Go自然時鐘的hour cadence已對拍。外層退出須丟棄原生army frame，不繼續caller。

- 兩側各自raw初始化，固定原始RNG與等價輸入；C不執行原版CPU指令，不借原版執行後狀態。
- 以O0／O2比較完整RAM、VGA、DAC、暫存器、FLAGS、SS／SP、IN／OUT、API及callee軌跡；正常返回與outer轉移分列。
- 每個原始入口都需實際進入；變異編譯須能拒絕cursor步長／尾批、方向字寬、財政／士氣及碰撞分派等錯誤。
- 初態模型需先獨立檢查原版，再比較C；未確認的玩家語意維持原推論等級。
- 正式Go引擎、存檔格式及玩家路徑不在本輪變更範圍；原版素材與私有dump不入Git。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原生C與完整矩陣 | 初態與模型契約已審查；通過原生驗收前不提升CONFORMED。 |
| 完整主排程、正常玩家長程及原C機器碼 | 本輪局部契約不代證。 |
