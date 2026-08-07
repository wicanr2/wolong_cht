#!/usr/bin/env python3
"""BATTLE.MAP／BATTLE.MDL 的戰場結構檢視器。

一格戰場地圖存的是一個圖塊編號，而圖塊本身是**一疊 1–7 層的子圖塊**
（docs/re/11 §4.2）。堆疊高度 ≥ 4 的格子在原版裡會讓站上去的單位
被設一個旗標，而把它畫出來之後，攻城用的戰場會長出一圈封閉的城牆——
這是驗證整個資料模型最便宜的辦法。

    tools/battlefield.py map 5              # 用堆疊高度畫一張戰場
    tools/battlefield.py gates 5            # 列出城門／梯子的格子
    tools/battlefield.py survey             # 214 張的統計，驗證 docs/re/11 §4.4

原版資產不隨本專案散布，預設從 workplace/orig/dosv/ 讀。
"""

import argparse
import collections
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
ORIG = ROOT / "workplace" / "orig" / "dosv"

# docs/formats/07：512 B 索引 ＋ 214 個戰場 × 4,096 B；
# 每個戰場是 64 B 表頭 ＋ 64 × 62 格 ＋ 64 B 尾段。
INDEX_SIZE, FIELD_SIZE, NUM_FIELDS = 512, 4096, 214
CELLS_OFF, WIDTH, HEIGHT = 0x40, 64, 62

# docs/re/11 §4.1：表頭 4,096 B ＋ 3 個圖塊組；每組前 2,048 B 是圖塊定義。
MDL_HEADER, TILESET_SIZE, NUM_TILESETS = 4096, 63488, 3

# 頂層子圖塊落在這個範圍的格子會生成物件（城門／梯子，docs/re/11 §4.5）。
OBJECT_LO, OBJECT_HI = 0xBA, 0xBF

# 站上去會被設旗標的堆疊高度（原版 `cmp al, 4`）。
TALL = 4


class Battle:
    def __init__(self, mapfile, mdlfile):
        self.map = mapfile.read_bytes()
        self.mdl = mdlfile.read_bytes()
        # 每個圖塊組 256 筆 × 8 B 的定義
        self.defs = [
            [
                self.mdl[MDL_HEADER + t * TILESET_SIZE + i * 8:
                         MDL_HEADER + t * TILESET_SIZE + i * 8 + 8]
                for i in range(256)
            ]
            for t in range(NUM_TILESETS)
        ]

    def index(self, n):
        """回傳 (圖塊組編號, 索引第二欄)。第二欄是城門附近的 X。"""
        return self.map[n * 2], self.map[n * 2 + 1]

    def cells(self, n):
        off = INDEX_SIZE + n * FIELD_SIZE + CELLS_OFF
        return self.map[off:off + WIDTH * HEIGHT]

    def stack(self, tileset, tile):
        """回傳這個圖塊由下往上的子圖塊串。"""
        rec = self.defs[tileset][tile]
        return rec[1:rec[0] + 1]

    def heights(self, n):
        t, _ = self.index(n)
        return [self.defs[t][c][0] for c in self.cells(n)]

    def objects(self, n):
        """回傳 [(x, y, 頂層子圖塊)]，也就是城門與梯子的格子。"""
        t, _ = self.index(n)
        out = []
        for i, c in enumerate(self.cells(n)):
            s = self.stack(t, c)
            if s and OBJECT_LO <= s[-1] <= OBJECT_HI:
                out.append((i % WIDTH, i // WIDTH, s[-1]))
        return out


def cmd_map(b, args):
    t, gate = b.index(args.field)
    h = b.heights(args.field)
    tall = sum(1 for v in h if v >= TALL)
    print(f"戰場 {args.field}：圖塊組 {t}，索引第二欄 {gate}"
          f"（{'野戰用' if gate == 0 else '攻城用'}），"
          f"堆疊 ≥{TALL} 的格 {tall}/{WIDTH * HEIGHT}")
    objs = {(x, y) for x, y, _ in b.objects(args.field)}
    for row in range(HEIGHT):
        line = []
        for x in range(WIDTH):
            v = h[row * WIDTH + x]
            if (x, row) in objs:
                line.append("O")           # 城門／梯子
            elif v >= TALL:
                line.append("#")           # 高處（攻城圖上是城牆）
            elif v:
                line.append(".")
            else:
                line.append(" ")           # 空圖塊
        print(f"{row:2d} {''.join(line)}")


def cmd_gates(b, args):
    objs = b.objects(args.field)
    _, gate = b.index(args.field)
    print(f"戰場 {args.field}：{len(objs)} 個物件格，索引第二欄 X ＝ {gate}")
    for x, y, v in objs:
        print(f"  ({x:2d}, {y:2d})  子圖塊 0x{v:02X}")


def cmd_survey(b, _args):
    """驗證 docs/re/11 §4.4 的三條統計。"""
    walled = zero_col = zero_walled = 0
    field_maps_nonzero = []
    dist = collections.Counter()
    for n in range(NUM_FIELDS):
        _, gate = b.index(n)
        h = b.heights(n)
        has_wall = any(v >= TALL for v in h)
        walled += has_wall
        if gate == 0:
            zero_col += 1
            zero_walled += has_wall
        else:
            if not has_wall:
                print(f"  ⚠ 戰場 {n} 第二欄非 0 卻沒有高處")
            xs = {x for x, _, _ in b.objects(n)}
            dist[min((abs(gate - x) for x in xs), default=99)] += 1
        if n >= 192 and gate != 0:
            field_maps_nonzero.append(n)

    print(f"214 張戰場：有高處 {walled}、沒有 {NUM_FIELDS - walled}")
    print(f"第二欄為 0 的 {zero_col} 張，其中有高處 {zero_walled} 張")
    print(f"戰場 192–213（野戰用）第二欄非 0 的：{field_maps_nonzero or '無'}")
    print("第二欄與最近物件格的 X 距離分佈：", dict(sorted(dist.items())))


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--orig", type=pathlib.Path, default=ORIG,
                    help="原版資產目錄（預設 workplace/orig/dosv）")
    sub = ap.add_subparsers(dest="cmd", required=True)
    p = sub.add_parser("map", help="用堆疊高度畫一張戰場")
    p.add_argument("field", type=int)
    p = sub.add_parser("gates", help="列出城門／梯子的格子")
    p.add_argument("field", type=int)
    sub.add_parser("survey", help="214 張的統計")
    args = ap.parse_args()

    mapfile, mdlfile = args.orig / "BATTLE.MAP", args.orig / "BATTLE.MDL"
    if not mapfile.exists() or not mdlfile.exists():
        sys.exit(f"找不到 {mapfile} 或 {mdlfile}（原版資產請自備）")
    b = Battle(mapfile, mdlfile)

    if getattr(args, "field", 0) not in range(NUM_FIELDS) and args.cmd != "survey":
        sys.exit(f"戰場編號要在 0–{NUM_FIELDS - 1}")
    {"map": cmd_map, "gates": cmd_gates, "survey": cmd_survey}[args.cmd](b, args)


if __name__ == "__main__":
    main()
