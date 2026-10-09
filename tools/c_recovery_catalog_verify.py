#!/usr/bin/env python3
"""核對軍團／武將／勢力／開局清單與獨立資料參照。"""
import argparse,hashlib,importlib.util,json,shutil,tempfile
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_1716D':59,'sub_171D3':68,'sub_175FA':66,'sub_17663':61,'sub_1770C':226,
 'sub_178A7':62,'sub_17906':62,'sub_1799C':222,'sub_17A7A':75,'sub_17B3C':51,'sub_17BC0':174}
RAW={0x171A8:43,0x17217:54,0x1724D:48,0x1727D:250,0x1763C:39,0x176A0:60,
 0x176DC:48,0x178E5:33,0x17944:40,0x1796C:48,0x17B6F:33,0x17B90:48}
GROUPS={'builder':168,'pages':168,'caller':140,'header':156,'select':56,'army':64,'person':24,'cache':16,'entry-ds':12}
MUTANTS={1:'builder',2:'builder',3:'builder',4:'builder',5:'caller',6:'army',7:'person',8:'person',9:'builder',10:'cache',11:'cache',12:'caller',13:'pages'}
BACKLINKS=[
 {'ida_linear':0x17663,'evidence_markers':['0x1768A'],'older_document':'docs/re/26-list-window-engine.md','required_markers':['MOV CL,CS:98AA','119-c-list-families-restoration.md']},
 {'ida_linear':0x17663,'evidence_markers':['0x1768A'],'older_document':'docs/mechanics/10-strategy.md','required_markers':['選武將同樣沿用','119-c-list-families-restoration.md']},
 {'ida_linear':0x1727D,'evidence_markers':['0x1735C'],'older_document':'docs/re/27-list-row-fields.md','required_markers':['`+0x06` u8','119-c-list-families-restoration.md']},
 {'ida_linear':0x1799C,'evidence_markers':['0x17A5C'],'older_document':'docs/spec/38-list-windows.md','required_markers':['**3 位**／**3 位**','119-c-list-families-restoration.md']},
 {'ida_linear':0x17BC0,'evidence_markers':['sub_17BC0'],'older_document':'docs/re/27-list-row-fields.md','required_markers':['~~開局選勢力的逐列','119-c-list-families-restoration.md']}]
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 asset_spec=importlib.util.spec_from_file_location('base_assets',repo/'tools/c_recovery_talk_verify.py');base=importlib.util.module_from_spec(asset_spec);asset_spec.loader.exec_module(base)
 inputs={**base.ASSETS,'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
 'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
 'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7')}
 for name,(size,digest) in inputs.items():path=repo/'workplace/orig/dosv'/name;assert path.stat().st_size==size and sha(path)==digest
 p=json.loads((out/'ida/ida-probe.json').read_text());assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 routines,blocks={},{}
 for target in p['targets']:
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in target['chunks'])
  digest=hashlib.sha256(data).hexdigest();assert digest==target['file_sha256'] and len(data)==NEW[target['name']];routines[target['name']]=digest
 for b in p['decoded_blocks']:
  data=b''.join(bytes.fromhex(x['bytes']) for x in b['instructions']);assert len(data)==RAW[b['ida_linear']] and hashlib.sha256(data).hexdigest()==b['file_sha256']
  assert all(raw[x['ida_linear']-0x10000+512:x['ida_linear']-0x10000+512+len(bytes.fromhex(x['bytes']))].hex()==x['bytes'] for x in b['instructions'])
  blocks[f"code_{b['ida_linear']:X}"]=b['file_sha256']
 assert set(routines)==set(NEW) and sum(NEW.values())==1126 and sum(RAW.values())==744
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks'])+sum(len(b['instructions']) for b in p['decoded_blocks']);assert count==902
 receipts,full={},{}
 headers=[family*10+column for _ in range(4) for family,columns in enumerate([5,5,6,6,6,6,5]) for column in range(columns)]
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-catalog-parity-v1' and r['passed'] and r['mismatch'] is None
  assert r['input_sha256']==EXPECTED
  assert r['cases']==804 and r['groups']==GROUPS and r['full_ram_plane_audits']==r['indexed_content_audits']==804
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in {**routines,**blocks}.items())
  assert r['independent_builder_audits']==168 and r['independent_sort_audits']==240 and r['independent_cache_audits']==16
  assert r['accepted_selections']==28 and r['header_columns']==headers
  assert r['original_font_calls']==r['c_font_calls'] and sum(r['original_font_calls'])>0
  assert r['original_missing_fonts']==r['c_missing_fonts']==0 and r['original_font_misses']==r['c_font_misses']=={}
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch'];assert not r['passed'] and m['group']==g and r['groups']=={g:r['cases']}
  assert 0<r['cases']<=GROUPS[g] and m['case']==r['cases']-1
  assert any(m['original'+k]!=m['c'+k] for k in ['','_trace','_ports','_device','_planes','_ram','_api','_sound','_mouse','_in','_ticks'])
  info=(out/'results'/f'buildinfo-mutant-{n}.txt').read_text();assert f'-DKI_CATALOG_MUTATION={n}' in info and '-DKI_CATALOG_SOURCE_DIGEST=0x'+sha(out/'results/c-source.sha256')[:16] in info
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path),'buildinfo_sha256':sha(out/'results'/f'buildinfo-mutant-{n}.txt')}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']:
  info=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_CATALOG_SOURCE_DIGEST=0x'+compiled[:16] in info and 'matching_catalog' in info
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='catalog-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  gen_spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_catalog_generate.py');gen=importlib.util.module_from_spec(gen_spec);gen_spec.loader.exec_module(gen)
  gen.generate(root,out/'ida/ida-probe.json');assert (root/'tools/c_recovery/catalog_generated.inc').read_bytes()==(repo/'tools/c_recovery/catalog_generated.inc').read_bytes()
 evidence=repo/'docs/re/119-c-list-families-restoration.md';evidence_text=evidence.read_text();assert EXPECTED in evidence_text
 backlinks=[]
 for row in BACKLINKS:
  assert all(x in evidence_text for x in row['evidence_markers'])
  older=(repo/row['older_document']).read_text();assert all(x in older for x in row['required_markers'])
  backlinks.append({**row,'platform_module':'dosv/KI.EXE','input_sha256':EXPECTED,'level':'proven','evidence':'docs/re/119-c-list-families-restoration.md'})
 result={'schema':'wolong-c-catalog-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(out/'ida/ida-probe.json'),
 'new_routine_sha256':routines,'new_raw_entry_sha256':blocks,'original_instruction_count':902,
 'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_catalog.sh','tools/c_recovery_catalog_verify.py','tools/ida_catalog_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py','tools/c_recovery_talk_verify.py']},
 'input_assets':inputs,'cases_per_optimization':804,'groups':GROUPS,'receipt_sha256':receipts,
 'negative_controls_rejected':13,'mutants':controls,'entries_seen':full['O2']['entries_seen'],
 'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'independent_builder_audits':168,'independent_sort_audits':240,'independent_cache_audits':16,
 'accepted_selections':28,'headers_exercised':156,'backlink_contracts_verified':5,'resolution_backlinks':backlinks,
 'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,
 'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'Seven original corps/general/faction/new-game list callers and exact builders/renderers, 126/128/22 record domains, patched X/Y, entry DS:CFF, byte morale, affiliation/duty/sentinels, all original descriptor sorts and real headers, select/cancel and relation cache loop reset. Controlled runtime corps and prisoner data are explicit fixtures on four original scenarios. Complete RAM/planes/registers/FLAGS/stack/IN/OUT/API; independent pinned DOS/font/mouse/VGA; IF/TF=0; pinned dosgolem flags model. Invalid empty nonzero sort, natural appointment/diplomacy/new-game long paths and original C machine code remain pending',
 'c_machine_code_match':False}
 supplement=repo/'docs/re/c-catalog-code.json';sp=json.loads(supplement.read_text());assert len(sp['instructions'])==170
 assert sp['assembler_verification']['status']=='byte-exact' and sp['assembler_verification']['bytes']==361 and sp['constant_symbols']=={'imm_at_17284':65535}
 assert sp['assembler_verification']['generator_sha256']==sha(repo/'tools/catalog_code_supplement.py')
 go_proof_path=repo/'docs/re/c-catalog-go-verification.json';gp=json.loads(go_proof_path.read_text())
 assert gp['vet_passed'] and gp['test_passed'] and gp['cold_test_cache'] and gp['original_assets_readonly'] and gp['network']=='none'
 for name,digest in gp['source_sha256'].items():assert sha(repo/name)==digest,name
 result.update(assembly_code_supplement_sha256=sha(supplement),new_instruction_bytes=361,new_instruction_match=sp['assembler_verification'],go_verification_sha256=sha(go_proof_path))
 result['verification_tools_sha256']['tools/catalog_code_supplement.py']=sha(repo/'tools/catalog_code_supplement.py')
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('11 functions + 12 raw / 902 instructions / 804 whole-device cases / 13 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
