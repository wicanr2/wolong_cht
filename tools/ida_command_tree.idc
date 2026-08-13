// IDA Pro 9.4：攤開戰略指令列的八個頂層 handler（DOS/V KI.EXE）。
//
// sub_161CA 是指令列分派器，表 funcs_161FE 的八筆由邊界檢查
// （cx-0x18 < 0x180，除以 0x30）定死。每個 handler 開一個子選單，
// 子選單的項數與 TALK 訊息索引寫在 sub_193E9 的參數裡：
//
//     mov ax, <項數>
//     mov cx, <TALK 訊息索引>
//     mov dx, <位置／旗標>
//     call sub_193E9
//
// 把這三個常數讀出來，再查已全解的 TALK.DAT，指令名就有一手證據，
// 不必從函式行為反猜。
//
// 輸出 /work/cmdtree.txt。不改名、不加型別、不寫回資料庫。
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

static main()
{
    auto out;
    Wait();
    out = fopen("/work/cmdtree.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "戰略指令列：八個頂層 handler 與選單常式\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    // funcs_161FE 的八筆，依序
    dump_func(out, "sub_16224");   // [0]
    dump_func(out, "sub_16265");   // [1]
    dump_func(out, "sub_1678D");   // [2]
    dump_func(out, "sub_16288");   // [3]
    dump_func(out, "sub_1628F");   // [4]
    dump_func(out, "sub_162FB");   // [5]
    dump_func(out, "sub_16366");   // [6]
    dump_func(out, "sub_163BF");   // [7]

    // 選單常式本身：想知道 ax/cx/dx 各是什麼就得讀它
    dump_func(out, "sub_193E9");
    // 鍵盤／游標輸入來源：sub_1678D 與 sub_18412 都靠它回傳 al
    dump_func(out, "sub_1E453");
    // 訊息顯示：cx 的語意是 CONTEXT.md 的未解項之一
    dump_func(out, "sub_18810");

    fclose(out);
}
