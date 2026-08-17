#!/usr/bin/env python3
"""在原版截圖裡量出一塊純色矩形（按鈕）的邊界。

    tools/py.sh tools/find_button.py 圖.png x0 y0 x1 y1

用途是把「按鈕在哪」變成量出來的數字，而不是照著截圖用眼睛估。
先在指定範圍內找出現次數最多的顏色，再取那個顏色的連通範圍。
"""
import sys
from collections import Counter

sys.path.insert(0, "tools")
from parity_diff import read_png


def compress(values):
    """把 [1,2,3,7,8] 印成 '1..3, 7..8'。"""
    if not values:
        return "（無）"
    out, start, prev = [], values[0], values[0]
    for v in values[1:]:
        if v != prev + 1:
            out.append("%d..%d" % (start, prev))
            start = v
        prev = v
    out.append("%d..%d" % (start, prev))
    return ", ".join(out)


def main():
    if len(sys.argv) != 6:
        raise SystemExit("用法: find_button.py 圖.png x0 y0 x1 y1")
    src = sys.argv[1]
    x0, y0, x1, y1 = (int(v) for v in sys.argv[2:6])
    w, h, px = read_png(src)
    x1 = min(x1, w - 1)
    y1 = min(y1, h - 1)
    counts = Counter()
    for y in range(y0, y1 + 1):
        for x in range(x0, x1 + 1):
            counts[px[y][x]] += 1
    top = counts.most_common(1)[0][0]
    print("範圍 (%d,%d)-(%d,%d) 逐列命中最多色 %s 的次數：" % (x0, y0, x1, y1, top))
    for color, _ in counts.most_common(4):
        rows = [(y, sum(1 for x in range(x0, x1 + 1) if px[y][x] == color))
                for y in range(y0, y1 + 1)]
        band = [y for y, n in rows if n > 4]
        print("  %-16s 有 >4 個的列：%s" % (str(color), compress(band)))
    print("範圍 (%d,%d)-(%d,%d) 前六個顏色：" % (x0, y0, x1, y1))
    for color, n in counts.most_common(6):
        xs = [x for y in range(y0, y1 + 1) for x in range(x0, x1 + 1)
              if px[y][x] == color]
        ys = [y for y in range(y0, y1 + 1) for x in range(x0, x1 + 1)
              if px[y][x] == color]
        print("  %-16s %5d 次  x %d..%d  y %d..%d"
              % (str(color), n, min(xs), max(xs), min(ys), max(ys)))


if __name__ == "__main__":
    main()
