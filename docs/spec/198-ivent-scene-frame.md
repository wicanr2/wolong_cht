# 198 — 事件插圖少了外面那個框

**狀態：CONFORMED。** 原版貼 `IVENTGRF` 插圖之前會先畫一個
**(48, 128) 304×192** 的框，插圖 (56, 136) 288×176 蓋在它中間，
四邊各露 8 px ＝ 看得見的黃色外框。remake 的 `drawIventScene`
**只貼圖不畫框**，說服場景與判決畫面因此各差 6,000–8,800 px。

- 日期：2026-09-11
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_13D68`
- 推論等級：confirmed（尺寸來自同一條 `sub_10C14` 呼叫的立即值）
- 相關：[`45`](45-advise-scene-layout.md) §1.0（插圖座標從 VRAM 位址反推）、
  [`26`](26-yes-no-dialog.md)（同一支 `sub_10C14` 的另一個呼叫端）

## 1. 原版做什麼

```asm
sub_13D68:
    dx = 3 / bx = 8 / cx = 0C13h / call sub_10C14   ; ★ (48, 128) 304×192 的底
    ds = cs:word_19876
    bx = 2A87h / ax = 0B012h / call sub_1FA37       ; 插圖 (56, 136) 288×176
```

`cx = 0C13h` ＝ `ch = 0Ch`（12 個 16 px ＝ 192 高）、`cl = 13h`
（19 個 16 px ＝ 304 寬）；`dx = 3`／`bx = 8` 是粗格 (3, 8) ⇒ (48, 128)。

⇒ 框比插圖大 8 px 一圈，這一圈就是畫面上看得到的邊。
**這一行早就寫在 [`45`](45-advise-scene-layout.md) §1.0 的組語裡**——
當時只用它反推插圖座標，沒有注意到同一段還畫了一個底。

## 2. remake 錯在哪

`cmd/wlgame/talkscene.go` 的 `drawIventScene` 直接 `DrawImage` 到
(56, 136)，沒有任何框。五個呼叫端（說服、外交、判決、事件通知、撥款）
共用這一支，所以五個畫面同時少那一圈。

## 3. remake 實作

`drawIventScene` 在貼圖前先 `g.chrome.Window(screen, 48, 128, 304, 192, …)`。
一處改動涵蓋五個呼叫端（`CLAUDE.md` §7 第 6 條）。

## 4. 驗證

`tools/parity_screens.sh advise-scene sortie` 的 `map` 區
（[`../playtest/121`](../playtest/121-screen-parity-gate.md) §3）：

| | 修之前 | 修之後 |
|---|---:|---:|
| `sortie` `map` | 6,998 | **1,509**（只剩 remake 自加的「Enter 繼續」）|
| `advise-scene` `map` | 8,843 | **3,354**（另有進言選單殘影 1,750，見 §5）|

插圖本體 `(56,136,288,176)` 兩組都是 **0 / 50,688**，框接上之後沒有動到它。

## 5. 未解

- **進言選單在說服場景上沒關掉**（remake 殘影 1,750 px，原版進場時清掉）。
  與本規格無關，但同一張畫面上，記在
  [`../playtest/121`](../playtest/121-screen-parity-gate.md) §7。
- 框的填色在原版是什麼：插圖把中間 288×176 全蓋住，露出來的只有邊框那
  一圈，所以填色在這個畫面上量不到。事件通知那一種（`messages.go`，
  插圖頁不是 0）可能露得出來，還沒對過。
