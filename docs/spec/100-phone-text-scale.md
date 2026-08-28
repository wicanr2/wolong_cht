# 100 — 手機版的字放大 2 倍

**狀態：CONFORMED（2026-08-28 實作、單測、桌面 Xvfb 截圖驗過）。**

- 日期：2026-08-28
- 出處：**remake 差異**（使用者裁定 2026-08-28：「android 版本字太小，希望字體大小適合手機」）。
  原版字型是 16×15 點陣（[`../re/29`](../re/29-font-service-int15.md)），
  手機版版面本來就不是原版的（[`../mobile/android-ux.md`](../mobile/android-ux.md) §1）。
- 推論等級：不涉及原版事實。

## 1. 為什麼

手機版的邏輯畫布是 960×540，Ebiten 等比縮到螢幕。16 px 的字在 6 吋
1080p 手機上只有 18 個實體像素、約 2.5 mm，比 Android 內文預設（14 sp ≈ 3.7 mm）
小得多。`android-ux.md` §6 早就標明「放大幾倍才讀得清楚沒量過」。

## 2. 決定

| 項目 | 值 | 理由 |
|---|---|---|
| 倍率 | **2**（`phone.TextScale`）| 只用整數倍——點陣字非整數縮放會糊（`android-ux.md` §3 捏合縮放同一條規則）。3 倍是 48 px，一列只剩 7 列清單，太少 |
| 放大在哪做 | `textdraw.Drawer.SetScale`，畫的那一刻 `GeoM.Scale` | 字模快取仍是原尺寸；**桌面版不呼叫**，逐像素對拍不受影響 |
| 設在哪 | `phone.Session.Draw` 每幀 `td.SetScale(TextScale)` | Drawer 由 `mobile/wolong` 與 `cmd/wlandroid` 各建一份，設在建構處就會有一邊忘記 |

## 3. 版面跟著字高長

所有列高與留白從 `FontH`（= 15 × 倍率 = 30）與 `LineH`（= 17 × 倍率 = 34）長出來，
不再各自寫死 16：

| 常數 | 以前 | 現在 | 備註 |
|---|---:|---:|---|
| `StatusH` | 56 | `FontH+24` = 54 | 文字垂直置中 |
| `CommandH` | 64 | `FontH+36` = 66 | 扣掉 8 px 間隙後 50 ≥ 48 dp 觸控下限（`layout_test.go`）|
| `rowH`（清單列）| 30 | `LineH+8` = 42 | 一頁 8 列（以前 12）|
| `tabH` | 44 | `FontH+20` = 50 | |
| `CardW`／`CardH` | 300／176 | 360／`CardPadY*2+LineH*5+8` = 202 | 標題 ＋ 四列 |
| `BattleRowH` | 56 | `FontH+24` = 54 | ⚠ 兩列各兩行字塞不下——戰場區要留給原版 480×368 的視野（`TestBattleFieldFitsTheOriginalViewport`）。六格改成「位置 兵數」**一行** |
| 事件提示條列距 | 26 | `LineH+2` = 36 | 點擊熱區同步（`input.go`）|

狀態列右側的「資金／預備兵」改成從右往左用 `td.Width` 排，不再寫死 x。

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 放大 | `internal/ui/textdraw`：`SetScale`／`Scale`／`Width`／`LineH`／`FontH` |
| 常數 | `internal/ui/phone/layout.go`、`sheet.go`、`battle.go`、`notice.go` |
| 差異 | 整份都是 remake 差異；桌面版倍率維持 1 |

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `internal/ui/phone` 全部（觸控下限、戰場視野、關於頁列數）、`internal/ui/textdraw` |
| 截圖 | `tools/phone_shot.sh`：主畫面、據點小卡、武將一覽、系統「關於」、攻城戰場——文字 32 px，六格與六命令一行放得下，清單四欄不撞 |
| 實機 | **未做**（沒有裝置，同 `docs/release/10` §7）|

## 6. 未解

| 項目 | 現況 |
|---|---|
| 倍率不能在遊戲內調 | 固定 2。平板或小手機可能要 1 或 3，得先有實機回饋 |
