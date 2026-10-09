#!/usr/bin/env python3
"""世界初始化五函式與有限原始入口片段的來源及原生對拍驗證。"""
import argparse,hashlib,importlib.util,json,shlex,shutil,struct,tempfile
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_1533D':27,'sub_189F0':46,'sub_18A1E':179,'sub_18AD1':25,'sub_18AEA':40}
RAW_NAME='raw_entry_11BE0_prefix'
RAW_SHA='0c5138467de03c907f023374d58d25b2f72aec8cb487581777966c1f39d2d07e'
RAW_AUTHORITY='bounded initialization prefix of original named function; excludes entire main loop; no RET'
GROUPS={'city':964,'fringe':2048,'army':1052,'bulk':60,'score':464,'prefix':12,'chain':408}
FULL_AUDITS=5416
AUDIT_FIELD='independent_bootstrap_audits'
NATIVE_TAGS='matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga'
REPORT_METADATA={
 'raw_prefix_instruction_count':7,'raw_prefix_bytes':22,'outer_prefix_snapshots':0,
 'independent_raw_initialization':True,'native_prefix_entry':'bootstrap_call_prefix',
 'original_prefix_call_ip':0x0067,'prefix_stop_ip':0x1BF6,'prefix_has_synthetic_ret':False,
 'generic_full_11be0_supported':False,'cross_machine_snapshot_initialization':False,
 'stage_audits':{'helper':4588,'prefix':420,'chain-child':408},
 'segment_fixtures':{'D44=5000/D850=5000/occupancy=B000':4996,'D44=5000/D850=5100/occupancy=B000':4,'D44=5000/D850=5000/occupancy=C000':8},
 'nonlocal_exit_cases':144,'normal_return_cases':264,
 'scope':'Five original city/army map and scoring/date routines plus the finite seven-instruction 11BE0..11BF6 prefix. Both sides independently initialize raw runtime state; native C computes its own map, score, date and saved stack. Chain cases separately audit prefix and existing normal/nonlocal child closures. Whole main loop and C machine-code matching remain outside this slice.',
 'platform':'Independent pinned DOS/font/VGA services, full RAM/planes/palette/DAC/registers/FLAGS/stack/IN/OUT/API, controlled INT61 status and mature 3DA IN contract; no wall-clock claim',
 'game_palette_size':384,'game_palette_sha256':'1f0119c75ea5cd333bd3ac75ef92030f93924f011728edc9b1f727c483263708',
}
RESULT_FIELDS=('raw_prefix_instruction_count','raw_prefix_bytes','prefix_state_sha256','prefix_map_sha256','outer_prefix_snapshots','independent_raw_initialization','native_prefix_entry','original_prefix_call_ip','prefix_stop_ip','prefix_has_synthetic_ret','generic_full_11be0_supported','cross_machine_snapshot_initialization','stage_audits','segment_fixtures','palette_dac_audits','nonlocal_exit_cases','normal_return_cases','game_palette_size','game_palette_sha256','platform')
REUSED_SOURCES={'tools/c_recovery_main_data.go','tools/c_recovery_strategy_data.go'}
SCOPE=REPORT_METADATA['scope']+' Independent common raw input is shared before execution, and no original-prefix snapshot initializes native C. The raw fragment is not the complete named sub_11BE0; startup/main-loop continuation, TSR/hardware timing and original C machine-code matching remain outside this proof.'
BACKLINKS=[
 {'ida_linear':0x189F0,'function_ida_linear':0x189F0,'older_document':'docs/re/125-c-nonlocal-exit-restoration.md','required_markers':['初始化前段已由原生 C 驗證','126-c-world-bootstrap-restoration.md']},
 {'ida_linear':0x189F0,'function_ida_linear':0x189F0,'older_document':'docs/spec/245-c-nonlocal-exit.md','required_markers':['初始化前段已由原生 C 驗證','246-c-world-bootstrap.md']},
]
MUTANTS={1:'city',2:'city',3:'city',4:'fringe',5:'army',6:'army',7:'army',8:'bulk',9:'score',10:'prefix'}
MUTANT_NAMES={1:'center-divisor',2:'neutral-owner-colour',3:'own-fringe-colour',4:'fringe-upper-bound',5:'active-threshold',6:'descriptor-segment',7:'occupancy-increment',8:'last-corps-slot',9:'score-date-ds-result',10:'saved-prefix-sp'}
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def native_build(binary,info_path,compiled,mutation,optimize):
 """Read Go 1.26.7 inline ELF buildinfo without executing the binary or Go."""
 data=binary.read_bytes();assert data[:6]==b'\x7fELF\x02\x01',binary
 shoff=struct.unpack_from('<Q',data,40)[0]
 entsize,count,names_index=struct.unpack_from('<HHH',data,58)
 assert entsize>=64 and count>0 and names_index<count,binary
 sections=[struct.unpack_from('<IIQQQQIIQQ',data,shoff+i*entsize) for i in range(count)]
 names_section=sections[names_index];names=data[names_section[4]:names_section[4]+names_section[5]]
 entries=[section for section in sections if names[section[0]:].split(b'\0',1)[0]==b'.go.buildinfo']
 assert len(entries)==1,binary
 section=entries[0];blob=data[section[4]:section[4]+section[5]]
 assert blob[:14]==b'\xff Go buildinf:' and blob[14]==8 and blob[15]==2,binary
 def string_at(at):
  length=0;shift=0
  while True:
   assert at<len(blob) and shift<64,binary
   value=blob[at];at+=1;length|=(value&0x7f)<<shift
   if value<128:break
   shift+=7
  assert at+length<=len(blob),binary
  return blob[at:at+length],at+length
 version,at=string_at(32);module,_=string_at(at)
 assert version==b'go1.26.7' and len(module)>=33 and module[-17]==10,binary
 lines=module[16:-16].decode('utf-8').splitlines()
 assert 'path\tgithub.com/wicanr2/dosgolem/wolongcbootstrap' in lines,binary
 info=info_path.read_text();saved=info.splitlines()
 assert saved[0].endswith(': go1.26.7') and Path(saved[0].split(':',1)[0]).name==binary.name,info_path
 def flags_from(lines):
  prefix='build\tCGO_CFLAGS='
  values=[line.removeprefix('\t')[len(prefix):] for line in lines if line.removeprefix('\t').startswith(prefix)]
  assert len(values)==1
  quoted=shlex.split(values[0]);assert len(quoted)==1
  return shlex.split(quoted[0])
 flags=flags_from(lines);assert flags==flags_from(saved[1:]),binary
 digest='-DKI_BOOTSTRAP_SOURCE_DIGEST=0x'+compiled[:16]
 assert [value for value in flags if value.startswith('-DKI_BOOTSTRAP_SOURCE_DIGEST=')]==[digest],binary
 mutations=[value for value in flags if value.startswith('-DKI_BOOTSTRAP_MUTATION=')]
 assert mutations==([] if mutation is None else [f'-DKI_BOOTSTRAP_MUTATION={mutation}']),binary
 assert '-'+optimize in flags and 'build\t-tags='+NATIVE_TAGS in lines,binary
 return {'c_binary_sha256':hashlib.sha256(data).hexdigest(),'buildinfo_sha256':sha(info_path),'go_version':version.decode(),'cgo_flags':flags}
