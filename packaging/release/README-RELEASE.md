# 臥龍傳 remake 發行說明

@PKG_BOUNDARY@

## 啟動

@PKG_LAUNCH@

`-orig` 指向包含 `SINARIO.DAT`、`MMAP.*`、`*GRF.DAT`、`TALK.DAT`、`KI.EXE`
等松崗 DOS/V 資料的目錄。`-font` 指向字型目錄（含
`STDFONT.15`、`SPCFONT.15`、`ASCFONT.15`）。

可選的持久化存檔必須放在另一個可寫位置：

```text
./wlgame -save-file /path/to/writable/SAVE.DAT
```

`wlsim` 是無頭規則模擬器，`wlshot` 是素材／截圖工具，`wlview` 是互動素材與
大地圖檢視器。每個封裝內的 `SHA256SUMS.txt` 可驗證解壓後檔案。

⚠ **只有 `wlgame` 會自己找旁邊的資料目錄。** 另外三支是工具，
它們的路徑一律要明講——那是刻意的，工具的呼叫端本來就知道自己在讀哪一份。

## 文字校訂與資料界線

`wlgame` 會先讀取原始 `TALK.DAT`，再以同包的公開校訂覆蓋驗證原文後
套用已定案修正。原文與覆蓋不符時會失敗並顯示原因，不會靜默改寫資料。完整
`talk-dosv-corrected.json` 僅供開發驗證，不在此封裝內。

## 平台注意事項

- Linux amd64 包與 AppImage 已完成 Docker／Xvfb 短 GUI 截圖驗收。
- Windows amd64 與 macOS Intel／Apple Silicon 已完成交叉編譯與目標 ABI 檔頭驗證；
  尚未在對應作業系統做原生 GUI、輸入、音訊、縮放與字型短 smoke，第一次使用前請
  在目標平台執行上述啟動命令確認。
- macOS 封裝以 `darwin-amd64/` 與 `darwin-arm64/` 分別提供 Intel 與 Apple Silicon
  執行檔；請進入符合本機架構的目錄執行。資料目錄在包根，兩個架構共用一份。

## 行為範圍

主策略畫面的閒置更新會依「據點／軍團／物件／時鐘」順序前進；游標移動、按鈕或命令
發生的同一畫面不推進時間。事件 10 的自然原版 producer 仍屬未知研究項；現行遊戲使用
受控的 remake 近似流程，不將其宣稱為原版已證實行為。

## 授權

引擎與文件的授權條款在同一個目錄的 `LICENSE`：**非商業用途免費**（可以用、
可以散布，也可以改了再散布，條件是保留條款、標明改動與出處、不對作品本身收費）；
**商業用途要先洽談**（wicanr2@gmail.com）。

⚠ **授權不涵蓋原版素材**——原版的執行檔、資料檔、美術、音樂與點陣字型
屬於各自的權利人，要玩得自備合法的原版副本。
