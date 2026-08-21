# 10 — 事件 TALK 通知 modal

**狀態：已完成通知資料接縫、Linux/Xvfb 視覺抽樣與 remake TALK 五行分頁；原版未知
事件流程／未定位 formatter 分支仍未完成。**

- 日期：2026-08-09
- 執行環境：Docker／Xvfb，`wolong-go:20260809`
- 原始輸入：DOS/V `SINARIO.DAT`、`TALK.DAT` 與倚天字型，唯讀掛載

本輪把已由反組譯確認的事件 11／12／13 通知接到 `wlgame`：玩家城市災害使用 TALK
#70／#71／#72，事件 13 使用 TALK #51。事件 9 的釋放武將也共用同一個通知佇列，使用
TALK #37（程式索引 `0x25`）。

## 1. 可重現驗收

在 Docker／Xvfb 中，以自備且唯讀掛載的 DOS/V 原版素材與倚天字型執行：

```text
wlgame -orig workplace/orig/dosv -font workplace/eten \
  -scenario 0 -player 0 -seed 17 -open-message \
  -shot /tmp/wlgame-event-modal.png -shot-frames 30
```

`-open-message` 是驗收旗標：它只建立已證實的玩家首都 TALK #70，方便固定抽樣，
不宣稱模擬未知的事件 producer。2026-08-09 的輸出是
[`wlgame-event-modal.png`](../images/wlgame-event-modal.png)：畫面顯示「許昌發生了暴風雨。」
與 Enter／Space 繼續提示；第 30 幀仍為 196 年 4 月 1 日，證明 modal 期間世界時間停止。

## 2. 實作邊界

- `Event.TalkNotices` 只保存原始 TALK index 與 city/general ID；呈現層才解 Big5、保留
  TALK.DAT 硬斷行並代換 ASCII marker `\\1`／`\\2`。
- 缺少 marker 時整則訊息不顯示，避免把半句或錯誤索引當成原版文字。
- 視窗位置遵守 640×400 原生畫布的戰略／命令區層次；一般通知已按原始 hard line、
  實際 ASCII／CJK 字寬與五行／16 px 分頁，TALK composite 另使用原版肖像／場景位置。
  這仍不等於未定位事件的原版肖像、逐頁動畫、音效或完整 formatter parity。

## 3. 證據

直接證據為 `KI.EXE.asm` 的 IDA 線性位址 `0001237E`、`000134A6`、`000134B1`、
`00013507`，以及唯讀 `TALK.DAT`；完整輸入雜湊、工具版本與推論等級記在
[`RESEARCH-LOG.md`](../../RESEARCH-LOG.md)。
事件 10 的 raw consumer、事件 6／7 次要 TALK raw 條件／索引與災害物件 timer
已在 2026-08-10 接入；事件 10 的受控 raw producer 與事件 6／7 的 raw
formatter word 邊界也已補上。

### 3.1 未完成項

- 事件 10 producer 仍未定位。**只有負證據**，不能把受控注入口寫成原版劇本來源。
- 事件 6 #72 的缺失 formatter payload 維持 fail-closed。
- 原版／remake 同狀態畫面對拍仍是剩餘驗收項。
