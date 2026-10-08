# 224：C 位元對齊貼圖與六按鈕 caller

**狀態：CONFORMED。O0/O2 各 4,954 組完整 RAM／四平面／ABI 相同，十一個錯版拒絕。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具與位址：IDA Pro 9.4 database linear
- DB SHA-256：`0d0bd6bded0fe104b1648c8fa3af852ef209b57ee7e1464b84b786dc0f2affa8`
- 出處：[`re/105`](../re/105-c-aligned-blit-restoration.md)、[`re/18`](../re/18-tactical-button-glyphs.md)

## 1. 原始契約

還原 `sub_1F888`、`sub_1F938`、`sub_1F999`、`sub_1C7F4`、`sub_1C673` 共 429 bytes。
保留 byte／word 溢位、原始寄存器／堆疊、左右 mask、有號 JL、word byte 交換、
dummy latch read、SI 來源前進與全部 VGA I/O 順序。六按鈕使用原始表與實際圖庫。
裝置保存／恢復使用既有 C，底層 FAR 服務採固定明示 fixture。

## 2. 驗證閘門

原版與 C 的 RAM、四 plane、latch 與 port 狀態獨立，初始狀態與來源在執行前固定。
O0/O2 逐 case 比完整 1 MB RAM、四個 65536-byte plane、十四暫存器、FLAGS、
GC／seq／latch／index、port 順序與入口 stack 快照。內容區像素從 VGA 第 40 列取 400 列。
八種對齊、遮罩、width byte wrap、BH signed 邊界、word 來源跨界、真實資產六 slot、
原始位置重畫 caller 與既有外框呼叫鏈都須有收據；改錯遮罩、步幅或來源必須被拒絕。
編譯來源摘要進 CGO flags，並由 binary build info 核對。

## 3. 產物與未解範圍

研究 C、header、fixture、driver、來源索引與驗證摘要進版控。
本機原版、資產、DB 與重建結果留在 `workplace/matching-decompilation/c-aligned/`。
本輪不修改正式 Go 規則、UI 或存檔格式。

| 項目 | 邊界 |
|---|---|
| 正常玩家戰術畫面 | caller／primitive 收據不足以完成正常路徑 |
| 自然裝置輸入與實機 wall-clock | 固定 FAR fixture 與 VGA 模型 |
| 其他字型與文字 renderer | 本輪範圍外，尚未完成 |
| C 機器碼匹配 | 未驗證 |
