#!/usr/bin/env python3
"""核對開機C的固定來源、乾淨重生、原生建置與同狀態收據。"""
import argparse
import hashlib
import importlib.util
import json
import re
import shlex
import struct
import tarfile
import tempfile
from pathlib import Path

INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PROBE_SHA = '63e12863b199542902260e5f6b917c00e0f268b923ec0149a7f321d3f07a8a4a'
DB_SHA = 'c36861d16e82ddc1fbee8a36e62b49441919fac37a3e5f106b2648b582b58dcd'
ROOTS = (0x145F8, 0x15AD1, 0x15AFC, 0x15C14, 0x19377, 0x193D7,
         0x1AB7C, 0x1C300, 0x1E567, 0x1E6FF)
RAW_BLOCKS = {0x1AB48: 1}
TAGS = 'matching_residual,matching_glyph,matching_vga'
SHARED_GO = ('residual_platform', 'vga_bus', 'glyph_platform')
# 完整矩陣跑完後凍結；凍結前 verify 只跑 source-only。
GROUPS = {'diplomacy': 5, 'draw-queue': 4, 'map-node': 4, 'picker-entry': 4, 'picker-loop': 5, 'picker-row': 12, 'season': 34, 'sidebar': 8, 'tactical-cap': 13}
STAGES = {'return': 89}
RETURNS = {'return': 89}
INDEPENDENT_AUDITS = 578
BACKLINKS = ()


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
    assert 'path\tgithub.com/wicanr2/dosgolem/wolongcresidual' in lines
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
    assert [v for v in actual if v.startswith('-DKI_RESIDUAL_SOURCE_DIGEST=')] == ['-DKI_RESIDUAL_SOURCE_DIGEST=0x' + compiled[:16]]
    assert [v for v in actual if v.startswith('-DKI_RESIDUAL_MUTATION=')] == ([] if mutation is None else [f'-DKI_RESIDUAL_MUTATION={mutation}'])
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
    required.update(('tools/c_recovery_residual_generate.py', 'tools/c_recovery_residual.go',
                     'tools/c_recovery_residual_data.go', 'tools/c_recovery_display_generate.py'))
    required.update('tools/c_recovery_' + name + '.go' for name in SHARED_GO)
    assert required <= records.keys()
    compiled = sha(manifest)
    assert (out / 'results/compiled-source-digest.txt').read_text().strip() == compiled
    return records, compiled


def source_snapshots(out, records, compiled, slug):
    verified = []
    manifest = (out / 'results/c-source.sha256').read_bytes()
    for path in sorted((out / 'results').glob('*.sources.tar.gz')):
        sidecar = path.with_name(path.name + '.sha256')
        if not sidecar.exists():
            continue  # Interrupted archives cannot establish source provenance.
        try:
            archive = tarfile.open(path, 'r:gz')
        except (tarfile.TarError, EOFError):
            continue
        with archive:
            members = {member.name.removeprefix('./'): member for member in archive.getmembers()}
            entry = members.get('manifest.sha256')
            if entry is None or not entry.isfile():
                continue
            content = archive.extractfile(entry).read()
            if hashlib.sha256(content).hexdigest() != compiled:
                continue
            assert content == manifest
            assert sha(path) == sidecar.read_text().split()[0]
            for name, digest in records.items():
                member = members[name]
                assert member.isfile()
                assert hashlib.sha256(archive.extractfile(member).read()).hexdigest() == digest, name
            module = members['module-go.mod']
            assert module.isfile()
            expected_module = ('module github.com/wicanr2/dosgolem/wolongcresidual\n\n'
                               'go 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\n'
                               'replace github.com/wicanr2/dosgolem => /golem\n').encode()
            assert archive.extractfile(module).read() == expected_module
            verified.append({'path': f'workplace/matching-decompilation/{slug}/results/' + path.name,
                             'sha256': sha(path), 'compiled_source_manifest_sha256': compiled,
                             'source_files_verified': len(records)})
    assert verified, 'No complete source snapshot matches the compiled manifest'
    return verified


