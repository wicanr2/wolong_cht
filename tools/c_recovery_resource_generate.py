#!/usr/bin/env python3
"""原始 resource／cue／allocation 到 native C，平台中斷保持原 ABI。"""
import argparse
import importlib.util
import json
from pathlib import Path


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('display', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec); spec.loader.exec_module(dl)
    report = json.loads(probe.read_text())
    assert report['input_sha256'] == 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    groups = [(t['ida_linear'], [r for c in t['chunks'] for r in c['instructions']], t['name']) for t in report['targets']]
    lines = ['/* Fixed IDA native C resource control; spec/233, re/113. Original locators retained. */']
    for ea, rows, name in groups:
        labels = {r['ida_linear'] for r in rows}
        lines += [f'/* Original {name}; IDA 0x{ea:X}. */', f'static void rs_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
        for row in rows:
            at, mn = row['ida_linear'], row['mnemonic']
            nextoff = (at + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mn in list(dl.COND) + ['jmp', 'loop']:
                assert row['operands'][0]['addr'] + 0x10000 in labels
            if mn == 'int': statement = f'gl_interrupt(m,{row["operands"][0]["value"]},{nextoff});'
            elif mn == 'lahf': statement = 'vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
            elif mn == 'sahf': statement = 'm->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
            else: statement = dl.statement(row).replace('dl_call(', 'rs_call(')
            mutation, bad = 0, None
            if at == 0x1F4B7: mutation, bad = 1, statement.replace('(61440)', '(32768)')
            elif at == 0x1F4C8: mutation, bad = 2, statement.replace(',3840,', ',4096,')
            elif at == 0x1F4C6: mutation, bad = 3, statement.replace('if(m->flags&1)', 'if(m->flags&(1|0x40))')
            elif at == 0x187C2: mutation, bad = 4, statement.replace(',2048,', ',1024,')
            elif at == 0x10246: mutation, bad = 5, '; /* lost cue same-value gate */'
            elif at == 0x10279: mutation, bad = 6, '; /* lost cue offset high-byte clear */'
            elif at == 0x10281: mutation, bad = 7, statement.replace('dl_read16(m,m->ds,(uint16_t)(m->bx))', 'dl_read16(m,m->ds,(uint16_t)(m->bx+2))')
            elif at == 0x102B4: mutation, bad = 8, statement.replace('if(m->flags&0x40)', 'if(!(m->flags&0x40))')
            elif at == 0x100ED: mutation, bad = 9, statement.replace(',3328,', ',3327,')
            elif at == 0x1011F: mutation, bad = 10, statement.replace(',24,', ',23,')
            elif at == 0x102C3: mutation, bad = 11, statement.replace(',255);', ',0);')
            elif at == 0x1F4D7: mutation, bad = 12, statement.replace('(65535)', '(0)')
            lines.append(f'L_{at:X}: ; /* IDA 0x{at:X}; {row["assembly"]} */')
            if bad is not None:
                assert bad != statement, (hex(at), statement)
                lines += [f'#if KI_RESOURCE_MUTATION == {mutation}', '    ' + bad, '#else', '    ' + statement, '#endif']
            else: lines.append('    ' + statement)
        lines.append('}')
    lines.append('int ki_resource_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for ea, rows, name in groups: lines.append(f'case 0x{ea - 0x10000:x}:rs_body_{ea:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    (repo / 'tools/c_recovery/resource_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated', len(groups), 'entries;', sum(len(g[1]) for g in groups), 'instructions')


if __name__ == '__main__':
    p = argparse.ArgumentParser(description=__doc__); p.add_argument('--repo', type=Path, required=True)
    p.add_argument('--probe', type=Path, required=True); a = p.parse_args(); generate(a.repo, a.probe)
