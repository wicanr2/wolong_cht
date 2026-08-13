# 臥龍傳 remake 發行說明

此封裝只包含 remake 程式與公開的 `translations/corrections.json`。
**不包含**松崗 DOS/V 的執行檔、資料檔、圖像、音樂、完整 TALK 文字表或倚天字型；
請自行準備合法的松崗繁中版資料與中文字型。

## 啟動

Linux／macOS：

```text
./wlgame -orig /path/to/songgang-cht -font /path/to/eten-fonts
```

Windows：

```text
wlgame.exe -orig C:\\path\\to\\songgang-cht -font C:\\path\\to\\eten-fonts
```

`-orig` 指向包含 `SINARIO.DAT`、`MMAP.*`、`*GRF.DAT`、`TALK.DAT`、`KI.EXE`
等松崗 DOS/V 資料的目錄。`-font` 指向玩家自備的字型目錄（含
`STDFONT.15`、`SPCFONT.15`、`ASCFONT.15`）。

可選的持久化存檔必須放在另一個可寫位置：

```text
./wlgame -orig /path/to/songgang-cht -font /path/to/eten-fonts \
  -save-file /path/to/writable/SAVE.DAT
```

`wlsim` 是無頭規則模擬器，`wlshot` 是素材／截圖工具，`wlview` 是互動素材與
大地圖檢視器。每個封裝內的 `SHA256SUMS.txt` 可驗證解壓後檔案。

## 文字校訂與資料界線

`wlgame` 會先讀取玩家自備的原始 `TALK.DAT`，再以同包的公開校訂覆蓋驗證原文後
套用已定案修正。原文與覆蓋不符時會失敗並顯示原因，不會靜默改寫資料。完整
`talk-dosv-corrected.json` 僅供開發驗證，不在此封裝內。

## 平台注意事項

- Linux amd64 包與 AppImage 已完成 Docker／Xvfb 短 GUI 截圖驗收。
- Windows amd64 與 macOS Intel／Apple Silicon 已完成交叉編譯與目標 ABI 檔頭驗證；
  尚未在對應作業系統做原生 GUI、輸入、音訊、縮放與字型短 smoke，第一次使用前請
  在目標平台執行上述啟動命令確認。
- macOS 封裝以 `darwin-amd64/` 與 `darwin-arm64/` 分別提供 Intel 與 Apple Silicon
  執行檔；請進入符合本機架構的目錄執行。

## 行為範圍

主策略畫面的閒置更新會依「據點／軍團／物件／時鐘」順序前進；游標移動、按鈕或命令
發生的同一畫面不推進時間。事件 10 的自然原版 producer 仍屬未知研究項；現行遊戲使用
受控的 remake 近似流程，不將其宣稱為原版已證實行為。
