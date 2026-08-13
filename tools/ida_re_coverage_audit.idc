// IDA Pro 9.4 DOS/V KI.EXE 結構覆蓋審計。
//
// 只輸出資料庫事實：函式邊界、指令數、直接呼叫與間接控制流候選。
// 不改名、不加型別、不寫註解，也不把符號名提升成遊戲語意。
#include <idc.idc>

static is_indirect_operand(t)
{
    // IDC operand types: o_reg=1, o_mem=2, o_phrase=3, o_displ=4,
    // o_imm=5, o_far=6, o_near=7. 只有 far/near 是靜態直接目標。
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

static audit_function(out, name)
{
    auto ea, start, end, p, m, insns, calls, jumps, indirect_calls;
    auto indirect_jumps, drefs_from, callers, t;
    ea = get_name_ea_simple(name);
    fprintf(out, "\nTARGET %s EA=%08X\n", name, ea);
    if (ea == BADADDR) {
        fprintf(out, "  STATUS=missing\n");
        return;
    }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR || end == BADADDR) {
        fprintf(out, "  STATUS=not-a-function\n");
        return;
    }
    insns = 0;
    calls = 0;
    jumps = 0;
    indirect_calls = 0;
    indirect_jumps = 0;
    drefs_from = 0;
    callers = count_callers(start);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        insns = insns + 1;
        m = print_insn_mnem(p);
        if (m == "call") {
            calls = calls + 1;
            t = get_operand_type(p, 0);
            if (is_indirect_operand(t)) {
                indirect_calls = indirect_calls + 1;
                fprintf(out, "  INDIRECT_CALL %08X TYPE=%d %s\n", p, t, GetDisasm(p));
            }
            auto c, ct;
            for (c = get_first_cref_from(p); c != BADADDR; c = get_next_cref_from(p, c)) {
                // IDA 同時回傳 call target 與下一指令的 ordinary-flow cref；
                // 保留原始 XrefType，不能把兩者都冒稱為呼叫邊。
                ct = XrefType();
                fprintf(out, "  CREF TYPE=%d %08X -> %08X %s\n", ct, p, c, get_func_name(c));
            }
        }
        if (substr(m, 0, 1) == "j") {
            jumps = jumps + 1;
            t = get_operand_type(p, 0);
            if (is_indirect_operand(t)) {
                indirect_jumps = indirect_jumps + 1;
                fprintf(out, "  INDIRECT_JUMP %08X TYPE=%d %s\n", p, t, GetDisasm(p));
            }
        }
        auto x;
        for (x = get_first_dref_from(p); x != BADADDR; x = get_next_dref_from(p, x)) {
            drefs_from = drefs_from + 1;
        }
    }
    fprintf(out, "  RANGE=%08X-%08X BYTES=%u INSNS=%u DIRECT_CALLERS=%u\n",
            start, end, end-start, insns, callers);
    fprintf(out, "  CALLS=%u INDIRECT_CALLS=%u JUMPS=%u INDIRECT_JUMPS=%u DREFS_FROM=%u\n",
            calls, indirect_calls, jumps, indirect_jumps, drefs_from);
}

static audit_address_owner(out, label, ea)
{
    fprintf(out, "\nADDRESS %s EA=%08X OWNER=%s START=%08X END=%08X DISASM=%s\n",
            label, ea, get_func_name(ea),
            get_func_attr(ea, FUNCATTR_START), get_func_attr(ea, FUNCATTR_END), GetDisasm(ea));
}

static audit_data_refs(out, name)
{
    auto ea, x, t;
    ea = get_name_ea_simple(name);
    fprintf(out, "\nDATA %s EA=%08X\n", name, ea);
    if (ea == BADADDR) { return; }
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        t = XrefType();
        fprintf(out, "  XREF TYPE=%d FROM=%08X OWNER=%s %s\n",
                t, x, get_func_name(x), GetDisasm(x));
    }
}

static dump_range(out, label, start, end)
{
    auto p;
    fprintf(out, "\nRANGE %s %08X-%08X\n", label, start, end);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "  %08X OWNER=%s %s\n", p, get_func_name(p), GetDisasm(p));
    }
}

static main()
{
    auto out, ea, end, fn_count, fn_bytes, fn_insns, no_callers, p;
    Wait();
    out = fopen("/out/ida-re-coverage-audit.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "IDA Pro 9.4 structural audit\n");
    fprintf(out, "Input: DOS/V KI.EXE.i64 SHA-256 7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26\n");
    fprintf(out, "Address space: IDA DOS/V linear addresses\n");
    fprintf(out, "Semantics: none; names are original database navigation labels only\n\n");

    fn_count = 0;
    fn_bytes = 0;
    fn_insns = 0;
    no_callers = 0;
    for (ea = get_next_func(0); ea != BADADDR; ea = get_next_func(ea)) {
        end = get_func_attr(ea, FUNCATTR_END);
        fn_count = fn_count + 1;
        if (end != BADADDR) {
            fn_bytes = fn_bytes + end - ea;
            for (p = ea; p != BADADDR && p < end; p = next_head(p, end)) {
                fn_insns = fn_insns + 1;
            }
        }
        if (count_callers(ea) == 0) { no_callers = no_callers + 1; }
    }
    fprintf(out, "DATABASE FUNCTIONS=%u FUNCTION_BYTES=%u FUNCTION_HEADS=%u NO_DIRECT_CALLERS=%u\n",
            fn_count, fn_bytes, fn_insns, no_callers);

    // 戰術入口、逐幀更新、腳本、投影、sprite／HUD 與畫面合成候選。
    audit_function(out, "sub_19A33");
    audit_function(out, "sub_19FA0");
    audit_function(out, "sub_1A1C5");
    audit_function(out, "sub_19CB3");
    audit_function(out, "sub_19E10");
    audit_function(out, "sub_1A065");
    audit_function(out, "sub_1A426");
    audit_function(out, "sub_1A4AB");
    audit_function(out, "sub_1ADC8");
    audit_function(out, "sub_1B240");
    audit_function(out, "sub_1B8AA");
    audit_function(out, "sub_1B941");
    audit_function(out, "sub_1BB10");
    audit_function(out, "sub_1BB6D");
    audit_function(out, "sub_1C51E");
    audit_function(out, "sub_1C7A9");
    audit_function(out, "sub_1C83E");
    audit_function(out, "sub_1C863");
    audit_function(out, "sub_1CBE5");
    audit_function(out, "sub_1DAAA");
    audit_function(out, "sub_1DA1C");
    audit_function(out, "sub_1DC03");
    audit_function(out, "sub_1DC9D");
    audit_function(out, "sub_1D958");
    audit_function(out, "sub_1D971");
    audit_function(out, "sub_1D9D1");
    audit_function(out, "sub_1DDB4");
    audit_function(out, "sub_1E085");
    audit_function(out, "sub_1E0E1");
    audit_address_owner(out, "loc_1A065", 0x1A065);
    audit_data_refs(out, "word_1E160");
    audit_data_refs(out, "word_1E162");
    audit_data_refs(out, "word_1E15C");
    dump_range(out, "self-modifying update/render", 0x1A065, 0x1A156);
    dump_range(out, "camera origin writer", 0x1DC9D, 0x1DD1E);
    dump_range(out, "display-list consumer", 0x1DDB4, 0x1DE95);
    fclose(out);
}
