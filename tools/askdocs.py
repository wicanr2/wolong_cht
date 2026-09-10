#!/usr/bin/env python3
"""用關鍵詞查 `docs/` 的標題與狀態，回答「這件事解過沒」。

    tools/py.sh tools/askdocs.py 召見 對話框
    tools/py.sh tools/askdocs.py --selftest

⛔ **這支要在「準備說卡住／要開新工作」的那一刻跑**，不是收尾時。
`grep -r` 也做得到同樣的事——差別在它只印**標題 ＋ 狀態行 ＋ 路徑**，
一眼看得出「那份說通過了」，不會被全文淹沒。

2026-09-10：同局面對拍撞到「孫乾　大人，主公有事召見。」停住，
於是我把「玩家事件怎麼處理」寫成三選一的決策呈給使用者——
而 `docs/playtest/118` 的**標題**就是
「劉備 90 天無人值守：君主召見用右鍵推進，選項用移動次數」，
狀態「通過」，連 90 天跑得完都驗過了。

`CLAUDE.md` §7 第 1 條與 `rules/00-rules-index.md` 都寫過「先查手上已有的」，
而它還是在「宣告卡住」那一刻沒被想起——**規則有了，缺的是觸發點**。
"""
import os
import re
import sys

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "docs")


def scan(root):
    """回傳 (路徑, 標題, 狀態行) 的清單。"""
    out = []
    for dirpath, _, names in os.walk(root):
        for n in sorted(names):
            if not n.endswith(".md"):
                continue
            p = os.path.join(dirpath, n)
            title, status = "", ""
            try:
                with open(p, encoding="utf-8") as fh:
                    for line in fh:
                        if not title and line.startswith("# "):
                            title = line[2:].strip()
                            continue
                        if title and not status and "狀態" in line:
                            status = re.sub(r"\s+", " ", line.strip())[:80]
                            break
                        if title and line.startswith("## "):
                            break
            except OSError:
                continue
            out.append((os.path.relpath(p, os.path.join(root, "..")), title, status))
    return out


def search(docs, words):
    hits = []
    for path, title, status in docs:
        hay = f"{path} {title} {status}"
        if all(w in hay for w in words):
            hits.append((path, title, status))
    return hits


def selftest():
    docs = scan(ROOT)
    ok = True

    def check(name, cond):
        nonlocal ok
        print(("  ok    " if cond else "  FAIL  ") + name)
        ok = ok and cond

    check("掃得到文件", len(docs) > 100)
    check("標題抓得到", sum(1 for _, t, _ in docs if t) > len(docs) * 0.9)
    check("狀態行抓得到", sum(1 for _, _, s in docs if s) > len(docs) * 0.5)
    # ⭐ 正對照：用一組真實的關鍵詞，確認它找得到那一份。
    hits = search(docs, ["召見"])
    check("關鍵詞『召見』找得到 playtest/118",
          any("118" in p for p, _, _ in hits))
    # 負對照：不存在的詞要回空，不能什麼都命中。
    check("負對照：亂打的詞回空", not search(docs, ["zzz不存在的詞zzz"]))
    print("askdocs selftest：" + ("通過" if ok else "失敗"))
    return 0 if ok else 1


def main():
    args = sys.argv[1:]
    if args[:1] == ["--selftest"]:
        return selftest()
    if not args:
        print(__doc__)
        return 2
    hits = search(scan(ROOT), args)
    if not hits:
        print("沒有標題或狀態行命中——**這不代表沒解過**，"
              "再查 docs/INDEX.md 的斷言總表與 docs/re/43 的缺口總表。")
        return 1
    print(f"{len(hits)} 份命中：")
    for path, title, status in hits:
        print(f"  {path}")
        print(f"    {title}")
        if status:
            print(f"    {status}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
