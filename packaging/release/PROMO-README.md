# 推廣影片

全部是 H.264／AAC 的 MP4。畫面只有 remake 自己畫出來的東西，配樂是本專案的原創合成
（`tools/promo_score.py`）；沒有原版的音樂、音效或錄影混進來。

| 檔案 | 規格 | 內容 |
|---|---|---|
| `wolong-remake-trailer.mp4` | 1280×720、60 秒 | 主預告。大地圖、野戰、攻城三段是**逐幀錄下來的實跑畫面**；事件視窗、校訂與存檔那幾段是截圖（那些畫面本來就不動）|
| `wolong-remake-android.mp4` | 1280×720、48 秒 | Android 版：縮放、據點小卡、四張一覽、軍團編成、進言、戰場。全片都是實跑錄製 |
| `wolong-remake-dosv-live-comparison.mp4` | 1280×720、60 秒 | 松崗 DOS/V 與 remake 的實機動態比較 |
| `wolong-remake-classic-revival.mp4` | 1280×720、60 秒 | 以「讓經典再現」為主題的代表幀比較 |
| `wolong-remake-yt-comparison.mp4` | 1280×400、24 秒 | 原版代表幀與 remake 的結構／色彩比較。**沒有配樂**——它是研究用的並排圖影片，不是宣傳片 |

比較片的原版側是**代表幀**，不是同日期、同輸入、同狀態的畫面，所以它們證明的是
畫面骨架與 HUD 幾何對得上，不是逐像素 parity。

各片的 SHA-256 見上層 `SHA256SUMS.txt`。製作紀錄與重播參數留在原始專案的 `docs/promo/`。
