#!/usr/bin/env python3
"""核對軍團C的固定來源、乾淨重生、原生建置與同狀態收據。"""
import argparse
import hashlib
import importlib.util
import json
import re
import shlex
import struct
import tempfile
from pathlib import Path

INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PARENT_SHA = '937f8738225ae0f9e5efc133245316dd075d845ffded97a22f5f211b4804f2a7'
DB_SHA = '4f01452cc2df29a9306aa4b444d1d4125e3f27f7730fcfdadbefd9f9da3a4667'
ROOTS = (0x125A3, 0x12600, 0x1264A, 0x12662, 0x126FF, 0x12708,
         0x127A2, 0x127F6, 0x12804, 0x12808, 0x12831, 0x12880,
         0x128F4, 0x12A7E, 0x142AB, 0x14300, 0x1562B)
TAGS = 'matching_army,matching_engagement,matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga'
SHARED_GO = ('engagement_data', 'engagement_control_data', 'engagement_ui_data',
             'engagement_front_data', 'engagement_extra_data', 'engagement_combat_data',
             'engagement_save', 'outcome_data', 'route_data', 'interaction_data',
             'main_data', 'strategy_data', 'engagement_platform', 'vga_bus',
             'glyph_platform', 'input_platform')
# 在工廠與完整原版first矩陣凍結後填入實際審查值；未定不得發布CONFORMED。
GROUPS = None
STAGES = None
RETURNS = None
INDEPENDENT_AUDITS = None


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def read(path):
    return json.loads(path.read_text(encoding='utf-8'))


def load_module(repo, name):
    spec = importlib.util.spec_from_file_location(name, repo / 'tools' / (name + '.py'))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def native_build(binary, info, compiled, mutation, optimization):
    """直接讀ELF的Go buildinfo，核對另存輸出與實際CGO旗標。"""
    data = binary.read_bytes()
    assert data[:6] == b'\x7fELF\x02\x01'
    shoff = struct.unpack_from('<Q', data, 40)[0]
    stride, count, ni = struct.unpack_from('<HHH', data, 58)
    assert stride >= 64 and 0 < ni < count and shoff + stride * count <= len(data)
    sections = [struct.unpack_from('<IIQQQQIIQQ', data, shoff + i * stride) for i in range(count)]
    ns = sections[ni]
    names = data[ns[4]:ns[4] + ns[5]]
    matches = [s for s in sections if names[s[0]:].split(b'\0', 1)[0] == b'.go.buildinfo']
    assert len(matches) == 1
    s = matches[0]
    blob = data[s[4]:s[4] + s[5]]
    assert blob[:16] == b'\xff Go buildinf:\x08\x02'

    def string(at):
        length = shift = 0
        while True:
            assert at < len(blob) and shift < 64
            value = blob[at]
            at += 1
            length |= (value & 127) << shift
            if value < 128:
                break
            shift += 7
        assert at + length <= len(blob)
        return blob[at:at + length], at + length

    version, at = string(32)
    module, _ = string(at)
    assert version == b'go1.26.7' and len(module) >= 33 and module[-17] == 10
    lines = module[16:-16].decode().splitlines()
    saved = info.read_text(encoding='utf-8').splitlines()
    assert 'path\tgithub.com/wicanr2/dosgolem/wolongcarmy' in lines
    assert saved[0].endswith(': go1.26.7') and Path(saved[0].split(':', 1)[0]).name == binary.name

    def flags(items):
        prefix = 'build\tCGO_CFLAGS='
        found = [x.removeprefix('\t')[len(prefix):] for x in items if x.removeprefix('\t').startswith(prefix)]
        assert len(found) == 1
        quoted = shlex.split(found[0])
        assert len(quoted) == 1
        return shlex.split(quoted[0])

    actual = flags(lines)
    assert actual == flags(saved[1:])
    assert [v for v in actual if re.fullmatch('-O[0-9sgz]', v)] == ['-' + optimization]
    assert [v for v in actual if v.startswith('-DKI_ARMY_SOURCE_DIGEST=')] == ['-DKI_ARMY_SOURCE_DIGEST=0x' + compiled[:16]]
    assert [v for v in actual if v.startswith('-DKI_ARMY_MUTATION=')] == ([] if mutation is None else [f'-DKI_ARMY_MUTATION={mutation}'])
    assert 'build\t-tags=' + TAGS in lines and 'build\tCGO_ENABLED=1' in lines
    return {'binary_sha256': sha(binary), 'buildinfo_sha256': sha(info), 'cgo_flags': actual}


