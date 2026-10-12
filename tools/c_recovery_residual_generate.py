#!/usr/bin/env python3
"""殘部群 10 函式靜態C圖；re/134、spec/254。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

PROBE_SHA = '63e12863b199542902260e5f6b917c00e0f268b923ec0149a7f321d3f07a8a4a'
INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
WANTED = {0x145F8, 0x15AD1, 0x15AFC, 0x15C14, 0x19377, 0x193D7,
          0x1AB7C, 0x1C300, 0x1E567, 0x1E6FF}
RAW_WANTED = {0x1AB48: 1}
HOOK_ENTRIES = {0x15B7B}
MUTANT_STATEMENTS = {
    0x14606: (1, 'dl_write8(m,m->ds,(uint16_t)(m->bx+25),254);'),
    0x15AD7: (2, 'm->cx=(uint16_t)(5);'),
    0x15B71: (3, 'ec_cmp(m,(uint8_t)m->ax,23,8);'),
    0x15C26: (4, 'ec_cmp(m,(uint8_t)m->bx,12,8);'),
    0x1938A: (5, 'ec_cmp(m,(uint8_t)m->ax,2,8);'),
    0x193DC: (6, 'ec_cmp(m,(uint8_t)(m->ax>>8),3,8);'),
    0x1AB82: (7, 'dl_write8(m,m->ds,(uint16_t)(m->si+25),41);'),
    0x1C307: (8, 'vg_al(m,(uint8_t)(6));'),
    0x1E577: (9, 'vg_al(m,(uint8_t)(dl_read8(m,m->cs,58495)));'),
    0x1E70F: (10, 'vg_al(m,(uint8_t)(dl_read8(m,m->cs,58494)));'),
}


def generate(repo, probe_path, output=None):
    output = output or repo / "tools/c_recovery"
    raw_probe = probe_path.read_bytes()
    assert hashlib.sha256(raw_probe).hexdigest() == PROBE_SHA
    probe = json.loads(raw_probe)
    assert probe['input_sha256'] == probe['ida_input_sha256'] == INPUT_SHA
    assert probe['function_count'] == 739
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    assert set(probe['recovery_targets']) == WANTED
    selected = [t for t in probe['targets'] if t['ida_linear'] in WANTED]
    assert len(selected) == 10
    rows = {}
    named_count = named_bytes = raw_count = raw_bytes = 0
    entries = set(WANTED) | set(RAW_WANTED) | set(HOOK_ENTRIES)
    originals = {}
    for target in selected:
        originals[target['ida_linear']] = target['name']
        for chunk in target['chunks']:
            assert hashlib.sha256(b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])).hexdigest() == target['file_sha256']
            for row in chunk['instructions']:
                assert row['ida_linear'] not in rows
                rows[row['ida_linear']] = row
                named_count += 1
                named_bytes += len(bytes.fromhex(row['bytes']))
            for relocation in chunk['loader_relocations']:
                insn = next(i for i in chunk['instructions'] if i['file_offset'] <= relocation['file_offset'] < i['file_offset'] + len(bytes.fromhex(i['bytes'])))
                assert insn['bytes'].startswith('9a'), ('non-far relocation', insn)
                assert relocation['original_word'] == 0x1000 and relocation['ida_word'] == 0x2000
    assert sum(len(c['loader_relocations']) for t in selected for c in t['chunks']) == 2
    for block in probe['decoded_blocks']:
        assert RAW_WANTED[block['ida_linear']] == len(block['instructions'])
        entries.add(block['ida_linear'])
        originals[block['ida_linear']] = f'raw_entry_{block["ida_linear"]:X}'
        for row in block['instructions']:
            assert row['ida_linear'] not in rows
            rows[row['ida_linear']] = row
            raw_count += 1
            raw_bytes += len(bytes.fromhex(row['bytes']))
    assert len(probe['decoded_blocks']) == 1
    assert probe['far_call_target'] == 0x20000 and len(probe['far_call_sites']) == 2
    for row in rows.values():
        if row['mnemonic'] == 'call' and row['operands'][0]['type'] == 7:
            target = row['operands'][0]['addr'] + 0x10000
            if target in rows:
                entries.add(target)
    assert named_count == 247
    assert raw_count == 1
    assert HOOK_ENTRIES <= set(rows)
    lines = ['/* Original native C residual graph; re/134, spec/254. No CPU.Step or guest decoder. */',
             'static int rs_known(uint16_t t) {switch(t) {']
    for ea in sorted(rows):
        lines.append(f'case 0x{ea-0x10000:x}:')
    lines.append('return 1;default:return 0;}}')
    lines.extend(['static void rs_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int rs_skip_first_enter) {',
                  'const uint16_t rs_core_cs=m->cs;int rs_invoke_old=0;m->ip=entry;',
                  'RS_DISPATCH:',
                  'if(m->cs!=rs_core_cs) {assert(!rs_invoke_old);return;}',
                  'switch(m->ip) {'])
    for ea in sorted(rows):
        lines.append(f'case 0x{ea-0x10000:x}:goto L_{ea:X};')
    lines.extend(['default:break;}',
                  'if(!rs_invoke_old)return;',
                  'ki_startup_invoke(m,m->ip,h);rs_invoke_old=0;goto RS_DISPATCH;'])
    seen_mutants = set()
    for ea, row in sorted(rows.items()):
        mn, ops = row['mnemonic'], row['operands']
        next_ea = ea + len(bytes.fromhex(row['bytes']))
        nextoff = next_ea - 0x10000
        lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
        lines.append(f'm->ip=0x{ea-0x10000:x};')
        if ea in entries:
            lines.append(f'/* Original entry: {originals.get(ea, "inner call target")} */')
            lines.append(f'if(rs_skip_first_enter)rs_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x{ea-0x10000:x},h->user);')
        else:
            lines.append('rs_skip_first_enter=0;')
        def jump(target):
            if target in rows:
                return f'goto L_{target:X};'
            return f'm->ip=0x{target-0x10000:x};rs_invoke_old=1;goto RS_DISPATCH;'
        terminal = False
        if mn == 'call':
            if row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);rs_invoke_old=0;goto RS_DISPATCH;'
            else:
                statement = '{uint16_t target=%s;ec_push(m,0x%x);m->ip=target;rs_invoke_old=1;goto RS_DISPATCH;}' % (dl.read(ops[0]), nextoff)
            terminal = True
        elif mn == 'retn':
            assert row['bytes'] == 'c3'
            statement = 'm->ip=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;'
            terminal = True
        elif mn == 'retf':
            statement = 'm->ip=ec_pop(m);m->cs=ec_pop(m);rs_invoke_old=0;goto RS_DISPATCH;'
            terminal = True
        elif mn == 'jmp':
            assert ops[0]['type'] == 7
            statement = jump(ops[0]['addr'] + 0x10000)
            terminal = True
        elif mn in dl.COND:
            statement = f'if({dl.COND[mn]}) {{{jump(ops[0]["addr"]+0x10000)}}}'
        elif mn == 'loop':
            statement = f'if(--m->cx) {{{jump(ops[0]["addr"]+0x10000)}}}'
        elif mn in ['movsw', 'stosw']:
            code = bytes.fromhex(row['bytes'])
            repeat = int(code[0] in [0xf2, 0xf3])
            statement = f'mc_string(m,{dict(movsw=0, stosw=2)[mn]},{repeat});'
        elif mn == 'int':
            assert ops[0]['type'] == 5
            statement = f'rs_interrupt(m,{ops[0]["value"]},0x{nextoff:x},h);'
        elif mn == 'in':
            assert ops[0]['original'].lower() == 'al' and ops[1]['original'].lower() == 'dx'
            statement = 'vg_al(m,(uint8_t)(wolong_residual_in8(m->dx)));'
        elif mn == 'pushf':
            statement = 'ec_push(m,m->flags);'
        elif mn == 'popf':
            statement = 'm->flags=ec_pop(m)|2;'
        elif mn == 'cli':
            statement = 'm->flags&=(uint16_t)~0x200u;'
        else:
            try:
                statement = dl.statement(row)
            except Exception as error:
                raise ValueError(f'{ea:X}: {row["assembly"]}: {error}') from error
        if ea in MUTANT_STATEMENTS:
            number, bad = MUTANT_STATEMENTS[ea]
            assert bad != statement
            seen_mutants.add(ea)
            lines.extend([f"#if KI_RESIDUAL_MUTATION == {number}", bad, "#else", statement, "#endif"])
        else:
            lines.append(statement)
        if not terminal:
            assert next_ea in rows, ('missing fallthrough instruction', hex(ea), hex(next_ea))
            lines.append(f'goto L_{next_ea:X};')
    assert seen_mutants == set(MUTANT_STATEMENTS)
    lines.append('}')
    lines.extend(['int ki_residual_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {',
                  'if(!rs_known(t))return 0;',
                  'uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);',
                  'rs_graph(m,t,h,1); /* Old invoke already emitted the entry hook. */',
                  '/* Old external ABI owns one final POP; reject abandoned native frames. */',
                  'assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);',
                  'ec_push(m,ret);return 1;}',
                  'void ki_residual_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(rs_known(t))rs_graph(m,t,h,0);else ki_startup_invoke(m,t,h);}'])
    declarations = []
    for target in selected:
        name, offset = target['name'], target['ida_linear'] - 0x10000
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{rs_graph(m,0x{offset:x},h,0);}}')
        declarations.append(f'void {name}(KiMachine16 *,const KiEconomyHooks *);')
    for block in probe['decoded_blocks']:
        ea = block['ida_linear']
        name = f'raw_entry_{ea:X}'
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{rs_graph(m,0x{ea-0x10000:x},h,0);}}')
        declarations.append(f'void {name}(KiMachine16 *,const KiEconomyHooks *);')
    (output / 'residual_generated.inc').write_text('\n'.join(lines)+'\n', encoding='utf-8')
    (output / 'residual.h').write_text('\n'.join(['#ifndef WOLONG_C_RESIDUAL_H', '#define WOLONG_C_RESIDUAL_H', '#include "startup.h"',
        'uint8_t wolong_residual_in8(uint16_t);',
        '_Noreturn void wolong_residual_exit(KiMachine16 *,uint8_t,void *);',
        'int ki_residual_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);', 'void ki_residual_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);', *declarations, '#endif'])+'\n', encoding='utf-8')
    print(f'Generated {len(selected)} named / {len(probe["decoded_blocks"])} raw / {len(rows)} instructions; {len(MUTANT_STATEMENTS)} controls')
    return {'named_functions': len(selected), 'named_instructions': named_count,
            'raw_instructions': raw_count, 'unique_instructions': len(rows)}


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', required=True, type=Path)
    parser.add_argument('--probe', required=True, type=Path)
    parser.add_argument('--output', type=Path)
    args = parser.parse_args()
    generate(args.repo, args.probe, args.output)
