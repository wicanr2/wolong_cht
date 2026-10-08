#!/usr/bin/env python3
"""將 IDA 的檔案定位 inventory 轉成可獨立組譯的 DOS MZ 來源。

指令以助憶碼重組。資料、未知 bytes 與未支援的指令分別記錄，不能混算。
所有 workload 須在 Docker 內執行。原始 KI.EXE 只作匯出／驗證輸入。
"""
import argparse
import hashlib
import json
import re
import struct
import subprocess
from pathlib import Path

EXPECTED = 'fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868'
PREFIXES = {0x26:'es',0x2e:'cs',0x36:'ss',0x3e:'ds',0xf0:'lock',0xf2:'repne',0xf3:'rep'}
REGS = re.compile(r'\b(?:[abcd][lh]|[abcd]x|[sd]i|[bs]p|e(?:[abcd]x|[sd]i|[bs]p))\b',re.I)
WIDTHS = {1:'byte',2:'word',4:'dword',8:'qword',10:'tbyte'}


def run(args, timeout=60):
    r=subprocess.run(args,capture_output=True,text=True,timeout=timeout)
    if r.returncode:raise RuntimeError(r.stderr)
    return r.stdout


def elf_sections(path, absolute_symbols=None):
    raw=path.read_bytes()
    if raw[:6]!=b'\x7fELF\x01\x01':raise ValueError('expected little-endian ELF32')
    shoff=struct.unpack_from('<I',raw,32)[0]
    entsize, count, names_idx=struct.unpack_from('<HHH',raw,46)
    if not count:count=struct.unpack_from('<I',raw,shoff+20)[0]
    if names_idx==0xffff:names_idx=struct.unpack_from('<I',raw,shoff+24)[0]
    sections=[struct.unpack_from('<10I',raw,shoff+i*entsize) for i in range(count)]
    names=raw[sections[names_idx][4]:sections[names_idx][4]+sections[names_idx][5]]
    result={}
    section_names=[]
    for s in sections:
        name=names[s[0]:].split(b'\0',1)[0].decode('ascii')
        section_names.append(name)
        if s[1]==1:result[name]=raw[s[4]:s[4]+s[5]]
    for s in sections:
        if s[1]!=9:continue
        symtab=sections[s[6]];strtab=sections[symtab[6]]
        strings=raw[strtab[4]:strtab[4]+strtab[5]]
        name=section_names[s[7]]
        data=bytearray(result[name])
        for at in range(s[4],s[4]+s[5],s[9]):
            offset,info=struct.unpack_from('<II',raw,at)
            if info&255!=20:raise ValueError(f'unsupported candidate relocation {info&255}')
            sym=struct.unpack_from('<IIIBBH',raw,symtab[4]+(info>>8)*symtab[9])
            symbol=strings[sym[0]:].split(b'\0',1)[0].decode('ascii')
            value=(absolute_symbols or {})[symbol]
            addend=int.from_bytes(data[offset:offset+2],'little')
            data[offset:offset+2]=((value+addend)&0xffff).to_bytes(2,'little')
        result[name]=bytes(data)
    return result


def numeric_operand(op, row):
    kind=op['type'];original=op['original']
    if kind==1:return original
    if kind==11 and re.fullmatch(r'st(?:\([0-7]\))?',original):return original
    if kind==5:
        value=op['value']
        for at in row['relocation_file_offsets']:
            offset=at-row['file_offset']
            value=int.from_bytes(bytes.fromhex(row['bytes'])[offset:offset+2],'little')
        value &= (1<<(op['dtype_size']*8))-1
        return hex(value)
    if kind in (2,3,4):
        regs=REGS.findall(original)
        disp=op['addr'] if kind in (2,4) else 0
        if kind==2:address=hex(disp)
        else:
            address='+'.join(r.lower() for r in regs)
            if not address:raise ValueError(f'cannot decode address: {original}')
            if disp:address+='+'+hex(disp)
        width=WIDTHS.get(op['dtype_size'])
        if width is None:raise ValueError(f'unsupported memory width: {op["dtype_size"]}')
        return f'{width} ptr [{address}]'
    raise ValueError(f'unsupported operand type {kind}: {original}')


