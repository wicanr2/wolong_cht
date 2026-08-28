# 推廣片產出紀錄

**狀態：五支影片都已產出並驗過媒體規格。主預告是 72 秒，
含語言切換與原版並排兩段；遊戲段落是逐幀錄下來的實跑畫面。
2026-08-28 因為戰術規則層動過，主預告、原版實機對照片與手機片三支全部重錄重剪。**

- 日期：2026-08-28（三支重錄重剪）／2026-08-26（主預告改成 72 秒的分鏡）／2026-08-21（其餘）

可交付推廣片已集中在 [`dist-all/promo`](../../dist-all/promo)。主預告輸出為
`dist-all/promo/wolong-remake-trailer.mp4`，長度 72 秒、1280×720、H.264／AAC。

⚠ **主預告有兩處原版衍生物，都是刻意的、也都標明了**：50–60 秒的並排段
左半是原版實機（`tools/dosv_capture.sh` 的受控擷取），配樂自 2026-08-26 起
是原版的曲子（由本專案的 OPL3 從使用者自備的 `BGM.DAT` 算出來，不是取樣
原版錄音）。細節見下面的「音訊與權利標記」。**原版的 `BGM.DAT`／`SOUND.DAT`
／執行檔本身沒有進發行素材**，其餘段落都是 remake 自己畫的畫面。

## 主預告的分鏡（2026-08-26 重剪）

| 秒 | 段落 | 來源 |
|---:|---|---|
| 0–4 | 標題卡 | 合成 |
| 4–14 | 大地圖 | 逐幀實跑（330 張，51 張不重複）|
| 14–18.5 | 事件訊息／財政 | 截圖（`promo-*.png`，**每次重剪都重拍**）|
| 18.5–30.5 | 野戰 | 逐幀實跑（390 張，84 張不重複）|
| 30.5–38.5 | 攻城 | 逐幀實跑（270 張，76 張不重複）|
| 38.5–41 | 戰果 | 截圖 |
| **41–50** | **語言切換**：繁中／简体／日本語／English | 四次啟動各 130 張 |
| **50–60** | **原版並排**：松崗 DOS/V 實機 ｜ remake | 原版是 `tools/dosv_capture.sh` 的受控擷取 |
| 60–65 | 一覽表／進言 | 截圖（`promo-*.png`，**每次重剪都重拍**）|
| 65–72 | 結尾卡 | 合成 |

輸出（2026-08-28 重錄重剪）：72.000 秒、1280×720、h264／aac 44.1 kHz 立體聲，
SHA-256 `131b2b0ee0b72ec9208b99148b44baa28807b6906da0c48456472a7e66f3ce43`。

⚠ **「不重複張數」會隨機器負載變**：同一條命令在閒置的機器上與在
14 核 load 90 的機器上量到的數字差三成（大地圖 71 → 51、野戰 119 → 84）。
它是「有沒有在動」的下限，**不是可比較的指標**——拿兩次的數字比大小沒有意義。

⚠ **程式改了就要重錄動態段**，不能只重剪。
踩過兩次：第一版的動態段錄在修「小兵只畫一半」與「月結框不消失」之前
（[`../spec/88`](../spec/88-display-polish-parity.md)），片子裡看得到那兩個問題；
2026-08-28 的戰術規則層修正（[`../release/10`](../release/10-full-20260828.md) §1）
同樣改變了攻城與野戰兩段的畫面。

⭐ **錄製與重拍的命令現在都在腳本裡**：動態段七 ＋ 三段在
`tools/promo_live_capture.sh`，五張靜態圖在 `tools/promo_stills.sh`。
先前這兩批的旗標沒有留下來，重拍要先重新猜一次——而猜錯的症狀是
「圖看起來很像」。

⚠ **靜態圖一律用 `promo-*.png`，不要借 playtest 的證據圖。**
那些圖的 SHA-256 記在 `docs/re/13`／`docs/playtest/12` 當 fixture 身分，
重拍會把紀錄弄壞；而不重拍，片子裡就會同時出現三個世代的 UI——
2026-08-26 踩過：片中同時有 2026-08-10 的純藍對話框、更早的舊版 HUD，
與當時的藍底龍紋，看起來像「同一款遊戲的對話框有三種樣子」
（[`../spec/88`](../spec/88-display-polish-parity.md) §1.1）。

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
  sh /src/tools/promo_video.sh /src/dist/promo/wolong-remake-trailer.mp4
