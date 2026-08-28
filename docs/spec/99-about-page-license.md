# 99 — 手機版「關於」頁要顯示授權條款

**狀態：CONFORMED（2026-08-28 實作並有單測）。**

- 日期：2026-08-28
- 出處：**remake 差異**，不對應原版機制（原版沒有「關於」頁）。
  授權政策見根目錄 `LICENSE` 與 [`../release/10`](../release/10-full-20260828.md) §6。
- 推論等級：不涉及原版事實。

## 1. 為什麼要做

授權條款要跟著每一份發行包走（`~/.claude/rulebook/85`）。桌面四個包與
AppImage 都帶了 `LICENSE`，**APK 帶不了**——gradle 的 assets 塞得進去但沒有
人看得到。拿到 APK 的人唯一看得到字的地方是遊戲畫面，所以條款摘要要顯示在
遊戲內。手機版系統面板本來就有「關於」分頁（`docs/spec/86` §4 的第三頁旁邊），
先前只有三列。

## 2. 內容

「關於」分頁改成下列各列（每列一句，不捲頁也放得下）：

| 列 | 內容 |
|---|---|
| 臥龍傳 Remake | 版本字串（`phone.BuildVersion`，APK 由 `tools/android_build.sh` 以 `-X` 注入；未注入時是 `dev`）|
| 原版 | `NEO･GETEN 1994 / 松崗 1995` |
| 資料來源 | 使用者自備，不隨程式散布（暗色）|
| 授權 | 專有授權：非商業免費 |
| 可以 | 使用、修改、再散布（非商業）|
| 商業使用 | 需書面授權，歡迎來談 |
| 聯絡 | `wicanr2@gmail.com` |
| 不涵蓋 | 原版執行檔、資料、美術、音樂、字型（暗色）|
| 條款全文 | `github.com/wicanr2/wolong_cht` 的 `LICENSE` |

⚠ **不叫 open source**——非商業限制不符合 OSI 定義，畫面上一律寫「專有授權」。
⚠ 文字要挑**倚天 Big5 有的字**（同 `languageRows` 的「●」那條），
英文與 `@` 走半形字型，沒問題。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 列內容 | `internal/ui/phone/sheet.go` 的 `aboutRows()`（`systemRows` 第 3 頁）|
| 版本注入 | `internal/ui/phone.BuildVersion`；`tools/android_build.sh` 的 ebitenmobile `-ldflags` 加 `-X` |
| 差異 | 整頁都是 remake 差異 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestAboutTabShowsLicenseTerms`（`internal/ui/phone`）：第 3 頁含「非商業」「wicanr2@gmail.com」「不涵蓋」，且不含 `open source` |
| 對原版 | 不適用（remake 專屬畫面）|

## 5. 未解

| 項目 | 現況 |
|---|---|
| 條款全文沒有在手機上顯示 | 只顯示摘要與全文出處。全文 104 行，要另做可捲動的文字頁；摘要已滿足「收到的人知道自己被授權了什麼」 |
