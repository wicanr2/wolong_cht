#!/usr/bin/env python3
"""找出兩張畫面之間最合適的整體位移（規格：`docs/spec/90-same-state-parity.md`）。

    tools/py.sh tools/parity_offset.py 原版.png remake.png [--region banner]

## 為什麼需要它

逐區差分只回答「差多少」，不回答「是不是整片平移」。而這兩件事的處置
完全不同：**平移是版面常數錯**（改一個數字），**散佈是內容或調色盤不同**
（要改素材或解碼）。把 −4…+4 的位移全掃一遍，看最小值落在哪裡，
一次就分得開。

最小值落在 (0,0) 而且明顯低於鄰居 → 幾何對齊了，剩下的是顏色或內容。
最小值落在別處 → 那個位移就是版面常數的誤差。
"""
import argparse
import sys

sys.path.insert(0, __file__.rsplit("/", 1)[0])
from parity_diff import REGIONS, read_png  # noqa: E402

SPAN = 4


def region_rect(name):
    for r in REGIONS:
        if r[0] == name:
            return r[1:]
    raise SystemExit("沒有這一區：%s（有 %s）" % (name, ", ".join(r[0] for r in REGIONS)))


def mismatch(a, b, rect, dx, dy, w, h):
    rx, ry, rw, rh = rect
    n = total = 0
    for y in range(ry, min(ry + rh, h)):
        sy = y + dy
        if sy < 0 or sy >= h:
            continue
        for x in range(rx, min(rx + rw, w)):
            sx = x + dx
            if sx < 0 or sx >= w:
                continue
            total += 1
            if a[y][x] != b[sy][sx]:
                n += 1
    return n, total


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("original")
    ap.add_argument("remake")
    ap.add_argument("--region", default="banner")
    ns = ap.parse_args()

    aw, ah, a = read_png(ns.original)
    bw, bh, b = read_png(ns.remake)
    if (aw, ah) != (bw, bh):
        print("✗ 尺寸不同：%dx%d ≠ %dx%d" % (aw, ah, bw, bh))
        return 2
    rect = region_rect(ns.region)

    best = None
    print("區 `%s` 的位移掃描（remake 相對原版）：" % ns.region)
    print("| dy \\ dx | " + " | ".join("%+d" % d for d in range(-SPAN, SPAN + 1)) + " |")
    print("|---|" + "---:|" * (2 * SPAN + 1))
    for dy in range(-SPAN, SPAN + 1):
        cells = []
        for dx in range(-SPAN, SPAN + 1):
            n, total = mismatch(a, b, rect, dx, dy, aw, ah)
            ratio = n / total if total else 1.0
            cells.append("%.1f%%" % (ratio * 100))
            if best is None or ratio < best[0]:
                best = (ratio, dx, dy)
        print("| **%+d** | " % dy + " | ".join(cells) + " |")
    ratio, dx, dy = best
    print("\n最小不同比例 %.2f%% 落在 dx=%+d dy=%+d" % (ratio * 100, dx, dy))
    if (dx, dy) == (0, 0):
        print("→ 幾何對齊；剩下的差異是顏色或內容，不是版面常數。")
    else:
        print("→ 版面常數差 (%+d, %+d)。" % (dx, dy))
    return 0


if __name__ == "__main__":
    sys.exit(main())
