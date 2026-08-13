# Android 版規劃

**狀態：規劃已啟動；觸控 shell 原型 debug APK 已產出，並已在 Android 模擬器完成安裝、啟動與有限觸控 smoke；完整遊戲核心尚未接入 Android。**

- 日期：2026-08-11

本規劃以現有 Go／Ebiten 核心為基礎，手機只新增平台殼、視口與輸入轉接層。規則、時鐘、存檔格式、TALK 排版與 DOS/V 640×400 原版模式共用同一份核心，避免 Android 分叉出另一套遊戲邏輯。

## 第一版產品決定

| 項目 | 初始決定 | 原因 |
|---|---|---|
| 姿態 | 橫向優先；第一個可玩版鎖定橫向 | DOS/V 原版是 640×400，橫向能保留完整地圖與右欄 |
| 畫面 | 保留 640×400 邏輯畫布，整數優先、等比例縮放、必要時留黑邊 | 原版外框、游標、按鈕 glyph 與 TALK 分頁不被手機比例拉伸 |
| 操作 | 手指觸控層覆蓋在畫面外圍；核心畫面不縮成桌面 HUD 的小副本 | 觸控 hitbox 可達 48 dp，避免誤觸 |
| 模式 | `原版模式`＋`手機模式` | 原版模式驗收 DOS/V；手機模式解決指令、地圖與 TALK 的手指操作 |
| 目標 ABI | 原型 APK 以 `arm64-v8a` 與 `x86_64` 模擬器路徑驗證；最低 Android API 固定為 29 | 先封閉可重現的模擬器／裝置路徑，再決定 release signing 與商店 ABI |
| 音訊／素材 | 不內嵌原版資料與倚天字型；由使用者選取合法資料目錄，存檔寫入 Android 可寫目錄 | 延續 deny-list 與原始資產唯讀邊界 |

## 2026-08-11 原型實際驗證

- `android/app/build/outputs/apk/debug/app-debug.apk` 已以固定 Docker 工具鏈建置；套件名為
  `com.wicanr2.wolong`，`minSdk=29`、`compileSdk/targetSdk=35`。APK SHA-256 為
  `1bac716a49dfbba486c7386dd30e4f9a04d0b49cd088cbdf30d1257c0cc02d76`。
- 使用 `wolong-android-emulator:20260811`、Android 35 `google_apis;x86_64`、KVM 加速的
  1080×1920 模擬器；應用程式依 manifest 鎖定橫向後，實際畫面為 1920×1080，保留黑邊與
  640×496（含手機底部 dock）的邏輯 shell。
- 已通過安裝、啟動與畫面截圖；長按 `CONTINUE` 顯示 `TALK page 1/3`，長按 `MENU`
  開啟 `MOBILE COMMAND DRAWER` 並顯示 `COMMAND DRAWER OPEN`。`adb tap` 的瞬間 down/up
  可能跨過 Ebiten frame，因此驗收使用有界長按重播，不代表玩家必須長按。
- 證據畫面：[`android-wolong-touch-prototype.png`](../images/android-wolong-touch-prototype.png)、
  [`android-wolong-touch-after-continue.png`](../images/android-wolong-touch-after-continue.png)、
  [`android-wolong-touch-after-menu.png`](../images/android-wolong-touch-after-menu.png)。
- 這是手機操作殼的可丟棄原型，不包含完整自然時鐘、事件 2–10、存檔／讀檔、原版資料匯入、
  實機效能或 release signing；不把本次 smoke 寫成 Android 完整支援。

固定工具鏈摘要：Go `1.25.12`、Ebiten `v2.9.9`、Android Gradle Plugin `8.7.3`、
SDK platform `35`、NDK `27.2.12479018`、`gomobile bind -androidapi 29`。模擬器映像與
依賴的建置說明保留在 [`tools/android-emulator.Dockerfile`](../../tools/android-emulator.Dockerfile)。

## 手機畫面配置

1. 先由安全區域計算 16:10 遊戲視口；系統瀏海、手勢列與狀態列不得被互動區覆蓋。
2. 寬度足夠時，在右側放可收合的指令／情報抽屜；寬度不足時改放底部抽屜。兩者都只送出既有遊戲指令，不直接改 `World` 欄位。
3. 地圖維持原版 27×21 格的邏輯座標。點一下等同原版游標選取；拖曳在沒有進入選取狀態時平移地圖；長按開啟該格的上下文指令。
4. 指令按鈕使用原版 glyph 的顏色與語意作視覺線索，但外層 hitbox 擴成至少 48×48 dp；數值 3×6 面板以整格點擊，不要求精準碰到 16×16 原始像素格。
5. TALK／通知使用左右滑動翻頁、點擊空白處前進；事件選項仍保留兩段式確認，避免滑動誤觸直接送出外交或撥款。
6. Android 返回鍵依序為「關閉上下文／數值／TALK 視窗 → 回到系統視窗 → 顯示離開確認」，不直接終止遊戲或偷偷推進時鐘。
7. 長按、滑動、拖曳與點擊使用明確的手勢狀態機；同一根手指的 `down/up` 不得同時觸發地圖選取與指令按鈕。