def source_manifest(repo, out):
    manifest = out / 'results/c-source.sha256'
    records = {}
    for line in manifest.read_text().splitlines():
        digest, name = line.split(None, 1)
        assert re.fullmatch('[0-9a-f]{64}', digest) and name.startswith('/repo/')
        path = repo / name.removeprefix('/repo/')
        assert sha(path) == digest, name
        records[str(path.relative_to(repo))] = digest
    required = {str(p.relative_to(repo)) for p in (repo / 'tools/c_recovery').iterdir() if p.suffix in ('.c', '.h', '.inc')}
    required.update(('tools/c_recovery_army_generate.py', 'tools/c_recovery_army.go',
                     'tools/c_recovery_army_data.go', 'tools/c_recovery_display_generate.py'))
    required.update('tools/c_recovery_' + name + '.go' for name in SHARED_GO)
    assert required <= records.keys()
    compiled = sha(manifest)
    assert (out / 'results/compiled-source-digest.txt').read_text().strip() == compiled
    return records, compiled


def primary(repo, out):
    assert sha(repo / 'workplace/orig/dosv/KI.EXE') == INPUT_SHA
    parent_path = repo / 'workplace/matching-decompilation/c-tick/census/ida-probe.json'
    assert sha(parent_path) == PARENT_SHA
    assert sha(parent_path.parent / 'input.exe.i64') == DB_SHA
    source = read(out / 'ida-source.json')
    parent = read(parent_path)
    assert source['fresh_ida_analysis'] is False
    assert source['input_sha256'] == source['ida_input_sha256'] == INPUT_SHA
    assert source['tool_version'] == '9.4' and source['ida_load_paragraph'] == 0x1000
    assert source['targets'] == sorted((t for t in parent['targets'] if t['ida_linear'] in ROOTS), key=lambda t: t['ida_linear'])
    assert tuple(t['ida_linear'] for t in source['targets']) == ROOTS
    assert source['decoded_blocks'] == [] and source['unresolved_call_boundaries'] == []
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    count = size = 0
    routines = {}
    for target in source['targets']:
        body = b''
        for chunk in target['chunks']:
            assert not chunk['loader_relocations']
            chunk_bytes = bytes.fromhex(chunk['file_bytes'])
            at = chunk['start'] - 0x10000 + 512
            assert raw[at:at + len(chunk_bytes)] == chunk_bytes
            assert b''.join(bytes.fromhex(row['bytes']) for row in chunk['instructions']) == chunk_bytes
            count += len(chunk['instructions'])
            size += len(chunk_bytes)
            body += chunk_bytes
        assert hashlib.sha256(body).hexdigest() == target['file_sha256']
        routines[target['name']] = target['file_sha256']
    assert (count, size) == (434, 1068)
    coverage = read(repo / 'docs/re/c-army-code.json')
    assert coverage['probe_sha256'] == sha(out / 'ida-source.json')
    assert len(coverage['instruction_coverage']) == 434 and coverage['instructions'] == []
    assert all(r['independent_native_assembly_match'] and r['fixed_prior_coverage'] for r in coverage['instruction_coverage'])
    original_rows = sorted((row for t in source['targets'] for ch in t['chunks'] for row in ch['instructions']), key=lambda row: row['ida_linear'])
    for row, published in zip(original_rows, coverage['instruction_coverage'], strict=True):
        assert row['ida_linear'] == published['ida_linear']
        assert row['file_offset'] == published['file_start']
        assert row['bytes'] == published['bytes'] == published['ida_loaded_bytes']
        assert published['file_end'] == row['file_offset'] + len(bytes.fromhex(row['bytes']))
    assembly = coverage['assembler_verification']
    assert assembly['status'] == 'byte-exact' and assembly['instructions'] == assembly['bytes'] == 0
    assert assembly['independent_reassembly']['instructions'] == assembly['published_slice_reassembly']['instructions'] == 434
    assert assembly['independent_reassembly']['bytes'] == assembly['published_slice_reassembly']['bytes'] == 1068
    assert assembly['instruction_byte_fallbacks'] == 0
    for key in ('published_manifest', 'published_code_source', 'published_linker'):
        assert sha(repo / assembly[key]) == assembly[key + '_sha256']
    assert sha(repo / 'tools/c_recovery_army_prepare.py') == assembly['generator_sha256']
    assert sha(repo / 'tools/assembly_rebuild.py') == assembly['translator_sha256']
    for stem, key in (('army-proof', 'independent_reassembly'), ('army-published-slice', 'published_slice_reassembly')):
        receipt = assembly[key]
        assert sha(out / (stem + '.S')) == receipt['source_sha256']
        linker = out / (stem + '.ld') if key == 'independent_reassembly' else repo / assembly['published_linker']
        assert sha(linker) == receipt['linker_sha256']
        binary = out / (stem + '.bin')
        assert sha(binary) == receipt['binary_sha256']
        data = binary.read_bytes()
        for row in original_rows:
            at, encoded = row['file_offset'], bytes.fromhex(row['bytes'])
            assert data[at:at + len(encoded)] == encoded
    generator = load_module(repo, 'c_recovery_army_generate')
    controls = ''.join(f'{number}\t{group}\t{name}\n' for number, group, name, ea, statement in generator.MUTANTS)
    assert (repo / 'tools/c_recovery_army_mutants.tsv').read_text() == controls
    with tempfile.TemporaryDirectory(prefix='army-regenerate-') as tmp:
        generator.generate(repo, out / 'ida-source.json', Path(tmp))
        generated = {}
        for name in ('army.h', 'army_generated.inc'):
            assert (Path(tmp) / name).read_bytes() == (repo / 'tools/c_recovery' / name).read_bytes(), name
            generated['tools/c_recovery/' + name] = sha(Path(tmp) / name)
    return source, routines, generated


