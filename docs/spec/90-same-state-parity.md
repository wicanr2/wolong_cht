# 90 — 同狀態畫面對拍

**狀態：CONFORMED。管線接起來了，也對原版跑過一輪
（[`../playtest/37`](../playtest/37-main-screen-parity.md)）。**
⭐ 開局用不到存檔——兩邊天然同狀態。

- 日期：2026-08-15
- 出處：[`docs/re/20`](../re/20-ida-re-coverage-audit.md) §5 第四道 gate（逐層差分）、
  [`docs/re/47`](../re/47-main-screen-window-registry.md)（主畫面的原版幾何）
- 推論等級：分區座標**confirmed**（從機器碼算出來的）；判定門檻是本規格的設計

## 1. 這份規格要解決什麼

「主畫面跟原版不一樣」目前是一句話。**它要變成一張逐區、逐像素的清單**，
否則只能繼續調常數然後憑感覺判斷有沒有變好。

`docs/re/20` §4 已經指出目前所有畫面證據的共同缺陷：

> 用不同戰況的 YouTube frame 做 layout-only 比較，
> **沒有同一戰場、同一雙方、同一命令、同一時刻的 DOS/V 對照**。

## 2. 同狀態怎麼達成

**用存檔定位狀態，不要用操作序列。** 即時制送同一串按鍵每次跑到的時間點都不同
（`CLAUDE.md` §3.1），所以：

```
① 原版在 DOSBox-X 裡玩到某個狀態 → 存檔（原版自己的存檔功能）
② 把那份 SAVE.DAT 複製出來（唯讀掛載，不動原始素材）
③ remake 載入同一份 SAVE.DAT 的同一槽
④ 兩邊各截一張 640×400
```

原生存檔格式（[`20`](20-save-format.md)）在這裡的作用是**反向**的：
remake 讀原版 `SAVE.DAT` 建立狀態，之後每一輪對拍都從同一份原生檔載入，
確保 remake 這一側每次都完全一樣。

### 2.1 時間要凍結

即時制的畫面每個 tick 都在動。兩邊都必須停在**不會變的畫面**：

| 側 | 怎麼凍 |
|---|---|
| 原版 | 用除錯器 `execution.pause`（[`../playtest/21`](../playtest/21-dosboxx-bridge-sampling.md)），停住之後再截圖 |
| remake | 固定亂數種子 ＋ 固定 tick 數，截圖前不再 `Tick`（`CLAUDE.md` §9）|

**沒有凍結就不要對拍**——差出來的像素會混進「時間差」，
而那種雜訊看起來跟真的版面錯誤一模一樣。

## 3. 分區

座標出自 [`docs/re/47`](../re/47-main-screen-window-registry.md)，
**每一條邊都是從機器碼算的**：

| 區 | 矩形（x, y, w, h）| 來源 |
|---|---|---|
| `banner` | 0, 0, 640, 32 | `sub_18755`：`sub_1E3D7(al=6, cx=0x0450)` → 640×32 |
| `command` | 0, 32, 432, 32 | `sub_1614A` → `sub_1895D(cx=0x021B)`，Y 由 `si=bx+2` 下移 32 |
| `map` | 0, 64, 432, 336 | 三個常駐視窗的補集 |
| `minimap` | 432, 32, 208, 160 | `sub_15A3A` → `sub_1895D(cx=0x0A0D)` |
| `faction` | 432, 192, 208, 208 | `sub_15E2D` → `sub_1895D(cx=0x0D0D)` |

右欄三段相加 `32 + 160 + 208 = 400`，剛好鋪滿畫面高度——
這是矩形讀對了的算術檢查。

## 4. 判定

每一區輸出三個數字：**不同像素數**、**佔該區比例**、**最大色差**。

| 等級 | 條件 | 意思 |
|---|---|---|
| `PASS` | 不同像素 ＝ 0 | 逐像素相同 |
| `NEAR` | ≤ 該區 0.5% | 差在少數像素，通常是字型或抗鋸齒 |
| `FAIL` | 其餘 | 版面或素材不對 |

**這是保存專案，`PASS` 才是目標，`NEAR` 只是分類不是及格。**
每一項 `NEAR`／`FAIL` 都要在 `docs/playtest/` 記下差在哪、為什麼。

## 5. remake 實作

| 項目 | 位置 |
|---|---|
| 原版擷取 | `tools/dosv_capture.sh`（主機端入口）→ `tools/dosv_live_capture.sh`（容器內），再 `tools/parity_crop.py` 切成 640×400 |
| remake 擷取 | `tools/parity_shot.sh`——用 `wlgame -shot` 直接寫**邏輯畫面** 640×400。`tools/shot.sh` 抓的是 1600×900 桌面，尺寸對不上 |
| 分辨「位移」與「畫錯」 | `tools/parity_shift.py`（平移搜尋）、`tools/parity_locate.py`（拿地標去另一張找）|
| 看差在哪一種 | `tools/patch_zoom.py`（同一小塊並排放大）、`tools/palette_compare.py`（先排除調色盤刻度）|
| 逐區差分 | `tools/parity_diff.py` ✅ 正對照已接進 `tools/check.sh` |
| 紀錄 | `docs/playtest/` 的一份新文件 ＋ 差分圖放 `docs/playtest/parity/` |

## 6. 驗證

| 方式 | 內容 |
|---|---|
| 自我檢查 ✅ | `--selftest`：同圖每區 0、平移 1 px 每區非 0。**沒有這個正對照，「全 PASS」可能只是工具沒在比**。已接進 `check.sh` |
| 尺寸不合 ✅ | 大小不同直接報錯——縮放後硬比會把「尺寸不對」這個最重要的線索洗掉 |
| 對原版 ✅ | 第一輪的實際結果在 [`../playtest/37`](../playtest/37-main-screen-parity.md)：地圖區 0.5%，`faction` 區逐像素 PASS |

## 7. 未解

| 項目 | 現況 |
|---|---|
| 各視窗**內部**的排版 | 分區的外框已由機器碼定死（§3），框內的頭像／文字列座標仍是影片估值（[`docs/spec/12`](12-strategy-chrome.md) §7）|
| 送點擊的座標 | DOSBox-X 的**視窗**是 640×480，遊戲的 640×400 在 y 偏移 40（`tools/parity_crop.py` 量的），而 INT 33 把整個視窗等比對映到遊戲畫面——**送 y 要乘 1.2，不是減 40**。這是本機設定的性質，把 `int33 max y` 改成 400 應該就 1:1，還沒試 |
| 主畫面的四窗狀態 | 開局四個視窗全關（`sub_11A6E` 結尾 `mov cs:byte_198A6, 0`）。要開得先移游標再按同一點（`docs/re/47` §3.1），單純 `click` 會被當成移動吃掉 |
| 調色盤季節組 | 兩側都要鎖同一組，否則整片顏色不同（`docs/formats/02`）|
