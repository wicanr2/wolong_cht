// 非破壞性 I/O 與中斷指令索引：用於確認舊程式實際碰到的硬體／TSR 介面。
// 僅輸出 IDA 已辨識的原始指令與原始函式名；不以字串搜尋代替資料流結論。
#include <idc.idc>

static main()
{
    auto out, seg, ea, end, mnem, fn;
    Wait();
    out = fopen("/work/ida-io-ports.txt", "w");
    if (out == 0) {
        return;
    }
    fprintf(out, "IDA Pro 9.4 raw I/O and interrupt index\n");
    fprintf(out, "Address space is the active IDA database linear space.\n");
    for (seg = FirstSeg(); seg != BADADDR; seg = NextSeg(seg)) {
        ea = SegStart(seg);
        end = SegEnd(seg);
        while (ea != BADADDR && ea < end) {
            if (isCode(GetFlags(ea))) {
                mnem = GetMnem(ea);
                if (mnem == "in" || mnem == "out" || mnem == "int") {
                    fn = get_func_name(ea);
                    if (fn == "") {
                        fn = "<no-function>";
                    }
                    fprintf(out, "%08X  %s  %s\n", ea, fn, GetDisasm(ea));
                }
            }
            ea = NextHead(ea, end);
        }
    }
    fclose(out);
}
