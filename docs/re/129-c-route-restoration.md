# 129：自我修改尋路、退卻與軍團潰散的 C

**狀態：CONFORMED。九函式與一raw入口，O0／O2各4,588組完整狀態與十六錯版通過。**

- 日期：2026-10-10
- 輸入：松崗DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為IDA database linear，段基址10000。
- Database SHA-256：`31091a8a266e73ba1e903c5bca69409dca513179fd3114e1ec7074af1d17546d`。
- Probe SHA-256：`ff99cfbd72b7211d14ac5d7544ad7ffa5af9fa8ac8c421d2f1d4311a438d1a66`。
- 探針：[ida_route_probe.py](../../tools/ida_route_probe.py)；規格：[spec/249](../spec/249-c-route.md)。

## 已證實原始閉包

| 原始定位 | 範圍 |
|---|---|
| `loc_1491B`，0x1491B至0x14A0F | 107指令／244bytes。原始raw入口有完整RET，保留原名稱；搜尋從AX目標開始，BX／CX是兩個終止節點。 |
| `sub_14A0F`，0x14A0F至0x14A7B | 原始完整命名函式，保留函式邊界。檢查port visited、連結方向、邊長與反向port，再寫入環形佇列。 |
| `sub_147BB`／`sub_1487B` | 行軍重算與退卻包裝，保留當前連結兩端、CF、步進±4及原先的DS恢復。 |
| `sub_1291A`／`sub_12977`／`sub_129C3` | 原軍團潰散、俘虜／轉屬判斷、佔用計數與訊息。 |
| `sub_12BA8`／`sub_19656`／`sub_196CF` | 原dirty flag清除與小地圖還原，沿用真實VGA與訊息依賴。 |

九個完整命名函式共436指令／1,016bytes；一個raw入口107指令／244bytes，總計543指令／1,260bytes。
`sub_125A3`、`sub_12662`及戰後處理`sub_1474A`只作導航。

## 自我修改與成本

0x14931、0x14936、0x1493B分別寫入CS:49B8、49BE、49D2。
0x149B6、0x149BC、0x149CD隨後消費這些live立即值，不能沿用靜態1234h／12h。
這些是原始指令區的CS offset；與上述IDA線性位址分開。

原成本包含節點+4、非己方據點+0A6h並設bit15，以及連結byte長度。
搜尋保留環形佇列、最低成本bucket與原訪問標記，不能用一般heap的相容結果代替。
既有規則與限制見[spec/192](../spec/192-route-cost-model.md)及[spec/46](../spec/46-post-battle-retreat.md)。

已證實的呼叫者有兩個：`sub_147BB`的IDA 0x147EA／0x1483A，以及`sub_1487B`的0x148D7／0x148E5。
同一包裝內的兩個site屬互斥分支，不能據此把每次入口算成兩次搜尋。

## 來源與驗證入口

[79條新增指令來源](route-handler-code.json)與[覆蓋收據](c-route-code.json)分開保存；來源使用完整檔案SHA-256，收據可綁定新的整檔manifest而不改寫來源。
[Go收據](c-route-go-verification.json)保留本輪39套件冷測及隔離自測，原生C完整狀態與十六錯版驗證已完成。

## 驗證結果與回填

- O0／O2各4,588組完整RAM／四plane／DAC／暫存器／FLAGS／SS／IN／OUT／API相同；群組為search 1,200、edge 536、replan 576、retreat 256、collapse 1,668、minimap 352。
- 原版先經獨立raw圖表／訪問標記／環形queue與固定raw RNG模型，再與C比較。兩側從各自raw初態開始，C不讀原版執行後Snapshot。
- 三個live立即值、雙終止節點、異成本／同成本、敵方高位元、192-node低frontier反覆回繞、early RET與DF、行軍／退卻及潰散分支均涵蓋；十六實編譯錯版均由狀態差異拒絕。
- 本輪圖表是可重生的synthetic packed graphs；本機未找到可用的完整原版raw port-table快照，不宣稱已驗所有原版道路或正常玩家路徑。
- 新小地圖blit的完整planes由獨立模型核對；既有196ED疊畫只驗真實callee參數與完整C比較，沒有另寫第二份舊renderer。
- [C來源](../../tools/c_recovery/route.c)、[獨立模型](../../tools/c_recovery_route_data.go)與[收據](c-route-verification.json)可回查；重跑`bash tools/c_recovery_route.sh`。
- 新增79指令／179bytes後，整檔25,078指令／57,045bytes來源重建67,099-byte EXE完全相同；來源與收據分開，所有supplement維持完整檔案SHA-256驗證。
- 已回填[re/65](65-ai-march-decision-chain.md)、[spec/192](../spec/192-route-cost-model.md)及[spec/46](../spec/46-post-battle-retreat.md)。舊成本模型列`oq-2ef4b20c0b00482c0266`原為merge-target／Issue#3，舊說明「佇列結構與其他成本項沒逐條讀」已由本輪來源取代；原列完整內容仍可由Git基準4324137回查。Issue#3的Go等距／tie-break問題與遠端Issue狀態未改。

## 戰後閉包的後續原生 C 證據

戰後閉包已由原生 C 驗證。`sub_196CF`的原始`0x196CF`與自動判定、易主或小地圖依賴已由固定raw狀態與原版／C完整比較核對，見[re/130](130-c-battle-outcome-restoration.md)。含玩家勢力滅亡的原SS／SP恢復及caller不續行；這不代證Go完整交戰、正常玩家長程或原C機器碼。較早收據保留原範圍。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整軍團更新、交戰／戰術、主排程與自然玩家長程 | 本輪閉包不代證 |
| 非法圖表、原硬體時間與原C機器碼 | 不由局部比較推廣 |
