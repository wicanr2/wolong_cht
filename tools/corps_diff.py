#!/usr/bin/env python3
"""逐 byte 比兩份 SAVE.DAT 的軍團表（規格：`docs/spec/172`）。

    tools/py.sh tools/corps_diff.py 原版.DAT remake.DAT [--only 35]

⭐ **只比「兩邊都活著」的軍團。** 一支只在單邊存在時整筆 64 byte 都會不同，
那是「編成時機不同」不是「欄位對不上」，混在一起看會把後者淹掉。
"""

import sys

BASE, SIZE, N = 0x22C0, 64, 127
# 已知語意的欄位，印出來給人看（docs/spec/172 §1）。
NAMES = {
    0x00: "旗標", 0x0A: "步進", 0x0B: "計時", 0x0C: "路徑指標",
    0x0E: "現在節點", 0x10: "X", 0x12: "Y", 0x14: "目標X", 0x16: "目標Y",
    0x18: "目標節點", 0x1A: "佔用偏移", 0x1C: "佔用段", 0x1E: "間隔",
    0x20: "意圖", 0x23: "階段",
}


def main():
    args = sys.argv[1:]
    only = None
    if "--only" in args:
        i = args.index("--only")
        only = int(args[i + 1])
        del args[i:i + 2]
    if len(args) < 2:
        print(__doc__)
        return 2
    a = open(args[0], "rb").read()
    b = open(args[1], "rb").read()
    cols, corps, single = 0, 0, 0
    for i in range(N):
        if only is not None and i != only:
            continue
        ra = a[BASE + i * SIZE:BASE + (i + 1) * SIZE]
        rb = b[BASE + i * SIZE:BASE + (i + 1) * SIZE]
        alive_a, alive_b = ra[0] & 0xC0, rb[0] & 0xC0
        if not alive_a and not alive_b:
            continue
        if bool(alive_a) != bool(alive_b):
            single += 1
            print(f"軍團 {i:3d} 只在{'原版' if alive_a else 'remake'}存在")
            continue
        d = [k for k in range(SIZE) if ra[k] != rb[k]]
        if not d:
            continue
        corps += 1
        cols += len(d)
        print(f"軍團 {i:3d}：" + "  ".join(
            f"+0x{k:02X}{'(' + NAMES[k] + ')' if k in NAMES else ''} "
            f"{ra[k]:02X}→{rb[k]:02X}" for k in d))
    print(f"\n合計 {cols} 個 byte／{corps} 支；只在單邊存在 {single} 支")
    return 0


if __name__ == "__main__":
    sys.exit(main())
