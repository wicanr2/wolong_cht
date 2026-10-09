# 123：軍團行軍指令、選點與狀態分派的 C

**狀態：CONFORMED。25函式，O0／O2各1,320完整裝置與十六錯版通過。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro9.4；位址均為IDA database linear，段基址10000。
- Database SHA-256：`248af036e3ccf80b8bf6df1712662a2f9ff599504a538ddc9c5cc7eb30ae3fb3`；probe SHA-256：`c864ebad098ef8c7ce49b34fd3c0d0a30e591922c01bbee15980a65d41ac8e3d`。
- 證據：[re/45](45-corps-command-mode.md)、[re/84](84-popup-row-band-and-world-cursor.md)、[re/85](85-march-target-hit-test.md)、[spec/149](../spec/149-march-target-map-picker.md)。規格：[spec/243](../spec/243-c-march-command.md)。

## 已證實原始定位

| 入口 | 契約 |
|---|---|
| `sub_17F90`／`sub_17FDB` | 軍團資訊接原選點與命令；取消後外層仍CALL14325並OR2。 |
| `sub_1703C`／`sub_1709E` | 熱區優先；BX==2取消，BX==3不直接取消；CB..D3圖塊且登記X/Y相同才命中，搜尋無上界。 |
| `sub_1804E` | 原游標格位置夾制與二／三項選單，SAHF還原CF，取消恢復原AL項數。 |
| `sub_11F7F` | 世界像素減像素鏡頭並夾制，再分別右移4相加；兩分支都回BX按鍵遮罩。 |
| `sub_159A6`／`sub_15AB6`／`sub_12151`／`sub_121B2` | 原表159D2[22]的間接入口與鏡頭夾制，使用真正mouse farcall。 |
| `sub_14325`與九個handler | CS4358十二項有限表，非玩家且stage<8時索引加4，保留門檻、隊列、入口DI、補兵與RND。 |
| `sub_14548`／`sub_1463E`／`sub_14651`／`sub_14689`／`sub_11C8D` | 路線目的地、解散、計數及重畫。 |

25函式共613指令／1477bytes；五個farcall的MZ段字1000在IDA載入為2000，依file_bytes核對原檔。
C farcall讀取執行期原operand，不寫死IDA段字。

## 原版邊界

- `0x17FB4`／`0x17FB7` 不檢查取消CF，仍分派並設bit1；最終stage可能繼續改變或立即解散。
- `0x17052` 只比較BX==2；`0x170D2` 據點搜尋無額外上限。
- `0x11FC6`／`0x12073` 都回BX=BP，picker不套主畫面98A3bit7的禁止分派條件。
- `0x143D3` 讀入口DI+18，14325沒有設定DI，不補成目標城。
- `0x1466F` 只清軍團+00；解散另退兵、清主將職務、減軍團數與occupancy，不清整筆記錄。
- 選點沒有11D8E時鐘呼叫；重畫仍消費物件旗標／動畫phase，不能宣稱world RAM不變。

## 驗證結果

- O0／O2各1,320組：384分派、576 handler、24解散、32路線、44圖塊命中、24命令、28完整controller、24選點／重畫、96游標、32鏡頭、24小地圖、28選單、4重畫。
- 原版先經獨立模型 15060 項核對；controller只由原版生成分派前snapshot，先驗命令四欄再跑狀態模型。C不生成此期望值。
- 每例在執行前設定原raw RNG：C=seed XOR A5、S=seed、table[i]=byte(i×73+i/2+seed)。固定矩陣seed為0／17／34／51／68／85／102／119，原版與C初態相同，不重擲。每版實際RND入口 16 次。
- 選點包含非據點重試、熱區31壓過有效據點、小地圖22經159A6→15AB6→121B2，以及左右鍵同按。日期不變且11D8E呼叫數0；不宣稱動畫RAM不變。
- 十六個實際編譯錯版全拒絕，兩側缺字0，完整RAM／四plane／FLAGS／SS frame／裝置及API軌跡一致。
- [版控收據](c-march-verification.json)、[指令覆蓋](c-march-code.json)、[Go冷測](c-march-go-verification.json)保存來源與結果；重跑 `bash tools/c_recovery_march.sh`。
- 613原指令／1477bytes全已有覆蓋，五筆MZ反向重定位與獨立重組通過；本輪新增組語0條。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然長程、非法stage／資料損毀與原C機器碼 | 不由局部閉包代證 |
