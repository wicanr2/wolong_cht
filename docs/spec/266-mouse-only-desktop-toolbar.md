# 266：桌面全滑鼠操作與外部工具列

**狀態：CONFORMED。桌面工具列與一條正常編成／行軍路徑通過；不代表完整 remake 驗收。**

- 日期：2026-11-02。
- 依據：本輪使用者指定參考 `/home/anr2/cht/rich2` 工具列，取消快捷鍵，並明選全滑鼠操作。
- 範圍：桌面 `cmd/wlgame`；HD 與手機版另在 remake 完成後處理。
- 原版輸入：DOS/V `KI.EXE` SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 原版證據沿用 [spec/73](73-right-click-cancel.md)、[spec/149](149-march-target-map-picker.md)、[spec/155](155-desktop-launcher-input.md) 的既有輸入契約；本規格不新增原版語意。

## 使用者裁定的 remake 差異

- 原版遊戲畫面維持 640×400。常駐工具列另占上方 40 像素，遊戲熱區仍使用原版座標。
- 工具列提供語言、設定、存檔、讀檔、返回、確認、離開。HD 按鈕等 HD 完成後才加入。
- 移除所有遊戲鍵盤捷徑，包含方向鍵、Enter、Space、Esc、F9、F10 與速度調整。姓名選字仍可用原版滑鼠字盤；文字輸入可另由文字欄位接收。
- 返回只取消目前的輸入擁有者。確認只交給目前的輸入擁有者，不能在背景另開指令。
- 語言、設定、存讀檔只在允許的穩定畫面啟用。模態流程不排隊執行工具列要求。
- 工具列擁有從按下到放開的完整滑鼠手勢。移出工具列後放開也不得穿透背景。
- 原版的游標移動暫停、停手恢復與持續外推捲動仍依原規格。工具列手勢期間不推進世界，單純停留仍在停手等待後恢復。

## 實作邊界

入口為 [`cmd/wlgame/desktop.go`](../../cmd/wlgame/desktop.go)；包裝既有遊戲的更新與繪圖，不改資料格式、存檔或規則。
按鈕文字來自 `translations/desktop-controls.json`，含四語。對拍 `-shot` 與錄影維持既有 640×400 產物，避免把外殼混入原版畫面比對。

2026-11-11增量：依[spec275頭像首段](275-hd-toolbar-portraits-and-tools.md)加入第八個頭像顯示按鈕，
單排維持上方40像素，按鈕命中與繪製共用八鈕幾何。原七項action編號與遊戲座標保留。
下方30項與舊binary雜湊仍是原七鈕的歷史收據，新八鈕正常操作收據沿spec275，不互相代證。

其後依spec275已確認50%～200%鏡頭加入三個縮放鈕，正式工具列共十一鈕，仍為40像素單排。
可變按鈕寬度、四語文字與命中共用幾何；舊action編號保留，縮放追加編號。
本段正常縮放／點城／存讀檔收據見[spec275第二段](275-hd-toolbar-portraits-and-tools.md)，不改遊戲規則與原版畫布。

前段裁定依[spec275第四段](275-hd-toolbar-portraits-and-tools.md)：工具列展開28高、收起6高，
靠上緣喚出、離開收起，仍為十一鈕單排並放在640×400畫布外。按住與放開幀不收展，
原版文字選單／拖曳持有輸入時不喚出。原七鈕及40高的幾何只保留歷史定位。

2026-11-12依[spec89續片](89-siege-damage-report-toggle.md)的使用者B選擇加入「戰報」，現行為十二鈕單排。
語言／設定縮至48寬，英文用`Lang.`／`Setup`；戰報56寬，按鈕範圍從x12至628。
保留28高展開、6高收起及畫布外置；戰報只有最近一場，明細關閉交回原視窗流程，不新增快捷鍵。
驗證入口沿spec89，不以新增入口代表AI、武將編輯器或整體HD完成。

### 移除額外「右鍵關閉」提示（CONFORMED，限定抽驗，2026-11-12）

使用者指出「右鍵關閉」訊息多餘，要求移除原版沒有的提示。
已查證據點／軍團情報窗旁的提示框由remake額外繪製，原版版面見
[spec23](23-city-info-window.md)、[spec24](24-corps-info-window.md)。
進言結束後另使用同一個remake提示字串，亦屬本次移除範圍。

