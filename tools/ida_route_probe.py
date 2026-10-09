#!/usr/bin/env python3
"""IDA：行軍尋路、自我修改區段與軍團潰散閉包。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_auto, ida_bytes, ida_funcs, ida_nalt, ida_pro, ida_ua, ida_name, idautils

ROOTS = (0x147BB, 0x1487B, 0x1291A, 0x14A0F)
NAVIGATION_ROOTS = (0x125A3, 0x12662, 0x1474A)
RAW_START, RAW_END = 0x1491B, 0x14A7B


try:
    ida_auto.auto_wait()
    ledger_path = Path('/evidence/c-recovery-status.json')
    ledger = json.loads(ledger_path.read_text(encoding='utf-8'))
    raw_name_before=ida_name.get_name(RAW_START)
    raw_helper_name_before=ida_name.get_name(0x14A0F)
    at=RAW_START
    raw_rows=[]
    while at<RAW_END:
        instruction=ida_ua.insn_t();size=ida_ua.decode_insn(instruction,at)
        assert size>0 and at+size<=RAW_END,hex(at)
        if not ida_bytes.is_code(ida_bytes.get_full_flags(at)):
            ida_bytes.del_items(at,ida_bytes.DELIT_SIMPLE,size)
            assert ida_ua.create_insn(at)==size
        raw_rows.append(probe.instruction(at));at+=size
    known = {RAW_START,0x14A0F}|{x['ida_linear'] for x in list(ledger['functions'].values())+list(ledger['code_blocks'].values())}
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
    result['decoded_blocks']=[]
    for a,b,name in [(RAW_START,0x14A0F,raw_name_before)]:
        rows=[x for x in raw_rows if a<=x['ida_linear']<b]
        file_bytes=raw[a-0x10000+header:b-0x10000+header]
        assert b''.join(bytes.fromhex(x['bytes']) for x in rows)==file_bytes
        result['decoded_blocks'].append({'ida_linear':a,'end_ida_linear':b,'original_name_before_analysis':name,
          'boundary_authority':'explicit raw decode through original RET and next named-function boundary; original IDA name preserved',
          'instructions':rows,'file_sha256':hashlib.sha256(file_bytes).hexdigest()})
    result['closure_scope']='Original march replanning, retreat-route wrapper and disband/capture helpers, with raw self-modifying route search and relax helper. Full army tick and battle outcome remain navigation only.'
    result['recovery_ledger_sha256'] = hashlib.sha256(ledger_path.read_bytes()).hexdigest()
    path.write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
