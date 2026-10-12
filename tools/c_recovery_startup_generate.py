#!/usr/bin/env python3
"""原開機、主排程與命名UI靜態C圖；保留guest stack、MZ段值與原始入口。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

PROBE_SHA = '3d3761a281fa2ef1b50d1e8088e29053335ec7317473448e5624681e87cd8ea7'
INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
WANTED = {0x10000, 0x1005B, 0x1006B, 0x10155, 0x10210, 0x10236, 0x1033B, 0x1034E,
          0x10612, 0x11A6E, 0x11AAE, 0x11AC3, 0x11BE0, 0x11CD0, 0x18B12, 0x18B40,
          0x18CAE, 0x18DC8, 0x18E5A, 0x18EA0, 0x18F37, 0x18F6D, 0x18F7C, 0x18FC9,
          0x1905E, 0x1906E, 0x1907F, 0x1908D, 0x190A3, 0x190C0, 0x1912D, 0x19144,
          0x1915B, 0x1917C, 0x19214, 0x19223, 0x1928A, 0x192E3, 0x192EE, 0x192F5,
          0x1EB6C, 0x1EB93}
RAW_WANTED = {0x10043: 11, 0x1006A: 1, 0x10356: 8}
PAUSE_ENTRIES = {0x11C1F}
MUTANT_STATEMENTS = {
    0x1002D: (1, 'vg_ah(m,(uint8_t)(76));vg_al(m,(uint8_t)(1));'),
    0x10024: (2, 'm->dx=(uint16_t)(6567);'),
    0x10217: (3, 'vg_ah(m,(uint8_t)(3));'),
    0x1EB78: (4, 'vg_al(m,(uint8_t)(21));'),
    0x10345: (5, 'm->ax=(uint16_t)(3073);'),
    0x19037: (6, '{uint16_t target=dl_read16(m,m->cs,(uint16_t)(m->bx+36761));ec_push(m,0x903c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}'),
    0x192E5: (7, 'ec_cmp(m,m->ax,5241,16);'),
    0x19136: (8, 'vg_al(m,(uint8_t)(147));'),
    0x1035E: (9, 'ec_cmp(m,dl_read8(m,m->cs,3373),254,8);'),
    0x11C5B: (10, 'ec_cmp(m,dl_read8(m,m->ds,39075),129,8);'),
    0x11BEC: (11, 'dl_write16(m,m->cs,39171,m->sp);'),
    0x10051: (12, 'm->ax=(uint16_t)(2);'),
    0x10226: (13, 'dl_write8(m,m->cs,526,254);'),
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
    assert len(selected) == 42
    rows = {}
    named_count = named_bytes = raw_count = raw_bytes = 0
    entries = set(WANTED) | set(RAW_WANTED) | set(PAUSE_ENTRIES)
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
    assert sum(len(c['loader_relocations']) for t in selected for c in t['chunks']) == 15
    for block in probe['decoded_blocks']:
        assert RAW_WANTED[block['ida_linear']] == len(block['instructions'])
        entries.add(block['ida_linear'])
        originals[block['ida_linear']] = f'raw_entry_{block["ida_linear"]:X}'
        for row in block['instructions']:
            assert row['ida_linear'] not in rows
            rows[row['ida_linear']] = row
            raw_count += 1
            raw_bytes += len(bytes.fromhex(row['bytes']))
    assert len(probe['decoded_blocks']) == 3
    for table in probe['fixed_dispatch_tables']:
        assert table['all_entries_in_recovery']
        assert rows[table['call_ida_linear']]['mnemonic'] == 'call'
        entries.update(0x10000 + word for word in table['words'])
    assert probe['far_call_target'] == 0x20000 and len(probe['far_call_sites']) == 15
    for row in rows.values():
        if row['mnemonic'] == 'call' and row['operands'][0]['type'] == 7:
            target = row['operands'][0]['addr'] + 0x10000
            if target in rows:
                entries.add(target)
    assert PAUSE_ENTRIES <= set(rows)
    assert named_count == 1021
    lines = ['/* Original native C startup graph; re/133, spec/253. No CPU.Step or guest decoder. */',
             'static int su_known(uint16_t t) {switch(t) {']
    for ea in sorted(rows):
        lines.append(f'case 0x{ea-0x10000:x}:')
    lines.append('return 1;default:return 0;}}')
    lines.extend(['static void su_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int su_skip_first_enter) {',
                  'const uint16_t su_core_cs=m->cs;int su_invoke_old=0;m->ip=entry;',
                  'SU_DISPATCH:',
                  'if(m->cs!=su_core_cs) {assert(!su_invoke_old);return;}',
                  'switch(m->ip) {'])
    for ea in sorted(rows):
        lines.append(f'case 0x{ea-0x10000:x}:goto L_{ea:X};')
    lines.extend(['default:break;}',
                  'if(!su_invoke_old)return;',
                  'ki_army_invoke(m,m->ip,h);su_invoke_old=0;goto SU_DISPATCH;'])
    seen_mutants = set()
    for ea, row in sorted(rows.items()):
        mn, ops = row['mnemonic'], row['operands']
        next_ea = ea + len(bytes.fromhex(row['bytes']))
        nextoff = next_ea - 0x10000
        lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
        lines.append(f'm->ip=0x{ea-0x10000:x};')
        if ea in entries:
            lines.append(f'/* Original entry: {originals.get(ea, "inner pause entry")} */')
            lines.append(f'if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x{ea-0x10000:x},h->user);')
        else:
            lines.append('su_skip_first_enter=0;')
        def jump(target):
            if target in rows:
                return f'goto L_{target:X};'
            return f'm->ip=0x{target-0x10000:x};su_invoke_old=1;goto SU_DISPATCH;'
        terminal = False
        if mn == 'call':
            if row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);su_invoke_old=0;goto SU_DISPATCH;'
            else:
                statement = '{uint16_t target=%s;ec_push(m,0x%x);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}' % (dl.read(ops[0]), nextoff)
            terminal = True
        elif mn == 'retn':
            assert row['bytes'] == 'c3'
            statement = 'm->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;'
            terminal = True
        elif mn == 'retf':
            statement = 'm->ip=ec_pop(m);m->cs=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;'
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
            statement = f'su_interrupt(m,{ops[0]["value"]},0x{nextoff:x},h);'
        elif mn == 'in':
            assert ops[0]['original'].lower() == 'al' and ops[1]['original'].lower() == 'dx'
            statement = 'vg_al(m,(uint8_t)(wolong_startup_in8(m->dx)));'
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
            lines.extend([f"#if KI_STARTUP_MUTATION == {number}", bad, "#else", statement, "#endif"])
        else:
            lines.append(statement)
        if not terminal:
            assert next_ea in rows, ('missing fallthrough instruction', hex(ea), hex(next_ea))
            lines.append(f'goto L_{next_ea:X};')
    assert seen_mutants == set(MUTANT_STATEMENTS)
    lines.append('}')
    lines.extend(['int ki_startup_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {',
                  'if(!su_known(t))return 0;',
                  'uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);',
                  'su_graph(m,t,h,1); /* Old invoke already emitted the entry hook. */',
                  '/* Old external ABI owns one final POP; reject abandoned native frames. */',
                  'if(t==0x356){uint16_t css=dl_read16(m,ss,(uint16_t)(sp+2));assert(m->ss==ss&&m->sp==(uint16_t)(sp+4)&&m->ip==ret&&m->cs==css);ec_push(m,css);}',
                  'else assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);',
                  'ec_push(m,ret);return 1;}',
                  'void ki_startup_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(su_known(t))su_graph(m,t,h,0);else ki_army_invoke(m,t,h);}'])
    declarations = []
    for target in selected:
        name, offset = target['name'], target['ida_linear'] - 0x10000
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{su_graph(m,0x{offset:x},h,0);}}')
        declarations.append(f'void {name}(KiMachine16 *,const KiEconomyHooks *);')
    for block in probe['decoded_blocks']:
        ea = block['ida_linear']
        name = f'raw_entry_{ea:X}'
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{su_graph(m,0x{ea-0x10000:x},h,0);}}')
        declarations.append(f'void {name}(KiMachine16 *,const KiEconomyHooks *);')
    (output / 'startup_generated.inc').write_text('\n'.join(lines)+'\n', encoding='utf-8')
    (output / 'startup.h').write_text('\n'.join(['#ifndef WOLONG_C_STARTUP_H', '#define WOLONG_C_STARTUP_H', '#include "army.h"',
        'uint8_t wolong_startup_in8(uint16_t);',
        '_Noreturn void wolong_startup_exit(KiMachine16 *,uint8_t,void *);',
        'int ki_startup_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);', 'void ki_startup_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);', *declarations, '#endif'])+'\n', encoding='utf-8')
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
