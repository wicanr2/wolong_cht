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

## 怎麼加一份

1. 先確認 `docs/re/` 有出處。**沒有 RE 就沒有規格**——不要從 remake 的
   現況反寫規格，那只會把既有的偏差固定下來。
2. 複製 `TEMPLATE.md`，編號沿用 `docs/mechanics/` 的分類（10 政略、
   20 軍事、30 戰鬥、40 經濟、50 外交、60 人事、70 AI、80 勝負）。
3. 「remake 實作」欄留白就是缺口，**不要為了讓表好看而填近似值**——
   近似要寫成「近似」並說明差在哪。
