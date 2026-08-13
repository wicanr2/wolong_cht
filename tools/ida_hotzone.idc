// IDA Pro 9.4：戰略層的滑鼠熱區查表與其周邊常式（DOS/V KI.EXE）。
//
// sub_1E453 由 (cx>>3) + 10*(dx & 0xF8) 取 byte，也就是把像素座標
// 換算成 80 欄格點再查一張區域圖。呼叫端隨後做 `sub al, <base>` 再分派，
// 所以那些減數是熱區編號的基底，不是 ASCII 鍵碼。
//
// 這支腳本把「誰填那張圖」與「分派前後的配套常式」一起讀出來，
// 免得從呼叫端的參數順序反推語意。
//
// 輸出 /work/hotzone.txt。不改名、不加型別、不寫回資料庫。
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

static dump_xrefs(out, name)
{
    auto ea, x, t;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n---- 資料 %s EA=%08X 的讀寫端 ----\n", name, ea);
    if (ea == BADADDR) { return; }
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        // XrefType: dr_O=1 取位址、dr_W=2 寫、dr_R=3 讀。
        t = XrefType();
        fprintf(out, "  TYPE=%d %08X %s  %s\n", t, x, get_func_name(x), GetDisasm(x));
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/hotzone.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "戰略層滑鼠熱區：查表、填表與配套常式\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    // 熱區圖的兩個指標：段與偏移
    dump_xrefs(out, "word_1E479");
    dump_xrefs(out, "word_1E47B");

    // 分派前的把關：sub_1E453 之前一定先呼叫它，CF 代表沒有輸入
    dump_func(out, "sub_121E7");
    // 狀態列提示：cx 是 TALK 索引，cx=0xFFFF 清除（待驗）
    dump_func(out, "sub_18853");
    // 進言選單本體的內層
    dump_func(out, "sub_193E9");
    // 選軍團／移動游標
    dump_func(out, "sub_1716D");
    dump_func(out, "sub_12151");
    // 成對出現的存／復原
    dump_func(out, "sub_12078");
    dump_func(out, "sub_120D6");

    fclose(out);
}
