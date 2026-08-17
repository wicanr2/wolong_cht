#!/usr/bin/env python3
"""列出 `MMAP.MCH` 的 64 個物件矩陣：尺寸與哪幾格有圖塊。

    tools/py.sh tools/mch_patterns.py workplace/orig/dosv/MMAP.MCH

metadata 從 `0xA000` 起，每項 4 byte ＝ 寬、高、矩陣位移（LE，相對 `0xA100`）；
矩陣裡 `0xFF` 表示那一格不畫（`internal/assets/world/mmapmch.go`）。

用途是回答「畫面上這幾格一起變的東西，是不是同一個物件矩陣」——
逐張圖塊比會失敗，因為物件是**一組**圖塊，中間大多是不畫的格子。
"""
import sys

TABLE = 0xA000
DATA = 0xA100
ENTRY = 4
ENTRIES = 0x100 // ENTRY


def main():
    if len(sys.argv) != 2:
        raise SystemExit("用法: mch_patterns.py MMAP.MCH")
    raw = open(sys.argv[1], "rb").read()
    for i in range(ENTRIES):
        e = TABLE + i * ENTRY
        w, h = raw[e], raw[e + 1]
        off = raw[e + 2] | raw[e + 3] << 8
        start = DATA + off
        if w == 0 or h == 0 or start + w * h > len(raw):
            print("%2d  —（寬 %d 高 %d 位移 %d，超出範圍）" % (i, w, h, off))
            continue
        tiles = raw[start:start + w * h]
        drawn = sum(1 for t in tiles if t != 0xFF)
        print("%2d  %d×%d  位移 %5d  有圖塊 %d/%d" % (i, w, h, off, drawn, w * h))
        for y in range(h):
            row = tiles[y * w:(y + 1) * w]
            print("      " + " ".join("  ." if t == 0xFF else "%3d" % t for t in row))


if __name__ == "__main__":
    main()
