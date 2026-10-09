# 234：C 原始大地圖顯示格閉包

**狀態：CONFORMED。O0／O2 各 2,622 組與 12 個錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[re/114](../re/114-c-world-map-cells-restoration.md)

## 原始契約

保留 40×23、每格 8 bytes、384-byte 來源列距、16-bit 段／位移回繞與原始 FLAGS。
推入端先裁切鏡頭再測 `0x10` 保護及四張上限。renderer 處理原髒格／快取條件，
保留第五個儲存槽與 `0xFF` 疊圖分支，不用 remake 的顯示結果回填原 C。
MDL 為 128 bytes，MCH 為 32-byte mask 加 128-byte 四平面 color。
原 `LODSW`／`MOVSW`／`STOSW` 與 REP／DF、VGA selector 及 memory bus 均保留。

`sub_1D7E7` 的有效輸入要求 DF=0；原 renderer 在呼叫前先 `CLD`。
受控 DF=1 會從 `CS:0xD858` 反向寫入自己的後續指令，原 guest 因改碼而改變控制流。
此破壞程式碼的狀態不當作局部 native C 函式驗收，失敗實驗另行保留。

## 驗證閘門

O0／O2 比完整 RAM、四 plane、暫存器、FLAGS、SS frame、函式入口與 I/O。
使用真實 MMAP.MAP／MDL／MCH，覆蓋初始化、鏡頭取樣、所有 256 個圖塊、
完整旗標值域、0–5 層記錄、推入裁切／容量、重畫與快取、DF／位移邊界。
圖塊合成獨立比來源 bytes；錯版必須拒絕，C 從固定 IDA 證據乾淨重生。

整幅原版測試的指令預算為 400 萬，源自 920 格的最大 6 次合成、
四平面搬移與控制指令上界；記錄實際最大步數。函式入口 trace 容量為 16,384，
覆蓋滿圖 5 個儲存槽加 `0xFF` 分支的 8,281 次入口；不調整原始分支或迴圈。

## 實作與收據

`tools/c_recovery/mapcells.c` 與生成內容保留十支原始函式。
[驗證收據](../re/c-mapcells-verification.json)與[DF 前提實驗](../re/c-mapcells-df-precondition.json)
記錄來源、原始素材、完整矩陣與限制；重跑 `tools/c_recovery_mapcells.sh`。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 正常玩家與上游 producer | 本閉包不代證完整場景退出與軍團／物件更新 |
| 原作者 C compiler／machine code | 尚未確認 |
