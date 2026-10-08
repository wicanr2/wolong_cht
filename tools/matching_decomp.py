#!/usr/bin/env python3
"""KI.EXE 的局部 matching 試點。必須在有 GNU binutils/GCC 的 Docker 內執行。

組語以指令重建，不複製原始 code bytes；C 探針只測試現代 GCC 的可匹配性。
局部結果不代表整檔、原始編譯器或 Go remake 已通過。
"""
import argparse
import hashlib
import json
import re
import struct
import subprocess
from pathlib import Path

EXPECTED_SHA256 = "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868"
TARGETS = {"sub_11D8E": (0x1D8E, 137), "sub_131AE": (0x31AE, 68),
           "sub_1ECE0": (0xECE0, 28)}

# {load} 選擇 x86 同語意指令的 load-direction 編碼。
# .org 只固定段內位置；不把補零算入已還原範圍。
ASSEMBLY = """.intel_syntax noprefix
.code16
.text
.equ sub_15358, 0x5358
.equ sub_19377, 0x9377
.equ sub_13E11, 0x3e11
.equ sub_11E17, 0x1e17
.org 0x1d8e
sub_11D8E:
    cmp byte ptr [0x0cf2], 8
    jb loc_11DF4
    cmp byte ptr [0x0cf3], 0x17
    jb loc_11DE0
    mov ax, word ptr [0x0cf0]
    {load} cmp al, ah
    jb loc_11DD7
    cmp byte ptr [0x0cf4], 12
    jb loc_11DC1
    cmp word ptr [0x0cf6], 1000
    jb loc_11DB8
    mov word ptr [0x0cf6], 998
loc_11DB8:
    inc word ptr [0x0cf6]
    mov byte ptr [0x0cf4], 0
loc_11DC1:
    inc byte ptr [0x0cf4]
    mov bl, byte ptr [0x0cf4]
    {load} xor bh, bh
    mov ah, byte ptr [bx-0x6755]
    {load} xor al, al
    mov word ptr [0x0cf0], ax
    call sub_15358
loc_11DD7:
    inc byte ptr [0x0cf0]
    mov byte ptr [0x0cf3], 0
loc_11DE0:
    inc byte ptr [0x0cf3]
    mov byte ptr [0x0cf2], 0
    call sub_19377
    call sub_13E11
    call sub_11E17
    jmp loc_11DF8
loc_11DF4:
    inc byte ptr [0x0cf2]
loc_11DF8:
    mov al, byte ptr [0x0cfa]
    {load} and al, al
    jz locret_11E16
loc_11DFF:
    cmp byte ptr [0x0d2c], 0
    jz loc_11DFF
loc_11E06:
    cmp byte ptr [0x0d2d], al
    jb loc_11E06
    mov byte ptr [0x0d2d], 0
    mov byte ptr [0x0d2c], 0
locret_11E16:
    ret
.org 0x31ae
sub_131AE:
    dec byte ptr [0x31ad]
    jz loc_131B5
    ret
loc_131B5:
    cmp word ptr [0x0d20], 0x100
    jb loc_131BE
    ret
loc_131BE:
    mov byte ptr [0x31ad], 10
    mov es, word ptr [0x0d56]
    mov bx, word ptr [0x0d20]
    mov ax, word ptr es:[bx]
    mov dx, word ptr es:[bx+2]
    add bx, 4
    mov word ptr [0x0d20], bx
    {load} and al, al
    jz locret_131F1
    mov ds, word ptr [0x0d52]
    {load} mov bl, al
    {load} xor bh, bh
    dec bx
    shl bx, 1
    call word ptr cs:[bx+0x31f2]
    mov ax, cs
    mov ds, ax
locret_131F1:
    ret
.org 0xece0
sub_1ECE0:
    push ds
    push bx
    mov ax, cs
    mov ds, ax
    mov bx, 0xecfe
    mov al, byte ptr [0xecfd]
    xlatb
    add al, byte ptr [0xecfc]
    add byte ptr [0xecfc], 0x89
    mov byte ptr [0xecfd], al
    pop bx
    pop ds
    ret
"""

