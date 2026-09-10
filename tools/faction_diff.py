#!/usr/bin/env python3
"""逐 byte 比兩份 SAVE.DAT 的勢力表（22 筆 × 64 B，區塊 `+0x0080` 起）。

    tools/py.sh tools/faction_diff.py 原版.DAT remake.DAT [--only 3]
    tools/py.sh tools/faction_diff.py --selftest

⭐ **這一張表原本不在對拍範圍裡。** `city_diff` 比據點、`corps_diff` 比軍團，
勢力表（預備兵、資金、本月支出、軍團數、侵攻目標…）**一格都沒比過**——
於是預備兵的偏差累積了幾千拍才被發現，而症狀出現在完全不相干的地方
（補兵的每槽分配差 4）。

欄位對映出自 `internal/state/state.go` 的寫回端與 `docs/formats/08`。
"""

import sys

BASE, SIZE, N = 0x0080, 64, 22

# 已解語意的欄位。word/i24 的第二、三個 byte 標成續，方便一眼看出是同一格。
NAMES = {
    0x00: "存在旗標", 0x01: "君主", 0x02: "軍師", 0x03: "首都",
    0x04: "預備兵(步)", 0x05: "預備兵(步)續",
    0x06: "預備兵(騎)", 0x07: "預備兵(騎)續",
    0x08: "預備兵(弓)", 0x09: "預備兵(弓)續",
    0x14: "軍團數", 0x16: "求援據點", 0x17: "失守據點", 0x18: "武將數",
    0x19: "侵攻目標",
    0x1A: "本月支出", 0x1B: "本月支出續", 0x1C: "本月支出續",
    0x1D: "士氣基準",
    0x20: "資金", 0x21: "資金續", 0x22: "資金續",
    0x23: "據點數", 0x28: "好戰", 0x2A: "外交官",
}
# word 欄位：印出十進位的值，差 1 與差 256 一眼就分得開。
WORDS = (0x04, 0x06, 0x08)


def diff(a, b, only=None):
    rows, cols, facs = [], 0, 0
    for i in range(N):
        if only is not None and i != only:
            continue
        ra = a[BASE + i * SIZE:BASE + (i + 1) * SIZE]
        rb = b[BASE + i * SIZE:BASE + (i + 1) * SIZE]
        d = [k for k in range(SIZE) if ra[k] != rb[k]]
        if not d:
            continue
        facs += 1
        cols += len(d)
        parts = []
        for k in d:
            name = NAMES.get(k, "")
            tag = f"({name})" if name else ""
            parts.append(f"+0x{k:02X}{tag} {ra[k]:02X}→{rb[k]:02X}")
        rows.append((i, parts, ra, rb))
    return rows, cols, facs


def main():
    args = sys.argv[1:]
    if "--selftest" in args:
        return selftest()
    only = None
    if "--only" in args:
        j = args.index("--only")
        only = int(args[j + 1])
        del args[j:j + 2]
    if len(args) < 2:
        print(__doc__)
        return 2
    a = open(args[0], "rb").read()
    b = open(args[1], "rb").read()
    rows, cols, facs = diff(a, b, only)
    for i, parts, ra, rb in rows:
        print(f"勢力 {i:2d}：" + "  ".join(parts))
        for off in WORDS:
            va = ra[off] | ra[off + 1] << 8
            vb = rb[off] | rb[off + 1] << 8
            if va != vb:
                print(f"          +0x{off:02X} {NAMES[off]}：{va} → {vb}（差 {vb - va:+}）")
    print(f"\n合計 {cols} 個 byte／{facs} 個勢力")
    return 0


def selftest():
    fails = []

    def check(label, cond):
        print(("  ok  " if cond else "  FAIL") + "  " + label)
        if not cond:
            fails.append(label)

    blank = bytearray(BASE + N * SIZE)
    same = bytes(blank)
    check("完全相同 → 0 個 byte", diff(same, same)[1] == 0)

    # 正對照：改一個 byte 就要抓到，而且落在對的那個勢力。
    m = bytearray(same)
    m[BASE + 3 * SIZE + 0x04] = 0x10
    rows, cols, facs = diff(same, bytes(m))
    check("改一個 byte → 抓到 1 個", cols == 1 and facs == 1)
    check("落在對的勢力上", rows and rows[0][0] == 3)

    # word 欄位的值要印成十進位差額（差 1 與差 256 分得開）。
    m2 = bytearray(same)
    m2[BASE + 3 * SIZE + 0x05] = 0x01  # 高 byte +1 ＝ 值 +256
    r2 = diff(same, bytes(m2))
    check("高 byte 的差也算數", r2[1] == 1)

    # 反對照：最後一個勢力的最後一個 byte 也要在範圍內。
    m3 = bytearray(same)
    m3[BASE + (N - 1) * SIZE + SIZE - 1] = 0xFF
    check("掃到最後一個勢力的最後一格", diff(same, bytes(m3))[1] == 1)

    # 反對照：勢力表**之外**的改動不能被算進來。
    m4 = bytearray(same)
    m4[BASE - 1] = 0xFF
    check("負對照：表前一個 byte 不算", diff(same, bytes(m4))[1] == 0)
    m5 = bytearray(same + b"\xff")
    check("負對照：表後一個 byte 不算", diff(same, bytes(m5))[1] == 0)

    print("正對照通過" if not fails else f"{len(fails)} 條沒過")
    return 1 if fails else 0


if __name__ == "__main__":
    sys.exit(main())
