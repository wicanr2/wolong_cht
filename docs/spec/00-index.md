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
| **推論等級** | confirmed／強證據／假說／未知（`CLAUDE.md` §9）。假說也可以實作，但要標 |
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
| 進言「請求君主出陣」（`sub_1699E`）| [`11-ai-sortie.md`](11-ai-sortie.md) | 可實作，**尚未實作** |
| 主畫面的視窗外框與指令列 | [`12-strategy-chrome.md`](12-strategy-chrome.md) | 版面常數全部改用機器碼值；各視窗**內部**排版仍有估值 |
| 主畫面四個視窗的開關 | [`13-main-window-toggles.md`](13-main-window-toggles.md) | 已實作並留下截圖；原版執行期未驗 |
| 編成時預備兵怎麼分配（`sub_14698`）| [`21-corps-formation-reserves.md`](21-corps-formation-reserves.md) | 已實作並有逐項單測 |
| 財政視窗 | [`14-finance-window.md`](14-finance-window.md) | 版面已照原版重寫；數值輸入器未接 |
| 軍團編成視窗 | [`22-corps-formation-window.md`](22-corps-formation-window.md) | 版面已照原版重寫；頭像與滑鼠未接 |
| 據點情報視窗 | [`23-city-info-window.md`](23-city-info-window.md) | 版面已照原版實作 |
| 軍團情報視窗 | [`24-corps-info-window.md`](24-corps-info-window.md) | 版面已照原版實作；指令流程未接 |
| 四槽選擇視窗 | [`25-slot-select-window.md`](25-slot-select-window.md) | 讀取／儲存已照原版；新遊戲未共用 |
| ＹＥＳ／ＮＯ 對話框 | [`26-yes-no-dialog.md`](26-yes-no-dialog.md) | 版面與命中算式已照原版 |
| 君主選擇視窗 | [`27-lord-select-window.md`](27-lord-select-window.md) | 版面已照原版；「自定」未接 |
| 劇本 JSON | [`28-scenario-json.md`](28-scenario-json.md) | 匯出／匯入／round-trip 已可用（`cmd/wlscen`）|
| 音樂與音效 | [`29-audio.md`](29-audio.md) | **DRAFT**：錄音鏈路已驗，逐曲觸發缺 RE，播放層未做 |
| remake 原生存檔格式 | [`20-save-format.md`](20-save-format.md) | 已接進遊戲並驗過；**只差放回 DOSBox 實測** |
| 戰術側欄的內容組成 | [`31-tactical-sidebar.md`](31-tactical-sidebar.md) | 七格已照原版實作；`▶▶` 列只畫美術不接行為 |
| 攻城的「門強度」條 | [`32-gate-strength-bar.md`](32-gate-strength-bar.md) | 已實作並有單測；右鍵提前收掉未接 |
| 底列六格是選部隊 | [`33-squad-selection.md`](33-squad-selection.md) | 已實作並有單測；命令圖示的來源段未定案 |
| 一覽表的欄位與版面 | [`38-list-windows.md`](38-list-windows.md) | 四個家族照原版重做；捲軸未解 |
| 戰術畫面的玩家操作 | [`37-tactical-player-controls.md`](37-tactical-player-controls.md) | 陣形選單與陣形線已接；說明書 4.2–4.6 的功能逐條對照過 |
| 兩個平面的地面圖與登城 | [`36-ground-planes-and-climbing.md`](36-ground-planes-and-climbing.md) | 已實作並有單測；攻城仍打不下來，卡點換成攻方不前進 |
| 縮小地圖的據點標記 | [`35-strategy-minimap.md`](35-strategy-minimap.md) | 已實作並有單測；22 勢力的選擇視窗用「點一下換下一個」代替 |
| 兩個速度設定的五檔 | [`34-speed-steps.md`](34-speed-steps.md) | 已實作並有單測；「最高速」的上限是 remake 差異 |
| 行軍指示的三選一 | [`39-march-order-menu.md`](39-march-order-menu.md) | 已接進畫面並有單測 |
| AI 軍團的決策鏈 | [`40-ai-march-decision.md`](40-ai-march-decision.md) | 已實作；「逐站前進」未移植 |
| 訊息框的版面 | [`41-message-box-geometry.md`](41-message-box-geometry.md) | 機器碼與影格兩條證據都對上，已實作 |
| 事件場景上誰在說話 | [`42-event-scene-speakers.md`](42-event-scene-speakers.md) | 兩個框已實作；結果階段的上框未解 |
| 回不了家就敗走 | [`43-rout-on-blocked-return.md`](43-rout-on-blocked-return.md) | 已實作並有單測 |
| 進言的原文 | [`44-advise-original-text.md`](44-advise-original-text.md) | 三個指令 × 64 則全部改查 `TALK.DAT` |
| 進言的畫面 | [`45-advise-scene-layout.md`](45-advise-scene-layout.md) | 插圖 ＋ 兩個框 ＋ 五列選單已實作並留下截圖 |
| 戰後退一站回家 | [`46-post-battle-retreat.md`](46-post-battle-retreat.md) | 已實作並有單測 |
| 據點易主之後守軍調頭 | [`47-city-fall-corps-redirect.md`](47-city-fall-corps-redirect.md) | 已實作並有單測 |
| 據點被攻陷時內政官被遣回 | [`48-governor-returns-on-city-fall.md`](48-governor-returns-on-city-fall.md) | 已實作並有單測，兩則訊息也接了 |
| 進言第四、五項（遷都／請求出陣）| [`49-advise-relocate-and-sortie.md`](49-advise-relocate-and-sortie.md) | 已接進畫面並留下截圖 |
| 軍費直接扣資金 | [`50-corps-upkeep-charges-funds.md`](50-corps-upkeep-charges-funds.md) | 已改正並有單測；長跑數字跟著變 |
| 結局：存活勢力數歸一 | [`30-victory.md`](30-victory.md) | 判定已實作並驗過；**結局的過場與訊息未做** |
| 顏色到不了滿刻度（6 bit DAC）| [`51-vga-dac-palette-scale.md`](51-vga-dac-palette-scale.md) | **READY**，尚未全面套用 |
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
| 結局的播放 | [`67-ending-playback.md`](67-ending-playback.md) | 已實作並有單測；十二幕 ＋ 逐字文字 ＋ 十七階淡入淡出 |
| 同狀態畫面對拍（方法）| [`90-same-state-parity.md`](90-same-state-parity.md) | 主畫面五區逐像素相同；§4.1 記下「參考影格本身會有東西」 |
| 戰場的逐區對拍（分區）| [`91-tactical-parity.md`](91-tactical-parity.md) | 九區裡六區逐像素相同，戰場區 0.17% |

## 怎麼加一份

1. 先確認 `docs/re/` 有出處。**沒有 RE 就沒有規格**——不要從 remake 的
   現況反寫規格，那只會把既有的偏差固定下來。
2. 複製 `TEMPLATE.md`，編號沿用 `docs/mechanics/` 的分類（10 政略、
   20 軍事、30 戰鬥、40 經濟、50 外交、60 人事、70 AI、80 勝負）。
3. 「remake 實作」欄留白就是缺口，**不要為了讓表好看而填近似值**——
   近似要寫成「近似」並說明差在哪。
