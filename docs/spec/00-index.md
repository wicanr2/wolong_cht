# 00 — 規格索引：已解的規則有沒有被實作、有沒有被驗過

**狀態：索引。規格是 `docs/re/`（程式碼在哪）與 `internal/`（我們寫了什麼）
之間的那一層——它回答「這條規則實作了沒、驗過沒」。**

- 日期：2026-08-14

## 為什麼要有這一層

三份既有文件都不回答這個問題：

| 文件 | 回答 |
|---|---|
| `docs/re/` | 原版的程式碼在哪、怎麼寫的 |
| `docs/mechanics/` | 這個遊戲怎麼運作 |
| `docs/INDEX.md` | 某個欄位／常數解了沒 |

**沒有一份回答「解出來的東西進到 remake 了嗎」。** 而那正是差距的來源：

> `formAICorps` 的門檻寫成「已經有一支軍團就不再編成」，
> 而原版是 `max(5, 資金 ÷ 8192)` 扣掉現有數。AI 因此永遠只有一支軍團、
> 資金再多也不擴軍。**這個 bug 不會讓任何測試變紅**——測試驗的是
> 我們寫下的行為，不是原版的行為。它是靠回頭讀 `sub_14575` 才發現的。

規格就是為了讓這種「靜態解出來了、但沒接進去」在**清單上看得見**。

## `[HARD]` 動工前先寫，polish 也算

**任何會改到 `internal/` 或 `cmd/` 的工作，動手之前先寫一份規格**
（`CLAUDE.md` §10，使用者裁定 2026-08-14）。

**「只是調一個常數」不是例外。** 版面常數正是最容易照參考影片猜的東西，
而猜出來的值與原版差幾像素，肉眼看不出來、測試也驗不到——
測試驗的是我們自己寫下的契約。**規格的第一格是「出處」，
光是要填那一格就會逼你先去找原版怎麼寫的。**

## 一份規格長什麼樣

見 [`TEMPLATE.md`](TEMPLATE.md)。五個欄位缺一不可：

| 欄位 | 為什麼要 |
|---|---|
| **出處** | 指回 `docs/re/` 的哪一節、哪個位址。沒有出處的規格是憑印象寫的 |
| **算式／流程** | 照抄原版，不整理、不簡化。原版的怪癖（死碼、位移不對稱）要標明「照抄」 |
| **推論等級** | 四級，見 `CLAUDE.md` §9。不確定的結論也可以實作，但要在規格裡標明等級 |
| **remake 實作** | Go 的檔案與函式。**空白 ＝ 已解但沒接**，那是缺口不是完成 |
| **驗證** | 單元測試名稱，或動態取樣紀錄（`docs/playtest/21`）。**測試綠只證明我們自洽**，要對原版就得有取樣 |

狀態只有三種：

- `DRAFT`——還在寫，不能照著實作
- `READY`——可以實作
- `CONFORMED`——已實作**而且**驗過（單測 ＋ 至少一項對原版的證據）

## 索引

