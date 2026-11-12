# 125：正常滑鼠出兵、城兵攻陷與玩家攻方戰術

**狀態：remake正常玩家路徑抽驗通過。** 城兵分支63項、駐城軍團分支67項。
原版同狀態碰撞與亂數矩陣仍未驗證，不關閉Issue #30。

- 日期：2026-11-03。
- 來源基線：`bd09fd40e344efee9d01e286a70f2a2303890caf` 加目前工作樹。
- 遊戲程式：本機 `workplace/scratch/wlgame-battle-hit`，SHA-256
  `bcd7904992f486eb64fe4513b0a02c78a3d0594a8a1ca303ab53dc647a2e64c0`。
- 工具鏈：`san1-go:1.26.7-ebiten2.9.9-r1`，image ID
  `sha256:aac706e67af842cea8a6ad6011a1828161722381e1d7e2f21398aa31f28c07a6`。
- 素材：松崗DOS/V `KI.EXE` SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；
  `SINARIO.DAT` SHA-256
  `21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`。
- 操作契約：[spec266](../spec/266-mouse-only-desktop-toolbar.md)、
  [spec267](../spec/267-original-choice-band.md)、[spec268](../spec/268-advise-mouse-lists.md)、
  [spec269](../spec/269-peaceful-turnback-step.md)、[spec270](../spec/270-tactical-hit-boundaries.md)。
- 遭遇契約：[spec175](../spec/175-encounter-standoff-countdown.md)、
  [spec186](../spec/186-standoff-corps-before-city.md)。本次不修改遊戲執行邏輯。

## 2026-11-13：登城首拍與高層讀取

目前程式`wlgame-cao-scale-r1`接入[spec36續片](../spec/36-ground-planes-and-climbing.md)，戰報提示維持預設關閉。
[正常滑鼠](../../workplace/scratch/ui-windows-20261104/cao-scale-remake-gui.json)66個功能／同步檢查點、15圖通過，
包含重複工具列同步。流程為曹操第一劇本、宣戰呂布、編夏侯淵、停手行軍、途中存讀檔、自然入戰、改線、突擊與主動退卻返回。
目視核對[戰術畫面](../../workplace/scratch/ui-windows-20261104/cao-scale-remake-battle-animation-1.png)的城牆、單兵及介面。
依舊在椎陽遭遇呂布，沒有強制入戰，也未驗到夏侯淵守許昌的指定局面。

[獨立規則收據](../../workplace/scratch/ui-windows-20261104/cao-scale-evidence.json)重播相同原版入口與seed，
逐兵27欄位及RNG限前138拍相同；腳本內部狀態限前131拍，後續差異明列保留。
七套件冷測／vet及兩項舊錯版負對照通過，來源見[接線收據](../../workplace/scratch/ui-windows-20261104/cao-scale-integration.json)。
不由這次正常流程抽樣宣稱整場同狀態、架梯、投射物或手機GUI已完成。

## 2026-11-12：曹操對呂布的正常路徑比較

當輪程式為`wlgame-battle-notice-r2`，接通使用者當時選定的三秒提示與按需戰報，現已改為預設關閉，契約見[spec89續片](../spec/89-siege-damage-report-toggle.md)。
[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/battle-notice-remake-gui.json)通過109個功能／同步檢查點、20圖，
包含重複同步，不當作獨立功能。兩次自然遭遇涵蓋提示自行消失、停手日期前進、工具列重開、
原版／HD明細、四種關閉方式、讀檔失效及點提示無地圖穿透；兩種存檔SHA未被介面重寫。
首場仍是夏侯淵攻、呂布守椎陽，不宣稱已跑到夏侯淵守許昌。
前輪圖示驗收有視野缺口，已補正常定位後的可見軍團及原版局部像素，見[playtest50補驗](50-corps-on-map.md)。
[實際畫面](../../workplace/scratch/ui-windows-20261104/battle-notice-game-review.html)與
[來源收據](../../workplace/scratch/ui-windows-20261104/battle-notice-integration.json)留本機。
本片只改介面，未提升前115拍以外的戰術規則或完整結果聲明。

