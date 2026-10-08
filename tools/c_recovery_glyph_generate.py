#!/usr/bin/env python3
"""固定 IDA glyph 指令到 native C；far／INT 為原始 ABI 的平台邊界。"""
import argparse
import importlib.util
import json
from pathlib import Path

def generate(repo,probe):
 spec=importlib.util.spec_from_file_location('dl',repo/'tools/c_recovery_display_generate.py');dl=importlib.util.module_from_spec(spec);spec.loader.exec_module(dl)
 r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
 targets={t['ida_linear']:t for t in r['targets']};groups=[]
 for ea in [0x1f720,0x1f7a4,0x106f5,0x106fd]:groups.append((ea,[x for c in targets[ea]['chunks'] for x in c['instructions']],targets[ea]['name']))
 for b in r['decoded_blocks']:groups.append((b['ida_linear'],b['instructions'],b['original_name_before_analysis']))
 lines=['/* Generated native C glyph control; locator=proven, semantic evidence re/108/spec/228.',' * Original names/addresses/operands retained. No opcode dispatch. */']
 for ea,rows,name in groups:
  lines.append(f'/* Original {name or "(no IDA function)"}; IDA 0x{ea:X}; raw control; source re/108 */')
  lines.append(f'static void gl_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h) {{')
  labels={x['ida_linear'] for x in rows}
  for x in rows:
   ea0=x['ida_linear'];size=len(bytes.fromhex(x['bytes']));nextoff=(ea0+size-0x10000)&65535
   if x['mnemonic'] in list(dl.COND)+['jmp','loop']:assert 0x10000+x['operands'][0]['addr'] in labels
   if x['mnemonic']=='int':statement=f'gl_interrupt(m,{x["operands"][0]["value"]},{nextoff});'
   elif x['mnemonic']=='call' and bytes.fromhex(x['bytes'])[0]==0x9a:statement=f'gl_far(m,0x{ea0-0x10000:x},0x{nextoff:x},h);'
   elif ea0 in [0x1071e,0x10735,0x1074c]:statement=f'm->ax=dl_read16(m,m->cs,0x{ea0-0x10000+1:x});'
   else:statement=dl.statement(x).replace('dl_call(', 'gl_call(')
   lines.append(f'L_{ea0:X}: ; /* IDA {ea0:#x}; {x["assembly"]} */')
   mutation=0;bad=None
   if ea0==0x1f774:mutation=1;bad=statement.replace(',624,',',623,')
   elif ea0==0x1f77a:mutation=2;bad=statement.replace(',384,',',383,')
   elif ea0 in [0x1f7f8,0x1f849]:mutation=3;bad='; /* lost byte exchange */'
   elif ea0 in [0x1f804,0x1f857]:mutation=4;bad='; /* lost narrow stride */'
   elif ea0 in [0x1f80d,0x1f812,0x1f817,0x1f860,0x1f865,0x1f86a]:mutation=5;bad='; /* lost latch read */'
   elif ea0 in [0x1f81b,0x1f86e]:mutation=6;bad=statement.replace(',78,',',77,')
   elif ea0==0x1f824:mutation=7;bad=statement.replace('&(2)','&(0)')
   elif ea0 in [0x1071e,0x10735,0x1074c]:mutation=8;bad='m->ax=0x9001;'
   elif ea0==0x10728:mutation=9;bad=statement.replace(',16,16',',8,16')
   if bad is not None:
    assert bad!=statement,(hex(ea0),statement)
    lines.extend([f'#if KI_GLYPH_MUTATION == {mutation}','    '+bad,'#else','    '+statement,'#endif'])
   else:lines.append('    '+statement)
  lines.append('}')
 lines.append('int ki_glyph_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
 for ea,rows,name in groups:lines.append(f'case 0x{ea-0x10000:x}:gl_body_{ea:X}(m,h);break;')
 lines.append('default:return 0;}return 1;}')
 (repo/'tools/c_recovery/glyph_generated.inc').write_text('\n'.join(lines)+'\n')
 print('Generated glyph entries',len(groups),'instructions',sum(len(g[1]) for g in groups))

if __name__=='__main__':
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True);a=p.parse_args();generate(a.repo,a.probe)
