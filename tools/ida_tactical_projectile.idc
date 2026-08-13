// 戰術投射物發射、移動、繪製與命中的非破壞性 IDA Pro 9.4 匯出器。
// 僅輸出既有資料庫的原始函式／xref／運算元；不改名、不加型別、不寫回輸入。
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

static dump_at(out, label, ea)
{
    auto start, end, p;
    fprintf(out, "\n=== ADDRESS %s @ %08X ===\n", label, ea);
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR || end == BADADDR) { fprintf(out, "NOT IN A DATABASE FUNCTION\n"); return; }
    fprintf(out, "FUNCTION %s RANGE %08X-%08X\n", get_func_name(start), start, end);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "%08X  %s\n", p, GetDisasm(p));
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/ida-tactical-projectile.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "IDA Pro 9.4 non-destructive tactical projectile export\n");
    fprintf(out, "Address space: DOS/V IDA linear addresses\n");
    fprintf(out, "Raw instructions/xrefs only; no semantic renaming is performed.\n");

    dump_function(out, "sub_1AB7C");
    dump_function(out, "sub_1AB9C");
    dump_function(out, "sub_1ABB2");
    dump_function(out, "sub_1ABFF");
    dump_function(out, "sub_1AC55");
    dump_function(out, "sub_1ACA4");
    dump_function(out, "sub_1ACD6");
    dump_function(out, "sub_1AD2D");
    dump_function(out, "sub_1AD7F");
    dump_function(out, "sub_1ADC8");
    dump_function(out, "sub_1B8AA");
    dump_function(out, "sub_1B941");
    dump_function(out, "sub_1B97E");
    dump_function(out, "sub_1BA2E");
    dump_function(out, "sub_1BAB7");
    dump_at(out, "special launch caller at 0001AC4A", 0x1AC4A);
    dump_at(out, "special launch caller at 0001ACA0", 0x1ACA0);
    fclose(out);
}
