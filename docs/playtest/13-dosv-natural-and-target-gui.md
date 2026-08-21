# DOS/V 自然畫面與目標平台 GUI 驗收

**狀態：歷史驗收紀錄（2026-08-11，影片 oracle）。**
⭐ **「嚴格同狀態逐像素差異」後來做到了**——主畫面五區、系統選單、
戰場九區裡六區都逐像素相同（[`37`](37-main-screen-parity.md)–[`40`](40-tactical-parity.md)）。
**Windows／macOS 原生 runtime 仍未在目標 OS 上跑過**，那一條至今成立。

- 日期：2026-08-11
- 目標：松崗 DOS/V；PC-98 不列入本輪畫面 oracle
- 範圍：短 smoke、自然啟動頁面紀錄、Linux GUI、Windows／macOS GUI 交叉建置、三平台候選封裝
- 限制：依使用者要求不跑完整長程遊戲測試

## 結果摘要

| 驗收項目 | 結果 | 證據／限制 |
|---|---|---|
| DOS/V 原版自然畫面 | 啟動路徑可取得 | 密碼頁可由空白確認進入原版開場；自然策略／戰術的既有影片仍是本輪視覺參考 |
| remake DOS/V 素材短 smoke | 通過 | [`wlgame-dosv-natural-remake-skeleton.png`](../images/wlgame-dosv-natural-remake-skeleton.png)，固定 `seed=17`、30 幀，右欄骨架已對齊 |
| Linux GUI | 通過 | Docker／Xvfb 啟動 `cmd/wlgame` 並產生截圖 |
| Windows amd64 GUI 建置 | 通過建置 gate | 輸出格式 `PE32+ executable x86-64`；沒有 Windows 原生桌面 runtime |
| macOS amd64 GUI 建置 | 通過建置 gate | 輸出格式 `Mach-O x86_64`；沒有 macOS 原生桌面 runtime |
| macOS arm64 GUI 建置 | 通過建置 gate | 輸出格式 `Mach-O arm64`；沒有 macOS 原生桌面 runtime |
| 三平台候選封裝 | 通過封裝 gate | `dist/release-20260811/packages/`；每包 SHA-256／deny-list／解包檢查通過 |
| 原版／remake 自然畫面視覺對拍 | 通過（影片參考） | 橫幅、命令列、27×21 地圖區、208 px 右欄與縮小地圖已對齊 |
| 原版／remake 嚴格同狀態逐像素 | 未宣稱 | 影片是 478×360、30 fps、壓縮／縮放且日期與鏡頭狀態不同 |

## 使用者提供的 YouTube 自然 oracle

