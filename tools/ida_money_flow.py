#!/usr/bin/env python3
"""列出資金與「本月支出」那四支帳務常式的每一個呼叫端。

    tools/ida.sh script dosv tools/ida_money_flow.py KI.EXE.i64
    輸出 /work/money-flow.txt（＝ workplace/ida/<版本>/census/money-flow.txt）

## 為什麼要這一支

錢的流向散在整份程式裡，而**四支收斂點是 24 位飽和加減**：

    sub_15609  資金 +=（+0x20/+0x22，飽和 0x09FE98）
    sub_1563B  資金 −=（下限 0xF60168）
    sub_15673  本月支出 +=（+0x1A/+0x1C，月結時扣掉再歸零）
    0x15663    上面那支的包裝：dh ＝ 勢力編號、ax:dl ＝ 金額

**問「誰花了錢」等於問「誰呼叫這四支」**，而那是交叉參考圖答得出來的，
grep `.asm` 答不出來（`0x15663` 連函式名都沒有）。
"""
import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import idautils
import idc

OUT = "/work/money-flow.txt"
TARGETS = [
    (0x15609, "資金 +="),
    (0x1562B, "資金 −=（包裝：dh ＝ 勢力）"),
    (0x1563B, "資金 −="),
    (0x15663, "本月支出 +=（包裝：dh ＝ 勢力）"),
    (0x15673, "本月支出 +="),
]


def owner(ea):
    f = ida_funcs.get_func(ea)
    if f is None:
        return "（無函式）"
    return idc.get_func_name(f.start_ea)


def main():
    ida_auto.auto_wait()
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("帳務常式的呼叫端（IDA DOS/V linear address）\n")
        fh.write("輸入檔 SHA-256：%s\n" % ida_nalt.retrieve_input_file_sha256().hex())
        fh.write("函式數：%d\n\n" % ida_funcs.get_func_qty())
        for ea, label in TARGETS:
            fh.write("==== %05X  %s ====\n" % (ea, label))
            n = 0
            for x in idautils.XrefsTo(ea):
                # fl_CN=16／fl_CF=17 才是呼叫邊；fl_F=21 是循序落下。
                if x.type not in (16, 17):
                    continue
                n += 1
                fh.write("  %08X  %-14s %s\n" % (x.frm, owner(x.frm), idc.GetDisasm(x.frm)))
            fh.write("  呼叫點共 %d 個\n\n" % n)
    ida_pro.qexit(0)


main()
