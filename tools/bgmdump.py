#!/usr/bin/env python3
"""把 `*BGM.DAT` 的聲軌事件攤開（規格：`docs/re/56-bgm-track-events.md`）。

    tools/py.sh tools/bgmdump.py workplace/orig/dosv/BGM.DAT
    tools/py.sh tools/bgmdump.py workplace/orig/dosv/BGM.DAT --song 0 --events 40
    tools/py.sh tools/bgmdump.py workplace/orig/dosv/OPENBGM.DAT --ynsound workplace/orig/dosv/YNSOUND.COM

## 這是什麼

事件是 2 bytes 一筆，由 `YNSOUND.COM` 的 INT 8 播放引擎解譯
（`docs/re/56` §2）。**三張查表在 TSR 裡不在資料裡**，所以要一起讀：

    0x0AB0  長度表（高 byte & 0x7F → tick 數）
    0x0AD0  B0 值表（低 byte → block ＋ F-Number 高 2 位）
    0x0B50  A0 值表（低 4 位 → F-Number 低 8 位）

## 為什麼不把表寫死在工具裡

那三張表是**原版資料**，寫進工具等於把原版內容 commit 進版控
（`CLAUDE.md` §9）。而且兩版的 TSR 不見得一樣——讀進來才驗得出差異。
"""
import argparse
import struct
import sys

# 三張表在 YNSOUND.COM 裡的位置。COM offset − 0x100 ＝ 檔案 offset。
TBL_LEN, TBL_B0, TBL_A0 = 0x0AB0, 0x0AD0, 0x0B50
NOTE_NAMES = ["休", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"]

EMPTY_TRACK = 0x22  # 未使用的聲軌一律指到這個 stub

# INT 8 的頻率：PIT divisor 256（`docs/re/57` §5）。
PIT_HZ = 1193182.0 / 256


def control(lo, hi):
    """控制事件（低 byte ≥ 0x80）。分派看 bit 4–6（`docs/re/56` §5）。"""
    kind = (lo >> 4) & 7
    if kind == 0:
        return "音量%d" % hi
    if kind == 1:
        return "漸%s(量%d,每%d)" % ("弱" if lo & 1 else "強", hi >> 4, (hi & 0x0F) * 4)
    if kind == 2:
        return "音色%d" % hi
    if kind == 3:
        n = (0xFF - hi) * 11 // 8
        return "速度%d(%.1f tick/s)" % (hi, PIT_HZ / n) if n else "速度%d" % hi
    if kind == 4:
        return {0xC1: "迴圈回跳", 0xC2: "呼叫子段", 0xC3: "子段返回"}.get(lo, "跳回記號")
    if kind == 5:
        return {0xD1: "迴圈起點×%d" % hi, 0xD2: "子段入口"}.get(lo, "記號")
    if kind == 6:
        return "無作用[%02X %02X]" % (lo, hi)
    return "旗標%d" % hi


def load_tables(path):
    com = open(path, "rb").read()
    if len(com) < 0xB60:
        sys.exit("%s 太小（%d B），不像 YNSOUND.COM" % (path, len(com)))

    def at(com_off, n):
        return com[com_off - 0x100:com_off - 0x100 + n]

    return at(TBL_LEN, 32), at(TBL_B0, 32), at(TBL_A0, 16)


def songs(data):
    """回傳 [(offset, length)]。沒有索引的單曲檔回傳整個檔。"""
    if struct.unpack_from("<I", data, 0)[0] != 0x0100:
        return [(0, len(data))]
    out = []
    for i in range(11):
        off, length = struct.unpack_from("<II", data, i * 8)
        out.append((off, length))
    return out


def describe(block, start, limit, tbl_len, tbl_b0, tbl_a0):
    """把一條聲軌解成人看得懂的事件列。"""
    out, i = [], start
    while len(out) < limit and i + 2 <= len(block):
        lo, hi = block[i], block[i + 1]
        i += 2
        if lo >= 0x80:
            out.append(control(lo, hi))
            continue
        note, octave = lo & 0x0F, (lo >> 4) & 0x07
        ticks = tbl_len[hi & 0x7F] if (hi & 0x7F) < len(tbl_len) else -1
        tie = "~" if hi & 0x80 else ""
        if note == 0:
            out.append("休%s/%d" % (tie, ticks))
        elif note < len(NOTE_NAMES):
            out.append("%s%d%s/%d" % (NOTE_NAMES[note], octave, tie, ticks))
        else:
            out.append("?%02X%s/%d" % (lo, tie, ticks))
    return out


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("path")
    ap.add_argument("--ynsound", default="workplace/orig/dosv/YNSOUND.COM")
    ap.add_argument("--song", type=int, default=-1, help="只印第 N 曲")
    ap.add_argument("--events", type=int, default=16, help="每軌印幾個事件")
    args = ap.parse_args()

    tbl_len, tbl_b0, tbl_a0 = load_tables(args.ynsound)
    data = open(args.path, "rb").read()
    items = songs(data)
    print("%s：%d bytes，%d 曲" % (args.path, len(data), len(items)))

    for n, (off, length) in enumerate(items):
        if args.song >= 0 and n != args.song:
            continue
        block = data[off:off + length]
        ptrs = struct.unpack_from("<6H", block, 0x10)
        print("\n==== 第 %d 曲  offset 0x%04X 長 %d ====" % (n, off, length))
        print("  聲軌指標：" + " ".join("%04X" % p for p in ptrs))
        for ch, p in enumerate(ptrs):
            if p == EMPTY_TRACK:
                print("  聲軌 %d：空" % ch)
                continue
            events = describe(block, p, args.events, tbl_len, tbl_b0, tbl_a0)
            print("  聲軌 %d @%04X：%s" % (ch, p, "  ".join(events)))


main()
