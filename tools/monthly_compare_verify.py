#!/usr/bin/env python3
"""月結三方比較收據核對；原版資料留本機，Docker 內執行。"""
import argparse
import hashlib
import json
from pathlib import Path

EXE = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
SCENARIO = '21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'


def same_state(original, go, original_rng, go_rng):
    return original == go and original_rng == go_rng


def selftest():
    block = bytes(22208)
    rng = bytes(258)
    assert same_state(block, block, rng, rng)
    for offset in [0x12, 0x80, 0x8c0, 0x42df, 0x52c0]:
        changed = bytearray(block)
        changed[offset] = 1
        assert not same_state(block, changed, rng, rng)
    changed = bytearray(rng)
    changed[0] = 1
    assert not same_state(block, block, rng, changed)
    print('whole block, globals, faction, city, general, queue and RNG rejection selftest: PASS')


def verify(output, repo):
    selftest()
    assert hashlib.sha256((repo/'workplace/orig/dosv/KI.EXE').read_bytes()).hexdigest() == EXE
    assert hashlib.sha256((repo/'workplace/orig/dosv/SINARIO.DAT').read_bytes()).hexdigest() == SCENARIO
    r = json.loads((output/'audit.json').read_text())
    assert r['schema'] == 'wolong-c-go-monthly-audit-v1'
    assert r['input_sha256'] == EXE and r['scenario_sha256'] == SCENARIO
    assert r['cases'] == r['original_c_cases'] == r['go_passed'] == 36 and r['go_diverged'] == 0
    expected = {f's{s}-p{p}-seed{seed}' for s in range(1,5) for p in [0,7,21] for seed in [0,77,255]}
    assert {v['name'] for v in r['vectors']} == expected
    original_digest, go_digest = hashlib.sha256(), hashlib.sha256()
    files = {}
    for v in r['vectors']:
        assert not v['precondition_diffs'] and not v['postcondition_diffs']
        assert v['original_go_equal'] and v['rng_equal'] and v['settled']
        assert v['go_runtime']['strategic_ai'] and not v['go_runtime']['approximate_event10']
        values = {}
        for kind in ['input.block','original.block','go.block','rng-before.bin','rng-original.bin','rng-go.bin']:
            p = output/'vectors'/f"{v['name']}.{kind}"
            values[kind] = p.read_bytes()
            assert len(values[kind]) == (22208 if kind.endswith('block') else 258)
            files[str(p.relative_to(output))] = hashlib.sha256(values[kind]).hexdigest()
        assert files[f"vectors/{v['name']}.input.block"] == v['input_sha256']
        assert same_state(values['original.block'],values['go.block'],values['rng-original.bin'],values['rng-go.bin'])
        assert values['rng-original.bin'].hex() == v['original_rng'] == v['go_rng']
        assert values['rng-before.bin'][0] == v['counter'] and values['rng-before.bin'][1] == v['state_index']
        original_digest.update(values['original.block']);original_digest.update(values['rng-original.bin'])
        go_digest.update(values['go.block']);go_digest.update(values['rng-go.bin'])
    assert original_digest.hexdigest() == r['original_state_sha256'] == r['go_state_sha256'] == go_digest.hexdigest()
    source = {}
    for file in ['go-source.sha256','c-source.sha256']:
        for line in (output/file).read_text().splitlines():
            digest, name = line.split(None,1);relative = name.strip().removeprefix('/repo/')
            assert hashlib.sha256((repo/relative).read_bytes()).hexdigest() == digest
            source[relative] = digest
    before = (output/'golem-source-before.sha256').read_bytes()
    assert before == (output/'golem-source-after.sha256').read_bytes()
    baseline_path = output/'before-score/audit.json'
    if baseline_path.exists():
        baseline = json.loads(baseline_path.read_text())
        assert baseline['original_c_cases'] == 36 and baseline['go_passed'] == 0 and baseline['go_diverged'] == 36
    result = {'schema':'wolong-c-go-monthly-verification-v1','status':'conformed-local-rules',
              'input_sha256':EXE,'scenario_sha256':SCENARIO,'cases':36,'original_c_cases':36,'go_passed':36,
              'whole_original_block_compared':True,'rng_compared':True,'comparison_masks':[],
              'audit_sha256':hashlib.sha256((output/'audit.json').read_bytes()).hexdigest(),
              'vector_sha256':files,'source_sha256':source,'rejection_selftests':6,
              'oracle_revision':(output/'golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256':hashlib.sha256(before).hexdigest(),
              'scope':r['scope'],'c_machine_code_match':False}
    (output/'verification.json').write_text(json.dumps(result,ensure_ascii=False,indent=2)+'\n')
    print('36 original/C/Go complete block and RNG comparisons with current source identity: PASS')


if __name__ == '__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--output',type=Path,required=True);p.add_argument('--repo',type=Path,required=True)
    a=p.parse_args();verify(a.output,a.repo)
