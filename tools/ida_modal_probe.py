#!/usr/bin/env python3
"""外交／金額視窗十七個原始函式與遠呼叫、旗標、stack 證據。"""
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
probe.TARGETS=tuple(0x10000+x for x in (
    0x2078,0x20d6,0x38c7,0x38e6,0x3902,0x39e8,0x3c3d,0x3b7e,
    0x3d09,0x3d45,0x3c99,0x3cdc,0x9321,0x87ff,0x3d68,0x1d46,0x2216))
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
            file_bytes=raw[start:end];expected=bytearray(file_bytes)
            rows=[]
            for at in relocations:
                if not start<=at<end:continue
                assert at+2<=end
                original=struct.unpack_from('<H',raw,at)[0]
                shifted=(original+0x1000)&0xffff
                struct.pack_into('<H',expected,at-start,shifted)
                rows.append({'file_offset':at,'original_word':original,'ida_word':shifted})
            assert bytes.fromhex(chunk['bytes'])==expected
            chunk.update(file_bytes=file_bytes.hex(),loader_relocations=rows)
            file_chunks.append(file_bytes);applied+=len(rows)
        target['file_sha256']=hashlib.sha256(b''.join(file_chunks)).hexdigest()
    assert len(report['targets'])==17 and applied==20
    report['target_relocations_applied']=applied
    path.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
