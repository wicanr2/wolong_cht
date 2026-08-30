# 原版實機遊玩 × remake 實機的對照推廣片

**狀態：已產出並驗過。原版側是自己跑的受控 DOSBox-X 實機遊玩，
只有戰術戰場那一格仍取自使用者提供的錄影。**

- 日期：2026-08-23
- 成品：`dist/promo/wolong-remake-dosv-realmachine.mp4`
- 長度／格式：72.000 秒、1280×720、30 fps、H.264 ＋ AAC（44.1 kHz 立體聲）
- 成品 SHA-256：`522313d127846d57a14ed732960fcc81635cac68c47c40f170da08091346de02`（2026-08-29 重錄重剪；原版側的 DOSBox-X 擷取沒動，重錄的是 remake 側）
- 合成腳本：[`../../tools/promo_dosv_realmachine.sh`](../../tools/promo_dosv_realmachine.sh)
- 取代：[`dosv-live-comparison.md`](dosv-live-comparison.md)（2026-08-12，原版畫面九成來自 YouTube 錄影）

## 1. 這一支跟前一支差在哪

前一支的原版畫面**九成來自使用者提供的 YouTube 錄影**，只有「NEW GAME」
那一格是自己的實機擷取。這一支反過來：**原版側是自己跑的**，
可以照 timeline 重跑，也不夾帶第三方素材。

| 段 | 原版側來源 |
|---|---|
| 開新遊戲、劇本與君主選擇 | 自己的 DOSBox-X 實機 |
| 大地圖與時鐘 | 自己的 DOSBox-X 實機 |
| 軍團編成 | 自己的 DOSBox-X 實機 |
| 指令與事件訊息 | 自己的 DOSBox-X 實機 |
| 行軍指示 | 自己的 DOSBox-X 實機 |
| **戰術戰場** | ⚠ **使用者提供的錄影**（見 §4），片上有標明 |

## 2. 分鏡

| 時間 | 畫面 | 來源 |
|---|---|---|
| 00–04 | 開場卡 | — |
| 04–11 | 原版實機：開新遊戲、劇本與君主選擇 | `promo-dosv-main/o1-newgame.mp4` |
| 11–20 | **並排**：大地圖 | 原版 `o2-strategy.mp4` ／ remake `live/map` |
| 20–28 | 原版實機：軍團編成 | `o3-corps.mp4` |
| 28–36 | **並排**：指令與事件 | 原版 `o5-battle.mp4` ／ remake `live/map` |
| 36–44 | 原版實機：行軍指示 | `o4-march.mp4` |
| 44–52 | 原版：戰術戰場 | 使用者錄影 317–325 秒，**片上標明** |
| 52–61 | remake 實機：野戰 | `live/battle` |
| 61–67 | remake 實機：攻城 | `live/siege` |
| 67–72 | 結尾卡 | — |

配樂是 DOSBox-X mixer 實錄的原版 AdLib（`original-audio/original-adlib.wav`，
2026-08-12 那次擷取）。只做音量、限制器與淡入淡出。

## 3. ⭐ `record:` 錄不出即時制的真實速度

原本的 `record:秒,fps,前綴` 是**逐幀 `import`**，一張要 280–330 ms。
量出來：要 8 秒 10 fps，拿到 80 張圖，但它們**橫跨 25.67 秒**——
實際只有 **3.08 fps**。

照標稱 fps 編出來，原版會播快 3.2 倍。這一款是即時制、畫面上有時鐘
（`CLAUDE.md` §3.1），對照片那樣做等於**謊報原版的速度**。

所以加了 `grab-start:fps,名稱` ／ `grab-stop`，用 `ffmpeg -f x11grab`
錄同一個 Xvfb：真實時間、固定幀率。實測 15.000 秒、425 幀。

兩件事讓它必須是 start／stop 兩步而不是「錄 N 秒」：

- 同步錄的話 timeline 會卡在 ffmpeg 上，**錄得到閒置畫面，錄不到操作**
  ——而推廣片要的正是「有人在玩」。
- 停止一定要 **SIGINT 不能 SIGKILL**：mp4 的 moov atom 在收尾時才寫，
  硬砍的話檔案存在、大小正常，**但播不動**。

⚠ **DOSBox-X 自己的 `ctrl+alt+F5` 沒有用。** 試過，log 只在啟動時印一行
`USING AVI+ZMBV`（那是編碼器宣告，不是開始錄），輸出目錄一個檔都沒有。
`record:` 保留給「幾張可辨識的代表畫面」。

## 4. ⚠ 原版的戰術戰場沒擷取到

原版要開出戰場的方法是「編成一支軍團，然後等 AI 來打」
（[`../playtest/40`](../playtest/40-tactical-parity.md) §1）。
**四次專門擷取都沒等到**：

