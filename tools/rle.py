#!/usr/bin/env python3
"""原版的 RLE 解壓（`MMAP.MAP` 用的那一種）。

出處：`KI.EXE` 的 `sub_1F5E7`（`MMAP.MAP` 專用載入器，
與其他檔走的 `sub_1F4A2` 不同支）。

演算法——**用「連續兩個相同的 byte」當 run 的觸發**，沒有逃脫字元：

    逐 byte 複製；
    一旦輸出的 byte 與前一個相同，下一個輸入 byte 是「再重複幾次」；
    次數 0 表示那兩個相同的 byte 就只是字面值，回到逐 byte 模式。

用法：

    rle.py decode <輸入> <輸出> [預期長度]

只用標準函式庫。
"""
import sys


def decode(src, expect=None):
    """回傳解壓後的 bytes。expect 給了就順便檢查長度。"""
    out = bytearray()
    i = 0
    n = len(src)
    while i < n:
        prev = src[i]
        out.append(prev)
        i += 1
        # 逐 byte 複製，直到出現與前一個相同的 byte
        while i < n:
            cur = src[i]
            out.append(cur)
            i += 1
            if cur == prev:
                break
            prev = cur
        else:
            break
        if i >= n:
            break
        count = src[i]
        i += 1
        # 次數 0 → 那兩個相同的 byte 只是字面值
        out.extend(bytes([prev]) * count)
    if expect is not None and len(out) != expect:
        print(f'⚠ 解出 {len(out)} B，預期 {expect} B（差 {len(out)-expect}）',
              file=sys.stderr)
    return bytes(out)


if __name__ == '__main__':
    if len(sys.argv) < 4 or sys.argv[1] != 'decode':
        sys.exit(__doc__)
    raw = open(sys.argv[2], 'rb').read()
    expect = int(sys.argv[4]) if len(sys.argv) > 4 else None
    out = decode(raw, expect)
    open(sys.argv[3], 'wb').write(out)
    print(f'{sys.argv[2]}: {len(raw)} B → {len(out)} B')
