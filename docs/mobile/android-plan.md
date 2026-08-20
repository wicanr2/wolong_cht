# Android 版規劃

**狀態：路線已定案，核心尚未接入。** 手機端**只共用規則層**，畫面與操作重新設計
（使用者裁定 2026-08-20）。目前 repo 裡的 APK 是**沒有任何遊戲核心的假畫面**——
`mobile/wolong/game.go` 畫的是寫死的矩形。

- 日期：2026-08-20
- UX 規格：[`android-ux.md`](android-ux.md)
- 現有工具鏈：`tools/android-emulator.Dockerfile`、`android/`、`mobile/wolong/`

## 1. 兩個決定（使用者裁定 2026-08-20）

| 決定 | 內容 |
|---|---|
| **UX** | 為手機重新設計 UI，**只共用規則層**（版面與操作重畫，見 [`android-ux.md`](android-ux.md)）|
| **驗證** | **只有 Docker 模擬器**；實機那一格保持未完成 |

⚠ **手機版因此不是 parity 參考。** `CLAUDE.md` §1 允許「現代化外殼」，
但要求每一項改動標記為 remake 差異——**桌面版仍是唯一的對拍基準**，
手機版的畫面差異不進 `docs/playtest/` 的 parity 數字。

⚠ **只有模擬器就照實寫。** 實機那一格保持未完成，不以模擬器截圖冒充
（同 M8 的 Windows／macOS：交叉建置的檔頭不能取代實機 smoke）。

## 2. ⭐ 為什麼「重新設計 UI」的相依關係反而最乾淨

桌面遊戲本體 `cmd/wlgame` 是 **16,587 行、53 個檔、253 個 `game` 方法，
而且是 `package main`**——手機綁定 import 不到它。
若要沿用原版版面，第一步就得把這一整包搬進可 import 的套件。

走重新設計那條路**不必搬**：手機端直接接 `internal/`。逐套件量過：

| 套件 | 平台相依 | 能不能直接給 Android 用 |
|---|---|---|
| `internal/state` | 只有 `os.ReadFile`／`os.WriteFile` | ✅ 給真實路徑就行 |
| `internal/rules/*` | 只有 `os.ReadFile`（`tactical/formation.go` 讀 `KI.EXE`）| ✅ |
| `internal/assets/*` | `os` ＋ `path/filepath`，解碼回 `image.RGBA` | ✅ **不認識 Ebiten** |
| `internal/savepath` | `os` ＋ `path/filepath`，刻意不依賴 Ebiten | ✅ |
| `internal/ui/textdraw` | Ebiten | ✅ gomobile 下 Ebiten 可用 |
| `cmd/wlgame` | `package main` | ❌ **不需要** |

⭐ **這是這條路線最實在的好處**：不必動那 16,587 行，也就不會在搬運過程中
改壞已經量到 0 px 的桌面版。

## 3. ⚠ 必須先解的一件事：原版資料怎麼進手機

`internal/*` 全部用 `os.ReadFile(路徑)` 讀資料，而 **Android 11 以上拿不到
使用者選的資料夾路徑**——SAF 給的是 `content://` URI，不是檔案系統路徑。

所以要有一個**匯入步驟**（Java 側）：

```
使用者用 SAF 選原版資料夾
      ↓  Java: DocumentFile 逐檔讀出
複製進 app 私有目錄 filesDir/orig/
      ↓  Go: os.ReadFile("<filesDir>/orig/KI.EXE") 照常運作
```

三件事跟著這個決定走：

| 項目 | 作法 |
|---|---|
| 要複製哪些檔 | 松崗版 69 檔全部；缺檔要**逐檔列出缺哪一個**，不能只說「載入失敗」 |
| 倚天字型 | 同一個匯入流程；沒有字型就中文顯示方框（桌面版已經是這個行為）|
| 存檔 | 寫 `filesDir/save/`，**不碰匯入來源**（延續「原版資產唯讀」）|

⚠ 現在的 `AndroidManifest.xml` **沒有任何儲存權限也沒有 picker**——
這一段是從零開始。

## 4. 分層

```text
Android Activity（Java）
   ├── SAF 匯入：content:// → filesDir/orig/           ← 新增
   └── EbitenView
            ↓ lifecycle、safe-area、touch
mobile/wolong（gomobile 綁定）
            ↓
internal/ui/mobile   ← 新的手機呈現層（畫面 ＋ 手勢）
            ↓
internal/state ＋ internal/rules ＋ internal/assets    ← 原封不動共用
```

