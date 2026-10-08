#!/usr/bin/env python3
"""核對 C 時鐘的原版、Go、來源身分及負對照收據；Docker 內執行。"""
import argparse
import hashlib
import json
from pathlib import Path


def verify(output, repo):
    expected = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    routine = '4e2b6ecf618756da0f51e5dfb29c34b502fc70ffbf2d7050c5f5f51cb1a311dd'
    groups = {'legal-date': 226665, 'wait-and-callback': 96, 'year-boundary': 65536}
    receipts = {}
    for level in ['O0', 'O2']:
        path = output / 'results' / f'{level}.json'
        result = json.loads(path.read_text(encoding='utf-8'))
        assert result['schema'] == 'wolong-c-clock-parity-v1'
        assert result['passed'] and result['mismatch'] is None
        assert result['input_sha256'] == expected and result['routine_sha256'] == routine
        assert result['cases'] == 292297 and result['go_cases'] == 292249
        assert result['groups'] == groups and result['full_memory_audits'] == 286
        assert result['original_trace_sha256'] == result['c_trace_sha256']
        assert result['c_machine_code_match'] is False
        receipts[level] = hashlib.sha256(path.read_bytes()).hexdigest()
    assert receipts['O0'] == receipts['O2']
    for mutation, count in [(1, 207), (2, 151110), (3, 9), (4, 226668)]:
        result = json.loads((output / 'results' / f'mutant-{mutation}.json').read_text())
        assert result['input_sha256'] == expected and result['routine_sha256'] == routine
        assert not result['passed'] and result['cases'] == count
        assert result['mismatch']['case'] == count - 1
        assert result['go_cases'] == count - 1
        mismatch = result['mismatch']
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_date', '_polls', '_trace'])
    before = (output / 'results/golem-source-before.sha256').read_bytes()
    assert before == (output / 'results/golem-source-after.sha256').read_bytes()
    source = {}
    for line in (output / 'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1)
        relative = name.strip().removeprefix('/repo/')
        assert hashlib.sha256((repo / relative).read_bytes()).hexdigest() == digest
        source[relative] = digest
    binary = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(binary).hexdigest() == expected
    assert hashlib.sha256(binary[0x1f8e:0x2017]).hexdigest() == routine
    index = json.loads((repo / 'docs/re/c-recovery-status.json').read_text())
    assert index['input_sha256'] == expected
    current = index['functions'].get('sub_11D8E')
    if current:
        assert current['source_sha256'] == source[current['source']]
        assert current['header_sha256'] == source[current['header']]
        assert current['fixture_sha256'] == source[current['fixture']]
    result = {
        'schema': 'wolong-c-clock-verification-v1', 'function': 'sub_11D8E',
        'status': 'semantic-conformed', 'input_sha256': expected,
        'routine_sha256': routine, 'cases_per_optimization': 292297,
        'go_cases_per_optimization': 292249, 'groups': groups,
        'full_memory_audits_per_optimization': 286,
        'receipt_sha256': receipts, 'source_sha256': source,
        'negative_controls_rejected': 4,
        'oracle_revision': (output / 'results/golem-revision.txt').read_text().strip(),
        'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
        'scope': 'Legal dates and u16 year rollover; explicit callee/poll fixtures; '
                 'IF/TF=0; disjoint stack; no full callee, hardware timing or player path claim',
        'c_machine_code_match': False,
    }
    (output / 'verification.json').write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print('clock C / original / Go receipts, mutations and source identity: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--repo', type=Path, required=True)
    args = parser.parse_args()
    verify(args.output, args.repo)
