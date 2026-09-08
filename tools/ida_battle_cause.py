"""對唯讀來源的資料庫副本匯出戰況根因候選；不修改名稱或註記。"""
import json
import ida_auto
import ida_bytes
import ida_funcs
import ida_kernwin
import ida_nalt
import ida_pro
import idautils
import idc

ida_auto.auto_wait()
result = {'tool': 'IDA Pro ' + ida_kernwin.get_kernel_version(),
          'address_space': 'IDA DOS/V linear',
          'input': ida_nalt.get_input_file_path(),
          'input_sha256': ida_nalt.retrieve_input_file_sha256().hex(),
          'function_count': ida_funcs.get_func_qty(), 'functions': []}
for ea in (0x14E5C, 0x15130, 0x15285, 0x152D7, 0x1A6FA, 0x1ADC8, 0x1A7B7, 0x1B240,
           0x1AF69, 0x1A398, 0x1A34F, 0x1AE56, 0x1A754, 0x1A785, 0x1A1C5,
           0x1A591, 0x1A5E1, 0x1A5F7, 0x1474A):
    f = ida_funcs.get_func(ea)
    assert f is not None, hex(ea)
    item = {'address': hex(ea), 'original_name': ida_funcs.get_func_name(f.start_ea),
            'start': hex(f.start_ea), 'end': hex(f.end_ea),
            'xrefs': [{'from': hex(x.frm), 'type': x.type} for x in idautils.XrefsTo(ea)],
            'instructions': []}
    for p in idautils.FuncItems(f.start_ea):
        item['instructions'].append({'address': hex(p),
            'bytes': ida_bytes.get_bytes(p, idc.get_item_size(p)).hex(),
            'operand': idc.GetDisasm(p), '附加語意': '未附加；根因候選待動態查證',
            '推論等級': '未知', '證據': '本輸入資料庫原始指令；不以候選標籤推定語意'})
    result['functions'].append(item)
with open('/work/battle-cause.json', 'w', encoding='utf-8') as f:
    json.dump(result, f, ensure_ascii=False, indent=2)
ida_pro.qexit(0)
