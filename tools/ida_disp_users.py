#!/usr/bin/env python3
# IDAPython：全庫掃「運算元位移等於指定值」的指令。
#
#   printf '858\n85B\n' > workplace/ida/dosv/census/disp_list.txt
#   tools/ida.sh script dosv tools/ida_disp_users.py KI.EXE.i64
#
# 為什麼要有這一支：`[si+858h]` 這種段內欄位存取**沒有交叉參考**——
# IDA 的 xref 只涵蓋直接位址參考，`基址暫存器 + 常數位移` 完全不在裡面。
# 想知道「誰讀寫據點 +0x18」，grep 反組譯文字會漏掉換了寫法的那些
# （`[bx+858h]`／`[di+858h]`／`es:[si+858h]`），而漏掉的通常正是寫入端。
# 解碼運算元才問得準。
#
# 輸出 /work/disp_users.txt。第一行是 probe（函式數 ＋ 輸入檔雜湊）；
# 沒有 probe 就分不出「沒找到」與「沒跑到」。
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_ua
import idautils
import idc

LIST = "/work/disp_list.txt"
OUT = "/work/disp_users.txt"

# 寫入類助憶碼。判讀寫不要看運算元順序以外的東西——
# `cmp` 與 `test` 讀而不寫，混進來會讓「寫入端」清單灌水。
WRITE = {"mov", "and", "or", "xor", "add", "sub", "adc", "sbb", "inc", "dec"}


def main():
    ida_auto.auto_wait()
    want = []
    for line in open(LIST, encoding="utf-8"):
        line = line.split("#")[0].strip()
        if line:
            want.append(int(line, 16))
    rows = []
    nfunc = 0
    for ea in idautils.Functions():
        nfunc += 1
        fn = idc.get_func_name(ea)
        f = ida_funcs.get_func(ea)
        for head in idautils.Heads(f.start_ea, f.end_ea):
            insn = ida_ua.insn_t()
            if not ida_ua.decode_insn(insn, head):
                continue
            for k, op in enumerate(insn.ops):
                if op.type != ida_ua.o_displ:
                    continue
                if op.addr not in want:
                    continue
                mnem = insn.get_canon_mnem()
                kind = "寫" if (k == 0 and mnem in WRITE) else "讀"
                rows.append((op.addr, fn, head, kind, idc.GetDisasm(head)))
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("probe: %d 支函式, sha256=%s, 位移 %s\n"
                 % (nfunc, ida_nalt.retrieve_input_file_sha256().hex()[:16],
                    "、".join("%X" % d for d in want)))
        for disp in want:
            sel = [r for r in rows if r[0] == disp]
            fh.write("\n==== +%03Xh：%d 處 ====\n" % (disp, len(sel)))
            for _, fn, head, kind, dis in sel:
                fh.write("  %-12s %08X %s  %s\n" % (fn, head, kind, dis.strip()))
    ida_pro.qexit(0)


main()
