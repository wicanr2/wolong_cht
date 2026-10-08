# 213：武將月度評分 +0x1F 寫回

**狀態：CONFORMED。月度評分 typed data、載入／保存、byte wrap 與 36 三方向量通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原版定位：IDA Pro 9.4 database linear `0x155A6`–`0x155EC`，70 bytes
- bytes SHA-256：`02cdbf0d15bd8b42736fcf248f786d502a0a62d2d0de9b8068e684488915b293`
- 原始向量：[`re/96`](../re/96-c-monthly-world-update-restoration.md)

## 1. 已證實契約

原版只處理存在旗標 >=0x80 的 127 武將，終止槽 127 不碰。
前三能力 byte 各右移 4，再加武力 byte 左移 1、統率 byte 左移 1，
全部 byte 相加以 8-bit 繞回，寫入記錄 +0x1F。此欄位的高階用途本輪不重新命名推測。

原版月結順序是生產力 → 武將入仕／俘虜 → 評分 → queue 月度初始化。
Go 在入仕處理後、compactEventQueue 前重算存活武將的同一個 byte，不加入 RNG。

## 2. typed data 與保存

`General.MonthlyScore` 保存原版 +0x1F，載入讀 byte、Bytes 寫同一位置。
未觸發月結時只保存原值；月結後按原版公式重算。既有能力與存在旗標沿用，不建立第二套規則。
停用／終止槽與未知 bytes 依原始 raw block 保留，名稱只描述可證實運算。

## 3. 驗證

新增月結更新、byte wrap、停用武將與載入／保存 round-trip 冷測。
36 個完整三方月結向量的載入基準仍需相同，general +0x1F 差異應歸零；
所有其他 raw 差異與 RNG 繼續報告，不以評分修正宣稱完整月結 parity。

## 4. 未解範圍

| 項目 | 邊界 |
|---|---|
| 評分的高階 UI／AI 用途 | 未驗證，本輪只接已證實計算與原版欄位 |
| 其餘 Go 月結、queue 與 RNG 差異 | 未驗證，另按三方 audit 定位 |
