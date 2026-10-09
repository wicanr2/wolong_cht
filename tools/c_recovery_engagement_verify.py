#!/usr/bin/env python3
"""核對交戰／戰術234函式與12raw的原始來源、裝置對拍與錯版反例。"""
import argparse
import hashlib
import importlib.util
import json
import re
import shlex
import shutil
import struct
import tempfile
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PROBE_SHA = '32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9'
DATABASE_SHA = 'cc4bf0da038a769fc7c12f8dcc47e884c1eb25ce959ceed0f86c3fb44182308f'
EVIDENCE = 'docs/re/131-c-engagement-tactical-restoration.md'
SPEC = 'docs/spec/251-c-engagement-tactical.md'
# 案例工廠已核對：4本體＋240控制＋48介面＋152交戰＋36存檔／世界介面＋268戰鬥。
GROUPS = {
    'combat-dispatch': 16, 'direction': 20, 'duel': 12, 'engagement': 4,
    'field': 16, 'front-helper': 40, 'grid': 4, 'melee': 16, 'patch': 28,
    'projectile': 36, 'projectile-render': 8, 'projectile-spawn': 24,
    'projection': 8, 'ranged': 24, 'save': 16, 'script': 96, 'settings': 28,
    'siege': 20, 'structure-collapse': 8, 'structure-ui': 12, 'terrain': 76,
    'tile-render': 12, 'ui-archive': 24, 'unit': 80, 'unit-render': 8,
    'unit-reposition': 32, 'unit-status': 16, 'unit-swap': 8,
    'vertical-move': 32, 'vertical-route': 4, 'world-ui': 20,
}
STAGES = {'engagement': 748, 'warmup-pause': 508, 'world-warmup': 236}
RETURNS = {'normal': 1492, 'nonlocal': 0}
TACTICAL_FRAME_AUDITS = 256
BACKLINKS = [
    (0x14A7B, 'docs/spec/105-encounter-goes-straight-to-battle.md', ('C交戰／戰術補證', 're/131', 'spec/251')),
    (0x14A7B, 'docs/spec/187-only-the-loser-is-judged.md', ('C交戰／戰術補證', 're/131', 'spec/251')),
    (0x14F8A, 'docs/spec/191-garrison-leader-is-general-127.md', ('C交戰／戰術補證', 're/131', 'spec/251')),
    (0x14A7B, 'docs/re/130-c-battle-outcome-restoration.md', ('C交戰／戰術補證', 're/131', 'spec/251')),
]
MUTANTS = {
    1: ('field', 'standoff-mask'), 2: ('siege', 'standoff-mask'),
    3: ('terrain', 'terrain-limit'), 4: ('patch', 'timer-live-opcode'),
    5: ('patch', 'unit-x-cutoff'), 6: ('patch', 'unit-y-cutoff'),
    7: ('patch', 'unit-dl-live'), 8: ('patch', 'target-live-compare'),
    9: ('patch', 'target-live-opcode'), 10: ('patch', 'path-live-opcode'),
    11: ('script', 'script-dispatch'), 12: ('script', 'direction-dispatch'),
    13: ('unit', 'unit-command'), 14: ('unit', 'unit-order'),
    15: ('settings', 'settings-dispatch'), 16: ('direction', 'direction-table'),
    17: ('projection', 'signed-shift'), 18: ('grid', 'grid-clear-count'),
    19: ('settings', 'six-button-count'), 20: ('direction', 'negative-x'),
}
NATIVE_TAGS = 'matching_engagement,matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga'


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def json_file(path):
    return json.loads(path.read_text(encoding='utf-8'))


def positive_counts(values):
    assert isinstance(values, dict) and values
    assert all(isinstance(k, str) and type(v) is int and v > 0 for k, v in values.items())


def digest(value):
    assert isinstance(value, str) and re.fullmatch('[0-9a-f]{64}', value), value
    return value


