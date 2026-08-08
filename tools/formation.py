#!/usr/bin/env python3
"""十六種陣形表的檢視器。

說明書 4.3 說陣形變更「共 16 種陣形」，而 `sub_1AA2C` 算兵的目標座標時
用的就是一張表：

    bx = 兵編號 × 2 + 陣形編號 × 96
    dx = cs:[bx - 0x331C]      ; dl = X 位移、dh = Y 位移（有號）

96 ＝ 48 個兵 × 2 byte，所以**一個陣形是 48 組 (dx, dy)**。
另一側取同一張表但 `neg dl`——**陣形左右鏡射**。

    tools/formation.py list           # 16 種陣形的尺寸
    tools/formation.py show 5         # 用 ASCII 畫一個陣形
    tools/formation.py show 5 --raw   # 連座標一起印

表在 `KI.EXE` 的 seg000 偏移 `0xCCE4`（檔案偏移 `0xCEE4`）。
⚠ 位址只對 `dosv` 版；`pc98` 是另一次編譯，要另外找。

原版資產不隨本專案散布，預設從 workplace/orig/dosv/ 讀。
"""

import argparse
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
ORIG = ROOT / "workplace" / "orig" / "dosv"

# seg000 偏移 → 檔案偏移是 +0x200（MZ 表頭 32 段 × 16）。
TABLE_OFF = 0xCEE4
NUM_FORMATIONS = 16
SOLDIERS = 48            # 一側 6 隊 × 8 人
STRIDE = SOLDIERS * 2    # 96，正是 `set1` 指令乘的那個數


def signed(b):
    return b - 256 if b > 127 else b


def load(path):
    d = path.read_bytes()
    out = []
    for f in range(NUM_FORMATIONS):
        o = TABLE_OFF + f * STRIDE
        out.append([(signed(d[o + i * 2]), signed(d[o + i * 2 + 1]))
                    for i in range(SOLDIERS)])
    return out


def cmd_list(forms, _args):
    print("陣形  X 範圍      Y 範圍      寬 × 高   隊長位置")
    for i, f in enumerate(forms):
        xs = [p[0] for p in f]
        ys = [p[1] for p in f]
        w, h = max(xs) - min(xs) + 1, max(ys) - min(ys) + 1
        print(f"{i:3d}   {min(xs):4d}..{max(xs):<4d}  {min(ys):4d}..{max(ys):<4d}"
              f"  {w:3d} × {h:3d}   {f[0]}")


def cmd_show(forms, args):
    f = forms[args.formation]
    xs = [p[0] for p in f]
    ys = [p[1] for p in f]
    x0, x1, y0, y1 = min(xs), max(xs), min(ys), max(ys)
    grid = [[" "] * (x1 - x0 + 1) for _ in range(y1 - y0 + 1)]
    for i, (x, y) in enumerate(f):
        # 每 8 個一隊，隊長是每隊的第一個（sub_1A754 的 1 + 7 結構）。
        ch = str(i // 8) if i % 8 else "★"
        cell = grid[y - y0][x - x0]
        grid[y - y0][x - x0] = ch if cell == " " else "+"
    print(f"陣形 {args.formation}：X {x0}..{x1}　Y {y0}..{y1}"
          f"　（★ ＝ 隊長，數字 ＝ 隊編號，+ ＝ 重疊）")
    print("    " + "".join("|" if x == 0 else "." for x in range(x0, x1 + 1)))
    for y, row in enumerate(grid, start=y0):
        print(f"{y:3d} " + "".join(row))
    if args.raw:
        for i, p in enumerate(f):
            end = "\n" if i % 8 == 7 else "  "
            print(f"{i:2d}{p}", end=end)


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--orig", type=pathlib.Path, default=ORIG)
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("list", help="16 種陣形的尺寸")
    p = sub.add_parser("show", help="用 ASCII 畫一個陣形")
    p.add_argument("formation", type=int)
    p.add_argument("--raw", action="store_true", help="連座標一起印")
    args = ap.parse_args()

    path = args.orig / "KI.EXE"
    if not path.exists():
        sys.exit(f"找不到 {path}（原版資產請自備）")
    forms = load(path)
    if args.cmd == "show" and not 0 <= args.formation < NUM_FORMATIONS:
        sys.exit(f"陣形編號要在 0–{NUM_FORMATIONS - 1}")
    {"list": cmd_list, "show": cmd_show}[args.cmd](forms, args)


if __name__ == "__main__":
    main()
