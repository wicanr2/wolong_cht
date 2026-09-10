#!/usr/bin/env python3
# IDAPython：查一支函式的**回傳值有沒有被消費**。
#
#   printf '1474A 壞滅判定\n15130 自動判定\n' \
#     > workplace/ida/dosv/census/consumer_list.txt
#   tools/ida.sh script dosv tools/ida_return_consumers.py KI.EXE.i64
#
# 為什麼要有這一支：追一支函式追到它「設了旗標／回了 CF／寫了 AX」，
# 很容易就下「所以會發生 X」的結論——**而回傳值可能被整個丟掉**。
#
# ⛔ 2026-09-10 連踩三次同一個形狀（教訓 `flag-set-without-consumer`）：
#   · `sub_1474A` 士氣 0 回 STC ⇒ 斷言「士氣不足 100 的軍團打贏也會散掉」。
#     實際上 `sub_14ADE` 拿「誰贏」分流，**勝方那一位根本不看**。
#   · 先斷言 `sub_14E5C`／`sub_14ED7` 不看 `sub_15130` 的回傳值 ⇒
#     「壞滅沒有消費端」。實際上看的是它們的**呼叫端**——只追一層不夠。
#   · 軍團 `+0x00` 位元 0 的設定端與清除端都記了、消費端沒記，
#     於是它被當成顯示狀態。
#
# 這支對每個呼叫端印**呼叫之後的 N 條指令**，一眼看得出：
#   · `jb`／`jnb`／`jc`／`jnc`      → 消費 CF
#   · `and al,al`／`test ah,2`／`cmp al,…` → 消費 AX
#   · 直接 `retn`／`pop`／下一個 `call` → **沒消費**（回傳值被丟掉）
#
# ⚠ **它只看直線的下幾條指令**，不做資料流分析：呼叫端把 AX 存起來
# 稍後再看的情形抓不到，跨基本區塊的也抓不到。**輸出是候選不是結論**——
# 「看起來沒消費」要再往上追一層（那正是第二次踩坑的地方）。
#
# 輸出 /work/return_consumers.txt。第一行是 probe（函式數 ＋ 輸入檔雜湊）。
import ida_auto
import ida_funcs
import ida_nalt
import ida_pro
import ida_segment
import ida_ua
import idautils
import idc

LIST = "/work/consumer_list.txt"
OUT = "/work/return_consumers.txt"
AFTER = 4  # 呼叫之後看幾條指令

# 消費 CF 的條件跳躍（16-bit：jb/jnb/jc/jnc/jae/jbe/ja…）
CF_JUMPS = {"jb", "jnb", "jc", "jnc", "jae", "jbe", "ja", "jnae", "jnb"}
# 消費 AX／旗標的比較
CMP_OPS = {"and", "test", "cmp", "or", "xor", "shr", "shl", "sar"}


def verdict(insns):
    """從呼叫之後那幾條指令猜有沒有消費回傳值。"""
    for _, mnem, dis in insns:
        if mnem in CF_JUMPS:
            return "消費 CF"
        if mnem in CMP_OPS and ("a" in dis.split(None, 1)[-1][:6]
                                or "al" in dis or "ah" in dis or "ax" in dis):
            return "消費 AX"
        if mnem in ("retn", "retf", "ret"):
            return "⚠ 沒消費（直接 retn）"
        if mnem == "call":
            return "⚠ 沒消費（下一個 call）"
    return "？（前 %d 條看不出來）" % AFTER


def main():
    ida_auto.auto_wait()
    want = []
    for line in open(LIST, encoding="utf-8"):
        line = line.split("#")[0].strip()
        if not line:
            continue
        parts = line.split(None, 1)
        want.append((int(parts[0], 16), parts[1] if len(parts) > 1 else ""))

    # 逐 segment 掃 call，收「呼叫目標 → 呼叫點」。
    # ⭐ 逐 segment 不逐函式：沒有呼叫端的常式不會變成函式，
    # 而 16-bit near call 的 `op.addr` 是段內 offset 不是 linear address，
    # 所以用「下一條指令的位址 ＋ 位移」還原不了——直接比 GetDisasm 的目標名。
    calls = {}
    nseg = 0
    for s_ea in idautils.Segments():
        seg = ida_segment.getseg(s_ea)
        if seg.type not in (ida_segment.SEG_CODE, ida_segment.SEG_NORM):
            continue
        nseg += 1
        head = seg.start_ea
        while head < seg.end_ea:
            insn = ida_ua.insn_t()
            n = ida_ua.decode_insn(insn, head)
            if n <= 0:
                head += 1
                continue
            if insn.get_canon_mnem() == "call":
                for tgt, _ in want:
                    name = idc.get_func_name(tgt)
                    if name and name in idc.GetDisasm(head):
                        calls.setdefault(tgt, []).append(head)
            head += n

    with open(OUT, "w", encoding="utf-8") as fh:
        fh.write("probe: %d 支函式 / %d 段, sha256=%s\n"
                 % (len(list(idautils.Functions())), nseg,
                    ida_nalt.retrieve_input_file_sha256().hex()[:16]))
        for tgt, note in want:
            sites = calls.get(tgt, [])
            fh.write("\n==== %08X %s（%s）：%d 個呼叫端 ====\n"
                     % (tgt, idc.get_func_name(tgt) or "?", note, len(sites)))
            if not sites:
                fh.write("  ⚠ 一個都沒有——下結論前先確認掃描本身有正對照\n")
            for site in sites:
                f = ida_funcs.get_func(site)
                fn = ida_funcs.get_func_name(f.start_ea) if f else "（無函式）"
                insns = []
                ea = site
                for _ in range(AFTER):
                    ea = idc.next_head(ea, idc.BADADDR)
                    if ea == idc.BADADDR:
                        break
                    ins = ida_ua.insn_t()
                    if ida_ua.decode_insn(ins, ea) <= 0:
                        break
                    insns.append((ea, ins.get_canon_mnem(),
                                  idc.GetDisasm(ea).strip()))
                fh.write("  %-12s %08X  → %s\n" % (fn, site, verdict(insns)))
                for ea, _, dis in insns:
                    fh.write("        %08X  %s\n" % (ea, dis))
    ida_pro.qexit(0)


main()
