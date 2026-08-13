# 推廣片產出紀錄

**狀態：remake 推廣片、代表幀「經典再現」比較片與 DOS/V／remake 實機動態比較片均已產出並完成媒體規格驗證；
三平台候選可執行包與 Linux AppImage 已產出，Windows／macOS 原生 GUI 仍是獨立 release gate。**

- 日期：2026-08-12

可交付推廣片已集中在 [`dist-all/promo`](../../dist-all/promo)。主預告輸出為 `dist/promo/wolong-remake-trailer.mp4`，長度 60 秒、1280×720、
H.264／AAC。內容只使用 remake 的 DOS/V 目標畫面截圖、事件／戰術驗收圖與
本專案原創合成聲；沒有把原版影片、`BGM.DAT`、`SOUND.DAT` 或原版執行檔放進
發行素材。

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

```text
sh tools/promo_video.sh dist/promo/wolong-remake-trailer.mp4
```

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
  tools/promo_classic_revival.sh dist/promo/wolong-remake-classic-revival.mp4
```

## DOS/V／remake 實機動態比較片

使用者要求「要有遊戲實際遊玩畫面」，因此新增
[`wolong-remake-dosv-live-comparison.mp4`](../../dist-all/promo/wolong-remake-dosv-live-comparison.mp4)。
這支 60 秒影片以松崗 DOS/V 原版動態遊玩錄影、受控 DOSBox-X 原版新遊戲畫面，以及
remake 實機 GUI 擷取組成；內容涵蓋自然策略、編成／行軍、事件訊息與戰術畫面。

完整來源界線、秒數、正常鍵盤路徑／獨立戰術 fixture 的區別、原版音訊排除、SHA-256、
抽樣驗收與離線重播命令見 [dosv-live-comparison.md](dosv-live-comparison.md)。它是使用者
指定的推廣比較媒體，不會被包入任何目標平台的遊戲封裝。
