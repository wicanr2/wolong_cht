#!/usr/bin/env python3
"""從固定舊IDA證據抽取軍團17函式，獨立組譯並核對已發布指令來源。"""
from __future__ import annotations

import argparse
import copy
import hashlib
import importlib.util
import json
import os
import re
import struct
from pathlib import Path

INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PARENT = 'workplace/matching-decompilation/c-tick/census/ida-probe.json'
PARENT_SHA = '937f8738225ae0f9e5efc133245316dd075d845ffded97a22f5f211b4804f2a7'
DATABASE = 'workplace/matching-decompilation/c-tick/census/input.exe.i64'
DATABASE_SHA = '4f01452cc2df29a9306aa4b444d1d4125e3f27f7730fcfdadbefd9f9da3a4667'
IMAGE = 'sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e'
EVIDENCE = 'docs/re/132-c-army-update-restoration.md'
BASELINE = 'workplace/matching-decompilation/c-army/ledger-baseline-665.json'
BASELINE_SHA = '93116486e133048e7bb5100cb33b54bd37b9503a9136ae9b3e28823f539f0953'
BASELINE_PROVENANCE = 'workplace/matching-decompilation/c-army/ledger-baseline-665-provenance.json'
ROOTS = (0x125A3, 0x12600, 0x1264A, 0x12662, 0x126FF, 0x12708,
         0x127A2, 0x127F6, 0x12804, 0x12808, 0x12831, 0x12880,
         0x128F4, 0x12A7E, 0x142AB, 0x14300, 0x1562B)
FALLTHROUGHS = {(0x12706, 0x12708), (0x12807, 0x12808)}
EXTERNAL = {0x14325, 0x12BA8, 0x147BB, 0x102F5, 0x14A7B,
            0x14ADE, 0x1291A, 0x10CDE, 0x18810, 0x1563B}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load(path):
    return json.loads(path.read_text(encoding='utf-8'))


def write_json(path, value):
    assert path.parent.is_dir() and path.parent.stat().st_uid == os.getuid()
    assert not path.exists() or path.is_file() and path.stat().st_uid == os.getuid()
    path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')


def validate_external_ledger(repo, baseline):
    """只驗外部callee原bytes身分及狀態；本輪函式日後發布不自證。"""
    current = load(repo / 'docs/re/c-recovery-status.json')
    assert current['input_sha256'] == INPUT_SHA
    for name, entry in baseline['functions'].items():
        if entry['ida_linear'] not in EXTERNAL:
            continue
        now = current['functions'][name]
        assert now['ida_linear'] == entry['ida_linear']
        assert now['routine_sha256'] == entry['routine_sha256']
        assert now['status'] == entry['status'] == 'semantic-conformed'


