// IDA Pro 9.4：一覽表的 23 個 callback（DOS/V KI.EXE）。
//
// 這些位址只被間接參考（呼叫端把它們當立即值傳給 sub_1820E，引擎再
// `call cs:word_181A6`），所以 IDA 的自動分析看不到任何 code xref，
// 一個函式都沒建。位址本身是確定的——來源見 docs/re/26 §4。
//
// ⚠ 這支會**改寫資料庫**（add_func）。所以一定要用 `tools/ida.sh script`，
// 它跑在 .i64 的唯讀副本上；原始資料庫的雜湊不受影響。
// 建函式只是讓 IDA 願意解析與追蹤流程，**不改名、不加型別、不寫註解**——
// 名字會把推論等級烙進之後每一次 dump。
//
// 輸出 /work/listcb.txt。
#include <idc.idc>

static CB(o) { return 0x10000 + o; }

static ensure_func(out, ea)
{
    auto s;
    s = get_func_attr(ea, FUNCATTR_START);
    if (s == ea) { return 1; }
    if (add_func(ea) != 0) { return 1; }
    // add_func 失敗多半是那一段還沒被當成程式碼——這些位址只有間接參考，
    // 自動分析沒有理由去解析它們。先清掉既有定義再強制解一條指令。
    del_items(ea, DELIT_EXPAND, 1);
    if (create_insn(ea) == 0) {
        fprintf(out, "  ⚠ %08X 連一條指令都解不出來\n", ea);
        return 0;
    }
    if (add_func(ea) != 0) { return 1; }
    fprintf(out, "  ⚠ %08X 解得出指令但建不了函式\n", ea);
    return 0;
}

static dump_func(out, ea, label)
{
    auto start, end, p, x, n;
    fprintf(out, "\n==== %s  %08X ====\n", label, ea);
    start = get_func_attr(ea, FUNCATTR_START);
    end = get_func_attr(ea, FUNCATTR_END);
    if (start == BADADDR || start != ea) {
        // 建不了函式（被鄰接函式吸收，或尾端不是乾淨的 retn）時，
        // 改成從該位址線性反組譯到第一個 retn。**這是線性走訪不是控制流分析**，
        // 中途的條件跳躍會跳過的區段也會被印出來，讀的時候要自己判斷。
        fprintf(out, "  ⚠ 無獨立函式邊界（所屬：%s），以下為線性反組譯至第一個 retn\n",
                get_func_name(ea));
        p = ea;
        for (n = 0; n < 80 && p != BADADDR; n = n + 1) {
            fprintf(out, "  %08X  %s\n", p, GetDisasm(p));
            if (print_insn_mnem(p) == "retn" || print_insn_mnem(p) == "retf") { break; }
            p = next_head(p, BADADDR);
        }
        return;
    }
    n = 0;
    fprintf(out, "直接呼叫者：");
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
    auto out, i, o, offs, labels, n;
    Wait();
    out = fopen("/work/listcb.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "一覽表 callback（建函式後）\n");
    fprintf(out, "位址空間：IDA DOS/V linear address；括號內是傳進 sub_1820E 的段內偏移\n");
    fprintf(out, "⚠ 本次在資料庫副本上執行 add_func，原始 .i64 未動\n");

    offs = object();
    // 依 docs/re/26 §4 的家族順序。註記寫在 label 裡，不寫進資料庫。
    auto o_list, l_list;
    o_list = "70EC,713D,71A8,7217,724D,"
             "7378,73CE,743B,745F,"
             "7550,75C8,763C,76A0,76DC,"
             "77EE,7875,78E5,7944,796C,"
             "7AC6,7B12,7B6F,7B90";
    l_list = "軍團-si,軍團-di,軍團-ax(1716D),軍團-ax(171D3),軍團-bx畫列,"
             "據點-si,據點-di,據點-ax(17400),據點-bx畫列,"
             "武將-si,武將-di,武將-ax(175FA),武將-ax(17663),武將-bx畫列,"
             "勢力-si,勢力-di,勢力-ax(178A7),勢力-ax(17906),勢力-bx畫列,"
             "開局-si,開局-di,開局-ax(17B3C),開局-bx畫列";

    // 先全部建函式，再全部 dump——後建的函式可能影響先前的邊界推導。
    n = 0;
    for (i = 0; i < 23; i = i + 1) {
        o = xtol(substr(o_list, i * 5, i * 5 + 4));
        if (ensure_func(out, CB(o))) { n = n + 1; }
    }
    Wait();
    fprintf(out, "\n建立／確認函式 %d / 23\n", n);

    auto pos, nextpos, lab;
    pos = 0;
    for (i = 0; i < 23; i = i + 1) {
        o = xtol(substr(o_list, i * 5, i * 5 + 4));
        nextpos = strstr(substr(l_list, pos, -1), ",");
        if (nextpos < 0) { lab = substr(l_list, pos, -1); }
        else { lab = substr(l_list, pos, pos + nextpos); pos = pos + nextpos + 1; }
        dump_func(out, CB(o), lab);
    }
    fclose(out);
}
