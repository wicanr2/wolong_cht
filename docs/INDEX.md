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
| [`docs/formats/03-grf-images.md`](formats/03-grf-images.md) | 03 — GRF.DAT 圖庫格式 | KAOGRF／KYOGRF／IVENTGRF READY，ICONGRF 部分解。 ⭐ 視窗底紋的點陣找到了——IC… | 2026-08-17 |
| [`docs/formats/04-map-sch-container.md`](formats/04-map-sch-container.md) | 04 — .MAP／.SCH／.MCH：兩種完全不同的東西 | 容器格式的索引層 READY，壓縮演算法未解。 | 2026-08-07 |
| [`docs/formats/05-mmap-worldmap.md`](formats/05-mmap-worldmap.md) | 05 — MMAP. 大地圖 | MMAP.MDL、地圖尺寸、自動連接與 MMAP.MCH 物件圖形入口 confirmed；MMAP.MAP 的 R… | 2026-08-07 |
| [`docs/formats/06-mmap-rle.md`](formats/06-mmap-rle.md) | 06 — MMAP.MAP 的 RLE 壓縮 | READY。 | 2026-08-07 |
| [`docs/formats/07-battle.md`](formats/07-battle.md) | 07 — BATTLE. 戰場資料 | 分段結構、圖塊定義、子圖塊與人物圖形的像素格式都 confirmed。 剩三項未解（§10）。 | 2026-08-07 |
| [`docs/formats/08-sinario-save.md`](formats/08-sinario-save.md) | 08 — SINARIO.DAT / SAVE.DAT：劇本與存檔 | 整體結構 confirmed，武將能力值 confirmed，其餘欄位進行中。 | 2026-08-07 |
| [`docs/formats/09-cutscene-images.md`](formats/09-cutscene-images.md) | 09 — OPEN_S.DAT／END_S.DAT：過場畫面 | READY。 | 2026-08-18 |
| [`docs/mechanics/00-index.md`](mechanics/00-index.md) | 00 — 遊戲機制索引 | 索引與推論等級定義，長期有效。 | 2026-08-08 |
| [`docs/mechanics/10-strategy.md`](mechanics/10-strategy.md) | 10 — 大地圖政略 | 指令清單完整；戰略數值與 AI 決策大多已由機器碼解出並實作。 剩六個位置的效果一項（§7）。 | 2026-08-13 |
| [`docs/mechanics/15-realtime.md`](mechanics/15-realtime.md) | 15 — 即時制的時間模型 | ✅ READY。整條時間鏈已在機器碼裡讀出來。 | 2026-08-08 |
| [`docs/mechanics/20-military.md`](mechanics/20-military.md) | 20 — 行軍與軍團 | 道路網與行軍已解並實作；AI 與玩家的編成流程都已解； ⭐ | 2026-08-13 |
| [`docs/mechanics/30-combat.md`](mechanics/30-combat.md) | 30 — 戰場（戰術） | 進場規則與戰略層的自動判定全解；戰術核心已接入並實作，完整結算與少數分支未完 | 2026-08-09 |
| [`docs/mechanics/40-economy.md`](mechanics/40-economy.md) | 40 — 經濟：資金與預備兵 | 機制與公式都已解 | 2026-08-08 |
| [`docs/mechanics/50-diplomacy.md`](mechanics/50-diplomacy.md) | 50 — 外交 | 成立條件與外交官的數值都已解 | 2026-08-13 |
| [`docs/mechanics/60-personnel.md`](mechanics/60-personnel.md) | 60 — 武將 | 三個能力值的作用與身分欄位已定案（說明書＋機器碼），數值公式部分已知。 ⭐ 政治的兩個實際用途也解了——內政效果 c… | 2026-08-13 |
| [`docs/mechanics/70-ai.md`](mechanics/70-ai.md) | 70 — 電腦 AI 的判斷邏輯 | 侵攻目標的決策鏈、友好度漂移、宣戰三閘與 AI 編成入口已由機器碼讀出；remake 已接上可重播的敵方出兵切片。 | 2026-08-09 |
| [`docs/mechanics/80-victory.md`](mechanics/80-victory.md) | 80 — 勝負判定 | 勝負條件全部定案。⭐ 結局的觸發已由機器碼定案 （[../re/59](../re/59-game-over-exi… | 2026-08-08 |
| [`docs/mobile/android-plan.md`](mobile/android-plan.md) | Android 版規劃 | 核心已接入，主畫面、四個入口、戰場與 SAF 匯入都可用。 | 2026-08-20 |
| [`docs/mobile/android-ux.md`](mobile/android-ux.md) | Android UX 規格 | 主畫面、四個入口、進言流程與戰場都已實作。 | 2026-08-20 |
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
| [`docs/playtest/13-dosv-natural-and-target-gui.md`](playtest/13-dosv-natural-and-target-gui.md) | DOS/V 自然畫面與目標平台 GUI 驗收 | 歷史驗收紀錄（2026-08-11，影片 oracle）。 | 2026-08-11 |
| [`docs/playtest/14-m7-review.md`](playtest/14-m7-review.md) | M7 校訂文字人工審查報告 | 60 筆定案校訂已完成逐筆語意、marker、硬換行、寬度與代表畫面抽樣。 | 2026-08-11 |
| [`docs/playtest/15-event2-5-talk-sampling.md`](playtest/15-event2-5-talk-sampling.md) | 事件 2–5 TALK 完整分支抽樣 | 36 個 raw TALK 頁面、18 組雙頁回應的分支、marker、硬換行、字寬與五列版面 抽樣通過；不宣稱完整… | 2026-08-11 |
| [`docs/playtest/16-event9-long-route.md`](playtest/16-event9-long-route.md) | 事件 9 長程通知流程 | 27 小時 bounded queue、玩家／非玩家／在野通知條件與 #409 no-op 已通過； 完整自然劇本依… | 2026-08-11 |
| [`docs/playtest/17-expert-dosbox-remake.md`](playtest/17-expert-dosbox-remake.md) | 17 — DOSBox 原版／remake 可玩性專家驗證 | remake 正常策略路徑與存檔／讀檔通過；DOS/V 原版密碼頁已可進入開場，尚未展開完整自然長程驗證；PC-98… | 2026-08-11 |
| [`docs/playtest/18-dosv-password-verification.md`](playtest/18-dosv-password-verification.md) | 18 — 松崗 DOS/V 密碼頁輸入驗證 | 已證實，在受控 DOSBox-X 重播中按「確定」即可越過密碼頁；密碼頁不再是 DOS/V 原版行為驗證的阻擋。 | 2026-08-12 |
| [`docs/playtest/19-tactical-minimap.md`](playtest/19-tactical-minimap.md) | 19 — DOS/V 戰術縮圖 raw producer 驗收 | PASS（已證實 producer 的 remake 實作）。底圖與部隊點都已接上； 陣形線、游標十字與城壁受損的局… | 2026-08-12 |
| [`docs/playtest/20-tactical-layout-parity.md`](playtest/20-tactical-layout-parity.md) | 20 — 松崗 DOS/V 戰術版面 parity 重開 | 歷史量測紀錄（2026-08-12，影片幀對幾何）。 | 2026-08-12 |
| [`docs/playtest/21-command-window-parity.md`](playtest/21-command-window-parity.md) | 21 — 松崗 DOS/V 指揮／事件／一覽畫面 parity 重開 | 歷史量測紀錄（2026-08-12，影片幀對幾何）。 | 2026-08-12 |
| [`docs/playtest/21-dosboxx-bridge-sampling.md`](playtest/21-dosboxx-bridge-sampling.md) | 21 — DOSBox-X AI Bridge：第一次動態取樣 | 三條斷言全部取到證據。+0x00 低 4 位 ＝ 敵方鄰居遮罩（192/192， 對照讀法只對 12/192）；+0… | 2026-08-14 |
| [`docs/playtest/22-field-siege-shared-layout.md`](playtest/22-field-siege-shared-layout.md) | 攻城／兩軍遭遇共用戰術骨架驗收 | PASS（共用幾何與原版指令面板已封口；不代表動畫逐像素 parity） | 2026-08-12 |
| [`docs/playtest/23-main-screen-geometry.md`](playtest/23-main-screen-geometry.md) | 23 — 主畫面幾何：從機器碼定死，第一次逐區對拍 | 版面常數全部換成機器碼算出來的值（外框四項 ＋ 右欄內部一整組）。 橫幅的位移掃描落在 (0,0)，幾何對齊。 | 2026-08-15 |
| [`docs/playtest/24-window-toggles.md`](playtest/24-window-toggles.md) | 24 — 四個常駐視窗的開關：實作驗收 | 開關可用，四窗全開與全關兩張截圖都拍到了。右欄的四條邊與原版 參考影片逐條對上（y = 168／184／192／19… | 2026-08-15 |
| [`docs/playtest/25-audio-capture-feasibility.md`](playtest/25-audio-capture-feasibility.md) | 25 — 音訊擷取的可行性：DOSBox 錄得到，鏈路已打通 | 可行性 confirmed。開場動畫的 15 秒錄音有內容（RMS 424）， 轉成 ogg 再解回來仍是 RMS … | 2026-08-15 |
| [`docs/playtest/26-bgm-render-vs-recording.md`](playtest/26-bgm-render-vs-recording.md) | 26 — 合成出來的音樂對得上原版錄音 | confirmed。internal/audio 渲染的 OPENBGM.DAT 與 DOSBox 錄的同一首， 包… | 2026-08-15 |
| [`docs/playtest/27-original-video-frame-parity.md`](playtest/27-original-video-frame-parity.md) | 27 — ⭐ 拿原版實錄影片對版面：主畫面與戰術畫面的幾何都落在 3 px 內 | 主畫面與戰術畫面的 | 2026-08-16 |
| [`docs/playtest/28-siege-breach-measurement.md`](playtest/28-siege-breach-measurement.md) | 28 — 量攻城：remake 的攻方打不進城，原因是城牆四格厚 | 歷史量測紀錄。 | 2026-08-16 |
| [`docs/playtest/29-strategy-minimap-markers.md`](playtest/29-strategy-minimap-markers.md) | 29 — 縮小地圖的據點標記接上原版的四種顏色 | ✅ 已實作並留下截圖。 | 2026-08-16 |
| [`docs/playtest/30-ground-planes-implemented.md`](playtest/30-ground-planes-implemented.md) | 30 — 兩個平面的地面圖接上規則層：攻方終於會打城牆了 | 通過。 | 2026-08-16 |
| [`docs/playtest/31-parity-inventory.md`](playtest/31-parity-inventory.md) | 31 — 原版 vs remake 逐畫面盤點（2026-08-16） | 盤點，不是量測。 | 2026-08-16 |
| [`docs/playtest/32-talk-layout-fit.md`](playtest/32-talk-layout-fit.md) | 32 — M7 排版 parity：1,022 則逐則量進訊息框 | 量完了。單行超寬 | 2026-08-16 |
| [`docs/playtest/33-ai-march-long-run.md`](playtest/33-ai-march-long-run.md) | 33 — AI 行軍鏈接上之後的長跑觀察 | 量完了。世界會動了——軍團有生有滅、不變量全程成立。 ⭐ 最重要的發現不是 AI 的行為，是 | 2026-08-17 |
| [`docs/playtest/34-advise-scene-screens.md`](playtest/34-advise-scene-screens.md) | 34 — 進言的畫面驗收：插圖 ＋ 兩個框 ＋ 五列選單 | 通過。 | 2026-08-17 |
| [`docs/playtest/35-advise-verdict-screens.md`](playtest/35-advise-verdict-screens.md) | 35 — 進言的五項與「請求君主出陣」的定案畫面 | 通過。 | 2026-08-17 |
| [`docs/playtest/36-window-texture.md`](playtest/36-window-texture.md) | 36 — 視窗底紋畫上去了 | 通過。 | 2026-08-17 |
| [`docs/playtest/37-main-screen-parity.md`](playtest/37-main-screen-parity.md) | 37 — 開局主畫面的逐區對拍：第一次真的對原版跑 | ⭐⭐ 五區全部 PASS。開局主畫面的 640×400 逐像素等於原版。 | 2026-08-17 |
| [`docs/playtest/38-window-parity.md`](playtest/38-window-parity.md) | 38 — 三個視窗開著時的對拍：三個視窗區逐像素相同 | 通過。 | 2026-08-17 |
| [`docs/playtest/39-system-window-parity.md`](playtest/39-system-window-parity.md) | 39 — 系統選單開著時的對拍：五區裡四區 PASS，選單本身也 PASS | 通過（2026-08-17 的量測），但 2026-08-23 之後 | 2026-08-17 |
| [`docs/playtest/40-tactical-parity.md`](playtest/40-tactical-parity.md) | 40 — 戰場的逐區對拍：六區逐像素相同，戰場區 0.17% | 部分通過。 | 2026-08-18 |
| [`docs/playtest/41-m7-corrected-text-on-screen.md`](playtest/41-m7-corrected-text-on-screen.md) | 41 — M7 校訂後的畫面抽樣：18 則實跑截圖，沒有一行超寬 | 通過。 | 2026-08-22 |
| [`docs/playtest/42-window-parity.md`](playtest/42-window-parity.md) | 42 — 四類視窗的實機對拍：財政收斂到 NEAR，一覽表與編成抓出五項實質差異 | （2026-08-25 第三輪）武將一覽（編成清單）、財政視窗、 數值輸入器、交戰目標勢力清單 | 2026-08-24 |
| [`docs/playtest/43-field-battle-parity.md`](playtest/43-field-battle-parity.md) | 43 — 野戰的同狀態對拍：九區裡七區逐像素相同，地形 0 差 | （2026-08-25 對白收斂後）九區裡七區 0 差 PASS、 field | 2026-08-24 |
| [`docs/promo/README.md`](promo/README.md) | 推廣片產出紀錄 | 五支影片都已產出並驗過媒體規格。主預告的遊戲段落是逐幀錄下來的實跑畫面； 三平台候選包與 Linux AppImag… | 2026-08-21 |
| [`docs/promo/android.md`](promo/android.md) | Android 版推廣片 | 已產出 48 秒片；素材只有 remake 自己的畫面與本專案原創的合成配樂。 | 2026-08-20 |
| [`docs/promo/classic-revival.md`](promo/classic-revival.md) | 「經典再現」推廣片 | 已產出研究／推廣用 60 秒比較片；不把代表幀比較宣稱為同狀態逐像素 parity。 | 2026-08-11 |
| [`docs/promo/dosv-adlib-and-tactical-review.md`](promo/dosv-adlib-and-tactical-review.md) | 推廣片原版 AdLib 與戰術骨架審查 | 歷史審查紀錄（2026-08-12）。 | 2026-08-12 |
| [`docs/promo/dosv-live-comparison.md`](promo/dosv-live-comparison.md) | DOS/V 原版／remake 實機動態比較推廣片 | 已產出並驗過，已被 dosv-realmachine 取代。 | 2026-08-12 |
| [`docs/promo/dosv-realmachine.md`](promo/dosv-realmachine.md) | 原版實機遊玩 × remake 實機的對照推廣片 | 已產出並驗過。原版側是自己跑的受控 DOSBox-X 實機遊玩， 只有戰術戰場那一格仍取自使用者提供的錄影。 | 2026-08-23 |
| [`docs/promo/yt-remake-pixel-review.md`](promo/yt-remake-pixel-review.md) | YouTube／remake 推廣片像素差異審查 | 已完成影片對照；確認像素差異，但不把不同遊戲狀態誤稱為同狀態 逐像素 parity。 | 2026-08-11 |
| [`docs/re/00-index.md`](re/00-index.md) | 00 — RE 知識庫入口 | 索引。44 份反組譯筆記按子系統分類，每份一句話說它回答什麼問題。 | 2026-08-14 |
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
| [`docs/re/17-dosv-audio-tsr.md`](re/17-dosv-audio-tsr.md) | 17 — 松崗 DOS/V 音源 TSR 與戰術效果碼 | INT 61h 介面、遊戲端效果碼與硬體 register 寫入已證實。 ⭐ 晶片是 | 2026-08-12 |
| [`docs/re/18-tactical-button-glyphs.md`](re/18-tactical-button-glyphs.md) | 18 — DOS/V 戰術底列按鈕 glyph 候選資產研究 | 六個命令 glyph、底板、右欄複合面板與選取矩形已解出並接入 remake。 | 2026-08-12 |
| [`docs/re/19-outcome.md`](re/19-outcome.md) | DOS/V 已證實敗北 outcome 接線 | READY（只涵蓋兩種敗北；不涵蓋勝利、君主死亡或原版返回標題）。 | 2026-08-12 |
| [`docs/re/20-ida-re-coverage-audit.md`](re/20-ida-re-coverage-audit.md) | 20 — DOS/V IDA 逆向覆蓋與 remake 差距審計 | REVIEWED。足以支撐可玩重製與多數規則，但不足以支撐高忠實度戰術呈現；主要缺口是原版顯示串列、相機狀態機、逐幀… | 2026-08-12 |
| [`docs/re/21-function-census.md`](re/21-function-census.md) | 21 — DOS/V KI.EXE 全函式覆蓋普查 | 量測完成， | 2026-08-14 |
| [`docs/re/22-strategy-command-tree.md`](re/22-strategy-command-tree.md) | 22 — 戰略指令列與滑鼠熱區分派 | 指令列的八個槽位、兩層子選單、財政四子項、編成入口、熱區查表定址式與 狀態列提示常式 confirmed。指令列的圖… | 2026-08-13 |
| [`docs/re/23-bgm-resource-format.md`](re/23-bgm-resource-format.md) | 23 — BGM.DAT 音樂資源格式 | 容器索引與曲塊的聲軌指標表 confirmed，兩版都逐 byte 對齊到檔尾且餘 0。 ⭐ | 2026-08-13 |
| [`docs/re/24-unread-function-catalogue.md`](re/24-unread-function-catalogue.md) | 24 — 未讀函式目錄：252 支的證據與下手順序 | 252 支未讀函式全部登記，附上由共用常式與 TALK 訊息取得的角色證據。 計數要排除本檔自己（§5），否則登記等… | 2026-08-13 |
| [`docs/re/25-message-variants-and-personnel.md`](re/25-message-variants-and-personnel.md) | 25 — 訊息變體展開與人事指令 | 訊息索引 ≥ 0x196 的 ×8 變體展開 confirmed，13 個呼叫值共 104 則變體已對出； 人事四支… | 2026-08-13 |
| [`docs/re/26-list-window-engine.md`](re/26-list-window-engine.md) | 26 — 一覽表視窗引擎 | 視窗幾何、五個一覽表家族的描述子、選取迴圈、持久化排序狀態、四個描述子欄位 （兩個 callback ＋ 標題字串 … | 2026-08-13 |
| [`docs/re/27-list-row-fields.md`](re/27-list-row-fields.md) | 27 — 一覽表的逐列繪製與外交關係等級 | 四個家族的逐列常式全部讀完，欄位對照 confirmed；兵力的 ×10 顯示、 三條換色規則與外交關係六級換算 c… | 2026-08-13 |
| [`docs/re/28-text-number-rendering.md`](re/28-text-number-rendering.md) | 28 — 文字與數字的繪製層 | 數字繪製、兩支字串繪製與 EGA 平面寫入方式 confirmed。 單字元 blitter loc_1F75E 與… | 2026-08-13 |
| [`docs/re/29-font-service-int15.md`](re/29-font-service-int15.md) | 29 — 原版怎麼顯示中文：INT 15h 字型服務與 END_S13/S14.DAT | 整條鏈 confirmed（靜態）。KI.EXE 側走 DOS/V 的 INT 15h AH=50h 向常駐服務要字… | 2026-08-13 |
| [`docs/re/30-corps-formation-ui.md`](re/30-corps-formation-ui.md) | 30 — 軍團編成畫面：兵員池、兵種切換與派生數值 | 編成畫面的主迴圈與被它呼叫的常式 confirmed。 兵員池的搬運方向、兵種切換的循環、兵力的自動分配、確定的前提… | 2026-08-13 |
| [`docs/re/31-faction-picker-screen.md`](re/31-faction-picker-screen.md) | 31 — 勢力一覽：22 格的兩欄版面與領地重繪 | 版面、命中判定、顏色規則、兩個守門條件與介面字串 confirmed。 分派表已印出，但 sub_15AD1 → s… | 2026-08-14 |
| [`docs/re/32-strategy-detail-panels.md`](re/32-strategy-detail-panels.md) | 32 — 戰略側的兩個詳情面板：據點與軍團 | 兩支面板函式與軍團面板的外層 confirmed，欄位與字串表全部對出。 畫肖像的 sub_107D2 見 [33]… | 2026-08-14 |
| [`docs/re/33-shared-draw-helpers.md`](re/33-shared-draw-helpers.md) | 33 — 共用繪圖層：字串包裝、肖像快取、小地圖上色 | 六支共用常式 confirmed。肖像快取的替換策略、小地圖的座標換算與 據點標記的畫法都定案。實際載入 bytes… | 2026-08-14 |
| [`docs/re/34-corps-status-bits.md`](re/34-corps-status-bits.md) | 34 — 軍團記錄 +0x00 的位元圖，與改用 IDAPython 之後的掃法 | 位元 1／2 的設定端、清除端與語意 confirmed （位元 1 ＝ 下一步要重算、位元 2 ＝ 委任）。 位元… | 2026-08-14 |
| [`docs/re/35-strategy-ui-module-map.md`](re/35-strategy-ui-module-map.md) | 35 — 戰略 UI 模組全圖：108 支函式的叢集歸屬 | 叢集歸屬 confirmed（呼叫圖是精確的，不是啟發式）。 各叢集的角色標籤是強證據——來自「它呼叫哪些已定案語意… | 2026-08-14 |
| [`docs/re/36-tactical-module-map.md`](re/36-tactical-module-map.md) | 36 — 戰術戰鬥模組全圖：主迴圈與它的十一個子系統 | 叢集歸屬 confirmed（呼叫圖精確）。角色標籤是強證據， 來自「呼叫哪些已定案語意的共用常式」與 I/O 埠使… | 2026-08-14 |
| [`docs/re/37-graphics-and-runtime-module-map.md`](re/37-graphics-and-runtime-module-map.md) | 37 — 圖庫、繪圖底層與 C runtime 兩個模組的全圖 | 叢集歸屬 confirmed。硬體層的角色由 I/O 埠直接判定（精確）， 其餘角色標籤是強證據。 | 2026-08-14 |
| [`docs/re/38-strategy-core-module-map.md`](re/38-strategy-core-module-map.md) | 38 — 戰略核心三個模組的全圖：指令樹、到站處理與月結 | 八個指令的 handler 全部定位，狀態列與選單索引一併對出（confirmed）。 到站處理與月結的叢集歸屬 c… | 2026-08-14 |
| [`docs/re/39-remaining-unread.md`](re/39-remaining-unread.md) | 39 — 剩餘未讀函式的逐支歸屬 | 清單。 | 2026-08-14 |
| [`docs/re/40-garrison-relief-request.md`](re/40-garrison-relief-request.md) | 40 — 據點求援與援軍派遣 | 整條鏈 confirmed（每一支都逐行讀過）。 sub_140C9 的距離算式裡有一處 | 2026-08-14 |
| [`docs/re/42-leaf-functions.md`](re/42-leaf-functions.md) | 42 — 戰術以外的 47 支葉節點 | 47 支全部逐行讀過。四件事因此定案：INT 61h 是音源 TSR 的介面、 byte_198A6 的位元圖完整、… | 2026-08-14 |
| [`docs/re/43-open-questions.md`](re/43-open-questions.md) | 43 — 未解缺口總表（生成的） | 生成的清單，跑 tools/py.sh tools/re_open_questions.py 重出。 這一份不下結論… | 2026-08-25 |
| [`docs/re/44-threat-and-reinforcement-ai.md`](re/44-threat-and-reinforcement-ai.md) | 44 — 威脅偵測與 AI 出兵：據點每 tick 掃一次 | 整條鏈逐行讀完。三件事定案：據點 +0x18 是佔用圖讀回來的軍團數、 +0x00 低 4 位是「哪幾個鄰居是敵方」… | 2026-08-14 |
| [`docs/re/45-corps-command-mode.md`](re/45-corps-command-mode.md) | 45 — 軍團的三種指令模式：戰鬥指揮／委任／解體 | 軍團 +0x00 位元 2 定案 ＝ | 2026-08-14 |
| [`docs/re/46-strategy-chrome-cell-layer.md`](re/46-strategy-chrome-cell-layer.md) | 46 — 主畫面的指令列沒有按鈕圖，外框取自 ICONGRF 段 3 | 指令列的繪製路徑逐支讀完。指令列 | 2026-08-15 |
| [`docs/re/47-main-screen-window-registry.md`](re/47-main-screen-window-registry.md) | 47 — 主畫面的四個常駐視窗：開關、位元集與各自的矩形 | 四個常駐視窗的開關熱區、位元集、分派表與四個視窗各自的像素矩形 全部 confirmed。繪圖／熱區常式的參數對應（… | 2026-08-15 |
| [`docs/re/48-window-display-list.md`](re/48-window-display-list.md) | 48 — 視窗內容是一份顯示清單，不是一張圖 | 清單的位置、記錄格式、場景切分、十個場景的歸屬與九個 opcode 的語意全部解出來了。⭐ 記錄的第六個 word … | 2026-08-15 |
| [`docs/re/49-corps-formation-window.md`](re/49-corps-formation-window.md) | 49 — 軍團編成視窗的版面與動態層 | 視窗矩形、靜態層、六個槽的圖示與數字、四個顯示值的來源與座標 全部解出來了。 | 2026-08-15 |
| [`docs/re/50-city-info-window.md`](re/50-city-info-window.md) | 50 — 據點情報視窗，以及據點 +0x16 高 4 位的用途 | 視窗矩形、靜態層、五個顯示值的來源與座標、左半那張 96×96 圖的 出處全部解出來了。⭐ 據點 +0x16 的高 … | 2026-08-15 |
| [`docs/re/51-corps-info-window.md`](re/51-corps-info-window.md) | 51 — 軍團情報視窗（顯示清單場景 4） | 視窗矩形、靜態層、九個顯示值的座標與來源全部解出來了。 ⭐ 這個視窗畫空槽時會取到 | 2026-08-15 |
| [`docs/re/52-slot-select-window.md`](re/52-slot-select-window.md) | 52 — 四槽選擇視窗：新遊戲、讀取、儲存共用同一個 | 視窗矩形、靜態層、四個槽的內容與座標、三個標題、 「哪些槽不能選」的判定全部解出來了。⭐ 原版的「新遊戲」不是另一個… | 2026-08-15 |
| [`docs/re/53-lord-select-window.md`](re/53-lord-select-window.md) | 53 — 君主選擇視窗（顯示清單場景 8） | 視窗矩形、靜態層、七個顯示值的座標與來源、兩個熱區的語意 全部解出來了。⭐ 「確定」在軍師還沒命名時會被導去「自定」… | 2026-08-15 |
| [`docs/re/54-advisor-naming-window.md`](re/54-advisor-naming-window.md) | 54 — 軍師命名視窗（顯示清單場景 9，松崗版特有） | 視窗矩形、靜態層、九個熱區的位置與語意解出來了。 選字表的資料來源與翻頁邏輯未讀（§4）。 | 2026-08-15 |
| [`docs/re/55-system-menu-window.md`](re/55-system-menu-window.md) | 55 — 系統選單視窗，以及 op 08 屬性的真正編碼 | 視窗矩形、六列版面、六個熱區、⭐ 中間四列的設定值與選項字串 全部解出來了（§4）。六個 handler 讀到四支（… | 2026-08-15 |
| [`docs/re/56-bgm-track-events.md`](re/56-bgm-track-events.md) | 56 — ⭐ BGM.DAT 的聲軌事件編碼解出來了 | 事件編碼、三張查表、播放引擎的迴圈、控制事件的分派全部 confirmed。 音符與控制事件都可以完整還原成語意（§… | 2026-08-15 |
| [`docs/re/57-opl3-register-map.md`](re/57-opl3-register-map.md) | 57 — ⭐ DOS/V 的音源是 OPL3，六個聲軌各佔一組 4-operator 通道 | 晶片型號、通道配置、音色記錄版面、音量與速度換算、SOUND.DAT 的記錄結構全部 confirmed。剩兩張表的… | 2026-08-15 |
| [`docs/re/58-bgm-scene-mapping.md`](re/58-bgm-scene-mapping.md) | 58 — ⭐ 哪一首配哪個場景：BGM.DAT 的 11 首全部對出來了 | 曲 0／2–10 全部 confirmed（呼叫端的立即值、查表或計算式）。 只剩曲 1 ——它在 DOS/V 的 … | 2026-08-15 |
| [`docs/re/59-game-over-exit-codes.md`](re/59-game-over-exit-codes.md) | 59 — ⭐ 結局與敗北是靠離開碼交出去的 | 三個離開碼與各自的觸發點 confirmed。 ⭐ 結局的閘門是 | 2026-08-16 |
| [`docs/re/60-tactical-sidebar.md`](re/60-tactical-sidebar.md) | 60 — 戰術側欄：那一欄畫了什麼，每一格由誰畫 | 側欄七格的內容全部解出（confirmed）。 | 2026-08-16 |
| [`docs/re/61-timer-tick-source.md`](re/61-timer-tick-source.md) | 61 — 計時中斷是誰發的：節流的頻率終於有數字了 | ✅ 解出來了。 | 2026-08-16 |
| [`docs/re/62-strategy-minimap.md`](re/62-strategy-minimap.md) | 62 — 主畫面縮小地圖：據點標記的四種顏色、視野框、勢力篩選 | ✅ 內容組成全解。 | 2026-08-16 |
| [`docs/re/63-ground-plane-map.md`](re/63-ground-plane-map.md) | 63 — ⭐ 登城機制解完了：地面層表在另一個段，城門那一格就是樓梯 | ✅ 解出來了。 | 2026-08-16 |
| [`docs/re/64-corps-arrival-state-machine.md`](re/64-corps-arrival-state-machine.md) | 64 — 軍團抵達時的狀態機：+0x23 的分派表與解體 | 分派表的結構、索引算式與 Stage 8–11 的語意 confirmed。 ⭐ 解體的消費端找到了——Stage … | 2026-08-16 |
| [`docs/re/65-ai-march-decision-chain.md`](re/65-ai-march-decision-chain.md) | 65 — 電腦勢力的行軍決策鏈：Stage 0–3 與整條軍團生命週期 | 四支 AI handler 逐行讀完，Stage 0–3／8–11 的轉移條件全部 confirmed。 ⭐ 目標選… | 2026-08-17 |
| [`docs/re/66-message-box-geometry.md`](re/66-message-box-geometry.md) | 66 — 訊息框的版面：一個框、一張肖像、四列字 | 全部 confirmed，而且機器碼與原版實錄影格兩條獨立證據對上。 ⭐ 框是固定的 (160, 160, 256,… | 2026-08-17 |
| [`docs/re/67-city-emblem-on-strategy-map.md`](re/67-city-emblem-on-strategy-map.md) | 67 — 大地圖上的據點徽記：位置就在記錄座標，顏色照勢力分三類 | ✅ 全解並實作。 | 2026-08-17 |
| [`docs/re/68-t3-frontier-functions.md`](re/68-t3-frontier-functions.md) | 68 — T3 那九支：只在狀態檔與程式碼裡出現過的函式 | 九支全部讀完。 | 2026-08-18 |
| [`docs/re/69-t2-cross-reference.md`](re/69-t2-cross-reference.md) | 69 — T2 那 18 支：逐支讀過，各自歸位 | 完成。 | 2026-08-18 |
| [`docs/re/70-d7end-ending-player.md`](re/70-d7end-ending-player.md) | 70 — D7END.EXE：結局播放器與結局全文 | 播放順序、版面、結尾文字與過場圖的格式都解出來了。 | 2026-08-18 |
| [`docs/re/71-strategy-hotspot-dispatch.md`](re/71-strategy-hotspot-dispatch.md) | 71 — 戰略層的兩張熱區分派表，以及點縮小地圖會發生什麼 | 左鍵表 off_159D2 的 32 筆全部攤開，索引就是熱區碼。⭐ 熱區 0x16 （點縮小地圖）＝ 把大地圖鏡頭… | 2026-08-22 |
| [`docs/re/72-world-map-display-list.md`](re/72-world-map-display-list.md) | 72 — 大地圖的顯示表：地形一層 ＋ 最多四層疊圖 | 已解。 | 2026-08-23 |
| [`docs/re/73-new-game-faction-list.md`](re/73-new-game-faction-list.md) | 73 — 新遊戲怎麼選君主：先一張清單，再一張卡 | 整條流程解出來了。⭐ 君主卡上沒有「換勢力」的熱區—— 換勢力是退回上一層的 | 2026-08-24 |
| [`docs/re/74-battle-opening-duel.md`](re/74-battle-opening-duel.md) | 74 — 開戰喊話是單挑狀態機的開頭：挑戰、拒戰、對打互嗆、決著 | 全段解出並實作（2026-08-25）——挑戰／拒戰／應戰、回合互嗆、 對打段與決著都在 internal/rule… | 2026-08-25 |
| [`docs/reference/01-jp-manual.md`](reference/01-jp-manual.md) | 01 — 日文原版說明書判讀紀錄 | 有實質機制的頁都讀完了，剩 p.6 啟動操作與 p.36–38 附錄。 | 2026-08-08 |
| [`docs/reference/02-jp-cht-diff.md`](reference/02-jp-cht-diff.md) | 02 — 日中對照：TALK.DAT 第一批發現 | 全量 1,022 則的 | 2026-08-16 |
| [`docs/reference/03-baked-japanese.md`](reference/03-baked-japanese.md) | 03 — 燒進美術裡的日文：松崗版沒重繪的部分 | 已確認的缺口：標題橫幅「臥竜伝」兩版相同（松崗沒重繪）。 | 2026-08-07 |
| [`docs/reference/04-first-survey.md`](reference/04-first-survey.md) | 04 — 首輪偵查紀錄（2026-08-07） | 歷史快照。已解的項目由 docs/formats/、docs/re/ 與 docs/INDEX.md 接手；本檔只保… | 2026-08-07 |
| [`docs/reference/05-eten-font-provenance.md`](reference/05-eten-font-provenance.md) | 05 — 松崗版的中文字型是倚天字型 | confirmed。END_S14.DAT 與倚天 ascfont.15 byte-for-byte 相同； END… | 2026-08-13 |
| [`docs/release/01-cross-build-gate.md`](release/01-cross-build-gate.md) | 01 — 發行閘重跑：五平台交叉建置 ＋ deny-list（2026-08-17） | 發行閘通過。五個平台建得出來、deny-list 掃過沒有原版資產。 ⭐ macOS 的 Ebiten 本體可以交叉… | 2026-08-17 |
| [`docs/release/02-three-platform-20260820.md`](release/02-three-platform-20260820.md) | 02 — 2026-08-20 三平台重新交付 | 已交付並驗過。 | 2026-08-20 |
| [`docs/release/03-three-platform-20260821.md`](release/03-three-platform-20260821.md) | 03 — 2026-08-21 三平台重新交付（含 Android 版） | 已交付並驗過，已被 wolong-remake-20260822 那一批取代。 | 2026-08-21 |
| [`docs/release/04-three-platform-20260822.md`](release/04-three-platform-20260822.md) | 04 — 2026-08-22 四平台完整版（內含遊戲檔案） | 已交付並驗過，已被 wolong-remake-20260823 那一批取代。 | 2026-08-22 |
| [`docs/release/05-full-20260823.md`](release/05-full-20260823.md) | 05 — 2026-08-23 四平台完整版（修掉玩家回報的三個問題） | 已交付並驗過。⛔ 內含原版資產，不可外流。 | 2026-08-23 |
| [`docs/release/06-appimage-20260824.md`](release/06-appimage-20260824.md) | 06 — 2026-08-24 只重打 Linux AppImage（同一天兩次） | 已被 07-full-20260824.md 取代 | 2026-08-24 |
| [`docs/release/07-full-20260824.md`](release/07-full-20260824.md) | 07 — 2026-08-24 四平台完整版（一致批次） | 已交付並驗過。⛔ 內含原版資產，不可外流。 | 2026-08-24 |
| [`docs/release/README-RELEASE.md`](release/README-RELEASE.md) | 臥龍傳 remake 可執行封裝 | 四平台完整包、Linux AppImage、推廣片與驗收紀錄已集中於 [dist-all](../../dist-a… | 2026-08-24 |
| [`docs/spec/00-index.md`](spec/00-index.md) | 00 — 規格索引：已解的規則有沒有被實作、有沒有被驗過 | 索引。規格是 docs/re/（程式碼在哪）與 internal/（我們寫了什麼） 之間的那一層——它回答「這條規則… | 2026-08-14 |
| [`docs/spec/10-city-tick.md`](spec/10-city-tick.md) | 10 — 據點整備、威脅偵測與求援 | CONFORMED。整條鏈已實作，並在 PC-98 原版的執行期記憶體上取樣驗過 （+0x18／+0x14 各 0/… | 2026-08-14 |
| [`docs/spec/11-ai-sortie.md`](spec/11-ai-sortie.md) | 11 — 進言「請求君主出陣」 | CONFORMED。 | 2026-08-14 |
| [`docs/spec/12-strategy-chrome.md`](spec/12-strategy-chrome.md) | 12 — 主畫面的視窗外框、指令列與右欄 | CONFORMED。主畫面的四個常駐視窗矩形、指令列版面與縮小地圖／勢力篩選鈕的 位置全部由機器碼定死（[docs/… | 2026-08-15 |
| [`docs/spec/13-main-window-toggles.md`](spec/13-main-window-toggles.md) | 13 — 主畫面四個視窗的開關 | CONFORMED。已實作並留下四窗全開／全關的截圖； 舊的 g.open[] 那一套已整個拿掉，主畫面視窗只剩一份… | 2026-08-15 |
| [`docs/spec/14-finance-window.md`](spec/14-finance-window.md) | 14 — 財政視窗 | CONFORMED。版面已照原版重寫並有契約測試； 數值輸入器已接上（[78](78-amount-input-ed… | 2026-08-15 |
| [`docs/spec/20-save-format.md`](spec/20-save-format.md) | 20 — remake 原生存檔格式 | CONFORMED。編解碼、路徑與遊戲接線都實作並驗過。 存檔一次寫兩份（原版格式 ＋ 原生檔），讀檔優先原生檔。 | 2026-08-14 |
| [`docs/spec/21-corps-formation-reserves.md`](spec/21-corps-formation-reserves.md) | 21 — 編成時預備兵怎麼分配 | CONFORMED。已實作並有逐項單測。 | 2026-08-15 |
| [`docs/spec/22-corps-formation-window.md`](spec/22-corps-formation-window.md) | 22 — 軍團編成視窗 | CONFORMED。版面、武將頭像與六個槽的滑鼠熱區都照原版實作， 並有契約測試。 | 2026-08-16 |
| [`docs/spec/23-city-info-window.md`](spec/23-city-info-window.md) | 23 — 據點情報視窗 | CONFORMED。版面已照原版實作並有契約測試。 | 2026-08-15 |
| [`docs/spec/24-corps-info-window.md`](spec/24-corps-info-window.md) | 24 — 軍團情報視窗 | CONFORMED。版面已照原版實作並有契約測試。 | 2026-08-15 |
| [`docs/spec/25-slot-select-window.md`](spec/25-slot-select-window.md) | 25 — 四槽選擇視窗（新遊戲／讀取／儲存） | CONFORMED。讀取／儲存已照原版版面實作； 新遊戲仍走 remake 自己的啟動殼層（§5）。 | 2026-08-15 |
| [`docs/spec/26-yes-no-dialog.md`](spec/26-yes-no-dialog.md) | 26 — ＹＥＳ／ＮＯ 對話框 | CONFORMED。版面與命中算式已照原版實作並有契約測試。 | 2026-08-15 |
| [`docs/spec/27-lord-select-window.md`](spec/27-lord-select-window.md) | 27 — 君主選擇視窗 | CONFORMED。版面已照原版實作並有契約測試； 輸入照原版收斂成兩個熱區（§2.1）。「自定」（軍師命名）還沒接… | 2026-08-15 |
| [`docs/spec/28-scenario-json.md`](spec/28-scenario-json.md) | 28 — 劇本的 JSON 匯出與匯入 | CONFORMED。四個區塊 round-trip 全過。 | 2026-08-15 |
| [`docs/spec/29-audio.md`](spec/29-audio.md) | 29 — 音樂與音效 | CONFORMED。音樂與音效都會出聲、與原版錄音比對過， ⭐ | 2026-08-15 |
| [`docs/spec/30-victory.md`](spec/30-victory.md) | 30 — 結局：存活勢力數歸一 | CONFORMED。⭐ 觸發條件在機器碼裡是 | 2026-08-16 |
| [`docs/spec/31-tactical-sidebar.md`](spec/31-tactical-sidebar.md) | 31 — 戰術側欄的內容組成 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/32-gate-strength-bar.md`](spec/32-gate-strength-bar.md) | 32 — 攻城的「門強度」條 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/33-squad-selection.md`](spec/33-squad-selection.md) | 33 — 底列六格是選部隊，不是第二套命令列 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/34-speed-steps.md`](spec/34-speed-steps.md) | 34 — 兩個速度設定：五檔、各檔的實際節奏 | READY。 | 2026-08-16 |
| [`docs/spec/35-strategy-minimap.md`](spec/35-strategy-minimap.md) | 35 — 縮小地圖的據點標記與視野框 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/36-ground-planes-and-climbing.md`](spec/36-ground-planes-and-climbing.md) | 36 — 兩個平面的地面圖、導航位元與登城 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/37-tactical-player-controls.md`](spec/37-tactical-player-controls.md) | 37 — 戰術畫面的玩家操作：陣形選單與陣形線 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/38-list-windows.md`](spec/38-list-windows.md) | 38 — 一覽表：視窗幾何、欄位與逐列格式 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/39-march-order-menu.md`](spec/39-march-order-menu.md) | 39 — 行軍指示的三選一：戰鬥指揮／委任／解體 | CONFORMED。 | 2026-08-16 |
| [`docs/spec/40-ai-march-decision.md`](spec/40-ai-march-decision.md) | 40 — 電腦勢力的行軍決策鏈（Stage 0–3／8／10） | CONFORMED。 | 2026-08-17 |
| [`docs/spec/41-message-box-geometry.md`](spec/41-message-box-geometry.md) | 41 — 訊息框的版面常數 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/42-event-scene-speakers.md`](spec/42-event-scene-speakers.md) | 42 — 事件場景上誰在說話 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/43-rout-on-blocked-return.md`](spec/43-rout-on-blocked-return.md) | 43 — 回不了家的軍團會敗走 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/44-advise-original-text.md`](spec/44-advise-original-text.md) | 44 — 進言用原版的原文，不用改寫的句子 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/45-advise-scene-layout.md`](spec/45-advise-scene-layout.md) | 45 — 進言的畫面：插圖 ＋ 兩個框輪流講話 ＋ 五列選單 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/46-post-battle-retreat.md`](spec/46-post-battle-retreat.md) | 46 — 戰後敗方退一站回家，退不了就壞滅 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/47-city-fall-corps-redirect.md`](spec/47-city-fall-corps-redirect.md) | 47 — 據點易主之後，舊主留在那一格的軍團調頭回家 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/48-governor-returns-on-city-fall.md`](spec/48-governor-returns-on-city-fall.md) | 48 — 據點被攻陷，派駐的內政官被遣回 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/49-advise-relocate-and-sortie.md`](spec/49-advise-relocate-and-sortie.md) | 49 — 進言的第四、五項：遷都與請求君主出陣 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/50-corps-upkeep-charges-funds.md`](spec/50-corps-upkeep-charges-funds.md) | 50 — 軍費直接扣資金，不進「本月支出」 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/51-vga-dac-palette-scale.md`](spec/51-vga-dac-palette-scale.md) | 51 — DOS/V 的顏色到不了滿刻度：4 bit → VGA 6 bit DAC → 8 bit | CONFORMED。 | 2026-08-17 |
| [`docs/spec/52-main-screen-camera-and-banner-date.md`](spec/52-main-screen-camera-and-banner-date.md) | 52 — 開局的鏡頭位置與橫幅日期 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/53-city-tile-by-ownership.md`](spec/53-city-tile-by-ownership.md) | 53 — 據點中心的圖塊跟著歸屬換，首都再疊一張 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/54-ui-colours-from-palette.md`](spec/54-ui-colours-from-palette.md) | 54 — 介面顏色一律查調色盤，命令列的底是黑的 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/55-minimap-view-box.md`](spec/55-minimap-view-box.md) | 55 — 縮小地圖的視野框是點陣，而且解開了「差四格」 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/56-battlefield-rotation.md`](spec/56-battlefield-rotation.md) | 56 — 戰場轉 180 度：什麼時候轉、轉的時候圖塊值要換 | CONFORMED。三段算式都接上了，並用原版的許昌攻防戰驗過： field 區 87.8% → 46.1%、小地圖… | 2026-08-17 |
| [`docs/spec/57-tactical-projection.md`](spec/57-tactical-projection.md) | 57 — 戰術畫面的等角投影：地形走出來、物件算出來，兩者不是同一條式子 | CONFORMED。 | 2026-08-17 |
| [`docs/spec/58-display-slot-depth-range.md`](spec/58-display-slot-depth-range.md) | 58 — 顯示格的深度範圍與 8 列的帶 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/59-battle-opening-orders.md`](spec/59-battle-opening-orders.md) | 59 — 開場的常令：腳本那一側不要先下「攻擊」 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/60-battle-talk-duration.md`](spec/60-battle-talk-duration.md) | 60 — 戰場對白顯示多久：60 個 tick，每側各一個 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/61-soldier-initial-hp-from-morale.md`](spec/61-soldier-initial-hp-from-morale.md) | 61 — 兵的開場體力 ＝ 軍團士氣 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/62-swapped-unit-skips-its-turn.md`](spec/62-swapped-unit-skips-its-turn.md) | 62 — 被換位的兵，這一幀不動 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/63-hit-stun.md`](spec/63-hit-stun.md) | 63 — 被打中的兵有硬直：當幀 ＋ 之後兩幀都不動 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/64-capital-relocation-report.md`](spec/64-capital-relocation-report.md) | 64 — 遷都之後說什麼：自國君主下令，他國要有外交官才報得回來 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/65-retreated-soldiers-survive.md`](spec/65-retreated-soldiers-survive.md) | 65 — 退到畫面外的兵算生還，不算戰死 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/66-broken-walls-repaint.md`](spec/66-broken-walls-repaint.md) | 66 — 城壁與門打壞之後，畫面上的地形要跟著換 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/67-ending-playback.md`](spec/67-ending-playback.md) | 67 — 結局的播放 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/68-death-animation.md`](spec/68-death-animation.md) | 68 — 倒地動畫：四幀，換一組圖 | CONFORMED。 | 2026-08-18 |
| [`docs/spec/69-world-fingerprint.md`](spec/69-world-fingerprint.md) | 69 — 世界指紋：同一個 seed 兩次跑出同一個值 | CONFORMED。⚠ 這一份沒有原版出處——它是 remake 的驗證設施，不是原版機制。 | 2026-08-20 |
| [`docs/spec/70-phone-chrome.md`](spec/70-phone-chrome.md) | 70 — 手機版的按鈕與面板用原版的底色與外框 | CONFORMED。 | 2026-08-21 |
| [`docs/spec/71-promo-live-capture.md`](spec/71-promo-live-capture.md) | 71 — 桌面版的逐幀錄製（推廣片的動態素材） | CONFORMED。 | 2026-08-21 |
| [`docs/spec/72-bundled-game-data.md`](spec/72-bundled-game-data.md) | 72 — 內含遊戲檔案的四平台完整版 | CONFORMED。四平台的完整包都內含原版資料與倚天字型並實跑驗過； dist-all 因此從「可散布」翻成「私人… | 2026-08-22 |
| [`docs/spec/73-right-click-cancel.md`](spec/73-right-click-cancel.md) | 73 — 右鍵取消是輸入層的語意，不是每個視窗各自的功能 | CONFORMED。七個面板改成問同一支 cancelled()，右鍵與 ESC 都退回上一層。 | 2026-08-23 |
| [`docs/spec/74-corps-on-world-map.md`](spec/74-corps-on-world-map.md) | 74 — 軍團要畫在大地圖上 | CONFORMED。軍團以 MMAP.MCH 的原版圖塊疊在大地圖上，桌面與手機共用同一條算式。 | 2026-08-23 |
| [`docs/spec/75-bundled-audio.md`](spec/75-bundled-audio.md) | 75 — 完整版要出得了聲 | CONFORMED。完整版收 32 個 ogg，wlgame 沒給 -audio 時會自己找執行檔旁邊的 audio/。 | 2026-08-23 |
| [`docs/spec/76-lord-not-in-formation.md`](spec/76-lord-not-in-formation.md) | 76 — 主君能不能編成：系統選單裡的開關 | CONFORMED。系統選單多一列「主君編成」，遊戲中可隨時切； ⚠ 預設「可」——與原版行為不同，是使用者裁定的 … | 2026-08-23 |
| [`docs/spec/77-rout-talk-messages.md`](spec/77-rout-talk-messages.md) | 77 — 敗走的兩段訊息：TALK #1F 與 #23 ＋ 八變體 | CONFORMED。 | 2026-08-23 |
| [`docs/spec/78-amount-input-editor.md`](spec/78-amount-input-editor.md) | 78 — 數值輸入器的上限語意，以及把它接進財政視窗 | CONFORMED。 | 2026-08-23 |
| [`docs/spec/79-new-game-faction-list.md`](spec/79-new-game-faction-list.md) | 79 — 新遊戲的勢力清單 | CONFORMED。 | 2026-08-24 |
| [`docs/spec/80-duel-opening.md`](spec/80-duel-opening.md) | 80 — 開戰單挑：挑戰、拒戰、應戰、回合互嗆、決著 | CONFORMED（2026-08-25）。狀態機在 internal/rules/tactical/duel.go… | 2026-08-25 |
| [`docs/spec/90-same-state-parity.md`](spec/90-same-state-parity.md) | 90 — 同狀態畫面對拍 | CONFORMED。管線接起來了，也對原版跑過一輪 （[../playtest/37](../playtest/37… | 2026-08-15 |
| [`docs/spec/91-tactical-parity.md`](spec/91-tactical-parity.md) | 91 — 戰場的逐區對拍：分區、同狀態怎麼達成 | CONFORMED。 | 2026-08-17 |

## 斷言（欄位／常數 → 推論等級 → 出處）

共 144 條。**要查「這件事解了沒」先看這裡**，
不要重讀整份文件，更不要重推一次。

### confirmed（78 條）

| 鍵 | 出處 |
|---|---|
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 1 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 2 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. ✅ 成功：地形類型的空間關係檢查 ▸ 類型 8 鄰接水 | `docs/playtest/03-verification-log.md` |
| 1. 戰場怎麼被選出來 ▸ 攻城戰 | `docs/mechanics/30-combat.md` |
| 1. 戰場怎麼被選出來 ▸ 野戰 | `docs/mechanics/30-combat.md` |
| 1.4 ⭐⭐ 兩個可直接寫成程式的定義 ▸ 國力 | `docs/mechanics/70-ai.md` |
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
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x02 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x08 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0A | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0C | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x0E | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x10 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x14 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x18 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x19 | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1A | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1B | `docs/formats/08-sinario-save.md` |
| 1.6 據點表：192 筆 × 32 B（區塊 +0x08C0） ▸ +0x1C–+0x1F | `docs/formats/08-sinario-save.md` |
| 1.7 軍團表：127 筆 × 64 B（區塊 +0x22C0） ▸ +0x09 | `docs/formats/08-sinario-save.md` |
| 2. 地形類型對映表（confirmed） ▸ 1 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 2 | `docs/mechanics/30-combat.md` |
| 2. 地形類型對映表（confirmed） ▸ 6 | `docs/mechanics/30-combat.md` |
| 2. 憑什麼說原版不行 ▸ 4 | `docs/spec/76-lord-not-in-formation.md` |
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

### 說明書（10 條）

| 鍵 | 出處 |
|---|---|
| 1. 三個能力值 ▸ +17 | `docs/mechanics/60-personnel.md` |
| 1. 三個能力值 ▸ +18 | `docs/mechanics/60-personnel.md` |
| 1. 三個能力值 ▸ +19 | `docs/mechanics/60-personnel.md` |
| 1. 改了什麼 ▸ 系統視窗 | `docs/playtest/24-window-toggles.md` |
| 3.4 其他戰術判定（說明書） ▸ 中央突破戰法 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 擊破狀態 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 突擊時機 | `docs/mechanics/70-ai.md` |
| 3.4 其他戰術判定（說明書） ▸ 陣形有利不利 | `docs/mechanics/70-ai.md` |
| 5.8 ⭐ 十一個命令處理常式，隊長與隊員各一套 ▸ 命令 | `docs/re/11-tactical-battle.md` |
| 6.1 圖庫段 ＝ ICONGRF 段 3 ＋ 0x9A0 ▸ 場景 0 的 0x1200–0x1440（4 張 24×16） | `docs/re/48-window-display-list.md` |

### 未解（30 條）

| 鍵 | 出處 |
|---|---|
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x0008 | `docs/formats/08-sinario-save.md` |
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x1EC0 | `docs/formats/08-sinario-save.md` |
| 1. 檔案 ＝ 4 個劇本區塊 × 22,208 B ▸ +0x3AC0…+0x42C0 | `docs/formats/08-sinario-save.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D2F8 | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D306 | `docs/re/11-tactical-battle.md` |
| 2. 記憶體佈局：sub_1CC31 ▸ ds:0D42 | `docs/re/11-tactical-battle.md` |
| 3. 曲塊內部 ▸ +0x00 | `docs/re/23-bgm-resource-format.md` |
| 3. 曲塊內部 ▸ +0x06–+0x0F | `docs/re/23-bgm-resource-format.md` |
| 3. 武將記錄（32 byte） ▸ +0 | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +16（0x10） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +20（0x14） | `docs/formats/08-sinario-save.md` |
| 3. 武將記錄（32 byte） ▸ +25,+27 | `docs/formats/08-sinario-save.md` |
| 3. 軍團記錄的欄位（部分） ▸ +0x08 | `docs/re/05-battle-selection.md` |
| 3.2 ⭐ 大地圖是 640×368，四個視窗蓋在它上面 ▸ 0x80 | `docs/re/47-main-screen-window-registry.md` |
| 4. 只在單邊存在的檔（confirmed） ▸ PASS.MAP／PASS.SCH | `docs/re/01-first-recon.md` |
| 4. 據點記錄（32 byte） ▸ +0 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +12 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +14 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +16 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +17 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +18 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +22 | `docs/formats/08-sinario-save.md` |
| 4. 據點記錄（32 byte） ▸ +27–+31 | `docs/formats/08-sinario-save.md` |
| 4. 數字字模在 ICONGRF 段 3 裡 ▸ +0x0000 | `docs/spec/52-main-screen-camera-and-banner-date.md` |
| 4. 數字字模在 ICONGRF 段 3 裡 ▸ +0x08F0 | `docs/spec/52-main-screen-camera-and-banner-date.md` |
| 檔案 ▸ 40-economy.md | `docs/mechanics/00-index.md` |
| 索引 ▸ 事件場景上誰在說話 | `docs/spec/00-index.md` |
| 軍團記錄（64 B，段內 2240h，127 筆） ▸ +0x0C | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x04 | `docs/re/08-hourly-update.md` |
| 連結記錄（16 byte） ▸ +0x08 | `docs/re/08-hourly-update.md` |

