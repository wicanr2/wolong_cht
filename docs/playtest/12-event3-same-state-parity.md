# event3 同狀態對拍（2026-08-10）

本輪只做短路徑，不跑完整劇本。目的為固定事件 3 的 raw fixture、原版 PC-98 三選一
畫面，以及 remake 從同一 fixture 進入選單後的 640×400 截圖。

**狀態：同狀態事件接線、原版式 composite、3×6 實際格位選取、TALK 五行分頁、
DOS/V 內框／按鍵 glyph／硬體游標資產與 pending 結束後消像已通過短 smoke；自然
DOS/V／remake 整張畫面的逐像素對拍仍未完成。**

- 日期：2026-08-10

## Fixture

- PC-98 `KI.EXE` SHA-256：`061917F9F3F5C03E29397A9C636D546052128A99B8C8CE31DED0E84CF2A481E8`
- fixture 事件字：`0x0303`（來源勢力 3、事件 3、Param 低 byte 為玩家 0）
- 來源勢力外交官：武將 17
- 原版畫面：[`pc98-oracle-event3-choice.png`](../images/pc98-oracle-event3-choice.png)
- remake composite 選項畫面：[`wlgame-event3-choice.png`](../images/wlgame-event3-choice.png)
  （SHA-256：`CA40B865B44A6EA13ED5B4F2C0B6AB913A0BC895EF48D7A19E1825501E535151`）
- remake 數值器畫面：[`wlgame-event3-amount.png`](../images/wlgame-event3-amount.png)
  （SHA-256：`27A5474EBA79C92C23B24A79938CA4E1D376B9FA52C0956AE3D3359C0404609D`）

## 重播摘要

原版以 PC-98 DOSBox-X harness 啟動後，先進入前置外交通知，再進入三選一。remake
以相同 raw save overlay 啟動，等待事件通知後送出一次 `Enter`，在第 600 個畫面更新
截圖。兩邊日期均為 196 年 4 月 1 日；remake 的截圖另外固定亂數種子 1。

原版與 remake 可直接比對地圖、日期、事件語意、肖像／IVENTGRF 位置與選項流程；兩張
remake 圖另外驗證：滑鼠命中 `(88,200)` 的第一個 raw 格位會更新數值器，完成 pending
後下一幀不再繪製事件場景與肖像，後續 TALK 才以五行一頁顯示。`KI.EXE` 內建 cursor
與 `ICONGRF` 96×64 下半部 glyph 已由 DOS/V 原始 bytes 解出並接線；這仍不是整張
原版／remake 截圖相等測試。

## 本輪結果

- 事件 3 raw dispatch：通過。
- 前置 TALK → 三選一 composite：通過。
- 3×6 數值編輯：通過 raw table、`AmountEdit` 單測與 X11 格位點擊 smoke；鍵盤 fallback
  與滑鼠共用同一狀態動作。
- TALK：通過原始 hard line、結構尾空行、中間空行與五行／16 px 分頁測試。
- 消像：通過 pending 結束後地圖重畫；事件場景／肖像不殘留。
- DOS/V cursor／button glyph 資產來源與解碼：通過；自然整張畫面对拍、其他事件物件動畫
  的逐像素 parity：未完成。
