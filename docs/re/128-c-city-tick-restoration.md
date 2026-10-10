# 128：據點輪轉與動畫物件更新的 C

**狀態：CONFORMED。18函式，O0／O2各8,144組完整狀態與十六錯版通過。**

- 日期：2026-10-10
- 輸入：松崗DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為IDA database linear，段基址10000。
- Database SHA-256：`a888707683a83c6a2f146506b042d771892b4af464949afa5396f3dfcb159c22`。
- Probe SHA-256：`1922512894237831508aa8d2ba041d8f6b389a9f473775010e3d1fb2d1d542f5`。
- 探針：[ida_tick_probe.py](../../tools/ida_tick_probe.py)；規格：[spec/248](../spec/248-c-city-tick.md)。

## 已證實原始閉包

| 入口與原始定位 | 範圍 |
|---|---|
| `sub_13EFD`，0x13EFD | 每次更新一個據點；游標加20，達1800回零。保留佔用表低7bits、+85A與+841的比較／寫入、非中立AI與整備／災害順序。 |
| `sub_13F74`至`sub_14155` | 鄰城四槽、外交門檻、16-byte SS候選、旗標、回報冷卻與同節點軍團分派。保留byte回繞、CF及原始搜尋。 |
| `sub_14194`，0x14194 | 據點整備、內政官計數、原亂數讀取與城兵上限。相關機制見[經濟](../mechanics/40-economy.md)。 |
| `sub_14269`，0x14269 | 原marker扣除+851，差額再影響+850、+84E與+853。+84E的16-bit減法不改成飽和減法。 |
| `sub_14575`／`sub_145C1` | 原請求容量、武將掃描與既有出陣編成依賴，保留原caller參數與部分成功的CF。 |
| `sub_12459`／`sub_1248A`（0x1248A）／`sub_124FF` | 32個動畫物件的存在門檻、倒數、dirty bit、後16槽座標／速度與亂數偏移。 |
| `sub_10CDE`／`sub_1EB11`／`sub_1EB5E` | 回報提示的呼叫順序及原始IN／OUT控制；以固定平台I/O契約驗證，不推導實機音高或時間。 |

共18函式、562指令／1,340bytes。主排程`sub_11CD0`、軍團更新`sub_125A3`與完整主迴圈`sub_11BE0`只作導航。
主排程的靜態直接呼叫普查列出180個函式，另有7筆未封閉邊界；間接分派仍需另查。
該普查使用同一KI.EXE、IDA9.4及IDA線性位址，保存在本機`workplace/matching-decompilation/c-tick/census/`，不代表已還原：

- Probe SHA-256：`937f8738225ae0f9e5efc133245316dd075d845ffded97a22f5f211b4804f2a7`。
- Database SHA-256：`4f01452cc2df29a9306aa4b444d1d4125e3f27f7730fcfdadbefd9f9da3a4667`。
- 探針來源SHA-256：`3d9b8f941dbc960b56b5ea2b965f259443b1e7b11b5254274784a6289997b247`。

## 驗證紀錄

[Go收據](c-tick-go-verification.json)保存本輪39套件冷測與隔離自測；原生C驗證通過，完整文件檢查另由Go收據記錄。

- O0／O2各8,144組完整RAM／四plane／DAC／暫存器／FLAGS／SS／IN／OUT／API一致；群組為city 816、neighbors 256、threat 192、growth 1,408、disaster 1,280、reinforce 676、objects 3,476、alert 40。
- 原版先通過獨立資料／亂數模型，再與各自raw初態的C比較。每個測試seed在執行前固定，258-byte初態雜湊與設定公式存入收據，不重擲或挑通過結果。
- 據點192游標、佔用低7bits、原先勢力寫入、鄰接／外交／哨兵、內政官、災害回繞、32物件門檻／倒數／速度與座標均有分辨性反例。
- 原14575→145C1→既有16E8F成功／部分成功分派保留原callee參數與CF。舊編成區域由既有出陣證據及本輪完整狀態比較承接，不把它重寫成新的獨立六槽模型。
- 保留14137原讀848的距離式；十六實編譯錯版均由狀態差異拒絕。提示I/O模型核對平台輸入及輸出順序，C只讀自己的獨立裝置，不消費原版IN stream。
- [C來源](../../tools/c_recovery/tick.c)、[獨立模型](../../tools/c_recovery_tick_data.go)、[收據](c-tick-verification.json)與[指令覆蓋](c-tick-code.json)可回查；重跑`bash tools/c_recovery_tick.sh`。
- 562指令／1,340bytes均已有組語覆蓋，獨立重組一致，完整67,099-byte EXE仍相同。
[指令覆蓋](c-tick-code.json)記錄562指令／1,340bytes的兩次獨立重組與既有組語來源對映。

## 物件速度與座標回繞勘誤

已證實，`sub_1248A` 的0x124A7至0x124B7及0x124D1至0x124E1依storm bounds產生DL／DH的±1，
0x124F8與0x124FB再加入速度高byte；座標硬回繞另在0x124B9至0x124CB及0x124E3至0x124F5。
X固定界線為−16／400，Y為−16／272。這與無storm時四個global為−16／−16／400／400的
[既有契約](../spec/216-monthly-storm-globals.md)分開，不把global的Y上限改成272。

[re/14](14-mmap-mch-objects.md)與[spec/146](../spec/146-map-cloud-objects.md)原先的立即夾制、反向與關進矩形說法撤回。
現有Go已區分速度調整與固定回繞，這次訂正文件，沒有新增玩家規則。

## 生產力扣減的計算界線

已證實，0x1428F的倍率是`[SI+84F]`，也就是生產力word `P`的高byte。
marker超過防災值時，差額`D`為1至255，0x14293／0x14295將乘積右移兩次：
`loss=floor(D*(P>>8)/4)`。

合法據點記錄中，`P>>8=0`時損失為0；其餘情況有`loss≤63.75*(P>>8)<P`。
因此0x14297的原word減法不會借位。所有65,536個`P`在最大`D=255`的界線亦已核對。
原先「借位夾零」錯版在這個輸入域與原版等價，現改用0x1428F錯讀低byte的錯版。
舊正向收據、等價錯版與來源另存本機`workplace/matching-decompilation/c-tick/equivalent-control-review/`，沒有改變正常演算法或挑選亂數。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 軍團更新、交戰／戰術閉包、完整主排程與自然玩家長程 | 後續還原，不能由本輪局部對拍代證 |
| 原硬體時間及原C機器碼 | 固定平台I/O與原生語意比較不代證 |

## 軍團輪轉與行軍C補證

先前僅作導航的 `sub_125A3`（IDA 線性 `0x125A3`）已納入17函式的軍團切片，通過固定raw局部矩陣。
最終狀態與收據見 [re/132](132-c-army-update-restoration.md) 與 [spec/252](../spec/252-c-army-update.md)。
本文件的歷史收據與推論等級維持原範圍；完整主排程、正常玩家長程及原C機器碼匹配仍未完成。