def verify(repo, out, source_only=False):
    source, routines, generated = primary(repo, out)
    if source_only:
        print('17 named / 434 instructions / 1068 bytes / clean generation: PASS')
        return
    records, compiled = source_manifest(repo, out)
    assert all(value is not None for value in (GROUPS, STAGES, RETURNS, INDEPENDENT_AUDITS)), 'Army semantic matrix has not been finalized'
    reports = {opt: read(out / 'results' / (opt + '.json')) for opt in ('O0', 'O2')}
    assert reports['O0'] == reports['O2'], 'Optimization reports differ'
    report = reports['O2']
    assert report['passed'] and report['mismatch'] is None
    assert report['input_sha256'] == INPUT_SHA
    assert report['new_probe_sha256'] == sha(out / 'ida-source.json')
    assert report['new_routine_sha256'] == routines
    assert report['source_manifest']['sha256'] == compiled
    assert report['cases'] > 0 and report['full_ram_plane_audits'] >= report['cases']
    assert report['groups'] == GROUPS and report['stage_audits'] == STAGES
    assert report['cases'] == sum(GROUPS.values())
    assert report['full_ram_plane_audits'] == sum(STAGES.values())
    assert report['palette_dac_audits'] == report['full_ram_plane_audits']
    assert (report['normal_return_cases'], report['nonlocal_exit_cases']) == RETURNS
    assert report['army_independent_audits'] == INDEPENDENT_AUDITS
    assert sum(report['groups'].values()) == report['cases']
    assert sum(report['stage_audits'].values()) == report['full_ram_plane_audits']
    assert report['original_state_sha256'] == report['c_state_sha256']
    assert report['independent_raw_initialization'] and not report['cross_machine_snapshot_initialization']
    assert report['fixed_raw_rng_state'] and report['rng_initial_state_sha256']
    assert not report['c_machine_code_match']
    for target in source['targets']:
        assert report['entries_seen'].get(target['name'], 0) > 0, target['name']
    builds = {opt: native_build(out / ('army-' + opt), out / 'results' / ('buildinfo-' + opt + '.txt'), compiled, None, opt) for opt in reports}
    mutants = {}
    for line in (repo / 'tools/c_recovery_army_mutants.tsv').read_text().splitlines():
        number, group, name = line.split('\t')
        path = out / 'results' / ('mutant-' + number + '.json')
        bad = read(path)
        assert (out / 'results' / ('mutant-' + number + '.exit')).read_text().strip() == '1'
        assert not bad['passed'] and bad['mismatch'] is not None
        assert set(bad['groups']) == {group} and bad['groups'][group] == bad['cases']
        assert bad['source_manifest']['sha256'] == compiled
        assert bad['input_sha256'] == INPUT_SHA and bad['new_probe_sha256'] == report['new_probe_sha256']
        assert bad['new_routine_sha256'] == routines
        difference = bad['mismatch']
        assert difference['group'] == group
        assert any(difference.get('original' + suffix) != difference.get('c' + suffix) for suffix in (
            '', '_ram', '_planes', '_trace', '_ports', '_in', '_api', '_palette_dac'))
        log = (out / 'results' / ('mutant-' + number + '.log')).read_text()
        assert 'pass=false' in log and not any(word in log for word in ('panic:', 'SIGABRT', 'SIGSEGV', 'unsupported='))
        build = native_build(out / ('mutant-' + number), out / 'results' / ('buildinfo-mutant-' + number + '.txt'), compiled, int(number), 'O0')
        mutants[number] = {'group': group, 'name': name, 'exit_status': 1, 'rejection_kind': 'state-mismatch', 'receipt_sha256': sha(path), **build}
    assert set(mutants) == {str(number) for number in range(1, 13)}
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    for line in before.decode().splitlines():
        digest, name = line.split(None, 1)
        assert sha(Path(name)) == digest
    assert source_manifest(repo, out) == (records, compiled)
    result = {
        'schema': 'wolong-c-army-verification-v1', 'status': 'semantic-conformed',
        'input_sha256': INPUT_SHA, 'ida_database_sha256': DB_SHA, 'parent_probe_sha256': PARENT_SHA,
        'scoped_source_sha256': sha(out / 'ida-source.json'), 'fresh_ida_analysis': False,
        'named_functions': 17, 'original_instruction_count': 434, 'original_instruction_bytes': 1068,
        'new_routine_sha256': routines, 'source_sha256': records, 'compiled_source_manifest_sha256': compiled,
        'native_builds': builds, 'cases_per_optimization': report['cases'],
        'whole_device_stage_audits_per_optimization': report['full_ram_plane_audits'],
        'groups': report['groups'], 'stage_audits': report['stage_audits'], 'entries_seen': report['entries_seen'],
        'receipt_sha256': {opt: sha(out / 'results' / (opt + '.json')) for opt in reports},
        'negative_controls_rejected': len(mutants), 'mutants': mutants,
        'assembly_coverage_receipt_sha256': sha(repo / 'docs/re/c-army-code.json'),
        'exact_clean_regeneration': True, 'clean_regenerated_files': generated,
        'compiled_sources_unchanged_after_execution': True,
        'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
        'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
        'tool_versions': (out / 'results/tool-versions.txt').read_text(),
        'normal_return_cases': report['normal_return_cases'], 'nonlocal_exit_cases': report['nonlocal_exit_cases'],
        'rng_initial_state_sha256': report['rng_initial_state_sha256'], 'rng_initialization': report['rng_initialization'],
        'scope': report['scope'], 'platform': report['platform'], 'c_machine_code_match': False,
        'verification_tools_sha256': {name: sha(repo / name) for name in (
            'tools/c_recovery_army.sh', 'tools/c_recovery_army_container.sh', 'tools/c_recovery_army_verify.py',
            'tools/c_recovery_army_prepare.py', 'tools/c_recovery_army_generate.py', 'tools/c_recovery_army_mutants.tsv')},
    }
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(f'17 named / {report["cases"]} cases / {len(mutants)} state mismatches: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--source-only', action='store_true')
    args = parser.parse_args()
    verify(args.repo, args.output, args.source_only)
