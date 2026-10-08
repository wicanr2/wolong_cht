#!/usr/bin/env python3
"""IDA 9.4 全 segment 程式碼 inventory；沿用試點的原始定位與分級語意匯出。"""
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_bytes
import ida_funcs
import ida_loader
import ida_nalt
import ida_pro
import ida_segment
import idautils
import idc


def main():
    probe.main()
    report = json.loads(Path('/output/ida-probe.json').read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    header_size = struct.unpack_from('<H', raw, 8)[0]*16
    reloc_count, reloc_start = struct.unpack_from('<H',raw,6)[0], struct.unpack_from('<H',raw,24)[0]
    relocation_files = []
    for i in range(reloc_count):
        offset, segment = struct.unpack_from('<HH',raw,reloc_start+i*4)
        relocation_files.append(header_size+segment*16+offset)
    code = []
    for seg_ea in idautils.Segments():
        seg = ida_segment.getseg(seg_ea)
        for ea in idautils.Heads(seg.start_ea, seg.end_ea):
            if not ida_bytes.is_code(ida_bytes.get_full_flags(ea)):
                continue
            row = probe.instruction(ea)
            at = row['file_offset']
            size = len(bytes.fromhex(row['bytes']))
            if at < header_size or at+size > len(raw):
                raise ValueError(f'code outside original file at {ea:#x}')
            row['ida_database_bytes'] = row['bytes']
            row['bytes'] = raw[at:at+size].hex()
            row['segment_base'] = ida_segment.get_segm_base(seg)
            row['code_references'] = list(idautils.CodeRefsFrom(ea, False))
            row['relocation_file_offsets'] = [x for x in relocation_files if at<=x<at+size]
            row['original_name'] = idc.get_name(ea)
            row['function_name'] = ida_funcs.get_func_name(ea)
            code.append(row)
    code.sort(key=lambda r:r['file_offset'])
    for previous, current in zip(code,code[1:]):
        if previous['file_offset']+len(bytes.fromhex(previous['bytes']))>current['file_offset']:
            raise ValueError('overlapping IDA code items')
    data_spans=[]
    start=header_size; last=None
    for at in range(header_size,len(raw)):
        ea=0x10000+at-header_size
        # FF_TAIL 是 item 的尾 byte；分類要回到 head，不能把指令尾算成未知。
        flags=ida_bytes.get_full_flags(ida_bytes.get_item_head(ea))
        kind='code' if ida_bytes.is_code(flags) else 'ida-data' if ida_bytes.is_data(flags) else 'unknown'
        if kind != last:
            if last is not None:data_spans.append({'file_start':start,'file_end':at,'kind':last})
            start,last=at,kind
    data_spans.append({'file_start':start,'file_end':len(raw),'kind':last})
    if sum(s['file_end']-s['file_start'] for s in data_spans if s['kind']=='code') != sum(len(bytes.fromhex(r['bytes'])) for r in code):
        raise ValueError('code byte classification differs from decoded instruction coverage')
    report.update(schema='wolong-assembly-inventory-v1',input_size=len(raw),
                  mz_header_size=header_size,relocation_file_offsets=relocation_files,
                  code=code,classification_spans=data_spans)
    Path('/output/ida-probe.json').write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')


try:
    main()
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
