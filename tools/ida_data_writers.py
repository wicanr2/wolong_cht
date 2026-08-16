#!/usr/bin/env python3
# IDAPython：列出誰讀、誰寫指定的資料位址，並把寫入端的立即值一起帶出來。
#
# 為什麼不用 `tools/ida_xref.idc`：那一支只印 xref 清單，不分讀寫也不帶值。
# 想回答「這個旗標的初值是多少」時，要的正是**寫入端寫進去的那個立即值**。
#
# ⚠ 讀寫判定用 `XrefType()`（dr_O=1 取位址／dr_W=2 寫／dr_R=3 讀），
#   不比對助憶碼字串（`CLAUDE.md` §4.1）。
# ⚠ xref 只涵蓋直接參考：`ptr = &x` 之後的間接寫入抓不到，
#   症狀是「讀很多處、寫一處」——看到寫入數異常少，先看「取位址」那幾筆。
#
# 目標寫在 TARGETS，一行一個符號名。
# 用法：tools/ida.sh script dosv tools/ida_data_writers.py KI.EXE.i64
# 輸出 /work/data_writers.txt
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_xref
import idc

OUT = "/work/data_writers.txt"
# 要查哪些符號寫在 /work/census/data_writers_list.txt（一行一個）；
# 檔案不在就用這裡的預設。**工具固定，要查什麼是輸入**（同 ida_dump.py）。
LIST = "/work/data_writers_list.txt"
TARGETS = ["word_10D50"]
BACK = 8


def fname(ea):
    f = ida_funcs.get_func(ea)
    return ida_funcs.get_func_name(f.start_ea) if f else "（無函式）"


def targets():
    try:
        with open(LIST, encoding="utf-8") as fh:
            names = [ln.split()[0] for ln in fh if ln.strip() and
                     not ln.lstrip().startswith("#")]
    except OSError:
        return TARGETS
    return names or TARGETS


def main():
    ida_auto.auto_wait()
    lines = ["資料讀寫端（IDA DOS/V linear address）",
             "輸入檔 SHA-256：%s" % ida_nalt.retrieve_input_file_sha256().hex(),
             "函式數：%d" % ida_funcs.get_func_qty()]
    for name in targets():
        ea = idc.get_name_ea_simple(name)
        lines.append("")
        lines.append("==== %s EA=%08X ====" % (name, ea))
        if ea == idc.BADADDR:
            lines.append("  找不到符號")
            continue
        n = 0
        x = ida_xref.get_first_dref_to(ea)
        while x != idc.BADADDR:
            t = idc.get_operand_type(x, 1)
            kind = {1: "取位址", 2: "寫", 3: "讀"}.get(ida_xref.xrefblk_t and 0, "")
            imm = ""
            if t == idc.o_imm:
                imm = "  立即值=%Xh" % idc.get_operand_value(x, 1)
            lines.append("  %08X  %-14s  %s%s" % (x, fname(x), idc.GetDisasm(x), imm))
            n += 1
            x = ida_xref.get_next_dref_to(ea, x)
        lines.append("  參考共 %d 處（%s）" % (n, kind or "讀寫請看反組譯"))
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(lines) + "\n")
    ida_pro.qexit(0)


main()
