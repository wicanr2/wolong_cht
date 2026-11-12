#!/usr/bin/env python3
"""Export targeting and formation command evidence without renaming symbols.

Run tools/ida.sh probe dosv tools/ida_tactical_targeting.py OUTPUT.
Use a distinct output directory to preserve earlier probe receipts.
"""
import json
import os
from pathlib import Path

import ida_auto
import ida_bytes
import ida_funcs
import ida_kernwin
import ida_nalt
import ida_pro
import idautils
import idc


def main():
    ida_auto.auto_wait()
    regions = []
    for start, end in [(0x1A754, 0x1A8DE), (0x1A92E, 0x1A953), (0x1AA2C, 0x1AACF)]:
        instructions = [{
            "address": hex(ea),
            "bytes": (ida_bytes.get_bytes(ea, idc.get_item_size(ea)) or b"").hex(),
            "function": ida_funcs.get_func_name(ea),
            "instruction": idc.GetDisasm(ea),
            "data_refs": [hex(x) for x in idautils.DataRefsFrom(ea)],
        } for ea in idautils.Heads(start, end)]
        regions.append({"start": hex(start), "end": hex(end), "instructions": instructions})
    result = {
        "tool": "IDA Pro " + ida_kernwin.get_kernel_version(),
        "address_space": "IDA DOS/V linear",
        "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
        "function_count": ida_funcs.get_func_qty(),
        "uid": os.getuid(), "gid": os.getgid(),
        "regions": regions,
        "xrefs": {hex(target): [{"from": hex(x.frm), "type": x.type}
                                 for x in idautils.XrefsTo(target)]
                  for target in [0x1A85B, 0x1AA2C]},
    }
    Path("/output/ida-probe.json").write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    ida_pro.qexit(0)


main()
