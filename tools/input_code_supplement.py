#!/usr/bin/env python3
"""原始 wait patch 區以助憶碼匹配，MZ file bytes 與 IDA bytes 分開保留。"""
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
    probe_path = repo / 'workplace/matching-decompilation/c-input/ida/ida-probe.json'
    probe = json.loads(probe_path.read_text())
    mapping = json.loads((repo / 'workplace/matching-decompilation/assembly/build/source-map.json').read_text())
    old = json.loads((repo / 'docs/re/rectangle-handler-code.json').read_text())
    existing = {r['file_start'] for r in old['instructions']}
    rows, originals = [], {}
    for block in probe['decoded_blocks']:
        for row in block['instructions']:
            start = row['ida_linear'] - 0x10000 + 512
            end = start + len(bytes.fromhex(row['bytes']))
            if start in existing:
                continue
            if not any(s['kind'] == 'data-or-unknown' and s['file_start'] <= start < end <= s['file_end'] for s in mapping):
                continue
            converted = {**row, 'file_offset': start, 'bytes': raw[start:end].hex(),
                         'relocation_file_offsets': [r for r in probe['relocation_file_offsets'] if start <= r < end]}
            rows.append(converted)
            originals[start] = row
    assert len(rows) == 21 and sum(len(bytes.fromhex(r['bytes'])) for r in rows) == 71
    output.mkdir(parents=True, exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, output)
    assert not unmatched and not symbols
    public = old['instructions'].copy()
    for row in rows:
        start = row['file_offset']
        original = originals[start]
        public.append({'ida_linear': row['ida_linear'], 'file_start': start, 'file_end': start + len(bytes.fromhex(row['bytes'])),
                       'kind': 'instruction', 'original_name': original['original_name_before_analysis'] or None,
                       'function_name': None, 'original_assembly': row['assembly'], 'gas': matched[start],
                       'bytes': row['bytes'], 'ida_bytes': original['bytes'],
                       'relocation_file_offsets': row['relocation_file_offsets'],
                       'original_ida_data_line': original['original_ida_data_line'],
                       'ida_code_classified': original['ida_code_classified_before'], 'operands': original['operands'],
                       'locator_level': 'proven', 'scope_sources': ['docs/re/111-c-input-mouse-restoration.md']})
    record = {**old, 'database_sha256': hashlib.sha256((probe_path.parent / 'input.exe.i64').read_bytes()).hexdigest(),
              'evidence': 'docs/re/111-c-input-mouse-restoration.md', 'instructions': sorted(public, key=lambda r: r['file_start']),
              'assembler_image_id': image_id, 'assembler_trials': trials}
    (repo / 'docs/re/rectangle-handler-code.json').write_text(json.dumps(record, ensure_ascii=False, indent=2) + '\n')
    print('Added 21 instructions / 71 bytes; 3 MZ fields; old supplement retained; no instruction-byte fallback')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--image-id', required=True)
    args = parser.parse_args()
    main(args.repo, args.output, args.image_id)
