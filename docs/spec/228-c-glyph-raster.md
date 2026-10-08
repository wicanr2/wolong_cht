# 228：C 原始字型向量與 glyph raster

**狀態：CONFORMED。O0/O2 各 252 真實 glyph／場景全 RAM／plane／ABI 相同，九個錯版拒絕。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 出處與 DB：[`re/108`](../re/108-c-glyph-raster-restoration.md)

## 1. 原始契約

還原字型向量、兩 patched far call、32-byte SS buffer、三字 self operand、near／far
返回與原始 register／FLAGS／VGA 副作用。移除 glyph raster no-op fixture。
字庫只從原本唯讀輸入載入，兩側獨立平台服務與 cache；C 只使用服務 API，不執行 guest CPU。

## 2. 驗證閘門

兩種最佳化比較原始全 RAM／plane／ABI／latch／port與入口快照。
所有單／雙 byte、對齊／色彩／背景／模式／邊界／三字／原始場景與 live operand 都須覆蓋。
字庫 count／size／hash、取得向量、字模 buffer 與 call順序由資料驗證。
錯版必須被拒絕，生成 C、服務 API adapter 與實際 binary 的來源摘要相符。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 原版 TSR／逐字磁碟時序 | 採用成熟字型平台契約 |
| 實機／正常玩家操作 | 局部 glyph／場景不足以代證 |
| C 原作者工具鏈與機器碼 | 尚未驗證 |
