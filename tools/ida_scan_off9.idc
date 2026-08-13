// IDA Pro 9.4：誰碰位移 +9（DOS/V KI.EXE）。
//
// docs/re/30 §6 的假說：軍團 +0x09（勢力×5）與 +0x08（朝向）相加得大地圖圖塊。
// 「同時出現 +8] 與 +9]」掃出零命中——而那正是相鄰兩欄被當成一個 word
// 讀走時會有的樣子（`mov ax,[si+8]` 之後 al ＝ 朝向、ah ＝ 基底）。
//
// 所以改掃兩件事：
//   A. 所有出現 `+9]` 的指令（誰在用這一欄，不論讀寫）
//   B. 所有 word 大小的 `[reg+8]` 存取（可能一次吃掉 +8 與 +9）
//
// ⚠ 樣式比對不是語意分析：+9／+8 可能屬於任何結構。命中是候選不是結論。
// ⚠ 只讀不寫。輸出 /work/scan_off9.txt。
#include <idc.idc>

static main()
{
    auto out, f, end, p, d, m, n9, n8;
    Wait();
    out = fopen("/work/scan_off9.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "位移 +9 與 word [reg+8] 的使用點\n");
    fprintf(out, "位址空間：IDA DOS/V linear address\n");

    // 正對照：函式總數與一條已知必中的指令。兩者任一為 0 就代表掃描本身壞了，
    // 而不是「沒有人用這一欄」。**沒有正對照的零命中不可信。**
    fprintf(out, "正對照：0001701A 的反組譯 = 「%s」（應含 +9]）\n", GetDisasm(0x1701A));

    fprintf(out, "\n---- A. 出現 `+9]` 的指令 ----\n");
    n9 = 0;
    for (f = get_next_func(0); f != BADADDR; f = get_next_func(f)) {
        end = get_func_attr(f, FUNCATTR_END);
        if (end == BADADDR) { continue; }
        for (p = f; p != BADADDR && p < end; p = next_head(p, end)) {
            d = GetDisasm(p);
            if (strstr(d, "+9]") >= 0) {
                fprintf(out, "  %-12s %08X  %s\n", get_func_name(f), p, d);
                n9 = n9 + 1;
            }
        }
    }
    fprintf(out, "  共 %d 處\n", n9);

    fprintf(out, "\n---- B. word 大小的 [reg+8]（`mov ax/bx/cx/dx/si/di, [..+8]`）----\n");
    n8 = 0;
    for (f = get_next_func(0); f != BADADDR; f = get_next_func(f)) {
        end = get_func_attr(f, FUNCATTR_END);
        if (end == BADADDR) { continue; }
        for (p = f; p != BADADDR && p < end; p = next_head(p, end)) {
            d = GetDisasm(p);
            if (strstr(d, "+8]") < 0) { continue; }
            m = print_insn_mnem(p);
            if (m != "mov" && m != "cmp" && m != "add" && m != "xchg") { continue; }
            // 只留運算元 0 是 16 位暫存器的（byte ptr 的不算）
            if (strstr(d, "byte ptr") >= 0) { continue; }
            if (strstr(d, " ax,") < 0 && strstr(d, " bx,") < 0 && strstr(d, " cx,") < 0 &&
                strstr(d, " dx,") < 0 && strstr(d, " si,") < 0 && strstr(d, " di,") < 0 &&
                strstr(d, " bp,") < 0) { continue; }
            fprintf(out, "  %-12s %08X  %s\n", get_func_name(f), p, d);
            n8 = n8 + 1;
        }
    }
    fprintf(out, "  共 %d 處\n", n8);
    fclose(out);
}
