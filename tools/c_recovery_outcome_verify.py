#!/usr/bin/env python3
"""核對自動戰鬥與易主21函式的原始來源、裝置對拍與錯版反例。"""
import argparse
import hashlib
import importlib.util
import json
import shlex
import shutil
import struct
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PROBE_SHA = 'b20f688289b06ca5384f3ac810cd4b9b93cc5ddb82d5be4d12556701a8988ce9'
DATABASE_SHA = '6efd571f0d44cf6d52e3d1eeecdc41079a082cfbf15a798a8fb1b4dae656d282'
NEW = {'sub_10CE7': 9, 'sub_14236': 51, 'sub_1474A': 113, 'sub_14CF3': 112, 'sub_14D63': 65, 'sub_14DA4': 76, 'sub_14DF0': 108, 'sub_14FCE': 166, 'sub_15074': 64, 'sub_150B4': 35, 'sub_15130': 131, 'sub_151B3': 210, 'sub_15285': 82, 'sub_152D7': 102, 'sub_15CA4': 34, 'sub_15CE0': 57, 'sub_17028': 20, 'sub_188CC': 62, 'sub_1890A': 83, 'sub_195C9': 118, 'sub_1963F': 23}
TARGETS = {86661, 100618, 85412, 89252, 94248, 86320, 86451, 86196, 82486, 103999, 103881, 83786, 100556, 85966, 86743, 89312, 85347, 68839, 85488, 85235, 86132}
GROUPS = {'battle':5376,'retreat':384,'capture':1376,'fall':96,'diplomat':32,'history':8,'army':320,'neighbors':512,'minimap':256,'alert':4}
MUTANTS = {1: ('battle', 'first-side-coefficient'), 2: ('battle', 'morale-divisor'), 3: ('battle', 'aptitude-shift'), 4: ('battle', 'general-rng-mask'), 5: ('battle', 'battle-ratio-offset'), 6: ('battle', 'battle-ratio-limit'), 7: ('battle', 'loser-base-loss'), 8: ('battle', 'leader-minimum'), 9: ('retreat', 'retreat-men-threshold'), 10: ('retreat', 'own-city-stay'), 11: ('capture', 'old-city-count'), 12: ('capture', 'new-city-count'), 13: ('capture', 'governor-duty'), 14: ('capture', 'failed-capital-faction-mask'), 15: ('capture', 'old-corps-redirection'), 16: ('neighbors', 'neighbor-city-count'), 17: ('fall', 'surviving-faction-count'), 18: ('diplomat', 'diplomat-duty'), 19: ('history', 'city-owner-history'), 20: ('army', 'army-occupancy')}
NATIVE_TAGS = 'matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga'


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def native_build(binary, info_path, compiled, mutation, optimize):
    """核對ELF內嵌Go建置資料與另存輸出，不執行受測程式。"""
    data=binary.read_bytes()
    assert data[:6]==b'\x7fELF\x02\x01',binary
    shoff=struct.unpack_from('<Q',data,40)[0]
    entsize,count,names_index=struct.unpack_from('<HHH',data,58)
    assert entsize>=64 and 0<names_index<count
    sections=[struct.unpack_from('<IIQQQQIIQQ',data,shoff+i*entsize) for i in range(count)]
    ns=sections[names_index];names=data[ns[4]:ns[4]+ns[5]]
    entries=[s for s in sections if names[s[0]:].split(b'\0',1)[0]==b'.go.buildinfo']
    assert len(entries)==1
    s=entries[0];blob=data[s[4]:s[4]+s[5]]
    assert blob[:14]==b'\xff Go buildinf:' and blob[14:16]==b'\x08\x02'
    def string_at(at):
        length=shift=0
        while True:
            assert at<len(blob) and shift<64
            value=blob[at];at+=1;length|=(value&127)<<shift
            if value<128:break
            shift+=7
        assert at+length<=len(blob)
        return blob[at:at+length],at+length
    version,at=string_at(32);module,_=string_at(at)
    assert version==b'go1.26.7' and len(module)>=33 and module[-17]==10
    lines=module[16:-16].decode().splitlines()
    assert 'path\tgithub.com/wicanr2/dosgolem/wolongcoutcome' in lines
    saved=info_path.read_text().splitlines()
    assert saved[0].endswith(': go1.26.7') and Path(saved[0].split(':',1)[0]).name==binary.name
    def flags_from(lines):
        prefix='build\tCGO_CFLAGS='
        values=[x.removeprefix('\t')[len(prefix):] for x in lines if x.removeprefix('\t').startswith(prefix)]
        assert len(values)==1
        quoted=shlex.split(values[0]);assert len(quoted)==1
        return shlex.split(quoted[0])
    flags=flags_from(lines);assert flags==flags_from(saved[1:])
    digest='-DKI_OUTCOME_SOURCE_DIGEST=0x'+compiled[:16]
    assert [v for v in flags if v.startswith('-DKI_OUTCOME_SOURCE_DIGEST=')]==[digest]
    assert [v for v in flags if v.startswith('-DKI_OUTCOME_MUTATION=')]==([] if mutation is None else [f'-DKI_OUTCOME_MUTATION={mutation}'])
    assert '-'+optimize in flags and 'build\t-tags='+NATIVE_TAGS in lines
    return {'c_binary_sha256':sha(binary),'buildinfo_sha256':sha(info_path),'go_version':version.decode(),'cgo_flags':flags}