def native_build(binary, info_path, compiled, mutation, optimize):
    """核對ELF內嵌Go建置資料與另存輸出，不執行受測程式。"""
    data = binary.read_bytes()
    assert data[:6] == b'\x7fELF\x02\x01', binary
    shoff = struct.unpack_from('<Q', data, 40)[0]
    entsize, count, names_index = struct.unpack_from('<HHH', data, 58)
    assert entsize >= 64 and 0 < names_index < count and shoff + entsize * count <= len(data)
    sections = [struct.unpack_from('<IIQQQQIIQQ', data, shoff + i * entsize) for i in range(count)]
    ns = sections[names_index]
    names = data[ns[4]:ns[4] + ns[5]]
    entries = [s for s in sections if names[s[0]:].split(b'\0', 1)[0] == b'.go.buildinfo']
    assert len(entries) == 1
    section = entries[0]
    blob = data[section[4]:section[4] + section[5]]
    assert blob[:14] == b'\xff Go buildinf:' and blob[14:16] == b'\x08\x02'

    def string_at(at):
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

    version, at = string_at(32)
    module, _ = string_at(at)
    assert version == b'go1.26.7' and len(module) >= 33 and module[-17] == 10
    lines = module[16:-16].decode().splitlines()
    assert 'path\tgithub.com/wicanr2/dosgolem/wolongcengagement' in lines
    saved = info_path.read_text(encoding='utf-8').splitlines()
    assert saved[0].endswith(': go1.26.7')
    assert Path(saved[0].split(':', 1)[0]).name == binary.name

    def flags_from(items):
        prefix = 'build\tCGO_CFLAGS='
        values = [x.removeprefix('\t')[len(prefix):] for x in items if x.removeprefix('\t').startswith(prefix)]
        assert len(values) == 1
        quoted = shlex.split(values[0])
        assert len(quoted) == 1
        return shlex.split(quoted[0])

    flags = flags_from(lines)
    assert flags == flags_from(saved[1:])
    expected_define = '-DKI_ENGAGEMENT_SOURCE_DIGEST=0x' + compiled[:16]
    assert [v for v in flags if v.startswith('-DKI_ENGAGEMENT_SOURCE_DIGEST=')] == [expected_define]
    assert [v for v in flags if v.startswith('-DKI_ENGAGEMENT_MUTATION=')] == ([] if mutation is None else [f'-DKI_ENGAGEMENT_MUTATION={mutation}'])
    assert [v for v in flags if re.fullmatch('-O[0-9sgz]', v)] == ['-' + optimize]
    assert 'build\t-tags=' + NATIVE_TAGS in lines and 'build\tCGO_ENABLED=1' in lines
    return {'c_binary_sha256': sha(binary), 'buildinfo_sha256': sha(info_path),
            'go_version': version.decode(), 'cgo_flags': flags}


