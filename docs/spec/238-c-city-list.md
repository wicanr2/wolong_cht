# 238：C 據點清單與遷都

**狀態：CONFORMED。O0／O2 各 520 完整裝置案例與十二錯版通過。**

- 日期：2026-10-09
- 證據：[re/118](../re/118-c-city-list-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約

保留原函式名稱、raw 入口、IDA 位址與運算元。原始近／遠呼叫、SS 清單與選中列的
特殊返回必須一致。排序採原描述子與原交換次序，保留 byte／word patch 和 JBE／JNB。
清單 renderer、字型、數字、滑鼠、事件 writer、對話與世界重畫沿用既有唯一 C 實作。

## 驗證閘門

- O0／O2 比完整 RAM、plane、暫存器／FLAGS、SS frame、呼叫與裝置軌跡。
- 原六欄排序、相等值、byte／word、有效清單筆數及原交換順序。
- 清單建立、空列 renderer、上下界、滑塊、拖動、選取與取消。
- 四劇本與多勢力；完整遷都接受／拒絕、選原首都重試。
- 真實素材與原始資料；C 不執行 guest CPU，受控 mouse input 明示。
- 編譯來源雜湊、錯版拒絕、固定 IDA 證據乾淨重生。

## 實作與收據

`tools/c_recovery/list.c` 與 `list_generated.inc` 保留原控制流、live patch 和特殊返回。
平台夾具以固定 mouse 事件與既有裝置服務執行，不呼叫 C-side guest CPU。
[版控收據](../re/c-list-verification.json)記錄 120 排序、36 builder、48 捲軸獨立核對、
24 表頭操作與 16 完整遷都。正式 Go 規則未變動。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 正常玩家長程與其他清單 | 本輪局部 caller 不代證 |
| 原作者 C 機器碼 | 尚未匹配 |
