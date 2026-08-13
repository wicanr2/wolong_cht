// 把一支函式連同「每個位址是什麼、誰還在用它」一起 dump 出來。
//
//   tools/ida.sh raw dosv idat -A "-S/work/tools/ida_func.idc sub_12E89" KI.EXE.i64
//   cat workplace/ida/dosv/func-sub_12E89.txt
//
// 為什麼不 grep .asm：`.asm` 是攤平的文字，**沒有交叉參考圖**。
// 讀組語時手上拿到的是位址（`[bx+23h]`、`word_1987C`），不是語意，
// 而「誰還在用它」正是判斷語意最便宜的證據——grep 只能從呼叫端的
// 參數順序反推，那是間接證據，會推錯。
//
// 這支把兩件事併在同一次輸出裡：
//   1. IDA 解析過的指令（運算元已經翻成名字）
//   2. 每個被參考到的資料位址：名字 ＋ **所有讀它／寫它的地方**
// 這樣「這個欄位是誰寫的」不必另外查一次。
//
// 三件在這個 image 上實測過的事（~/.claude/knowledge-base/retro/ida-pro-9.4.md）：
//   1. **IDAPython 可用**，但要 ida-pro-9.4-idapython:py312-v1；
//      基底 image 跑 .py 是零輸出的靜默失敗。tools/ida.sh 依副檔名自動選。
//      新腳本優先寫 IDAPython（骨架見 tools/ida_scan.py），IDC 只當退路——
//      IDC 缺一半內建函式（get_func_qty 就不存在），而缺函式在 headless
//      底下是靜默中止：已寫的輸出一起消失，exit code 仍是 0。
//   2. headless 的 print / Message() **不進 stdout**，一律 fopen 寫檔。
//      不寫檔就等於沒跑，而且 exit code 還是 0。
//   3. 讀寫判定用 XrefType()（dr_W=2／dr_R=3／dr_O=1），
//      不要比對助憶碼字串，也不要看 print_operand。

#include <idc.idc>

// 把一個位址的所有 data xref 印出來，標明讀／寫／取位址。
static dump_refs(out, ea, indent) {
    auto x, t, kind, fn, n;
    n = 0;
    for (x = get_first_dref_to(ea); x != BADADDR; x = get_next_dref_to(ea, x)) {
        t = XrefType();
        kind = "?";
        if (t == 1) { kind = "取位址"; }
        if (t == 2) { kind = "寫"; }
        if (t == 3) { kind = "讀"; }
        fn = get_func_name(x);
        if (fn == "") { fn = "（不在函式內）"; }
        fprintf(out, "%s%s  %s  %08X  %s\n", indent, kind, fn, x, GetDisasm(x));
        n = n + 1;
        // 太多的話截斷——列出前 24 筆已經足夠判斷語意，
        // 全列會把輸出淹掉（`word_10D52` 這種段暫存器有上百筆）。
        if (n >= 24) {
            fprintf(out, "%s…（還有更多，已截斷於 24 筆）\n", indent);
            break;
        }
    }
    if (n == 0) {
        fprintf(out, "%s（沒有 data xref —— ⚠ 立即值形式的參考 IDA 不建 xref，\n", indent);
        fprintf(out, "%s  回零筆不等於沒人用，見 docs/re/03 §0）\n", indent);
    }
}

static main() {
    auto name, ea, end, out, path, dis, x, seen, i, target, tname;

    Wait();                                   // 一定要等自動分析跑完

    name = ARGV[1];
    if (name == "") { name = "start"; }
    path = "/work/func-" + name + ".txt";
    out = fopen(path, "w");
    if (out == 0) { return; }

    ea = get_name_ea_simple(name);
    if (ea == BADADDR) {
        fprintf(out, "找不到符號 %s\n", name);
        fclose(out);
        return;
    }

    end = get_func_attr(ea, FUNCATTR_END);
    if (end == BADADDR) { end = ea + 0x200; }   // 不是函式就往後看一段

    fprintf(out, "=== %s  %08X–%08X ===\n\n", name, ea, end);

    // ① 誰呼叫這一支
    fprintf(out, "--- 呼叫者 ---\n");
    i = 0;
    for (x = get_first_cref_to(ea); x != BADADDR; x = get_next_cref_to(ea, x)) {
        fprintf(out, "  %s  %08X  %s\n", get_func_name(x), x, GetDisasm(x));
        i = i + 1;
    }
    if (i == 0) { fprintf(out, "  （沒有 code xref）\n"); }

    // ② 本體
    fprintf(out, "\n--- 反組譯 ---\n");
    seen = "";
    while (ea != BADADDR && ea < end) {
        dis = GetDisasm(ea);
        fprintf(out, "%08X  %s\n", ea, dis);
        // 收集這一行參考到的資料位址，等一下一次把 xref 印出來，
        // **不要邊走邊印**——那會讓反組譯本體被 xref 淹沒讀不下去。
        for (x = get_first_dref_from(ea); x != BADADDR; x = get_next_dref_from(ea, x)) {
            tname = get_name(x);
            if (tname != "" && strstr(seen, "|" + tname + "|") < 0) {
                seen = seen + "|" + tname + "|";
            }
        }
        ea = next_head(ea, end);
    }

    // ③ 這支碰到的每個資料位址：誰還在用它
    fprintf(out, "\n--- 這支碰到的資料，以及誰還在用 ---\n");
    i = 0;
    while (1) {
        auto p, q;
        p = strstr(seen, "|");
        if (p < 0) { break; }
        seen = substr(seen, p + 1, -1);
        q = strstr(seen, "|");
        if (q < 0) { break; }
        tname = substr(seen, 0, q);
        seen = substr(seen, q + 1, -1);
        target = get_name_ea_simple(tname);
        if (target != BADADDR) {
            fprintf(out, "\n%s  (%08X)\n", tname, target);
            dump_refs(out, target, "    ");
        }
        i = i + 1;
        if (i > 40) { break; }
    }

    fclose(out);
}
