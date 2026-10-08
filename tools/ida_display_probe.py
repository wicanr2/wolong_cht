#!/usr/bin/env python3
"""矩形／計量條／選取與原始 C 繪圖閉包的 IDA 證據。"""
import sys
import traceback
import hashlib
import json
import struct
from pathlib import Path
sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_pro
import ida_nalt
import ida_bytes
import ida_ua
import idautils
probe.TARGETS=(0x1030F,0x10337,0x1E993,0x1E9A7,0x1E9C1,0x1F020,0x1F1A3,0x1F888,0x1F938,0x1F999,0x1FB29,0x1FBA7,0x1F6DC,0x1F878)

try:
    probe.main()
    path=Path('/output/ida-probe.json')
    report=json.loads(path.read_text(encoding='utf-8'))
    raw=Path(ida_nalt.get_input_file_path()).read_bytes()
    header=struct.unpack_from('<H',raw,8)[0]*16
    count=struct.unpack_from('<H',raw,6)[0]
    table=struct.unpack_from('<H',raw,24)[0]
    relocations=[header+off+seg*16 for off,seg in
                 (struct.unpack_from('<HH',raw,table+i*4) for i in range(count))]
    report.update(mz_header_size=header,relocation_file_offsets=relocations,ida_load_paragraph=0x1000)
    applied=0
    for target in report['targets']:
        file_chunks=[]
        for chunk in target['chunks']:
            start=chunk['start']-0x10000+header
            end=start+chunk['end']-chunk['start']
            file_bytes=raw[start:end];expected=bytearray(file_bytes);rows=[]
            for at in relocations:
                if not start<=at<end:continue
                assert at+2<=end
                original=struct.unpack_from('<H',raw,at)[0];shifted=(original+0x1000)&0xffff
                struct.pack_into('<H',expected,at-start,shifted)
                rows.append({'file_offset':at,'original_word':original,'ida_word':shifted})
            assert bytes.fromhex(chunk['bytes'])==expected
            chunk.update(file_bytes=file_bytes.hex(),loader_relocations=rows)
            file_chunks.append(file_bytes);applied+=len(rows)
        target['file_sha256']=hashlib.sha256(b''.join(file_chunks)).hexdigest()
    assert len(report['targets'])==14
    report['target_relocations_applied']=applied
    report['raw_tables']=[]
    for ea,size in [(0x1d2e4,6),(0x1d2ea,12)]:
        at=ea-0x10000+header;data=raw[at:at+size]
        assert ida_bytes.get_bytes(ea,size)==data
        report['raw_tables'].append({'ida_linear':ea,'file_offset':at,'size':size,
                                     'file_bytes':data.hex(),'sha256':hashlib.sha256(data).hexdigest()})
    import idc
    report['decoded_blocks']=[]
    ranges=[(0x1ea1f,0x1ea26),(0x1ea26,0x1ea30),(0x1ea30,0x1ea3a),(0x1ea3a,0x1ea44),
            (0x1ea44,0x1ea5a),(0x1ea5a,0x1ea70),(0x1ea70,0x1ea96),(0x1ea96,0x1eab3),(0x1eab3,0x1ead0),
            (0x1ead0,0x1eae9),(0x1edfe,0x1f020),(0x1f26e,0x1f45b),(0x1f465,0x1f4a2)]
    for start,end in ranges:
        ea=start;rows=[];original_name=idc.get_name(start)
        while ea<end:
            if ea==0x1ee62:
                rows.append({'ida_linear':ea,'bytes':ida_bytes.get_bytes(ea,2).hex(),'inline_data':True,'evidence':'RET at 1EE61; word stores/reads target 1EE62'});ea+=2;continue
            insn=ida_ua.insn_t();size=ida_ua.decode_insn(insn,ea)
            if size<=0 or ea+size>end:raise ValueError(f'cannot decode boundary {ea:X}')
            original_line=idc.generate_disasm_line(ea,0);classified=ida_bytes.is_code(ida_bytes.get_full_flags(ea));name=idc.get_name(ea)
            raw_before=ida_bytes.get_bytes(ea,size)
            if not classified:
                ida_bytes.del_items(ea,ida_bytes.DELIT_SIMPLE,size)
                assert ida_ua.create_insn(ea)==size
            row=probe.instruction(ea);assert bytes.fromhex(row['bytes'])==raw_before
            row.update(original_ida_data_line=original_line,original_name_before_analysis=name,ida_code_classified_before=classified)
            rows.append(row);ea+=size
        report['decoded_blocks'].append({'ida_linear':start,'end':end,'original_name_before_analysis':original_name,'instructions':rows,'analysis_scope':'temporary database; pristine executable and existing project databases unchanged'})
    report['inline_state']=[{'ida_linear':0x1ee62,'bytes':ida_bytes.get_bytes(0x1ee62,2).hex()},{'ida_linear':0x1f45b,'bytes':ida_bytes.get_bytes(0x1f45b,10).hex()}]
    base=0x1ea0d
    report['dispatch_table']={'ida_linear':base,'offsets':[int.from_bytes(ida_bytes.get_bytes(base+i*2,2),'little') for i in range(9)]}
    path.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
