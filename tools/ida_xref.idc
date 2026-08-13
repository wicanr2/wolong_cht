// 查 IDA 資料庫的交叉參考圖。
//
//   tools/ida.sh raw dosv idat -A "-S/work/tools/ida_xref.idc aKaogrfDat" KI.EXE.i64
//   cat workplace/ida/dosv/xref-aKaogrfDat.txt
//
// 為什麼不 grep .asm：.asm 是攤平的文字，沒有交叉參考圖。想知道
// 「誰用了這個東西」，grep 只能從呼叫端的參數順序反推，那是間接證據。
//
// 三件在這個 image 上實測過的事（見 ~/.claude/knowledge-base/retro/ida-pro-9.4.md）：
//   1. **IDAPython 可用**，但要 ida-pro-9.4-idapython:py312-v1；
//      基底 image 跑 .py 是零輸出的靜默失敗。tools/ida.sh 依副檔名自動選。
//      新腳本優先寫 IDAPython（骨架見 tools/ida_scan.py），IDC 只當退路——
//      IDC 缺一半內建函式（get_func_qty 就不存在），而缺函式在 headless
//      底下是靜默中止：已寫的輸出一起消失，exit code 仍是 0。
//   2. headless 的 print / Message() 不進 stdout，一律 fopen 寫檔。
//      不寫檔就等於沒跑，而且 exit code 還是 0。
//   3. 讀寫判定用 XrefType()，不要比對助憶碼字串。

#include <idc.idc>

static main() {
    auto name, ea, x, t, kind, out, path, n;

    Wait();                                  // 一定要等自動分析跑完

    name = ARGV[1];
    if (name == "") {
        name = "start";
    }
    path = "/work/xref-" + name + ".txt";
    out = fopen(path, "w");

    ea = get_name_ea_simple(name);
    if (ea == BADADDR) {
        fprintf(out, "找不到符號: %s\n", name);
        fclose(out);
        return;
    }
    fprintf(out, "符號 %s @ 0x%X\n", name, ea);

    n = 0;
    fprintf(out, "\n--- 資料參考（誰讀寫這個位址）---\n");
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        t = XrefType();
        if (t == 2) {
            kind = "寫";
        } else if (t == 1) {
            kind = "取址";
        } else {
            kind = "讀";
        }
        fprintf(out, "%-4s 0x%X  %-24s  在 %s\n",
                kind, x, GetDisasm(x), get_func_name(x));
        n = n + 1;
    }
    fprintf(out, "（共 %d 筆）\n", n);

    n = 0;
    fprintf(out, "\n--- 程式碼參考（誰跳到／呼叫這裡）---\n");
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "0x%X  %-24s  在 %s\n",
                x, GetDisasm(x), get_func_name(x));
        n = n + 1;
    }
    fprintf(out, "（共 %d 筆）\n", n);

    fclose(out);
}
