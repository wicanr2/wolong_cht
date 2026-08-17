#!/usr/bin/env python3
"""印出 `MMAP.MAP` 某一塊的圖塊編號。

    tools/py.sh tools/map_tile.py workplace/orig/dosv/MMAP.MAP 200 110 8 6

用途是回答「畫面上那一格是幾號圖塊」。從截圖只看得到**畫出來的樣子**，
看不到編號；而要判斷「原版是不是換了圖塊」就得先知道檔案裡寫的是幾號。

RLE 與 `internal/assets/rle` 同一支（`KI.EXE` 的 `sub_1F5E7`）：
沒有逃脫字元，連續兩個相同的 byte 之後那一個 byte 是「再重複幾次」。
"""
import sys

WIDTH, HEIGHT = 384, 256


def decode(src):
    out = bytearray()
    i = 0
    n = len(src)
    while i < n:
        prev = src[i]
        out.append(prev)
        i += 1
        matched = False
        while i < n:
            cur = src[i]
            out.append(cur)
            i += 1
            if cur == prev:
                matched = True
                break
            prev = cur
        if not matched or i >= n:
            break
        count = src[i]
        i += 1
        out.extend(bytes([prev]) * count)
    return bytes(out)


def main():
    if len(sys.argv) not in (4, 6):
        raise SystemExit("用法: map_tile.py MMAP.MAP x y [寬 高]")
    path = sys.argv[1]
    x0, y0 = int(sys.argv[2]), int(sys.argv[3])
    w = int(sys.argv[4]) if len(sys.argv) == 6 else 8
    h = int(sys.argv[5]) if len(sys.argv) == 6 else 8
    raw = decode(open(path, "rb").read())
    print("解出 %d B（地圖本體 %d B，尾巴 %d B）"
          % (len(raw), WIDTH * HEIGHT, len(raw) - WIDTH * HEIGHT))
    print("    " + " ".join("%3d" % (x0 + dx) for dx in range(w)))
    for dy in range(h):
        y = y0 + dy
        row = [raw[y * WIDTH + x0 + dx] for dx in range(w)]
        print("%3d " % y + " ".join("%3d" % v for v in row))


if __name__ == "__main__":
    main()
