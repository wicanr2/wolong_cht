#!/usr/bin/env python3
"""由固定 IDA 還原四類清單；共用引擎保持唯一實作。"""
import argparse,importlib.util,json
from pathlib import Path

def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 groups=[(t['ida_linear'],[x for c in t['chunks'] for x in c['instructions']],t['name']) for t in r['targets']]
 groups += [(t['ida_linear'],t['instructions'],t['original_name_before_analysis']) for t in r['decoded_blocks']]
 lines=['/* Original native C list families; re/119, spec/239. */']
 for at,rows,name in groups:
  labels={x['ida_linear'] for x in rows}
  lines += [f'/* Original {name or "unnamed code"}; IDA 0x{at:X}. */',f'static void cc_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{']
  for x in rows:
   ea,mn=x['ida_linear'],x['mnemonic'];nextoff=(ea+len(bytes.fromhex(x['bytes']))-0x10000)&65535
   if ea in [0x1722a,0x17232]:statement=f'm->dx=dl_read16(m,m->cs,0x{ea-0x10000+1:x});'
   elif ea==0x17643:statement='vg_al(m,dl_read8(m,m->cs,0x7644));'
   elif mn=='call':statement=f'cc_call(m,{dl.read(x["operands"][0])},{nextoff},h);'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert x['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(x)
   mutation,bad=0,None
   if ea==0x171c9:mutation,bad=1,statement.replace(',16832,',',16896,')
   elif ea==0x1722a:mutation,bad=2,'m->dx=0x1234;'
   elif ea==0x17643:mutation,bad=3,'vg_al(m,0);'
   elif ea==0x176c7:mutation,bad=4,'; /* lost exclusion of ruler */'
   elif ea==0x1768a:mutation,bad=5,'vg_cl(m,0);'
   elif ea==0x1735c:mutation,bad=6,'m->ax=dl_read16(m,m->ds,(uint16_t)(m->si+6));'
   elif ea==0x17361:mutation,bad=6,'; /* erroneous word morale retained */'
   elif ea==0x17741:mutation,bad=7,'; /* lost affiliation fallback gate */'
   elif ea==0x17776:mutation,bad=8,'; /* lost duty priority */'
   elif ea==0x17957:mutation,bad=9,'; /* lost exclusion of player */'
   elif ea==0x17aaa:mutation,bad=10,'; /* lost strict relation boundary */'
   elif ea==0x17a9a:mutation,bad=11,'if(m->cx==22) {vg_dl(m,0);ec_logic(m,0,8);} /* erroneous cross-row state */'
   elif ea==0x17b44:mutation,bad=12,statement.replace('(136)','(24)')
   elif ea==0x17a5c:mutation,bad=13,statement.replace('(36867)','(36868)')
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {x["assembly"]} */')
   if bad is not None:
    assert bad!=statement,(hex(ea),statement);lines += [f'#if KI_CATALOG_MUTATION == {mutation}','    '+bad,'#else','    '+statement,'#endif']
   else:lines.append('    '+statement)
  assert rows[-1]['mnemonic']=='retn';lines.append('}')
 lines.append('int ki_catalog_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,rows,name in groups:lines.append(f'case 0x{at-0x10000:x}:cc_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int cc_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,rows,name in groups)+' return 1;default:return 0;}}')
 (repo/'tools/c_recovery/catalog_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated',len(r['targets']),'functions,',len(r['decoded_blocks']),'raw entries;',sum(len(x[1]) for x in groups),'instructions')
if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
