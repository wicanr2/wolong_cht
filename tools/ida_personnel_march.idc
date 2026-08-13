// IDA Pro 9.4：人事四支、行軍指示鏈與頂層模式分派（DOS/V KI.EXE）。
//
// 挑這一組的理由（docs/re/24 §4 的順序）：
//   - 人事四支已被 docs/re/22 §3.2 依 TALK #78 命名，但**函式本體沒讀過**。
//     外交官是停戰與協力的前置條件，解開就閉合一個既有迴路。
//   - 行軍指示是軍團指令的另一半，狀態列訊息 #3／#21 已指出用途。
//   - off_159D2 的 [1][2][3] 是頂層模式分派表裡僅有的三個未讀槽位。
//
// 輸出 /work/personnel.txt。不改名、不加型別、不寫回資料庫。
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

static dump_at(out, ea, count)
{
    auto p, i, x;
    fprintf(out, "\n==== 位址 %08X（無函式邊界時仍可讀）====\n", ea);
    fprintf(out, "  參考它的：");
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "%08X(%s) ", x, get_func_name(x));
    }
    fprintf(out, "\n");
    p = ea;
    for (i = 0; i < count && p != BADADDR; i = i + 1) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
        p = next_head(p, BADADDR);
    }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/personnel.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "人事、行軍指示與頂層模式分派\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    // 人事：funcs_16279 的四筆
    dump_func(out, "sub_16A9B");   // 內政官任命
    dump_func(out, "sub_16B08");   // 內政官解任
    dump_func(out, "sub_16B71");   // 外交官任命
    dump_func(out, "sub_16BE3");   // 外交官解任
    dump_func(out, "sub_16C92");   // 三支共用；也是編成的內層

    // 行軍指示：sub_1628F 的 al=1 分支
    dump_func(out, "sub_17F90");
    dump_func(out, "sub_17FDB");

    // 頂層模式分派 off_159D2 的未讀槽位
    dump_at(out, 0x1614A, 30);
    dump_at(out, 0x15E1E, 12);
    dump_at(out, 0x15A3A, 45);

    fclose(out);
}
