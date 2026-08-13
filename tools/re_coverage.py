#!/usr/bin/env python3
"""把 IDA 全函式普查與專案文件交叉比對，產出 RE 覆蓋地圖。

回答的問題是「哪些函式從來沒有人讀過」，而不是「我記得解過哪些」。
輸入是 tools/ida_function_census.idc 產出的 census.tsv；
比對對象是 docs/ 底下的文件與 internal/ 的實作。

    tools/py.sh tools/re_coverage.py workplace/ida/dosv/census/census.tsv

分級（只看證據在哪，不看結論對不對）：

    T1  docs/re/ 有寫              反組譯筆記層級的記錄
    T2  docs/ 其他地方有寫          格式、機制、playtest 提到
    T3  只出現在狀態檔或程式碼註解    CONTEXT.md / CLAUDE.md / internal/
    T4  全專案沒有任何一處提過        未觸及

T4 才是真正的缺口。T3 要當心：狀態檔提過不代表讀懂了。
"""
import collections
import os
import re
import sys

SYM = re.compile(r"\b(?:sub|loc|word|byte|off|unk|dword|asc|a)_([0-9A-Fa-f]{4,8})\b")

# KI.EXE 的模組分佈。位址是 IDA DOS/V linear address，
# 邊界由 census 的函式起點聚類與既有筆記的主題對照得出，屬強推論不是 confirmed。
SEGMENTS = [
    (0x10000, 0x11800, "啟動、C runtime 與低階 I/O"),
    (0x11800, 0x13800, "戰略資料存取：勢力、據點、武將表"),
    (0x13800, 0x15800, "訊息格式化與 TALK 輸出"),
    (0x15800, 0x17800, "月結、經濟、AI 決策"),
    (0x17800, 0x19800, "戰略 UI：選單、視窗、數值輸入"),
    (0x19800, 0x1D800, "戰術戰鬥：主迴圈、腳本、單位更新"),
    (0x1D800, 0x1E800, "戰術呈現：顯示串列、相機、合成"),
    (0x1E800, 0x21000, "圖庫解碼、繪圖底層與程式庫"),
]


def segment_of(ea):
    for lo, hi, name in SEGMENTS:
        if lo <= ea < hi:
            return name
    return "未分類"


def load_census(path):
    funcs = {}
    callers = collections.defaultdict(set)
    icalls = collections.defaultdict(list)
    with open(path, encoding="utf-8", errors="replace") as fh:
        for line in fh:
            f = line.rstrip("\n").split("\t")
            if f[0] == "FUNC" and len(f) >= 12:
                start = int(f[1], 16)
                funcs[start] = {
                    "start": start,
                    "bytes": int(f[3]),
                    "insns": int(f[4]),
                    "name": f[5],
                    "callers": int(f[6]),
                    "icalls": int(f[8]),
                }
            elif f[0] == "CALL" and len(f) >= 4:
                callers[int(f[2], 16)].add(int(f[1], 16))
            elif f[0] == "ICALL" and len(f) >= 5:
                icalls[int(f[1], 16)].append((int(f[2], 16), f[4]))
    return funcs, callers, icalls


def scan_mentions(repo, exclude=()):
    """回傳 symbol -> 提到它的檔案集合。

    `exclude` 是路徑片段清單，命中的檔案不計入。用途只有一個：
    **目錄型文件會把自己算進覆蓋率**。`docs/re/24` 逐支登記了每個未讀函式，
    掃描一跑，287 支全部從「沒被提過」升級成「已寫進 docs/re/」，
    指標歸零而實際上一支都沒多讀。排除它才量得到真正讀過的部分。
    """
    hits = collections.defaultdict(set)
    roots = ["docs", "internal", "cmd", "tools", "translations"]
    files = []
    for root in roots:
        for dirpath, dirnames, filenames in os.walk(os.path.join(repo, root)):
            dirnames[:] = [d for d in dirnames if d not in (".git", "images", "__pycache__")]
            for fn in filenames:
                if fn.endswith((".md", ".go", ".py", ".idc", ".sh", ".json", ".txt")):
                    files.append(os.path.join(dirpath, fn))
    for fn in ("CONTEXT.md", "CLAUDE.md", "AGENTS.md", "MEMORY.md",
               "WORKLIST.md", "RESEARCH-LOG.md", "VERIFICATION-MATRIX.md",
               "REMAKE-PLAN.md", "README.md"):
        p = os.path.join(repo, fn)
        if os.path.exists(p):
            files.append(p)
    for path in files:
        try:
            with open(path, encoding="utf-8", errors="replace") as fh:
                text = fh.read()
        except OSError:
            continue
        rel = os.path.relpath(path, repo)
        if any(x in rel.replace(os.sep, "/") for x in exclude):
            continue
        for m in SYM.finditer(text):
            hits[int(m.group(1), 16)].add(rel)
    return hits


