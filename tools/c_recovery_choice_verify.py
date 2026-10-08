#!/usr/bin/env python3
"""核對原始 selector、live operand、進言 caller 與完整裝置收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import struct
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW = {'sub_1036F': 84, 'sub_103C3': 35, 'sub_103E6': 46, 'sub_10414': 161,
       'sub_104B5': 74, 'sub_104FF': 78, 'sub_1054D': 78, 'sub_1059B': 82,
       'nullsub_5': 1, 'sub_1061F': 16, 'sub_10B46': 105, 'sub_10BAF': 30,
       'sub_19479': 70, 'sub_194BF': 80}
GROUPS = {'selector': 48, 'bands': 8, 'scroll': 16, 'helpers': 27, 'live': 3,
          'corpus': 18, 'speech': 8, 'advise-choice': 4, 'xor-roundtrip': 6}
MUTANTS = {1: 'corpus', 2: 'helpers', 3: 'bands', 4: 'scroll', 5: 'selector', 6: 'selector',
           7: 'helpers', 8: 'helpers', 9: 'selector', 10: 'scroll', 11: 'scroll', 12: 'live'}
BACKLINKS = {'docs/re/84-popup-row-band-and-world-cursor.md': ('後續原始 C selector', '112-c-choice-selector-restoration.md'),
             'docs/re/88-mouse-cursor-visibility.md': ('後續原始 C selector', '112-c-choice-selector-restoration.md'),
             'docs/spec/45-advise-scene-layout.md': ('後續原始 C selector', '232-c-choice-selector.md')}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify(repo, out):
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == EXPECTED
    spec = importlib.util.spec_from_file_location('talk_contract', repo / 'tools/c_recovery_talk_verify.py')
    assets = importlib.util.module_from_spec(spec); spec.loader.exec_module(assets)
    for name, (size, digest) in assets.ASSETS.items():
        path = repo / 'workplace/orig/dosv' / name
        assert path.stat().st_size == size and sha(path) == digest
    probe = json.loads((out / 'ida/ida-probe.json').read_text())
    assert probe['tool_version'] == '9.4' and probe['function_count'] == 739
    assert probe['input_sha256'] == probe['ida_input_sha256'] == EXPECTED
    assert len(probe['targets']) == 14 and len(probe['decoded_blocks']) == 2
    routines, blocks = {}, {}
    for target in probe['targets']:
        data = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert len(data) == NEW[target['name']]
        assert data == b''.join(raw[c['start'] - 0x10000 + 512:c['end'] - 0x10000 + 512] for c in target['chunks'])
        routines[target['name']] = hashlib.sha256(data).hexdigest()
        assert routines[target['name']] == target['file_sha256']
    assert sum(NEW.values()) == 940
    for block in probe['decoded_blocks']:
        at, end = block['ida_linear'] - 0x10000 + 512, block['end'] - 0x10000 + 512
        data = b''.join(bytes.fromhex(r['bytes']) for r in block['instructions'])
        assert data == raw[at:end]
        blocks[f'code_{block["ida_linear"]:X}'] = hashlib.sha256(data).hexdigest()
    receipts, full = {}, {}
    for level in ['O0', 'O2']:
        path = out / 'results' / f'{level}.json'; r = json.loads(path.read_text())
        assert r['schema'] == 'wolong-c-choice-parity-v1' and r['input_sha256'] == EXPECTED
        assert r['passed'] and r['mismatch'] is None and r['cases'] == 138 and r['groups'] == GROUPS
        assert r['full_ram_plane_audits'] == r['indexed_content_audits'] == 138
        assert r['original_state_sha256'] == r['c_state_sha256'] and not r['c_machine_code_match']
        assert r['xor_roundtrips_restored'] == 3
        assert r['selected_rows']['bands'] == [0, 0, 1, 0, 1, 3, 0, 0]
        assert r['selected_rows']['live'] == [1, 1, 1]
        assert r['selected_rows']['advise-choice'] == [1, 1, 1, 1]
        assert all(r['routine_sha256'][n] == h and r['entries_seen'].get(n, 0) > 0 for n, h in {**routines, **blocks}.items())
        assert r['original_font_calls'] == r['c_font_calls'] and r['original_missing_fonts'] == r['c_missing_fonts'] == 0
        assert sum(r['original_font_calls']) > 100
        receipts[level], full[level] = sha(path), r
    assert receipts['O0'] == receipts['O2']
    controls = {}
    for number, group in MUTANTS.items():
        path = out / 'results' / f'mutant-{number}.json'; r = json.loads(path.read_text())
        assert not r['passed'] and r['input_sha256'] == EXPECTED and r['groups'] == {group: r['cases']}
        assert 0 < r['cases'] <= GROUPS[group]
        mismatch = r['mismatch']; assert mismatch['group'] == group and mismatch['case'] == r['cases'] - 1
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_trace', '_ports', '_device', '_planes', '_ram', '_api', '_mouse', '_in', '_ticks'])
        controls[str(number)] = {'group': group, 'first_rejected': r['cases'], 'receipt_sha256': sha(path)}
    source = {}
    for line in (out / 'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1); relative = name.strip().removeprefix('/repo/')
        assert sha(repo / relative) == digest; source[relative] = digest
    compiled = sha(out / 'results/c-source.sha256')
    assert (out / 'results/compiled-source-digest.txt').read_text().strip() == compiled
    for level in ['O0', 'O2']:
        build = (out / 'results' / f'buildinfo-{level}.txt').read_text()
        assert '-DKI_CHOICE_SOURCE_DIGEST=0x' + compiled[:16] in build and 'matching_choice' in build
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='choice-repro-') as directory:
        root = Path(directory); (root / 'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo / 'tools/c_recovery_display_generate.py', root / 'tools/c_recovery_display_generate.py')
        spec = importlib.util.spec_from_file_location('choice_generate', repo / 'tools/c_recovery_choice_generate.py')
        generator = importlib.util.module_from_spec(spec); spec.loader.exec_module(generator)
        generator.generate(root, out / 'ida/ida-probe.json')
        assert (root / 'tools/c_recovery/choice_generated.inc').read_bytes() == (repo / 'tools/c_recovery/choice_generated.inc').read_bytes()
    supplement = json.loads((repo / 'docs/re/rectangle-handler-code.json').read_text())
    assert len(supplement['instructions']) == 369
    new_code = [r for r in supplement['instructions'] if 0x19440 <= r['ida_linear'] < 0x19444 or 0x1945a <= r['ida_linear'] < 0x19469]
    assert len(new_code) == 10 and sum(r['file_end'] - r['file_start'] for r in new_code) == 19
    for name, (marker, link) in BACKLINKS.items():
        text = (repo / name).read_text(); assert marker in text and link in text
    result = {'schema': 'wolong-c-choice-verification-v1', 'status': 'semantic-conformed', 'input_sha256': EXPECTED,
              'ida_database_sha256': sha(out / 'ida/input.exe.i64'), 'input_assets': assets.ASSETS,
              'new_routine_sha256': routines, 'code_entry_sha256': blocks, 'source_sha256': source,
              'compiled_source_manifest_sha256': compiled, 'cases_per_optimization': 138, 'groups': GROUPS,
              'xor_roundtrips_restored': 3, 'selected_rows': full['O2']['selected_rows'],
              'receipt_sha256': receipts, 'negative_controls_rejected': len(controls), 'mutants': controls,
              'font_calls_per_optimization': full['O2']['original_font_calls'], 'font_missing': 0,
              'entries_seen': full['O2']['entries_seen'], 'scope_backlinks': BACKLINKS,
              'exact_clean_regeneration': True, 'supplement_sha256': sha(repo / 'docs/re/rectangle-handler-code.json'),
              'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
              'scope': 'Original SS selector frame, absolute/visible rows, bands, live callback/immediates, XOR, cursor save/restore, map/pixel protection and existing speech/advise caller; independent DOS/VGA and fixed mouse/counter inputs. Current nullsub_5 image retained. Normal player flow and original runtime replacement/C machine code remain outside evidence',
              'c_machine_code_match': False}
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n')
    print('14 named C routines; 2 raw entries; 138 whole-device cases; 3 XOR roundtrips; 12 mutants: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__); parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True); args = parser.parse_args(); verify(args.repo, args.output)
