#!/usr/bin/env python3
"""核對數字 C、固定 IDA、實際編譯來源與獨立 VGA 對照收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW = {'sub_1062F': 107, 'sub_1069A': 68, 'sub_106DE': 23, 'sub_11E17': 47}
GROUPS = {'number': 312, 'colors': 1280, 'digit': 2816, 'blank': 256,
          'primitive-boundary': 192, 'date-caller': 4, 'parameter-caller': 21}
MUTANTS = {1: 'number', 2: 'number', 3: 'number', 4: 'digit', 5: 'digit',
           6: 'primitive-boundary', 7: 'digit', 8: 'colors'}
BACKLINKS = {
    'docs/re/28-text-number-rendering.md': ('後續 C 數字 raster', '109-c-number-raster-restoration.md'),
    'docs/re/102-c-numeric-editor-restoration.md': ('後續 C 數字 raster', '109-c-number-raster-restoration.md'),
    'docs/spec/221-c-numeric-editor.md': ('後續 C 數字 raster', '229-c-number-raster.md'),
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify(repo, out):
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == EXPECTED
    probe = json.loads((out / 'ida/ida-probe.json').read_text())
    assert probe['tool_version'] == '9.4'
    assert probe['input_sha256'] == probe['ida_input_sha256'] == EXPECTED
    assert probe['function_count'] == 739
    assert len(probe['targets']) == 4 and len(probe['decoded_blocks']) == 2
    routines = {}
    for target in probe['targets']:
        data = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert data == b''.join(raw[c['start'] - 0x10000 + 512:c['end'] - 0x10000 + 512]
                               for c in target['chunks'])
        digest = hashlib.sha256(data).hexdigest()
        assert digest == target['file_sha256']
        routines[target['name']] = digest
        if target['name'] in NEW:
            assert len(data) == NEW[target['name']]
    for block in probe['decoded_blocks']:
        data = b''.join(bytes.fromhex(r['bytes']) for r in block['instructions'])
        assert data == raw[block['ida_linear'] - 0x10000 + 512:block['end'] - 0x10000 + 512]
        name = 'code_10984' if block['ida_linear'] == 0x10984 else 'sub_11E17'
        routines[name] = hashlib.sha256(data).hexdigest()
    date = next(f for f in probe['inventory'] if f['ida_linear'] == 0x11E17)
    assert date['name'] == 'sub_11E17' and date['chunks'] == [[0x11E17, 0x11E46]]
    assert not any(f['ida_linear'] == 0x10984 for f in probe['inventory'])
    receipts, full = {}, {}
    for level in ['O0', 'O2']:
        path = out / 'results' / f'{level}.json'
        result = json.loads(path.read_text())
        assert result['schema'] == 'wolong-c-numbers-parity-v1'
        assert result['input_sha256'] == EXPECTED and result['passed'] and result['mismatch'] is None
        assert result['groups'] == GROUPS and result['cases'] == sum(GROUPS.values())
        assert result['full_ram_plane_audits'] == result['indexed_content_audits'] == result['cases']
        assert result['original_state_sha256'] == result['c_state_sha256']
        assert result['c_machine_code_match'] is False
        assert all(result['routine_sha256'][n] == h and result['entries_seen'].get(n, 0) > 0
                   for n, h in routines.items())
        assert result['original_font_calls'] == result['c_font_calls'] == [0, 0]
        receipts[level], full[level] = sha(path), result
    assert receipts['O0'] == receipts['O2']
    controls = {}
    for number, group in MUTANTS.items():
        path = out / 'results' / f'mutant-{number}.json'
        result = json.loads(path.read_text())
        assert not result['passed'] and result['input_sha256'] == EXPECTED
        assert result['groups'] == {group: result['cases']}
        assert 0 < result['cases'] <= GROUPS[group]
        mismatch = result['mismatch']
        assert mismatch['group'] == group and mismatch['case'] == result['cases'] - 1
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_trace', '_ports', '_device', '_planes', '_ram'])
        controls[str(number)] = {'group': group, 'first_rejected': result['cases'], 'receipt_sha256': sha(path)}
    source = {}
    for line in (out / 'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1)
        relative = name.strip().removeprefix('/repo/')
        assert sha(repo / relative) == digest
        source[relative] = digest
    for name in ['tools/c_recovery/numbers.c', 'tools/c_recovery/numbers.h',
                 'tools/c_recovery/numbers_fixture.h', 'tools/c_recovery/numbers_generated.inc',
                 'tools/c_recovery_numbers.go', 'tools/c_recovery_numbers_generate.py',
                 'tools/c_recovery_vga_bus.go', 'tools/c_recovery_display_generate.py']:
        assert name in source
    compiled = sha(out / 'results/c-source.sha256')
    assert (out / 'results/compiled-source-digest.txt').read_text().strip() == compiled
    for level in ['O0', 'O2']:
        build = (out / 'results' / f'buildinfo-{level}.txt').read_text()
        assert '-DKI_NUMBERS_SOURCE_DIGEST=0x' + compiled[:16] in build
        assert '-' + level + ' ' in build and 'matching_numbers' in build
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='numbers-repro-') as directory:
        root = Path(directory)
        (root / 'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo / 'tools/c_recovery_display_generate.py', root / 'tools/c_recovery_display_generate.py')
        spec = importlib.util.spec_from_file_location('numbers_generate', repo / 'tools/c_recovery_numbers_generate.py')
        generator = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(generator)
        generator.generate(root, out / 'ida/ida-probe.json')
        assert (root / 'tools/c_recovery/numbers_generated.inc').read_bytes() == (repo / 'tools/c_recovery/numbers_generated.inc').read_bytes()
    regenerated = {'exact_clean_regeneration': True,
                   'generated_source_sha256': sha(repo / 'tools/c_recovery/numbers_generated.inc'),
                   'ida_probe_sha256': sha(out / 'ida/ida-probe.json')}
    (out / 'results/generated-repro.json').write_text(json.dumps(regenerated, indent=2) + '\n')
    for path, (marker, link) in BACKLINKS.items():
        text = (repo / path).read_text()
        assert marker in text and link in text
    result = {'schema': 'wolong-c-numbers-verification-v1', 'status': 'semantic-conformed',
              'input_sha256': EXPECTED, 'ida_database_sha256': sha(out / 'ida/input.exe.i64'),
              'new_routine_sha256': {n: routines[n] for n in NEW},
              'code_entry_sha256': {'0x10984': routines['code_10984']},
              'source_sha256': source, 'compiled_source_manifest_sha256': compiled,
              'cases_per_optimization': sum(GROUPS.values()), 'groups': GROUPS,
              'full_ram_plane_audits_per_optimization': sum(GROUPS.values()),
              'indexed_content_audits_per_optimization': sum(GROUPS.values()),
              'receipt_sha256': receipts, 'entries_seen': full['O2']['entries_seen'],
              'negative_controls_rejected': len(controls), 'mutants': controls,
              'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
              'clean_regeneration': regenerated,
              'scope_backlinks': BACKLINKS,
              'scope': 'Real ICONGRF segment3 digits; native C; independent VGA/RAM; IF/TF=0; legal DIV quotient; primitive DF/CH and two original callers; no normal-player or C machine-code claim',
              'c_machine_code_match': False}
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n')
    print('4 named C routines; 1 raw code entry; 4881 full RAM/plane/ABI cases per optimization; 8 mutants: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    verify(args.repo, args.output)
