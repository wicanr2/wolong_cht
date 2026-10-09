#!/usr/bin/env python3
"""由固定 IDA 證據還原君主出陣／自動編成；不執行 guest opcode。"""
import argparse
import importlib.util
import json
from pathlib import Path

TARGETS=(0x1461D,0x14698,0x14717,0x15E80,0x15EB7,0x15F27,0x15F5D,
         0x15F7F,0x1699E,0x16E8F,0x16EC9,0x16F26,0x16F86,0x16FD2)

def generate(repo,probe):
 s=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(s);s.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets'] if t['ida_linear'] in TARGETS]
 assert {a for a,b,c in groups}==set(TARGETS)
 lines=['/* Original native C ruler sortie / formation / sidebar; re/117, spec/237. */']
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows}
  lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void vd_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if mn=='call' and bytes.fromhex(x['bytes'])[0]==0x9a:statement=f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
   elif mn=='cbw':statement='m->ax=(uint16_t)(int16_t)(int8_t)m->ax;'
   elif mn=='lahf':statement='vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
   elif mn=='sahf':statement='m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x).replace('dl_call(', 'vd_call(')
   mutation,bad=0,None
   if ea==0x169e7:mutation,bad=1,statement.replace(',600,',',601,')
   elif ea==0x169c7:mutation,bad=2,'if(m->flags&1) goto L_169D2;'
   elif ea==0x16efb:mutation,bad=3,statement.replace(',50,',',51,')
   elif ea==0x16f1e:mutation,bad=4,'; /* lost carry restoration */'
   elif ea==0x146f0:mutation,bad=5,'; /* lost division remainder */'
   elif ea==0x146f7:mutation,bad=6,statement.replace('(100)','(99)')
   elif ea==0x16fca:mutation,bad=7,'; /* lost occupancy increment */'
   elif ea==0x16ffb:mutation,bad=8,statement.replace('m->bx','m->ax')
   elif ea==0x16a08:mutation,bad=9,statement.replace('&(251)','&(255)')
   elif ea==0x15e80:mutation,bad=10,statement.replace('&(2)','&(0)')
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if bad is not None:
    assert bad!=statement,(hex(ea),statement);lines += [f'#if KI_VERDICT_MUTATION == {mutation}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  lines.append('}')
 lines.append('int ki_verdict_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:vd_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 (repo/'tools/c_recovery/verdict_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated',len(groups),'routines;',sum(len(x[1]) for x in groups),'original instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
