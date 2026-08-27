# 發行驗收紀錄

本目錄由 Docker-only 發行流程產生。驗收範圍與結論如下：

- 公開 `corrections.json` 已對玩家自備松崗 `TALK.DAT` fail-closed 套用，並在開發驗收中
  與既有 1,022 則校訂表逐則位元組一致；完整原版文字表未收錄於任何包檔。
- @VERIFY_BOUNDARY@
- `linux-tar-smoke-@RELEASE_STAMP@.png` 是解壓 Linux amd64 包後、從不同工作目錄執行
  `wlgame` 的 Xvfb 固定種子截圖；`appimage-smoke-@RELEASE_STAMP@.png` 是以
  `--appimage-extract-and-run` 啟動 AppImage 的相同 smoke，
  `appimage-ending-@RELEASE_STAMP@.png` 走結局過場那一條路。
- 截圖都是 640×400、固定種子 7、劇本 0、第 120 幀；**同一個局面的三張本來就會同雜湊**，
  雜湊列在發行根目錄的 `SHA256SUMS.txt`，不在這裡重抄一份。
- `BUILD-ABI.md` 列出 Linux／Windows／macOS 的實際交叉建置產物雜湊與目標 ABI 摘要。Linux
  已做 GUI smoke；Windows／macOS 只完成交叉建置與檔頭驗證，尚未完成目標系統 GUI runtime。
- `SHA256SUMS.txt` 位於發行根目錄，列出所有交付檔案（不含其自身）的 SHA-256；請在其上層
  目錄以 `sha256sum -c dist-all/SHA256SUMS.txt` 再驗一次。

> **為什麼這裡不抄雜湊。** 截圖的雜湊每一批都會變（批次日期進檔名、畫面隨程式改），
> 抄一份在說明裡等於保證它過期——而**過期的雜湊比沒有雜湊更糟**，
> 照著驗的人會得出「檔案被動過」的錯結論。單一來源是 `SHA256SUMS.txt`。
