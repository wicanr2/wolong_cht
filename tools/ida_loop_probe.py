#!/usr/bin/env python3
"""IDA：主迴圈及其尚未還原直接依賴的原始定位。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_nalt, ida_pro, idautils

ROOTS = (0x11BE0, 0x109D0, 0x11CD0, 0x11E46, 0x159B7,
         0x10CDE, 0x12459, 0x125A3, 0x13EFD, 0x11F0E,
         0x159D0, 0x15AA2, 0x15E4C, 0x161B6)
RECOVERY_ROOTS = (0x109D0, 0x11E46, 0x11F0E, 0x159B7,
                  0x159D0, 0x15AA2, 0x15E4C, 0x161B6)

try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    known = {x['ida_linear'] for x in ledger['functions'].values()}
    for ea in ROOTS:
        f = ida_funcs.get_func(ea)
        assert f and f.start_ea == ea, hex(ea)
    probe.TARGETS = ROOTS
    probe.main()
    path = Path('/output/ida-probe.json')
    result = json.loads(path.read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    assert result['input_sha256'] == ledger['input_sha256']
    header = struct.unpack_from('<H', raw, 8)[0] * 16
    count, table = struct.unpack_from('<H', raw, 6)[0], struct.unpack_from('<H', raw, 24)[0]
    relocations = [header + off + segment * 16 for off, segment in
                   (struct.unpack_from('<HH', raw, table + i * 4) for i in range(count))]
    result.update(mz_header_size=header, relocation_file_offsets=relocations, ida_load_paragraph=0x1000)
    for target in result['targets']:
        parts = []
        for chunk in target['chunks']:
            start = chunk['start'] - 0x10000 + header
            data = raw[start:start + chunk['end'] - chunk['start']]
            expected, normalized = bytearray(data), []
            for at in relocations:
                if start <= at < start + len(data):
                    assert at + 2 <= start + len(data)
                    original = struct.unpack_from('<H', raw, at)[0]
                    loaded = (original + 0x1000) & 65535
                    struct.pack_into('<H', expected, at - start, loaded)
                    normalized.append({'file_offset': at, 'original_word': original, 'ida_word': loaded})
            assert expected == bytes.fromhex(chunk['bytes'])
            chunk.update(file_bytes=data.hex(), loader_relocations=normalized)
            parts.append(data)
        callees = sorted({x.to for a, b in idautils.Chunks(target['ida_linear'])
                          for ea in idautils.Heads(a, b) for x in idautils.XrefsFrom(ea) if x.type in (16, 17)})
        target.update(file_sha256=hashlib.sha256(b''.join(parts)).hexdigest(), direct_callees=callees,
                      unrestored_direct_callees=[ea for ea in callees if ea not in known])
    result['navigation_only_targets'] = [ea for ea in ROOTS if ea not in RECOVERY_ROOTS]
    result['recovery_targets'] = list(RECOVERY_ROOTS)
    result['right_button_table'] = {
        'ida_linear': 0x15A06,
        'bytes': ida_bytes.get_bytes(0x15A06, 52).hex(),
        'entries': [ida_bytes.get_word(0x15A06 + i * 2) for i in range(26)],
        'fallthrough_instruction': probe.instruction(0x159D0),
    }
    result['closure_scope'] = 'Eight original world-interaction, right-button and fade-in functions. Complete 11BE0 and its update scheduler are navigation only. The right-button dispatch table and fallthrough to original nullsub_1 are explicitly retained.'
    result['recovery_ledger_sha256'] = hashlib.sha256(ledger_path.read_bytes()).hexdigest()
    path.write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
