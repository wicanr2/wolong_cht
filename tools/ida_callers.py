#!/usr/bin/env python3
# IDAPython：查「誰呼叫這個位址」，含 IDA 沒認成函式的區段。
#
#   printf '1D98B 設 bit5\n1D9AF 清 bit5\n' \
#     > workplace/ida/dosv/census/callers_list.txt
#   tools/ida.sh script dosv tools/ida_callers.py KI.EXE.i64
#
# 為什麼不用 tools/ida_dump.py 的 callers()：那一支要先 get_func()，
# 而**目標本身沒被 IDA 認成函式時整支就 dump 不到**——
# `0x1D98B`／`0x1D9AF` 這兩支（顯示格旗標 bit 5 的設定端與清除端）
# 正是這種情形，逐函式掃描一個都看不到，於是「沒有人設 bit 5」
# 這個**假的負證據**差一點成立。
#
# 三種來源都收：
#   - code xref（IDA 認得的 call／jmp 邊）
#   - 全 segment 逐條解碼，目標位址等於指定值的 call／jmp（含未認成函式的區段）
#   - 立即值等於該位址 offset 的指令（取址後間接呼叫）
#
# 輸出 /work/callers.txt。probe 在第一行。
import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import ida_segment
import ida_ua
import idautils
import idc

LIST = "/work/callers_list.txt"
OUT = "/work/callers.txt"


def owner(ea):
    f = ida_funcs.get_func(ea)
    return ida_funcs.get_func_name(f.start_ea) if f else "（無函式）"


def main():
    ida_auto.auto_wait()
    want = []
    for line in open(LIST, encoding="utf-8"):
        line = line.split("#")[0].strip()
        if not line:
            continue
        parts = line.split(None, 1)
        want.append((int(parts[0], 16), parts[1] if len(parts) > 1 else ""))

    # 全 segment 逐條掃：call/jmp 的目標，以及立即值。
    hits = {ea: {"call": [], "imm": [], "data": []} for ea, _ in want}
    offs = {}
    for ea, _ in want:
        seg = ida_segment.getseg(ea)
        offs[ea] = (ea - seg.start_ea) if seg else None

    nseg = 0
    ninsn = 0
    for s in idautils.Segments():
        seg = ida_segment.getseg(s)
        if seg.type not in (ida_segment.SEG_CODE, ida_segment.SEG_NORM):
            continue
        nseg += 1
        ea = seg.start_ea
        while ea < seg.end_ea:
            insn = ida_ua.insn_t()
            n = ida_ua.decode_insn(insn, ea)
            if n <= 0:
                ea = ida_bytes.next_head(ea, seg.end_ea)
                if ea == idc.BADADDR:
                    break
                continue
            ninsn += 1
            mnem = insn.get_canon_mnem()
            for op in insn.ops:
                if op.type in (ida_ua.o_near, ida_ua.o_far):
                    # ⚠ 16-bit real mode 的 near call，`op.addr` 是**段內
                    # offset**，不是 linear address。拿它直接跟 linear 目標
                    # 比對永遠不成立——而那是一個**假零**：正對照
                    # （已知有 3 個呼叫者的 sub_1DD22）也會回 0。
                    tgt = op.addr
                    cand = {tgt, seg.start_ea + tgt}
                    if mnem in ("call", "jmp"):
                        for t in hits:
                            if t in cand:
                                hits[t]["call"].append((ea, mnem))
                elif op.type == ida_ua.o_imm:
                    for t, off in offs.items():
                        if off is not None and (op.value & 0xFFFF) == off:
                            hits[t]["imm"].append((ea, mnem))
            ea += n

    # 第四層：資料段裡存著這個 offset 的 word（跳表／函式指標）。
    # 指令層的立即值掃不到「編譯期就寫死在表裡」的那一種——
    # `funcs_131E8` 那張跳表就是這個形狀，而 xref 只在 IDA 已經把
    # 那一格標成 offset 時才建得起來。
    for s_ea in idautils.Segments():
        seg2 = ida_segment.getseg(s_ea)
        blob = ida_bytes.get_bytes(seg2.start_ea, seg2.end_ea - seg2.start_ea)
        if not blob:
            continue
        for t, off in offs.items():
            if off is None:
                continue
            pat = bytes((off & 0xFF, (off >> 8) & 0xFF))
            i = blob.find(pat)
            while i != -1:
                hits[t]["data"].append(seg2.start_ea + i)
                i = blob.find(pat, i + 1)

    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("probe: %d 段 / %d 條指令, sha256=%s\n"
                 % (nseg, ninsn, ida_nalt.retrieve_input_file_sha256().hex()[:16]))
        for target, note in want:
            fh.write("\n==== %08X %s ====\n" % (target, note))
            fh.write("  所屬：%s\n" % owner(target))
            xr = list(idautils.CodeRefsTo(target, 0)) + \
                list(idautils.CodeRefsTo(target, 1))
            fh.write("  -- IDA code xref（%d）--\n" % len(set(xr)))
            for x in sorted(set(xr)):
                fh.write("     %08X  %-12s %s\n"
                         % (x, owner(x), idc.GetDisasm(x).strip()))
            c = hits[target]["call"]
            fh.write("  -- 全 segment 掃到的 call／jmp（%d）--\n" % len(c))
            for x, mnem in c:
                fh.write("     %08X  %-12s %s\n"
                         % (x, owner(x), idc.GetDisasm(x).strip()))
            i = hits[target]["imm"]
            fh.write("  -- 立即值 ＝ offset %s（%d）--\n"
                     % ("%04X" % offs[target] if offs[target] is not None else "?",
                        len(i)))
            for x, mnem in i:
                fh.write("     %08X  %-12s %s\n"
                         % (x, owner(x), idc.GetDisasm(x).strip()))
            d = hits[target]["data"]
            fh.write("  -- 資料段裡的 word ＝ offset（%d，含偶然撞號）--\n"
                     % len(d))
            for x in d[:40]:
                fh.write("     %08X  %-12s %s\n"
                         % (x, idc.get_name(x) or "-", idc.GetDisasm(x).strip()))
            if len(d) > 40:
                fh.write("     …（另 %d 筆）\n" % (len(d) - 40))
    ida_pro.qexit(0)


main()
