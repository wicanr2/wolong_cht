# YouTube／remake 推廣片像素差異審查

**狀態：已完成影片對照；確認像素差異，但不把不同遊戲狀態誤稱為同狀態
逐像素 parity。**

- 日期：2026-08-11
- 原版參考：使用者提供的 YouTube 錄影
  [af6xqcicXoI](https://www.youtube.com/watch?v=af6xqcicXoI) 的已保存代表幀。
- remake 參考：`docs/images/wlgame-dosv-natural-remake-skeleton.png`，640×400；這是本輪
  依原版右欄骨架調整後的最新基準幀。
- 比較條件：原版 80 秒代表幀已去黑邊／還原為 640×400；remake 使用同尺寸自然策略
  基準幀。原版日期為 196/4/5，remake 為 196/4/1，因此只作視覺差異量測。

## 產物

- [remake 推廣片](../../dist/promo/wolong-remake-trailer.mp4)：1280×720，H.264／AAC（本文寫作時是 60 秒版；現行主預告是 72 秒，見 [`README.md`](README.md)）。
- [YouTube／remake 對照片](../../dist/promo/wolong-remake-yt-comparison.mp4)：24 秒、
  1280×400、H.264；六組原版代表幀與 remake 驗收畫面並排，無原版音訊。
- [自然畫面並排圖](yt-remake-natural-side-by-side.png)。
- [自然畫面差異圖](yt-remake-natural-difference.png)。
- [DOS/V 骨架對照圖](dosv-skeleton-compare.png)。
- 可重現腳本：[tools/promo_yt_compare.sh](../../tools/promo_yt_compare.sh)。

對照片只使用已保存的 YouTube 代表幀，沒有把原版影片、原版音樂或原版執行檔放進
發行素材。它是研究／驗收產物，不取代正式推廣片。

## 實測結果

以 ImageMagick `compare` 對 640×400 RGB 畫面量測：

| 區域 | 不同像素數 | 區域像素數 | 不同比例 | RMSE（正規化） |
|---|---:|---:|---:|---:|
| 全畫面 | 249,178 | 256,000 | 97.34% | 0.329145 |
| 上方 32 px 橫幅 | 19,792 | 20,480 | 96.64% | 0.256230 |
| 下方 32 px 命令列 | 18,980 | 20,480 | 92.68% | 0.369600 |
| 左側 432×336 地圖 | 145,024 | 145,152 | 99.91% | 0.329224 |
| 右側 208×336 資訊欄 | 65,382 | 69,888 | 93.55% | 0.335341 |

原始輸出為 `AE=249178`、`RMSE=21570.5 (0.329145)`。右欄與命令列的不同像素數
在骨架調整後下降；這證實使用者要求的「把
推廣片／remake 畫面與 YouTube 原版遊玩畫面並列，看出像素差異」已完成；差異圖也
保留了具體可回看的像素層結果。

## 本輪 DOS/V 骨架調整

- 常駐畫布固定為 640×400：32 px banner、32 px 命令列、左側 432×336 地圖、右側
  208 px 欄位。
- minimap 下方改成原版 16 px 紅／藍勢力色標列；下方情報框與 minimap 共用一條
  8 px 分隔邊，避免原 remake 產生雙倍 16 px 假分隔。
- 情報框改為原版常駐順序：64×64 頭像、君主／首都／軍師三列、信賴度、紅色分隔線，
  以及黑底資金／預備兵區。（中央那一直排圖示後來解出來了，見檔尾。）

## 解讀與範圍決定

這些數值不能直接當成 renderer 的錯誤率，原因是：

- YouTube 來源是 478×360、有損壓縮／縮放的錄影；
- 兩張畫面日期不同，輸入與地圖／部隊狀態也不同；
- 推廣片另有 1280×720 的縮放、黑邊與字幕，因此 raw metric 使用推廣片所採用的
  640×400 remake 原始幀，而不是把字幕與黑邊算進去。

因此本專案採用的結論是：YouTube／推廣片比較已足以驗收 DOS/V 畫面骨架、色彩層級、
HUD 位置與整體視覺方向；嚴格同狀態逐像素 diff 仍維持 **未宣稱**，不再以它阻擋
remake，也不繞過 DOS/V 密碼保護取得原版畫面。

## 雜湊

| 產物 | SHA-256 |
|---|---|
| `dist/promo/wolong-remake-trailer.mp4` | `72ffc20f0ab22dc5a43771d115b944d57c45477973640e83304fddb1545292fd`（2026-08-11 當時的版本；現行雜湊見 `dist-all/SHA256SUMS.txt`）|
| `dist/promo/wolong-remake-yt-comparison.mp4` | `3efe64fc7ff903ed15850a1c44c09d8ec82ff2d3a1a38681ae376e46ac84aa9a` |
| `wlgame-dosv-natural-remake-skeleton.png` | `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24` |
| `yt-remake-natural-difference.png` | `a16d3878766d4428dc9bf09f2de1ee83802486df3a3ca532e60af8e6f5cb6872` |
| `yt-remake-natural-side-by-side.png` | `df76ffc1cdb21be11b309d746a8526d99ba20e04fbb124fae8384e3e27b5cd4e` |
| `dosv-skeleton-compare.png` | `7e78ff252d1579a29f2ce93e10c9feabfdb1a392d4c67306812acc44d594042a` |

## 這一輪之後怎麼了

當時唯一列為未解的「中央 raw reserve glyph」已經解出來了：資金／預備兵欄
左邊那一直排是 `ICONGRF` 段 3 `0x1BA0` 起的四張 24×16（天秤／馬／弓／步），
`library.DOSVResourceIcon` 用的是原版素材
（[`../re/48`](../re/48-window-display-list.md) §6）。那一區後來逐像素 PASS
（[`../playtest/37`](../playtest/37-main-screen-parity.md)）。

<!-- 缺口：無 -->
