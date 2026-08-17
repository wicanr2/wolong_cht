#!/usr/bin/env python3
"""在 `MMAP.MDL` 的 256 張地形圖塊裡找出「畫面上那一格畫的是哪一張」。

    tools/py.sh tools/mdl_match.py MMAP.MDL 截圖.png x y [GAMEPAL.BRG] [組]

`MMAP.MCH` 的物件圖塊帶遮罩、要疊在底圖上（`tools/mch_match.py`）；
`MMAP.MDL` 的地形圖塊是**不透明**的 128 B ＝ 4 個色平面，直接比就好。

用途：畫面上某一格與 `MMAP.MAP` 寫的圖塊編號對不起來時，
先問「原版畫的是這一組裡的哪一張」——如果是其中一張，
那就是**換了圖塊編號**，不是疊了東西上去。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png

TILE = 16
TILE_BYTES = 128
TILE_COUNT = 256
PLANE = 32


def mdl_tile(raw, idx):
    base = idx * TILE_BYTES
    out = []
    for y in range(TILE):
        row = []
        for x in range(TILE):
            bit = 7 - x % 8
            colour = 0
            for plane in range(4):
                b = raw[base + plane * PLANE + y * 2 + x // 8]
                colour |= ((b >> bit) & 1) << plane
            row.append(colour)
        out.append(row)
    return out


def dac(v):
    d = v * 4
    return (d << 2) | (d >> 4)


def bank_colours(brg, bank):
    off = bank * 16 * 3
    return [(dac(brg[off + i * 3 + 1]), dac(brg[off + i * 3 + 2]), dac(brg[off + i * 3]))
            for i in range(16)]


def main():
    if len(sys.argv) < 5:
        raise SystemExit("用法: mdl_match.py MMAP.MDL 截圖.png x y [GAMEPAL.BRG] [組]")
    raw = open(sys.argv[1], "rb").read()
    _, _, img = read_png(sys.argv[2])
    x0, y0 = int(sys.argv[3]), int(sys.argv[4])
    brg = sys.argv[5] if len(sys.argv) > 5 else "workplace/orig/dosv/GAMEPAL.BRG"
    bank = int(sys.argv[6]) if len(sys.argv) > 6 else 0
    pal = bank_colours(open(brg, "rb").read(), bank)
    want = [[img[y0 + dy][x0 + dx] for dx in range(TILE)] for dy in range(TILE)]

    scored = []
    for idx in range(TILE_COUNT):
        t = mdl_tile(raw, idx)
        bad = sum(1 for dy in range(TILE) for dx in range(TILE)
                  if pal[t[dy][dx]] != want[dy][dx])
        scored.append((bad, idx))
    scored.sort()
    print("%s (%d,%d) 最像的圖塊：" % (sys.argv[2], x0, y0))
    for bad, idx in scored[:8]:
        print("  圖塊 %3d：不同 %d / 256" % (idx, bad))


if __name__ == "__main__":
    main()
