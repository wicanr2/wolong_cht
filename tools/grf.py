#!/usr/bin/env python3
"""`*GRF.DAT` 圖庫解碼。

格式見 `docs/formats/03-grf-images.md`，出處是 `docs/re/03` 的反組譯：

    4 bpp planar，**plane-major**——plane0 整張、plane1 整張、plane2、plane3
    每平面每列 width/8 byte，列跨距在 VRAM 裡是 80 byte（螢幕 640 寬）
    像素值 = plane0 的 bit ＋ plane1<<1 ＋ plane2<<2 ＋ plane3<<3

`KAOGRF.DAT` 的尺寸是 confirmed 的：載入器用「編號 × 2048」定位、
`di=800h` 一次讀 2,048 B，繪製常式跑 64 列 × 每列 8 byte × 4 平面。

用法：

    grf.py sheet <檔.DAT> <寬> <高> <調色盤.BRG> <組號> <out.png> [每列張數]
    grf.py one   <檔.DAT> <寬> <高> <調色盤.BRG> <組號> <索引> <out.png> [放大]

只用標準函式庫，不裝任何套件。
"""
import os
import struct
import sys
import zlib

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import brg as brgmod                                        # noqa: E402


def frame_bytes(width, height):
    """一張圖佔幾個 byte。"""
    return width * height // 2


def decode(blob, offset, width, height):
    """回傳 height × width 的調色盤索引二維串列。"""
    stride = width // 8                                     # 每平面每列的 byte 數
    plane = stride * height
    rows = []
    for y in range(height):
        row = []
        for x in range(width):
            byte_i = y * stride + x // 8
            bit = 7 - (x & 7)                               # 高位在左
            value = 0
            for p in range(4):
                b = blob[offset + p * plane + byte_i]
                value |= ((b >> bit) & 1) << p
            row.append(value)
        rows.append(row)
    return rows


def _png(path, width, height, rows):
    raw = b''.join(b'\x00' + bytes(c for px in row for c in px) for row in rows)

    def chunk(tag, data):
        body = tag + data
        return (struct.pack('>I', len(data)) + body
                + struct.pack('>I', zlib.crc32(body) & 0xffffffff))

    with open(path, 'wb') as fp:
        fp.write(b'\x89PNG\r\n\x1a\n')
        fp.write(chunk(b'IHDR',
                       struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0)))
        fp.write(chunk(b'IDAT', zlib.compress(raw, 9)))
        fp.write(chunk(b'IEND', b''))


def _palette(path, bank):
    colors = brgmod.load(path)
    bank = int(bank)
    return colors[bank * brgmod.BANK:(bank + 1) * brgmod.BANK]


def cmd_sheet(src, width, height, palpath, bank, out, per_row=15):
    width, height, per_row = int(width), int(height), int(per_row)
    blob = open(src, 'rb').read()
    size = frame_bytes(width, height)
    count = len(blob) // size
    pal = _palette(palpath, bank)
    cols = min(per_row, count)
    sheet_rows = (count + cols - 1) // cols
    canvas = [[(0, 0, 0)] * (cols * width) for _ in range(sheet_rows * height)]
    for i in range(count):
        cx, cy = (i % cols) * width, (i // cols) * height
        for y, row in enumerate(decode(blob, i * size, width, height)):
            for x, idx in enumerate(row):
                canvas[cy + y][cx + x] = pal[idx]
    _png(out, cols * width, sheet_rows * height, canvas)
    print(f'{out}: {count} 張 {width}×{height}，'
          f'排成 {cols}×{sheet_rows}（餘 {len(blob) % size} B）')


def cmd_region(src, offset, width, height, palpath, bank, out, per_row=16):
    """解組合檔裡的一段。`ICONGRF.DAT` 是四段拼起來的，每段尺寸不同。

    offset 可以寫 0x2800 或十進位。count 由段長推算——
    段長是載入器的 `di` 值，不是檔案長度。
    """
    width, height, per_row = int(width), int(height), int(per_row)
    offset = int(str(offset), 0)
    blob = open(src, 'rb').read()[offset:]
    size = frame_bytes(width, height)
    count = len(blob) // size
    pal = _palette(palpath, bank)
    cols = min(per_row, count)
    sheet_rows = (count + cols - 1) // cols
    canvas = [[(0, 0, 0)] * (cols * width) for _ in range(sheet_rows * height)]
    for i in range(count):
        cx, cy = (i % cols) * width, (i // cols) * height
        for y, row in enumerate(decode(blob, i * size, width, height)):
            for x, idx in enumerate(row):
                canvas[cy + y][cx + x] = pal[idx]
    _png(out, cols * width, sheet_rows * height, canvas)
    print(f'{out}: 位移 0x{offset:X} 起 {count} 張 {width}×{height}'
          f'（餘 {len(blob) % size} B）')


def cmd_one(src, width, height, palpath, bank, index, out, zoom=4):
    width, height, index, zoom = int(width), int(height), int(index), int(zoom)
    blob = open(src, 'rb').read()
    pal = _palette(palpath, bank)
    grid = decode(blob, index * frame_bytes(width, height), width, height)
    canvas = [[pal[grid[y // zoom][x // zoom]] for x in range(width * zoom)]
              for y in range(height * zoom)]
    _png(out, width * zoom, height * zoom, canvas)
    print(f'{out}: #{index} {width}×{height} ×{zoom}')


if __name__ == '__main__':
    if len(sys.argv) < 3:
        sys.exit(__doc__)
    fn = {'sheet': cmd_sheet, 'region': cmd_region, 'one': cmd_one}.get(sys.argv[1])
    if not fn:
        sys.exit(__doc__)
    sys.exit(fn(*sys.argv[2:]) or 0)
