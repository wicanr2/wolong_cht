# 118：C 據點清單與遷都

**狀態：CONFORMED。18 函式與 5 個 raw 入口通過同狀態驗證。**

- 日期：2026-10-09
- 輸入：松崗 DOS/V `KI.EXE`，67,099 bytes。
- SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- IDA 資料庫 SHA-256：`a61cd29b9daae8a6446f91e9fbe5e82f29101ad0beaac327dcfc6314283ab712`。
- 工具與位址：IDA Pro 9.4 database linear，段基址 `0x10000`。
- 匯出：`tools/ida_list_probe.py`，本機證據在 `workplace/matching-decompilation/c-list/ida/`。
- 既有證據：[re/26](26-list-window-engine.md)、[re/27](27-list-row-fields.md)。
- 規格：[spec/238](../spec/238-c-city-list.md)。
- 83 個原基準未覆蓋指令的組語補充：[c-list-code.json](c-list-code.json)，由固定 GNU binutils 逐 byte 驗證。
- 補充產生器：[tools/list_code_supplement.py](../../tools/list_code_supplement.py)，在既有固定 binutils image 內以 `--repo /repo --output /output` 執行，輸入為本輪 IDA 與原基準 source-map；完整冷重建入口為 `tools/matching_code_record.sh`。

## 原始契約

| 原始定位 | 已確認的行為 | 等級與來源 |
|---|---|---|
| `sub_16909` | 選原首都會顯示訊息並重選；比較兩城欄位後走三句對話、遷都 writer 與重畫 | 已證實，固定輸入指令 |
| `sub_17400`、`0x1743B`、`0x1745F`、`sub_1748F` | 192 城依所屬勢力建清單，512-byte SS 陣列保存記錄位址；每頁十列 | 已證實，固定輸入指令及 re/26 |
| `sub_1820E`、`0x1828F` | 寫入 builder／row callback、標題與欄位描述子，建立框、熱區、捲軸 | 已證實，固定輸入指令 |
| `sub_18412`、`sub_18463` | 選中列時丟掉一層返回位址，直接返回外層清單 | 已證實，`0x184B7` 的 `pop ax` 與後續 RET |
| `0x1857F`、`0x185B2` | 原描述子選欄位、byte／word 寬度與 JBE／JNB；排序會改寫原指令 bytes | 已證實，原六欄描述子與 patch writer |
| `sub_184DD`、`sub_1851A`、`sub_18546`、`sub_18607`～`sub_18713` | 拖動、上下捲動、比例滑塊與四 plane renderer | 已證實，固定輸入指令及 re/26 |

## 驗證結果

O0／O2 各 520 組完整 RAM、plane、暫存器／FLAGS、SS frame、呼叫與裝置 API 一致。
包含四劇本、三勢力的清單、120 組獨立交換排序參照、36 組獨立 builder 與 48 組捲軸公式核對。
六表頭實際點擊 24 組，完整遷都 16 組包含八次首都寫入與四次原首都重選；十二錯版全拒絕。
835 條原指令可乾淨重生 C；七個 patch 點帶出的 83 個原基準未覆蓋指令共 215 bytes 經固定 GNU binutils 匹配。
來源、工具版本、輸入、每個入口與反例見 [版控收據](c-list-verification.json)，
重跑 `bash tools/c_recovery_list.sh`。

## 標題指標勘誤

`0x1821E` 寫入 `CS:83D3`，即 `0x183D2` 的 MOV SI 立即值。
執行指令即讀取該值。原始 bytes、live patch、標題 renderer 與錯版驗證支持此結論。
re/26 舊查詢只搜尋地址文字，漏掉指令取立即值；原說法與勘誤保留在 WORKLOG／RESEARCH-LOG。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 其他清單家族 | 本輪據點描述子的結果不代證武將、軍團或勢力的完整 caller |
| 自然玩家長程路徑 | 局部遷都入口不代證自然選單或通關 |
| 原作者 C 與工具鏈 | 尚未確認 C 機器碼一致 |
