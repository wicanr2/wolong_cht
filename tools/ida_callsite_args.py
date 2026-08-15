#!/usr/bin/env python3
# IDAPython（跑在 ida-pro-9.4-idapython:py312-v1；tools/ida.sh 依 .py 副檔名自動選）
#
# 把「呼叫某支常式的每一個呼叫點」連同**呼叫前餵進暫存器的立即值**列出來。
#
# 為什麼要這一支：`sub_10337(al = 場景編號)` 這種「用一個數字選內容」的
# 共用常式，光看被呼叫端解不出誰用哪一個場景；而場景的身分若改用
# **字串內容**去認，會得到看似合理的錯配——場景 4 與 5 都有「將軍／
# 總兵力／六個編成位置」，靠字串分不出哪一個是編成畫面。
# **呼叫端的立即值是一手證據，字串長相是二手推論。**
#
# 目標寫在 /work/call_targets.txt（一行一個符號名），工具本身不必改。
#
#   printf 'sub_10337\nsub_1895D\n' > workplace/ida/dosv/census/call_targets.txt
#   tools/ida.sh script dosv tools/ida_callsite_args.py KI.EXE.i64
#
# ⚠ headless 的 print 不進 stdout，一律寫檔；輸出檔第一行是 probe
#   （函式數 ＋ 輸入檔 SHA-256），否則分不出「沒找到」與「沒跑到」。
#
# 輸出 /work/callsite_args.txt
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_xref
import idc

LIST = "/work/call_targets.txt"
OUT = "/work/callsite_args.txt"
BACK = 24

REGS = ("al", "ah", "ax", "bl", "bh", "bx", "cl", "ch", "cx",
        "dl", "dh", "dx", "si", "di")


def fname(ea):
    f = ida_funcs.get_func(ea)
    return ida_funcs.get_func_name(f.start_ea) if f else "?"


def preceding_immediates(call_ea, func_start):
    """往回走，收集最後一次寫進各暫存器的值。

    遇到別的 call 就停——跨過一個 call 之後暫存器不再由這一段負責，
    硬讀會得到看似合理的假值。
    """
    seen = {}
    ea = idc.prev_head(call_ea, func_start)
    steps = 0
    while ea != idc.BADADDR and ea >= func_start and steps < BACK:
        steps += 1
        mnem = idc.print_insn_mnem(ea).lower()
        if mnem == "call":
            break
        if mnem == "mov":
            dst = idc.print_operand(ea, 0).lower()
            if dst in REGS and dst not in seen:
                if idc.get_operand_type(ea, 1) == idc.o_imm:
                    seen[dst] = idc.get_operand_value(ea, 1)
                else:
                    seen[dst] = idc.print_operand(ea, 1)
        elif mnem == "xor":
            dst = idc.print_operand(ea, 0).lower()
            if dst in REGS and dst == idc.print_operand(ea, 1).lower() \
                    and dst not in seen:
                seen[dst] = 0
        ea = idc.prev_head(ea, func_start)
    return seen


def main():
    ida_auto.auto_wait()
    try:
        with open(LIST, encoding="utf-8") as fh:
            targets = [ln.split("#")[0].strip() for ln in fh]
    except OSError as exc:
        targets = []
    targets = [t for t in targets if t]

    lines = ["呼叫點立即值（IDA DOS/V linear address）",
             "輸入檔 SHA-256：%s" % ida_nalt.retrieve_input_file_sha256().hex(),
             "函式數：%d" % ida_funcs.get_func_qty(),
             "目標數：%d" % len(targets)]
    for tgt in targets:
        ea = idc.get_name_ea_simple(tgt)
        lines.append("")
        lines.append("==== %s EA=%08X ====" % (tgt, ea))
        if ea == idc.BADADDR:
            lines.append("  找不到符號")
            continue
        n = 0
        x = ida_xref.get_first_cref_to(ea)
        while x != idc.BADADDR:
            if idc.print_insn_mnem(x).lower() == "call":
                n += 1
                f = ida_funcs.get_func(x)
                imm = preceding_immediates(x, f.start_ea if f else x)
                lines.append("  %08X  %-14s  %s" % (
                    x, fname(x),
                    " ".join("%s=%s" % (k, ("%Xh" % v) if isinstance(v, int) else v)
                             for k, v in sorted(imm.items()))))
            x = ida_xref.get_next_cref_to(ea, x)
        lines.append("  呼叫點共 %d 個" % n)
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(lines) + "\n")
    ida_pro.qexit(0)


main()