| 主題 | 規格 | 狀態 |
|---|---|---|
| 亂數消費順序對齊（長時間對拍前置）| [`160-rng-consumption-alignment.md`](160-rng-consumption-alignment.md) | DRAFT；原版側軌跡已取得，remake 側與比對工具未做 |
| 劉備不作為長時間對拍管線 | [`161-liubei-idle-longrun-parity.md`](161-liubei-idle-longrun-parity.md) | DRAFT；原版側 90 天已跑通（26 秒、兩次逐 byte 相同），remake 側未做 |
| 戰術命令批次與移動體力順序 | [`159-tactical-command-movement-order.md`](159-tactical-command-movement-order.md) | READY；窄批次修正已測，完整排程仍待對拍 |
| 委任將領混合值保留原統率 | [`158-delegated-leader-mix.md`](158-delegated-leader-mix.md) | READY；完整驗收進行中 |
| 據點整備、威脅偵測與求援（`sub_13EFD` 鏈）| [`10-city-tick.md`](10-city-tick.md) | 已實作並對原版取樣驗過 |
| 進言「請求君主出陣」（`sub_1699E`）| [`11-ai-sortie.md`](11-ai-sortie.md) | 已實作並有單測；兩道閘都從機器碼讀出來 |
| 主畫面的視窗外框與指令列 | [`12-strategy-chrome.md`](12-strategy-chrome.md) | 版面與各視窗內部排版都照機器碼；主畫面五區逐像素對過 |
| 主畫面四個視窗的開關 | [`13-main-window-toggles.md`](13-main-window-toggles.md) | 已實作；主畫面逐像素對過。原版執行期的**開關行為**仍未驗 |
| 編成時預備兵怎麼分配（`sub_14698`）| [`21-corps-formation-reserves.md`](21-corps-formation-reserves.md) | 已實作並有逐項單測 |
| 財政視窗 | [`14-finance-window.md`](14-finance-window.md) | 版面已照原版重寫；數值輸入器已接（[`78`](78-amount-input-editor.md)） |
| 軍團編成視窗 | [`22-corps-formation-window.md`](22-corps-formation-window.md) | 版面、武將頭像與六個槽的滑鼠熱區都照原版 |
| 據點情報視窗 | [`23-city-info-window.md`](23-city-info-window.md) | 版面已照原版實作 |
| 軍團情報視窗 | [`24-corps-info-window.md`](24-corps-info-window.md) | 版面已照原版實作；指令流程走 [`39`](39-march-order-menu.md) 的選單，入口與原版不同（§5） |
| 四槽選擇視窗 | [`25-slot-select-window.md`](25-slot-select-window.md) | 讀取／儲存已照原版；新遊戲未共用 |
| ＹＥＳ／ＮＯ 對話框 | [`26-yes-no-dialog.md`](26-yes-no-dialog.md) | 版面與命中算式已照原版 |
| 君主選擇視窗 | [`27-lord-select-window.md`](27-lord-select-window.md) | 版面已照原版；「自定」未接 |
| 劇本 JSON | [`28-scenario-json.md`](28-scenario-json.md) | 匯出／匯入／round-trip 已可用（`cmd/wlscen`）|
| 音樂與音效 | [`29-audio.md`](29-audio.md) | 已實作；場景對應已解（`docs/re/58`），音色的諧波結構未量化比對 |
| remake 原生存檔格式 | [`20-save-format.md`](20-save-format.md) | 已接進遊戲並驗過；**只差放回 DOSBox 實測** |
| 戰術側欄的內容組成 | [`31-tactical-sidebar.md`](31-tactical-sidebar.md) | 七格已照原版實作；`▶▶` 列只畫美術不接行為 |
| 攻城的「門強度」條 | [`32-gate-strength-bar.md`](32-gate-strength-bar.md) | 已實作並有單測；右鍵提前收掉未接 |
| 底列六格是選部隊 | [`33-squad-selection.md`](33-squad-selection.md) | 已實作並有單測；六張命令圖示照 `ICONGRF` 段 3 的 `碼 × 0xC0` |
| 一覽表的欄位與版面 | [`38-list-windows.md`](38-list-windows.md) | 四個家族、捲軸的四個熱區與欄寬定義都照原版 |
| 戰術畫面的玩家操作 | [`37-tactical-player-controls.md`](37-tactical-player-controls.md) | 陣形選單與陣形線已接；說明書 4.2–4.6 的功能逐條對照過 |
| 兩個平面的地面圖與登城 | [`36-ground-planes-and-climbing.md`](36-ground-planes-and-climbing.md) | 已實作並有單測，拿三張原版攻城圖對過數字 |
| 縮小地圖的據點標記 | [`35-strategy-minimap.md`](35-strategy-minimap.md) | 已實作並有單測；22 勢力的選擇視窗用「點一下換下一個」代替 |
| 兩個速度設定的五檔 | [`34-speed-steps.md`](34-speed-steps.md) | 已實作並有單測；「最高速」的上限是 remake 差異 |
| 行軍指示的三選一 | [`39-march-order-menu.md`](39-march-order-menu.md) | 已接進畫面並有單測 |
| AI 軍團的決策鏈 | [`40-ai-march-decision.md`](40-ai-march-decision.md) | 已實作；「逐站前進」未移植 |
| 訊息框的版面 | [`41-message-box-geometry.md`](41-message-box-geometry.md) | 機器碼與影格兩條證據都對上，已實作 |
| 事件場景上誰在說話 | [`42-event-scene-speakers.md`](42-event-scene-speakers.md) | 兩個框已實作；結果句走的是一般通知框（`sub_13C3D` → `sub_18810`），不是上下框，肖像取派駐外交官（§2） |
| 回不了家就敗走 | [`43-rout-on-blocked-return.md`](43-rout-on-blocked-return.md) | 已實作並有單測 |
| 進言的原文 | [`44-advise-original-text.md`](44-advise-original-text.md) | 三個指令 × 64 則全部改查 `TALK.DAT` |
| 進言的畫面 | [`45-advise-scene-layout.md`](45-advise-scene-layout.md) | 插圖 ＋ 兩個框 ＋ 五列選單已實作並留下截圖 |
| 戰後退一站回家 | [`46-post-battle-retreat.md`](46-post-battle-retreat.md) | 已實作並有單測 |
| 據點易主之後守軍調頭 | [`47-city-fall-corps-redirect.md`](47-city-fall-corps-redirect.md) | 已實作並有單測 |
| 據點被攻陷時內政官被遣回 | [`48-governor-returns-on-city-fall.md`](48-governor-returns-on-city-fall.md) | 已實作並有單測，兩則訊息也接了 |
| 進言第四、五項（遷都／請求出陣）| [`49-advise-relocate-and-sortie.md`](49-advise-relocate-and-sortie.md) | 已接進畫面並留下截圖 |
| 軍費直接扣資金 | [`50-corps-upkeep-charges-funds.md`](50-corps-upkeep-charges-funds.md) | 已改正並有單測；長跑數字跟著變 |
| 結局：存活勢力數歸一 | [`30-victory.md`](30-victory.md) | 判定與兩則結局訊息已實作並驗過；結局的過場見 [`67`](67-ending-playback.md) |
| 顏色到不了滿刻度（6 bit DAC）| [`51-vga-dac-palette-scale.md`](51-vga-dac-palette-scale.md) | 已全面套用（介面顏色一律查調色盤）；主畫面逐像素對過 |
| 開局鏡頭與橫幅日期 | [`52-main-screen-camera-and-banner-date.md`](52-main-screen-camera-and-banner-date.md) | 已實作，主畫面逐像素對過 |
| 據點圖塊跟著歸屬換 | [`53-city-tile-by-ownership.md`](53-city-tile-by-ownership.md) | 已實作，主畫面逐像素對過 |
| 介面顏色一律查調色盤 | [`54-ui-colours-from-palette.md`](54-ui-colours-from-palette.md) | 已實作 |
| 縮小地圖的視野框 | [`55-minimap-view-box.md`](55-minimap-view-box.md) | 已實作 |
| 戰場轉 180 度 | [`56-battlefield-rotation.md`](56-battlefield-rotation.md) | 已實作；地形、小地圖、子圖塊與**旗**都餵同一份 |
| 戰術的等角投影 | [`57-tactical-projection.md`](57-tactical-projection.md) | 已實作；地形走出來、物件算出來 |
| 顯示格的深度範圍與 8 列的帶 | [`58-display-slot-depth-range.md`](58-display-slot-depth-range.md) | 已實作 |
| 開場的常令 | [`59-battle-opening-orders.md`](59-battle-opening-orders.md) | 已實作並有單測 |
| 戰場對白的壽命與每側一個框 | [`60-battle-talk-duration.md`](60-battle-talk-duration.md) | 已實作並有單測 |
| 兵的開場體力 ＝ 軍團士氣 | [`61-soldier-initial-hp-from-morale.md`](61-soldier-initial-hp-from-morale.md) | 已實作；要與 `62` 一起看 |
| 被換位的兵這一幀不動 | [`62-swapped-unit-skips-its-turn.md`](62-swapped-unit-skips-its-turn.md) | 已實作並有迴歸閘 |
| 挨打的三幀硬直 | [`63-hit-stun.md`](63-hit-stun.md) | 已實作並有單測 |
| 遷都之後說什麼 | [`64-capital-relocation-report.md`](64-capital-relocation-report.md) | 已實作並有單測；他國遷都要有外交官才報得回來 |
| 退到畫面外的兵算生還 | [`65-retreated-soldiers-survive.md`](65-retreated-soldiers-survive.md) | 已實作並有單測 |
| 打壞的城壁與門要在畫面上換掉 | [`66-broken-walls-repaint.md`](66-broken-walls-repaint.md) | 已實作並有單測；繪圖層跟著規則層的圖塊版本走 |
| 結局的播放 | [`67-ending-playback.md`](67-ending-playback.md) | 已實作並有單測；十二幕 ＋ 逐字文字 ＋ 十七階淡入淡出。**三段的節拍各自照原版重做**，整段 3 分 21 秒（[`67`](67-ending-playback.md) §8）|
| 倒地動畫 | [`68-death-animation.md`](68-death-animation.md) | 已實作並有單測；四幀、三個兵種組、後兩幀換第二張 |
| 世界指紋（**remake 設施，無原版出處**）| [`69-world-fingerprint.md`](69-world-fingerprint.md) | 已實作並有單測（含 15 欄的正對照）；決定性迴歸 ＋ Android 里程碑 A 的判準 |
| 手機版的底色與外框取自原版 | [`70-phone-chrome.md`](70-phone-chrome.md) | 已實作並有單測；顏色與外框與桌面版共用 `internal/ui/chrome` |
| 桌面版逐幀錄製（**remake 工具，無原版出處**）| [`71-promo-live-capture.md`](71-promo-live-capture.md) | 已實作；推廣主片的大地圖與兩場戰鬥改成實跑錄製 |
| 內含遊戲檔案的四平台完整版（**發行設施，無原版出處**）| [`72-bundled-game-data.md`](72-bundled-game-data.md) | 已實作並實跑驗過；`dist-all` 是私人批次，不可外流 |
| 右鍵取消是輸入層的語意 | [`73-right-click-cancel.md`](73-right-click-cancel.md) | 已實作並有單測；七個面板共用同一支 `cancelled()` |
| 軍團要畫在大地圖上 | [`74-corps-on-world-map.md`](74-corps-on-world-map.md) | 已實作並有單測；用 `MMAP.MCH` 的原版圖塊，桌面與手機共用算式 |
| 完整版要出得了聲（**發行設施**）| [`75-bundled-audio.md`](75-bundled-audio.md) | 已實作；沒給 `-audio` 時自己找執行檔旁邊的 `audio/` |
| 主君能不能編成 | [`76-lord-not-in-formation.md`](76-lord-not-in-formation.md) | 已實作並有單測；⚠ **預設「可」是 remake 差異**，原版不能 |
| 敗走的兩段訊息 | [`77-rout-talk-messages.md`](77-rout-talk-messages.md) | 已實作並有單測；#1F 與 #23 ＋ 組 `0x198`。⚠ 敗走走不到 #20 |
| 數值輸入器的上限語意 ＋ 財政 | [`78-amount-input-editor.md`](78-amount-input-editor.md) | 已實作並有單測；`sub_17C6E` 的 `ax` 是**上限**不是初值，錨點由呼叫端給 |
| 新遊戲的勢力清單 | [`79-new-game-faction-list.md`](79-new-game-faction-list.md) | 已實作並有單測；視窗 (136,104)、五欄、一頁 10 列。⚠ 點擊路徑無頭驗不到（§3.1）|
| 開戰單挑：挑戰、拒戰、應戰、回合互嗆、決著 | [`80-duel-opening.md`](80-duel-opening.md) | 已實作並有單測；狀態機照機器碼 |
| 災害的實際數值：機率、marker 量、持續與距離衰減 | [`81-disaster-quantities.md`](81-disaster-quantities.md) | 已實作並有單測；機率、marker 量、距離衰減都有出處 |
| 應戰軍團的挑選：兵數 × 士氣 × 評價 | [`82-defender-selection.md`](82-defender-selection.md) | 已實作並有單測；兵數 × 士氣 × 評價 |
| 新遊戲的開局政略評估（sub_12BD9 的第二個呼叫點） | [`83-initial-strategy-pass.md`](83-initial-strategy-pass.md) | 已實作並有單測；`sub_12BD9` 的開局呼叫點，孫策攻劉繇的分歧因此收掉 |
| 多語系：簡體中文、日文、英文 | [`84-multilanguage.md`](84-multilanguage.md) | 四個語系端到端可玩；簡體與英文已第二人覆核（§6） |
| 半形語系的清單欄界（英文版的姓名欄） | [`85-latin-list-layout.md`](85-latin-list-layout.md) | 已實作並實跑驗過；半形語系的清單欄界另排 |
| 執行期切換語言（含手機版） | [`86-runtime-language-switch.md`](86-runtime-language-switch.md) | 已實作並實跑驗過；F9／殼層／手機面板三個入口 |
| 半形語系的畫面調整（清單以外） | [`87-latin-screen-layout.md`](87-latin-screen-layout.md) | 已實作並實跑驗過；原版美術上的中文不翻（§2） |
| 三處顯示與原版對不上 | [`88-display-polish-parity.md`](88-display-polish-parity.md) | 已實作並實跑驗過 |
| 戰後的損害報告改成可關的選項 | [`89-siege-damage-report-toggle.md`](89-siege-damage-report-toggle.md) | 已實作並實跑驗過；remake 差異，預設照原版 |
| Android 也要有原版的音樂 | [`92-android-music.md`](92-android-music.md) | 已實作並實跑驗過；APK 的音檔走 `ImportActivity` |
| 攻城「一撞歸零」的面向常數要跟著戰場翻轉 | [`93-siege-wall-instant-break-facing.md`](93-siege-wall-instant-break-facing.md) | 已實作並有單測；面向跟著戰場翻轉 |
| 退卻的繞路點不可以每幀清掉 | [`94-retreat-path-not-cleared-every-frame.md`](94-retreat-path-not-cleared-every-frame.md) | 已實作並有單測 |
| 開場擺兵的高度要用地面層表，不是堆疊高度 | [`95-spawn-height-uses-ground-plane.md`](95-spawn-height-uses-ground-plane.md) | 已實作並有單測 |
| 守陣不可以在回陣的那一步被降級成「就位」 | [`96-guard-command-not-downgraded.md`](96-guard-command-not-downgraded.md) | 已實作並有單測 |
| 登城的觸發：X 與 Y 都走不動就試 Z，不必先走到目標格 | [`97-climb-when-both-axes-blocked.md`](97-climb-when-both-axes-blocked.md) | 已實作並有單測 |
| 爬不上去的那一下要打門：未破的門是這樣被打開的 | [`98-climb-into-a-gate-hits-it.md`](98-climb-into-a-gate-hits-it.md) | 已實作並有單測 |
| 同狀態畫面對拍（方法）| [`90-same-state-parity.md`](90-same-state-parity.md) | 主畫面五區逐像素相同；§4.1 記下「參考影格本身會有東西」 |
| 戰場的逐區對拍（分區）| [`91-tactical-parity.md`](91-tactical-parity.md) | 九區裡六區逐像素相同（2026-08-18）；⚠ 取樣點已不等價，見 `playtest/49` |
| 手機版「關於」頁顯示授權條款（**remake 差異，發行設施**）| [`99-about-page-license.md`](99-about-page-license.md) | 已實作並有單測；APK 帶不了 `LICENSE` 檔，摘要顯示在遊戲內 |
| 手機版的字放大 2 倍（**remake 差異**）| [`100-phone-text-scale.md`](100-phone-text-scale.md) | 已實作、單測與截圖驗過；版面常數從字高長出來，桌面版倍率不變 |
| 手機版放大的字用 Scale2x 去鋸齒（**remake 差異**）| [`101-phone-glyph-scale2x.md`](101-phone-glyph-scale2x.md) | 已實作、單測與放大對照過；只在倍率 2 套，桌面版不經過 |
| 戰場的 `▶▶` 是快轉 | [`102-battle-fast-forward.md`](102-battle-fast-forward.md) | 已實作、單測、實機對照過行為；底紋未對上（§5）|
| 手機版外交提案的「提示金額」鍵盤（**remake 差異**）| [`103-phone-diplomacy-amount-keypad.md`](103-phone-diplomacy-amount-keypad.md) | 已實作並有單測；輸入動作與上限沿用 `78` |
| 「自定」軍師命名視窗 | [`104-advisor-naming-window.md`](104-advisor-naming-window.md) | 已實作、單測與截圖；肖像位置是假說（§5）|
| 遭遇直接進戰場，沒有「戰鬥指揮／委任」選單 | [`105-encounter-goes-straight-to-battle.md`](105-encounter-goes-straight-to-battle.md) | 已實作並有單測；機器碼 ＋ 實機兩條證據。⚠ 遭遇訊息本身 remake 還沒有（§5）|
| 訊息框的臉是固定的通報者 | [`106-message-box-reporter-portrait.md`](106-message-box-reporter-portrait.md) | 已實作並有單測；順帶解掉「同一個標記出現兩次是兩個值」|
| 啟動殼層的 UI 顏色也要查調色盤 | [`107-launcher-ui-colours.md`](107-launcher-ui-colours.md) | 已實作、單測（含正對照）；君主卡 0 px、勢力清單本體 0 px。捲軸滑塊差 1 px 未解（§7）|
| 進言問理由之前君主要先講那一句 | [`108-advise-ask-reason-line.md`](108-advise-ask-reason-line.md) | 已實作並有單測；實機對照過（`../playtest/56` §4.3）|
| 編成成功之後主將要講一句 | [`109-formation-leader-line.md`](109-formation-leader-line.md) | 已實作並有單測；主公型那三格是空的，取到空字串不開框 |
| 指令列的「軍團」是兩項彈出選單 | [`110-corps-command-menu.md`](110-corps-command-menu.md) | 已實作並有單測；選單框本身還沒逐像素比（§5）|
| 君主帶著軍團時進言關掉（**remake 差異**）| [`111-lord-with-corps-blocks-advise.md`](111-lord-with-corps-blocks-advise.md) | 使用者裁定 2026-09-01；判準與「請求君主出陣」共用一支 |
| **游標停下之後的恢復延遲**（即時制的反應時間）| [`112-cursor-idle-resume-delay.md`](112-cursor-idle-resume-delay.md) | 已實作並有單測；游標移動中世界完全停住，停下後等 160 個回呼（0.549 秒）|
| **武將的心向勢力**（`+0x19`）：在野出仕與俘虜歸降 | [`114-general-affinity.md`](114-general-affinity.md) | 已實作並有單測；在野武將每月 25% 兌現，俘虜要關押方就是心向的勢力才歸降。隨機投靠那一條還沒接 |
| **兵的戰力來自統率力**（不是士氣）| [`115-soldier-power.md`](115-soldier-power.md) | 已實作並有單測；同一場攻城的勝負跟著翻面。戰術九區對拍待重跑 |
| 驗收戰場少了子圖塊表，打破的門反而封城 | [`116-retreat-cannot-leave-the-city.md`](116-retreat-cannot-leave-the-city.md) | 已修：fixture 改用 `NewFieldFromTileLayers`。**正式路徑本來就沒問題** |
| **RLE 資料檔的 4 byte 長度頭** | [`113-rle-length-header.md`](113-rle-length-header.md) | 原版三個執行檔都 `LSEEK` 跳過它才解壓；`rle.DecodeFile` 已接，19 個過場檔逐檔解到宣告長度 |
| 驗收捷徑要先武裝開場喊話再推戰場 | [`117-fixture-arms-duel-before-stepping.md`](117-fixture-arms-duel-before-stepping.md) | 已修並有突變測試；野戰對拍的 `field` 從 11.24% 回到 95 px |
| 截圖的時機用局面條件（`-shot-when`／`-auto-messages`）| [`118-shot-when-condition.md`](118-shot-when-condition.md) | 驗收設施，預設關閉。自然流程與捷徑截出的畫面**九區逐像素相同**；條件不成立時不寫檔、回非零 |
| **`\1`／`\4` 代入的是呼び名，不是姓名** | [`119-talk-marker-fields.md`](119-talk-marker-fields.md) | 已實作、單測與突變測試。四個劇本裡受影響的只有四人，而三個是名人（諸葛亮→孔明、司馬懿→仲達、龐統→鳳雛）|
| **尋路是反應式的**：全域佇列每幀兩個兵 | [`120-pathfind-request-queue.md`](120-pathfind-request-queue.md) | 已實作並有單測。取代 `replanInterval` 那個明示的 remake 差異；三場迴歸戰仍打得完 |
| 碼頭與棧道的野戰用錯戰場（`SelectWater` 沒接）| [`121-water-battlefield-selection.md`](121-water-battlefield-selection.md) | 已修並補值域斷言。214／215 超出 0–213，呼叫端一直退回合成戰場 |
| 音效的 TYPE 1–4 是四段主衰減 | [`122-sound-type-levels.md`](122-sound-type-levels.md) | 已接。原版每段加 4 個 OPL TL 單位（約 3 dB）；remake 播 OGG，接成播放增益 |
| 武將下場的五則訊息（脫身／被擒／自刎）| [`123-captive-talk-messages.md`](123-captive-talk-messages.md) | 已接。敗方與勝方看到不同的一則，兩邊都不是就不出；被擒那一則比的是**舊主** |
| 反白是色號 XOR 12（選單列 ＋ 指令列格）| [`124-menu-highlight-xor.md`](124-menu-highlight-xor.md) | 已修。黑底白字 ⇒ 黃底藍字；指令列先前根本沒畫反白 |
| 選單框寬要用沒切補位空白的那一行 | [`125-menu-box-width-from-padding.md`](125-menu-box-width-from-padding.md) | 已修。標籤改讀 `TALK #79`，解碼走 `DecodeKeepPad` |
| 指令列的四個彈出選單走同一支常式 | [`126-command-popup-menus.md`](126-command-popup-menus.md) | 已接三張（人事／軍團／據點）。人事從一覽表換成選單，據點多了「首都確認」 |
| 被俘的主公型武將當場變成臣下型 | [`127-captured-sovereign-becomes-retainer.md`](127-captured-sovereign-becomes-retainer.md) | 已接。清旗標 bit 6 ＋ 說話類型 `+3`；⚠ 順帶訂正「以舊主已滅為條件」那個說法 |
| 隊長倒下不會清掉那一隊的待機兵 | [`128-squad-leader-gone-keeps-reserve.md`](128-squad-leader-gone-keeps-reserve.md) | 已改。原版沒有清待機的地方，而且「隊長不在」是**每幀重新施加**的 |
| 戰術戰鬥打完，兩側的士氣都要按兵力比縮 | [`129-post-battle-morale-scaling.md`](129-post-battle-morale-scaling.md) | 已改。remake 兩行的**分母都寫成新值**，等於沒有縮 |
| 沒有心向的在野武將：隨機投靠，偏向武將最少的勢力 | [`130-freelance-random-join.md`](130-freelance-random-join.md) | 已接。⭐ **不經過 25% 那道閘**；骰面 ≥ 48 是玩家專屬的救濟 |
| 原版 oracle 換成 dosgolem，不再需要 DOSBox | [`131-dosgolem-oracle.md`](131-dosgolem-oracle.md) | 已接 `tools/dosgolem.sh`。五格 × 五區全 0 px，整條鏈 0.99 秒；⭐ 即時制的取樣點可以寫成**遊戲日期** |
| **行軍走到和平勢力的邊界要掉頭** | [`132-march-turnback-at-peace.md`](132-march-turnback-at-peace.md) | CONFORMED。⭐ **和平就出不了兵**：`sub_142AB` 在野外每一步都問前方據點是誰的，是和平勢力就折返。已實作，三組突變測試 |
| **開場擺位：邊界那一欄 ＋ 亂數 Y** | [`133-opening-deployment.md`](133-opening-deployment.md) | CONFORMED。⭐ 原版把每個兵放在戰場邊界（側 0 `X=1`、側 1 `X=62`），`Y ＝ 亂數 & 0x1F + 0x10`，**不查佔用**；之後才走進陣形。remake 已接上，量到 X 全在邊界欄、Y 16–47、44 個兵重疊 |
| **尋路是波數佇列，不是先進先出** | [`134-pathfind-wave-order.md`](134-pathfind-wave-order.md) | CONFORMED。⭐ 成本 ≥ 目前波數的格子要**推回佇列**，有兵的格子等於多躺八波；純 FIFO 之下「繞路成本 8」從來沒有生效過。修好之後同一場第 70 拍**96 個兵逐槽全等**，小地圖 8 px → 0 px |
| **腳本指令 16 是「掛一個對白框」** | [`135-script-message-command.md`](135-script-message-command.md) | CONFORMED。攻城戰開場的勸降對白：TALK 組 ＝ `0x1CE + 運算元`、側別 ＝ 參數 bit 0，另有一道「現在是哪個收尾階段」的閘。remake 的 `opMessage` 是空的 stub，於是那兩個框一個都沒畫 |
| **對白框的參數是一條共用串流** | [`136-battle-talk-parameters.md`](136-battle-talk-parameters.md) | CONFORMED。`sub_1C315` 推 `[對手, 說話者]` 兩個參數，而 **`\6` 也吃一個**——有 `\6` 的變體 `\1` ＝ 說話者，沒有的 `\1` ＝ 對手。變體用武將 `+0x1E` 的**原始值 0–7**，不是進言那條路的 0–2 |
| **全形標點用原版內建的 408 格** | [`137-builtin-symbol-font.md`](137-builtin-symbol-font.md) | CONFORMED。`END_S13.DAT` 開頭的 408 格才是遊戲用的全形符號字型，與倚天 `SPCFONT.15` **不同**（逗號差 (+3,−2)）。順帶讓繁中不必自備倚天字型 |
| **狀態層對拍：另外三張表** | [`138-state-table-parity.md`](138-state-table-parity.md) | CONFORMED。**5,502 個欄位逐欄比對通過**。勢力／據點／武將表逐欄比，作法照搬軍團表那一輪。名字兩邊都印 hex（dosgolem 零相依），補白是全形空格 `A1 40` |
| **AI 決策軌跡對拍** | [`139-ai-decision-trace.md`](139-ai-decision-trace.md) | CONFORMED。第一輪比過：事件種類、參數形狀與每月量級一致。攔 `sub_12FBF`（所有事件的共用出口）取軌跡，比**種類／發起方／節奏**不比數值——亂數只影響時機與對象，種類與產生條件是規則決定的 |
| **左下角的狀態列提示框** | [`140-status-message-box.md`](140-status-message-box.md) | CONFORMED。`sub_18853` 每個指令流程都掛的那個框 ＝ `(0, 320, 256, 80)`、肖像 `0x93`，與一般訊息框是**同一個框換位置**。已接三條：行軍（[`39`](39-march-order-menu.md) §3.7）、財政（#16）、編成（#0）|
| **全軍退卻之後的 120 拍倒數** | [`141-retreat-countdown.md`](141-retreat-countdown.md) | CONFORMED。`sub_1A6FA` 的三條出口裡排最前面的那一條：退卻中每拍減 1，走完由**沒退卻的那一側**獲勝。⭐ 實測正常打完走的是「補不出兵」——退卻的兵八拍就走完了，倒數是**兜底** |
| **解任不先過濾** | [`142-personnel-dismiss-flow.md`](142-personnel-dismiss-flow.md) | CONFORMED。原版的內政官／外交官解任**照列全部**，選到沒派人的才跳 TALK #54／#55。⭐ `sub_16B4F` 是**先寫 0xFF 再看舊值**；成功或失敗都回清單繼續選 |
| **武將職務 `+0x17`** | [`143-general-duty-field.md`](143-general-duty-field.md) | CONFORMED。`+0x17` 是 0–4 的職務值（`－－－`／軍團長／內政官／外交官／俘虜），一覽表的身分欄直接查 `cs:75A4h`；`5 ＝ 君主`是「職務 0 ＋ bit 6」臨時算的。remake 收成 `Posted bool`，任命不寫職務、解任不清經費 |
| **選了軍師就從武將表消失** | [`144-advisor-leaves-general-table.md`](144-advisor-leaves-general-table.md) | CONFORMED。新遊戲定案時 `loc_11AF8` 把選中的軍師記錄 `+0x00` 寫 0（存在旗標一起沒了）並把勢力武將數減一，所以他不出現在任何清單裡。⭐ **不是寫 `+0x17`**——`spec/76` §2 那條「寫入者未讀」的缺口就此關閉 |
| **指令列最後兩格** | [`145-general-and-faction-cells.md`](145-general-and-faction-cells.md) | CONFORMED。武將那格是**迴圈**，選完那位自陳擅長哪一種戰場（俘虜／城塞／野戰／海戰四組，由三個適性挑，**平手歸前面那一個**）；勢力那格把鏡頭移到該勢力首都並開那個據點的情報卡。兩格各有狀態列 #24／#25。八格到此全部有原版擷取 |
| **大地圖上會飄的雲** | [`146-map-cloud-objects.md`](146-map-cloud-objects.md) | CONFORMED。`MMAP.MCH` 物件 **type 0** ＝ 16 朵常駐的雲（16×9 格），座標存在劇本與存檔的 `0x21C0`，每次 map-loop 移動、起暴風雨才被關進那 11×11 格。⭐ 順帶訂正**物件型別查表整體差一格**——remake 先前把火災畫成雲的圖形 |
| **受控存檔** | [`147-controlled-parity-save.md`](147-controlled-parity-save.md) | CONFORMED。`tools/parity_save.py` 從既有存檔出發只改指定欄位，做出兩邊都載得進去的受控起點。⭐ **存檔本身也可以是實驗器材**——`--no-clouds` 把會飄的雲關掉之後，跑滿一個遊戲日的 `map` 殘差從 22,210 掉到 66（只剩原版的地圖游標框）|
| **共用候選過濾** | [`148-shared-candidate-filter.md`](148-shared-candidate-filter.md) | CONFORMED。內政官任命／外交官任命／編成選武將三條流程在原版是**同一支 `sub_17663`**：存在 ∧ 勢力＝玩家 ∧ 職務＝0 ∧ 不是君主。⭐ 原版沒有「俘虜」也沒有「軍師」那一關——**兩個都由別的機制吃掉**（職務 4／整筆不在武將表裡）|
| **行軍目標的地圖選點** | [`149-march-target-map-picker.md`](149-march-target-map-picker.md) | CONFORMED。原版選目標據點是**在大地圖上點**（`sub_1703C`），不是一覽表；選完軍團先開軍團情報面板。⭐ 第一次把原版自己畫的游標（15×15 空心框 ＋ 黑影）畫對——[`../playtest/94`](../playtest/94-march-map-picker.md) 五區全 0 px |
| **外交的兩道前置閘** | [`150-diplomacy-preconditions.md`](150-diplomacy-preconditions.md) | CONFORMED。停戰與請求協助**要先派外交官到對方**（`sub_165EF`），而且同型的使者不能已經在路上。⚠ remake 的規則層把外交官那一條寫反了（要求「沒有」），而且有單元測試把錯的行為釘住——[`../playtest/95`](../playtest/95-diplomacy-preconditions.md) |
| **大地圖點擊** | [`151-map-click.md`](151-map-click.md) | CONFORMED。原版的大地圖左鍵是一條完整的分派（`sub_11E46`）：據點開情報卡、軍團開情報面板、兩者都有先跳一張兩項選單（TALK #80）。remake 先前點地圖什麼都不做 |
| **畫面模式（液晶）** | [`152-video-mode-lcd-palette.md`](152-video-mode-lcd-palette.md) | CONFORMED。系統選單第 1 列切的是 `GAMEPAL.BRG` 的 bank 0–3 ↔ 4–7。⭐ **bank 就是調色盤不是美術**，所以整個功能＝把所有取 bank 的地方換成一支 `paletteBank()`；順帶抓到「視窗外框開局載一次就不再重畫」——換季也一樣漏 |
| **「遊戲結束」的確認選單** | [`153-quit-confirm-menu.md`](153-quit-confirm-menu.md) | CONFORMED。系統選單第 6 列跳的是 TALK #81 兩項選單「終　了」／「取　消」，⭐ **位置寫死在 (352, 272)**——六個 handler 裡唯一不跟游標走的。F10 那條保留 remake 的 ＹＥＳ／ＮＯ（[`26`](26-yes-no-dialog.md)）|
| **原版的滑鼠游標** | [`154-mouse-cursor.md`](154-mouse-cursor.md) | CONFORMED（對拍端）。逐像素殘差裡那 95 px 就是它：14×14 紅箭頭，圖樣從兩張不同背景的原版擷取抽出來。接上 `-cursor X,Y` 之後**請求協助第二步那一張整個 640×400 一個像素都不差**。⚠ 遊玩端還不自繪——原版在清單等待時不畫游標，而那是不是 oracle 的限制沒定案 |

