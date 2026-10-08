# 222：C 熱區、視窗 wrapper 與真實 query 接線

**狀態：CONFORMED。十函式 O0/O2 各261,367組原版/C與80組Go scalar相同，八個負對照拒絕。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 出處：[`re/103`](../re/103-c-hotspot-restoration.md)、[`re/22`](../re/22-strategy-command-tree.md)、[`re/47`](../re/47-main-screen-window-registry.md)
- DB `workplace/matching-decompilation/c-hotspot/ida/input.exe.i64` SHA-256：`e5b4874f382093fb14953bc6bb25fe8072338061f7ce72010621a37b8be71d5a`

## 1. 原始契約

十函式 491 bytes，定位在 re/103。原始 map 80×50、ES／DI init、段與 base、word arith、
原始零尺寸 byte loop、非零覆寫 CF、清除、AL query 與保存保持。
BL／DL low-byte F8h 對完整 word 等價 FFF8h，不能丟 high byte。
原始 window wrapper、tile flags、border caller 均以真正 C 接線，底層 VGA blit 保留 primitive。

## 2. 接線與驗證

原始十八 key table → 真實 C map 登記 → 固定 pixel／cancel leaf → 真實原版/C map query
→ 數值主迴圈 → 原始財政 globals。座標從實際 raw table 找，不手寫按鍵位置。
原版/C 比完整 register、FLAGS、map／flags／data／stack、callee 入口及抽查 1 MB。
全部 640×400 query、Y≥256、高 word／base wrap、重疊、原始零尺寸、wrapper／border、
十八鍵與有限數值計畫、四劇本接線，八個錯誤版本全拒絕，每版2042次完整1 MB核對相同，見re/103。
固定 input 預先宣告，有明示完成／取消及尾端取消。原始與 C 無 mask／結果挑選。

## 3. 產物與未解範圍

正式 Go 不直接猜補；研究 C／工具進 Git，原版／DB／RAM 收據只留 ignored
`workplace/matching-decompilation/c-hotspot/`。舊 glyph 名稱與 word 公式按同版直接證據勘誤，
保留歷史收據及原始定位，不把舊 C ABI 成功改寫成像素匹配。

| 項目 | 邊界 |
|---|---|
| 真實 VGA plane／blit／保存像素 | primitive 不代證實際 memory device |
| 自然滑鼠與玩家視窗 | 固定 pixel fixture 不作正常玩家證據 |
| 完整 Go 視窗三方 | 未驗證 |
| C 機器碼匹配 | 未驗證 |