# 這些是已知規則的 C 探針，不是已匹配或可執行的 DOS 程式。
# GCC -m16 使用 32-bit ABI，無法以這些宣告表達原版 DS/ES 與旗標契約。
C_COMMON = """typedef unsigned char u8;
typedef unsigned short u16;
#define B(a) (*(volatile u8 *)(a))
#define W(a) (*(volatile u16 *)(a))
"""
C_PROBES = {
    "sub_1ECE0": """u8 sub_1ECE0(void) {
    u8 value = B(0xecfe + B(0xecfd)) + B(0xecfc);
    B(0xecfc) += 0x89;
    B(0xecfd) = value;
    return value;
}
""",
    "sub_11D8E": """extern void sub_15358(void), sub_19377(void), sub_13E11(void), sub_11E17(void);
void sub_11D8E(void) {
    if (B(0xcf2) < 8) { ++B(0xcf2); goto pace; }
    if (B(0xcf3) < 23) goto hour;
    if (B(0xcf0) < B(0xcf1)) goto day;
    if (B(0xcf4) < 12) goto month;
    if (W(0xcf6) >= 1000) W(0xcf6) = 998;
    ++W(0xcf6); B(0xcf4) = 0;
month:
    ++B(0xcf4); W(0xcf0) = (u16)B(0x98ab + B(0xcf4)) << 8;
    sub_15358();
day:
    ++B(0xcf0); B(0xcf3) = 0;
hour:
    ++B(0xcf3); B(0xcf2) = 0;
    sub_19377(); sub_13E11(); sub_11E17();
pace:
    { u8 speed = B(0xcfa); if (!speed) return;
      while (!B(0xd2c)) {} while (B(0xd2d) < speed) {}
      B(0xd2d) = 0; B(0xd2c) = 0; }
}
""",
    "sub_131AE": """/* 分派器的平面記憶體模型；沒有重建 DS/ES 或 handler 的暫存器 ABI。 */
typedef void (*handler)(u16 event, u16 parameter);
void sub_131AE(void) {
    if (--B(0x31ad)) return;
    if (W(0xd20) >= 0x100) return;
    B(0x31ad) = 10;
    u16 cursor = W(0xd20);
    unsigned queue = (unsigned)W(0xd56) << 4;
    u16 event = W(queue + cursor), parameter = W(queue + cursor + 2);
    W(0xd20) = cursor + 4;
    if ((u8)event) {
        ((handler *)(0x31f2))[(u8)event - 1](event, parameter);
    }
}
""",
}


def run(command, *, cwd=None):
    result = subprocess.run(command, cwd=cwd, capture_output=True, text=True, timeout=45)
    if result.returncode:
        raise RuntimeError(f"{command[0]} failed ({result.returncode}): {result.stderr}")
    return result.stdout


def compare(expected, actual):
    positions = [i for i in range(max(len(expected), len(actual)))
                 if expected[i:i+1] != actual[i:i+1]]
    return {"exact": not positions, "expected_size": len(expected),
            "actual_size": len(actual), "different_bytes": len(positions),
            "first_difference_relative": positions[0] if positions else None,
            "original_sha256": hashlib.sha256(expected).hexdigest(),
            "rebuilt_sha256": hashlib.sha256(actual).hexdigest(),
            "relocations_masked": 0}


def validate_probe(probe, original):
    if probe.get('schema') != 'wolong-matching-ida-probe-v1':
        raise ValueError('unexpected IDA probe schema')
    sha = hashlib.sha256(original).hexdigest()
    if sha != EXPECTED_SHA256 or probe.get('input_sha256') != sha:
        raise ValueError('input SHA-256 is not the selected original')
    if not probe.get('semantic_index_applied'):
        raise ValueError('graded semantic index was not applied to the selected original')
    found = {f['name']: f for f in probe['targets']}
    for name, (offset, length) in TARGETS.items():
        f = found[name]
        if len(f['chunks']) != 1:
            raise ValueError(f'{name}: multiple chunks require explicit coverage')
        c = f['chunks'][0]
        if c['start'] != 0x10000+offset or c['end']-c['start'] != length:
            raise ValueError(f'{name}: function boundaries differ')
        at = c['instructions'][0]['file_offset']
        if original[at:at+length] != bytes.fromhex(c['bytes']):
            raise ValueError(f'{name}: IDA bytes differ from original file')
        cursor = 0
        for i in c['instructions']:
            data = bytes.fromhex(i['bytes'])
            if i['ida_linear'] != c['start']+cursor or i['file_offset'] != at+cursor:
                raise ValueError(f'{name}: discontinuous IDA instruction coverage')
            if original[at+cursor:at+cursor+len(data)] != data:
                raise ValueError(f'{name}: instruction bytes differ')
            cursor += len(data)
        if cursor != length:
            raise ValueError(f'{name}: incomplete IDA instruction coverage')
    return found


