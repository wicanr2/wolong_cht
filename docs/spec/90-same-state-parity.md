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
| 原版 | ⭐ **開系統選單**——它開著時遊戲時間停止（說明書 3.1），畫面也不再變（[`../playtest/39`](../playtest/39-system-window-parity.md)）。比接除錯器便宜得多。要更精細才用 `execution.pause`（[`../playtest/21`](../playtest/21-dosboxx-bridge-sampling.md)）|
| remake | 固定亂數種子 ＋ 固定 tick 數，截圖前不再 `Tick`（`CLAUDE.md` §9）|

**沒有凍結就不要對拍**——差出來的像素會混進「時間差」，
而那種雜訊看起來跟真的版面錯誤一模一樣。

### 2.2 ⭐ 從原版存檔開局：整條管線可跑了

`readSave` 找不到原生存檔就退回 `state.LoadScenario`，而那一支吃的
**就是原版的四區塊格式**——所以原版存的 `SAVE.DAT` 直接載得動，
不必另外寫轉檔。

原版側（座標是視窗座標 ＝ 遊戲座標 × 1.2）：

```
click:46,0;press      系統選單（開著時時間停止，§2.1）
click:376,192         資料儲存 → 跳出 SAVE DATA 四格
press                 存到第 1 格（游標還停在那一列）
rclick:376,192        關掉 SAVE DATA，系統選單留著
savefile:SAVE.DAT     把 guest 寫的存檔複製到輸出目錄
```

remake 側：

```
tools/parity_shot.sh out.png -direct -scenario 0 -player 0     -save-file <那份 SAVE.DAT> -load-slot 0 -open-window -3 -cam X,Y
```

⚠ **鏡頭不在存檔裡**：remake 載入後把鏡頭移到首都，原版停在玩家最後
捲到的地方，所以還是要用 `-cam` 對齊（`tools/find_camera.py` 反推）。

⚠ 存檔也裝不下**執行期狀態**（天候物件、事件佇列的部分內容）。
存得越晚，這一類差異越多——所以要對「載入後的畫面」就存得早一點，
要量「存檔漏了什麼」就故意存晚一點。

### 2.3 戰術畫面的同狀態：用存檔裡**現成的軍團**，不要現編

戰場那一側對到剩下的差異全是「雙方是誰、各帶多少兵」
（[`../playtest/40`](../playtest/40-tactical-parity.md) §2），
而驗收捷徑 `demoBattle` 是**現編**兩支軍團（各 2 騎 2 弓 2 步、預備兵灌 9000），
所以名字與兵力一定對不上。

原版那一側沒有辦法「在戰鬥中存檔」，但**戰鬥前一刻可以**：

```
編成一支軍團 → 等 AI 來打 → 開戰之前開系統選單存檔 → 關掉視窗 → 等它開打
```

系統選單開著時時間停止（§2.1），所以那一刻的存檔就是開戰前的局面；
`SAVE.DAT` 裡的軍團記錄含**六個位置的兵種與人數**
（[`../formats/08`](../formats/08-sinario-save.md)），開戰時不會再變。

remake 這一側因此需要「拿既有的兩支軍團開戰」而不是現編：

| 旗標 | 作用 |
|---|---|
| `-load-slot N` | 從 `-save-file` 的第 N 槽載入（§2.2）|
| `-list-corps` | 把載入後**還活著的軍團**印出來（編號、勢力、主將、兵力），用來挑下一個旗標的參數 |
| `-siege-corps 攻,守` | 指定兩支**既有**軍團打攻城戰；不呼叫 `FormCorps`，兵種與人數照存檔 |
| `-battle-steps N` | 截圖前推進幾個戰術 tick（預設 120）。原版的開場對白那一幀是 **0** |

`-siege-node` 仍然決定打哪一張戰場；`-siege-defend` 仍然決定玩家站哪一邊
（守方要把戰場轉 180 度，[`56`](56-battlefield-rotation.md)）。

### 2.4 ⭐ 野戰的同狀態：改造存檔 ＋ LOAD DATA

野戰要兩支軍團在野外同格，而原版的「行軍指示」在選單第二列、
滑鼠時間線點不到（[`../playtest/40`](../playtest/40-tactical-parity.md) §1.2）。
繞法是**直接改存檔**（格式已全解）把玩家軍團擺到 AI 軍團的行進路徑上，
原版從標題畫面 LOAD DATA 讀入（NEW GAME 答 NO；槽的熱區是日期鈕），
幾個遊戲時辰後自然遭遇。配方、欄位與「佔用圖快取欄 `+0x1A`/`+0x1C`
不搬就是幽靈」的坑見 [`../playtest/43`](../playtest/43-field-battle-parity.md) §2–§3。

