# 04 — 2026-08-22 三平台重新交付（含 Android 完整版）

**狀態：已交付並驗過。** `dist-all` 是一致的 `wolong-remake-20260822` 批次：
五個桌面包與 Android APK 都由這一輪的原始碼建出來，雜湊、deny-list 與
Linux GUI smoke 都過。

- 日期：2026-08-22
- 工具：`tools/release_all.sh 20260822` ＋ `tools/android_build.sh`
  ＋ `tools/release_smoke.sh 20260822`
  ＋ `tools/py.sh tools/release_all_fs.py refresh`
- 前一批：`wolong-remake-20260821`（[`03`](03-three-platform-20260821.md)，已被本批取代）

## 1. 產物

| 檔 | 大小 (B) | SHA-256 |
|---|---:|---|
| `wolong-remake-linux-amd64-20260822.AppImage` | 9,942,208 | `5755504dc09bf3915966abf7caa14a92f7ddab123ec3dbdb1f6f5484296dc176` |
| `wolong-remake-linux-amd64-20260822.tar.gz` | 9,593,572 | `79157e8ac0118dd82d87fd970d5683a1e8e0c76fb301bf9cd7134a7f54387c51` |
| `wolong-remake-linux-arm64-tools-20260822.tar.gz` | 2,156,200 | `a598c88e1666f2c2f83a09cd7b765996ca1f4db5b2e1772a6017629bfe6ee649` |
| `wolong-remake-macos-universal-20260822.tar.gz` | 18,131,358 | `837c0b8d0b6136f495bb8e9822977df3a2540d458f23c0fd3af662f9a563c5a6` |
| `wolong-remake-windows-amd64-20260822.tar.gz` | 9,548,723 | `9734da94ad7a98dfc59a9478f309d2f2f64b51a29e61bc209040e34bd3984de0` |
| `wolong-remake-android-debug-20260822.apk` | 22,718,998 | `2797a4edd321305fb1b3990ac3e83be0db467f5d897d7e3bf6b4a6dc55bcfdec` |

Android 版仍列在 `experimental/android/`，界線是**簽章與驗收**不是功能：
它是完整的遊戲，規則層與桌面版是同一份程式碼，但只有 debug 簽章，
而且只在 Docker 的模擬器上驗過。

## 2. 這一批比 `20260821` 多了什麼

| 項目 | 出處 |
|---|---|
| 結局的節拍照原版重做：每一幕多了 11 秒看圖時間，整段 1 分 47 秒 → 3 分 21 秒 | [`../spec/67`](../spec/67-ending-playback.md) §8 |
| 戰場的右鍵熱區層：熱區 `0x1D` 右鍵提前收掉門強度條 | [`../spec/32`](../spec/32-gate-strength-bar.md) |
| 點縮小地圖捲鏡頭（熱區 `0x16`）、22 勢力選擇視窗（熱區 `0x17`） | [`../re/71`](../re/71-strategy-hotspot-dispatch.md)、[`../spec/35`](../spec/35-strategy-minimap.md) §2.5 |
| M7 的校訂文字實跑抽樣：18 則全部在框內 | [`../playtest/41`](../playtest/41-m7-corrected-text-on-screen.md) |
| 九個文件目錄逐列稽核，未解總表 570 → 431 列 | [`../re/43`](../re/43-open-questions.md) |

## 3. ⭐ APK 沒被複製進交付目錄——`refresh` 少了一步

重建 APK 之後只跑 `refresh`，交付目錄的 `experimental/android/` 會**只剩
`README.md`**，而 `release-manifest.json` 照樣寫著上一批的檔名。
`SHA256SUMS.txt` 也照樣產得出來，因為它只雜湊真的存在的檔——
**少一個檔不會讓任何一步失敗**。

成因是 APK 的複製寫在 `stage()` 裡，而 `refresh` 不呼叫 `stage`；
`android_experimental` 這個欄位又寫在 `finalise()`，`refresh` 也不碰。
Android 是另一條管線（`tools/android_build.sh`）建的，
於是「重建 APK ＋ refresh」這條最自然的路徑正好兩步都跳過。

修法兩層，都在 `tools/release_all_fs.py`：

1. `sync_android()` 抽成一支，`stage()` 與 `refresh()` 都呼叫它——
   複製當批 APK、清掉別批的、重寫 `README.md`，回傳發行檔名。
   `refresh()` 拿回傳值回填 manifest。
2. `verify_manifest_paths()`：manifest 裡列到的每一個路徑都要真的存在，
   不存在就讓 `finalise` 與 `refresh` 當場失敗。

第二層才是真正的閘。第一層修的是這一次的成因，第二層擋的是
**下一次任何一個檔案漏掉**——沉默的成功比失敗難發現。

## 4. 驗收

| 項目 | 結果 |
|---|---|
| AppImage 啟動 ＋ 大地圖 | ✅ `verification/appimage-smoke-20260822.png`（196年4月2日）|
| AppImage 結局過場 | ✅ `verification/appimage-ending-20260822.png` |
| Linux tar 解開直接跑 | ✅ `verification/linux-tar-smoke-20260822.png`，與 AppImage 那張逐位元組相同 |
| 交叉建置檔頭 | ✅ ELF x86-64／PE32+ x86-64／Mach-O x86_64／Mach-O arm64 |
| `.so` 的 16 KB 對齊 | ✅ ELF LOAD 段與 zip 內位移都過（[`03`](03-three-platform-20260821.md) §4）|
| deny-list | ✅ 沒有原版資產 |
| 雜湊 | ✅ `sha256sum -c dist-all/SHA256SUMS.txt` 二十筆全部相符 |
| manifest 路徑 | ✅ 六個交付路徑逐一存在（§3 的新檢查）|
| Windows／macOS 原生 GUI | ❌ **仍未實機驗收**，是 M8 唯一的閘 |

⚠ **結局那一幕的截圖上沒有文字**，那不是缺陷：容器裡沒有倚天字型，
`textdraw` 不可用時整段文字就不畫。字型與原版資料一律由玩家自備。

## 5. 未解

| 項目 | 現況 |
|---|---|
| Windows／macOS 原生 GUI 實機驗收 | 沒有硬體，Docker 代不了 |
| Android 實機驗收 | 只有模擬器；16 KB page size 的裝置也還沒實際載過 |
| Android 正式簽章 | 還沒決定 keystore 怎麼保管 |
