#!/usr/bin/env python3
"""固定 IDA 行軍25函式轉 C，保留原始定位；不複製共用規則。"""
import argparse,importlib.util,json
from pathlib import Path
def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
 assert len(groups)==25 and not r['decoded_blocks']
 lines=['/* Original native C march; re/123, spec/243. */']
 mutants={
  0x17052:(1,'ec_cmp(m,m->bx,3,16);'),0x170C5:(2,'ec_cmp(m,(uint8_t)m->ax,0xd3,8);'),
  0x18028:(3,'; /* lost delegation clear */'),0x18032:(4,'dl_write8(m,m->ds,(uint16_t)(m->si+35),0);'),
  0x18075:(5,'m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|3);'),0x1434C:(6,'m->bx=ec_math(m,m->bx,4,16,0,0);'),
  0x143BB:(7,'ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si+4)),299,16);'),0x143DE:(8,'; /* lost random +1 */'),
  0x14439:(9,'ec_cmp(m,(uint8_t)m->ax,0,8);'),0x14470:(10,'ec_cmp(m,29,dl_read8(m,m->ds,m->bx),8);'),
  0x1466F:(11,'dl_write16(m,m->ds,m->si,0);'),0x17FB7:(12,'; /* lost ordered bit */'),
  0x170B3:(13,'m->ax=ec_math(m,m->ax,0,16,0,0);'),0x12069:(14,'dl_write16(m,m->ds,0x989c,m->ax);'),
  0x143D3:(15,'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->bx+24)),2,8);'),0x14561:(16,'; /* lost node comparison */'),
 }
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows};lines += [f'/* Original {name}; IDA 0x{at:X}. */',f'static void mch_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if mn=='call' and x['bytes'].startswith('9a'):statement=f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
   elif mn=='call':statement=f'mch_call(m,{dl.read(x["operands"][0])},{nextoff},h);'
   elif mn=='sahf':statement='m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
   elif mn=='les':
    op=x['operands'];off=dl.address(op[1]);seg=dl.segment(op[1]);statement='{uint16_t a='+off+',s='+seg+';uint16_t v=dl_read16(m,s,a);'+dl.write(op[0],'v')+'m->es=dl_read16(m,s,(uint16_t)(a+2));}'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x)
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if ea in mutants:
    n,bad=mutants[ea];assert bad!=statement;lines += [f'#if KI_MARCH_MUTATION == {n}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  assert rows[-1]['mnemonic']=='retn';lines.append('}')
 lines.append('int ki_march_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:mch_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int mch_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,rows,name in groups)+' return 1;default:return 0;}}')
 (repo/'tools/c_recovery/march_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated 25 functions /',sum(len(x[1]) for x in groups),'instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
