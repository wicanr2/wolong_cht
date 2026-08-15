#!/usr/bin/env python3
"""把原版實錄影片的一格對齊成 640×400，並量出版面地標。

    tools/py.sh tools/parity_frame.py workplace/parity/f008.ppm \\
        --crop 0,56,956,661 --out workplace/parity/f008-640x400.ppm

## 為什麼要這一支

remake 的版面常數是從機器碼讀出來的（`docs/re/47`），但一直**沒有
辦法在原版畫面上核對**——模擬器進不到「四窗全開」那個狀態
（`docs/playtest/23` §4.1）。實錄影片繞過了那個閘：它就是那個狀態。

⚠ **這不是逐像素 parity。** 影片是 YouTube 的再編碼，956×605 的內容
是 640×400 縮放上來的（垂直還多拉了 1.1%），色彩經過 4:2:0 取樣。
能比的是**幾何**：邊界落在哪一格。像素級的比對要另外找路。

## 怎麼量

沿著 x 與 y 各算一次「相鄰像素的差」，取行／列平均。版面的框線是
整條的直線，所以會在那一列（欄）形成尖峰；地圖與人物的紋理是散的，
平均之後壓不過框線。**只回報尖峰的位置，不做任何配對**——配對是
呼叫端的事，工具不下結論。
"""
import argparse
import sys


def read_ppm(path):
    d = open(path, "rb").read()
    parts, i = [], 0
    while len(parts) < 4:
        while d[i] in b" \t\r\n":
            i += 1
        j = i
        while d[j] not in b" \t\r\n":
            j += 1
        parts.append(d[i:j])
        i = j
    if parts[0] != b"P6":
        sys.exit("只吃 P6 的 PPM（用 ffmpeg 轉：-pix_fmt rgb24 out.ppm）")
    return int(parts[1]), int(parts[2]), d[i + 1:]


def write_ppm(path, w, h, px):
    with open(path, "wb") as fh:
        fh.write(b"P6\n%d %d\n255\n" % (w, h))
        fh.write(bytes(px))


def resample(w, h, px, box, ow, oh):
    """最近鄰重取樣。**不做內插**——內插會把 1 px 的框線抹掉，
    而框線正是要量的東西。"""
    x0, y0, x1, y1 = box
    out = bytearray(ow * oh * 3)
    for oy in range(oh):
        sy = y0 + (y1 - y0) * oy // oh
        for ox in range(ow):
            sx = x0 + (x1 - x0) * ox // ow
            s = (sy * w + sx) * 3
            o = (oy * ow + ox) * 3
            out[o:o + 3] = px[s:s + 3]
    return out


def edges(w, h, px, axis, lo, hi):
    """回傳每一列（或每一欄）的相鄰差平均。"""
    prof = []
    if axis == "y":
        for y in range(h):
            s = 0
            for x in range(lo, hi):
                a = (y * w + x) * 3
                b = (min(y + 1, h - 1) * w + x) * 3
                s += abs(px[a] - px[b]) + abs(px[a + 1] - px[b + 1]) + abs(px[a + 2] - px[b + 2])
            prof.append(s / max(hi - lo, 1))
    else:
        for x in range(w):
            s = 0
            for y in range(lo, hi):
                a = (y * w + x) * 3
                b = (y * w + min(x + 1, w - 1)) * 3
                s += abs(px[a] - px[b]) + abs(px[a + 1] - px[b + 1]) + abs(px[a + 2] - px[b + 2])
            prof.append(s / max(hi - lo, 1))
    return prof


def peaks(prof, top):
    """取局部極大值，由強到弱。"""
    out = []
    for i in range(1, len(prof) - 1):
        if prof[i] >= prof[i - 1] and prof[i] > prof[i + 1]:
            out.append((prof[i], i))
    out.sort(reverse=True)
    return out[:top]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("ppm")
    ap.add_argument("--crop", required=True, help="x0,y0,x1,y1（含頭不含尾）")
    ap.add_argument("--out", help="輸出重取樣後的 PPM")
    ap.add_argument("--rows", default="40,360", help="量垂直邊界時取樣的列範圍")
    ap.add_argument("--cols", default="8,420", help="量水平邊界時取樣的欄範圍")
    ap.add_argument("--top", type=int, default=8)
    args = ap.parse_args()

    w, h, px = read_ppm(args.ppm)
    box = tuple(int(v) for v in args.crop.split(","))
    out = resample(w, h, px, box, 640, 400)
    if args.out:
        write_ppm(args.out, 640, 400, out)

    r0, r1 = (int(v) for v in args.rows.split(","))
    c0, c1 = (int(v) for v in args.cols.split(","))
    print("輸入 %dx%d，裁切 %s → 640x400" % (w, h, box))
    print("垂直邊界（x，取 y %d–%d）：" % (r0, r1))
    for v, i in peaks(edges(640, 400, out, "x", r0, r1), args.top):
        print("   x=%3d  強度 %.0f" % (i, v))
    print("水平邊界（y，取 x %d–%d）：" % (c0, c1))
    for v, i in peaks(edges(640, 400, out, "y", c0, c1), args.top):
        print("   y=%3d  強度 %.0f" % (i, v))


main()
