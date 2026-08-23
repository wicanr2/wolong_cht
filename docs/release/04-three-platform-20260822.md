# 04 — 2026-08-22 四平台完整版（內含遊戲檔案）

**狀態：已交付並驗過，已被 `wolong-remake-20260823` 那一批取代。**
⛔ 這一批內含原版資產，不可外流。
`dist-all` 是一致的 `wolong-remake-20260822` 批次：四個平台的包裡都有
松崗 DOS/V 的 69 個原版檔與倚天點陣字，解開或裝上去就能玩。

- 日期：2026-08-22
- 使用者裁定：`dist-all` 就是四平台完整版（規格 [`../spec/72`](../spec/72-bundled-game-data.md)）
- 工具：`tools/release_all.sh 20260822` ＋ `WOLONG_BUNDLE_DATA=1 tools/android_build.sh`
  ＋ `tools/release_smoke.sh 20260822` ＋ `WOLONG_BUNDLED=1 tools/android_smoke.sh`
- 後續：`wolong-remake-20260823`（[`05`](05-full-20260823.md)）
- 前一批：`wolong-remake-20260821`（[`03`](03-three-platform-20260821.md)，已被本批取代）

## 1. ⛔ 散布界線

這一批**不可外流**——不上傳、不轉傳、不放進會被同步到公開位置的目錄。
`CLAUDE.md` §9 的 `[HARD]` 規則寫著「不得 commit 或**打包**任何原版
`.EXE/.COM/.DAT/…`，也不得散布倚天字型」；使用者在我指出衝突後裁定
打包這半邊要做，**commit 那半邊沒有例外**：`dist-all/` 在 `.gitignore`，
deny-list 掃的是 git 追蹤的檔案，兩者都碰不到這裡。

保存專案的定位沒有改變：公開產出仍然只有引擎程式碼與譯文校訂紀錄。

**怎麼分辨手上這一批是哪一種**（⚠ 看檔名分不出來，兩種的名字一模一樣）：

| 判準 | 完整版 | 可散布 |
|---|---|---|
| `DO-NOT-DISTRIBUTE.md` | 存在 | 不存在 |
| manifest 的 `distributable` | `false` | `true` |
| manifest 的 `original_assets_included` | `true` | `false` |
| 包裡有沒有 `gamedata/` | 有 | 沒有 |

要一份可散布的：`WOLONG_BUNDLE_DATA=0 tools/release_all.sh <YYYYMMDD>`。

## 2. 產物

| 檔 | 大小 (B) | SHA-256 |
|---|---:|---|
| `wolong-remake-linux-amd64-20260822.AppImage` | 12,133,568 | `2fd66558e92f2e1bc09a16c8fbd653362b4bc952f252a576d9e592ec420487dc` |
| `wolong-remake-linux-amd64-20260822.tar.gz` | 11,755,451 | `1a38f835a5e840b4a5a26118a5d5242d9c18b730ae8d3ecc7848adaeb0785c45` |
| `wolong-remake-windows-amd64-20260822.tar.gz` | 11,710,024 | `c243bd86f4e91e0ef9f5907b7216d660c80c349391ebee3f73df5bdb65e5cb64` |
| `wolong-remake-macos-universal-20260822.tar.gz` | 20,296,128 | `fc8a6959b56ea45ee310e13f64c3d347c028dced8a2706fbfc17398517f774c2` |
| `wolong-remake-android-debug-20260822.apk` | 25,136,585 | `33e2c55f933553b409bd39d3dffa7669cda339a0b379bc2416f9bf59fbc7d076` |
| `wolong-remake-linux-arm64-tools-20260822.tar.gz` | 2,156,295 | `5aff2ea654f7a7dac7f4e5bb7acb1fe3b29b8be48bd598690a6b5884f1b9de02` |

⭐ **Android 移到 `packages/`**，與另外三個平台並列。它先前在
`experimental/android/`，界線是簽章與驗收——那是**驗收狀態**，
寫在說明裡就夠，用目錄分級表達會讓人以為它功能不完整。
它仍然只有 debug 簽章、也還沒實機驗收。

⚠ **`linux-arm64-tools` 不含資料**：它只有 `wlsim`／`wlshot`，不是完整遊戲。

## 3. 資料在包裡的位置

| 平台 | 版面 |
|---|---|
| Linux／Windows tar | 包根的 `gamedata/`（69）＋ `fonts/`（3）|
| macOS tar | 同上，**兩個架構共用一份**（執行檔在 `darwin-<arch>/`）|
| AppImage | `usr/share/wolong-remake/gamedata`、`.../fonts` |
| APK | `assets/gamedata/{orig,eten}`，第一次啟動由 `ImportActivity` 解到私有目錄 |

