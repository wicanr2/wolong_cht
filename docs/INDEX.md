# 文件索引（自動產生，不要手改）

由 `tools/index.py generate` 從 `docs/**/*.md` 生出來。
**狀態與日期是從各文件的內文讀的，不是另外維護的一份摘要**——
手寫的索引會過時，而過時的索引比沒有索引更糟（`docs/playtest/04` 的「四件事被擋住」整張表在寫下的當下就是錯的）。

提交前跑 `tools/index.py check`，它會擋下：狀態行與內文矛盾、
同一斷言在兩份文件等級不同、連結壞掉、`CONTEXT.md` 漏登記。

## 文件

| 文件 | 標題 | 狀態 | 日期 |
|---|---|---|---|
| [`docs/formats/01-talk-dat.md`](formats/01-talk-dat.md) | 01 — TALK.DAT 訊息表格式 | READY | 2026-08-07 |
| [`docs/formats/02-brg-palette.md`](formats/02-brg-palette.md) | 02 — .BRG 調色盤格式 | READY | 2026-08-07 |
| [`docs/formats/03-grf-images.md`](formats/03-grf-images.md) | 03 — GRF.DAT 圖庫格式 | KAOGRF／KYOGRF／IVENTGRF READY，ICONGRF 部分解。 | 2026-08-13 |
| [`docs/formats/04-map-sch-container.md`](formats/04-map-sch-container.md) | 04 — .MAP／.SCH／.MCH：兩種完全不同的東西 | 容器格式的索引層 READY，壓縮演算法未解。 | 2026-08-07 |
| [`docs/formats/05-mmap-worldmap.md`](formats/05-mmap-worldmap.md) | 05 — MMAP. 大地圖 | MMAP.MDL、地圖尺寸、自動連接與 MMAP.MCH 物件圖形入口 confirmed；MMAP.MAP 的 R… | 2026-08-07 |
| [`docs/formats/06-mmap-rle.md`](formats/06-mmap-rle.md) | 06 — MMAP.MAP 的 RLE 壓縮 | READY。 | 2026-08-07 |
| [`docs/formats/07-battle.md`](formats/07-battle.md) | 07 — BATTLE. 戰場資料 | 分段結構、圖塊定義與像素格式都 confirmed。 | 2026-08-07 |
| [`docs/formats/08-sinario-save.md`](formats/08-sinario-save.md) | 08 — SINARIO.DAT / SAVE.DAT：劇本與存檔 | 整體結構 confirmed，武將能力值 confirmed，其餘欄位進行中。 | 2026-08-07 |
| [`docs/mechanics/00-index.md`](mechanics/00-index.md) | 00 — 遊戲機制索引 | 索引與推論等級定義，長期有效。 | 2026-08-08 |
| [`docs/mechanics/10-strategy.md`](mechanics/10-strategy.md) | 10 — 大地圖政略 | 指令清單完整；部分戰略數值與 AI 決策已由機器碼解出並實作，仍有未解公式。 | 2026-08-13 |
| [`docs/mechanics/15-realtime.md`](mechanics/15-realtime.md) | 15 — 即時制的時間模型 | ✅ READY。整條時間鏈已在機器碼裡讀出來。 | 2026-08-08 |
| [`docs/mechanics/20-military.md`](mechanics/20-military.md) | 20 — 行軍與軍團 | 道路網與行軍已解並實作；AI 與玩家的編成流程都已解； 六個位置怎麼換算成戰力仍未完。 | 2026-08-13 |
| [`docs/mechanics/30-combat.md`](mechanics/30-combat.md) | 30 — 戰場（戰術） | 進場規則與戰略層的自動判定全解；戰術核心已接入並實作，完整結算與少數分支未完 | 2026-08-09 |
| [`docs/mechanics/40-economy.md`](mechanics/40-economy.md) | 40 — 經濟：資金與預備兵 | 機制與公式都已解 | 2026-08-08 |
| [`docs/mechanics/50-diplomacy.md`](mechanics/50-diplomacy.md) | 50 — 外交 | 成立條件與外交官的數值都已解 | 2026-08-13 |
| [`docs/mechanics/60-personnel.md`](mechanics/60-personnel.md) | 60 — 武將 | 三個能力值的作用與身分欄位已定案（說明書＋機器碼），數值公式部分已知。 政治如何影響內政效果與外交官要價仍未解。 | 2026-08-13 |
| [`docs/mechanics/70-ai.md`](mechanics/70-ai.md) | 70 — 電腦 AI 的判斷邏輯 | 侵攻目標的決策鏈、友好度漂移、宣戰三閘與 AI 編成入口已由機器碼讀出；remake 已接上可重播的敵方出兵切片。 | 2026-08-09 |
| [`docs/mechanics/80-victory.md`](mechanics/80-victory.md) | 80 — 勝負判定 | 勝負條件全部定案（說明書 1.1、1.2 節）。 | 2026-08-08 |
| [`docs/mobile/android-plan.md`](mobile/android-plan.md) | Android 版規劃 | 規劃已啟動；觸控 shell 原型 debug APK 已產出，並已在 Android 模擬器完成安裝、啟動與有限觸… | 2026-08-11 |
| [`docs/playtest/01-dosbox-dosv.md`](playtest/01-dosbox-dosv.md) | 01 — 松崗 DOS/V 版首次實跑（DOSBox） | DOS/V 版可開機並可由密碼頁進入開場；字型懸案結案。 | 2026-08-07 |
| [`docs/playtest/02-dosboxx-pc98.md`](playtest/02-dosboxx-pc98.md) | 02 — PC-98 日文原版實跑（DOSBox-X）：oracle 建立完成 | PC-98 oracle 已建立，沒有防拷。 | 2026-08-07 |
| [`docs/playtest/03-verification-log.md`](playtest/03-verification-log.md) | 03 — 補驗紀錄：哪些結論撐得住，哪些撐不住 | 持續累加。專記驗證與失敗的驗證嘗試。 | 2026-08-07 |
| [`docs/playtest/04-mouse-automation-blocked.md`](playtest/04-mouse-automation-blocked.md) | 04 — 受阻：PC-98 oracle 動不了遊戲內的滑鼠游標 | ⛔ 已作廢——本文的兩個核心判斷都被推翻，解法見 docs/playtest/06。 | 2026-08-07 |
| [`docs/playtest/05-viewport-and-city-coords.md`](playtest/05-viewport-and-city-coords.md) | 05 — 一張靜止的截圖就把據點座標定案了 | SINARIO.DAT 的據點 X／Y 座標升到 confirmed。 | 2026-08-08 |
| [`docs/playtest/06-mouse-solved.md`](playtest/06-mouse-solved.md) | 06 — PC-98 oracle 的滑鼠通了 | 可以自動操作遊戲。 | 2026-08-08 |
| [`docs/playtest/07-in-game-oracle.md`](playtest/07-in-game-oracle.md) | 07 — 進到遊戲本體：三個欄位靠實機畫面定案 | 可自動玩進遊戲本體；首都與頭像編號兩欄位升 confirmed。 | 2026-08-08 |
| [`docs/playtest/08-wlgame-normal-strategy-path.md`](playtest/08-wlgame-normal-strategy-path.md) | 08 — wlgame 正常編成／行軍路徑 | 真實 SINARIO.DAT 下的編成、行軍、城兵攻城，以及敵方 AI 軍團遭遇選單 都已由鍵盤正常操作驗收；使用者… | 2026-08-09 |
| [`docs/playtest/09-wlgame-normal-tactical-path.md`](playtest/09-wlgame-normal-tactical-path.md) | 09 — wlgame 正常遭遇到戰術戰場 | 正常開局、正常 AI 遭遇、戰鬥指揮選單、戰術戰場、戰後結果報告與 GUI 戰後回戰略已驗收；原版同狀態逐像素對拍仍… | 2026-08-09 |
| [`docs/playtest/10-event-message-modal.md`](playtest/10-event-message-modal.md) | 10 — 事件 TALK 通知 modal | 已完成通知資料接縫、Linux/Xvfb 視覺抽樣與 remake TALK 五行分頁；原版未知 事件流程／未定位 … | 2026-08-09 |
| [`docs/playtest/11-event6-original-fixture.md`](playtest/11-event6-original-fixture.md) | 11 — 原版事件 6 fixture oracle | 事件 6 主要結果畫面已由原版 fixture 證實；不是自然長程存檔，也不封閉次要 formatter。 | 2026-08-10 |
| [`docs/playtest/12-event3-same-state-parity.md`](playtest/12-event3-same-state-parity.md) | event3 同狀態對拍（2026-08-10） | 同狀態事件接線、原版式 composite、3×6 實際格位選取、TALK 五行分頁、 DOS/V 內框／按鍵 gl… | 2026-08-10 |
| [`docs/playtest/13-dosv-natural-and-target-gui.md`](playtest/13-dosv-natural-and-target-gui.md) | DOS/V 自然畫面與目標平台 GUI 驗收 | 影片 oracle 對拍通過；嚴格同狀態逐像素差異仍不宣稱，Windows／macOS 原生 runtime 仍受阻。 | 2026-08-11 |
| [`docs/playtest/14-m7-review.md`](playtest/14-m7-review.md) | M7 校訂文字人工審查報告 | 60 筆定案校訂已完成逐筆語意、marker、硬換行、寬度與代表畫面抽樣；不宣稱 1,022 則全文重譯或密碼保護下… | 2026-08-11 |
| [`docs/playtest/15-event2-5-talk-sampling.md`](playtest/15-event2-5-talk-sampling.md) | 事件 2–5 TALK 完整分支抽樣 | 36 個 raw TALK 頁面、18 組雙頁回應的分支、marker、硬換行、字寬與五列版面 抽樣通過；不宣稱完整… | 2026-08-11 |
| [`docs/playtest/16-event9-long-route.md`](playtest/16-event9-long-route.md) | 事件 9 長程通知流程 | 27 小時 bounded queue、玩家／非玩家／在野通知條件與 #409 no-op 已通過； 完整自然劇本依… | 2026-08-11 |
| [`docs/playtest/17-expert-dosbox-remake.md`](playtest/17-expert-dosbox-remake.md) | 17 — DOSBox 原版／remake 可玩性專家驗證 | remake 正常策略路徑與存檔／讀檔通過；DOS/V 原版密碼頁已可進入開場，尚未展開完整自然長程驗證；PC-98… | 2026-08-11 |
| [`docs/playtest/18-dosv-password-verification.md`](playtest/18-dosv-password-verification.md) | 18 — 松崗 DOS/V 密碼頁輸入驗證 | 已證實，在受控 DOSBox-X 重播中按「確定」即可越過密碼頁；密碼頁不再是 DOS/V 原版行為驗證的阻擋。 | 2026-08-12 |
| [`docs/playtest/19-tactical-minimap.md`](playtest/19-tactical-minimap.md) | 19 — DOS/V 戰術縮圖 raw producer 驗收 | PASS（已證實 producer 的 remake 實作）；局部更新與原版精確外框素材仍為 unknown。 | 2026-08-12 |
| [`docs/playtest/20-tactical-layout-parity.md`](playtest/20-tactical-layout-parity.md) | 20 — 松崗 DOS/V 戰術版面 parity 重開 | PARTIAL（主要幾何、右欄命令面板、底列 glyph、原版初始相機、 32×30 display grid、鄰格… | 2026-08-12 |
| [`docs/playtest/21-command-window-parity.md`](playtest/21-command-window-parity.md) | 21 — 松崗 DOS/V 指揮／事件／一覽畫面 parity 重開 | PARTIAL（事件 TALK、系統面板與一覽第一層主要幾何已修正；一覽詳細層與捲軸未完成）。 | 2026-08-12 |
| [`docs/playtest/22-field-siege-shared-layout.md`](playtest/22-field-siege-shared-layout.md) | 攻城／兩軍遭遇共用戰術骨架驗收 | PASS（共用幾何與原版指令面板已封口；不代表動畫逐像素 parity） | 2026-08-12 |
| [`docs/promo/README.md`](promo/README.md) | 推廣片產出紀錄 | remake 推廣片、代表幀「經典再現」比較片與 DOS/V／remake 實機動態比較片均已產出並完成媒體規格驗證… | 2026-08-12 |
| [`docs/promo/classic-revival.md`](promo/classic-revival.md) | 「經典再現」推廣片 | 已產出研究／推廣用 60 秒比較片；不把代表幀比較宣稱為同狀態逐像素 parity。 | 2026-08-11 |
| [`docs/promo/dosv-adlib-and-tactical-review.md`](promo/dosv-adlib-and-tactical-review.md) | 推廣片原版 AdLib 與戰術骨架審查 | 原版音軌來源與影片縮放鏈已修正；戰術固定骨架的 16 px viewport 偏差已修正。完整同狀態戰術 parit… | 2026-08-12 |
| [`docs/promo/dosv-live-comparison.md`](promo/dosv-live-comparison.md) | DOS/V 原版／remake 實機動態比較推廣片 | 已產出並完成畫面、媒體規格與來源界線驗收；這是同類型畫面／流程的推廣比較，不是同日期、同輸入、同狀態的逐像素 par… | 2026-08-12 |
| [`docs/promo/yt-remake-pixel-review.md`](promo/yt-remake-pixel-review.md) | YouTube／remake 推廣片像素差異審查 | 已完成影片對照；確認像素差異，但不把不同遊戲狀態誤稱為同狀態 逐像素 parity。 | 2026-08-11 |
| [`docs/re/01-first-recon.md`](re/01-first-recon.md) | 01 — 首輪偵查：兩版檔案清單與比對 | READY。兩版檔案清單、執行結構、逐檔比對都完成。 | 2026-08-07 |
| [`docs/re/02-palette-routine.md`](re/02-palette-routine.md) | 02 — 調色盤常式：.BRG 的通道順序與亮度縮放 | READY。兩版調色盤常式互證，.BRG 格式定案。 | 2026-08-07 |
| [`docs/re/03-image-blitter.md`](re/03-image-blitter.md) | 03 — 圖庫載入器與 VGA 繪製常式 | READY。圖庫載入器與四平面繪製常式都讀完。 | 2026-08-07 |
| [`docs/re/04-mmap-entry-points.md`](re/04-mmap-entry-points.md) | 04 — 大地圖：入口點與記憶體佈局 | 入口、尺寸、自動連接與 MMAP.MCH 物件圖形入口已定案。 | 2026-08-07 |
| [`docs/re/05-battle-selection.md`](re/05-battle-selection.md) | 05 — 政略與戰術的接縫：戰場怎麼被選出來 | 政略↔戰術的接縫已解；戰場記錄的部分欄位未解。 | 2026-08-07 |
| [`docs/re/06-game-clock.md`](re/06-game-clock.md) | 06 — 遊戲時鐘：sub_11D8E | 整條時間鏈已解，confirmed。 | 2026-08-08 |
| [`docs/re/07-monthly-settlement.md`](re/07-monthly-settlement.md) | 07 — 月結：sub_15358 與它的九支子程式 | 經濟公式全部讀出來，confirmed。 | 2026-08-08 |
| [`docs/re/08-hourly-update.md`](re/08-hourly-update.md) | 08 — 每「時」的世界更新：sub_13E11 | 主結構已解。軍團編成、外交官效果、事件分派都在這裡。 | 2026-08-08 |
| [`docs/re/09-combat.md`](re/09-combat.md) | 09 — 戰鬥：觸發、自動判定、傷亡與武將的下場 | 戰略層自動判定與政略↔戰術入口已解；戰術本體的核心規則已在 [docs/re/11](11-tactical-bat… | 2026-08-09 |
| [`docs/re/10-rng.md`](re/10-rng.md) | 10 — 亂數產生器：sub_1ECE0 與 sub_1EC82 | 全解，已實作成 internal/rules/rng。 | 2026-08-08 |
| [`docs/re/11-tactical-battle.md`](re/11-tactical-battle.md) | 11 — 戰術戰鬥：模組結構與戰場資料模型 | 模組骨架、戰場資料模型、核心移動／命中／傷害規則已大致解出並接入測試； 正常玩家已可由遭遇選單進入攻城戰術畫面並送出… | 2026-08-09 |
| [`docs/re/12-diplomacy-dialogue.md`](re/12-diplomacy-dialogue.md) | 12 — 停戰說服訊息索引：#190–#198 | 三變體槽位與停戰說服這條索引路徑已證實；事件 6／7 的次要呼叫已定位， 但 formatter 參數契約與完整可見… | 2026-08-09 |
| [`docs/re/13-pc98-numeric-window.md`](re/13-pc98-numeric-window.md) | 13 — DOS/V 數字輸入視窗量測與 CJK 版面決策 | 以 DOS/V 為唯一畫面基準；原始座標、DOS/V 96×64 內框 blit、3×6 每格操作、 實際按鍵 gl… | 2026-08-10 |
| [`docs/re/14-mmap-mch-objects.md`](re/14-mmap-mch-objects.md) | 14 — MMAP.MCH 戰略地圖物件 | 資產格式、事件 12 的火災／暴動圖形鏈與 typed 動畫／移動時序 confirmed； type 3 的事件語… | 2026-08-10 |
| [`docs/re/15-event10-producer.md`](re/15-event10-producer.md) | 15 — 事件 10 producer 深度逆向 | 事件 10 dispatcher／consumer／queue writer 已證實；原版自然 producer 仍… | 2026-08-11 |
| [`docs/re/16-idle-clock-event10.md`](re/16-idle-clock-event10.md) | 16 — DOS/V 無輸入自動時鐘與事件 10 關係 | 無輸入時的自動時鐘／軍團行軍已由 IDA .i64 證實；事件 10 是該路徑 中的受節流 queue consum… | 2026-08-11 |
| [`docs/re/17-dosv-audio-tsr.md`](re/17-dosv-audio-tsr.md) | 17 — 松崗 DOS/V 音源 TSR 與戰術效果碼 | INT 61h 介面、遊戲端效果碼與硬體 register 寫入已證實；BGM.DAT 的容器與 聲軌指標結構已解（… | 2026-08-12 |
| [`docs/re/18-tactical-button-glyphs.md`](re/18-tactical-button-glyphs.md) | 18 — DOS/V 戰術底列按鈕 glyph 候選資產研究 | 六個命令 glyph、底板、右欄複合面板與選取矩形已解出並接入 remake。 | 2026-08-12 |
| [`docs/re/19-outcome.md`](re/19-outcome.md) | DOS/V 已證實敗北 outcome 接線 | READY（只涵蓋兩種敗北；不涵蓋勝利、君主死亡或原版返回標題）。 | 2026-08-12 |
| [`docs/re/20-ida-re-coverage-audit.md`](re/20-ida-re-coverage-audit.md) | 20 — DOS/V IDA 逆向覆蓋與 remake 差距審計 | REVIEWED。足以支撐可玩重製與多數規則，但不足以支撐高忠實度戰術呈現；主要缺口是原版顯示串列、相機狀態機、逐幀… | 2026-08-12 |
| [`docs/re/21-function-census.md`](re/21-function-census.md) | 21 — DOS/V KI.EXE 全函式覆蓋普查 | 量測完成。739 支函式中 252 支（34%，佔程式碼 24%）從未被讀過； 缺口集中在戰略 UI、戰術戰鬥與月結… | 2026-08-14 |
| [`docs/re/22-strategy-command-tree.md`](re/22-strategy-command-tree.md) | 22 — 戰略指令列與滑鼠熱區分派 | 指令列的八個槽位、兩層子選單、財政四子項、編成入口、熱區查表定址式與 狀態列提示常式 confirmed。指令列圖形… | 2026-08-13 |
| [`docs/re/23-bgm-resource-format.md`](re/23-bgm-resource-format.md) | 23 — BGM.DAT 音樂資源格式 | 容器索引與曲塊的聲軌指標表 confirmed，兩版都逐 byte 對齊到檔尾且餘 0。 聲軌資料本身的事件編碼、+… | 2026-08-13 |
| [`docs/re/24-unread-function-catalogue.md`](re/24-unread-function-catalogue.md) | 24 — 未讀函式目錄：252 支的證據與下手順序 | 252 支未讀函式全部登記，附上由共用常式與 TALK 訊息取得的角色證據。 計數要排除本檔自己（§5），否則登記等… | 2026-08-13 |
| [`docs/re/25-message-variants-and-personnel.md`](re/25-message-variants-and-personnel.md) | 25 — 訊息變體展開與人事指令 | 訊息索引 ≥ 0x196 的 ×8 變體展開 confirmed，13 個呼叫值共 104 則變體已對出； 人事四支… | 2026-08-13 |
| [`docs/re/26-list-window-engine.md`](re/26-list-window-engine.md) | 26 — 一覽表視窗引擎 | 視窗幾何、五個一覽表家族的描述子、選取迴圈、持久化排序狀態與四個描述子欄位 （兩個 callback ＋ 標題字串 … | 2026-08-13 |
| [`docs/re/27-list-row-fields.md`](re/27-list-row-fields.md) | 27 — 一覽表的逐列繪製與外交關係等級 | 四個家族的逐列常式全部讀完，欄位對照 confirmed；兵力的 ×10 顯示、 三條換色規則與外交關係六級換算 c… | 2026-08-13 |
| [`docs/re/28-text-number-rendering.md`](re/28-text-number-rendering.md) | 28 — 文字與數字的繪製層 | 數字繪製、兩支字串繪製與 EGA 平面寫入方式 confirmed。 單字元 blitter loc_1F75E 與… | 2026-08-13 |
| [`docs/re/29-font-service-int15.md`](re/29-font-service-int15.md) | 29 — 原版怎麼顯示中文：INT 15h 字型服務與 END_S13/S14.DAT | 整條鏈 confirmed（靜態）。KI.EXE 側走 DOS/V 的 INT 15h AH=50h 向常駐服務要字… | 2026-08-13 |
| [`docs/re/30-corps-formation-ui.md`](re/30-corps-formation-ui.md) | 30 — 軍團編成畫面：兵員池、兵種切換與派生數值 | 編成畫面的主迴圈與被它呼叫的常式 confirmed。 兵員池的搬運方向、兵種切換的循環、兵力的自動分配、確定的前提… | 2026-08-13 |
| [`docs/re/31-faction-picker-screen.md`](re/31-faction-picker-screen.md) | 31 — 勢力一覽：22 格的兩欄版面與領地重繪 | 版面、命中判定、顏色規則、兩個守門條件與介面字串 confirmed。 分派表已印出，但 sub_15AD1 → s… | 2026-08-14 |
| [`docs/re/32-strategy-detail-panels.md`](re/32-strategy-detail-panels.md) | 32 — 戰略側的兩個詳情面板：據點與軍團 | 兩支面板函式與軍團面板的外層 confirmed，欄位與字串表全部對出。 畫肖像的 sub_107D2 與兩支收尾常… | 2026-08-14 |
| [`docs/re/33-shared-draw-helpers.md`](re/33-shared-draw-helpers.md) | 33 — 共用繪圖層：字串包裝、肖像快取、小地圖上色 | 五支共用常式 confirmed。肖像快取的替換策略與小地圖的座標換算定案。 實際搬 bytes 的 sub_1E3… | 2026-08-14 |
| [`docs/reference/01-jp-manual.md`](reference/01-jp-manual.md) | 01 — 日文原版說明書判讀紀錄 | 有實質機制的頁都讀完了，剩 p.6 啟動操作與 p.36–38 附錄。 | 2026-08-08 |
| [`docs/reference/02-jp-cht-diff.md`](reference/02-jp-cht-diff.md) | 02 — 日中對照：TALK.DAT 第一批發現 | 全量 1,022 則的 | 2026-08-13 |
| [`docs/reference/03-baked-japanese.md`](reference/03-baked-japanese.md) | 03 — 燒進美術裡的日文：松崗版沒重繪的部分 | 已確認的缺口：標題橫幅「臥竜伝」兩版相同（松崗沒重繪）。 | 2026-08-07 |
| [`docs/reference/04-first-survey.md`](reference/04-first-survey.md) | 04 — 首輪偵查紀錄（2026-08-07） | 歷史快照。已解的項目由 docs/formats/、docs/re/ 與 docs/INDEX.md 接手；本檔只保… | 2026-08-07 |
| [`docs/reference/05-eten-font-provenance.md`](reference/05-eten-font-provenance.md) | 05 — 松崗版的中文字型是倚天字型 | confirmed。END_S14.DAT 與倚天 ascfont.15 byte-for-byte 相同； END… | 2026-08-13 |
| [`docs/release/README-RELEASE.md`](release/README-RELEASE.md) | 臥龍傳 remake 可執行封裝 | 三平台候選封裝、Linux AppImage、推廣片與驗收紀錄已集中於 [dist-all](../../dist-… | 2026-08-12 |

## 斷言（欄位／常數 → 推論等級 → 出處）

共 145 條。**要查「這件事解了沒」先看這裡**，
不要重讀整份文件，更不要重推一次。

### confirmed（76 條）

| 鍵 | 出處 |
|---|---|
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 1 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 2 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 8 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. 戰場怎麼被選出來 ▸ 攻城戰 | `docs/mechanics/30-combat.md` |
| 1. 戰場怎麼被選出來 ▸ 野戰 | `docs/mechanics/30-combat.md` |
| 1.4 ⭐⭐ 兩個可直接寫成程式的定義 ▸ 國力 | `docs/mechanics/70-ai.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x00 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x01 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x02 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x06 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x08 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x14 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x18 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x1A | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x1C | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x1D | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x20 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x23 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x3E | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x3F | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x00 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x02 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x08 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0A | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0C | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0E | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x10 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x19 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1A | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1C–+0x1F | `docs/formats/08-sinario-save.md` |
| 1.7 軍團表：127 筆 × 64 B（區塊 +0x22C0） ▸ +0x09 | `docs/formats/08-sinario-save.md` |
| 2. 地形類型對映表（confirmed） ▸ 1 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 2 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 6 | `docs/mechanics/30-combat.md` |
| 3. 武將記錄（32 byte） ▸ +1 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +14（0x0E） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +15（0x0F） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +17 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +18 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +19 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +2 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +23（0x17） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +24（0x18） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +28（0x1C） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +29（0x1D） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +31（0x1F） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +8 | `docs/formats/08-sinario-save.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x1A | `docs/re/05-battle-selection.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x1C | `docs/re/05-battle-selection.md` |
| 4. 據點記錄（32 byte） ▸ +2 | `docs/formats/08-sinario-save.md` |
| 5.8l 兵士記錄（32 B）目前解出來的欄位 ▸ +0x1E | `docs/re/11-tactical-battle.md` |
| 已接入的原版資料流 ▸ sub_12C52／sub_12CDF | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_12DB8 | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_12DF3 | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_12EFB | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_13526／sub_13639 | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_145C1 | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_14698 | `docs/mechanics/70-ai.md` |
| 已接入的原版資料流 ▸ sub_16EC9 | `docs/mechanics/70-ai.md` |
| 據點記錄再解出四個欄位 ▸ +19h | `docs/re/07-monthly-settlement.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x00 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x01 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x02 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x06 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0A | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0B | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0E | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x10 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x12 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x14 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x16 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x18 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x1A | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x1C | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x1E | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x20 | `docs/re/08-hourly-update.md` |

### 強證據（26 條）

| 鍵 | 出處 |
|---|---|
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x03 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x11 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x12 | `docs/formats/08-sinario-save.md` |
| 2. 地形類型對映表（confirmed） ▸ 0 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 3 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 4 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 5 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 7 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 8 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 9 | `docs/mechanics/30-combat.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x01 | `docs/re/05-battle-selection.md` |
| 4. 據點記錄（32 byte） ▸ +1 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +10 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +26 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +8 | `docs/formats/08-sinario-save.md` |
| 據點記錄再解出四個欄位 ▸ +11h | `docs/re/07-monthly-settlement.md` |
| 據點記錄再解出四個欄位 ▸ +12h | `docs/re/07-monthly-settlement.md` |
| 據點記錄再解出四個欄位 ▸ +13h | `docs/re/07-monthly-settlement.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x04 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x00 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x02 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x06 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0A | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0C | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0E | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x0F | `docs/re/08-hourly-update.md` |

### 說明書（8 條）

| 鍵 | 出處 |
|---|---|
| 1. 三個能力值 ▸ +17 | `docs/mechanics/60-personnel.md` |
| 1. 三個能力值 ▸ +18 | `docs/mechanics/60-personnel.md` |
| 1. 三個能力值 ▸ +19 | `docs/mechanics/60-personnel.md` |
| 3.4 其他戰術判定（說明書） ▸ 中央突破戰法 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 擊破狀態 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 突擊時機 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 陣形有利不利 | `docs/mechanics/70-ai.md` |
| 5.8 ⭐ 十一個命令處理常式，隊長與隊員各一套 ▸ 命令 | `docs/re/11-tactical-battle.md` |

### 假說（1 條）

| 鍵 | 出處 |
|---|---|
| 2.1 執行結構（已驗證） ▸ YNSOUND.COM | `docs/reference/04-first-survey.md` |

### 未解（34 條）

| 鍵 | 出處 |
|---|---|
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x0008 | `docs/formats/08-sinario-save.md` |
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x1EC0 | `docs/formats/08-sinario-save.md` |
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x3AC0…+0x42C0 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x16 | `docs/formats/08-sinario-save.md` |
| 1.5 勢力表：22 筆 × 64 B（區塊 +0x80） ▸ +0x19 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x17 | `docs/formats/08-sinario-save.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D2F8 | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D2FC | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D306 | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D42 | `docs/re/11-tactical-battle.md` |
| 2.7 圖像與調色盤 ▸ KYOGRF.DAT | `docs/reference/04-first-survey.md` |
| 3. 曲塊內部 ▸ +0x00 | `docs/re/23-bgm-resource-format.md` |
| 3. 曲塊內部 ▸ +0x06–+0x0F | `docs/re/23-bgm-resource-format.md` |
| 3. 武將記錄（32 byte） ▸ +0 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +16（0x10） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +20（0x14） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +25,+27 | `docs/formats/08-sinario-save.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x08 | `docs/re/05-battle-selection.md` |
| 4. 只在單邊存在的檔（confirmed） ▸ PASS.MAP／PASS.SCH | `docs/re/01-first-recon.md` |
| 4. 據點記錄（32 byte） ▸ +0 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +12 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +14 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +16 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +17 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +18 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +22 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +27–+31 | `docs/formats/08-sinario-save.md` |
| 6. 月結與季節掛在時鐘的哪裡（confirmed） ▸ 世界更新 | `docs/mechanics/15-realtime.md` |
| 檔案 ▸ 40-economy.md | `docs/mechanics/00-index.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x08 | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0C | `docs/re/08-hourly-update.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x23 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x04 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x08 | `docs/re/08-hourly-update.md` |

