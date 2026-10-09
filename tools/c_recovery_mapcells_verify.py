#!/usr/bin/env python3
"""核對世界顯示格原版／C、完整來源、獨立素材合成與錯版收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import tempfile
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_1D46A':25,'sub_1D483':32,'sub_1D4C7':88,'sub_1D615':85,'sub_1D66A':257,
     'sub_1D76B':23,'sub_1D782':20,'sub_1D796':81,'sub_1D7E7':29,'sub_1D804':70}
GROUPS={'init':12,'copy':10,'background':768,'overlay':512,'flush':16,
        'flags':1024,'push':256,'render':8,'pipeline':16}
MUTANTS={1:'copy',2:'copy',3:'flags',4:'flush',5:'background',6:'overlay',7:'overlay',
         8:'flush',9:'flags',10:'flags',11:'push',12:'push'}
ASSETS={'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
        'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
        'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7')}
RAW_MAP='740708c27a89db0a7f82865be623b5732099ec97aabf399e7dbde1450c87c861'

def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
    for name,(size,digest) in ASSETS.items():
        p=repo/'workplace/orig/dosv'/name;assert p.stat().st_size==size and sha(p)==digest
    spec=importlib.util.spec_from_file_location('rle',repo/'tools/rle.py');rle=importlib.util.module_from_spec(spec);spec.loader.exec_module(rle)
    decoded=rle.decode_file((repo/'workplace/orig/dosv/MMAP.MAP').read_bytes())
    assert len(decoded)==98304 and hashlib.sha256(decoded).hexdigest()==RAW_MAP
    assert decoded==(out/'fixtures/MMAP.raw').read_bytes()
    source_probe=repo/'workplace/matching-decompilation/c-scene/ida/ida-probe.json'
    p=json.loads(source_probe.read_text());assert p['tool_version']=='9.4' and p['function_count']==739
    assert p['input_sha256']==p['ida_input_sha256']==EXPECTED
    routines={};instructions=0
    for t in p['targets']:
        if t['name'] not in NEW:continue
        data=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks'])
        assert len(data)==NEW[t['name']] and hashlib.sha256(data).hexdigest()==t['file_sha256']
        assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
        routines[t['name']]=t['file_sha256'];instructions+=sum(len(c['instructions']) for c in t['chunks'])
    assert set(routines)==set(NEW) and sum(NEW.values())==710 and instructions==328
    receipts,full={},{}
    for level in ['O0','O2']:
        path=out/'results'/f'{level}.json';v=json.loads(path.read_text())
        assert v['schema']=='wolong-c-mapcells-parity-v1' and v['input_sha256']==EXPECTED
        assert v['passed'] and v['mismatch'] is None and v['cases']==2622 and v['groups']==GROUPS
        assert v['full_ram_plane_audits']==2622 and v['source_composite_audits']==512
        assert v['tile_blit_audits']==256 and v['original_step_bound']==4000000
        assert all(0<n<4000000 for n in v['max_original_steps'].values())
        assert v['original_state_sha256']==v['c_state_sha256'] and not v['c_machine_code_match']
        assert all(v['routine_sha256'][n]==h and v['entries_seen'].get(n,0)>0 for n,h in routines.items())
        assert v['asset_sha256']=={'MMAP.MDL':ASSETS['MMAP.MDL'][1],'MMAP.MCH':ASSETS['MMAP.MCH'][1],'MMAP.raw':RAW_MAP}
        receipts[level]=sha(path);full[level]=v
    assert receipts['O0']==receipts['O2']
    controls={}
    for number,group in MUTANTS.items():
        path=out/'results'/f'mutant-{number}.json';v=json.loads(path.read_text());m=v['mismatch']
        assert not v['passed'] and v['input_sha256']==EXPECTED and v['groups']=={group:v['cases']}
        assert 0<v['cases']<=GROUPS[group] and m['group']==group and m['case']==v['cases']-1
        assert any(m['original'+k]!=m['c'+k] for k in ['', '_trace','_ports','_device','_ram','_planes'])
        controls[str(number)]={'group':group,'first_rejected':v['cases'],'receipt_sha256':sha(path)}
    source={}
    for line in (out/'results/c-source.sha256').read_text().splitlines():
        digest,path=line.split(None,1);relative=path.strip().removeprefix('/repo/')
        assert sha(repo/relative)==digest;source[relative]=digest
    compiled=sha(out/'results/c-source.sha256')
    assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
    for name in ['O0','O2',*[f'mutant-{n}' for n in MUTANTS]]:
        text=(out/'results'/f'buildinfo-{name}.txt').read_text()
        assert '-DKI_MAPCELLS_SOURCE_DIGEST=0x'+compiled[:16] in text and 'matching_mapcells' in text
        if name.startswith('mutant-'):assert '-O0' in text and '-DKI_MAPCELLS_MUTATION='+name.split('-')[1] in text
    before=(out/'results/golem-source-before.sha256').read_bytes()
    assert before==(out/'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='mapcells-repro-') as d:
        root=Path(d);(root/'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
        spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_mapcells_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
        gen.generate(root,source_probe)
        assert (root/'tools/c_recovery/mapcells_generated.inc').read_bytes()==(repo/'tools/c_recovery/mapcells_generated.inc').read_bytes()
    result={'schema':'wolong-c-mapcells-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
            'ida_database_sha256':sha(source_probe.parent/'input.exe.i64'),'ida_probe_sha256':sha(source_probe),
            'input_assets':ASSETS,'decoded_map_sha256':RAW_MAP,'decoded_map_bytes':98304,
            'new_routine_sha256':routines,'original_instruction_count':instructions,
            'source_sha256':source,'compiled_source_manifest_sha256':compiled,
            'verification_tools_sha256':{p:sha(repo/p) for p in ['tools/c_recovery_mapcells.sh','tools/c_recovery_mapcells_verify.py','tools/c_recovery_mapcells_prepare.py','tools/ida_scene_probe.py','tools/rle.py']},
            'cases_per_optimization':2622,'groups':GROUPS,'source_composite_audits':512,
            'tile_blit_audits':256,'max_original_steps':full['O2']['max_original_steps'],
            'original_step_bound':4000000,
            'helper_preconditions':'sub_1D7E7 requires DF=0 as established by original renderer CLD; DF=1 overwrites its own executing code. Other safe DF branches remain sampled.',
            'receipt_sha256':receipts,'negative_controls_rejected':12,'mutants':controls,
            'entries_seen':full['O2']['entries_seen'],'exact_clean_regeneration':True,
            'tool_versions':(out/'results/tool-versions.txt').read_text(),
            'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope':'Original 40x23 display records: segment init, map sampling, clipped four-overlay producer, all 256 flag values and tile IDs, cached/direct/background/masked/capital paths, four-plane blit and linked source/init/copy/push/render. Independent MMAP decoder/source composites and VGA buses; full RAM/planes/ABI/I/O. Normal player scene, army/object producer chains and original C machine code remain outside evidence',
            'c_machine_code_match':False}
    (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print('10 map-cell functions; 328 instructions; 2622 cases; 512 source composites; 12 mutants: PASS')

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    a=p.parse_args();verify(a.repo,a.output)
