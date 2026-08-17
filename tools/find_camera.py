#!/usr/bin/env python3
"""從一張原版截圖反推大地圖的鏡頭在哪一格。

    tools/py.sh tools/find_camera.py 截圖.png MMAP.MDL MMAP.MAP [欄 列 邊長]

作法：把畫面上某一小塊逐格認成 `MMAP.MDL` 的圖塊編號，
再拿那個編號矩陣去 `MMAP.MAP` 裡找唯一的落點。

⭐ **比「用眼睛認地形」可靠**：地圖上長得像的地方很多，
但一個 3×3 的圖塊編號矩陣通常只出現一次。找到多個落點就會印出來，
不會挑一個當答案。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png
from mdl_match import mdl_tile, bank_colours
from map_tile import decode

WIDTH, HEIGHT = 384, 256
TILE = 16
BANNER = 32


def main():
    if len(sys.argv) < 4:
        raise SystemExit("用法: find_camera.py 截圖.png MMAP.MDL MMAP.MAP [欄 列 邊長]")
    _, _, img = read_png(sys.argv[1])
    mdl = open(sys.argv[2], "rb").read()
    raw = decode(open(sys.argv[3], "rb").read())
    col0 = int(sys.argv[4]) if len(sys.argv) > 4 else 2
    row0 = int(sys.argv[5]) if len(sys.argv) > 5 else 4
    side = int(sys.argv[6]) if len(sys.argv) > 6 else 3
    pal = bank_colours(open("workplace/orig/dosv/GAMEPAL.BRG", "rb").read(), 0)
    tiles = [[[pal[c] for c in row] for row in mdl_tile(mdl, i)] for i in range(256)]

    want = []
    for r in range(side):
        line = []
        for c in range(side):
            sx, sy = (col0 + c) * TILE, BANNER + (row0 + r) * TILE
            best, bad = min(
                ((i, sum(1 for y in range(TILE) for x in range(TILE)
                         if tiles[i][y][x] != img[sy + y][sx + x]))
                 for i in range(256)), key=lambda t: t[1])
            if bad:
                raise SystemExit("螢幕格 (%d,%d) 認不出圖塊（最接近 %d，差 %d）"
                                 % (col0 + c, row0 + r, best, bad))
            line.append(best)
        want.append(line)
    print("認出來的圖塊矩陣：")
    for line in want:
        print("   " + " ".join("%3d" % v for v in line))

    hits = []
    for y in range(HEIGHT - side + 1):
        for x in range(WIDTH - side + 1):
            if all(raw[(y + r) * WIDTH + x + c] == want[r][c]
                   for r in range(side) for c in range(side)):
                hits.append((x - col0, y - row0))
    if not hits:
        print("在地圖上找不到這個矩陣")
        return
    print("鏡頭候選（%d 個）：" % len(hits))
    for cx, cy in hits[:8]:
        print("   camX=%d camY=%d" % (cx, cy))


if __name__ == "__main__":
    main()
