#!/usr/bin/env python3
"""政治／俘虜二十一函式與完整月結規則接線的原版收據核對，Docker only。"""
import argparse
import hashlib
import json
from pathlib import Path

GROUPS={'assignment':1587,'writer':60,'power':125,'relation':5808,'captives':2048,
        'generals':1536,'frontier':352,'politics':80,'initialize':16,'scenario-monthly':12,'cooperation':8}
FUNCTIONS={'sub_1585F':('generals',58),'sub_15899':('captives',167),'sub_15940':('captives',80),
           'sub_12AD2':('assignment',34),'sub_15990':('politics',22),'sub_1301C':('writer',50),
           'sub_12BD9':('initialize',121),'sub_12C52':('frontier',141),'sub_12CDF':('frontier',91),
           'sub_12D3A':('politics',30),'sub_12FB1':('politics',14),'sub_12D58':('politics',96),
           'sub_12DB8':('politics',59),'sub_12DF3':('politics',64),'sub_130F0':('relation',26),
           'sub_1310A':('relation',15),'sub_13091':('power',58),'sub_12E33':('politics',86),
           'sub_12E89':('politics',114),'sub_12EFB':('politics',118),'sub_12F71':('politics',64)}


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
        assert r['schema']=='wolong-c-politics-parity-v1' and r['input_sha256']==expected and r['scenario_sha256']==sh
        assert r['routine_sha256']==routines and r['passed'] and r['mismatch'] is None
        assert r['cases']==11632 and r['groups']==GROUPS and r['full_memory_audits']==91
        assert r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
        assert r['rng_fixture']['table_seed_hms']==[12,34,56] and r['rng_fixture']['rerolls']==0
        rng_fixture=r['rng_fixture'];receipts[level]=hashlib.sha256(q.read_bytes()).hexdigest()
    assert receipts['O0']==receipts['O2']
    for n,group,count in [(1,'writer',4),(2,'relation',4),(3,'initialize',1),(4,'assignment',1),(5,'power',3),(6,'generals',8),(7,'frontier',1),(8,'cooperation',1)]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text());assert r['input_sha256']==expected and r['routine_sha256']==routines
        assert not r['passed'] and r['cases']==count and r['groups']=={group:count}
        m=r['mismatch'];assert m['group']==group and m['case']==count-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_bank','_queue','_scratch'])
    before=(output/'results/golem-source-before.sha256').read_bytes();assert before==(output/'results/golem-source-after.sha256').read_bytes()
    sources={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert hashlib.sha256((repo/relative).read_bytes()).hexdigest()==digest;sources[relative]=digest
    idx=json.loads((repo/'docs/re/c-recovery-status.json').read_text());assert idx['input_sha256']==expected
    for name in FUNCTIONS:
        row=idx['functions'].get(name)
        if row:
            for key in ['source','header','fixture']:assert row[key+'_sha256']==sources[row[key]]
    result={'schema':'wolong-c-politics-verification-v1','status':'semantic-conformed','input_sha256':expected,
            'scenario_sha256':sh,'functions':FUNCTIONS,'routine_sha256':routines,
            'ida_database_sha256':hashlib.sha256((output/'ida/input.exe.i64').read_bytes()).hexdigest(),
            'cases_per_optimization':11632,'groups':GROUPS,'full_memory_audits_per_optimization':91,
            'receipt_sha256':receipts,'source_sha256':sources,'negative_controls_rejected':8,'rng_fixture':rng_fixture,
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope':'21 local functions and full C monthly rules on four scenarios; UI/sound/redraw RET fixtures; '
                    'valid indices/pointers/frontier lists and at least one faction; IF/TF=0; no complete Go monthly claim',
            'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print('21 politics/prisoner functions, complete monthly rules, mutations and source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--output',type=Path,required=True);p.add_argument('--repo',type=Path,required=True);a=p.parse_args();verify(a.output,a.repo)
