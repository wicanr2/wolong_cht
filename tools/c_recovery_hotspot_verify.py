#!/usr/bin/env python3
"""十函式熱區／query／memory ABI 與真實 pixel 輸入核對。"""
import argparse
import hashlib
import json
import struct
from pathlib import Path

GROUPS={'query-pixels':256000,'query-word':240,'init':5,'rectangle':432,'register-cells':4000,'wrapper':450,'numeric-pixels':80,'finance-pixels':64,'scenario-hotspot':96}
GO_GROUPS={'numeric-pixels':80}
GO_FUNCTIONS={}
SIZES={0xe3c0:23,0xe3d7:68,0xe41b:56,0xe453:38,0x895d:71,0x89de:18,0xd5d4:65,0xc14:76,0xc60:23,0xc77:53}
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
SCENARIO='21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'
BACKLINKS={
    'docs/re/22-strategy-command-tree.md':('原始 word 定址勘誤','103-c-hotspot-restoration.md'),
    'docs/re/47-main-screen-window-registry.md':('原始 word 定址勘誤','103-c-hotspot-restoration.md'),
    'docs/re/102-c-numeric-editor-restoration.md':('熱區 callee 勘誤','103-c-hotspot-restoration.md'),
    'docs/spec/78-amount-input-editor.md':('熱區與 glyph 勘誤','103-c-hotspot-restoration.md'),
    'docs/spec/221-c-numeric-editor.md':('後續 C 熱區證據見','222-c-hotspot-map.md'),
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
    assert reloc_count==116 and p['relocation_file_offsets']==relocations and p['target_relocations_applied']==0
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
    assert len(routines)==10 and sum(SIZES.values())==491 and applied==0
    receipts={};passed={};total=sum(GROUPS.values())
    for level in ['O0','O2']:
        q=output/'results'/f'{level}.json';r=json.loads(q.read_text())
        assert r['schema']=='wolong-c-hotspot-parity-v1'
        assert r['runtime_load_paragraph']==0x110 and r['far_runtime_segment']==0x1110 and r['trace_entry_size']==292
        assert r['input_sha256']==EXPECTED and r['scenario_sha256']==SCENARIO
        assert r['routine_sha256']==routines and r['passed'] and r['mismatch'] is None
        assert r['go_groups']==GO_GROUPS and r['go_cases']==sum(GO_GROUPS.values())
        assert r['groups']==GROUPS and r['cases']==total and r['full_memory_audits']==total//128+1
        assert set(r['entries_seen'])==set(routines) and all(v>0 for v in r['entries_seen'].values())
        assert r['original_state_sha256']==r['c_state_sha256'] and not r['c_machine_code_match']
        assert r['rng_fixture']['table_seed_hms']==[12,34,56] and r['rng_fixture']['rerolls']==0
        receipts[level]=sha(q);passed[level]=r
    assert receipts['O0']==receipts['O2']
    mutants={}
    for n,group in [(1,'init'),(2,'rectangle'),(3,'rectangle'),(4,'rectangle'),(5,'query-pixels'),(6,'wrapper'),(7,'wrapper'),(8,'wrapper')]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text())
        assert r['input_sha256']==EXPECTED and r['routine_sha256']==routines and not r['passed']
        assert 0<r['cases']<=GROUPS[group] and r['groups']=={group:r['cases']}
        m=r['mismatch'];assert m['group']==group and m['case']==r['cases']-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_bank','_queue','_scratch','_globals','_cadence','_ui','_map','_flags_map'])
        mutants[str(n)]={'group':group,'first_rejected':r['cases'],'receipt_sha256':sha(output/'results'/f'mutant-{n}.json')}
    before=(output/'results/golem-source-before.sha256').read_bytes()
    assert before==(output/'results/golem-source-after.sha256').read_bytes()
    sources={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/')
        assert sha(repo/relative)==digest;sources[relative]=digest
    for name in ['hotspot.c','hotspot.h','hotspot_fixture.h','numeric.c','modal.c','events.c','hourly.c','politics.c','world_update.c','settlement.c','economy.c','rng.c']:
        assert 'tools/c_recovery/'+name in sources
    assert 'tools/c_recovery_hotspot.go' in sources
    for path,(marker,link) in BACKLINKS.items():
        text=(repo/path).read_text()
        assert marker in text and link in text
    assert '(y & 0xFFF8)' in (repo/'docs/re/22-strategy-command-tree.md').read_text()
    assert '(bx & 0xFFF8)' in (repo/'docs/re/47-main-screen-window-registry.md').read_text()
    go_sources={}
    for line in (output/'results/go-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert sha(repo/relative)==digest;go_sources[relative]=digest
    assert 'internal/state/diplomacy_event.go' in go_sources
    digest=hashlib.sha256((output/'results/c-source.sha256').read_bytes()).hexdigest()
    assert (output/'results/compiled-source-digest.txt').read_text().strip()==digest
    for level in ['O0','O2']:
        assert '-DKI_HOTSPOT_SOURCE_DIGEST=0x'+digest[:16] in (output/'results'/f'buildinfo-{level}.txt').read_text()
    idx=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
    assert idx['input_sha256']==EXPECTED
    for name in routines:
        row=idx['functions'].get(name)
        if row:
            assert row['go_cases_per_optimization']==0
            for key in ['source','header','fixture']:
                assert row[key+'_sha256']==sources[row[key]]
    result={'schema':'wolong-c-hotspot-verification-v1','status':'semantic-conformed',
            'input_sha256':EXPECTED,'scenario_sha256':SCENARIO,'routine_sha256':routines,
            'ida_database_sha256':sha(output/'ida/input.exe.i64'),'source_sha256':sources,'go_source_sha256':go_sources,
            'cases_per_optimization':total,'groups':GROUPS,'full_memory_audits_per_optimization':total//128+1,
            'receipt_sha256':receipts,'entries_seen':passed['O2']['entries_seen'],
            'negative_controls_rejected':8,'mutants':mutants,'rng_fixture':passed['O2']['rng_fixture'],
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'go_cases_per_optimization':sum(GO_GROUPS.values()),'go_groups':GO_GROUPS,'go_function_cases':GO_FUNCTIONS,'scope_backlinks':BACKLINKS,'runtime_load_paragraph':passed['O2']['runtime_load_paragraph'],'far_runtime_segment':passed['O2']['far_runtime_segment'],
            'target_relocations_applied':0,'compiled_source_manifest_sha256':digest,'input_terminal_cancel':True,'c_machine_code_match':False,'scope':passed['O2']['scope']}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print(f'10 hotspot functions; {total} cases per optimization; 8 mutations; current source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    a=p.parse_args();verify(a.repo,a.output)