| 次 | 做法 | 結果 |
|---|---|---|
| 1 | 編成 ＋ 行軍，等 150 秒 | 到 5月28日，沒開戰 |
| 2 | 只編成，等 270 秒 | 6月5日出現「張遼的兵馬，向許昌進攻過來了！！」，但訊息框要點掉才進戰場，錄影已停 |
| 3 | 編成，等 280 秒後點訊息框 | 6月2日還沒到那則訊息 |
| 4 | 照 `probe-march` 的成功配方（編成 ＋ 行軍 ＋ 純等 250 秒）| 到 6月18日，沒開戰 |

成因不是腳本錯：**原版以時鐘播種**，每次跑的 RNG 都不同
（`CLAUDE.md` §3.1「原版沒有固定 tick rate」）。同一組操作在
`probe-march` 那次第 45 天開戰，在這四次都沒有。

⭐ **這一格因此用錄影補，並在片上標明**「此段取自使用者提供的松崗 DOS/V
遊玩錄影，非本次實機擷取」。使用者 2026-08-23 裁定「以實機為主，錄影補不足」。

## 5. 驗收

| 項目 | 結果 |
|---|---|
| 媒體規格 | ✅ 72.000 秒、2160 幀、1280×720、30 fps、單一 H.264 ＋ 單一 AAC |
| 音軌 | ✅ 平均 −21.9 dB、峰值 −4.5 dB，非靜音也沒削波 |
| 中文可讀 | ✅ 六個時間點抽幀檢視，標題、標籤與註記都完整（Noto Serif CJK）|
| 畫面沒被字條蓋住 | ✅ 遊戲畫面縮成 944×590 放在 y=65，上下字條不重疊 |
| 原版側真的在動 | ✅ 五段各取樣 7–51 張，相異率 7/8 – 40/51 |
| remake 側真的在動 | ✅ map 330 張中 38 張不重複、battle 390 中 119、siege 270 中 28 |

⚠ **這不是逐像素 parity 證明。** 兩側不是同一局面、同一輸入；
remake 的大地圖是**最高速檔**錄的（低速檔八秒幾乎不動），
原版是預設速度——**兩側的時鐘推進速度不可比**，片上已註明。
逐像素的判定在 [`../playtest/37`](../playtest/37-main-screen-parity.md) 與
[`../playtest/40`](../playtest/40-tactical-parity.md)。

## 6. 怎麼重跑

原版擷取（每次約 6 分鐘，四段）：

```sh
OPEN="click:320,215;wait:3;click:300,190;wait:4;click:450,154;wait:2;click:352,336;wait:3;click:352,336;wait:4"
CORPS="click:0,0;wait:2;click:37,0;press;wait:2;click:20,5;press;wait:2;click:200,154;press;wait:2;click:324,336;press;wait:2;rclick:200,154;wait:1;rclick:200,154;wait:2"
MARCH="click:25,5;press;wait:2;click:220,115;press;wait:2;click:200,134;press;wait:2;click:308,236;press;wait:2;click:312,250;press;wait:3"
tools/dosv_capture.sh promo-dosv-main \
  "wait:125;grab-start:30,o1-newgame;$OPEN;grab-stop;grab-start:30,o2-strategy;wait:20;grab-stop;grab-start:30,o3-corps;$CORPS;grab-stop;grab-start:30,o4-march;$MARCH;grab-stop;grab-start:30,o5-battle;wait:150;grab-stop"
```

remake 擷取：`tools/promo_live_capture.sh`（`docs/spec/71`）。

合成（容器不開網路、素材唯讀）：

```sh
docker run --rm --network none --memory 4g --cpus 2 --pids-limit 256 \
  -u "$(id -u):$(id -g)" \
  -v "$PWD":/src:ro -v "$PWD/workplace/promo-live":/capture:ro \
  -v "$PWD/workplace/promo/live":/remake:ro -v "$PWD/dist/promo":/out \
  u5cht/video:latest \
  bash /src/tools/promo_dosv_realmachine.sh /out/wolong-remake-dosv-realmachine.mp4
```

## 7. 未解

| 項目 | 現況 |
|---|---|
| 原版戰術戰場的實機擷取 | 四次未觸發（§4）。要嘛接受原版 RNG 的變異多跑幾次，要嘛從存檔直接進戰場——後者要先解「怎麼從存檔載入到開戰的那一刻」 |
| 原版 AdLib 的同場錄音 | `ctrl+F6` 的 WAV 擷取這次沒生效，配樂沿用 2026-08-12 那次的實錄 |
| 兩側時鐘速度可比 | remake 用最高速檔才看得到動靜；要真的可比，得先量原版預設檔的每日實時秒數 |
