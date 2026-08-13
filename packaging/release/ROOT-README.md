# 臥龍傳 remake 完整交付目錄

這是統一的可交付根目錄。所有發行檔、推廣影片、雜湊與驗收紀錄都放在這裡；建置中間檔
不保留。版本：`@RELEASE_VERSION@`。

## 內容

- `packages/`：Linux amd64、Windows amd64、macOS（Intel＋Apple Silicon）完整包，以及
  Linux amd64 AppImage。另有 Linux arm64 的無頭工具伴隨包，並非完整 GUI 平台包。
- `promo/`：三支 60 秒、1280×720 的推廣主片與一支 24 秒、1280×400 的代表幀比較短片；
  其中 `dosv-live-comparison.mp4` 以「讓經典再現」為主軸，包含使用者指定的松崗 DOS/V
  與 remake 實機比較。影片不含原版音訊。
- `verification/`：Linux tar 與 AppImage 的 Xvfb GUI smoke 截圖、ABI 檔頭摘要與封裝檢查。
- `experimental/android/`：已驗證可啟動的 Android 觸控 shell 原型；不是完整遊戲，故不列為
  三平台完整發行之一。
- `SHA256SUMS.txt`：本根目錄所有可交付檔案的 SHA-256。

## 重要界線

桌面包不含任何原版執行檔、資料、美術、音樂、字型或完整原版文字表；玩家必須自行提供
合法松崗 DOS/V 資料與字型。Windows／macOS 是交叉建置候選，尚未在目標作業系統完成原生
GUI 短 smoke；此限制已寫入各包的說明，不以封裝存在取代實機驗證。

Android 原型特意分列為實驗性附件，避免將觸控介面 proof-of-concept 誤稱為完整 Android
遊戲。完整桌面遊戲的三平台定義為 Linux amd64、Windows amd64 與 macOS 雙架構包。
