// IDA Pro 9.4 全函式普查（DOS/V KI.EXE）。
//
// 目的：把「哪些函式已經有人讀過、哪些從來沒被碰過」變成可查的表，
// 而不是靠印象。只輸出資料庫事實——函式邊界、指令數、呼叫邊、
// 間接控制流與資料參考；不改名、不加型別、不寫註解，
// 也不把 IDA 的導航符號當成遊戲語意。
//
// 輸出 /work/census.tsv，每行一筆，欄位以 tab 分隔：
//
//   FUNC  start end bytes insns name callers calls indirect_calls jumps indirect_jumps drefs
//   CALL  caller_start callee_ea callee_name         直接呼叫邊（XrefType 16/17）
//   ICALL caller_start site_ea operand_type disasm   間接呼叫點
//   DREF  func_start site_ea target_ea xref_type     函式碰到的資料位址
//   ORPHAN ea name                                    沒有直接呼叫者的函式
//   NOFUNC ea                                         被 call 但沒有函式邊界的目標
//
// 用法見 tools/ida_census.sh。
#include <idc.idc>

static is_indirect_operand(t)
{
    // o_reg=1, o_mem=2, o_phrase=3, o_displ=4, o_imm=5, o_far=6, o_near=7。
    // 只有 far/near 是靜態可解析的直接目標，其餘都要靠 runtime 才知道去哪。
    return t != 6 && t != 7;
}

static count_callers(ea)
{
    auto x, n;
    n = 0;
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        n = n + 1;
    }
    return n;
}

static census_one(out, start)
{
    auto end, p, m, t, c, ct, x;
    auto insns, calls, icalls, jumps, ijumps, drefs;
    end = get_func_attr(start, FUNCATTR_END);
    if (end == BADADDR) {
        fprintf(out, "FUNC\t%08X\tBADADDR\t0\t0\t%s\t%u\t0\t0\t0\t0\t0\n",
                start, get_func_name(start), count_callers(start));
        return;
    }
    insns = 0; calls = 0; icalls = 0; jumps = 0; ijumps = 0; drefs = 0;
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        insns = insns + 1;
        m = print_insn_mnem(p);
        if (m == "call") {
            calls = calls + 1;
            t = get_operand_type(p, 0);
            if (is_indirect_operand(t)) {
                icalls = icalls + 1;
                fprintf(out, "ICALL\t%08X\t%08X\t%d\t%s\n", start, p, t, GetDisasm(p));
            }
            for (c = get_first_cref_from(p); c != BADADDR; c = get_next_cref_from(p, c)) {
                // IDA 對 call 同時回傳呼叫目標與下一指令的 ordinary flow。
                // fl_CN=16 / fl_CF=17 才是呼叫邊；fl_F=21 是循序落下，不能算。
                ct = XrefType();
                if (ct == 16 || ct == 17) {
                    fprintf(out, "CALL\t%08X\t%08X\t%s\n", start, c, get_func_name(c));
                    if (get_func_attr(c, FUNCATTR_START) == BADADDR) {
                        fprintf(out, "NOFUNC\t%08X\n", c);
                    }
                }
            }
        }
        if (substr(m, 0, 1) == "j") {
            jumps = jumps + 1;
            t = get_operand_type(p, 0);
            if (is_indirect_operand(t)) {
                ijumps = ijumps + 1;
                fprintf(out, "IJUMP\t%08X\t%08X\t%d\t%s\n", start, p, t, GetDisasm(p));
            }
        }
        for (x = get_first_dref_from(p); x != BADADDR; x = get_next_dref_from(p, x)) {
            drefs = drefs + 1;
            fprintf(out, "DREF\t%08X\t%08X\t%08X\t%d\n", start, p, x, XrefType());
        }
    }
    fprintf(out, "FUNC\t%08X\t%08X\t%u\t%u\t%s\t%u\t%u\t%u\t%u\t%u\t%u\n",
            start, end, end - start, insns, get_func_name(start),
            count_callers(start), calls, icalls, jumps, ijumps, drefs);
    if (count_callers(start) == 0) {
        fprintf(out, "ORPHAN\t%08X\t%s\n", start, get_func_name(start));
    }
}

static main()
{
    auto out, ea, n;
    Wait();
    out = fopen("/work/census.tsv", "w");
    if (out == 0) { return; }
    fprintf(out, "# IDA Pro 9.4 全函式普查\n");
    fprintf(out, "# 位址空間：IDA DOS/V linear address\n");
    fprintf(out, "# 符號名是資料庫導航標籤，不帶遊戲語意\n");
    n = 0;
    for (ea = get_next_func(0); ea != BADADDR; ea = get_next_func(ea)) {
        census_one(out, ea);
        n = n + 1;
    }
    fprintf(out, "TOTAL\t%u\n", n);
    fclose(out);
}
