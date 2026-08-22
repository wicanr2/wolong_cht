# 03 — 2026-08-21 三平台重新交付（含 Android 版）

**狀態：已交付並驗過，已被 `wolong-remake-20260822` 那一批取代。**
當時的 `dist-all` 是一致的 `wolong-remake-20260821` 批次：
五個桌面包與 Android APK 都由這一輪的原始碼建出來，雜湊、deny-list 與
Linux GUI smoke 都過。

- 日期：2026-08-21
- 工具：`tools/release_all.sh 20260821` ＋ `tools/release_smoke.sh 20260821`
  ＋ `tools/py.sh tools/release_all_fs.py refresh`
- 後續：`wolong-remake-20260822`（[`04`](04-three-platform-20260822.md)）
- 前一批：`wolong-remake-20260820`（[`02`](02-three-platform-20260820.md)，已被本批取代）

## 1. 產物

| 檔 | 大小 (B) | SHA-256 |
|---|---:|---|
| `wolong-remake-linux-amd64-20260821.AppImage` | 9,942,208 | `3c17162b8860b49fc9b6773e690a80fa596e7d9a84f2c330d11b88f945b3cb18` |
| `wolong-remake-linux-amd64-20260821.tar.gz` | 9,591,200 | `e8a0cf323045f461f41b1c9f7e4cbb8e471fa7a15830ccd4b0c170e850bc1321` |
| `wolong-remake-linux-arm64-tools-20260821.tar.gz` | 2,156,137 | `77bd23c4f455c9cb35fdc32827fdf2cd979c6759990c9743913bae75456abe74` |
| `wolong-remake-macos-universal-20260821.tar.gz` | 18,123,650 | `81aef61c038bd7b6d6dcca3df0d12e268faec7a0914de358fe782f91230ea086` |
| `wolong-remake-windows-amd64-20260821.tar.gz` | 9,548,428 | `3d12ee3f18e76d513c2cba3be6887384d0e3f32fdc316c41d88777e3a31739ed` |
| `wolong-remake-android-debug-20260821.apk` | 22,718,566 | `c6b083b419968a78cb79ab9ade32d596abea690e3686e6465f68944777dc0f4a` |

Android 版列在 `experimental/android/`，界線是**簽章與驗收**不是功能：
它是完整的遊戲，規則層與桌面版是同一份程式碼，但只有 debug 簽章，
而且只在 Docker 的模擬器上驗過。

## 2. 這一批比 `20260820` 多了什麼

| 項目 | 出處 |
|---|---|
| Android 版從觸控原型變成完整遊戲：四個入口、戰場、存檔、事件訊息、三種擋世界的決定 | [`../mobile/android-plan.md`](../mobile/android-plan.md) |
| Android 的按鈕與面板改用原版的底色與外框 | [`../spec/70`](../spec/70-phone-chrome.md) |
| 存檔寫的是**當下**的世界，不是載入時的那一份 | `internal/savefile`、[`../../CONTEXT.md`](../../CONTEXT.md) §6 |
| 戰術側欄六個命令的熱區從 128 px 收回原版的 48 px | [`../re/60`](../re/60-tactical-sidebar.md) §6 |
| 桌面與手機共用 `internal/ui/isoview`、`internal/battlesetup` 等七個套件 | 同上 |
| `.so` 改成 16 KB 對齊，Android 15 的 16 KB 裝置載得起來 | §4 |
| 推廣主片的大地圖與兩場戰鬥改成實跑逐幀錄製 | [`../spec/71`](../spec/71-promo-live-capture.md) |

## 3. 驗收

| 項目 | 結果 |
|---|---|
| AppImage 啟動 ＋ 大地圖 | ✅ `verification/appimage-smoke-20260821.png`（196年4月2日）|
| AppImage 結局過場 | ✅ `verification/appimage-ending-20260821.png` |
| Linux tar 解開直接跑 | ✅ `verification/linux-tar-smoke-20260821.png` |
| 交叉建置檔頭 | ✅ ELF x86-64／PE32+ x86-64／Mach-O x86_64／Mach-O arm64 |
| deny-list | ✅ 沒有原版資產 |
| 雜湊 | ✅ `sha256sum -c dist-all/SHA256SUMS.txt` 二十筆全部相符 |
| Windows／macOS 原生 GUI | ❌ **仍未實機驗收**，是 M8 唯一的閘 |

⚠ **結局那一幕的截圖上沒有文字**，那不是缺陷：容器裡沒有倚天字型，
`textdraw` 不可用時整段文字就不畫。字型與原版資料一律由玩家自備。

## 4. `.so` 的 16 KB 對齊

Android 15 起有 16 KB page size 的裝置。Go 產出的 `libgojni.so` 預設是
4 KB（LOAD 段 `align=0x1000`），那種 `.so` 在 16 KB 的機器上**載不起來**——
症狀是一啟動就掛，而 4 KB 的機器上完全正常。

⚠ **`zipalign -P 16` 驗不到這一層。** 它看的是 zip 裡的檔案位移，
`.so` 內部的 ELF 段對齊是另一回事；兩者都要對。

做法：`ebitenmobile bind` 帶
`-ldflags "-extldflags=-Wl,-z,max-page-size=16384"`，
並在 `tools/android_build.sh` 建完之後用 `readelf` 逐段檢查，
不是 `0x4000` 就讓建置失敗。**旗標打錯字、工具鏈換版都會讓對齊悄悄掉回去，
而那個世界與正常世界在 4 KB 的機器上長得一模一樣。**

## 5. 這一批當時的未解

> **這一批已被 [`04`](04-three-platform-20260822.md) 取代**，現行的閘看那一份。
> 三條當時都成立，也全部原封不動延續到 `20260822`，登記在 `04` §5，
> 不在這裡重複數一次。

| 項目 | 現況 |
|---|---|
| ~~Windows／macOS 原生 GUI 實機驗收~~ | 見 [`04`](04-three-platform-20260822.md) |
| ~~Android 實機驗收~~ | 見 [`04`](04-three-platform-20260822.md) |
| ~~Android 正式簽章~~ | 見 [`04`](04-three-platform-20260822.md) |

§4 的 16 KB 對齊做法**仍然是現行做法**，`20260822` 的 APK 也照它建、照它驗。
