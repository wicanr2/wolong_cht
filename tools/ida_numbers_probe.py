#!/usr/bin/env python3
"""數字 raster 與日期／訊息參數入口的固定 IDA 證據。"""
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

probe.TARGETS = (0x1062F, 0x1069A, 0x106DE, 0x10CAC)
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
    for start, end in [(0x10984, 0x109AF), (0x11E17, 0x11E46)]:
        rows = []
        ea = start
        name = idc.get_name(start)
        while ea < end:
            insn = ida_ua.insn_t()
            size = ida_ua.decode_insn(insn, ea)
            assert size > 0 and ea + size <= end
            before = ida_bytes.get_bytes(ea, size)
            assert before == raw[ea - 0x10000 + 512:ea - 0x10000 + 512 + size]
            row = probe.instruction(ea)
            row.update(original_name_before_analysis=idc.get_name(ea),
                       ida_code_classified_before=ida_bytes.is_code(ida_bytes.get_full_flags(ea)))
            rows.append(row)
            ea += size
        report['decoded_blocks'].append({'ida_linear': start, 'end': end,
                                        'original_name_before_analysis': name,
                                        'instructions': rows})
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