前輪程式為`wlgame-cao-swap-r1`，依[spec62](../spec/62-swapped-unit-skips-its-turn.md)修正坡地換位的StepZ交換與重畫旗標時機。
原版117共同快照重現，467次換位前後記錄符合規則；正式來源前115個可比邊界，96槽27欄位、RNG與腳本狀態吻合。
[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/cao-swap-remake-gui.json)66個功能／同步檢查點、15圖通過，
含曹操對呂布開局、編成、停手行軍、途中存讀檔、自然入戰、改線、突擊與正常主動退卻返回。
[來源收據](../../workplace/scratch/ui-windows-20261104/cao-swap-integration.json)與[本機頁](../../workplace/scratch/ui-windows-20261104/cao-swap-game-review.html)
保存程式與證據；本輪不作高平面導航、投射物、架梯、完整結果或原版同拍像素聲明。以下保留前輪收據。

前輪程式為`wlgame-cao-ground-r1`，依[spec36](../spec/36-ground-planes-and-climbing.md)補回規則地圖表頭列、
Y1–62、斜坡的StepZ增減與低／高平面旗標，並讓動態物件直接用原版投影。
原版117個共同快照重現，第73拍96筆移動完成點及原始地圖／模型吻合，低平面4096bytes全等。
正式來源前84拍，每拍96槽27欄位、完整RNG與腳本狀態相同；第85拍槽73／79的StepZ仍有差異。
[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/cao-ground-remake-gui.json)66個功能／同步檢查點、15圖通過，
涵蓋曹操開局、宣戰、編夏侯淵、停手行軍、途中存讀檔、自然入戰、改線、突擊及主動退卻返回。
[來源收據](../../workplace/scratch/ui-windows-20261104/cao-ground-integration.json)與[本機頁](../../workplace/scratch/ui-windows-20261104/cao-ground-game-review.html)
保存受測程式與證據。畫面為正常流程檢視，沒有原版同拍像素或完整高平面、投射物、架梯、整場結果聲明。
以下保留先前各輪收據。

前輪程式為`wlgame-cao-phase-r1`，補回[spec59](../spec/59-battle-opening-orders.md)的攻守腳本四槽。
原版實際段1，舊Go誤讀段0。新唯讀觀測重現117個共同快照；正式來源的前73拍，
每拍96槽27欄位、完整RNG及腳本游標／等待／條件吻合。第74拍仍有鎖敵、地形及移動差異。
[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/cao-phase-remake-gui.json)通過66個功能／同步檢查點、15圖，
從曹操開局、宣戰呂布、夏侯淵編成、下令徐州、停手行軍與途中存讀檔，自然入戰、改線、突擊及主動退卻返回。
第一次取樣在戰術已退出後仍讀戰場欄位，第二次開場持續前進但超過舊時限；
保留兩份原紀錄、修正腳本競態及驗收時限後，同一binary及拍點條件重跑通過。
[來源收據](../../workplace/scratch/ui-windows-20261104/cao-phase-integration.json)與[本機頁](../../workplace/scratch/ui-windows-20261104/cao-phase-game-review.html)
保存正式程式、來源及限制，不宣稱完整攻城勝負、投射物或架梯已一致。
下方保留先前各輪收據。

前輪程式為`wlgame-cao-path-r1`，依[spec134](../spec/134-pathfind-wave-order.md)補回反向擴散與成本回溯。
同一入戰狀態前40拍，每拍96槽27欄位與完整RNG相同；12筆尋路含11筆成功及1筆失敗均吻合。
新原版唯讀觀測的17個共同快照未變，第4拍槽29差異已修正。
正式GUI重新通過66個功能／同步檢查點、15圖，含改線後停手移動、突擊與正常主動退卻返回，
並核對途中存讀檔。沒有原版整場結果或自然GUI同拍聲明。
[最新收據](../../workplace/scratch/ui-windows-20261104/cao-path-integration.json)及
[本機頁](../../workplace/scratch/ui-windows-20261104/cao-path-game-review.html)保存程式與來源；
全導航／成本表的映射差異保留證據限制，投射物、架梯、後續命令與完整結果仍沿Issues #13／#16／#24。
以下各輪數字與程式為歷史收據。

