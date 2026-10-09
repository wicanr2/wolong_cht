# 131：交戰入口與戰術引擎的 C

**狀態：原始來源已核對，spec251 READY；原生C驗證尚未完成。**

- 日期：2026-10-10。
- 輸入：松崗DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4，既有py312-v1映像；位址為IDA database linear，段基址10000。
- 探針：[ida_engagement_probe.py](../../tools/ida_engagement_probe.py)；規格：[spec/251](../spec/251-c-engagement-tactical.md)。
- 探索閉包：179個命名函式5863指令／13967bytes，另有7個raw區段。這是來源普查，尚非C完成率。
- 探索Probe SHA-256：`191942cf78cdab19431b876ee82376dac6551ed25d6cf2f0bbe9198abe2832b6`。
- 探索Database SHA-256：`8646b160fb9bdb85085cf36d13f0f4309a8342eb0d9eabcf87ce43f1dccbe5d8`。

## 原始流程與定位

| 原始定位 | 原始證據及目前判定 |
|---|---|
| `sub_14A7B`，0x14A7B | 已證實野戰入口先選戰場、找對手，再清雙方對峙及處理敗方。 |
| `sub_14ADE`，0x14ADE | 已證實攻城入口區分已有守軍與臨時城兵，接入結果及據點易主。 |
| `sub_14B63`／`sub_14C72` | 已證實戰場取樣、名單與最高評分代表；保留原始運算元，語意沿用較新spec。 |
| `sub_11B5A`，0x11B5A | 已證實玩家未委任的戰術入口；0x11B76–0x11BDF為原落下尾端，不能補假RET。 |
| `sub_14F8A`，0x14F8A | 原城兵臨時軍團建立入口，納入本輪C來源；城兵主將127的既有證據見spec/191，C執行結果仍待完整矩陣。 |
| `loc_1A065`，0x1A065 | 已證實原始呼叫目標，IDA原分析混有資料；在一次性DB解碼原bytes，不修改原版。 |
| 0x1C0C5–0x1C30C | 已證實兩張原回呼表的raw handler與共享尾段；0x1C300位於0x1C2FD之TEST運算元內，不能把該候選當有效函式入口。 |

野戰／攻城結果沿用[re/130](130-c-battle-outcome-restoration.md)已驗證C。
原戰術支線納入本研究範圍，不以自動判定代替。
完整軍團更新、主排程、原C機器碼匹配與正常玩家長程仍是後續整合範圍。

## 間接分派與IDA候選審查

已證實，原`0x1C031`是`FF9786C0`間接CALL，必須讀原DS的C086表。
原表後方的padding與下一條指令可被誤讀為C300及168B，IDA候選xref不能直接當呼叫邊。
工具保留被拒絕的候選及原因，固定表與原始有效索引另行核對。
跨入中段的`0x1A562→0x1A531`與`0x1DBEB→0x1DB98`必須保留原位置及堆疊。

## 固定原始閉包

- 最終probe SHA-256：`32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9`；IDA database SHA-256：`cc4bf0da038a769fc7c12f8dcc47e884c1eb25ce959ceed0f86c3fb44182308f`。
- 234個命名函式7080指令／16696指令bytes，12個raw區段737指令／1867bytes；合計7817指令／18563bytes，沒有重複指令。
- 逐條核對直接CALL、JMP與落下後繼，缺失數均為0。八個間接CALL及一個間接JMP均有固定原表，表後程式碼不當作回呼。
- 已證實，六鍵UI由`0x160E6`之CX=6、`0x160EC`之AL=20及原LOOP產生20–25；`0x1ACA4`產生DL0–3。六項及四項表可封閉，15FAA的舊BX=A分支不擴張有效UI輸入。
- 已證實，`0x19FE0`以CS:D340恢復SP，原RET回`0x11B76`；SS不變。完整11B5A的cleanup最後於`0x11BDF`返回原交戰caller。
- 七個直接code patch消費端包含`0x1B55F`、`0x1B567`、`0x1A06A`、`0x1AB4E`、`0x1C601`、`0x1C604`與`0x1BE33`。保留原opcode及operand，需由live memory驗證，不能把掃描當作間接寫入已窮盡。

