#!/usr/bin/env python3
"""原始訊息等待、mouse far table 與游標閉包的 IDA 證據。"""
import hashlib
import json
import struct
import sys
import traceback
from pathlib import Path

sys.path.insert(0, '/tools')
import ida_matching_probe as probe
import ida_bytes
import ida_nalt
import ida_pro
import ida_ua
import idc

probe.TARGETS = (0x18810, 0x121E7, 0x20000, 0x2002E, 0x2006F, 0x20070, 0x2009A,
                 0x200BD, 0x200C0, 0x20137, 0x2016B, 0x2019D, 0x201C6,
                 0x201E4, 0x2020C, 0x20249, 0x202A0, 0x202BD, 0x202FE)
try:
    probe.main()
    path = Path('/output/ida-probe.json')
    report = json.loads(path.read_text(encoding='utf-8'))
    raw = Path(ida_nalt.get_input_file_path()).read_bytes()
    header = struct.unpack_from('<H', raw, 8)[0] * 16
    relocations = [header + off + seg * 16 for off, seg in
                   (struct.unpack_from('<HH', raw, struct.unpack_from('<H', raw, 24)[0] + i * 4)
                    for i in range(struct.unpack_from('<H', raw, 6)[0]))]
    for target in report['targets']:
        parts = []
        for chunk in target['chunks']:
            start = chunk['start'] - 0x10000 + header
            data = raw[start:start + chunk['end'] - chunk['start']]
            expected = bytearray(data)
            shifted = []
            for offset in relocations:
                if start <= offset < start + len(data):
                    before = struct.unpack_from('<H', raw, offset)[0]
                    struct.pack_into('<H', expected, offset - start, (before + 0x1000) & 65535)
                    shifted.append({'file_offset': offset, 'original_word': before, 'ida_word': (before + 0x1000) & 65535})
            assert bytes(expected) == bytes.fromhex(chunk['bytes'])
            chunk.update(file_bytes=data.hex(), loader_relocations=shifted)
            parts.append(data)
        target['file_sha256'] = hashlib.sha256(b''.join(parts)).hexdigest()
    report['decoded_blocks'] = []
    for start, end in [(0x1222B, 0x12286), (0x20101, 0x20137)]:
        rows = []
        ea = start
        while ea < end:
            insn = ida_ua.insn_t()
            size = ida_ua.decode_insn(insn, ea)
            assert size > 0 and ea + size <= end
            original = idc.generate_disasm_line(ea, 0)
            classified = ida_bytes.is_code(ida_bytes.get_full_flags(ea))
            before = ida_bytes.get_bytes(ea, size)
            if not classified:
                ida_bytes.del_items(ea, ida_bytes.DELIT_SIMPLE, size)
                assert ida_ua.create_insn(ea) == size
            row = probe.instruction(ea)
            assert bytes.fromhex(row['bytes']) == before
            row.update(original_ida_data_line=original,
                       original_name_before_analysis=idc.get_name(ea),
                       ida_code_classified_before=classified)
            rows.append(row)
            ea += size
        report['decoded_blocks'].append({'ida_linear': start, 'end': end,
                                        'original_name_before_analysis': idc.get_name(start),
                                        'instructions': rows})
    table = raw[0x1020E:0x1022E]
    assert ida_bytes.get_bytes(0x2000E, 32) == table
    report['mouse_table'] = {'ida_linear': 0x2000E, 'file_offset': 0x1020E,
                             'file_bytes': table.hex(),
                             'offsets': [int.from_bytes(table[i:i + 2], 'little') for i in range(0, 32, 2)]}
    report['patch_word'] = {'ida_linear': 0x12239, 'file_offset': 0x2439, 'file_bytes': raw[0x2439:0x243B].hex()}
    report['relocation_file_offsets'] = relocations
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