前輪程式為`wlgame-cao-move-r1`，補回鎖敵Y方向分支及隊長命令7的陣形處理，
來源與原始位址見[spec159](../spec/159-tactical-command-movement-order.md)。
沿用上輪正常原版開局的完整入戰記錄，局部Go比較前3拍96槽的27個欄位全等，40個RNG狀態相同。
第4拍槽29的中繼X仍為原版2、Go5，不由本片關閉Issues #13、#16或#24。

正式GUI另從首頁跑第一劇本曹操，正常宣戰、編夏侯淵、下令徐州、停手行軍及途中存讀檔，
自然遇到呂布後改陣形線。隊長在178拍處於命令7、位置及目標(3,32)，
點中央線後目標成(26,32)，指標固定(300,180)時位置X依序3、4、5、6。
再正常改回原線、突擊，於240拍後點退卻返回戰略地圖；沒有戰術直入或強制勝負。
66個功能／同步檢查點及15張截圖通過，這是主動退卻路徑，不代表完整突擊勝負已與原版相同。
[來源／測試收據](../../workplace/scratch/ui-windows-20261104/cao-move-integration.json)與
[本機實際畫面](../../workplace/scratch/ui-windows-20261104/cao-move-game-review.html)保存目前binary及抽驗範圍。

上一輪續查程式為`wlgame-cao-delay-r1`。已修正兵士開場額外的一拍硬直，
契約與原始定位見[spec200訂正](../spec/200-opening-move-delay.md)。
新原版觀測保留擺位前RNG、初始化及逐拍96筆32-byte記錄；兩次共同230個快照全等。
先前自然開局的兩側兵種、血量與戰力相同，但94槽Y不同、完整RNG不同，
不能把總拍數差異直接解讀為速度錯誤。使用相同原版擺位輸入後，修正前僅96個倒數不同，
修正後9項初始化欄位及首次移動的X／Y、血量、餘力、倒數、步行相位全部吻合。
這份獨立規則比較只隔離開場，未搬入正式遊戲路徑。

正式程式另由正常首頁跑完編成、停手行軍、存讀檔、夏侯淵對呂布自然入戰與戰後返回，
61個功能／同步檢查點及14張截圖通過；沒有戰術直入、注入座標或強制勝負。
第一個腳本錯誤要求喊話消失才下令，第二個在同一小時內首格移動便要求小時改變，
都屬驗收取樣問題。保留兩份失敗收據，改成按既有輸入／子刻契約等待後，同一程式重跑通過。
[本機新畫面](../../workplace/scratch/ui-windows-20261104/cao-delay-game-review.html)及
[完整收據](../../workplace/scratch/ui-windows-20261104/cao-delay-integration.json)保存來源、工具與binary SHA。
自然GUI的完整入戰RNG及下令拍點仍不同，後段結算、投射物及架梯仍不作parity聲明。
下方為前一輪通知修正的歷史收據，不能用其舊binary代表目前程式。

使用者指定第一劇本曹操開局對呂布，檢查操作、動畫及滑鼠靜止後的時間推進。
現行Go程式為本機wlgame-cao-lubu-r2；[接線收據](../../workplace/scratch/ui-windows-20261104/cao-lubu-integration.json)保存50檔來源、二進位與工具SHA。
松崗KI.EXE及SINARIO.DAT沿上方相同雜湊，原版資料唯讀；dosgolem來源提交為f79eab6db8c64a524a06535e53e5d32b0bbcf7c8。
正常原版流程由新遊戲、第一劇本及曹操選取開始，沒有直接呼叫戰鬥、修改日期、軍團或勝負。
兩側固定seed1：原版在首次sub_1ECE0取數前設定既有NewFixed(1)的258-byte狀態；Go使用-seed 1。
這是明示亂數測試條件，不是原版RTC自然播種或整段亂數呼叫次數相同的聲明。

