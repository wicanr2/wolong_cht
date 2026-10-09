#!/usr/bin/env python3
"""固定 IDA 八入口轉 C，保留 live DS 分派表與原 RET 落下邊界。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
TARGETS = {
    0x109D0: ('sub_109D0', 76), 0x11E46: ('sub_11E46', 200),
    0x11F0E: ('sub_11F0E', 34), 0x159B7: ('sub_159B7', 25),
    0x159D0: ('nullsub_1', 1), 0x15AA2: ('sub_15AA2', 20),
    0x15E4C: ('sub_15E4C', 20), 0x161B6: ('sub_161B6', 20),
}
MUTATIONS = {
    1: ('fade', 'fade-final-level'), 2: ('fade', 'second-vblank-wait'),
    3: ('popup', 'popup-y-clamp'), 4: ('popup', 'popup-x-clamp'),
    5: ('point', 'city-tile-lower-bound'), 6: ('point', 'city-tile-upper-bound'),
    7: ('point', 'popup-carry-cancel'), 8: ('right', 'right-index-reset'),
    9: ('right', 'live-table-read'), 10: ('clear', 'clear-bit-two'),
    11: ('clear', 'clear-bit-one'), 12: ('clear', 'clear-bit-zero'),
}
MUTANT_STATEMENTS = {
    0x10A11: (1, 'ec_cmp(m,(uint8_t)m->cx,15,8);'),
    0x109E9: (2, 'goto L_109F3; /* Omit the second retrace wait. */'),
    0x11F12: (3, 'ec_cmp(m,(uint8_t)m->ax,21,8);'),
    0x11F16: (3, 'vg_al(m,21);'),
    0x11F1D: (4, 'ec_cmp(m,(uint8_t)m->dx,34,8);'),
    0x11F22: (4, 'vg_dl(m,34);'),
    0x11E7F: (5, 'ec_cmp(m,(uint8_t)m->cx,0xcc,8);'),
    0x11E84: (6, 'ec_cmp(m,(uint8_t)m->cx,0xd3,8);'),
    0x11EBF: (7, '; /* Ignore the city/corps popup cancellation carry. */'),
    0x159BB: (8, 'm->bx=0; /* Wrong but valid live table[0] index. */'),
    0x159C0: (9, 'ix_call(m,0x59d0,0x59c4,h); /* Hard-coded nullsub_1. */'),
    0x15AB0: (10, '{uint16_t result=dl_read8(m,m->ds,0x98a6)&0xfd;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}'),
    0x15E4C: (11, '{uint16_t result=dl_read8(m,m->ds,0x98a6)&0xfe;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}'),
    0x161B6: (12, '{uint16_t result=dl_read8(m,m->ds,0x98a6)&0xfb;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}'),
}


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    data = json.loads(probe.read_text(encoding='utf-8'))
    assert data['input_sha256'] == EXPECTED
    assert set(data['recovery_targets']) == set(TARGETS)
    groups = []
    for target in sorted(data['targets'], key=lambda row: row['ida_linear']):
        at = target['ida_linear']
        if at not in TARGETS:
            continue
        original = b''.join(bytes.fromhex(chunk['file_bytes']) for chunk in target['chunks'])
        assert (target['name'], len(original)) == TARGETS[at]
        assert hashlib.sha256(original).hexdigest() == target['file_sha256']
        assert not any(chunk['loader_relocations'] for chunk in target['chunks'])
        rows = [row for chunk in target['chunks'] for row in chunk['instructions']]
        groups.append((at, rows, target['name']))
    assert len(groups) == 8 and sum(len(rows) for _, rows, _ in groups) == 159
    lines = ['/* Original world interaction; re/127, spec/247. */']
    lines.extend(f'static void ix_body_{at:X}(KiMachine16 *,const KiEconomyHooks *);' for at, _, _ in groups)
    seen = set()
    for at, rows, name in groups:
        labels = {row['ida_linear'] for row in rows}
        lines.extend([f'/* Original {name}; IDA 0x{at:X}. */',
                      f'static void ix_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
        for row in rows:
            ea, mnemonic = row['ida_linear'], row['mnemonic']
            nextoff = (ea + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mnemonic == 'call' and row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
            elif mnemonic == 'call':
                statement = f'ix_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            elif mnemonic == 'in':
                assert row['operands'][0]['dtype_size'] == 1
                statement = dl.write(row['operands'][0], f'wolong_main_in8({dl.read(row["operands"][1])})')
            else:
                if mnemonic in list(dl.COND) + ['jmp', 'loop']:
                    assert row['operands'][0]['addr'] + 0x10000 in labels
                statement = dl.statement(row)
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if ea in MUTANT_STATEMENTS:
                number, bad = MUTANT_STATEMENTS[ea]
                assert bad != statement
                seen.add(ea)
                lines.extend([f'#if KI_INTERACTION_MUTATION == {number}', '    ' + bad,
                              '#else', '    ' + statement, '#endif'])
            else:
                lines.append('    ' + statement)
        if at == 0x159B7:
            assert rows[-1]['ida_linear'] == 0x159CE and rows[-1]['mnemonic'] != 'retn'
            lines.extend(['m->ip=0x59d0; /* Original fall-through to the separate nullsub_1 RET. */',
                          'if(h&&h->enter)h->enter(m,0x59d0,h->user);',
                          'ix_body_159D0(m,h); /* Body tail only; the caller consumes one guest return address. */'])
        else:
            assert rows[-1]['mnemonic'] == 'retn'
        lines.append('}')
    assert seen == set(MUTANT_STATEMENTS)
    lines.append('int ki_interaction_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at, _, _ in groups:
        lines.append(f'case 0x{at-0x10000:x}:ix_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    lines.append('static int ix_known(uint16_t t) {switch(t) {' +
                 ' '.join(f'case 0x{at-0x10000:x}:' for at, _, _ in groups) + ' return 1;default:return 0;}}')
    for at, _, name in groups:
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_interaction_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
    (repo / 'tools/c_recovery/interaction_generated.inc').write_text('\n'.join(lines) + '\n', encoding='utf-8')
    print('Generated 8 original entries / 159 instructions / 396 bytes; 12 mutations at 14 sites')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
