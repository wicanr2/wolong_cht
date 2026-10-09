# 243：C 軍團行軍選點與命令分派

**狀態：CONFORMED。O0／O2各1,320完整裝置與十六錯版通過。**

- 日期：2026-10-09
- 出處：[re/123](../re/123-c-march-command-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與閘門

- 原軍團資訊到命令、目標選點、二／三項選單、十二表項九handler、路線與解散保留原始ABI。
- mouse／TALK／VGA／RND／補兵共用唯一C；五筆farcall保留執行期operand，SAHF及LES依原指令。
- O0／O2比完整RAM、四plane、暫存器／FLAGS、SS frame、呼叫、IN／OUT、DOS／font／mouse／sound。
- 四劇本、接受／取消／重選、兩鍵遮罩、熱區優先、雙floor與鏡頭夾制；全部12stage與9handler、門檻、隊列、DI residual、補兵與解散。
- 原版先經獨立模型；controller另捕捉原版分派前snapshot。每例預置明示raw RNG表/C/S，不調seed挑結果。
- 日期與時鐘入口計數分開於動畫RAM；不把取消建模為所有資料不變。
- 零缺字、來源／工具／素材雜湊、C重生、實編譯錯版拒絕、MZ反向重定位及指令獨立重組。
- 正式Go不改；Go冷測與全專案檢查，既有失敗獨立列出。

## 實作與收據

[march.c](../../tools/c_recovery/march.c)與產生來源保留原函式與指令，RNG呼叫路由接既有ABI。
[獨立模型](../../tools/c_recovery_march_data.go)驗原版狀態、隊列、DI、補兵／解散及RNG，C另比完整裝置。
[版控收據](../re/c-march-verification.json)、[原檔指令／MZ](../re/c-march-code.json)與[Go冷測](../re/c-march-go-verification.json)記錄範圍。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然長程／非法資料／原C機器碼 | 本輪不代證 |
