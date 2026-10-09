#!/usr/bin/env python3
"""IDA：君主出陣、六格自動編成與側欄的固定直接呼叫閉包。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path
sys.path.insert(0,'/tools')
import ida_matching_probe as probe
import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_name
import ida_pro
import ida_ua
import idc
import idautils

ROOTS=(0x1461D,0x14698,0x14717,0x15E80,0x15EB7,0x15F27,0x15F5D,
       0x15F7F,0x1699E,0x16E8F,0x16EC9,0x16F26,0x16F86,0x16FD2)
RAW=()
try:
    ida_auto.auto_wait()
    ledger=json.loads(Path('/evidence/c-recovery-status.json').read_text(encoding='utf-8'))
    known={x['ida_linear'] for x in ledger['functions'].values()}
    def calls(at):
        return sorted({x.to for a,b in idautils.Chunks(at) for ea in idautils.Heads(a,b)
                       for x in idautils.XrefsFrom(ea) if x.type in (16,17)})
    wanted=set(ROOTS);pending=list(ROOTS);raw_targets={a for a,b in RAW}
    while pending:
        at=pending.pop()
        assert ida_funcs.get_func(at) and ida_funcs.get_func(at).start_ea==at,hex(at)
        for target in calls(at):
            if target in known or target in wanted or target in raw_targets:continue
            f=ida_funcs.get_func(target)
            if f and f.start_ea==target:wanted.add(target);pending.append(target)
        assert len(wanted)<=64
    probe.TARGETS=tuple(sorted(wanted));probe.main()
    p=Path('/output/ida-probe.json');r=json.loads(p.read_text(encoding='utf-8'))
    raw=Path(ida_nalt.get_input_file_path()).read_bytes()
    assert r['input_sha256']==ledger['input_sha256']
    header=struct.unpack_from('<H',raw,8)[0]*16;count=struct.unpack_from('<H',raw,6)[0];table=struct.unpack_from('<H',raw,24)[0]
    reloc=[header+off+seg*16 for off,seg in (struct.unpack_from('<HH',raw,table+i*4) for i in range(count))]
    r.update(mz_header_size=header,relocation_file_offsets=reloc,ida_load_paragraph=0x1000)
    for t in r['targets']:
        parts=[]
        for c in t['chunks']:
            a=c['start']-0x10000+header;data=raw[a:a+c['end']-c['start']];expected=bytearray(data);rows=[]
            for at in reloc:
                if not a<=at<a+len(data):continue
                assert at+2<=a+len(data);old=struct.unpack_from('<H',raw,at)[0];shift=(old+0x1000)&65535
                struct.pack_into('<H',expected,at-a,shift);rows.append({'file_offset':at,'original_word':old,'ida_word':shift})
            assert expected==bytes.fromhex(c['bytes']);c.update(file_bytes=data.hex(),loader_relocations=rows);parts.append(data)
        t['file_sha256']=hashlib.sha256(b''.join(parts)).hexdigest();t['direct_callees']=calls(t['ida_linear'])
    blocks=[]
    for a,b in RAW:
        rows=[];at=a
        while at<b:
            original_line=idc.generate_disasm_line(at,0)
            originally_code=ida_bytes.is_code(ida_bytes.get_full_flags(at))
            if not originally_code:
                # Decode only in this disposable DB; preserve the original data view.
                insn=ida_ua.insn_t();size=ida_ua.decode_insn(insn,at)
                assert size>0 and at+size<=b
                ida_bytes.del_items(at,ida_bytes.DELIT_SIMPLE,size)
                assert ida_ua.create_insn(at)==size
            row=probe.instruction(at);rows.append(row);at+=len(bytes.fromhex(row['bytes']))
            if not originally_code:
                row['original_ida_data_line']=original_line
                row['originally_code']=False
        assert at==b
        data=raw[a-0x10000+header:b-0x10000+header]
        assert data==b''.join(bytes.fromhex(x['bytes']) for x in rows)
        blocks.append({'ida_linear':a,'end':b,'original_name_before_analysis':ida_name.get_name(a),
                       'file_sha256':hashlib.sha256(data).hexdigest(),'instructions':rows,
                       'xrefs_to':[{'from':x.frm,'type':x.type} for x in idautils.XrefsTo(a)]})
    r['decoded_blocks']=blocks
    r['scene_callers']=[]
    for x in idautils.XrefsTo(0x13B08):
        if x.type not in (16,17):continue
        f=ida_funcs.get_func(x.frm)
        assert f is not None
        rows=[probe.instruction(ea) for ea in idautils.Heads(max(f.start_ea,x.frm-40),x.frm+5)
              if ida_bytes.is_code(ida_bytes.get_full_flags(ea))]
        r['scene_callers'].append({'from':x.frm,'function':ida_funcs.get_func_name(f.start_ea),'instructions':rows})
    r['closure_scope']='Ruler sortie, automatic formation and faction sidebar fixed roots; transitive unrestored direct callees. Relocation list UI is a separate pending closure.'
    r['recovery_ledger_sha256']=hashlib.sha256(Path('/evidence/c-recovery-status.json').read_bytes()).hexdigest()
    p.write_text(json.dumps(r,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8');ida_pro.qexit(1)
else:ida_pro.qexit(0)
