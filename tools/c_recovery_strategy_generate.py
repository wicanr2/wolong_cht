#!/usr/bin/env python3
"""固定 IDA 政略29函式轉 C，保留原始定位與 byte 語意。"""
import argparse
import importlib.util
import json
from pathlib import Path


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    data = json.loads(probe.read_text())
    assert data['input_sha256'] == 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    groups = [(t['ida_linear'], [x for c in t['chunks'] for x in c['instructions']], t['name']) for t in data['targets']]
    assert len(groups) == 29 and not data['decoded_blocks']
    lines = ['/* Original native C strategy; re/124, spec/244. */']
    mutants = {
        0x161CC: (1, 'm->ax=ec_math(m,m->ax,23,16,0,1);'),
        0x1622E: (2, 'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+0x4257)),1,8);'),
        0x1307B: (3, 'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+3)),(uint8_t)(m->dx>>8),8);'),
        0x13C26: (4, 'ec_cmp(m,(uint8_t)m->dx,0xe1,8);'),
        0x13BE4: (5, '; /* lost reason budget decrement */'),
        0x13BC9: (6, '; /* lost attempted reason bit */'),
        0x13BF0: (7, 'm->flags&=(uint16_t)~1u;'),
        0x13DA1: (8, 'dl_write8(m,m->cs,0xd00,0xfe);'),
        0x16854: (9, 'm->ax=(uint8_t)m->ax;'),
        0x16A37: (10, 'vg_ah(m,ec_math(m,(uint8_t)(m->ax>>8),19,8,0,0));'),
        0x16498: (11, 'vg_al(m,ec_math(m,(uint8_t)m->ax,0x7f,8,0,1));'),
        0x16592: (12, 'ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->si+0x28)),8);'),
        0x166E5: (13, 'ec_cmp(m,m->di,m->di,16);'),
        0x163A0: (14, 'if(m->flags&1) goto L_163A4;'),
        0x163BC: (15, 'm->es=ec_pop(m);'),
        0x13890: (16, '; /* lost accepted-reason trust halving */'),
    }
    seen = set()
    for at, rows, name in groups:
        labels = {x['ida_linear'] for x in rows}
        lines.extend([f'/* Original {name}; IDA 0x{at:X}. */', f'static void sty_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
        for row in rows:
            ea, mn = row['ida_linear'], row['mnemonic']
            nextoff = (ea + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mn == 'call' and row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
            elif mn == 'call':
                statement = f'sty_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            elif mn == 'cbw':
                statement = 'm->ax=(uint16_t)(int16_t)(int8_t)m->ax;'
            elif mn == 'int':
                statement = f'gl_interrupt(m,{dl.read(row["operands"][0])},0x{nextoff:x});'
            elif mn == 'pushf':
                statement = 'ec_push(m,m->flags);'
            elif mn == 'popf':
                statement = 'm->flags=ec_pop(m)|2;'
            elif mn == 'lahf':
                statement = 'vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
            elif mn == 'sahf':
                statement = 'm->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
            else:
                if mn in list(dl.COND) + ['jmp', 'loop']:
                    assert row['operands'][0]['addr'] + 0x10000 in labels
                statement = dl.statement(row)
                if mn == 'jmp' and row['operands'][0]['addr'] + 0x10000 == at:
                    statement = f'm->ip=0x{at-0x10000:x};if(h&&h->enter)h->enter(m,m->ip,h->user);' + statement
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if ea in mutants:
                number, bad = mutants[ea]
                assert bad != statement
                seen.add(ea)
                lines.extend([f'#if KI_STRATEGY_MUTATION == {number}', '    '+bad, '#else', '    '+statement, '#endif'])
            else:
                lines.append('    '+statement)
        assert rows[-1]['mnemonic'] == 'retn' or at == 0x16366
        lines.append('}')
    assert seen == set(mutants)
    lines.append('int ki_strategy_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at, rows, name in groups:
        lines.append(f'case 0x{at-0x10000:x}:sty_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    lines.append('static int sty_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at, rows, name in groups)+' return 1;default:return 0;}}')
    for at, rows, name in groups:
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_strategy_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
    (repo / 'tools/c_recovery/strategy_generated.inc').write_text('\n'.join(lines)+'\n')
    print('Generated 29 functions /', sum(len(x[1]) for x in groups), 'instructions')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
