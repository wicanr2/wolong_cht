#!/usr/bin/env python3
"""固定 IDA 自動戰鬥及據點易主閉包轉 C，保留字寬、FLAGS 與非區域退出。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PROBE_SHA256 = 'b20f688289b06ca5384f3ac810cd4b9b93cc5ddb82d5be4d12556701a8988ce9'
TARGETS = {
    0x10CE7: ('sub_10CE7', 9), 0x14236: ('sub_14236', 51),
    0x1474A: ('sub_1474A', 113), 0x14CF3: ('sub_14CF3', 112),
    0x14D63: ('sub_14D63', 65), 0x14DA4: ('sub_14DA4', 76),
    0x14DF0: ('sub_14DF0', 108), 0x14FCE: ('sub_14FCE', 166),
    0x15074: ('sub_15074', 64), 0x150B4: ('sub_150B4', 35),
    0x15130: ('sub_15130', 131), 0x151B3: ('sub_151B3', 210),
    0x15285: ('sub_15285', 82), 0x152D7: ('sub_152D7', 102),
    0x15CA4: ('sub_15CA4', 34), 0x15CE0: ('sub_15CE0', 57),
    0x17028: ('sub_17028', 20), 0x188CC: ('sub_188CC', 62),
    0x1890A: ('sub_1890A', 83), 0x195C9: ('sub_195C9', 118),
    0x1963F: ('sub_1963F', 23),
}
MUTANT_STATEMENTS = {
    0x1513F: (1, 'vg_al(m,2);'),
    0x152CA: (2, '; /* Omitted third morale SHR, giving division by four. */'),
    0x1530E: (3, 'vg_cl(m,3);'),
    0x152F6: (4, '{uint8_t result=(uint8_t)m->ax&1;ec_logic(m,result,8);vg_al(m,result);}'),
    0x15155: (5, 'm->cx=ec_math(m,m->cx,7,16,0,0);'),
    0x15158: (5, 'm->dx=ec_math(m,m->dx,7,16,0,0);'),
    0x15173: (6, 'ec_cmp(m,m->ax,99,16);'),
    0x15178: (6, 'm->ax=99;'),
    0x15223: (7, 'vg_al(m,ec_math(m,(uint8_t)m->ax,7,8,0,0));'),
    0x1520E: (8, 'vg_ah(m,0);'),
    0x15233: (8, 'vg_ah(m,0);'),
    0x1478F: (9, 'ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),299,16);'),
    0x14775: (10, '; /* Omitted original own-city stay branch. */'),
    0x14D0A: (11, '; /* Omitted old city-count DEC. */'),
    0x14D2A: (12, '; /* Omitted new city-count INC. */'),
    0x14D7E: (13, 'dl_write8(m,m->ds,(uint16_t)(m->bx+0x17),1);'),
    0x14E54: (14, '{uint8_t result=dl_read8(m,m->ds,m->bx)&0xbf;ec_logic(m,result,8);dl_write8(m,m->ds,m->bx,result);}'),
    0x14D18: (15, '; /* Omitted bounded old-corps redirection call. */'),
    0x1892B: (16, '; /* Omitted neighboring city-count INC. */'),
    0x14FE8: (17, '; /* Omitted surviving-faction-count DEC. */'),
    0x15091: (18, 'dl_write8(m,m->ds,(uint16_t)(m->bx+0x17),1);'),
    0x1425B: (19, '; /* Omitted city owner-history write. */'),
    0x17036: (20, '; /* Omitted army occupancy DEC. */'),
}


def generate(repo, probe):
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    assert hashlib.sha256(probe.read_bytes()).hexdigest() == PROBE_SHA256
    data = json.loads(probe.read_text(encoding='utf-8'))
    assert data['input_sha256'] == EXPECTED and set(data['recovery_targets']) == set(TARGETS)
    assert not data.get('decoded_blocks')
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
    assert len(groups) == 21 and sum(len(rows) for _, rows, _ in groups) == 796
    assert sum(size for _, size in TARGETS.values()) == 1721
    lines = ['/* Original native C battle/city outcome closure; re/130, spec/250. Nonlocal 11CB1 uses the existing transfer runner. */']
    seen = set()
    for at, rows, name in groups:
        labels = {row['ida_linear'] for row in rows}
        lines.extend([f'/* Original {name}; IDA 0x{at:X}. */',
                      f'static void oc_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
        for row in rows:
            ea, mn = row['ida_linear'], row['mnemonic']
            nextoff = (ea + len(bytes.fromhex(row['bytes'])) - 0x10000) & 65535
            if mn == 'call' and row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
            elif mn == 'call':
                statement = f'oc_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
            elif mn == 'pushf':
                statement = 'ec_push(m,m->flags);'
            elif mn == 'popf':
                statement = 'm->flags=ec_pop(m)|2;'
            elif mn == 'rcl':
                assert row['operands'][1]['value'] == 1
                statement = dl.write(row['operands'][0], 'tx_rcl(m,' + dl.read(row['operands'][0]) + ',' + str(row['operands'][0]['dtype_size'] * 8) + ')')
            elif mn == 'rcr':
                assert row['operands'][0]['dtype_size'] == 2 and row['operands'][1]['value'] == 1
                statement = dl.write(row['operands'][0], 'hr_rcr(m,' + dl.read(row['operands'][0]) + ')')
            elif mn == 'xlat':
                assert ea == 0x1530D and row['bytes'] == 'd7'
                statement = 'vg_al(m,dl_read8(m,m->ds,(uint16_t)(m->bx+(uint8_t)m->ax)));'
            elif mn == 'movsw':
                statement = f'mc_string(m,0,{1 if row["bytes"].startswith("f3") else 0});'
            else:
                if mn in list(dl.COND) + ['jmp', 'loop']:
                    assert row['operands'][0]['addr'] + 0x10000 in labels
                statement = dl.statement(row)
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
            if ea in MUTANT_STATEMENTS:
                number, bad = MUTANT_STATEMENTS[ea]
                assert bad != statement
                seen.add(ea)
                lines.extend([f'#if KI_OUTCOME_MUTATION == {number}', '    ' + bad, '#else', '    ' + statement, '#endif'])
            else:
                lines.append('    ' + statement)
        assert rows[-1]['mnemonic'] == 'retn'
        lines.append('}')
    assert seen == set(MUTANT_STATEMENTS) and {n for n, _ in MUTANT_STATEMENTS.values()} == set(range(1, 21))
    lines.append('int ki_outcome_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at, _, _ in groups:
        lines.append(f'case 0x{at-0x10000:x}:oc_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    lines.append('static int oc_known(uint16_t t) {switch(t) {' + ' '.join(f'case 0x{at-0x10000:x}:' for at, _, _ in groups) + ' return 1;default:return 0;}}')
    for at, _, name in groups:
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_outcome_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
    (repo / 'tools/c_recovery/outcome_generated.inc').write_text('\n'.join(lines) + '\n', encoding='utf-8')
    print('Generated 21 original functions / 796 instructions / 1721 bytes; 20 mutations at 23 sites')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo', type=Path, required=True)
    parser.add_argument('--probe', type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.probe)
