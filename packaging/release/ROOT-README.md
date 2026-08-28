# 臥龍傳 remake 完整交付目錄

這是統一的根目錄。所有發行檔、推廣影片、雜湊與驗收紀錄都放在這裡；建置中間檔
不保留。版本：`@RELEASE_VERSION@`。

## 內容

- `packages/`：**四個平台的完整包**——Linux amd64、Windows amd64、
  macOS（Intel＋Apple Silicon）與 Android，另有 Linux amd64 的 AppImage。
  Linux arm64 是無頭工具伴隨包（只有 `wlsim`／`wlshot`），不是完整平台包。
- `promo/`：72 秒的主預告、72 秒的原版實機對照片、60 秒的代表幀比較片、
  48 秒的 Android 版推廣片，與 24 秒、1280×400 的並排短片。主預告與 Android 片的
  遊戲段落是逐幀錄下來的實跑畫面；`dosv-realmachine.mp4` 的原版側是自己跑的受控
  DOSBox-X 實機遊玩。**主預告的配樂是原版的曲子**（由本專案的 OPL3 從使用者自備的
  `BGM.DAT` 算出來，不是取樣原版錄音），50–60 秒的並排段左半是原版畫面——
  兩處都是原版衍生物，見 `promo/README.md`。
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

## 授權

`LICENSE`：**非商業用途免費**（含修改與再散布），**商業用途要先洽談**
（wicanr2@gmail.com）。**不涵蓋原版素材**——那些屬於各自的權利人。
