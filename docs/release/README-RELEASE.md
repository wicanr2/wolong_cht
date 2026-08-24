# 臥龍傳 remake 可執行封裝

**狀態：四平台完整包、Linux AppImage、推廣片與驗收紀錄已集中於
[`dist-all`](../../dist-all)，目前是一致的 `wolong-remake-20260824` 批次
（`07-full-20260824.md`）；Windows／macOS 原生 GUI 與 Android 實機都尚未驗證。
⛔ 本機那一批內含原版資產，不可外流。**

- 日期：2026-08-24
- 每一批的產物、雜湊與驗收：`docs/release/` 逐批一份，最新是
  [`07`](07-full-20260824.md)

## ⚠ 包裡那份說明的唯一來源是模板

使用者實際拿到的 `README-RELEASE.md` 由
[`packaging/release/README-RELEASE.md`](../../packaging/release/README-RELEASE.md)
產生——它帶 `@PKG_BOUNDARY@`／`@PKG_LAUNCH@` 兩個佔位符，
`tools/release_all_fs.py` 依「這一批含不含原版資料」填不同的文字。

**這一份不再複製那些內容。** 先前這裡有一份手抄版，日期停在 `20260812`、
散布界線還寫著「原版資料不隨包散布」——而 `dist-all` 從 2026-08-22 起
就是內含遊戲檔案的完整版。⭐ **同一段文字放兩個地方，過期的那一份不會有人發現**。

## 啟動

⚠ **完整版不必給任何資料旗標**——`wlgame` 會自己找執行檔旁邊的
`gamedata/`、`fonts/` 與 `audio/`（`docs/spec/72` §3、`docs/spec/75`）。
不含資料的可散布版才要 `-orig`／`-font`：

```text
wlgame -orig /path/to/dosv -font /path/to/eten-font \
  -save-file /path/to/writable/SAVE.DAT
```

**完整的旗標說明以包裡那一份為準**（來源是 `packaging/release/README-RELEASE.md`，
見上一節）。`wlsim` 是無頭規則模擬器、`wlshot` 是素材／截圖工具、
`wlview` 是互動素材與大地圖檢視器。

## 事件 10 與自然時間

在主策略畫面，游標連續穩定且玩家沒有下達命令／按下滑鼠按鈕時，世界才會自然更新；
已下達目的地的軍團會繼續行軍，據點、物件與日期依「據點／軍團／物件／時鐘」順序前進。
游標移動或命令發生的同一畫面不會偷跑時間。已排入的事件 10 仍由每時 queue 節拍轉成
TALK；月結俘虜通知是可關閉的 remake substitute，不宣稱已找到松崗原版的自然 `0x0A`
producer。

## 平台內容

- Linux：amd64 主封裝含原生 `wlgame`／`wlview`、`wlsim`、`wlshot`；另有 Linux arm64 邏輯工具伴隨包，含 `wlsim`／`wlshot`。另提供 AppImage（檔名跟著批次走，最新一批見 `07-full-20260824.md`）；arm64 的 Ebiten GUI 需在目標 Linux 原生工具鏈建置。
- Windows：amd64 封裝含 PE32+ `wlgame.exe`／`wlview.exe`、`wlsim.exe`、`wlshot.exe`。
- macOS：封裝同時含 Intel (`darwin-amd64`) 與 Apple Silicon (`darwin-arm64`) 目錄，各含 `wlgame`、`wlview`、`wlsim`、`wlshot`。macOS 的 Ebiten 本體要 cgo，由 osxcross 工具鏈交叉建置，`tools/release.sh` 與 `tools/release_all.sh` 都會做；沒有那顆映像時只會少掉 `wlgame`／`wlview`，其餘平台照跑。

最近一次五平台重跑與 deny-list 結果見 [`01-cross-build-gate.md`](01-cross-build-gate.md)。

Linux tar 封裝與 AppImage 已在 Docker／Xvfb 完成啟動與短截圖 smoke；兩者都由解壓／掛載後的執行檔自動找到同包公開校訂覆蓋。Windows／macOS 本輪完成目標 ABI 交叉建置與檔頭驗證，但沒有把原生視窗、輸入、音訊與字型載入寫成已在目標作業系統實機驗證；第一次啟動前請在目標平台做短 smoke。

## 完整性與界線

每個封裝旁的 `SHA256SUMS.txt` 是對解包檔案的雜湊紀錄；[`dist-all/SHA256SUMS.txt`](../../dist-all/SHA256SUMS.txt) 也列出 AppImage、影片與驗收輸出。AppImage 需要 Linux 主機的 X11／GL／ALSA 執行環境，並仍須由玩家以參數指定合法原版資料與中文字型。封裝不含原始資產，deny-list 已掃描候選產物。DOS/V 密碼保護不再是影片比較的技術阻擋，但沒有同狀態無損原始畫面，故逐像素 parity 仍不可宣稱。

## 統一交付目錄

[`dist-all`](../../dist-all) 是本輪唯一的可交付根目錄，內含三平台桌面包、AppImage、五支推廣影片、雜湊、Linux GUI smoke 截圖與 Android 版的 debug 包。Android 的界線是簽章與實機驗收，不是功能：它是完整的遊戲，但只有 debug 簽章，也只在模擬器上驗過。

Android 版是**為手機重新設計的介面**（960×540 邏輯畫布、狀態列 ＋ 四個入口的指令列），不是把 640×400 的桌面版面搬過去；規則層與桌面版共用同一份程式碼。操作規格與模擬器 smoke 見 [`docs/mobile/android-plan.md`](../mobile/android-plan.md) 與 [`docs/mobile/android-ux.md`](../mobile/android-ux.md)。

## 未解

| 項目 | 現況 |
|---|---|
| Windows／macOS 原生 GUI | 交叉建置的產物只驗了檔頭，沒有在目標作業系統跑過。M8 唯一的閘 |
| Android 實機驗收 | 只有 Docker 模擬器；觸控手感、真實 GPU、高 DPI 上的點陣字可讀性都驗不到 |
| Android 正式簽章 | 出的是 debug 簽章，keystore 怎麼保管還沒決定 |
| 16 KB page size 裝置 | `.so` 的 LOAD 段已是 `0x4000`，但沒有那種裝置或 AVD 實際載過 |
