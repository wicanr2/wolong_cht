#!/usr/bin/env python3
"""原始 selector／row callback 到 native C，保留 patched immediate 與 far ABI。"""
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
    groups = [(t['ida_linear'], [r for c in t['chunks'] for r in c['instructions']], t['name']) for t in report['targets']]
    groups += [(b['ida_linear'], b['instructions'], b['original_name_before_analysis']) for b in report['decoded_blocks']]
    lines = ['/* Fixed IDA native C choice control; original locators; spec/232, re/112. */']
    for ea, rows, name in groups:
        labels = {r['ida_linear'] for r in rows}
        lines += [f'/* Original {name or "(raw code)"}; IDA 0x{ea:X}. */', f'static void ch_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
        for row in rows:
            at, mn = row['ida_linear'], row['mnemonic']
            nextoff = (at + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mn in list(dl.COND) + ['jmp', 'loop']:
                assert row['operands'][0]['addr'] + 0x10000 in labels
            if mn == 'call' and bytes.fromhex(row['bytes'])[0] == 0x9a:
                statement = f'mx_far(m,0x{at - 0x10000:x},0x{nextoff:x},h);'
            elif mn == 'call':
                statement = f'ch_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            elif mn == 'lahf':
                statement = 'vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
            elif mn == 'sahf':
                statement = 'm->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
            elif at in [0x19440, 0x19442, 0x19464]:
                statement = dl.write(row['operands'][0], f'dl_read8(m,m->cs,0x{at - 0x10000 + 1:x})')
            elif at == 0x19466:
                statement = 'm->si=dl_read16(m,m->cs,0x9467);'
            else:
                statement = dl.statement(row)
            mutation, bad = 0, None
            if at == 0x1942F:
                mutation, bad = 1, '; /* lost row stride doubling */'
            elif at == 0x19466:
                mutation, bad = 2, 'm->si=0x1234;'
            elif at == 0x10464:
                mutation, bad = 3, statement.replace('if(!(m->flags&(1|0x40)))', 'if(!(m->flags&1))')
            elif at == 0x1059E:
                mutation, bad = 4, statement.replace('dl_read8(m,m->ss,(uint16_t)(m->bp))', 'dl_read8(m,m->ss,(uint16_t)(m->bp+1))')
            elif at == 0x104A2:
                mutation, bad = 5, 'm->flags|=1;'
            elif at == 0x104B3:
                mutation, bad = 6, 'm->flags&=(uint16_t)~1u;'
            elif at == 0x10BB2 or at == 0x10BBF or at == 0x10BC7:
                if at == 0x10BBF:
                    statement = dl.statement(row)
                mutation, bad = 7, '; /* lost XOR latch read */'
            elif at == 0x104FA:
                mutation, bad = 8, '; /* lost old DX restore */'
            elif at == 0x194BD or at == 0x1950D:
                mutation, bad = 9, '; /* lost selector carry restore */'
            elif at == 0x10586:
                mutation, bad = 10, '; /* lost absolute selected-row decrement */'
            elif at == 0x105D8:
                mutation, bad = 11, '; /* lost absolute selected-row increment */'
            elif at == 0x19440:
                mutation, bad = 12, 'vg_dl(m,9);'
            lines.append(f'L_{at:X}: ; /* IDA 0x{at:X}; {row["assembly"]} */')
            if bad is not None:
                assert bad != statement, (hex(at), statement)
                lines += [f'#if KI_CHOICE_MUTATION == {mutation}', '    ' + bad, '#else', '    ' + statement, '#endif']
            else:
                lines.append('    ' + statement)
        lines.append('}')
    lines.append('int ki_choice_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for ea, rows, name in groups:
        lines.append(f'case 0x{ea - 0x10000:x}:ch_body_{ea:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    (repo / 'tools/c_recovery/choice_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated', len(groups), 'entries;', sum(len(g[1]) for g in groups), 'instructions')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
