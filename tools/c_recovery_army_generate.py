#!/usr/bin/env python3
"""將固定軍團證據產生為靜態C控制流，保留原始入口與guest stack。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

PROBE_SHA = 'ed6d851a394023ab58809aa504391bd20e7c649aa721ac57b8ef99ee925c26b7'
INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
ROOTS = (0x125A3, 0x12600, 0x1264A, 0x12662, 0x126FF, 0x12708,
         0x127A2, 0x127F6, 0x12804, 0x12808, 0x12831, 0x12880,
         0x128F4, 0x12A7E, 0x142AB, 0x14300, 0x1562B)
# number, expected runner group, name, original instruction, bad C statement.
MUTANTS = (
    (1, 'batch', 'batch-count', 0x125B2, 'm->cx=15;'),
    (2, 'batch', 'timer-reload', 0x125C9, '; /* omitted original timer reload */'),
    (3, 'upkeep', 'upkeep-road-branch', 0x12611, 'if(!(m->flags&1)) {goto L_12624;}'),
    (4, 'upkeep', 'upkeep-quarter-shift', 0x12617, '; /* omitted original quarter shift */'),
    (5, 'upkeep', 'morale-increment', 0x1263D,
     'dl_write8(m,m->ds,(uint16_t)(m->si+6),ec_math(m,dl_read8(m,m->ds,(uint16_t)(m->si+6)),9,8,0,0));'),
    (6, 'direction', 'signed-direction', 0x12807, 'm->ax=(uint8_t)m->ax;'),
    (7, 'movement', 'occupancy-decrement', 0x1269C, '; /* omitted original occupancy decrement */'),
    (8, 'encounter', 'collision-last-slot', 0x1283A, 'm->cx=126;'),
    (9, 'encounter', 'standoff-countdown', 0x12866, 'dl_write8(m,m->ds,(uint16_t)(m->si+3),11);'),
    (10, 'peace', 'peace-boundary', 0x142E3, 'ec_cmp(m,(uint8_t)m->ax,0x81,8);'),
    (11, 'cleanup', 'collapse-duty', 0x12A8E, '; /* omitted original general duty clear */'),
    (12, 'arrival', 'ai-city-gate', 0x14304,
     'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+0x858)),0,8);'),
)


def generate(repo, probe_path, output=None):
    output = output or repo / 'tools/c_recovery'
    assert output.is_dir()
    encoded = probe_path.read_bytes()
    assert hashlib.sha256(encoded).hexdigest() == PROBE_SHA
    probe = json.loads(encoded)
    assert probe['input_sha256'] == probe['ida_input_sha256'] == INPUT_SHA
    assert probe['source_kind'] == 'scoped-extraction-of-pinned-ida-probe'
    assert probe['fresh_ida_analysis'] is False
    assert not probe['unresolved_call_boundaries'] and not probe['decoded_blocks']
    assert tuple(probe['recovery_targets']) == ROOTS
    targets = probe['targets']
    assert tuple(t['ida_linear'] for t in targets) == ROOTS
    rows = {i['ida_linear']: i for t in targets for chunk in t['chunks'] for i in chunk['instructions']}
    assert len(rows) == 434
    assert sum(len(bytes.fromhex(i['bytes'])) for i in rows.values()) == 1068
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    entries = {t['ida_linear']: t['name'] for t in targets}
    mutations = {ea: (number, bad) for number, group, name, ea, bad in MUTANTS}
    assert len(mutations) == 12 and set(mutations) <= rows.keys()
    lines = ['/* Fixed original army C graph; re/132, spec/252. No guest instruction decoder. */',
             'static int am_known(uint16_t t) {switch(t) {']
    lines += [f'case 0x{ea-0x10000:x}:' for ea in sorted(rows)]
    lines += ['return 1;default:return 0;}}',
              'static void am_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int skip_first_enter) {',
              'const uint16_t core_cs=m->cs;int invoke_external=0;m->ip=entry;',
              'AM_DISPATCH:',
              'if(m->cs!=core_cs) {assert(!invoke_external);return;}',
              'switch(m->ip) {']
    lines += [f'case 0x{ea-0x10000:x}:goto L_{ea:X};' for ea in sorted(rows)]
    lines += ['default:break;}', 'if(!invoke_external)return;',
              'ki_engagement_invoke(m,m->ip,h);invoke_external=0;goto AM_DISPATCH;']
    for ea, row in sorted(rows.items()):
        mnemonic, ops = row['mnemonic'], row['operands']
        next_ea = ea + len(bytes.fromhex(row['bytes']))
        nextoff = next_ea - 0x10000
        lines += [f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */', f'm->ip=0x{ea-0x10000:x};']
        if ea in entries:
            lines += [f'/* Original entry: {entries[ea]} */',
                      f'if(skip_first_enter)skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x{ea-0x10000:x},h->user);']
        else:
            lines += ['skip_first_enter=0;']

        def jump(target):
            if target in rows:
                return f'goto L_{target:X};'
            return f'm->ip=0x{target-0x10000:x};invoke_external=1;goto AM_DISPATCH;'

        terminal = False
        if mnemonic == 'call':
            assert ops[0]['type'] == 7
            statement = '{uint16_t target=%s;ec_push(m,0x%x);m->ip=target;invoke_external=1;goto AM_DISPATCH;}' % (dl.read(ops[0]), nextoff)
            terminal = True
        elif mnemonic == 'retn':
            assert row['bytes'] == 'c3'
            statement = 'm->ip=ec_pop(m);invoke_external=0;goto AM_DISPATCH;'
            terminal = True
        elif mnemonic == 'jmp':
            assert ops[0]['type'] == 7
            statement, terminal = jump(ops[0]['addr'] + 0x10000), True
        elif mnemonic in dl.COND:
            statement = f'if({dl.COND[mnemonic]}) {{{jump(ops[0]["addr"]+0x10000)}}}'
        elif mnemonic == 'loop':
            statement = f'if(--m->cx) {{{jump(ops[0]["addr"]+0x10000)}}}'
        elif mnemonic == 'cbw':
            statement = 'm->ax=(uint16_t)(int16_t)(int8_t)m->ax;'
        elif mnemonic == 'xlat':
            statement = 'vg_al(m,dl_read8(m,m->ds,(uint16_t)(m->bx+(uint8_t)m->ax)));'
        else:
            try:
                statement = dl.statement(row)
            except Exception as error:
                raise ValueError(f'{ea:X}: {row["assembly"]}: {error}') from error
        if ea in mutations:
            number, bad = mutations[ea]
            assert bad != statement
            lines += [f'#if KI_ARMY_MUTATION == {number}', bad, '#else', statement, '#endif']
        else:
            lines += [statement]
        if not terminal:
            assert next_ea in rows, ('missing fallthrough', hex(ea), hex(next_ea))
            lines += [f'goto L_{next_ea:X};']
    lines += ['}', 'int ki_army_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {',
              'if(!am_known(t))return 0;',
              'uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);',
              'am_graph(m,t,h,1); /* External resolver already emitted the entry hook. */',
              'assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);',
              'ec_push(m,ret);return 1;}',
              'void ki_army_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(am_known(t))am_graph(m,t,h,0);else ki_engagement_invoke(m,t,h);}']
    declarations = []
    for target in targets:
        name, offset = target['name'], target['ida_linear'] - 0x10000
        lines += [f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{am_graph(m,0x{offset:x},h,0);}}']
        declarations += [f'void {name}(KiMachine16 *,const KiEconomyHooks *);']
    (output / 'army_generated.inc').write_text('\n'.join(lines) + '\n', encoding='utf-8')
    header = ['#ifndef WOLONG_C_ARMY_H', '#define WOLONG_C_ARMY_H', '#include "engagement.h"',
              'int ki_army_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);',
              'void ki_army_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);',
              'uint64_t ki_army_compiled_source(void);', 'unsigned ki_army_mutation(void);',
              *declarations, '#endif']
    (output / 'army.h').write_text('\n'.join(header) + '\n', encoding='utf-8')
    summary = {'scope': 'spec252 original army static C graph; semantic verification pending',
               'probe_sha256': PROBE_SHA, 'named_functions': 17, 'named_instructions': 434,
               'named_instruction_bytes': 1068, 'raw_blocks': 0, 'raw_instructions': 0,
               'raw_instruction_bytes': 0, 'unique_instructions': 434,
               'hook_entries': [f'0x{ea:X}' for ea in sorted(entries)],
               'instruction_decoder': False, 'c_side_cpu_step': False,
               'external_invoke': 'ki_engagement_invoke with unchanged hooks',
               'mutants': {str(n): {'group': group, 'name': name, 'ida_linear': ea}
                           for n, group, name, ea, bad in MUTANTS}}
    print('Generated 17 named /434 instructions /1068bytes;12 controls; original guest stack and fallthroughs')
    return summary


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', required=True, type=Path)
    parser.add_argument('--probe', required=True, type=Path)
    parser.add_argument('--output', type=Path)
    args = parser.parse_args()
    generate(args.repo, args.probe, args.output)
