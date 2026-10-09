#!/usr/bin/env python3
"""原軍團／物件 producer 與 live capacity 矩陣 handler 的 native C。"""
import argparse
import importlib.util
import json
from pathlib import Path

def generate(repo,probe):
    spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py')
    dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
    r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
    groups += [(b['ida_linear'],b['instructions'],b['original_name_before_analysis']) for b in r['decoded_blocks']]
    lines=['/* Native C original overlay producers; re/115, spec/235. Original locators retained. */']
    for at,rows,name in groups:
        labels={x['ida_linear'] for x in rows}
        lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void ov_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
        for x in rows:
            ea=x['ida_linear']
            if x['mnemonic'] in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
            if ea==0x1D5A9:statement='ec_cmp(m,(uint8_t)m->bx,dl_read8(m,m->cs,0xd5ab),8);'
            else:statement=dl.statement(x).replace('dl_call(', 'ov_call(')
            mutation,bad=0,None
            if ea==0x12B18:mutation,bad=1,statement.replace(',192,',',191,')
            elif ea==0x12B2D:mutation,bad=2,statement.replace('m->si+18','m->si+19')
            elif ea==0x12B35:mutation,bad=3,statement.replace('m->si+8','m->si+7')
            elif ea==0x12B4C:mutation,bad=4,'; /* lost formation pattern contribution */'
            elif ea==0x12560:mutation,bad=5,statement.replace('&(7)','&(3)')
            elif ea==0x1256C:mutation,bad=6,'; /* lost metadata stride shift */'
            elif ea==0x1257B:mutation,bad=7,statement.replace(',256,',',128,')
            elif ea==0x1D5A9:mutation,bad=8,'ec_cmp(m,(uint8_t)m->bx,4,8);'
            elif ea==0x1D59A:mutation,bad=9,statement.replace(',255,',',254,')
            elif ea==0x1D5A0:mutation,bad=10,statement.replace('&(16)','&(0)')
            elif ea==0x1D5BB:mutation,bad=11,statement.replace('|(32)','|(0)')
            elif ea==0x15D47:mutation,bad=12,statement.replace('(96)','(64)')
            lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
            if bad is not None:
                assert bad!=statement,(hex(ea),statement)
                lines += [f'#if KI_OVERLAY_MUTATION == {mutation}','    '+bad,'#else','    '+statement,'#endif']
            else:lines.append('    '+statement)
        lines.append('}')
    lines.append('int ki_overlay_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:ov_body_{at:X}(m,h);break;')
    lines.append('default:return 0;}return 1;}')
    (repo/'tools/c_recovery/overlay_generated.inc').write_text('\n'.join(lines)+'\n')
    print('Generated',len(r['targets']),'named and',len(r['decoded_blocks']),'raw entries;',sum(len(x[1]) for x in groups),'instructions')

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True)
    a=p.parse_args();generate(a.repo,a.probe)
