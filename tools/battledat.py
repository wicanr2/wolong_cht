#!/usr/bin/env python3
"""BATTLE.DAT 的戰場腳本反組譯器。

`BATTLE.DAT` 是 **32 段 × 256 byte 的腳本**，每段是給一種武將在一種戰場上
用的行動腳本。原版每幀執行一個 word（`sub_1A426`）：

    低 5 位 = 指令碼（0–18）      高 3 位 = 參數（0–7）
    第二個 byte = 指令的運算元

段編號 = 武將記錄 `+0x16` × 4 + 戰場類別（`docs/re/11` §3.1）。

    tools/battledat.py list                 # 32 段各用到哪些指令
    tools/battledat.py dump 8               # 反組譯第 8 段
    tools/battledat.py verify               # 驗證整個檔案都是合法指令

原版資產不隨本專案散布，預設從 workplace/orig/dosv/ 讀。
"""

import argparse
import collections
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
ORIG = ROOT / "workplace" / "orig" / "dosv"

SCRIPT_SIZE = 256          # 一段（＝ sub_1CBE5 的 di）
NUM_SCRIPTS = 32           # 8192 ÷ 256
NUM_OPCODES = 19           # funcs_1A457 的表長

# 19 個指令。名稱與作用出自 docs/re/11 §3.5。
OPS = [
    ("wait",     "等待 N 幀"),
    ("form",     "切換陣形：word_1D344 ← 陣形編號 × 96"),
    ("line",     "陣形線：0→58 自軍側、1→36 中央、≥2→16 敵軍側"),
    ("order",    "下命令：參數 7＝全軍、0–5＝指定隊"),
    ("q.d346",   "條件 ← byte_1D346"),
    ("q.d33c",   "條件 ← word_1D33C 與 0x1C 比大小（0/1/2）"),
    ("q.mine",   "條件 ← 我方六隊命令的最小值"),
    ("q.theirs", "條件 ← 敵方六隊命令的最小值"),
    ("q.d31e",   "條件 ← byte_1D31E（word_1D31A 高位 ≤ 0x20 時取 2）"),
    ("q.rand",   "條件 ← 亂數 mod N（N＝0 時 mod 6）"),
    ("branch",   "條件分支，下一個 word 是目標"),
    ("q.a24",    "條件 ← word_1D30A:0x24（上限 255）"),
    ("q.a04",    "條件 ← word_1D30A:0x04（上限 255）"),
    ("order.by", "依 [+0x24] == 參數×18 的隊下命令"),
    ("q.d31c",   "條件 ← word_1D31C 高位"),
    ("q.min18",  "條件 ← 0xC00 那張表 [+0x18] 的最小值 × 4"),
    ("msg",      "訊息 0x1CE + N（依 byte_1D349 與參數決定要不要出）"),
    ("q.cmd9",   "條件 ← 我方大將的命令 ≥ 9"),
    ("q.b03",    "條件 ← 我方大將的 [+0x03]"),
]

# 分支指令的五種比較（`sub_1A591` 的 switch）。
CMP = ["無條件", "==", "!=", "<", ">"]

BRANCH = 10


def decode(word):
    """回傳 (指令碼, 參數, 運算元)。"""
    lo, hi = word[0], word[1]
    return lo & 0x1F, (lo & 0xE0) >> 5, hi


def scripts(data):
    for i in range(NUM_SCRIPTS):
        yield i, data[i * SCRIPT_SIZE:(i + 1) * SCRIPT_SIZE]


