#!/usr/bin/env python3
"""原交戰與戰術完整靜態 C 圖；保留 guest stack、live patch 與原始入口。"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path

PROBE_SHA = '32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9'
INPUT_SHA = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
LIVE = {
    0x1B55F: 'ec_cmp(m,m->di,dl_read16(m,m->cs,0xb561),16);',
    0x1B567: 'ec_cmp(m,m->di,dl_read16(m,m->cs,0xb569),16);',
    0x1AB4E: 'vg_dl(m,dl_read8(m,m->cs,0xab4f));',
    0x1C601: 'ec_cmp(m,(uint8_t)m->cx,dl_read8(m,m->cs,0xc603),8);',
}
LIVE_JUMP = {0x1A06A, 0x1C604, 0x1BE33}
MUTANT_STATEMENTS = {84645: (1, '{uint8_t result=dl_read8(m,m->ds,m->si)&0xef;ec_logic(m,result,8);dl_write8(m,m->ds,m->si,result);}'), 84716: (2, '{uint8_t result=dl_read8(m,m->ds,m->si)&0xef;ec_logic(m,result,8);dl_write8(m,m->ds,m->si,result);}'), 84930: (3, 'ec_cmp(m,(uint8_t)m->bx,7,8);'), 106602: (4, 'if(m->flags&0x40)goto L_1A07C;'), 111967: (5, 'ec_cmp(m,m->di,0x600,16);'), 111975: (6, 'ec_cmp(m,m->di,0x600,16);'), 109390: (7, 'vg_dl(m,0x1f);'), 116225: (8, 'ec_cmp(m,(uint8_t)m->cx,0,8);'), 116228: (9, 'if(m->flags&0x40)goto L_1C60B;'), 114227: (10, 'if(m->flags&0x40)goto L_1BE38;'), 107607: (11, '{uint16_t target=dl_read16(m,m->cs,0xa466);ec_push(m,0xa45c);m->ip=target;eg_invoke_old=1;goto EG_DISPATCH;}'), 107927: (12, 'm->ip=dl_read16(m,m->cs,0xa59c);eg_invoke_old=1;goto EG_DISPATCH;'), 108513: (13, '{uint16_t target=dl_read16(m,m->cs,0xa7f7);ec_push(m,0xa7e6);m->ip=target;eg_invoke_old=1;goto EG_DISPATCH;}'), 108583: (14, '{uint16_t target=dl_read16(m,m->cs,0xa83b);ec_push(m,0xa82c);m->ip=target;eg_invoke_old=1;goto EG_DISPATCH;}'), 90076: (15, '{uint16_t target=dl_read16(m,m->ds,0x605c);ec_push(m,0x5fe0);m->ip=target;eg_invoke_old=1;goto EG_DISPATCH;}'), 109814: (16, '{uint16_t target=dl_read16(m,m->cs,0xad01);ec_push(m,0xacfb);m->ip=target;eg_invoke_old=1;goto EG_DISPATCH;}'), 121659: (17, 'm->bx=st_shr(m,m->bx,16);'), 121213: (18, 'm->cx=0x3bff;'), 90342: (19, 'm->cx=5;'), 109748: (20, 'vg_dl(m,3);')}


def generate(repo, probe_path, output=None):
    output = output or repo / "tools/c_recovery"
    raw_probe = probe_path.read_bytes()
    assert hashlib.sha256(raw_probe).hexdigest() == PROBE_SHA
    probe = json.loads(raw_probe)
    assert probe['input_sha256'] == probe['ida_input_sha256'] == INPUT_SHA
    assert not probe['unresolved_call_boundaries']
    spec = importlib.util.spec_from_file_location('dl', repo / 'tools/c_recovery_display_generate.py')
    dl = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(dl)
    selected = [t for t in probe['targets'] if t['ida_linear'] in set(probe['recovery_targets'])]
    rows = {}
    named_count = named_bytes = raw_count = raw_bytes = 0
    entries = {t['ida_linear'] for t in selected}
    originals = {}
    for target in selected:
        originals[target['ida_linear']] = target['name']
        for chunk in target['chunks']:
            assert hashlib.sha256(b''.join(bytes.fromhex(c['file_bytes']) for c in target['chunks'])).hexdigest() == target['file_sha256']
            for row in chunk['instructions']:
                if row['ida_linear'] in rows:
                    assert rows[row['ida_linear']]['bytes'] == row['bytes']
                rows[row['ida_linear']] = row
                named_count += 1
                named_bytes += len(bytes.fromhex(row['bytes']))
            for relocation in chunk['loader_relocations']:
                insn = next(i for i in chunk['instructions'] if i['file_offset'] <= relocation['file_offset'] < i['file_offset'] + len(bytes.fromhex(i['bytes'])))
                assert insn['bytes'].startswith('9a'), ('non-far relocation', insn)
    for block in probe['decoded_blocks']:
        entries.add(block['ida_linear'])
        originals[block['ida_linear']] = block.get('original_name_before_analysis') or f'raw_entry_{block["ida_linear"]:X}'
        for row in block['instructions']:
            if row['ida_linear'] in rows:
                assert rows[row['ida_linear']]['bytes'] == row['bytes']
            rows[row['ida_linear']] = row
            raw_count += 1
            raw_bytes += len(bytes.fromhex(row['bytes']))
    for table in probe['fixed_dispatch_tables']:
        if rows[table['call_or_jump_ida_linear']]['mnemonic'] == 'call':
            entries.update(0x10000 + word for word in table['words'] if 0x10000 + word in rows)
    for row in rows.values():
        if row['mnemonic'] == 'call' and row['operands'][0]['type'] == 7:
            target = row['operands'][0]['addr'] + 0x10000
            if target in rows:
                entries.add(target)
    assert set(LIVE) | LIVE_JUMP <= set(rows)
    lines = ['/* Original native C engagement graph; re/131, spec/251. No CPU.Step or guest decoder. */',
             'static int eg_known(uint16_t t) {switch(t) {']
    for ea in sorted(rows):
        lines.append(f'case 0x{ea-0x10000:x}:')
    lines.append('return 1;default:return 0;}}')
    lines.extend(['static void eg_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int eg_skip_first_enter) {',
                  'const uint16_t eg_core_cs=m->cs;int eg_invoke_old=0;m->ip=entry;',
                  'EG_DISPATCH:',
                  'if(m->cs!=eg_core_cs) {assert(!eg_invoke_old);return;}',
                  'switch(m->ip) {'])
    for ea in sorted(rows):
        lines.append(f'case 0x{ea-0x10000:x}:goto L_{ea:X};')
    lines.extend(['default:break;}',
                  'if(!eg_invoke_old)return;',
                  'ki_outcome_invoke(m,m->ip,h);eg_invoke_old=0;goto EG_DISPATCH;'])
    seen_mutants = set()
    for ea, row in sorted(rows.items()):
        mn, ops = row['mnemonic'], row['operands']
        next_ea = ea + len(bytes.fromhex(row['bytes']))
        nextoff = next_ea - 0x10000
        lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
        lines.append(f'm->ip=0x{ea-0x10000:x};')
        if ea in entries:
            lines.append(f'/* Original entry: {originals.get(ea, "raw callback entry")} */')
            lines.append(f'if(eg_skip_first_enter)eg_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x{ea-0x10000:x},h->user);')
        else:
            lines.append('eg_skip_first_enter=0;')
        def jump(target):
            if target in rows:
                return f'goto L_{target:X};'
            return f'm->ip=0x{target-0x10000:x};eg_invoke_old=1;goto EG_DISPATCH;'
        terminal = False
        if ea in LIVE:
            statement = LIVE[ea]
        elif ea in LIVE_JUMP:
            target = ops[0]['addr'] + 0x10000
            statement = '{uint8_t opcode=dl_read8(m,m->cs,0x%x);assert(opcode==0x74||opcode==0xeb);if(opcode==0xeb||(m->flags&0x40)) {%s}}' % (ea-0x10000, jump(target))
        elif mn == 'call':
            if row['bytes'].startswith('9a'):
                statement = f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);eg_invoke_old=0;goto EG_DISPATCH;'
            else:
                statement = '{uint16_t target=%s;ec_push(m,0x%x);m->ip=target;eg_invoke_old=1;goto EG_DISPATCH;}' % (dl.read(ops[0]), nextoff)
            terminal = True
        elif mn == 'retn':
            assert row['bytes'] == 'c3'
            statement = 'm->ip=ec_pop(m);eg_invoke_old=0;goto EG_DISPATCH;'
            terminal = True
        elif mn == 'jmp':
            statement = jump(ops[0]['addr'] + 0x10000) if ops[0]['type'] == 7 else f'm->ip={dl.read(ops[0])};eg_invoke_old=1;goto EG_DISPATCH;'
            terminal = True
        elif mn in dl.COND:
            statement = f'if({dl.COND[mn]}) {{{jump(ops[0]["addr"]+0x10000)}}}'
        elif mn == 'loop':
            statement = f'if(--m->cx) {{{jump(ops[0]["addr"]+0x10000)}}}'
        elif mn in ['movsw', 'lodsw', 'stosw']:
            repeat = int(bytes.fromhex(row['bytes'])[0] in [0xf2, 0xf3])
            statement = f'mc_string(m,{dict(movsw=0,lodsw=1,stosw=2)[mn]},{repeat});'
        elif mn in ['movsb', 'lodsb', 'stosb']:
            repeat = int(bytes.fromhex(row['bytes'])[0] in [0xf2, 0xf3])
            statement = f'eg_string8(m,{dict(movsb=0,lodsb=1,stosb=2)[mn]},{repeat});'
        elif mn == 'sar':
            assert ops[1]['value'] == 1
            statement = dl.write(ops[0], f'eg_sar1(m,{dl.read(ops[0])},{ops[0]["dtype_size"]*8})')
        elif mn == 'imul':
            statement = f'eg_imul(m,{dl.read(ops[-1])},{ops[-1]["dtype_size"]*8});'
        elif mn == 'cbw':
            statement = 'm->ax=(uint16_t)(int16_t)(int8_t)m->ax;'
        elif mn == 'les':
            statement = '{uint16_t off=%s,seg=%s;uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));%sm->es=bank;}' % (dl.address(ops[1]), dl.segment(ops[1]), dl.write(ops[0], 'value'))
        elif mn == 'int':
            statement = f'gl_interrupt(m,{dl.read(ops[0])},0x{nextoff:x});'
        elif mn == 'pushf':
            statement = 'ec_push(m,m->flags);'
        elif mn == 'popf':
            statement = 'm->flags=ec_pop(m)|2;'
        elif mn == 'lahf':
            statement = 'vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
        elif mn == 'sahf':
            statement = 'm->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
        elif mn == 'cli':
            statement = 'm->flags&=(uint16_t)~0x200u;'
        elif mn == 'sti':
            statement = 'm->flags|=0x200;'
        else:
            try:
                statement = dl.statement(row)
            except Exception as error:
                raise ValueError(f'{ea:X}: {row["assembly"]}: {error}') from error
        if ea in MUTANT_STATEMENTS:
            number,bad = MUTANT_STATEMENTS[ea]
            assert bad != statement
            seen_mutants.add(ea)
            lines.extend([f"#if KI_ENGAGEMENT_MUTATION == {number}",bad,"#else",statement,"#endif"])
        else:
            lines.append(statement)
        if not terminal:
            assert next_ea in rows, ('missing fallthrough instruction', hex(ea), hex(next_ea))
            lines.append(f'goto L_{next_ea:X};')
    assert seen_mutants == set(MUTANT_STATEMENTS)
    lines.append('}')
    lines.extend(['int ki_engagement_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {',
                  'if(!eg_known(t))return 0;',
                  'uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);',
                  'eg_graph(m,t,h,1); /* Old invoke already emitted the entry hook. */',
                  '/* Old external ABI owns one final POP; reject abandoned native frames. */',
                  'assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);',
                  'ec_push(m,ret);return 1;}',
                  'void ki_engagement_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(eg_known(t))eg_graph(m,t,h,0);else ki_outcome_invoke(m,t,h);}'])
    declarations = []
    for target in selected:
        name, offset = target['name'], target['ida_linear'] - 0x10000
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{eg_graph(m,0x{offset:x},h,0);}}')
        declarations.append(f'void {name}(KiMachine16 *,const KiEconomyHooks *);')
    for block in probe['decoded_blocks']:
        ea = block['ida_linear']
        name = f'raw_entry_{ea:X}'
        lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{eg_graph(m,0x{ea-0x10000:x},h,0);}}')
        declarations.append(f'void {name}(KiMachine16 *,const KiEconomyHooks *);')
    (output / 'engagement_generated.inc').write_text('\n'.join(lines)+'\n', encoding='utf-8')
    (output / 'engagement.h').write_text('\n'.join(['#ifndef WOLONG_C_ENGAGEMENT_H', '#define WOLONG_C_ENGAGEMENT_H', '#include "outcome.h"',
        'int ki_engagement_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);', 'void ki_engagement_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);', *declarations, '#endif'])+'\n', encoding='utf-8')
    summary = {'scope':'spec251 original static C instruction graph; semantic verification pending',
        'probe_sha256':PROBE_SHA, 'named_functions':len(selected), 'named_instructions':named_count,
        'named_instruction_bytes':named_bytes, 'raw_blocks':len(probe['decoded_blocks']),
        'raw_instructions':raw_count, 'raw_instruction_bytes':raw_bytes, 'unique_instructions':len(rows),
        'hook_entries':[f'0x{x:X}' for x in sorted(entries)],
        'live_instruction_consumers':[f'0x{x:X}' for x in sorted(set(LIVE)|LIVE_JUMP)],
        'instruction_decoder':False, 'c_side_cpu_step':False}
    print(f'Generated {len(selected)} named / {len(probe["decoded_blocks"])} raw / {len(rows)} instructions; 20 controls, 7 live consumers')
    return summary


if __name__ == '__main__':
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repo',required=True,type=Path)
    parser.add_argument('--probe',required=True,type=Path)
    parser.add_argument('--output',type=Path)
    args=parser.parse_args()
    generate(args.repo,args.probe,args.output)
