# 55 — 縮小地圖的視野框是點陣，而且解開了「差四格」

**狀態：CONFORMED。** 框的點陣、尺寸與位置算式都解出來了，
remake 畫上去之後縮小地圖那一區**逐像素 PASS**。
⭐ **順帶解掉一個掛了兩份規格的謎**：那個「差 4 格」不是原版的怪癖，
是我們的地圖解碼少跳了開頭 4 byte（[`../formats/05`](../formats/05-mmap-worldmap.md) §2.1）。

- 日期：2026-08-17
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 出處：`sub_15C58`（算位置）、`sub_196ED` ＋ `sub_19752`（畫）、
  [`../re/62`](../re/62-strategy-minimap.md) §4
- 實機證據：`workplace/promo-live/parity-win/w3-all.png`（鏡頭 (4,0)）、
  `workplace/promo-live/parity-main2/m3-all.png`（鏡頭 (214,10)）
- 推論等級：**confirmed**（點陣直接從檔案解出來，位置兩個不同鏡頭各驗一次）

## 1. 框是一張 24×11 的小圖

`sub_196ED` 把 `ds` 設成 `cs:word_10D4C`——`ICONGRF` 段 3 的 **`+0x8F0`**，
緊接在數字字模（`+0x840`）後面的另一塊 176 byte
（[`52`](52-main-screen-camera-and-banner-date.md) §4 的切法）。

`sub_19752` 每列讀 **3 byte ＝ 24 px**、跑 **11 列**（`bp = 0Bh`），
被 `sub_196ED` 呼叫**五次**，set/reset 值依序是 `0, 1, 2, 4, 8`：

```
第一趟 al=0：把遮罩內的像素清成色 0
後四趟      ：以 OR 把四個色平面疊上去（1|2|4|8 ＝ 15 ＝ 白）
```

⇒ 5 × 33 ＝ **165 byte**，一張 4 bpp 的小圖 ＋ 一張遮罩。

解出來的圖形（`.` ＝ 不畫、`F` ＝ 白、`0` ＝ 黑）：

```
.FFFFFFFFFFFFFFFFF0.....
FF000000000000000FF0....
F00..............0F0....
F0................F0....
F0................F0....
F0................F0....
F0................F0....
F0................F0....
FF0..............FF0....
0FFFFFFFFFFFFFFFFF00....
.000000000000000000.....
```

⭐ **白邊 ＋ 右下黑影的立體框**，實際只用到左邊 20 px。
用 `vector.StrokeRect` 描一個矩形永遠對不上。

## 2. 位置：`camX/2 + 440`，沒有偏移

`sub_15C58` 算位置：

```
dx = ds:988Eh (camX) >> 1 + 1B8h   ; 440
bx = ds:9890h (camY) >> 1 + 28h    ;  40
```

`word_1988E` 就是畫面上 x=0 那一欄——`sub_11F7F` 的
`mov ax, ds:9882h / shr ax,4 ×4 / mov ds:988Eh, ax` 把鏡頭像素直接除以 16，
而 `sub_1D615` 拿它當 `di` 去索引地圖列。**中間沒有任何加減。**

把兩張實機截圖的鏡頭用 `tools/find_camera.py`（從畫面認圖塊、回地圖找落點）
反推，再量框的位置：

| 鏡頭 | 框的白邊左上角 | `440 + camX/2` |
|---|---|---:|
| (0, 0) | (441, 40) | 440 ＋ 白邊在第 1 欄 ✅ |
| (210, 10) | (546, 45) | 545 ＋ 白邊在第 1 欄 ✅ |

點陣的白邊從第 1 欄開始，所以圖的原點是 440 與 545。

⚠ **量鏡頭要從 offset 4 讀地圖。** `MMAP.MAP` 解壓後開頭那 4 byte 是
長度欄位（[`../formats/05`](../formats/05-mmap-worldmap.md) §2.1）；
少跳它，反推出來的鏡頭會一律多 4 欄，這一格就會看起來差 −2 px。
`tools/find_camera.py` 已經跳過。

Y 一樣是 `40 + camY/2`，一格不差。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 解點陣 | `gfx.DecodeViewBox` → `Library.ViewBox` |
| 畫 | `cmd/wlgame` 的 `drawMinimapViewBox`，位置 `(minimapX + (camX−4)/2, minimapY + camY/2)` |
| 那個 4 | `minimapCamBias` |

## 4. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestViewBoxShape`：24×11、只有色 0／15 與透明，第 0 列的白邊 17 px |
| 對拍 ✅ | [`../playtest/38`](../playtest/38-window-parity.md)：縮小地圖那一區 **0.00% PASS** |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 剩下的 11 byte | `+0x8F0` 那一塊有 176 byte，框只用 165 |
