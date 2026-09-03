#!/usr/bin/env python3
"""反解一個顯示格：（欄, 列）在每一層深度分別對到哪一格戰場。

    tools/py.sh tools/slot_cells.py <戰場編號> <欄> <列> [camX camY]

`internal/ui/isoview` 的 `cellOffset` 是正向的（戰場格 → 顯示格）：

    欄 ＝ (x + y + 1) − (camX + camY)
    列 ＝ floorDiv2((y + 1 − x) − (camY − camX)) − z

逐區對拍定位到某一塊像素之後，要問的下一句是「那一格該畫什麼」，
而這需要反過來走。同一個顯示格在不同深度對到**不同的格子**——
所以輸出是一張七列的表，不是一個答案。

配 `tools/subtile_match.py` 使用：那一支給出子圖塊編號，
這一支給出「哪一格、哪一層、原本的圖塊是什麼」。

⚠ 玩家守城時原版把戰場轉 180 度（`docs/spec/56` §3），所以兩種框都印。
"""
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
ORIG = ROOT / "workplace" / "orig" / "dosv"

INDEX_SIZE, FIELD_SIZE = 512, 4096
CELLS_OFF, WIDTH, HEIGHT = 0x40, 64, 62
MDL_HEADER, TILESET_SIZE, TILE_DEF = 4096, 63488, 8
MAX_STACK = 7

# `sub_199F3`：word_1D328 ＝ 0x24、word_1D32A ＝ 0x0E。
CAM_INIT_X, CAM_INIT_Y = 0x24, 0x0E


def floordiv2(v):
    """原版的 `sar bx, 1`——負數往下取整。"""
    return -((-v + 1) // 2) if v < 0 else v // 2


def rotate_tile(v):
    """`sub_1CBBC` 的值對映（`internal/assets/battle.RotateTile`）。"""
    if v < 0x30:
        return v
    if v <= 0xCF:
        return ((v - 0x30) ^ 0x10) + 0x30
    if v <= 0xEF:
        return v ^ 3 if (v & 3) in (0, 3) else v
    return v ^ 1


def load(field):
    """回傳 (未翻轉的格子, 翻轉的格子, 每個圖塊的子圖塊堆疊)。"""
    bmap = (ORIG / "BATTLE.MAP").read_bytes()
    mdl = (ORIG / "BATTLE.MDL").read_bytes()
    off = INDEX_SIZE + field * FIELD_SIZE + CELLS_OFF
    flat = bmap[off:off + WIDTH * HEIGHT]
    plain = [list(flat[y * WIDTH:(y + 1) * WIDTH]) for y in range(HEIGHT)]
    rot = [[rotate_tile(plain[HEIGHT - 1 - y][WIDTH - 1 - x])
            for x in range(WIDTH)] for y in range(HEIGHT)]
    base = MDL_HEADER + bmap[field * 2] * TILESET_SIZE
    stacks = []
    for i in range(256):
        r = mdl[base + i * TILE_DEF: base + (i + 1) * TILE_DEF]
        stacks.append(list(r[1:1 + min(r[0], MAX_STACK)]))
    return plain, rot, stacks


def cell_for(col, row, z, camx, camy):
    """哪一格戰場在深度 z 落在這個顯示格。沒有就回 None。"""
    for y in range(HEIGHT):
        x = (camx + camy) + col - (y + 1)
        if not 0 <= x < WIDTH:
            continue
        if floordiv2(((y + 1) - x) - (camy - camx)) - z == row:
            return x, y
    return None


def main():
    field = int(sys.argv[1])
    col, row = int(sys.argv[2]), int(sys.argv[3])
    camx, camy = (int(sys.argv[4]), int(sys.argv[5])) if len(sys.argv) > 5 \
        else (CAM_INIT_X, CAM_INIT_Y)

    plain, rot, stacks = load(field)
    print(f"戰場 {field}　鏡頭 ({camx},{camy})　顯示格（欄 {col}, 列 {row}）")
    for label, cells in (("未翻轉", plain), ("翻轉 180°", rot)):
        print(f"\n== {label} ==")
        for z in range(MAX_STACK):
            at = cell_for(col, row, z, camx, camy)
            if at is None:
                print(f"  z={z}：沒有格子對到")
                continue
            x, y = at
            tile = cells[y][x]
            st = stacks[tile]
            top = st[z] if z < len(st) else None
            print(f"  z={z}　戰場格 ({x},{y})　圖塊 0x{tile:02X}　"
                  f"堆疊 {st}　這一層 {top}")


if __name__ == "__main__":
    main()
