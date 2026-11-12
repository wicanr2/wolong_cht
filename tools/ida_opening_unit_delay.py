#!/usr/bin/env python3
"""Read-only evidence for the opening soldier delay and adjacent object records.

Run through tools/ida.sh probe dosv tools/ida_opening_unit_delay.py OUTPUT.
The wrapper builds a disposable database from the original executable.
"""
import json
import os
from pathlib import Path

import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import ida_kernwin
import idautils
import idc


def main():
    ida_auto.auto_wait()
    regions = []
    for start, end in [(0x19AC0, 0x19CE2), (0x19CE2, 0x19E70)]:
        instructions = []
        for ea in idautils.Heads(start, end):
            instructions.append({
                "address": hex(ea),
                "bytes": (ida_bytes.get_bytes(ea, idc.get_item_size(ea)) or b"").hex(),
                "function": ida_funcs.get_func_name(ea),
                "instruction": idc.GetDisasm(ea),
                "data_refs": [hex(x) for x in idautils.DataRefsFrom(ea)],
            })
        regions.append({"start": hex(start), "end": hex(end), "instructions": instructions})
    # Include the direct initialization sources of the delay byte. The adjacent
    # word store alone cannot establish which allocation receives it.
    refs = {}
    for target in [0x19AF4, 0x19C45, 0x19CE2]:
        refs[hex(target)] = [{"from": hex(x.frm), "type": x.type}
                             for x in idautils.XrefsTo(target)]
    result = {
        "tool": "IDA Pro " + ida_kernwin.get_kernel_version(),
        "address_space": "IDA DOS/V linear",
        "input_sha256": ida_nalt.retrieve_input_file_sha256().hex(),
        "function_count": ida_funcs.get_func_qty(),
        "uid": os.getuid(), "gid": os.getgid(),
        "regions": regions, "xrefs": refs,
    }
    Path("/output/ida-probe.json").write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    ida_pro.qexit(0)


main()