- 移除據點與軍團情報窗外的額外提示框，以及進言結果的「右鍵關閉」提示。
- 刪除無用的提示座標、共用列舉及四語字串，同步其餘提示的索引。
- 當時保留原版TALK #4「以滑鼠的右鍵回復。」；2026-11-13依使用者新裁定不再顯示，見下方增補。
- 驗收採既有桌面、情報窗、進言與取消冷測，正常滑鼠抽驗情報窗關閉及進言返回；本機截圖不作整款原版parity聲明。

本段先列READY再實作。[冷測、vet及建置](../../workplace/scratch/ui-windows-20261104/no-close-hint-tests.log)通過。
[正常滑鼠抽驗](../../workplace/scratch/ui-windows-20261104/no-close-hint-gui.json)由首頁走新遊戲，
包含原版／HD據點窗、正常編成後的軍團窗、軍團與據點同格選擇、進言結果及正常存讀檔。
69個功能與同步檢查點包含重複工具列同步，不當作69個獨立功能；兩種存檔未被顯示或關窗重寫。
進言說服分支結束提示的移除經程式審查與既有測試，未由本輪正常GUI代證。
[實際畫面頁](../../workplace/scratch/ui-windows-20261104/no-close-hint-game-review.html)附四張修正後截圖及一張前版比較，
[離線頁面抽驗](../../workplace/scratch/ui-windows-20261104/no-close-hint-review-qa.json)核對五圖、窄視窗、JS錯誤及外部請求。
[接線收據](../../workplace/scratch/ui-windows-20261104/no-close-hint-integration.json)保存45檔來源覆蓋包、受測程式及截圖SHA；
流程與建頁工具分別為[滑鼠工具](../../workplace/scratch/ui-windows-20261104/no-close-hint-gui.py)、
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-no-close-hint-review.py)及
[頁面工具](../../workplace/scratch/ui-windows-20261104/qa-no-close-hint-review.cjs)。

### 編成與存讀檔成功不另開通知（CONFORMED，限定抽驗，2026-11-12）

使用者指出「夏侯淵編成完畢」「已儲存第一槽」是額外訊息，要求檢視必要性與玩家體驗。
已查證兩者同由lastEvent額外事件框呈現，不阻塞輸入，但遮住地圖六秒；沿[spec88](88-display-polish-parity.md)的remake差異。
編成成功已有[原版主將肖像對白](109-formation-leader-line.md)，存讀檔成功亦會關閉選槽窗。

- 編成、存檔及讀檔成功不再另設lastEvent成功通知，並清除舊通知，避免成功後仍顯示上次失敗。
- 保留編成後原版對白、部隊狀態、選槽關窗與存檔內容，作為既有結果回饋；不新增提示音或工具列動畫。
- 存讀檔失敗仍留在選槽窗，沿既有錯誤提示允許重試；編成失敗保留原版TALK #62。
- 其他遊戲事件與錯誤通知維持既有路徑，不以刪除成功通知代替完整事件稽核。
- 驗收採既有編成／存檔／桌面冷測，並抽驗正常編成原版對白、成功關窗與存讀檔；不擴大原版parity聲明。

本段先列READY再實作。[冷測、vet與建置](../../workplace/scratch/ui-windows-20261104/quiet-success-tests.log)通過。
[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/quiet-success-gui.json)通過50個功能與同步檢查點，包含重複工具列同步，
不當作50個獨立功能。原版編成主將對白保留，成功編成／存檔／讀檔後無額外通知；兩種存檔存在且正常讀回，未被顯示操作重寫。
失敗處理依既有測試與程式審查保留，未由本輪正常GUI代證全部儲存裝置錯誤。
[本機七圖比較頁](../../workplace/scratch/ui-windows-20261104/quiet-success-game-review.html)附五張本輪畫面及兩張前版比較。
入口為[流程工具](../../workplace/scratch/ui-windows-20261104/quiet-success-gui.py)、
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-quiet-success-review.py)、
[頁面抽驗工具](../../workplace/scratch/ui-windows-20261104/qa-quiet-success-review.cjs)及
[接線收據](../../workplace/scratch/ui-windows-20261104/quiet-success-integration.json)，保留45檔來源覆蓋包、程式、測試與截圖SHA。
[離線抽驗](../../workplace/scratch/ui-windows-20261104/quiet-success-review-qa.json)通過七圖載入及窄視窗，JS錯誤與外部請求皆0。

