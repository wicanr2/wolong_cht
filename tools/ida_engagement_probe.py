#!/usr/bin/env python3
"""IDA：野戰、攻城與戰術引擎的直接及固定表分派閉包。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_nalt, ida_pro, ida_ua, ida_name, idautils

ROOTS = (0x14A7B, 0x14ADE)
NAVIGATION_ROOTS = (0x125A3,)
RAW_REGIONS = ((0x11B76, 0x11BE0), (0x1A065, 0x1A12A), (0x1C0C5, 0x1C30D),
               (0x1A7E6, 0x1A7E7), (0x1A82C, 0x1A82D))
TABLES = ((0x1A457, 0x1A466, 19), (0x1A597, 0x1A59C, 5),
          (0x1C006, 0x1C048, 31), (0x1C031, 0x1C086, 31),
          (0x1A7E1, 0x1A7E7, 11), (0x1A827, 0x1A82D, 9),
          (0x11A5E, 0x159D2, 5), (0x15FDC, 0x16056, 6),
          (0x1ACF6, 0x1AD01, 4))

try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    known = {x['ida_linear'] for x in list(ledger['functions'].values())+list(ledger['code_blocks'].values())}
    raw_regions = []
    for begin, end in RAW_REGIONS:
        name = ida_name.get_name(begin)
        rows, at = [], begin
        while at < end:
            insn = ida_ua.insn_t()
            size = ida_ua.decode_insn(insn, at)
            assert size > 0 and at+size <= end, hex(at)
            if not ida_bytes.is_code(ida_bytes.get_full_flags(at)):
                ida_bytes.del_items(at, ida_bytes.DELIT_SIMPLE, size)
                assert ida_ua.create_insn(at) == size
            rows.append(probe.instruction(at))
            at += size
        raw_regions.append({'ida_linear': begin, 'end_ida_linear': end,
          'original_name_before_analysis': name, 'instructions': rows,
          'boundary_authority': 'explicit primary decode; adjacent existing function or fixed dispatch-table boundary'})
    tables = [{'call_or_jump_ida_linear': call, 'table_ida_linear': begin, 'entries': count,
               'words': [ida_bytes.get_word(begin+2*i) for i in range(count)]}
              for call, begin, count in TABLES]
    recovery, pending, unresolved = set(ROOTS), list(ROOTS), []
    rejected, indirect, interior = [], [], []
    rows_by_raw = {b['ida_linear']: b['instructions'] for b in raw_regions}
    pending.extend(rows_by_raw)
    visited = set()
    def decode_raw_call(begin):
        name = ida_name.get_name(begin)
        rows, queue, seen = [], [begin], set()
        while queue:
            at = queue.pop()
            if at in seen:
                continue
            seen.add(at)
            decoded = ida_ua.insn_t()
            size = ida_ua.decode_insn(decoded, at)
            assert size > 0 and 0x10000 <= at < 0x20000, hex(at)
            if not ida_bytes.is_code(ida_bytes.get_full_flags(at)):
                ida_bytes.del_items(at, ida_bytes.DELIT_SIMPLE, size)
                assert ida_ua.create_insn(at) == size
            row = probe.instruction(at)
            rows.append(row)
            mn = row['mnemonic']
            if mn in ('retn','retf','iret'):
                continue
            if mn.startswith('j') or mn.startswith('loop'):
                assert row['operands'][0]['type'] == 7, hex(at)
                queue.append(0x10000+row['operands'][0]['addr'])
                if mn == 'jmp':
                    continue
            queue.append(at+size)
            assert len(seen) < 4096
        rows.sort(key=lambda r:r['ida_linear'])
        block = {'ida_linear':begin,'end_ida_linear':max(r['ida_linear']+len(bytes.fromhex(r['bytes'])) for r in rows),
                 'original_name_before_analysis':name,'instructions':rows,
                 'boundary_authority':'primary control-flow decode from original direct CALL; both conditional edges followed through RET'}
        raw_regions.append(block)
        rows_by_raw[begin] = rows
        pending.append(begin)

    def add_target(at, origin):
        if at in known or at in recovery or at in rows_by_raw:
            return
        if any(b['ida_linear'] <= at < b['end_ida_linear'] for b in raw_regions):
            assert any(r['ida_linear'] == at for b in raw_regions for r in b['instructions']), hex(at)
            return
        f = ida_funcs.get_func(at)
        if f and f.start_ea == at:
            recovery.add(at)
            pending.append(at)
        elif f and f.start_ea in recovery:
            assert ida_bytes.is_code(ida_bytes.get_full_flags(at)), hex(at)
            interior.append({'from':origin,'to':at,'containing_function':f.start_ea,
                             'reason':'original table selects an existing instruction inside the recovered function'})
        elif not f:
            decode_raw_call(at)
        else:
            unresolved.append({'from': origin, 'to': at, 'containing_function': f.start_ea if f else None})
    for table in tables:
        if table['call_or_jump_ida_linear'] == 0x1A597:
            continue
        for word in table['words']:
            add_target(word+0x10000, table['call_or_jump_ida_linear'])
    while pending:
        ea = pending.pop()
        if ea in visited:
            continue
        visited.add(ea)
        if ea in rows_by_raw:
            rows = rows_by_raw[ea]
        else:
            f = ida_funcs.get_func(ea)
            assert f and f.start_ea == ea, hex(ea)
            rows = [probe.instruction(at) for a,b in idautils.Chunks(ea) for at in idautils.Heads(a,b)
                    if ida_bytes.is_code(ida_bytes.get_full_flags(at))]
        for row in rows:
            if row['mnemonic'] != 'call':
                continue
            candidates = [x.to for x in idautils.XrefsFrom(row['ida_linear']) if x.type in (16,17)]
            if row['operands'][0]['type'] not in (6,7):
                indirect.append({'ida_linear': row['ida_linear'], 'bytes': row['bytes'],
                                 'original_assembly': row['assembly'], 'xref_candidates': candidates})
                for target in candidates:
                    rejected.append({'from': row['ida_linear'], 'to': target,
                      'reason': 'indirect operand; fixed original table replaces synthetic call xref'})
                continue
            for target in candidates:
                add_target(target, row['ida_linear'])
        assert len(recovery) < 512
    probe.TARGETS = tuple(sorted(recovery | set(NAVIGATION_ROOTS)))
    probe.main()
    path = Path('/output/ida-probe.json')
    result = json.loads(path.read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    assert result['input_sha256'] == ledger['input_sha256']
    header = struct.unpack_from('<H', raw, 8)[0]*16
    count, offset = struct.unpack_from('<H', raw, 6)[0], struct.unpack_from('<H', raw, 24)[0]
    relocations = [header+off+seg*16 for off,seg in
                   (struct.unpack_from('<HH', raw, offset+i*4) for i in range(count))]
    result.update(mz_header_size=header, relocation_file_offsets=relocations, ida_load_paragraph=0x1000)
    for target in result['targets']:
        parts = []
        for chunk in target['chunks']:
            start = chunk['start']-0x10000+header
            data = raw[start:start+chunk['end']-chunk['start']]
            expected, normalized = bytearray(data), []
            for at in relocations:
                if start <= at < start+len(data):
                    assert at+2 <= start+len(data)
                    old = struct.unpack_from('<H', raw, at)[0]
                    loaded = (old+0x1000)&65535
                    struct.pack_into('<H', expected, at-start, loaded)
                    normalized.append({'file_offset': at, 'original_word': old, 'ida_word': loaded})
            assert expected == bytes.fromhex(chunk['bytes'])
            chunk.update(file_bytes=data.hex(), loader_relocations=normalized)
            parts.append(data)
        target['file_sha256'] = hashlib.sha256(b''.join(parts)).hexdigest()
    for block in raw_regions:
        begin, end = block['ida_linear'], block['end_ida_linear']
        data = raw[begin-0x10000+header:end-0x10000+header]
        loaded = b''.join(bytes.fromhex(r['bytes']) for r in block['instructions'])
        expected = bytearray(data)
        normalized = []
        for at in relocations:
            if begin-0x10000+header <= at < end-0x10000+header:
                old = struct.unpack_from('<H', raw, at)[0]
                new = (old+0x1000)&65535
                struct.pack_into('<H', expected, at-(begin-0x10000+header), new)
                normalized.append({'file_offset':at,'original_word':old,'ida_word':new})
        for row in block['instructions']:
            at = row['ida_linear']-begin
            assert bytes.fromhex(row['bytes']) == expected[at:at+len(bytes.fromhex(row['bytes']))]
        block.update(file_sha256=hashlib.sha256(data).hexdigest(),file_bytes=data.hex(),loader_relocations=normalized)
    result.update(navigation_only_targets=list(NAVIGATION_ROOTS), recovery_targets=sorted(recovery),
                  decoded_blocks=raw_regions, fixed_dispatch_tables=tables,
                  unresolved_call_boundaries=unresolved, rejected_indirect_xref_candidates=rejected,
                  indirect_calls=indirect,
                  interior_entry_targets=interior,
                  closure_scope='Field/siege engagement and tactical engine with raw fallthrough/update/dispatch handlers; full army/main scheduling and natural player paths still require verification.',
                  recovery_ledger_sha256=hashlib.sha256(ledger_path.read_bytes()).hexdigest())
    path.write_text(json.dumps(result,ensure_ascii=False,indent=2)+chr(10),encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
