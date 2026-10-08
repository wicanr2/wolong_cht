#!/usr/bin/env python3
"""原始彈出 selector、row callback、選取 XOR 與 frame 的 IDA 證據。"""
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

probe.TARGETS = (0x1036F, 0x103C3, 0x103E6, 0x10414, 0x104B5, 0x104FF,
                 0x1054D, 0x1059B, 0x105ED, 0x1061F, 0x10B46, 0x10BAF,
                 0x19479, 0x194BF)
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
                    shifted.append({'file_offset': offset, 'original_word': before})
            assert bytes(expected) == bytes.fromhex(chunk['bytes'])
            chunk.update(file_bytes=data.hex(), loader_relocations=shifted)
            parts.append(data)
        target['file_sha256'] = hashlib.sha256(b''.join(parts)).hexdigest()
    report['decoded_blocks'] = []
    for start, end in [(0x19409, 0x1945A), (0x1945A, 0x19479)]:
        rows = []
        ea = start
        while ea < end:
            if ea in (0x19441, 0x19465):
                raise ValueError('inline operand reached as instruction')
            insn = ida_ua.insn_t()
            size = ida_ua.decode_insn(insn, ea)
            assert size > 0 and ea + size <= end
            before = ida_bytes.get_bytes(ea, size)
            classified = ida_bytes.is_code(ida_bytes.get_full_flags(ea))
            original = idc.generate_disasm_line(ea, 0)
            if not classified:
                ida_bytes.del_items(ea, ida_bytes.DELIT_SIMPLE, size)
                assert ida_ua.create_insn(ea) == size
            row = probe.instruction(ea)
            assert bytes.fromhex(row['bytes']) == before
            row.update(original_name_before_analysis=idc.get_name(ea),
                       original_ida_data_line=original, ida_code_classified_before=classified)
            rows.append(row)
            ea += size
        report['decoded_blocks'].append({'ida_linear': start, 'end': end,
                                        'original_name_before_analysis': idc.get_name(start), 'instructions': rows})
    report['live_operands'] = [{'ida_linear': ea, 'size': size,
                               'file_bytes': raw[ea - 0x10000 + header:ea - 0x10000 + header + size].hex()}
                              for ea, size in [(0x1036D, 2), (0x19441, 1), (0x19443, 1), (0x19465, 1), (0x19467, 2)]]
    report['relocation_file_offsets'] = relocations
    path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
except Exception:
    Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
    ida_pro.qexit(1)
else:
    ida_pro.qexit(0)
