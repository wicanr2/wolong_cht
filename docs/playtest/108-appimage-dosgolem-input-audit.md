# 108 — AppImage 與 dosgolem 正常滑鼠操作抽驗

**狀態：未通過操作一致性驗收。** 畫面局部相同，不代表點擊位置與確認流程相同。

- 日期：2026-09-07；本輪只驗證，未修改遊戲、重打包或發布。
- AppImage：`dist-all/packages/wolong-remake-linux-amd64-20260907.AppImage`。
  SHA-256：`f1e5239232f4f7d88e0a4d6ab22cb2ce9bf0bc58d6e72830198a8ba30b0ba9b0`。

- 包內 `usr/bin/wlgame` SHA-256：
  `ba2c6a9b57762483462e0b0e3121d51ae4b3921ffa0078da7a40d8896b5b10c7`。

- 原版：dosv `KI.EXE`，SHA-256：
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- dosgolem 工作副本提交：`55067dbb030510ef56723885fff7b617e31d0338`。
  本專案檢查時提交：`8f0aa4adf3946db2993ba5f2c2c5db322e7febd4`。
- 本機證據根目錄：`workplace/appimage-audit-20260907/`；原版截圖由本輪 dosgolem 自行生成，未使用 DOSBox 圖片。
- 方法：兩端均從正常啟動選單經滑鼠進遊戲，不用直啟、座標注入或強制戰鬥。
  AppImage 自解後執行原封包的 `AppRun`，沒有掛入儲存庫素材；Xvfb 1280×800，視窗幾何記在 `app-trace.json`，最近鄰回到 640×400。

## 1. 已證實的差異

| 項目 | dosgolem 執行 dosv | AppImage | 判定 |
|---|---|---|---|
| **第一章日期點擊 `(300,152)`** | 196 年第一章 | **208 年第二章** | **操作缺陷，阻擋一致性驗收** |
| 勢力清單初次點第一列 `(450,128)` | 第一列反白，再點一次才進君主卡 | 點一次便進君主卡 | 確認步驟不同；尚未在本輪裁定修改方式 |
| 君主卡按確定 | 直接進大地圖 | 多一頁「確認新遊戲／開始／返回」 | 流程差異 |
| 啟動首頁 | NEW GAME 的 YES／NO | 額外主選單，再進 YES／NO；無存檔時未顯示 LOAD DATA | 已有首頁保留擴充選項的裁定，見 [106](106-launcher-parity.md) §2.2，不當成新缺陷 |
| 編成候選 | 不含曹操君主 | 預設含曹操，且候選順序不同 | 既有可切換的重製差異，見 [規格 76](../spec/76-lord-not-in-formation.md) |

第一列不是模擬器速度造成：兩邊在靜止的劇本選單點同一座標，再各走到主畫面，年份分別為 196 與 208。
dosgolem 的 `hotspots` 實測第一章日期熱區為 `x=256..375,y=144..159`，四章的日期熱區相隔 48 px。
AppImage 畫面已採原版四章版面，操作卻未對齊。現行 `cmd/wlgame/launcher.go` 的
`launcherRowRect(launcherScenario, row)` 仍取殼層列幾何，是成因線索；本輪未逆向封包執行檔，故不把原始碼線索當成封包資料流的完整證明。

可回查截圖：

- `orig/scenario.png` 與 `app-same-click/scenario.png`：點之前。
- `orig/main.png` 與 `app-same-click/main.png`：相同日期位置點擊後，第一章與第二章。
- `orig/faction-selected.png` 與 `app-same-click/lord.png`：第一次點勢力列後。
- `app/extra-confirm.png`：AppImage 額外確認頁。

**訂正前期試探：** 最初只觀察到 `(220,128)` 在原版無作用、AppImage 會選章，曾暫稱「熱區較大」。
`(300,152)` 的正對照證明是**顯示與選章命中幾何錯位**，不能只記成較寬鬆的現代化操作。
另一次試探原版 `click:300,158` 選到第四章，不拿那份輸入工具試探作產品證據；正式結論用原版已回讀熱區的 `(300,152)`，完整成功序列另存 `oracle-steps.txt`。

