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
// 輸出 /work/helpers.txt。
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


static dump_words(out, ea, n, label)
{
    auto i, v;
    fprintf(out, "\n==== %s  %08X（%d 筆 word）====\n", label, ea, n);
    for (i = 0; i < n; i = i + 1) {
        v = get_wide_word(ea + i * 2);
        fprintf(out, "  [%d] %04X  -> %08X %s\n", i, v, 0x10000 + v,
                get_func_name(0x10000 + v));
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/helpers.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "文字繪製共用常式 ＋ 領地圖 ＋ 分派表\n位址空間：IDA DOS/V linear address\n");
    fprintf(out, "正對照：0001701A = 「%s」\n", GetDisasm(0x1701A));

    dump_words(out, 0x159D2, 6, "off_159D2 主畫面重繪分派表");

    dump_func(out, 0x106F5, "sub_106F5 畫字串（ax=屬性?）");
    dump_func(out, 0x106FD, "sub_106FD 畫字串");
    dump_func(out, 0x107D2, "sub_107D2 疑似肖像");
    dump_func(out, 0x188B0, "sub_188B0 畫勢力名");
    dump_func(out, 0x15CE0, "sub_15CE0 領地圖畫單一據點");
    dump_func(out, 0x15A3A, "sub_15A3A");
    fclose(out);
}
