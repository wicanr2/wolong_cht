#!/usr/bin/env python3
"""固定 IDA 指令到 native C 運算／位址 label；不執行或嵌入 opcode。"""
import argparse
import hashlib
import json
import re
from pathlib import Path

REG8={'al':'vg_al','ah':'vg_ah','bl':'gw_bl','bh':'gw_bh','cl':'vg_cl','ch':'vg_ch','dl':'vg_dl','dh':'vg_dh'}
LOW={'al':'ax','bl':'bx','cl':'cx','dl':'dx'}
HIGH={'ah':'ax','bh':'bx','ch':'cx','dh':'dx'}
COND={'jz':'m->flags&0x40','jnz':'!(m->flags&0x40)','jb':'m->flags&1','jnb':'!(m->flags&1)',
      'jbe':'m->flags&(1|0x40)','ja':'!(m->flags&(1|0x40))','jl':'ec_less(m)','jge':'!ec_less(m)',
      'jle':'!ec_greater(m)','jg':'ec_greater(m)'}

def segment(op):
    s=op['original'].lower()
    for reg in ['cs','es','ss','ds']:
        if reg+':' in s:return 'm->'+reg
    return 'm->ss' if re.search(r'\bbp\b',s) else 'm->ds'

def address(op):
    if op['type']==2:return str(op['addr']&65535)
    regs=re.findall(r'\b(?:bx|bp|si|di)\b',op['original'].lower())
    if not regs:raise ValueError('unknown address '+op['original'])
    terms=['m->'+r for r in regs]
    if op['type']==4 and op['addr']&65535:terms.append(str(op['addr']&65535))
    return '(uint16_t)('+ '+'.join(terms)+')'

def read(op):
    t=op['type'];s=op['original'].lower()
    if t==1:
        if s in LOW:return '(uint8_t)m->'+LOW[s]
        if s in HIGH:return '(uint8_t)(m->'+HIGH[s]+'>>8)'
        return 'm->'+s
    if t==5:return str(op['value']&((1<<(op['dtype_size']*8))-1))
    if t in [2,3,4]:return ('dl_read8' if op['dtype_size']==1 else 'dl_read16')+'(m,'+segment(op)+','+address(op)+')'
    if t==7:return str(op['addr']&65535)
    raise ValueError('unsupported read '+str(op))

def write(op,expr):
    if op['type']==1:
        r=op['original'].lower()
        if r in REG8:return REG8[r]+'(m,(uint8_t)('+expr+'));'
        return 'm->'+r+'=(uint16_t)('+expr+');'
    if op['type'] in [2,3,4]:return ('dl_write8' if op['dtype_size']==1 else 'dl_write16')+'(m,'+segment(op)+','+address(op)+','+expr+');'
    raise ValueError('unsupported destination '+str(op))

def statement(row):
    mn=row['mnemonic'];o=row['operands'];ea=row['ida_linear'];size=len(bytes.fromhex(row['bytes']));nextoff=(ea+size-0x10000)&65535
    if mn=='mov':return write(o[0],read(o[1]))
    if mn=='push':return 'ec_push(m,'+read(o[0])+');'
    if mn=='pop':return write(o[0],'ec_pop(m)')
    if mn=='xchg':return '{uint16_t a='+read(o[0])+',b='+read(o[1])+';'+write(o[0],'b')+write(o[1],'a')+'}'
    if mn in ['add','sub','adc','sbb']:
        return write(o[0],'ec_math(m,'+read(o[0])+','+read(o[1])+','+str(o[0]['dtype_size']*8)+','+('m->flags&1' if mn in ['adc','sbb'] else '0')+','+('1' if mn in ['sub','sbb'] else '0')+')')
    if mn=='cmp':return 'ec_cmp(m,'+read(o[0])+','+read(o[1])+','+str(o[0]['dtype_size']*8)+');'
    if mn in ['and','or','xor','test']:
        op={'and':'&','or':'|','xor':'^','test':'&'}[mn];expr='('+read(o[0])+')'+op+'('+read(o[1])+')';width=o[0]['dtype_size']*8
        return '{uint16_t result=('+expr+')&'+str((1<<width)-1)+';ec_logic(m,result,'+str(width)+');'+('' if mn=='test' else write(o[0],'result'))+'}'
    if mn in ['inc','dec']:return write(o[0],('ec_inc' if mn=='inc' else 'st_dec')+'(m,'+read(o[0])+','+str(o[0]['dtype_size']*8)+')')
    if mn=='not':return write(o[0],'~('+read(o[0])+')')
    if mn=='neg':return write(o[0],'ec_math(m,0,'+read(o[0])+','+str(o[0]['dtype_size']*8)+',0,1)')
    if mn in ['shl','shr','rol','ror']:
        return write(o[0],'dl_shift(m,'+read(o[0])+','+read(o[1])+','+str(o[0]['dtype_size']*8)+','+str({'shl':4,'shr':5,'rol':0,'ror':1}[mn])+')')
    if mn in ['mul','div']:
        visible=[op for op in o if op['original']];assert len(visible)==1
        return 'dl_'+mn+'(m,'+read(visible[0])+','+str(visible[0]['dtype_size']*8)+');'
    if mn in COND:return 'if('+COND[mn]+') goto L_'+format(o[0]['addr']+0x10000,'X')+';'
    if mn=='jmp':return 'goto L_'+format(o[0]['addr']+0x10000,'X')+';'
    if mn=='loop':return 'if(--m->cx) goto L_'+format(o[0]['addr']+0x10000,'X')+';'
    if mn=='call':return 'dl_call(m,'+read(o[0])+','+str(nextoff)+',h);'
    if mn=='retn':return 'return;'
    if mn=='cld':return 'm->flags&=(uint16_t)~0x400u;'
    if mn=='clc':return 'm->flags&=(uint16_t)~1u;'
    if mn=='stc':return 'm->flags|=1;'
    if mn=='out':return 'dl_out(m,'+read(o[0])+','+read(o[1])+','+str(o[1]['dtype_size']*8)+');'
    if mn=='movsb':return 'vg_movsb(m);'
    if mn=='lds':
        return '{uint16_t off='+address(o[1])+',seg='+segment(o[1])+';uint16_t value=dl_read16(m,seg,off),bank=dl_read16(m,seg,(uint16_t)(off+2));'+write(o[0],'value')+'m->ds=bank;}'
    raise ValueError('unsupported instruction '+row['assembly'])

