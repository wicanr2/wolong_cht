#!/usr/bin/env python3
"""IDA 9.4：matching decompilation 試點的非破壞性證據匯出。

由 tools/ida.sh probe 執行；原始名稱、位址、bytes 與 xref 均保留。
輸出只進 /output，目標 EXE 與既有資料庫不改寫。
"""
import hashlib
import json
import os
import re
import sys
import traceback
from pathlib import Path

import ida_auto
import ida_bytes
import ida_funcs
import ida_ida
import ida_kernwin
import ida_loader
import ida_nalt
import ida_pro
import ida_segment
import ida_ua
import idautils
import idc

TARGETS = (0x1ECE0, 0x11D8E, 0x131AE)
SEMANTICS = {}


def evidence_path(source, docs_root=Path('/documents')):
    path = Path(source)
    if path.is_absolute() or '..' in path.parts or not path.parts or path.parts[0] != 'docs':
        raise ValueError(f'invalid semantic evidence path: {source}')
    return docs_root / path.relative_to('docs')


def annotation(kind, key, original=None):
    row = SEMANTICS.get(kind, {}).get(key)
    if row and original is not None and row.get('original') != original:
        raise ValueError(f"semantic index operand differs at {key}: {original}")
    result = dict(row) if row else {'meaning': None, 'level': 'unknown', 'sources': []}
    result['warning'] = None if result['level'] == 'proven' else '⚠ 語意尚未證實'
    return result


def instruction(ea):
    insn = ida_ua.insn_t()
    size = ida_ua.decode_insn(insn, ea)
    if size <= 0:
        raise ValueError(f"cannot decode IDA linear address {ea:#x}")
    return {
        "ida_linear": ea,
        "file_offset": ida_loader.get_fileregion_offset(ea),
        "bytes": ida_bytes.get_bytes(ea, size).hex(),
        "mnemonic": idc.print_insn_mnem(ea),
        "assembly": idc.generate_disasm_line(ea, 0),
        "operands": [
            {"type": op.type, "addr": op.addr, "value": op.value, "dtype": op.dtype,
             "dtype_size": ida_ua.get_dtype_size(op.dtype), "reg": op.reg,
             "original": idc.print_operand(ea, i),
             "semantic": annotation('operands', f'0x{ea:X}:{i}',
                                    idc.print_operand(ea, i))}
            for i, op in enumerate(insn.ops) if op.type != ida_ua.o_void
        ],
    }


def main():
    global SEMANTICS
    ida_auto.auto_wait()
    input_path = Path(ida_nalt.get_input_file_path())
    raw = input_path.read_bytes()
    digest = hashlib.sha256(raw).hexdigest()
    database_digest = ida_nalt.retrieve_input_file_sha256().hex()
    if digest != database_digest:
        raise ValueError("IDA input SHA-256 differs from actual input")
    index_path = Path('/evidence/matching-semantic-index.json')
    index_raw = index_path.read_bytes()
    index = json.loads(index_raw)
    if index.get('schema') != 'wolong-matching-semantic-index-v1':
        raise ValueError('unexpected semantic index schema')
    if digest == index['input_sha256']:
        SEMANTICS = index
        for kind in ['functions', 'operands']:
            for row in index[kind].values():
                if row['level'] not in ('proven', 'strong inference', 'hypothesis', 'unknown'):
                    raise ValueError('ungraded semantic index row')
                for source in row['sources']:
                    if not evidence_path(source).is_file():
                        raise ValueError(f'missing semantic evidence: {source}')
    functions = list(idautils.Functions())
    inventory = []
    for ea in functions:
        f = ida_funcs.get_func(ea)
        inventory.append({
            "name": ida_funcs.get_func_name(ea), "ida_linear": ea,
            "chunks": list(idautils.Chunks(ea)), "flags": f.flags,
            "library_flag": bool(f.flags & ida_funcs.FUNC_LIB),
            "bp_frame_prefix": ida_bytes.get_bytes(ea, 3) in
            (b"\x55\x8b\xec", b"\x55\x89\xe5"),
        })
    targets = []
    for ea in TARGETS if input_path.name.upper() == 'KI.EXE' else ():
        f = ida_funcs.get_func(ea)
        if f is None or f.start_ea != ea:
            raise ValueError(f"target is not an IDA function entry: {ea:#x}")
        chunks = []
        for start, end in idautils.Chunks(ea):
            body = []
            for address in idautils.Heads(start, end):
                if ida_bytes.is_code(ida_bytes.get_full_flags(address)):
                    body.append(instruction(address))
            chunks.append({"start": start, "end": end,
                           "bytes": ida_bytes.get_bytes(start, end-start).hex(),
                           "instructions": body})
        targets.append({
            "name": ida_funcs.get_func_name(ea), "ida_linear": ea,
            "semantic": annotation('functions', f'0x{ea:X}'),
            "chunks": chunks,
            "xrefs_to": [{"from": x.frm, "type": x.type}
                         for x in idautils.XrefsTo(ea)],
        })
    segments = []
    for ea in idautils.Segments():
        s = ida_segment.getseg(ea)
        segments.append({"name": ida_segment.get_segm_name(s),
                         "start": s.start_ea, "end": s.end_ea,
                         "selector": s.sel, "bitness": s.bitness})
    report = {
        "schema": "wolong-matching-ida-probe-v1",
        "tool": "IDA Pro", "tool_version": ida_kernwin.get_kernel_version(),
        "python_version": sys.version.split()[0],
        "input_name": input_path.name, "input_sha256": digest,
        "image_id": os.environ.get('WOLONG_IDA_IMAGE_ID'),
        "semantic_index_sha256": hashlib.sha256(index_raw).hexdigest(),
        "semantic_index_applied": bool(SEMANTICS),
        "ida_input_sha256": database_digest,
        "address_space": "IDA database linear; file_offset stored separately",
        "compiler_id_navigation_only": ida_ida.inf_get_cc_id(),
        "function_count": len(functions), "segments": segments,
        "inventory": inventory, "targets": targets,
        "printable_toolchain_markers": [
            {"file_offset": m.start(), "text": m.group().decode('ascii')}
            for m in re.finditer(rb'[\x20-\x7e]{6,}', raw)
            if re.search(rb'borland|turbo|watcom|microsoft|copyright|tlink|masm|tasm',
                         m.group(), re.I)
        ],
        "entry_instructions": [instruction(ea) for ea in
                               list(idautils.Heads(0x10000, 0x1006B))
                               if ida_bytes.is_code(ida_bytes.get_full_flags(ea))],
    }
    Path('/output/ida-probe.json').write_text(
        json.dumps(report, ensure_ascii=False, indent=2)+"\n", encoding='utf-8')


if __name__ == '__main__':
    try:
        main()
    except Exception:
        Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(), encoding='utf-8')
        ida_pro.qexit(1)
    else:
        ida_pro.qexit(0)