| 比較點 | 原版／remake觀測 |
|---|---|
| 出兵 | 曹操0對呂布13宣戰，夏侯淵37編成6000兵，目的地徐州79；主君編成以正常設定關閉 |
| 停手 | 指標不動時，日期及軍團座標皆前進；Go另以正常小地圖把我軍放入可見區後抽驗，不注入鏡頭或軍團狀態 |
| 首次阻塞 | 4/1 9時、(210,114)出現陳留援軍訊息；原版右鍵可繼續，Go漏接已修。新的多頁／續行測試通過，舊版本被負對照拒絕 |
| 入戰 | 兩側皆4/6 6時，夏侯淵37在(241,108)攻呂布35在(243,108)，據點椎陽76；目標徐州並不表示必須等到徐州才遇敵 |
| 原版重現 | 134步的日期、指令數與結果相同，33張既有截圖SHA相同；加入唯讀節拍觀測後也相同 |
| 走位 | 原版0／4／16／32拍各48兵一側，後續可見移位；Go正常開場、突擊、戰鬥後段與退場有實際截圖 |
| 後段差異 | 原版第932拍到sub_19F58結算入口，Go本次最後活動幀440。突擊下令完成時原版181拍、Go首次生效178拍；完整入戰亂數與單兵狀態未鎖定，不能將差距單獨歸因於速度 |
| 架梯／攀城 | 原版本案可見梯子；Go本輪未取得同一事件的對應畫面，不宣稱其不存在或已等價 |

原版計數器1D318的末值為A4h，會回繞；932來自對該計數器寫入的唯讀累計，不把164當成整場耗時。
原版取樣停在結算寫回入口，未在此輪新增完整戰後世界返回聲明。
七種int61音效呼叫尚未實作，音訊排除於本次對拍完成範圍。
本案沒有關閉Issue #13、#16、#24；下一步須固定入戰狀態，分開比較走位、投射物與HP時序，不為讓畫面接近而直接改速度。

入口：[比較矩陣](../../workplace/scratch/ui-windows-20261104/cao-lubu-comparison.json)、
[本機實際畫面](../../workplace/scratch/ui-windows-20261104/cao-lubu-game-review.html)、
[原版重播指令](../../workplace/scratch/ui-windows-20261104/cao-original-replay-final.commands)、
[原版完整觀測](../../workplace/scratch/ui-windows-20261104/cao-original-replay-observed.log)、
[Go正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/cao-lubu-remake-gui.json)與
[視野內行軍收據](../../workplace/scratch/ui-windows-20261104/cao-lubu-visible-gui.json)。

## 正常操作與前提

兩條路徑均從正常首頁開第一章曹操，remake seed固定1。先由設定關閉主君編成，
編成軍團37，再從進言的敵對提案選袁術12，用「我國較有利」說服君主。
交戰關係實際成立後，從軍團清單點地圖的汝南100 `(216,137)`，選戰鬥指揮。
沒有使用direct-entry、軍團注入、修改時鐘、強制遭遇或強制勝負。

開局對張繡的提案會先被拒絕，因此不把「想攻宛」直接當作已交戰。
本次選袁術，原始鄰接為許昌82→舞陽87→汝南100，途中是自己或目標勢力的領土。
這個資料查核只選擇驗收場景，不代證原版完整行軍。

存檔全部在容器暫存目錄。原版素材唯讀；滑鼠按下／放開與位移使用既有唯讀trace握手。
等待條件採實際軍團位置、節點、路徑欄位、戰鬥攻守側與介面狀態；牆鐘逾時只用來限制腳本。
下表的日期是本次觀測值，不是要求所有即時GUI重播都在同一日期完成輸入。

