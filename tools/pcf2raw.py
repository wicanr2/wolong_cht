#!/usr/bin/env python3
"""把 X11 的 PCF 點陣字型轉成本專案用的區位 raw 格式（docs/spec/84）。

    tools/py.sh tools/pcf2raw.py <in.pcf[.gz]> <out> [--charset jis|gb]

用途：簡體要 `HZK16`（GB2312）、日文要 `JISKAN16`（JIS X 0208），
這兩種 DOS 時代的字型檔不一定拿得到，但等價的 X11 點陣字型很常見
（Debian 的 `xfonts-intl-*`：`jiskan16.pcf.gz`、`gb16st.pcf.gz`）。
轉出來的檔放進 `-font` 目錄就能用。

**產物不進版控**（與倚天、HZK16 同一個政策：字型是別人的資產）。

輸出格式：`(區−1)×94 ＋ (位−1)` 為索引，每字 32 byte、
16 列 × 每列 2 byte、MSB-first——與 `internal/assets/cjk` 的區位載入器一致。
"""

import argparse
import gzip
import struct
import sys
from pathlib import Path

PCF_BITMAPS = 8
PCF_BDF_ENCODINGS = 32
GLYPH_H, ROW_BYTES = 16, 2
GLYPH_STRIDE = GLYPH_H * ROW_BYTES
ZONE = 94


def read(path):
    data = Path(path).read_bytes()
    if data[:2] == b"\x1f\x8b":
        data = gzip.decompress(data)
    if data[:4] != b"\x01fcp":
        sys.exit(f"{path} 不是 PCF（開頭 {data[:4]!r}）")
    return data


class Table:
    """PCF 的一張表。format 的 bit 2 決定表內整數的位元組順序。"""

    def __init__(self, blob, offset, size):
        self.blob = blob
        self.base = offset
        self.size = size
        self.fmt = struct.unpack_from("<I", blob, offset)[0]
        self.big = bool((self.fmt >> 2) & 1)
        self.msb_bits = bool((self.fmt >> 3) & 1)
        self.pos = offset + 4

    def i32(self):
        v = struct.unpack_from(">i" if self.big else "<i", self.blob, self.pos)[0]
        self.pos += 4
        return v

    def i16(self):
        v = struct.unpack_from(">h" if self.big else "<h", self.blob, self.pos)[0]
        self.pos += 2
        return v

    def u8(self):
        v = self.blob[self.pos]
        self.pos += 1
        return v


def tables(blob):
    count = struct.unpack_from("<i", blob, 4)[0]
    out = {}
    for i in range(count):
        t, fmt, size, off = struct.unpack_from("<4i", blob, 8 + i * 16)
        out[t] = (off, size)
    return out


REVERSE = bytes(int(f"{b:08b}"[::-1], 2) for b in range(256))


def convert(src, dst, charset):
    blob = read(src)
    tabs = tables(blob)
    for need, name in ((PCF_BITMAPS, "PCF_BITMAPS"),
                       (PCF_BDF_ENCODINGS, "PCF_BDF_ENCODINGS")):
        if need not in tabs:
            sys.exit(f"{src} 缺 {name}")

    bm = Table(blob, *tabs[PCF_BITMAPS])
    nglyph = bm.i32()
    offsets = [bm.i32() for _ in range(nglyph)]
    sizes = [bm.i32() for _ in range(4)]
    data_start = bm.pos
    total = sizes[bm.fmt & 3]

    enc = Table(blob, *tabs[PCF_BDF_ENCODINGS])
    min2, max2, min1, max1, _default = (enc.i16(), enc.i16(), enc.i16(),
                                        enc.i16(), enc.i16())
    if max1 <= 0:
        sys.exit("這是單位元組字型，不是區位字集")
    # 區位的偏移：GB 系列的第一位元組是 區+0xA0，JIS 系列是 區+0x20。
    bias = 0xA0 if min1 >= 0xA1 else 0x20

    out = bytearray()
    written = 0
    for b1 in range(min1, max1 + 1):
        for b2 in range(min2, max2 + 1):
            gi = enc.i16()
            if gi < 0 or gi >= nglyph:
                continue
            ku, ten = b1 - bias, b2 - bias
            if not (1 <= ku <= ZONE and 1 <= ten <= ZONE):
                continue
            start = data_start + offsets[gi]
            end = data_start + (offsets[gi + 1] if gi + 1 < nglyph else total)
            raw = blob[start:end]
            # ⚠ **每一列都補齊到 glyph pad**（format 的低 2 bit：1／2／4／8
            # byte）。16 寬的字用 int 對齊就是一列 4 byte、一個字 64 byte，
            # 照 32 byte 硬讀會得到「前 8 列各取前半」——畫面上是有字形但
            # 全部糊掉，比缺字難認得多。列寬要從實際的每字長度反推。
            if len(raw) % GLYPH_H:
                continue
            stride = len(raw) // GLYPH_H
            if stride < ROW_BYTES:
                continue
            if not bm.msb_bits:
                raw = bytes(REVERSE[b] for b in raw)
            unit = 1 << ((bm.fmt >> 4) & 3)
            if not bm.big and unit > 1:
                # scan unit 以 LSB 存放時，單位內的 byte 是反的。
                raw = b"".join(raw[i:i + unit][::-1]
                               for i in range(0, len(raw), unit))
            glyph = b"".join(raw[r * stride:r * stride + ROW_BYTES]
                             for r in range(GLYPH_H))
            idx = (ku - 1) * ZONE + (ten - 1)
            need = (idx + 1) * GLYPH_STRIDE
            if len(out) < need:
                out.extend(b"\x00" * (need - len(out)))
            out[idx * GLYPH_STRIDE:need] = glyph
            written += 1

    Path(dst).write_bytes(bytes(out))
    probe = "あ" if charset == "jis" else "啊"
    print(f"{dst}：{written} 個字、{len(out)} byte")
    print(f"  自我檢查交給 internal/assets/cjk 的載入器（探針字「{probe}」）")


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("src")
    ap.add_argument("dst")
    ap.add_argument("--charset", choices=["jis", "gb"], default="jis")
    a = ap.parse_args()
    convert(a.src, a.dst, a.charset)
