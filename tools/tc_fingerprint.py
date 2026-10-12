#!/usr/bin/env python3
"""KI.EXE 工具鏈指紋普查：讀 census＋原檔 bytes，不跑 IDA。spec/255、re/135。"""
import argparse
import collections
import hashlib
import json
import re
import struct
import sys
from pathlib import Path

INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', required=True, type=Path)
    args = parser.parse_args()
    raw = (args.repo / 'workplace/orig/dosv/KI.EXE').read_bytes()
    assert hashlib.sha256(raw).hexdigest() == INPUT_SHA, 'input identity'
    header = struct.unpack_from('<H', raw, 8)[0] * 16
    ledger = json.loads((args.repo / 'workplace/matching-decompilation/assembly/build/function-ledger.json').read_text())
    assert len(ledger) == 739, len(ledger)

    prologue = collections.Counter()
    bp_frame, bp_save = [], []
    for entry in ledger:
        start = entry['ida_linear'] - 0x10000 + header
        head = raw[start:start + 6]
        prologue[head[:2]] += 1
        if head[:1] == b'\x55':
            if head[1:3] == b'\x83\xec' or head[1:3] == b'\x81\xec':
                bp_frame.append(entry['original_name'])
            else:
                bp_save.append(entry['original_name'])
    print(f'prologue top: {prologue.most_common(5)}')
    print(f'bp-frame funcs ({len(bp_frame)}): {sorted(bp_frame)}')
    print(f'push-bp-save funcs ({len(bp_save)}): {sorted(bp_save)}')
    assert len(bp_frame) == 6 and len(bp_save) == 2

    single, tails, multi = 0, collections.Counter(), []
    for entry in ledger:
        if len(entry['chunks']) == 1:
            single += 1
            tails[raw[entry['chunks'][0][1] - 0x10000 + header - 1]] += 1
        else:
            multi.append(entry['original_name'])
    print(f'single-chunk: {single}, tail C3: {tails[0xC3]}, tail CB: {tails[0xCB]}, multi: {multi}')
    assert single == 737 and tails[0xC3] == 717 and len(multi) == 2
    cb = [e['original_name'] for e in ledger if len(e['chunks']) == 1
          and raw[e['chunks'][0][1] - 0x10000 + header - 1] == 0xCB]
    print(f'tail-CB owner: {cb} (far-call driver segment, retf expected)')
    assert cb == ['sub_20000']

    asm = (args.repo / 'tools/c_recovery/KI.code.S').read_text(encoding='utf-8').splitlines()
    insns = [l for l in asm if l.startswith('# IDA')]
    print(f'indexed instructions: {len(insns)}')
    assert len(insns) == 25097

    ints = collections.Counter()
    for line in insns:
        m = re.search(r'; int\s+(?:([0-9A-Fa-f]+)h|(3))\b', line)
        if m:
            ints[(m.group(1) or m.group(2)).upper()] += 1
    print(f'int sites: {dict(sorted(ints.items()))} total={sum(ints.values())}')
    assert dict(ints) == {'10': 2, '15': 2, '1A': 1, '21': 34, '33': 8, '50': 4, '61': 16, '3': 1}

    def owner(ea: int):
        for entry in ledger:
            for a, b in entry['chunks']:
                if a <= ea < b:
                    return entry['original_name']
        return None

    framed = {'sub_17FDB': 2, 'sub_1820E': 0x200, 'sub_19479': 6,
              'sub_194BF': 6, 'sub_1F4DF': 6, 'sub_1F593': 6}
    stack_args = []
    for line in insns:
        code = line.split(';')[-1]
        for m in re.finditer(r'\[(?:bp|BP)\+([0-9A-Fa-f]+)(h)?\]', code):
            value = int(m.group(1), 16 if m.group(2) else 10)
            if value < 8:
                continue
            ea = int(re.match(r'# IDA (0x[0-9A-Fa-f]+);', line).group(1), 16)
            name = owner(ea)
            if name in framed and value < framed[name]:
                continue
            if name in framed:
                stack_args.append((hex(ea), name, code.strip()))
    print(f'stack-arg refs in framed funcs: {len(stack_args)}')
    assert not stack_args

    rep_funcs = set()
    for line in insns:
        if re.search(r'; (rep|repe|repne) ', line):
            ea = int(re.match(r'# IDA (0x[0-9A-Fa-f]+);', line).group(1), 16)
            rep_funcs.add(owner(ea))
    reps = sum(1 for line in insns if re.search(r'; (rep|repe|repne) ', line))
    print(f'rep string ops: {reps} in {len(rep_funcs)} funcs')
    assert reps == 27 and len(rep_funcs) == 19

    gpr32 = [l for l in insns if re.search(r'\se(ax|bx|cx|dx|si|di|bp|sp)[,\s\]]', l.split(';')[-1])]
    print(f'32-bit GPR refs: {len(gpr32)}')
    assert not gpr32
    print('toolchain fingerprint: PASS (hand-written 16-bit assembly shape)')
    return 0


if __name__ == '__main__':
    sys.exit(main())
