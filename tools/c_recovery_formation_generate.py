#!/usr/bin/env python3
"""固定 IDA 編成七函式轉 C，保留原始定位；不複製共用規則。"""
import argparse,importlib.util,json
from pathlib import Path
def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
 assert len(groups)==7 and not r['decoded_blocks']
 lines=['/* Original native C formation; re/120, spec/240. */']
 mutants={
  0x16D56:(1,'dl_write8(m,m->ds,(uint16_t)(m->si+42),2);'),
  0x16D81:(2,'dl_mul(m,9,16);'),
  0x16DC7:(3,'m->ax=dl_shift(m,m->ax,2,16,5);'),
  0x16E3F:(4,'vg_al(m,dl_read8(m,m->ds,(uint16_t)(m->di+24)));'),
  0x16D42:(5,'; /* lost cancellation pool return */'),
  0x16D0A:(6,'vg_al(m,2);'),
  0x16D18:(7,'goto L_16D24; /* lost empty leader refusal */'),
  0x16D24:(8,'; /* lost committed corps initialization */'),
  0x16C7C:(9,'; /* lost outer sidebar refresh */'),
  0x16D45:(10,'fc_call(m,0x6fd2,0x6d48,h);m->di=ec_math(m,m->di,4,16,0,1);'),
 }
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows};lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void fc_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if mn=='call':statement=f'fc_call(m,{dl.read(x["operands"][0])},{nextoff},h);'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x)
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if ea in mutants:
    n,bad=mutants[ea];assert bad!=statement;lines += [f'#if KI_FORMATION_MUTATION == {n}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  assert rows[-1]['mnemonic']=='retn';lines.append('}')
 lines.append('int ki_formation_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:fc_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int fc_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,rows,name in groups)+' return 1;default:return 0;}}')
 (repo/'tools/c_recovery/formation_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated 7 functions /',sum(len(x[1]) for x in groups),'instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
