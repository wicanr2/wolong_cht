#!/usr/bin/env python3
"""核對原據點清單、特殊返回、排序與完整遷都的原版／C 收據。"""
import argparse,hashlib,importlib.util,json,shutil,tempfile
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_15E60':32,'sub_16909':149,'sub_17400':59,'sub_1748F':193,'sub_181C0':78,
 'sub_1820E':129,'sub_18412':70,'nullsub_2':1,'sub_18463':89,'sub_184BC':33,'sub_184DD':61,
 'sub_1851A':44,'sub_18546':57,'sub_18607':91,'sub_18662':103,'sub_186C9':50,'sub_18713':66,'sub_18755':90}
RAW={0x1743B:36,0x1745F:48,0x1828F:387,0x1857F:51,0x185B2:85}
GROUPS={'rows':108,'list-cancel':72,'sort':120,'scroll':132,'scroll-map':36,'selection':12,'relocation':16,'header':24}
MUTANTS={1:'rows',2:'rows',3:'list-cancel',4:'sort',5:'sort',6:'selection',7:'scroll',8:'scroll-map',9:'relocation',10:'relocation',11:'relocation',12:'selection'}
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def verify(repo,out):
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 asset_spec=importlib.util.spec_from_file_location('assets',repo/'tools/c_recovery_talk_verify.py');assets=importlib.util.module_from_spec(asset_spec);asset_spec.loader.exec_module(assets)
 inputs={**assets.ASSETS,
  'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
  'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
  'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7'),
  'BGM.DAT':(20826,'7a51c8b9a349b9e088f3796b70c268181c60bcebead70942f00e1621523dedc9'),
  'IVENTGRF.DAT':(76032,'23fffa03b2bba3b4c920c5db49b32b7a3fa26729523b0092a6be5ad3a0737f2f')}
 for name,(size,digest) in inputs.items():
  path=repo/'workplace/orig/dosv'/name;assert path.stat().st_size==size and sha(path)==digest,name
 p=json.loads((out/'ida/ida-probe.json').read_text());assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 routines={};blocks={}
 for target in p['targets']:
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in target['chunks'])
  digest=hashlib.sha256(data).hexdigest();assert digest==target['file_sha256']
  assert target['name'] in NEW and len(data)==NEW[target['name']];routines[target['name']]=digest
 for b in p['decoded_blocks']:
  rows=b['instructions'];data=b''.join(bytes.fromhex(x['bytes']) for x in rows)
  assert len(data)==RAW[b['ida_linear']] and hashlib.sha256(data).hexdigest()==b['file_sha256']
  assert all(raw[x['ida_linear']-0x10000+512:x['ida_linear']-0x10000+512+len(bytes.fromhex(x['bytes']))].hex()==x['bytes'] for x in rows)
  blocks[f"code_{b['ida_linear']:X}"]=b['file_sha256']
 assert set(routines)==set(NEW) and sum(NEW.values())==1395 and sum(RAW.values())==607
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks'])+sum(len(b['instructions']) for b in p['decoded_blocks']);assert count==835
 receipts,full={},{}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-list-parity-v1' and r['input_sha256']==EXPECTED
  assert r['passed'] and r['mismatch'] is None and r['cases']==520 and r['groups']==GROUPS
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==520
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in {**routines,**blocks}.items())
  assert r['entries_seen']['sub_16909']==16 and r['entries_seen']['sub_133FD']==8
  assert r['independent_sort_audits']==120 and r['independent_builder_audits']==36 and r['independent_scroll_audits']==48
  assert r['header_columns']==list(range(6))*4
  assert r['original_font_calls']==r['c_font_calls'] and sum(r['original_font_calls'])>0
  assert r['original_missing_fonts']==r['c_missing_fonts']==0
  assert r['original_font_misses']==r['c_font_misses']=={}
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 controls={}
 for n,g in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch']
  assert not r['passed'] and r['groups']=={g:r['cases']} and m['group']==g
  assert 0<r['cases']<=GROUPS[g] and m['case']==r['cases']-1
  assert any(m['original'+k]!=m['c'+k] for k in ['','_trace','_ports','_device','_planes','_ram','_api','_sound','_mouse','_in','_ticks'])
  buildinfo=(out/'results'/f'buildinfo-mutant-{n}.txt').read_text();assert f'-DKI_LIST_MUTATION={n}' in buildinfo
  assert '-DKI_LIST_SOURCE_DIGEST=0x'+sha(out/'results/c-source.sha256')[:16] in buildinfo
  controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path),'buildinfo_sha256':sha(out/'results'/f'buildinfo-mutant-{n}.txt')}
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest;source[path]=digest
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 for level in ['O0','O2']:
  text=(out/'results'/f'buildinfo-{level}.txt').read_text();assert '-DKI_LIST_SOURCE_DIGEST=0x'+compiled[:16] in text and 'matching_list' in text
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='list-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_list_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,out/'ida/ida-probe.json');assert (root/'tools/c_recovery/list_generated.inc').read_bytes()==(repo/'tools/c_recovery/list_generated.inc').read_bytes()
 result={'schema':'wolong-c-list-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(out/'ida/ida-probe.json'),
 'new_routine_sha256':routines,'new_raw_entry_sha256':blocks,'original_instruction_count':count,
 'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_list.sh','tools/c_recovery_list_verify.py','tools/ida_list_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py','tools/list_code_supplement.py','tools/assembly_rebuild.py','tools/c_recovery_talk_verify.py']},
 'cases_per_optimization':520,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':12,'mutants':controls,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'independent_sort_audits':120,'independent_builder_audits':36,'independent_scroll_audits':48,
 'input_assets':inputs,'tool_versions':(out/'results/tool-versions.txt').read_text(),
 'exact_clean_regeneration':True,'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
 'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':'Original city list callbacks and six descriptors (header 0 rebuilds, 1..5 sort), live byte/word/JBE/JNB sort, exact exchange order, 512-byte SS list, special return-address removal, full row/font/number/scrollbar rendering, select/cancel/drag, and four-scenario relocation accept/refuse/current-capital retry through real mouse input and original world/sidebar/scene writers. Complete RAM/planes/registers/FLAGS/stack/IN/OUT/API; independent pinned DOS/font/mouse/sound/VGA; IF/TF=0; isolated stack; FLAGS from pinned dosgolem. Natural player menus/long paths, other list families and original C machine code remain pending',
 'c_machine_code_match':False}
 supplement=repo/'docs/re/c-list-code.json';sp=json.loads(supplement.read_text())
 assert sp['assembler_verification']['status']=='byte-exact' and sp['assembler_verification']['instructions']==83 and sp['assembler_verification']['bytes']==215
 assert {0x1831f,0x183d2,0x1858f,0x185e8,0x185ef,0x185f1,0x185f3}.issubset({x['ida_linear'] for x in sp['instructions']})
 result.update(assembly_code_supplement_sha256=sha(supplement),new_instruction_bytes=215,new_instruction_match=sp['assembler_verification'])
 ledger=json.loads((repo/'docs/re/c-recovery-status.json').read_text());backlinks=ledger['resolution_backlinks']
 assert len(backlinks)==3 and {x['ida_linear'] for x in backlinks}=={0x183d2,0x16909}
 for row in backlinks:
  assert row['platform_module']=='dosv/KI.EXE' and row['input_sha256']==EXPECTED and row['level']=='proven'
  evidence=(repo/row['evidence']).read_text();older=(repo/row['older_document']).read_text()
  assert EXPECTED in evidence and all(x in evidence for x in row['evidence_markers'])
  assert all(x in older for x in row['required_markers']),row['older_document']
 result.update(backlink_contracts_verified=3,resolution_backlinks=backlinks)
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print('18 functions + 5 raw / 835 instructions / 520 whole-device cases / 12 mutants: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
