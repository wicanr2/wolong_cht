#!/usr/bin/env python3
"""核對 TALK 真實字形、肖像／DOS ABI、固定 IDA 與編譯來源收據。"""
import argparse
import ast
import hashlib
import importlib.util
import json
import shutil
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW = {'sub_106F9': 4, 'sub_1075B': 119, 'sub_107D2': 115, 'sub_1084A': 90,
       'sub_18853': 48, 'sub_189A4': 58, 'sub_1E38C': 26, 'sub_1F4DF': 73}
GROUPS = {'markers': 56, 'strings': 20, 'corpus': 1022, 'portrait': 150,
          'cache': 13, 'status': 7, 'file': 9}
MUTANTS = {1: 'corpus', 2: 'strings', 3: 'strings', 4: 'strings', 5: 'markers',
           6: 'cache', 7: 'portrait', 8: 'portrait', 9: 'strings', 10: 'corpus'}
ASSETS = {'TALK.DAT': (34182, '08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384'),
          'SINARIO.DAT': (88832, '21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c'),
          'KAOGRF.DAT': (307200, 'b9c7745e3ed9b32f0c12003fe81f82756136f738427dc421ec96f6c4ed5c4ac8'),
          'ICONGRF.DAT': (47776, '2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d'),
          'END_S13.DAT': (404992, '0ab4c43920787bf87c91bba5730078e451cb6125b878f201d3e70d648412569b'),
          'END_S14.DAT': (3840, '1d0cf09d0a319a9e7039190688c6a905ba4370bd369fbfcfa43b2078480d6918')}
