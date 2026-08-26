# 推廣片產出紀錄

**狀態：五支影片都已產出並驗過媒體規格。主預告 2026-08-26 重剪為 72 秒，
加入語言切換與原版並排兩段；遊戲段落是逐幀錄下來的實跑畫面。**

- 日期：2026-08-26（主預告重剪）／2026-08-21（其餘）

可交付推廣片已集中在 [`dist-all/promo`](../../dist-all/promo)。主預告輸出為
`dist-all/promo/wolong-remake-trailer.mp4`，長度 72 秒、1280×720、H.264／AAC。
內容只使用 remake 自己的畫面與本專案原創合成聲；沒有把原版影片、`BGM.DAT`、
`SOUND.DAT` 或原版執行檔放進發行素材。

## 主預告的分鏡（2026-08-26 重剪）

| 秒 | 段落 | 來源 |
|---:|---|---|
| 0–4 | 標題卡 | 合成 |
| 4–14 | 大地圖 | 逐幀實跑（330 張，71 張不重複）|
| 14–18.5 | 事件 TALK ／ 撥款 | 截圖（畫面本來就不動）|
| 18.5–30.5 | 野戰 | 逐幀實跑（390 張，119 張不重複）|
| 30.5–38.5 | 攻城 | 逐幀實跑（270 張，56 張不重複）|
| 38.5–41 | 戰果 | 截圖 |
| **41–50** | **語言切換**：繁中／简体／日本語／English | 四次啟動各 130 張 |
| **50–60** | **原版並排**：松崗 DOS/V 實機 ｜ remake | 原版是 `tools/dosv_capture.sh` 的受控擷取 |
| 60–65 | M7 校訂／存檔 | 截圖 |
| 65–72 | 結尾卡 | 合成 |

輸出（2026-08-26）：72.000 秒、1280×720、h264／aac 44.1 kHz 立體聲，
SHA-256 `9b2e2c0ad2ad25499e008acd96b0e7d8e805d13beeef5dd2a2acccc136cb3e86`。

**段落長度加起來要等於配樂長度**（`tools/promo_score.py` 的 `DURATION`，
現在是 72）。對不上時音樂的轉折會落在畫面中間——那種不對齊說不出哪裡怪，
只會覺得片子鬆散。

### 語言段為什麼是四次啟動不是按 F9

F9 在遊戲中即時切換，切出來的畫面與 `-lang` 啟動**逐像素相同**
（[`../playtest/46`](../playtest/46-runtime-language-switch.md)）。
四次啟動可以逐幀重現，按鍵注入不行——同一顆種子、同一個時間點，
四段拍到的是同一個局面的四種語言。

### 並排段的界線

兩側**不是同一局面、不是同一輸入，也不是同一個時鐘速度**（remake 這一段
是最高速檔錄的）。片上已註明，逐像素的判定在
[`../playtest/37`](../playtest/37-main-screen-parity.md) 與
[`../playtest/40`](../playtest/40-tactical-parity.md)。

⚠ 原版擷取是 640×480、上下各 40 px 黑邊。不裁掉的話等比縮放後原版那一格
會小一圈，看起來像「remake 比較大」而不是「錄影帶黑邊」。

## 主預告的遊戲段落是實跑錄製

大地圖、野戰與攻城三段由 [`tools/promo_live_capture.sh`](../../tools/promo_live_capture.sh)
逐幀錄出來（機制見 [`../spec/71`](../spec/71-promo-live-capture.md)）：程式自己把每一張
畫出來的圖寫成 PNG，不做螢幕擷取，所以同一組旗標跑兩次得到同一批圖。
2026-08-26 那一批是大地圖 330 張、野戰 390 張、攻城 270 張，
不重複的張數分別是 71／119／56。

⚠ **戰術段的鏡頭與推進步數會隨規則改動失效。** 那兩個參數綁的是「這一局
打到哪裡」，規則層一動就要重掃——症狀是錄出一段空景，**不會有任何錯誤訊息**。
2026-08-26 重掃時攻城段就中了：舊設定（`-battle-steps 200 -battle-cam 20,15`）
錄到 270 張只有 12 張不一樣，而且那 12 張是循環的動畫不是戰況；
掃過 0／40／90／120／150／170 之後 120 步是尖峰，鏡頭改到 44,8（城內）。
`tools/promo_live_capture.sh` 有一道 20 張的下限擋這件事。

事件視窗、繁中校訂與存檔那幾段仍是截圖——**那些畫面本來就不動**，
錄成影片與截圖沒有差別。

