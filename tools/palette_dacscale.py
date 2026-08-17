#!/usr/bin/env python3
"""驗證「DOS/V 的 4 bit 通道值走 VGA 6 bit DAC」這條算式。

    tools/py.sh tools/palette_dacscale.py workplace/orig/dosv/GAMEPAL.BRG 原版截圖.png

`docs/re/02` §5：檔案裡的通道值 `v` 是 4 bit，DOS/V 額外 `shl ah,1` 兩次
才寫進 VGA DAC——也就是 DAC 值 ＝ `4v`，**上限 60 不是 63**。
所以整條類比輸出永遠到不了滿刻度，白色是 #F3F3F3 不是 #FFFFFF。

這一支把三種換算各算一次，再看哪一種的顏色集合真的出現在實機截圖裡：

    v*255/15         4 bit 直接鋪滿 0–255（PC-98 的類比調色盤是這樣）
    (4v<<2)|(4v>>4)  VGA 的 6→8 bit 位元複製
    round(4v*255/63) 類比刻度的四捨五入

判準是**實機截圖裡有沒有那個顏色**，不是哪個公式比較漂亮。
"""
import sys
from collections import Counter

sys.path.insert(0, "tools")
from parity_diff import read_png

BANK = 16
BPC = 3


def variants(v):
    d = v * 4                       # DOS/V：shl ah,1 兩次
    return {
        "v*255/15": v * 0xFF // 0x0F,
        "6→8 位元複製": (d << 2) | (d >> 4),
        "類比四捨五入": (d * 255 + 31) // 63,
    }


def main():
    if len(sys.argv) != 3:
        raise SystemExit("用法: palette_dacscale.py GAMEPAL.BRG 原版截圖.png")
    data = open(sys.argv[1], "rb").read()
    w, h, px = read_png(sys.argv[2])
    seen = Counter(px[y][x] for y in range(h) for x in range(w))

    print("截圖用到 %d 色。" % len(seen))
    names = list(variants(0))
    for bank in range(min(4, len(data) // (BANK * BPC))):
        off = bank * BANK * BPC
        hit = {n: 0 for n in names}
        for i in range(BANK):
            c = data[off + i * BPC: off + i * BPC + BPC]
            for n in names:
                rgb = (variants(c[1])[n], variants(c[2])[n], variants(c[0])[n])
                if rgb in seen:
                    hit[n] += 1
        print("第 %d 組：%s" % (bank, "  ".join("%s %d/16" % (n, hit[n]) for n in names)))

    print()
    print("| 4 bit v | %s |" % " | ".join(names))
    print("|---:|%s" % ("---:|" * len(names)))
    for v in range(16):
        t = variants(v)
        print("| %d | %s |" % (v, " | ".join(str(t[n]) for n in names)))


if __name__ == "__main__":
    main()
