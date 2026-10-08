# 216：月度暴風雨四個 globals word 保存

**狀態：CONFORMED。四個 storm globals 保存、原始 raw 錨點、round-trip 與 36 三方向量通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 9.4 database linear：`sub_122DB`，globals `0x10D22`／`0x10D24`／`0x10D26`／`0x10D28`
- 出處：[`re/96`](../re/96-c-monthly-world-update-restoration.md)，36 向量 `s3-p0-seed255`

## 1. 已證實契約

四個矩形 word 存在原版保存區塊 +0x32／34／36／38。
月結先清舊區域為 FFF0／FFF0／0190／0190；新 storm writer 成功後寫 MinX／MinY／MaxX／MaxY。
入佇列失敗保留清空區域，不保留舊 stormArea。Go 原有 StormArea 公式與 raw writer 次數沿用。

Go 在月結時將四個原始 word 存入 typed 保存欄位，Bytes 寫回同一原版位置。
原始 raw 保存錨點保持不變，RawBlock 仍回載入時的原始區塊。
未發生月結的載入／匯出仍保留原始 bytes，不猜補原版區域或改存檔格式。

## 2. 驗證與未解範圍

冷測 storm／無 storm 的四 word、原版格式 round-trip；36 三方向量的 globals 差異需歸零。
其他 queue／RNG 差異繼續報告，完整畫面／音效／長程仍未驗證。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 暴風雨正常 UI／音效與長程流程 | 未驗證，保存 word 的一致不替代玩家垂直鏈 |
