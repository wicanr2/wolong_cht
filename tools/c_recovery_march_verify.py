#!/usr/bin/env python3
"""行軍 C 原始指令、完整裝置比較、錯版與來源綁定驗證。"""
import argparse,hashlib,importlib.util,json,tempfile,shutil
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_11C8D': 36, 'sub_11F7F': 249, 'sub_12151': 97, 'sub_121B2': 53, 'sub_14325': 51, 'sub_14370': 45, 'sub_1439D': 18, 'sub_143AF': 96, 'sub_1440F': 87, 'sub_14466': 29, 'sub_14483': 22, 'sub_14499': 16, 'sub_144A9': 45, 'sub_144D6': 44, 'sub_14548': 45, 'sub_1463E': 19, 'sub_14651': 56, 'sub_14689': 15, 'sub_159A6': 17, 'sub_15AB6': 27, 'sub_1703C': 98, 'sub_1709E': 77, 'sub_17F90': 75, 'sub_17FDB': 115, 'sub_1804E': 45}
GROUPS={'dispatch': 384, 'handler': 576, 'dissolve': 24, 'route': 32, 'hit': 44, 'command': 24, 'controller': 28, 'picker': 24, 'cursor': 96, 'camera': 32, 'minimap': 24, 'menu': 28, 'redraw': 4}
MUTANTS={1: 'picker', 2: 'hit', 3: 'command', 4: 'command', 5: 'menu', 6: 'dispatch', 7: 'handler', 8: 'handler', 9: 'handler', 10: 'handler', 11: 'dissolve', 12: 'controller', 13: 'hit', 14: 'cursor', 15: 'handler', 16: 'route'}
BACKLINKS=[{'ida_linear': 98228, 'function_ida_linear': 98192, 'older_document': 'docs/spec/149-march-target-map-picker.md', 'required_markers': ['取消仍分派', '123-c-march-command-restoration.md']}, {'ida_linear': 83567, 'function_ida_linear': 83537, 'older_document': 'docs/spec/39-march-order-menu.md', 'required_markers': ['只清軍團+00', '123-c-march-command-restoration.md']}, {'ida_linear': 94418, 'function_ida_linear': 94366, 'older_document': 'docs/re/85-march-target-hit-test.md', 'required_markers': ['資料完整性前提', '123-c-march-command-restoration.md']}, {'ida_linear': 73670, 'function_ida_linear': 73599, 'older_document': 'docs/re/47-main-screen-window-registry.md', 'required_markers': ['picker不套用', '123-c-march-command-restoration.md']}]
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 asset_spec=importlib.util.spec_from_file_location('assets',repo/'tools/c_recovery_talk_verify.py');assets=importlib.util.module_from_spec(asset_spec);asset_spec.loader.exec_module(assets)
 inputs={**assets.ASSETS,'KYOGRF.DAT':(69120,'e086f526bdada5baf751c41d2f73a78a0ba70002f282f63a9f33114542ed933f'),'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
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
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks']);assert count==613
 receipts={};full={}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-march-parity-v1' and r['passed'] and r['mismatch'] is None
  assert r['input_sha256']==EXPECTED and r['groups']==GROUPS and r['cases']==sum(GROUPS.values())==1320
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==1320
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in routines.items())
  assert r['independent_march_audits']>1320
  assert r['entries_seen'].get('sub_11D8E',0)==0 and r['entries_seen'].get('sub_1ECE0',0)>0
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
  info=(out/'results'/f'buildinfo-mutant-{n}.txt').read_text();assert f'-DKI_MARCH_MUTATION={n}' in info
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path),'buildinfo_sha256':sha(out/'results'/f'buildinfo-mutant-{n}.txt')}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest,path;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']+[f'mutant-{n}' for n in MUTANTS]:
  info=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_MARCH_SOURCE_DIGEST=0x'+compiled[:16] in info and 'matching_march' in info
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='march-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_march_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,probe);assert (root/'tools/c_recovery/march_generated.inc').read_bytes()==(repo/'tools/c_recovery/march_generated.inc').read_bytes()
 result={'schema':'wolong-c-march-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(probe),'new_routine_sha256':routines,
 'original_instruction_count':613,'original_instruction_bytes':1477,'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_march.sh','tools/c_recovery_march_verify.py','tools/ida_march_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py','tools/c_recovery_talk_verify.py','tools/march_code_supplement.py']},
 'cases_per_optimization':1320,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':16,'mutants':controls,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'independent_march_audits':full['O2']['independent_march_audits'],
 'input_assets':inputs,'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,
 'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'25 original corps command/map picker/stage dispatcher routines; four scenarios, 12 stages/9 handlers, documented DI residual, routes, refill and dissolve, real mouse/far services and mode-popup retry, outer cancellation still dispatches and sets bit1. Original-first full state/command/stage snapshot models, complete RAM/planes/registers/FLAGS/stack/IN/OUT/API. Every case initializes original raw RNG C=seed^A5, S=seed and table[i]=byte(i*73+i/2+seed), with fixed enumerated seeds and no reroll; original and C receive identical raw state. IF/TF=0 and pinned dosgolem flags; natural long player paths, invalid stage/damaged city map and original C machine code remain pending.',
 'c_machine_code_match':False}
 code=repo/'docs/re/c-march-code.json';sp=json.loads(code.read_text());ap=sp['assembler_verification']
 assert sp['input_sha256']==EXPECTED and sp['probe_sha256']==sha(probe) and sp['database_sha256']==sha(out/'ida/input.exe.i64')
 assert not sp['instructions'] and sp['constant_symbols']=={}
 assert ap['status']=='byte-exact' and ap['instructions']==ap['bytes']==0 and ap['unique_input_instructions']==ap['covered_input_instructions']==613
 assert ap['independent_reassembly']['status']=='byte-exact' and ap['independent_reassembly']['instructions']==613 and ap['independent_reassembly']['bytes']==1477
 assert len(sp['instruction_coverage'])==613 and all(x['fixed_prior_coverage'] and x['independent_native_assembly_match'] for x in sp['instruction_coverage'])
 assert ap['published_manifest_sha256']==sha(repo/ap['published_manifest']) and ap['published_code_source_sha256']==sha(repo/ap['published_code_source'])
 assert ap['generator_sha256']==sha(repo/'tools/march_code_supplement.py')
 relocations=sp['loader_relocation_normalization'];assert len(relocations)==5
 for row in relocations:
  assert row['locator_level']=='proven' and row['file_bytes']=='9a00000010' and row['ida_loaded_bytes']=='9a00000020'
  assert raw[row['file_start']:row['file_end']].hex()==row['file_bytes']
  assert len(row['loader_relocations'])==1 and row['loader_relocations'][0]['original_word']==0x1000 and row['loader_relocations'][0]['ida_word']==0x2000
 result['mz_far_call_normalization']=relocations
 result['rng_calls_per_optimization']=full['O2']['entries_seen']['sub_1ECE0']
 result['game_clock_calls_per_optimization']=0
 result['rng_setup']='Before each run: raw CS:ECFC C=seed^A5, S=seed, table[i]=byte(i*73+i/2+seed); fixed matrix seeds 0,17,34,51,68,85,102,119; identical original/native initial state; no reroll.'
 result['assembly_coverage_receipt_sha256']=sha(code)
 go=repo/'docs/re/c-march-go-verification.json'
 if go.exists():
  gp=json.loads(go.read_text());assert gp['vet_passed'] and gp['test_passed'] and gp['cold_test_cache']
  assert gp['original_assets_readonly'] and gp['network']=='none' and gp['tested_packages']==39 and gp['cached_test_packages']==0
  for name,digest in gp['source_sha256'].items():assert sha(repo/name)==digest,name
  result['go_verification_sha256']=sha(go)
 evidence='docs/re/123-c-march-command-restoration.md';evidence_text=(repo/evidence).read_text();assert EXPECTED in evidence_text
 result['resolution_backlinks']=[]
 result['backlink_status']='pending-publication'
 if '**狀態：CONFORMED' in evidence_text:
  for row in BACKLINKS:
   assert f"0x{row['ida_linear']:X}" in evidence_text,row
   older=(repo/row['older_document']).read_text();assert all(marker in older for marker in row['required_markers']),row
   result['resolution_backlinks'].append({**row,'platform_module':'dosv/KI.EXE','input_sha256':EXPECTED,'level':'proven','evidence':evidence})
  result['backlink_status']='verified'
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('25 functions / 613 instructions / 1320 whole-device cases / 16 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
