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

# 19 個處理常式的位址。**名稱還沒解讀**，這裡只放位址，
# 免得用臆測的名字污染筆記（CLAUDE.md §8）。
HANDLERS = [
    "sub_1A48C", "sub_1A495", "sub_1A4AB", "sub_1A4BF", "sub_1A50D",
    "sub_1A516", "sub_1A52E", "sub_1A560", "sub_1A564", "sub_1A57A",
    "sub_1A591", "sub_1A5E1", "sub_1A5F7", "sub_1A60D", "sub_1A654",
    "sub_1A65D", "sub_1A69F", "sub_1A6CF", "sub_1A6E8",
]


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
    print(f"若位元組隨機，期望約 {total * (32 - NUM_OPCODES) // 32} 個")
    for i, w in bad[:10]:
        print(f"  位移 0x{i:04X}: {w.hex(' ')}")
    return 0 if not bad else 1


def cmd_list(data, _args):
    print("段  武將+0x16  戰場類別  用到的指令")
    for i, s in scripts(data):
        ops = sorted({decode(s[o:o + 2])[0] for o in range(0, SCRIPT_SIZE, 2)})
        print(f"{i:2d}      {i // 4}        {i % 4}      {ops}")
    hist = collections.Counter(decode(data[o:o + 2])[0]
                               for o in range(0, len(data), 2))
    print("\n全檔指令分佈（指令碼: 次數）：")
    for op, n in sorted(hist.items()):
        print(f"  {op:2d} {HANDLERS[op]:>10}  {n:5d}")


def cmd_dump(data, args):
    n = args.script
    s = data[n * SCRIPT_SIZE:(n + 1) * SCRIPT_SIZE]
    print(f"段 {n}（武將 +0x16 ＝ {n // 4}，戰場類別 {n % 4}）")
    for o in range(0, SCRIPT_SIZE, 2):
        op, par, arg = decode(s[o:o + 2])
        print(f"  {o:3d}  {s[o]:02x} {s[o+1]:02x}   op={op:2d} "
              f"par={par}  arg=0x{arg:02x}  {HANDLERS[op]}")


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
