#!/usr/bin/env python3
"""稽核：同一個 IDA 符號在不同文件裡被標成**互斥的兩種角色**。

    tools/py.sh tools/symbol_roles.py
    tools/py.sh tools/symbol_roles.py --selftest

⭐ **為什麼要這個。** `phantom_scan` 擋「指向不存在的東西」、
`stale_scan` 擋「值不對」，但**兩份文件對同一支函式各說各話**時
兩支都是綠的——而那正是最貴的一種錯：兩邊都 CONFORMED，
照著任一份實作都「有出處」。

> 2026-09-07：`sub_10241` 在 `docs/spec/45` 是「播第 6 曲」、
> 在 `docs/spec/79` 是「調色盤組 0」。remake 照後者把啟動殼層的背景
> 載成第 0 組，而正確答案是第 1 組——逐像素才抓到
> （`docs/playtest/106` §4）。

⚠ **這是稽核工具不是閘**：一支函式本來就可能既畫圖又讀輸入，
所以誤報是預期的。它的用途是**列出候選讓人去比對**，
不接進 `tools/check.sh`。
"""

import os
import re
import sys
from collections import defaultdict

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SYMBOL = re.compile(r"\bsub_([0-9A-F]{4,5})\b")

# 互斥的角色類別。**只放真的互斥的**——「畫圖」與「讀輸入」不互斥，
# 而「播曲」與「載調色盤」是。
ROLES = {
    "音訊": ("BGM", "播曲", "播第", "音效", "音源", "曲）", "曲、", "第 6 曲"),
    "調色盤": ("調色盤", "色盤", "palette", "GAMEPAL"),
    "存檔": ("存檔", "SAVE.DAT", "槽位", "四槽"),
    "亂數": ("亂數", "rand(", "擲骰"),
}


def docs():
    for root, _, names in os.walk(os.path.join(REPO, "docs")):
        for n in sorted(names):
            if n.endswith(".md"):
                yield os.path.join(root, n)


def scan(paths):
    """回傳 {符號: {角色: {(檔, 行號)}}}。"""
    seen = defaultdict(lambda: defaultdict(set))
    for path in paths:
        rel = os.path.relpath(path, REPO)
        with open(path, encoding="utf-8") as fh:
            for i, line in enumerate(fh, 1):
                syms = set(SYMBOL.findall(line))
                if not syms:
                    continue
                for role, words in ROLES.items():
                    if any(w in line for w in words):
                        for sym in syms:
                            seen[sym][role].add((rel, i))
    return seen


def conflicts(seen):
    return {s: r for s, r in seen.items() if len(r) > 1}


def selftest():
    import tempfile
    ok = True

    def want(label, cond):
        nonlocal ok
        print(f"  {'✓' if cond else '✗'} {label}")
        ok = ok and cond

    with tempfile.TemporaryDirectory() as tmp:
        a = os.path.join(tmp, "a.md")
        b = os.path.join(tmp, "b.md")
        open(a, "w", encoding="utf-8").write("`sub_10241` 播第 6 曲\n")
        open(b, "w", encoding="utf-8").write("`sub_10241` 載調色盤組 0\n")
        c = conflicts(scan([a, b]))
        want("擋下「同一支既是播曲又是調色盤」", "10241" in c)
        open(b, "w", encoding="utf-8").write("`sub_10241` 也是播曲\n")
        want("同一種角色不報", not conflicts(scan([a, b])))
        open(b, "w", encoding="utf-8").write("`sub_9999` 載調色盤組 0\n")
        want("不同符號不報（負對照）", not conflicts(scan([a, b])))
    print("正對照" + ("通過" if ok else "失敗"))
    return 0 if ok else 1


def main():
    if "--selftest" in sys.argv:
        print("符號角色衝突掃描自我測試（正對照）")
        return selftest()
    found = conflicts(scan(list(docs())))
    if not found:
        print("符號角色衝突：0 筆")
        return 0
    print(f"符號角色衝突：{len(found)} 筆（**稽核用，誤報是預期的**）\n")
    for sym in sorted(found):
        roles = found[sym]
        print(f"## sub_{sym}（{len(roles)} 種角色）")
        for role in sorted(roles):
            where = "、".join(f"{d}:{n}" for d, n in sorted(roles[role])[:3])
            print(f"  {role}：{where}")
        print()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