def primary_source(repo, out):
    """以固定IDA圖、原始file bytes與MZ載入規則驗證來源。"""
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == EXPECTED and len(raw) == 67099
    probe = out / 'ida-closed-v5/ida-probe.json'
    p = json_file(probe)
    assert sha(probe) == PROBE_SHA
    assert sha(out / 'ida-closed-v5/input.exe.i64') == DATABASE_SHA
    assert p['tool'] == 'IDA Pro' and p['tool_version'] == '9.4'
    assert p['input_sha256'] == p['ida_input_sha256'] == EXPECTED
    assert p['mz_header_size'] == 512 and p['ida_load_paragraph'] == 0x1000
    assert set(p['navigation_only_targets']) == {0x125A3}
    targets = set(p['recovery_targets'])
    assert len(targets) == len(p['recovery_targets']) == 234
    selected = [t for t in p['targets'] if t['ida_linear'] in targets]
    assert len(selected) == 234 and len(p['decoded_blocks']) == 12
    assert set(t['ida_linear'] for t in p['targets']) == targets | {0x125A3}
    relocation_count = struct.unpack_from('<H', raw, 6)[0]
    relocation_table = struct.unpack_from('<H', raw, 24)[0]
    relocation_offsets = set()
    for i in range(relocation_count):
        offset, segment = struct.unpack_from('<HH', raw, relocation_table + i * 4)
        relocation_offsets.add(512 + segment * 16 + offset)
    assert relocation_offsets == set(p['relocation_file_offsets'])
    rows, routines, blocks, relocations = {}, {}, {}, {}
    named_count = named_bytes = raw_count = raw_bytes = extent = 0

    def check_extent(start, end, file_hex, loaded_hex, chunk_relocations, instructions):
        assert 0x10000 <= start < end <= 0x20000
        file_start = start - 0x10000 + 512
        original = bytes.fromhex(file_hex)
        loaded = bytearray(original)
        assert len(original) == end - start
        assert original == raw[file_start:file_start + len(original)]
        applicable = sorted(x for x in relocation_offsets if file_start <= x < file_start + len(original))
        assert sorted(x['file_offset'] for x in chunk_relocations) == applicable
        for rel in chunk_relocations:
            at = rel['file_offset'] - file_start
            assert at + 2 <= len(original)
            word = struct.unpack_from('<H', original, at)[0]
            assert word == rel['original_word']
            assert (word + 0x1000) & 0xffff == rel['ida_word']
            struct.pack_into('<H', loaded, at, rel['ida_word'])
            assert rel['file_offset'] not in relocations
            relocations[rel['file_offset']] = rel
        if loaded_hex is not None:
            assert loaded == bytes.fromhex(loaded_hex)
        occupied = set()
        for row in instructions:
            ea = row['ida_linear']
            position = ea - start
            data = bytes.fromhex(row['bytes'])
            assert row['file_offset'] == ea - 0x10000 + 512
            assert start <= ea and ea + len(data) <= end
            assert data == loaded[position:position + len(data)]
            assert ea not in rows
            occupied_now = set(range(ea, ea + len(data)))
            assert not occupied & occupied_now
            occupied |= occupied_now
            rows[ea] = row
        return len(instructions), sum(len(bytes.fromhex(x['bytes'])) for x in instructions)

    for target in selected:
        original = b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])
        assert hashlib.sha256(original).hexdigest() == target['file_sha256']
        assert target['name'] not in routines
        routines[target['name']] = target['file_sha256']
        for chunk in target['chunks']:
            n, size = check_extent(chunk['start'], chunk['end'], chunk['file_bytes'], chunk['bytes'],
                                   chunk['loader_relocations'], chunk['instructions'])
            named_count += n
            named_bytes += size
            extent += chunk['end'] - chunk['start']
    for block in p['decoded_blocks']:
        original = bytes.fromhex(block['file_bytes'])
        assert hashlib.sha256(original).hexdigest() == block['file_sha256']
        name = f'raw_entry_{block["ida_linear"]:X}'
        assert name not in blocks
        blocks[name] = block['file_sha256']
        n, size = check_extent(block['ida_linear'], block['end_ida_linear'], block['file_bytes'], None,
                               block['loader_relocations'], block['instructions'])
        raw_count += n
        raw_bytes += size
    assert (named_count, named_bytes, extent, raw_count, raw_bytes) == (7080, 16696, 16706, 737, 1867)
    assert len(rows) == 7817 and named_bytes + raw_bytes == 18563
    assert len(relocations) == 24 and p['unresolved_call_boundaries'] == []
    ledger = json_file(repo / 'docs/re/c-recovery-status.json')
    old_entries = {x['ida_linear'] for x in ledger['functions'].values()}
    old_entries |= {x['ida_linear'] for x in ledger['code_blocks'].values()}
    terminals = {'retn', 'retf', 'iret', 'jmp'}
    indirect = {}
    for ea, row in rows.items():
        mn, operands = row['mnemonic'], row['operands']
        if mn in ('call', 'jmp') and operands[0]['type'] not in (6, 7):
            indirect[ea] = mn
        if operands and operands[0]['type'] == 7 and (mn in ('call', 'jmp', 'loop') or mn.startswith('j')):
            assert 0x10000 + operands[0]['addr'] in rows.keys() | old_entries, (hex(ea), operands[0])
        if operands and mn == 'call' and operands[0]['type'] == 6:
            assert bytes.fromhex(row['bytes'])[0] == 0x9a and 0x20000 in old_entries
        if mn not in terminals:
            assert ea + len(bytes.fromhex(row['bytes'])) in rows, ('missing fallthrough', hex(ea))
    tables = p['fixed_dispatch_tables']
    assert len(tables) == 9 and len(indirect) == 9
    assert list(indirect.values()).count('call') == 8 and list(indirect.values()).count('jmp') == 1
    assert set(indirect) == {t['call_or_jump_ida_linear'] for t in tables}
    assert {x['ida_linear'] for x in p['indirect_calls']} == {ea for ea, kind in indirect.items() if kind == 'call'}
    for table in tables:
        at = table['table_ida_linear'] - 0x10000 + 512
        words = table['words']
        assert len(words) == table['entries']
        assert list(struct.unpack_from('<' + 'H' * len(words), raw, at)) == words
        assert all(0x10000 + value in rows.keys() | old_entries for value in words)
    return probe, p, rows, routines, blocks


