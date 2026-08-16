# 41 — 訊息框的版面常數

**狀態：CONFORMED。** 框、肖像、文字的座標全部有機器碼出處，
並在原版實錄影格上逐項覆驗。

- 日期：2026-08-17
- 出處：[`../re/66`](../re/66-message-box-geometry.md)（`sub_18810` → `sub_1895D`／
  `sub_1075B` → `sub_10BCD`，＋ `workplace/parity/f008-640.ppm` 覆驗）
- 推論等級：**confirmed（靜態 ＋ 影格）**

## 1. 原版的數字

| 項目 | 值 | 出處 |
|---|---|---|
| 框 | **(160, 160, 256, 80)** | `sub_1895D(bx=8, dx=0Ah, cx=0510h)`，16 px 粗格、Y 加 2 格 |
| 內容區 | (168, 168) 起 **240 × 64** | `sub_10BCD` 的四邊各內縮 8 px |
| 肖像 | **(168, 168)**，64 × 64 | `sub_1075B` 的 `dx×16 + 72` 是右緣；`KAOGRF` 每張 2,048 B |
| 文字起點 | **(240, 176)** | `dx×16 + 80`／`bx×16 + 16` |
| 行距 | **16 px** | 全形字高，`sub_106F9` 每字前進 16 |
| 每列可用寬 | **160 px ＝ 10 全形字** | 框右內緣 416−8 減文字起點 240 |
| 列數 | **4** | 內容高 64 ÷ 16 |
| 一般通知的肖像 | **`KAOGRF` 第 147 張**（`al = 0x93`）| 渲染出來與影格上那張臉相同 |

⭐ **原版只有一個訊息框。** 有沒有「講話的人」不改變版面——
一般通知也走同一個框，只是肖像固定用第 147 張。

## 2. remake 改了什麼

| 常數 | 舊值 | 新值 |
|---|---:|---:|
| `talkBoxX` | 24 | **160** |
| `talkBoxY` | 80 | **160** |
| `talkBoxW` | 232 | **256** |
| `talkBoxH` | 80 | 80 |
| `talkPortraitX` | 32 | **168** |
| `talkPortraitY` | 88 | **168** |
| `talkTextX` | 80 | **240** |
| `talkTextY` | 96 | **176** |
| `talkTextWidth` | 160 | 160 |
| `talkLinePitch` | 16 | 16 |
| `messagePageRows` | 5 | **4**（＝ `talkBoxRows`）|

⭐ **舊的文字起點會蓋到肖像**：肖像從框內 +8 起、64 px 寬，
到框內 +72 才結束，而文字從框內 +56 開始。新的 +80 剛好讓開。

沒有肖像可用時（`portraitPage < 0`）改用 `defaultPortraitPage = 147`，
版面與有肖像時完全一致——這是把 remake 先前「置中、依內容決定大小」
的第二種框收掉。

## 3. 為什麼 `talkTextWidth` 不用改

[`../playtest/32`](../playtest/32-talk-layout-fit.md) §4 記著「影格上量到約 275 px」，
那量的是**框**（256 px）不是文字區。文字區 160 px 本來就對——
`TALK.DAT` 的原文有 825 行剛好是 160 px（10 全形字），
訊息在檔案裡就折好了，遊戲不做換行（[`../re/66`](../re/66-message-box-geometry.md) §6）。

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestMessageBoxMatchesOriginalGeometry`：11 個常數逐一比對 §1 |
| 單元測試 | `TestTextDoesNotOverlapPortrait`：文字起點在肖像右緣之後 |
| 單元測試 | `TestAllTalkLinesFitTheirBox`（既有）：門檻改成 4 列 |
| 截圖 | `-open-talk` 與原版 `f008` 對位 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 事件場景的兩個講話位置 | (8, 88) 與 (136, 296) 由機器碼算出，**沒有影格覆驗**（[`../re/66`](../re/66-message-box-geometry.md) §8）。remake 的事件場景版面暫不動 |
| 框的底紋 | 龍紋底紋仍未解（[`../formats/03`](../formats/03-grf-images.md) §5.5），remake 用純色 |
