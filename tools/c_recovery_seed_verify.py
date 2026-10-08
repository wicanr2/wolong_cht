#!/usr/bin/env python3
"""固定 RTC 的 C 播種收據驗證；Docker 內執行。"""
import argparse
import hashlib
import json
from pathlib import Path


def verify(output, repo):
    expected='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    hashes={}
    for level in ['O0','O2']:
        p=output/'results'/f'{level}.json';r=json.loads(p.read_text(encoding='utf-8'))
        assert r['schema']=='wolong-c-seed-parity-v1' and r['passed']
        assert r['input_sha256']==expected and r['cases']==86420 and r['go_cases']==86408
        assert r['steps']=={'3369':86420} and r['full_memory_audits']==85
        assert r['original_state_sha256']==r['c_state_sha256']
        hashes[level]=hashlib.sha256(p.read_bytes()).hexdigest()
    assert hashes['O0']==hashes['O2']
    for n,count in [(1,1),(2,3),(3,2)]:
        r=json.loads((output/'results'/f'mutant-{n}.json').read_text())
        assert not r['passed'] and r['mismatch'] is not None and r['cases']==count
    assert (output/'results/golem-source-before.sha256').read_bytes()==(output/'results/golem-source-after.sha256').read_bytes()
    source={}
    for line in (output/'results/c-source.sha256').read_text().splitlines():
        digest,name=line.split(None,1);relative=name.strip().removeprefix('/repo/')
        assert hashlib.sha256((repo/relative).read_bytes()).hexdigest()==digest
        source[relative]=digest
    assert hashlib.sha256((repo/'workplace/orig/dosv/KI.EXE').read_bytes()).hexdigest()==expected
    index_path=repo/'docs/re/c-recovery-status.json'
    if index_path.exists():
        index=json.loads(index_path.read_text(encoding='utf-8'))
        current=index['functions'].get('sub_1EC82')
        if current:
            assert index['input_sha256']==expected
            assert current['source_sha256']==source[current['source']]
            assert current['header_sha256']==source[current['header']]
    result={'schema':'wolong-c-seed-verification-v1','function':'sub_1EC82','status':'semantic-conformed',
            'input_sha256':expected,'cases_per_optimization':86420,'go_cases_per_optimization':86408,
            'receipt_sha256':hashes,'source_sha256':source,'negative_controls_rejected':3,
            'oracle_revision':(output/'results/golem-revision.txt').read_text().strip(),
            'oracle_compiled_tree_sha256':hashlib.sha256((output/'results/golem-source-before.sha256').read_bytes()).hexdigest(),
            'scope':'fixed successful RTC reply, all legal BCD times, IF/TF=0, disjoint stack; XOR AF matches dosgolem model',
            'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    print('seed C / original / Go receipt and source identity: PASS')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--output',type=Path,required=True);p.add_argument('--repo',type=Path,required=True)
    a=p.parse_args();verify(a.output,a.repo)
