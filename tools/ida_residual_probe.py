#!/usr/bin/env python3
"""IDA：殘部群（census 有、還原帳上沒有的 15 支）的 inventory。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_name, ida_nalt, ida_pro, ida_ua, idautils, idc

ROOTS = (0x145F8, 0x15AD1, 0x15AFC, 0x15C14, 0x19377, 0x193D7,
         0x1AB7C, 0x1C300, 0x1E567, 0x1E6FF)

try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    blocks = ledger.get('code_blocks', {})
    known = {x['ida_linear'] for x in list(ledger['functions'].values()) + list(blocks.values())}
    overlap = sorted(set(ROOTS) & known)
    assert not overlap, [hex(ea) for ea in overlap]
    parents = {ea: ida_funcs.get_func(ea) for ea in ROOTS}
    missing_roots = [hex(ea) for ea, f in parents.items() if f is None or f.start_ea != ea]
    assert not missing_roots, missing_roots
    probe.TARGETS = tuple(sorted(ROOTS))
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
                      unrestored_direct_callees=[ea for ea in callees if ea not in known and ea not in set(ROOTS)])
    result['seed_entries'] = [{'ida_linear': ea, 'original_name': ida_name.get_name(ea),
                              'containing_function': parents[ea].start_ea if parents[ea] else None,
                              'xrefs_to': [{'from': x.frm, 'type': x.type} for x in idautils.XrefsTo(ea)]}
                             for ea in ROOTS]
    frags = [(0x1AB48, 0x1AB4A, 'tail jmp EB 32 into sub_1AB7C; code, no function')]
    result['decoded_blocks'] = []
    for fstart, fend, note in frags:
        rows, at = [], fstart
        while at < fend:
            instruction = ida_ua.insn_t()
            size = ida_ua.decode_insn(instruction, at)
            assert size > 0, hex(at)
            if not ida_bytes.is_code(ida_bytes.get_full_flags(at)):
                ida_bytes.del_items(at, ida_bytes.DELIT_SIMPLE, size)
                assert ida_ua.create_insn(at) == size, hex(at)
            rows.append(probe.instruction(at))
            at += size
        assert at == fend, (hex(at), hex(fend))
        frag_bytes = b''.join(bytes.fromhex(row['bytes']) for row in rows)
        file_start = fstart - 0x10000 + header
        assert raw[file_start:file_start + len(frag_bytes)] == frag_bytes
        frag_relocations = []
        for at in relocations:
            if file_start <= at < file_start + len(frag_bytes):
                original = struct.unpack_from('<H', raw, at)[0]
                frag_relocations.append({'file_offset': at, 'original_word': original,
                                         'ida_word': (original + 0x1000) & 65535})
        result['decoded_blocks'].append({
            'ida_linear': fstart, 'end_ida_linear': fend,
            'boundary_authority': note,
            'instructions': rows,
            'file_bytes': frag_bytes.hex(),
            'loader_relocations': frag_relocations,
            'file_sha256': hashlib.sha256(frag_bytes).hexdigest(),
        })
    assert [len(b['instructions']) for b in result['decoded_blocks']] == [1]
    boundaries = []
    for target in result['targets']:
        for chunk in target['chunks']:
            for row in chunk['instructions']:
                if row['mnemonic'] != 'call':
                    continue
                if row['bytes'].startswith('9a'):
                    code = bytes.fromhex(row['bytes'])
                    assert len(code) == 5 and int.from_bytes(code[1:3], 'little') == 0
                    assert int.from_bytes(code[3:5], 'little') == 0x2000
                    assert row['ida_linear'] not in boundaries
                    boundaries.append(row['ida_linear'])
    result['far_call_sites'] = sorted(boundaries)
    result['far_call_target'] = 0x20000
    result['recovery_targets'] = sorted(ROOTS)
    result['closure_scope'] = ('Residual census functions absent from the C recovery ledger: faction-picker UI, '
                               'season/BGM, tactical attack cap, map-node/draw-queue, diplomacy cleanup, '
                               'tactical sidebar hotzone entry, VGA GC leaf helpers, two null subs.')
    result['recovery_ledger_sha256'] = hashlib.sha256(ledger_path.read_bytes()).hexdigest()
    path.write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
