#!/usr/bin/env python3
# IDAPython：把指定的位址區間逐條反組譯出來。
#
# 為什麼需要它：`tools/ida_dump.py` 是以**函式**為單位，而 IDA 沒認出函式的
# 區段（`get_func()` 回 None）就 dump 不到——熱區登記表裡那四個
# 連續呼叫正好落在這種地方。
#
# 區間寫在 RANGES，不必為了換位址改工具。
#
# ⚠ headless 的 print 不進 stdout，一律寫檔；輸出帶 probe（輸入檔 SHA-256）。
# 用法：tools/ida.sh script dosv tools/ida_range.py KI.EXE.i64
# 輸出 /work/range.txt
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import idc

OUT = "/work/range.txt"
RANGES = [
    (0x10300, 0x10480, "聲軌狀態 0A02 的三個使用點（播放引擎候選）"),
    (0x10270, 0x102A0, "099A 的使用點"),
    (0x10540, 0x10660, "099A 的另兩個使用點"),
]


def main():
    ida_auto.auto_wait()
    lines = ["位址區間反組譯（IDA DOS/V linear address）",
             "輸入檔 SHA-256：%s" % ida_nalt.retrieve_input_file_sha256().hex(),
             "函式數：%d" % ida_funcs.get_func_qty()]
    for start, end, note in RANGES:
        lines.append("")
        lines.append("==== %08X–%08X  %s ====" % (start, end, note))
        ea = start
        while ea != idc.BADADDR and ea < end:
            f = ida_funcs.get_func(ea)
            owner = ida_funcs.get_func_name(f.start_ea) if f else "（無函式）"
            lines.append("  %08X  %-10s  %s" % (ea, owner, idc.GetDisasm(ea)))
            ea = idc.next_head(ea, end)
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(lines) + "\n")
    ida_pro.qexit(0)


main()
