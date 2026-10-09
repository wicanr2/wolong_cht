# 250：C 自動戰鬥與據點易主

**狀態：CONFORMED。O0／O2各8,364組完整狀態與二十錯版通過。**

- 日期：2026-10-10
- 出處：[re/130](../re/130-c-battle-outcome-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與驗證

- 原生C保留21個原始函式、byte／word／32-bit乘除與截斷、XLAT及真實既有依賴，不改寫成新的戰鬥規則。
- 自動判定保留AL=0的第一方特殊係數、各加8、比例上限100、同值分支、六槽亂數損害、大將最低1、兩種士氣式與壞滅bits。
- 退卻保留己方據點停留、原搜尋caller、300／首都門檻及CF，不在士氣或大將歸零時猜補階段。
- 易主保留舊owner、內政官／外交官、據點數、遷都、先轉向後滅亡、鄰接快照與小地圖。
- 玩家勢力滅亡必須真正恢復保存SS／SP並退出原caller；C runner不續行被丟棄的native stack。
- 每個測試以固定raw RNG初態及兩側各自raw資料開始，C不讀原版執行後Snapshot，不挑seed。
- 原版先經獨立資料／算式／呼叫與固定亂數模型，再做O0／O2完整RAM／四plane／DAC／暫存器／FLAGS／SS／IN／OUT／API比較。
- 有效分母與完整raw記錄是測試前提；原非法DIV不包裝成正常成功結果。
- 至少20個可分辨的實編譯錯版由狀態比較拒絕，涵蓋算式、傷亡、退卻、易主、官員、滅亡及顯示依賴。
- 正式Go不改。完整交戰caller、戰術、軍團更新及C機器碼匹配仍待後續。

[Go收據](../re/c-outcome-go-verification.json)記錄正式Go冷測及隔離自測，原生C完整狀態與二十錯版均已通過。
[指令覆蓋](../re/c-outcome-code.json)保存原始21函式的位址、bytes與獨立組譯結果。

## 實作與收據

[outcome.c](../../tools/c_recovery/outcome.c)保留21函式與真實既有依賴；[獨立模型](../../tools/c_recovery_outcome_data.go)從raw初態計算戰鬥、易主、退卻及原caller是否退出。
[收據](../re/c-outcome-verification.json)保存完整狀態、二十錯版與來源綁定；[指令覆蓋](../re/c-outcome-code.json)保存原始bytes與位址。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整交戰／戰術／軍團／主排程與正常玩家長程 | 本輪不代證 |
| 非法DIV／記錄、原硬體時間與原C機器碼 | 不由局部原生C比較推廣 |
