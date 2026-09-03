#!/usr/bin/env python3
# IDAPython：全庫掃「對記憶體做位元運算，而立即值含指定位元」的指令。
#
#   printf '20 顯示格旗標 bit5\n40 顯示格旗標 bit6\n' \
#     > workplace/ida/dosv/census/bitflag_list.txt
#   tools/ida.sh script dosv tools/ida_bitflag_users.py KI.EXE.i64
#
# 為什麼要有這一支：`or byte ptr [si], 20h` 這種**段內欄位的位元設定**
# 既沒有交叉參考（基址暫存器 + 位移不在 xref 裡），grep 反組譯文字又會
# 漏掉換了寫法的（`or [bx],20h`／`or es:[di],20h`／`or dl,20h` 後再寫回）。
# `tools/ida_disp_users.py` 是照**位移**掃，這一支是照**立即值的位元**掃——
# 問的是「誰設這個旗標」，不是「誰碰這個欄位」。
#
# 兩種都收：
#   - 目標是記憶體 → 直接寫旗標
#   - 目標是暫存器 → 可能是先讀進來改再寫回，附近幾條指令一起印出來判斷
#
# ⭐ 逐 **segment** 掃，不是逐函式。IDA 靠 xref 建函式，**沒有呼叫端的
# 常式就不會變成函式**，於是 `idautils.Functions()` 一條都看不到它。
# 2026-09-03：顯示格旗標 bit 5 的**唯一**設定端（`0x1D98B`）與清除端
# （`0x1D9AF`）正是這種情形——逐函式版本掃出「沒有人設 bit 5」，
# 而那是**假零**。所屬函式印成「（無函式）」，那一欄本身就是線索。
#
# 輸出 /work/bitflag_users.txt。第一行是 probe（函式數 ＋ 輸入檔雜湊）；
# 沒有 probe 就分不出「沒找到」與「沒跑到」。
import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import ida_segment
import ida_ua
import idautils
import idc

LIST = "/work/bitflag_list.txt"
OUT = "/work/bitflag_users.txt"

# 位元運算：只有這些會「設／清／測」旗標。`mov` 也收——把整個旗標
# 一次寫成含該位元的常數也是設定端。
BITOPS = {"or", "and", "test", "xor", "mov", "btr", "bts"}
MEMTYPES = {ida_ua.o_displ, ida_ua.o_phrase, ida_ua.o_mem}


def ctx(head, before=2, after=2):
    """前後幾條指令，判斷「讀進暫存器改完寫回」用。"""
    out = []
    ea = head
    for _ in range(before):
        prev = idc.prev_head(ea)
        if prev == idc.BADADDR:
            break
        out.insert(0, prev)
        ea = prev
    out.append(head)
    ea = head
    for _ in range(after):
        nxt = idc.next_head(ea, idc.BADADDR)
        if nxt == idc.BADADDR:
            break
        out.append(nxt)
        ea = nxt
    return out


def main():
    ida_auto.auto_wait()
    want = []
    for line in open(LIST, encoding="utf-8"):
        line = line.split("#")[0].strip()
        if not line:
            continue
        parts = line.split(None, 1)
        want.append((int(parts[0], 16), parts[1] if len(parts) > 1 else ""))

    rows = []
    nfunc = len(list(idautils.Functions()))
    nseg = 0
    ninsn = 0
    for s_ea in idautils.Segments():
        seg = ida_segment.getseg(s_ea)
        if seg.type not in (ida_segment.SEG_CODE, ida_segment.SEG_NORM):
            continue
        nseg += 1
        head = seg.start_ea
        while head < seg.end_ea:
            insn = ida_ua.insn_t()
            n = ida_ua.decode_insn(insn, head)
            if n <= 0:
                nxt = ida_bytes.next_head(head, seg.end_ea)
                if nxt == idc.BADADDR:
                    break
                head = nxt
                continue
            ninsn += 1
            mnem = insn.get_canon_mnem()
            if mnem in BITOPS:
                imm = None
                for op in insn.ops:
                    if op.type == ida_ua.o_imm:
                        imm = op.value & 0xFFFF
                        break
                if imm is not None:
                    f = ida_funcs.get_func(head)
                    fn = (ida_funcs.get_func_name(f.start_ea) if f
                          else "（無函式）")
                    dst = insn.ops[0]
                    kind = "記憶體" if dst.type in MEMTYPES else "暫存器"
                    rows.append((imm, fn, head, mnem, kind,
                                 idc.GetDisasm(head)))
            head += n

    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("probe: %d 支函式 / %d 段 / %d 條指令, sha256=%s, 遮罩 %s\n"
                 % (nfunc, nseg, ninsn,
                    ida_nalt.retrieve_input_file_sha256().hex()[:16],
                    "、".join("%02X" % m for m, _ in want)))
        for mask, note in want:
            sel = [r for r in rows if (r[0] & mask) == mask and r[0] <= 0xFF]
            fh.write("\n==== 遮罩 %02Xh（%s）：%d 處 ====\n"
                     % (mask, note, len(sel)))
            for imm, fn, head, mnem, kind, dis in sel:
                fh.write("  %-12s %08X %s %-4s imm=%02Xh  %s\n"
                         % (fn, head, kind, mnem, imm, dis.strip()))
                if kind == "暫存器":
                    for c in ctx(head):
                        tag = ">>" if c == head else "  "
                        fh.write("        %s %08X  %s\n"
                                 % (tag, c, idc.GetDisasm(c).strip()))
    ida_pro.qexit(0)


main()
