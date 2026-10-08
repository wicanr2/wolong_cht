#!/usr/bin/env python3
"""十一個世界更新 C 函式、真實事件 writer 與來源身分的 Docker 收據核對。"""
import argparse
import hashlib
import json
from pathlib import Path

GROUPS = {'growth':560,'score':131072,'writer':960,'trust':480,'governor':768,
          'diplomat':3388,'relation':968,'disaster':1536,'storm':3072,'marker':48,'scenario-monthly':12}
FUNCTIONS = {'sub_15695':('growth',128),'sub_155A6':('score',70),'sub_12FBF':('writer',79),
             'sub_157FE':('trust',42),'sub_15715':('governor',122),'sub_1578F':('diplomat',111),
             'sub_130CB':('relation',8),'sub_13119':('relation',31),'sub_122DB':('storm',163),
             'sub_12286':('disaster',85),'sub_1237E':('marker',129)}


def verify(output, repo):
    expected='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    sh='21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest()==expected
    assert hashlib.sha256((repo/'workplace/orig/dosv/SINARIO.DAT').read_bytes()).hexdigest()==sh
    probe=json.loads((output/'ida/ida-probe.json').read_text())
    assert probe['schema']=='wolong-matching-ida-probe-v1' and probe['tool_version']=='9.4'
    assert probe['input_sha256']==probe['ida_input_sha256']==expected and probe['function_count']==739
    routines={}
    for t in probe['targets']:
        name=t['name'];_,size=FUNCTIONS[name];b=b''.join(bytes.fromhex(c['bytes']) for c in t['chunks']);off=t['ida_linear']-0xfe00
        assert len(b)==size and b==raw[off:off+size];routines[name]=hashlib.sha256(b).hexdigest()
    assert set(routines)==set(FUNCTIONS)
    receipts={};rng_fixture=None
    for level in ['O0','O2']:
        p=output/'results'/f'{level}.json';r=json.loads(p.read_text())
        assert r['schema']=='wolong-c-world-parity-v1'
        assert r['input_sha256']==expected and r['scenario_sha256']==sh and r['routine_sha256']==routines
        assert r['passed'] and r['mismatch'] is None and r['go_mismatch'] is None
        assert r['cases']==142864 and r['go_cases']==560 and r['groups']==GROUPS and r['go_groups']=={'growth':560}
        assert r['full_memory_audits']==1117 and r['original_state_sha256']==r['c_state_sha256']
        assert r['rng_fixture']['table_seed_hms']==[12,34,56] and r['rng_fixture']['rerolls']==0
        assert r['c_machine_code_match'] is False
        rng_fixture=r['rng_fixture'];receipts[level]=hashlib.sha256(p.read_bytes()).hexdigest()
    assert receipts['O0']==receipts['O2']
    for n,group,count in [(1,'growth',429),(2,'score',257),(3,'writer',25),(4,'governor',1),
                           (5,'diplomat',8),(6,'disaster',1),(7,'marker',7)]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text())
        assert r['input_sha256']==expected and r['routine_sha256']==routines
        assert not r['passed'] and r['go_mismatch'] is None and r['cases']==count and r['groups']=={group:count}
        mismatch=r['mismatch'];assert mismatch['group']==group and mismatch['case']==count-1
        assert any(mismatch['original'+s]!=mismatch['c'+s] for s in ['', '_trace','_bank','_queue'])
    before=(output/'results/golem-source-before.sha256').read_bytes()
    assert before==(output/'results/golem-source-after.sha256').read_bytes()
    sources={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/')
        assert hashlib.sha256((repo/relative).read_bytes()).hexdigest()==digest;sources[relative]=digest
    index=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
    assert index['input_sha256']==expected
    for name in FUNCTIONS:
        current=index['functions'].get(name)
        if current:
            for key in ['source','header','fixture']:assert current[key+'_sha256']==sources[current[key]]
    report={'schema':'wolong-c-world-verification-v1','status':'semantic-conformed','input_sha256':expected,
            'scenario_sha256':sh,'functions':FUNCTIONS,'routine_sha256':routines,
            'ida_database_sha256':hashlib.sha256((output/'ida/input.exe.i64').read_bytes()).hexdigest(),
            'cases_per_optimization':142864,'go_cases_per_optimization':560,'groups':GROUPS,
            'full_memory_audits_per_optimization':1117,'receipt_sha256':receipts,'source_sha256':sources,
            'negative_controls_rejected':7,'rng_fixture':rng_fixture,
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope':'Eleven world-update functions; real linked writer and RNG; four scenario inputs; '
                    'sub_1585F/sub_12BD9 plus UI/sound/redraw have RET fixtures; IF/TF=0; legal inputs',
            'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
    print('eleven world-update functions, real writer, mutations and current source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--output',type=Path,required=True);p.add_argument('--repo',type=Path,required=True)
    a=p.parse_args();verify(a.output,a.repo)
