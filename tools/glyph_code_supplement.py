#!/usr/bin/env python3
"""將固定 IDA 解出的未分類 code 加入已驗證組語補充，保留資料槽。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

def main(repo,output,image_id):
    spec=importlib.util.spec_from_file_location('asm',repo/'tools/assembly_rebuild.py');asm=importlib.util.module_from_spec(spec);spec.loader.exec_module(asm)
    original=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(original).hexdigest()==asm.EXPECTED
    probe_path=repo/'workplace/matching-decompilation/c-glyph/ida/ida-probe.json';probe=json.loads(probe_path.read_text())
    mapping=json.loads((repo/'workplace/matching-decompilation/assembly/build/source-map.json').read_text());rows=[];locators={}
    for block in probe['decoded_blocks']:
        for row in block['instructions']:
            if row.get('inline_data'):continue
            start=row['ida_linear']-0x10000+512;end=start+len(bytes.fromhex(row['bytes']))
            assert original[start:end]==bytes.fromhex(row['bytes'])
            if not any(p['kind']=='data-or-unknown' and p['file_start']<=start<end<=p['file_end'] for p in mapping):continue
            converted={**row,'file_offset':start,'relocation_file_offsets':[]}
            rows.append(converted);locators[start]=(block,row)
    assert rows
    output.mkdir(parents=True,exist_ok=True)
    matched,unmatched,trials,symbols=asm.compile_candidates(rows,output);assert not unmatched
    assert not symbols,'supplement must not introduce unrecorded linker constants'
    old=json.loads((repo/"docs/re/rectangle-handler-code.json").read_text())
    public=old["instructions"].copy()
    existing={r["file_start"] for r in public}
    for row in rows:
        at=row['file_offset'];block,rawrow=locators[at];opcode=None
        if 0x1ea1f<=block['ida_linear']<0x1ead0:
            opcode=probe['dispatch_table']['offsets'].index(block['ida_linear']-0x10000)+1
        entry={'ida_linear':row['ida_linear'],'file_start':at,'file_end':at+len(bytes.fromhex(row['bytes'])),'kind':'instruction',
               'original_name':rawrow['original_name_before_analysis'] or None,'function_name':None,'original_assembly':row['assembly'],
               'gas':matched[at],'bytes':row['bytes'],'original_ida_data_line':rawrow['original_ida_data_line'],
               'ida_code_classified':rawrow['ida_code_classified_before'],'operands':row['operands'],'locator_level':'proven',
               'scope_sources':['docs/re/107-c-display-interpreter-restoration.md']}
        if opcode is not None:entry['opcode']=opcode
        entry["scope_sources"]=["docs/re/108-c-glyph-raster-restoration.md"]
        if at not in existing:public.append(entry)
    record={'schema':'wolong-matching-code-supplement-v1','input_sha256':probe['input_sha256'],'tool':'IDA Pro','tool_version':'9.4',
            'database_sha256':hashlib.sha256((probe_path.parent/'input.exe.i64').read_bytes()).hexdigest(),'image_id':probe['image_id'],
            'evidence':'docs/re/108-c-glyph-raster-restoration.md','instructions':sorted(public,key=lambda r:r['file_start']),'preserved_inline_data':old.get('preserved_inline_data',[])+probe['inline_state'],
            'assembler_image_id':image_id,'assembler_trials':trials}
    (repo/'docs/re/rectangle-handler-code.json').write_text(json.dumps(record,ensure_ascii=False,indent=2)+'\n')
    print('Supplement instructions',len(public),'bytes',sum(r['file_end']-r['file_start'] for r in public),'all assembled; inline data excluded')

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);p.add_argument('--image-id',required=True)
    a=p.parse_args();main(a.repo,a.output,a.image_id)
