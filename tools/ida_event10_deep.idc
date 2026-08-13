// 事件 10 producer 追查的非破壞性 IDA Pro 9.4 匯出器。
//
// 本工具只讀目前資料庫的函式邊界與 code/data xref，輸出原始函式名、
// 線性位址與運算元；不改名、不加型別、不寫回輸入檔。它把所有已知
// queue writer 的直接 caller 一次匯出，讓「尚未找到 0x0A producer」
// 能被檢查，而不是由散文結論代替證據。
#include <idc.idc>

static dump_callers(out, name)
{
    auto ea, x, n;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n=== CALLERS %s @ %08X ===\n", name, ea);
    if (ea == BADADDR) {
        fprintf(out, "UNKNOWN SYMBOL\n");
        return;
    }
    n = 0;
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "%08X  %s  %s\n", x, get_func_name(x), GetDisasm(x));
        n = n + 1;
    }
    if (n == 0) {
        fprintf(out, "(no direct code xref)\n");
    }
}

static dump_function(out, name)
{
    auto ea, start, end, p;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n=== FUNCTION %s @ %08X ===\n", name, ea);
    if (ea == BADADDR) {
        fprintf(out, "UNKNOWN SYMBOL\n");
        return;
    }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR || end == BADADDR) {
        fprintf(out, "NOT A DATABASE FUNCTION\n");
        return;
    }
    fprintf(out, "RANGE %08X-%08X\n", start, end);
    dump_callers(out, name);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "%08X  %s\n", p, GetDisasm(p));
    }
}

static dump_refs(out, name)
{
    auto ea, x, n, t, kind;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n=== DATA REFS %s @ %08X ===\n", name, ea);
    if (ea == BADADDR) {
        fprintf(out, "UNKNOWN SYMBOL\n");
        return;
    }
    n = 0;
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        t = XrefType();
        kind = "other";
        if (t == 1) { kind = "address"; }
        if (t == 2) { kind = "write"; }
        if (t == 3) { kind = "read"; }
        fprintf(out, "%s  %08X  %s  %s\n", kind, x, get_func_name(x), GetDisasm(x));
        n = n + 1;
    }
    if (n == 0) {
        fprintf(out, "(no direct data xref; this does not exclude register/pointer access)\n");
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/ida-event10-deep.txt", "w");
    if (out == 0) {
        return;
    }
    fprintf(out, "IDA Pro 9.4 non-destructive event 10 producer export\n");
    fprintf(out, "Address space: DOS/V IDA linear addresses\n");
    fprintf(out, "Raw instructions/xrefs only; no semantic renaming is performed.\n");

    // 初始化與每月／每時入口。
    dump_function(out, "sub_11AC3");
    dump_function(out, "sub_15358");
    dump_function(out, "sub_131AE");
    dump_function(out, "sub_13496");

    // sub_12FBF 的全部已知 direct callers 及它們的 wrapper callers。
    dump_function(out, "sub_12286");
    dump_function(out, "sub_122DB");
    dump_function(out, "sub_12D3A");
    dump_function(out, "sub_12E33");
    dump_function(out, "sub_12E89");
    dump_function(out, "sub_12EFB");
    dump_function(out, "sub_12F71");
    dump_function(out, "sub_15715");
    dump_function(out, "sub_1578F");
    dump_function(out, "sub_157FE");

    // sub_1301C 的全部已知 direct caller。
    dump_function(out, "sub_1300E");
    dump_function(out, "sub_134B1");
    dump_function(out, "sub_15940");
    dump_function(out, "sub_16623");

    dump_refs(out, "word_10D20");
    dump_refs(out, "word_10D56");
    dump_refs(out, "byte_131AD");
    fclose(out);
}