def clean_regeneration(repo, probe):
    spec = importlib.util.spec_from_file_location('engagement_generator', repo / 'tools/c_recovery_engagement_generate.py')
    gen = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(gen)
    with tempfile.TemporaryDirectory(prefix='engagement-repro-') as temporary:
        root = Path(temporary)
        output = root / 'tools/c_recovery'
        output.mkdir(parents=True)
        shutil.copyfile(repo / 'tools/c_recovery_display_generate.py', root / 'tools/c_recovery_display_generate.py')
        gen.generate(root, probe)
        names = ['engagement.h', 'engagement_generated.inc']
        for name in names:
            assert (output / name).read_bytes() == (repo / 'tools/c_recovery' / name).read_bytes(), name
    return names


def assembly_source(repo, rows, p, publication):
    code = repo / 'docs/re/c-engagement-code.json'
    coverage = json_file(code)
    assembly = coverage['assembler_verification']
    assert coverage['schema'] == 'wolong-matching-code-supplement-v1'
    assert coverage['input_sha256'] == EXPECTED and coverage['probe_sha256'] == PROBE_SHA
    assert coverage['database_sha256'] == DATABASE_SHA and coverage['evidence'] == EVIDENCE
    assert (coverage['named_instruction_count'], coverage['named_instruction_bytes'], coverage['named_chunk_extent_bytes'],
            coverage['raw_instruction_count'], coverage['raw_instruction_bytes']) == (7080, 16696, 16706, 737, 1867)
    assert set(coverage['recovery_targets']) == set(p['recovery_targets'])
    assert set(coverage['navigation_only_targets']) == {0x125A3}
    assert len(coverage['original_functions']) == 234 and len(coverage['original_code_blocks']) == 12
    assert assembly['status'] == 'byte-exact'
    assert assembly['input_instruction_rows'] == assembly['unique_input_instructions'] == 7817
    assert assembly['covered_input_instructions'] == 7786
    assert assembly['instructions'] == 31 and assembly['bytes'] == 84 and assembly['instruction_byte_fallbacks'] == 0
    independent = assembly['independent_reassembly']
    assert independent['status'] == 'byte-exact' and independent['instructions'] == 7817 and independent['bytes'] == 18563
    assert independent['trials'][-1]['remaining'] == 0
    assert assembly['sources_before'] == assembly['sources_after']
    for name, value in assembly['sources_before'].items():
        assert sha(repo / name) == value, name
    actual = coverage['instruction_coverage']
    assert len(actual) == 7817 and {x['ida_linear'] for x in actual} == set(rows)
    assert all(x['independent_native_assembly_match'] for x in actual)
    assert sum(not x['fixed_prior_coverage'] for x in actual) == 31
    assert len(coverage['instructions']) == 31 and sum(x['file_end'] - x['file_start'] for x in coverage['instructions']) == 84
    assert len(coverage['retired_instructions']) == 12
    assert sum(x['file_end'] - x['file_start'] for x in coverage['retired_instructions']) == 28
    assert coverage['retirement_summary']['predicted_total_instructions'] == 25097
    assert coverage['retirement_summary']['predicted_total_instruction_bytes'] == 57101
    assert len(coverage['loader_relocation_normalization']) == 24
    assert len(coverage['preserved_inline_data']) == 1
    inline = coverage['preserved_inline_data'][0]
    assert inline['instruction'] is False and len(bytes.fromhex(inline['bytes'])) == 10
    assert inline['ida_start'] == 0x1A59C and inline['ida_end'] == 0x1A5A6
    for row in actual:
        original = rows[row['ida_linear']]
        assert row['file_start'] == original['file_offset']
        assert row['file_end'] - row['file_start'] == len(bytes.fromhex(original['bytes']))
        assert row['ida_loaded_bytes'] == original['bytes']
        loaded = bytearray.fromhex(row['bytes'])
        for relocation in row['loader_relocations']:
            at = relocation['file_offset'] - row['file_start']
            assert struct.unpack_from('<H', loaded, at)[0] == relocation['original_word']
            struct.pack_into('<H', loaded, at, relocation['ida_word'])
        assert loaded == bytes.fromhex(original['bytes'])
    # The supplement binds S/linker. The manifest binds this supplement, avoiding a hash cycle.
    assert 'published_manifest_sha256' not in assembly
    if publication:
        assert assembly['published_code_source'] == 'tools/c_recovery/KI.code.S'
        assert assembly['published_linker'] == 'tools/c_recovery/KI.code.ld'
        assert assembly['published_code_source_sha256'] == sha(repo / assembly['published_code_source'])
        assert assembly['published_linker_sha256'] == sha(repo / assembly['published_linker'])
        manifest = json_file(repo / 'docs/re/assembly-code-record.json')
        assert manifest['input_sha256'] == EXPECTED and manifest['input_size'] == 67099
        assert manifest['instruction_count'] == 25097 and manifest['instruction_bytes'] == 57101
        supplement = [x for x in manifest['additional_supplements'] if x['path'] == 'docs/re/c-engagement-code.json']
        assert len(supplement) == 1 and supplement[0]['sha256'] == sha(code)
        assert supplement[0]['instructions'] == 31 and supplement[0]['bytes'] == 84
        assert supplement[0]['retired_instructions'] == 12 and supplement[0]['retired_bytes'] == 28
        build = json_file(repo / 'workplace/matching-decompilation/assembly/code-record/code-record-verification.json')
        assert build['whole_file_exact'] and build['input_sha256'] == build['rebuilt_sha256'] == EXPECTED
        assert build['rebuilt_size'] == 67099 and build['instruction_count'] == 25097
        assert build['compiled_instruction_bytes'] == 57101 and build['imported_instruction_bytes'] == 0
        assert build['manifest_sha256'] == sha(repo / 'docs/re/assembly-code-record.json')
        assert build['source_sha256'] == manifest['source_sha256'] == assembly['published_code_source_sha256']
        assert build['linker_sha256'] == manifest['linker_sha256'] == assembly['published_linker_sha256']
        assert len(build['negative_controls']) == 3 and all(x['rejected'] for x in build['negative_controls'])
    return sha(code)


