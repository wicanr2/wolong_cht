# 225：C 矩形與戰術計量條閉包

**狀態：CONFORMED。十六函式 O0/O2 各 134,958 組原版／C 相同，十二個錯版拒絕。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具與位址：IDA Pro 9.4 database linear
- DB SHA-256：`3f3a95793967805c9df4f87de5db577b6eadb0aa3a39d4837be8c9d43a8dc9a5`
- 出處：[`re/106`](../re/106-c-rectangle-bars-restoration.md)

## 1. 原始契約

矩形描框／實心、clip flags、水平／垂直 mask、REP STOSB、write mode、暫存器／FLAGS
與堆疊按原始控制流還原。原始 caller 連成兩層選取、六隊清除、六待機條、四側欄條與窗口填黑。
資料 bank 使用已核對的 ICONGRF 第 1／3 段；X／Y 與資料偏移的位址空間分開。
計量計算保留 word 右移、AL／CL byte 取值與最後上限，體力不改成乘三除四。

## 2. 驗證閘門

原版與 C 初始狀態在執行前固定，獨立裝置與記憶體。O0/O2 比完整 1 MB RAM／
四 plane／原始 ABI／port／latch 與入口快照。內容區從 VGA 第 40 列取 400 列。
涵蓋反向／裁切／signed 端點、首尾 mask、clip 邊消失、DF、真實圖庫與各 caller。
計量原版／C／Go 比較完整 raw word 範圍，錯版必須被拒絕，編譯摘要綁定實際 binary。

## 3. 產物與未解範圍

研究來源與驗證摘要進 GitHub，原版、資料與狀態只留本機 `workplace/matching-decompilation/c-rect/`。
正常玩家、自然時序、其餘 renderer 與 C 機器碼仍未完成。Go 顯示修正另查
[spec/226](226-sidebar-bar-widths.md)，不改規則與存檔格式。

| 項目 | 邊界 |
|---|---|
| 正常玩家戰術／視窗 | caller 矩陣不代表自然操作完成 |
| 自然滑鼠與硬體時序 | 固定 FAR fixture 與平台模型 |
| 完整 C 還原 | 本輪只涵蓋十六函式閉包 |
| C 機器碼匹配 | 未驗證 |
