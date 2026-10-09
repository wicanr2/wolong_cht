#!/usr/bin/env python3
"""核對軍團／物件 producer、live capacity、來源矩陣與完整裝置收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import tempfile
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW={'sub_11CC9':7,'sub_12533':112,'sub_12AF4':54,'sub_12B2A':18,'sub_12B3C':108,'sub_15D19':162}
GROUPS={'point':990,'army':280,'objects':896,'matrix':12320,'minimap':48,'scan':56,'corpus':16,'pipeline':16}
MUTANTS={1:'scan',2:'point',3:'point',4:'army',5:'objects',6:'objects',7:'objects',8:'matrix',9:'matrix',10:'matrix',11:'matrix',12:'minimap'}
ASSETS={'SINARIO.DAT':(88832,'21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'),
        'MMAP.MAP':(80716,'51b6fcaa390c80bd8a358dcacbf7d8dbb6dfeb0e048d8c3bdd86329df3401bcf'),
        'MMAP.MDL':(32768,'2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
        'MMAP.MCH':(43058,'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7')}
RAW_MAP='740708c27a89db0a7f82865be623b5732099ec97aabf399e7dbde1450c87c861'
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,out):
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
    for n,(size,digest) in ASSETS.items():
        p=repo/'workplace/orig/dosv'/n;assert p.stat().st_size==size and sha(p)==digest
    assert (out/'fixtures/MMAP.raw').stat().st_size==98304 and sha(out/'fixtures/MMAP.raw')==RAW_MAP
    mch=(repo/'workplace/orig/dosv/MMAP.MCH').read_bytes()
    patterns=[i for i in range(64) if mch[0xa000+i*4] and mch[0xa001+i*4]];assert len(patterns)==55
    p=json.loads((out/'ida/ida-probe.json').read_text());assert p['tool_version']=='9.4' and p['function_count']==739
    assert p['input_sha256']==p['ida_input_sha256']==EXPECTED
    routines={}
    for t in p['targets']:
        data=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks']);assert len(data)==NEW[t['name']]
        assert data==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
        routines[t['name']]=hashlib.sha256(data).hexdigest();assert routines[t['name']]==t['file_sha256']
    assert set(routines)==set(NEW) and sum(NEW.values())==461
    block=p['decoded_blocks'][0];assert len(p['decoded_blocks'])==1 and block['ida_linear']==0x1d51f and block['end']==0x1d5d4
    data=b''.join(bytes.fromhex(x['bytes']) for x in block['instructions']);assert data==raw[0xd51f+512:0xd5d4+512]
    code={'code_1D51F':hashlib.sha256(data).hexdigest()};assert len(data)==181 and code['code_1D51F']==block['file_sha256']
    patch=next(x for x in block['instructions'] if x['ida_linear']==0x1d5a9)
    assert patch['bytes']=='80fb05' and patch['mnemonic']=='cmp' and patch['originally_code'] is False
    assert 'db' in patch['original_ida_data_line']
    instruction_count=sum(len(c['instructions']) for t in p['targets'] for c in t['chunks'])+len(block['instructions']);assert instruction_count==267
    receipts,full={},{}
    for level in ['O0','O2']:
        path=out/'results'/f'{level}.json';r=json.loads(path.read_text())
        assert r['schema']=='wolong-c-overlay-parity-v1' and r['input_sha256']==EXPECTED
        assert r['passed'] and r['mismatch'] is None and r['cases']==14622 and r['groups']==GROUPS
        assert r['full_ram_plane_audits']==14622 and r['source_matrix_audits']==12320
        assert r['valid_metadata_patterns']==patterns and not r['c_machine_code_match']
        assert r['original_state_sha256']==r['c_state_sha256']
        assert all(r['routine_sha256'][n]==h and r['entries_seen'].get(n,0)>0 for n,h in {**routines,**code}.items())
        assert r['entries_seen']['sub_1D66A']==4 and r['entries_seen']['sub_11CC9']>=20
        receipts[level]=sha(path);full[level]=r
    assert receipts['O0']==receipts['O2']
    controls={}
    for n,g in MUTANTS.items():
        path=out/'results'/f'mutant-{n}.json';r=json.loads(path.read_text());m=r['mismatch']
        assert not r['passed'] and r['input_sha256']==EXPECTED and r['groups']=={g:r['cases']}
        assert 0<r['cases']<=GROUPS[g] and m['group']==g and m['case']==r['cases']-1
        assert any(m['original'+k]!=m['c'+k] for k in ['', '_trace','_ports','_device','_ram','_planes'])
        controls[str(n)]={'group':g,'first_rejected':r['cases'],'receipt_sha256':sha(path)}
    source={}
    for line in (out/'results/c-source.sha256').read_text().splitlines():
        digest,path=line.split(None,1);relative=path.strip().removeprefix('/repo/');assert sha(repo/relative)==digest;source[relative]=digest
    compiled=sha(out/'results/c-source.sha256');assert compiled==(out/'results/compiled-source-digest.txt').read_text().strip()
    for name in ['O0','O2',*[f'mutant-{n}' for n in MUTANTS]]:
        text=(out/'results'/f'buildinfo-{name}.txt').read_text();assert '-DKI_OVERLAY_SOURCE_DIGEST=0x'+compiled[:16] in text and 'matching_overlay' in text
        if name.startswith('mutant-'):assert '-O0' in text and '-DKI_OVERLAY_MUTATION='+name.split('-')[1] in text
    before=(out/'results/golem-source-before.sha256').read_bytes();assert before==(out/'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='overlay-repro-') as d:
        root=Path(d);(root/'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo/'tools/c_recovery_display_generate.py',root/'tools/c_recovery_display_generate.py')
        spec=importlib.util.spec_from_file_location('gen',repo/'tools/c_recovery_overlay_generate.py');gen=importlib.util.module_from_spec(spec);spec.loader.exec_module(gen)
        gen.generate(root,out/'ida/ida-probe.json');assert (root/'tools/c_recovery/overlay_generated.inc').read_bytes()==(repo/'tools/c_recovery/overlay_generated.inc').read_bytes()
    result={'schema':'wolong-c-overlay-verification-v1','status':'semantic-conformed','input_sha256':EXPECTED,
            'ida_database_sha256':sha(out/'ida/input.exe.i64'),'input_assets':ASSETS,'decoded_map_sha256':RAW_MAP,
            'new_routine_sha256':routines,'code_entry_sha256':code,'original_instruction_count':267,
            'source_sha256':source,'compiled_source_manifest_sha256':compiled,
            'verification_tools_sha256':{n:sha(repo/n) for n in ['tools/c_recovery_overlay_verify.py','tools/c_recovery_overlay.sh','tools/ida_world_overlay_probe.py','tools/c_recovery_mapcells_prepare.py','tools/rle.py']},
            'cases_per_optimization':14622,'groups':GROUPS,'source_matrix_audits':12320,'valid_metadata_patterns':patterns,
            'receipt_sha256':receipts,'negative_controls_rejected':12,'mutants':controls,'entries_seen':full['O2']['entries_seen'],
            'max_original_steps':full['O2']['max_original_steps'],'exact_clean_regeneration':True,
            'tool_versions':(out/'results/tool-versions.txt').read_text(),'oracle_revision':(out/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope':'Original army/object records to clipped slot matrices, live 3/5 capacity, transparent/protected/cache slots, old-frame phase update, minimap VGA point, four scenario overlay passes and full native C renderer pipeline. Independent MCH source matrix audit; complete RAM/plane/ABI/I/O. IF/TF=0; disjoint stack; FLAGS from pinned dosgolem model. Type 3 natural producer, normal player scene and C machine code remain outside evidence',
            'c_machine_code_match':False}
    (out/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print('6 named / 1 raw C entries; 267 instructions; 14622 whole-device cases; 12320 source matrices; 12 mutants: PASS')

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    a=p.parse_args();verify(a.repo,a.output)
