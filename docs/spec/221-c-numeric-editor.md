# 221：C 數值輸入、原始按鍵與財政返回

**狀態：CONFORMED。十七函式 O0/O2 各 788,931 組原版/C、787,839 組 Go scalar 相同，十個負對照拒絕。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 出處：[`re/102`](../re/102-c-numeric-editor-restoration.md)、[`spec/78`](78-amount-input-editor.md)、[`re/101`](../re/101-c-modal-restoration.md)
- DB `workplace/matching-decompilation/c-numeric/ida/input.exe.i64` SHA-256：`8a4d51f1057f192aba3f3e24369737121d9155cf28ed2329b3805e67a5e4f59e`

## 1. 原始契約

十七函式、598 bytes 的定位見 re/102。保留原始word cap、SI、AX、DX、multiply carry／
飽和、DIV 模型、按鍵表、LODSB／CLD、LAHF／SAHF、stack、裝置保護與返回。
完整主迴圈使用真正六個 handler。完成 handler STC 使主迴圈結束，主迴圈回 CLC；
右鍵取消回 STC。0x5E 的清零按鍵與真正取消不混成一件事。
財政上限為 100／10000，後三列確定後除十再寫原始 CS globals；取消不寫回。

## 2. 受控 leaf 與 Go 對照

輸入 stream 為有限 raw key／cancel 事件，事先固定，至少含一個完成／取消出口，另追加明示尾端取消。
裝置 FAR leaf 提供固定有號狀態與 pointer、會 clobber FLAGS，原始 PUSHF／POPF 保留。
圖形、文字與 popup leaf 為明示 fixture。合法 segment／pointer、IF/TF=0；非法高 key 不送入跳表。
Go 用正式 `state.EditAmountValue` 比相同有效 current／cap／action 的結果值，
bool 只表示 action 有效，不比較為原版 CF。原版/C 的 raw u16 與完整 ABI 另驗。

## 3. 驗證閘門與垂直鏈

原始 table → C 按鍵／主迴圈 → 財政 caller／既有 modal → 原始 globals／scalar。
O0/O2 全 register、FLAGS、SS、callee trace、控制欄位與完整 RAM 抽查相同；十七函式都有實際入口。
全 u16 digit／multiply／delete、相鄰 cap、十八格、完成／取消、原始 FAR／table 與四劇本接線。
十個改錯版本全拒絕，每版 6164 次完整 1 MB 核對相同，見 re/102。Go 可比數值矩陣與完整 C ABI 不混稱 normal-player parity。
來源與工具進 Git；原版、DB、memory 收據留在 ignored `workplace/matching-decompilation/c-numeric/`。

## 4. 未解範圍

| 項目 | 邊界 |
|---|---|
| 真實圖形／input／popup 與自然玩家數值視窗 | primitive fixture 不取代 UI／像素證據 |
| 完整 Go 財政／視窗三方 | 局部數值核心不作完整玩家證據 |
| 自然執行時間與硬體 | 未驗證 |
| C 機器碼匹配 | 未驗證 |

## 熱區 callee 後續證據

後續 C 熱區證據見 [spec/222](222-c-hotspot-map.md)。本規格原始 bytes／ABI 收據保留，
`sub_17D5F` 的十八格導覽更正為 raw 熱區登記，舊 `glyph` 群組名不作語意證據。