| **桌面啟動選章與兩段確認** | [`155-desktop-launcher-input.md`](155-desktop-launcher-input.md) | CONFORMED。日期欄與命中共用幾何，勢力清單兩段確認與右鍵分層返回；Android 不在本輪範圍 |
| **桌面偏好保存** | [`156-desktop-preferences.md`](156-desktop-preferences.md) | CONFORMED。正常啟動還原使用者偏好，原版存檔與 Android 不變；明示參數優先，直接驗收入口不讀寫偏好。 |

| **戰術 AI 分支比較** | [`157-battle-script-comparisons.md`](157-battle-script-comparisons.md) | CONFORMED。參數 3／4 為無號大於等於／小於等於；修正呂布兵力充足卻誤退卻。 |

## 怎麼加一份


1. 先確認 `docs/re/` 有出處。**沒有 RE 就沒有規格**——不要從 remake 的
   現況反寫規格，那只會把既有的偏差固定下來。
2. 複製 `TEMPLATE.md`，編號沿用 `docs/mechanics/` 的分類（10 政略、
   20 軍事、30 戰鬥、40 經濟、50 外交、60 人事、70 AI、80 勝負）。
3. 「remake 實作」欄留白就是缺口，**不要為了讓表好看而填近似值**——
   近似要寫成「近似」並說明差在哪。
