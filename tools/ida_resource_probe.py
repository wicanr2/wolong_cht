#!/usr/bin/env python3
"""完整檔案載入、BGM cue 和原始記憶體分段設定的 IDA 證據。"""
import hashlib
import json
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_nalt
import ida_pro

probe.TARGETS = (0x187AF, 0x1E378, 0x1F4A2, 0x10241, 0x102C2, 0x100DF)
try:
    probe.main()
    path = Path('/output/ida-probe.json')
    report = json.loads(path.read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    assert report['input_sha256'] == 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    for target in report['targets']:
        parts = []
        for chunk in target['chunks']:
            at = chunk['start'] - 0x10000 + 512
            data = raw[at:at + chunk['end'] - chunk['start']]
            assert data == bytes.fromhex(chunk['bytes'])
            chunk['file_bytes'] = data.hex()
            parts.append(data)
        target['file_sha256'] = hashlib.sha256(b''.join(parts)).hexdigest()
    report['resource_names'] = []
    for off in [0xDB8, 0xDC1, 0xDF0, 0xE07]:
        name = raw[off + 512:].split(b'\0', 1)[0]
        report['resource_names'].append({'ida_linear': off + 0x10000, 'file_offset': off + 512,
                                         'name': name.decode('ascii'), 'file_bytes': (name + b'\0').hex()})
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
