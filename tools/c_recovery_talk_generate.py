#!/usr/bin/env python3
"""固定 IDA TALK 指令產生 native C，保留原始 near table 與 DOS ABI。"""
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
    groups = [(t['ida_linear'], [r for c in t['chunks'] for r in c['instructions']], t['name'])
              for t in report['targets']]
    groups.extend((b['ida_linear'], b['instructions'], b['original_name_before_analysis'])
                  for b in report['decoded_blocks'])
    lines = ['/* Fixed IDA native C TALK control; original locators; spec/230, re/110. */']
    for ea, rows, name in groups:
        labels = {r['ida_linear'] for r in rows}
        lines += [f'/* Original {name or "(unnamed code)"}; IDA 0x{ea:X}. */',
                  f'static void tx_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
        for row in rows:
            at, mnemonic = row['ida_linear'], row['mnemonic']
            nextoff = (at + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mnemonic in list(dl.COND) + ['jmp', 'loop']:
                assert row['operands'][0]['addr'] + 0x10000 in labels
            if mnemonic == 'int':
                statement = f'gl_interrupt(m,{row["operands"][0]["value"]},{nextoff});'
            elif mnemonic == 'lahf':
                statement = 'vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
            elif mnemonic == 'sahf':
                statement = 'm->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
            elif mnemonic == 'rcl':
                assert row['operands'][1]['value'] == 1
                statement = dl.write(row['operands'][0], 'tx_rcl(m,' + dl.read(row['operands'][0]) + ',' + str(row['operands'][0]['dtype_size'] * 8) + ')')
            else:
                statement = dl.statement(row).replace('dl_call(', 'tx_call(')
            mutation, bad = 0, None
            if at == 0x10771:
                mutation, bad = 1, '; /* wrong variant stride */'
            elif at == 0x1086A:
                mutation, bad = 2, '; /* lost half-width rewind */'
            elif at == 0x10881:
                mutation, bad = 3, statement.replace(',16,', ',8,')
            elif at == 0x1089F:
                mutation, bad = 4, statement.replace(',48,', ',32,')
            elif at == 0x108D0:
                mutation, bad = 5, statement.replace(',8,', ',2,')
            elif at == 0x107FB:
                mutation, bad = 6, statement.replace('&(3)', '&(1)')
            elif at == 0x1081D:
                mutation, bad = 7, '; /* wrong portrait offset stride */'
            elif at == 0x10821:
                mutation, bad = 8, statement.replace('(2048)', '(1024)')
            elif at == 0x10980:
                mutation, bad = 9, '; /* lost no-output marker cursor cancellation */'
            elif at == 0x10773:
                mutation, bad = 10, '; /* lost variant byte addition */'
            lines.append(f'L_{at:X}: ; /* IDA 0x{at:X}; {row["assembly"]} */')
            if bad is not None:
                assert bad != statement, (hex(at), statement)
                lines += [f'#if KI_TALK_MUTATION == {mutation}', '    ' + bad,
                          '#else', '    ' + statement, '#endif']
            else:
                lines.append('    ' + statement)
        lines.append('}')
    lines.append('int ki_talk_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for ea, rows, name in groups:
        lines.append(f'case 0x{ea - 0x10000:x}:tx_body_{ea:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    (repo / 'tools/c_recovery/talk_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated', len(groups), 'entries;', sum(len(g[1]) for g in groups), 'original instructions')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
