# 246：C 城市與軍團初始化

**狀態：CONFORMED。O0／O2各5,008組、5,416次全狀態與十錯版通過。**

- 日期：2026-10-09
- 出處：[re/126](../re/126-c-world-bootstrap-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與閘門

- 原生 C 還原五個完整函式及 11BE0 的七指令前段；不合成 RET，不計為完整主迴圈。
- 城市依原 byte 算式／低 nibble／中立 24 與玩家色碼重畫中心及四鄰格。
- 軍團以 flags≥80、127 槽與 byte 累加為準，覆寫原 descriptor，保留 inactive 與槽 127。
- 評分及日期使用既有唯一 C 實作，驗原 register／DS 契約、合法日期及存在武將槽 126／127。
- 原版與 C 各自從 raw 初態產生初始化結果，不把原版 prefix snapshot 當作 C 輸入。
- 獨立模型先驗原版 map／occupancy／descriptor／評分／日期，O0／O2 再比完整 RAM／VGA／DAC／FLAGS／SS frame／IN／OUT／API。
- 另驗原生 prefix→信賴度／進言退出鏈；C bridge 不續行被丟棄的 caller。
- 十個實編譯錯版必須由狀態比較拒絕，來源／bytes／原始邊界與獨立重組需一致。
- 正式 Go 不改。完整 11BE0 與 C 機器碼匹配仍留待後續。

## 實作與收據

[bootstrap.c](../../tools/c_recovery/bootstrap.c)實作五函式；七指令前綴使用明示入口，普通函式分派仍不支援完整 11BE0。
[獨立模型](../../tools/c_recovery_bootstrap_data.go)計算 map／occupancy／descriptor／評分／日期；[收據](../re/c-bootstrap-verification.json)、[指令覆蓋](../re/c-bootstrap-code.json)、[Go 冷測](../re/c-bootstrap-go-verification.json)記錄來源及驗證。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整主迴圈與正常玩家路徑 | 本輪不代證 |
| 非法座標、原硬體時間及原 C 機器碼 | 不由局部原生 C 比較推廣 |