```

⚠ **輸出要寫到 `dist/promo/`，不是 `dist-all/promo/`。**
`release_all_fs.py` 的 `promo_source()` **先找 `dist/promo/`**，找到才輪到
`dist-all/promo/`——只寫後者的話，下一次打包會用 `dist/promo/` 裡的舊片
把新片蓋掉，而 manifest、`SHA256SUMS.txt` 與檔名全部正常
（2026-08-26 踩過，[`../release/08`](../release/08-full-20260826.md) §4）。
產完之後複製一份到 `dist-all/promo/` 再跑
`WOLONG_BUNDLE_DATA=1 tools/py.sh tools/release_all_fs.py refresh` 重算雜湊。

⚠ **容器裡沒有 python3**，所以配樂要先在 `tools/py.sh` 那邊產好再用
`PROMO_SCORE` 指進來。

`PROMO_SCORE` 沒給的時候腳本才會退回 `tools/promo_score.py` 的原創合成配樂；
**現行主預告用的是原版曲子的 72 秒版**（下一節）。段落長度加起來要剛好等於
配樂長度，否則 ffmpeg 會在最後一段截斷或留白。

## 音訊與權利標記

主預告的配樂 2026-08-26 起**改用原版的曲子**（使用者裁定）。

- 來源：`tools/bgm2ogg.sh` ＋ `cmd/wlaudio` 的 **OPL3 合成**，從使用者自備的
  `BGM.DAT` 算出來的（[`../spec/29`](../spec/29-audio.md)）——**不是從原版
  錄音取樣**，也不是把 `BGM.DAT` 放進影片。
- 曲目照場景挑（[`../re/58`](../re/58-bgm-scene-mapping.md)）：
  大地圖段用曲 2（春）、戰場段用曲 9（平原野戰）。剪點對齊影片段落。
- ⚠ **它演奏的是原版的曲子**，所以這條音軌是原版衍生物：
  只用在推廣片，**不進任何遊戲包**，也不宣稱是本專案的創作。
- 重跑：`tools/promo_score_original.sh workplace/promo/score/original-72.wav`。

前一版的原創合成配樂（動機 D4–C4–F4–E4–A4，`tools/promo_score.py`）還在，
給「不想用原版曲子的場合」備著。

影片中的畫面全部是 remake 自己畫的；原版實機錄影只出現在標明的並排段。

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
72 秒，SHA-256 `1128867fe16cebaa2b1ae4e87371673b240f0643395bb6bc4ecbed4a5fcacb6d`。
**原版側是自己跑的受控 DOSBox-X 實機遊玩**：開新遊戲、劇本與君主選擇、
大地圖與時鐘、軍團編成、事件訊息、行軍指示，全部照 timeline 可以重跑。
只有戰術戰場那一格取自使用者提供的錄影，並在片上標明。

⭐ **2026-08-26 改成幾乎全片並排**（使用者建議）：原本開新遊戲、軍團編成、
行軍指示、戰術戰場四段只有原版單邊，現在都配上 remake 的同一個流程——
八段裡六段是並排。remake 側的三段新素材（`newgame`／`form`／`march`）
由 `tools/wlgame_frames.sh` 錄，與其他動態段同一條路。

⚠ **並排的來源尺寸不一定是 640×480**：受控擷取是 640×480 上下各 40 px 黑邊，
使用者提供的錄影是 956×720。所以 `split_clip` 一律先 `scale=640:480` 再裁——
少了那一步，956 寬的來源會被裁掉右邊一大半，而畫面看起來只是「構圖怪」
不像出錯。

分鏡、來源界線、SHA-256、抽樣驗收與重跑命令見
[dosv-realmachine.md](dosv-realmachine.md)。

前一支 [`dosv-live-comparison.md`](dosv-live-comparison.md)（2026-08-12，60 秒）
的原版畫面九成來自 YouTube 錄影，已被取代，不再放在 `dist-all/promo/`。

兩支都是推廣比較媒體，不會被包入任何目標平台的遊戲封裝。
