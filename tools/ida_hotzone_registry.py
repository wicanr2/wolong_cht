#!/usr/bin/env python3
# IDAPython（跑在 ida-pro-9.4-idapython:py312-v1；tools/ida.sh 依 .py 副檔名自動選）
#
# 把**每一個註冊出來的滑鼠熱區**列成表：誰註冊、編號、座標、尺寸。
#
# 為什麼要這一支：主畫面上的可點區域（說明書說的「四個開關」）在畫面上
# 是四個小圖示，靠肉眼估座標點了六次都沒反應——而熱區是**登記出來的**，
# 登記端的立即值就是答案，不必用滑鼠去掃。
#
# `sub_1E3D7(al = 熱區編號, bx = X, dx = Y, cx = 尺寸)`：
# 座標與尺寸的單位在 `docs/re/46` §2 已解（8 px 格），這裡只負責把
# **呼叫端餵進去的立即值**取出來，不做單位換算——換算是文件那一層的事。
#
# 方法：對每一個 call sub_1E3D7，往回走最多 24 條指令，記下最後一次
# 寫進 al/ax/bx/cx/dx 的**立即值**。遇到別的 call 就停——跨過一個
# call 之後暫存器的內容不再由這一段負責，硬讀會得到看似合理的假值。
#
# ⚠ headless 的 print 不進 stdout，一律寫檔；輸出檔要帶 probe
#   （函式數 ＋ 輸入檔 SHA-256），否則分不出「沒找到」與「沒跑到」。
#
# 用法：tools/ida.sh script dosv tools/ida_hotzone_registry.py KI.EXE.i64
# 輸出 /work/hotzone_registry.txt
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_ua
import ida_xref
import idautils
import idc

OUT = "/work/hotzone_registry.txt"
TARGETS = ["sub_1E38C"]
BACK = 24

REGS = {"al": "al", "ah": "ah", "ax": "ax", "bx": "bx", "bl": "bl",
        "cx": "cx", "cl": "cl", "dx": "dx", "dl": "dl", "si": "si", "di": "di"}


def fname(ea):
    f = ida_funcs.get_func(ea)
    return ida_funcs.get_func_name(f.start_ea) if f else "?"


def preceding_immediates(call_ea, func_start):
    """往回走，收集最後一次寫進各暫存器的立即值。"""
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
                # 立即值記數值，其他來源記運算元字面——找「誰把某個段變數
                # 餵進載入器」時，來源是變數名而不是立即值。
                if idc.get_operand_type(ea, 1) == idc.o_imm:
                    seen[dst] = idc.get_operand_value(ea, 1)
                else:
                    seen[dst] = idc.print_operand(ea, 1)
        elif mnem == "xor":
            dst = idc.print_operand(ea, 0).lower()
            src = idc.print_operand(ea, 1).lower()
            if dst in REGS and dst == src and dst not in seen:
                seen[dst] = 0
        ea = idc.prev_head(ea, func_start)
    return seen


def main():
    ida_auto.auto_wait()
    lines = []
    lines.append("熱區登記表（IDA DOS/V linear address）")
    lines.append("輸入檔 SHA-256：%s" % ida_nalt.retrieve_input_file_sha256().hex())
    lines.append("函式數：%d" % ida_funcs.get_func_qty())
    for tgt in TARGETS:
        ea = idc.get_name_ea_simple(tgt)
        lines.append("")
        lines.append("==== %s EA=%08X ====" % (tgt, ea))
        if ea == idc.BADADDR:
            lines.append("  找不到符號")
            continue
        n = 0
        x = ida_xref.get_first_cref_to(ea)
        while x != idc.BADADDR:
            # fl_CN=16／fl_CF=17 才是呼叫邊，fl_F=21 是循序落下。
            if ida_xref.xrefblk_t and idc.print_insn_mnem(x).lower() == "call":
                n += 1
                f = ida_funcs.get_func(x)
                start = f.start_ea if f else x
                imm = preceding_immediates(x, start)
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
