#!/usr/bin/env python3
"""原場景主入口／鏡頭／小地圖 C，保留跨函式尾跳與 far ABI。"""
import argparse
import importlib.util
import json
from pathlib import Path

def generate(repo,probe):
 s=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(s);s.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets'] if t['ida_linear']!=0x11d46]
 lines=['/* Original native C scene/camera/minimap; re/116, spec/236. */']
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows}
  lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void sr_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if ea==0x11f3b:statement='if(m->flags&1) {sr_tail(m,0x1f5a,h);return;}'
   elif mn=='call' and bytes.fromhex(x['bytes'])[0]==0x9a:statement=f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
   elif mn in ['movsw','lodsw']:statement=f'mc_string(m,{0 if mn=="movsw" else 1},0);'
   elif mn=='rcr':statement='m->dx=hr_rcr(m,m->dx);'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x).replace('dl_call(', 'sr_call(')
   mutation,bad=0,None
   if ea==0x11f38:mutation,bad=1,statement.replace(',2,',',1,')
   elif ea==0x11f71:mutation,bad=2,statement.replace(',39050,',',39052,')
   elif ea==0x1d4c2:mutation,bad=3,statement.replace('|(224)','|(96)')
   elif ea==0x15c8a:mutation,bad=4,'; /* lost remembered minimap X */'
   elif ea==0x19555:mutation,bad=5,statement.replace('(607)','(608)')
   elif ea==0x195bf:mutation,bad=6,statement.replace(',20,',',19,')
   elif ea==0x195c4:mutation,bad=7,statement.replace(',2808,',',2807,')
   elif ea==0x1976b:mutation,bad=8,'; /* lost rotation through carry */'
   elif ea==0x19758:mutation,bad=9,statement.replace('(11)','(10)')
   elif ea==0x13b2c:mutation,bad=10,statement.replace(',3,',',2,')
   elif ea==0x13b40:mutation,bad=11,'; /* lost original scene restore */'
   elif ea==0x13b37:mutation,bad=12,statement.replace(',3,',',0,')
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if bad is not None:
    assert bad!=statement,(hex(ea),statement);lines += [f'#if KI_RESUME_MUTATION == {mutation}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  lines.append('}')
 lines.append('int ki_resume_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:sr_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 (repo/'tools/c_recovery/resume_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated',len(groups),'new routines;',sum(len(x[1]) for x in groups),'original instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
