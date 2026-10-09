# 247：C 世界點選、右鍵分派與淡入

**狀態：CONFORMED。O0／O2各2,148組完整狀態與十二錯版通過。**

- 日期：2026-10-09
- 出處：[re/127](../re/127-c-world-interaction-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與閘門

- 還原 `sub_109D0`、`sub_11E46`、`sub_11F0E`、`sub_159B7`、原始 `nullsub_1`、`sub_15AA2`、`sub_15E4C`、`sub_161B6`，共159指令／396bytes。
- 淡入使用原16色、強度0至16、兩次垂直回掃等待與既有唯一DAC算式；平台時間只驗固定I/O契約。
- 點選保留原map／occupancy指標、城市搜尋、四個word共八byte的SS區域，以及取消／城市／多軍團選擇迴圈。
- `sub_11F0E`只限制座標低byte，Y上限22、X上限35，原始原點與選單參數照原值。
- 右鍵以AL乘2索引原26-word表。完整讀取live表，覆蓋全部26個合法索引，另以修改表項的反例分辨硬編碼路由。
- `sub_159B7`末端落入原`nullsub_1`的RET，只消費一次返回位址，不合成額外RET。
- 三個清除動作分別清除98A6的bit2／bit1／bit0，並保留原顯示清單呼叫參數與DS契約。
- C共用已驗證的顯示、選單、城市／軍團面板與滑鼠函式，不以no-op取代這些依賴。
- 原版先驗獨立資料與呼叫參數模型，再做O0／O2完整RAM／四plane／DAC／暫存器／FLAGS／SS／I/O／API比較。
- 實編譯錯版須由狀態比較拒絕；來源雜湊、原始邊界、live table與獨立重組需一致。
- 正式Go不改。完整11BE0與更新排程11CD0仍待後續還原。

[Go檢查收據](../re/c-interaction-go-verification.json)記錄本輪正式Go冷測及隔離自測；原生C的完整裝置比較及十二錯版已通過。

## 實作與收據

[interaction.c](../../tools/c_recovery/interaction.c)保留原八函式與單一落下RET；[獨立模型](../../tools/c_recovery_interaction_data.go)由raw資料驗城市／軍團選擇、原DS表、旗標與DAC。
[收據](../re/c-interaction-verification.json)與[指令覆蓋](../re/c-interaction-code.json)保存來源、2,148組完整狀態及十二錯版。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整主迴圈、更新排程與正常玩家長程 | 本輪不代證 |
| 非法表索引、原硬體時間及原C機器碼 | 不由局部原生C比較推廣 |