## 結果

| 分支 | 正常路徑 | 結果 |
|---|---|---|
| 城兵守城 | 宣戰後立即出兵，途中存讀檔，再到汝南 | 63項通過；4/8 6時城池歸屬由12改為0，4/8 8時已到 `(216,137)`、節點100、連結0；沒有駐城軍團，走城兵自動判定 |
| 軍團守城 | 留37在許昌，正常等待敵軍團3到汝南駐守，才下令出兵；途中存讀檔 | 67項通過；4/17 19時軍團37主攻、軍團3守、戰場100、玩家側0；滑鼠突擊後第900幀正常結束，汝南歸玩家，戰後存讀檔保留歸屬與首都82 |

兩條路徑途中存檔均位於 `(202,114)`；原生讀檔保留下面的道路欄位。

| 欄位 | 值 |
|---|---|
| 出發節點 | 82 |
| 連結位址 | `0x0F30` |
| 路徑點 | `0x4348` |
| 步進量 | 4 |
| 指定及目前目標 | 100，座標 `(216,137)` |

重播器的入口檢查點為frame4，完整trace從frame0開始，開場尚未完成；待開場及對話結束才點側欄突擊 `(496,280)`。
trace確認常令2、大將命令2／2。戰後軍團37仍在城門北側 `(216,134)`，未把取得城池說成已站在城中心。
城兵分支另有實際到站檢查，兩者分開驗證。

原版矩形、和平折返及局部更新helper已有各自收據；本篇只新增Go remake正常攻方路徑。
63與67項包含重複的共用操作，不視為130個獨立遊戲情境。

## 收據與重播

本機收據根目錄為 `workplace/shots/attacker-march-20261103/`。

| 路徑 | 內容 |
|---|---|
| `capture/` | 63項城兵分支，含當時的重播器來源、完整trace、JSON及截圖 |
| `battle-attempt2/` | 67項軍團分支，含當時的重播器來源、完整trace、JSON及截圖 |
| `verification.json`、`source-manifest.json` | 驗證範圍、受測程式、來源封存與逐檔收據雜湊 |
| `replay.sh` | 有界Docker重播入口，使用既有image、原版唯讀與獨立空輸出目錄 |

現行工具是 [`tools/desktop_mouse_smoke.py`](../../tools/desktop_mouse_smoke.py)，
模式 `capture` 與 `attack` 都從正常首頁開始。執行器需求與命令見本機 `replay.sh`；
此腳本由Docker啟動Xvfb後呼叫上述工具，不在主機執行Python、遊戲或GUI。

`attack-charge.png` 已目視確認汝南戰術畫面；`attack-post-load.png` 確認回到戰略。
`attack-normal-entry.png` 擷取在PendingBattle首次觀測時，畫面仍為戰略過渡，不能當成戰術畫面證據。
城兵分支的攻陷訊息與返回畫面亦已目視；城池歸屬、到站及道路欄位由trace與存讀檔檢查確認。

## 失敗分類與限制

- `attempt1/` 已正常攻陷只有城兵的汝南，腳本卻要求必有戰術守軍，因此失敗。
  此結果促使測試分開覆蓋城兵與駐城軍團，沒有修改遊戲來配合測試。
- `battle-attempt1/` 改選呂布時，實際提案已直接同意；腳本仍等待理由選單，因此失敗。
  後續依實際回覆分派；軍團分支改回已有自然駐軍軌跡的袁術場景。
### 未解範圍

- 本次沒有原版相同存檔、raw RNG、游標、佔用圖與逐拍碰撞的新對拍收據。
  Issue #30的同狀態閘門仍開著，不以本篇取代原版驗證。
- 不涵蓋完整戰役、其他劇本、全部戰術命令或平台打包，也不重驗音訊及整幅像素一致性。
