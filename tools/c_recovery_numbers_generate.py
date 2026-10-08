#!/usr/bin/env python3
"""固定 IDA 數字指令產生 native C；不在執行期讀取 opcode。"""
import argparse
import importlib.util
import json
from pathlib import Path


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('display', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    report = json.loads(probe.read_text())
    assert report['input_sha256'] == 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    targets = {t['ida_linear']: t for t in report['targets']}
    groups = []
    for ea in [0x1062F, 0x1069A, 0x106DE]:
        target = targets[ea]
        groups.append((ea, [r for c in target['chunks'] for r in c['instructions']], target['name']))
    groups.extend((b['ida_linear'], b['instructions'], b['original_name_before_analysis'])
                  for b in report['decoded_blocks'])
    lines = ['/* Fixed IDA native C number control; spec/229, re/109. Original locators retained. */']
    for ea, rows, name in groups:
        labels = {r['ida_linear'] for r in rows}
        lines += [f'/* Original {name}; IDA 0x{ea:X}. */',
                  f'static void nr_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
        for row in rows:
            at, mnemonic = row['ida_linear'], row['mnemonic']
            if mnemonic in list(dl.COND) + ['jmp', 'loop']:
                assert row['operands'][0]['addr'] + 0x10000 in labels
            if mnemonic == 'lodsb':
                statement = 'vg_al(m,dl_read8(m,m->ds,m->si));m->si=(uint16_t)(m->si+((m->flags&0x400)?-1:1));'
            elif mnemonic == 'stosb':
                statement = 'dl_write8(m,m->es,m->di,(uint8_t)m->ax);m->di=(uint16_t)(m->di+((m->flags&0x400)?-1:1));'
            elif mnemonic == 'cwd':
                statement = 'm->dx=(m->ax&0x8000)?0xffff:0;'
            else:
                statement = dl.statement(row).replace('dl_call(', 'nr_call(')
            mutation, bad = 0, None
            if at == 0x10655:
                mutation, bad = 1, '; /* lost carry for negative multiples of 65536 */'
            elif at == 0x10644:
                mutation, bad = 2, '; /* lost right-edge width adjustment */'
            elif at == 0x1063B:
                mutation, bad = 3, statement.replace('3412', '3404')
            elif at == 0x106AA:
                mutation, bad = 4, statement.replace('&(15)', '&(0)')
            elif at in [0x106AF, 0x106CC]:
                mutation, bad = 5, '; /* lost latch read */'
            elif at == 0x106E9:
                mutation, bad = 6, statement.replace('((m->flags&0x400)?-1:1)', '1')
            elif at in [0x106B3, 0x106D3]:
                mutation, bad = 7, statement.replace(',79,', ',78,')
            elif at == 0x10674:
                mutation, bad = 8, statement.replace('(m->bx>>8)', '(m->ax>>8)')
            lines.append(f'L_{at:X}: ; /* IDA 0x{at:X}; {row["assembly"]} */')
            if bad is not None:
                assert bad != statement, (hex(at), statement)
                lines += [f'#if KI_NUMBERS_MUTATION == {mutation}', '    ' + bad,
                          '#else', '    ' + statement, '#endif']
            else:
                lines.append('    ' + statement)
        lines.append('}')
    lines.append('int ki_numbers_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for ea, rows, name in groups:
        lines.append(f'case 0x{ea - 0x10000:x}:nr_body_{ea:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    (repo / 'tools/c_recovery/numbers_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated', len(groups), 'entries;', sum(len(g[1]) for g in groups), 'original instructions')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
