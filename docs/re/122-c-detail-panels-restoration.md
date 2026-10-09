# 122：據點資訊與軍團面板的原始 C

**狀態：CONFORMED。十函式，O0／O2各700完整裝置與十二錯版通過。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為 IDA database linear，段基址 `0x10000`。
- Database SHA-256：`4c02d01d142286b9ba9ece847b292ddfbf85cff9901abf3db91d3e30aa0443b1`；probe SHA-256：`d4b28311a9c26af5c7cacc19f696603cebfe3666e7d823f5b3c02c5276c9945d`。
- KYOGRF.DAT：69,120 bytes，SHA-256 `e086f526bdada5baf751c41d2f73a78a0ba70002f282f63a9f33114542ed933f`。
- 證據：[re/32](32-strategy-detail-panels.md)、[re/50](50-city-info-window.md)、[re/51](51-corps-info-window.md)。規格：[spec/242](../spec/242-c-detail-panels.md)。

## 已證實定位

| 原始入口 | 契約 |
|---|---|
| `sub_17E1F` | SI為據點索引×32，加840後完整畫框、欄位、景觀，等右鍵並關閉；未保存SI。 |
| `sub_17E4A` | 據點名、城主／中立、類型／首都及四個原數字；中立仍讀DS:0603判首都。 |
| `sub_17F1A` | 高nibble×1200，CX清0，只以16-bit AX seek原KYOGRF；張15回繞到0E00。 |
| `sub_17F61`、`sub_17F81` | 原據點框、場景3與關閉。 |
| `sub_1807B` | 原軍團框、場景4、肖像、三名、總兵力×10與byte士氣；數字BX為0F04／0F03。 |
| `sub_1812A` | 六槽圖示與兵數；僅保存BP／DS／SI，ES與DI不恢復。 |
| `sub_1817D` | 關軍團框；CS:98A6 bit02開啟時恢復自勢力情報。 |
| `sub_15E1E`、`sub_15E2D` | 自勢力情報框與場景0；後者寫DS:98A6 OR2，入口DS契約為CS。 |

十函式共310指令／711bytes，所有callee沿用既有唯一C。
本輪包含據點完整局部視窗與軍團繪製／退出；軍團行軍指令上游仍另列範圍。

## 原始邊界與勘誤

- `0x17E2C` 只加840，入口SI為據點索引×32；原入口不只地圖，還有指令列的 `sub_162FB`／`sub_163BF`，見re/50的原始callers。

- `0x17EBA` CWD保留上昇值的符號，負值BX=0A04。
- `0x17EF3`–`0x17EF7` 左移取得據點索引的DH；neutral並未跳過後續首都比較。
- `0x17F31`／`0x17F33` 張14起FC00可讀至10E00；張15只傳0E00，不能把DX高半部補進seek。
- `0x18116`／`0x1811B` 士氣只讀byte並清AH；不可改為word。
- `0x18150` 的基底12C0令type4讀1500綠色天秤；type0下溢至D200，不是騎馬空隊。
- `0x1818B` 測bit02，`0x1819C` 呼叫15E1E恢復自勢力情報，並非小地圖bit04。

## 驗證結果

- O0／O2各700組：400據點欄位、64景觀、48完整據點視窗、各4據點框／關閉、80軍團欄位、48合法六槽、4受控type0、各16軍團關閉／自勢力恢復／自勢力框。
- 原版先經獨立參數、世界資料、KYOGRF buffer與ABI模型 7324 項核對，再與C比較完整1MiB RAM、四plane、FLAGS／SS frame與裝置軌跡。
- 十二個實際編譯錯版全拒絕，兩側缺字0；來源乾淨重生一致。
- 中立首都覆蓋、超出自然兵力的u16上界、type0圖庫下溢與景觀張15均明示受控資料邊界，不冒充自然玩家狀態。
- [版控收據](c-details-verification.json)、[指令覆蓋](c-details-code.json)與[Go冷測](c-details-go-verification.json)保存來源及結果；重跑 `bash tools/c_recovery_details.sh`。
- 原310指令／711bytes全已有覆蓋，另以固定binutils獨立重組完全一致，本輪新增組語0條。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 軍團行軍指令上游、自然長程與原作者C機器碼 | 不由資訊繪製代證 |
