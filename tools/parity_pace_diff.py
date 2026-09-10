#!/usr/bin/env python3
"""把 dosgolem 的 `-watch` 記錄切成子刻，逐拍與 remake 的取數序列比對。

    tools/py.sh tools/parity_pace_diff.py 原版.log remake_seq.txt [--skip-orig N]
    tools/py.sh tools/parity_pace_diff.py 原版.log remake_seq.txt --tick 2967
    tools/py.sh tools/parity_pace_diff.py --selftest

⛔ **這支只比「這一拍取了幾次」，不比取數的來源。** 順序錯位而次數相同時
它一片綠——2026-09-10 的月結就是這樣：`sub_12BD9` 的政略在原版排在災害
之前，remake 排在之後，兩邊都是 647 次，`--tick 2967` 並排才看得出來
（docs/spec/185）。**下「這一拍一致」的結論之前，先用 `--tick` 看一次序列。**

⚠ **兩邊的第 1 拍不一定是同一拍。** 原版側的 `-watch` 常常比快照早開始
記錄，實測差 **27 個子刻**。沒有對齊就比，會得到一個「看起來很像規則差異」
的 18.8%，而第一個分歧落在第 11 拍——**形狀與真的規則分歧一模一樣**。

⭐ **對齊不靠人工填位移，靠據點巡迴游標。** 兩邊都說得出「這一拍處理的是
哪一個據點」——原版是 `sub_14194` 的 `SI ÷ 32`，remake 是 `Tick` 之前的
`CityCursor`（`rng_pace.go -seq-out` 的第三欄）。本工具用它自動對齊，
並且**逐拍檢查兩邊的游標一直同步**；一旦脫節就當場停下來說在第幾拍。

沒有游標欄的舊序列檔退回**位移掃描**（0–48）並印警告；
`--skip-orig N` 可以手動指定原版要先丟掉幾拍。

⭐ **切拍要照原版一拍內的順序**（`sub_13EFD`，docs/re/44 §1）：

    sub_13F74（威脅／求援）→ sub_14194（內政）→ sub_14269（災害）

`sub_14194` 是唯一每拍必被呼叫一次的錨點，但**威脅／求援在它之前**，
所以 `14060`（`sub_14057`）與 `14171`（`sub_14155`）取的數要歸到
**下一個**錨點，不是上一個。照字面順序歸位會讓整條求援時間軸平移一拍，
而平移之後的形狀仍然「看起來像一種規則差異」。

⚠ **同一次呼叫可能被攔兩次。** 實測 `#30740174` 與 `#30740182` 是同一個
`sub_14194`（SI 相同、相差 8 道指令）。不去重就會多出一拍取數 0 的
幽靈子刻——而那長得就像「原版這一拍什麼都沒做」。

⚠ **`sub_1ECE0` 也會被攔兩次**，而且更難認：多出來的那一筆長得像
「原版在這一拍多取一個亂數」，也就是**一種規則差異的形狀**。判準是
指令距離——`sub_1ECE0` 本體 12 條指令，兩次真呼叫之間不可能更短
（實測 17,214 次裡有 2 次相差 8 條、暫存器完全相同）。
去掉幾筆會印在報表上，不會靜靜吃掉。
"""

import re
import sys

# ⭐ 這幾個來源在該拍的 `sub_14194` 之**前**執行，要歸到下一個錨點。
# 全部都在威脅／求援那一段（`sub_13F74` → `sub_14028`／`sub_14057`
# → `sub_140C9`／`sub_14155`），而錨點是內政的 `sub_14194`。
#
# ⛔ **漏一個位址的症狀跟真的規則分歧一模一樣**：`140F9`
# （`sub_140C9` 的玩家求援冷卻擲骰）漏掉時，報表說「原版拍 5654 求援、
# remake 拍 5655 求援」，而實測那一筆的指令序號比同拍的 `14194`
# **早 43 道**、`SI=0820`（據點 65）——remake 是對的。
# 新增位址前先看一次 log 的序號，不要照函式名猜。
PRE = ("14060", "14171", "140F9")
# 兩次錨點相差不到這麼多道指令就當成同一次呼叫。一拍是兩萬多道，差距懸殊。
DEDUPE_STEPS = 1000
# `sub_1ECE0` 本體 12 條指令（`0001ECE0`–`0001ECFB`），所以兩次**真**呼叫
# 之間至少差這麼多。比這更近而且來源相同的，是同一次被攔了兩下。
RNG_MIN_STEPS = 12

