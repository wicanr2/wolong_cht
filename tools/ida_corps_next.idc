// IDA Pro 9.4：docs/re/30 §7 留下的缺口（DOS/V KI.EXE）。
//
// 1. dump 四支：sub_16E8F（疑似兵力分配）、sub_16F86、sub_155EC、
//    sub_15AFC（未讀清單裡最大的一支）與它的呼叫者。
// 2. 掃全庫找同時碰到 `+8]` 與 `+9]` 的函式——docs/re/30 §6 的假說是
//    「軍團 +0x09（勢力×5）＋ +0x08（朝向）＝ 大地圖圖塊編號」，
//    要證偽就得找到同時讀這兩欄的常式。
//
// ⚠ 第 2 項是**在反組譯文字上做樣式比對**，不是語意分析：
//    `+8]`／`+9]` 可能屬於任何結構，命中只是候選，要逐支讀過才算數。
//
// ⚠ 只讀不寫：不 add_func、不改名、不加型別、不寫註解。
// 輸出 /work/corps_next.txt。
#include <idc.idc>

static dump_func(out, ea, label)
{
    auto start, end, p, x, n;
    fprintf(out, "\n==== %s  %08X ====\n", label, ea);
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start != ea || end == BADADDR) {
        fprintf(out, "  ⚠ 無獨立函式邊界（所屬：%s），線性反組譯至第一個 retn\n",
                get_func_name(ea));
        p = ea;
        for (n = 0; n < 200 && p != BADADDR; n = n + 1) {
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
}

static scan_pair(out)
{
    auto i, f, start, end, p, d, hit8, hit9, lines;
    fprintf(out, "\n==== 掃描：同時出現 +8] 與 +9] 的函式 ====\n");
    fprintf(out, "（樣式比對，不是語意分析；命中只是候選）\n");
    for (i = 0; i < get_func_qty(); i = i + 1) {
        f = get_func_attr(getn_func(i), FUNCATTR_START);
        if (f == BADADDR) { continue; }
        end = get_func_attr(f, FUNCATTR_END);
        hit8 = 0; hit9 = 0; lines = "";
        for (p = f; p != BADADDR && p < end; p = next_head(p, end)) {
            d = GetDisasm(p);
            if (strstr(d, "+8]") >= 0) { hit8 = 1; lines = lines + sprintf("      %08X  %s\n", p, d); }
            if (strstr(d, "+9]") >= 0) { hit9 = 1; lines = lines + sprintf("      %08X  %s\n", p, d); }
        }
        if (hit8 && hit9) {
            fprintf(out, "\n  %s  %08X  bytes=%u\n", get_func_name(f), f, end - f);
            fprintf(out, "%s", lines);
        }
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/corps_next.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "docs/re/30 §7 的缺口\n位址空間：IDA DOS/V linear address\n");

    dump_func(out, 0x16E8F, "sub_16E8F 疑似兵力分配");
    dump_func(out, 0x16F86, "sub_16F86");
    dump_func(out, 0x155EC, "sub_155EC 預備兵相加");
    dump_func(out, 0x15AFC, "sub_15AFC 最大的未讀函式");

    scan_pair(out);
    fclose(out);
}
