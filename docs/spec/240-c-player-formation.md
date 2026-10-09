# 240：C 玩家軍團編成介面

**狀態：CONFORMED。O0／O2 各 200 完整裝置案例與十錯版通過。**

- 日期：2026-10-09
- 證據：[re/120](../re/120-c-player-formation-restoration.md)、[re/30](../re/30-corps-formation-ui.md)、[re/49](../re/49-corps-formation-window.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約

七個原始函式與全部 247 指令轉為可重生 C 敘述，保留原始名稱、位址、運算元與 FLAGS／堆疊。
使用既有唯一的軍團規則、清單、TALK、滑鼠、字型與 VGA。
初始化只改 type；切換先退再分；主將 0 拒絕；取消不重算派生欄，池飽和與合法熱區同原版。
本輪只新增研究用 C；正式 Go 引擎行為不變。

## 驗證閘門

- O0／O2 比較完整 RAM、四 plane、暫存器／FLAGS、SS frame、呼叫、IN／OUT、DOS／font／mouse／sound。
- 四劇本、三池邊界、六槽 type 循環、空主將拒絕、確定／取消、外層選將後返回清單。
- 直接 renderer 與初始化的非零哨兵欄位，依獨立資料模型核對 slots／reserves。
- 原始素材雜湊、字型零缺字、編譯來源 digest、具體錯版拒絕與乾淨來源重生。
- 全部原指令核對既有組語覆蓋；有新增時先獨立重組再整合。
- 正式 Go 冷測與完整文件／資產檢查；既有失敗獨立列出。

## 實作與收據

[formation.c](../../tools/c_recovery/formation.c)與生成來源接回既有唯一的清單、規則及 renderer。
獨立分兵模型在 [資料參照](../../tools/c_recovery_formation_data.go)，原版先經模型檢查，再與 C 比較。
[版控收據](../re/c-formation-verification.json)、[指令覆蓋](../re/c-formation-code.json)與
[正式 Go 冷測](../re/c-formation-go-verification.json)記錄實際範圍。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然玩家長程／原作者 C 機器碼 | 本輪不代證 |
