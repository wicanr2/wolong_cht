# 72 — 內含遊戲檔案的四平台完整版

**狀態：CONFORMED。四平台的完整包都內含原版資料與倚天字型並實跑驗過；
`dist-all` 因此從「可散布」翻成「私人批次」。**

- 日期：2026-08-22
- 出處：使用者裁定 2026-08-22（口頭），我先指出與 `CLAUDE.md` §9 `[HARD]` 的衝突後獲重申
- 推論等級：**決定**，不是從原版推出來的

## 1. 這一條與 `[HARD]` 規則的關係

`CLAUDE.md` §9 寫著：

> 不得 commit 或**打包**任何原版 `.EXE/.COM/.DAT/.MAP/.MDL/.SCH/.MCH/.BRG/.O`，
> 也不得散布倚天字型。

**這份規格明確違反那一條的「打包」半邊，是使用者裁定的例外。**
「commit」那半邊不受影響——`dist-all/` 在 `.gitignore:54`，
`tools/denylist.py` 掃的是 git 追蹤的檔案，兩者都碰不到。

⚠ **代價要寫清楚**：`dist-all` 之後**不可外流**。它原本的定位是
「可交付根目錄」，`release-manifest.json` 有 `original_assets_included: false`
這個欄位，`README.md` 明寫「桌面包不含任何原版執行檔、資料、美術、音樂、字型」。
這些敘述全部要翻面，否則會有人照著舊敘述把它傳出去。

**保存專案的界線沒有改變**：公開產出仍然只有引擎程式碼與校訂紀錄。
改的只是「本機這一份 `dist-all` 是給自己玩的」。

## 2. 交付結構

```
dist-all/
  DO-NOT-DISTRIBUTE.md          ← 新增，第一眼就看得到
  packages/
    wolong-remake-linux-amd64-<stamp>.tar.gz      ← 含 gamedata/ ＋ fonts/
    wolong-remake-linux-amd64-<stamp>.AppImage    ← 資料在 AppImage 內
    wolong-remake-windows-amd64-<stamp>.tar.gz    ← 含 gamedata/ ＋ fonts/
    wolong-remake-macos-universal-<stamp>.tar.gz  ← 含 gamedata/ ＋ fonts/
    wolong-remake-android-<stamp>.apk             ← 資料在 APK 的 assets/
    wolong-remake-linux-arm64-tools-<stamp>.tar.gz ← 工具伴隨包，**不含資料**
```

⭐ **Android 從 `experimental/` 提到 `packages/`**，與另外三個平台對等——
它現在也是「裝上去就能玩」。debug 簽章與沒有實機驗收這兩件事仍然成立，
但那是**驗收狀態**，不是「能不能玩」，寫在 README 而不是用目錄分級表達。

⚠ **`linux-arm64-tools` 仍然不含資料**：它只有 `wlsim`／`wlshot`，
不是完整遊戲，塞 4.4 MB 進去只是讓包變大。

| 包內路徑 | 內容 | 來源 |
|---|---|---|
| `gamedata/` | 69 個原版檔 | `workplace/orig/dosv/` |
| `fonts/` | `STDFONT.15`、`ASCFONT.15`、`SPCFONT.15` | `workplace/eten/` |
| `translations/` | `corrections.json` | 不變 |

## 3. 程式怎麼找到資料

`wlgame` 的兩個旗標預設值是 **repo 相對路徑**
（`-orig workplace/orig/dosv`、`-font workplace/eten`），
解開的包裡不成立——使用者跑 `./wlgame` 會得到「載不到字型，中文變方框」。

做法**照抄既有的 `bundledTalkCorrectionsPath()`**（`cmd/wlgame/main.go:1298`）：
使用者明講就用他講的，沒講才依序找。

```
resolveDataDir(值, 預設值, 包內名稱):
    值 != 預設值            → 值            # 使用者明講，一律尊重
    isDir(預設值)           → 預設值        # repo 內開發，行為完全不變
    <exe 目錄>/包內名稱      → 它            # tar/zip 解開後
    <exe 目錄>/../包內名稱   → 它            # AppImage 的 usr/bin/
    <exe 目錄>/../share/wolong-remake/包內名稱 → 它
    否則                    → 預設值        # 讓載入器噴可診斷的錯，不要靜默
```

`-orig` 用 `gamedata`，`-font` 用 `fonts`。

⭐ **三條性質**，缺一條這個改動就不安全：

1. **不覆蓋明講的旗標**——對拍與驗收全部明講 `-orig`，一個字都不受影響。
2. **repo 內行為不變**——`workplace/orig/dosv` 存在時第二條就命中了。
3. **找不到時回預設值**，由既有的載入器報錯。**不要靜默跳過**——
   靜默的成功比失敗難發現（`CLAUDE.md` §7 第 21 條）。

⚠ **只改 `wlgame`。** `wlsim`／`wlshot`／`wlview` 是工具，它們的呼叫端
一律明講路徑；跟著改只會擴大這次的受影響面。

## 4. Android 怎麼內嵌

現況是 SAF 匯入：`ImportActivity` 把使用者選的資料夾複製進
`getFilesDir()/orig` 與 `/eten`，`MainActivity.dataRoot()` 再交給 Go
（`Wolongmobile.setDataRoot`）。**Go 端只認檔案系統路徑。**

所以內嵌**不動 Go 一個字**：資料放進 APK 的 `assets/`，
`ImportActivity` 在沒有資料時先試「從 assets 解開」，成功就直接開遊戲。

