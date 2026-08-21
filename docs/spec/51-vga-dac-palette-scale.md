# 51 — DOS/V 的顏色到不了滿刻度：4 bit → VGA 6 bit DAC → 8 bit

**狀態：CONFORMED。** 算式有機器碼出處，也在原版實機截圖上驗過
（`GAMEPAL.BRG` 春組 **16 色全中**），整條解碼路徑都走它，
主畫面因此逐像素相同（§6）。
⭐ **松崗版的白色是 `#F3F3F3`，不是 `#FFFFFF`。**

- 日期：2026-08-17
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 出處：[`../re/02`](../re/02-palette-routine.md) §5（`shl ah,1` 兩次）、
  [`../formats/02`](../formats/02-brg-palette.md) §2
- 實機證據：`workplace/promo-live/parity-probe/p2-ok2.png`
  （松崗 DOS/V，196年4月1日的主畫面，DOSBox-X `machine=vgaonly`）
- 工具：`tools/palette_dacscale.py`（新增）、`tools/palette_compare.py`（新增）
- 推論等級：**confirmed**

## 1. 問題

`.BRG` 每個通道是 **4 bit**（0–15）。remake 一直用

```
srgb = v * 255 / 15          # 15 → 255
```

把它鋪滿 0–255。**DOS/V 不是這樣走的。**
[`../re/02`](../re/02-palette-routine.md) §5 的亮度常式在 PC-98 到 `ah` 為止，
DOS/V 多了兩次 `shl ah,1`：

```
DAC 值 = v << 2 = 4v          v ≤ 15  ⇒  DAC ≤ 60
```

VGA DAC 是 **6 bit（0–63）**，而這條路徑最大只給到 **60**。
換句話說**類比輸出永遠到不了滿刻度**：色號 15 是 63 的 60/63 ≈ 95.2%。

## 2. 算式

```
DAC   = Scale(v, brightness) << 2            # 0–60
srgb  = (DAC << 2) | (DAC >> 4)              # VGA 的 6→8 bit 位元複製
```

| 4 bit `v` | 舊：`v*255/15` | **新：位元複製** | 差 |
|---:|---:|---:|---:|
| 0 | 0 | 0 | 0 |
| 2 | 34 | 32 | −2 |
| 4 | 68 | 65 | −3 |
| 6 | 102 | 97 | −5 |
| 8 | 136 | 130 | −6 |
| 10 | 170 | 162 | −8 |
| 12 | 204 | 195 | −9 |
| 13 | 221 | 211 | −10 |
| 14 | 238 | 227 | −11 |
| 15 | 255 | **243** | −12 |

## 3. 為什麼是位元複製，不是類比四捨五入

「照類比刻度算」會得到 `round(4v × 255/63)`，與位元複製只差在 `v = 3, 12, 13`
三個值。兩種都講得通，所以**用實機截圖裁決**：

```
tools/py.sh tools/palette_dacscale.py workplace/orig/dosv/GAMEPAL.BRG \
    workplace/parity/orig-closed.png
```

| 換算 | 春組 16 色裡出現在實機截圖的 |
|---|---:|
| `v*255/15`（舊）| **1/16**（只有黑色）|
| 位元複製 | **16/16** |
| 類比四捨五入 | 12/16 |

⭐ **判準是「那個顏色有沒有真的出現在實機畫面上」**，不是哪條式子漂亮。
截圖只用了 15 色，而位元複製那一組整組都在裡面。

> ⚠ 這條只適用 **DOS/V**。PC-98 的類比調色盤本來就是 4 bit／通道，
> 15 就是滿刻度——那一版的白是 `#FFFFFF`。
> **兩版的顏色本來就不完全一樣**，不是同一條轉換式。

## 4. 影響範圍

每一個像素。整張畫面的每個通道都亮了 4%，所以
[`../playtest/37`](../playtest/37-main-screen-parity.md) 的逐區對拍
在修之前**五區全 FAIL、最大色差 255**——看起來像版面全錯，
實際上版面沒問題。

## 5. remake 實作

| 項目 | 位置 |
|---|---|
| 換算 | `internal/assets/palette/palette.go` 的 `toSRGB` |
| 亮度 | `Scale` 不動（淡入淡出仍是 0–16 那條）|

## 6. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestSeasonColor`：四季的色 14 改成 `#82A261`／`#51A210`／`#D38200`／`#F3F3F3` |
| 單元測試 | `TestDACScaleNeverReachesFullScale`：`v=15` 必須是 243，且任何 `v` 都 ≤ 243 |
| 實機 | [`../playtest/37`](../playtest/37-main-screen-parity.md) 的逐區對拍 |
