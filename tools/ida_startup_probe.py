#!/usr/bin/env python3
"""IDA：開機、開新遊戲／讀檔與 idle 主迴圈的閉包。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_name, ida_nalt, ida_pro, ida_ua, idautils, idc

ROOTS = (0x10000, 0x1005B, 0x1006B, 0x10155, 0x10210, 0x10236, 0x1033B,
         0x1034E, 0x10612, 0x11A6E, 0x11AAE, 0x11AC3, 0x11BE0, 0x11CD0)
RECOVERY_ROOTS = ROOTS
EXPANSION_CAP = 96

try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    blocks = ledger.get('code_blocks', {})
    known = {x['ida_linear'] for x in list(ledger['functions'].values()) + list(blocks.values())}
    parents = {ea: ida_funcs.get_func(ea) for ea in ROOTS}
    missing_roots = [hex(ea) for ea, f in parents.items() if f is None or f.start_ea != ea]
    assert not missing_roots, missing_roots
    wanted = {f.start_ea for f in parents.values() if f}
    pending, recovery = list(RECOVERY_ROOTS), set(RECOVERY_ROOTS)
    indirect_call_sites, fallthrough_edges = [], []
    while pending:
        start = pending.pop()
        f = ida_funcs.get_func(start)
        assert f and f.start_ea == start, hex(start)
        for a, b in idautils.Chunks(start):
            heads = list(idautils.Heads(a, b))
            for ea in heads:
                for x in idautils.XrefsFrom(ea):
                    if x.type not in (16, 17) or x.to in known or x.to in recovery:
                        continue
                    other = ida_funcs.get_func(x.to)
                    if other is None or other.start_ea != x.to:
                        indirect_call_sites.append({'from': ea, 'to': x.to, 'type': x.type})
                        continue
                    wanted.add(x.to)
                    recovery.add(x.to)
                    pending.append(x.to)
            if heads:
                last = heads[-1]
                if idc.print_insn_mnem(last) == 'jmp':
                    insn = ida_ua.insn_t()
                    if ida_ua.decode_insn(insn, last) > 0 and insn.ops[0].type == ida_ua.o_near:
                        target_func = ida_funcs.get_func(insn.ops[0].addr)
                        if target_func is not None and target_func.start_ea != start:
                            fallthrough_edges.append(
                                {'from': last, 'to': insn.ops[0].addr,
                                 'target_function': target_func.start_ea})
        assert len(recovery) <= EXPANSION_CAP, len(recovery)
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
    frags = [(0x10043, 0x1005B, 'start exit tail; ends INT 21h AH=4Ch, no RETN'),
             (0x1006A, 0x1006B, 'shared RETN byte; sub_1005B and nonlocal RET land here'),
             (0x10356, 0x1036D, 'INT61 periodic tick callback; cli/pushf, D2C=1, D2D++ saturating, popf/retf')]
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
    assert [len(b['instructions']) for b in result['decoded_blocks']] == [11, 1, 8]
    call_site, table_base, table_words = 0x19037, 0x18F97, 8
    insn = ida_ua.insn_t()
    assert ida_ua.decode_insn(insn, call_site) > 0
    assert idc.print_insn_mnem(call_site) == 'call' and insn.ops[0].type == ida_ua.o_displ
    assert insn.ops[0].addr & 65535 == table_base - 0x10000
    words = [ida_bytes.get_word(table_base + i * 2) for i in range(table_words)]
    assert words == [0x905E, 0x906E, 0x907F, 0x908D, 0x90A3, 0x90C0, 0x912D, 0x9144]
    for w in words:
        other = ida_funcs.get_func(0x10000 + w)
        assert other is not None and other.start_ea == 0x10000 + w and (0x10000 + w) in recovery
    result['fixed_dispatch_tables'] = [{
        'call_ida_linear': call_site, 'table_ida_linear': table_base,
        'words': words, 'entry_count': table_words,
        'all_entries_in_recovery': True,
    }]
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
    far_sites = sorted(boundaries)
    assert len(far_sites) == 15
    result['far_call_sites'] = far_sites
    result['far_call_target'] = 0x20000
    result['recovery_targets'] = sorted(recovery)
    result['navigation_only_targets'] = sorted(wanted - recovery)
    result['indirect_call_sites'] = indirect_call_sites
    result['fallthrough_edges'] = fallthrough_edges
    result['closure_scope'] = ('Original DOS entry, command-line dispatch, new-game/load and idle '
                               'world main loop, expanded through unrestored direct callees. '
                               'Fixtures and full player path remain separate verification targets.')
    result['recovery_ledger_sha256'] = hashlib.sha256(ledger_path.read_bytes()).hexdigest()
    path.write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