```
ImportActivity.onCreate:
    hasData()          → startGame()
    unpackBundled()    → startGame()      # 新增這一步
    否則                → 顯示 SAF 選資料夾（原本的畫面）
```

| 決定 | 理由 |
|---|---|
| 解開到 `getFilesDir()`，不直接從 assets 讀 | Go 端用 `os.ReadFile`。要讓它讀 assets 得走 `AssetManager`，那會改到 `internal/assets` 的載入路徑——**平台細節不該滲進解碼層** |
| 多花約 4.8 MB 儲存 | APK 裡一份、解開後一份。換到的是零 Go 改動 |
| SAF 那條路**留著** | 沒有內嵌資料的建置（給別人的版本）仍然要能用。刪掉它等於讓這個 APK 只有一種建法 |
| 用 `SINARIO.DAT` 判定齊不齊 | 沿用 `ImportActivity.REQUIRED`，不新增第二套判準 |

⚠ **`assets/` 有沒有資料是建置時的事**，`tools/android_build.sh` 用
`WOLONG_BUNDLE_DATA=1` 開關；預設**不**內嵌，避免有人不小心建出一個
帶原版資料的 APK 又送出去。

## 5. remake 實作

| 項目 | 位置 |
|---|---|
| 資料目錄解析 | `cmd/wlgame/main.go` 的 `resolveDataDir` |
| Android 解包 | `android/app/src/main/java/com/wicanr2/wolong/ImportActivity.java` |
| 打包 | `tools/release_all_fs.py` 的 `stage`／`appdir`／`finalise` |
| Android 建置開關 | `tools/android_build.sh` 的 `WOLONG_BUNDLE_DATA` |
| 警告文件 | `packaging/release/DO-NOT-DISTRIBUTE.md` |
| 差異 | 這整份都是 remake 差異——原版沒有「發行包」這個概念 |

## 6. 驗證（2026-08-22 全部跑過）

| 方式 | 結果 |
|---|---|
| 單元測試 | ✅ `TestResolveDataDir`／`TestBundledDirNames`（`cmd/wlgame`）：明講的旗標、repo 路徑、exe 旁、全都找不到，四條各一個子測試 |
| 解開就能跑 | ✅ tar 解到空目錄，`./wlgame` 不帶任何資料旗標，**沒有「載不到字型」那行警告**（接進 `tools/release_smoke.sh`）|
| AppImage | ✅ `--appimage-extract-and-run` 不帶旗標，中文正常（最近一次是 `20260824` 批次）|
| Android | ✅ 模擬器乾淨安裝、**一個 byte 都沒推**，`ImportActivity` 自己解開 69 個檔並轉進 `MainActivity`（最近一次是 `20260824` 批次的 APK）|

⚠ **這裡刻意不寫 `dist-all/verification/…` 的路徑。** 那個目錄
**每次重打發行都會被 `promote` 整個換掉**（`tools/release_smoke.sh` 的檔頭
寫過這件事），所以拿它當長期證據的引用，下一批一產出就指向不存在的檔。
⭐ **會被清掉的路徑不能當證據**——要留就寫批次號，重跑的指令寫在 §6。
| APK 的資料與桌面一致 | ✅ 指紋 frame 1／60／120 ＝ `2b58e7b5…`／`5b3585cf…`／`36eb02d3…`，**與桌面逐位元組相同**。⭐ 資料只要差一個 byte，世界狀態就會分岔——這比逐檔比對便宜也更強 |
| 四個平台同一個檔案集 | ✅ 桌面三包各 69 ＋ 3，APK 也是 69 ＋ 3 |
| 沒有內嵌時不退步 | ✅ `WOLONG_BUNDLE_DATA` 沒開時 `tools/android_build.sh` 會驗 APK 裡**沒有**原版資產，有就讓建置失敗 |
| 雜湊 | ✅ `sha256sum -c dist-all/SHA256SUMS.txt` 二十二筆相符 |

### 6.1 兩個踩過的坑

**一、複製原版資料會把「唯讀」一起複製過去。** `workplace/orig/` 整個是
`chmod a-w`，`copytree` 連目錄的 mode 一起帶——staging 目錄自己變成不可寫，
收尾的 `shutil.rmtree` 在 `unlink` 吃 `PermissionError`
（unlink 要的是**父目錄**的寫入權，不是檔案的）。
現在目錄放寬到 0755、檔案維持 0444。

**二、兩條管線挑了不同的檔案集。** `copytree` 收隱藏項，
於是桌面包多了 `workplace/orig/dosv/.jsdos/`（js-dos 的 dosbox.conf），
73 個檔；而 `tools/android_build.sh` 的 `cp "$ORIG_DIR"/*` 不收 dotfile，
69 個。⚠ **這種差異不會當場壞掉**，但之後任何「兩邊對不起來」的問題
都會先被誤判成別的原因。現在兩邊都只收最上層、不收隱藏項。

## 7. 未解

| 項目 | 現況 |
|---|---|
| Windows／macOS 上「解開就能跑」 | ⛔ 沒有那兩個平台的機器。包內版面驗過（`gamedata/`、`fonts/` 位置正確），但 `resolveDataDir` 在那兩個 OS 上沒實跑過 |
| APK 內嵌後的實機驗收 | ⛔ 沒有裝置。模擬器驗到了解包與指紋，驗不到真實儲存空間與 DPI |
| 25 MB 的 APK 在低容量裝置上 | 解包後 app 私有目錄再佔約 4.8 MB，總共約 30 MB。**沒有量過安裝失敗的門檻** |
