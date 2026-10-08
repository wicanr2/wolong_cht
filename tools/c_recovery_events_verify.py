#!/usr/bin/env python3
"""十三碼事件、三十函式來源、實際入口覆蓋與八個負對照收據核對。"""
import argparse
import hashlib
import json
from pathlib import Path

GROUPS={'presence':5632,'relation':4620,'target':240,'official':2560,
        'diplomacy':3072,'disaster':136,'disaster-event':17408,'capital':2112,
        'cleanup':640,'handlers':5184,'scenario-dispatch':1872}
SIZES={0x320c:20,0x3220:66,0x3262:71,0x32a9:64,0x32e9:62,0x3327:97,
       0x3388:98,0x33ea:19,0x33fd:136,0x3485:17,0x34a6:11,0x34b1:86,
       0x351a:12,0x3526:133,0x35ab:66,0x35ed:76,0x3639:48,0x3669:46,
       0x3697:45,0x36c4:78,0x3712:95,0x3771:103,0x37d8:29,0x37f5:59,
       0x3138:55,0x4502:70,0x6a3d:94,0x23ff:57,0x2438:33,0x50d7:73}
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
SCENARIO='21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'
BACKLINKS={
    'docs/spec/218-c-hourly-update.md':('後續 C handler 證據見','219-c-event-handlers.md'),
    'docs/re/99-c-hourly-update-restoration.md':('後續 C handler 證據見','100-c-event-handlers-restoration.md'),
    'docs/spec/64-capital-relocation-report.md':('後續 C 控制流證據見','100-c-event-handlers-restoration.md'),
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify(repo,output):
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest()==EXPECTED
    assert sha(repo/'workplace/orig/dosv/SINARIO.DAT')==SCENARIO
    p=json.loads((output/'ida/ida-probe.json').read_text())
    assert p['schema']=='wolong-matching-ida-probe-v1' and p['tool_version']=='9.4'
    assert p['input_sha256']==p['ida_input_sha256']==EXPECTED and p['function_count']==739
    routines={}
    for t in p['targets']:
        off=t['ida_linear']-0x10000;assert off in SIZES
        b=b''.join(bytes.fromhex(c['bytes']) for c in t['chunks'])
        assert len(b)==SIZES[off] and b==raw[off+512:off+512+len(b)]
        assert t['name']==f'sub_{0x10000+off:X}'
        routines[t['name']]=hashlib.sha256(b).hexdigest()
    assert len(routines)==30 and sum(SIZES.values())==1919
    receipts={};passed={};total=sum(GROUPS.values())
    for level in ['O0','O2']:
        q=output/'results'/f'{level}.json';r=json.loads(q.read_text())
        assert r['schema']=='wolong-c-events-parity-v1'
        assert r['input_sha256']==EXPECTED and r['scenario_sha256']==SCENARIO
        assert r['routine_sha256']==routines and r['passed'] and r['mismatch'] is None
        assert r['groups']==GROUPS and r['cases']==total and r['full_memory_audits']==total//128+1
        assert set(r['entries_seen'])==set(routines) and all(v>0 for v in r['entries_seen'].values())
        assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
        assert r['rng_fixture']['table_seed_hms']==[12,34,56] and r['rng_fixture']['rerolls']==0
        receipts[level]=sha(q);passed[level]=r
    assert receipts['O0']==receipts['O2']
    mutants={}
    for n,group in [(1,'presence'),(2,'relation'),(3,'target'),(4,'disaster'),(5,'capital'),(6,'disaster'),(7,'cleanup'),(8,'diplomacy')]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text())
        assert r['input_sha256']==EXPECTED and r['routine_sha256']==routines and not r['passed']
        assert 0<r['cases']<=GROUPS[group] and r['groups']=={group:r['cases']}
        m=r['mismatch'];assert m['group']==group and m['case']==r['cases']-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_bank','_queue','_scratch','_globals','_cadence'])
        mutants[str(n)]={'group':group,'first_rejected':r['cases'],'receipt_sha256':sha(output/'results'/f'mutant-{n}.json')}
    before=(output/'results/golem-source-before.sha256').read_bytes()
    assert before==(output/'results/golem-source-after.sha256').read_bytes()
    sources={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/')
        assert sha(repo/relative)==digest;sources[relative]=digest
    for name in ['events.c','events.h','events_fixture.h','hourly.c','politics.c','world_update.c','settlement.c','economy.c','rng.c']:
        assert 'tools/c_recovery/'+name in sources
    assert 'tools/c_recovery_events.go' in sources
    for path,(marker,link) in BACKLINKS.items():
        text=(repo/path).read_text()
        assert marker in text and link in text
    idx=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
    assert idx['input_sha256']==EXPECTED
    for name in routines:
        row=idx['functions'].get(name)
        if row:
            for key in ['source','header','fixture']:
                assert row[key+'_sha256']==sources[row[key]]
    result={'schema':'wolong-c-events-verification-v1','status':'semantic-conformed',
            'input_sha256':EXPECTED,'scenario_sha256':SCENARIO,'routine_sha256':routines,
            'ida_database_sha256':sha(output/'ida/input.exe.i64'),'source_sha256':sources,
            'cases_per_optimization':total,'groups':GROUPS,'full_memory_audits_per_optimization':total//128+1,
            'receipt_sha256':receipts,'entries_seen':passed['O2']['entries_seen'],
            'negative_controls_rejected':8,'mutants':mutants,'rng_fixture':passed['O2']['rng_fixture'],
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope_backlinks':BACKLINKS,'c_machine_code_match':False,'scope':passed['O2']['scope']}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print(f'30 event functions; {total} cases per optimization; 8 mutations; current source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    a=p.parse_args();verify(a.repo,a.output)
