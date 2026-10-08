# 226：側欄計量條的原始位移與 byte 寬度

**狀態：CONFORMED。O0/O2 各 131,072 組原版／C／正式 Go 長度相同，UI 冷測通過。**

- 日期：2026-10-09
- 出處：[`re/106`](../re/106-c-rectangle-bars-restoration.md)、原始 `sub_1C775`／`sub_1C78E`
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear；DB 身分見 re/106

## 1. 玩家顯示契約

兵力 word 右移兩次後取 AL，再限制為 124。
體力 word 先右移一次，另右移一次後相加，再取 CL，最後限制為 124。
因此體力 3 的長度為 1；兵力 1024 的右移結果取 byte 為 0。
Go 的 `battleSideBarLengths` 應保留原始 u16／u8 計算；負數由既有顯示入口正規化為 0。
位置、顏色、高度、原始資料、遊戲規則與存檔格式不變。

## 2. 驗證閘門與未解範圍

對兩種計量各比較 0..65535，原版、C 與正式 Go 函式的相同長度。
Go package 冷測與現有 sidebar／selection／layout 測試應通過。
新測試使用原始相鄰邊界與動態結果，不能只把實作公式重抄到測試。
局部計量的 framebuffer 由 spec/225 的真實 C／VGA 矩陣驗證。

| 項目 | 邊界 |
|---|---|
| 正常玩家整體畫面 | 局部長度與 C plane 收據不足以代證 |
| 全場戰術／長流程 | 本輪未驗證 |
