#!/usr/bin/env python3
"""非局部退出與色盤 C 原始指令、完整裝置比較、錯版與來源綁定驗證。"""
import argparse,hashlib,importlib.util,json,tempfile,shutil,shlex,struct
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_10A1C':73,'sub_11CB1':24,'sub_1EBDC':82}
GROUPS={'palette':960,'fade':16,'escape':96,'trust':160,'scene':248,'reason':960,'reason-loop':80}
PREFIX_MAP_EXPECTED={
 '0':{'changed_bytes':180,'sha256':'51daebff50e31bfc3acb4f5b8735e7a60e0b070c99e23d7858081663661a45ce'},
 '1':{'changed_bytes':450,'sha256':'10a4be4d2e630bd58242d517c6602fe4a1659a5c8b39c7bad86c541096b90ed3'},
 '2':{'changed_bytes':538,'sha256':'b3b04766894b95987c7ae23cd186ec0f62a59bec9172bda27a0ac8d007e6c9c4'},
 '3':{'changed_bytes':584,'sha256':'e445787934f9724aee1b7b927166c25439029a60c4fe8c81bf03404b02dd66a4'},
}
MUTANTS={1:'escape',2:'escape',3:'escape',4:'escape',5:'fade',6:'palette',7:'palette',8:'trust'}
MUTANT_NAMES={1:'saved-stack-segment',2:'saved-stack-pointer',3:'exit-ax-push',4:'exit-ax-restore',5:'fade-first-level',6:'dac-rounding',7:'dac-channel-order',8:'nonlocal-unwind'}
BACKLINKS=[
 {'ida_linear':0x11CB1,'function_ida_linear':0x11CB1,'older_document':'docs/re/124-c-strategy-command-restoration.md','required_markers':['非區域返回已由三函式閉包驗證','125-c-nonlocal-exit-restoration.md']},
 {'ida_linear':0x11CB1,'function_ida_linear':0x11CB1,'older_document':'docs/spec/244-c-strategy-command.md','required_markers':['非區域返回已由三函式閉包驗證','245-c-nonlocal-exit.md']},
 {'ida_linear':0x159C0,'function_ida_linear':0x159B7,'older_document':'docs/re/71-strategy-hotspot-dispatch.md','required_markers':['2026-10-09 勘誤（re/125，已證實）','125-c-nonlocal-exit-restoration.md','FF 97 06 5A','0x15A06','26 個 word','0x159D0 RET']},
 {'ida_linear':0x159C0,'function_ida_linear':0x159B7,'older_document':'docs/re/47-main-screen-window-registry.md','required_markers':['2026-10-09 勘誤（re/125，已證實）','125-c-nonlocal-exit-restoration.md','FF 97 06 5A','0x15A06','26 個 word']},
 {'ida_linear':0x159C0,'function_ida_linear':0x159B7,'older_document':'CONTEXT.md','required_markers':['2026-10-09 勘誤（re/125，已證實）','125-c-nonlocal-exit-restoration.md','右鍵表基址與落下返回勘誤','FF 97 06 5A','0x15A06']},
 {'ida_linear':0x159CE,'function_ida_linear':0x159B7,'older_document':'docs/mechanics/10-strategy.md','required_markers':['2026-10-09 勘誤（re/125，已證實）','125-c-nonlocal-exit-restoration.md','0x159CE','0x159D0 RET']},
]
AUDIT_FIELD='independent_main_audits'
DRIVER_SCOPE='Three original saved-stack exit/fade/palette routines with original CALL10067 prefix; four scenarios, independent deep snapshots, real nonlocal transfer and normal returns, complete RAM/VGA/palette/DAC/ports/state comparison. Direct routine closure; no wall-clock or complete natural player-path claim.'
NONLOCAL_SCOPE='Original CALL10067 prefix establishes CS9901/9903 and outer return006A. Trust, scene and reason underflow restore that saved stack and leave old callers through the live C runner; normal CF paths return to their independent entry frames.'
PLATFORM_SCOPE='Independent pinned DOS/font/VGA services; native C saved-stack exit and palette control; mature Machine.In8 3DA samples; IF/TF=0; pinned CPU-model undefined flags; explicit DAC and INT61 status fixtures; no wall-clock or TSR hardware claim'
SCOPE=DRIVER_SCOPE+' Original prefixes are supplied as common initial snapshots; navigation routines, complete startup/main-loop execution, TSR/hardware timing and original C machine-code matching are outside this recovery claim.'
RNG_SETUP='Raw RNG state is preserved from the fixed original KI.EXE image and scenario header for both original/native; no additional seed injection, no rerolls and no deliberately exercised RNG path. Actual RNG and game-clock entries are reported without a nonzero-call requirement.'
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
 assert 'path\tgithub.com/wicanr2/dosgolem/wolongcmain' in lines,binary
 info=info_path.read_text();saved=info.splitlines()
 assert saved[0].endswith(': go1.26.7') and Path(saved[0].split(':',1)[0]).name==binary.name,info_path
 def flags_from(lines):
  prefix='build\tCGO_CFLAGS='
  values=[line.removeprefix('\t')[len(prefix):] for line in lines if line.removeprefix('\t').startswith(prefix)]
  assert len(values)==1
  quoted=shlex.split(values[0]);assert len(quoted)==1
  return shlex.split(quoted[0])
 flags=flags_from(lines);assert flags==flags_from(saved[1:]),binary
 digest='-DKI_MAIN_SOURCE_DIGEST=0x'+compiled[:16]
 assert [value for value in flags if value.startswith('-DKI_MAIN_SOURCE_DIGEST=')]==[digest],binary
 mutations=[value for value in flags if value.startswith('-DKI_MAIN_MUTATION=')]
 assert mutations==([] if mutation is None else [f'-DKI_MAIN_MUTATION={mutation}']),binary
 assert '-'+optimize in flags and 'build\t-tags=matching_main,matching_strategy,matching_input,matching_glyph,matching_vga' in lines,binary
 return {'c_binary_sha256':hashlib.sha256(data).hexdigest(),'buildinfo_sha256':sha(info_path),'go_version':version.decode(),'cgo_flags':flags}
