#!/usr/bin/env python3
"""IDA：主迴圈保存堆疊與非區域退出的原始定位。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_name, ida_nalt, ida_pro, ida_ua, idautils, idc

ROOTS = (0x11BE0, 0x11CB1, 0x10A1C, 0x189F0, 0x1533D, 0x159B7)
RECOVERY_ROOTS = (0x11CB1, 0x10A1C)

try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    known = {x['ida_linear'] for x in list(ledger['functions'].values()) + list(ledger['code_blocks'].values())}
    parents = {ea: ida_funcs.get_func(ea) for ea in ROOTS}
    wanted = {f.start_ea for f in parents.values() if f}
    pending, recovery = list(RECOVERY_ROOTS), set(RECOVERY_ROOTS)
    while pending:
        start = pending.pop()
        f = ida_funcs.get_func(start)
        assert f and f.start_ea == start, hex(start)
        for a, b in idautils.Chunks(start):
            for ea in idautils.Heads(a, b):
                for x in idautils.XrefsFrom(ea):
                    if x.type not in (16, 17) or x.to in known or x.to in recovery:
                        continue
                    other = ida_funcs.get_func(x.to)
                    assert other and other.start_ea == x.to, (hex(ea), hex(x.to))
                    wanted.add(x.to)
                    recovery.add(x.to)
                    pending.append(x.to)
        assert len(recovery) < 32
    probe.TARGETS = tuple(sorted(wanted))
    assert probe.TARGETS
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
    result['seed_entries'] = [{'ida_linear': ea, 'original_name': ida_name.get_name(ea),
                              'containing_function': parents[ea].start_ea if parents[ea] else None,
                              'xrefs_to': [{'from': x.frm, 'type': x.type} for x in idautils.XrefsTo(ea)]}
                             for ea in ROOTS]
    result['recovery_targets'] = sorted(recovery)
    result['navigation_only_targets'] = sorted(wanted - recovery)
    result['right_button_table'] = {
        'ida_linear': 0x15A06,
        'bytes': ida_bytes.get_bytes(0x15A06, 52).hex(),
        'entries': [ida_bytes.get_word(0x15A06 + i * 2) for i in range(26)],
        'fallthrough_instruction': probe.instruction(0x159D0),
    }
    result['nearby_code'] = []
    at = 0x11BE0
    while at < 0x11CD0:
        instruction = ida_ua.insn_t()
        size = ida_ua.decode_insn(instruction, at)
        assert size > 0
        if not ida_bytes.is_code(ida_bytes.get_full_flags(at)):
            ida_bytes.del_items(at, ida_bytes.DELIT_SIMPLE, size)
            assert ida_ua.create_insn(at) == size
        result['nearby_code'].append(probe.instruction(at))
        at += size
    result['closure_scope'] = 'Original containing-function boundaries and nearby control flow for 11BE0 saved main stack and 11CB1 nonlocal exit. No dependency bodies replaced or assumed.'
    result['recovery_ledger_sha256'] = hashlib.sha256(ledger_path.read_bytes()).hexdigest()
    path.write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
