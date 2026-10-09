#!/usr/bin/env python3
"""原始世界顯示格閉包到 native C；保留全部位址與運算元。"""
import argparse
import importlib.util
import json
from pathlib import Path

WANTED = (0x1D46A, 0x1D483, 0x1D4C7, 0x1D615, 0x1D66A,
          0x1D76B, 0x1D782, 0x1D796, 0x1D7E7, 0x1D804)

def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('display', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec); spec.loader.exec_module(dl)
    report = json.loads(probe.read_text()); targets = {t['ida_linear']: t for t in report['targets']}
    assert report['input_sha256'] == 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    lines = ['/* Original map-cell native C; re/114, spec/234. No opcode fetch. */']
    count = 0
    for at in WANTED:
        t = targets[at]; rows = [r for c in t['chunks'] for r in c['instructions']]
        assert all(x in WANTED for x in t['direct_callees'])
        labels = {r['ida_linear'] for r in rows}; count += len(rows)
        lines += [f'/* Original {t["name"]}; IDA 0x{at:X}. */',
                  f'static void mc_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
        for row in rows:
            ea, mn = row['ida_linear'], row['mnemonic']
            if mn in list(dl.COND) + ['jmp', 'loop']:
                assert row['operands'][0]['addr'] + 0x10000 in labels
            if mn in ['movsw', 'lodsw', 'stosw']:
                statement = f'mc_string(m,{ {"movsw": 0, "lodsw": 1, "stosw": 2}[mn]},{int(row["assembly"].lstrip().startswith("rep "))});'
            else: statement = dl.statement(row).replace('dl_call(', 'mc_call(')
            mutation, bad = 0, None
            if ea == 0x1D65E: mutation, bad = 1, statement.replace(',344,', ',345,')
            elif ea == 0x1D655: mutation, bad = 2, statement.replace('|(96)', '|(64)')
            elif ea == 0x1D698: mutation, bad = 3, statement.replace('&(16)', '&(0)')
            elif ea == 0x1D775: mutation, bad = 4, statement.replace('(55384)', '(55386)')
            elif ea == 0x1D7FC: mutation, bad = 5, 'mc_string(m,0,0); /* lost REP */'
            elif ea == 0x1D81E: mutation, bad = 6, '; /* lost mask inversion */'
            elif ea == 0x1D81A: mutation, bad = 7, statement.replace('(16)', '(15)')
            elif ea == 0x1D7A3: mutation, bad = 8, statement.replace(',78,', ',79,')
            elif ea == 0x1D705: mutation, bad = 9, statement.replace(',8,', ',7,')
            elif ea == 0x1D735: mutation, bad = 10, statement.replace('(255)', '(254)')
            elif ea == 0x1D507: mutation, bad = 11, statement.replace(',4,', ',5,')
            elif ea == 0x1D4FE: mutation, bad = 12, statement.replace('&(16)', '&(0)')
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if bad is not None:
                assert bad != statement, (hex(ea), statement)
                lines += [f'#if KI_MAPCELLS_MUTATION == {mutation}', '    '+bad, '#else', '    '+statement, '#endif']
            else: lines.append('    '+statement)
        lines.append('}')
    lines.append('int ki_mapcells_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at in WANTED: lines.append(f'case 0x{at-0x10000:x}:mc_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    (repo / 'tools/c_recovery/mapcells_generated.inc').write_text('\n'.join(lines)+'\n')
    print('Generated',len(WANTED),'map-cell functions;',count,'original instructions')

if __name__ == '__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True)
    p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
