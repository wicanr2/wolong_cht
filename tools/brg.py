#!/usr/bin/env python3
"""`.BRG` 調色盤解碼與色票輸出。

格式見 `docs/formats/02-brg-palette.md`，出處是 `docs/re/02` 的反組譯：

    每色 3 byte，順序 B, R, G（副檔名就是通道順序）
    每個通道 4 bit，值域 0–15
    每 16 色為一組（bank），原版一次只切一組進硬體

亮度縮放（原版的淡入淡出）：

    out = ((v << 4) * brightness + 0x80) >> 8      brightness 0–16

PC-98 直接把 out 當 4 bit 類比調色盤值寫進 0AAh/0ACh/0AEh；
DOS/V 再左移 2 位變成 VGA DAC 的 6 bit。**兩者的顏色是同一個**，
只是硬體位寬不同——所以轉成 8 bit sRGB 時要用同一條式子。

用法：

    brg.py info  <檔.BRG>                    印出每組 16 色的數值
    brg.py swatch <檔.BRG> <out.png> [每格px]  輸出色票 PNG

只用標準函式庫，不裝任何套件。
"""
import os
import struct
import sys
import zlib

BANK = 16                                   # 一組 16 色（`cmp ch, 10h`）
FULL = 16                                   # 亮度上限（`mov cl, 10h`）


def scale(value, brightness=FULL):
    """原版的亮度縮放，回傳 0–15 的通道值。"""
    return (((value << 4) * brightness) + 0x80) >> 8


def to_srgb(value, brightness=FULL):
    """4 bit 通道 → 8 bit sRGB。15 → 255。"""
    v = scale(value, brightness)
    return v * 0xff // 0x0f


def load(path):
    """回傳 [(r, g, b), …]，已經是 0–255。"""
    blob = open(path, 'rb').read()
    if len(blob) % 3:
        raise ValueError(f'{path}: {len(blob)} B 不是 3 的倍數')
    out = []
    for i in range(0, len(blob), 3):
        b, r, g = blob[i], blob[i + 1], blob[i + 2]     # ← BRG，不是 RGB
        out.append((to_srgb(r), to_srgb(g), to_srgb(b)))
    return out


def _png(path, width, height, rows):
    """最小 PNG 輸出（truecolor 8-bit），不依賴任何套件。"""
    raw = b''.join(b'\x00' + bytes(px for pixel in row for px in pixel)
                   for row in rows)

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


def cmd_info(path):
    colors = load(path)
    banks = len(colors) // BANK
    print(f'{os.path.basename(path)}: {len(colors)} 色 = {banks} 組 × {BANK}')
    blob = open(path, 'rb').read()
    for bank in range(banks):
        print(f'  bank {bank}:')
        for i in range(BANK):
            idx = bank * BANK + i
            b, r, g = blob[idx * 3], blob[idx * 3 + 1], blob[idx * 3 + 2]
            rr, gg, bb = colors[idx]
            print(f'    {i:2d}  檔案 B={b:2d} R={r:2d} G={g:2d}'
                  f'   → #{rr:02x}{gg:02x}{bb:02x}')


def cmd_swatch(path, out, cell=32):
    cell = int(cell)
    colors = load(path)
    banks = len(colors) // BANK
    width, height = BANK * cell, banks * cell
    rows = []
    for y in range(height):
        bank = y // cell
        rows.append([colors[bank * BANK + x // cell] for x in range(width)])
    _png(out, width, height, rows)
    print(f'{out}: {width}×{height}，{banks} 組 × {BANK} 色')


if __name__ == '__main__':
    if len(sys.argv) < 3:
        sys.exit(__doc__)
    fn = {'info': cmd_info, 'swatch': cmd_swatch}.get(sys.argv[1])
    if not fn:
        sys.exit(__doc__)
    sys.exit(fn(*sys.argv[2:]) or 0)