def tier(files):
    if not files:
        return "T4"
    if any(f.startswith("docs/re/") for f in files):
        return "T1"
    if any(f.startswith("docs/") for f in files):
        return "T2"
    return "T3"


def main():
    census = sys.argv[1]
    repo = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    funcs, callers, icalls = load_census(census)
    hits = scan_mentions(repo)

    for ea, f in funcs.items():
        f["files"] = hits.get(ea, set())
        f["tier"] = tier(f["files"])
        f["seg"] = segment_of(ea)

    order = ["T1", "T2", "T3", "T4"]
    by_tier = collections.Counter(f["tier"] for f in funcs.values())
    bytes_by_tier = collections.Counter()
    for f in funcs.values():
        bytes_by_tier[f["tier"]] += f["bytes"]
    total_bytes = sum(bytes_by_tier.values())

    print("## 全域覆蓋\n")
    print("| 分級 | 函式數 | 佔比 | 程式碼 bytes | 佔比 |")
    print("|---|---:|---:|---:|---:|")
    label = {"T1": "T1 `docs/re/` 有反組譯筆記",
             "T2": "T2 其他 `docs/` 提到",
             "T3": "T3 只在狀態檔／程式碼",
             "T4": "T4 未觸及"}
    for t in order:
        n = by_tier[t]
        b = bytes_by_tier[t]
        print(f"| {label[t]} | {n} | {n*100//len(funcs)}% | {b:,} | {b*100//total_bytes}% |")
    print(f"| **合計** | **{len(funcs)}** | | **{total_bytes:,}** | |")

    print("\n## 分模組覆蓋\n")
    print("| 位址區間 | 模組 | 函式 | T1 | T2 | T3 | T4 | T4 bytes |")
    print("|---|---|---:|---:|---:|---:|---:|---:|")
    for lo, hi, name in SEGMENTS:
        sel = [f for f in funcs.values() if lo <= f["start"] < hi]
        if not sel:
            continue
        c = collections.Counter(f["tier"] for f in sel)
        t4b = sum(f["bytes"] for f in sel if f["tier"] == "T4")
        print(f"| `{lo:05X}`–`{hi:05X}` | {name} | {len(sel)} | "
              f"{c['T1']} | {c['T2']} | {c['T3']} | {c['T4']} | {t4b:,} |")

    print("\n## T4 未觸及函式：依程式碼大小排序（前 40）\n")
    print("| 函式 | bytes | 指令 | 直接呼叫者 | 已記錄的呼叫者 | 模組 |")
    print("|---|---:|---:|---:|---|---|")
    t4 = sorted((f for f in funcs.values() if f["tier"] == "T4"),
                key=lambda f: -f["bytes"])
    for f in t4[:40]:
        known = sorted(funcs[c]["name"] for c in callers.get(f["start"], ())
                       if c in funcs and funcs[c]["tier"] in ("T1", "T2"))
        kn = "、".join(f"`{k}`" for k in known[:3]) or "—"
        if len(known) > 3:
            kn += f" 等 {len(known)} 支"
        print(f"| `{f['name']}` | {f['bytes']} | {f['insns']} | {f['callers']} | {kn} | {f['seg']} |")
    print(f"\nT4 共 {len(t4)} 支，"
          f"合計 {sum(f['bytes'] for f in t4):,} bytes。")

    orphans = [f for f in funcs.values() if f["callers"] == 0]
    print(f"\n## 沒有直接呼叫者的函式：{len(orphans)} 支\n")
    print("直接 xref 抓不到間接分派與取址後呼叫，所以「零呼叫者」不等於死碼。"
          "下列各支要嘛由跳表選中，要嘛由取址後間接呼叫，要嘛真的沒被用到——"
          "三者要分開驗證。\n")
    print("| 函式 | bytes | 分級 | 模組 |")
    print("|---|---:|---|---|")
    for f in sorted(orphans, key=lambda f: -f["bytes"])[:20]:
        print(f"| `{f['name']}` | {f['bytes']} | {f['tier']} | {f['seg']} |")

    nic = sum(len(v) for v in icalls.values())
    print(f"\n## 間接呼叫點：{nic} 處，分佈在 {len(icalls)} 支函式\n")
    print("每一處都是靜態 xref 的斷點：IDA 列得出候選表，卻不知道 runtime 選哪支。\n")
    print("| 所在函式 | 位址 | 指令 | 分級 |")
    print("|---|---|---|---|")
    for start in sorted(icalls, key=lambda s: -len(icalls[s])):
        for site, dis in icalls[start]:
            t = funcs[start]["tier"] if start in funcs else "?"
            print(f"| `{funcs[start]['name']}` | `{site:08X}` | `{dis.strip()}` | {t} |")


if __name__ == "__main__":
    main()
