# 249：C 尋路、退卻與軍團潰散

**狀態：CONFORMED。O0／O2各4,588組完整狀態與十六錯版通過。**

- 日期：2026-10-10
- 出處：[re/129](../re/129-c-route-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與驗證

- 還原九完整命名函式436指令／1,016bytes與raw入口`loc_1491B`的107指令／244bytes；兩類分開計數。
- 保留三個CS code patch及其live消費端、DS圖表／ES世界資料、visited及環形佇列記憶體，完整比較正常與失敗CF及ABI。
- 原端點、連結方向、byte長度、節點與敵方成本、高位元、同成本處理及queue wrap均以原指令為準。
- 行軍重算、退卻、俘虜／轉屬、佔用圖與小地圖更新接既有真實C依賴，不以no-op替代。
- 原版先經獨立圖表／資料／固定raw RNG模型，再做O0／O2完整RAM／四plane／DAC／暫存器／FLAGS／SS／IN／OUT／API比較。
- 必須包含雙端點可區分、路徑不通、同成本、多成本、敵方／中立、環形佇列回繞、原code patch及成功／失敗重算與退卻案例。
- 實編譯錯版至少16個，均由狀態差異拒絕；避免等價錯版與無界非法圖表。
- 新解出的組語指令另記補充來源，獨立重組後納入整檔逐byte驗證；不把原版指令當非指令資料匯入。
- 正式Go不改。完整軍團更新、交戰／戰術、主排程與C機器碼匹配仍待後續。

[新增組語來源](../re/route-handler-code.json)、[指令覆蓋](../re/c-route-code.json)與[Go收據](../re/c-route-go-verification.json)保存可重生入口及證據；原生C完整狀態與十六錯版均已通過。

## 實作與收據

[route.c](../../tools/c_recovery/route.c)保留九函式與完整RET的raw入口；[獨立模型](../../tools/c_recovery_route_data.go)逐bucket模擬原訪問表／環形queue，另驗潰散及小地圖寫入。
[收據](../re/c-route-verification.json)保存4,588組完整狀態、十六錯版與來源綁定。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整軍團更新、交戰／戰術、主排程與自然玩家長程 | 本輪不代證 |
| 非法圖表／原硬體時間與原C機器碼 | 不由局部原生C比較推廣 |
