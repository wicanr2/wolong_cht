#!/usr/bin/env python3
"""把 dosgolem 的 `ipeek:1ECFC:258` 輸出轉成 `-rng-state` 吃的 258 byte 檔。

    tools/dosgolem.sh out "…;ipeek:1ECFC:258" > peek.log
    tools/py.sh tools/rng_state.py peek.log out.bin
    tools/py.sh tools/rng_state.py --selftest

⚠ **吃檔名不吃 stdin**：`tools/py.sh` 的 `docker run` 沒有 `-i`，
管線進來的東西到不了容器裡。

原版的產生器狀態只有三樣東西，而且都在固定位址（`docs/re/10` §1）：
計數器 `cs:1ECFCh`、狀態 `cs:1ECFDh`、256 byte 置換表 `cs:1ECFEh`。
讀出來灌進 remake，兩邊的亂數流從那一刻起完全相同（`docs/spec/147` §5）。
"""

import re
import sys

WANT = 258


def parse(text: str) -> bytes:
    """從 dosgolem 的輸出裡撈出 `1ECFC = XX XX …` 那一行。"""
    for line in text.splitlines():
        m = re.match(r"\s*1ECFC\s*=\s*((?:[0-9A-Fa-f]{2}\s*)+)$", line.strip())
        if m:
            return bytes(int(x, 16) for x in m.group(1).split())
    raise SystemExit("找不到 `1ECFC = …` 那一行；`ipeek:1ECFC:258` 跑了嗎？")


def selftest() -> int:
    ok = True

    def check(label, cond):
        nonlocal ok
        print(f"  {'✓' if cond else '✗'} {label}")
        ok = ok and cond

    body = " ".join("%02X" % (i & 0xFF) for i in range(WANT))
    got = parse(f"noise\n   1ECFC = {body}\nmore noise\n")
    check("撈得到那一行", len(got) == WANT)
    check("順序沒被打亂", got[0] == 0 and got[1] == 1 and got[257] == 1)
    # 負對照：位址不對的不能撈。
    try:
        parse("   1ECFD = 00 01\n")
        check("擋下位址不對的（負對照）", False)
    except SystemExit:
        check("擋下位址不對的（負對照）", True)
    # 負對照：長度不足要看得出來。
    short = parse("   1ECFC = 00 01 02\n")
    check("長度不足回得出來（負對照）", len(short) == 3)
    return 0 if ok else 1


def main() -> int:
    if "--selftest" in sys.argv[1:]:
        return selftest()
    if len(sys.argv) != 3:
        raise SystemExit(__doc__)
    with open(sys.argv[1], encoding="utf-8", errors="replace") as f:
        data = parse(f.read())
    if len(data) != WANT:
        raise SystemExit(f"讀到 {len(data)} byte，預期 {WANT}——`ipeek` 的長度給對了嗎？")
    open(sys.argv[2], "wb").write(data)
    print(f"{sys.argv[2]}：{len(data)} byte（c={data[0]:02X} s={data[1]:02X}）")
    return 0


if __name__ == "__main__":
    sys.exit(main())
