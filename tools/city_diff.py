#!/usr/bin/env python3
"""兩份 SAVE.DAT 的據點表逐 byte 比——`corps_diff.py` 的據點版。

    tools/py.sh tools/city_diff.py 原版.DAT remake.DAT [據點編號]
    tools/py.sh tools/city_diff.py --selftest

⭐ **差異的形狀會說話**：吃亂數的欄位（上昇值、防災值、城兵）岔開時是
**對稱雙向的小幅偏移**（偏高的座數 ≈ 偏低的座數）；規則錯的話是單向。
所以最後那張「依欄位」的表比總數有用（docs/playtest/119 §34.3）。
"""

import sys

CB, CS, NC = 0x08C0, 32, 192

# 欄位名照 docs/formats/08。⚠ 上昇值在 `+0x10`、防災值在 `+0x11`、
# **城兵上限在 `+0x12`**、城兵在 `+0x13`——四個連號，抄錯一個看不出來。
NAMES = {
    0x00: "旗標", 0x01: "所屬", 0x08: "X", 0x0A: "Y",
    0x10: "上昇值", 0x11: "防災值", 0x12: "城兵上限", 0x13: "城兵",
    0x18: "軍團數", 0x19: "內政官", 0x1A: "舊主",
}


def diff(a: bytes, b: bytes, only=None):
    """回 (逐座差異, 依欄位的差值分佈)。"""
    rows, hist = [], {}
    for i in range(NC):
        if only is not None and i != only:
            continue
        ra = a[CB + i * CS:CB + (i + 1) * CS]
        rb = b[CB + i * CS:CB + (i + 1) * CS]
        d = [k for k in range(CS) if ra[k] != rb[k]]
        if not d:
            continue
        rows.append((i, [(k, ra[k], rb[k]) for k in d]))
        for k in d:
            hist.setdefault(k, []).append(rb[k] - ra[k])
    return rows, hist


def report(rows, hist):
    total = sum(len(cols) for _, cols in rows)
    for i, cols in rows:
        s = "  ".join(f'+{k:#04x}({NAMES.get(k, "?")}) {x:02X}→{y:02X}'
                      for k, x, y in cols)
        print(f"據點 {i:>3}：{s}")
    print(f"\n合計 {total} 個 byte／{len(rows)} 座")
    if hist:
        print("\n依欄位：")
        for k in sorted(hist):
            d = hist[k]
            plus = sum(1 for v in d if v > 0)
            minus = sum(1 for v in d if v < 0)
            print(f'  +{k:#04x}({NAMES.get(k, "?")})：{len(d):>3} 座'
                  f'　remake 偏高 {plus}／偏低 {minus}'
                  f'　差值 {min(d):+d} … {max(d):+d}')
    return total


def selftest():
    """⚠ 正反對照：只驗「它印得出差異」證明不了它在比對的地方。"""
    ok = True

    def check(name, cond):
        nonlocal ok
        print(f'  {"✓" if cond else "✗"} {name}')
        ok = ok and cond

    base = bytearray(CB + NC * CS)
    same, _ = diff(bytes(base), bytes(base))
    check("完全相同 → 0 座", not same)

    # 據點 5 的城兵差 3。
    mod = bytearray(base)
    mod[CB + 5 * CS + 0x13] = 3
    rows, hist = diff(bytes(base), bytes(mod))
    check("改一個 byte → 剛好一座、一欄",
          len(rows) == 1 and rows[0][0] == 5 and len(rows[0][1]) == 1)
    check("欄位認得出是城兵（+0x13）", 0x13 in hist and hist[0x13] == [3])

    # ⚠ 據點表之外的 byte 不該被算進來——否則軍團表一動就誤報。
    outside = bytearray(base)
    outside[CB - 1] = 0xFF
    outside += b"\xff" * 64
    check("據點表以外的差異不算", not diff(bytes(base), bytes(outside))[0])

    # 只看一座時，別座的差異要被濾掉。
    two = bytearray(base)
    two[CB + 5 * CS + 0x13] = 1
    two[CB + 9 * CS + 0x13] = 1
    check("指定據點只回那一座",
          len(diff(bytes(base), bytes(two), only=9)[0]) == 1)

    # 差值的方向要分得出來（偏高／偏低是判斷「亂數岔開 vs 規則錯」的依據）。
    down = bytearray(base)
    down[CB + 5 * CS + 0x10] = 0
    up = bytearray(base)
    up[CB + 5 * CS + 0x10] = 7
    base2 = bytearray(base)
    base2[CB + 5 * CS + 0x10] = 3
    _, h1 = diff(bytes(base2), bytes(down))
    _, h2 = diff(bytes(base2), bytes(up))
    check("差值帶正負號（分得出偏高／偏低）",
          h1[0x10] == [-3] and h2[0x10] == [4])

    print("city_diff 自我測試：" + ("通過" if ok else "**失敗**"))
    return 0 if ok else 1


def main(argv):
    if "--selftest" in argv:
        return selftest()
    args = [a for a in argv[1:] if not a.startswith("--")]
    if len(args) < 2:
        print(__doc__)
        return 2
    a = open(args[0], "rb").read()
    b = open(args[1], "rb").read()
    only = int(args[2]) if len(args) > 2 else None
    report(*diff(a, b, only))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
