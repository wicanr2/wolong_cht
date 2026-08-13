// IDA Pro 9.4：人事指令的四支處理函式（DOS/V KI.EXE）。
//
// docs/re/22 §3.2 解出 funcs_16279 的四筆分派：
//   sub_16A9B 內政官任命 / sub_16B08 內政官解任
//   sub_16B71 外交官任命 / sub_16BE3 外交官解任
// 四支都還沒讀。這一支把它們連同同區未讀的鄰居一起 dump，
// 目標是補 docs/mechanics/60-personnel.md 的任免規則。
//
// 一併帶出每支碰到的資料位址還有誰讀寫——「誰還在用它」是最便宜的證據。
//
// ⚠ 只讀不寫：不 add_func、不改名、不加型別、不寫註解。
// 輸出 /work/personnel.txt。
#include <idc.idc>

static xref_kind(t)
{
    if (t == 1) { return "取址"; }
    if (t == 2) { return "寫"; }
    if (t == 3) { return "讀"; }
    return "?";
}

static dump_data_users(out, ea)
{
    auto x, n, t;
    n = 0;
    fprintf(out, "      %04X 的其他使用者：", ea);
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        // 逐筆判定讀／寫，不比對助憶碼字串（CLAUDE.md §4.1）
        t = 0;
        auto y;
        for (y = get_first_dref_from(x); y != BADADDR; y = get_next_dref_from(x, y)) {
            if (y == ea) { t = XrefType(); }
        }
        fprintf(out, "%s(%s) ", get_func_name(x), xref_kind(t));
        n = n + 1;
        if (n > 12) { fprintf(out, "…"); break; }
    }
    fprintf(out, " 共 %d 處\n", n);
}

static dump_func(out, ea, label)
{
    auto start, end, p, x, n, y, seen, k;
    fprintf(out, "\n==== %s  %08X ====\n", label, ea);
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start != ea || end == BADADDR) {
        fprintf(out, "  ⚠ 無獨立函式邊界（所屬：%s），線性反組譯至第一個 retn\n",
                get_func_name(ea));
        p = ea;
        for (n = 0; n < 140 && p != BADADDR; n = n + 1) {
            fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
            if (print_insn_mnem(p) == "retn" || print_insn_mnem(p) == "retf") { break; }
            p = next_head(p, BADADDR);
        }
        return;
    }
    n = 0;
    fprintf(out, "呼叫者：");
    for (x = get_first_cref_to(start); x != BADADDR; x = get_next_cref_to(start, x)) {
        fprintf(out, "%s ", get_func_name(x));
        n = n + 1;
    }
    fprintf(out, " 共 %d 處，bytes=%u\n", n, end - start);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
    }
    // 這支碰到的資料位址，逐一列出還有誰用
    fprintf(out, "  ── 資料參考 ──\n");
    seen = object();
    k = 0;
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        for (y = get_first_dref_from(p); y != BADADDR; y = get_next_dref_from(p, y)) {
            if (get_func_attr(y, FUNCATTR_START) != BADADDR) { continue; }  // 略過函式
            fprintf(out, "    %08X → %04X %s\n", p, y, get_name(y));
            dump_data_users(out, y);
            k = k + 1;
            if (k > 40) { fprintf(out, "    …（截斷）\n"); return; }
        }
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/personnel.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "人事指令：內政官／外交官的任免\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    dump_func(out, 0x16A9B, "sub_16A9B 內政官任命");
    dump_func(out, 0x16B08, "sub_16B08 內政官解任");
    dump_func(out, 0x16B71, "sub_16B71 外交官任命");
    dump_func(out, 0x16BE3, "sub_16BE3 外交官解任");

    // 同區未讀的鄰居：可能是共用的挑人／檢核常式
    dump_func(out, 0x16B4F, "sub_16B4F");
    dump_func(out, 0x16C2A, "sub_16C2A");
    dump_func(out, 0x16C92, "sub_16C92 本區最大的未讀函式");

    fclose(out);
}
