#!/usr/bin/env python3
"""量出一張畫面裡水平／垂直邊線的位置（規格：`docs/spec/90-same-state-parity.md`）。

    tools/py.sh tools/parity_edges.py 畫面.png [--top 12] [--min-run 0.5]

## 為什麼需要它

參考來源是**壓縮過的影片幀**時，逐像素差分每一區都會是 ~100%，
那個數字什麼都不告訴你。但**視窗外框的位置壓縮壞不掉**——
框線是整列／整欄的高對比邊界，在梯度能量上是明顯的尖峰。

所以拿它回答「原版的橫幅底在第幾列」這種問題：量到的列號可以直接
與機器碼算出來的值（32／64／192／432）對照，對得上就是幾何相符，
不必要求壓縮過的像素逐點相同。

## 判讀

輸出是「梯度能量最高的列／欄」。**外框是兩條靠在一起的邊**
（框的上緣與下緣），所以一個 8 px 高的框線通常給出相鄰的兩個尖峰；
把它們當成一個邊界看。
"""
import argparse
import sys

sys.path.insert(0, __file__.rsplit("/", 1)[0])
from parity_diff import read_png  # noqa: E402


def luma(px):
    return (px[0] * 299 + px[1] * 587 + px[2] * 114) // 1000


def energy(rows, w, h, axis):
    """axis='y' 回傳每一列與上一列的平均亮度差；'x' 回傳每一欄與左一欄的。

    ⚠ 用**亮度差的大小**，不是「有幾個像素不同」。參考幀壓縮過，
    幾乎每個像素都與鄰居不同，數個數會每一列都是 100%——
    那個指標在這裡是飽和的，看起來有輸出其實沒有訊號。
    """
    lum = [[luma(px) for px in row] for row in rows]
    out = []
    if axis == "y":
        for y in range(1, h):
            s = sum(abs(lum[y][x] - lum[y - 1][x]) for x in range(w))
            out.append((y, s / w / 255))
    else:
        for x in range(1, w):
            s = sum(abs(lum[y][x] - lum[y][x - 1]) for y in range(h))
            out.append((x, s / h / 255))
    return out


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("image")
    ap.add_argument("--top", type=int, default=12)
    ap.add_argument("--min-run", type=float, default=0.10,
                    help="只列出平均亮度差至少這個比例的列／欄")
    ap.add_argument("--rect", help="只看這塊矩形，格式 x,y,w,h。"
                                   "整張畫面掃時大地圖的紋理會蓋過右欄的框線")
    ns = ap.parse_args()

    w, h, rows = read_png(ns.image)
    print("畫面 %dx%d" % (w, h))
    if ns.rect:
        rx, ry, rw, rh = (int(v) for v in ns.rect.split(","))
        rows = [row[rx:rx + rw] for row in rows[ry:ry + rh]]
        w, h = rw, rh
        print("只看 (%d,%d,%d,%d)：底下的列／欄號是**相對這塊矩形**的，"
              "加回 (%d,%d) 才是畫面座標" % (rx, ry, rw, rh, rx, ry))
    for axis, label in (("y", "列"), ("x", "欄")):
        vals = [t for t in energy(rows, w, h, axis) if t[1] >= ns.min_run]
        vals.sort(key=lambda t: -t[1])
        vals = sorted(vals[:ns.top])
        print("\n%s（變化比例 ≥ %.0f%%，取前 %d 名後依位置排序）：" %
              (label, ns.min_run * 100, ns.top))
        print("| %s | 變化比例 |" % label)
        print("|---:|---:|")
        for pos, ratio in vals:
            print("| %d | %.0f%% |" % (pos, ratio * 100))
    return 0


if __name__ == "__main__":
    sys.exit(main())
