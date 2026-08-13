// IDA Pro 9.4：追 ICONGRF 段 1 的圖塊偏移從哪來（DOS/V KI.EXE）。
//
// docs/formats/03 §5.3 指出段 1 不是密排的圖塊陣列：`sub_1CA3B` 用的 `si`
// 是 `3E00h`／`3E80h`／`3DC0h` 這類段內偏移，中間夾著別的東西。
// 該節寫明下手點是「逐一追每個 si 常數的來源」，不要整段當同尺寸陣列切。
//
// 這支腳本做兩件事：
//   1. 列出 sub_1CA3B 的每個呼叫者，以及呼叫前對 si 的賦值
//   2. 掃全資料庫，找出所有把立即值寫進 si 且落在段 1 值域的指令
//
// 輸出 /work/icongrf-seg1.txt。不改名、不加型別、不寫回資料庫。
#include <idc.idc>

static dump_func(out, name)
{
    auto ea, start, end, p, x;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n==== %s EA=%08X ====\n", name, ea);
    if (ea == BADADDR) { fprintf(out, "  找不到符號\n"); return; }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR) { fprintf(out, "  沒有函式邊界\n"); return; }
    fprintf(out, "呼叫者：");
    for (x = get_first_cref_to(start); x != BADADDR; x = get_next_cref_to(start, x)) {
        fprintf(out, "%s(%08X) ", get_func_name(x), x);
    }
    fprintf(out, "  bytes=%u\n", end - start);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
    }
}

// 印呼叫點前 n 條指令，看 si／ax／cx 是在哪裡被設定的。
static dump_before(out, site, n)
{
    auto p, i, seq;
    fprintf(out, "\n---- 呼叫點 %08X（%s）之前 %d 條 ----\n",
            site, get_func_name(site), n);
    // prev_head 往回走，先蒐集再倒著印。
    seq = site;
    for (i = 0; i < n; i = i + 1) {
        p = prev_head(seq, 0);
        if (p == BADADDR) { break; }
        seq = p;
    }
    for (p = seq; p != BADADDR && p <= site; p = next_head(p, site + 1)) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
    }
}

static main()
{
    auto out, ea, target, x, p, m, op, v, seg1_hits;
    Wait();
    out = fopen("/work/icongrf-seg1.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "ICONGRF 段 1：圖塊偏移的來源追蹤\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    dump_func(out, "sub_1CA3B");
    // 段 3 走的另一支繪製常式，docs/formats/03 §5.3 標為還沒讀
    dump_func(out, "sub_1F888");
    dump_func(out, "sub_1CAA8");
    dump_func(out, "sub_1C7F4");
    dump_func(out, "sub_1C863");
    dump_func(out, "sub_1FA37");

    target = get_name_ea_simple("sub_1CA3B");
    if (target != BADADDR) {
        for (x = get_first_cref_to(target); x != BADADDR;
             x = get_next_cref_to(target, x)) {
            dump_before(out, x, 12);
        }
    }

    // 全庫掃 `mov si, <立即值>`，挑落在 ICONGRF 段 1 值域的。
    // 值域取寬一點（0x3000–0x4000），寧可多列也不要漏。
    fprintf(out, "\n==== 全庫 `mov si, imm` 落在 3000h–4000h ====\n");
    seg1_hits = 0;
    for (ea = get_next_func(0); ea != BADADDR; ea = get_next_func(ea)) {
        auto endea;
        endea = get_func_attr(ea, FUNCATTR_END);
        if (endea == BADADDR) { continue; }
        for (p = ea; p != BADADDR && p < endea; p = next_head(p, endea)) {
            m = print_insn_mnem(p);
            if (m != "mov") { continue; }
            if (print_operand(p, 0) != "si") { continue; }
            if (get_operand_type(p, 1) != 5) { continue; }  // o_imm
            v = get_operand_value(p, 1);
            if (v >= 0x3000 && v < 0x4000) {
                fprintf(out, "  %08X  %-28s  %s\n", p, get_func_name(p), GetDisasm(p));
                seg1_hits = seg1_hits + 1;
            }
        }
    }
    fprintf(out, "  共 %d 處\n", seg1_hits);

    fclose(out);
}