來源：[臥龍傳 呂布開局滅曹操](https://www.youtube.com/watch?v=af6xqcicXoI)，由使用者
提供。影片 metadata：長度 567 秒、30 fps、478×360、上傳日期 2018-07-01；下載的
低畫質格式為 YouTube format `18`（約 28.64 MiB），只在 Docker 暫存，未把影片放入
儲存庫。

代表幀保留在 `docs/images/`，供日後不依賴 YouTube 即可回看：

| 時間 | 影像 |
|---:|---|
| 20 s | [`yt-wolong-natural-20s.png`](../images/yt-wolong-natural-20s.png) |
| 80 s | [`yt-wolong-natural-80s.png`](../images/yt-wolong-natural-80s.png)；去黑邊／還原 640×400 版 [`yt-wolong-natural-80s-640x400.png`](../images/yt-wolong-natural-80s-640x400.png) |
| 160 s | [`yt-wolong-natural-160s.png`](../images/yt-wolong-natural-160s.png) |
| 240 s | [`yt-wolong-natural-240s.png`](../images/yt-wolong-natural-240s.png) |
| 320 s | [`yt-wolong-natural-320s.png`](../images/yt-wolong-natural-320s.png) |
| 400 s | [`yt-wolong-natural-400s.png`](../images/yt-wolong-natural-400s.png) |
| 480 s | [`yt-wolong-natural-480s.png`](../images/yt-wolong-natural-480s.png) |
| 550 s | [`yt-wolong-natural-550s.png`](../images/yt-wolong-natural-550s.png) |

80 秒幀原始 SHA-256 為
`d33fff8d664e24321274310287dce38b82c82cfb62f3d0427e70dfd5bd301e08`；去黑邊／縮放後
640×400 參考幀 SHA-256 為
`c0217b8722bd44a22a112a2981b626126d5ee53d3e9f00498c6cbd018e08d6`。

影片實際可量得的共同 GUI 骨架是：橫幅 32 px、命令列 32 px、左側地圖 432×336 px
（27×21 個 16 px 格），右側 208 px；右欄上方是 192×128 縮小地圖，下方是君主／
軍師、信賴度、資金與三種預備兵。`cmd/wlgame/strategyhud.go` 現在以這組幾何畫出
常駐 DOS/V 自然策略 HUD，浮動命令／列表／事件視窗仍維持原有路徑。

這裡的「對拍通過」是影片 oracle 的畫面結構、資產、色彩層級與數值欄位置對照；影片
畫面在 80 秒為 196 年 4 月 5 日，remake smoke 為 196 年 4 月 1 日，且來源經過
478×360 縮放與有損壓縮，所以不把兩張圖的每個像素宣稱為同狀態 parity。

## DOS/V 原版短 smoke

在隔離 Docker／Xvfb／DOSBox 中使用松崗 DOS/V 素材，固定 `machine=vgaonly`、
`cputype=486`、`cycles=fixed 20000` 啟動 `START`。結果是可重現的密碼頁，畫面上顯示
「密碼輸入：第 09 頁」；既有截圖見 [`dosv-copy-protection.png`](../images/dosv-copy-protection.png)。

### 2026-08-12 勘誤：空白確認已通過

後續以 DOSBox-X 的 `mouse_emulation=integration` 重新取得真正的 INT 33 點選後，
空白、`0000`、`1234` 都能點「確定」並進入原版開場。沒有修改執行檔、沒有猜測
說明書答案，也不把密碼頁加入 remake；完整三組可追溯結果見
[`18-dosv-password-verification.md`](18-dosv-password-verification.md)。

因此密碼頁不再是 DOS/V 自然流程的技術阻擋。這份文件仍不宣稱同日期／同輸入／同鏡頭的
逐像素 parity：那是尚未執行的同狀態驗收，不是密碼頁無法越過。

## Remake 可重播基準

使用 DOS/V `workplace/orig/dosv`、倚天字型、`scenario=0`、`player=0`、`seed=17`、
`speed=0`，在 640×400 Xvfb 畫布執行：

```text
wlgame -orig workplace/orig/dosv -font workplace/eten \
  -speed 0 -seed 17 -shot /out/remake-natural.png -shot-frames 30
```

輸出已複製到 [`wlgame-dosv-natural-remake-skeleton.png`](../images/wlgame-dosv-natural-remake-skeleton.png)，
SHA-256：

```text
45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24
```

它確認 DOS/V 的 640×400 外框、命令列、右欄 HUD、地圖／banner、字型與固定種子路徑
可由 remake 重播；影片提供了自然遊戲狀態的結構與色彩 oracle，但不提供同一個日期／
鏡頭／輸入狀態的無損原始像素。

## 目標平台 GUI 建置

Windows 使用 `moo2-ebiten:latest`、`GOOS=windows GOARCH=amd64 CGO_ENABLED=0`，
輸出為 `PE32+ executable (console) x86-64`。macOS 使用 `u5cht/osxcross:latest`
與 `MacOSX15.5.sdk`，分別以 `x86_64-apple-darwin24.5-clang`、
`aarch64-apple-darwin24.5-clang` 建置 `GOOS=darwin` 的 amd64／arm64 binary。
本輪暫存產物的 SHA-256 為：Windows amd64
`ea7bbb9747cb37bc47ae716ac1795f9cd15cce5b110eb8f78d3ed602489b8816`、macOS amd64
`a8894b8245cac1895c8341b18082efb9a8e76cbc9fa4dc795f19484fa76de840`、macOS arm64
`50f8c9378b41a595df99e8495a06b78aeda23b05251bb4ae3744557f2366b8e2`。

這些結果封閉的是「程式能以目標 ABI 編譯」、「GUI 入口沒有被 Linux 專用程式碼卡住」與
「候選包不含原版資產」的 gate；尚未封閉 Windows／macOS 真機視窗、輸入法、縮放、音訊
與資產路徑 runtime。候選包已建立，但不把交叉建置寫成平台 parity。

## 後續 gate

若日後取得原版同一日期／鏡頭的無損影像，才可進一步做嚴格同狀態逐像素 diff。若取得
Windows／macOS 原生測試機，則補跑 GUI 啟動、視窗尺寸、鍵鼠輸入、字型載入與短 smoke；
影片視覺 gate 已通過；候選包可供目標平台短 smoke，正式發行仍須等原生 runtime gate。

## 2026-08-11 DOS/V 自然策略骨架調整

依 YouTube 80 秒參考幀重拍 remake，修正右欄常駐結構：

- minimap 下方改成 16 px 紅／藍勢力色標列；
- 上下兩個右欄框共用 8 px 分隔邊，移除原 remake 多出的一格空白；
- 情報框改成 64×64 頭像、君主／首都／軍師三列、信賴度、紅色分隔線與黑底資源區；
- 左側仍固定 432×336 地圖，整體維持 DOS/V 640×400。

最新基準幀為 [`wlgame-dosv-natural-remake-skeleton.png`](../images/wlgame-dosv-natural-remake-skeleton.png)，
並排證據為 [`dosv-skeleton-compare.png`](../promo/dosv-skeleton-compare.png)。此修正把
畫面骨架對回原版參考幀；日期／輸入／地圖狀態仍不同，因此不宣稱同狀態逐像素 parity。

## 2026-08-11 AppImage、經典再現影片與 Android shell

- Linux amd64 AppImage `dist/release-20260811/packages/wolong-remake-linux-amd64-20260811.AppImage`
  以 AppImage tool 建置；AppDir 根目錄 `.desktop`／`AppRun`、deny-list 與
  `APPIMAGE_EXTRACT_AND_RUN=1` Docker／Xvfb 啟動 smoke 通過。AppImage 仍需玩家自備合法
  DOS/V 資料與中文字型，不包含任何原版資產。
- 「經典再現」影片為
  [`wolong-remake-classic-revival.mp4`](../../dist-all/promo/wolong-remake-classic-revival.mp4)，
  以使用者 YouTube 代表幀對照 remake 固定 `seed=17` 代表幀；影片標示原版／remake 與
  「非同狀態逐像素 parity」界線。重播參數沿用 `core=normal`、`cputype=486`、
  `cycles=20000`，`machine=pc98` 只適用 PC-98 oracle，不改寫 DOS/V 結論。
- Android shell debug APK 在 API 35 `google_apis;x86_64`、KVM 模擬器完成安裝／啟動；
  `android-wolong-touch-prototype.png` 是初始畫面，長按 `CONTINUE` 與 `MENU` 的結果分別見
  `android-wolong-touch-after-continue.png`、`android-wolong-touch-after-menu.png`。這只證明
  橫向視口與 dock 輸入殼，不證明完整 Android 核心、自然 clock、事件或存檔 parity。