def verify(repo,out):
    if GROUPS is None:
        raise RuntimeError('Reviewed outcome case matrix has not been frozen')
    count=sum(GROUPS.values())
    assert all(n>0 for n in GROUPS.values()) and len(MUTANTS)==20
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest()==EXPECTED
    probe=out/'ida/ida-probe.json';p=json.loads(probe.read_text())
    assert sha(probe)==PROBE_SHA and sha(out/'ida/input.exe.i64')==DATABASE_SHA
    assert p['tool_version']=='9.4' and p['input_sha256']==p['ida_input_sha256']==EXPECTED
    assert set(p['recovery_targets'])==TARGETS and set(p['navigation_only_targets'])=={0x14A7B,0x14ADE,0x125A3}
    routines={};rows=[]
    for target in p['targets']:
        if target['ida_linear'] not in TARGETS:continue
        data=b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in target['chunks'])
        assert len(data)==NEW[target['name']] and hashlib.sha256(data).hexdigest()==target['file_sha256']
        assert not any(c['loader_relocations'] for c in target['chunks'])
        routines[target['name']]=target['file_sha256']
        rows.extend(x for c in target['chunks'] for x in c['instructions'])
    assert set(routines)==set(NEW) and len(rows)==796 and sum(NEW.values())==1721
    assert p['unresolved_call_boundaries']==[]
    source={}
    for line in (out/'results/c-source.sha256').read_text().splitlines():
        digest,path=line.split(None,1);path=path.strip().removeprefix('/repo/')
        assert sha(repo/path)==digest,path
        source[path]=digest
    assert {'tools/c_recovery/outcome.c','tools/c_recovery/outcome.h',
      'tools/c_recovery/outcome_fixture.h','tools/c_recovery/outcome_generated.inc',
      'tools/c_recovery_outcome.go','tools/c_recovery_outcome_data.go'}.issubset(source)
    compiled=sha(out/'results/c-source.sha256')
    assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
    manifest={'path':'workplace/matching-decompilation/c-outcome/results/c-source.sha256','sha256':compiled}
    reports={};receipts={};native={}
    for level in ('O0','O2'):
        path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
        assert r['schema']=='wolong-c-outcome-parity-v1' and r['passed'] and r['mismatch'] is None
        assert r['input_sha256']==EXPECTED and r['cases']==count and r['groups']==GROUPS
        assert r['full_ram_plane_audits']==r['indexed_content_audits']==r['palette_dac_audits']==count
        assert r['independent_outcome_audits']>=count and r['stage_audits']=={'outcome':count}
        assert r['new_routine_sha256']==routines and r['source_manifest']==manifest
        assert r['new_probe_sha256']==PROBE_SHA
        assert all(r['entries_seen'].get(name,0)>0 for name in NEW)
        assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
        assert r['independent_raw_initialization'] and not r['cross_machine_snapshot_initialization']
        assert r['fixed_raw_rng_state'] and r['actual_nonlocal_exit']
        assert r['nonlocal_exit_cases']==144 and r['normal_return_cases']==8220
        assert r['normal_return_cases']+r['nonlocal_exit_cases']==count
        assert r['entries_seen'].get('sub_11CB1',0)>=r['nonlocal_exit_cases']
        assert r['rng_initial_state_sha256']
        for seed,digest in r['rng_initial_state_sha256'].items():
            value=int(seed,16)
            initial=bytes([value^0xA5,value])+bytes((i*73+i//2+value)&255 for i in range(256))
            assert len(initial)==258 and hashlib.sha256(initial).hexdigest()==digest
        assert r['original_missing_fonts']==r['c_missing_fonts']==0
        assert r['original_font_misses']==r['c_font_misses']=={}
        assert r['original_font_calls']==r['c_font_calls']
        reports[level]=r;receipts[level]=sha(path)
        native[level]=native_build(out/f'outcome-{level}',out/'results'/f'buildinfo-{level}.txt',compiled,None,level)
    assert receipts['O0']==receipts['O2']
    names={}
    for line in (out/'results/mutant-names.tsv').read_text().splitlines():
        n,group,name=line.split('\t');assert int(n) not in names
        names[int(n)]=(group,name)
    assert names==MUTANTS
    controls={}
    for n,(group,name) in MUTANTS.items():
        path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch']
        assert (out/'results'/f'mutant-{n}.exit').read_text().strip()=='1'
        assert r['schema']=='wolong-c-outcome-parity-v1' and not r['passed']
        assert r['input_sha256']==EXPECTED and r['source_manifest']==manifest
        assert r['groups']=={group:r['cases']} and 0<r['cases']<=GROUPS[group]
        assert m['group']==group and m['case']==r['cases']-1
        assert any(m['original'+key]!=m['c'+key] for key in ('','_ram','_planes','_trace','_ports','_in','_api','_palette_dac'))
        build=native_build(out/f'mutant-{n}',out/'results'/f'buildinfo-mutant-{n}.txt',compiled,n,'O0')
        native[f'mutant-{n}']=build
        controls[str(n)]={'group':group,'name':name,'rejection_kind':'state-mismatch','exit_status':1,
          'first_rejected':r['cases'],'receipt_sha256':sha(path),'log_sha256':sha(out/'results'/f'mutant-{n}.log'),**build}
    before=(out/'results/golem-source-before.sha256').read_bytes()
    assert before==(out/'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='outcome-repro-') as temporary:
        root=Path(temporary);(root/'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
        spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_outcome_generate.py')
        gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen);gen.generate(root,probe)
        assert (root/'tools/c_recovery/outcome_generated.inc').read_bytes()==(repo/'tools/c_recovery/outcome_generated.inc').read_bytes()
    code=repo/'docs/re/c-outcome-code.json';coverage=json.loads(code.read_text());assembly=coverage['assembler_verification']
    assert coverage['input_sha256']==EXPECTED and coverage['probe_sha256']==PROBE_SHA
    assert assembly['unique_input_instructions']==assembly['covered_input_instructions']==796
    assert assembly['independent_reassembly']['instructions']==796 and assembly['independent_reassembly']['bytes']==1721
    assert assembly['independent_reassembly']['status']=='byte-exact'
    assert len(coverage['instruction_coverage'])==796
    assert all(x['fixed_prior_coverage'] and x['independent_native_assembly_match'] for x in coverage['instruction_coverage'])
    assert assembly['published_manifest_sha256']==sha(repo/assembly['published_manifest'])
    assert assembly['published_code_source_sha256']==sha(repo/assembly['published_code_source'])
    result={'schema':'wolong-c-outcome-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
      'ida_database_sha256':DATABASE_SHA,'ida_probe_sha256':PROBE_SHA,'new_routine_sha256':routines,
      'original_instruction_count':796,'original_instruction_bytes':1721,'source_sha256':source,
      'compiled_source_manifest_sha256':compiled,'native_builds':native,'cases_per_optimization':count,
      'whole_device_stage_audits_per_optimization':count,'groups':GROUPS,'receipt_sha256':receipts,
      'negative_controls_rejected':20,'state_mismatch_controls_rejected':20,'mutants':controls,
      'entries_seen':reports['O2']['entries_seen'],'independent_outcome_audits':reports['O2']['independent_outcome_audits'],
      'assembly_coverage_receipt_sha256':sha(code),'exact_clean_regeneration':True,
      'tool_versions':(out/'results/tool-versions.txt').read_text(),'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
      'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),'c_machine_code_match':False,
      'scope':reports['O2']['scope'],'platform':reports['O2']['platform'],
      'rng_initial_state_sha256':reports['O2']['rng_initial_state_sha256'],
      'rng_initialization':reports['O2']['rng_initialization'],
      'nonlocal_exit_cases':reports['O2']['nonlocal_exit_cases'],'normal_return_cases':reports['O2']['normal_return_cases'],
      'verification_tools_sha256':{name:sha(repo/name) for name in ('tools/c_recovery_outcome.sh','tools/c_recovery_outcome_verify.py','tools/c_recovery_outcome_generate.py','tools/ida_outcome_probe.py','tools/outcome_code_supplement.py')}}
    go=repo/'docs/re/c-outcome-go-verification.json'
    if go.exists():
        proof=json.loads(go.read_text())
        assert proof['vet_passed'] and proof['test_passed'] and proof['cold_test_cache']
        assert proof['tested_packages']==39 and proof['cached_test_packages']==0
        assert proof['original_assets_readonly'] and proof['network']=='none'
        for name,digest in proof['source_sha256'].items():assert sha(repo/name)==digest
        result['go_verification_sha256']=sha(go)
    evidence=repo/'docs/re/130-c-battle-outcome-restoration.md'
    result['backlink_status']='pending-publication'
    if '**狀態：CONFORMED' in evidence.read_text():
        assert 'sub_1474A' in evidence.read_text() and '0x1474A' in evidence.read_text()
        ledger=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
        assert ledger['input_sha256']==EXPECTED and len(ledger['functions'])==431 and len(ledger['code_blocks'])==50
        assert all(name in ledger['functions'] for name in NEW) and 'sub_11BE0' not in ledger['functions']
        for name in NEW:
            entry=ledger['functions'][name]
            assert entry['routine_sha256']==routines[name] and not entry['c_machine_code_match']
        for address,document in [(0x1474A,'docs/spec/46-post-battle-retreat.md'),
                                 (0x14CF3,'docs/spec/47-city-fall-corps-redirect.md'),
                                 (0x196CF,'docs/re/129-c-route-restoration.md')]:
            text=(repo/document).read_text()
            assert '戰後閉包已由原生 C 驗證' in text and '130-c-battle-outcome-restoration.md' in text
            matches=[x for x in ledger['resolution_backlinks'] if x['older_document']==document and x['ida_linear']==address]
            assert len(matches)==1
            row=matches[0]
            assert row['evidence']=='docs/re/130-c-battle-outcome-restoration.md' and row['input_sha256']==EXPECTED
            assert row['platform_module']=='dosv/KI.EXE' and row['function_ida_linear']==address
            assert row['level']=='proven' and all(x in text for x in row['required_markers'])
            assert all(marker in evidence.read_text() for marker in row['evidence_markers'])
        result['backlink_status']='verified'
    (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    print(f'21 original functions / 796 instructions / {count} full-device cases / 20 state mismatches: PASS')


if __name__=='__main__':
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo',type=Path,required=True);parser.add_argument('--output',type=Path,required=True)
    args=parser.parse_args();verify(args.repo,args.output)