def source_manifest(repo, out):
    path = out / 'results/c-source.sha256'
    source = {}
    for line in path.read_text(encoding='utf-8').splitlines():
        value, name = line.split(None, 1)
        assert name.startswith('/repo/')
        name = name.removeprefix('/repo/')
        assert '..' not in Path(name).parts and not Path(name).is_absolute()
        digest(value)
        assert name not in source or source[name] == value
        assert sha(repo / name) == value, name
        source[name] = value
    required = {
        'tools/c_recovery/engagement.c', 'tools/c_recovery/engagement.h',
        'tools/c_recovery/engagement_fixture.h', 'tools/c_recovery/engagement_generated.inc',
        'tools/c_recovery_engagement.go', 'tools/c_recovery_engagement_generate.py',
        'tools/c_recovery_engagement_data.go', 'tools/c_recovery_engagement_control_data.go',
        'tools/c_recovery_engagement_ui_data.go', 'tools/c_recovery_engagement_front_data.go',
        'tools/c_recovery_engagement_extra_data.go', 'tools/c_recovery_engagement_platform.go',
        'tools/c_recovery_engagement_combat_data.go', 'tools/c_recovery_engagement_save.go',
        'tools/c_recovery_display_generate.py', 'tools/c_recovery_vga_bus.go',
        'tools/c_recovery_glyph_platform.go', 'tools/c_recovery_input_platform.go',
        'tools/c_recovery_outcome_data.go', 'tools/c_recovery_route_data.go',
        'tools/c_recovery_interaction_data.go', 'tools/c_recovery_main_data.go',
        'tools/c_recovery_strategy_data.go',
    }
    assert required <= source.keys(), sorted(required - source.keys())
    assert {str(x.relative_to(repo)) for x in (repo / 'tools/c_recovery').iterdir()
            if x.is_file() and x.suffix in ('.c', '.h', '.inc')} <= source.keys()
    compiled = sha(path)
    assert compiled == (out / 'results/compiled-source-digest.txt').read_text().strip()
    return source, compiled