BACKLINKS = {
    'docs/re/79-talk-marker-handlers.md': ('後續原始 C TALK', '110-c-talk-rendering-restoration.md'),
    'docs/re/33-shared-draw-helpers.md': ('後續原始 C TALK', '110-c-talk-rendering-restoration.md'),
    'docs/spec/119-talk-marker-fields.md': ('後續原始 C TALK', '230-c-talk-rendering.md'),
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def rendering_coverage(result, nonempty_slots):
    # A clipped corpus can still produce equal RAM/C digests. Require actual font work.
    assert sum(result['original_font_calls']) >= nonempty_slots
    assert result['original_font_calls'] == result['c_font_calls']
    assert result['original_missing_fonts'] == result['c_missing_fonts'] == 0


def verify(repo, out):
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == EXPECTED
    for name, (size, digest) in ASSETS.items():
        path = repo / 'workplace/orig/dosv' / name
        assert path.stat().st_size == size and sha(path) == digest
    talk = (repo / 'workplace/orig/dosv/TALK.DAT').read_bytes()
    offsets = [int.from_bytes(talk[i:i + 2], 'little') for i in range(0, 2048, 2)]
    assert offsets[0] == 2048 and offsets[1022] == len(talk) - 1 and offsets[1023] == 0
    assert offsets[:1022] == sorted(offsets[:1022])
    nonempty_slots = sum(talk[o] != 0 for o in offsets[:1022])
    probe = json.loads((out / 'ida/ida-probe.json').read_text())
    assert probe['tool_version'] == '9.4' and probe['function_count'] == 739
    assert probe['input_sha256'] == probe['ida_input_sha256'] == EXPECTED
    assert len(probe['targets']) == 8 and len(probe['decoded_blocks']) == 6
    routines, blocks = {}, {}
    for target in probe['targets']:
        data = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert data == b''.join(raw[c['start'] - 0x10000 + 512:c['end'] - 0x10000 + 512]
                               for c in target['chunks'])
        assert len(data) == NEW[target['name']]
        digest = hashlib.sha256(data).hexdigest()
        assert digest == target['file_sha256']
        routines[target['name']] = digest
    assert sum(NEW.values()) == 533
    for block in probe['decoded_blocks']:
        data = b''.join(bytes.fromhex(r['bytes']) for r in block['instructions'])
        assert data == raw[block['ida_linear'] - 0x10000 + 512:block['end'] - 0x10000 + 512]
        name = f'code_{block["ida_linear"]:X}'
        blocks[name] = hashlib.sha256(data).hexdigest()
    assert sum(b['end'] - b['ida_linear'] for b in probe['decoded_blocks']) == 210
    assert probe['marker_table']['offsets'] == [0x8b2, 0x8db, 0x904, 0x939, 0x95b, 0x97e, 0x984]
    assert bytes.fromhex(probe['portrait_filename']['file_bytes']) == b'KAOGRF.DAT\0'
    receipts, full = {}, {}
    for level in ['O0', 'O2']:
        path = out / 'results' / f'{level}.json'
        result = json.loads(path.read_text())
        assert result['schema'] == 'wolong-c-talk-parity-v1'
        assert result['input_sha256'] == EXPECTED and result['passed'] and result['mismatch'] is None
        assert result['groups'] == GROUPS and result['cases'] == sum(GROUPS.values())
        assert result['full_ram_plane_audits'] == result['indexed_content_audits'] == result['cases']
        assert result['original_state_sha256'] == result['c_state_sha256'] and not result['c_machine_code_match']
        assert result['corpus_slots_seen'] == {str(i): 1 for i in range(1022)}
        rendering_coverage(result, nonempty_slots)
        assert result['cache_sequence_api_counts'] == [4, 4, 4, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0]
        assert result['dos_api_calls'] == {'INT21/AH=3D': 1195, 'INT21/AH=42': 1194, 'INT21/AH=3F': 1194, 'INT21/AH=3E': 1194}
        assert all(result['routine_sha256'][n] == h and result['entries_seen'].get(n, 0) > 0
                   for n, h in {**routines, **blocks}.items())
        receipts[level], full[level] = sha(path), result
    assert receipts['O0'] == receipts['O2']
    clipped = dict(full['O2'], original_font_calls=[209, 76], c_font_calls=[209, 76])
    try:
        rendering_coverage(clipped, nonempty_slots)
    except AssertionError:
        pass
    else:
        raise AssertionError('clipped corpus accepted by coverage guard')
    controls = {}
    for number, group in MUTANTS.items():
        path = out / 'results' / f'mutant-{number}.json'
        result = json.loads(path.read_text())
        assert not result['passed'] and result['input_sha256'] == EXPECTED
        assert result['groups'] == {group: result['cases']} and 0 < result['cases'] <= GROUPS[group]
        mismatch = result['mismatch']
        assert mismatch['group'] == group and mismatch['case'] == result['cases'] - 1
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_trace', '_ports', '_device', '_planes', '_ram', '_api'])
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
        assert '-DKI_TALK_SOURCE_DIGEST=0x' + compiled[:16] in build and 'matching_talk' in build
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    with tempfile.TemporaryDirectory(prefix='talk-repro-') as directory:
        root = Path(directory)
        (root / 'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo / 'tools/c_recovery_display_generate.py', root / 'tools/c_recovery_display_generate.py')
        spec = importlib.util.spec_from_file_location('talk_generate', repo / 'tools/c_recovery_talk_generate.py')
        generator = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(generator)
        generator.generate(root, out / 'ida/ida-probe.json')
        assert (root / 'tools/c_recovery/talk_generated.inc').read_bytes() == (repo / 'tools/c_recovery/talk_generated.inc').read_bytes()
        resolver = next(n for n in ast.parse((repo / 'tools/ida_matching_probe.py').read_text()).body
                        if isinstance(n, ast.FunctionDef) and n.name == 'evidence_path')
        scope = {'Path': Path}
        exec(compile(ast.Module(body=[resolver], type_ignores=[]), 'evidence_path', 'exec'), scope)
        (root / 're').mkdir(); (root / 'spec').mkdir(); (root / 'spec/x.md').write_text('evidence')
        assert scope['evidence_path']('docs/spec/x.md', root).is_file()
        assert not scope['evidence_path']('docs/re/x.md', root).exists()
        for invalid in ['docs/../x.md', '/docs/spec/x.md', 'spec/x.md', '../docs/spec/x.md']:
            try:
                scope['evidence_path'](invalid, root)
            except ValueError:
                pass
            else:
                raise AssertionError('invalid evidence path accepted: ' + invalid)
    for name, (marker, link) in BACKLINKS.items():
        text = (repo / name).read_text()
        assert marker in text and link in text
    result = {'schema': 'wolong-c-talk-verification-v1', 'status': 'semantic-conformed',
              'input_sha256': EXPECTED, 'ida_database_sha256': sha(out / 'ida/input.exe.i64'),
              'input_assets': ASSETS,
              'new_routine_sha256': routines, 'code_entry_sha256': blocks, 'source_sha256': source,
              'compiled_source_manifest_sha256': compiled, 'cases_per_optimization': sum(GROUPS.values()),
              'groups': GROUPS, 'corpus_slots_audited': 1022, 'nonempty_slot_font_lower_bound': nonempty_slots,
              'font_calls_per_optimization': full['O2']['original_font_calls'], 'font_missing': 0,
              'dos_api_calls_per_optimization': full['O2']['dos_api_calls'],
              'cache_sequence_api_counts': full['O2']['cache_sequence_api_counts'],
              'receipt_sha256': receipts, 'entries_seen': full['O2']['entries_seen'],
              'negative_controls_rejected': len(controls), 'mutants': controls,
              'clipped_corpus_guard_rejected': True, 'exact_clean_regeneration': True,
              'scope_backlinks': BACKLINKS,
              'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
              'scope': 'Original TALK slots/variants/live near table, raw SS/name fields, real Font and portraits, DOS API/stack/CF and four-slot cache; independent services and native C; IF/TF=0; INT50 media UI, normal player flow and C machine code remain outside evidence',
              'c_machine_code_match': False}
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n')
    print('8 named C routines; 6 raw handlers; 1277 full-state cases; 1022 audited corpus slots; 10 mutants: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    verify(args.repo, args.output)
