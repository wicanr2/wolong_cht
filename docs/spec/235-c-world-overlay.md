# 235：C 原始軍團／物件顯示 producer

**狀態：CONFORMED。O0／O2 各 14,622 組、12,320 次 source 矩陣與 12 個錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[re/115](../re/115-c-world-overlay-restoration.md)

## 原始契約

保留兩次軍團巡覽、bit `0x20` 矩陣分支、`>=0xC0` 存在分支、X word／Y byte、
圖塊 byte 加法、四個 formation 相位、MCH metadata 與 `0x100` source 基址。
物件保留 32 筆、舊相位查表、bit 0 清除與下次相位 `&7`，不替 type 3 命名。
矩陣保留 signed 裁切、來源列距、透明 `0xFF`、保護 `0x10`、live capacity 3／5 與髒格 `0x20`。
小地圖點按原 VGA read／write 順序產生，flags 依固定 dosgolem 模型。

## 驗證閘門

O0／O2 比完整 RAM、四 plane、FLAGS／暫存器／SS frame、函式入口與 I/O。
真實四劇本、MMAP.MDL／MCH，覆蓋軍團存在／dirty／方向／formation、物件相位／type、
矩陣裁切／透明／容量／保護／重入、八種小地圖 bit 對齊與兩色分支。
閉包接原 C 顯示格與 renderer，錯版必須拒絕，C 從固定 IDA 證據乾淨重生。

## 實作與收據

`tools/c_recovery/overlay.c` 與生成內容保留六函式及一個 raw entry。
[收據](../re/c-overlay-verification.json)記錄來源、完整同狀態矩陣、live capacity
與錯版，重跑 `tools/c_recovery_overlay.sh`。正常玩家與 type 3 自然來源仍依未解範圍判讀。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 正常玩家／完整場景與 type 3 | 局部輸入不代證完整自然路徑 |
| 原作者 C compiler／machine code | 尚未確認 |
