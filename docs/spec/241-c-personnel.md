# 241：C 人事選單與任免流程

**狀態：CONFORMED。O0／O2各2,332完整裝置與十二錯版通過。**

- 日期：2026-10-09
- 出處：[re/121](../re/121-c-personnel-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約

七個原始函式以可重生 C 保留原始位址、運算元、FLAGS與SS frame。
選單查 live CS table；四條流程共用既有唯一清單、TALK、mouse、font、VGA及世界重畫。
任命只寫兩個欄位並保留經費；解任 XCHG FF、清經費與職務。
候選清單保留原篩選，已任官員不被第一層過濾；第二層取消回第一層。
本輪新增研究 C，正式 Go 行為不變。

## 驗證閘門

- O0／O2 完整 RAM、四 plane、暫存器／FLAGS、SS frame、呼叫、IN／OUT、DOS／font／mouse／sound相同。
- 四劇本、全合法官員索引與FF、非零經費正對照、任命保留／解任清零。
- 第一層取消、第二層取消、已有人／空缺拒絕與重試、成功後再選同目標、四項真實選單分派。
- 原版先經獨立欄位模型，C再同狀態比較；來源與素材雜湊、字型零缺字、錯版拒絕及乾淨 C 重生。
- 全部185原指令核對版控組語覆蓋並獨立重組；正式Go冷測與完整專案檢查，舊失敗獨立列出。

## 實作與收據

[personnel.c](../../tools/c_recovery/personnel.c)與產生來源保留原始控制流。
[資料參照](../../tools/c_recovery_personnel_data.go)先核對原版任免欄位與世界資料保存，外交流程僅排除既有22個顯示cache byte。
[版控收據](../re/c-personnel-verification.json)、[指令覆蓋](../re/c-personnel-code.json)與[Go冷測](../re/c-personnel-go-verification.json)記錄實際範圍。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 正常長程與原作者 C 機器碼 | 本輪不代證 |
