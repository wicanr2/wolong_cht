#!/usr/bin/env python3
"""IDA：原事件場景入口與尚未還原的直接 callee，保留原始定位。"""
import hashlib
import json
import sys
import struct
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import idautils

ROOTS = (0x13B08, 0x13D09, 0x13D45, 0x13D68, 0x11D46,
         0x1D46A, 0x1D483, 0x1D4C7, 0x1D615, 0x1D66A, 0x1D76B,
         0x1D782, 0x1D796, 0x1D7E7, 0x1D804)
try:
    ida_auto.auto_wait()
    ledger = json.loads(Path('/evidence/c-recovery-status.json').read_text(encoding='utf-8'))
    known = {r['ida_linear'] for r in ledger['functions'].values()}
    def callees(at):
        return sorted({x.to for start, end in idautils.Chunks(at)
                       for ea in idautils.Heads(start, end)
                       for x in idautils.XrefsFrom(ea)
                       if x.type in (16, 17)})
    extra = {target for at in ROOTS for target in callees(at) if target not in known}
    assert all(ida_funcs.get_func(at) and ida_funcs.get_func(at).start_ea == at
               for at in {*ROOTS, *extra})
    probe.TARGETS = tuple(sorted({*ROOTS, *extra}))
    probe.main()
    path = Path('/output/ida-probe.json')
    report = json.loads(path.read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    assert report['input_sha256'] == ledger['input_sha256']
    header = struct.unpack_from('<H', raw, 8)[0] * 16
    count = struct.unpack_from('<H', raw, 6)[0]
    table = struct.unpack_from('<H', raw, 24)[0]
    relocations = [header + off + seg * 16 for off, seg in
                   (struct.unpack_from('<HH', raw, table + i * 4) for i in range(count))]
    report.update(mz_header_size=header, relocation_file_offsets=relocations,
                  ida_load_paragraph=0x1000)
    for target in report['targets']:
        parts = []
        for chunk in target['chunks']:
            at = chunk['start'] - 0x10000 + 512
            data = raw[at:at + chunk['end'] - chunk['start']]
            expected = bytearray(data); applied = []
            for file_at in relocations:
                if not at <= file_at < at + len(data): continue
                assert file_at + 2 <= at + len(data)
                old = struct.unpack_from('<H', raw, file_at)[0]
                shifted = (old + 0x1000) & 0xffff
                struct.pack_into('<H', expected, file_at - at, shifted)
                applied.append({'file_offset': file_at, 'original_word': old,
                                'ida_word': shifted})
            assert expected == bytes.fromhex(chunk['bytes'])
            chunk.update(file_bytes=data.hex(), loader_relocations=applied); parts.append(data)
        target['file_sha256'] = hashlib.sha256(b''.join(parts)).hexdigest()
        target['direct_callees'] = callees(target['ida_linear'])
        target['unrestored_direct_callees'] = [at for at in target['direct_callees'] if at not in known]
    report['closure_scope'] = 'Scene/resume and map-cell roots plus one layer of unrestored direct call xrefs; indirect dispatch is not covered'
    report['recovery_ledger_sha256'] = hashlib.sha256(Path('/evidence/c-recovery-status.json').read_bytes()).hexdigest()
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