def verify(repo,out):
 if any(value is None for value in (GROUPS,AUDIT_FIELD,NATIVE_TAGS,REPORT_METADATA,RESULT_FIELDS,REUSED_SOURCES,SCOPE)):
  raise RuntimeError('bootstrap matrix/metadata contract pending; do not substitute guessed counts')
 assert GROUPS and all(isinstance(n,int) and n>0 for n in GROUPS.values())
 assert set(MUTANTS)==set(MUTANT_NAMES)==set(range(1,11))
 assert all(group in GROUPS for group in MUTANTS.values())
 case_count=sum(GROUPS.values())
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 asset_spec=importlib.util.spec_from_file_location('assets',repo/'tools/c_recovery_talk_verify.py');assets=importlib.util.module_from_spec(asset_spec);asset_spec.loader.exec_module(assets)
 inputs={**assets.ASSETS,'GAMEPAL.BRG':(384, '1f0119c75ea5cd333bd3ac75ef92030f93924f011728edc9b1f727c483263708'),'KYOGRF.DAT':(69120,'e086f526bdada5baf751c41d2f73a78a0ba70002f282f63a9f33114542ed933f'),'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
 'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
 'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7')}
 for name,(size,digest) in inputs.items():
  path=repo/'workplace/orig/dosv'/name;assert path.stat().st_size==size and sha(path)==digest,name
 probe=out/'ida/ida-probe.json';p=json.loads(probe.read_text())
 assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 assert p['image_id']=='sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab'
 assert set(p['recovery_targets'])=={0x1533D,0x189F0,0x18A1E,0x18AD1,0x18AEA}
 assert p['navigation_only_targets']==[0x11BE0]
 routines={}
 for target in p['targets']:
  if target['ida_linear'] not in p['recovery_targets']:continue
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in target['chunks'])
  assert len(data)==NEW[target['name']] and hashlib.sha256(data).hexdigest()==target['file_sha256']
  assert not any(c['loader_relocations'] for c in target['chunks'])
  routines[target['name']]=target['file_sha256']
 assert set(routines)==set(NEW) and sum(NEW.values())==317
 named_count=sum(len(c['instructions']) for t in p['targets'] if t['ida_linear'] in p['recovery_targets'] for c in t['chunks'])
 assert named_count==140
 assert len(p['decoded_blocks'])==1
 block=p['decoded_blocks'][0];rows=block['instructions']
 assert block['ida_linear']==0x11BE0 and block['end_ida_linear']==0x11BF6
 assert block['original_name_before_analysis']=='sub_11BE0' and block['boundary_authority']==RAW_AUTHORITY
 assert len(rows)==7 and not any(row['mnemonic'] in ('ret','retn','retf','iret') for row in rows)
 fragment=b''.join(bytes.fromhex(row['bytes']) for row in rows)
 assert fragment==raw[0x1BE0+512:0x1BF6+512] and len(fragment)==22
 assert hashlib.sha256(fragment).hexdigest()==block['file_sha256']==RAW_SHA
 raw_hashes={RAW_NAME:RAW_SHA}
 receipts={};full={}
 for level in ('O0','O2'):
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-bootstrap-parity-v1' and r['passed'] and r['mismatch'] is None
  assert r['input_sha256']==EXPECTED and r['groups']==GROUPS and r['cases']==case_count
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==r['palette_dac_audits']==FULL_AUDITS
  assert r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
  assert r['new_routine_sha256']==routines and r['new_code_block_sha256']==raw_hashes
  assert all(r['routine_sha256'][name]==digest and r['entries_seen'].get(name,0)>0 for name,digest in routines.items())
  assert r[AUDIT_FIELD]>=case_count
  for key,want in REPORT_METADATA.items():assert r[key]==want,(level,key)
  assert r['new_raw_sha256']==raw_hashes
  assert sum(r['stage_audits'].values())==FULL_AUDITS
  assert sum(r['segment_fixtures'].values())==case_count
  assert r['nonlocal_exit_cases']+r['normal_return_cases']==GROUPS['chain']
  assert len(r['prefix_state_sha256'])==len(r['prefix_map_sha256'])==GROUPS['prefix']==12
  assert set(r['prefix_state_sha256'])==set(r['prefix_map_sha256'])
  for key in r['prefix_state_sha256']:
   assert len(bytes.fromhex(r['prefix_state_sha256'][key]))==len(bytes.fromhex(r['prefix_map_sha256'][key]))==32
  assert r['original_font_calls']==r['c_font_calls']
  assert r['original_missing_fonts']==r['c_missing_fonts']==0 and r['original_font_misses']==r['c_font_misses']=={}
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/')
  assert sha(repo/path)==digest,path;source[path]=digest
 required={'tools/c_recovery/bootstrap.c','tools/c_recovery/bootstrap.h','tools/c_recovery/bootstrap_fixture.h','tools/c_recovery/bootstrap_generated.inc','tools/c_recovery_bootstrap.go','tools/c_recovery_bootstrap_data.go'}
 assert required.union(REUSED_SOURCES).issubset(source)
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 manifest_binding={'path':'workplace/matching-decompilation/c-bootstrap/results/c-source.sha256','sha256':compiled}
 assert all(r['source_manifest']==manifest_binding for r in full.values())
 native={}
 for level in ['O0','O2']+[f'mutant-{n}' for n in MUTANTS]:
  mutation=None if level in ('O0','O2') else int(level.split('-')[1])
  binary=out/(f'bootstrap-{level}' if mutation is None else level)
  native[level]=native_build(binary,out/'results'/f'buildinfo-{level}.txt',compiled,mutation,level if mutation is None else 'O0')
 names_path=out/'results/mutant-names.tsv';names={}
 for line in names_path.read_text().splitlines():
  fields=line.split('\t');assert len(fields)==3
  n=int(fields[0]);assert n not in names;names[n]=(fields[1],fields[2])
 assert names=={n:(group,MUTANT_NAMES[n]) for n,group in MUTANTS.items()}
 controls={}
 for n,group in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json';log=out/'results'/f'mutant-{n}.log'
  assert (out/'results'/f'mutant-{n}.exit').read_text().strip()=='1'
  r=json.loads(path.read_text());m=r['mismatch']
  assert r['schema']=='wolong-c-bootstrap-parity-v1' and r['input_sha256']==EXPECTED and r['source_manifest']==manifest_binding
  assert not r['passed'] and m['group']==group and r['groups']=={group:r['cases']}
  assert 0<r['cases']<=GROUPS[group] and m['case']==r['cases']-1
  pairs=['','_trace','_ports','_planes','_ram','_api','_in','_palette_dac']
  different=any(m['original'+key]!=m['c'+key] for key in pairs)
  transfer_difference=m['native_transferred']!=int(m['expected_escape'])
  if m['expected_escape']:
   transfer_difference|=any(m['native_transfer_'+key.lower()]!=m['original'][key] for key in ('IP','SS','SP'))
  assert different or transfer_difference
  controls[str(n)]={'group':group,'name':MUTANT_NAMES[n],'rejection_kind':'state-mismatch','exit_status':1,'first_rejected':r['cases'],'receipt_sha256':sha(path),'log_sha256':sha(log),**native[f'mutant-{n}']}
 before=(out/'results/golem-source-before.sha256').read_bytes()
 assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='bootstrap-repro-') as temporary:
  root=Path(temporary);(root/'tools/c_recovery').mkdir(parents=True)
  shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_bootstrap_generate.py')
  gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen);gen.generate(root,probe)
  assert (root/'tools/c_recovery/bootstrap_generated.inc').read_bytes()==(repo/'tools/c_recovery/bootstrap_generated.inc').read_bytes()
 code=repo/'docs/re/c-bootstrap-code.json';sp=json.loads(code.read_text());ap=sp['assembler_verification']
 assert sp['input_sha256']==EXPECTED and sp['probe_sha256']==sha(probe) and sp['database_sha256']==sha(out/'ida/input.exe.i64')
 assert set(sp['recovery_targets'])==set(p['recovery_targets']) and sp['navigation_only_targets']==[0x11BE0]
 assert not sp['instructions'] and sp['constant_symbols']=={} and sp['loader_relocation_normalization']==[]
 assert sp['named_instruction_count']==140 and sp['named_instruction_bytes']==317
 assert sp['raw_instruction_count']==7 and sp['raw_instruction_bytes']==22
 assert ap['status']=='byte-exact' and ap['instructions']==ap['bytes']==0
 assert ap['unique_input_instructions']==ap['covered_input_instructions']==147
 independent=ap['independent_reassembly'];assert independent['status']=='byte-exact' and independent['instructions']==147 and independent['bytes']==339
 assert ap['fixed_prior_instruction_count']==24999 and ap['fixed_prior_instruction_bytes']==56866
 assert len(sp['instruction_coverage'])==147 and all(row['fixed_prior_coverage'] and row['independent_native_assembly_match'] for row in sp['instruction_coverage'])
 for row in sp['instruction_coverage']:assert raw[row['file_start']:row['file_end']].hex()==row['bytes'] and row['file_start']==row['ida_linear']-0x10000+512
 assert len(sp['original_functions'])==5 and {target['name']:target['file_sha256'] for target in sp['original_functions']}==routines
 assert len(sp['original_code_blocks'])==1
 recorded=sp['original_code_blocks'][0]
 assert recorded['ida_linear']==0x11BE0 and recorded['end_ida_linear']==0x11BF6 and recorded['original_name_before_analysis']=='sub_11BE0'
 assert recorded['source_entry']==RAW_NAME and recorded['file_sha256']==RAW_SHA and recorded['boundary_authority']==RAW_AUTHORITY
 assert ap['published_manifest_sha256']==sha(repo/ap['published_manifest']) and ap['published_code_source_sha256']==sha(repo/ap['published_code_source'])
 assert ap['generator_sha256']==sha(repo/'tools/bootstrap_code_supplement.py') and ap['translator_sha256']==sha(repo/'tools/assembly_rebuild.py')
 assert ap['image_id']=='sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e'
 for relative,binding in ap['coverage_inputs'].items():assert sha(repo/relative)==binding['sha256'],relative
 assert {name:binding['instruction_count'] for name,binding in ap['coverage_inputs'].items()}=={'workplace/matching-decompilation/assembly/build/source-map.json':24376,'docs/re/rectangle-handler-code.json':370,'docs/re/c-list-code.json':83,'docs/re/c-catalog-code.json':170}
 result={'schema':'wolong-c-bootstrap-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(probe),
 'new_routine_sha256':routines,'new_code_block_sha256':raw_hashes,
 'original_instruction_count':147,'original_instruction_bytes':339,'named_instruction_count':140,'named_instruction_bytes':317,'raw_instruction_count':7,'raw_instruction_bytes':22,
 'source_sha256':source,'compiled_source_manifest_sha256':compiled,'native_builds':native,
 'cases_per_optimization':case_count,'whole_device_stage_audits_per_optimization':FULL_AUDITS,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':10,'state_mismatch_controls_rejected':10,'mutants':controls,
 'mutant_names_sha256':sha(names_path),'entries_seen':full['O2']['entries_seen'],'independent_bootstrap_audits':full['O2'][AUDIT_FIELD],
 'input_assets':inputs,'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,
 'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'assembly_coverage_receipt_sha256':sha(code),'recovery_targets':p['recovery_targets'],'navigation_only_targets':p['navigation_only_targets'],
 'raw_code_blocks':{'0x11BE0':recorded},'scope':SCOPE,'c_machine_code_match':False,
 'verification_tools_sha256':{name:sha(repo/name) for name in ('tools/c_recovery_bootstrap.sh','tools/c_recovery_bootstrap_verify.py','tools/c_recovery_bootstrap_generate.py','tools/ida_bootstrap_probe.py','tools/bootstrap_code_supplement.py','tools/c_recovery_display_generate.py','tools/c_recovery_talk_verify.py')}}
 for key in RESULT_FIELDS:result[key]=full['O2'][key]
 go=repo/'docs/re/c-bootstrap-go-verification.json'
 if go.exists():
  proof=json.loads(go.read_text());assert proof['vet_passed'] and proof['test_passed'] and proof['cold_test_cache']
  assert proof['original_assets_readonly'] and proof['network']=='none' and proof['tested_packages']==39 and proof['cached_test_packages']==0
  for name,digest in proof['source_sha256'].items():assert sha(repo/name)==digest,name
  result['go_verification_sha256']=sha(go)
 evidence_name='docs/re/126-c-world-bootstrap-restoration.md';evidence=repo/evidence_name;assert EXPECTED in evidence.read_text()
 result['resolution_backlinks']=[];result['backlink_status']='pending-publication'
 if '**狀態：CONFORMED' in evidence.read_text():
  assert BACKLINKS is not None,'fixed bootstrap correction markers are pending'
  ledger=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
  assert ledger['input_sha256']==EXPECTED and len(ledger['functions'])==375 and len(ledger['code_blocks'])==49
  assert all(name in ledger['functions'] for name in NEW) and 'sub_11BE0' not in ledger['functions']
  record=ledger['code_blocks']['0x11BE0']
  assert record['ida_linear']==0x11BE0 and record['original_name_before_analysis']=='sub_11BE0'
  assert record['boundary_authority']==RAW_AUTHORITY and record['c_machine_code_match'] is False
  for row in BACKLINKS:
   assert f"0x{row['ida_linear']:X}" in evidence.read_text(),row
   older=(repo/row['older_document']).read_text();assert all(marker in older for marker in row['required_markers']),row
   want={**row,'platform_module':'dosv/KI.EXE','input_sha256':EXPECTED,'level':'proven','evidence':evidence_name,
         'evidence_markers':[f"sub_{row['function_ida_linear']:X}",f"0x{row['ida_linear']:X}"]}
   assert all(marker in evidence.read_text() for marker in want['evidence_markers']),row
   matches=[actual for actual in ledger['resolution_backlinks'] if all(actual.get(key)==want[key] for key in ('platform_module','input_sha256','ida_linear','older_document'))]
   assert len(matches)==1 and all(matches[0].get(key)==value for key,value in want.items()),row
   result['resolution_backlinks'].append(want)
  result['backlink_status']='verified'
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
 print(f'5 named functions + 7-instruction raw prefix / 147 instructions / {case_count} whole-device cases / 10 state mismatches: PASS')
if __name__=='__main__':
 parser=argparse.ArgumentParser(description=__doc__);parser.add_argument('--repo',type=Path,required=True);parser.add_argument('--output',type=Path,required=True)
 args=parser.parse_args();verify(args.repo,args.output)
