#!/usr/bin/env python3
"""外交／金額視窗、MZ 重定位、十七函式來源及負對照核對。"""
import argparse
import hashlib
import json
import struct
from pathlib import Path

GROUPS={'window':1250,'helper':2304,'selector':192,'reply':200,'diplomacy':40320,'trust-rng':3072,'amount':2016,'wrapper':2240,'scenario-modal':1728}
SIZES={0x2078:94,0x20d6:123,0x38c7:31,0x38e6:28,0x3902:230,0x39e8:288,0x3c3d:92,0x3b7e:43,0x3d09:60,0x3d45:35,0x3c99:39,0x3cdc:45,0x9321:21,0x87ff:17,0x3d68:41,0x1d46:72,0x2216:21}
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
SCENARIO='21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'
BACKLINKS={
    'docs/spec/219-c-event-handlers.md':('後續 C 視窗證據見','220-c-modal-control.md'),
    'docs/re/100-c-event-handlers-restoration.md':('後續 C 視窗證據見','101-c-modal-restoration.md'),
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
    assert p['mz_header_size']==512 and p['ida_load_paragraph']==0x1000
    reloc_table,reloc_count=struct.unpack_from('<H',raw,24)[0],struct.unpack_from('<H',raw,6)[0]
    relocations=[512+off+seg*16 for off,seg in (struct.unpack_from('<HH',raw,reloc_table+i*4) for i in range(reloc_count))]
    assert reloc_count==116 and p['relocation_file_offsets']==relocations and p['target_relocations_applied']==20
    routines={};applied=0
    for t in p['targets']:
        off=t['ida_linear']-0x10000;assert off in SIZES
        file_chunks=[]
        for c in t['chunks']:
            at=c['start']-0xfe00;file=bytes.fromhex(c['file_bytes']);assert file==raw[at:at+len(file)]
            expected=bytearray(file);wanted=[x for x in relocations if at<=x<at+len(file)]
            assert [r['file_offset'] for r in c['loader_relocations']]==wanted
            for row in c['loader_relocations']:
                original=struct.unpack_from('<H',raw,row['file_offset'])[0]
                assert row['original_word']==original and row['ida_word']==(original+0x1000)&65535
                struct.pack_into('<H',expected,row['file_offset']-at,row['ida_word']);applied+=1
            assert bytes(expected)==bytes.fromhex(c['bytes']);file_chunks.append(file)
        b=b''.join(file_chunks);assert len(b)==SIZES[off] and t['name']==f'sub_{0x10000+off:X}'
        routines[t['name']]=hashlib.sha256(b).hexdigest();assert routines[t['name']]==t['file_sha256']
    assert len(routines)==17 and sum(SIZES.values())==1280 and applied==20
    receipts={};passed={};total=sum(GROUPS.values())
    for level in ['O0','O2']:
        q=output/'results'/f'{level}.json';r=json.loads(q.read_text())
        assert r['schema']=='wolong-c-modal-parity-v1'
        assert r['runtime_load_paragraph']==0x110 and r['far_runtime_segment']==0x1110 and r['trace_entry_size']==252
        assert r['input_sha256']==EXPECTED and r['scenario_sha256']==SCENARIO
        assert r['routine_sha256']==routines and r['passed'] and r['mismatch'] is None
        assert r['groups']==GROUPS and r['cases']==total and r['full_memory_audits']==total//128+1
        assert set(r['entries_seen'])==set(routines) and all(v>0 for v in r['entries_seen'].values())
        assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
        assert r['rng_fixture']['table_seed_hms']==[12,34,56] and r['rng_fixture']['rerolls']==0
        receipts[level]=sha(q);passed[level]=r
    assert receipts['O0']==receipts['O2']
    mutants={}
    for n,group in [(1,"window"),(2,"trust-rng"),(3,"amount"),(4,"amount"),(5,"amount"),(6,"selector"),(7,"helper"),(8,"reply"),(9,"helper"),(10,"window")]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text())
        assert r['input_sha256']==EXPECTED and r['routine_sha256']==routines and not r['passed']
        assert 0<r['cases']<=GROUPS[group] and r['groups']=={group:r['cases']}
        m=r['mismatch'];assert m['group']==group and m['case']==r['cases']-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_bank','_queue','_scratch','_globals','_cadence','_ui'])
        mutants[str(n)]={'group':group,'first_rejected':r['cases'],'receipt_sha256':sha(output/'results'/f'mutant-{n}.json')}
    before=(output/'results/golem-source-before.sha256').read_bytes()
    assert before==(output/'results/golem-source-after.sha256').read_bytes()
    sources={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/')
        assert sha(repo/relative)==digest;sources[relative]=digest
    for name in ['modal.c','modal.h','modal_fixture.h','events.c','hourly.c','politics.c','world_update.c','settlement.c','economy.c','rng.c']:
        assert 'tools/c_recovery/'+name in sources
    assert 'tools/c_recovery_modal.go' in sources
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
    result={'schema':'wolong-c-modal-verification-v1','status':'semantic-conformed',
            'input_sha256':EXPECTED,'scenario_sha256':SCENARIO,'routine_sha256':routines,
            'ida_database_sha256':sha(output/'ida/input.exe.i64'),'source_sha256':sources,
            'cases_per_optimization':total,'groups':GROUPS,'full_memory_audits_per_optimization':total//128+1,
            'receipt_sha256':receipts,'entries_seen':passed['O2']['entries_seen'],
            'negative_controls_rejected':10,'mutants':mutants,'rng_fixture':passed['O2']['rng_fixture'],
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'scope_backlinks':BACKLINKS,'runtime_load_paragraph':passed['O2']['runtime_load_paragraph'],'far_runtime_segment':passed['O2']['far_runtime_segment'],
            'target_relocations_applied':20,'c_machine_code_match':False,'scope':passed['O2']['scope']}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print(f'17 modal functions; {total} cases per optimization; 10 mutations; current source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    a=p.parse_args();verify(a.repo,a.output)
