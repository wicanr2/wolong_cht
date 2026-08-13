// IDA Pro 9.4：把間接分派跳表的內容攤開（DOS/V KI.EXE）。
//
// 普查（tools/re_coverage.py）指出 30 個間接呼叫點是靜態分析的斷點：
// IDA 看得到 `call cs:funcs_XXXX[bx]`，卻不知道 runtime 選哪一支。
// 但表本身是靜態資料——把表讀出來，候選集合就從「未知」變成「有界的清單」。
//
// 這不等於知道 runtime 走哪條；它把問題從「不知道有哪些可能」
// 降級成「知道有哪些可能，還要驗哪一條會走到」。兩者的推論等級不同。
//
// 輸出 /work/dispatch.txt。不改名、不加型別、不寫回資料庫。
#include <idc.idc>

static seg_base()
{
    // DOS/V KI.EXE 在 IDA 裡的 segment base；表項存的是 segment offset。
    return 0x10000;
}

static dump_table(out, label, ea, n)
{
    auto i, w, target, fname, fstart;
    fprintf(out, "\nTABLE %s EA=%08X 讀 %d 筆\n", label, ea, n);
    if (ea == BADADDR) { fprintf(out, "  未解析\n"); return; }
    for (i = 0; i < n; i = i + 1) {
        w = get_wide_word(ea + i * 2);
        target = seg_base() + w;
        fname = get_func_name(target);
        fstart = get_func_attr(target, FUNCATTR_START);
        fprintf(out, "  [%2d] +%04X word=%04X -> %08X %s%s\n",
                i, i * 2, w, target,
                fname == "" ? "(非函式起點)" : fname,
                (fstart != BADADDR && fstart != target) ? " ※落在函式中段" : "");
    }
}

static dump_func(out, name)
{
    auto ea, start, end, p, x;
    ea = get_name_ea_simple(name);
    fprintf(out, "\n==== FUNC %s EA=%08X ====\n", name, ea);
    if (ea == BADADDR) { fprintf(out, "  找不到符號\n"); return; }
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR) { fprintf(out, "  沒有函式邊界\n"); return; }
    fprintf(out, "呼叫者：");
    for (x = get_first_cref_to(start); x != BADADDR; x = get_next_cref_to(start, x)) {
        fprintf(out, "%s(%08X) ", get_func_name(x), x);
    }
    fprintf(out, "\n");
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
    }
}

static dump_bytes(out, label, start, n)
{
    auto i, line, b;
    fprintf(out, "\nBYTES %s %08X 共 %d\n", label, start, n);
    line = "";
    for (i = 0; i < n; i = i + 1) {
        b = get_wide_byte(start + i);
        line = line + sprintf("%02X ", b);
        if (i % 16 == 15) {
            fprintf(out, "  %08X  %s\n", start + i - 15, line);
            line = "";
        }
    }
    if (line != "") { fprintf(out, "  %s\n", line); }
}

static main()
{
    auto out;
    Wait();
    out = fopen("/work/dispatch.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "DOS/V KI.EXE 間接分派表內容\n");
    fprintf(out, "位址空間：IDA DOS/V linear address，segment base 0x10000\n");
    fprintf(out, "表項是 16-bit segment offset，加 base 後為 linear\n");

    // === 月結／經濟／AI 模組的四張表：普查標為 T4，全專案沒有文件提過 ===
    dump_func(out, "sub_161CA");
    dump_table(out, "funcs_161FE", get_name_ea_simple("funcs_161FE"), 16);
    dump_func(out, "sub_16224");
    dump_table(out, "funcs_16255", get_name_ea_simple("funcs_16255"), 16);
    dump_func(out, "sub_16265");
    dump_table(out, "funcs_16279", get_name_ea_simple("funcs_16279"), 16);
    dump_func(out, "sub_1678D");
    dump_table(out, "funcs_167AE", get_name_ea_simple("funcs_167AE"), 16);

    // === 戰略 UI 模組：兩張表加一組共用的函式指標 ===
    dump_func(out, "sub_18412");
    dump_table(out, "funcs_18450", get_name_ea_simple("funcs_18450"), 16);
    dump_func(out, "sub_18FC9");
    dump_table(out, "funcs_19037", get_name_ea_simple("funcs_19037"), 16);

    // word_181A6／word_181A8 被四支 T4 函式當函式指標呼叫。
    dump_bytes(out, "word_181A6 附近", 0x181A0, 32);

    // === 未解的 runtime 分派：sub_159A6／sub_159B7 用 bx 相對定址 ===
    dump_func(out, "sub_159A6");
    dump_func(out, "sub_159B7");
    dump_table(out, "funcs_159C0", get_name_ea_simple("funcs_159C0"), 16);
    dump_bytes(out, "off_159D2 附近", 0x159C0, 48);

    // === sub_15FAA 的 [bx+6056h]、sub_1E9C1 的 [bx-15F3h] ===
    dump_func(out, "sub_15FAA");
    dump_bytes(out, "0x16056 附近", 0x16056, 48);

    // === loc_1A065：自我修改碼，audit 20 §2.2 標為沒有可靠函式邊界 ===
    dump_func(out, "sub_19FA0");
    dump_bytes(out, "loc_1A065 自我修改區", 0x1A065, 256);

    // === 唯一的孤兒函式 ===
    dump_func(out, "sub_14A0F");

    fclose(out);
}
