#!/usr/bin/env python3
"""核對原基準覆蓋後，以固定 binutils 匹配清單閉包的新增指令。"""
import argparse,hashlib,importlib.util,json
from pathlib import Path

def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def generate(repo,out):
 spec=importlib.util.spec_from_file_location('assembly',repo/'tools/assembly_rebuild.py');asm=importlib.util.module_from_spec(spec);spec.loader.exec_module(asm)
 probe_path=out/'ida/ida-probe.json';p=json.loads(probe_path.read_text())
 original=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(original).hexdigest()==asm.EXPECTED==p['input_sha256']
 baseline=json.loads((repo/'workplace/matching-decompilation/assembly/build/source-map.json').read_text())
 historical_path=repo/'docs/re/rectangle-handler-code.json';historical=json.loads(historical_path.read_text());assert len(historical['instructions'])==370
 coverage=[(x['file_start'],x['file_end']) for x in baseline if x['kind']=='instruction']+[(x['file_start'],x['file_end']) for x in historical['instructions']]
 all_rows=[x for t in p['targets'] for c in t['chunks'] for x in c['instructions']]+[x for b in p['decoded_blocks'] for x in b['instructions']]
 rows=[]
 for x in all_rows:
  at=x['ida_linear']-0x10000+512;z=at+len(bytes.fromhex(x['bytes']))
  if any(a<=at and z<=b for a,b in coverage):continue
  assert original[at:z].hex()==x['bytes']
  unknown=[b for b in baseline if b['kind']=='data-or-unknown' and b['file_start']<=at<z<=b['file_end']];assert len(unknown)==1
  rows.append({**x,'file_offset':at,'relocation_file_offsets':[],'baseline_unknown_span':[unknown[0]['file_start'],unknown[0]['file_end']]})
 rows.sort(key=lambda x:x['file_offset']);assert len(rows)==83 and sum(len(bytes.fromhex(x['bytes'])) for x in rows)==215
 trials_path=out/'assembly-candidates';trials_path.mkdir(exist_ok=True)
 matched,unmatched,trials,symbols=asm.compile_candidates(rows,trials_path);assert not unmatched and not symbols,(unmatched,symbols)
 source=['.intel_syntax noprefix','.code16','.section .image,"ax",@progbits'];entries=[]
 for x in rows:
  at=x['file_offset'];end=at+len(bytes.fromhex(x['bytes']));gas=matched[at]
  entries.append({'ida_linear':x['ida_linear'],'file_start':at,'file_end':end,'kind':'instruction',
   'original_name':'','function_name':None,'original_assembly':x['assembly'],'gas':gas,'bytes':x['bytes'],
   'original_ida_data_line':x.get('original_ida_data_line'),'ida_code_classified':x.get('originally_code',True),
   'baseline_unknown_span':x['baseline_unknown_span'],'operands':x['operands'],'locator_level':'proven',
   'scope_sources':['docs/re/118-c-city-list-restoration.md']})
  source.extend([f'.org {at:#x}',gas])
 source_path=out/'list-patch.S';source_path.write_text('\n'.join(source)+'\n')
 asm.run(['as','--32','-o',str(out/'list-patch.o'),str(source_path)])
 asm.run(['ld','-m','elf_i386','--section-start','.image=0','--entry','0','-o',str(out/'list-patch.elf'),str(out/'list-patch.o')])
 asm.run(['objcopy','-O','binary','--only-section=.image',str(out/'list-patch.elf'),str(out/'list-patch.bin')])
 b=(out/'list-patch.bin').read_bytes()
 for x in entries:assert b[x['file_start']:x['file_end']]==original[x['file_start']:x['file_end']]
 verification={'status':'byte-exact','instructions':83,'bytes':215,'patch_points':7,'source_sha256':sha(source_path),
  'binary_sha256':hashlib.sha256(b).hexdigest(),'image_id':'sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e',
  'assembler':asm.run(['as','--version']).splitlines()[0],
  'generator_sha256':sha(Path(__file__)),'translator_sha256':sha(repo/'tools/assembly_rebuild.py'),
  'historical_supplement_sha256':sha(historical_path),'trials':trials}
 result={'schema':'wolong-matching-code-supplement-v1','input_sha256':asm.EXPECTED,'tool':'IDA Pro','tool_version':'9.4',
  'database_sha256':sha(out/'ida/input.exe.i64'),'probe_sha256':sha(probe_path),
  'evidence':'docs/re/118-c-city-list-restoration.md','instructions':entries,'assembler_verification':verification}
 (out/'list-code-supplement.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('83 previously uncovered instructions / 215 bytes: native GNU assembly exact')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();generate(a.repo,a.output)
