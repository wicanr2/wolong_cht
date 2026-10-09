#!/usr/bin/env python3
"""核對完整場景主入口、尾跳、鏡頭與小地圖的原版／C 收據。"""
import argparse
import hashlib
import importlib.util
import json
import struct
import shutil
import tempfile
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_11F30':42,'sub_11F5A':37,'sub_13B08':82,'sub_15C58':76,'sub_19541':112,
     'sub_195B1':24,'sub_196ED':101,'sub_19752':68,'sub_1D4A3':36}
GROUPS={'minimap':96,'minimap-clamp':3,'camera':10,'camera-helper':18,'panel':4,'resume':24,'scene':16}
MUTANTS={1:'camera',2:'camera-helper',3:'camera-helper',4:'panel',5:'minimap-clamp',6:'minimap',
         7:'minimap',8:'minimap',9:'minimap',10:'scene',11:'scene',12:'scene'}
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 spec=importlib.util.spec_from_file_location('assets',repo/'tools/c_recovery_talk_verify.py');assets=importlib.util.module_from_spec(spec);spec.loader.exec_module(assets)
 inputs={**assets.ASSETS,'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
 'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
 'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7'),
 'BGM.DAT':(20826,'7a51c8b9a349b9e088f3796b70c268181c60bcebead70942f00e1621523dedc9'),
 'IVENTGRF.DAT':(76032,'23fffa03b2bba3b4c920c5db49b32b7a3fa26729523b0092a6be5ad3a0737f2f')}
 for n,(size,digest) in inputs.items():p=repo/'workplace/orig/dosv'/n;assert p.stat().st_size==size and sha(p)==digest
 p=json.loads((out/'ida/ida-probe.json').read_text());assert p['tool_version']=='9.4' and p['function_count']==739 and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 routines={};revisited={}
 for target in p['targets']:
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in target['chunks'])
  digest=hashlib.sha256(data).hexdigest();assert digest==target['file_sha256']
  if target['name'] in NEW:assert len(data)==NEW[target['name']];routines[target['name']]=digest
  else:assert target['name']=='sub_11D46';revisited[target['name']]=digest
 assert set(routines)==set(NEW) and sum(NEW.values())==578
 count=sum(len(c['instructions']) for t in p['targets'] if t['name'] in NEW for c in t['chunks']);assert count==255
 talk=(repo/'workplace/orig/dosv/TALK.DAT').read_bytes();worlds=(repo/'workplace/orig/dosv/SINARIO.DAT').read_bytes()
 first_line_minimum=0
 for scenario in range(4):
  temperament=worlds[scenario*22208+0x80+0x4240+0x1e]
  if temperament>=3:temperament-=3
  for flag,base in [(0,0x182),(1,0x18c)]:
   for index in [base+temperament,base+3,base+4+(3 if flag else 0)+temperament]:
    at=struct.unpack_from('<H',talk,index*2)[0];line=talk[at:].split(b'\0',1)[0].decode('cp950')
    # Only literal fullwidth characters from the first line; later rows and template names add work.
    first_line_minimum+=sum(ord(c)>127 for c in line)*2
 assert first_line_minimum==432
 receipts,full={},{}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-resume-parity-v1' and r['input_sha256']==EXPECTED
  assert r['passed'] and r['mismatch'] is None and r['cases']==171 and r['groups']==GROUPS
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==171
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in {**routines,**revisited}.items())
  assert r['entries_seen']['sub_13B08']==16 and r['entries_seen']['sub_11D46']==40
  assert r['entries_seen']['sub_13C99']==32 and r['entries_seen']['sub_13CDC']==16 and r['entries_seen']['sub_1075B']==48
  assert r['original_font_calls']==r['c_font_calls'] and r['original_font_calls'][0]>=first_line_minimum
  assert r['original_missing_fonts']==r['c_missing_fonts']==0
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch']
  assert not r['passed'] and r['input_sha256']==EXPECTED and r['groups']=={g:r['cases']}
  assert 0<r['cases']<=GROUPS[g] and m['group']==g and m['case']==r['cases']-1
  assert any(m['original'+k]!=m['c'+k] for k in ['','_trace','_ports','_device','_planes','_ram','_api','_sound','_mouse','_in','_ticks'])
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path)}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']:
  text=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_RESUME_SOURCE_DIGEST=0x'+compiled[:16] in text and 'matching_resume' in text
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='resume-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_resume_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,out/'ida/ida-probe.json');assert (root/'tools/c_recovery/resume_generated.inc').read_bytes()==(repo/'tools/c_recovery/resume_generated.inc').read_bytes()
 result={'schema':'wolong-c-resume-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'input_assets':inputs,'new_routine_sha256':routines,'reverified_existing_routine_sha256':revisited,
 'original_instruction_count':255,'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_resume.sh','tools/c_recovery_resume_verify.py','tools/ida_scene_resume_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py']},
 'cases_per_optimization':171,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':12,'mutants':controls,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'source_first_line_fullwidth_minimum':first_line_minimum,
 'scene_callers':[(c['function'],c['from']) for c in p['scene_callers']],
 'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
 'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'Original scene entry with caller TALK bases 0x182/0x18C, world DS and real records, three dialogue/press-counter gates, stop/play, IVENT/MMAP, full world redraw, camera tail transfer and minimap four-plane restore/bit-aligned frame. Independent pinned DOS/font/mouse/sound/VGA services; complete RAM/plane/register/stack/IN/OUT/API. IF/TF=0, isolated stack; FLAGS from pinned dosgolem. Natural event/long player path and original C machine code remain outside evidence',
 'c_machine_code_match':False}
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('9 new / existing world resume; 255 instructions; 171 whole-device cases; 12 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
