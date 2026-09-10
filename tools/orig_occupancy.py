#!/usr/bin/env python3
"""從 dosgolem 的 `peek` 輸出抽出原版的**單位佔用圖**（docs/spec/183）。

    # ① 先算要 peek 哪幾列（段基底與 Y 從快照的軍團記錄反推）
    tools/py.sh tools/orig_occupancy.py plan 快照.DAT
    # ② 把印出來的那串接在 dosgolem 的步驟後面跑，再解析
    tools/py.sh tools/orig_occupancy.py parse peek.log 輸出.json

⭐ **為什麼快照非帶它不可**：那張圖是**增量維護的歷史**，不是現況的函數
（`sub_12662` 移動時 `dec`／`inc`、`sub_16F86` 編成時 `inc`、`sub_12977`
敗走時 `dec`、`sub_189F0` 開局重建一次）。原版側的對拍是「同一次執行
繼續往下跑」，圖上任何一次漂掉都留著；remake 從快照重建只能得到
「每一支活著的軍團都在自己那一格」這個理想值。兩者在快照當下就可能不同。

版面：每列 384 B（`Y × 24` paragraph ＝ `Y × 384`），X 是列內偏移。
段基底是 `cs:word_19872`，從任一支活著的軍團的 `+0x1C − Y × 24` 反推。
"""

import json
import re
import sys

CORPS_BASE, CORPS_SIZE, CORPS_N = 0x22C0, 64, 127
ROW_BYTES = 384          # 一列 384 B ＝ 24 paragraph
ALIVE = 0x80


def u16(b, off):
    return b[off] | b[off + 1] << 8


def survey(block: bytes):
    """回傳（段基底, {(x, y): 軍團數}）。"""
    seg, cells = None, {}
    for i in range(CORPS_N):
        r = block[CORPS_BASE + i * CORPS_SIZE:CORPS_BASE + (i + 1) * CORPS_SIZE]
        if r[0] < ALIVE:
            continue
        x, y = u16(r, 0x10), u16(r, 0x12)
        if seg is None:
            seg = u16(r, 0x1C) - y * 24
        cells[(x, y)] = cells.get((x, y), 0) + 1
    return seg, cells


def plan(block: bytes) -> str:
    """產生 peek 步驟：有軍團的每一列各讀 384 B。"""
    seg, cells = survey(block)
    if seg is None:
        raise SystemExit("快照裡沒有活著的軍團，反推不出段基底")
    rows = sorted({y for _, y in cells})
    return ";".join(f"peek:{seg + y * 24:04X}:0:{ROW_BYTES}" for y in rows), seg, rows


def parse(text: str, seg: int) -> dict:
    """把 `peek` 的每一行 `段:偏移 = XX XX …` 拼回「座標 → 計數」。"""
    out = {}
    for line in text.splitlines():
        m = re.match(r"\s*([0-9A-Fa-f]{1,4}):([0-9A-Fa-f]{1,4})\s*=\s*((?:[0-9A-Fa-f]{2}\s*)+)$",
                     line.strip())
        if not m:
            continue
        s, off = int(m.group(1), 16), int(m.group(2), 16)
        y, rem = divmod((s - seg) * 16, ROW_BYTES)
        if rem or y < 0:
            continue  # 不是佔用圖的列
        for k, byte in enumerate(m.group(3).split()):
            v = int(byte, 16) & 0x7F  # `sub_13EFD` 抄的時候也去掉 bit 7
            if v:
                out[f"{off + k},{y}"] = v
    return out


def main() -> int:
    args = sys.argv[1:]
    if len(args) >= 2 and args[0] == "plan":
        steps, seg, rows = plan(open(args[1], "rb").read())
        print(f"# 段基底 {seg:#06X}，{len(rows)} 列：{rows}", file=sys.stderr)
        print(steps)
        return 0
    if len(args) >= 3 and args[0] == "parse":
        seg = int(args[3], 16) if len(args) > 3 else None
        if seg is None:
            raise SystemExit("parse 要帶段基底（十六進位）")
        cells = parse(open(args[1], encoding="utf-8", errors="replace").read(), seg)
        with open(args[2], "w", encoding="utf-8") as fh:
            json.dump({"schema": "wolong-occupancy/1", "cells": cells}, fh,
                      ensure_ascii=False, indent=2, sort_keys=True)
        print(f"寫出 {args[2]}：{len(cells)} 格非零")
        return 0
    if args and args[0] == "--selftest":
        return selftest()
    print(__doc__)
    return 2


def selftest() -> int:
    fails = []

    def check(label, cond):
        print(("  ok  " if cond else "  FAIL") + "  " + label)
        if not cond:
            fails.append(label)

    # 造一份只有兩支軍團的區塊：座標 (5,15) 與 (7,15)，段基底 0x46B3。
    b = bytearray(CORPS_BASE + CORPS_N * CORPS_SIZE)
    for i, (x, y) in enumerate([(5, 15), (7, 15)]):
        r = CORPS_BASE + i * CORPS_SIZE
        b[r] = 0xC0
        b[r + 0x10], b[r + 0x12] = x, y
        segv = 0x46B3 + y * 24
        b[r + 0x1C], b[r + 0x1D] = segv & 0xFF, segv >> 8
    seg, cells = survey(bytes(b))
    check("段基底反推得出來", seg == 0x46B3)
    check("兩支軍團在同一列的兩格", cells == {(5, 15): 1, (7, 15): 1})

    steps, _, rows = plan(bytes(b))
    check("只 peek 有軍團的那一列", rows == [15] and steps.count("peek:") == 1)
    check("列的段算對", f"{0x46B3 + 15 * 24:04X}" in steps)

    # 解析：第 15 列的第 5 格是 2、第 7 格是 1。
    row = ["00"] * ROW_BYTES
    row[5], row[7] = "82", "01"  # 82h ＝ 2 ＋ bit 7，抄的時候要去掉
    text = f"{0x46B3 + 15 * 24:04X}:0000 = " + " ".join(row)
    got = parse(text, 0x46B3)
    check("解析出兩格，bit 7 被去掉", got == {"5,15": 2, "7,15": 1})
    # 負對照：不屬於佔用圖的段要被忽略。
    check("負對照：別的段不算", parse("2754:0000 = 01 02 03", 0x46B3) == {})

    print("正對照通過" if not fails else f"{len(fails)} 條沒過")
    return 1 if fails else 0


if __name__ == "__main__":
    sys.exit(main())
