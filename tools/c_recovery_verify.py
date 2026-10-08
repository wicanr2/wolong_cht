#!/usr/bin/env python3
"""驗證 C 還原收據、來源身分與負對照；只能在 Docker 內執行。"""
import argparse
import hashlib
import json
from pathlib import Path


def verify(root, repo):
    expected = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    receipts = {}
    for level in ['O0', 'O2']:
        p = root/'results'/f'{level}.json'
        r = json.loads(p.read_text(encoding='utf-8'))
        assert r['schema'] == 'wolong-c-rng-parity-v1'
        assert r['input_sha256'] == expected and r['passed']
        assert r['cases'] == 263680 and r['go_cases'] == 263168
        assert r['full_memory_audits'] == 258
        assert r['original_trace_sha256'] == r['c_trace_sha256']
        assert r['original_go_state_trace_sha256'] == r['go_state_trace_sha256']
        assert r['fixed_inputs'] == [
            {'Name':'identity','Mul':1,'Add':0}, {'Name':'reverse','Mul':255,'Add':255},
            {'Name':'affine-197-13','Mul':197,'Add':13}, {'Name':'affine-137-39','Mul':137,'Add':39}]
        receipts[level] = hashlib.sha256(p.read_bytes()).hexdigest()
    assert receipts['O0'] == receipts['O2']
    for n in range(1, 5):
        r = json.loads((root/'results'/f'mutant-{n}.json').read_text(encoding='utf-8'))
        assert not r['passed'] and r['register_or_flag_mismatch'] is not None
        assert r['cases'] == (1281 if n == 4 else 1)
        if n == 4:
            assert r['register_or_flag_mismatch']['group'] == 'stack-alias'
    assert (root/'results/golem-source-before.sha256').read_bytes() == (root/'results/golem-source-after.sha256').read_bytes()
    source_hashes = {}
    for line in (root/'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1)
        relative = name.strip().removeprefix('/repo/')
        assert hashlib.sha256((repo/relative).read_bytes()).hexdigest() == digest
        source_hashes[relative] = digest
    assert hashlib.sha256((repo/'workplace/orig/dosv/KI.EXE').read_bytes()).hexdigest() == expected
    index_path = repo/'docs/re/c-recovery-status.json'
    if index_path.exists():
        index=json.loads(index_path.read_text(encoding='utf-8'))
        assert index['input_sha256']==expected
        current=index['functions']['sub_1ECE0']
        assert current['source_sha256']==source_hashes[current['source']]
        assert current['header_sha256']==source_hashes[current['header']]
    result = {'schema':'wolong-c-recovery-verification-v1', 'function':'sub_1ECE0',
              'input_sha256':expected, 'status':'semantic-conformed',
              'cases_per_optimization':263680, 'go_cases_per_optimization':263168,
              'source_sha256':source_hashes, 'receipt_sha256':receipts,
              'oracle_revision':(root/'results/golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256':hashlib.sha256((root/'results/golem-source-before.sha256').read_bytes()).hexdigest(),
              'scope':'four controlled permutation tables, 16-bit final state, IF/TF=0, no clock or normal-player-path claim',
              'c_machine_code_match':False, 'negative_controls_rejected':4}
    (root/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    print('C / dosgolem / Go receipts and source identity: PASS')


if __name__ == '__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--output',type=Path,required=True)
    p.add_argument('--repo',type=Path,required=True)
    args=p.parse_args();verify(args.output,args.repo)