## YouTube 原版／remake 畫面比較

依使用者要求，已把 YouTube 原版代表幀與推廣片所使用的 remake 畫面製成研究用
[對照片](../../dist-all/promo/wolong-remake-yt-comparison.mp4)，並保留[自然畫面並排圖](yt-remake-natural-side-by-side.png)
與[像素差異圖](yt-remake-natural-difference.png)。量測與解讀見
[yt-remake-pixel-review.md](yt-remake-pixel-review.md)。

這項比較確認的是可見像素差異、畫面骨架與 HUD 幾何，不冒充同日期／同輸入／同狀態
的逐像素 parity。這一支舊對照片的原版側仍只有代表幀；使用者後續要求的實機動態
對照另見下節，且不隨三平台遊戲包發行。

## 可重現產出

在掛載本儲存庫、倚天字型與 `/usr/share/fonts` 的 ffmpeg Docker 容器內執行：

```sh
tools/py.sh tools/promo_score.py workplace/promo/score/score-72.wav
docker run --rm --network none --memory 4g --cpus 2 --pids-limit 256 \
  -u "$(id -u):$(id -g)" -e HOME=/tmp \
  -e PROMO_FONTFILE=/usr/share/fonts/opentype/noto/NotoSerifCJK-Bold.ttc \
  -e PROMO_SCORE=workplace/promo/score/score-72.wav \
  -v "$PWD":/src u5cht/video:latest \
  sh /src/tools/promo_video.sh /src/dist-all/promo/wolong-remake-trailer.mp4
```

⚠ **容器裡沒有 python3**，所以配樂要先在 `tools/py.sh` 那邊產好再用
`PROMO_SCORE` 指進來。

腳本會先由 `tools/promo_score.py` 產生 60 秒 WAV，再把最新 DOS/V 自然策略骨架、
事件畫面、
事件 2–5 TALK、戰術戰鬥／投射物／戰果、事件 9、M7、存檔與結尾卡串成影片。

## 音訊與權利標記

- 動機：D4–C4–F4–E4–A4（MIDI 62–60–65–64–69），為本專案本輪新作，沒有從
  原版資料或參考影片取樣。
- 音色：程式合成三角波／方波與簡單鼓聲；44.1 kHz、立體聲、60 秒。
- 音訊用途：remake 推廣片預告配樂；不是原版音樂還原，也不應標記成原版 OST。
- 影片中的畫面標記為 remake 驗收截圖；YouTube 原版錄影僅作研究 oracle，未被
  混入影片。

影片完成後仍需在正式發行前保存 ffprobe 摘要與 SHA-256，並由發行者在目標平台
實際播放一次；這些是媒體封裝 gate，不是遊戲規則 parity 證據。

## 「經典再現」比較片

依使用者提供的 YouTube 遊玩影片與 `retro` deterministic DOSBox 技巧，新增
[`classic-revival.md`](classic-revival.md) 與
[`dist-all/promo/wolong-remake-classic-revival.mp4`](../../dist-all/promo/wolong-remake-classic-revival.mp4)。
影片的原版側使用代表幀，remake 側使用固定 `seed=17` 的驗收畫面；片中清楚標示
`core=normal`、`cputype=486`、`cycles=20000` 的重播基準與「非同狀態逐像素 parity」界線。
重現命令：

```text
PROMO_FONTFILE=/fonts/NotoSansTC-Regular.otf \
  tools/promo_classic_revival.sh dist-all/promo/wolong-remake-classic-revival.mp4
```

## DOS/V／remake 實機對照片

現行的是
[`wolong-remake-dosv-realmachine.mp4`](../../dist-all/promo/wolong-remake-dosv-realmachine.mp4)，
72 秒。**原版側是自己跑的受控 DOSBox-X 實機遊玩**：開新遊戲、劇本與君主選擇、
大地圖與時鐘、軍團編成、事件訊息、行軍指示，全部照 timeline 可以重跑。
只有戰術戰場那一格取自使用者提供的錄影，並在片上標明。

分鏡、來源界線、SHA-256、抽樣驗收與重跑命令見
[dosv-realmachine.md](dosv-realmachine.md)。

前一支 [`dosv-live-comparison.md`](dosv-live-comparison.md)（2026-08-12，60 秒）
的原版畫面九成來自 YouTube 錄影，已被取代，不再放在 `dist-all/promo/`。

兩支都是推廣比較媒體，不會被包入任何目標平台的遊戲封裝。