### 下令及已有對白的結算不再重複通知（CONFORMED，限定抽驗，2026-11-12）

接續已確認的成功通知清理。已查證行軍與戰術下令額外寫入lastEvent；外交及撥款亦在原版TALK結果之外另寫結算摘要。
原版操作依[行軍選單](39-march-order-menu.md)、[戰術下令](33-squad-selection.md)及
[事件說話順序](42-event-scene-speakers.md)，本段不改命令或結算規則。

- 行軍與戰術命令成功時清除額外事件列，依部隊狀態及動作呈現結果；下令失敗保留原處理。
- 外交成功或玩家拒絕，保留原版對白並移除重複摘要；規則回傳未成立時保留失敗提醒。
- 撥款成功、玩家拒絕與已知的零金額對白分支，保留TALK並移除額外摘要；其他未成立狀態保留提醒。
- 不清空原版訊息佇列，不移除宣戰、戰鬥結果、災害等重要事件通知，不改遊戲狀態或存檔。
- 受影響冷測及正常曹操開局的編成／行軍抽驗；外交／撥款結算以既有規則與對白測試核對，不冒稱完整正常事件GUI通過。

受影響冷測、vet與正常曹操對呂布流程通過；收據沿下段。原版TALK與規則呼叫保留，只移除額外摘要。

### 曹操對呂布的訊息右鍵與停手前進（CONFORMED，限定輸入修正與正常流程抽驗，2026-11-12）

使用者指定完成介面驗證後，由第一劇本曹操開局對呂布，比較正常操作、動畫及滑鼠靜止時的時間推進。
原版輸入沿[re111](../re/111-c-input-mouse-restoration.md)：訊息等待常式保留右鍵分支，
原`0x12239`的`EB 0F`只略過左鍵詢問。滑鼠恢復時間沿[spec112](112-cursor-idle-resume-delay.md)。

固定seed1的正常曹操出兵已重現：196/4/1 9時、軍團37在(210,114)停在一則訊息，反覆右鍵未推進。
`updateMessageOnly`只接受左鍵或確認，漏接右鍵；本段補回右鍵推進目前一頁，不取消整條流程或一次清空佇列。
最後一頁才執行續行，訊息關閉後仍沿既有idleGate.Pause等待停手，再讓世界前進。
保留左鍵與工具列確認；不改原版對白、亂數、行軍路線或動畫速度。
驗收包含多頁訊息與續行、無輸入不跳頁，以及相同seed／正常曹操對呂布路徑重跑；動畫逐項記錄觀測範圍，不預先宣稱全部一致。

[訊息分頁回歸](../../cmd/wlgame/messages_test.go)驗證右鍵逐頁與最後一頁續行，
[舊行為負對照](../../workplace/scratch/ui-windows-20261104/message-right-negative.log)會拒絕原本只接受左鍵的版本。
[受影響冷測、vet與建置](../../workplace/scratch/ui-windows-20261104/cao-lubu-tests-r2.log)通過。
[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/cao-lubu-remake-gui.json)通過60個功能／同步檢查點及14張截圖，
包含重複同步，不當作60個獨立功能；完成曹操宣戰呂布、夏侯淵編成、徐州下令、停手行軍、途中存讀檔與椎陽戰鬥退場。
[可見我軍補充抽驗](../../workplace/scratch/ui-windows-20261104/cao-lubu-visible-gui.json)以正常小地圖移到我軍，證明同一靜止指標下日期及可見軍團都前進。

[原版／remake比較](../../workplace/scratch/ui-windows-20261104/cao-lubu-comparison.json)記錄相同4/6 6時、37對35、椎陽76與兩軍座標；
原版134步與33張截圖重播相同。戰術後段有差異：原版932拍到結算入口，本次remake約440拍退場。
完整入戰RNG、單兵狀態與精確下令拍未鎖定，不把這個差距直接當成動畫速度錯誤，也不宣稱全部動畫已一致。
完整案例入口為[playtest125](../playtest/125-mouse-attacker-march.md)，
[本機比較頁](../../workplace/scratch/ui-windows-20261104/cao-lubu-game-review.html)提供兩側實際畫面，
[頁面抽驗](../../workplace/scratch/ui-windows-20261104/cao-lubu-review-qa.json)驗證離線切換。
[接線收據](../../workplace/scratch/ui-windows-20261104/cao-lubu-integration.json)固定50檔來源覆蓋包、受測程式及工具身分。
[收尾核對](../../workplace/scratch/ui-windows-20261104/cao-lubu-final-checks.json)核對來源、原版重播、頁面、素材與擁有權。

