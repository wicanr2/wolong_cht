// IDA Pro 9.4：YNSOUND.COM 的 command 表與資源 layout parser（松崗 DOS/V）。
//
// docs/re/17 已定案 INT 61h 介面：handler 在 COM offset 0x0103，
// 以 AH 當 command index 查 COM offset 0x0115 的 word table。
// 該文件把 command 0x07（COM offset 0x01F6）標為「資源 layout parser，格式未知」，
// 並列為下一步。這支腳本把 command 表與 parser 本體攤開。
//
// COM 檔在 IDA 裡的 segment base 是 0x10000，COM offset 0xNNNN → linear 0x1NNNN。
// 表項存的是 COM offset。
//
// 輸出 /work/ynsound.txt。不改名、不加型別、不寫回資料庫。
#include <idc.idc>

static COM(off) { return 0x10000 + off; }

static dump_range(out, label, start, end)
{
    auto p;
    fprintf(out, "\n---- %s %08X-%08X ----\n", label, start, end);
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
        if (i % 16 == 15) { fprintf(out, "  %08X  %s\n", start + i - 15, line); line = ""; }
    }
    if (line != "") { fprintf(out, "  %s\n", line); }
}

static main()
{
    auto out, i, w, ea, n;
    Wait();
    out = fopen("/work/ynsound.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "YNSOUND.COM：INT 61h command 表與資源 parser\n");
    fprintf(out, "SHA-256 e2c6a6a8576c4f2a96b7e3f156d7f48c9570ae03539fe9367adb78aebb364fa1\n");
    fprintf(out, "位址空間：IDA linear，COM offset + 0x10000\n");

    // 資料庫整體：COM 只有 3,463 B，全部列出函式邊界比較快
    fprintf(out, "\n==== 資料庫函式 ====\n");
    n = 0;
    for (ea = get_next_func(0); ea != BADADDR; ea = get_next_func(ea)) {
        fprintf(out, "  %08X-%08X  %s\n", ea,
                get_func_attr(ea, FUNCATTR_END), get_func_name(ea));
        n = n + 1;
    }
    fprintf(out, "  共 %d 支\n", n);

    // INT 61h handler
    dump_range(out, "INT 61h handler（COM 0x0103）", COM(0x0103), COM(0x0115));

    // command word table（COM 0x0115），讀 24 筆
    fprintf(out, "\n==== command 表（COM 0x0115）====\n");
    for (i = 0; i < 24; i = i + 1) {
        w = get_wide_word(COM(0x0115) + i * 2);
        fprintf(out, "  AH=%02X  word=%04X -> %08X %s\n",
                i, w, COM(w), get_func_name(COM(w)));
    }

    // command 0x07：docs/re/17 標為資源 layout parser
    dump_range(out, "command 0x07 資源 parser（COM 0x01F6）", COM(0x01F6), COM(0x02C0));
    dump_bytes(out, "COM 0x01F6 起", COM(0x01F6), 208);

    fclose(out);
}
