// IDA Pro 9.4：軍團編成畫面（DOS/V KI.EXE）。
//
// sub_16C92 是編成畫面的主迴圈（狀態列 TALK #1「請下達各部隊編成之指示。」，
// 完成時走 cx=0x19B 的八變體「編成完成」）。這一支把它的上下游一次讀完：
//   sub_16C5E   進入點
//   sub_16D56 / sub_16DFD / sub_16FD2 / sub_16D6F / sub_16DA8 / sub_16E80  畫面各區
//   sub_14717   熱區編號 → 槽索引（四個呼叫者，語意未定）
//   sub_16F26   實際編成（已被多份文件提及，這裡取全文對照）
//
// ⚠ 只讀不寫：不 add_func、不改名、不加型別、不寫註解。
// 輸出 /work/corps_form.txt。
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
        for (n = 0; n < 160 && p != BADADDR; n = n + 1) {
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

static main()
{
    auto out;
    Wait();
    out = fopen("/work/corps_form.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "軍團編成畫面\n位址空間：IDA DOS/V linear address\n");

    dump_func(out, 0x16C5E, "sub_16C5E 進入點");
    dump_func(out, 0x14717, "sub_14717 熱區→槽索引？");
    dump_func(out, 0x16D56, "sub_16D56");
    dump_func(out, 0x16DFD, "sub_16DFD");
    dump_func(out, 0x16FD2, "sub_16FD2");
    dump_func(out, 0x16D6F, "sub_16D6F 印數字");
    dump_func(out, 0x16DA8, "sub_16DA8 印數字");
    dump_func(out, 0x16E80, "sub_16E80");
    dump_func(out, 0x16F26, "sub_16F26 實際編成");

    fclose(out);
}