def verify(repo,out):
 if any(v is None for v in (GROUPS,MUTANTS,BACKLINKS,AUDIT_FIELD,SCOPE,RNG_SETUP)):
  raise RuntimeError('main driver matrix contract pending; do not substitute guessed counts')
 assert isinstance(GROUPS,dict) and GROUPS and all(isinstance(g,str) and isinstance(n,int) and n>0 for g,n in GROUPS.items())
 assert isinstance(MUTANTS,dict) and MUTANTS and all(isinstance(n,int) and n>0 and g in GROUPS for n,g in MUTANTS.items())
 case_count=sum(GROUPS.values())
 raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
 asset_spec=importlib.util.spec_from_file_location('assets',repo/'tools/c_recovery_talk_verify.py');assets=importlib.util.module_from_spec(asset_spec);asset_spec.loader.exec_module(assets)
 inputs={**assets.ASSETS,'GAMEPAL.BRG':(384, '1f0119c75ea5cd333bd3ac75ef92030f93924f011728edc9b1f727c483263708'),'KYOGRF.DAT':(69120,'e086f526bdada5baf751c41d2f73a78a0ba70002f282f63a9f33114542ed933f'),'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
 'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
 'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7')}
 for name,(size,digest) in inputs.items():
  path=repo/'workplace/orig/dosv'/name;assert path.stat().st_size==size and sha(path)==digest,name
 game_palette=(repo/'workplace/orig/dosv/GAMEPAL.BRG').read_bytes()
 assert len(game_palette)==384 and hashlib.sha256(game_palette).hexdigest()==inputs['GAMEPAL.BRG'][1]
 working_palette=game_palette[48:96];assert len(working_palette)==48
 probe=out/'ida/ida-probe.json';p=json.loads(probe.read_text());assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
 assert p['image_id']=='sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab'
 full_probe=p
 assert set(p["recovery_targets"])=={0x10A1C,0x11CB1,0x1EBDC}
 assert set(p["navigation_only_targets"])=={0x11BE0,0x1533D,0x159B7,0x189F0}
 p={**p,"targets":[t for t in p["targets"] if t["ida_linear"] in p["recovery_targets"]],"decoded_blocks":p.get("decoded_blocks",[])}
 routines={}
 for t in p['targets']:
  data=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks'])
  assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
  assert len(data)==NEW[t['name']] and hashlib.sha256(data).hexdigest()==t['file_sha256'];routines[t['name']]=t['file_sha256']
 assert set(routines)==set(NEW) and len(routines)==3 and sum(NEW.values())==179 and not p['decoded_blocks']
 count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks']);assert count==89
 receipts={};full={}
 for level in ['O0','O2']:
  path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
  assert r['schema']=='wolong-c-main-parity-v1' and r['passed'] and r['mismatch'] is None
  assert r['input_sha256']==EXPECTED and r['groups']==GROUPS and r['cases']==case_count
  assert r['scope']==DRIVER_SCOPE and 'returning_scene_only' not in r and r['nonlocal_scope']==NONLOCAL_SCOPE
  assert r['platform']==PLATFORM_SCOPE and r['controlled_sound_status_fixture'] is True and r['matrix_roots']==3
  assert r['full_ram_plane_audits']==r['indexed_content_audits']==case_count
  assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
  assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in routines.items())
  assert r[AUDIT_FIELD]>=case_count
  assert r["new_routine_sha256"]==routines
  assert len(r["navigation_only_targets"])==4 and set(r["navigation_only_targets"])=={"sub_11BE0","sub_1533D","sub_159B7","sub_189F0"}
  assert r["palette_dac_audits"]==case_count
  assert r["outer_prefix_snapshots"]==12 and len(r["prefix_instructions"])==12
  assert all(isinstance(n,int) and n>0 for n in r["prefix_instructions"].values())
  assert r["raw_map_verified_before_prefix"] is True and r["prefix_selected_player"]==0
  assert set(r["prefix_map_sha256"])==set(r["prefix_map_changed_bytes"])==set(r["prefix_instructions"])
  scenario_counts={scenario:0 for scenario in PREFIX_MAP_EXPECTED}
  for key,digest in r["prefix_map_sha256"].items():
   scenario=key.split(':')[0];assert scenario in PREFIX_MAP_EXPECTED,key
   expected=PREFIX_MAP_EXPECTED[scenario];scenario_counts[scenario]+=1
   assert digest==expected["sha256"] and r["prefix_map_changed_bytes"][key]==expected["changed_bytes"],key
  assert all(count==3 for count in scenario_counts.values())
  assert r["nonlocal_exit_cases"]>0 and r["normal_return_cases"]>0
  assert r["nonlocal_exit_cases"]+r["normal_return_cases"]==case_count
  assert r["dac_initialization"]=="Explicit fixture: DAC[index]=byte(index*13+seed*7)&63 before the deep snapshot; mainBeforeDAC preserves the full 768-byte input"
  assert r["palette_profiles"]==["GAMEPAL.BRG[48:96]","zero","all nibble15","nibble ramp byte(i*7+3)&15"]
  assert r["game_palette_size"]==inputs['GAMEPAL.BRG'][0]==384 and r["game_palette_sha256"]==inputs['GAMEPAL.BRG'][1]
  assert r['original_font_calls']==r['c_font_calls'] and sum(r['original_font_calls'])>0
  assert r['original_missing_fonts']==r['c_missing_fonts']==0 and r['original_font_misses']==r['c_font_misses']=={}
  receipts[level]=sha(path);full[level]=r
 assert receipts['O0']==receipts['O2']
 source={}
 for line in (out/'results/c-source.sha256').read_text().splitlines():
  digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/');assert sha(repo/path)==digest,path;source[path]=digest
 assert {"tools/c_recovery/main.c","tools/c_recovery/main.h","tools/c_recovery/main_fixture.h","tools/c_recovery/main_generated.inc","tools/c_recovery_main.go","tools/c_recovery_main_data.go","tools/c_recovery_strategy_data.go"}.issubset(source)
 compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
 manifest_binding={'path':'workplace/matching-decompilation/c-main/results/c-source.sha256','sha256':compiled}
 assert all(r['source_manifest']==manifest_binding for r in full.values())
 native={}
 for level in ['O0','O2']+[f'mutant-{n}' for n in MUTANTS]:
  mutation=None if level in ['O0','O2'] else int(level.split('-')[1])
  binary=out/(f'main-{level}' if mutation is None else level)
  native[level]=native_build(binary,out/'results'/f'buildinfo-{level}.txt',compiled,mutation,level if mutation is None else 'O0')
 controls={}
 names_path=out/'results/mutant-names.tsv';names={}
 for line in names_path.read_text().splitlines():
  fields=line.split('\t');assert len(fields)==3
  n=int(fields[0]);assert n not in names;names[n]=(fields[1],fields[2])
 assert set(names)==set(MUTANTS)==set(MUTANT_NAMES)
 assert all(names[n]==(group,MUTANT_NAMES[n]) for n,group in MUTANTS.items())
 for n,g in MUTANTS.items():
  path=out/'results'/f'mutant-{n}.json'
  exit_path=out/'results'/f'mutant-{n}.exit';log_path=out/'results'/f'mutant-{n}.log'
  exit_status=int(exit_path.read_text().strip())
  assert exit_status==1
  r=json.loads(path.read_text());m=r['mismatch']
  assert r['schema']=='wolong-c-main-parity-v1' and r['input_sha256']==EXPECTED
  assert r['source_manifest']==manifest_binding
  assert not r['passed'] and m['group']==g and r['groups']=={g:r['cases']}
  assert 0<r['cases']<=GROUPS[g] and m['case']==r['cases']-1
  different=any(m["original"+k]!=m["c"+k] for k in ["","_trace","_ports","_device","_planes","_ram","_api","_sound","_mouse","_in","_ticks","_palette_dac"])
  transfer_difference=m["native_transferred"]!=int(m["expected_escape"])
  if m["expected_escape"]:
   transfer_difference|=any(m["native_transfer_"+key.lower()]!=m["original"][key] for key in ("IP","SS","SP"))
  assert different or transfer_difference
  controls[str(n)]={'group':g,'name':MUTANT_NAMES[n],'rejection_kind':'state-mismatch','exit_status':exit_status,'first_rejected':r['cases'],'receipt_sha256':sha(path),'log_sha256':sha(log_path),**native[f'mutant-{n}']}
 before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
 with tempfile.TemporaryDirectory(prefix='main-repro-') as d:
  root=Path(d);(root/'tools/c_recovery').mkdir(parents=True);shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
  spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_main_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
  gen.generate(root,probe);assert (root/'tools/c_recovery/main_generated.inc').read_bytes()==(repo/'tools/c_recovery/main_generated.inc').read_bytes()
 result={'schema':'wolong-c-main-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
 'ida_database_sha256':sha(out/'ida/input.exe.i64'),'ida_probe_sha256':sha(probe),'new_routine_sha256':routines,
 'original_instruction_count':89,'original_instruction_bytes':179,'source_sha256':source,'compiled_source_manifest_sha256':compiled,
 'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_main.sh','tools/c_recovery_main_verify.py','tools/ida_main_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py','tools/c_recovery_talk_verify.py','tools/main_code_supplement.py']},
 'cases_per_optimization':case_count,'groups':GROUPS,'receipt_sha256':receipts,'negative_controls_rejected':len(MUTANTS),'mutants':controls,
 'state_mismatch_controls_rejected':len(MUTANTS),'native_builds':native,
 'entries_seen':full['O2']['entries_seen'],'font_calls_per_optimization':full['O2']['original_font_calls'],'missing_fonts':0,
 'independent_main_audits':full['O2'][AUDIT_FIELD],
 'input_assets':inputs,'tool_versions':(out/'results/tool-versions.txt').read_text(),'exact_clean_regeneration':True,
 'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
 'scope':SCOPE,
 'source_scope':DRIVER_SCOPE,'nonlocal_scope':NONLOCAL_SCOPE,'platform':PLATFORM_SCOPE,'controlled_sound_status_fixture':True,
 'c_machine_code_match':False}
 code=repo/'docs/re/c-main-code.json';sp=json.loads(code.read_text());ap=sp['assembler_verification']
 assert sp['input_sha256']==EXPECTED and sp['probe_sha256']==sha(probe) and sp['database_sha256']==sha(out/'ida/input.exe.i64')
 assert not sp['instructions'] and sp['constant_symbols']=={}
 assert set(sp['recovery_targets'])=={0x10A1C,0x11CB1,0x1EBDC} and set(sp['navigation_only_targets'])=={0x11BE0,0x1533D,0x159B7,0x189F0}
 assert ap['fixed_prior_instruction_count']==24999 and ap['fixed_prior_instruction_bytes']==56866
 assert ap['status']=='byte-exact' and ap['instructions']==ap['bytes']==0 and ap['unique_input_instructions']==ap['covered_input_instructions']==89
 assert ap['independent_reassembly']['status']=='byte-exact' and ap['independent_reassembly']['instructions']==89 and ap['independent_reassembly']['bytes']==179
 assert len(sp['instruction_coverage'])==89 and all(x['fixed_prior_coverage'] and x['independent_native_assembly_match'] for x in sp['instruction_coverage'])
 assert ap['published_manifest_sha256']==sha(repo/ap['published_manifest']) and ap['published_code_source_sha256']==sha(repo/ap['published_code_source'])
 assert ap['generator_sha256']==sha(repo/'tools/main_code_supplement.py')
 assert ap['translator_sha256']==sha(repo/'tools/assembly_rebuild.py')
 assert ap['image_id']=='sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e'
 for relative,binding in ap['coverage_inputs'].items():assert sha(repo/relative)==binding['sha256'],relative
 assert {name:binding['instruction_count'] for name,binding in ap['coverage_inputs'].items()}=={'workplace/matching-decompilation/assembly/build/source-map.json':24376,'docs/re/rectangle-handler-code.json':370,'docs/re/c-list-code.json':83,'docs/re/c-catalog-code.json':170}
 assert len(sp['original_functions'])==3 and {t['name']:t['file_sha256'] for t in sp['original_functions']}==routines
 for row in sp['instruction_coverage']:
  assert raw[row['file_start']:row['file_end']].hex()==row['bytes']
  assert row['file_start']==row['ida_linear']-0x10000+512
 relocations=sp['loader_relocation_normalization'];assert len(relocations)==1 and relocations[0]['ida_linear']==0x11CBF
 for row in relocations:
  assert row['locator_level']=='proven' and row['file_bytes']=='9a00000010' and row['ida_loaded_bytes']=='9a00000020'
  assert raw[row['file_start']:row['file_end']].hex()==row['file_bytes']
  assert len(row['loader_relocations'])==1 and row['loader_relocations'][0]['original_word']==0x1000 and row['loader_relocations'][0]['ida_word']==0x2000
 result['mz_far_call_normalization']=relocations
 result['recovery_targets']=full_probe['recovery_targets']
 result['navigation_only_targets']=full_probe['navigation_only_targets']
 for key in ('palette_dac_audits','prefix_instructions','outer_prefix_snapshots','prefix_map_sha256','prefix_map_changed_bytes','raw_map_verified_before_prefix','prefix_selected_player','nonlocal_exit_cases','normal_return_cases','dac_initialization','palette_profiles','game_palette_sha256','game_palette_size'):result[key]=full['O2'][key]
 result['independent_prefix_map_expectation']=PREFIX_MAP_EXPECTED
 result['working_palette_offset']=48
 result['working_palette_size']=len(working_palette)
 result['working_palette_sha256']=hashlib.sha256(working_palette).hexdigest()
 result['mutant_names_sha256']=sha(names_path)
 result['rng_calls_per_optimization']=full['O2']['entries_seen'].get('sub_1ECE0',0)
 result['game_clock_calls_per_optimization']=full['O2']['entries_seen'].get('sub_11D8E',0)
 result['rng_setup']=RNG_SETUP
 result['assembly_coverage_receipt_sha256']=sha(code)
 go=repo/'docs/re/c-main-go-verification.json'
 if go.exists():
  gp=json.loads(go.read_text());assert gp['vet_passed'] and gp['test_passed'] and gp['cold_test_cache']
  assert gp['original_assets_readonly'] and gp['network']=='none' and gp['tested_packages']==39 and gp['cached_test_packages']==0
  for name,digest in gp['source_sha256'].items():assert sha(repo/name)==digest,name
  result['go_verification_sha256']=sha(go)
 evidence='docs/re/125-c-nonlocal-exit-restoration.md';evidence_text=(repo/evidence).read_text();assert EXPECTED in evidence_text
 result['resolution_backlinks']=[]
 result['backlink_status']='pending-publication'
 if '**狀態：CONFORMED' in evidence_text:
  ledger=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
  assert ledger['input_sha256']==EXPECTED and len(ledger['functions'])==370 and len(ledger['code_blocks'])==48
  assert all(name in ledger['functions'] for name in NEW)
  for row in BACKLINKS:
   assert f"0x{row['ida_linear']:X}" in evidence_text,row
   older=(repo/row['older_document']).read_text();assert all(marker in older for marker in row['required_markers']),row
   want={**row,'platform_module':'dosv/KI.EXE','input_sha256':EXPECTED,'level':'proven','evidence':evidence,
         'evidence_markers':[f"sub_{row['function_ida_linear']:X}",f"0x{row['ida_linear']:X}"]}
   assert all(marker in evidence_text for marker in want['evidence_markers']),row
   matches=[actual for actual in ledger['resolution_backlinks'] if all(actual.get(key)==want[key] for key in ('platform_module','input_sha256','ida_linear','older_document'))]
   assert len(matches)==1 and all(matches[0].get(key)==value for key,value in want.items()),row
   result['resolution_backlinks'].append(want)
  result['backlink_status']='verified'
 (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
 print(f'3 functions / 89 instructions / {case_count} whole-device cases / {len(MUTANTS)} state mismatches: PASS')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True);a=p.parse_args();verify(a.repo,a.output)
