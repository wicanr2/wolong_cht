#!/usr/bin/env python3
"""據點結算 C 收據、完整經濟接線、負對照與來源／IDA 身分核對；Docker only。"""
import argparse
import hashlib
import json
from pathlib import Path

GROUPS = {'income': 589824, 'recruit': 1179648, 'player-tax': 3232,
          'ai-gate': 6144, 'city-settlement': 440, 'scenario-monthly': 12}
GO_GROUPS = {**GROUPS, 'income': 196608, 'player-tax': 3030}
FUNCTIONS = {'sub_153C6': ('city-settlement', 144), 'sub_15456': ('ai-gate', 57),
             'sub_1548F': ('player-tax', 109), 'sub_15538': ('income', 15),
             'sub_15547': ('recruit', 95)}


def verify(output, repo):
    expected = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    scenario_hash = '21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == expected
    assert hashlib.sha256((repo / 'workplace/orig/dosv/SINARIO.DAT').read_bytes()).hexdigest() == scenario_hash
    probe = json.loads((output / 'ida/ida-probe.json').read_text())
    assert probe['schema'] == 'wolong-matching-ida-probe-v1'
    assert probe['tool_version'] == '9.4' and probe['function_count'] == 739
    assert probe['input_sha256'] == probe['ida_input_sha256'] == expected
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
        path = output / 'results' / f'{level}.json'
        result = json.loads(path.read_text())
        assert result['schema'] == 'wolong-c-settlement-parity-v1'
        assert result['input_sha256'] == expected and result['scenario_sha256'] == scenario_hash
        assert result['routine_sha256'] == routines
        assert result['passed'] and result['mismatch'] is None and result['go_mismatch'] is None
        assert result['cases'] == 1779300 and result['go_cases'] == 1385882
        assert result['groups'] == GROUPS and result['go_groups'] == GO_GROUPS
        assert result['full_memory_audits'] == 435
        assert result['original_state_sha256'] == result['c_state_sha256']
        assert result['rng_fixture']['table_seed_hms'] == [12, 34, 56]
        assert result['rng_fixture']['rerolls'] == 0 and result['c_machine_code_match'] is False
        receipts[level] = hashlib.sha256(path.read_bytes()).hexdigest()
    assert receipts['O0'] == receipts['O2']
    rng_fixture = result['rng_fixture']
    for n, group, count in [(1, 'income', 65539), (2, 'recruit', 65),
                            (3, 'player-tax', 2421), (4, 'ai-gate', 1),
                            (5, 'city-settlement', 1)]:
        result = json.loads((output / 'results' / f'mutant-{n}.json').read_text())
        assert result['input_sha256'] == expected and result['routine_sha256'] == routines
        assert not result['passed'] and result['go_mismatch'] is None
        assert result['cases'] == count and result['rng_fixture'] == rng_fixture
        mismatch = result['mismatch']
        assert mismatch['group'] == group and mismatch['case'] == result['cases'] - 1
        assert result['groups'] == {group: result['cases']}
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_trace', '_locals'])
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
        'schema': 'wolong-c-settlement-verification-v1', 'status': 'semantic-conformed',
        'input_sha256': expected, 'scenario_sha256': scenario_hash,
        'functions': FUNCTIONS, 'routine_sha256': routines,
        'ida_database_sha256': hashlib.sha256((output / 'ida/input.exe.i64').read_bytes()).hexdigest(),
        'cases_per_optimization': 1779300, 'go_cases_per_optimization': 1385882,
        'groups': GROUPS, 'go_groups': GO_GROUPS, 'full_memory_audits_per_optimization': 435,
        'receipt_sha256': receipts, 'source_sha256': sources, 'negative_controls_rejected': 5,
        'rng_fixture': rng_fixture,
        'oracle_revision': (output / 'results/golem-revision.txt').read_text().strip(),
        'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
        'scope': 'Five local functions and linked monthly economic calculation; four raw '
                 'scenarios; nine tail callees and redraw have RET fixtures; legal inputs; IF/TF=0',
        'c_machine_code_match': False,
    }
    (output / 'verification.json').write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
    print('five settlement functions, original scenarios, mutations and source identity: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--repo', type=Path, required=True)
    args = parser.parse_args()
    verify(args.output, args.repo)