def prepare_source(repo):
    assert sha(repo / PARENT) == PARENT_SHA
    assert sha(repo / DATABASE) == DATABASE_SHA
    parent = load(repo / PARENT)
    assert parent['tool'] == 'IDA Pro' and parent['tool_version'] == '9.4'
    assert parent['input_sha256'] == parent['ida_input_sha256'] == INPUT_SHA
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert len(raw) == 67099 and hashlib.sha256(raw).hexdigest() == INPUT_SHA
    header = struct.unpack_from('<H', raw, 8)[0] * 16
    count, table = struct.unpack_from('<H', raw, 6)[0], struct.unpack_from('<H', raw, 24)[0]
    relocations = [header + off + seg * 16 for off, seg in
                   (struct.unpack_from('<HH', raw, table + n * 4) for n in range(count))]
    assert header == parent['mz_header_size'] == 512
    assert parent['ida_load_paragraph'] == 0x1000
    assert relocations == parent['relocation_file_offsets']
    targets = sorted((copy.deepcopy(t) for t in parent['targets'] if t['ida_linear'] in ROOTS),
                     key=lambda t: t['ida_linear'])
    assert tuple(t['ida_linear'] for t in targets) == ROOTS
    rows, owners, cross_falls = {}, {}, []
    for target in targets:
        original = b''.join(bytes.fromhex(ch['file_bytes']) for ch in target['chunks'])
        assert hashlib.sha256(original).hexdigest() == target['file_sha256']
        assert isinstance(target['xrefs_to'], list)
        for chunk in target['chunks']:
            start, end = chunk['start'], chunk['end']
            file_start = start - 0x10000 + header
            data, loaded = bytes.fromhex(chunk['file_bytes']), bytes.fromhex(chunk['bytes'])
            assert len(data) == end - start and raw[file_start:file_start + len(data)] == data
            expected = [r for r in relocations if file_start <= r < file_start + len(data)]
            assert sorted(r['file_offset'] for r in chunk['loader_relocations']) == sorted(expected)
            reconstructed = bytearray(data)
            for relocation in chunk['loader_relocations']:
                at = relocation['file_offset'] - file_start
                word = struct.unpack_from('<H', data, at)[0]
                assert word == relocation['original_word']
                assert (word + 0x1000) & 0xFFFF == relocation['ida_word']
                struct.pack_into('<H', reconstructed, at, relocation['ida_word'])
            assert bytes(reconstructed) == loaded
            assert expected == [], 'This seventeen-function scope has no MZ relocations'
            cursor = start
            for instruction in chunk['instructions']:
                ea = instruction['ida_linear']
                encoded = bytes.fromhex(instruction['bytes'])
                at = ea - 0x10000 + header
                assert ea == cursor and at == instruction['file_offset']
                assert raw[at:at + len(encoded)] == encoded
                assert ea not in rows and encoded == loaded[ea - start:ea - start + len(encoded)]
                rows[ea], owners[ea] = instruction, target['name']
                cursor += len(encoded)
            assert cursor == end
    assert len(rows) == 434 and sum(len(bytes.fromhex(i['bytes'])) for i in rows.values()) == 1068
    ledger_path = repo / BASELINE
    assert sha(ledger_path) == BASELINE_SHA
    provenance = load(repo / BASELINE_PROVENANCE)
    assert provenance['source_path'] == 'docs/re/c-recovery-status.json'
    assert provenance['source_date'] == '2026-10-10'
    assert provenance['source_sha256'] == BASELINE_SHA
    hashfile = ledger_path.with_suffix('.sha256')
    assert hashfile.read_text().strip() == BASELINE_SHA + '  ledger-baseline-665.json'
    ledger = load(ledger_path)
    assert ledger['input_sha256'] == INPUT_SHA
    assert len(ledger['functions']) == 665 and len(ledger['code_blocks']) == 62
    validate_external_ledger(repo, ledger)
    confirmed = {row['ida_linear']: {'original_name': name, **row}
                 for name, row in ledger['functions'].items() if row['status'] == 'semantic-conformed'}
    confirmed.update({row['ida_linear']: {'ledger_key': key, **row}
                      for key, row in ledger['code_blocks'].items() if row['status'] == 'semantic-conformed'})
    assert not set(ROOTS) & confirmed.keys()
    controls, external, indirect = [], {}, []
    for ea, instruction in sorted(rows.items()):
        mnemonic, operands = instruction['mnemonic'], instruction['operands']
        if mnemonic.startswith('j') or mnemonic in ('call', 'loop', 'loope', 'loopne'):
            op = operands[0]
            if op['type'] != 7:
                indirect.append({'ida_linear': ea, 'assembly': instruction['assembly']})
                continue
            target = 0x10000 + op['addr']
            assert target in rows or target in confirmed, (hex(ea), hex(target))
            controls.append({'ida_linear': ea, 'target_ida_linear': target,
                             'mnemonic': mnemonic, 'original_assembly': instruction['assembly'],
                             'target_kind': 'scope-instruction' if target in rows else 'confirmed-callee'})
            if target not in rows:
                external[target] = confirmed[target]
        if mnemonic not in ('retn', 'retf', 'iret', 'jmp'):
            successor = ea + len(bytes.fromhex(instruction['bytes']))
            assert successor in rows, ('missing fallthrough', hex(ea), hex(successor))
            if owners[ea] != owners[successor]:
                cross_falls.append({'from_ida_linear': ea, 'to_ida_linear': successor,
                                    'from_original_function': owners[ea],
                                    'to_original_function': owners[successor]})
    assert len(controls) == 109 and not indirect and set(external) == EXTERNAL
    assert {(x['from_ida_linear'], x['to_ida_linear']) for x in cross_falls} == FALLTHROUGHS
    caller_context = []
    for target in targets:
        for xref in target['xrefs_to']:
            origin = xref['from']
            if origin in rows:
                continue
            containers = [entry for entry in parent['inventory']
                          if any(start <= origin < end for start, end in entry['chunks'])]
            assert len(containers) == 1
            entry = containers[0]
            status = confirmed.get(entry['ida_linear'])
            caller_context.append({'entry_ida_linear': target['ida_linear'], 'xref': copy.deepcopy(xref),
                                   'caller_original_name': entry['name'], 'caller_ida_linear': entry['ida_linear'],
                                   'caller_status': 'semantic-conformed' if status else 'navigation-only; C pending',
                                   'caller_ledger_entry': status})
    assert [(r['xref']['from'], r['entry_ida_linear'], r['caller_ida_linear']) for r in caller_context] == [
        (0x11D16, 0x125A3, 0x11CD0)]
    scoped = copy.deepcopy(parent)
    scoped.update(targets=targets, decoded_blocks=[], recovery_targets=list(ROOTS), navigation_only_targets=[],
                  unresolved_call_boundaries=[], recovery_ledger_sha256=sha(ledger_path))
    scoped['source_kind'] = 'scoped-extraction-of-pinned-ida-probe'
    scoped['fresh_ida_analysis'] = False
    scoped['closure_scope'] = 'Seventeen original army update functions rooted at125A3; direct controls and fallthroughs are checked against these rows and the pinned confirmed-callee ledger.'
    scoped['scoped_extraction'] = {
        'parent_probe': PARENT, 'parent_probe_sha256': PARENT_SHA,
        'parent_database': DATABASE, 'parent_database_sha256': DATABASE_SHA,
        'parent_recovery_ledger_sha256': parent['recovery_ledger_sha256'],
        'parent_closure_scope': parent['closure_scope'],
        'parent_unresolved_call_boundaries': parent['unresolved_call_boundaries'],
        'preparation_tool': 'tools/c_recovery_army_prepare.py',
        'preparation_tool_sha256': sha(Path(__file__)), 'root_ida_linear': 0x125A3,
        'instruction_count': 434, 'instruction_bytes': 1068, 'direct_control_edges': controls,
        'direct_control_edge_count': 109, 'cross_function_fallthroughs': cross_falls,
        'indirect_control_flow': [], 'mz_relocated_instructions': [],
        'external_callee_ledger': BASELINE,
        'external_callee_ledger_sha256': sha(ledger_path),
        'external_callee_ledger_provenance': BASELINE_PROVENANCE,
        'external_callee_ledger_source_date': provenance['source_date'],
        'current_ledger_validation': {'path': 'docs/re/c-recovery-status.json',
                                      'fields': ['original_name', 'ida_linear', 'routine_sha256', 'status'],
                                      'scope': 'Only the ten external callees; this slice may later be published without changing the extracted source.'},
        'external_callee_ledger_counts': {'functions': 665, 'code_blocks': 62},
        'external_callees': [external[key] for key in sorted(external)],
        'external_caller_context': caller_context,
        'grade_policy': 'Original names, operands, xrefs and graded meanings are copied unchanged; XLAT is data access.',
        'scope_note': 'The retained tool/version/address fields describe the pinned parent IDA analysis; this extraction runs no IDA and creates no new database.',
    }
    return scoped, [{**rows[ea], 'relocation_file_offsets': []} for ea in sorted(rows)], owners, raw


