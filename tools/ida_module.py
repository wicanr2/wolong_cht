#!/usr/bin/env python3
# IDAPython：模組級掃描。給一個位址區間，把區間內每支函式的證據一次抽齊。
#
# 為什麼要這個：逐支 dump 讀不完 200 多支。但一支函式的「它在做什麼」
# 多半由**它呼叫誰**決定——這一款的 UI 與流程函式都走「設常數 → 呼叫共用常式」，
# 而共用常式的語意已經定案（docs/re/33、docs/re/24 §2）。
# 所以把「呼叫誰 ＋ 傳什麼立即值 ＋ 碰哪張表 ＋ 引用哪些字串」抽出來，
# 大部分函式的角色就能判定，只有少數要逐行讀。
#
# ⚠ 立即值配對是啟發式的（同一個暫存器可能被覆寫、分支可能跳過賦值）。
#   輸出的是**觀察到的組合**，不是「它一定這樣呼叫」。
#
# 區間寫在 /work/module_range.txt，一行「起 迄 名稱」（十六進位 IDA 線性位址）。
# 用法：tools/ida.sh script dosv tools/ida_module.py KI.EXE.i64
# 輸出 /work/module.txt（給人讀）＋ /work/module.json
import json

import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import ida_ua
import ida_xref
import idautils
import idc

RANGE = "/work/module_range.txt"
OUT_TXT = "/work/module.txt"
OUT_JSON = "/work/module.json"

# ⚠ 不要把 0x4200（城兵臨時軍團的基址）放進來：`mov ax, 4200h` 也是
#   DOS 的 LSEEK 功能碼（AH=42h），檔案 I/O 常式全部會被誤標成碰軍團表。
#   基址與 DOS 功能碼撞號時，以「這個立即值還可能是什麼」為準捨棄。
BASES = {
    0x0840: "據點", 0x0842: "據點名", 0x2240: "軍團", 0x2242: "軍團+2",
    0x4240: "武將", 0x4242: "武將名",
    0x4257: "武將職務", 0x425E: "武將說話型",
}
# 已定案語意的共用常式 → (標籤, 帶訊息索引的暫存器)
KNOWN = {
    "sub_18853": ("狀態列", "cx"), "sub_18810": ("訊息", "cx"),
    "sub_193E9": ("選單", "cx"), "sub_13B08": ("君主反應訊息", "cx"),
    "sub_1895D": ("繪框", "al"), "sub_189A4": ("訊息框", "al"),
    "sub_1E3D7": ("寫熱區", "al"), "sub_1E453": ("讀熱區", None),
    "sub_106FD": ("畫三字名稱", None), "sub_106F5": ("畫變長字串", None),
    "sub_1062F": ("印數字", None), "sub_107D2": ("畫肖像", None),
    "sub_188B0": ("畫勢力名", None), "sub_1FA37": ("圖塊blit", "ax"),
    "sub_121E7": ("等按鍵", None), "sub_12078": ("游標存", None),
    "sub_120D6": ("游標復", None), "sub_15E80": ("重畫狀態畫面", "al"),
    "sub_17400": ("選據點", None), "sub_17663": ("選武將", None),
    "sub_17906": ("選勢力", None), "sub_1716D": ("選軍團", None),
    "sub_1820E": ("一覽表選取", None), "sub_181C0": ("開視窗", None),
    "sub_1ECE0": ("亂數", None), "sub_20000": ("系統服務", "ax"),
    "sub_1075B": ("排版繪訊息", "cx"),
}
ARGS = ("ax", "al", "ah", "bx", "bl", "bh", "cx", "cl", "ch", "dx", "dl", "si", "di")


def scan(f):
    """回傳 (呼叫序列, 碰到的表, INT/PORT, 立即值配對)。"""
    calls, tabs, special, pend = [], set(), [], {}
    ea = f.start_ea
    while ea < f.end_ea:
        insn = ida_ua.insn_t()
        if ida_ua.decode_insn(insn, ea) > 0:
            m = insn.get_canon_mnem()
            for op in insn.ops:
                if op.type == ida_ua.o_void:
                    break
                v = (op.value & 0xFFFF) if op.type == ida_ua.o_imm else (
                    (op.addr & 0xFFFF) if op.type == ida_ua.o_displ else None)
                if v in BASES:
                    tabs.add(BASES[v])
            if m == "mov" and insn.ops[1].type == ida_ua.o_imm:
                r = idc.print_operand(ea, 0)
                if r in ARGS:
                    pend[r] = insn.ops[1].value & 0xFFFF
            elif m in ("int",):
                special.append("INT %02Xh @%08X" % (insn.ops[0].value, ea))
            elif m in ("in", "out"):
                special.append("%s @%08X %s" % (m, ea, idc.GetDisasm(ea)))
            elif m == "call":
                t = idc.print_operand(ea, 0)
                calls.append((t, dict(pend)))
                pend = {}
        ea = idc.next_head(ea, f.end_ea)
    return calls, sorted(tabs), special


def evidence(calls):
    out, seen = [], set()
    for target, args in calls:
        info = KNOWN.get(target)
        if not info:
            continue
        label, reg = info
        piece = label
        if reg and reg in args:
            piece += "=%X" % args[reg]
        if piece not in seen:
            seen.add(piece)
            out.append(piece)
    return out


def main():
    ida_auto.auto_wait()
    try:
        ranges = [l.split(None, 2) for l in open(RANGE, encoding="utf-8")
                  if l.strip() and not l.startswith("#")]
    except OSError:
        ranges = []
    res, txt = [], []
    txt.append("probe: %d 支函式, sha256=%s, 區間 %d 個"
               % (len(list(idautils.Functions())),
                  ida_nalt.retrieve_input_file_sha256().hex()[:16], len(ranges)))
    for r in ranges:
        lo, hi = int(r[0], 16), int(r[1], 16)
        name = r[2].strip() if len(r) > 2 else ""
        txt.append("\n======== %05X–%05X %s ========" % (lo, hi, name))
        for fea in idautils.Functions():
            if not (lo <= fea < hi):
                continue
            f = ida_funcs.get_func(fea)
            fn = ida_funcs.get_func_name(fea)
            calls, tabs, special = scan(f)
            callers = []
            x = ida_xref.get_first_cref_to(fea)
            while x != idc.BADADDR:
                g = ida_funcs.get_func(x)
                if g:
                    callers.append(ida_funcs.get_func_name(g.start_ea))
                x = ida_xref.get_next_cref_to(fea, x)
            ev = evidence(calls)
            callees = []
            for t, _ in calls:
                if t.startswith(("sub_", "loc_", "nullsub")) and t not in callees:
                    callees.append(t)
            rec = {"name": fn, "ea": "%08X" % fea, "bytes": f.end_ea - f.start_ea,
                   "callers": sorted(set(callers)), "callees": callees,
                   "tables": tabs, "evidence": ev, "special": special}
            res.append(rec)
            txt.append("%-12s %5dB 呼叫者:%-28s 表:%-16s %s"
                       % (fn, rec["bytes"], ",".join(rec["callers"])[:28] or "—",
                          "、".join(tabs) or "—", " | ".join(ev) or "—"))
            if callees:
                txt.append("             → %s" % " ".join(callees[:14]))
            for s in special:
                txt.append("             ⚑ %s" % s)
    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)
    with open(OUT_TXT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(txt))
    ida_pro.qexit(0)


main()
