#!/usr/bin/env python3
"""驗證完整讀檔、原始 cue 與 allocation 的原版／C 收據。"""
import argparse
import hashlib
import importlib.util
import json
import shutil
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
NEW = {'sub_187AF': 29, 'sub_1E378': 20, 'sub_1F4A2': 59,
       'sub_10241': 129, 'sub_102C2': 14, 'sub_100DF': 118}
GROUPS = {'stream': 16, 'assets': 5, 'mmap-caller': 2, 'cue': 64,
          'stop': 5, 'allocation': 1}
MUTANTS = {1: 'stream', 2: 'stream', 3: 'stream', 4: 'mmap-caller',
           5: 'cue', 6: 'cue', 7: 'cue', 8: 'cue', 9: 'allocation',
           10: 'allocation', 11: 'stop', 12: 'stream'}
ASSETS = {
    'MMAP.MDL': (32768, '2fa1dd1b1ec7c426cf22583334a62dc61bdb1f1cefc59cbe8e9827480d118d1d'),
    'MMAP.MCH': (43058, 'b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7'),
    'BGM.DAT': (20826, '7a51c8b9a349b9e088f3796b70c268181c60bcebead70942f00e1621523dedc9'),
    'IVENTGRF.DAT': (76032, '23fffa03b2bba3b4c920c5db49b32b7a3fa26729523b0092a6be5ad3a0737f2f'),
}
BACKLINKS = {
    'docs/re/04-mmap-entry-points.md': '113-c-resource-cue-restoration.md',
    'docs/re/58-bgm-scene-mapping.md': '113-c-resource-cue-restoration.md',
    'docs/re/106-c-rectangle-bars-restoration.md': '113-c-resource-cue-restoration.md',
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify(repo, out):
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == EXPECTED
    spec = importlib.util.spec_from_file_location('talk_contract', repo / 'tools/c_recovery_talk_verify.py')
    assets = importlib.util.module_from_spec(spec); spec.loader.exec_module(assets)
    inputs = {**assets.ASSETS, **ASSETS}
    for name, (size, digest) in inputs.items():
        path = repo / 'workplace/orig/dosv' / name
        assert path.stat().st_size == size and sha(path) == digest
    probe = json.loads((out / 'ida/ida-probe.json').read_text())
    assert probe['tool_version'] == '9.4' and probe['function_count'] == 739
    assert probe['input_sha256'] == probe['ida_input_sha256'] == EXPECTED
    assert {t['name'] for t in probe['targets']} == set(NEW)
    routines, instructions = {}, 0
    for target in probe['targets']:
        data = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert len(data) == NEW[target['name']]
        assert data == b''.join(raw[c['start'] - 0x10000 + 512:c['end'] - 0x10000 + 512] for c in target['chunks'])
        routines[target['name']] = hashlib.sha256(data).hexdigest()
        assert routines[target['name']] == target['file_sha256']
        instructions += sum(len(c['instructions']) for c in target['chunks'])
    assert sum(NEW.values()) == 369 and instructions == 165
    receipts, full = {}, {}
    for level in ['O0', 'O2']:
        path = out / 'results' / f'{level}.json'; r = json.loads(path.read_text())
        assert r['schema'] == 'wolong-c-resource-parity-v1' and r['input_sha256'] == EXPECTED
        assert r['passed'] and r['mismatch'] is None and r['cases'] == 93 and r['groups'] == GROUPS
        assert r['full_ram_plane_audits'] == r['indexed_content_audits'] == 93
        assert r['complete_buffer_audits'] == 21
        assert r['asset_banks'] == 'Stream buffers at 3000:0800; MMAP at 3000; BGM at 6100; SS at F000:8000'
        assert r['original_state_sha256'] == r['c_state_sha256'] and not r['c_machine_code_match']
        assert all(r['routine_sha256'][n] == h and r['entries_seen'].get(n, 0) > 0 for n, h in routines.items())
        assert r['entries_seen']['sub_10241'] == 64 and r['entries_seen']['sub_100DF'] == 1
        assert r['entries_seen']['sub_1F4A2'] == 25 and r['entries_seen']['sub_1E378'] == 17
        receipts[level], full[level] = sha(path), r
    assert receipts['O0'] == receipts['O2']
    controls = {}
    for number, group in MUTANTS.items():
        path = out / 'results' / f'mutant-{number}.json'; r = json.loads(path.read_text())
        assert not r['passed'] and r['input_sha256'] == EXPECTED and r['groups'] == {group: r['cases']}
        assert 0 < r['cases'] <= GROUPS[group]
        mismatch = r['mismatch']; assert mismatch['group'] == group and mismatch['case'] == r['cases'] - 1
        assert any(mismatch['original' + suffix] != mismatch['c' + suffix]
                   for suffix in ['', '_trace', '_ports', '_device', '_planes', '_ram', '_api', '_sound', '_mouse', '_in', '_ticks'])
        controls[str(number)] = {'group': group, 'first_rejected': r['cases'], 'receipt_sha256': sha(path)}
    source = {}
    for line in (out / 'results/c-source.sha256').read_text().splitlines():
        digest, name = line.split(None, 1); relative = name.strip().removeprefix('/repo/')
        assert sha(repo / relative) == digest; source[relative] = digest
    compiled = sha(out / 'results/c-source.sha256')
    assert (out / 'results/compiled-source-digest.txt').read_text().strip() == compiled
    for level in ['O0', 'O2']:
        build = (out / 'results' / f'buildinfo-{level}.txt').read_text()
        assert '-DKI_RESOURCE_SOURCE_DIGEST=0x' + compiled[:16] in build and 'matching_resource' in build
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    versions = (out / 'results/tool-versions.txt').read_text()
    assert 'go version go1.26.7 linux/amd64' in versions and '12.2.0' in versions
    with tempfile.TemporaryDirectory(prefix='resource-repro-') as directory:
        root = Path(directory); (root / 'tools/c_recovery').mkdir(parents=True)
        shutil.copyfile(repo / 'tools/c_recovery_display_generate.py', root / 'tools/c_recovery_display_generate.py')
        spec = importlib.util.spec_from_file_location('resource_generate', repo / 'tools/c_recovery_resource_generate.py')
        generator = importlib.util.module_from_spec(spec); spec.loader.exec_module(generator)
        generator.generate(root, out / 'ida/ida-probe.json')
        assert (root / 'tools/c_recovery/resource_generated.inc').read_bytes() == (repo / 'tools/c_recovery/resource_generated.inc').read_bytes()
    for name, link in BACKLINKS.items():
        assert '後續原始 C 資源' in (repo / name).read_text() and link in (repo / name).read_text()
    result = {'schema': 'wolong-c-resource-verification-v1', 'status': 'semantic-conformed', 'input_sha256': EXPECTED,
              'ida_database_sha256': sha(out / 'ida/input.exe.i64'), 'input_assets': inputs,
              'new_routine_sha256': routines, 'original_instruction_count': instructions,
              'source_sha256': source, 'compiled_source_manifest_sha256': compiled,
              'cases_per_optimization': 93, 'groups': GROUPS, 'complete_buffer_audits': 21,
              'receipt_sha256': receipts, 'negative_controls_rejected': len(controls), 'mutants': controls,
              'entries_seen': full['O2']['entries_seen'], 'scope_backlinks': BACKLINKS,
              'exact_clean_regeneration': True,
              'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
              'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
              'tool_versions': versions,
              'verification_tools_sha256': {p: sha(repo / p) for p in [
                  'tools/c_recovery_resource_verify.py', 'tools/c_recovery_resource.sh',
                  'tools/ida_resource_probe.py']},
              'scope': 'Original full-file DOS reads and retries on successful media, 0xF000 boundaries, complete source buffers and tail guards, MMAP pair caller, 16 cue indices with repeat/mute and stop, successful two-allocation segment layout. Independent DOS/INT61 services; full RAM/VGA/register/stack/API/sound state. INT50 error GUI, allocation failure exit, normal player flow, original TSR audio and C machine code remain outside evidence',
              'c_machine_code_match': False}
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n')
    print('6 named C routines; 165 original instructions; 93 whole-device cases; 21 buffers; 12 mutants: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__); parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True); args = parser.parse_args(); verify(args.repo, args.output)
