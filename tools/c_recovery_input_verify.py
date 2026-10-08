#!/usr/bin/env python3
"""核對兩 segment input、cursor、counter 與指令重建的固定來源收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import struct
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW = {'sub_18810': 67, 'sub_121E7': 47, 'sub_20000': 14, 'sub_2002E': 65, 'nullsub_4': 1,
       'sub_20070': 42, 'sub_2009A': 35, 'sub_200BD': 3, 'sub_200C0': 54, 'sub_20137': 52,
       'sub_2016B': 50, 'sub_2019D': 41, 'sub_201C6': 30, 'sub_201E4': 40, 'sub_2020C': 61,
       'sub_20249': 87, 'sub_202A0': 29, 'sub_202BD': 65, 'sub_202FE': 29}
GROUPS = {'services': 48, 'cursor': 51, 'edges': 45, 'alignment': 8, 'wait': 18,
          'poll': 14, 'message': 15, 'patched-caller': 2, 'roundtrip': 9}
MUTANTS = {1: 'poll', 2: 'wait', 3: 'wait', 4: 'wait', 5: 'cursor', 6: 'cursor',
           7: 'edges', 8: 'edges', 9: 'cursor', 10: 'services', 11: 'cursor', 12: 'message'}
BACKLINKS = {'docs/spec/45-advise-scene-layout.md': ('後續原始 C 等待', '231-c-input-mouse.md'),
             'docs/re/42-leaf-functions.md': ('後續原始 C 等待', '111-c-input-mouse-restoration.md'),
             'docs/re/25-message-variants-and-personnel.md': ('後續原始 C 等待', '111-c-input-mouse-restoration.md')}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify(repo, out):
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == EXPECTED
    asset_spec = importlib.util.spec_from_file_location('talk_contract', repo / 'tools/c_recovery_talk_verify.py')
    asset_module = importlib.util.module_from_spec(asset_spec)
    asset_spec.loader.exec_module(asset_module)
    for name, (size, digest) in asset_module.ASSETS.items():
        path = repo / 'workplace/orig/dosv' / name
        assert path.stat().st_size == size and sha(path) == digest
    probe = json.loads((out / 'ida/ida-probe.json').read_text())
    assert probe['tool_version'] == '9.4' and probe['function_count'] == 739
    assert probe['input_sha256'] == probe['ida_input_sha256'] == EXPECTED
    assert len(probe['targets']) == 19 and len(probe['decoded_blocks']) == 2
    routines, blocks = {}, {}
    for target in probe['targets']:
        data = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert data == b''.join(raw[c['start'] - 0x10000 + 512:c['end'] - 0x10000 + 512]
                               for c in target['chunks'])
        assert len(data) == NEW[target['name']]
        routines[target['name']] = hashlib.sha256(data).hexdigest()
        assert routines[target['name']] == target['file_sha256']
    assert sum(NEW.values()) == 812
    for block in probe['decoded_blocks']:
        start, end = block['ida_linear'] - 0x10000 + 512, block['end'] - 0x10000 + 512
        blocks[f'code_{block["ida_linear"]:X}'] = hashlib.sha256(raw[start:end]).hexdigest()
        for row in block['instructions']:
            at = row['ida_linear'] - 0x10000 + 512
            expected = bytearray(raw[at:at + len(bytes.fromhex(row['bytes']))])
            for offset in probe['relocation_file_offsets']:
                if at <= offset < at + len(expected):
                    value = struct.unpack_from('<H', raw, offset)[0]
                    struct.pack_into('<H', expected, offset - at, (value + 0x1000) & 65535)
            assert expected.hex() == row['bytes']
    assert probe['patch_word']['file_bytes'] == 'eb0f'
    assert probe['mouse_table']['offsets'] == [46, 112, 154, 189, 189, 189, 189, 189, 189, 192, 111, 111, 111, 111, 111, 46]
    receipts, full = {}, {}
    for level in ['O0', 'O2']:
        path = out / 'results' / f'{level}.json'
        result = json.loads(path.read_text())
        assert result['schema'] == 'wolong-c-input-parity-v1'
        assert result['input_sha256'] == EXPECTED and result['passed'] and result['mismatch'] is None
        assert result['groups'] == GROUPS and result['cases'] == 210
        assert result['full_ram_plane_audits'] == result['indexed_content_audits'] == 210
        assert result['original_state_sha256'] == result['c_state_sha256'] and not result['c_machine_code_match']
        assert result['cursor_roundtrips_restored'] == 3
        assert result['counter_checkpoints']['wait'] == 140 and result['counter_checkpoints']['message'] == 110
        assert result['mouse_queries']['wait'] == 220 and result['mouse_queries']['poll'] == 63
        assert all(result['routine_sha256'][n] == h and result['entries_seen'].get(n, 0) > 0
                   for n, h in {**routines, **blocks}.items())
        assert result['original_font_calls'] == result['c_font_calls'] and result['original_missing_fonts'] == result['c_missing_fonts'] == 0
        receipts[level], full[level] = sha(path), result
    assert receipts['O0'] == receipts['O2']
    controls = {}
    for number, group in MUTANTS.items():
        path = out / 'results' / f'mutant-{number}.json'
        result = json.loads(path.read_text())
        assert not result['passed'] and result['input_sha256'] == EXPECTED
        assert result['groups'] == {group: result['cases']} and 0 < result['cases'] <= GROUPS[group]
        mismatch = result['mismatch']
        assert mismatch['group'] == group and mismatch['case'] == result['cases'] - 1
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_trace', '_ports', '_device', '_planes', '_ram', '_api', '_mouse', '_in', '_ticks'])
        controls[str(number)] = {'group': group, 'first_rejected': result['cases'], 'receipt_sha256': sha(path)}
    source = {}
    for line in (out / 'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1)
        relative = name.strip().removeprefix('/repo/')
        assert sha(repo / relative) == digest
        source[relative] = digest
    compiled = sha(out / 'results/c-source.sha256')
    assert (out / 'results/compiled-source-digest.txt').read_text().strip() == compiled
    for level in ['O0', 'O2']:
        build = (out / 'results' / f'buildinfo-{level}.txt').read_text()
        assert '-DKI_INPUT_SOURCE_DIGEST=0x' + compiled[:16] in build and 'matching_input' in build
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='input-repro-') as directory:
        root = Path(directory)
        (root / 'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo / 'tools/c_recovery_display_generate.py', root / 'tools/c_recovery_display_generate.py')
        spec = importlib.util.spec_from_file_location('input_generate', repo / 'tools/c_recovery_input_generate.py')
        generator = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(generator)
        generator.generate(root, out / 'ida/ida-probe.json')
        assert (root / 'tools/c_recovery/input_generated.inc').read_bytes() == (repo / 'tools/c_recovery/input_generated.inc').read_bytes()
    supplement = json.loads((repo / 'docs/re/rectangle-handler-code.json').read_text())
    assert len(supplement['instructions']) == 359
    new_code = [r for r in supplement['instructions'] if 0x12239 <= r['ida_linear'] < 0x12280]
    assert len(new_code) == 21 and sum(r['file_end'] - r['file_start'] for r in new_code) == 71
    assert sum(len(r.get('relocation_file_offsets', [])) for r in new_code) == 3
    for name, (marker, link) in BACKLINKS.items():
        text = (repo / name).read_text()
        assert marker in text and link in text
    result = {'schema': 'wolong-c-input-verification-v1', 'status': 'semantic-conformed',
              'input_sha256': EXPECTED, 'ida_database_sha256': sha(out / 'ida/input.exe.i64'),
              'input_assets': asset_module.ASSETS,
              'new_routine_sha256': routines, 'code_entry_sha256': blocks, 'source_sha256': source,
              'compiled_source_manifest_sha256': compiled, 'cases_per_optimization': 210, 'groups': GROUPS,
              'cursor_roundtrips_restored': 3, 'counter_checkpoints': full['O2']['counter_checkpoints'],
              'mouse_queries': full['O2']['mouse_queries'], 'receipt_sha256': receipts,
              'negative_controls_rejected': len(controls), 'mutants': controls, 'entries_seen': full['O2']['entries_seen'],
              'scope_backlinks': BACKLINKS, 'exact_clean_regeneration': True,
              'supplement_sha256': sha(repo / 'docs/re/rectangle-handler-code.json'),
              'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
              'scope': 'Original two-segment near/far mouse APIs, live patch/dispatch table, cursor/memory/VGA/IN/OUT/IF and message render-wait-clear; explicit fixed press/counter inputs; no C-side CPU.Step. Direct callback is not automatic player callback dispatch; no timer wall-clock or C machine-code claim',
              'c_machine_code_match': False}
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n')
    print('19 named C routines; 2 raw entries; 210 whole-device cases; 3 cursor roundtrips; 12 mutants: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    verify(args.repo, args.output)
