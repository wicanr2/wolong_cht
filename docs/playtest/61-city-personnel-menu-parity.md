# 61 — 「據點」與「人事」兩張彈出選單的對拍：四區全部 0 px

**狀態：通過。** 兩張都是指令列、大地圖、縮小地圖、自勢力情報**四區逐像素相同**，
選單框與反白的指令列格也各 0 px。`banner` 各剩 116 px，
與 [`60`](60-corps-menu-parity.md) 同一個成因（原版 4月20日、remake 4月1日）。

⭐ **這一輪沒有抓到任何缺陷**——[`60`](60-corps-menu-parity.md) 修的那三件
（框寬、反白 XOR、指令列反白）在這兩張上一次就對。三張走同一份
`popupMenu` 實作（[`../spec/126`](../spec/126-command-popup-menus.md)），
而**位置與字數是各自的立即值**，所以這一次比的是那些立即值，不是繪製程式碼。

- 日期：2026-09-03
- 規格：[`../spec/126`](../spec/126-command-popup-menus.md)（四張走同一支常式）、
  [`../spec/124`](../spec/124-menu-highlight-xor.md)（反白）、
  [`../spec/125`](../spec/125-menu-box-width-from-padding.md)（框寬）
- 原版側：`WOLONG_DOSV_SEED_SAVE=SAVE-B.DAT tools/dosv_capture.sh
  parity-menu-city "…;shot:cmd;tap:30,5,5;wait:2;shot:menu"`
  （人事那一張把 `tap:30,5,5` 換成 `tap:10,5,5`），
  再用 `tools/parity_crop.py` 切成 640×400
- remake 側：`tools/parity_shot.sh out.png -direct -scenario 0 -player 0
  -seed 7 -open-command-menu city|personnel -cam 0,0 -shot-frames 1`

## 1. 結果

| 區 | 據點 | 人事 | 判定 |
|---|---:|---:|---|
| `banner` | 116 px（0.57%）| 116 px（0.57%）| FAIL＝日期 |
| `command` | **0** | **0** | **PASS**（含反白的那一格）|
| `map` | **0** | **0** | **PASS**（選單本身就在這一區裡）|
| `minimap` | **0** | **0** | **PASS** |
| `faction` | **0** | **0** | **PASS** |
| 選單框 | **0**（`240,64,112,48`）| **0**（`48,64,144,80`）| **PASS** |
| 反白的指令列格 | **0**（`264,40,48,16`）| **0**（`72,40,48,16`）| **PASS** |

反白格的 x ＝ `24 + 指令索引 × 48`：據點是第 5 格 ⇒ 264、人事是第 1 格 ⇒ 72。
與 [`60`](60-corps-menu-parity.md) 的軍團（第 4 格 ⇒ 216）是同一條算式，
三個獨立取樣點互相印證（[`../spec/126`](../spec/126-command-popup-menus.md) §1）。

## 2. 原版側的擷取

照 [`54`](54-menu-second-row-tap.md)：timeline 完全沿用 `parity-tap5` 那一輪，
**只換最後一下 `tap` 的 x**。進到大地圖之後 INT 33 的範圍是整個世界，
一個主機像素 ≈ 9.6 個遊戲像素，所以

| 指令 | 遊戲座標 x | `tap` 的 x |
|---|---:|---:|
| 人事（第 1 格）| 72–120 | **10**（≈ 96）|
| 軍團（第 4 格）| 216–264 | 25（≈ 240）|
| 據點（第 5 格）| 264–312 | **30**（≈ 288）|

⚠ **主機端的 `timeout` 殺不掉容器。** 第一次跑用 `timeout 900 tools/dosv_capture.sh`
包在前景，被工具層的兩分鐘上限先砍掉——腳本死了，容器還在跑，
輸出目錄只留下一個空的 `manifest.txt` 與停在 `step begin=click:320,245`
的 trace。**看起來像擷取失敗，實際上是被中斷。**
一次擷取含 `wait:130` 要三分半，要嘛背景跑，要嘛把上限開到夠。

## 3. 未解

| 項目 | 現況 | 下手點 |
|---|---|---|
| 日期對不上 | 原版跑到 4月20日才截到，`banner` 因此永遠 116 px | 同 [`39`](39-system-window-parity.md) §4：用存檔定位，或加一個「跑到指定日期」的驗收旗標 |
| 選完之後的流程 | 據點的「首都確認」／「據點一覽」、人事的四項，**選下去之後的畫面都沒有對拍** | 原版側在 `tap` 之後再 `tap` 一次那一列 |
| 進言那一張 | 五項選單**確定拍不到**（[`42`](42-window-parity.md) §5：30 ms 瞬按與右鍵回退兩條路都試過），而且 remake 還沒併進 `popupMenu` | — |