## 原生C與重組來源

[C來源](../../tools/c_recovery/engagement.c)、[介面](../../tools/c_recovery/engagement.h)與[產生器](../../tools/c_recovery_engagement_generate.py)保留全部原始指令及入口。
[原版/C runner](../../tools/c_recovery_engagement.go)、[平台橋接](../../tools/c_recovery_engagement_platform.go)及[fixture](../../tools/c_recovery/engagement_fixture.h)比較完整裝置狀態，尚未完成原生矩陣。
[隔離執行入口](../../tools/c_recovery_engagement.sh)及[驗證器](../../tools/c_recovery_engagement_verify.py)在完整矩陣與入口覆蓋核對前不發布CONFORMED。
[資料與流程審查](../../tools/c_recovery_engagement_data.go)、[邊界案例](../../tools/c_recovery_engagement_control_data.go)及[UI案例](../../tools/c_recovery_engagement_ui_data.go)保留固定原始初態。局部函式由兩側各自真初始化到首次A156入口後執行，不交叉複製執行後Snapshot。
[交戰前端案例與模型](../../tools/c_recovery_engagement_front_data.go)另審查地形、選軍、委任、自動結果及真正內層戰術返回。需要世界畫面的函式先真跑完退卻及世界還原，再從各側自有狀態執行。
[存檔與世界UI補充案例](../../tools/c_recovery_engagement_extra_data.go)、[戰鬥補充案例](../../tools/c_recovery_engagement_combat_data.go)及[獨立存檔模型](../../tools/c_recovery_engagement_save.go)補足實際入口與文件副作用；新增案例在實跑前維持未驗證。
[指令重組工具](../../tools/engagement_code_supplement.py)與[來源收據](c-engagement-code.json)記錄全部7817條指令的獨立重組，以及31條新增／替換指令和12條舊錯邊界。
[Go驗證紀錄](c-engagement-go-verification.json)保存既有正式Go的39套件冷測；C原生矩陣與專案檢查各自記錄，不互相替代。

## 已驗證範圍與fixture勘誤

- 四劇本各一例完整`sub_11B5A`真退卻、`sub_19FDC`恢復SP、`0x11B76`世界還原及真正RET，原版／C完整RAM、VGA、DAC、FLAGS、SS及API一致。
- 擴充抽樣304例／604階段通過；234個命名入口已有194個實際進入，12個raw入口已有10個。剩餘入口與錯版仍待驗，不由此提升為全閉包CONFORMED。
- 組語更新為25097指令／57101指令bytes；完整67099-byte EXE與原版相同，三個錯版拒絕。舊12條解碼及其來源保存在補充收據，不覆寫歷史source-map。

初次完整退卻後，原版停在`sub_1E81C`道路建圖。已證實根因是fixture沒有向DOS預留手動arena，並非城市座標或RNG。
原`0x1F66C`要求0800 paragraphs，MCB1153之後給出buffer1154，32768-byte讀檔區11540–1953F覆蓋了解碼目的12000。
失敗RAM在12000的64bytes恰等於壓縮MMAP.MAP的10AC4起點，符合第三次讀檔10004加AC0的地址計算。
兩側改以真DOS AH4A預留PSP至9B00，原decoder再依其既有分支申請剩餘空間；原decoder、seed及步數上限均未修改。

設定選單取消會由`sub_11C8D`重畫世界。停在A156的戰術狀態仍保留BATTLE資料，不能直接用來驗該世界UI。
此類案例改為兩側各自完成真退卻及世界還原，再用各側自有狀態執行；所有設定分支保留。

### 2026-10-10：垂直尋路案例的段位址勘誤

