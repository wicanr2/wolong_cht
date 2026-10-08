#!/usr/bin/env python3
"""八個VGA函式、獨立裝置、全部plane／latch／RAM與source digest核對。"""
import argparse
import hashlib
import json
from pathlib import Path

EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
ASSET='b9c7745e3ed9b32f0c12003fe81f82756136f738427dc421ec96f6c4ed5c4ac8'
SIZES={0x9796:45,0x97c3:45,0xf9b0:107,0xfa1b:28,0xfa37:107,0xfaa2:32,0xfac2:79,0xfb11:24}
GROUPS={'draw':2800,'leaf':5880,'save':200,'wrapper':1000,'real-assets':8,'roundtrip':54}
BACKLINKS={
    'docs/re/103-c-hotspot-restoration.md':('後續 C VGA 證據見','104-c-vga-blit-restoration.md'),
    'docs/spec/222-c-hotspot-map.md':('後續 C VGA 證據見','223-c-vga-blit.md'),
    'docs/re/03-image-blitter.md':('後續 C 控制流證據見','104-c-vga-blit-restoration.md'),
}

def sha(p):
    return hashlib.sha256(p.read_bytes()).hexdigest()

def verify(repo,output):
    raw=(repo/'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest()==EXPECTED
    portraits=(repo/'workplace/orig/dosv/KAOGRF.DAT').read_bytes()
    assert len(portraits)==150*2048 and hashlib.sha256(portraits).hexdigest()==ASSET
    p=json.loads((output/'ida/ida-probe.json').read_text())
    assert p['schema']=='wolong-matching-ida-probe-v1' and p['tool_version']=='9.4'
    assert p['input_sha256']==p['ida_input_sha256']==EXPECTED and p['function_count']==739
    assert p['target_relocations_applied']==0
    routines={}
    for t in p['targets']:
        off=t['ida_linear']-0x10000;assert off in SIZES
        b=b''.join(bytes.fromhex(c['file_bytes']) for c in t['chunks'])
        assert b==b''.join(bytes.fromhex(c['bytes']) for c in t['chunks'])
        assert len(b)==SIZES[off] and b==raw[off+512:off+512+len(b)]
        assert t['name']==f'sub_{off+0x10000:X}'
        routines[t['name']]=hashlib.sha256(b).hexdigest()
    assert len(routines)==8 and sum(SIZES.values())==467
    receipts={};passed={};total=sum(GROUPS.values())
    for level in ['O0','O2']:
        f=output/'results'/f'{level}.json';r=json.loads(f.read_text())
        assert r['schema']=='wolong-c-vga-parity-v1'
        assert r['input_sha256']==EXPECTED and r['portrait_sha256']==ASSET
        assert r['routine_sha256']==routines and r['passed'] and r['mismatch'] is None
        assert r['cases']==total and r['groups']==GROUPS and r['full_ram_plane_audits']==total
        assert r['indexed_pixel_audits']==total//32 and r['trace_entry_size']==90
        assert set(r['entries_seen'])==set(routines) and all(v>0 for v in r['entries_seen'].values())
        assert r['original_state_sha256']==r['c_state_sha256'] and r['c_machine_code_match'] is False
        receipts[level]=sha(f);passed[level]=r
    assert receipts['O0']==receipts['O2']
    controls={}
    for n,group in [(1,'leaf'),(2,'draw'),(3,'leaf'),(4,'draw'),(5,'draw'),(6,'save'),(7,'wrapper'),(8,'draw')]:
        f=output/'results'/f'mutant-{n}.json';r=json.loads(f.read_text())
        assert r['input_sha256']==EXPECTED and r['routine_sha256']==routines and not r['passed']
        assert 0<r['cases']<=GROUPS[group] and r['groups']=={group:r['cases']}
        m=r['mismatch'];assert m['group']==group and m['case']==r['cases']-1
        assert any(m['original'+s]!=m['c'+s] for s in ['', '_trace','_ports','_device','_planes','_ram'])
        controls[str(n)]={'group':group,'first_rejected':r['cases'],'receipt_sha256':sha(f)}
    source={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/')
        assert sha(repo/relative)==digest;source[relative]=digest
    for name in ['tools/c_recovery/vga.c','tools/c_recovery/vga.h','tools/c_recovery/vga_fixture.h','tools/c_recovery_vga.go','tools/c_recovery_vga_bus.go']:
        assert name in source
    compiled=sha(output/'results/c-source.sha256')
    assert (output/'results/compiled-source-digest.txt').read_text().strip()==compiled
    for level in ['O0','O2']:
        assert '-DKI_VGA_SOURCE_DIGEST=0x'+compiled[:16] in (output/'results'/f'buildinfo-{level}.txt').read_text()
    before=(output/'results/golem-source-before.sha256').read_bytes()
    assert before==(output/'results/golem-source-after.sha256').read_bytes()
    for path,(marker,link) in BACKLINKS.items():
        text=(repo/path).read_text();assert marker in text and link in text
    idx=json.loads((repo/'docs/re/c-recovery-status.json').read_text())
    for name in routines:
        row=idx['functions'].get(name)
        if row:
            for key in ['source','header','fixture']:
                assert row[key+'_sha256']==source[row[key]]
    images={}
    for index in [0,5,49,149]:
        for dest in [0,1280]:
            a=output/'results'/f'asset-{index}-{dest}-original.png'
            b=output/'results'/f'asset-{index}-{dest}-c.png'
            assert a.read_bytes()==b.read_bytes()
            assert a.read_bytes().startswith(b'\x89PNG\r\n\x1a\n')
            images[a.name]=sha(a);images[b.name]=sha(b)
    result={'schema':'wolong-c-vga-verification-v1','status':'semantic-conformed',
            'input_sha256':EXPECTED,'portrait_sha256':ASSET,'routine_sha256':routines,
            'ida_database_sha256':sha(output/'ida/input.exe.i64'),'source_sha256':source,
            'compiled_source_manifest_sha256':compiled,'cases_per_optimization':total,
            'groups':GROUPS,'full_ram_plane_audits_per_optimization':total,
            'indexed_pixel_audits_per_optimization':total//32,'receipt_sha256':receipts,
            'negative_controls_rejected':8,'mutants':controls,'entries_seen':passed['O2']['entries_seen'],
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
            'local_debug_palette_png_sha256':images,'scope_backlinks':BACKLINKS,
            'platform_contract':passed['O2']['platform'],'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print(f'8 VGA functions; {total} complete RAM/plane cases; independent latch/port states; 8 mutations; source identity: PASS')

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--repo',type=Path,required=True);p.add_argument('--output',type=Path,required=True)
    args=p.parse_args();verify(args.repo,args.output)
