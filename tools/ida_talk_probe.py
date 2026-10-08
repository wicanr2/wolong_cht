#!/usr/bin/env python3
"""TALK、七項參數分派與肖像讀檔閉包的固定 IDA 證據。"""
import hashlib
import json
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_bytes
import ida_nalt
import ida_pro
import ida_ua
import idc

probe.TARGETS = (0x106F9, 0x1075B, 0x107D2, 0x1084A, 0x18853, 0x189A4, 0x1E38C, 0x1F4DF)
try:
    probe.main()
    path = Path('/output/ida-probe.json')
    report = json.loads(path.read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    assert report['input_sha256'] == 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    for target in report['targets']:
        parts = []
        for chunk in target['chunks']:
            start = chunk['start'] - 0x10000 + 512
            data = raw[start:start + chunk['end'] - chunk['start']]
            assert data == bytes.fromhex(chunk['bytes'])
            chunk['file_bytes'] = data.hex()
            parts.append(data)
        target['file_sha256'] = hashlib.sha256(b''.join(parts)).hexdigest()
    report['decoded_blocks'] = []
    bounds = [0x108B2, 0x108DB, 0x10904, 0x10939, 0x1095B, 0x1097E, 0x10984]
    for start, end in zip(bounds, bounds[1:]):
        rows = []
        ea = start
        while ea < end:
            insn = ida_ua.insn_t()
            size = ida_ua.decode_insn(insn, ea)
            assert size > 0 and ea + size <= end
            assert ida_bytes.get_bytes(ea, size) == raw[ea - 0x10000 + 512:ea - 0x10000 + 512 + size]
            row = probe.instruction(ea)
            row.update(original_name_before_analysis=idc.get_name(ea),
                       ida_code_classified_before=ida_bytes.is_code(ida_bytes.get_full_flags(ea)))
            rows.append(row)
            ea += size
        report['decoded_blocks'].append({'ida_linear': start, 'end': end,
                                        'original_name_before_analysis': idc.get_name(start),
                                        'instructions': rows})
    table = raw[0xAA4:0xAB2]
    assert ida_bytes.get_bytes(0x108A4, 14) == table
    report['marker_table'] = {'ida_linear': 0x108A4, 'file_offset': 0xAA4,
                             'file_bytes': table.hex(),
                             'offsets': [int.from_bytes(table[i:i + 2], 'little') for i in range(0, 14, 2)]}
    report['portrait_filename'] = {'ida_linear': 0x10D79, 'file_offset': 0xF79,
                                   'file_bytes': raw[0xF79:0xF84].hex()}
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
