#!/usr/bin/env python3
"""把視窗底紋的排法拿實機截圖驗一次。

    tools/py.sh tools/window_texture_check.py

## 結論（docs/formats/03 §5.5）

底紋是 `ICONGRF.DAT` 檔尾的 128 byte（32×32、1 bpp），
而且**釘在螢幕上，不是釘在視窗左上角**：

    像素 (x, y) ＝ tile[y mod 32][x mod 32]

## 為什麼要留這一支

前幾輪都從視窗角落開始平鋪，怎麼試都只有 63–83%，於是把
「垂直週期 96」寫進文件當成待解的謎。**那個 96 是量錯的**——
截圖裡最高的一塊純底紋只有 50 列，位移 96 的自相關根本沒有重疊。
換成「逐列比回磚塊、讀出實際列序」之後，序列是乾淨的 0,1,2,…,31,0,…，
週期就是 32。

**這一支存在的意義是讓那個結論可以重跑**，不是給人讀數字。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png

SHOT = "workplace/dosbox/shots/x04.png"
ICON = "workplace/orig/dosv/ICONGRF.DAT"
OFF, SIZE = 0xBA20, 32
BLACK, BLUE = (0, 0, 0), (0, 32, 101)
# 截圖裡三塊確認是純底紋的區域（其餘地方有文字與圖示，同色但不是底紋）。
CLEAN = ((377, 409, 166, 216, "右側 32×50"),
         (103, 209, 138, 166, "左上乾淨帶 106×28"),
         (103, 209, 186, 214, "左下乾淨帶 106×28"))


def tile():
    blob = open(ICON, "rb").read()[OFF:OFF + SIZE * SIZE // 8]
    if len(blob) != SIZE * SIZE // 8:
        raise SystemExit("ICONGRF.DAT 太短，取不到檔尾那 128 byte")
    return [[(blob[r * 4 + (i >> 3)] >> (7 - (i & 7))) & 1 for i in range(SIZE)]
            for r in range(SIZE)]


def main():
    t = tile()
    w, h, px = read_png(SHOT)
    worst = 100.0
    for x0, x1, y0, y1, name in CLEAN:
        same = tot = 0
        for y in range(y0, y1):
            for x in range(x0, x1):
                c = px[y][x]
                if c != BLACK and c != BLUE:
                    continue
                tot += 1
                same += (1 if c == BLUE else 0) == t[y % SIZE][x % SIZE]
        rate = 100.0 * same / max(tot, 1)
        worst = min(worst, rate)
        print("  %-22s %.1f%%（%d 點）" % (name, rate, tot))
    if worst < 95.0:
        raise SystemExit("最差 %.1f%% < 95%%：排法不對" % worst)
    print("✅ 螢幕對齊平鋪：三塊都 ≥ 95%%（最差 %.1f%%）" % worst)


if __name__ == "__main__":
    main()
