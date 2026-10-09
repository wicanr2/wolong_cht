#!/usr/bin/env python3
"""固定 IDA 資訊十函式轉 C，保留原始定位；不複製共用規則。"""
import argparse,importlib.util,json
from pathlib import Path
def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
 assert len(groups)==10 and not r['decoded_blocks']
 lines=['/* Original native C details; re/122, spec/242. */']
 mutants={
  0x17E64:(1,'; /* lost neutral name branch */'),
  0x17E94:(2,'vg_al(m,9);'),
  0x17EC1:(3,'goto L_17EC5; /* lost negative color */'),
  0x17EFE:(4,'m->ax=24;'),
  0x17F31:(5,'m->ax=(m->ax+1)%15;dl_mul(m,m->dx,16);'),
  0x1810D:(6,'m->bx=0x0904;'),
  0x18116:(7,'m->ax=dl_read16(m,m->ds,(uint16_t)(m->si+6));'),
  0x1811B:(7,'; /* retained word morale */'),
  0x18150:(8,'m->si=ec_math(m,m->si,0x15c0,16,0,0);'),
  0x18156:(9,'m->bx=ec_math(m,m->bx,4,16,0,1);'),
  0x18191:(10,'goto L_181A2; /* lost own-faction restoration */'),
  0x15E2D:(11,'{uint16_t result=dl_read8(m,m->ds,0x98a6)|4;ec_logic(m,result,8);dl_write8(m,m->ds,0x98a6,result);}'),
  0x17E2C:(12,'m->si=ec_math(m,m->si,0x860,16,0,0);'),
 }
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows};lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void dtl_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if mn=='call':statement=f'dtl_call(m,{dl.read(x["operands"][0])},{nextoff},h);'
   elif mn=='cwd':statement='m->dx=(m->ax&0x8000)?0xffff:0;'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x)
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if ea in mutants:
    n,bad=mutants[ea];assert bad!=statement;lines += [f'#if KI_DETAILS_MUTATION == {n}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  assert rows[-1]['mnemonic']=='retn';lines.append('}')
 lines.append('int ki_details_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:dtl_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int dtl_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,rows,name in groups)+' return 1;default:return 0;}}')
 (repo/'tools/c_recovery/details_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated 10 functions /',sum(len(x[1]) for x in groups),'instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