580案例／1,156階段的歷史抽樣通過234個命名入口，12個raw入口僅進入11個，缺`0x1BFBF`。
已證實，該案例把四筆`F8`輸入寫到`D2FE=arena+0x900`；原`0x1BD84`卻由`CS:D2FC`載入ES，`0x1BE10`讀取`ES:[BX]`。
真初始化將`D2FC`設為`arena+0x700`，因此舊patch沒有改到尋路消費端。
案例只訂正四筆patch共用段位址為`0x1200+0x700`，暫存器、種子與格位不變；原C及原版指令不變。
580組收據保留原範圍，不據此宣稱`0x1BFBF`已驗證；修正後仍需重跑並取得實際入口證據。

修正後`vertical-route`四劇本／八個比較階段通過，`raw_entry_1BFBF`實際進入四次。
收據在`workplace/matching-decompilation/c-engagement/results/probe-vertical-route.json`。
目前來源正在執行完整O2矩陣，O0及二十個錯版仍待；此窄結果不提升整個閉包狀態。

### 2026-10-10：第十個反例的輸入勘誤

第十個`path-live-opcode`反例曾以28案例／56階段回傳PASS。
該結果只代表舊輸入沒有觸發差異，不能當作錯版已拒絕。
已證實，`control_data.go`的三筆尋路輸入仍寫到`D2FE=arena+0x900`；原`0x1BE10`讀取的是`0x1BD84`載入的`CS:D2FC`。
原`0x1BD47`將CL寫入`0x1BE33`；案例CL為`EB`時，正常版無條件跳過垂直邊，錯版固定JZ則會在`F8`的bit8設立時進入`0x1BFBF`。

測試段位址標記`FFFD`改由各側自己的`CS:D2FC`解析。
runner共用既有`engagementApplyPatches`，兩側各自讀取初始化後的指標，不交叉複製記憶體。
三筆反例輸入及四筆垂直路由輸入均改用此標記；暫存器、種子、格位、資料bytes與案例數不變。
原C、錯版定義與比較器維持原狀；來源雜湊改變後，二十個錯版與O0／O2均須重新執行。
舊28／56及580組收據保留，不由這次輸入修正提升驗證狀態。

## 原生驗證執行入口


[包裝器](../../tools/c_recovery_engagement.sh)啟動[容器工作腳本](../../tools/c_recovery_engagement_container.sh)。設定 `WOLONG_ENGAGEMENT_DETACHED=1` 可在背景執行；容器內的 `timeout` 管理工作期限，包裝器印出唯一的 `job.log` 與 `job.exit` 路徑。

各模式共用來源清單，`mutants`、`normal-o2`、`normal-o0` 須依序執行。`job.exit` 記錄原生工作退出碼；背景 `full` 或 `controls` 結束後仍須執行驗證器。前景 `full` 與 `controls` 保留自動驗證。此執行方式不提升 C 驗證狀態。

## 新來源05ce的後續驗證

來源清單SHA-256為`05ce526895dcda6d13ae9b2ccd842e6922410c1601ac2a8fa282de4f6a5fd6e8`。
O2完整748例／1492階段通過，234個named與12個raw皆有實際入口；16組完整SAVE及兩側各256次MMAP還原通過。
256次戰術frame還原與outer退出0分列。1492個未發生outer轉移的階段包含508個warmup暫停，真正返回caller的階段為984個。
O2收據位於`workplace/matching-decompilation/c-engagement/results/O2.json`，SHA-256為`3e3f7fdc1c50b839795d20e4122ce1e351a53532c80fd8ddf04ada17408d4eb1`。
同來源二十個實編譯錯版都由狀態差異拒絕，沒有panic、unsupported或SIG代替比較。
O0完整矩陣仍在背景執行；驗證器要求它與O2收據SHA-256相同，完成前維持READY，不回填665／62台帳。
前述580例、錯段位址及當輪待驗狀態保留為歷史紀錄。

## 未解範圍


| 項目 | 狀態 |
|---|---|
| code patch與原生C完整執行 | 原始直接寫入已定位；仍需實作及動態比較。 |
| 完整原生C、固定狀態矩陣及錯版拒絕 | 尚未完成。 |
| 完整軍團／主排程、正常玩家長程及原C機器碼 | 本輪局部閉包不代證。 |
