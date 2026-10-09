#!/usr/bin/env python3
"""固定 IDA 尋路與潰散閉包轉 C，保留 live code patch、佇列及 raw RET。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PROBE_SHA256 = 'ff99cfbd72b7211d14ac5d7544ad7ffa5af9fa8ac8c421d2f1d4311a438d1a66'
TARGETS = {
    0x1291A: ('sub_1291A', 93), 0x12977: ('sub_12977', 76),
    0x129C3: ('sub_129C3', 187), 0x12BA8: ('sub_12BA8', 49),
    0x147BB: ('sub_147BB', 192), 0x1487B: ('sub_1487B', 160),
    0x14A0F: ('sub_14A0F', 108), 0x19656: ('sub_19656', 121),
    0x196CF: ('sub_196CF', 30),
}
LIVE_STATEMENTS = {
    0x149B6: 'ec_cmp(m,m->si,dl_read16(m,m->cs,0x49b8),16);',
    0x149BC: 'ec_cmp(m,m->si,dl_read16(m,m->cs,0x49be),16);',
    0x149CD: 'ec_cmp(m,dl_read8(m,m->es,(uint16_t)(m->si+0x841)),dl_read8(m,m->cs,0x49d2),8);',
}
MUTANT_STATEMENTS = {
    0x149B6: (1, 'ec_cmp(m,m->si,dl_read16(m,m->cs,0x49be),16);'),
    0x149BC: (2, 'ec_cmp(m,m->si,dl_read16(m,m->cs,0x49b8),16);'),
    0x149CD: (3, 'ec_cmp(m,dl_read8(m,m->es,(uint16_t)(m->si+0x841)),0x12,8);'),
    0x149DC: (4, 'm->dx=ec_math(m,m->dx,5,16,0,0);'),
    0x149D5: (5, 'm->dx=ec_math(m,m->dx,0xa5,16,0,0);'),
    0x149D9: (6, '{uint8_t result=(uint8_t)(m->dx>>8)|0x40;ec_logic(m,result,8);vg_dh(m,result);}'),
    0x14A3D: (7, 'vg_dl(m,ec_math(m,(uint8_t)m->dx,dl_read8(m,m->ds,(uint16_t)(m->bx+5)),8,0,0));'),
    0x14A0F: (8, 'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x8000)),1,8);'),
    0x14985: (9, 'if(m->flags&(1|0x40)) goto L_149B6;'),
    0x14A3A: (10, 'm->bp=0xfffc;'),
    0x147F8: (11, 'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x23)),11,8);'),
    0x14848: (11, 'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x23)),11,8);'),
    0x14816: (12, '; /* Omitted original leg bit0 OR. */'),
    0x14823: (12, '; /* Omitted original leg bit0 OR. */'),
    0x1482A: (12, '; /* Omitted original leg bit0 OR. */'),
    0x12965: (13, 'vg_ah(m,ec_math(m,(uint8_t)(m->ax>>8),0x27,8,0,0));'),
    0x1298C: (14, '; /* Omitted original occupancy DEC. */'),
    0x129ED: (15, 'dl_write8(m,m->ds,m->si,8);'),
    0x12BA8: (16, '{uint8_t result=dl_read8(m,m->ds,m->si)&0xdf;ec_logic(m,result,8);dl_write8(m,m->ds,m->si,result);}'),
}


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    assert hashlib.sha256(probe.read_bytes()).hexdigest() == PROBE_SHA256
    data = json.loads(probe.read_text(encoding='utf-8'))
    assert data['input_sha256'] == EXPECTED and set(data['recovery_targets']) == set(TARGETS)
    groups = []
    for target in sorted(data['targets'], key=lambda row: row['ida_linear']):
        at = target['ida_linear']
        if at not in TARGETS:
            continue
        raw = b''.join(bytes.fromhex(chunk['file_bytes']) for chunk in target['chunks'])
        assert (target['name'], len(raw)) == TARGETS[at]
        assert hashlib.sha256(raw).hexdigest() == target['file_sha256']
        assert not any(chunk['loader_relocations'] for chunk in target['chunks'])
        groups.append((at, [row for chunk in target['chunks'] for row in chunk['instructions']], target['name']))
    assert len(groups) == 9 and sum(len(rows) for _, rows, _ in groups) == 436
    assert sum(size for _, size in TARGETS.values()) == 1016
    assert len(data['decoded_blocks']) == 1
    block = data['decoded_blocks'][0]
    assert (block['ida_linear'], block['end_ida_linear'], block['original_name_before_analysis']) == (0x1491B, 0x14A0F, 'loc_1491B')
    raw = b''.join(bytes.fromhex(row['bytes']) for row in block['instructions'])
    assert len(block['instructions']) == 107 and len(raw) == 244
    assert hashlib.sha256(raw).hexdigest() == block['file_sha256'] == 'beaa3cc04dbc47e5de38fff064e9299e240c1f31af0b54cab795a0ccc1001133'
    groups.append((0x1491B, block['instructions'], 'loc_1491B raw entry, complete through original RET'))
    groups.sort()
    lines = ['/* Original native C route/collapse closure; re/129, spec/249. Raw loc_1491B retains its original name. */']
    seen = set()
    for at, rows, name in groups:
        labels = {row['ida_linear'] for row in rows}
        lines.extend([f'/* Original {name}; IDA 0x{at:X}. */',
                      f'static void rt_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
        for row in rows:
            ea, mn = row['ida_linear'], row['mnemonic']
            nextoff = (ea + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if ea in LIVE_STATEMENTS:
                statement = LIVE_STATEMENTS[ea]
            elif mn == 'call' and row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
            elif mn == 'call':
                statement = f'rt_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            elif mn == 'stosw':
                statement = f'mc_string(m,2,{1 if row["bytes"].startswith("f3") else 0});'
            else:
                if mn in list(dl.COND) + ['jmp', 'loop']:
                    assert row['operands'][0]['addr'] + 0x10000 in labels
                statement = dl.statement(row)
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if ea in MUTANT_STATEMENTS:
                number, bad = MUTANT_STATEMENTS[ea]
                assert bad != statement
                seen.add(ea)
                lines.extend([f'#if KI_ROUTE_MUTATION == {number}', '    ' + bad, '#else', '    ' + statement, '#endif'])
            else:
                lines.append('    ' + statement)
        assert rows[-1]['mnemonic'] == 'retn'
        lines.append('}')
    assert seen == set(MUTANT_STATEMENTS) and {n for n, _ in MUTANT_STATEMENTS.values()} == set(range(1, 17))
    lines.append('int ki_route_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at, _, _ in groups:
        lines.append(f'case 0x{at-0x10000:x}:rt_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    lines.append('static int rt_known(uint16_t t) {switch(t) {' + ' '.join(f'case 0x{at-0x10000:x}:' for at, _, _ in groups) + ' return 1;default:return 0;}}')
    for at, _, name in groups:
        symbol = 'raw_entry_1491B' if at == 0x1491B else name
        lines.append(f'void {symbol}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_route_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
    (repo / 'tools/c_recovery/route_generated.inc').write_text('\n'.join(lines) + '\n', encoding='utf-8')
    print('Generated 9 named functions / 436 instructions / 1016 bytes + raw 107 / 244; 16 mutations at 19 sites')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