def published_rows(source):
    """解析已發布S的逐指令定位，同時保留真正的GAS文字。"""
    entries, offset, marker = {}, None, None
    for line in source.read_text(encoding='utf-8').splitlines():
        stripped = line.strip()
        if stripped.startswith('.org '):
            offset, marker = int(stripped.split()[1], 0), None
        elif stripped.startswith('# IDA '):
            match = re.match(r'# IDA (0x[0-9a-f]+); file (0x[0-9a-f]+);', stripped)
            assert match and int(match[2], 0) == offset
            marker = int(match[1], 0)
        elif stripped and not stripped.startswith(('#', '.')) and not stripped.endswith(':'):
            assert marker is not None and marker not in entries
            entries[marker] = {'file_offset': offset, 'gas': stripped}
            marker = None
    assert len(entries) == 25097
    return entries


def compile_slice(asm, out, stem, rows, gas, symbols=None, linker=None):
    source = out / (stem + '.S')
    source.write_text('# Independent army instruction slice; gaps are zero placeholders.\n'
                      '.intel_syntax noprefix\n.code16\n.section .image,"ax",@progbits\n' +
                      ''.join(f'.org {row["file_offset"]:#x}\n{gas[row["file_offset"]]}\n' for row in rows),
                      encoding='utf-8')
    if linker is None:
        linker = out / (stem + '.ld')
        linker.write_text(''.join(f'{name} = {value:#x};\n' for name, value in (symbols or {}).items()),
                          encoding='utf-8')
    asm.run(['as', '--32', '-o', str(out / (stem + '.o')), str(source)])
    asm.run(['ld', '-m', 'elf_i386', '-T', str(linker), '--section-start', '.image=0', '--entry', '0',
             '-o', str(out / (stem + '.elf')), str(out / (stem + '.o'))])
    binary = out / (stem + '.bin')
    asm.run(['objcopy', '-O', 'binary', '--only-section=.image', str(out / (stem + '.elf')), str(binary)])
    return source, linker, binary