def verify_rng(report):
    assert report['rng_initial_state_sha256']
    for seed, value in report['rng_initial_state_sha256'].items():
        assert re.fullmatch('[0-9A-F]{2}', seed)
        number = int(seed, 16)
        initial = bytes([number ^ 0xA5, number]) + bytes((i * 73 + i // 2 + number) & 255 for i in range(256))
        assert len(initial) == 258 and hashlib.sha256(initial).hexdigest() == value


def verify_save_and_restore(repo, report):
    original = repo / 'workplace/orig/dosv/SAVE.DAT'
    assert original.is_file() and original.stat().st_size == 88832
    assert report['save_input_sha256'] == sha(original)
    assert report['separate_save_scratch'] is True
    receipts = report['save_file_audits']
    assert isinstance(receipts, list) and len(receipts) == 16
    pairs = set()
    for receipt in receipts:
        scenario, slot = receipt['scenario'], receipt['slot']
        assert type(scenario) is int and type(slot) is int
        assert 0 <= scenario < 4 and 0 <= slot < 4
        assert (scenario, slot) not in pairs
        pairs.add((scenario, slot))
        assert receipt['bytes'] == 88832
        assert digest(receipt['original_sha256']) == digest(receipt['c_sha256'])
        assert receipt['independent_expected_match'] is True
        assert receipt['unchanged_gaps_and_other_slots'] is True
        assert receipt['separate_scratch'] is True
    assert pairs == {(scenario, slot) for scenario in range(4) for slot in range(4)}
    restored = report['mmap_restore_audits']
    assert isinstance(restored, dict) and set(restored) == {'original', 'c'}
    assert type(restored['original']) is int and type(restored['c']) is int
    assert restored['original'] == restored['c'] and restored['original'] > 0


def backlinks(repo, routines, blocks):
    evidence = (repo / EVIDENCE).read_text(encoding='utf-8')
    if '**狀態：CONFORMED' not in evidence:
        return 'pending-publication'
    assert BACKLINKS is not None, 'Final publication backlink matrix is not reviewed'
    assert '**狀態：CONFORMED' in (repo / SPEC).read_text(encoding='utf-8')
    ledger = json_file(repo / 'docs/re/c-recovery-status.json')
    assert ledger['input_sha256'] == EXPECTED
    assert len(ledger['functions']) == 665 and len(ledger['code_blocks']) == 62
    for name, value in routines.items():
        entry = ledger['functions'][name]
        assert entry['routine_sha256'] == value and entry['status'] == 'semantic-conformed'
        assert entry['spec'] == SPEC and entry['evidence'] == EVIDENCE and not entry['c_machine_code_match']
    for name, value in blocks.items():
        entry = ledger['code_blocks']['0x' + name.removeprefix('raw_entry_')]
        assert entry['routine_sha256'] == value and entry['status'] == 'semantic-conformed'
        assert entry['evidence'] == EVIDENCE and not entry['c_machine_code_match']
    for address, document, markers in BACKLINKS:
        text = (repo / document).read_text(encoding='utf-8')
        assert all(marker in text for marker in markers)
        matches = [x for x in ledger['resolution_backlinks']
                   if x['older_document'] == document and x['ida_linear'] == address]
        assert len(matches) == 1
        row = matches[0]
        assert row['evidence'] == EVIDENCE and row['input_sha256'] == EXPECTED
        assert row['platform_module'] == 'dosv/KI.EXE' and row['level'] == 'proven'
        assert all(marker in text for marker in row['required_markers'])
        assert all(marker in evidence for marker in row['evidence_markers'])
    return 'verified'


def verify(repo, out, source_only=False):
    probe, p, rows, routines, blocks = primary_source(repo, out)
    regenerated = clean_regeneration(repo, probe)
    assembly_sha = assembly_source(repo, rows, p, publication=not source_only)
    if source_only:
        print(f'234 named / 12 raw / 7817 instructions / 18563 bytes: source PASS; clean regeneration {regenerated}')
        return
    if GROUPS is None or STAGES is None or RETURNS is None or TACTICAL_FRAME_AUDITS is None:
        raise RuntimeError('Final engagement matrix, stage counts and returns are pending; publication is refused')
    positive_counts(GROUPS)
    positive_counts(STAGES)
    count, stages = sum(GROUPS.values()), sum(STAGES.values())
    assert STAGES['engagement'] == count and set(STAGES) <= {'engagement', 'warmup-pause', 'world-warmup'}
    assert set(RETURNS) == {'normal', 'nonlocal'}
    assert all(type(x) is int and x >= 0 for x in RETURNS.values()) and sum(RETURNS.values()) == stages
    paused = STAGES.get('warmup-pause', 0)
    assert RETURNS['normal'] >= paused
    assert len(MUTANTS) == 20 and set(MUTANTS) == set(range(1, 21))
    source, compiled = source_manifest(repo, out)
    manifest = {'path': 'workplace/matching-decompilation/c-engagement/results/c-source.sha256', 'sha256': compiled}
    reports, receipts, native = {}, {}, {}
    for level in ('O0', 'O2'):
        path = out / 'results' / f'{level}.json'
        r = json_file(path)
        assert r['schema'] == 'wolong-c-engagement-parity-v1' and r['passed'] and r['mismatch'] is None
        assert r['input_sha256'] == EXPECTED and r['cases'] == count and r['groups'] == GROUPS
        assert r['full_ram_plane_audits'] == r['indexed_content_audits'] == r['palette_dac_audits'] == stages
        assert r['independent_outcome_audits'] >= stages and r['stage_audits'] == STAGES
        assert r['new_routine_sha256'] == routines and r['new_code_block_sha256'] == blocks
        assert r['source_manifest'] == manifest and r['new_probe_sha256'] == PROBE_SHA
        positive_counts(r['entries_seen'])
        assert all(r['entries_seen'].get(name, 0) > 0 for name in routines.keys() | blocks.keys())
        assert digest(r['original_state_sha256']) == digest(r['c_state_sha256']) and not r['c_machine_code_match']
        assert r['independent_raw_initialization'] and not r['cross_machine_snapshot_initialization']
        assert r['fixed_raw_rng_state']
        assert r['nonlocal_exit_cases'] == RETURNS['nonlocal'] and r['normal_return_cases'] == RETURNS['normal']
        assert r['normal_return_cases'] + r['nonlocal_exit_cases'] == stages
        # Tactical guest-frame recovery is distinct from the outer KiMainFixture transfer.
        assert r['actual_nonlocal_exit'] == (RETURNS['nonlocal'] > 0)
        assert r['tactical_frame_restore_audits'] == TACTICAL_FRAME_AUDITS and TACTICAL_FRAME_AUDITS > 0
        assert r['actual_tactical_frame_restore'] is True
        verify_rng(r)
        verify_save_and_restore(repo, r)
        assert r['original_missing_fonts'] == r['c_missing_fonts'] == 0
        assert r['original_font_misses'] == r['c_font_misses'] == {}
        assert r['original_font_calls'] == r['c_font_calls']
        assert 'incomplete' in r['scope']
        reports[level], receipts[level] = r, sha(path)
        native[level] = native_build(out / f'engagement-{level}', out / 'results' / f'buildinfo-{level}.txt', compiled, None, level)
    assert receipts['O0'] == receipts['O2']
    names = {}
    for line in (out / 'results/mutant-names.tsv').read_text().splitlines():
        number, group, name = line.split('\t')
        assert int(number) not in names
        names[int(number)] = (group, name)
    assert names == MUTANTS
    controls = {}
    mismatch_fields = ('', '_ram', '_planes', '_trace', '_ports', '_in', '_api', '_palette_dac')
    for number, (group, name) in MUTANTS.items():
        path = out / 'results' / f'mutant-{number}.json'
        r = json_file(path)
        m = r['mismatch']
        assert (out / 'results' / f'mutant-{number}.exit').read_text().strip() == '1'
        assert r['schema'] == 'wolong-c-engagement-parity-v1' and not r['passed']
        assert r['input_sha256'] == EXPECTED and r['source_manifest'] == manifest
        assert r['new_probe_sha256'] == PROBE_SHA
        assert r['new_routine_sha256'] == routines and r['new_code_block_sha256'] == blocks
        assert r['groups'] == {group: r['cases']} and 0 < r['cases'] <= GROUPS[group]
        assert m['group'] == group and m['case'] == r['cases'] - 1 and m['stage'] in STAGES
        assert any(m['original' + key] != m['c' + key] for key in mismatch_fields)
        assert r['original_state_sha256'] == r['c_state_sha256']  # accepted preceding stages only
        assert r['normal_return_cases'] + r['nonlocal_exit_cases'] == sum(r['stage_audits'].values())
        assert r['full_ram_plane_audits'] == r['indexed_content_audits'] == r['palette_dac_audits'] == sum(r['stage_audits'].values())
        log = (out / 'results' / f'mutant-{number}.log').read_text(encoding='utf-8')
        assert 'pass=false' in log and not any(x in log for x in ('panic:', 'SIGABRT', 'SIGSEGV', 'unsupported='))
        build = native_build(out / f'mutant-{number}', out / 'results' / f'buildinfo-mutant-{number}.txt', compiled, number, 'O0')
        native[f'mutant-{number}'] = build
        controls[str(number)] = {'group': group, 'name': name, 'rejection_kind': 'state-mismatch',
                                 'exit_status': 1, 'first_rejected': r['cases'], 'stage': m['stage'],
                                 'receipt_sha256': sha(path), 'log_sha256': sha(out / 'results' / f'mutant-{number}.log'), **build}
    before = (out / 'results/golem-source-before.sha256').read_bytes()
    assert before and before == (out / 'results/golem-source-after.sha256').read_bytes()
    # The current compiled-source tree must still match the pre-build manifest.
    final_source, final_digest = source_manifest(repo, out)
    assert final_source == source and final_digest == compiled
    result = {
        'schema': 'wolong-c-engagement-verification-v1', 'status': 'semantic-conformed', 'input_sha256': EXPECTED,
        'ida_database_sha256': DATABASE_SHA, 'ida_probe_sha256': PROBE_SHA,
        'new_routine_sha256': routines, 'new_code_block_sha256': blocks,
        'original_instruction_count': 7817, 'original_instruction_bytes': 18563,
        'named_instruction_count': 7080, 'named_instruction_bytes': 16696, 'named_chunk_extent_bytes': 16706,
        'raw_instruction_count': 737, 'raw_instruction_bytes': 1867, 'source_sha256': source,
        'compiled_source_manifest_sha256': compiled, 'native_builds': native, 'cases_per_optimization': count,
        'whole_device_stage_audits_per_optimization': stages, 'stage_audits': STAGES, 'groups': GROUPS,
        'receipt_sha256': receipts, 'negative_controls_rejected': 20, 'state_mismatch_controls_rejected': 20,
        'mutants': controls, 'entries_seen': reports['O2']['entries_seen'],
        'independent_outcome_audits': reports['O2']['independent_outcome_audits'],
        'assembly_coverage_receipt_sha256': assembly_sha, 'exact_clean_regeneration': True,
        'clean_regenerated_files': regenerated, 'compiled_sources_unchanged_after_execution': True,
        'tool_versions': (out / 'results/tool-versions.txt').read_text(),
        'oracle_revision': (out / 'results/golem-revision.txt').read_text().strip(),
        'oracle_compiled_tree_sha256': hashlib.sha256(before).hexdigest(), 'c_machine_code_match': False,
        'scope': reports['O2']['scope'] + ' normal_return_cases includes paused warmup stages; returned_to_caller_stages excludes those pauses.',
        'platform': reports['O2']['platform'],
        'rng_initial_state_sha256': reports['O2']['rng_initial_state_sha256'],
        'rng_initialization': reports['O2']['rng_initialization'],
        'nonlocal_exit_cases': RETURNS['nonlocal'], 'normal_return_cases': RETURNS['normal'],
        'paused_stages': paused, 'returned_to_caller_stages': RETURNS['normal'] - paused,
        'return_count_scope': 'normal_return_cases 相容欄位包含 warmup-pause；returned_to_caller_stages 排除該暫停階段。',
        'tactical_frame_restore_audits': TACTICAL_FRAME_AUDITS, 'actual_tactical_frame_restore': True,
        'save_file_audits': reports['O2']['save_file_audits'],
        'separate_save_scratch': True, 'save_input_sha256': reports['O2']['save_input_sha256'],
        'mmap_restore_audits': reports['O2']['mmap_restore_audits'],
        'verification_tools_sha256': {name: sha(repo / name) for name in (
            'tools/c_recovery_engagement.sh', 'tools/c_recovery_engagement_container.sh', 'tools/c_recovery_engagement_verify.py',
            'tools/c_recovery_engagement_generate.py', 'tools/ida_engagement_probe.py', 'tools/engagement_code_supplement.py')},
    }
    go = repo / 'docs/re/c-engagement-go-verification.json'
    if go.exists():
        proof = json_file(go)
        assert proof['vet_passed'] and proof['test_passed'] and proof['cold_test_cache']
        assert proof['tested_packages'] == 39 and proof['cached_test_packages'] == 0
        assert proof['original_assets_readonly'] and proof['network'] == 'none'
        for name, value in proof['source_sha256'].items():
            assert sha(repo / name) == value
        result['go_verification_sha256'] = sha(go)
    result['backlink_status'] = backlinks(repo, routines, blocks)
    (out / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    print(f'234 named + 12 raw / 7817 instructions / {count} cases / {stages} stages / 20 state mismatches: PASS')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--source-only', action='store_true', help='Only primary bytes, CFG, clean generation and assembly source; no semantic publication')
    args = parser.parse_args()
    verify(args.repo, args.output, args.source_only)
