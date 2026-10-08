# 232：C 原始選列、捲動與 popup caller

**狀態：CONFORMED。O0/O2 各 138 全裝置／ABI 相同，三個 XOR 還原與十二錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[`re/112`](../re/112-c-choice-selector-restoration.md)

## 原始契約

保留 SS 24-byte selector frame、總列／可見列／絕對 index／可見 index、XOR 反白、
標題與 row callback、帶狀 Y=100..131、游標狀態 0／1／2 保存恢復、press count 與 CF。
原 `0x9465` stride、`0x9467` 字串、`0x9441`／`0x9443` 位置是 live operand。
`sub_193E9`、進言 caller 與已還原 mouse／TALK／VGA 組成原始 C 閉包，不改 Go 操作。

## 驗證閘門

O0/O2 比完整 RAM／四 plane／GC／seq／latch／IN／OUT／DOS Mouse／十四 register／FLAGS、
near／far／API／SS frame。覆蓋 accept／cancel、可見與頁邊界、live callback／stride／座標、
原始選單內容與游標狀態、map／pixel 保存還原，進言 speaker／advisor caller 的真實字形與等待。
保留原 `nullsub_5`，不以猜測畫面搬移取代；固定輸入不冒稱自然玩家操作。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然玩家／進言 producer | 局部 consumer 與固定序列不代證 |
| scroll helper 的執行期改寫 | 本輪依目前原始 image bytes，是否被改寫另追來源 |
| 原作者 C 工具鏈與機器碼 | 尚未確認 |
