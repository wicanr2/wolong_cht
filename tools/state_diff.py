#!/usr/bin/env python3
"""逐欄比對原版與 remake 的狀態表（規格：`docs/spec/138-state-table-parity.md`）。

    tools/py.sh tools/state_diff.py 原版.log remake.log
    tools/py.sh tools/state_diff.py --selftest

⭐ **畫面相同只證明畫得一樣，讀得一樣要另外證。** 三張表裡有一半的欄位
畫面上根本看不到（威脅量、求援冷卻、心向、原主、經費餘額）。

兩邊的輸出是同一個版面：`勢力 NN 君NNN 師NNN …`，欄位是
「中文標籤 ＋ 數字」，名字欄位是 `名<hex>`。這支把每一行拆成
`{標籤: 值}`，再逐鍵比。**不比對行的順序**，比對編號。

⚠ 名字用 `cp950` 解回來顯示，比對仍然比 hex——解碼只影響人看得懂，
不影響判定（`docs/spec/138` §3）。
"""

import re
import sys

# 一行的樣子：`勢力  0 君 16 師 62 …`／`據點 56 主 0 名C0D9B3B1 (243,100) …`
KINDS = ("勢力", "據點", "武將")
FIELD = re.compile(r"([一-鿿]{1,2})(-?\d+|[0-9A-F]*)(?![0-9A-F])")


def parse(text):
    """把一份輸出解析成 {表: {編號: {欄位: 值}}}。"""
    out = {k: {} for k in KINDS}
    for line in text.splitlines():
        # 去掉 Go log 的時間戳（`2026/09/06 16:12:27 `）與 dosgolem 的縮排。
        # ⚠ 樣式要**釘住日期時間的形狀**，不能寫成「砍掉前兩個 token」——
        # 那會把沒有時間戳的那一邊砍掉「勢力  0」，於是一列都解不出來
        # 而比對照樣回「零差異」。
        line = re.sub(r"^\s*\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} ", "", line).strip()
        m = re.match(r"^(勢力|據點|武將)\s+(\d+)\s+(.*)$", line)
        if not m:
            continue
        kind, idx, rest = m.group(1), int(m.group(2)), m.group(3)
        # 表頭（`共 22 個勢力（…）：`）不是資料列。兩邊的表頭都改成
        # 「共 N 個…」就不會撞上這個樣式，這一道是雙保險。
        if rest.startswith(("個", "人")) or "（" in rest:
            continue
        row = {}
        # 座標 `(243,100)` 另外抓，其餘走 `標籤 值`。
        xy = re.search(r"\((\s*\d+),(\s*\d+)\)", rest)
        if xy:
            row["X"], row["Y"] = int(xy.group(1)), int(xy.group(2))
            rest = rest[: xy.start()] + " " + rest[xy.end():]
        # `產 5142/10300`、`適 6/ 4/ 0`、`鄰 43, 74, 65,255` 這種多值欄位
        # 拆成 `產0`、`產1`…
        for label, value in re.findall(r"([一-鿿])((?:\s*-?\d+\s*[/,])*\s*-?\d+|[0-9A-F]+)", rest):
            parts = re.split(r"[/,]", value)
            if len(parts) == 1:
                row[label] = parts[0].strip()
            else:
                for i, p in enumerate(parts):
                    row[f"{label}{i}"] = p.strip()
        out[kind][idx] = row
    return out


def big5(h):
    try:
        return bytes.fromhex(h).decode("cp950")
    except Exception:
        return h


def compare(a, b):
    """回傳 (差異列表, 統計)。"""
    diffs = []
    stats = {}
    for kind in KINDS:
        ra, rb = a[kind], b[kind]
        stats[kind] = (len(ra), len(rb))
        for idx in sorted(set(ra) | set(rb)):
            if idx not in ra:
                diffs.append(f"{kind} {idx}：只有 remake 有")
                continue
            if idx not in rb:
                diffs.append(f"{kind} {idx}：只有原版有")
                continue
            for key in sorted(set(ra[idx]) | set(rb[idx])):
                va, vb = ra[idx].get(key), rb[idx].get(key)
                if va == vb:
                    continue
                if key in ("名", "呼"):
                    va, vb = f"{va}({big5(va)})", f"{vb}({big5(vb)})"
                diffs.append(f"{kind} {idx} 欄「{key}」：原版 {va}、remake {vb}")
    return diffs, stats


SAMPLE = """
   勢力  0 君 16 師 62 都 82 備  200/  400/  800 將 11 城 14 團 1 金    73772 戰14 氣200 敵255
   據點 56 主  0 名C0D9B3B1     (243,100) 產 5142/10300 昇107 災106 兵130/184 類1 官255 原  0 鄰 43, 74, 65,255
   武將 16 名B1E4BEDE     呼B1E4BEDE     勢  0 職0 武 9 統13 政13 適 6/ 4/ 0 本2 說1 向255 舊255
"""


def selftest():
    """正對照：同一份比自己必須零差異，改一個欄位必須報得出來。

    ⚠ **少了後半，一支「永遠回零差異」的壞工具會讓所有對拍看起來全過**
    ——那正是最想避免的失敗模式（`CLAUDE.md` §10）。
    """
    ok = True
    a = parse(SAMPLE)
    n = sum(len(a[k]) for k in KINDS)
    if n != 3:
        print(f"  ✗ 解析：抓到 {n} 列，應為 3"); ok = False
    else:
        print("  ✓ 解析：三張表各抓到一列")
    diffs, _ = compare(a, parse(SAMPLE))
    if diffs:
        print(f"  ✗ 同一份比自己回了 {len(diffs)} 個差異"); ok = False
    else:
        print("  ✓ 同一份比自己：零差異")
    for label, mutated in (
        ("數字欄位", SAMPLE.replace("統13", "統12")),
        ("名字欄位", SAMPLE.replace("名B1E4BEDE", "名B1E4BEDF")),
        ("座標", SAMPLE.replace("(243,100)", "(243,101)")),
        ("多值欄位", SAMPLE.replace("鄰 43, 74, 65,255", "鄰 43, 74, 66,255")),
    ):
        diffs, _ = compare(a, parse(mutated))
        if len(diffs) != 1:
            print(f"  ✗ 改 {label}：報了 {len(diffs)} 個差異 {diffs}"); ok = False
        else:
            print(f"  ✓ 改 {label}：抓得到")
    return ok


def main():
    if "--selftest" in sys.argv:
        print("狀態表對拍自我測試（正對照）")
        return 0 if selftest() else 1
    if len(sys.argv) != 3:
        print(__doc__)
        return 2
    a = parse(open(sys.argv[1], encoding="utf-8", errors="replace").read())
    b = parse(open(sys.argv[2], encoding="utf-8", errors="replace").read())
    diffs, stats = compare(a, b)
    for kind in KINDS:
        na, nb = stats[kind]
        print(f"{kind}：原版 {na} 列、remake {nb} 列")
    if not diffs:
        print("逐欄比對通過：三張表每一欄都相同")
        return 0
    print(f"\n不同的欄位共 {len(diffs)} 個：")
    for d in diffs[:80]:
        print(" ", d)
    if len(diffs) > 80:
        print(f"  …另外 {len(diffs) - 80} 個沒列出")
    return 1


if __name__ == "__main__":
    sys.exit(main())
