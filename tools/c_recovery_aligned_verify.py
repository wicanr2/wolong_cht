#!/usr/bin/env python3
"""核對五個新 C 函式、四平面／hit 消費鏈、來源身分與負對照。"""
import argparse
import hashlib
import json
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
ASSET='2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d'
NEW={0xf888:176,0xf938:97,0xf999:23,0xc7f4:74,0xc673:59}
GROUPS={'aligned':2496,'row':1344,'zero-height':56,'plane':512,'real-assets':96,
        'redraw':216,'buttons':6,'button-query':36,'window':192}
BACKLINKS={
    'docs/re/18-tactical-button-glyphs.md':('後續 C 控制流證據見','105-c-aligned-blit-restoration.md'),
    'docs/re/103-c-hotspot-restoration.md':('後續 C 裝置接線見','105-c-aligned-blit-restoration.md'),
    'docs/re/104-c-vga-blit-restoration.md':('後續 C 位元對齊證據見','105-c-aligned-blit-restoration.md'),
    'docs/spec/223-c-vga-blit.md':('後續 C 位元對齊證據見','224-c-aligned-blit.md'),
}

def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()

def verify(repo,output):
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes();assert hashlib.sha256(raw).hexdigest()==EXPECTED
    icons=(repo/'workplace/orig/dosv/ICONGRF.DAT').read_bytes();assert hashlib.sha256(icons).hexdigest()==ASSET
    assert len(icons)>=0x6700 and len(icons[0x2800:0x6700])==0x3f00
    p=json.loads((output/'ida/ida-probe.json').read_text())
    assert p['schema']=='wolong-matching-ida-probe-v1' and p['tool_version']=='9.4'
    assert p['input_sha256']==p['ida_input_sha256']==EXPECTED and p['function_count']==739
    assert len(p['targets'])==18 and p['target_relocations_applied']==3
    routines={};new={}
    for t in p['targets']:
        off=t['ida_linear']-0x10000
        b=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks'])
        assert b==b''.join(raw[c['start']-0x10000+512:c['end']-0x10000+512] for c in t['chunks'])
        routines[t['name']]=hashlib.sha256(b).hexdigest()
        assert routines[t['name']]==t['file_sha256']
        if off in NEW:
            assert len(b)==NEW[off] and all(not c['loader_relocations'] for c in t['chunks'])
            new[t['name']]=routines[t['name']]
    assert len(new)==5 and sum(NEW.values())==429
    for table in p['raw_tables']:
        b=bytes.fromhex(table['file_bytes']);at=table['file_offset'];assert b==raw[at:at+table['size']]
        assert hashlib.sha256(b).hexdigest()==table['sha256']
    receipts={};full={};total=sum(GROUPS.values())
    for level in ['O0','O2']:
        f=output/'results'/f'{level}.json';r=json.loads(f.read_text())
        assert r['schema']=='wolong-c-aligned-parity-v1' and r['input_sha256']==EXPECTED and r['icon_sha256']==ASSET
        assert all(r['routine_sha256'][name]==digest for name,digest in routines.items())
        assert r['passed'] and r['mismatch'] is None and r['cases']==total and r['groups']==GROUPS
        assert r['full_ram_plane_audits']==total and r['indexed_content_audits']==683
        assert r['content_top']==40 and r['content_height']==400 and r['button_queries']==36
        assert all(r['entries_seen'].get(name,0)>0 for name in new)
        assert r['entries_seen']['fixture-far']>0 and r['trace_entry_size']==90 and r['c_machine_code_match'] is False
        assert r['original_state_sha256']==r['c_state_sha256']
        receipts[level]=sha(f);full[level]=r
    assert receipts['O0']==receipts['O2']
    controls={}
    for n,group in [(1,'aligned'),(2,'aligned'),(3,'row'),(4,'aligned'),(5,'row'),(6,'aligned'),(7,'row'),(8,'buttons'),(9,'redraw'),(10,'aligned'),(11,'window')]:
        f=output/'results'/f'mutant-{n}.json';r=json.loads(f.read_text())
        assert not r['passed'] and r['input_sha256']==EXPECTED and r['icon_sha256']==ASSET
        assert 0<r['cases']<=GROUPS[group] and r['groups']=={group:r['cases']}
        m=r['mismatch'];assert m['group']==group and m['case']==r['cases']-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_ports','_device','_planes','_ram'])
        controls[str(n)]={'group':group,'first_rejected':r['cases'],'receipt_sha256':sha(f)}
    source={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/');assert sha(repo/relative)==digest;source[relative]=digest
    for name in ['tools/c_recovery/aligned.c','tools/c_recovery/aligned.h','tools/c_recovery/aligned_fixture.h','tools/c_recovery_aligned.go','tools/c_recovery_vga_bus.go']:assert name in source
    compiled=sha(output/'results/c-source.sha256');assert (output/'results/compiled-source-digest.txt').read_text().strip()==compiled
    for level in ['O0','O2']:assert '-DKI_ALIGNED_SOURCE_DIGEST=0x'+compiled[:16] in (output/'results'/f'buildinfo-{level}.txt').read_text()
    before=(output/'results/golem-source-before.sha256').read_bytes();assert before==(output/'results/golem-source-after.sha256').read_bytes()
    for path,(marker,link) in BACKLINKS.items():text=(repo/path).read_text();assert marker in text and link in text
    idx=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
    for r in full.values():assert r['routine_sha256']['sub_1E453']==idx['functions']['sub_1E453']['routine_sha256']
    for name in new:
        if name in idx['functions']:
            row=idx['functions'][name]
            for key in ['source','header','fixture']:assert row[key+'_sha256']==source[row[key]]
    counter=output/'results/word-linear-counterexample-smoke.json';r=json.loads(counter.read_text());m=r['mismatch']
    assert not r['passed'] and m['input']['SI']==65535 and m['original']==m['c'] and m['original_ram']==m['c_ram']
    assert m['original_device']==m['c_device'] and m['original_planes']!=m['c_planes']
    result={'schema':'wolong-c-aligned-verification-v1','status':'semantic-conformed',
            'input_sha256':EXPECTED,'icon_sha256':ASSET,'new_routine_sha256':new,
            'dependency_routine_sha256':routines,'ida_database_sha256':sha(output/'ida/input.exe.i64'),
            'source_sha256':source,'compiled_source_manifest_sha256':compiled,'cases_per_optimization':total,
            'groups':GROUPS,'full_ram_plane_audits_per_optimization':total,
            'indexed_content_audits_per_optimization':full['O2']['indexed_content_audits'],'button_queries':36,
            'receipt_sha256':receipts,'negative_controls_rejected':11,'mutants':controls,
            'entries_seen':full['O2']['entries_seen'],'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),'scope_backlinks':BACKLINKS,
            'word_wrap_counterexample_sha256':sha(counter),'platform_contract':full['O2']['platform'],
            'content_top':40,'content_height':400,'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print(f'5 new C functions, {total} complete RAM/plane cases, 36 hit consumers, 11 mutations, source identity: PASS')

if __name__=='__main__':
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo',type=Path,required=True);parser.add_argument('--output',type=Path,required=True)
    args=parser.parse_args();verify(args.repo,args.output)
