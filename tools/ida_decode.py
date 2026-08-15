#!/usr/bin/env python3
# IDAPython：把指定區間**強制當成程式碼**解碼出來。
#
# 為什麼需要它：`YNSOUND.COM` 這種 TSR，IDA 的自動分析只認得出 5 支函式，
# 其餘整片是 `db`。`tools/ida_range.py` 對資料項只印得出一行 db——
# 看起來就像「那裡沒有程式碼」，而實際上播放引擎就在裡面。
#
# ⚠ 這支會**改資料庫**（del_items + create_insn）。`tools/ida.sh script`
# 是對副本跑的，原始 `.i64` 的雜湊不變——要引用資料庫身分的分析仍然安全。
#
# 區間寫在 /work/decode_list.txt，一行一筆「起 迄 說明」（十六進位，
# IDA linear address）：
#
#   printf '1032F 10400 播放引擎\n' > workplace/ida/dosv/census/decode_list.txt
#   tools/ida.sh script dosv tools/ida_decode.py YNSOUND.COM.i64
#
# 輸出 /work/decode.txt；第一行是 probe（輸入檔 SHA-256）。
import ida_auto
import ida_bytes
import ida_funcs
import ida_nalt
import ida_pro
import idc

LIST = "/work/decode_list.txt"
OUT = "/work/decode.txt"


def main():
    ida_auto.auto_wait()
    try:
        with open(LIST, encoding="utf-8") as fh:
            rows = [ln.split("#")[0].split() for ln in fh]
    except OSError:
        rows = []
    rows = [r for r in rows if len(r) >= 2]

    lines = ["強制解碼（IDA linear address）",
             "輸入檔 SHA-256：%s" % ida_nalt.retrieve_input_file_sha256().hex(),
             "函式數：%d" % ida_funcs.get_func_qty(),
             "區間數：%d" % len(rows)]
    for row in rows:
        start, end = int(row[0], 16), int(row[1], 16)
        label = " ".join(row[2:]) if len(row) > 2 else ""
        lines.append("")
        lines.append("==== %08X–%08X  %s ====" % (start, end, label))
        # 先把整段拆成未定義，否則舊的資料項會擋住 create_insn。
        ida_bytes.del_items(start, ida_bytes.DELIT_SIMPLE, end - start)
        ea = start
        while ea < end:
            n = idc.create_insn(ea)
            if n <= 0:
                # 解不動就當一個 byte 印出來再往前——**不要跳過**，
                # 跳過會讓後面整段錯位而看起來像另一段程式碼。
                lines.append("  %08X  db %02X" % (ea, ida_bytes.get_byte(ea)))
                ea += 1
                continue
            lines.append("  %08X  %-40s ; %s" % (
                ea, idc.GetDisasm(ea),
                " ".join("%02X" % ida_bytes.get_byte(ea + i) for i in range(n))))
            ea += n
    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("\n".join(lines) + "\n")
    ida_pro.qexit(0)


main()
