#!/usr/bin/env python3
"""原版的 RLE 解壓（`MMAP.MAP` 用的那一種）。

出處：`KI.EXE` 的 `sub_1F5E7`（`MMAP.MAP` 專用載入器，
與其他檔走的 `sub_1F4A2` 不同支）。

演算法——**用「連續兩個相同的 byte」當 run 的觸發**，沒有逃脫字元：

    逐 byte 複製；
    一旦輸出的 byte 與前一個相同，下一個輸入 byte 是「再重複幾次」；
    次數 0 表示那兩個相同的 byte 就只是字面值，回到逐 byte 模式。

檔案前面還有 **4 byte 小端 u32 ＝ 解壓後的長度**，原版的載入器
`LSEEK` 跳過它才開始解（docs/spec/113）。整個檔要走 `decode_file`；
`decode` 只解裸資料。

用法：

    rle.py decode      <輸入> <輸出> [預期長度]   只解裸資料
    rle.py decode-file <輸入> <輸出>              跳過長度頭並核對

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


HEADER = 4                                  # 檔頭那 4 byte ＝ 小端 u32 解壓長度


def decode_file(src):
    """解一整個 RLE 資料檔：跳過 4 byte 長度頭，解剩下的，核對長度。

    原版的載入器在讀第一個 byte 之前先 LSEEK 到位移 4
    （`KI.EXE` 的 `sub_1F655`、`D7OPEN.EXE` 的 `sub_10E04`、`D7END.EXE`
    同一段——docs/re/76 §5、docs/spec/113）。

    ⭐ **長度是驗收條件不是參考值**：從 0 開始解會在某一處掉相位，
    差幾十個 byte、畫面整體位移。`MMAP.MAP` 是唯一躲過的檔。
    """
    if len(src) < HEADER:
        raise ValueError(f'只有 {len(src)} B，連 {HEADER} B 的長度頭都不夠')
    want = int.from_bytes(src[:HEADER], 'little')
    out = decode(src[HEADER:])
    if len(out) != want:
        raise ValueError(f'檔頭宣告 {want} B，解出來 {len(out)} B——'
                         '差一個 byte 就是解錯，不是尾巴沒編進去')
    return out


if __name__ == '__main__':
    if len(sys.argv) < 4 or sys.argv[1] not in ('decode', 'decode-file'):
        sys.exit(__doc__)
    raw = open(sys.argv[2], 'rb').read()
    if sys.argv[1] == 'decode-file':
        out = decode_file(raw)
    else:
        expect = int(sys.argv[4]) if len(sys.argv) > 4 else None
        out = decode(raw, expect)
    open(sys.argv[3], 'wb').write(out)
    print(f'{sys.argv[2]}: {len(raw)} B → {len(out)} B')
