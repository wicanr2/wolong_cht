#!/usr/bin/env python3
# IDAPython：替每支函式標出它碰哪幾張表（DOS/V KI.EXE）。
#
# 這一款把所有戰略資料放在同一個段的固定基址上，位址算術都是
# 「基址 ＋ 編號 × 記錄長度」，而基址是**立即值**。所以掃立即值就能知道
# 一支函式在動哪張表——這是 `si` 指向什麼的最便宜證據，
# 比從呼叫端往回推可靠（呼叫端的參數順序是間接證據，會推錯）。
#
# ⚠ 立即值可能只是巧合的常數。命中是**候選**，要配合該函式的其他證據才算數；
#   一支函式可以同時碰多張表（例如編成要同時動軍團、武將、勢力）。
#
# 用法：tools/ida.sh script dosv tools/ida_tables.py KI.EXE.i64
# 輸出 /work/tables.json ＋ /work/tables.txt
import json

import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_ua
import idautils
import idc

OUT_JSON = "/work/tables.json"
OUT_TXT = "/work/tables.txt"

# 段內基址 → 名稱。出處是 docs/formats/08-sinario-save.md。
# `+2` 的變體是「直接指到名稱欄」的寫法（docs/re/32 §1），一併認出來。
BASES = {
    0x0840: "據點表",
    0x0842: "據點表+2（名稱）",
    0x2240: "軍團表",
    0x2242: "軍團表+2",
    0x4200: "城兵臨時軍團",
    0x4240: "武將表",
    0x4242: "武將表+2（名稱）",
    0x4257: "武將表+0x17（職務）",
    0x425E: "武將表+0x1E（說話類型）",
}
# 想追的位移：對這些位移寫立即值的指令會被單獨列出來（docs/re/34 的位元圖）
WATCH = {0x00}
WRITE = {"mov", "and", "or", "xor", "add", "sub", "adc", "sbb"}


def main():
    ida_auto.auto_wait()
    funcs = {}
    for fea in idautils.Functions():
        f = ida_funcs.get_func(fea)
        name = ida_funcs.get_func_name(f.start_ea)
        tabs, writes = {}, []
        ea = f.start_ea
        while ea < f.end_ea:
            insn = ida_ua.insn_t()
            if ida_ua.decode_insn(insn, ea) > 0:
                mnem = insn.get_canon_mnem()
                for op in insn.ops:
                    if op.type == ida_ua.o_void:
                        break
                    if op.type == ida_ua.o_imm:
                        v = op.value & 0xFFFF
                        if v in BASES:
                            tabs.setdefault(BASES[v], []).append("%08X" % ea)
                    # 位移形式也算：mov al, [bx+4257h] 的 4257h 在 op.addr
                    if op.type == ida_ua.o_displ and (op.addr & 0xFFFF) in BASES:
                        tabs.setdefault(BASES[op.addr & 0xFFFF], []).append("%08X" % ea)
                d, s = insn.ops[0], insn.ops[1]
                mem = (d.addr if d.type == ida_ua.o_displ
                       else (0 if d.type == ida_ua.o_phrase else None))
                if mem in WATCH and s.type == ida_ua.o_imm and mnem in WRITE:
                    writes.append({"ea": "%08X" % ea, "text": idc.GetDisasm(ea)})
            ea = idc.next_head(ea, f.end_ea)
        if tabs or writes:
            funcs[name] = {"ea": "%08X" % f.start_ea,
                           "bytes": f.end_ea - f.start_ea,
                           "tables": {k: v[:4] for k, v in tabs.items()},
                           "watch_writes": writes}

    res = {
        "probe": {
            "functions": len(list(idautils.Functions())),
            "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
            "matched": len(funcs),
        },
        "funcs": funcs,
    }
    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)

    with open(OUT_TXT, "w", encoding="utf-8") as fh:
        p = res["probe"]
        fh.write("probe: %d 支函式, sha256=%s, 命中 %d 支\n\n"
                 % (p["functions"], p["input_sha256"][:16], p["matched"]))
        fh.write("---- 對 +0x00 寫立即值的函式，附它碰的表 ----\n")
        for name, d in sorted(funcs.items(), key=lambda kv: kv[1]["ea"]):
            if not d["watch_writes"]:
                continue
            fh.write("\n%s  %s  bytes=%d\n" % (name, d["ea"], d["bytes"]))
            fh.write("  表：%s\n" % ("、".join(d["tables"]) or "（無立即值線索）"))
            for w in d["watch_writes"]:
                fh.write("    %s  %s\n" % (w["ea"], w["text"]))
    ida_pro.qexit(0)


main()
