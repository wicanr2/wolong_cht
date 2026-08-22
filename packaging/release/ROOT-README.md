# 臥龍傳 remake 完整交付目錄

這是統一的根目錄。所有發行檔、推廣影片、雜湊與驗收紀錄都放在這裡；建置中間檔
不保留。版本：`@RELEASE_VERSION@`。

## 內容

- `packages/`：**四個平台的完整包**——Linux amd64、Windows amd64、
  macOS（Intel＋Apple Silicon）與 Android，另有 Linux amd64 的 AppImage。
  Linux arm64 是無頭工具伴隨包（只有 `wlsim`／`wlshot`），不是完整平台包。
- `promo/`：三支 60 秒、1280×720 的主片，一支 48 秒的 Android 版推廣片，與一支 24 秒、
  1280×400 的代表幀比較短片。主預告與 Android 片的遊戲段落是逐幀錄下來的實跑畫面；
  `dosv-live-comparison.mp4` 是松崗 DOS/V 與 remake 的實機動態比較。影片不含原版音訊。
- `verification/`：Linux tar 與 AppImage 的 Xvfb GUI smoke 截圖、ABI 檔頭摘要與封裝檢查。
- `SHA256SUMS.txt`：本根目錄所有檔案的 SHA-256。
- `release-manifest.json`：機器可讀的清單，含 `distributable` 與
  `original_assets_included` 兩個旗標。

## 重要界線

@ROOT_BOUNDARY@

Windows／macOS 是交叉建置候選，尚未在目標作業系統完成原生 GUI 短 smoke；
此限制已寫入各包的說明，不以封裝存在取代實機驗證。

Android 版與另外三個平台**並列**，它是完整的遊戲、規則層與桌面版是同一份程式碼。
它仍然只有 debug 簽章、也還沒做實機驗收——那是**驗收狀態**，不是功能缺口。