def build_and_compare(original_path, probe_path, output, image_id):
    original = original_path.read_bytes()
    probe = json.loads(probe_path.read_text())
    targets = validate_probe(probe, original)
    database_sha = hashlib.sha256((probe_path.parent/'input.exe.i64').read_bytes()).hexdigest()
    if database_sha != (probe_path.parent/'database.sha256').read_text().split()[0]:
        raise ValueError('IDA database hash sidecar differs from actual database')
    control_path = probe_path.parent.parent/'borland-control/ida-probe.json'
    control = json.loads(control_path.read_text())
    control_raw = (original_path.parent/'LOGO.EXE').read_bytes()
    if control['input_sha256'] != hashlib.sha256(control_raw).hexdigest():
        raise ValueError('Borland control input hash differs')
    markers = control['printable_toolchain_markers']
    bp_count = sum(f['bp_frame_prefix'] for f in control['inventory'])
    if not any('Borland C++' in m['text'] for m in markers) or not bp_count:
        raise ValueError('Borland control did not detect known positive anchors')
    output.mkdir(parents=True, exist_ok=True)
    if re.search(r'^\s*\.(byte|word|short|long|incbin)\b', ASSEMBLY, re.M):
        raise ValueError('assembly baseline may not inject raw code bytes')
    source = output/'reconstructed.S'
    source.write_text(ASSEMBLY, encoding='utf-8')
    run(['as','--32','-o',str(output/'reconstructed.o'),str(source)])
    run(['ld','-m','elf_i386','-Ttext','0','--entry','0','-o',
         str(output/'reconstructed.elf'),str(output/'reconstructed.o')])
    relocation_text = run(['readelf','-r',str(output/'reconstructed.elf')])
    if not 'There are no relocations' in relocation_text:
        raise ValueError('assembly result contains unresolved relocations')
    run(['objcopy','-O','binary','--only-section=.text',
         str(output/'reconstructed.elf'),str(output/'reconstructed.bin')])
    rebuilt = (output/'reconstructed.bin').read_bytes()
    changed_call_source = output/'changed-call.S'
    changed_call_source.write_text(ASSEMBLY.replace('sub_13E11, 0x3e11', 'sub_13E11, 0x3e12'),
                                   encoding='utf-8')
    run(['as','--32','-o',str(output/'changed-call.o'),str(changed_call_source)])
    run(['ld','-m','elf_i386','-Ttext','0','--entry','0','-o',
         str(output/'changed-call.elf'),str(output/'changed-call.o')])
    run(['objcopy','-O','binary','--only-section=.text',
         str(output/'changed-call.elf'),str(output/'changed-call.bin')])
    changed_call = (output/'changed-call.bin').read_bytes()
    clock_at = targets['sub_11D8E']['chunks'][0]['instructions'][0]['file_offset']
    changed_call_rejected = not compare(original[clock_at:clock_at+137],
                                        changed_call[0x1d8e:0x1d8e+137])['exact']
    matches, controls, compiler_trials = [], [], []
    for name, (offset, length) in TARGETS.items():
        c = targets[name]['chunks'][0]
        at = c['instructions'][0]['file_offset']
        expected = original[at:at+length]
        actual = rebuilt[offset:offset+length]
        matches.append({'name':name,'ida_linear':0x10000+offset,
                        'segment_offset':offset,'file_offset':at,
                        'scope':'instruction-authored assembly reassembly',
                        **compare(expected,actual)})
        mutant = bytearray(actual); mutant[0] ^= 1
        controls.append({'name':name,'one_byte_mutation_rejected':
                         not compare(expected,bytes(mutant))['exact']})
        csource = output/f'{name}.c'
        csource.write_text(C_COMMON+C_PROBES[name],encoding='utf-8')
        for optimization in ['-O0','-O1','-O2','-Os']:
            prefix = output/f'{name}-{optimization[1:]}'
            command = ['gcc','-m16',optimization,'-ffreestanding','-fno-pic','-fno-pie',
                       '-fno-stack-protector','-fno-asynchronous-unwind-tables',
                       '-fno-unwind-tables','-c',str(csource),'-o',str(prefix)+'.o']
            run(command)
            link = ['ld','-m','elf_i386','-Ttext',hex(offset),'--entry',name,
                    '--defsym','sub_15358=0x5358','--defsym','sub_19377=0x9377',
                    '--defsym','sub_13E11=0x3e11','--defsym','sub_11E17=0x1e17',
                    '-o',str(prefix)+'.elf',str(prefix)+'.o']
            run(link)
            run(['objcopy','-O','binary','--only-section=.text',
                 str(prefix)+'.elf',str(prefix)+'.bin'])
            trial = {'name':name,'flags':command[1:-4],
                     'scope':'modern GCC flat-memory C probe; original ABI unrecovered',
                     **compare(expected,Path(str(prefix)+'.bin').read_bytes())}
            compiler_trials.append(trial)
    wrong_probe = dict(probe); wrong_probe['input_sha256'] = '0'*64
    try:
        validate_probe(wrong_probe,original)
    except ValueError:
        wrong_hash_rejected = True
    else:
        wrong_hash_rejected = False
    header = struct.unpack_from('<14H',original)
    report = {
        'schema':'wolong-matching-result-v1','input_name':'dosv/KI.EXE',
        'input_sha256':hashlib.sha256(original).hexdigest(),
        'input_size':len(original), 'mz_header_size':header[4]*16,
        'mz_relocation_count':header[3],
        'ida_version':probe['tool_version'],
        'ida_function_count_navigation_only':probe['function_count'],
        'bp_frame_prefix_count':sum(f['bp_frame_prefix'] for f in probe['inventory']),
        'library_flag_count_navigation_only':sum(f['library_flag'] for f in probe['inventory']),
        'build_image_id':image_id,
        'assembler_version':run(['as','--version']).splitlines()[0],
        'linker_version':run(['ld','--version']).splitlines()[0],
        'compiler_version':run(['gcc','--version']).splitlines()[0],
        'assembly_matches':matches, 'c_compiler_trials':compiler_trials,
        'matched_assembly_bytes':sum(m['expected_size'] for m in matches if m['exact']),
        'matched_high_level_c_functions':len({t['name'] for t in compiler_trials if t['exact']}),
        'negative_controls':controls,
        'wrong_input_hash_rejected':wrong_hash_rejected,
        'changed_external_call_rejected':changed_call_rejected,
        'compiler_probe_control':{
            'input_name':control['input_name'], 'input_sha256':control['input_sha256'],
            'function_count_navigation_only':control['function_count'],
            'bp_frame_prefix_count':bp_count, 'markers':markers,
            'compiler_id_navigation_only':control['compiler_id_navigation_only'],
            'probe_sha256':hashlib.sha256(control_path.read_bytes()).hexdigest(),
        },
        'whole_executable_match':False, 'go_remake_verified':False,
        'original_compiler':'unknown','original_source_language':'unknown',
        'source_sha256':hashlib.sha256(source.read_bytes()).hexdigest(),
        'ida_probe_sha256':hashlib.sha256(probe_path.read_bytes()).hexdigest(),
        'ida_image_id':probe['image_id'],
        'semantic_index_sha256':probe['semantic_index_sha256'],
        'database_sha256':database_sha,
    }
    (output/'report.json').write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    print(f"Assembly: {sum(m['exact'] for m in matches)}/3 functions, {report['matched_assembly_bytes']} bytes")
    print(f"C probes: {sum(t['exact'] for t in compiler_trials)}/12 matching trials")
    if not all(m['exact'] for m in matches):
        for m in matches:
            if not m['exact']:print(m)
        raise ValueError('assembly reconstruction does not match')
    if not wrong_hash_rejected or not changed_call_rejected or not all(x['one_byte_mutation_rejected'] for x in controls):
        raise ValueError('negative controls did not reject invalid inputs')


if __name__ == '__main__':
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--original',type=Path,required=True)
    parser.add_argument('--probe',type=Path,required=True)
    parser.add_argument('--output',type=Path,required=True)
    parser.add_argument('--image-id',required=True)
    args=parser.parse_args()
    build_and_compare(args.original,args.probe,args.output,args.image_id)