def cmd_verify(data, _args):
    """整個檔案是不是都是合法指令。

    這是「BATTLE.DAT 就是腳本」最硬的一條證據：4,096 個 word 的低 5 位
    **沒有一個超過 18**，而位元組若是隨機的，期望會有約 1,664 個超過。
    """
    bad = [(i, w) for i in range(0, len(data), 2)
           if (data[i] & 0x1F) >= NUM_OPCODES for w in [data[i:i + 2]]]
    total = len(data) // 2
    print(f"{total} 個 word：低 5 位 ≥ {NUM_OPCODES}（無效指令）的有 {len(bad)} 個")
    # 分支指令後面那個 word 一定是目標（低位元組 0）。
    br = tgt = 0
    for _, s in scripts(data):
        for i in range(0, len(s) - 3, 2):
            if s[i] & 0x1F == BRANCH:
                br += 1
                tgt += s[i + 2] == 0
    print(f"分支指令 {br} 個，後面那個 word 低位為 0 的 {tgt} 個")
    print(f"若位元組隨機，期望約 {total * (32 - NUM_OPCODES) // 32} 個")
    for i, w in bad[:10]:
        print(f"  位移 0x{i:04X}: {w.hex(' ')}")
    return 0 if not bad else 1


def walk(script):
    """依序走一段腳本，回傳 [(位移, 指令碼, 參數, 運算元, 目標或 None)]。

    ⚠ **不能每兩個 byte 當一個指令。** 分支指令（10）後面那個 word 是
    **跳躍目標**，不是指令——它長得像「等待」（低位元組 0），
    照直線切會把 564 個目標誤讀成 564 次等待（docs/re/11 §3.4）。
    """
    out, i = [], 0
    while i < len(script) - 1:
        lo, hi = script[i], script[i + 1]
        op, par = lo & 0x1F, (lo & 0xE0) >> 5
        target = None
        if op == BRANCH and i + 3 < len(script) and script[i + 2] == 0:
            target = script[i + 3]
        out.append((i, op, par, hi, target))
        i += 4 if target is not None else 2
    return out


def cmd_list(data, _args):
    print("段  武將+0x16  戰場類別  指令數  用到的指令")
    for i, s in scripts(data):
        prog = walk(s)
        ops = sorted({p[1] for p in prog})
        print(f"{i:2d}      {i // 4}        {i % 4}      {len(prog):4d}   {ops}")
    hist = collections.Counter(p[1] for _, s in scripts(data) for p in walk(s))
    print("\n全檔指令分佈：")
    for op, n in sorted(hist.items()):
        print(f"  {op:2d} {OPS[op][0]:>9}  {n:5d}   {OPS[op][1]}")


def cmd_dump(data, args):
    n = args.script
    s = data[n * SCRIPT_SIZE:(n + 1) * SCRIPT_SIZE]
    prog = walk(s)
    labels = {p[4] * 2 for p in prog if p[4] is not None}
    print(f"段 {n}（武將 +0x16 ＝ {n // 4}，戰場類別 {n % 4}）")
    for off, op, par, arg, target in prog:
        mark = ">" if off in labels else " "
        name = OPS[op][0]
        if op == BRANCH:
            text = f"{name} {CMP[par] if par < len(CMP) else par} {arg} → {target * 2}"
        else:
            text = f"{name} 參數={par} 運算元={arg}"
        print(f"{mark} {off:3d}  {s[off]:02x} {s[off+1]:02x}  {text}")


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--orig", type=pathlib.Path, default=ORIG)
    sub = ap.add_subparsers(dest="cmd", required=True)
    sub.add_parser("verify", help="驗證整個檔案都是合法指令")
    sub.add_parser("list", help="32 段各用到哪些指令")
    p = sub.add_parser("dump", help="反組譯一段")
    p.add_argument("script", type=int)
    args = ap.parse_args()

    path = args.orig / "BATTLE.DAT"
    if not path.exists():
        sys.exit(f"找不到 {path}（原版資產請自備）")
    data = path.read_bytes()
    if len(data) != SCRIPT_SIZE * NUM_SCRIPTS:
        sys.exit(f"{path} 是 {len(data)} B，預期 {SCRIPT_SIZE * NUM_SCRIPTS}")
    if args.cmd == "dump" and not 0 <= args.script < NUM_SCRIPTS:
        sys.exit(f"段編號要在 0–{NUM_SCRIPTS - 1}")
    return {"verify": cmd_verify, "list": cmd_list, "dump": cmd_dump}[args.cmd](data, args)


if __name__ == "__main__":
    sys.exit(main() or 0)