LINE = re.compile(r"#(\d+) 呼叫 (\w+) .*SI=([0-9A-F]{4}).*近=([0-9A-F]+)")
# dosgolem 的 `clock` 步驟會印這一行。它是**對齊點的一手證據**：
# 數一數它前面有幾個錨點，就是原版比 remake 早跑的拍數。
CLOCK = re.compile(r"遊戲時鐘：(\d+)年(\d+)月(\d+)日 (\d+)時")


def parse(text, anchor="14194"):
    """回傳（子刻串, 時鐘標記串）。

    子刻是 `{'city':編號,'step':指令數,'rng':[來源…]}`；時鐘標記是
    `(前面有幾個錨點, (年,月,日,時))`。
    """
    ticks, pend, clocks = [], [], []
    last_rng, dropped = None, 0
    for line in text.splitlines():
        c = CLOCK.search(line)
        if c:
            clocks.append((len(ticks), tuple(int(g) for g in c.groups())))
            continue
        m = LINE.search(line)
        if not m:
            continue
        step, label, si, near = int(m.group(1)), m.group(2), int(m.group(3), 16), m.group(4)
        if label == anchor:
            if ticks and step - ticks[-1]["step"] < DEDUPE_STEPS and ticks[-1]["city"] == si // 32:
                continue  # 同一次呼叫被攔兩次
            ticks.append({"city": si // 32, "step": step, "rng": list(pend)})
            pend = []
        elif label == "1ECE0":
            if last_rng and last_rng[0] == near and step - last_rng[1] < RNG_MIN_STEPS:
                dropped += 1  # 同一次呼叫被攔兩下
                continue
            last_rng = (near, step)
            if near in PRE:
                pend.append(near)
            elif ticks:
                ticks[-1]["rng"].append(near)
    if dropped:
        print(f"⚠ 去掉 {dropped} 筆重複攔截的 `sub_1ECE0`"
              f"（來源相同而且相差不到 {RNG_MIN_STEPS} 道指令）")
    return ticks, clocks


def selftest():
    fails = []

    def check(label, cond):
        print(("  ok  " if cond else "  FAIL") + "  " + label)
        if not cond:
            fails.append(label)

    log = "\n".join([
        "  #100 呼叫 14194 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=13F5D 遠=0",
        "  #110 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=141D8 遠=0",
        "  #120 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=14216 遠=0",
        "  #900 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1220 DI=0 來自 近=14060 遠=0",
        "  #999 呼叫 14194 AX=0 BX=0 CX=0 DX=0 SI=1220 DI=0 來自 近=13F5D 遠=0",
        "  #1002 呼叫 14194 AX=0 BX=0 CX=0 DX=0 SI=1220 DI=0 來自 近=13F5D 遠=0",
        "  #1010 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1220 DI=0 來自 近=141D8 遠=0",
    ])
    t, _ = parse(log)
    check("兩拍（重複的錨點被去掉）", len(t) == 2)
    check("第一拍是據點 144、取 2 個", t[0]["city"] == 144 and len(t[0]["rng"]) == 2)
    check("14060 歸到下一拍", t[1]["rng"][0] == "14060" and len(t[1]["rng"]) == 2)

    # ⭐ `140F9`（求援冷卻）同樣在錨點之前——這一條是照實測 log 的形狀寫的：
    # `#220195316 … 近=140F9 SI=0820` 排在 `#220195359 呼叫 14194 SI=0820` 前面。
    log2 = "\n".join([
        "  #100 呼叫 14194 AX=0 BX=0 CX=0 DX=0 SI=0800 DI=0 來自 近=13F5D 遠=0",
        "  #110 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=0800 DI=0 來自 近=141D8 遠=0",
        "  #900 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=0820 DI=0 來自 近=140F9 遠=0",
        "  #950 呼叫 14194 AX=0 BX=0 CX=0 DX=0 SI=0820 DI=0 來自 近=13F5D 遠=0",
        "  #960 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=0820 DI=0 來自 近=141D8 遠=0",
    ])
    t2, _ = parse(log2)
    check("140F9 歸到下一拍（據點 65 那一拍）",
          len(t2) == 2 and t2[0]["city"] == 64 and len(t2[0]["rng"]) == 1
          and t2[1]["city"] == 65 and t2[1]["rng"] == ["140F9", "141D8"])
    # 負對照：不去重就會多一拍，14060 照字面歸位就會落在第一拍。
    naive, pend = [], []
    for line in log.splitlines():
        m = LINE.search(line)
        if m and m.group(2) == "14194":
            naive.append(0)
        elif m and m.group(2) == "1ECE0" and naive:
            naive[-1] += 1
    check("負對照：照字面切會多一拍", len(naive) == 3)
    check("負對照：照字面切會把 14060 算進前一拍", naive[0] == 3)

    # 對齊自檢：正對照要抓到 27 拍的位移，負對照（本來就對齊）要閉嘴。
    import random
    rnd = random.Random(7)
    body = [rnd.choice((2, 3, 3, 3, 4)) for _ in range(600)]
    shifted = [rnd.choice((2, 3, 4)) for _ in range(27)] + body
    hit = align_check(shifted, body, 0)
    check("對齊自檢：抓到 27 拍的位移", hit is not None and hit[0] == 27)
    check("對齊自檢：本來就對齊時不吭聲", align_check(body, body, 0) is None)
    check("對齊自檢：位移已經給對時不吭聲", align_check(shifted, body, 27) is None)

    # 據點游標對齊：正對照要找出 27，負對照（游標中途脫節）要抓得到。
    ocity = [(144 + i) % 192 for i in range(600)]
    rcity = ocity[27:]
    # ⚠ 游標每 192 拍繞一圈，所以 27 與 219 都對得上——**游標對齊本身
    # 是多解的**，它只能當驗證，不能當唯一依據（時鐘標記才是一手證據）。
    k, why = align_by_city(ocity, rcity)
    check("游標對齊：多解時不敢認", k is None and "不唯一" in why)
    k2, _ = align_by_city(ocity[:200], rcity[:173])
    check("游標對齊：只有一解時算得出 27", k2 == 27)
    check("游標對齊：對齊之後沒有脫節", city_drift(ocity, rcity, 27) is None)
    bad = list(rcity)
    bad[100] = (bad[100] + 1) % 192
    check("負對照：游標中途脫節抓得到", city_drift(ocity, bad, 27) == (101, rcity[100], bad[100]))
    check("負對照：據點對不上時說找不到位移",
          align_by_city(ocity, [999] + rcity[1:])[0] is None)

    # 序列檔的兩種版面都讀得到（第三欄是後來才加的）。
    import tempfile, os as _os
    fd, path = tempfile.mkstemp(suffix=".txt")
    with _os.fdopen(fd, "w") as fh:
        fh.write("# 起點 196年4月16日 16時 子刻 3\n"
                 "1 3 171 a.go:1,a.go:1,a.go:1\n2 2 172 a.go:1,a.go:1\n")
    c, wsrc, city, start = read_remake(path)
    check("新版面：讀得到據點欄", city == [171, 172] and c == [3, 2])
    check("新版面：讀得到起點時鐘", start == (196, 4, 16, 16))
    check("新版面：來源沒被吃掉", wsrc[0] == ["a.go:1"] * 3)
    with open(path, "w") as fh:
        fh.write("1 3 a.go:1,a.go:1,a.go:1\n2 2 a.go:1,a.go:1\n")
    c, wsrc, city, start = read_remake(path)
    check("舊版面：沒有據點欄也讀得到", city == [] and c == [3, 2] and start is None)
    check("舊版面：來源沒被當成據點", wsrc[1] == ["a.go:1"] * 2)
    _os.unlink(path)

    # 時鐘標記對齊：正對照要算出 27，兩個負對照要說不敢認。
    clocks = [(27, (196, 4, 16, 16)), (5507, (196, 5, 13, 4))]
    k, why = align_by_clock(clocks, (196, 4, 16, 16, 3))
    check("時鐘對齊：算出原版先丟 27 拍", k == 27)
    check("負對照：時鐘對不上就不敢認",
          align_by_clock(clocks, (196, 4, 16, 15, 0))[0] is None)
    check("負對照：同一個時刻出現兩次就不敢認",
          align_by_clock(clocks + [(300, (196, 4, 16, 16))], (196, 4, 16, 16, 3))[0] is None)
    check("負對照：沒有 clock 標記時說沒有", align_by_clock([], (196, 4, 16, 16, 3))[0] is None)

    # `sub_1ECE0` 的重複攔截：來源相同而且相差不到 12 道指令的才去掉。
    dup = "\n".join([
        "  #100 呼叫 14194 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=13F5D 遠=0",
        "  #110 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=14216 遠=0",
        "  #118 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=14216 遠=0",
    ])
    t, _ = parse(dup)
    check("重複攔截的亂數被去掉", len(t) == 1 and t[0]["rng"] == ["14216"])
    # 反對照 ①：距離夠遠就是兩次真的取數。
    far = dup.replace("#118", "#130")
    t, _ = parse(far)
    check("負對照：距離夠遠的兩次都算數", t[0]["rng"] == ["14216", "14216"])
    # 反對照 ②：來源不同就不是同一次呼叫，再近也都算數。
    other = dup.replace("#118 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=14216",
                        "#118 呼叫 1ECE0 AX=0 BX=0 CX=0 DX=0 SI=1200 DI=0 來自 近=141F1")
    t, _ = parse(other)
    check("負對照：來源不同就都算數", t[0]["rng"] == ["14216", "141F1"])
    return 1 if fails else 0


def read_remake(path):
    """讀 remake 的逐子刻序列。

    版面是 `拍 個數 [據點] 來源1,來源2,…`，前面可以有一行 `# 起點 …` 檔頭。
    **據點欄與檔頭都是後來才加的**，舊檔沒有——據點欄的分辨方式是
    「第三欄是不是純數字」，因為來源那一欄一定含 `.go:`。
    回傳（個數, 來源, 據點, 起點時鐘）。
    """
    counts, where, city, start = [], [], [], None
    for l in open(path):
        if not l.strip():
            continue
        if l.startswith("#"):
            m = re.search(r"(\d+)年(\d+)月(\d+)日 (\d+)時", l)
            if m:
                start = tuple(int(g) for g in m.groups())
            continue
        f = l.split()
        counts.append(int(f[1]))
        rest = f[2:]
        if rest and rest[0].isdigit():
            city.append(int(rest[0]))
            rest = rest[1:]
        where.append(rest[0].split(",") if rest else [])
    if city and len(city) != len(counts):
        raise SystemExit("序列檔有的行有據點欄、有的沒有")
    return counts, where, city, start


def align_by_clock(clocks, start):
    """用 dosgolem 的 `clock` 標記對齊：原版要先丟掉幾拍。

    ⭐ **這是一手證據，不是統計猜測。** 標記印的就是原版在那一刻的
    遊戲時鐘，而 remake 的序列檔檔頭印的是它的起點時鐘；兩個相同的
    那一個標記前面有幾個錨點，就是原版早跑的拍數。

    ⚠ 時鐘只到「時」，同一個小時有 9 個子刻——所以**同一個時刻可能
    印過不只一次**（每個 `clock` 步驟一次）。出現多筆就不敢認。
    """
    want = tuple(start[:4])
    hits = [n for n, c in clocks if c == want]
    if not hits:
        if not clocks:
            return None, "原版 log 裡沒有 `clock` 標記，對不了時鐘"
        return None, (f"原版 log 的 `clock` 標記裡沒有 {want[0]}年{want[1]}月"
                      f"{want[2]}日 {want[3]}時（有 {len(clocks)} 個標記）")
    if len(set(hits)) > 1:
        return None, f"同一個時刻的 `clock` 標記出現在第 {hits} 拍，位移不唯一"
    n = hits[0]
    return n, (f"原版先丟 {n} 拍——log 的 `clock` 標記 "
               f"{want[0]}年{want[1]}月{want[2]}日 {want[3]}時 對上 remake 的起點")


def align_by_city(ocity, rcity, confirm=200):
    """用據點巡迴游標對齊：找原版要先丟掉幾拍。

    ⭐ 游標每 192 拍繞一圈，所以「第一拍的據點相同」有很多解——
    要再看**後面連續幾拍也一致**才算數（`confirm` 拍）。
    """
    if not ocity or not rcity:
        return None, "兩邊都要有據點欄才對得起來"
    n = min(confirm, len(rcity))
    hits = []
    for k in range(len(ocity) - 1):
        if ocity[k] != rcity[0]:
            continue
        m = min(n, len(ocity) - k)
        if all(ocity[k + i] == rcity[i] for i in range(m)):
            hits.append(k)
        if len(hits) > 1:
            break
    if not hits:
        return None, f"找不到讓據點游標對得上的位移（remake 第 1 拍是據點 {rcity[0]}）"
    if len(hits) > 1:
        return None, f"位移不唯一（{hits[:2]} …），把 confirm 拉大或用 --skip-orig"
    return hits[0], f"原版先丟 {hits[0]} 拍（據點游標 {ocity[hits[0]]} 對 {rcity[0]}，前 {n} 拍全同）"


def city_drift(ocity, rcity, base):
    """對齊之後逐拍檢查游標有沒有脫節。回第一個不一致的（拍, 原版, remake）。"""
    n = min(len(ocity) - base, len(rcity))
    for i in range(n):
        if ocity[base + i] != rcity[i]:
            return i + 1, ocity[base + i], rcity[i]
    return None


def align_check(orig, remake, base, span=48, factor=2.0):
    """掃 0–span 的位移；找到比 `base` 好一倍以上的就回報。

    ⭐ **這是防「安靜地比錯對齊」的閘。** 沒有它，起點差 27 拍的兩條序列
    會得到 18.8% 的不一致與一個落在第 11 拍的「第一個分歧」，
    而那個形狀與真的規則分歧分不出來。
    """
    def rate(k):
        n = min(len(orig) - k, len(remake))
        if n <= 0:
            return None, 0
        return sum(1 for i in range(n) if orig[k + i] != remake[i]), n

    cur, curn = rate(base)
    if cur is None or curn == 0:
        return None
    best, bestk, bestn = cur, base, curn
    for k in range(span + 1):
        bad, n = rate(k)
        if bad is None or n == 0:
            continue
        if bad * bestn < best * n:      # 比率比較，避免可比區間不同時失真
            best, bestk, bestn = bad, k, n
    if bestk == base or best * curn * factor >= cur * bestn:
        return None
    return bestk, best, bestn, cur, curn


def main():
    args = sys.argv[1:]
    if args[:1] == ["--selftest"]:
        return selftest()
    if len(args) < 2:
        print(__doc__)
        return 2
    from_step, base = 0, None
    if "--from-step" in args:
        i = args.index("--from-step")
        from_step = int(args[i + 1])
        del args[i:i + 2]
    if "--skip-orig" in args:
        i = args.index("--skip-orig")
        base = int(args[i + 1])
        del args[i:i + 2]
    only = None
    if "--tick" in args:
        i = args.index("--tick")
        only = int(args[i + 1])
        del args[i:i + 2]
    ticks, clocks = parse(open(args[0], encoding="utf-8", errors="replace").read())
    if from_step:
        base = next(i for i, t in enumerate(ticks) if t["step"] > from_step)
    remake, rwhere, rcity, rstart = read_remake(args[1])

    if base is None and rstart is not None:
        base, why = align_by_clock(clocks, rstart)
        if base is None:
            print("⚠ " + why)
        else:
            print(f"對齊：{why}")
    if base is None and rcity:
        base, why = align_by_city([t["city"] for t in ticks], rcity)
        if base is None:
            print("⚠ " + why + "——退回位移掃描")
        else:
            print(f"對齊（游標，⚠ 每 192 拍一循環）：{why}")
    if base is None:
        base = 0
        hint = align_check([len(t["rng"]) for t in ticks], remake, base)
        if hint is not None:
            print(f"⚠ 對齊自檢：原版先丟 {hint[0]} 拍會降到 {hint[1] * 100.0 / hint[2]:.1f}%"
                  f"（現在的 {base} 拍是 {hint[3] * 100.0 / hint[4]:.1f}%）——"
                  f"先確認起點對齊再讀下面的數字")
    elif rcity:
        drift = city_drift([t["city"] for t in ticks], rcity, base)
        if drift is not None:
            print(f"⚠ 據點游標在第 {drift[0]} 拍脫節：原版 {drift[1]}、remake {drift[2]}"
                  f"——之後的逐拍比對沒有意義")
    n = min(len(ticks) - base, len(remake))
    if only is not None:
        k = only - 1
        if not 0 <= k < n:
            print(f"拍 {only} 不在可比區間（1–{n}）")
            return 2
        o = ticks[base + k]
        a, b = o["rng"], rwhere[k] if k < len(rwhere) else []
        print(f"拍 {only} 據點 {o['city']}：原版 {len(a)} 次、remake {len(b)} 次")
        for i in range(max(len(a), len(b))):
            x = a[i] if i < len(a) else "—"
            y = b[i] if i < len(b) else "—"
            print(f"  {i:3d}  {x:<8} {y}")
        return 0
    bad = []
    for k in range(n):
        o = ticks[base + k]
        if len(o["rng"]) != remake[k]:
            bad.append((k + 1, o["city"], len(o["rng"]), remake[k], o["rng"],
                        rwhere[k] if k < len(rwhere) else []))
    print(f"原版 {len(ticks)} 拍（從第 {base + 1} 拍起比），比 {n} 拍，"
          f"不一致 {len(bad)}（{len(bad) * 100.0 / n:.1f}%）")
    for b in bad[:40]:
        print(f"  拍 {b[0]:5d} 據點 {b[1]:3d} 原版 {b[2]:3d} remake {b[3]:3d}")
        print(f"        原版 {b[4][:8]}")
        print(f"        remake {b[5][:8]}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
