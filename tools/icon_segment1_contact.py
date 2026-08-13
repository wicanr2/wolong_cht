#!/usr/bin/env python3
"""Render DOS/V ICONGRF segment-1 candidates without assigning UI meaning.

All offsets in the labels are relative to ICONGRF segment 1, the segment loaded
by sub_1C7A9 from file offset 0x2800.  This is a research visualizer, not a
runtime decoder: it only applies the planar stride/height given by the
corresponding IDA blit call and never wires a candidate into the game.
"""
from __future__ import annotations

import argparse
import struct
import zlib
from pathlib import Path


FONT = {
    "0": ("11111", "10001", "10011", "10101", "11001", "10001", "11111"),
    "1": ("00100", "01100", "00100", "00100", "00100", "00100", "01110"),
    "2": ("11110", "00001", "00001", "00110", "01000", "10000", "11111"),
    "3": ("11110", "00001", "00001", "01110", "00001", "00001", "11110"),
    "4": ("00010", "00110", "01010", "10010", "11111", "00010", "00010"),
    "5": ("11111", "10000", "10000", "11110", "00001", "00001", "11110"),
    "6": ("00110", "01000", "10000", "11110", "10001", "10001", "01110"),
    "7": ("11111", "00001", "00010", "00100", "01000", "01000", "01000"),
    "8": ("01110", "10001", "10001", "01110", "10001", "10001", "01110"),
    "9": ("01110", "10001", "10001", "01111", "00001", "00010", "11100"),
    "A": ("01110", "10001", "10001", "11111", "10001", "10001", "10001"),
    "B": ("11110", "10001", "10001", "11110", "10001", "10001", "11110"),
    "C": ("01111", "10000", "10000", "10000", "10000", "10000", "01111"),
    "D": ("11110", "10001", "10001", "10001", "10001", "10001", "11110"),
    "E": ("11111", "10000", "10000", "11110", "10000", "10000", "11111"),
    "F": ("11111", "10000", "10000", "11110", "10000", "10000", "10000"),
    "x": ("00000", "10001", "01010", "00100", "01010", "10001", "00000"),
    "×": ("00000", "10001", "01010", "00100", "01010", "10001", "00000"),
    " ": ("00000",) * 7,
}


def png(path: Path, width: int, height: int, pixels: list[list[tuple[int, int, int]]]) -> None:
    def chunk(tag: bytes, data: bytes) -> bytes:
        body = tag + data
        return struct.pack(">I", len(data)) + body + struct.pack(">I", zlib.crc32(body) & 0xFFFFFFFF)

    raw = b"".join(b"\0" + bytes(c for px in row for c in px) for row in pixels)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("wb") as f:
        f.write(b"\x89PNG\r\n\x1a\n")
        f.write(chunk(b"IHDR", struct.pack(">IIBBBBB", width, height, 8, 2, 0, 0, 0)))
        f.write(chunk(b"IDAT", zlib.compress(raw, 9)))
        f.write(chunk(b"IEND", b""))


def load_palette(path: Path, bank: int) -> list[tuple[int, int, int]]:
    raw = path.read_bytes()
    if len(raw) % 3 or len(raw) < (bank + 1) * 48:
        raise ValueError(f"palette length {len(raw)} cannot provide bank {bank}")
    out = []
    base = bank * 48
    for i in range(16):
        # BRG is stored as B, R, G with 4-bit channel values.
        b, r, g = raw[base + i * 3 : base + i * 3 + 3]
        out.append((r * 16, g * 16, b * 16))
    return out


def decode_planar(seg: bytes, offset: int, width: int, height: int) -> list[list[int]]:
    if width % 8:
        raise ValueError("planar candidate width must be a multiple of 8")
    stride = width // 8
    plane = stride * height
    size = plane * 4
    if offset < 0 or offset + size > len(seg):
        raise ValueError(f"candidate 0x{offset:04X} needs 0x{size:X} bytes")
    rows: list[list[int]] = []
    for y in range(height):
        row = []
        for x in range(width):
            byte_at = offset + y * stride + x // 8
            bit = 7 - (x & 7)
            value = 0
            for plane_no in range(4):
                value |= ((seg[byte_at + plane_no * plane] >> bit) & 1) << plane_no
            row.append(value)
        rows.append(row)
    return rows


def scale(rows: list[list[int]], factor: int, palette: list[tuple[int, int, int]]) -> list[list[tuple[int, int, int]]]:
    out = []
    for row in rows:
        expanded = [palette[v] for v in row for _ in range(factor)]
        for _ in range(factor):
            out.append(expanded[:])
    return out


