#!/usr/bin/env python3
"""固定 IDA 據點輪轉與動畫完整閉包轉 C，保留候選及軍團搜尋原迴圈。"""
import argparse,hashlib,importlib.util,json
from pathlib import Path
EXPECTED='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
TARGETS={68830: ('sub_10CDE', 9), 74841: ('sub_12459', 49), 74890: ('sub_1248A', 117), 75007: ('sub_124FF', 52), 81661: ('sub_13EFD', 119), 81780: ('sub_13F74', 53), 81833: ('sub_13FA9', 127), 81960: ('sub_14028', 47), 82007: ('sub_14057', 92), 82099: ('sub_140B3', 22), 82121: ('sub_140C9', 140), 82261: ('sub_14155', 63), 82324: ('sub_14194', 162), 82537: ('sub_14269', 66), 83317: ('sub_14575', 76), 83393: ('sub_145C1', 55), 125713: ('sub_1EB11', 77), 125790: ('sub_1EB5E', 14)}
MUTANT_STATEMENTS={74851: (1, 'ec_cmp(m,dl_read8(m,m->ds,m->si),0x81,8);'), 74848: (2, 'm->bp=31;'), 74870: (3, 'ec_cmp(m,m->si,0x2150,16);'), 75010: (4, '{uint8_t result=(uint8_t)m->ax&3;ec_logic(m,result,8);vg_al(m,result);}'), 75020: (5, 'ec_cmp(m,(uint8_t)(m->dx>>8),0xf2,8);'), 75025: (5, 'vg_dh(m,0xf2);'), 75038: (6, 'ec_cmp(m,(uint8_t)m->dx,0x10,8);'), 81763: (7, 'ec_cmp(m,m->si,0x17e0,16);'), 81705: (8, 'dl_write8(m,m->ds,(uint16_t)(m->bx+0x16),(uint8_t)(m->ax>>8));'), 81900: (9, 'ec_cmp(m,(uint8_t)m->ax,0x17,8);'), 81908: (10, 'ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->di+0x600)),0x81,8);'), 81960: (11, '{uint8_t result=dl_read8(m,m->ds,(uint16_t)(m->si+0x840))&0x7f;ec_logic(m,result,8);dl_write8(m,m->ds,(uint16_t)(m->si+0x840),result);}'), 82324: (12, 'vg_cl(m,7);'), 82404: (13, 'ec_cmp(m,(uint8_t)m->ax,0xc7,8);'), 82408: (13, 'vg_al(m,0xc7);'), 82575: (14, 'dl_mul(m,dl_read8(m,m->ds,(uint16_t)(m->si+0x84e)),8);'), 82231: (15, 'm->dx=ec_math(m,m->dx,dl_read16(m,m->ds,(uint16_t)(m->bx+0x84a)),16,0,1);'), 125725: (16, '{uint8_t result=(uint8_t)m->ax|1;ec_logic(m,result,8);vg_al(m,result);}')}
def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py')
 dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 data=json.loads(probe.read_text(encoding='utf-8'))
 assert data['input_sha256']==EXPECTED and set(data['recovery_targets'])==set(TARGETS)
 groups=[]
 for target in sorted(data['targets'],key=lambda row:row['ida_linear']):
  at=target['ida_linear']
  if at not in TARGETS:continue
  raw=b''.join(bytes.fromhex(chunk['file_bytes']) for chunk in target['chunks'])
  assert (target['name'],len(raw))==TARGETS[at] and hashlib.sha256(raw).hexdigest()==target['file_sha256']
  assert not any(chunk['loader_relocations'] for chunk in target['chunks'])
  groups.append((at,[row for chunk in target['chunks'] for row in chunk['instructions']],target['name']))
 assert len(groups)==18 and sum(len(rows) for _,rows,_ in groups)==562
 lines=['/* Original native C city/object tick; re/128, spec/248. No scheduler or full main loop. */']
 seen=set()
 for at,rows,name in groups:
  labels={row['ida_linear'] for row in rows}
  lines.extend([f'/* Original {name}; IDA 0x{at:X}. */',f'static void ct_body_{at:X}(KiMachine16 *m,const KiEconomyHooks *h) {{'])
  for row in rows:
   ea,mn=row['ida_linear'],row['mnemonic'];nextoff=(ea+len(bytes.fromhex(row['bytes']))-0x10000)&65535
   if mn=='call' and row['bytes'].startswith('9a'):statement=f'mx_far(m,0x{ea-0x10000:x},0x{nextoff:x},h);'
   elif mn=='call':statement=f'ct_call(m,{dl.read(row["operands"][0])},{nextoff},h);'
   elif mn=='in':
    assert row['operands'][0]['dtype_size']==1
    statement=dl.write(row['operands'][0],f'wolong_main_in8({dl.read(row["operands"][1])})')
   elif mn=='stosw':
    statement=f'mc_string(m,2,{1 if row["bytes"].startswith("f3") else 0});'
   else:
    if mn in list(dl.COND)+['jmp','loop']:assert row['operands'][0]['addr']+0x10000 in labels
    statement=dl.statement(row)
   lines.append(f'L_{ea:X}: ; /* IDA 0x{ea:X}; {row["assembly"]} */')
   if ea in MUTANT_STATEMENTS:
    number,bad=MUTANT_STATEMENTS[ea];assert bad!=statement;seen.add(ea)
    lines.extend([f'#if KI_TICK_MUTATION == {number}','    '+bad,'#else','    '+statement,'#endif'])
   else:lines.append('    '+statement)
  assert rows[-1]['mnemonic']=='retn';lines.append('}')
 assert seen==set(MUTANT_STATEMENTS) and {n for n,_ in MUTANT_STATEMENTS.values()}==set(range(1,17))
 lines.append('int ki_tick_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for at,_,_ in groups:lines.append(f'case 0x{at-0x10000:x}:ct_body_{at:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 lines.append('static int ct_known(uint16_t t) {switch(t) {'+' '.join(f'case 0x{at-0x10000:x}:' for at,_,_ in groups)+' return 1;default:return 0;}}')
 for at,_,name in groups:lines.append(f'void {name}(KiMachine16 *m,const KiEconomyHooks *h) {{ki_tick_body(m,0x{at-0x10000:x},h);m->ip=ec_pop(m);}}')
 (repo/'tools/c_recovery/tick_generated.inc').write_text('\n'.join(lines)+'\n',encoding='utf-8')
 print('Generated 18 original functions / 562 instructions / 1340 bytes; 16 mutations at 18 sites')
if __name__=='__main__':
 parser=argparse.ArgumentParser(description=__doc__);parser.add_argument('--repo',type=Path,required=True);parser.add_argument('--probe',type=Path,required=True)
 args=parser.parse_args();generate(args.repo,args.probe)
