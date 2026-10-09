#!/usr/bin/env python3
"""核對君主出陣／自動編成的完整原版／C 同狀態收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import tempfile
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_1461D':33,'sub_14698':127,'sub_14717':51,'sub_15E80':55,'sub_15EB7':112,
     'sub_15F27':54,'sub_15F5D':34,'sub_15F7F':43,'sub_1699E':138,'sub_16E8F':58,
     'sub_16EC9':93,'sub_16F26':96,'sub_16F86':76,'sub_16FD2':86}
GROUPS={'formation':160,'troops':96,'sidebar':128,'sortie':64,'derived':3}
MUTANTS={1:'sortie',2:'sortie',3:'formation',4:'formation',5:'troops',6:'troops',
         7:'formation',8:'derived',9:'sortie',10:'sidebar'}
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 p=json.loads((out/'ida/ida-probe.json').read_text());assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 routines={}
 for target in p['targets']:
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in target['chunks'])
  digest=hashlib.sha256(data).hexdigest();assert digest==target['file_sha256']
  assert target['name'] in NEW and len(data)==NEW[target['name']];routines[target['name']]=digest
 assert set(routines)==set(NEW)
 assert sum(NEW.values())==1056
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks']);assert count==501
 receipts,full={},{}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-verdict-parity-v1' and r['input_sha256']==EXPECTED
  assert r['passed'] and r['mismatch'] is None and r['cases']==451 and r['groups']==GROUPS
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==451
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in routines.items())
  assert r['entries_seen']['sub_1699E']==64 and r['entries_seen']['sub_13B08']==64
  assert r['original_font_calls']==r['c_font_calls'] and sum(r['original_font_calls'])>0
  assert r['original_missing_fonts']==r['c_missing_fonts']==0
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch']
  assert not r['passed'] and r['groups']=={g:r['cases']} and m['group']==g
  assert 0<r['cases']<=GROUPS[g] and m['case']==r['cases']-1
  assert any(m['original'+k]!=m['c'+k] for k in ['','_trace','_ports','_device','_planes','_ram','_api','_sound','_mouse','_in','_ticks'])
  buildinfo=(out/'results'/f'buildinfo-mutant-{n}.txt').read_text();assert f'-DKI_VERDICT_MUTATION={n}' in buildinfo
  assert '-DKI_VERDICT_SOURCE_DIGEST=0x'+sha(out/'results/c-source.sha256')[:16] in buildinfo
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path),'buildinfo_sha256':sha(out/'results'/f'buildinfo-mutant-{n}.txt')}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']:
  text=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_VERDICT_SOURCE_DIGEST=0x'+compiled[:16] in text and 'matching_verdict' in text
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='verdict-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_verdict_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,out/'ida/ida-probe.json');assert (root/'tools/c_recovery/verdict_generated.inc').read_bytes()==(repo/'tools/c_recovery/verdict_generated.inc').read_bytes()
 result={'schema':'wolong-c-verdict-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(out/'ida/ida-probe.json'),
 'new_routine_sha256':routines,'original_instruction_count':count,'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_verdict.sh','tools/c_recovery_verdict_verify.py','tools/ida_verdict_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py']},
 'cases_per_optimization':451,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':10,'mutants':controls,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'input_assets':full['O2']['additional_assets'],'tool_versions':(out/'results/tool-versions.txt').read_text(),
 'exact_clean_regeneration':True,'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
 'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'Original ruler sortie entry with real three-line TALK/scene/cursor/world restoration, funds sign and 600-point/carry gate, six-slot selection/partial failure, reserve redistribution/remainder/cap, occupancy increment, original BX comparison and four-bit faction sidebar. Complete RAM/plane/register/FLAGS/stack/IN/OUT/API comparison; independent pinned DOS/font/mouse/sound/VGA services; IF/TF=0; isolated stack and occupancy bank; FLAGS from pinned dosgolem. Relocation list UI, natural player menus/long paths and original C machine code remain pending',
 'c_machine_code_match':False}
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('14 new / 501 instructions / 451 whole-device cases / 10 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
