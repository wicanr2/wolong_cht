# 223：C 四平面 blit 與保存／恢復

**狀態：CONFORMED。八函式O0/O2各9,942組完整RAM／四plane／latch／I/O相同，八個負對照拒絕。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 出處：[`re/104`](../re/104-c-vga-blit-restoration.md)、[`re/03`](../re/03-image-blitter.md)
- DB `workplace/matching-decompilation/c-vga/ida/input.exe.i64` SHA-256：`7e11f676288faef8bf0e105d9a6ee8c338cf8a77fe8fe6064c867fc64abd7118`

## 1. 原始契約

八函式467 bytes，定位見re/104。保留全register、DF／FLAGS、source連續性、row pitch80、
dummy ES read、MOVSB／REP MOVSB、GC I/O序列、plane major保存與原始近返回。
原版MZ與CPU真實執行；C算法獨立，標準VGA由兩個獨立dosgolem裝置狀態提供。
標準平台契約只引用IBM手冊與成熟模擬器，結果限平台模型，不冒稱實機逐週期一致。

## 2. 驗證閘門

固定初始RAM／四plane／GC／seq／latch、來源／目的pointer、尺寸與DF，兩側一致且執行前固定。
O0/O2比較全部14register、FLAGS、stack／RAM、四個65536-byteplane、GC／seq／latch與原始I/O序列。
涵蓋mask／row／word wrap、dummy latch、save／restore／wrapper與有限回圈，八個錯版全拒絕，另310次indexed pixels與18條三步roundtrip相同，見re/104。
I/O event順序與port/value比較，原版Step原樣保留作定位；C沒有相同CPU指令數，不混稱時間parity。
source manifest digest進實際CGO build flags並由binary build info核對。

## 3. 產物與未解範圍

後續 C 位元對齊證據見 [spec/224](224-c-aligned-blit.md)，追加非 byte 對齊貼圖、
caller 與第 40 列起的內容區抽查；原先完整 plane／ABI 收據仍保留。

研究來源與工具進Git，原版、DB、RAM、plane收據只留本機ignored `workplace/matching-decompilation/c-vga/`。
正式Go／玩家路徑不猜補，原版資料不散布。

| 項目 | 邊界 |
|---|---|
| 正常玩家完整畫面／視窗 | 局部blit不代表production route |
| 實機時序／wall-clock | 只使用固定VGA模型 |
| 完整Go圖形三方 | 未驗證 |
| C機器碼匹配 | 未驗證 |