## 軟體分層

```text
Android Activity / Ebiten mobile shell
        ↓ lifecycle、safe-area、touch、IME
mobile input/layout adapter
        ↓ InputIntent、ViewportTransform、TouchHit
cmd/wlgame 共用的 presenter／TALK／畫面元件
        ↓
internal/state + internal/rules + internal/assets
```

- 桌面滑鼠、鍵盤與 Android 觸控都轉成同一組 `InputIntent`；規則層不接觸 Android API。
- `ViewportTransform` 保存螢幕座標 ↔ 640×400 邏輯座標的縮放、黑邊與安全區；所有觸控測試都以這個轉換器為唯一入口。
- `TouchHit` 只產生「按哪個既有控制項」的結果；原版模式可關閉手機抽屜，手機模式才顯示額外操作提示。
- app 進入背景時先停止繪製與輸入，不呼叫 `World.Tick`；回到前景後由同一個固定 tick 來源恢復。存檔仍走現有 overlay，不能寫回唯讀原始資料。
- 文字輸入、檔案選取與存檔位置由平台殼提供 callback；`internal` 套件不引入 Android framework。

## 里程碑與短驗收

### A. 可丟棄的版面原型

- 只畫 640×400 代表畫面、右側／底部抽屜、數值面板與 TALK 分頁。
- 用固定假資料驗證安全區、48 dp hitbox、最小／最大常見手機比例。
- 產出一張 Android wireframe 截圖與座標表；不進核心規則。

### B. 觸控轉接層

- 實作 `ViewportTransform`、`TouchGesture`、`TouchHit` 與返回鍵狀態機。
- 單元測試涵蓋黑邊點擊不命中、縮放後地圖格命中、長按／滑動互斥、TALK 翻頁與事件 2–5 的二段確認。
- 桌面測試仍能用合成 `InputIntent` 重播同一組操作，確保不必依賴 Android 才能驗收。

### C. Android 建置橋

- 先固定 Go／Ebiten／Android SDK／NDK／`gomobile` 版本，建立可重現的 Docker 工具鏈；工具鏈未固定前不提交產生的 APK。
- 先出未簽署 debug APK，確認 `arm64-v8a` 安裝、啟動、橫向鎖定與 app pause/resume。
- 原始資料與字型以使用者選取的外部目錄／應用程式可寫目錄掛接；啟動缺檔要顯示可理解的中文錯誤，不得 panic。

### D. Android 短 smoke

- 新遊戲 → 地圖點擊 → 編成／行軍 → 無輸入時日期流逝 → TALK 逐頁 → 系統儲存／讀取 → 返回鍵關閉視窗。
- 同一固定 seed 在桌面與 Android 產生相同 state hash；畫面只驗收 640×400 原版模式與手機控制層，不要求手機 GPU 像素與 Xvfb 相同。
- 旋轉、背景、低電量與觸控中斷都必須不重複送出指令、不遺失存檔、不額外推進 clock。

### E. 實機與發布

- 至少一台窄螢幕 Android 手機與一台大螢幕 Android 平板各跑一次；記錄解析度、DPI、GPU、Android 版本與 APK hash。
- 完成 debug → release signing、版本號、隱私／素材聲明與三平台共用的玩家資料說明後，才建立 Android release gate。

## 暫不做的項目

- 第一版不做直向版、不把 640×400 畫面裁成全螢幕、不改成自由縮放地圖。
- 不用觸控便利性改寫原版即時時鐘、事件順序、數值邊界或存檔格式；便利操作只能在輸入層實現。
- 不把 Android 原型寫成「已完整支援」；在完整核心 smoke、實機驗收與 release signing 通過前，
  文件與發行包只標為觸控 shell 原型。

## Android 驗收清單

- [x] `ViewportTransform` 與安全區測試
- [x] 觸控 dock hitbox 與基本輸入轉接測試（完整手勢／返回鍵仍待做）
- [x] 手機模式 wireframe 與 640×400 邏輯畫布截圖
- [x] Docker 固定 Android SDK／NDK／`gomobile` 工具鏈
- [x] `arm64-v8a`／`x86_64` debug APK 安裝、啟動與橫向畫面
- [x] `CONTINUE`／`MENU` 長按觸控 smoke
- [ ] 觸控、TALK、無輸入 clock、存檔／讀檔 smoke
- [ ] 實機／平板驗收與 release signing