def text_pixels(text: str, fg: tuple[int, int, int], bg: tuple[int, int, int]) -> list[list[tuple[int, int, int]]]:
    width = len(text) * 6
    out = [[bg] * width for _ in range(7)]
    for n, ch in enumerate(text):
        glyph = FONT.get(ch, FONT[" "])
        for y, line in enumerate(glyph):
            for x, bit in enumerate(line):
                if bit == "1":
                    out[y][n * 6 + x] = fg
    return out


def paste(dst, src, x: int, y: int) -> None:
    for yy, row in enumerate(src):
        if 0 <= y + yy < len(dst):
            for xx, pixel in enumerate(row):
                if 0 <= x + xx < len(dst[0]):
                    dst[y + yy][x + xx] = pixel


def panel(rows, label: str, palette, factor: int, pad: int = 8):
    body = scale(rows, factor, palette)
    label_px = text_pixels(label, (255, 255, 255), (0, 0, 0))
    width = max(len(body[0]), len(label_px[0])) + pad * 2
    height = len(body) + len(label_px) + pad * 2
    canvas = [[(0, 0, 0)] * width for _ in range(height)]
    paste(canvas, label_px, pad, pad)
    paste(canvas, body, pad, pad + len(label_px) + pad // 2)
    return canvas


def join_panels(panels, cols: int, gap: int = 12):
    rows = (len(panels) + cols - 1) // cols
    widths = [0] * cols
    heights = [0] * rows
    for i, p in enumerate(panels):
        c, r = i % cols, i // cols
        widths[c] = max(widths[c], len(p[0]))
        heights[r] = max(heights[r], len(p))
    width = sum(widths) + gap * (cols - 1)
    height = sum(heights) + gap * (rows - 1)
    canvas = [[(18, 18, 18)] * width for _ in range(height)]
    y = 0
    for r in range(rows):
        x = 0
        for c in range(cols):
            i = r * cols + c
            if i < len(panels):
                paste(canvas, panels[i], x, y)
            x += widths[c] + gap
        y += heights[r] + gap
    return canvas


def render_sheet(seg: bytes, palette, specs, out: Path, factor: int, cols: int) -> None:
    panels = []
    for offset, width, height in specs:
        rows = decode_planar(seg, offset, width, height)
        panels.append(panel(rows, f"0x{offset:04X}", palette, factor))
    canvas = join_panels(panels, cols)
    png(out, len(canvas[0]), len(canvas), canvas)


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("icon", type=Path)
    ap.add_argument("palette", type=Path)
    ap.add_argument("outdir", type=Path)
    ap.add_argument("--bank", type=int, default=0)
    args = ap.parse_args()
    raw = args.icon.read_bytes()
    seg = raw[0x2800 : 0x2800 + 0x3F00]
    palette = load_palette(args.palette, args.bank)

    exact = [
        (0x0000, 128, 32),
        (0x0800, 128, 32),
        (0x1000, 128, 32),
        (0x1800, 128, 96),
        (0x3000, 80, 32),
        (0x3500, 128, 16),
        (0x3900, 24, 16),
        (0x3D80, 16, 8),
        (0x3DC0, 16, 8),
        (0x3E00, 16, 16),
        (0x3E80, 16, 16),
    ]
    render_sheet(seg, palette, exact, args.outdir / "icongrf-seg1-blit-candidates.png", 4, 2)

    # sub_1F888 不保留 SI；每平面消費 3 bytes × 16 列，四平面
    # 共前進 0xC0。sub_1C7F4 連呼叫六次，所以這裡展開為六張
    # 24×16，而不是把 0x3900 的第一張重複六次。
    commands = [(0x3900 + i * 0xC0, 24, 16) for i in range(6)]
    render_sheet(seg, palette, commands, args.outdir / "icongrf-seg1-battle-command-glyphs.png", 8, 3)

    tail8 = [(off, 16, 8) for off in range(0x3D80, 0x3F00, 0x40)]
    render_sheet(seg, palette, tail8, args.outdir / "icongrf-seg1-tail-16x8-windows.png", 8, 3)

    tail16 = [(off, 16, 16) for off in (0x3D80, 0x3E00, 0x3E80)]
    render_sheet(seg, palette, tail16, args.outdir / "icongrf-seg1-tail-16x16-candidates.png", 6, 3)
    print("generated 4 DOS/V ICONGRF segment-1 candidate sheets")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