**規則層不會知道自己跑在手機上。** 手機端唯一能對規則層做的事，
是送出與桌面同一組指令；不得為了觸控方便改時鐘、事件順序、數值邊界或存檔格式。

## 5. 里程碑

每一項都有「做完的判準」，判準寫不出來就不算里程碑。

| # | 里程碑 | 做完的判準 |
|---|---|---|
| **A** | **核心在 Android 上真的跑起來** | 同一個 seed 跑 N tick，Android 與桌面算出**相同的指紋**（`World.Fingerprint`，[`../spec/69`](../spec/69-world-fingerprint.md)）。⭐ 這一條不需要畫面就能驗，是整條路線最強的驗收 |
| **B** | SAF 匯入 | 選資料夾 → 69 檔複製進 `filesDir/orig/` → Go 端讀得到；缺檔逐檔列出 |
| **C** | 手機 UI v1 | 大地圖可縮放拖曳、頂部狀態列、底部指令列；日期會走（[`android-ux.md`](android-ux.md) §2–§4）|
| **D** | 進言／一覽／編成 | 三個全螢幕 sheet；送出的是既有指令，不直接改 `World` |
| **E** | 戰場 | 45 度視角沿用原版資產；控制改成「點部隊 → 底部指令」（§5）|
| **F** | 存檔／讀檔 | 四槽，寫 `filesDir/save/`，來源目錄不變；跨啟動可讀回 |
| **G** | 模擬器 smoke | 見 §6 |
| **H** | 實機驗收 | ⛔ **沒有裝置，這一格保持未完成** |

⚠ **A 之前不做任何 UI。** 先證明「核心能在 Android 上算出與桌面相同的結果」，
再談畫面——反過來做的話，畫面出錯時分不出是 UI 還是核心。

## 6. 驗收：模擬器能驗什麼、驗不到什麼

沿用 `wolong-android-emulator:20260811`（Android 35 `google_apis;x86_64`、KVM）。

| 驗得到 | 驗不到 |
|---|---|
| 安裝、啟動、橫向鎖定 | 真實觸控延遲與多點手勢的手感 |
| state hash 與桌面相同 | 真實 GPU driver（模擬器走 SwiftShader／host GL）|
| SAF 匯入流程 | 電量、發熱、長時間執行的記憶體壓力 |
| 存檔跨啟動 | 各家廠商的瀏海／手勢列差異 |
| 不同 AVD 比例（16:9／20:9／平板）| 真實螢幕的可讀性（點陣字在高 DPI 上的觀感）|

⚠ **`adb tap` 的瞬間 down/up 會跨過 Ebiten 的一幀**，2026-08-11 那次就踩到——
驗收要用有界長按重播，**那是量測手段，不代表玩家必須長按**。

### 自動化的形狀

建置與 smoke 目前**都是手動步驟**，工具鏈固定在
[`tools/android-emulator.Dockerfile`](../../tools/android-emulator.Dockerfile)。
里程碑 A 之前要各補一支放進 `tools/`：

| 要寫的 | 做什麼 |
|---|---|
| `android_build.sh` | `gomobile bind` → AAR → `gradle assembleDebug` |
| `android_smoke.sh` | 起 AVD → 安裝 → 推測試資料 → adb 操作 → 截圖 → 比對 |

## 7. 不做的事

- 不做直向版。
- 不為了觸控便利改寫規則層的任何東西。
- 不把手機版寫成「已完整支援」，在 A–G 都通過之前一律標為原型。
- 不內嵌原版資料與倚天字型（deny-list 邊界不變）。
- **不把手機版的畫面差異算進 parity 數字**——那是桌面版的職責。

## 8. 固定工具鏈

Go `1.25.12`、Ebiten `v2.9.9`、Android Gradle Plugin `8.7.3`、
SDK platform `35`、NDK `27.2.12479018`、`gomobile bind -androidapi 29`、
`minSdk 29`／`targetSdk 35`、套件名 `com.wicanr2.wolong`。

## 9. 未解

| 項目 | 現況 |
|---|---|
| 實機驗收 | ⛔ 沒有裝置。里程碑 H 保持未完成 |
| `gomobile bind` 的可重現建置腳本 | 手動步驟，還沒寫成 `android_build.sh` |
| 高 DPI 下的點陣字 | 原版字型是 16×15 點陣，手機上要整數放大幾倍才讀得清楚沒量過 |
| release signing | A–G 之前不談 |
