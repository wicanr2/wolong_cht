#!/usr/bin/env python3
"""把 dosgolem 的 `-watch` 記錄切成子刻，逐拍與 remake 的取數序列比對。

    tools/py.sh tools/parity_pace_diff.py 原版.log remake_seq.txt [--from-step N]
    tools/py.sh tools/parity_pace_diff.py --selftest

⭐ **切拍要照原版一拍內的順序**（`sub_13EFD`，docs/re/44 §1）：

    sub_13F74（威脅／求援）→ sub_14194（內政）→ sub_14269（災害）

`sub_14194` 是唯一每拍必被呼叫一次的錨點，但**威脅／求援在它之前**，
所以 `14060`（`sub_14057`）與 `14171`（`sub_14155`）取的數要歸到
**下一個**錨點，不是上一個。照字面順序歸位會讓整條求援時間軸平移一拍，
而平移之後的形狀仍然「看起來像一種規則差異」。

⚠ **同一次呼叫可能被攔兩次。** 實測 `#30740174` 與 `#30740182` 是同一個
`sub_14194`（SI 相同、相差 8 道指令）。不去重就會多出一拍取數 0 的
幽靈子刻——而那長得就像「原版這一拍什麼都沒做」。
"""

import re
import sys

# 這兩個來源在該拍的 `sub_14194` 之**前**執行，要歸到下一個錨點。
PRE = ("14060", "14171")
# 兩次錨點相差不到這麼多道指令就當成同一次呼叫。一拍是兩萬多道，差距懸殊。
DEDUPE_STEPS = 1000

LINE = re.compile(r"#(\d+) 呼叫 (\w+) .*SI=([0-9A-F]{4}).*近=([0-9A-F]+)")


def parse(text, anchor="14194"):
    """回傳 [{'city':編號,'step':指令數,'rng':[來源…]}]，一項一個子刻。"""
    ticks, pend = [], []
    for line in text.splitlines():
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
            if near in PRE:
                pend.append(near)
            elif ticks:
                ticks[-1]["rng"].append(near)
    return ticks


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
    t = parse(log)
    check("兩拍（重複的錨點被去掉）", len(t) == 2)
    check("第一拍是據點 144、取 2 個", t[0]["city"] == 144 and len(t[0]["rng"]) == 2)
    check("14060 歸到下一拍", t[1]["rng"][0] == "14060" and len(t[1]["rng"]) == 2)
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
    return 1 if fails else 0


def main():
    args = sys.argv[1:]
    if args[:1] == ["--selftest"]:
        return selftest()
    if len(args) < 2:
        print(__doc__)
        return 2
    from_step = 0
    if "--from-step" in args:
        i = args.index("--from-step")
        from_step = int(args[i + 1])
        del args[i:i + 2]
    ticks = parse(open(args[0], encoding="utf-8", errors="replace").read())
    base = 0
    if from_step:
        base = next(i for i, t in enumerate(ticks) if t["step"] > from_step)
    remake, rwhere = [], []
    for l in open(args[1]):
        if not l.strip():
            continue
        f = l.split()
        remake.append(int(f[1]))
        rwhere.append(f[2].split(",") if len(f) > 2 else [])
    n = min(len(ticks) - base, len(remake))
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
