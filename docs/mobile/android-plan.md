# Android 版規劃

**狀態：核心已接入，手機 UI 的主畫面與四個入口可用；戰場與 SAF 匯入未做。**
手機端**只共用規則層**，畫面與操作重新設計（使用者裁定 2026-08-20）。

- 日期：2026-08-20
- UX 規格：[`android-ux.md`](android-ux.md)
- 建置：`tools/android_build.sh`（`docker/android/Dockerfile` → `wolong-go-android:20260820`）
- 驗收：`tools/android_smoke.sh`（`tools/android-emulator.Dockerfile` → `wolong-android-emulator:20260820`）
- 專案檔：`android/`、`mobile/wolong/`、`internal/ui/phone/`

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
匯入這一段是從零開始。

**Go 端的介面已經備好了**：`Wolongmobile.setDataRoot(String)` 收一個根目錄，
底下 `orig/` 放 69 個原版檔、`eten/` 放點陣字、`save/` 是存檔。
Java 端優先給 `getFilesDir()`（`orig/` 存在時），否則退回外部目錄。

⚠ 開局**延後到第一次 Update**，不在 `init()` 做：Android 的資料路徑要由 Java
算出來，而 `mobile.SetGame` 必須在 `init()` 呼叫。這也讓「還沒匯入資料」
變成一個畫得出來的狀態，不是崩潰。

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
| **C** | 手機 UI v1 | ✅ 大地圖可縮放拖曳、頂部狀態列、底部指令列；日期會走（[`android-ux.md`](android-ux.md) §2–§4）|
| **D** | 進言／一覽／編成 | 🔵 進言與一覽做完了，**編成未做**；送出的是既有指令，不直接改 `World` |
| **E** | 戰場 | 45 度視角沿用原版資產；控制改成「點部隊 → 底部指令」（§5）|
| **F** | 存檔／讀檔 | ✅ 四槽，寫 `<root>/save/`，來源目錄不變；區塊 byte-for-byte 與游標都驗過 |
| **G** | 模擬器 smoke | 見 §6 |
| **H** | 實機驗收 | ⛔ **沒有裝置，這一格保持未完成** |

⚠ **畫面在桌面上開發與驗收**（`tools/phone_shot.sh`，一輪約 30 秒），
Android 上的第一個驗收是 A 的指紋。**A 通過之前不宣稱手機版能跑**——
沒有那條斷言的話，畫面出錯時分不出是 UI 還是核心。

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

| 腳本 | 做什麼 |
|---|---|
| [`tools/android_build.sh`](../../tools/android_build.sh) | `ebitenmobile bind` → AAR → `gradle assembleDebug` |
| [`tools/android_smoke.sh`](../../tools/android_smoke.sh) | 起 AVD → 安裝 → 推原版資料 → 抓 logcat 指紋 → 截圖 |

四個踩過的坑，都寫在腳本裡：

| 症狀 | 成因 |
|---|---|
| Java 檔每個中文字一行 `unmappable character` | gobind 把 Go 註解抄進 Java，javac 沒帶 `-encoding`，用**平台預設字集**——容器是 US-ASCII |
| `missing go.sum entry for golang.org/x/tools` | `ebitenmobile` 自己的相依。要用 `go install pkg@version` 裝（在主模組之外解析），不能 `go run` |
| adb 推進 `/sdcard/Android/data/<pkg>/` 之後 app `permission denied` | Android 11 以上那條路徑是 FUSE 掛的。改推 `/data/local/tmp` 再 `run-as` 複製進內部目錄 |
| logcat 只留 `The application was killed due to context loss` | 機器忙時 System UI 被判 ANR，對話框搶走焦點導致 GL context 遺失。smoke 先關掉錯誤對話框 |

## 7. 不做的事

- 不做直向版。
- 不為了觸控便利改寫規則層的任何東西。
- 不把手機版寫成「已完整支援」，在 A–G 都通過之前一律標為原型。
- 不內嵌原版資料與倚天字型（deny-list 邊界不變）。
- **不把手機版的畫面差異算進 parity 數字**——那是桌面版的職責。

## 8. 固定工具鏈

Go `1.25.12`、Ebiten `v2.9.9`、Android Gradle Plugin `8.7.3`、
SDK platform `35`、NDK `27.2.12479018`、`ebitenmobile bind -androidapi 29`、
`minSdk 29`／`targetSdk 35`、套件名 `com.wicanr2.wolong`、
AAR 的 Java 套件 `com.wicanr2.wolong.mobile`。

模擬器用 `system-images;android-34;google_apis;x86_64`：這台只有 x86_64，
arm64 的 image 要模擬指令集，慢到不適合當驗收迴圈。
**API 34 不等於 targetSdk 35 的實機行為**，差異只有里程碑 H 擋得到。

預設只建 `arm64` 與 `amd64` 兩個 ABI（實機與模擬器各一）。
32 位的 `arm`／`386` 要的話自己加 `WOLONG_ABIS`，但那會讓 bind 的時間加倍，
而這個專案沒有 32 位的驗收對象。

## 9. 未解

| 項目 | 現況 |
|---|---|
| 實機驗收 | ⛔ 沒有裝置。里程碑 H 保持未完成 |
| SAF 匯入 | 沒做。目前資料靠 `adb` 推進 app 內部目錄，那是驗收路徑不是玩家路徑 |
| 軍團編成與戰場 | 沒做（里程碑 D 的後半與 E）|
| 高 DPI 下的點陣字 | 原版字型是 16×15 點陣，手機上要整數放大幾倍才讀得清楚沒量過 |
| release signing | A–G 之前不談 |