def backlinks(repo, routines):
    markers = ('殘部群的C補證', '134-c-residual-restoration.md', '254-c-residual-faction-season.md')
    evidence = 'docs/re/134-c-residual-restoration.md'
    evidence_text = (repo / evidence).read_text(encoding='utf-8')
    checked = []
    for address, name in BACKLINKS:
        original = f'sub_{address:X}'
        assert original in routines and original in evidence_text
        text = (repo / name).read_text(encoding='utf-8')
        assert original in text and all(marker in text for marker in markers), name
        checked.append({'platform_module': 'dosv/KI.EXE', 'input_sha256': INPUT_SHA,
                        'ida_linear': address, 'function_ida_linear': address, 'level': 'proven',
                        'evidence': evidence, 'evidence_markers': [original],
                        'older_document': name, 'required_markers': list(markers),
                        'scope': '固定raw狀態的局部C補證；保留原推論等級與未驗證邊界。'})
    return checked


def primary(repo, out):
    assert sha(repo / 'workplace/orig/dosv/KI.EXE') == INPUT_SHA
    probe_path = repo / 'workplace/matching-decompilation/c-residual/ida/ida-probe.json'
    assert sha(probe_path) == PROBE_SHA
    assert sha(probe_path.parent / 'input.exe.i64') == DB_SHA
    source = read(probe_path)
    assert source['input_sha256'] == source['ida_input_sha256'] == INPUT_SHA
    assert source['tool'] == 'IDA Pro' and source['tool_version'] == '9.4'
    assert source['mz_header_size'] == 512 and source['ida_load_paragraph'] == 0x1000
    assert source['function_count'] == 739
    assert set(source['recovery_targets']) == set(ROOTS) and len(source['recovery_targets']) == 10
    for target in source['targets']:
        assert target['unrestored_direct_callees'] == [], hex(target['ida_linear'])
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    count = int.from_bytes(raw[6:8], 'little')
    table = int.from_bytes(raw[24:26], 'little')
    relocations = {512 + int.from_bytes(raw[table + i * 4 + 2:table + i * 4 + 4], 'little') * 16 +
                   int.from_bytes(raw[table + i * 4:table + i * 4 + 2], 'little') for i in range(count)}
    assert relocations == set(source['relocation_file_offsets'])
    assert len(source['far_call_sites']) == 2 and source['far_call_target'] == 0x20000
    assert 'fixed_dispatch_tables' not in source
    rows, routines, blocks = {}, {}, {}
    named_count = named_bytes = raw_count = raw_bytes = 0
    seen_relocations = set()

    def check_extent(start, end, file_hex, chunk_relocations, instructions):
        file_start = start - 0x10000 + 512
        original = bytes.fromhex(file_hex)
        assert len(original) == end - start
        assert original == raw[file_start:file_start + len(original)]
        applicable = sorted(x for x in relocations if file_start <= x < file_start + len(original))
        assert sorted(x['file_offset'] for x in chunk_relocations) == applicable
        for rel in chunk_relocations:
            at = rel['file_offset'] - file_start
            word = int.from_bytes(original[at:at + 2], 'little')
            assert word == rel['original_word'] == 0x1000 and rel['ida_word'] == 0x2000
            assert rel['file_offset'] not in seen_relocations
            seen_relocations.add(rel['file_offset'])
        for row in instructions:
            ea = row['ida_linear']
            data = bytes.fromhex(row['bytes'])
            assert row['file_offset'] == ea - 0x10000 + 512
            assert start <= ea and ea + len(data) <= end
            assert ea not in rows
            rows[ea] = row
        return len(instructions), len(original)

    for target in source['targets']:
        body = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert hashlib.sha256(body).hexdigest() == target['file_sha256']
        assert target['name'] not in routines
        routines[target['name']] = target['file_sha256']
        for chunk in target['chunks']:
            n, size = check_extent(chunk['start'], chunk['end'], chunk['file_bytes'],
                                   chunk['loader_relocations'], chunk['instructions'])
            named_count += n
            named_bytes += size
    for block in source['decoded_blocks']:
        assert RAW_BLOCKS[block['ida_linear']] == len(block['instructions'])
        assert hashlib.sha256(bytes.fromhex(block['file_bytes'])).hexdigest() == block['file_sha256']
        name = f'raw_entry_{block["ida_linear"]:X}'
        assert name not in blocks
        blocks[name] = block['file_sha256']
        n, size = check_extent(block['ida_linear'], block['end_ida_linear'], block['file_bytes'],
                               block['loader_relocations'], block['instructions'])
        raw_count += n
        raw_bytes += size
    assert (named_count, named_bytes, raw_count, raw_bytes) == (247, 619, 1, 2)
    assert len(source['decoded_blocks']) == 1 and len(seen_relocations) == 2
    ledger = read(repo / 'docs/re/c-recovery-status.json')
    old_entries = {x['ida_linear'] for x in ledger['functions'].values()}
    old_entries |= {x['ida_linear'] for x in ledger['code_blocks'].values()}
    terminals = {'retn', 'retf', 'iret', 'jmp'}
    for ea, row in rows.items():
        mn, operands = row['mnemonic'], row['operands']
        assert not (mn in ('call', 'jmp') and operands[0]['type'] not in (6, 7)), hex(ea)
        if operands and operands[0]['type'] == 7 and (mn in ('call', 'jmp', 'loop') or mn.startswith('j')):
            assert 0x10000 + operands[0]['addr'] in rows.keys() | old_entries, (hex(ea), operands[0])
        if operands and mn == 'call' and operands[0]['type'] == 6:
            assert bytes.fromhex(row['bytes'])[0] == 0x9a and 0x20000 in old_entries
        if mn not in terminals:
            assert ea + len(bytes.fromhex(row['bytes'])) in rows.keys() | old_entries, ('missing fallthrough', hex(ea))
    generator = load_module(repo, 'c_recovery_residual_generate')
    assert generator.PROBE_SHA == PROBE_SHA and generator.WANTED == set(ROOTS)
    assert generator.RAW_WANTED == RAW_BLOCKS and len(generator.MUTANT_STATEMENTS) == 10
    controls = {}
    for ea, (number, _) in generator.MUTANT_STATEMENTS.items():
        controls[number] = ea
    assert sorted(controls) == list(range(1, 11))
    names = {}
    for line in (repo / 'tools/c_recovery_residual_mutants.tsv').read_text().splitlines():
        number, group, name = line.split('\t')
        assert int(number) not in names
        names[int(number)] = (group, name)
    assert sorted(names) == list(range(1, 11))
    with tempfile.TemporaryDirectory(prefix='residual-regenerate-') as tmp:
        generator.generate(repo, probe_path, Path(tmp))
        generated = {}
        for name in ('residual.h', 'residual_generated.inc'):
            assert (Path(tmp) / name).read_bytes() == (repo / 'tools/c_recovery' / name).read_bytes(), name
            generated['tools/c_recovery/' + name] = sha(Path(tmp) / name)
    return source, routines, blocks, generated


