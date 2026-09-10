#!/usr/bin/env python3
"""逐 byte 比兩份 SAVE.DAT 的軍團表（規格：`docs/spec/172`）。

    tools/py.sh tools/corps_diff.py 原版.DAT remake.DAT [--only 35]

⭐ **只比「兩邊都活著」的軍團。** 一支只在單邊存在時整筆 64 byte 都會不同，
那是「編成時機不同」不是「欄位對不上」，混在一起看會把後者淹掉。

⭐ **被免除的欄位會逐筆印出來，不會靜靜跳過。** 判一個欄位「remake 不必
建模」等於為它關掉所有檢查——那種決定要跟斷言同級，所以 `EXEMPT` 每一條
都帶出處，而且報表最後會列出這一次實際豁免了幾筆。沒有出處就不要加進去。
"""

import sys

BASE, SIZE, N = 0x22C0, 64, 127
# 已知語意的欄位，印出來給人看（docs/spec/172 §1）。
NAMES = {
    0x00: "旗標", 0x08: "朝向", 0x09: "圖塊基底", 0x0A: "步進",
    0x0B: "計時", 0x0C: "路徑指標",
    0x0E: "現在節點", 0x10: "X", 0x12: "Y", 0x14: "目標節點", 0x16: "目標X",
    0x18: "目標Y", 0x1A: "佔用偏移", 0x1C: "佔用段", 0x1E: "間隔",
    0x20: "意圖", 0x23: "階段",
}


# 已知的「原版有、remake 的規則層結構上不會有」的位元。
#
# ⚠ 加進來的每一條都要寫明**出處**與**為什麼規則層不會有**。
# 這是豁免不是註解——寫進去的那一刻，這個 byte 就不再被任何檢查看到。
EXEMPT = [
    (0x00, 0x10,
     "位元 4 ＝ 繪圖狀態：`sub_12B3C` 開頭 `or byte ptr [si], 10h`，"
     "同一支接著算螢幕座標、用 `[si+3]` 低 2 位當動畫幀，然後 `loc_1D51F` "
     "blit（docs/re/34 §2）。原版跑起來一定會設它，規則層不畫圖就永遠是 0"),
]


def exempt_bits(off):
    """這個位移上被豁免的位元遮罩。"""
    m = 0
    for o, mask, _ in EXEMPT:
        if o == off:
            m |= mask
    return m


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
    waived = []
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
        d = []
        for k in range(SIZE):
            if ra[k] == rb[k]:
                continue
            if m := exempt_bits(k):
                if (ra[k] ^ rb[k]) & ~m & 0xFF == 0:
                    waived.append((i, k, ra[k], rb[k]))
                    continue
            d.append(k)
        if not d:
            continue
        corps += 1
        cols += len(d)
        print(f"軍團 {i:3d}：" + "  ".join(
            f"+0x{k:02X}{'(' + NAMES[k] + ')' if k in NAMES else ''} "
            f"{ra[k]:02X}→{rb[k]:02X}" for k in d))
    print(f"\n合計 {cols} 個 byte／{corps} 支；只在單邊存在 {single} 支")
    if waived:
        print(f"\n豁免 {len(waived)} 筆（差異只落在下列位元上）：")
        for off, mask, why in EXEMPT:
            hit = [w for w in waived if w[1] == off and (w[2] ^ w[3]) & mask]
            if hit:
                print(f"  +0x{off:02X} 遮罩 {mask:02X}h：{', '.join(str(h[0]) for h in hit)}")
                print(f"      {why}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
