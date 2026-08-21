# 臥龍傳 remake 完整交付目錄

這是統一的可交付根目錄。所有發行檔、推廣影片、雜湊與驗收紀錄都放在這裡；建置中間檔
不保留。版本：`@RELEASE_VERSION@`。

## 內容

- `packages/`：Linux amd64、Windows amd64、macOS（Intel＋Apple Silicon）完整包，以及
  Linux amd64 AppImage。另有 Linux arm64 的無頭工具伴隨包，並非完整 GUI 平台包。
- `promo/`：三支 60 秒、1280×720 的主片，一支 48 秒的 Android 版推廣片，與一支 24 秒、
  1280×400 的代表幀比較短片。主預告與 Android 片的遊戲段落是逐幀錄下來的實跑畫面；
  `dosv-live-comparison.mp4` 是松崗 DOS/V 與 remake 的實機動態比較。影片不含原版音訊。
- `verification/`：Linux tar 與 AppImage 的 Xvfb GUI smoke 截圖、ABI 檔頭摘要與封裝檢查。
- `experimental/android/`：Android 版的 **debug 包**。它是完整的遊戲，規則層與桌面版
  是同一份程式碼；列在這裡而不是併進三平台，是因為它只有 debug 簽章、也還沒做實機驗收。
- `SHA256SUMS.txt`：本根目錄所有可交付檔案的 SHA-256。

## 重要界線

桌面包不含任何原版執行檔、資料、美術、音樂、字型或完整原版文字表；玩家必須自行提供
合法松崗 DOS/V 資料與字型。Windows／macOS 是交叉建置候選，尚未在目標作業系統完成原生
GUI 短 smoke；此限制已寫入各包的說明，不以封裝存在取代實機驗證。

Android 版的界線是**簽章與驗收**，不是功能：它只在 Docker 的模擬器上驗過，沒有實機
驗收，也還沒有上架用的正式簽章。桌面三平台仍定義為 Linux amd64、Windows amd64 與
macOS 雙架構包。
