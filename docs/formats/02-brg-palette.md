# 02 — `.BRG` 調色盤格式

**狀態：READY**（可以照這份實作）

- 日期：2026-08-07
- 出處：`docs/re/02-palette-routine.md`（兩版 `KI.EXE` 的機器碼）
- 工具：`tools/brg.py`
- 檔案：四個，**兩版 byte-for-byte 完全相同**

## 1. 結構

```
每色 3 byte：   [0] 藍   [1] 紅   [2] 綠        ← 副檔名 BRG 就是通道順序
每通道         4 bit，值域 0–15（高 4 位固定為 0）
每 16 色       一組（bank），硬體一次只吃一組
```

| 檔案 | 大小 | 色數 | 組數 |
|---|---:|---:|---:|
| `GAMEPAL.BRG` | 384 | 128 | 8 |
| `OPENPAL.BRG` | 288 | 96 | 6 |
| `ENDPAL.BRG` | 576 | 192 | 12 |
| `OVERPAL.BRG` | 48 | 16 | 1 |

## 2. 轉成 8 bit sRGB

原版的亮度縮放（淡入淡出也走同一條）：

```
scaled = ((v << 4) * brightness + 0x80) >> 8        brightness 0–16，16 = 全亮
```

`brightness = 16` 時 `scaled == v`。轉 8 bit **要看是哪一版的硬體**：

```
DOS/V：DAC  = scaled << 2          # 6 bit，上限 60（不是 63）
       srgb = (DAC << 2) | (DAC >> 4)   # VGA 的 6→8 位元複製
PC-98：srgb = scaled * 255 // 15   # 類比調色盤本來就是 4 bit／通道
```

⭐ **DOS/V 的顏色到不了滿刻度**：`shl ah,1` 兩次讓 15 變成 DAC 60，
出來是 `#F3` 不是 `#FF`。**松崗版畫面上沒有純白。**
規格與實機證據（春組 16 色全中）見
[`../spec/51`](../spec/51-vga-dac-palette-scale.md)。

> **兩版的顏色因此不完全一樣**，不是同一條轉換式。
> remake 走 DOS/V 那一條。

## 3. 顯示是 16 色 planar

兩版都是 4 平面 16 色（DOS/V 走 VGA Graphics Controller ＋ Sequencer，
PC-98 走 GRCG）。**`*GRF.DAT` 的圖像因此是 4 bpp。**

畫面上同時只有 16 色；`.BRG` 的多組是**切換**用的，不是同時生效。

## 4. `GAMEPAL.BRG` 的分組

| 組 | 用途 |
|---|---|
| 0 | 春（色 14 ＝ `#82A261` 灰綠） |
| 1 | 夏（`#51A210` 鮮綠） |
| 2 | 秋（`#D38200` 橙褐） |
| 3 | 冬（`#F3F3F3` 雪白） |
| 4–7 | **「液晶」畫面模式**的四季（`docs/re/02` §6.2）|

**季節只換色號 14 這一個顏色**，其餘 15 色四季共用。
地表植被統一畫成色號 14，換調色盤就換季。

> **remake 要照這個機制做**，不要改成四套素材——
> 換素材與換調色盤的視覺行為不一樣（換調色盤是整個畫面所有色 14 的像素同時變）。

## 5. 工具

```sh
tools/py.sh tools/brg.py info   workplace/orig/dosv/GAMEPAL.BRG
tools/py.sh tools/brg.py swatch workplace/orig/dosv/GAMEPAL.BRG out.png 24
```

四個調色盤的色票在 `docs/images/palette-*.png`。

![GAMEPAL 八組色票](../images/palette-gamepal.png)

上四列是四季（肉眼幾乎看不出差別，因為只差一個顏色），下四列是未解的那四組。

## 6. 還沒解的

- 誰載入、誰選組（`docs/re/02` §7）。
- `OPENPAL` 6 組與 `ENDPAL` 12 組各自對應哪些畫面。
