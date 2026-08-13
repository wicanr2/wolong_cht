#!/usr/bin/env python3
# IDAPython（跑在 ida-pro-9.4-idapython:py312-v1；tools/ida.sh 依 .py 副檔名自動選）
#
# 為什麼改用 IDAPython：IDC 少了一半內建函式（`get_func_qty` 就不存在），
# 而缺函式在 headless 底下是**靜默中止** ——已寫的緩衝一起消失，exit code 仍是 0。
# IDAPython 有 idautils／ida_ua，指令與運算元是解碼出來的，不是比對反組譯文字。
#
# 這一支做三件事：
#   probe  函式數與輸入檔身分（沒有這一項就無法分辨「沒找到」與「沒跑到」）
#   A      所有「對 [reg+disp] 寫入立即值」的指令，用解碼後的運算元判定，
#          給定 disp 清單。用來找誰清軍團 +0x00 的位元。
#   B      指定函式的反組譯 ＋ 呼叫者
#
# ⚠ headless 的 print 不進 stdout，一律寫檔；收工要 ida_pro.qexit(0)。
# ⚠ **容器的 locale 不是 UTF-8**：`open(p, "w")` 配 `ensure_ascii=False` 會在
#    寫到第一個中文字時炸，而前面已經寫出去的幾十 KB 會留在磁碟上——
#    看起來像「有輸出但被截斷」。所有寫檔一律明寫 encoding="utf-8"。
#
# 用法：tools/ida.sh script dosv tools/ida_scan.py KI.EXE.i64
# 輸出 /work/scan.json（結構化）＋ /work/scan.txt（給人讀）
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

OUT_JSON = "/work/scan.json"
OUT_TXT = "/work/scan.txt"

# A：想追的結構位移。0 是軍團記錄 +0x00（存在旗標＋狀態位元，docs/re/32 §1.1）
WATCH_DISP = [0x00]
# B：要逐支 dump 的函式
DUMP = [
    (0x15D19, "sub_15D19 小地圖畫單一據點（re/33 §4 未解）"),
    (0x1F1A3, "sub_1F1A3 T4 第一名，6 個呼叫者"),
    (0x15FAA, "sub_15FAA off_159D2[3]"),
    (0x1614A, "sub_1614A off_159D2[0]"),
    (0x15E1E, "sub_15E1E off_159D2[1]"),
]


def fname(ea):
    f = ida_funcs.get_func(ea)
    return ida_funcs.get_func_name(f.start_ea) if f else "?"


def callers(ea):
    out = []
    x = ida_xref.get_first_cref_to(ea)
    while x != idc.BADADDR:
        out.append(fname(x))
        x = ida_xref.get_next_cref_to(ea, x)
    return out


def disasm(start, end, limit=400):
    rows, ea, n = [], start, 0
    while ea != idc.BADADDR and ea < end and n < limit:
        rows.append("%08X  %s" % (ea, idc.GetDisasm(ea)))
        ea = idc.next_head(ea, end)
        n += 1
    return rows


def scan_imm_writes(disps):
    """找『目的地是 [reg+disp]、來源是立即值』的指令。

    用 ida_ua 解碼，不比對反組譯文字——文字會因為 IDA 的顯示設定而變，
    而且 `[si]` 與 `[si+0]` 印出來不一樣但語意相同（o_phrase vs o_displ）。
    """
    hits = []
    want = set(disps)
    # 只有這些會改寫記憶體。不篩的話 `cmp word ptr [bp+0], 0` 也會被算成寫入。
    WRITE = {"mov", "and", "or", "xor", "add", "sub", "adc", "sbb"}
    for fea in idautils.Functions():
        f = ida_funcs.get_func(fea)
        ea = f.start_ea
        while ea < f.end_ea:
            insn = ida_ua.insn_t()
            if ida_ua.decode_insn(insn, ea) > 0:
                d, s = insn.ops[0], insn.ops[1]
                mem = None
                if d.type == ida_ua.o_displ:
                    mem = d.addr
                elif d.type == ida_ua.o_phrase:
                    mem = 0            # [si] 這種沒有位移的形式
                mnem = insn.get_canon_mnem()
                if (mem is not None and mem in want
                        and s.type == ida_ua.o_imm and mnem in WRITE):
                    hits.append({
                        "func": ida_funcs.get_func_name(f.start_ea),
                        "mnem": mnem,
                        "ea": "%08X" % ea,
                        "disp": mem,
                        "imm": s.value & 0xFFFF,
                        "text": idc.GetDisasm(ea),
                    })
            ea = idc.next_head(ea, f.end_ea)
    return hits


def main():
    ida_auto.auto_wait()

    res = {
        # probe：這三個欄位任一為空／為 0 就代表腳本沒真的跑到底
        "probe": {
            "functions": len(list(idautils.Functions())),
            "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
            "min_ea": "%08X" % idc.get_inf_attr(idc.INF_MIN_EA),
        },
        "imm_writes": scan_imm_writes(WATCH_DISP),
        "dumps": [],
    }
    for ea, label in DUMP:
        f = ida_funcs.get_func(ea)
        res["dumps"].append({
            "label": label,
            "ea": "%08X" % ea,
            "bytes": (f.end_ea - f.start_ea) if f else 0,
            "callers": callers(ea),
            "asm": disasm(f.start_ea, f.end_ea) if f else [],
        })

    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)

    with open(OUT_TXT, "w", encoding="utf-8") as fh:
        p = res["probe"]
        fh.write("probe: %d 支函式, sha256=%s, min_ea=%s\n"
                 % (p["functions"], p["input_sha256"][:16], p["min_ea"]))
        fh.write("\n---- A. 對 [reg+disp] 寫立即值（disp ∈ %s）----\n"
                 % [hex(d) for d in WATCH_DISP])
        for h in res["imm_writes"]:
            fh.write("  %-12s %s  %s\n" % (h["func"], h["ea"], h["text"]))
        fh.write("  共 %d 處\n" % len(res["imm_writes"]))
        for d in res["dumps"]:
            fh.write("\n==== %s  %s  bytes=%d ====\n" % (d["label"], d["ea"], d["bytes"]))
            fh.write("呼叫者：%s\n" % " ".join(d["callers"]))
            for r in d["asm"]:
                fh.write("  %s\n" % r)

    ida_pro.qexit(0)


main()
