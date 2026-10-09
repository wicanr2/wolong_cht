#!/usr/bin/env python3
"""固定 IDA 三函式退出閉包轉 C，保留原始 RET 與 SS/SP 語意。"""
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
    wanted = set(data['recovery_targets'])
    assert wanted == {0x11CB1, 0x10A1C, 0x1EBDC}
    groups = [(t['ida_linear'], [x for c in t['chunks'] for x in c['instructions']], t['name'])
              for t in data['targets'] if t['ida_linear'] in wanted]
    lines = ['/* Original native C nonlocal exit; re/125, spec/245. Navigation-only main loop is not generated. */']
    mutants = {
        0x11CB1: (1, 'm->ss=dl_read16(m,m->cs,0x9901);'),
        0x11CB6: (2, 'm->sp=dl_read16(m,m->cs,0x9903);'),
        0x11CBB: (3, '; /* lost AX push on restored outer stack */'),
        0x11CC7: (4, 'm->ax=ec_pop(m)^1;'),
        0x10A21: (5, 'vg_cl(m,15);'),
        0x1EBF2: (6, 'm->ax=ec_math(m,m->ax,0,16,0,0);'),
        0x1EC25: (7, 'dl_out(m,m->dx,(uint8_t)(m->bx>>8),8);'),
    }
    seen = set()
    for at, rows, name in groups:
        labels = {x['ida_linear'] for x in rows}
        lines.extend([f'/* Original {name}; IDA 0x{at:X}. */', f'static void nle_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
        for row in rows:
            ea, mn = row['ida_linear'], row['mnemonic']
            nextoff = (ea + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mn == 'call' and row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
            elif mn == 'call':
                statement = f'nle_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            elif mn == 'in':
                assert row['operands'][0]['dtype_size'] == 1
                statement = dl.write(row['operands'][0], f'wolong_main_in8({dl.read(row["operands"][1])})')
            else:
                if mn in list(dl.COND) + ['jmp', 'loop']:
                    assert row['operands'][0]['addr'] + 0x10000 in labels
                statement = dl.statement(row)
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if ea in mutants:
                number, bad = mutants[ea]
                assert bad != statement
                seen.add(ea)
                lines.extend([f'#if KI_MAIN_MUTATION == {number}', '    ' + bad, '#else', '    ' + statement, '#endif'])
            else:
                lines.append('    ' + statement)
        assert rows[-1]['mnemonic'] == 'retn'
        lines.append('}')
    assert seen == set(mutants)
    lines.append('int ki_main_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at, rows, name in groups:
        lines.append(f'case 0x{at-0x10000:x}:nle_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}')
    lines.extend(['if(t==0x1cb1) {', 'm->ip=ec_pop(m); /* Original RET consumes the restored outer frame. */',
                  '#if KI_MAIN_MUTATION == 8', 'return 1; /* Incorrectly resume old native caller after virtual nonlocal RET. */',
                  '#else', 'wolong_main_transfer(m,h?h->user:0);', '#endif', '}', 'return 1;}'])
    lines.append('static int nle_known(uint16_t t) {switch(t) {' + ' '.join(f'case 0x{at-0x10000:x}:' for at, rows, name in groups) + ' return 1;default:return 0;}}')
    for at, rows, name in groups:
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_main_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
    (repo / 'tools/c_recovery/main_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated 3 functions /', sum(len(x[1]) for x in groups), 'instructions')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
