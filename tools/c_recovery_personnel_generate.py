#!/usr/bin/env python3
"""固定 IDA 人事七函式轉 C，保留原始定位；不複製共用規則。"""
import argparse,importlib.util,json
from pathlib import Path
def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
 assert len(groups)==7 and not r['decoded_blocks']
 lines=['/* Original native C personnel; re/121, spec/241. */']
 mutants={
  0x16AD8:(1,'dl_write8(m,m->ds,(uint16_t)(m->bx+23),3);'),
  0x16AE7:(2,'dl_write8(m,m->ds,(uint16_t)(m->si+25),(uint8_t)m->ax);'),
  0x16BB3:(3,'dl_write8(m,m->ds,(uint16_t)(m->bx+23),2);'),
  0x16BC2:(4,'dl_write8(m,m->ds,(uint16_t)(m->si+41),(uint8_t)(m->ax>>8));'),
  0x16B67:(5,'; /* lost administrator budget clear */'),
  0x16C42:(6,'; /* lost diplomat budget clear */'),
  0x16B4F:(7,'m->bx&=0x00ff; /* wrong sentinel */'),
  0x16C4A:(8,'m->flags|=1;'),
  0x16AD6:(9,'if(m->flags&1) goto L_16AFA;'),
  0x16279:(10,'psn_call(m,dl_read16(m,m->cs,(uint16_t)(0x6280+(m->bx^2))),0x627e,h);'),
  0x16AED:(11,'dl_write8(m,m->ds,(uint16_t)(m->bx+26),0);m->di=m->sp;'),
  0x16B2E:(12,'m->si=ec_pop(m);'),
 }
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows};lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void psn_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if mn=='call':statement=f'psn_call(m,{dl.read(x["operands"][0])},{nextoff},h);'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x)
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if ea in mutants:
    n,bad=mutants[ea];assert bad!=statement;lines += [f'#if KI_PERSONNEL_MUTATION == {n}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  assert rows[-1]['mnemonic']=='retn';lines.append('}')
 lines.append('int ki_personnel_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:psn_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int psn_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,rows,name in groups)+' return 1;default:return 0;}}')
 (repo/'tools/c_recovery/personnel_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated 7 functions /',sum(len(x[1]) for x in groups),'instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
