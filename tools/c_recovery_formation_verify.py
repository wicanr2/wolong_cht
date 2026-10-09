#!/usr/bin/env python3
"""編成 C 原始指令、完整裝置比較、錯版與來源綁定驗證。"""
import argparse,hashlib,importlib.util,json,tempfile,shutil
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_16C5E':52,'sub_16C92':196,'sub_16D56':25,'sub_16D6F':57,'sub_16DA8':85,'sub_16DFD':131,'sub_16E80':15}
GROUPS={'init':4,'numbers':28,'icons':16,'frame':4,'close':4,'formation':56,'switch':72,'offpanel':4,'caller':12}
MUTANTS={1:'init',2:'numbers',3:'icons',4:'frame',5:'formation',6:'switch',7:'formation',8:'formation',9:'caller',10:'formation'}
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 asset_spec=importlib.util.spec_from_file_location('assets',repo/'tools/c_recovery_talk_verify.py');assets=importlib.util.module_from_spec(asset_spec);asset_spec.loader.exec_module(assets)
 inputs={**assets.ASSETS,'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
 'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
 'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7')}
 for name,(size,digest) in inputs.items():
  path=repo/'workplace/orig/dosv'/name;assert path.stat().st_size==size and sha(path)==digest,name
 probe=out/'ida/ida-probe.json';p=json.loads(probe.read_text());assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 routines={}
 for t in p['targets']:
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
  assert len(data)==NEW[t['name']] and hashlib.sha256(data).hexdigest()==t['file_sha256'];routines[t['name']]=t['file_sha256']
 assert set(routines)==set(NEW) and not p['decoded_blocks']
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks']);assert count==247
 receipts={};full={}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-formation-parity-v1' and r['passed'] and r['mismatch'] is None
  assert r['input_sha256']==EXPECTED and r['groups']==GROUPS and r['cases']==sum(GROUPS.values())==200
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==200
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in routines.items())
  assert r['independent_formation_audits']>200
  assert r['original_font_calls']==r['c_font_calls'] and sum(r['original_font_calls'])>0
  assert r['original_missing_fonts']==r['c_missing_fonts']==0 and r['original_font_misses']==r['c_font_misses']=={}
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch']
  assert not r['passed'] and m['group']==g and r['groups']=={g:r['cases']}
  assert 0<r['cases']<=GROUPS[g] and m['case']==r['cases']-1
  assert any(m['original'+k]!=m['c'+k] for k in ['','_trace','_ports','_device','_planes','_ram','_api','_sound','_mouse','_in','_ticks'])
  info=(out/'results'/f'buildinfo-mutant-{n}.txt').read_text();assert f'-DKI_FORMATION_MUTATION={n}' in info
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path),'buildinfo_sha256':sha(out/'results'/f'buildinfo-mutant-{n}.txt')}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest,path;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']+[f'mutant-{n}' for n in MUTANTS]:
  info=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_FORMATION_SOURCE_DIGEST=0x'+compiled[:16] in info and 'matching_formation' in info
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='formation-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_formation_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,probe);assert (root/'tools/c_recovery/formation_generated.inc').read_bytes()==(repo/'tools/c_recovery/formation_generated.inc').read_bytes()
 result={'schema':'wolong-c-formation-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(probe),'new_routine_sha256':routines,
 'original_instruction_count':247,'original_instruction_bytes':561,'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_formation.sh','tools/c_recovery_formation_verify.py','tools/ida_formation_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py','tools/c_recovery_talk_verify.py','tools/formation_code_supplement.py']},
 'cases_per_optimization':200,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':10,'mutants':controls,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'independent_formation_audits':full['O2']['independent_formation_audits'],
 'input_assets':inputs,'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,
 'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'Seven original player formation routines and complete transitive C render/input/list/TALK closure. Four scenarios, six-slot cycling, original quotient/remainder/cap distribution, initialize-only types, cancellation stale total, empty-leader refusal, confirmation and outer select/form/reselect/cancel. Controlled input checkpoints and legal unassigned generals; separate pinned DOS/font/mouse/VGA, complete RAM/planes/registers/FLAGS/stack/IN/OUT/API. IF/TF=0, pinned dosgolem flags; original TSR/wall-clock, natural long player paths, invalid hotspots and original C machine code remain pending.',
 'c_machine_code_match':False}
 code=repo/'docs/re/c-formation-code.json';sp=json.loads(code.read_text());ap=sp['assembler_verification']
 assert sp['input_sha256']==EXPECTED and sp['probe_sha256']==sha(probe) and sp['database_sha256']==sha(out/'ida/input.exe.i64')
 assert not sp['instructions'] and sp['constant_symbols']=={}
 assert ap['status']=='byte-exact' and ap['instructions']==ap['bytes']==0 and ap['unique_input_instructions']==ap['covered_input_instructions']==247
 assert ap['independent_reassembly']['status']=='byte-exact' and ap['independent_reassembly']['instructions']==247 and ap['independent_reassembly']['bytes']==561
 assert len(sp['instruction_coverage'])==247 and all(x['fixed_prior_coverage'] and x['independent_native_assembly_match'] for x in sp['instruction_coverage'])
 assert ap['published_manifest_sha256']==sha(repo/ap['published_manifest']) and ap['published_code_source_sha256']==sha(repo/ap['published_code_source'])
 assert ap['generator_sha256']==sha(repo/'tools/formation_code_supplement.py')
 result['assembly_coverage_receipt_sha256']=sha(code)
 go=repo/'docs/re/c-formation-go-verification.json'
 if go.exists():
  gp=json.loads(go.read_text());assert gp['vet_passed'] and gp['test_passed'] and gp['cold_test_cache']
  assert gp['original_assets_readonly'] and gp['network']=='none' and gp['tested_packages']==39 and gp['cached_test_packages']==0
  for name,digest in gp['source_sha256'].items():assert sha(repo/name)==digest,name
  result['go_verification_sha256']=sha(go)
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('7 functions / 247 instructions / 200 whole-device cases / 10 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