## 2026-11-13：省略右鍵操作提示（CONFORMED，限本片）

使用者指出「以滑鼠右鍵恢復」訊息不必要。已核對為原版TALK #4，
由`openMinimapFactionPicker`與`openCorpsInfo`掛左下角狀態框，原版證據在[re31](../re/31-faction-picker-screen.md) §1.3與[spec140](140-status-message-box.md)接線表。
本裁定取代先前「保留TALK #4顯示」的選擇，明列為使用者指定的remake差異。
勢力選擇與軍團情報開啟時清除舊狀態框，不顯示右鍵操作教學；保留TALK原文、資料索引、選擇規則與右鍵返回。
其他指令的任務提示不受此裁定影響。驗收調整既有狀態框測試，正常滑鼠抽驗兩個入口與右鍵關閉。

先列READY再實作。[受影響冷測與vet](../../workplace/scratch/ui-windows-20261104/quiet-notice-tests-r2-20261113.log)通過。
最終`wlgame-quiet-notice-r2`的[正常滑鼠抽驗](../../workplace/scratch/ui-windows-20261104/quiet-controls-gui-20261113.json)
通過44個功能／同步檢查點、三張實際截圖，包含重複工具列同步，不當作獨立功能數。
正常勢力選擇、編成後軍團情報與兩處右鍵返回已目視確認，左下角沒有教學框；其他任務提示沿原流程。
來源、受測程式與收據身分見[本輪收據](../../workplace/scratch/ui-windows-20261104/quiet-notice-integration-20261113.json)。

## 驗收

- 單元測試驗證工具列按鈕、停用條件、完整手勢隔離、返回一層及確認路由。
- 冷跑桌面套件與既有規則測試，實跑正常啟動、編成、行軍、存讀檔及返回。
- 正常 GUI 驗證所有鍵盤捷徑失效，工具列沒有覆蓋遊戲、沒有點擊穿透。
- 軍團行軍另依 [spec/74](74-corps-on-world-map.md) 驗證位置與朝向。不得把全程動畫或原版 parity 當作工具列測試的結論。

## 本輪驗證與限制

[`tools/desktop_mouse_smoke.py`](../../tools/desktop_mouse_smoke.py) 以真實 X11 輸入走正常首頁、四語切換、開新遊戲、設定、存讀檔、編成、行軍與離開取消，共 30 項通過。
使用 Go 1.26.7、Ebiten 2.9.9 與既有 `san1-go:1.26.7-ebiten2.9.9-r1`，Docker／Xvfb、UID/GID 1000:1000、固定 seed 1。
沒有 `-direct`、強制行軍、狀態注入或自動關對話旗標。

本機收據：`workplace/scratch/mouse-toolbar-gui.json`，二進位 SHA-256
`fa362568d528734a10205a474e8958e1c3556a402b35f9596ae78f6b311b0cc8`。
軍團 16 正常編成後下令前往宛，固定鏡頭 `(188,104)`，位置從 `(206,114)` 到 `(204,114)`，
朝向從 4 到 0；`mouse-toolbar-march-start.png` 與 `mouse-toolbar-marching.png` 已目視確認軍團在可見地圖內移動。
這是 remake 正常路徑的移動證據，不是原版同狀態、動畫幀或時間 parity。

正式範圍的測試套件冷跑與 `go vet` 通過。重播失敗與清理收據見 [`WORKLOG.md`](../../WORKLOG.md)。

## 未解範圍

| 項目 | 證據限制 |
|---|---|
| 原版彈出選單手感 | 2026-11-03 使用者選定原版上下推動反白，實作與限定驗收已由 [spec/267](267-original-choice-band.md) 接續完成；上方 30 項是修改前的歷史收據 |
| 戰術、完整戰役與原生平台 | 此次正常 GUI 抽驗不代證完整戰鬥、Windows／macOS 實機或整體 remake 完成 |
| HD／手機 | 使用者指定 remake 完成後依序製作；本輪沒有進行這兩條分支 |
