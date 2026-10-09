# 251：C 交戰入口與戰術引擎

**狀態：READY。原始控制流與九張固定分派表已核對，原生C驗證尚未完成。**

- 日期：2026-10-10。
- 證據：[re/131](../re/131-c-engagement-tactical-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 原始來源與有效輸入

- 固定probe SHA-256：`32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9`。
- 固定IDA database SHA-256：`cc4bf0da038a769fc7c12f8dcc47e884c1eb25ce959ceed0f86c3fb44182308f`。
- 234個命名函式與12個raw區段，合計7817條唯一指令／18563指令bytes。另10個chunk資料bytes保留為資料，不算指令。
- 原始直接CALL、JMP及落下後繼均已在本圖或已驗證callee內，沒有缺失後繼。間接CALL讀原DS／CS的live word，間接JMP保留五個原分支。
- 六鍵回呼的有效UI輸出由`sub_160CC`建立AL20–25；四方向回呼的DL由`sub_1ACA4`產生0–3。各表後面的指令不算額外回呼。
- C可用靜態原始指令圖及原SS堆疊保留所有控制流。原`0x19FE0`改SP後按實際RET續行到`0x11B76`，不能返回丟棄的frame。
- 自我修改的opcode與立即值必須從各側live code memory消費，不以初始literal取代。初始合法資料、合法表索引與有效DIV為fixture前提。
- 完整`sub_11B5A`由原`sub_19946`初始化戰術狀態，再跑原`sub_19FA0`至真實終了及世界還原。C不得把此支線換成自動判定或假返回。

## 契約

- 以`sub_14A7B`／`sub_14ADE`為入口，保留未委任時原`sub_11B5A`戰術分支；自動結果沿用已驗證C。
- 原始函式名、位址、bytes、chunks、運算元與資料表不改名覆蓋。附加語意保留證據等級。
- 原生C保存16-bit分段存取、字寬、FLAGS、SS／SP、FAR與原始分派表。跨函式中段與落下尾段不補假RET。
- 戰術結束改變SP後，執行實際RET位置，不能續行已丟棄的caller。
- 各側由唯讀原始資料及明示固定raw狀態初始化。C不執行原版CPU指令，不從原版執行後Snapshot初始化。
- 固定RNG及輸入在執行前選定；timer／裝置邊界明示為fixture，不宣稱實機wall-clock。
- 比較O0／O2完整RAM、VGA planes、DAC、暫存器、FLAGS、SS、IN／OUT及API；正常返回與非區域返回分列。
- 以實編譯錯版驗證邊界、分派、傷亡及戰術狀態差異可被拒絕。所有已實作函式需有實際入口證據。
- 正式Go規則維持既有範圍，既有冷測另記；不得由C局部比較提升完整Go戰後Issue狀態。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 原生C及完整狀態收據 | 驗證完成前不提升為CONFORMED。 |
| 完整軍團／主排程與正常玩家長程 | 本輪不代證。 |
| 原C機器碼與實機硬體時間 | 不能由原生語意比較推廣。 |
