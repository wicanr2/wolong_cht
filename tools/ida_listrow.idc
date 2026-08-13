// IDA Pro 9.4：一覽表的畫列常式與捲動／選取 handler（DOS/V KI.EXE）。
//
// docs/re/26 已定案引擎與四個描述子。這一支往下讀兩層：
//   - `sub_1748F`：畫一列。它決定每個欄位怎麼取值與格式化，
//     解開就等於拿到五種清單的完整欄位定義。
//   - `funcs_18450` 的五個 handler：捲動與選取的實際行為。
// 另外把其餘四個家族的建清單／畫列 callback 也建函式讀出來，
// 好跟據點家族（docs/re/26 §8）比對。
//
// ⚠ 會 add_func 改寫資料庫，所以一定要用 `tools/ida.sh script`（跑在唯讀副本上）。
// 不改名、不加型別、不寫註解。
//
// 輸出 /work/listrow.txt。
#include <idc.idc>

static ensure_func(out, ea)
{
    if (get_func_attr(ea, FUNCATTR_START) == ea) { return 1; }
    if (add_func(ea) != 0) { return 1; }
    del_items(ea, DELIT_EXPAND, 1);
    if (create_insn(ea) == 0) { return 0; }
    return add_func(ea) != 0;
}

static dump_at(out, ea, label)
{
    auto start, end, p, x, n;
    fprintf(out, "\n==== %s  %08X ====\n", label, ea);
    ensure_func(out, ea);
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == ea && end != BADADDR) {
        n = 0;
        fprintf(out, "呼叫者：");
        for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
            fprintf(out, "%s ", get_func_name(x));
            n = n + 1;
        }
        fprintf(out, " 共 %d 處，bytes=%u\n", n, end - start);
        for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
            fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
        }
        return;
    }
    // 建不了函式時線性走到第一個 retn。**這是線性走訪不是控制流分析。**
    fprintf(out, "  ⚠ 無獨立函式邊界（所屬：%s），線性反組譯至第一個 retn\n",
            get_func_name(ea));
    p = ea;
    for (n = 0; n < 120 && p != BADADDR; n = n + 1) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
        if (print_insn_mnem(p) == "retn" || print_insn_mnem(p) == "retf") { break; }
        p = next_head(p, BADADDR);
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/listrow.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "一覽表：畫列與捲動／選取\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    // 先把所有 callback 建成函式，讓後續解析有完整的流程資訊。
    ensure_func(out, 0x171A8); ensure_func(out, 0x17217); ensure_func(out, 0x1724D);
    ensure_func(out, 0x1743B); ensure_func(out, 0x1745F);
    ensure_func(out, 0x1763C); ensure_func(out, 0x176A0); ensure_func(out, 0x176DC);
    ensure_func(out, 0x178E5); ensure_func(out, 0x17944); ensure_func(out, 0x1796C);
    ensure_func(out, 0x17B6F); ensure_func(out, 0x17B90);
    Wait();

    dump_at(out, 0x1748F, "sub_1748F 畫一列");
    dump_at(out, 0x18607, "sub_18607 畫列前置");

    // funcs_18450 的五個 handler（熱區編號 0x3D..0x41）
    dump_at(out, 0x18458, "funcs_18450[0] nullsub_2");
    dump_at(out, 0x18463, "funcs_18450[1]");
    dump_at(out, 0x184DD, "funcs_18450[2]");
    dump_at(out, 0x1851A, "funcs_18450[3]");
    dump_at(out, 0x18546, "funcs_18450[4]");
    dump_at(out, 0x1857F, "loc_1857F 熱區 >= 0x42");

    // 其餘家族的畫列 callback，跟據點版比對欄位
    dump_at(out, 0x1724D, "軍團-畫列");
    dump_at(out, 0x176DC, "武將-畫列");
    dump_at(out, 0x1796C, "勢力-畫列");
    dump_at(out, 0x17B90, "開局-畫列");

    dump_at(out, 0x11D46, "sub_11D46");

    fclose(out);
}
