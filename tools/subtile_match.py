#!/usr/bin/env python3
"""在戰場的 192 個子圖塊裡找出「畫面上那幾個差異點該是哪一張」。

    tools/py.sh tools/subtile_match.py <戰場編號> <原版png> <remake png> x0 x1 y0 y1

逐區對拍剩下一小群像素、而「多畫一塊」與「少畫一層」長得一樣時，
與其去查繪圖幾何，不如**把每一張候選圖套進去比一次**：
192 張 × 所有可能的落點，留下「蓋住全部差異點而且顏色都對」的那些。
答案落在唯一一張上的話，連帶就指出它是哪個圖塊、哪一層、哪一格
（配 `tools/slot_cells.py` 反解顯示格）。

`docs/playtest/40` §13.1 用手工版定位到破門的柱子（子圖塊 179），
§13.2 用這一支再問一次，同一張圖出現在戰場的另一端。

⚠ 落點是**連續的 16 × 32**，而 `internal/ui/isoview` 的顯示格會把
四段 16 × 8 重排成 32 × 16（`unfoldDisplayTile`）。差異點落在同一段
之內時兩者只差一個平移，超過一段就要自己換算。
"""
import pathlib
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png  # noqa: E402

ROOT = pathlib.Path(__file__).resolve().parent.parent
ORIG = ROOT / "workplace" / "orig" / "dosv"

# docs/formats/07：512 B 索引 ＋ 214 個戰場 × 4,096 B。
INDEX_SIZE, FIELD_SIZE = 512, 4096
# docs/re/11 §4.1：表頭 4,096 B ＋ 3 個圖塊組；每組前 2,048 B 是圖塊定義，
# 後面是 192 個 320 B 的子圖塊。
MDL_HEADER, TILESET_SIZE = 4096, 63488
SUBTILE_BASE, SUBTILE_SIZE, NUM_SUBTILES = 2048, 320, 192
SUB_W, SUB_H, PLANE = 16, 32, 64
TRANSPARENT = -1


def bank0(path):
    """`GAMEPAL.BRG` 的第 0 組，回傳 16 個 4 bit 的 (R, G, B)。"""
    raw = path.read_bytes()
    return [(raw[i * 3 + 1], raw[i * 3 + 2], raw[i * 3]) for i in range(16)]


def decode(b):
    """320 B → 16 × 32 的色號陣列，−1 是透明。

    五個 64 B 位元平面，第一個是遮罩，其餘四個是 4bpp 的色號
    （`internal/assets/battle` 的 `decodePlanar`）。
    """
    pix = [[TRANSPARENT] * SUB_W for _ in range(SUB_H)]
    for y in range(SUB_H):
        for x in range(SUB_W):
            i = y * 2 + x // 8
            bit = 7 - x % 8
            if not (b[i] >> bit) & 1:
                continue
            v = 0
            for p in range(4):
                v |= ((b[PLANE * (1 + p) + i] >> bit) & 1) << p
            pix[y][x] = v
    return pix


def subtiles(field):
    """第 field 張戰場所屬圖塊組的 192 個子圖塊。"""
    bmap = (ORIG / "BATTLE.MAP").read_bytes()
    mdl = (ORIG / "BATTLE.MDL").read_bytes()
    base = MDL_HEADER + bmap[field * 2] * TILESET_SIZE + SUBTILE_BASE
    return [decode(mdl[base + k * SUBTILE_SIZE: base + (k + 1) * SUBTILE_SIZE])
            for k in range(NUM_SUBTILES)]


def colour_map(images):
    """色號 → 螢幕 RGB。截圖經過 DAC 縮放，所以拿畫面上出現過的顏色配對。"""
    seen = set()
    for a in images:
        seen |= {a[y][x] for y in range(len(a)) for x in range(len(a[0]))}
    out = {}
    for i, (r4, g4, b4) in enumerate(bank0(ORIG / "GAMEPAL.BRG")):
        want = (r4 * 17, g4 * 17, b4 * 17)
        out[i] = min(seen, key=lambda c: sum((a - b) ** 2
                                             for a, b in zip(c, want)))
    return out


def main():
    field = int(sys.argv[1])
    orig_png, remake_png = sys.argv[2], sys.argv[3]
    x0, x1, y0, y1 = (int(v) for v in sys.argv[4:8])

    subs = subtiles(field)
    _, _, o = read_png(orig_png)
    _, _, r = read_png(remake_png)
    idx2rgb = colour_map((o, r))

    diffs = [(x, y) for y in range(y0, y1 + 1) for x in range(x0, x1 + 1)
             if o[y][x] != r[y][x]]
    print(f"差異點 {len(diffs)} 個：{diffs}")

    hits = []
    for k, pix in enumerate(subs):
        for py in range(y1 - SUB_H + 1, y0 + 1):
            for px in range(x1 - SUB_W + 1, x0 + 1):
                if any(pix[dy - py][dx - px] == TRANSPARENT or
                       idx2rgb[pix[dy - py][dx - px]] != o[dy][dx]
                       for dx, dy in diffs):
                    continue
                wrong = 0
                for yy in range(SUB_H):
                    for xx in range(SUB_W):
                        c = pix[yy][xx]
                        sy, sx = py + yy, px + xx
                        if c == TRANSPARENT or not (
                                0 <= sy < len(o) and 0 <= sx < len(o[0])):
                            continue
                        if idx2rgb[c] != o[sy][sx]:
                            wrong += 1
                hits.append((wrong, k, px, py))

    hits.sort()
    print(f"\n蓋住全部差異點的候選：{len(hits)} 組")
    for wrong, k, px, py in hits[:25]:
        print(f"  子圖塊 {k:3d}  左上 ({px},{py})  足跡內畫錯 {wrong:3d} px")
    if not hits:
        print("  沒有。差異點可能不是地形，或跨了顯示格的兩段。")


if __name__ == "__main__":
    main()