def translate(row):
    raw=bytes.fromhex(row['bytes']);p=0;prefix=[]
    while p<len(raw) and raw[p] in (*PREFIXES,0x66,0x67):
        if raw[p] in PREFIXES:prefix.append(PREFIXES[raw[p]])
        p+=1
    mnemonic={'retn':'ret','xlat':'xlatb'}.get(row['mnemonic'],row['mnemonic'])
    operands=row['operands']
    if raw[p] in (0x9a,0xea):
        offset,segment=struct.unpack_from('<HH',raw,p+1)
        instruction=f'{"lcall" if raw[p]==0x9a else "ljmp"} {hex(segment)}, {hex(offset)}'
    elif operands and operands[0]['type']==7:
        relative=int.from_bytes(raw[p+1:], 'little', signed=True)
        delta=relative+len(raw)-len(prefix)
        instruction=f'{mnemonic} .{delta:+d}'
        if raw[p]==0xe9:instruction='{disp16} '+instruction
    elif mnemonic in ('movsb','movsw','stosb','stosw','lodsb','lodsw','cmpsb'):
        instruction=mnemonic
    else:
        if raw[p]==0xff and len(raw)>p+1 and (raw[p+1]>>3)&7 in (3,5):
            mnemonic='lcall' if (raw[p+1]>>3)&7==3 else 'ljmp'
        ops=', '.join(numeric_operand(op,row) for op in operands if op['original'])
        instruction=mnemonic+(' '+ops if ops else '')
    # 用前綴助憶碼保留明示的段／REP 前綴，不注入原始 instruction bytes。
    return '; '.join(prefix+[instruction])


def compile_candidates(rows, output):
    prefixes=['','{load} ','{store} ','{nooptimize} ',
              '{load} {nooptimize} ','{store} {nooptimize} ',
              '{disp16} ','{disp8} ','{disp16} {load} ','{disp16} {store} ', 'imm16-symbol']
    expressions={};errors={}
    for row in rows:
        try:expressions[row['file_offset']]=translate(row)
        except ValueError as e:errors[row['file_offset']]=str(e)
    pending={r['file_offset']:r for r in rows if r['file_offset'] in expressions}
    matched={};trials=[];constant_symbols={}
    for n,hint in enumerate(prefixes):
        if not pending:break
        candidates={};lines=['.intel_syntax noprefix','.code16']
        line_to_at={}
        for at,row in pending.items():
            text=expressions[at]
            if hint=='imm16-symbol':
                immediates=[op for op in row['operands'] if op['type']==5 and op['dtype_size']==2]
                if len(immediates)!=1:continue
                value=immediates[0]['value']&0xffff
                symbol=f'imm_at_{row["ida_linear"]:X}'
                constant_symbols[symbol]=value
                text=text.rsplit(', ',1)[0]+', offset '+symbol
            # 段前綴保持在 hint 之前。
            parts=text.split('; ');parts[-1]=('' if hint=='imm16-symbol' else hint)+parts[-1];text='; '.join(parts)
            candidates[at]=text
            lines.append(f'.section .i_{at:x},"ax",@progbits')
            lines.append(text);line_to_at[len(lines)]=at
        src=output/f'candidate-{n}.S';obj=output/f'candidate-{n}.o'
        active=set(candidates)
        # 診斷必須指到個別來源行；未識別的全域錯誤直接失敗。
        for attempt in range(3):
            src.write_text('\n'.join(lines)+'\n',encoding='utf-8')
            r=subprocess.run(['as','--32','-o',str(obj),str(src)],capture_output=True,text=True,timeout=45)
            if r.returncode==0:break
            bad=set()
            for line in r.stderr.splitlines():
                m=re.search(r':(\d+): Error:',line)
                if m and int(m[1]) in line_to_at:bad.add(line_to_at[int(m[1])])
            if not bad:raise RuntimeError(r.stderr[:2000])
            for line_no,at in line_to_at.items():
                if at in bad:lines[line_no-1]='';active.discard(at);errors[at]=r.stderr.splitlines()[0]
        else:raise RuntimeError('assembler errors did not converge')
        if r.returncode:raise RuntimeError(r.stderr[:2000])
        sections=elf_sections(obj,constant_symbols)
        found=0
        for at in active:
            got=sections.get(f'.i_{at:x}',b'')
            if got==bytes.fromhex(pending[at]['bytes']):
                matched[at]=candidates[at];found+=1
        for at in matched:pending.pop(at,None)
        trials.append({'hint':hint,'matched':found,'remaining':len(pending)})
        print(trials[-1],flush=True)
    used_symbols={s:v for s,v in constant_symbols.items() if any(s in expr for expr in matched.values())}
    return matched,{at:errors.get(at,'encoding differs') for at in set(r['file_offset'] for r in rows)-set(matched)},trials,used_symbols