def generate(repo, out, image):
    assert os.getuid() == 1000 and image == IMAGE
    assert out.is_dir() and out.stat().st_uid == os.getuid()
    translator = repo / 'tools/assembly_rebuild.py'
    protected = [PARENT, DATABASE, BASELINE, BASELINE_PROVENANCE,
                 str(Path(BASELINE).with_suffix('.sha256')), 'workplace/orig/dosv/KI.EXE',
                 'docs/re/assembly-code-record.json',
                 'tools/c_recovery/KI.code.S', 'tools/c_recovery/KI.code.ld',
                 'tools/assembly_rebuild.py', 'tools/c_recovery_army_prepare.py']
    before = {name: sha(repo / name) for name in protected}
    scoped, rows, owners, original = prepare_source(repo)
    source_path = out / 'ida-source.json'
    write_json(source_path, scoped)
    spec = importlib.util.spec_from_file_location('assembly', translator)
    asm = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(asm)
    candidates = out / 'army-assembly-candidates'
    candidates.mkdir(exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, candidates)
    assert not unmatched, 'Instruction-byte fallback is forbidden'
    independent = compile_slice(asm, out, 'army-proof', rows, matched, symbols=symbols)
    manifest_path = repo / 'docs/re/assembly-code-record.json'
    manifest = load(manifest_path)
    assert manifest['input_sha256'] == INPUT_SHA
    published = repo / manifest['source']
    assert sha(published) == manifest['source_sha256']
    published_linker = repo / manifest['linker']
    assert sha(published_linker) == manifest['linker_sha256']
    indexed = published_rows(published)
    published_gas, coverage = {}, []
    for row in rows:
        ea, at = row['ida_linear'], row['file_offset']
        end = at + len(bytes.fromhex(row['bytes']))
        assert indexed[ea]['file_offset'] == at
        ranges = [r for r in manifest['code_ranges'] if r['file_start'] <= at and end <= r['file_end']]
        assert len(ranges) == 1, ('uncovered instruction', hex(ea))
        code_range = ranges[0]
        assert hashlib.sha256(original[code_range['file_start']:code_range['file_end']]).hexdigest() == code_range['sha256']
        published_gas[at] = indexed[ea]['gas']
        coverage.append({'ida_linear': ea, 'file_start': at, 'file_end': end, 'bytes': row['bytes'],
                         'ida_loaded_bytes': row['bytes'], 'loader_relocations': [],
                         'function_name': owners[ea], 'original_assembly': row['assembly'],
                         'operands': copy.deepcopy(row['operands']), 'gas': matched[at],
                         'published_gas': indexed[ea]['gas'], 'published_code_range': code_range,
                         'fixed_prior_coverage': True, 'independent_native_assembly_match': True})
    published_slice = compile_slice(asm, out, 'army-published-slice', rows, published_gas, linker=published_linker)
    for source, linker, binary in (independent, published_slice):
        data, mask = binary.read_bytes(), bytearray(len(binary.read_bytes()))
        for row in rows:
            at, encoded = row['file_offset'], bytes.fromhex(row['bytes'])
            assert data[at:at + len(encoded)] == encoded
            mask[at:at + len(encoded)] = b'\1' * len(encoded)
        assert not any(value and not occupied for value, occupied in zip(data, mask))
    after = {name: sha(repo / name) for name in protected}
    assert before == after, 'Protected input changed during preparation'
    validate_external_ledger(repo, load(repo / BASELINE))
    result = {
        'schema': 'wolong-matching-code-supplement-v1', 'input_sha256': INPUT_SHA,
        'tool': scoped['tool'], 'tool_version': scoped['tool_version'], 'fresh_ida_analysis': False,
        'source_kind': scoped['source_kind'], 'database_sha256': DATABASE_SHA,
        'probe_sha256': sha(source_path), 'scoped_source': 'workplace/matching-decompilation/c-army/ida-source.json',
        'parent_probe': PARENT, 'parent_probe_sha256': PARENT_SHA,
        'evidence': EVIDENCE, 'constant_symbols': {}, 'instructions': [],
        'instruction_coverage': coverage, 'loader_relocation_normalization': [],
        'named_instruction_count': 434, 'named_instruction_bytes': 1068,
        'recovery_targets': list(ROOTS), 'navigation_only_targets': [],
        'original_functions': [{'name': t['name'], 'ida_linear': t['ida_linear'],
                                'file_sha256': t['file_sha256'], 'xrefs_to': t['xrefs_to'],
                                'chunks': [[ch['start'], ch['end']] for ch in t['chunks']]}
                               for t in scoped['targets']],
        'assembler_verification': {
            'status': 'byte-exact', 'instructions': 0, 'bytes': 0,
            'input_instruction_rows': 434, 'unique_input_instructions': 434,
            'covered_input_instructions': 434, 'instruction_byte_fallbacks': 0,
            'coverage_source': 'Current published source instruction markers and independently compiled published GAS slice',
            'published_manifest': 'docs/re/assembly-code-record.json', 'published_manifest_sha256': sha(manifest_path),
            'published_code_source': manifest['source'], 'published_code_source_sha256': sha(published),
            'published_linker': manifest['linker'], 'published_linker_sha256': sha(published_linker),
            'independent_reassembly': {'status': 'byte-exact', 'instructions': 434, 'bytes': 1068,
                                      'source_sha256': sha(independent[0]), 'linker_sha256': sha(independent[1]),
                                      'binary_sha256': sha(independent[2]), 'constant_symbols': symbols, 'trials': trials},
            'published_slice_reassembly': {'status': 'byte-exact', 'instructions': 434, 'bytes': 1068,
                                          'source_sha256': sha(published_slice[0]),
                                          'linker_sha256': sha(published_slice[1]), 'binary_sha256': sha(published_slice[2])},
            'image_id': IMAGE, 'assembler': asm.run(['as', '--version']).splitlines()[0],
            'linker': asm.run(['ld', '--version']).splitlines()[0],
            'generator_sha256': sha(Path(__file__)), 'translator_sha256': sha(translator),
            'sources_before': before, 'sources_after': after,
        },
        'address_space': 'IDA database linear base0x10000; separate original-file offsets include512-byte MZ header',
        'notes': ['Scoped extraction preserves the pinned parent IDA names, grades, operands, chunks and xrefs.',
                  'No fresh IDA run or database mutation; no semantic C conformance claim.',
                  'No new instructions and no global matching manifest or assembly source modifications.'],
    }
    write_json(out / 'army-code-supplement.json', result)
    public = repo / 'docs/re/c-army-code.json'
    write_json(public, result)
    print(f'17 original functions /434 instructions /1068bytes independently exact; new0; published slice exact')
    print('scoped_source_sha256', sha(source_path))
    print('public_receipt_sha256', sha(public))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--image-id', required=True)
    args = parser.parse_args()
    generate(args.repo, args.output, args.image_id)
