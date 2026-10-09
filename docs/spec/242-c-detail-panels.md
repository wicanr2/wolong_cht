# 242：C 據點資訊與軍團面板

**狀態：CONFORMED。O0／O2各700完整裝置與十二錯版通過。**

- 日期：2026-10-09
- 出處：[re/122](../re/122-c-detail-panels-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約

保留據點完整局部視窗、軍團欄位／六槽與退出、自勢力情報恢復的十個原始函式。
真實數字、字型、肖像、KYOGRF讀檔、VGA、熱區與mouse共用現有唯一C。
保留中立首都比較、上昇值符號／換色、景觀位移截斷、空槽天秤及ES／DI保存差異。
正式Go引擎行為不變。軍團指令上游仍另列，不以繪製函式代證整個行軍流程。

## 驗證閘門

- O0／O2完整RAM、四plane、暫存器／FLAGS、SS frame、呼叫、IN／OUT、DOS／font／mouse／sound相同。
- 四劇本、據點原始名稱與圖、neutral／capital、五類型、上昇值正負、全16個景觀nibble與張15回繞。
- 原版先經獨立欄位、數字／字串參數及KYOGRF buffer模型，再比C。
- 軍團總兵力高半部、byte士氣、六槽種類1–4與獨立標示的type0受控邊界。
- 面板02／04旗標、控制器左鍵重試後右鍵退出及原始ABI。
- 零缺字、固定來源／工具／素材雜湊、可編譯錯版拒絕、C乾淨重生、原指令覆蓋與獨立重組。
- 正式Go冷測與完整專案檢查，既有文件失敗獨立列出。

## 實作與收據

[details.c](../../tools/c_recovery/details.c)與產生來源保留原指令流程。
[資料參照](../../tools/c_recovery_details_data.go)先核對原版的數字／字串／貼圖參數、圖檔內容與世界資料不變。
[版控收據](../re/c-details-verification.json)、[指令覆蓋](../re/c-details-code.json)與[Go冷測](../re/c-details-go-verification.json)記錄實際範圍。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 軍團指令上游／自然長程／原C機器碼 | 本輪不代證 |
