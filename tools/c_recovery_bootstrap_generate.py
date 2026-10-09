#!/usr/bin/env python3
"""固定 IDA 初始化閉包與無 RET 的有限前綴轉原生 C。"""
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
    assert wanted == {0x1533D, 0x189F0, 0x18A1E, 0x18AD1, 0x18AEA}
    groups = [(t['ida_linear'], [x for c in t['chunks'] for x in c['instructions']], t['name'])
              for t in data['targets'] if t['ida_linear'] in wanted]
    fragment = data['decoded_blocks'][0]
    assert fragment['ida_linear'] == 0x11BE0 and fragment['end_ida_linear'] == 0x11BF6
    groups.append((0x11BE0, fragment['instructions'], 'sub_11BE0 bounded prefix, not complete function'))
    mutants = {
        0x18A46: (1, 'vg_cl(m,2);'),
        0x18A56: (2, 'vg_al(m,ec_math(m,(uint8_t)m->ax,1,8,0,0));'),
        0x18A72: (3, 'vg_ah(m,0);'),
        0x18AD8: (4, 'ec_cmp(m,(uint8_t)m->ax,21,8);'),
        0x18AEA: (5, 'ec_cmp(m,dl_read8(m,m->ds,m->si),0x81,8);'),
        0x18B08: (6, 'dl_write16(m,m->ds,(uint16_t)(m->si+28),0xb000);'),
        0x18B0E: (7, 'dl_write8(m,m->es,m->bx,1);'),
        0x18A09: (8, 'm->cx=128;'),
        0x15357: (9, 'm->ds=entry_ds;return;'),
        0x11BF1: (10, 'dl_write16(m,m->cs,0x9901,(uint16_t)(m->sp+2));'),
    }
    lines = ['/* Original native C world bootstrap; re/126, spec/246. Full 11BE0 is not registered. */']
    seen = set()
    for at, rows, name in groups:
        labels = {x['ida_linear'] for x in rows}
        lines.extend([f'/* Original {name}; IDA 0x{at:X}. */', f'static void bs_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
        if at == 0x1533D:
            lines.extend(['#if KI_BOOTSTRAP_MUTATION == 9', 'const uint16_t entry_ds=m->ds;', '#endif'])
        for row in rows:
            ea, mn = row['ida_linear'], row['mnemonic']
            nextoff = (ea + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mn == 'call':
                assert not row['bytes'].startswith('9a')
                statement = f'bs_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            else:
                if mn in list(dl.COND) + ['jmp', 'loop']:
                    assert row['operands'][0]['addr'] + 0x10000 in labels
                statement = dl.statement(row)
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if ea in mutants:
                number, bad = mutants[ea]
                assert bad != statement
                seen.add(ea)
                lines.extend([f'#if KI_BOOTSTRAP_MUTATION == {number}', '    ' + bad, '#else', '    ' + statement, '#endif'])
            else:
                lines.append('    ' + statement)
        if at == 0x11BE0:
            assert rows[-1]['ida_linear'] == 0x11BF1
            lines.append('m->ip=0x1bf6; /* Explicit control point; no virtual RET or POP. */')
        else:
            assert rows[-1]['mnemonic'] == 'retn'
        lines.append('}')
    assert seen == set(mutants)
    lines.append('int ki_bootstrap_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at, rows, name in groups:
        if at != 0x11BE0:
            lines.append(f'case 0x{at-0x10000:x}:bs_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    lines.append('static int bs_known(uint16_t t) {switch(t) {' + ' '.join(f'case 0x{at-0x10000:x}:' for at, rows, name in groups if at != 0x11BE0) + ' return 1;default:return 0;}}')
    for at, rows, name in groups:
        if at != 0x11BE0:
            lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_bootstrap_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
    lines.append('void raw_entry_11BE0_prefix(KiMachine16 *m,const KiEconomyHooks *h) {bs_body_11BE0(m,h);}')
    (repo / 'tools/c_recovery/bootstrap_generated.inc').write_text('\n'.join(lines) + '\n')
    print('Generated 5 functions / 140 instructions + prefix 7 instructions; total 339 bytes')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
