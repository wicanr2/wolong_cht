// MMAP.MCH type 3／物件更新追查的非破壞性 IDA Pro 9.4 匯出器。
// 僅匯出既有資料庫的原始函式、xref 與位址；不改名、不加型別、不寫回輸入。
#include <idc.idc>

static dump_callers(out, name)
{
    auto ea, x, n;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n=== CALLERS %s @ %08X ===\n", name, ea);
    if (ea == BADADDR) { fprintf(out, "UNKNOWN SYMBOL\n"); return; }
    n = 0;
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "%08X  %s  %s\n", x, get_func_name(x), GetDisasm(x));
        n = n + 1;
    }
    if (n == 0) { fprintf(out, "(no direct code xref)\n"); }
}

static dump_function(out, name)
{
    auto ea, start, end, p;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n=== FUNCTION %s @ %08X ===\n", name, ea);
    if (ea == BADADDR) { fprintf(out, "UNKNOWN SYMBOL\n"); return; }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR || end == BADADDR) { fprintf(out, "NOT A DATABASE FUNCTION\n"); return; }
    fprintf(out, "RANGE %08X-%08X\n", start, end);
    dump_callers(out, name);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "%08X  %s\n", p, GetDisasm(p));
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/ida-mch-deep.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "IDA Pro 9.4 non-destructive MCH object export\n");
    fprintf(out, "Address space: DOS/V IDA linear addresses\n");
    fprintf(out, "Raw instructions/xrefs only; no semantic renaming is performed.\n");

    dump_function(out, "sub_11CC9");
    dump_function(out, "sub_11CD0");
    dump_function(out, "sub_12286");
    dump_function(out, "sub_123FF");
    dump_function(out, "sub_12438");
    dump_function(out, "sub_12459");
    dump_function(out, "sub_1248A");
    dump_function(out, "sub_124FF");
    dump_function(out, "sub_12533");
    dump_function(out, "sub_134B1");
    dump_function(out, "sub_16475");
    dump_function(out, "sub_164F1");
    dump_function(out, "sub_165EF");
    dump_function(out, "sub_1676F");
    fclose(out);
}
