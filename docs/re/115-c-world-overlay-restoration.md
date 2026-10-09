# 115：原軍團／物件到世界顯示格的 C

**狀態：CONFORMED。6 支函式與原始矩陣 handler 的局部 C 行為通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 固定 IDA DB SHA-256：`0eaddc302d905f7522730d63dd6822f6aa366bd8c34e413c60fa174c2993152d`
- 工具／位址：IDA Pro 9.4 database linear；file offset 為 linear 減 `0x10000` 加 512
- 入口：[ida_world_overlay_probe.py](../../tools/ida_world_overlay_probe.py)
- 前置：[re/72](72-world-map-display-list.md)、[re/14](14-mmap-mch-objects.md)、[re/114](114-c-world-map-cells-restoration.md)
- 契約：[spec/235](../spec/235-c-world-overlay.md)

`sub_11CC9` 依序呼叫軍團 `sub_12AF4` 與物件 `sub_12533`。
軍團先處理 bit `0x20` 的矩陣，再以存在值 `>=0xC0` 推入單格圖塊；
物件處理 32 筆 16-byte 記錄，以舊相位查表，髒旗標成立時才推進下次相位。
`loc_1D51F` 在裁切後推入 source 矩陣，保留 `0xFF` 透明格、保護格、容量及髒格。

原始 `0x1D5A9` 被 IDA 標成資料，檔案 bytes `80 FB 05` 是 `cmp bl,5`。
入口先把 AH 寫入 `CS:0xD5AB`，軍團傳 3，物件傳 5。C 讀取 live byte，
不把矩陣容量改成單格 producer 的 4。此控制流為已證實，來自固定原始指令；
首次 IDA 建指令遭既有資料項阻擋，只在一次性 DB 清除該項後解碼，保留原資料行。

## 來源與收據

六函式 461 bytes、raw handler 181 bytes，共 267 條原始指令。
[overlay.c](../../tools/c_recovery/overlay.c)與[生成內容](../../tools/c_recovery/overlay_generated.inc)
保留原名稱、位址與運算元。重跑 [c_recovery_overlay.sh](../../tools/c_recovery_overlay.sh)，
[驗證收據](c-overlay-verification.json)與[函式台帳](c-recovery-status.json)受版控。

| 矩陣 | O0／O2 每版組數 |
|---|---:|
| 勢力／五方向／座標的單格推入 | 990 |
| 軍團矩陣／四相位／小地圖 | 280 |
| 物件 type／八相位／dirty／座標 | 896 |
| 55 個真實 pattern 的裁切／容量／保護／槽位 | 12,320 |
| 小地圖八 bit 對齊／兩色／Y | 48 |
| 存在／dirty 門檻與兩次軍團巡覽 | 56 |
| 四真實劇本的軍團／物件巡覽 | 16 |
| 初始化 → 地圖 → 原 producer → renderer | 16 |

O0／O2 各 14,622 組完整 RAM／四 plane／FLAGS／暫存器／SS frame／I/O 與入口相同。
12,320 次 source pattern 與格子表另以高層裁切／透明／保護／容量規則核對，
沒有以兩側同時指錯來源的結果當作通過。12 個 O0 錯版全拒絕，C 乾淨重生相同，
實際編譯 flags 綁定完整 source digest。IF／TF=0、堆疊獨立，FLAGS 依固定 dosgolem 模型。

原版巡覽維持舊相位取圖與下次相位寫入，type 3 只作 raw fixture；不補自然 producer。
raw 容量立即數的 3 bytes 已以 GNU 組語匹配，總組語為 24,746 指令／56,290 code bytes。
原檔 67,099 bytes 的資料缺口仍由自備原版匯入；不散布原版檔。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 物件自然 producer 與完整場景 | 本輪不代證 type 3 的自然來源或正常玩家完整場景 |
| 原作者 C 工具鏈／機器碼 | 尚未確認 |
