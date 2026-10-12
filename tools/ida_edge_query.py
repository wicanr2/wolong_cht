#!/usr/bin/env python3
"""查三個殘部邊位址的身分（code/data、所屬函式、前後 bytes）。"""
import json
import traceback
from pathlib import Path

import ida_auto, ida_bytes, ida_funcs, ida_name, ida_pro, ida_ua, idautils

try:
    ida_auto.auto_wait()
    out = {'function_count': sum(1 for _ in idautils.Functions()), 'sites': []}
    for ea in (0x1AB48, 0x1C0C4, 0x1A82D, 0x1C300, 0x15A00):
        f = ida_funcs.get_func(ea)
        insn = ida_ua.insn_t()
        size = ida_ua.decode_insn(insn, ea)
        row = {
            'ea': ea,
            'is_code': bool(ida_bytes.is_code(ida_bytes.get_full_flags(ea))),
            'is_data': bool(ida_bytes.is_data(ida_bytes.get_full_flags(ea))),
            'func': f.start_ea if f else None,
            'name': ida_name.get_name(ea),
            'decode_size': size,
            'bytes16': ida_bytes.get_bytes(ea - 8, 24).hex() if ida_bytes.get_bytes(ea - 8, 24) else None,
            'word': ida_bytes.get_word(ea),
            'xrefs_from': [{'to': x.to, 'type': x.type} for x in idautils.XrefsFrom(ea)],
        }
        out['sites'].append(row)
    Path('/work/edge_query.json').write_text(json.dumps(out, ensure_ascii=False, indent=2), encoding='utf-8')
except Exception:
    Path('/work/edge_query_error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
