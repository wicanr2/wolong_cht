#!/usr/bin/env python3
"""資訊介面 C 原始指令、完整裝置比較、錯版與來源綁定驗證。"""
import argparse,hashlib,importlib.util,json,tempfile,shutil
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={"sub_15E1E":15,"sub_15E2D":31,"sub_17E1F":43,"sub_17E4A":208,"sub_17F1A":71,"sub_17F61":32,"sub_17F81":15,"sub_1807B":175,"sub_1812A":83,"sub_1817D":38}
GROUPS={"city-values":400,"city-image":64,"city-window":48,"city-frame":4,"city-close":4,"army-values":80,"army-slots":48,"army-raw":4,"army-close":16,"own-panel":16,"own-frame":16}
MUTANTS={1:"city-values",2:"city-values",3:"city-values",4:"city-values",5:"city-image",6:"army-values",7:"army-values",8:"army-slots",9:"army-slots",10:"army-close",11:"own-frame",12:"city-window"}
BACKLINKS=[
 {"ida_linear":0x17E2C,"function_ida_linear":0x17E1F,"older_document":"docs/re/50-city-info-window.md","required_markers":["索引×32","122-c-detail-panels-restoration.md"]},
 {"ida_linear":0x17EF3,"function_ida_linear":0x17E4A,"older_document":"docs/re/32-strategy-detail-panels.md","required_markers":["左移","122-c-detail-panels-restoration.md"]},
 {"ida_linear":0x18116,"function_ida_linear":0x1807B,"older_document":"docs/re/32-strategy-detail-panels.md","required_markers":["0F04","0F03","122-c-detail-panels-restoration.md"]},
 {"ida_linear":0x1818B,"function_ida_linear":0x1817D,"older_document":"docs/re/32-strategy-detail-panels.md","required_markers":["恢復自勢力情報","122-c-detail-panels-restoration.md"]},
 {"ida_linear":0x18150,"function_ida_linear":0x1812A,"older_document":"docs/spec/24-corps-info-window.md","required_markers":["D200","122-c-detail-panels-restoration.md"]},
 {"ida_linear":0x17E2C,"function_ida_linear":0x17E1F,"older_document":"docs/spec/23-city-info-window.md","required_markers":["入口不只地圖","122-c-detail-panels-restoration.md"]}]
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
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks']);assert count==310
 receipts={};full={}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-details-parity-v1' and r['passed'] and r['mismatch'] is None
  assert r['input_sha256']==EXPECTED and r['groups']==GROUPS and r['cases']==sum(GROUPS.values())==700
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==700
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in routines.items())
  assert r['independent_details_audits']>700
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
  info=(out/'results'/f'buildinfo-mutant-{n}.txt').read_text();assert f'-DKI_DETAILS_MUTATION={n}' in info
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path),'buildinfo_sha256':sha(out/'results'/f'buildinfo-mutant-{n}.txt')}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest,path;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']+[f'mutant-{n}' for n in MUTANTS]:
  info=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_DETAILS_SOURCE_DIGEST=0x'+compiled[:16] in info and 'matching_details' in info
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='details-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_details_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,probe);assert (root/'tools/c_recovery/details_generated.inc').read_bytes()==(repo/'tools/c_recovery/details_generated.inc').read_bytes()
 result={'schema':'wolong-c-details-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(probe),'new_routine_sha256':routines,
 'original_instruction_count':310,'original_instruction_bytes':711,'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_details.sh','tools/c_recovery_details_verify.py','tools/ida_details_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py','tools/c_recovery_talk_verify.py','tools/details_code_supplement.py']},
 'cases_per_optimization':700,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':12,'mutants':controls,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'independent_details_audits':full['O2']['independent_details_audits'],
 'input_assets':inputs,'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,
 'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'Ten original city-information controller, city/corps renderers and own-faction restoration functions. Four scenarios, controlled fields, neutral/capital comparison, signed growth, all16 picture nibbles with seek truncation, raw type0 separated from legal1..4, byte morale and total high word. Original-first independent argument/buffer/world and ABI audits; complete RAM/planes/registers/FLAGS/stack/IN/OUT/API, separate pinned DOS/font/mouse/VGA, IF/TF=0 and pinned dosgolem flags. Upstream military commands, natural long player paths and original C machine code remain pending.',
 'c_machine_code_match':False}
 code=repo/'docs/re/c-details-code.json';sp=json.loads(code.read_text());ap=sp['assembler_verification']
 assert sp['input_sha256']==EXPECTED and sp['probe_sha256']==sha(probe) and sp['database_sha256']==sha(out/'ida/input.exe.i64')
 assert not sp['instructions'] and sp['constant_symbols']=={}
 assert ap['status']=='byte-exact' and ap['instructions']==ap['bytes']==0 and ap['unique_input_instructions']==ap['covered_input_instructions']==310
 assert ap['independent_reassembly']['status']=='byte-exact' and ap['independent_reassembly']['instructions']==310 and ap['independent_reassembly']['bytes']==711
 assert len(sp['instruction_coverage'])==310 and all(x['fixed_prior_coverage'] and x['independent_native_assembly_match'] for x in sp['instruction_coverage'])
 assert ap['published_manifest_sha256']==sha(repo/ap['published_manifest']) and ap['published_code_source_sha256']==sha(repo/ap['published_code_source'])
 assert ap['generator_sha256']==sha(repo/'tools/details_code_supplement.py')
 result['assembly_coverage_receipt_sha256']=sha(code)
 go=repo/'docs/re/c-details-go-verification.json'
 if go.exists():
  gp=json.loads(go.read_text());assert gp['vet_passed'] and gp['test_passed'] and gp['cold_test_cache']
  assert gp['original_assets_readonly'] and gp['network']=='none' and gp['tested_packages']==39 and gp['cached_test_packages']==0
  for name,digest in gp['source_sha256'].items():assert sha(repo/name)==digest,name
  result['go_verification_sha256']=sha(go)
 evidence='docs/re/122-c-detail-panels-restoration.md';evidence_text=(repo/evidence).read_text();assert EXPECTED in evidence_text
 result['resolution_backlinks']=[]
 result['backlink_status']='pending-publication'
 if '**狀態：CONFORMED' in evidence_text:
  for row in BACKLINKS:
   assert f"0x{row['ida_linear']:X}" in evidence_text,row
   older=(repo/row['older_document']).read_text();assert all(marker in older for marker in row['required_markers']),row
   result['resolution_backlinks'].append({**row,'platform_module':'dosv/KI.EXE','input_sha256':EXPECTED,'level':'proven','evidence':evidence})
  result['backlink_status']='verified'
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('10 functions / 310 instructions / 700 whole-device cases / 12 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
