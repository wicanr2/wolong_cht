#!/usr/bin/env python3
"""從 PC-98 FDI 磁片映像抽出檔案。

FDI = 4096 byte 檔頭 + 裸磁片映像。臥竜伝的五片都是 2HD 1232 KB
（1024 byte/sector × 1232 sectors），檔案系統是 FAT12，只有開機磁區的
OEM 欄位被換成自製的 `FGDOS1.0`。

只用標準函式庫，不裝任何套件。
"""
import os
import struct
import sys

FDI_HEADER = 4096


def read_image(path):
    with open(path, 'rb') as fp:
        return fp.read()[FDI_HEADER:]


def parse_bpb(img):
    return {
        'bps': struct.unpack('<H', img[0x0b:0x0d])[0],
        'spc': img[0x0d],
        'reserved': struct.unpack('<H', img[0x0e:0x10])[0],
        'nfat': img[0x10],
        'root_entries': struct.unpack('<H', img[0x11:0x13])[0],
        'total': struct.unpack('<H', img[0x13:0x15])[0],
        'media': img[0x15],
        'spf': struct.unpack('<H', img[0x16:0x18])[0],
        'oem': img[0x03:0x0b].decode('ascii', 'replace'),
    }


def fat12_chain(img, bpb, start):
    """跟著 FAT12 的簇鏈走，回傳簇編號串列。"""
    fat_off = bpb['reserved'] * bpb['bps']
    fat = img[fat_off:fat_off + bpb['spf'] * bpb['bps']]
    chain, cur = [], start
    while 2 <= cur < 0xff0:
        chain.append(cur)
        idx = cur + cur // 2
        pair = struct.unpack('<H', fat[idx:idx + 2])[0]
        cur = pair >> 4 if cur & 1 else pair & 0x0fff
        if len(chain) > 4096:                      # 壞鏈的保險絲
            raise RuntimeError('cluster chain too long')
    return chain


def root_dir(img, bpb):
    off = (bpb['reserved'] + bpb['nfat'] * bpb['spf']) * bpb['bps']
    for i in range(bpb['root_entries']):
        e = img[off + i * 32: off + i * 32 + 32]
        if e[0] == 0:
            break
        if e[0] == 0xe5 or e[11] & 0x08:           # 已刪除 / 磁碟標籤
            continue
        name = e[0:8].decode('shift_jis', 'replace').rstrip()
        ext = e[8:11].decode('shift_jis', 'replace').rstrip()
        yield {
            'name': f'{name}.{ext}' if ext else name,
            'attr': e[11],
            'cluster': struct.unpack('<H', e[26:28])[0],
            'size': struct.unpack('<I', e[28:32])[0],
        }


def extract(path, outdir):
    img = read_image(path)
    bpb = parse_bpb(img)
    data_off = (bpb['reserved'] + bpb['nfat'] * bpb['spf']) * bpb['bps'] \
        + bpb['root_entries'] * 32
    os.makedirs(outdir, exist_ok=True)
    print(f'{os.path.basename(path)}  oem={bpb["oem"]!r} '
          f'bps={bpb["bps"]} total={bpb["total"]}')
    for ent in root_dir(img, bpb):
        if ent['attr'] & 0x10:                     # 子目錄，這五片都沒有
            continue
        blob = b''
        for clus in fat12_chain(img, bpb, ent['cluster']):
            beg = data_off + (clus - 2) * bpb['spc'] * bpb['bps']
            blob += img[beg:beg + bpb['spc'] * bpb['bps']]
        blob = blob[:ent['size']]
        with open(os.path.join(outdir, ent['name']), 'wb') as fp:
            fp.write(blob)
        flag = '' if len(blob) == ent['size'] else '  ← 長度不足'
        print(f'   {ent["name"]:14s} {ent["size"]:8d}{flag}')


if __name__ == '__main__':
    if len(sys.argv) != 3:
        sys.exit('用法: fdi_extract.py <image.fdi> <輸出目錄>')
    extract(sys.argv[1], sys.argv[2])
