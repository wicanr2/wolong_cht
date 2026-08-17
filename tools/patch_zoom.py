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
        raise SystemExit("用法: patch_zoom.py a.png b.png x y 邊長|寬x高 倍率 out.png")
    pa, pb = sys.argv[1], sys.argv[2]
    x, y = int(sys.argv[3]), int(sys.argv[4])
    # 邊長可以寫成 `寬x高`——底列與橫幅都是又寬又扁的，方形切不出來。
    if "x" in sys.argv[5]:
        pw, ph = (int(v) for v in sys.argv[5].split("x"))
    else:
        pw = ph = int(sys.argv[5])
    zoom = int(sys.argv[6])
    dst = sys.argv[7]
    wa, ha, a = read_png(pa)
    wb, hb, b = read_png(pb)
    gap = 8
    if x + pw > min(wa, wb) or y + ph > min(ha, hb):
        raise SystemExit("(%d,%d) %dx%d 超出圖（%dx%d）" % (x, y, pw, ph, wa, ha))
    ow = (pw * 2 + gap) * zoom
    oh = ph * zoom
    rows = [[(0, 0, 0)] * ow for _ in range(oh)]
    for dy in range(ph):
        for dx in range(pw):
            ca = a[y + dy][x + dx]
            cb = b[y + dy][x + dx]
            for zy in range(zoom):
                row = rows[dy * zoom + zy]
                for zx in range(zoom):
                    row[(dx) * zoom + zx] = ca
                    row[(pw + gap + dx) * zoom + zx] = cb
    write_png(dst, ow, oh, rows)
    print("左＝%s 右＝%s，(%d,%d) 起 %d×%d，放大 %d 倍 → %s"
          % (pa, pb, x, y, pw, ph, zoom, dst))


if __name__ == "__main__":
    main()
