#!/usr/bin/env python3
"""找出某一區在兩張圖之間差了幾個像素的位移。

    tools/py.sh tools/parity_shift.py 原版.png remake.png [區名] [半徑]

逐區差分只回答「差多少」，不回答「差在哪一種」。**整片 60% 不同**
最常見的成因有兩種，而它們的處置完全相反：

- 卷軸位置差了幾格 → 版面沒錯，對齊就好
- 圖塊解錯 → 對齊也沒用

這一支在 ±半徑 內平移 remake 那一側，回報「不同像素最少」的位移。
最佳位移不是 (0,0) 就是第一種；最佳位移是 (0,0) 而差異仍高就是第二種。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png, REGIONS


def main():
    if len(sys.argv) < 3:
        raise SystemExit("用法: parity_shift.py a.png b.png [區名] [半徑]")
    name = sys.argv[3] if len(sys.argv) > 3 else "map"
    radius = int(sys.argv[4]) if len(sys.argv) > 4 else 24
    wa, ha, a = read_png(sys.argv[1])
    wb, hb, b = read_png(sys.argv[2])
    if (wa, ha) != (wb, hb):
        raise SystemExit("尺寸不同：%dx%d vs %dx%d" % (wa, ha, wb, hb))
    region = next((r for r in REGIONS if r[0] == name), None)
    if region is None:
        raise SystemExit("沒有這一區：%s（有 %s）"
                         % (name, "、".join(r[0] for r in REGIONS)))
    _, rx, ry, rw, rh = region
    # 取區的內縮部分，讓平移之後兩側都還在圖內。
    x0, y0 = rx + radius, ry + radius
    x1, y1 = rx + rw - radius, ry + rh - radius
    if len(sys.argv) > 6:
        # 明確給位移：搜尋半徑搜不到的大位移用這個確認。
        fx, fy = int(sys.argv[5]), int(sys.argv[6])
        n = total = 0
        for y in range(max(ry, -fy), min(ry + rh, ha - fy)):
            ra, rb = a[y], b[y + fy]
            for x in range(max(rx, -fx), min(rx + rw, wa - fx)):
                total += 1
                if ra[x] != rb[x + fx]:
                    n += 1
        print("區 %s 在位移 (%+d, %+d)：不同 %d / %d（%.1f%%）"
              % (name, fx, fy, n, total, 100.0 * n / total))
        # 對齊之後剩下的差異落在哪裡，比「剩多少」更有資訊：
        # 集中成幾塊 ＝ 有東西畫錯位置，散開 ＝ 圖塊本身差一點。
        cell = 32
        grid = {}
        for y in range(max(ry, -fy), min(ry + rh, ha - fy)):
            ra, rb = a[y], b[y + fy]
            for x in range(max(rx, -fx), min(rx + rw, wa - fx)):
                if ra[x] != rb[x + fx]:
                    grid[(x // cell, y // cell)] = grid.get((x // cell, y // cell), 0) + 1
        hot = sorted(grid.items(), key=lambda kv: -kv[1])[:10]
        print("  最集中的 %d 個 %d×%d 方塊：" % (len(hot), cell, cell))
        for (gx, gy), c in hot:
            print("    (%d,%d)-(%d,%d)：%d 點"
                  % (gx * cell, gy * cell, gx * cell + cell - 1, gy * cell + cell - 1, c))
        return
    best = []
    for dy in range(-radius, radius + 1):
        for dx in range(-radius, radius + 1):
            n = 0
            for y in range(y0, y1, 2):          # 取樣兩格一步，夠分辨勝負
                ra, rb = a[y], b[y + dy]
                for x in range(x0, x1, 2):
                    if ra[x] != rb[x + dx]:
                        n += 1
            best.append((n, dx, dy))
    best.sort()
    total = len(range(y0, y1, 2)) * len(range(x0, x1, 2))
    print("區 %s，取樣 %d 點，半徑 ±%d" % (name, total, radius))
    for n, dx, dy in best[:5]:
        print("  位移 (%+d, %+d)：不同 %d（%.1f%%）" % (dx, dy, n, 100.0 * n / total))
    n0 = next(n for n, dx, dy in best if dx == 0 and dy == 0)
    print("  位移 ( 0,  0)：不同 %d（%.1f%%）" % (n0, 100.0 * n0 / total))


if __name__ == "__main__":
    main()
