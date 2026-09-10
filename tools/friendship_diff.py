#!/usr/bin/env python3
"""兩份 SAVE.DAT 的**交友度矩陣**逐 byte 比（區塊 `+0x680`，22 列 × 24 B）。

    tools/py.sh tools/friendship_diff.py 原版.DAT remake.DAT
    tools/py.sh tools/friendship_diff.py --selftest

⭐ **這張表在快照裡一直有，只是沒人比。** `tools/orig_snapshot.py` 的
`TABLE_LEN = 0x5220` 從段內 0 起算，交友度在段內 `0x600` ⇒ 早就 peek 進來了；
而 `faction_diff.py` 只比勢力記錄那 22 × 64 B（段內 0x000–0x580），
`city_diff` 從 0x840 起——**中間那 704 個 byte 沒有任何工具在看**。

宣戰、停戰、協力的效果主要就落在這裡（`sub_13639`：雙向取小值右移一位
並清和平位元），所以少了它，「兩邊都宣戰了嗎」這個問題答不出來。

位址算法出自 `sub_13119`（docs/formats/08 §1.55）：

    交友度[觀察者][對象] ＝ 段內 0x600 + 觀察者 × 24 + 對象

每列 24 byte，實際只用前 22 欄；第 22–23 欄是填充（值 `0x80`）。
**最高位元 1 ＝ 和平**（docs/re/44 §2.2），所以差 `0x80` 是「開戰／停戰」，
差小數字是好感度漂移。
"""
import sys

BASE, ROW, COLS, ROWS = 0x680, 24, 22, 22


def diff(a, b):
    out = []
    for i in range(ROWS):
        for j in range(COLS):
            off = BASE + i * ROW + j
            if a[off] != b[off]:
                out.append((i, j, a[off], b[off]))
    return out


def show(rows):
    for i, j, x, y in rows:
        war = ""
        if (x ^ y) & 0x80:
            war = "　⚠ **和平位元翻了**（開戰／停戰）"
        print(f"  勢力 {i:2d} → {j:2d}：{x:02X} → {y:02X}"
              f"（{x & 0x7F} → {y & 0x7F}）{war}")
    print(f"\n合計 {len(rows)} 個 byte")


def selftest():
    a = bytearray(0x1000)
    b = bytearray(0x1000)
    ok = True

    def check(name, cond):
        nonlocal ok
        print(("  ok    " if cond else "  FAIL  ") + name)
        ok = ok and cond

    check("相同就回空", not diff(a, b))
    b[BASE + 3 * ROW + 5] = 0x7F
    a[BASE + 3 * ROW + 5] = 0xFF
    d = diff(a, b)
    check("抓得到差異", d == [(3, 5, 0xFF, 0x7F)])
    check("和平位元翻轉標得出來", (d[0][2] ^ d[0][3]) & 0x80 != 0)
    # ⚠ 負對照：**範圍外的 byte 不能被算進來**——填充欄與下一張表相鄰，
    # 少一個邊界檢查就會把不相干的差異算成交友度。
    a2, b2 = bytearray(0x1000), bytearray(0x1000)
    a2[BASE + ROWS * ROW] = 0xFF          # 矩陣的下一個 byte
    check("負對照：矩陣之後的 byte 不算", not diff(a2, b2))
    a3, b3 = bytearray(0x1000), bytearray(0x1000)
    a3[BASE + 5 * ROW + COLS] = 0xFF      # 第 22 欄（填充）
    check("負對照：填充欄不算", not diff(a3, b3))
    print("friendship_diff selftest：" + ("通過" if ok else "失敗"))
    return 0 if ok else 1


def main():
    args = sys.argv[1:]
    if args[:1] == ["--selftest"]:
        return selftest()
    if len(args) != 2:
        print(__doc__)
        return 2
    a = open(args[0], "rb").read()
    b = open(args[1], "rb").read()
    show(diff(a, b))
    return 0


if __name__ == "__main__":
    sys.exit(main())
