#!/usr/bin/env python3
"""做一份**受控存檔**給對拍用：從既有的 SAVE.DAT 出發，只改指定欄位。

    tools/py.sh tools/parity_save.py <來源 SAVE.DAT> <輸出目錄> [選項]
    tools/py.sh tools/parity_save.py --selftest

⭐ **為什麼要這個。** 即時制的對拍最怕「兩邊各自擲骰」：天災、雲的漂移、
AI 的判斷都吃亂數，取樣點一拉長兩邊就分開，而分開之後看到的殘差
**不是規則分歧，是亂數分歧**。改狀態比等亂數可靠——把會動的東西
先關掉，剩下的差異才有意義（`CLAUDE.md` §3.1、`docs/spec/90`）。

目前支援的改法：

  --no-clouds   把大地圖上那 16 朵會飄的雲停用（物件記錄 +0x00 的 bit 7）。
                原版 `sub_12459`／`sub_12533` 都以 `cmp [si], 80h / jb` 跳過
                停用的槽，所以清掉就等於整個關掉——**不移動也不繪製**
                （`docs/spec/146`）。

⚠ **輸出目錄要自己備妥其餘原版檔案**（dosgolem 的 gamedir 需要整套）。
這一支只寫 `SAVE.DAT`，不碰來源，也不碰 `workplace/orig/`。
"""

import argparse
import os
import shutil
import sys

BLOCK = 22208
BLOCKS = 4
FILE_SIZE = BLOCK * BLOCKS

# 32 筆物件記錄在區塊裡的位移（段內 0x2040 ＋ 0x80），每筆 16 B。
# 後 16 筆是常駐的雲（docs/spec/146 §2）。
MAP_OBJECT_BASE = 0x20C0
MAP_OBJECT_SIZE = 16
CLOUD_FIRST, CLOUD_COUNT = 16, 16


def disable_clouds(data: bytearray, block: int) -> int:
    """清掉雲的存在旗標，回傳實際關掉幾朵。"""
    base = block * BLOCK + MAP_OBJECT_BASE
    n = 0
    for slot in range(CLOUD_FIRST, CLOUD_FIRST + CLOUD_COUNT):
        off = base + slot * MAP_OBJECT_SIZE
        if data[off] & 0x80:
            data[off] &= ~0x80 & 0xFF
            n += 1
    return n


def selftest() -> int:
    """正對照 ＋ 負對照。少了負對照，一支「什麼都沒改」的工具也會通過。"""
    ok = True

    def check(label, cond):
        nonlocal ok
        print(f"  {'✓' if cond else '✗'} {label}")
        ok = ok and cond

    data = bytearray(FILE_SIZE)
    base = MAP_OBJECT_BASE
    for slot in range(32):
        data[base + slot * MAP_OBJECT_SIZE] = 0x80
    n = disable_clouds(data, 0)
    check("關掉 16 朵雲", n == 16)
    check("前 16 槽（火災／暴動）沒被動到",
          all(data[base + s * MAP_OBJECT_SIZE] == 0x80 for s in range(16)))
    check("後 16 槽的 bit 7 清掉了",
          all(data[base + s * MAP_OBJECT_SIZE] == 0x00 for s in range(16, 32)))
    check("其餘 byte 一個都沒動",
          sum(1 for i, b in enumerate(data)
              if b != 0 and i % MAP_OBJECT_SIZE == 0) == 16)
    # 負對照：本來就沒有雲的區塊要回 0，不能假裝有做事。
    empty = bytearray(FILE_SIZE)
    check("空區塊回 0 朵（負對照）", disable_clouds(empty, 0) == 0)
    # 負對照：只改指定的區塊。
    two = bytearray(FILE_SIZE)
    for slot in range(32):
        two[BLOCK + base + slot * MAP_OBJECT_SIZE] = 0x80
    check("只動指定的區塊（負對照）", disable_clouds(two, 0) == 0)
    return 0 if ok else 1


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("source", nargs="?")
    ap.add_argument("outdir", nargs="?")
    ap.add_argument("--slot", type=int, default=0, help="改第幾個區塊（0–3）")
    ap.add_argument("--no-clouds", action="store_true")
    ap.add_argument("--selftest", action="store_true")
    ns = ap.parse_args()

    if ns.selftest:
        return selftest()
    if not ns.source or not ns.outdir:
        ap.error("要給來源 SAVE.DAT 與輸出目錄")
    if not 0 <= ns.slot < BLOCKS:
        ap.error("--slot 要在 0–3")

    data = bytearray(open(ns.source, "rb").read())
    if len(data) != FILE_SIZE:
        ap.error(f"{ns.source} 是 {len(data)} B，預期 {FILE_SIZE}")

    changes = []
    if ns.no_clouds:
        changes.append(f"停用 {disable_clouds(data, ns.slot)} 朵雲")
    if not changes:
        ap.error("沒有指定任何改法")

    os.makedirs(ns.outdir, exist_ok=True)
    out = os.path.join(ns.outdir, "SAVE.DAT")
    if os.path.abspath(out) == os.path.abspath(ns.source):
        ap.error("輸出會蓋掉來源；換一個目錄")
    with open(out, "wb") as f:
        f.write(data)
    print(f"{out}（第 {ns.slot + 1} 槽）：" + "、".join(changes))
    return 0


if __name__ == "__main__":
    sys.exit(main())
