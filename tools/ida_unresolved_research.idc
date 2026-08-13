// DOS/V 未解項目的非破壞性 IDA 9.4 證據匯出器。
//
// 只讀 IDA 資料庫；不改名、不加型別、不寫回原始輸入。輸出保留原始函式名、
// 線性位址、運算元與 IDA xref 型別，供事件 10、事件 6/7、MCH、戰術與音源
// 的分級研究使用。
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

static dump_at(out, label, ea)
{
    auto start, end, p;
    fprintf(out, "\n=== ADDRESS %s @ %08X ===\n", label, ea);
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR || end == BADADDR) {
        fprintf(out, "NOT IN A DATABASE FUNCTION\n");
        return;
    }
    fprintf(out, "FUNCTION %s RANGE %08X-%08X\n", get_func_name(start), start, end);
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
        if (n == 128) {
            fprintf(out, "TRUNCATED AFTER 128 DIRECT DATA XREFS\n");
            break;
        }
    }
    if (n == 0) {
        fprintf(out, "(no direct data xref; this does not exclude register/pointer access)\n");
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/ida-unresolved-research.txt", "w");
    if (out == 0) {
        return;
    }
    fprintf(out, "IDA Pro 9.4 non-destructive research export\n");
    fprintf(out, "Address space: DOS/V IDA linear addresses\n");
    fprintf(out, "Raw instructions/xrefs only; semantic conclusions require separate review.\n");

    fprintf(out, "\n######## EVENT 10 / QUEUE ########\n");
    dump_function(out, "sub_12BD9");
    dump_function(out, "sub_12FB1");
    dump_function(out, "sub_12FBF");
    dump_function(out, "sub_1300E");
    dump_function(out, "sub_1301C");
    dump_function(out, "sub_1304E");
    dump_function(out, "sub_131AE");
    dump_function(out, "sub_13496");
    dump_function(out, "sub_15358");
    dump_function(out, "sub_15940");
    dump_refs(out, "word_10D20");
    dump_refs(out, "word_10D56");
    dump_refs(out, "byte_131AD");
    dump_refs(out, "funcs_131E8");

    fprintf(out, "\n######## EVENT 6/7 FORMATTER ########\n");
    dump_function(out, "sub_13C3D");
    dump_function(out, "sub_13220");
    dump_function(out, "sub_13262");
    dump_function(out, "sub_13327");
    dump_function(out, "sub_13388");
    dump_function(out, "sub_137D8");
    dump_function(out, "sub_13138");
    dump_function(out, "sub_1084A");
    dump_function(out, "sub_18810");
    dump_at(out, "formatter handler lookup", 0x108DB);

    fprintf(out, "\n######## MMAP.MCH OBJECTS ########\n");
    dump_function(out, "sub_1237E");
    dump_function(out, "sub_123FF");
    dump_function(out, "sub_12438");
    dump_function(out, "sub_12459");
    dump_function(out, "sub_1248A");
    dump_function(out, "sub_124FF");
    dump_function(out, "sub_12533");
    dump_function(out, "sub_134A6");
    dump_function(out, "sub_134B1");
    dump_refs(out, "word_1987A");

    fprintf(out, "\n######## TACTICAL / PROJECTILES ########\n");
    dump_function(out, "sub_1B941");
    dump_function(out, "sub_1B97E");
    dump_function(out, "sub_1BA2E");
    dump_function(out, "sub_1BAB7");
    dump_function(out, "sub_1AD2D");
    dump_function(out, "sub_1AD7F");
    dump_function(out, "sub_1B8AA");
    dump_function(out, "sub_1B240");
    dump_function(out, "sub_1BB10");

    fprintf(out, "\n######## AUDIO CANDIDATES ########\n");
    dump_function(out, "sub_102F5");
    dump_function(out, "sub_10CDE");
    dump_function(out, "sub_1EB11");
    fclose(out);
}
