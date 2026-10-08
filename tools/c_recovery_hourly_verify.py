#!/usr/bin/env python3
"""每時更新八函式、原始分派表與 source identity 收據核對；Docker only。"""
import argparse
import hashlib
import json
from pathlib import Path

GROUPS={'expense':100,'maintenance':125,'dispatch':14336,'diplomat':6400,'trust':30,'world':616,'scenario-hourly':264}
FUNCTIONS={'sub_13E11':('world',84),'sub_13E65':('maintenance',41),'sub_13E8E':('diplomat',111),
           'sub_15673':('expense',34),'sub_131AE':('dispatch',68),'sub_13496':('trust',16),
           'sub_13507':('trust',19),'sub_13DC9':('trust',72)}


def verify(output,repo):
    expected='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    sh='21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==expected
    assert hashlib.sha256((repo/'workplace/orig/dosv/SINARIO.DAT').read_bytes()).hexdigest()==sh
    p=json.loads((output/'ida/ida-probe.json').read_text());assert p['schema']=='wolong-matching-ida-probe-v1' and p['tool_version']=='9.4'
    assert p['input_sha256']==p['ida_input_sha256']==expected and p['function_count']==739
    routines={}
    for t in p['targets']:
        name=t['name'];_,size=FUNCTIONS[name];b=b''.join(bytes.fromhex(c['bytes']) for c in t['chunks']);off=t['ida_linear']-0xfe00
        assert len(b)==size and b==raw[off:off+size];routines[name]=hashlib.sha256(b).hexdigest()
    assert set(routines)==set(FUNCTIONS)
    receipts={};rng_fixture=None
    for level in ['O0','O2']:
        q=output/'results'/f'{level}.json';r=json.loads(q.read_text())
        assert r['schema']=='wolong-c-hourly-parity-v1' and r['input_sha256']==expected and r['scenario_sha256']==sh
        assert r['routine_sha256']==routines and r['passed'] and r['mismatch'] is None
        assert r['cases']==21871 and r['groups']==GROUPS and r['full_memory_audits']==171
        assert r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
        assert r['rng_fixture']['table_seed_hms']==[12,34,56] and r['rng_fixture']['rerolls']==0
        receipts[level]=hashlib.sha256(q.read_bytes()).hexdigest();rng_fixture=r['rng_fixture']
    assert receipts['O0']==receipts['O2']
    for n,group,count in [(1,'expense',14),(2,'maintenance',46),(3,'dispatch',57),(4,'dispatch',57),(5,'trust',2),(6,'world',1),(7,'diplomat',231)]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text());assert r['input_sha256']==expected and r['routine_sha256']==routines
        assert not r['passed'] and r['cases']==count and r['groups']=={group:count}
        m=r['mismatch'];assert m['group']==group and m['case']==count-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_bank','_queue','_scratch','_globals','_cadence'])
    before=(output/'results/golem-source-before.sha256').read_bytes();assert before==(output/'results/golem-source-after.sha256').read_bytes()
    sources={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert hashlib.sha256((repo/relative).read_bytes()).hexdigest()==digest;sources[relative]=digest
    idx=json.loads((repo/'docs/re/c-recovery-status.json').read_text());assert idx['input_sha256']==expected
    for name in FUNCTIONS:
        row=idx['functions'].get(name)
        if row:
            for key in ['source','header','fixture']:assert row[key+'_sha256']==sources[row[key]]
    result={'schema':'wolong-c-hourly-verification-v1','status':'semantic-conformed','input_sha256':expected,
            'scenario_sha256':sh,'functions':FUNCTIONS,'routine_sha256':routines,
            'ida_database_sha256':hashlib.sha256((output/'ida/input.exe.i64').read_bytes()).hexdigest(),
            'cases_per_optimization':21871,'groups':GROUPS,'full_memory_audits_per_optimization':171,
            'receipt_sha256':receipts,'source_sha256':sources,'negative_controls_rejected':7,'rng_fixture':rng_fixture,
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope':'Eight hourly/dispatch functions with real maintenance/diplomat/RNG; other 11 event '
                    'handlers and UI/sound/exit/redraw RET fixtures; legal indices/code.low 0..13; IF/TF=0',
            'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print('8 hourly functions, dispatch table, mutations and current source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--output',type=Path,required=True);p.add_argument('--repo',type=Path,required=True);a=p.parse_args();verify(a.output,a.repo)
