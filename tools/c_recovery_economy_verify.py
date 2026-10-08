#!/usr/bin/env python3
"""月結 C 原版／Go 收據、六組負對照與 IDA／來源身分核對；Docker 內執行。"""
import argparse
import hashlib
import json
from pathlib import Path

GROUPS = {'credit': 393440, 'debit': 393440, 'reserve': 393216,
          'distance': 393216, 'deficit': 196608, 'monthly': 5648}
FUNCTIONS = {'sub_15358': ('monthly', 110), 'sub_15609': ('credit', 34),
             'sub_1563B': ('debit', 40), 'sub_154FC': ('distance', 54),
             'sub_155EC': ('reserve', 13), 'sub_15828': ('deficit', 55)}


def verify(output, repo):
    expected = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == expected
    probe = json.loads((output / 'ida/ida-probe.json').read_text())
    assert probe['schema'] == 'wolong-matching-ida-probe-v1'
    assert probe['input_sha256'] == probe['ida_input_sha256'] == expected
    assert probe['tool_version'] == '9.4' and probe['function_count'] == 739
    routines = {}
    for target in probe['targets']:
        name = target['name']
        _, size = FUNCTIONS[name]
        data = b''.join(bytes.fromhex(chunk['bytes']) for chunk in target['chunks'])
        offset = target['ida_linear'] - 0xfe00
        assert len(data) == size and data == raw[offset:offset + size]
        routines[name] = hashlib.sha256(data).hexdigest()
    assert set(routines) == set(FUNCTIONS)
    receipts = {}
    for level in ['O0', 'O2']:
        p = output / 'results' / f'{level}.json'
        result = json.loads(p.read_text())
        assert result['schema'] == 'wolong-c-economy-parity-v1'
        assert result['passed'] and result['mismatch'] is None and result['go_mismatch'] is None
        assert result['input_sha256'] == expected and result['routine_sha256'] == routines
        assert result['groups'] == GROUPS and result['cases'] == 1775568
        assert result['go_cases'] == 1326175 and result['full_memory_audits'] == 867
        assert result['original_state_sha256'] == result['c_state_sha256']
        assert result['c_machine_code_match'] is False
        receipts[level] = hashlib.sha256(p.read_bytes()).hexdigest()
    assert receipts['O0'] == receipts['O2']
    for n, group, count in [(1, 'credit', 262145), (2, 'debit', 1),
                            (3, 'reserve', 131072), (4, 'distance', 321),
                            (5, 'deficit', 98308), (6, 'monthly', 129)]:
        result = json.loads((output / 'results' / f'mutant-{n}.json').read_text())
        assert result['input_sha256'] == expected and result['routine_sha256'] == routines
        assert not result['passed'] and result['go_mismatch'] is None
        assert result['cases'] == count and result['groups'] == {group: count}
        mismatch = result['mismatch']
        assert mismatch['case'] == count - 1 and mismatch['group'] == group
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_state', '_trace'])
    before = (output / 'results/golem-source-before.sha256').read_bytes()
    assert before == (output / 'results/golem-source-after.sha256').read_bytes()
    sources = {}
    for line in (output / 'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1)
        relative = name.strip().removeprefix('/repo/')
        assert hashlib.sha256((repo / relative).read_bytes()).hexdigest() == digest
        sources[relative] = digest
    index = json.loads((repo / 'docs/re/c-recovery-status.json').read_text())
    assert index['input_sha256'] == expected
    for name in FUNCTIONS:
        current = index['functions'].get(name)
        if current:
            for key in ['source', 'header', 'fixture']:
                assert current[key + '_sha256'] == sources[current[key]]
    report = {
        'schema': 'wolong-c-economy-verification-v1', 'status': 'semantic-conformed',
        'input_sha256': expected, 'functions': FUNCTIONS, 'routine_sha256': routines,
        'ida_database_sha256': hashlib.sha256((output / 'ida/input.exe.i64').read_bytes()).hexdigest(),
        'cases_per_optimization': 1775568, 'go_cases_per_optimization': 1326175,
        'groups': GROUPS, 'full_memory_audits_per_optimization': 867,
        'receipt_sha256': receipts, 'source_sha256': sources, 'negative_controls_rejected': 6,
        'oracle_revision': (output / 'results/golem-revision.txt').read_text().strip(),
        'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
        'rng_fixture': {
            'table_seed_hms': [12, 34, 56],
            'original_and_c_counter': 'uint8(case_index * 19)',
            'original_and_c_state_index': 'uint8(case_index * 31)',
            'setup': 'Write identical 258-byte state to CS:ECFC before each near call',
            'go_setup': 'rng.FromRaw with the identical pre-call 258-byte state',
            'rerolls': 0,
            'normal_game_seed_modified': False,
        },
        'scope': 'Six local functions, linked money/deficit/RNG; other monthly callees '
                 'use explicit fixtures; IF/TF=0; disjoint stack; AF follows dosgolem model',
        'c_machine_code_match': False,
    }
    (output / 'verification.json').write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
    print('six-function economy receipts, mutations, IDA and current source identity: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--repo', type=Path, required=True)
    args = parser.parse_args()
    verify(args.output, args.repo)
