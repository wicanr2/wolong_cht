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
probe.TARGETS=tuple(0x10000+x for x in (0xf020, 0xf140, 0xf17b, 0xf1a3, 0xaaa, 0xad9, 0xcac, 0xcc3, 0xc61f, 0xc6bf, 0xc6ae, 0xc6f6, 0xc74c, 0xc775, 0xc78e, 0xbcd, 0xf888, 0xf938, 0xf999, 0xc7f4, 0xc673, 0xf9b0, 0xfa1b, 0xfa37, 0xfaa2, 0xe3d7, 0x1b4, 0x1db, 0x895d, 0x89de, 0xd5d4, 0xc14, 0xc60, 0xc77, 0xe453, 0xdf))

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
    assert len(report['targets'])==36
    report['target_relocations_applied']=applied
    report['raw_tables']=[]
    for ea,size in [(0x1d2e4,6),(0x1d2ea,12)]:
        at=ea-0x10000+header;data=raw[at:at+size]
        assert ida_bytes.get_bytes(ea,size)==data
        report['raw_tables'].append({'ida_linear':ea,'file_offset':at,'size':size,
                                     'file_bytes':data.hex(),'sha256':hashlib.sha256(data).hexdigest()})
    report['display_rectangle_handlers']=[probe.instruction(ea) for ea in idautils.Heads(0x1ea26,0x1ea66) if ida_bytes.is_code(ida_bytes.get_full_flags(ea))]
    path.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
