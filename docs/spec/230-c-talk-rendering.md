# 230：C 原始 TALK 與肖像快取

**狀態：CONFORMED。O0/O2 各 1,277 組全 RAM／VGA／ABI 相同，1,022 槽與十錯版通過。**

- 日期：2026-10-09
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具／位址：IDA Pro 9.4 database linear
- 證據：[`re/110`](../re/110-c-talk-rendering-restoration.md)

## 原始契約

使用原始 TALK.DAT、SINARIO.DAT、KAOGRF.DAT 與字庫，保留原始索引／八變體、
NUL 換行、CF 碼寬、七標記、SS 參數消耗、direct／0xFF ID 模式、姓名／呼び名欄位、
顏色、VGA 及四格快取替換。C 讀取 live near table，不用自行拼出的字串取代原模板。
快取 miss 以原始 DOS ABI 讀取 2,048 bytes，命中保持原 slot／cursor 與讀檔次數。

DOS／字型服務採固定 dosgolem API，兩側擁有獨立服務與 RAM／VGA／檔案 handle。
C 呼叫平台服務鉤子，不執行 guest CPU，不回填原版輸出。低階讀檔保留 LAHF／SAHF、
CF、stack frame 與原始錯誤返回，不修補原版控制流。

## 驗證閘門

O0/O2 比較全 RAM／四 plane／GC／seq／latch／port、十四 register／FLAGS、near／far／
DOS API 順序與 stack。覆蓋全部 1,022 個訊息槽、七標記、四劇本姓名來源、直接位址、
人物肖像及替換序列，包含字色、半形／全形、換行和控制標記。錯版必須被拒絕，
來源與實際編譯 flags 的 digest 相符。

## 未解範圍

| 項目 | 限制 |
|---|---|
| INT 50h 與實機磁碟時序 | 保留原始呼叫，不以平台 API 代證原媒體錯誤 UI 或 wall-clock |
| 正常玩家事件流程 | renderer 接線不取代自然輸入驗收 |
| 原作者 C 工具鏈與機器碼 | 尚未確認 |
