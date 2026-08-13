// IDA Pro 9.4：三支一覽表與它們共用的視窗常式（DOS/V KI.EXE）。
//
// 挑這一組的理由：`sub_17400`（選據點）、`sub_17906`（選勢力）、
// `sub_17663`（選武將）是戰略層每個指令都要走的入口，而
// `sub_181C0`（開視窗）與 `sub_1820E`（一覽表選取）各有九個呼叫點卻
// 一份文件都沒提過。解開一次，所有走一覽表的指令都跟著有了。
//
// 一併讀 `sub_11D46`（17 個呼叫點，人事四支離開前都呼叫）與
// 解任的兩支實作 `sub_16B4F`／`sub_16C2A`。
//
// 輸出 /work/listwin.txt。不改名、不加型別、不寫回資料庫。
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
        fprintf(out, "%s(%08X) ", get_func_name(x), x);
        n = n + 1;
    }
    fprintf(out, " 共 %d 處，bytes=%u\n", n, end - start);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/listwin.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "一覽表與視窗常式\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    dump_func(out, "sub_17400");   // 選據點
    dump_func(out, "sub_17906");   // 選勢力
    dump_func(out, "sub_17663");   // 選武將
    dump_func(out, "sub_181C0");   // 開視窗（9 個呼叫點，未記錄）
    dump_func(out, "sub_1820E");   // 一覽表選取（9 個呼叫點）
    dump_func(out, "sub_11D46");   // 17 個呼叫點
    dump_func(out, "sub_16B4F");   // 內政官解任
    dump_func(out, "sub_16C2A");   // 外交官解任

    fclose(out);
}
