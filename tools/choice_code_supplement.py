#!/usr/bin/env python3
"""原 selector 的 live operand 區以助憶碼匹配，保留原 byte／定位。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path


def main(repo, output, image_id):
    spec = importlib.util.spec_from_file_location('assembly', repo / 'tools/assembly_rebuild.py')
    asm = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(asm)
    raw = (repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == asm.EXPECTED
    probe_path = repo / 'workplace/matching-decompilation/c-choice/ida/ida-probe.json'
    probe = json.loads(probe_path.read_text())
    mapping = json.loads((repo / 'workplace/matching-decompilation/assembly/build/source-map.json').read_text())
    old = json.loads((repo / 'docs/re/rectangle-handler-code.json').read_text())
    existing = {r['file_start'] for r in old['instructions']}
    rows = []
    for block in probe['decoded_blocks']:
        for row in block['instructions']:
            at = row['ida_linear'] - 0x10000 + 512
            end = at + len(bytes.fromhex(row['bytes']))
            if at in existing or not any(s['kind'] == 'data-or-unknown' and s['file_start'] <= at < end <= s['file_end'] for s in mapping):
                continue
            assert raw[at:end] == bytes.fromhex(row['bytes'])
            rows.append({**row, 'file_offset': at, 'relocation_file_offsets': []})
    assert len(rows) == 10 and sum(len(bytes.fromhex(r['bytes'])) for r in rows) == 19
    output.mkdir(parents=True, exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, output)
    assert not unmatched and not symbols
    public = old['instructions'].copy()
    for row in rows:
        at = row['file_offset']
        public.append({'ida_linear': row['ida_linear'], 'file_start': at, 'file_end': at + len(bytes.fromhex(row['bytes'])),
                       'kind': 'instruction', 'original_name': row['original_name_before_analysis'] or None,
                       'function_name': None, 'original_assembly': row['assembly'], 'gas': matched[at], 'bytes': row['bytes'],
                       'original_ida_data_line': row['original_ida_data_line'], 'ida_code_classified': row['ida_code_classified_before'],
                       'operands': row['operands'], 'locator_level': 'proven', 'scope_sources': ['docs/re/112-c-choice-selector-restoration.md']})
    record = {**old, 'database_sha256': hashlib.sha256((probe_path.parent / 'input.exe.i64').read_bytes()).hexdigest(),
              'evidence': 'docs/re/112-c-choice-selector-restoration.md', 'instructions': sorted(public, key=lambda r: r['file_start']),
              'assembler_image_id': image_id, 'assembler_trials': trials}
    (repo / 'docs/re/rectangle-handler-code.json').write_text(json.dumps(record, ensure_ascii=False, indent=2) + '\n')
    print('Added 10 instructions / 19 bytes; prior 359 retained; no instruction-byte fallback')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--image-id', required=True)
    args = parser.parse_args()
    main(args.repo, args.output, args.image_id)
