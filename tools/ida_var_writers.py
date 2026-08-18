#!/usr/bin/env python3
# IDAPython（跑在 ida-pro-9.4-idapython:py312-v1；tools/ida.sh 依 .py 副檔名自動選）
#
# 問題：某個全域變數「中途被改了」，要找出**誰改的**。
#
# 做法是逐個名字取交叉參考，用 XrefType() 分讀／寫／取位址三類
# （不要比對助憶碼字串，`or [x], 1` 與 `cmp [x], 1` 在文字上都含 x）。
# 每一筆寫入再往前後各印幾條指令，才看得出它是在什麼條件下發生的。
#
# ⚠ 交叉參考只涵蓋**直接**參考：`ptr = &x` 之後透過指標寫入抓不到。
#    所以「取位址」那幾筆要單獨列出來看——寫入數異常少就是這個成因。
#
# 與 `tools/ida_xref.idc` 的分工：那一支一次查**一個**符號、要走 `ida.sh raw`
# （會改寫 .i64），也不印上下文。這一支一次掃一組名字、走 `script`（先複製再跑），
# 而且每一筆寫入都帶前後幾條指令——判斷「它是在什麼條件下被改的」需要那幾行。
#
# 用法：tools/ida.sh script dosv tools/ida_var_writers.py KI.EXE.i64
# 輸出 /work/var-writers.txt（給人讀）＋ /work/var-writers.json
import json

import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_xref
import idautils
import idc

OUT_JSON = "/work/var-writers.json"
OUT_TXT = "/work/var-writers.txt"

# 要追的全域。戰術鏡頭與它的 dirty flag（docs/re/60 §9、docs/spec/57）
NAMES = [
    "word_1D318",   # 戰場的節拍計數器？
    "word_1D322",   # 側 0 的訊息框到期時刻？
    "word_1D324",   # 側 1
    "byte_1D349",   # 訊息要不要顯示的閘
]

CONTEXT_BEFORE = 6
CONTEXT_AFTER = 3

XT = {0: "unknown", 1: "offset(取位址)", 2: "write(寫)", 3: "read(讀)"}


def fname(ea):
    f = ida_funcs.get_func(ea)
    return ida_funcs.get_func_name(f.start_ea) if f else "?"


def context(ea):
    """往前 CONTEXT_BEFORE 條、往後 CONTEXT_AFTER 條指令。"""
    f = ida_funcs.get_func(ea)
    lo = f.start_ea if f else ea
    hi = f.end_ea if f else ea + 16
    back = []
    cur = ea
    for _ in range(CONTEXT_BEFORE):
        prev = idc.prev_head(cur, lo)
        if prev == idc.BADADDR or prev < lo:
            break
        back.append(prev)
        cur = prev
    rows = []
    for a in reversed(back):
        rows.append("     %08X  %s" % (a, idc.GetDisasm(a)))
    rows.append("  >> %08X  %s" % (ea, idc.GetDisasm(ea)))
    cur = ea
    for _ in range(CONTEXT_AFTER):
        nxt = idc.next_head(cur, hi)
        if nxt == idc.BADADDR or nxt >= hi:
            break
        rows.append("     %08X  %s" % (nxt, idc.GetDisasm(nxt)))
        cur = nxt
    return rows


def typed_refs(ea):
    """用 xrefblk_t 才拿得到 XrefType()。"""
    rows = []
    xb = ida_xref.xrefblk_t()
    ok = xb.first_to(ea, ida_xref.XREF_DATA)
    while ok:
        rows.append((xb.frm, xb.type))
        ok = xb.next_to()
    return rows


def main():
    ida_auto.auto_wait()
    res = {
        "probe": {
            "functions": len(list(idautils.Functions())),
            "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
        },
        "vars": [],
    }
    lines = []
    p = res["probe"]
    lines.append("probe: %d 支函式, sha256=%s"
                 % (p["functions"], p["input_sha256"][:16]))

    for name in NAMES:
        ea = idc.get_name_ea_simple(name)
        entry = {"name": name, "ea": None, "refs": []}
        lines.append("\n================ %s ================" % name)
        if ea == idc.BADADDR:
            lines.append("  找不到這個名字")
            res["vars"].append(entry)
            continue
        entry["ea"] = "%08X" % ea
        lines.append("位址 %08X" % ea)
        refs = typed_refs(ea)
        writes = [r for r in refs if r[1] == 2]
        others = [r for r in refs if r[1] != 2]
        lines.append("  合計 %d 筆參考：寫 %d、其他 %d"
                     % (len(refs), len(writes), len(others)))
        for frm, t in refs:
            entry["refs"].append({
                "from": "%08X" % frm,
                "func": fname(frm),
                "type": XT.get(t, str(t)),
                "text": idc.GetDisasm(frm),
            })
        lines.append("\n-- 寫入 --")
        for frm, _t in writes:
            lines.append("  [%s] %08X  %s" % (fname(frm), frm, idc.GetDisasm(frm)))
            lines.extend(context(frm))
            lines.append("")
        lines.append("-- 讀取／取位址 --")
        for frm, t in others:
            lines.append("  [%s] %08X  %-40s (%s)"
                         % (fname(frm), frm, idc.GetDisasm(frm), XT.get(t, t)))
        res["vars"].append(entry)

    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)
    with open(OUT_TXT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(lines) + "\n")
    ida_pro.qexit(0)


main()
