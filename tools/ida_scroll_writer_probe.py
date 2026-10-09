#!/usr/bin/env python3
"""檢查 nullsub_5 與近旁 operand 的直接讀／寫／取址，保留間接限制。"""
import json
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_bytes
import ida_funcs
import ida_pro
import ida_ua
import idautils
import idc

probe.TARGETS = (0x105ED, 0x100DF)
try:
    probe.main()
    path = Path('/output/ida-probe.json')
    report = json.loads(path.read_text(encoding='utf-8'))
    report['xref_queries'] = []
    for ea in [0x105ED, 0x105EE, 0x1036D, 0x19465, 0x19467]:
        xrefs = [{'from': x.frm, 'type': x.type, 'assembly': idc.generate_disasm_line(x.frm, 0)}
                 for x in idautils.XrefsTo(ea)]
        report['xref_queries'].append({'ida_linear': ea, 'original_name': idc.get_name(ea),
                                       'bytes': ida_bytes.get_bytes(ea, 3).hex(), 'xrefs': xrefs})
    report['segment_candidates'] = []
    for seg in idautils.Segments():
        import ida_segment
        end = ida_segment.getseg(seg).end_ea
        for ea in idautils.Heads(seg, end):
            if not ida_bytes.is_code(ida_bytes.get_full_flags(ea)):
                continue
            insn = ida_ua.insn_t()
            if ida_ua.decode_insn(insn, ea) <= 0:
                continue
            for i, op in enumerate(insn.ops):
                if op.type == ida_ua.o_void:
                    break
                value = op.value if op.type == ida_ua.o_imm else op.addr
                if op.type in [ida_ua.o_imm, ida_ua.o_mem, ida_ua.o_displ] and value & 65535 in range(0x5EC, 0x5F1):
                    report['segment_candidates'].append({'ida_linear': ea, 'operand_index': i,
                                                         'function': ida_funcs.get_func_name(ea),
                                                         'instruction': probe.instruction(ea)})
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
