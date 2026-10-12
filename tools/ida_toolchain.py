#!/usr/bin/env python3
# IDAPython toolchain probe for KI.EXE (runs via tools/ida.sh probe).
#
# Answers from inside the database, not from re-reading bytes:
#   probe     functions count + input SHA-256 (empty/zero means script died)
#   compiler  IDA's own compiler/cc/procname fields (each API guarded;
#             missing API is recorded as data, never a silent abort)
#   flirt     FUNC_LIB/FUNC_LUMINA counts + every non-sub_ function name.
#             A compiled C binary shows runtime/library names here; a
#             hand-written one shows (almost) only sub_XXXX.
#   entry     entry point + first instructions of the entry function
#   segments  segment list with names
#   imports   import-module count (MZ EXE: expect 0)
#
# Output: /output/ida-probe.json (gate: must exist and be non-empty)
#         /output/ida-toolchain.txt (human readable)
# ASCII only in outputs (container locale is not UTF-8).
import json
import traceback

import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import ida_segment
import idautils
import idc

OUT_JSON = "/output/ida-probe.json"
OUT_TXT = "/output/ida-toolchain.txt"


def attempt(label, fn):
    try:
        return {"ok": True, "value": fn()}
    except Exception as e:  # noqa: BLE001 - headless must not die silently
        return {"ok": False, "error": "%s: %s" % (type(e).__name__, e)}


def compiler_info():
    # IDA 9.4 names (verified by tools/ida_toolchain_api.py discovery probe):
    # INF_COMPILER (compiler_info_t), INF_CC_ID, INF_CC_CM, INF_CC_SIZE_*.
    # get_inf_structure() does not exist in 9.4.
    info = {}
    info["procname"] = attempt(
        "PROCNAME", lambda: idc.get_inf_attr(idc.INF_PROCNAME))
    info["filetype"] = attempt(
        "FILETYPE", lambda: idc.get_inf_attr(idc.INF_FILETYPE))

    def compiler_struct():
        ci = idc.get_inf_attr(idc.INF_COMPILER)
        if ci is None:
            return None
        out = {}
        for k in ("id", "cm", "size_i", "size_b", "size_e"):
            try:
                out[k] = getattr(ci, k)
            except AttributeError:
                pass
        return out

    info["compiler"] = attempt("COMPILER", compiler_struct)
    info["cc_id"] = attempt(
        "CC_ID", lambda: idc.get_inf_attr(idc.INF_CC_ID))
    info["cc_cm"] = attempt(
        "CC_CM", lambda: idc.get_inf_attr(idc.INF_CC_CM))
    info["cc_size_i"] = attempt(
        "CC_SIZE_I", lambda: idc.get_inf_attr(idc.INF_CC_SIZE_I))
    info["cc_size_l"] = attempt(
        "CC_SIZE_L", lambda: idc.get_inf_attr(idc.INF_CC_SIZE_L))

    def comp_name():
        import idaapi  # noqa: PLC0415
        cid = info["cc_id"].get("value")
        if not info["cc_id"].get("ok"):
            return {"note": "no cc id"}
        return {"cc_id": cid,
                "name": idaapi.get_compiler_name(cid)}

    info["compiler_name"] = attempt("COMP_NAME", comp_name)
    return info


def flirt_info():
    total, lib, lumina = 0, 0, 0
    named = []
    sub = 0
    for fea in idautils.Functions():
        total += 1
        f = ida_funcs.get_func(fea)
        name = ida_funcs.get_func_name(f.start_ea) or ""
        flags = f.flags if f else 0
        if flags & ida_funcs.FUNC_LIB:
            lib += 1
        if flags & getattr(ida_funcs, "FUNC_LUMINA", 0):
            lumina += 1
        if name.startswith("sub_"):
            sub += 1
        else:
            named.append({"ea": "0x%X" % f.start_ea, "name": name,
                          "flags": flags})
    return {"total": total, "func_lib": lib, "func_lumina": lumina,
            "sub_named": sub, "non_sub": named}


def entry_info():
    ep = attempt("ENTRY", lambda: idc.get_inf_attr(idc.INF_START_EA))
    out = {"start_ea": ep}
    start_ip = attempt("START_IP", lambda: idc.get_inf_attr(idc.INF_START_IP))
    out["start_ip"] = start_ip
    start_cs = attempt("START_CS", lambda: idc.get_inf_attr(idc.INF_START_CS))
    out["start_cs"] = start_cs
    if ep.get("ok") and ep.get("value") not in (None, idc.BADADDR):
        ea = ep["value"]
        rows = []
        cur, n = ea, 0
        while cur != idc.BADADDR and n < 24:
            rows.append("0x%X  %s" % (cur, idc.GetDisasm(cur)))
            cur = idc.next_head(cur, idc.BADADDR)
            n += 1
        out["entry_disasm"] = rows
    return out


def segments_info():
    segs = []
    for ea in idautils.Segments():
        s = ida_segment.getseg(ea)
        # get_segm_name(seg_t) works; get_segment_name() takes ea, not seg_t.
        segs.append({"name": ida_segment.get_segm_name(s) or "",  # noqa: DEP001
                     "start": "0x%X" % s.start_ea,
                     "end": "0x%X" % s.end_ea,
                     "bitness": s.bitness})
    return segs


def main():
    ida_auto.auto_wait()
    res = {
        "probe": {
            "functions": len(list(idautils.Functions())),
            "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
            "min_ea": "0x%X" % idc.get_inf_attr(idc.INF_MIN_EA),
        },
        "compiler": compiler_info(),
        "flirt": flirt_info(),
        "entry": entry_info(),
        "segments": segments_info(),
        "imports": attempt("IMPORTS",
                           lambda: ida_nalt.get_import_module_qty()),
    }
    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        json.dump(res, fh, indent=2)
        fh.write("\n")
    with open(OUT_TXT, "w", encoding="utf-8") as fh:
        fh.write("functions=%d sha=%s min_ea=%s\n" % (
            res["probe"]["functions"], res["probe"]["input_sha256"],
            res["probe"]["min_ea"]))
        fh.write("compiler=%s\n" % json.dumps(res["compiler"]))
        fh.write("flirt total=%d lib=%d lumina=%d sub=%d non_sub=%d\n" % (
            res["flirt"]["total"], res["flirt"]["func_lib"],
            res["flirt"]["func_lumina"], res["flirt"]["sub_named"],
            len(res["flirt"]["non_sub"])))
        for n in res["flirt"]["non_sub"][:60]:
            fh.write("  named %s %s flags=%d\n" % (n["ea"], n["name"],
                                                   n["flags"]))
        fh.write("segments=%s\n" % json.dumps(res["segments"]))
        fh.write("entry=%s\n" % json.dumps(res["entry"])[:2000])
        fh.write("imports=%s\n" % json.dumps(res["imports"]))
    ida_pro.qexit(0)


try:
    main()
except Exception:  # noqa: BLE001 - last-resort evidence, not silence
    with open("/output/ida-probe-error.txt", "w", encoding="utf-8") as fh:
        fh.write(traceback.format_exc())
    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        fh.write(json.dumps({"probe": {"functions": 0, "fatal": True}}))
    ida_pro.qexit(1)
