#!/usr/bin/env python3
"""逐 byte 比兩份區塊的全域欄位（`+0x00`–`+0x3A`，59 B）。

    tools/py.sh tools/global_diff.py 原版.DAT remake.DAT
    tools/py.sh tools/global_diff.py --selftest

⭐ **只比前 59 個 byte。** 那一段對映到原版執行期的 `cs:0CF0h`–`cs:0D2Ah`
（`internal/state/state.go` 的欄位註解），檢查點就是從那裡 `ipeek` 出來的。

⛔ **`+0x3B` 之後不能比。** 存檔檔案那裡放的是劇本標題（`titleOffset`），
而執行期記憶體的同一段（`cs:0D2Bh` 起）是別的變數——兩者只是剛好落在
同一個位移上。拿檢查點去比會看到「原版一直在改標題」，那是假的。
"""

import sys

SPAN = 0x3B  # 59 B：cs:0CF0h–cs:0D2Ah
NAMES = {
    0x00: "日", 0x01: "該月天數", 0x02: "子刻", 0x03: "時",
    0x04: "月", 0x05: "月續", 0x06: "年", 0x07: "年續",
    0x0D: "玩家記錄位址", 0x0E: "玩家記錄位址續", 0x0F: "玩家",
    0x10: "信賴度",
    0x12: "顯示用本月收入", 0x13: "同續", 0x14: "同續",
    0x15: "顯示用本月支出", 0x16: "同續", 0x17: "同續",
    0x18: "稅率",
    0x1A: "募兵數(騎)", 0x1B: "同續", 0x1C: "募兵數(弓)", 0x1D: "同續",
    0x1E: "募兵數(步)", 0x1F: "同續",
    0x20: "來月稅率",
    0x28: "軍團游標", 0x29: "軍團游標續",
    0x2C: "每時勢力游標", 0x2D: "同續",
    0x2E: "據點游標", 0x2F: "據點游標續",
    0x30: "事件佇列游標", 0x31: "同續",
    0x3A: "存活勢力數",
}


def diff(a, b):
    return [k for k in range(SPAN) if a[k] != b[k]]


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
    for k in d:
        name = NAMES.get(k, "")
        tag = f"({name})" if name else ""
        print(f"+0x{k:02X}{tag} {a[k]:02X}→{b[k]:02X}")
    print(f"\n合計 {len(d)} 個 byte（比的是 +0x00–+0x{SPAN - 1:02X}）")
    return 0


def selftest():
    fails = []

    def check(label, cond):
        print(("  ok  " if cond else "  FAIL") + "  " + label)
        if not cond:
            fails.append(label)

    base = bytes(0x100)
    check("完全相同 → 0 個 byte", diff(base, base) == [])

    m = bytearray(base)
    m[0x2E] = 0x20
    check("改據點游標 → 抓到", diff(base, bytes(m)) == [0x2E])

    m2 = bytearray(base)
    m2[SPAN - 1] = 0xFF
    check("最後一格（+0x3A）也在範圍內", diff(base, bytes(m2)) == [SPAN - 1])

    # ⛔ 負對照：59 B 之後是執行期變數，不能被算進來。
    m3 = bytearray(base)
    m3[SPAN] = 0xFF
    check("負對照：+0x3B 不算", diff(base, bytes(m3)) == [])
    m4 = bytearray(base)
    m4[0x42] = 0xFF
    check("負對照：標題那一段不算", diff(base, bytes(m4)) == [])

    print("正對照通過" if not fails else f"{len(fails)} 條沒過")
    return 1 if fails else 0


if __name__ == "__main__":
    sys.exit(main())
