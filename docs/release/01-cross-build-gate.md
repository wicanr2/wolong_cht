# 01 — 發行閘重跑：五平台交叉建置 ＋ deny-list（2026-08-17）

**狀態：歷史紀錄——2026-08-17 那一輪的發行閘結果。⭐ macOS 的 Ebiten 本體
可以交叉建，`tools/release.sh` 現在自己會做；目標 OS 實跑仍然缺（沒有
Mac／Windows 主機）。後續各批次的產物與驗收見 `02`–`10`，目前是 `10`（`20260828`）。**

- 日期：2026-08-17
- 指令：`tools/release.sh`（五平台）
- 工具：`tools/go.sh`（純 Go）、`wolong-osxcross-go:20260811-event10-r4`（darwin 的 cgo）

## 1. 為什麼要重跑

今天改了不少會進執行檔的東西——行軍指示的三選一、抵達狀態機、
一覽表的捲軸與換色、編成畫面的滑鼠、戰術底列的命令圖示。
**交叉建置是這些改動的第一道體檢**：純 Go 的部分在哪個平台都一樣，
但只要碰到 cgo（Ebiten 的視窗層），平台之間就會分歧。

## 2. 結果

| 平台 | 邏輯工具（`wlsim`／`wlshot`）| 遊戲本體（`wlgame`／`wlview`）|
|---|---|---|
| linux/amd64 | ✅ | ✅（本機 cgo）|
| linux/arm64 | ✅ | ❌ 要在目標平台建 |
| windows/amd64 | ✅ | ✅（純 Go，Ebiten 走 purego）|
| darwin/amd64 | ✅ | ✅ **osxcross 交叉建**，`Mach-O 64-bit x86_64` |
| darwin/arm64 | ✅ | ✅ **osxcross 交叉建**，`Mach-O 64-bit arm64` |

deny-list（發行閘）：**掃了 20 個檔、比對 120 個原版檔的內容雜湊，
沒有原版資產**。

## 3. 兩支發行腳本現在說同一件事

`tools/release.sh`（日常）與 `tools/release_all.sh`（完整交付）都用
osxcross 工具鏈交叉建 macOS 的 Ebiten 本體。映像名走
`$WOLONG_MAC_IMAGE`（預設 `wolong-osxcross-go:20260811-event10-r4`），
**沒有那顆映像就印一行跳過**——不讓整條發行流程卡在一顆選用的工具鏈上。

⚠ 有兩支腳本做同一件事時，**兩邊的能力說明要一起改**。
日常會跑的那一支說「做不到」，等於整個 repo 都以為做不到。

## 4. 還缺什麼

| 項目 | 現況 |
|---|---|
| 目標 OS 實跑 | **做不到**：這台是 Linux，沒有 Mac／Windows。檔頭驗過（PE32+／Mach-O），但視窗、輸入、音訊、字型載入都沒有在目標系統上跑過 |
| linux/arm64 的本體 | 要在 arm64 的 Linux 上建（Ebiten 的 cgo 沒有交叉工具鏈）|
| Windows 的 smoke | 同第一項 |
