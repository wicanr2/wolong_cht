#!/usr/bin/env python3
"""兩個原始 code segment 的 native C、far table 與已知 live patch。"""
import argparse
import importlib.util
import json
from pathlib import Path


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('display', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    report = json.loads(probe.read_text())
    targets = {t['ida_linear']: t for t in report['targets']}
    groups = []
    for target in report['targets']:
        ea = target['ida_linear']
        if ea == 0x2006F:
            continue
        rows = [r for c in target['chunks'] for r in c['instructions']]
        aliases = [ea]
        if ea == 0x2002E:
            rows += [r for c in targets[0x2006F]['chunks'] for r in c['instructions']]
            aliases.append(0x2006F)
        groups.append((ea, rows, target['name'], aliases))
    groups += [(b['ida_linear'], b['instructions'], b['original_name_before_analysis'], [b['ida_linear']])
               for b in report['decoded_blocks']]
    lines = ['/* Fixed IDA native C input; original segment locators; spec/231, re/111. */']
    for ea, rows, name, aliases in groups:
        base = 0x20000 if ea >= 0x20000 else 0x10000
        labels = {r['ida_linear'] for r in rows}
        lines += [f'/* Original {name or "(raw code)"}; IDA 0x{ea:X}. */',
                  f'static void mx_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h,uint32_t entry) {{',
                  'switch(entry) {' + ' '.join(f'case 0x{a:x}:goto L_{a:X};' for a in aliases) + ' default:assert(0);}']
        for row in rows:
            at, mn = row['ida_linear'], row['mnemonic']
            nextoff = (at + len(bytes.fromhex(row['bytes'])) - base) & 65535
            if at == 0x12239:
                statement = '{uint16_t patch=dl_read16(m,m->cs,0x2239);if(patch==0x0feb) goto L_1224A;assert(patch==0x9090);}'
            elif mn in dl.COND or mn in ['jmp', 'loop']:
                target = row['operands'][0]['addr'] + base
                assert target in labels, (hex(at), hex(target))
                condition = dl.COND.get(mn, '--m->cx' if mn == 'loop' else None)
                statement = (f'if({condition}) ' if condition else '') + f'goto L_{target:X};'
            elif mn == 'call' and bytes.fromhex(row['bytes'])[0] == 0x9A:
                statement = f'mx_far(m,0x{at - base:x},0x{nextoff:x},h);'
            elif mn == 'call':
                statement = f'mx_call(m,{dl.read(row["operands"][0])},{nextoff},0x{base:x},h);'
            elif mn in ['retn', 'retf']:
                statement = 'return;'
            elif mn == 'cli':
                statement = 'm->flags&=(uint16_t)~0x200u;'
            elif mn == 'sti':
                statement = 'm->flags|=0x200;'
            elif mn == 'pushf':
                statement = 'ec_push(m,m->flags);'
            elif mn == 'popf':
                statement = 'm->flags=ec_pop(m);'
            elif mn == 'int':
                statement = f'gl_interrupt(m,{row["operands"][0]["value"]},{nextoff});'
            elif mn == 'in':
                statement = dl.write(row['operands'][0], 'wolong_input_in(' + dl.read(row['operands'][1]) + ')')
            else:
                statement = dl.statement(row)
            mutation, bad = 0, None
            if at == 0x121EC:
                mutation, bad = 1, statement.replace('(1)', '(0)')
            elif at == 0x12239:
                mutation, bad = 2, '{uint16_t patch=dl_read16(m,m->cs,0x2239);assert(patch==0x0feb||patch==0x9090);}'
            elif at == 0x12236:
                mutation, bad = 3, statement.replace('(10)', '(9)')
            elif at == 0x12279:
                mutation, bad = 4, '; /* lost pending-left cleanup */'
            elif at == 0x201E4:
                mutation, bad = 5, statement.replace('(795)', '(827)')
            elif at == 0x20212:
                mutation, bad = 6, statement.replace('if(!(m->flags&1))', 'if(m->flags&1)')
            elif at == 0x20223:
                mutation, bad = 7, statement.replace(',632,', ',631,')
            elif at == 0x2024C:
                mutation, bad = 8, statement.replace(',616,', ',624,')
            elif at in [0x2021F, 0x20229, 0x20235]:
                mutation, bad = 9, '; /* lost cursor latch read */'
            elif at == 0x2000C:
                mutation, bad = 10, '; /* lost original STI */'
            elif at == 0x20135:
                mutation, bad = 11, '; /* lost callback POPF */'
            elif at == 0x18847:
                mutation, bad = 12, '; /* lost message close marker */'
            lines.append(f'L_{at:X}: ; /* IDA 0x{at:X}; {row["assembly"]} */')
            if at == 0x12259:
                lines.append('    m->ip=0x2259;if(h&&h->enter) h->enter(m,0x2259,h->user);')
            if at == 0x2006F:
                lines.append('    if(entry!=0x2006f) {m->ip=0x6f;if(h&&h->enter) h->enter(m,0x6f,h->user);}')
            if bad is not None:
                assert bad != statement, (hex(at), statement)
                lines += [f'#if KI_INPUT_MUTATION == {mutation}', '    ' + bad, '#else', '    ' + statement, '#endif']
            else:
                lines.append('    ' + statement)
        lines.append('}')
    lines.append('int ki_input_body(KiMachine16 *m,uint32_t loc,const KiEconomyHooks *h) {switch(loc) {')
    for ea, rows, name, aliases in groups:
        for a in aliases:
            lines.append(f'case 0x{a:x}:mx_body_{ea:X}(m,h,loc);break;')
    lines.append('default:return 0;}return 1;}')
    (repo / 'tools/c_recovery/input_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated', sum(len(g[3]) for g in groups), 'entries;', sum(len(g[1]) for g in groups), 'instructions')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
