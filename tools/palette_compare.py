#!/usr/bin/env python3
"""列出兩張圖各自用了哪些顏色，並嘗試把它們配對。

    tools/py.sh tools/palette_compare.py 原版.png remake.png

對拍全區 FAIL 時，第一個要排除的是**調色盤刻度不同**：
同一個 VGA 色在兩條管線可能寫成 (195,130,32) 與 (196,130,33)。
那種差異每個像素都不一樣，看起來像「整張都錯」，
但實際上版面可能一格都沒偏——先看這一支再決定要不要改版面。
"""
import sys
from collections import Counter

sys.path.insert(0, "tools")
from parity_diff import read_png


def palette(path):
    w, h, px = read_png(path)
    c = Counter(px[y][x] for y in range(h) for x in range(w))
    return w, h, c


def main():
    if len(sys.argv) != 3:
        raise SystemExit("用法: palette_compare.py a.png b.png")
    wa, ha, ca = palette(sys.argv[1])
    wb, hb, cb = palette(sys.argv[2])
    print("%s：%dx%d，%d 色" % (sys.argv[1], wa, ha, len(ca)))
    print("%s：%dx%d，%d 色" % (sys.argv[2], wb, hb, len(cb)))
    common = set(ca) & set(cb)
    print("完全相同的顏色：%d 個" % len(common))
    print()
    print("| 原版色 | 次數 | 最近的 remake 色 | 距離 |")
    print("|---|---:|---|---:|")
    for color, n in ca.most_common(20):
        best = min(cb, key=lambda o: sum(abs(o[i] - color[i]) for i in range(3)))
        d = sum(abs(best[i] - color[i]) for i in range(3))
        print("| %s | %d | %s | %d |" % (str(color), n, str(best), d))


if __name__ == "__main__":
    main()