def data_source(data):
    if data and not any(data):return ['.zero '+str(len(data))]
    return ['.byte '+','.join(hex(v) for v in data[i:i+16]) for i in range(0,len(data),16)]


def rebuild(original_path, inventory_path, output, image_id):
    original=original_path.read_bytes();inventory=json.loads(inventory_path.read_text(encoding='utf-8'))
    sha=hashlib.sha256(original).hexdigest()
    if sha!=EXPECTED or inventory['input_sha256']!=sha or inventory['schema']!='wolong-assembly-inventory-v1':
        raise ValueError('original / inventory identity differs')
    rows=inventory['code'];output.mkdir(parents=True,exist_ok=True)
    for row in rows:
        at=row['file_offset'];raw=bytes.fromhex(row['bytes'])
        if original[at:at+len(raw)]!=raw:raise ValueError('inventory bytes differ from original')
    matched,unmatched,trials,symbols=compile_candidates(rows,output)
    if unmatched:
        (output/'unmatched-instructions.json').write_text(json.dumps(unmatched,indent=2)+'\n',encoding='utf-8')
        raise ValueError(f'{len(unmatched)} instructions cannot be assembled; no code-byte fallback allowed')
    source=['# Local-only reconstructed KI.EXE; data and unknown spans retain original bytes.',
            '.intel_syntax noprefix','.code16','.section .image,"ax",@progbits']
    fields=struct.unpack_from('<14H',original)
    field_names=['e_magic','e_cblp','e_cp','e_crlc','e_cparhdr','e_minalloc','e_maxalloc',
                 'e_ss','e_sp','e_csum','e_ip','e_cs','e_lfarlc','e_ovno']
    for name,value in zip(field_names,fields):source.append(f'.word {hex(value)} # MZ {name}')
    header_size=fields[4]*16;reloc_at=fields[12];reloc_end=reloc_at+fields[3]*4
    source.extend(data_source(original[28:reloc_at]))
    for n in range(fields[3]):
        offset,segment=struct.unpack_from('<HH',original,reloc_at+n*4)
        source.append(f'.word {hex(offset)}, {hex(segment)} # MZ relocation {n}')
    source.extend(data_source(original[reloc_end:header_size]))
    cursor=header_size;mapping=[{'file_start':0,'file_end':header_size,'kind':'mz-header'}]
    for row in rows:
        at=row['file_offset'];size=len(bytes.fromhex(row['bytes']))
        if at>cursor:
            source.extend(data_source(original[cursor:at]))
            mapping.append({'file_start':cursor,'file_end':at,'kind':'data-or-unknown'})
        source.append(f'.org {hex(at)}')
        if row['original_name'] and re.fullmatch(r'[A-Za-z_][A-Za-z_0-9]*',row['original_name']):
            source.append(row['original_name']+':')
        source.append(f'# IDA {row["ida_linear"]:#x}; file {at:#x}; {row["assembly"]}')
        line_no=len(source)+1
        source.append(matched[at]);kind='instruction'
        mapping.append({'file_start':at,'file_end':at+size,'ida_linear':row['ida_linear'],
                        'original_name':row['original_name'],'function_name':row['function_name'],
                        'assembly_line':line_no,'kind':kind,
                        'original_assembly':row['assembly'],'operands':row['operands'],
                        'locator_level':'proven'})
        cursor=at+size
    source.extend(data_source(original[cursor:]));source.append(f'.org {hex(len(original))}')
    if cursor<len(original):mapping.append({'file_start':cursor,'file_end':len(original),'kind':'data-or-unknown'})
    expected_cursor=0
    for part in mapping:
        if part['file_start']!=expected_cursor:raise ValueError('source map contains a hole or overlap')
        expected_cursor=part['file_end']
    if expected_cursor!=len(original):raise ValueError('source map does not cover the complete executable')
    src=output/'KI.reconstructed.S';src.write_text('\n'.join(source)+'\n',encoding='utf-8')
    run(['as','--32','-o',str(output/'KI.o'),str(src)])
    linker_source=output/'KI.ld'
    linker_source.write_text('\n'.join(f'{name} = {hex(value)};' for name,value in symbols.items())+
                             '\nSECTIONS { .image 0 : { *(.image) } }\n',encoding='utf-8')
    run(['ld','-m','elf_i386','-T',str(linker_source),'--entry','0','-o',str(output/'KI.elf'),str(output/'KI.o')])
    run(['objcopy','-O','binary','--only-section=.image',str(output/'KI.elf'),str(output/'KI.EXE')])
    rebuilt=(output/'KI.EXE').read_bytes()
    # .image 沒有 .text 的 linker 地址語意；資料與 branch 都須在最終檔再次比。
    exact=rebuilt==original
    report={'schema':'wolong-assembly-rebuild-v1','input_sha256':sha,'input_size':len(original),
            'rebuilt_sha256':hashlib.sha256(rebuilt).hexdigest(),'rebuilt_size':len(rebuilt),
            'whole_file_exact':exact,'instruction_count':len(rows),'matched_instruction_count':len(matched),
            'matched_instruction_bytes':sum(len(bytes.fromhex(r['bytes'])) for r in rows if r['file_offset'] in matched),
            'unrecovered_instructions':[{'file_offset':at,'reason':reason} for at,reason in sorted(unmatched.items())],
            'classification_spans':inventory['classification_spans'],'trials':trials,
            'mz_relocation_count':fields[3],'masked_bytes':0,
            'source_map_covers_whole_file':True,
            'ida_version':inventory['tool_version'],'ida_image_id':inventory['image_id'],
            'build_image_id':image_id,
            'database_sha256':hashlib.sha256((inventory_path.parent/'input.exe.i64').read_bytes()).hexdigest(),
            'inventory_sha256':hashlib.sha256(inventory_path.read_bytes()).hexdigest(),
            'source_sha256':hashlib.sha256(src.read_bytes()).hexdigest(),
            'linker_source_sha256':hashlib.sha256(linker_source.read_bytes()).hexdigest(),
            'imm16_linker_symbols':symbols,
            'assembler_version':run(['as','--version']).splitlines()[0],
            'c_reconstruction_status':'not-started','semantic_completeness':'unknown',
            'negative_control_one_byte_mutation_rejected':bytes([rebuilt[0]^1])+rebuilt[1:]!=original}
    (output/'report.json').write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    (output/'source-map.json').write_text(json.dumps(mapping,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    source_by_ea={r['ida_linear']:r for r in mapping if 'ida_linear' in r}
    ledger=[]
    for function in inventory['inventory']:
        chunks=function['chunks']
        members=[r for r in rows if any(start<=r['ida_linear']<end for start,end in chunks)]
        entry=source_by_ea.get(function['ida_linear'])
        semantics=next((t['semantic'] for t in inventory['targets'] if t['ida_linear']==function['ida_linear']),
                       {'meaning':None,'level':'unknown','sources':[],'warning':'⚠ 語意尚未證實'})
        ledger.append({'original_name':function['name'],'ida_linear':function['ida_linear'],
                       'chunks':chunks,'boundary_authority':'IDA navigation; not original source boundary proof',
                       'assembly_line':entry['assembly_line'] if entry else None,
                       'decoded_instruction_bytes':sum(len(bytes.fromhex(r['bytes'])) for r in members),
                       'instruction_bytes_exact':all(r['file_offset'] in matched for r in members),
                       'semantic':semantics,'c_status':'not-started'})
    (output/'function-ledger.json').write_text(json.dumps(ledger,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    # 從來源檔與 linker script 冷重建；這一階段不讀 original 或 inventory。
    standalone=output/'standalone';standalone.mkdir(exist_ok=True)
    run(['as','--32','-o',str(standalone/'KI.o'),str(src)])
    run(['ld','-m','elf_i386','-T',str(linker_source),'--entry','0','-o',str(standalone/'KI.elf'),str(standalone/'KI.o')])
    run(['objcopy','-O','binary','--only-section=.image',str(standalone/'KI.elf'),str(standalone/'KI.EXE')])
    report['standalone_source_exact']=hashlib.sha256((standalone/'KI.EXE').read_bytes()).hexdigest()==EXPECTED
    controls=[]
    for kind in ['instruction','linker-constant','mz-relocation']:
        mutated_source=source.copy();mutated_linker=linker_source.read_text(encoding='utf-8')
        if kind=='instruction':
            idx=next(i for i,line in enumerate(mutated_source) if line=='clc')
            mutated_source[idx]='stc'
        elif kind=='linker-constant':
            name=next(iter(symbols));mutated_linker=mutated_linker.replace(f'{name} = {hex(symbols[name])};',f'{name} = {hex(symbols[name]^1)};')
        else:
            idx=next(i for i,line in enumerate(mutated_source) if line.endswith('# MZ relocation 0'))
            offset,segment=struct.unpack_from('<HH',original,reloc_at)
            mutated_source[idx]=f'.word {hex(offset^1)}, {hex(segment)} # changed MZ relocation'
        msrc=standalone/f'{kind}.S';mld=standalone/f'{kind}.ld'
        msrc.write_text('\n'.join(mutated_source)+'\n',encoding='utf-8');mld.write_text(mutated_linker,encoding='utf-8')
        run(['as','--32','-o',str(standalone/'mutant.o'),str(msrc)])
        run(['ld','-m','elf_i386','-T',str(mld),'--entry','0','-o',str(standalone/'mutant.elf'),str(standalone/'mutant.o')])
        run(['objcopy','-O','binary','--only-section=.image',str(standalone/'mutant.elf'),str(standalone/f'{kind}.EXE')])
        controls.append({'kind':kind,'rejected':hashlib.sha256((standalone/f'{kind}.EXE').read_bytes()).hexdigest()!=EXPECTED})
    report['negative_controls']=controls
    report['function_ledger_count_navigation_only']=len(ledger)
    (output/'report.json').write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    print('whole file exact',exact,'instructions',len(matched),'/',len(rows),'unrecovered',len(unmatched),flush=True)
    if not exact or not report['standalone_source_exact'] or not all(c['rejected'] for c in controls):
        raise ValueError('reconstruction or independent controls failed')


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--original',type=Path,required=True);p.add_argument('--inventory',type=Path,required=True)
    p.add_argument('--output',type=Path,required=True)
    p.add_argument('--image-id',required=True)
    a=p.parse_args();rebuild(a.original,a.inventory,a.output,a.image_id)
