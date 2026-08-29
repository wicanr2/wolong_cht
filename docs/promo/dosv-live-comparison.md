# DOS/V 原版／remake 實機動態比較推廣片

**狀態：已產出並驗過，已被 `dosv-realmachine` 取代。**
當時的原版畫面九成來自使用者提供的 YouTube 錄影；新的那一支原版側是自己跑的實機。
這是同類型畫面／流程的推廣比較，不是同日期、同輸入、同狀態的逐像素 parity 證明。

- 日期：2026-08-12
- 成品：`wolong-remake-dosv-live-comparison.mp4`（**檔案已刪**，2026-08-29 清舊產物時連 `dist/promo/` 的那份一起；換成 `dosv-realmachine`。要重播照 §離線重播 再產一次）
- 長度／格式：60.000 秒、1280×720、30 fps、H.264／AAC（44.1 kHz、雙聲道）
- 成品 SHA-256：`a00a221bf3c4213f4d9777c66a3e58a06dcae93aab2beecc146de67a7447973d`（**當時的值**，檔案已刪，沒有工具在驗它）
- 合成腳本：[`tools/promo_dosv_live_comparison.sh`](../../tools/promo_dosv_live_comparison.sh)
- 後續：[`dosv-realmachine.md`](dosv-realmachine.md)（2026-08-23，原版側改成自己的實機擷取）

## 影片主軸與素材界線

這支影片以「經典再現」為主軸，畫面主體是實際遊玩／實際 GUI 擷取，而不是把代表幀
拼成假裝遊玩的影片。

- 原版動態片段來自使用者指定的松崗 DOS/V 遊玩錄影
  `https://www.youtube.com/watch?v=af6xqcicXoI`。本輪重新取得的檔案為 956×720、
  567.281 秒；暫存檔 SHA-256 為
  `e516632af2c4e7676bcfcdf33042ae1ccbe2df3cc74c1ec2083671851a50af2e`。
- 原版開場的「NEW GAME」畫面則由專案的受控 DOSBox-X 實機擷取取得，條件為
  `machine=vgaonly`、`core=normal`、`cputype=486`、`cycles=fixed 20000`，並由正常啟動
  路徑到達。
- remake 策略、編成、目的地與行軍畫面由 `wlgame` 的正常鍵盤路徑錄取；固定
  `seed=17`，沒有傳入 `-open-*` 旗標。正常段落的輸入順序是編成、選目的地、下達行軍。
- remake 戰術段是明確標示的獨立可視化 fixture（`-open-siege`），只用於展示戰術 GUI，
  不當成正常自然路徑的驗收證據。
- 原版畫面段落在個別編碼時以 `-an` 去除同步音訊，最後再鋪入同一使用者錄影
  `300–375 秒`的原版遊戲 AdLib 音軌；只做音量、限制器與淡入淡出，沒有呼叫
  `tools/promo_score.py` 或重作音色。完整來源與權利界線見
  [`dosv-adlib-and-tactical-review.md`](dosv-adlib-and-tactical-review.md)。
- 成品是使用者指定的宣傳比較媒體，不隨 Windows、macOS、Linux 或 Android 遊戲包發行，
  且不包含原版可執行檔、資料檔、字型或原始錄影檔。

## 60 秒分鏡

| 時間 | 畫面 | 實機來源與標示 |
|---|---|---|
| 00–04 秒 | 開場卡 | 說明是松崗 DOS/V 原版實機 × remake 實機。 |
| 04–09 秒 | 原版新遊戲 | 受控 DOSBox-X 正常啟動後的原版畫面。 |
| 09–16 秒 | 自然策略並排 | 原版錄影 76–83 秒與 remake 正常策略 GUI；片上標示「非同狀態逐像素判定」。 |
| 16–22 秒 | 原版策略 | 原版錄影 154–160 秒的策略地圖與時鐘。 |
| 22–28 秒 | remake 編成／行軍 | 正常鍵盤路徑的實機擷取。 |
| 28–33 秒 | 指令與事件並排 | 原版錄影 236–241 秒的事件訊息與 remake 目的地選擇；明示不是同一局面。 |
| 33–40 秒 | 原版戰術 | 原版錄影 317–324 秒；先正規化／裁成 640×400。 |
| 40–45 秒 | remake 戰術 | 新 viewport 的攻城 fixture；片上明示不同戰況、只作 layout-only 比較。 |
| 45–50 秒 | 原版戰鬥／結果 | 原版錄影 398–403 秒。 |
| 50–55 秒 | remake 行軍、時鐘與自然畫面 | 正常鍵盤路徑下達命令後的實機擷取。 |
| 55–60 秒 | 結尾卡 | 明示配樂來自使用者錄影中的原版 AdLib。 |

## 驗收結果

1. 影片十個時間點（2、6、12、20、27、34、42、48、54、59 秒）已抽成接觸表檢視；
   原版／remake 標籤、字幕、戰術界線與結尾卡皆完整可讀，沒有裁切。
2. remake 原始擷取序列以 SHA-256 去重：策略 idle 20 幀有 15 種內容、編成後通知 15 幀有
   10 種、目的地選擇 15 幀有 11 種、行軍 25 幀有 22 種、戰術 idle 15 幀有 15 種、
   戰術操作後 20 幀有 18 種。這確認影片取自持續渲染的實機畫面，而非重複一張靜態圖。
3. `ffprobe` 已確認成品只有一條 H.264 視訊與一條 AAC 音軌，長度正好 60.000 秒、視訊 1800 幀。
4. 42 秒代表幀曾發現 overlay 字型缺 CJK 而顯示方框；正式成品已改用
   Noto CJK 重建，中文標題與字幕可讀。
5. 原版畫面先正規化為 640×480，再裁成 640×400；remake 原始擷取也是 640×400，
   兩側使用同一條最近鄰縮放鏈。戰術離屏 viewport 已由 496×384 修為 480×368。
6. 最終 AAC 音軌平均 `-21.3 dB`、峰值 `-4.2 dB`，非靜音且無削波。

## 離線重播

原版原始錄影只放在 `.gitignore` 排除的 `workplace/`，不進版本庫或遊戲包。若要重播，
須由使用者提供該錄影，將其暫放在
`workplace/promo-live/original-video/original.mp4`，再使用下列完全離線的 Docker 指令：

```text
docker run --rm --network none --memory 4g --cpus 2 --pids-limit 256 \
  -u "$(id -u):$(id -g)" \
  -v "$PWD":/src:ro \
  -v "$PWD/workplace/promo-live":/capture:ro \
  -v "$PWD/dist/promo":/out \
  -e PROMO_FONTFILE=/usr/share/fonts/opentype/noto/NotoSerifCJK-Bold.ttc \
  -e WOLONG_PROMO_TACTICAL=/capture/remake-tactical-new \
  u5cht/video:latest \
  bash /src/tools/promo_dosv_live_comparison.sh /out/wolong-remake-dosv-live-comparison.mp4
```

容器不開網路、原版來源唯讀掛載、只允許寫入 `dist/promo`。完成後應重新以 `ffprobe`、
音量與 SHA-256 驗收。
