#!/usr/bin/env python3
"""核對 native C 顯示分派、實際場景、glyph 邊界与組語補充。"""
import argparse
import hashlib
import json
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_1030F':40,'sub_10337':4,'sub_1E993':20,'sub_1E9A7':26,'sub_1E9C1':76,
     'sub_1FB29':126,'sub_1FBA7':84,'sub_1F6DC':68,'sub_1F878':16}
GROUPS={'register':12,'scene':120,'axis':144,'line-box':36,'texture':5,'width':6,'text':4,'opcode':54,'table-swap':1}
BACKLINKS={
 'docs/re/48-window-display-list.md':('後續 C 分派證據見','107-c-display-interpreter-restoration.md'),
 'docs/re/106-c-rectangle-bars-restoration.md':('後續顯示清單證據見','107-c-display-interpreter-restoration.md'),
 'docs/spec/225-c-rectangle-bars.md':('後續 C 顯示清單','227-c-display-interpreter.md'),
}

def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 p=json.loads((out/'ida/ida-probe.json').read_text());assert p['input_sha256']==p['ida_input_sha256']==EXPECTED and p['tool_version']=='9.4' and p['function_count']==739
 assert len(p['targets'])==14 and len(p['decoded_blocks'])==13
 routines={}
 for t in p['targets']:
  b=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks']);assert b==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
  routines[t['name']]=hashlib.sha256(b).hexdigest();assert routines[t['name']]==t['file_sha256']
  if t['name'] in NEW:assert len(b)==NEW[t['name']]
 new={n:routines[n] for n in NEW};assert sum(NEW.values())==460
 for b in p['decoded_blocks']:
  for r in b['instructions']:
   at=r['ida_linear']-0x10000+512;assert bytes.fromhex(r['bytes'])==raw[at:at+len(bytes.fromhex(r['bytes']))]
 data=p['inline_state'];assert [(r['ida_linear'],len(bytes.fromhex(r['bytes']))) for r in data]==[(0x1ee62,2),(0x1f45b,10)]
 full={};receipts={}
 for level in ['O0','O2']:
  f=out/'results'/f'{level}.json';r=json.loads(f.read_text());assert r['schema']=='wolong-c-display-parity-v1' and r['input_sha256']==EXPECTED
  assert r['passed'] and r['mismatch'] is None and r['cases']==382 and r['groups']==GROUPS
  assert r['full_ram_plane_audits']==382 and r['indexed_content_audits']==382 and r['go_bar_cases']==0
  assert r['trace_entry_size']==90 and r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
  assert all(r['routine_sha256'][n]==v for n,v in routines.items())
  assert all(r['entries_seen'].get(n,0)>0 for n in NEW) and r['entries_seen']['sub_1F75E']>0
  receipts[level]=sha(f);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in [(1,'scene'),(2,'opcode'),(3,'table-swap'),(4,'opcode'),(5,'text'),(6,'text'),(7,'texture'),(8,'axis'),(9,'line-box'),(10,'register')]:
  f=out/'results'/f'mutant-{n}.json';r=json.loads(f.read_text());assert not r['passed'] and r['input_sha256']==EXPECTED and 0<r['cases']<=GROUPS[g] and r['groups']=={g:r['cases']}
  m=r['mismatch'];assert m['case']==r['cases']-1 and m['group']==g
  assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_ports','_device','_planes','_ram'])
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(f)}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert sha(repo/relative)==digest;source[relative]=digest
 for name in ['tools/c_recovery/display.c','tools/c_recovery/display.h','tools/c_recovery/display_fixture.h','tools/c_recovery/display_generated.inc','tools/c_recovery_display_generate.py','tools/c_recovery_display.go']:assert name in source
 generated=(repo/'tools/c_recovery/display_generated.inc').read_text();assert 'opcode fetch' in generated and 'static void dl_body_' in generated
 compiled=sha(out/'results/c-source.sha256');assert (out/'results/compiled-source-digest.txt').read_text().strip()==compiled
 for level in ['O0','O2']:assert '-DKI_DISPLAY_SOURCE_DIGEST=0x'+compiled[:16] in (out/'results'/f'buildinfo-{level}.txt').read_text()
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 for path,(marker,link) in BACKLINKS.items():s=(repo/path).read_text();assert marker in s and link in s
 supplement=json.loads((repo/'docs/re/rectangle-handler-code.json').read_text());assert supplement['input_sha256']==EXPECTED
 original_display=[r for r in supplement['instructions'] if 0x1EA1F<=r['ida_linear']<0x1F4A2]
 assert len(original_display)==317 and sum(r['file_end']-r['file_start'] for r in original_display)==746
 result={'schema':'wolong-c-display-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,'new_routine_sha256':new,
         'ida_database_sha256':sha(out/'ida/input.exe.i64'),'source_sha256':source,'compiled_source_manifest_sha256':compiled,'cases_per_optimization':382,'groups':GROUPS,
         'full_ram_plane_audits_per_optimization':382,'indexed_content_audits_per_optimization':382,'receipt_sha256':receipts,'negative_controls_rejected':10,'mutants':controls,
         'entries_seen':full['O2']['entries_seen'],'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
         'glyph_boundary':'Fixed cmp CX,CX; RET fixture at original 1F75E; real text loop/width/shadow/call inputs; no raster parity claim',
         'scope_backlinks':BACKLINKS,'supplement_source_sha256':sha(repo/'docs/re/rectangle-handler-code.json'),'supplement_instructions':317,'supplement_code_bytes':746,'c_machine_code_match':False}
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('9 named C functions; 17 additional code entries; 382 full RAM/plane scenes/opcodes; 10 mutations; source identity: PASS')

if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
 a=p.parse_args();verify(a.repo,a.output)
