# 131 — 原版 oracle 換成 dosgolem：不再需要 DOSBox

**狀態：CONFORMED。** 開機到三視窗全開共六格，對原版 DOSBox-X 擷取
**五區全部 0 px**（[`../playtest/65`](../playtest/65-dosgolem-oracle.md)）。
⭐ 除了畫面，**原版的變數與控制流也問得到**（`peek`／`watch`）。

- 日期：2026-09-05
- 出處：[dosgolem](https://github.com/wicanr2/dosgolem) 的
  [`spec/007` VGA 平面模式](https://github.com/wicanr2/dosgolem/blob/master/docs/spec/007-vga-planar.md)、
  [`spec/008` 臥龍傳要的服務](https://github.com/wicanr2/dosgolem/blob/master/docs/spec/008-wolong-services.md)、
  [`findings/003` 五格對拍](https://github.com/wicanr2/dosgolem/blob/master/docs/findings/003-wolong-boots-and-matches.md)。
  工作副本在 `~/cht/dosgolem-wolong`
- 推論等級：confirmed（逐點對拍，不是推的）

## 1. 原版做什麼

`KI.EXE` 開機 → `int 10h AX=0012h`（640×480、16 色、四平面）→
向 `INT 15h AH=50h` 要兩個字型常式的遠位址 → 用 EGA Set/Reset 畫字 →
主迴圈輪詢 `int 33h AX=5`（先問右鍵再問左鍵）→
遊戲時鐘由 `YNSOUND.COM` 用 291.3 Hz 的回呼推
（[`../re/61`](../re/61-timer-tick-source.md)）。

這四條路 dosgolem 都走得通，所以「原版跑到某個狀態的畫面」不必再經過
DOSBox ＋ Xvfb ＋ xdotool。

## 2. 兩條路的差別

| | DOSBox-X（`tools/dosv_live_capture.sh`）| dosgolem（`tools/dosgolem.sh`）|
|---|---|---|
| 畫面 | X11 視窗 640×480，遊戲在 y 偏移 40，要跑 `parity_crop.py` | 直接輸出裁好的 640×400 |
| 取樣點 | `wait:3` 之類的**秒數**，即時制之下每次停在不同的遊戲日期 | `wait`（畫面停住）或 **`until:196/4/9`（遊戲日期）** |
| 座標 | 視窗座標，`int 33h` 把整個視窗等比對映到 640×400 | **遊戲座標**（0–639 × 0–399），送什麼就是什麼 |
| 原版的內部狀態 | 只能從像素反推 | **直接讀記憶體**（`peek:0CF0:8`）|
| 原版的控制流 | 問不到 | **攔任一支常式**（`WOLONG_DOSGOLEM_WATCH=11D8E`），印暫存器與呼叫端 |
| 哪裡可以點 | 只能一格一格試 | **讀熱區圖**（`hotspots`）——`sub_1E453` 查的那張 80×50 格圖 |
| 右鍵／瞬按 | `rclick`／`tap` | 都有（瞬按是獨立動作——彈出選單長按會當場選走第一列）|
| 一條五格的時間軸 | 146 秒的 `wait` 加總 | **0.99 秒** |
| 決定性 | 靠 `cycles=fixed 20000` ＋ 固定 sleep，仍會漂 | 以指令數計時，同輸入同輸出 |

⚠ **視窗座標換算的分母是 479 不是 480。**
`遊戲 y ＝ 視窗 y × 399 ÷ 479`。大部分點與 `× 400 ÷ 480` 算出來一樣，
**只有少數差 1**——視窗 336 是 279 不是 280，而那 1 個像素讓主畫面的
游標整塊對不上（62 點），其餘四區全 0。`tools/dosgolem.sh` 的第三個參數
給 `dosbox` 就自動換算，舊腳本因此搬得過來。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 包裝器 | `tools/dosgolem.sh` |
| 上游 | dosgolem 的 `apps/wolong/`（畫面幾何、開機流程、座標換算、遊戲時鐘）|
| 差異 | 無。**這一支不改 remake 任何東西**，它換掉的是原版側的取樣方式 |

`tools/dosgolem.sh` 吐出來的 PNG 已經是 640×400，直接餵
`tools/parity_diff.py`，**不必再跑 `tools/parity_crop.py`**。

DOSBox-X 那條路**不刪**：它是這一條的正對照，而且 PC-98 版只有它跑得動
（dosgolem 沒有 PC-98 的機器層）。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 對原版 | [`../playtest/65`](../playtest/65-dosgolem-oracle.md)：五格 × 五區 ＝ 25 個比較全部 0 px |
| 單元測試 | dosgolem 側（下面那一行），每一支都做過突變測試 |

```sh
# 在 dosgolem 的工作副本裡
tools/go.sh test ./internal/... -run 'TestWriteMode3AppliesALU|TestFullIndex|TestSoundTimerCallbackActuallyFires'
```

## 5. 未解

| 項目 | 現況 |
|---|---|
| PC-98 版 | dosgolem 沒有 PC-98 的機器層（不同的顯示與字型架構）。那一版仍走 DOSBox-X |
| ~~載入既有存檔的路徑~~ | **已通**：`WOLONG_DOSGOLEM_GAMEDIR` 指到帶那份 `SAVE.DAT` 的目錄，走 NEW GAME → NO → LOAD DATA。對同一份存檔的原版擷取三區 0 px（[`../playtest/66`](../playtest/66-dosgolem-load-save.md)）|
| ⚠ 遊戲中的座標 | 大地圖有一層**捲動原點**（`畫面 ＝ 滑鼠 − 原點`），所以遊戲中要用 `sclick`／`stap`，選單畫面才用 `click`（[`../playtest/66`](../playtest/66-dosgolem-load-save.md) §2）|
| ⚠ `wait` 在遊戲中不適用 | 即時制的畫面永遠不會靜止，`wait` 會跑到預算上限。遊戲中用 `steps:` 或 `until:` |
| 戰術畫面 | 還沒試。戰場走的是同一套繪製層，預期沒有新的機器層缺口，但**沒驗過就是沒驗過** |
| 視窗 x → 遊戲 x 的換算 | 視窗 416 對到遊戲 415，而同一批的 300／360／450 都是 1:1。成因未查，只影響游標位置（[`../playtest/65`](../playtest/65-dosgolem-oracle.md) §3.1）|
| 音源 | `int 61h` 只記錄不模擬（時鐘回呼除外）。音訊 parity 仍走 [`29`](29-audio.md) 的錄音比對 |
