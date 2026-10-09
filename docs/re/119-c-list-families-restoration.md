# 119：C 軍團、武將與勢力清單

**狀態：CONFORMED。11 函式／12 raw，O0／O2 各 804 組完整狀態一致。**

- 日期：2026-10-09
- 輸入：松崗 DOS/V `KI.EXE`，67,099 bytes。
- SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- IDA 資料庫 SHA-256：`63967f83d99ba47b84f6ce73d721f13e8f22b514bbf7212abfdd9f6fe3b0fec8`。
- 工具與位址：IDA Pro 9.4 database linear，段基址 `0x10000`。
- 匯出：`tools/ida_catalog_probe.py`，本機證據在 `workplace/matching-decompilation/c-catalog/ida/`。
- 既有證據：[re/26](26-list-window-engine.md)、[re/27](27-list-row-fields.md)、[re/73](73-new-game-faction-list.md)。
- 規格：[spec/239](../spec/239-c-list-families.md)。
- 新增指令：[c-catalog-code.json](c-catalog-code.json)，170 指令／361 bytes；產生器為 [tools/catalog_code_supplement.py](../../tools/catalog_code_supplement.py)，在既有固定 binutils image 內以 `--repo /repo --output /output` 執行。
- 正式 Go 冷測：[c-catalog-go-verification.json](c-catalog-go-verification.json)；沿用 `eob-remake-go:1.26.7-ebiten2.9.9-20261008-r2`、`tools/go.sh` 與本輪專用快取，原始素材實際掛載為唯讀。

## 原始契約

| 入口與回呼 | 控制流 | 等級與來源 |
|---|---|---|
| `sub_1716D`、`sub_171D3`、`0x171A8`、`0x17217` | 掃 126 軍團；玩家清單或精確 X/Y 清單，座標版本不限玩家勢力 | 已證實，固定原始指令 |
| `sub_175FA`、`sub_17663`、`0x1763C`、`0x176A0` | 掃 128 武將；指定勢力或玩家無職且排除君主；175FA 的勢力來自入口 DS:0CFF | 已證實，固定原始指令 |
| `sub_178A7`、`sub_17906`、`0x178E5`、`0x17944` | 掃 22 勢力；後者排除玩家，兩入口都先更新外交顯示 cache | 已證實，固定原始指令 |
| `sub_17B3C`、`0x17B6F` | 開局勢力清單原點 (136,104)，每次排序狀態為 0 | 已證實，固定原始指令及 re/73 |
| `0x1727D`、`sub_1770C`、`sub_1799C`、`sub_17BC0` | 原名稱、數字、身分、哨兵、換色、空列與十列 renderer | 已證實，固定原始指令 |
| `sub_17A7A` | 原和平位與 19/20、99/100/101 邊界；LOOP 回到 `0x17A9A`，每筆清 DL | 已證實，固定原始指令與回跳目標 |

各家族狀態 0 重建自然順序，其餘欄位依原 byte／word 與 JBE／JNB 描述子排序。
共用引擎沿用 [re/118](118-c-city-list-restoration.md) 的唯一 C 實作。

## 舊結論勘誤

- `0x1768A` 的 MOV CL,CS:98AA 覆寫前面的 XOR；選武將保留排序狀態。
- `0x1735C` 的 MOV AL,[SI+6] 讀取 u8 士氣，接著清 AH。
- `0x17A5C` 的 BX=9003 是三位數字；四格欄寬仍沿用原幾何。
- `sub_17BC0` 已逐欄核對，開局列包含君主、軍師、首都與兩個三位數字。

## 驗證結果

804 組包含 168 builder、168 頁 renderer、140 caller、156 表頭、56 選取／取消、
64 軍團邊界、24 武將身分、16 外交 cache、12 入口 DS 案例。
168 builder／240 排序／16 cache 另以獨立資料參照核對；28 次成功選取回傳原記錄位移。
十三個編譯錯版全拒絕，兩側缺字 0，902 原指令可乾淨重生 C。
來源、工具、素材與回鏈護欄見 [c-catalog-verification.json](c-catalog-verification.json)，
重跑 `bash tools/c_recovery_catalog.sh`。受控軍團、俘虜與邊界欄位只作測試初始狀態。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 空清單加非零排序狀態 | 原 DEC CH 可能下溢；不加入推測性修正 |
| 正常玩家長程流程 | 局部清單入口不代證任命、外交或完整新遊戲流程 |
| 原作者 C 與工具鏈 | 尚未確認 C 機器碼一致 |
