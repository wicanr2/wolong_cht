#!/usr/bin/env python3
"""逐筆比兩份區塊的事件佇列（`+0x52C0` 起 256 筆 × 4 B）。

    tools/py.sh tools/event_diff.py 原版.DAT remake.DAT
    tools/py.sh tools/event_diff.py --selftest

⭐ **這是第五張表。** 據點、勢力、軍團、全域四張都比過之後，佇列仍然
一格都沒比——而它決定「哪一天會發生什麼」：宣戰、遷都、合作、停戰、
撥款請求全部排在這裡，由 `sub_131AE` 每 10 個「時」派發一筆。

一筆 4 B：`Code` 低 byte ＝ 事件種類、高 byte ＝ 發起勢力；
`Param` 依事件種類而異（事件 1 的低 byte 是宣戰對象）。
`sub_131AE` 只走前 64 筆（`docs/re/15` §2.1）。

⛔ **同局面對拍現在還不能用它下判斷。** `tools/orig_snapshot.py` 只拼
「全域 59 B ＋ 四張表」兩段，事件佇列那一段是從原版 `SAVE.DAT` **模板**
來的——拿它比 remake 的活狀態，看到的是「模板 vs 執行期」而不是
「原版 vs remake」。佇列在另一個段（`cs:word_10D56`），要先 peek 出來
（worklist `event-queue-not-peeked`）。
"""

import sys

BASE, SIZE, N, DISPATCH = 0x52C0, 4, 0x100, 0x40
KIND = {
    1: "宣戰", 2: "合作", 3: "停戰", 4: "據點撥款", 5: "外交官撥款",
    6: "外交回報", 7: "合作回報", 8: "遷都", 9: "釋放俘虜", 10: "訊息",
    11: "天災", 12: "天災清除", 13: "信賴度",
}


def show(r):
    if r == b"\0\0\0\0":
        return "（空）"
    kind, fac = r[0], r[1]
    name = KIND.get(kind, "?")
    return f"事件{kind:02X}({name}) 勢力{fac:>3} param={r[2] | r[3] << 8:#06x}"


def diff(a, b):
    return [i for i in range(N)
            if a[BASE + i * SIZE:BASE + (i + 1) * SIZE] != b[BASE + i * SIZE:BASE + (i + 1) * SIZE]]


def main():
    args = sys.argv[1:]
    if "--selftest" in args:
        return selftest()
    if len(args) < 2:
        print(__doc__)
        return 2
    a = open(args[0], "rb").read()
    b = open(args[1], "rb").read()
    d = diff(a, b)
    within = [i for i in d if i < DISPATCH]
    for i in d:
        ra = a[BASE + i * SIZE:BASE + (i + 1) * SIZE]
        rb = b[BASE + i * SIZE:BASE + (i + 1) * SIZE]
        tag = "" if i < DISPATCH else "（派發範圍外）"
        print(f"[{i:3d}]{tag} {show(ra)}  →  {show(rb)}")
    print(f"\n合計 {len(d)} 筆，其中 {len(within)} 筆在派發範圍（前 {DISPATCH} 筆）內")
    return 0


def selftest():
    fails = []

    def check(label, cond):
        print(("  ok  " if cond else "  FAIL") + "  " + label)
        if not cond:
            fails.append(label)

    base = bytes(BASE + N * SIZE)
    check("完全相同 → 0 筆", diff(base, base) == [])

    m = bytearray(base)
    m[BASE + 3 * SIZE] = 0x01
    check("改一筆 → 抓到那一筆", diff(base, bytes(m)) == [3])

    m2 = bytearray(base)
    m2[BASE + (N - 1) * SIZE + SIZE - 1] = 0xFF
    check("最後一筆也在範圍內", diff(base, bytes(m2)) == [N - 1])

    # 負對照：佇列前後的 byte 不能被算進來。
    m3 = bytearray(base)
    m3[BASE - 1] = 0xFF
    check("負對照：佇列前一個 byte 不算", diff(base, bytes(m3)) == [])
    m4 = bytearray(base + b"\xff")
    check("負對照：佇列後一個 byte 不算", diff(base, bytes(m4)) == [])

    check("事件種類翻得出名字", "宣戰" in show(b"\x01\x07\xff\x08"))
    check("空的印成（空）", show(b"\0\0\0\0") == "（空）")

    print("正對照通過" if not fails else f"{len(fails)} 條沒過")
    return 1 if fails else 0


if __name__ == "__main__":
    sys.exit(main())
