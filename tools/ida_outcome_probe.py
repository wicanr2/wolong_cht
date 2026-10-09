#!/usr/bin/env python3
"""IDA：自動戰鬥、戰後退卻與據點易主的原始閉包。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_nalt, ida_pro, idautils

ROOTS = (0x15130, 0x14CF3)
NAVIGATION_ROOTS = (0x14A7B, 0x14ADE, 0x125A3)


try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    known = {x['ida_linear'] for x in list(ledger['functions'].values())+list(ledger['code_blocks'].values())}
    recovery=set(ROOTS)
    pending=list(ROOTS)
    unresolved=[]
    while pending:
        ea=pending.pop()
        f=ida_funcs.get_func(ea)
        assert f and f.start_ea==ea,hex(ea)
        for a,b in idautils.Chunks(ea):
            for at in idautils.Heads(a,b):
                for x in idautils.XrefsFrom(at):
                    if x.type not in (16,17) or x.to in known or x.to in recovery:continue
                    other=ida_funcs.get_func(x.to)
                    if not other or other.start_ea!=x.to:
                        unresolved.append({'from':at,'to':x.to,'containing_function':other.start_ea if other else None})
                        continue
                    recovery.add(x.to);pending.append(x.to)
        assert len(recovery)<512
    probe.TARGETS=tuple(sorted(recovery|set(NAVIGATION_ROOTS)))
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
    result['navigation_only_targets']=list(NAVIGATION_ROOTS)
    result['recovery_targets']=sorted(recovery)
    result['unresolved_call_boundaries']=unresolved
    result['closure_scope']='Complete original automatic-combat result and city ownership-transfer direct-call closures. Field/siege battle entry and full army tick remain navigation only; tactical battle and normal-player paths are not claimed.'
    result['recovery_ledger_sha256'] = hashlib.sha256(ledger_path.read_bytes()).hexdigest()
    path.write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