def generate(repo,probe):
    r=json.loads(probe.read_text());assert r['input_sha256']=='fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
    targets={t['ida_linear']:t for t in r['targets']}
    wanted=[0x1030f,0x10337,0x1e993,0x1e9a7,0x1e9c1,0x1fb29,0x1fba7,0x1f6dc,0x1f878]
    groups=[]
    for ea in wanted:
        t=targets[ea];rows=[x for c in t['chunks'] for x in c['instructions']];groups.append((ea,rows,t['name']))
    for b in r['decoded_blocks']:
        rows=[x for x in b['instructions'] if not x.get('inline_data')]
        if b['ida_linear']==0x1f26e:
            limits=[0x1f26e,0x1f3cf,0x1f426,0x1f43b,0x1f45b]
            for a,z in zip(limits,limits[1:]):groups.append((a,[x for x in rows if a<=x['ida_linear']<z],b['original_name_before_analysis'] if a==b['ida_linear'] else None))
        else:groups.append((b['ida_linear'],rows,b['original_name_before_analysis']))
    lines=['/* Generated from fixed IDA evidence; native C statements, no opcode fetch.',' * Original names/address/operands remain in comments; spec/227, re/107. */']
    aliases={}
    for ea,rows,name in groups:
        labels={x['ida_linear'] for x in rows}
        aliases[ea]={ea}|{0x10000+x['operands'][0]['addr'] for x in rows if x['mnemonic']=='call' and x['operands'][0]['type']==7 and 0x10000+x['operands'][0]['addr'] in labels}
    for ea,rows,name in groups:lines.append(f'static void dl_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry);')
    for ea,rows,name in groups:
        assert rows
        labels={x['ida_linear'] for x in rows}
        lines.append(f'/* original {name or "(no IDA function)"}; IDA 0x{ea:X}; raw control flow; evidence re/107 */')
        lines.append(f'static void dl_body_{ea:X}(KiMachine16 *m,const KiEconomyHooks *h,uint16_t entry) {{')
        lines.append('switch(entry) {'+' '.join(f'case 0x{a-0x10000:x}:goto L_{a:X};' for a in sorted(aliases[ea]))+' default:assert(0);}')
        for row in rows:
            if row['mnemonic'] in list(COND)+['jmp','loop']:
                assert 0x10000+row['operands'][0]['addr'] in labels,(hex(ea),row)
            lines.append(f'L_{row["ida_linear"]:X}: ; /* IDA {row["ida_linear"]:#x}; {row["assembly"]} */')
            normal=statement(row);mutated=None;mutation=0
            if row['ida_linear']==0x1e9ea:mutation=1;mutated=normal.replace(',12,',',24,')
            elif row['ida_linear']==0x1eada:mutation=2;mutated=normal.replace('m->si+2','m->si+4')
            elif row['ida_linear']==0x1e9f9:mutation=3;mutated=normal.replace('m->ds','m->cs')
            elif row['ida_linear']==0x1eac6:mutation=4;mutated=normal.replace('60143','60141')
            elif row['ida_linear']==0x1f6f5:mutation=5;mutated='; /* lost byte rewind */'
            elif row['ida_linear']==0x1f6f7:mutation=6;mutated=normal.replace('&(4)','&(0)')
            elif row['ida_linear']==0x1f3b2:mutation=7;mutated=normal.replace('dl_read16(m,m->cs,62555)','(dl_read16(m,m->cs,62555)+1)')
            elif row['ida_linear']==0x1fbf0:mutation=8;mutated=normal.replace(',80,',',79,')
            elif row['ida_linear']==0x1f477:mutation=9;mutated=normal.replace('m->si','m->dx')
            elif row['ida_linear']==0x1e9b0:mutation=10;mutated=normal.replace(',60145,',',60146,')
            if mutated is not None:
                assert mutated!=normal,(hex(row['ida_linear']),normal)
                lines.extend([f'#if KI_DISPLAY_MUTATION == {mutation}','    '+mutated,'#else','    '+normal,'#endif'])
            else:lines.append('    '+normal)
        lines.append('}')
    lines.append('int ki_display_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {switch(t) {')
    for ea,rows,name in groups:
        for a in sorted(aliases[ea]):lines.append(f'case 0x{a-0x10000:x}:dl_body_{ea:X}(m,h,t);break;')
    lines.extend(['default:return 0;}return 1;}'])
    path=repo/'tools/c_recovery/display_generated.inc';path.write_text('\n'.join(lines)+'\n')
    print('Native C entries',sum(len(x) for x in aliases.values()),'instructions',sum(len(x[1]) for x in groups),'SHA256',hashlib.sha256(path.read_bytes()).hexdigest())

if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--repo',type=Path,required=True);p.add_argument('--probe',type=Path,required=True)
    a=p.parse_args();generate(a.repo,a.probe)