原版側用 `WOLONG_DOSV_SEED_SAVE` 預置存檔；remake 側
`-save-file … -load-slot 0 -encounter-choose`（遭遇出現時自動選戰鬥指揮，
照自然流程進戰場、不重擺軍團——`-siege-corps` 那條會重擺，戰場格會跑掉）。
⚠ 像素比對一律用 `shot:` 的 PNG，`grab-start:` 的 mp4 是失真編碼，
逐區 diff 會全部 99%。

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

### 4.1 ⭐ 參考影格本身會有東西

原版側是**實錄影格**，所以「原版有、remake 沒有」不自動等於
remake 的缺陷。已經遇過兩種：

| 來源 | 樣子 | 處置 |
|---|---|---|
| **滑鼠游標** | 紅白箭頭停在畫面中間（戰場對拍吃了 95 px）| 換一張沒有游標的影格，或標成不可消 |
| **每次擲骰的初值** | 旗的揮舞相位（`rand & 3` 起手，116 px）| 標成不可消 |

⚠ **每一群差異都要放大看過再歸類**，不要只看像素數就往 remake 身上算。
到 0.2% 這個量級，錄影本身的東西會變成主要成分
（[`../playtest/40`](../playtest/40-tactical-parity.md) §13）。

## 5. remake 實作

| 項目 | 位置 |
|---|---|
| 原版擷取 | `tools/dosv_capture.sh`（主機端入口）→ `tools/dosv_live_capture.sh`（容器內），再 `tools/parity_crop.py` 切成 640×400 |
| 原版存檔 | `dosv_live_capture.sh` 的 `savefile:檔名` 步驟——把 guest 寫的 `SAVE.DAT` 複製到輸出目錄（**原版資產，只落在 gitignore 的 workplace/**）|
| remake 擷取 | `tools/parity_shot.sh`——用 `wlgame -shot` 直接寫**邏輯畫面** 640×400。`tools/shot.sh` 抓的是 1600×900 桌面，尺寸對不上 |
| 對齊鏡頭 | `wlgame -cam X,Y`（驗收旗標）。**進到大地圖之後移動滑鼠就會捲動**（`docs/re/47` §3.1），所以原版側只要點過視窗開關，鏡頭就不在開局位置了；要比縮小地圖的視野框就得讓 remake 跟著搬 |
| 分辨「位移」與「畫錯」 | `tools/parity_shift.py`（平移搜尋）、`tools/parity_locate.py`（拿地標去另一張找）|
| 看差在哪一種 | `tools/patch_zoom.py`（同一小塊並排放大）、`tools/palette_compare.py`（先排除調色盤刻度）|
| 逐區差分 | `tools/parity_diff.py` ✅ 正對照已接進 `tools/check.sh` |
| 紀錄 | `docs/playtest/` 的一份新文件 ＋ 差分圖放 `docs/playtest/parity/` |

### 5.1 remake 側的視窗 fixture 旗標

對拍要讓 remake 停在「原版剛開某個視窗」的同一個狀態，
每個視窗一支 `wlgame` 旗標（只供截圖，不進正常操作路徑）：

| 旗標 | 停在 | 對應原版 |
|---|---|---|
| `-open-list` | 武將一覽剛開、無選取 | 指令列 #4 |
| `-open-form-pick`（＋`-form-pick-row N`） | 編成的武將一覽 | 指令列 #3 |
| `-open-finance`（＋`-finance-amount N`） | 財政視窗（＋數值輸入器）| 指令列 #2 |
| `-advise-menu` | 進言五項選單 | 指令列 #1 |
| **`-advise-target`** | **進言 → 交戰的目標勢力清單**（選單第 0 列剛按下：清單開著、軍師框問「請選擇交戰之勢力。」）| 指令列 #1 → 第 0 列 |
| `-open-corps` | 軍團一覽 | 指令列 #4（軍團）→ 一覽 |
| **`-open-cities`** | 據點一覽（自勢力據點）| 指令列 #5（據點）|
| **`-open-factions`** | 勢力一覽（他勢力）| 指令列 #7（勢力）|
| **`-open-cityinfo N`** | 據點情報卡（`docs/spec/23`；N＝據點編號，−1＝玩家首都）| 指令列 #5（據點）的預設卡／地圖點據點 |

新的視窗對拍先在這張表登一列再動 `cmd/wlgame`。

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
