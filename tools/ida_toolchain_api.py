#!/usr/bin/env python3
# Minimal follow-up: enumerate IDA 9.4 compiler-id APIs (pure discovery).
# Output: /output/ida-probe.json (probe gate requires it).
import json
import traceback

import ida_auto
import ida_nalt
import ida_pro
import idautils
import idc

OUT_JSON = "/output/ida-probe.json"


def main():
    ida_auto.auto_wait()
    res = {"probe": {
        "functions": len(list(idautils.Functions())),
        "input_sha256": ida_nalt.retrieve_input_file_sha256().hex()}}
    res["idc_inf_comp"] = [x for x in dir(idc)
                            if x.startswith("INF_") and "COMP" in x]
    res["idc_inf_cc"] = [x for x in dir(idc) if x.startswith("INF_CC")]
    try:
        import ida_ida  # noqa: PLC0415
        res["ida_ida_comp"] = [x for x in dir(ida_ida)
                               if "compil" in x.lower()]
        res["ida_ida_cc"] = [x for x in dir(ida_ida)
                             if x.startswith("inf_get_cc")]
    except Exception as e:  # noqa: BLE001
        res["ida_ida_error"] = str(e)
    try:
        import idaapi  # noqa: PLC0415
        cands = {}
        if hasattr(idaapi, "get_inf_structure"):
            try:
                inf = idaapi.get_inf_structure()
                cands["get_inf_structure.compiler_id"] = getattr(
                    inf, "compiler_id", "NOATTR")
            except Exception as e:  # noqa: BLE001
                cands["get_inf_structure_error"] = str(e)
        else:
            cands["get_inf_structure"] = "MISSING"
        res["idaapi_cands"] = cands
    except Exception as e:  # noqa: BLE001
        res["idaapi_error"] = str(e)
    with open(OUT_JSON, "w", encoding="utf-8") as fh:
        json.dump(res, fh, indent=2)
        fh.write("\n")
    ida_pro.qexit(0)


try:
    main()
except Exception:  # noqa: BLE001
    with open("/output/ida-probe-error.txt", "w", encoding="utf-8") as fh:
        fh.write(traceback.format_exc())
    ida_pro.qexit(1)
