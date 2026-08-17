#!/usr/bin/env python3
"""在另一張圖裡找出某一小塊出現在哪。

    tools/py.sh tools/parity_locate.py 樣板.png x y 邊長 目標.png

用途是量「兩邊的卷軸差多少」。逐區差分只說差多少，
`parity_shift.py` 只搜得動小半徑；差到整片畫面之外時，
直接拿一塊認得出來的地標去目標圖裡找，比擴大搜尋半徑便宜得多。

回報**完全相同**的位置（找不到就回報最接近的），
因為調色盤對上之後，同一個圖塊應該逐像素相同——
「只找得到近似位置」本身就是一條結論：圖塊沒有畫成一樣。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png


def main():
    if len(sys.argv) != 6:
        raise SystemExit("用法: parity_locate.py 樣板.png x y 邊長 目標.png")
    tpl, x, y, side, dst = (sys.argv[1], int(sys.argv[2]), int(sys.argv[3]),
                            int(sys.argv[4]), sys.argv[5])
    _, _, a = read_png(tpl)
    w, h, b = read_png(dst)
    patch = [row[x:x + side] for row in a[y:y + side]]
    exact, best, bestn = [], None, side * side + 1
    for oy in range(h - side + 1):
        for ox in range(w - side + 1):
            n = 0
            for dy in range(side):
                rb, rp = b[oy + dy], patch[dy]
                for dx in range(side):
                    if rb[ox + dx] != rp[dx]:
                        n += 1
                        if n >= bestn:
                            break
                if n >= bestn:
                    break
            if n == 0:
                exact.append((ox, oy))
            if n < bestn:
                bestn, best = n, (ox, oy)
    print("樣板 %s (%d,%d) %d×%d" % (tpl, x, y, side, side))
    if exact:
        print("在 %s 找到 %d 個完全相同的位置：%s"
              % (dst, len(exact), "、".join("(%d,%d)" % p for p in exact[:8])))
        for ox, oy in exact[:8]:
            print("  位移 (%+d, %+d)" % (ox - x, oy - y))
    else:
        print("在 %s 沒有完全相同的位置；最接近的是 (%d,%d)，差 %d/%d 點"
              % (dst, best[0], best[1], bestn, side * side))
        print("  位移 (%+d, %+d)" % (best[0] - x, best[1] - y))


if __name__ == "__main__":
    main()
