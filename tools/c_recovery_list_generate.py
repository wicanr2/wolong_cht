#!/usr/bin/env python3
"""由 IDA 原 CFG 還原清單、排序與遷都，保留特殊返回。"""
import argparse,importlib.util,json
from pathlib import Path

def generate(repo,probe):
 s=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(s);s.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
 groups += [(t['ida_linear'],t['instructions'],t['original_name_before_analysis']) for t in r['decoded_blocks']]
 lines=['/* Original native C city list / sort / relocation; re/118, spec/238. */']
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows}
  lines += [f'/* Original {name or "unnamed code"}; IDA 0x{at:X}. */',f'static int ll_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if ea in [0x1831f,0x183d2]:statement=f'm->si=dl_read16(m,m->cs,0x{ea-0x10000+1:x});'
   elif ea==0x1858f:statement='m->bx=ec_math(m,m->bx,dl_read16(m,m->cs,0x8591),16,0,0);'
   elif ea in [0x185e8,0x185f3]:statement=f'll_sort_load(m,0x{ea-0x10000:x});'
   elif ea==0x185ef:statement='ll_sort_compare(m);'
   elif ea==0x185f1:statement='if(ll_sort_skip(m)) goto L_185F8;'
   elif ea==0x18455:statement='m->ip=0x8412;if(h&&h->enter) h->enter(m,0x8412,h->user);goto L_18412;'
   elif mn=='call' and bytes.fromhex(x['bytes'])[0]==0x9a:statement=f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
   elif mn=='call':statement=f'if(ll_call(m,{dl.read(x["operands"][0])},{nextoff},h)) return 1;'
   elif mn=='cwd':statement='m->dx=m->ax&0x8000?0xffff:0;'
   elif mn=='lahf':statement='vg_ah(m,(uint8_t)((m->flags&0xd5)|2));'
   elif mn=='sahf':statement='m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5)|2);'
   elif mn=='pushf':statement='ec_push(m,m->flags);'
   elif mn=='popf':statement='m->flags=ec_pop(m)|2;'
   elif mn in ['stosw','lodsw']:statement=f'mc_string(m,{2 if mn=="stosw" else 1},{1 if bytes.fromhex(x["bytes"])[0] in [0xf2,0xf3] else 0});'
   elif mn=='retn':statement='return 0;'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x)
   mutation,bad=0,None
   if ea==0x17456:mutation,bad=1,statement.replace(',8256,',',8224,')
   elif ea==0x17496:mutation,bad=2,'ec_cmp(m,0xffff,0xffff,16);'
   elif ea==0x183d2:mutation,bad=3,'m->si=0x7378+2;'
   elif ea==0x185e8:mutation,bad=4,'m->ax=dl_read16(m,m->ds,(uint16_t)(m->bx+m->si));'
   elif ea==0x185f1:mutation,bad=5,'if(m->flags&1) goto L_185F8;'
   elif ea==0x184b7:mutation,bad=6,'; /* lost removal of selection-loop return address */'
   elif ea==0x18566:mutation,bad=7,'; /* lost scroll increment */'
   elif ea==0x18734:mutation,bad=8,'dl_div(m,(uint8_t)m->dx+1,8);'
   elif ea==0x1695d:mutation,bad=9,'if(!(m->flags&1)) goto L_16969;'
   elif ea==0x16965:mutation,bad=10,'if(m->flags&1) goto L_16969;'
   elif ea==0x1698d:mutation,bad=11,'; /* lost actual relocation writer */'
   elif ea==0x184d2:mutation,bad=12,statement.replace('(16)','(15)')
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if bad is not None:
    assert bad!=statement,(hex(ea),statement);lines += [f'#if KI_LIST_MUTATION == {mutation}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  if at==0x18412:lines += ['m->ip=0x8458;if(h&&h->enter) h->enter(m,0x8458,h->user);return ll_body_18458(m,h);']
  else:assert rows[-1]['mnemonic']=='retn',(hex(at),rows[-1])
  lines.append('}')
 # The cancellation path falls into nullsub_2; declare every body before definitions.
 lines[1:1]=[f'static int ll_body_{at:X}(KiMachine16 *,const KiEconomyHooks *);' for at,rows,name in groups]
 lines.append('int ki_list_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h,int *escaped) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:*escaped=ll_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int ll_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,rows,name in groups)+' return 1;default:return 0;}}')
 (repo/'tools/c_recovery/list_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated',len(r['targets']),'functions,',len(r['decoded_blocks']),'raw entries;',sum(len(x[1]) for x in groups),'instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
