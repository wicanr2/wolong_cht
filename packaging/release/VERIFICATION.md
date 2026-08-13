# 發行驗收紀錄

本目錄由 Docker-only 發行流程產生。驗收範圍與結論如下：

- 公開 `corrections.json` 已對玩家自備松崗 `TALK.DAT` fail-closed 套用，並在開發驗收中
  與既有 1,022 則校訂表逐則位元組一致；完整原版文字表未收錄於任何包檔。
- `packages/` 已經過封裝內容檢查與 deny-list 掃描；原版執行檔、資料檔、美術、音樂、字型與
  `talk-dosv-corrected.json` 都不應存在。
- `linux-package-smoke.png` 是解壓 Linux amd64 包後、從不同工作目錄執行 `wlgame` 的 Xvfb
  固定種子截圖，驗證同包 `translations/corrections.json` 可由執行檔位置自動找到。
- `appimage-smoke.png` 是以 `APPIMAGE_EXTRACT_AND_RUN=1` 啟動 AppImage 的相同 GUI smoke。
- 兩張 smoke 截圖皆為 640×400，SHA-256 都是
  `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。
- `BUILD-ABI.md` 列出 Linux／Windows／macOS 的實際交叉建置產物雜湊與目標 ABI 摘要。Linux
  已做 GUI smoke；Windows／macOS 只完成交叉建置與檔頭驗證，尚未完成目標系統 GUI runtime。
- `SHA256SUMS.txt` 位於發行根目錄，列出所有交付檔案（不含其自身）的 SHA-256；請在其上層
  目錄以 `sha256sum -c dist-all/SHA256SUMS.txt` 再驗一次。