## 2. 相同劇本的補驗

為隔離上述缺陷，AppImage 改點 `(220,128)` 進第一章，再以正常滑鼠流程選曹操。
這是明列的操作繞路，不能反過來宣稱第一章日期熱區正常。

| 抽驗 | 結果與界線 |
|---|---|
| 第一章曹操開局、196 年 4 月 1 日、視窗關閉 | `orig/main.png` 對 `app/main.png`，**五區合計 0 / 256,000 px**；僅證此開局畫面，未對齊隱藏 RNG |
| 軍團彈出選單 | `orig/popup.png` 對 `app/popup.png`，框 `(192,56,112,56)` **0 / 6,272 px**；背景鏡頭不同，不能宣稱整張同狀態 |
| 君主卡右鍵 | 兩端皆返回勢力清單 |
| 編成候選、軍團彈出選單、系統選單右鍵 | 兩端皆關閉當前視窗 |
| 系統選單阻擋時間 | 原版再跑一千萬道指令，時鐘仍為 196 年 4 月 2 日 7 時；AppImage 等兩秒，全畫面 0 px 變動。AppImage 只取畫面，不宣稱讀到了內部子刻 |

`parity_diff.py --selftest` 的同圖與平移負對照通過；數字保存於
`main-diff.txt`、`popup-diff.txt`、`pause-diff.txt`。
第一次正常序列中的 `corps-menu.png` 實際是**編成候選**（點 x=208），名稱不代表畫面語意；
真正軍團選單是另外點 x=240 所得 `popup.png`，兩者未混算。

## 3. 重播與環境

- `oracle-steps.txt` 是原版正常啟動、選章、選君主、返回、編成候選與系統選單完整序列。
  以既有 `tools/dosgolem.sh` 接該檔內容即可重生；原版輸入目錄維持唯讀。
- `same-click-steps.json`／`same-click-trace.json` 保存錯章那一輪的 AppImage 操作。
- `app-steps.json`／`app-trace.json` 保存繞過錯章後的同劇本操作；
  `app-probe.py`／`app-run.sh` 是本輪 X11 操作與擷取程式。
- dosgolem 使用既有 `golang:1.24-bookworm`；AppImage 使用既有 `demonwinter-go`，
  皆 Docker、`--rm`、無網路、資源受限、UID/GID 1000:1000。
- 初次 AppImage 因容器無 ALSA 裝置退出，改用容器 `/tmp/.asoundrc` 的 null 裝置後乾淨重跑。
  中途試用不存在的 `-no-audio` 旗標及缺少 Pillow 的擷取失敗均屬驗證腳本問題；正式流程無該旗標，影像轉換使用既有 ImageMagick。
- 原版初次輸出目錄未建立、原版無效熱區的失敗沒有列入通過證據；正式完整序列 exit 0。
- 未 commit、push、重編 remake、變更 AppImage 或修改原始素材。

## 4. 未解與驗收範圍

本輪足以否定「整體操作一致」，但不是全遊戲完整測試。
尚未重驗存檔／重啟讀檔、正常進戰場後的戰術操作、長時間拖曳、雙擊、邊緣捲動節奏及實體滑鼠延遲。
dosgolem 的指令預算不是實機牆鐘；`int 61h` 音訊服務也不是完整聲音模擬，不能由本輪推定聽感與即時速度一致。
地圖鏡頭差異與游標顯示仍有既有邊界（[規格 149](../spec/149-march-target-map-picker.md)、[規格 154](../spec/154-mouse-cursor.md)）。
下一個最小工作是依既有劇本視窗規格審查並修正啟動選章命中範圍，再以重建後 AppImage 重跑第一章日期點擊；本輪不以原始碼截圖代替封包驗收。

## 5. 後續修正（2026-09-07）

使用者授權修正後，兩項桌面輸入缺陷已在 `v.1.0.1-20260907` 的 AppImage 重播通過，
見 [109](109-desktop-launcher-fix.md)。上面的失敗結論保留給原 `20260907.AppImage`，
不套用到修正版；額外確認頁仍待取捨。
