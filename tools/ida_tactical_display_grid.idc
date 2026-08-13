// DOS/V KI.EXE 戰術顯示格證據匯出（IDA Pro 9.4）。
//
// 非破壞性：只讀函式邊界、原始名稱、指令、code/data xref；不改名、不加型別。
// 位址皆為 IDA DOS/V linear address。語意分級留給版控文件審查。
#include <idc.idc>

static dump_xrefs_to(out, ea)
{
    auto x;
    fprintf(out, "  XREFS_TO\n");
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "    CREF TYPE=%d FROM=%08X OWNER=%s %s\n",
                XrefType(), x, get_func_name(x), GetDisasm(x));
    }
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        fprintf(out, "    DREF TYPE=%d FROM=%08X OWNER=%s %s\n",
                XrefType(), x, get_func_name(x), GetDisasm(x));
    }
}

static dump_function(out, name)
{
    auto ea, start, end, p, x;
    ea = get_name_ea_simple(name);
    fprintf(out, "\nFUNCTION ORIGINAL=%s EA=%08X\n", name, ea);
    if (ea == BADADDR) {
        fprintf(out, "  STATUS=missing\n");
        return;
    }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    fprintf(out, "  RANGE=%08X-%08X BYTES=%u\n", start, end, end-start);
    dump_xrefs_to(out, start);
    fprintf(out, "  INSTRUCTIONS\n");
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "    %08X %s\n", p, GetDisasm(p));
        for (x = get_first_dref_from(p); x != BADADDR; x = get_next_dref_from(p, x)) {
            fprintf(out, "      DREF_FROM TYPE=%d TO=%08X NAME=%s\n",
                    XrefType(), x, get_name(x));
        }
    }
}

static dump_data(out, name)
{
    auto ea, x;
    ea = get_name_ea_simple(name);
    fprintf(out, "\nDATA ORIGINAL=%s EA=%08X\n", name, ea);
    if (ea == BADADDR) { return; }
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        fprintf(out, "  DREF TYPE=%d FROM=%08X OWNER=%s %s\n",
                XrefType(), x, get_func_name(x), GetDisasm(x));
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/out/ida-tactical-display-grid.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "IDA Pro 9.4 tactical display-grid export\n");
    fprintf(out, "Input: DOS/V KI.EXE SHA-256 fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868\n");
    fprintf(out, "Address space: IDA DOS/V linear addresses\n");
    fprintf(out, "Contract: original locator + raw instruction/xref only; no semantic renames\n");

    dump_function(out, "sub_19946");
    dump_function(out, "sub_1B240");
    dump_function(out, "sub_1B360");
    dump_function(out, "sub_1B3B2");
    dump_function(out, "sub_1BAB7");
    dump_function(out, "sub_1BB10");
    dump_function(out, "sub_1D958");
    dump_function(out, "sub_1D971");
    dump_function(out, "sub_1D9D1");
    dump_function(out, "sub_1DA1C");
    dump_function(out, "sub_1DAAA");
    dump_function(out, "sub_1DB34");
    dump_function(out, "sub_1DB9B");
    dump_function(out, "sub_1DC03");
    dump_function(out, "sub_1DC9D");
    dump_function(out, "sub_1DD22");
    dump_function(out, "sub_1DDB4");
    dump_function(out, "sub_1DE95");
    dump_function(out, "sub_1DFA6");
    dump_function(out, "sub_1DFBB");
    dump_function(out, "sub_1DFE8");
    dump_function(out, "sub_1E011");
    dump_function(out, "sub_1E085");
    dump_function(out, "sub_1E0E1");

    dump_data(out, "word_1E158");
    dump_data(out, "word_1E15A");
    dump_data(out, "word_1E15C");
    dump_data(out, "word_1E15E");
    dump_data(out, "word_1E160");
    dump_data(out, "word_1E162");
    fclose(out);
    qexit(0);
}
