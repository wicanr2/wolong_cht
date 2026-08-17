#!/usr/bin/env python3
"""把兩張圖的同一小塊並排放大，用來看「差在哪一種」。

    tools/py.sh tools/patch_zoom.py a.png b.png x y 邊長 倍率 out.png

整片 60% 不同但肉眼看起來一樣時，通常是網點相位或圖塊選錯。
兩者在縮圖上分不出來，放大到一個像素一格就分得出來。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png, write_png


def main():
    if len(sys.argv) != 8:
        raise SystemExit("用法: patch_zoom.py a.png b.png x y 邊長 倍率 out.png")
    pa, pb = sys.argv[1], sys.argv[2]
    x, y, side, zoom = (int(v) for v in sys.argv[3:7])
    dst = sys.argv[7]
    wa, ha, a = read_png(pa)
    wb, hb, b = read_png(pb)
    gap = 8
    ow = (side * 2 + gap) * zoom
    oh = side * zoom
    rows = [[(0, 0, 0)] * ow for _ in range(oh)]
    for dy in range(side):
        for dx in range(side):
            ca = a[y + dy][x + dx]
            cb = b[y + dy][x + dx]
            for zy in range(zoom):
                row = rows[dy * zoom + zy]
                for zx in range(zoom):
                    row[(dx) * zoom + zx] = ca
                    row[(side + gap + dx) * zoom + zx] = cb
    write_png(dst, ow, oh, rows)
    print("左＝%s 右＝%s，(%d,%d) 起 %d×%d，放大 %d 倍 → %s"
          % (pa, pb, x, y, side, side, zoom, dst))


if __name__ == "__main__":
    main()