def verify(repo, out, source_only=False):
    source, routines, blocks, generated = primary(repo, out)
    if source_only:
        print('10 named / 1 raw / 248 instructions / 621 bytes / clean generation: PASS')
        return
    records, compiled = source_manifest(repo, out)
    assert all(value is not None for value in (GROUPS, STAGES, RETURNS, INDEPENDENT_AUDITS)), 'Residual semantic matrix has not been finalized'
    reports = {opt: read(out / 'results' / (opt + '.json')) for opt in ('O0', 'O2')}
    assert reports['O0'] == reports['O2'], 'Optimization reports differ'
    report = reports['O2']
    assert report['passed'] and report['mismatch'] is None
    assert report['input_sha256'] == INPUT_SHA
    assert report['new_probe_sha256'] == PROBE_SHA
    assert report['new_routine_sha256'] == routines
    assert report['new_code_block_sha256'] == blocks
    assert report['source_manifest']['sha256'] == compiled
    assert report['cases'] > 0 and report['full_ram_plane_audits'] >= report['cases']
    assert report['groups'] == GROUPS and report['stage_audits'] == STAGES
    assert report['cases'] == sum(GROUPS.values())
    assert report['full_ram_plane_audits'] == sum(STAGES.values())
    assert report['palette_dac_audits'] == report['full_ram_plane_audits']
    assert report['indexed_content_audits'] == report['full_ram_plane_audits']
    assert report['stage_audits'] == STAGES
    assert report['independent_residual_audits'] == INDEPENDENT_AUDITS
    assert report['original_state_sha256'] == report['c_state_sha256']
    assert report['independent_raw_initialization'] and not report['cross_machine_snapshot_initialization']
    assert not report['c_machine_code_match']
    for target in source['targets']:
        assert report['entries_seen'].get(target['name'], 0) > 0, target['name']
    for block in source['decoded_blocks']:
        assert report['entries_seen'].get(f'raw_entry_{block["ida_linear"]:X}', 0) > 0
    builds = {opt: native_build(out / ('residual-' + opt), out / 'results' / ('buildinfo-' + opt + '.txt'), compiled, None, opt) for opt in reports}
    mutants = {}
    for line in (repo / 'tools/c_recovery_residual_mutants.tsv').read_text().splitlines():
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
            '', '_ram', '_planes', '_trace', '_ports', '_in', '_api', '_console', '_allocs'))
        log = (out / 'results' / ('mutant-' + number + '.log')).read_text()
        assert 'pass=false' in log and not any(word in log for word in ('panic:', 'SIGABRT', 'SIGSEGV', 'unsupported='))
        build = native_build(out / ('mutant-' + number), out / 'results' / ('buildinfo-mutant-' + number + '.txt'), compiled, int(number), 'O0')
        mutants[number] = {'group': group, 'name': name, 'exit_status': 1, 'rejection_kind': 'state-mismatch', 'receipt_sha256': sha(path), **build}
    assert set(mutants) == {str(number) for number in range(1, 11)}
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before == (out / 'results/golem-source-after.sha256').read_bytes()
    for line in before.decode().splitlines():
        digest, name = line.split(None, 1)
        assert sha(Path(name)) == digest
    assert source_manifest(repo, out) == (records, compiled)
    snapshots = source_snapshots(out, records, compiled, 'c-residual')
    result = {
        'schema': 'wolong-c-residual-verification-v1', 'status': 'semantic-conformed',
        'input_sha256': INPUT_SHA, 'ida_database_sha256': DB_SHA, 'probe_sha256': PROBE_SHA,
        'named_functions': 10, 'raw_fragments': 1, 'original_instruction_count': 248, 'original_instruction_bytes': 621,
        'new_routine_sha256': routines, 'new_code_block_sha256': blocks,
        'source_sha256': records, 'compiled_source_manifest_sha256': compiled,
        'native_builds': builds, 'cases_per_optimization': report['cases'],
        'whole_device_stage_audits_per_optimization': report['full_ram_plane_audits'],
        'groups': report['groups'], 'stage_audits': report['stage_audits'], 'entries_seen': report['entries_seen'],
        'receipt_sha256': {opt: sha(out / 'results' / (opt + '.json')) for opt in reports},
        'negative_controls_rejected': len(mutants), 'mutants': mutants,
        'exact_clean_regeneration': True, 'clean_regenerated_files': generated,
        'compiled_sources_unchanged_after_execution': True,
        'compiled_source_snapshots': snapshots,
        'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
        'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(),
        'tool_versions': (out / 'results/tool-versions.txt').read_text(),
        
        'independent_audits': INDEPENDENT_AUDITS,
        'scope': report['scope'], 'platform': report['platform'], 'c_machine_code_match': False,
        'verification_tools_sha256': {name: sha(repo / name) for name in (
            'tools/c_recovery_residual.sh', 'tools/c_recovery_residual_container.sh', 'tools/c_recovery_residual_verify.py',
            'tools/c_recovery_residual_generate.py', 'tools/c_recovery_residual_mutants.tsv')},
        'backlink_status': backlinks(repo, routines),
    }
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(f'10 named / 1 raw / {report["cases"]} cases / {len(mutants)} state mismatches: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    verify(args.repo, args.output)
