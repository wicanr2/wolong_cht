// IDA Pro 9.4：被大量未讀函式呼叫的共用常式（DOS/V KI.EXE）。
//
// 挑選依據不是直覺，是 tools/re_classify.py 算出來的「被 T4 函式呼叫次數」。
// 認出一支共用常式，所有呼叫它的函式就同時多一條證據——
// 比從最大的函式開始讀划算得多。
//
// 另外 dump 三個「文件引用了、但目前資料庫沒有函式邊界」的位址
// （tools/phantom_scan.py 抓到的）。沒有邊界不代表位址是錯的，
// 要看那裡究竟是什麼才能判斷該修文件還是該修分析狀態。
//
// 輸出 /work/hothelpers.txt。不改名、不加型別、不寫回資料庫。
#include <idc.idc>

static dump_func(out, name)
{
    auto ea, start, end, p, x, n;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n==== %s EA=%08X ====\n", name, ea);
    if (ea == BADADDR) { fprintf(out, "  找不到符號\n"); return; }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR) { fprintf(out, "  沒有函式邊界\n"); return; }
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

// 沒有函式邊界的位址：印出前後脈絡，並看誰參考它。
static dump_loose(out, ea, before, after)
{
    auto p, x, i;
    fprintf(out, "\n==== 無邊界位址 %08X ====\n", ea);
    fprintf(out, "  所屬函式：%s（START=%08X）\n",
            get_func_name(ea), get_func_attr(ea, FUNCATTR_START));
    fprintf(out, "  參考它的：");
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "%08X(%s) ", x, get_func_name(x));
    }
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        fprintf(out, "[data]%08X(%s) ", x, get_func_name(x));
    }
    fprintf(out, "\n");
    p = ea;
    for (i = 0; i < before; i = i + 1) {
        auto q;
        q = prev_head(p, 0);
        if (q == BADADDR) { break; }
        p = q;
    }
    for (i = 0; i < before + after && p != BADADDR; i = i + 1) {
        fprintf(out, "  %08X %s %s\n", p, p == ea ? "→" : " ", GetDisasm(p));
        p = next_head(p, BADADDR);
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/hothelpers.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "被未讀函式大量呼叫的共用常式\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    dump_func(out, "sub_1895D");   // 20 支 T4 呼叫，最高槓桿
    dump_func(out, "sub_1F9B0");   // 4 支
    dump_func(out, "sub_10AAA");
    dump_func(out, "sub_10CAC");
    dump_func(out, "sub_106FD");   // 4 B，疑似 thunk
    dump_func(out, "sub_10337");
    dump_func(out, "sub_106F5");
    dump_func(out, "sub_167CD");   // 財政子項第 0 個，補齊 docs/re/22 §5
    dump_func(out, "sub_16C5E");   // 編成（指令列 #3）的本體

    dump_loose(out, 0x1D51F, 6, 10);
    dump_loose(out, 0x1EF24, 6, 10);
    dump_loose(out, 0x160BD, 6, 10);

    fclose(out);
}
