#!/usr/bin/env python3
"""IDA：遷都與據點清單，保留 raw CFG、間接呼叫與自我修改運算元。"""
import hashlib,json,struct,sys,traceback
from pathlib import Path
sys.path.insert(0,'/tools')
import ida_matching_probe as probe
import ida_auto,ida_bytes,ida_funcs,ida_name,ida_nalt,ida_pro,ida_ua,idautils,idc

ROOTS=(0x15E60,0x16909,0x17400,0x1748F,0x181C0,0x1820E,0x18412,0x18458,0x18463,
       0x184BC,0x184DD,0x1851A,0x18546,0x18607,0x18662,0x186C9,0x18713,0x18755)
RAW_ROOTS=(0x1743B,0x1745F,0x1828F,0x1857F,0x185B2)
try:
 ida_auto.auto_wait()
 ledger=json.loads(Path('/evidence/c-recovery-status.json').read_text(encoding='utf-8'))
 known={x['ida_linear'] for x in list(ledger['functions'].values())+list(ledger['code_blocks'].values())}
 def calls(at):
  return sorted({x.to for a,b in idautils.Chunks(at) for ea in idautils.Heads(a,b)
                 for x in idautils.XrefsFrom(ea) if x.type in (16,17)})
 wanted=set(ROOTS);pending=list(ROOTS);raw_entries=set(RAW_ROOTS);raw_rows={}
 def add_target(t):
  if t in known or t in wanted or t==0x11234:return
  f=ida_funcs.get_func(t)
  if f and f.start_ea==t:wanted.add(t);pending.append(t)
  elif 0x10000<=t<0x20000:raw_entries.add(t)
 def walk(start):
  queue=[start];seen=set();rows=[]
  while queue:
   at=queue.pop()
   if at in seen:continue
   assert 0x10000<=at<0x20000 and len(seen)<512,hex(at)
   seen.add(at);original=idc.generate_disasm_line(at,0);iscode=ida_bytes.is_code(ida_bytes.get_full_flags(at))
   if not iscode:
    insn=ida_ua.insn_t();size=ida_ua.decode_insn(insn,at);assert size>0
    ida_bytes.del_items(at,ida_bytes.DELIT_SIMPLE,size);assert ida_ua.create_insn(at)==size
   row=probe.instruction(at);rows.append(row)
   if not iscode:row.update(original_ida_data_line=original,originally_code=False)
   mn=row['mnemonic'];ops=row['operands'];nextat=at+len(bytes.fromhex(row['bytes']))
   if mn=='call':
    for x in idautils.XrefsFrom(at):
     if x.type in (16,17):add_target(x.to)
   if mn in ['retn','retf','iret']:continue
   if mn=='jmp':
    assert ops[0]['type']==7,(hex(at),row['assembly']);queue.append(ops[0]['addr']+0x10000);continue
   if mn.startswith('j') or mn in ['loop','loope','loopne']:queue.append(ops[0]['addr']+0x10000)
   queue.append(nextat)
  return sorted(rows,key=lambda x:x['ida_linear'])
 done=set()
 while pending or raw_entries-done:
  while pending:
   at=pending.pop();f=ida_funcs.get_func(at);assert f and f.start_ea==at,hex(at)
   for t in calls(at):add_target(t)
  for at in sorted(raw_entries-done):raw_rows[at]=walk(at);done.add(at)
  assert len(wanted)<96 and len(raw_rows)<48
 probe.TARGETS=tuple(sorted(wanted));probe.main()
 p=Path('/output/ida-probe.json');r=json.loads(p.read_text(encoding='utf-8'));raw=Path(ida_nalt.get_input_file_path()).read_bytes()
 assert r['input_sha256']==ledger['input_sha256']
 header=struct.unpack_from('<H',raw,8)[0]*16;count=struct.unpack_from('<H',raw,6)[0];table=struct.unpack_from('<H',raw,24)[0]
 reloc=[header+off+seg*16 for off,seg in (struct.unpack_from('<HH',raw,table+i*4) for i in range(count))]
 r.update(mz_header_size=header,relocation_file_offsets=reloc,ida_load_paragraph=0x1000)
 for t in r['targets']:
  parts=[]
  for c in t['chunks']:
   a=c['start']-0x10000+header;data=raw[a:a+c['end']-c['start']];expected=bytearray(data);rr=[]
   for at in reloc:
    if not a<=at<a+len(data):continue
    assert at+2<=a+len(data);old=struct.unpack_from('<H',raw,at)[0];shift=(old+0x1000)&65535
    struct.pack_into('<H',expected,at-a,shift);rr.append({'file_offset':at,'original_word':old,'ida_word':shift})
   assert expected==bytes.fromhex(c['bytes']);c.update(file_bytes=data.hex(),loader_relocations=rr);parts.append(data)
  t['file_sha256']=hashlib.sha256(b''.join(parts)).hexdigest();t['direct_callees']=calls(t['ida_linear'])
 r['decoded_blocks']=[]
 for at,rows in sorted(raw_rows.items()):
  for row in rows:
   a=row['ida_linear']-0x10000+header;assert raw[a:a+len(bytes.fromhex(row['bytes']))]==bytes.fromhex(row['bytes'])
  r['decoded_blocks'].append({'ida_linear':at,'original_name_before_analysis':ida_name.get_name(at),'instructions':rows,
      'file_sha256':hashlib.sha256(b''.join(bytes.fromhex(x['bytes']) for x in rows)).hexdigest(),
      'xrefs_to':[{'from':x.frm,'type':x.type} for x in idautils.XrefsTo(at)]})
 r['list_descriptors']={hex(a):{'bytes':ida_bytes.get_bytes(a,z-a).hex(),
      'heads':[{'ida_linear':ea,'code':ida_bytes.is_code(ida_bytes.get_full_flags(ea)),'assembly':idc.generate_disasm_line(ea,0)} for ea in idautils.Heads(a,z)]}
      for a,z in [(0x17378,0x17400),(0x181A6,0x181C0),(0x183F0,0x18463),(0x1857F,0x18607)]}
 r['closure_scope']='Original relocation caller and city list callbacks, viewport, selection, scrollbar and sort; raw CFG entries retain original addresses and live call operands.'
 r['recovery_ledger_sha256']=hashlib.sha256(Path('/evidence/c-recovery-status.json').read_bytes()).hexdigest()
 p.write_text(json.dumps(r,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
except Exception:
 Path('/output/ida-probe-error.txt').write_text(traceback.format_exc(),encoding='utf-8');ida_pro.qexit(1)
else:ida_pro.qexit(0)
