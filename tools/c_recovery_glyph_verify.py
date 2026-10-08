#!/usr/bin/env python3
"""核對真实字形、独立字体服务、far／self operand 与全 raster 收據。"""
import argparse
import hashlib
import json
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
FONTS={'END_S13.DAT':(404992,'0ab4c43920787bf87c91bba5730078e451cb6125b878f201d3e70d648412569b'),
       'END_S14.DAT':(3840,'1d0cf09d0a319a9e7039190688c6a905ba4370bd369fbfcfa43b2078480d6918')}
NEW={'sub_1F720':62,'sub_1F7A4':212,'sub_106F5':4,'sub_106FD':4}
GROUPS={'vectors':2,'glyph':140,'alignment':64,'text-name':16,'scene':30}
BACKLINKS={
 'docs/re/29-font-service-int15.md':('後續原始 C glyph','108-c-glyph-raster-restoration.md'),
 'docs/re/107-c-display-interpreter-restoration.md':('後續真實 glyph','108-c-glyph-raster-restoration.md'),
 'docs/spec/227-c-display-interpreter.md':('真正 glyph raster','228-c-glyph-raster.md'),
}

def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 for name,(size,digest) in FONTS.items():p=repo/'workplace/orig/dosv'/name;assert p.stat().st_size==size and sha(p)==digest
 probe=json.loads((out/'ida/ida-probe.json').read_text());assert probe['tool_version']=='9.4' and probe['input_sha256']==probe['ida_input_sha256']==EXPECTED and probe['function_count']==739
 assert len(probe['targets'])==6 and len(probe['decoded_blocks'])==2
 routines={}
 for t in probe['targets']:
  b=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks']);assert b==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
  routines[t['name']]=hashlib.sha256(b).hexdigest();assert routines[t['name']]==t['file_sha256']
  if t['name'] in NEW:assert len(b)==NEW[t['name']]
 assert sum(NEW.values())==282
 for b in probe['decoded_blocks']:
  for x in b['instructions']:
   at=x['ida_linear']-0x10000+512;v=bytes.fromhex(x['bytes']);assert v==raw[at:at+len(v)]
 receipts={};full={}
 for level in ['O0','O2']:
  f=out/'results'/f'{level}.json';r=json.loads(f.read_text());assert r['schema']=='wolong-c-glyph-parity-v1' and r['input_sha256']==EXPECTED
  assert r['passed'] and r['mismatch'] is None and r['cases']==252 and r['groups']==GROUPS
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==252 and r['go_bar_cases']==0
  assert r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
  assert r['original_font_calls']==r['c_font_calls']==[740,150] and r['original_missing_fonts']==r['c_missing_fonts']==0
  assert r['entries_seen']['dosv-font-full@0080:0410']==740 and r['entries_seen']['dosv-font-half@0080:0414']==150
  assert all(r['entries_seen'].get(n,0)>0 and r['routine_sha256'][n]==v for n,v in routines.items())
  assert r['entries_seen']['sub_1F75E']==996 and r['entries_seen']['sub_1F7A4']==890
  receipts[level]=sha(f);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in [(1,'glyph'),(2,'glyph'),(3,'glyph'),(4,'glyph'),(5,'glyph'),(6,'glyph'),(7,'glyph'),(8,'text-name'),(9,'text-name')]:
  f=out/'results'/f'mutant-{n}.json';r=json.loads(f.read_text());assert not r['passed'] and r['input_sha256']==EXPECTED and r['groups']=={g:r['cases']} and 0<r['cases']<=GROUPS[g]
  m=r['mismatch'];assert m['case']==r['cases']-1 and m['group']==g
  assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_ports','_device','_planes','_ram'])
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(f)}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert sha(repo/relative)==digest;source[relative]=digest
 for name in ['tools/c_recovery/glyph.c','tools/c_recovery/glyph.h','tools/c_recovery/glyph_fixture.h','tools/c_recovery/glyph_generated.inc','tools/c_recovery_glyph_generate.py','tools/c_recovery_glyph.go','tools/c_recovery_glyph_platform.go','tools/c_recovery_display_generate.py']:assert name in source
 compiled=sha(out/'results/c-source.sha256');assert (out/'results/compiled-source-digest.txt').read_text().strip()==compiled
 for level in ['O0','O2']:assert '-DKI_GLYPH_SOURCE_DIGEST=0x'+compiled[:16] in (out/'results'/f'buildinfo-{level}.txt').read_text()
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 platform=(repo/'tools/c_recovery_glyph_platform.go').read_text();assert '.IntHook(c,' in platform and '.Step(' not in platform
 for path,(marker,link) in BACKLINKS.items():s=(repo/path).read_text();assert marker in s and link in s
 supplement=json.loads((repo/'docs/re/rectangle-handler-code.json').read_text());original_subset=[x for x in supplement['instructions'] if not (0x12239<=x['ida_linear']<0x12280 or 0x19440<=x['ida_linear']<0x19444 or 0x1945a<=x['ida_linear']<0x19469)];assert len(original_subset)==338 and sum(x['file_end']-x['file_start'] for x in original_subset)==805
 result={'schema':'wolong-c-glyph-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,'fonts':FONTS,'new_routine_sha256':{n:routines[n] for n in NEW},
         'ida_database_sha256':sha(out/'ida/input.exe.i64'),'source_sha256':source,'compiled_source_manifest_sha256':compiled,'cases_per_optimization':252,'groups':GROUPS,
         'full_ram_plane_audits_per_optimization':252,'indexed_content_audits_per_optimization':252,'font_calls_per_optimization':[740,150],'font_missing':0,
         'receipt_sha256':receipts,'negative_controls_rejected':9,'mutants':controls,'entries_seen':full['O2']['entries_seen'],'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
         'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),'scope_backlinks':BACKLINKS,'glyph_raster_fixture':False,
         'platform_contract':'Independent pinned DOS/V Font services/caches; C-side calls IntHook API without CPU.Step; actual local font bytes; no original TSR/I/O-time claim',
         'supplement_source_sha256':sha(repo/'docs/re/rectangle-handler-code.json'),'c_machine_code_match':False}
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('4 named C functions; 2 original code entries; 252 real glyph/scene cases; 740/150 font reads; 9 mutations; source identity: PASS')

if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
