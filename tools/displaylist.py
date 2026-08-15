#!/usr/bin/env python3
"""把 `KI.EXE` 裡的視窗顯示清單攤開（規格：`docs/re/48-window-display-list.md`）。

    tools/py.sh tools/displaylist.py workplace/orig/dosv/KI.EXE
    tools/py.sh tools/displaylist.py workplace/orig/dosv/KI.EXE --scene 0

## 這是什麼

`sub_10337(al = 場景編號)` 是一支直譯器，不是一次 blit。清單在 `cs:0E16`
（檔案位移 `0x1016`），每筆 12 bytes ＝ 六個 word：

    op  X  Y  arg1  arg2  0

第一個 byte 為 0 的記錄是**場景分隔**，`al` 決定跳過幾個。
`X`／`Y` 相對於呼叫端傳進去的原點。

## 為什麼要一支工具而不是逐場景讀組語

同一份清單餵 11 個視窗建構常式。要解「某個視窗的內部排版」時，
把它的場景印出來就好——**逐場景讀組語是把已經是資料的東西再推一次**。

## 判讀

opcode 的用途見 `docs/re/48` §2。`08` 的 `arg1` 是字串位址，
這支工具會直接把 cp950 解出來；`09` 的 `arg1` 是圖庫位移、
`arg2` 是「高<<8 ｜ 寬」。
"""
import argparse
import struct
import sys

LIST_CS = 0x0E16          # 清單在程式段的位移（sub_1030F 的 ax）
EXE_HEADER = 0x200        # cs:XXXX → 檔案位移 ＝ XXXX + 0x200
REC = 12
MAX_SCENES = 64
# opcode 實際只用到 0x03–0x09。設 0x0F 是留餘裕，同時擋住把 Big5 文字
# 當成記錄讀進來（那些 word 動輒 0xA1A1 以上）。
MAX_OP = 0x0F

OPS = {
    0x03: "填矩形",
    0x06: "未解",
    0x07: "場景範圍",
    0x08: "字串",
    0x09: "貼圖",
}


def cstr(data, cs_off, limit=40):
    """把 cs 位移上的字串取出來（cp950，遇 0 結束）。"""
    off = cs_off + EXE_HEADER
    if off < 0 or off + limit > len(data):
        return None
    raw = data[off:off + limit].split(b"\x00")[0]
    try:
        return raw.decode("cp950")
    except UnicodeDecodeError:
        return repr(raw)


def scenes(data):
    """切出所有場景。回傳 [(場景編號, [記錄…])]，記錄是六個 word 的 tuple。"""
    off = LIST_CS + EXE_HEADER
    out, cur, started = [], [], False
    while len(out) < MAX_SCENES and off + REC <= len(data):
        rec = struct.unpack("<6H", data[off:off + REC])
        off += REC
        if rec[0] == 0:
            # 分隔記錄。第一筆之前沒有內容，不要產出一個空場景。
            if started:
                out.append(cur)
            cur, started = [], True
            continue
        if not started:
            # 清單開頭理論上就是分隔記錄；不是的話別硬解。
            raise SystemExit("清單開頭不是分隔記錄，格式假設不成立")
        # ⚠ 清單**沒有結束標記**，後面接的是 Big5 字串資料。
        # 沒有這道閘會多解出一個「場景」，內容是文字被當成記錄讀的亂碼——
        # 而那看起來就像一個還沒解開的場景，不像讀過頭。
        if rec[0] > MAX_OP:
            break
        cur.append(rec)
    return list(enumerate(out))


def describe(data, rec):
    op, x, y, a1, a2, _ = rec
    name = OPS.get(op, "未知 op %02X" % op)
    if op == 0x08:
        return "%-6s (%3d,%3d)  字串 cs:%04X ＝ %r  屬性 %04X" % (
            name, x, y, a1, cstr(data, a1), a2)
    if op == 0x09:
        # 尺寸的編碼與 `sub_1E3D7` 的 cx 同一套：高 byte ＝ 高、低 byte ＝ 寬。
        # 場景 0 的四筆列距 16 px 正好等於解出來的高，是這個讀法的檢查。
        return "%-6s (%3d,%3d)  圖庫位移 %04X  %d×%d（寬×高）" % (
            name, x, y, a1, a2 & 0xFF, a2 >> 8)
    if op in (0x03, 0x07):
        return "%-6s (%3d,%3d)–(%3d,%3d)  ＝ %d×%d" % (
            name, x, y, a1, a2, a1 - x + 1, a2 - y + 1)
    return "%-6s (%3d,%3d)  arg %04X %04X" % (name, x, y, a1, a2)


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("exe")
    ap.add_argument("--scene", type=int, help="只印這一個場景")
    ns = ap.parse_args()
    data = open(ns.exe, "rb").read()

    all_scenes = scenes(data)
    print("清單在 cs:%04X（檔案位移 0x%04X），共 %d 個場景\n"
          % (LIST_CS, LIST_CS + EXE_HEADER, len(all_scenes)))
    for i, recs in all_scenes:
        if ns.scene is not None and i != ns.scene:
            continue
        print("==== 場景 %d（%d 筆）====" % (i, len(recs)))
        for r in recs:
            print("  " + describe(data, r))
        print()
    return 0


if __name__ == "__main__":
    sys.exit(main())