`wlgame` 的 `-orig`／`-font` 預設是 repo 相對路徑，解開的包裡不成立。
新的 `resolveDataDir` 照 `bundledTalkCorrectionsPath` 的形狀補上退路：
**明講的旗標一律優先**，沒講才依序找 repo 路徑、執行檔旁、AppImage 版面
（[`../spec/72`](../spec/72-bundled-game-data.md) §3）。

## 4. 驗收

| 項目 | 結果 |
|---|---|
| AppImage 啟動 ＋ 大地圖 | ✅ `verification/appimage-smoke-20260822.png` |
| AppImage 結局過場 | ✅ `verification/appimage-ending-20260822.png` |
| Linux tar 解開直接跑 | ✅ `verification/linux-tar-smoke-20260822.png` |
| **不帶任何資料旗標** | ✅ `verification/bundled-nodflags-20260822.png`，而且**沒有「載不到字型」那行警告** |
| **Android 乾淨安裝、不推資料** | ✅ 自己解出 69 個檔並轉進遊戲，`verification/android-bundled-20260822.png` |
| APK 的資料與桌面一致 | ✅ 指紋 1／60／120 三幀與桌面**逐位元組相同** |
| 四平台同一個檔案集 | ✅ 各 69 ＋ 3 |
| 交叉建置檔頭 | ✅ ELF x86-64／PE32+ x86-64／Mach-O x86_64／Mach-O arm64 |
| `.so` 的 16 KB 對齊 | ✅ ELF LOAD 段與 zip 內位移都過 |
| manifest 路徑 | ✅ 六個交付路徑逐一存在 |
| 雜湊 | ✅ `sha256sum -c dist-all/SHA256SUMS.txt` 二十二筆相符 |
| Windows／macOS 原生 GUI | ❌ **仍未實機驗收**，是 M8 唯一的閘 |

⚠ **結局那一幕的截圖上沒有文字**：容器裡跑的是 AppImage 內建的字型路徑，
而那一張是在 `-orig` 明講的舊路徑下拍的。文字本身在 §4 第四列那張裡是正常的。

## 5. ⭐ 三個安靜的失敗

**一、`refresh` 不會複製 APK。** 重建 APK 之後只跑 `refresh`，
`experimental/android/` 會只剩 `README.md`，而 manifest 照樣寫著上一批的檔名、
`SHA256SUMS.txt` 也照樣產得出來——**少一個檔不會讓任何一步失敗**。
修法兩層：`sync_android()` 給 `stage` 與 `refresh` 共用，
以及 `verify_manifest_paths()` 讓 manifest 指到不存在的路徑當場失敗。
第二層才是閘。

**二、模擬器截圖睡著時會拍到全黑，而那是一張合法的 PNG。**
指紋那一段跑 24 秒，中途螢幕關掉並上鎖；叫醒之後停在 Android 桌面，
拍到的圖不黑、大小也正常，**只是拍到別的東西**。
現在：開機就把 `screen_off_timeout` 拉到 30 分鐘、截圖前查 `mCurrentFocus`
是不是我們的 activity、拍完檢查檔案大小，全黑就重試。

**三、我自己寫的醒沒醒檢查是壞的。**
`set -o pipefail` 配 `grep -q`——命中就提早結束，上游的 `tr` 吃 SIGPIPE，
**整條 pipeline 回非零，明明命中卻被判成沒命中**。
實測 `mWakefulness=Awake`、`Display State=ON`，而檢查照樣報「叫不醒」。
現在 grep 在裝置上跑，不在 host 這端接管線。

## 6. 這一批當時的未解

> **這一批已被 [`05`](05-full-20260823.md) 取代**，現行的閘看那一份。
> 四條當時都成立，也全部原封不動延續下去，登記在 `05` §5，
> 不在這裡重複數一次。

| 項目 | 現況 |
|---|---|
| ~~Windows／macOS 原生 GUI 實機驗收~~ | 見 [`05`](05-full-20260823.md) |
| ~~Android 實機驗收~~ | 見 [`05`](05-full-20260823.md) |
| ~~Android 正式簽章~~ | 見 [`05`](05-full-20260823.md) |
| ~~完整版在低容量裝置上~~ | 見 [`05`](05-full-20260823.md) |

<!-- 缺口：無 -->
