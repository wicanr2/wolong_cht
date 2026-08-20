# 02 — 2026-08-20 三平台重新交付

**狀態：已交付並驗過。** `dist-all` 是一致的 `wolong-remake-20260820` 批次：
五個包全部由同一次建置產出，雜湊、deny-list 與 Linux GUI smoke 都過。

- 日期：2026-08-20
- 工具：`tools/release_all.sh 20260820` ＋ `tools/release_smoke.sh 20260820`
- 前一批：`wolong-remake-20260812`（已被本批取代）

## 1. 產物

| 檔 | 大小 (B) | SHA-256 |
|---|---:|---|
| `wolong-remake-linux-amd64-20260820.AppImage` | 9,872,576 | `4b122c854e50f2afa1e098f38487ac6007080949503ce6b87b13d9987f17b792` |
| `wolong-remake-linux-amd64-20260820.tar.gz` | 9,523,162 | `cd7dfa984548ec6ba30e11c8457569e818941dde2395a2be5ef81e138403f2f5` |
| `wolong-remake-windows-amd64-20260820.tar.gz` | 9,481,578 | `defef9c063b565fc1ea3579796fd7d9139f9e8bd3b55c5c9aee7d0383f7297d4` |
| `wolong-remake-macos-universal-20260820.tar.gz` | 18,003,888 | `a8211a5b1cd06cdc8adaf31808bffdcc1f452772771cde1bb897e159c3735da5` |
| `wolong-remake-linux-arm64-tools-20260820.tar.gz` | 2,105,526 | `d9921aa096ec020a3099304cfc8a18134f68419d7ddd5bf77ba9d975e3503c59` |

Android 觸控原型是**實驗性附件**，不計入三平台完整發行。
⭐ 它的檔名帶的是**它自己的建置日期**（`…-20260811.apk`）而不是發行日期——
`release_all.sh` 只是把既有的 `app-debug.apk` 複製過來，沒有重編；
掛上發行日期會讓檔名宣稱一個沒發生過的建置，而雜湊會與上一批一模一樣。

## 2. 這一批比 `20260812` 多了什麼

| 項目 | 出處 |
|---|---|
| 結局過場：十二幕 ＋ 逐字打出來的十行結尾文字 ＋ 十七階淡入淡出 | [`../spec/67`](../spec/67-ending-playback.md) |
| 過場圖 `OPEN_S*`／`END_S*` 的解碼器 | [`../formats/09`](../formats/09-cutscene-images.md) |
| 倒地動畫（四幀、三個兵種組）| [`../spec/68`](../spec/68-death-animation.md) |
| 打壞的城壁與門會在畫面上換掉 | [`../spec/66`](../spec/66-broken-walls-repaint.md) |
| 兵的開場體力 ＝ 軍團士氣、被換位的兵這一幀不動、挨打硬直、退卻算生還 | [`../spec/61`](../spec/61-soldier-initial-hp-from-morale.md)–[`65`](../spec/65-retreated-soldiers-survive.md) |
| 遷都報告要有外交官才收得到 | [`../spec/64`](../spec/64-capital-relocation-report.md) |
| 六階外交的門檻收斂成一份實作 | [`../mechanics/50`](../mechanics/50-diplomacy.md) |
| 戰場逐區對拍：九區裡六區逐像素相同、`field` 0.17% | [`../playtest/40`](../playtest/40-tactical-parity.md) |

## 3. 驗收

| 項目 | 結果 |
|---|---|
| AppImage 啟動 ＋ 大地圖 | ✅ `verification/appimage-smoke-20260820.png`（196年4月2日）|
| AppImage 結局過場 | ✅ `verification/appimage-ending-20260820.png`（第一幕的月下亭子）|
| Linux tar 解開直接跑 | ✅ `verification/linux-tar-smoke-20260820.png` |
| 交叉建置檔頭 | ✅ ELF x86-64／PE32+ x86-64／Mach-O x86_64／Mach-O arm64 |
| deny-list | ✅ 掃 20 個檔、比對 120 個原版檔的內容雜湊，沒有原版資產 |
| 雜湊 | ✅ `sha256sum -c dist-all/SHA256SUMS.txt` 全部相符 |

⚠ **結局那一幕的截圖上沒有文字**，那不是缺陷：容器裡沒有倚天字型，
`textdraw` 不可用時整段文字就不畫。跑起來會先印「載不到倚天全形字型」——
**字型與原版資料一律由玩家自備**。

## 4. 這一輪修掉的三個發行流程問題

| 問題 | 症狀 | 修法 |
|---|---|---|
| `tools/release_all.sh` 沒有執行權限 | 直接 `exit 126`，什麼都沒動 | `chmod +x` |
| 版本字串硬寫在十幾處 | 換版本要改十幾個字串，漏改一處就產出「名字對不上內容」的檔案 | 只在 `release_all.sh` 定一次 `STAMP`，其餘由 `RELEASE_STAMP` 推導 |
| 推廣片找不到來源 | `stage` 只看 `dist/promo/`，而四支影片只剩 `dist-all/promo/`；**要編完三平台才會發現** | `promo_source()` 加退路（`dist/promo` 優先、退回 `dist-all/promo`）|

⭐ 第三項的教訓是**驗證的順序**：整批重建把「輸入齊不齊」放在最貴的步驟之後，
於是一個複製檔案的錯誤要花十幾分鐘才浮現。

順帶把兩個建置容器改用共用的 `wl-gobuild` 快取（原本每次 `/tmp/gocache` 冷啟動），
第二次重跑從十幾分鐘降到約一分半。

## 5. 未解

| 項目 | 現況 |
|---|---|
| Windows／macOS 的實機驗收 | 仍未做（M8 唯一的閘）。這一批只有 Linux 有 GUI smoke，另兩個平台只驗了檔頭 |
| `verification/` 的截圖不在管線裡 | `promote` 每次都會清掉，要另外跑 `tools/release_smoke.sh` 再 `release_all_fs.py refresh` |
| Android 原型沒有重編 | 內容仍是 2026-08-11 那次（檔名已如實反映）|
