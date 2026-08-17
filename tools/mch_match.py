#!/usr/bin/env python3
"""在 `MMAP.MCH` 的 256 張圖塊裡找出「疊上去之後會變成原版那一格」的那一張。

    tools/py.sh tools/mch_match.py MMAP.MCH 原版.png remake.png x y [調色盤.BRG] [組]

大地圖的每一格除了地形圖塊，還可以再疊最多五張 MCH 圖塊
（`sub_1D66A` 的 `[si+3]`–`[si+7]`，`docs/re/67` §4）。MCH 圖塊自帶遮罩，
所以疊完的結果 ＝ **遮罩處用 MCH 的顏色、其餘保留底下的地形**。

這一支就照這個規則把 256 張各合成一次，跟原版截圖的同一格比：
完全相同的那一張就是答案。**比的是合成結果不是圖塊本身**——
只比圖塊會被透明區騙過去。

MCH 圖塊 160 B ＝ 32 B 遮罩 ＋ 4 × 32 B 色平面（`internal/assets/world/mmapmch.go`）。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png

TILE = 16
TILE_BYTES = 160
TILE_COUNT = 256
MASK_BYTES = 32


def mch_tile(raw, idx):
    """回傳 16×16 的 (調色盤索引 or None)。None ＝ 透明。"""
    base = idx * TILE_BYTES
    out = []
    for y in range(TILE):
        row = []
        for x in range(TILE):
            bit = 7 - x % 8
            if (raw[base + y * 2 + x // 8] >> bit) & 1 == 0:
                row.append(None)
                continue
            colour = 0
            for plane in range(4):
                b = raw[base + MASK_BYTES + plane * MASK_BYTES + y * 2 + x // 8]
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
    if len(sys.argv) < 6:
        raise SystemExit("用法: mch_match.py MMAP.MCH 原版.png remake.png x y "
                         "[GAMEPAL.BRG] [組]")
    mch = open(sys.argv[1], "rb").read()
    _, _, orig = read_png(sys.argv[2])
    _, _, base = read_png(sys.argv[3])
    x0, y0 = int(sys.argv[4]), int(sys.argv[5])
    brg_path = sys.argv[6] if len(sys.argv) > 6 else "workplace/orig/dosv/GAMEPAL.BRG"
    bank = int(sys.argv[7]) if len(sys.argv) > 7 else 0
    pal = bank_colours(open(brg_path, "rb").read(), bank)

    want = [[orig[y0 + dy][x0 + dx] for dx in range(TILE)] for dy in range(TILE)]
    under = [[base[y0 + dy][x0 + dx] for dx in range(TILE)] for dy in range(TILE)]

    scored = []
    for idx in range(TILE_COUNT):
        t = mch_tile(mch, idx)
        bad = 0
        for dy in range(TILE):
            for dx in range(TILE):
                c = t[dy][dx]
                got = under[dy][dx] if c is None else pal[c]
                if got != want[dy][dx]:
                    bad += 1
        scored.append((bad, idx))
    scored.sort()
    print("底圖 %s (%d,%d)，目標 %s" % (sys.argv[3], x0, y0, sys.argv[2]))
    for bad, idx in scored[:8]:
        print("  MCH 圖塊 %3d：不同 %d / 256" % (idx, bad))


if __name__ == "__main__":
    main()
