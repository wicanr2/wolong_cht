#!/usr/bin/env python3
# IDAPython：通用函式 dumper。要 dump 哪幾支寫在 /work/dump_list.txt，
# 一行一筆「位址 說明」（位址是 IDA linear address，十六進位、可帶 0x）。
#
#   printf '1699E 位元2清除端\n144A9 or [si],2\n' > workplace/ida/dosv/census/dump_list.txt
#   tools/ida.sh script dosv tools/ida_dump.py KI.EXE.i64
#
# 這樣不必為了換幾個位址就改工具檔——工具是固定的，要看什麼是輸入。
# 每支附呼叫者、它碰到的表基址（同 tools/ida_tables.py 的 BASES）與資料 xref。
#
# 輸出 /work/dump.txt。probe 在第一行；沒有 probe 就分不出「沒找到」與「沒跑到」。
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_ua
import ida_xref
import idautils
import idc

LIST = "/work/dump_list.txt"
OUT = "/work/dump.txt"

BASES = {
    0x0840: "據點表", 0x0842: "據點表+2", 0x2240: "軍團表", 0x2242: "軍團表+2",
    0x4200: "城兵臨時軍團", 0x4240: "武將表", 0x4242: "武將表+2",
    0x4257: "武將+0x17職務", 0x425E: "武將+0x1E說話型",
}


def callers(ea):
    out, x = [], ida_xref.get_first_cref_to(ea)
    while x != idc.BADADDR:
        f = ida_funcs.get_func(x)
        out.append(ida_funcs.get_func_name(f.start_ea) if f else "%08X" % x)
        x = ida_xref.get_next_cref_to(ea, x)
    return out


def tables_in(f):
    hit = set()
    ea = f.start_ea
    while ea < f.end_ea:
        insn = ida_ua.insn_t()
        if ida_ua.decode_insn(insn, ea) > 0:
            for op in insn.ops:
                if op.type == ida_ua.o_void:
                    break
                v = None
                if op.type == ida_ua.o_imm:
                    v = op.value & 0xFFFF
                elif op.type == ida_ua.o_displ:
                    v = op.addr & 0xFFFF
                if v in BASES:
                    hit.add(BASES[v])
        ea = idc.next_head(ea, f.end_ea)
    return sorted(hit)


def main():
    ida_auto.auto_wait()
    try:
        rows = [l.split(None, 1) for l in open(LIST, encoding="utf-8")
                if l.strip() and not l.startswith("#")]
    except OSError:
        rows = []
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("probe: %d 支函式, sha256=%s, 清單 %d 筆\n"
                 % (len(list(idautils.Functions())),
                    ida_nalt.retrieve_input_file_sha256().hex()[:16], len(rows)))
        for r in rows:
            ea = int(r[0], 16)
            if ea < 0x10000:
                ea += 0x10000          # 允許只寫段內偏移
            label = r[1].strip() if len(r) > 1 else ""
            f = ida_funcs.get_func(ea)
            fh.write("\n==== %08X %s ====\n" % (ea, label))
            if not f:
                fh.write("  ⚠ 這個位址沒有函式（可能只被間接參考）\n")
                continue
            fh.write("呼叫者：%s\n" % " ".join(callers(f.start_ea)))
            fh.write("碰到的表：%s\n" % ("、".join(tables_in(f)) or "（無立即值線索）"))
            fh.write("bytes=%d\n" % (f.end_ea - f.start_ea))
            p = f.start_ea
            while p != idc.BADADDR and p < f.end_ea:
                fh.write("  %08X  %s\n" % (p, idc.GetDisasm(p)))
                p = idc.next_head(p, f.end_ea)
    ida_pro.qexit(0)


main()
