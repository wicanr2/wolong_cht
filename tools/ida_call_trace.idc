// IDA Pro 9.4：全函式的「立即值 → 呼叫」軌跡（DOS/V KI.EXE）。
//
// 為什麼要這個：355 支函式沒被任何文件提過，逐支讀組語太慢。
// 但這支程式有一個很好用的性質——UI 與流程函式幾乎都會走
// 「設常數 → 呼叫共用常式」的形式，而那些常數多半是 TALK 訊息索引：
//
//     mov cx, 16h / call sub_18853     → 狀態列顯示訊息 #22
//     mov ax, 5 / mov cx, 4Dh / call sub_193E9  → 5 項選單，選項文字在 #77
//     mov al, 93h / call sub_18810     → 顯示訊息 #147
//
// TALK.DAT 已全解，所以把這些常數抽出來查表，就能用**遊戲自己的文字**
// 判斷一支函式在做什麼——比從呼叫端的參數順序反推可靠得多。
//
// 這支只輸出事實，不做配對也不猜語意；配對與分類交給 tools/re_classify.py。
//
// 每行一筆，tab 分隔：
//   FN    start end name                          函式起點
//   IMM   start ea reg value                      mov <reg>, <立即值>
//   CALL  start ea target targetname              直接呼叫
//   ICALL start ea disasm                         間接呼叫
//   INT   start ea number                         軟體中斷
//   PORT  start ea mnem operands                  in / out
//   STR   start ea target text                    指到字串的資料參考
//
// 輸出 /work/calltrace.tsv。不改名、不加型別、不寫回資料庫。
#include <idc.idc>

static main()
{
    auto out, ea, end, p, m, t, c, ct, x, v, s, o0;
    Wait();
    out = fopen("/work/calltrace.tsv", "w");
    if (out == 0) { return; }
    fprintf(out, "# 立即值與呼叫軌跡\n");
    fprintf(out, "# 位址空間：IDA DOS/V linear address\n");

    for (ea = get_next_func(0); ea != BADADDR; ea = get_next_func(ea)) {
        end = get_func_attr(ea, FUNCATTR_END);
        if (end == BADADDR) { continue; }
        fprintf(out, "FN\t%08X\t%08X\t%s\n", ea, end, get_func_name(ea));
        for (p = ea; p != BADADDR && p < end; p = next_head(p, end)) {
            m = print_insn_mnem(p);

            if (m == "mov" && get_operand_type(p, 1) == 5) {  // o_imm
                o0 = print_operand(p, 0);
                // 只記暫存器目標；記憶體目標的立即值另有語意，這裡不混。
                if (get_operand_type(p, 0) == 1) {
                    fprintf(out, "IMM\t%08X\t%08X\t%s\t%d\n",
                            ea, p, o0, get_operand_value(p, 1));
                }
            }
            else if (m == "call") {
                t = get_operand_type(p, 0);
                if (t == 6 || t == 7) {
                    for (c = get_first_cref_from(p); c != BADADDR;
                         c = get_next_cref_from(p, c)) {
                        ct = XrefType();
                        // fl_CN=16 / fl_CF=17 才是呼叫邊；fl_F=21 是循序落下。
                        if (ct == 16 || ct == 17) {
                            fprintf(out, "CALL\t%08X\t%08X\t%08X\t%s\n",
                                    ea, p, c, get_func_name(c));
                        }
                    }
                } else {
                    fprintf(out, "ICALL\t%08X\t%08X\t%s\n", ea, p, GetDisasm(p));
                }
            }
            else if (m == "int") {
                fprintf(out, "INT\t%08X\t%08X\t%d\n", ea, p, get_operand_value(p, 0));
            }
            else if (m == "in" || m == "out") {
                fprintf(out, "PORT\t%08X\t%08X\t%s\t%s\n", ea, p, m, GetDisasm(p));
            }

            // 指到字串的資料參考：字串是最直接的語意線索。
            for (x = get_first_dref_from(p); x != BADADDR;
                 x = get_next_dref_from(p, x)) {
                s = get_strlit_contents(x, -1, 0);
                if (s != "") {
                    fprintf(out, "STR\t%08X\t%08X\t%08X\t%s\n", ea, p, x, s);
                }
            }
        }
    }
    fclose(out);
}
