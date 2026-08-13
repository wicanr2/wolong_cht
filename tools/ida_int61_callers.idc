// INT 61h 呼叫端與 TSR handler 原始位元組的非破壞性 IDA Pro 9.4 匯出器。
// 適用於松崗 DOS/V KI.EXE 與 YNSOUND.COM；只讀資料庫，不建立函式、
// 不改名、不加型別、不寫回輸入。
#include <idc.idc>

static dump_function(out, start)
{
    auto end, p;
    end = get_func_attr(start, FUNCATTR_END);
    fprintf(out, "\n=== FUNCTION %s @ %08X ===\n", get_func_name(start), start);
    fprintf(out, "RANGE %08X-%08X\n", start, end);
    for (p = start; p != BADADDR && p < end; p = next_head(p, end)) {
        fprintf(out, "%08X  %s\n", p, GetDisasm(p));
    }
}

static dump_bytes(out, label, start, count)
{
    auto p, end, q;
    fprintf(out, "\n=== RAW BYTES %s @ %08X count=%d ===\n", label, start, count);
    end = start + count;
    for (p = start; p < end; p = p + 16) {
        fprintf(out, "%08X", p);
        for (q = 0; q < 16 && p + q < end; q = q + 1) {
            fprintf(out, " %02X", get_wide_byte(p + q));
        }
        fprintf(out, "\n");
    }
}

static main()
{
    auto out, seg, ea, end, start, last;
    Wait();
    out = fopen("/work/ida-int61-callers.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "IDA Pro 9.4 non-destructive INT 61h export\n");
    fprintf(out, "Address space: active IDA database linear addresses\n");
    fprintf(out, "Raw instructions/xrefs only; no semantic renaming is performed.\n");
    dump_bytes(out, "linear 00010103", 0x10103, 0x80);
    last = BADADDR;
    for (seg = FirstSeg(); seg != BADADDR; seg = NextSeg(seg)) {
        ea = SegStart(seg);
        end = SegEnd(seg);
        while (ea != BADADDR && ea < end) {
            if (isCode(GetFlags(ea)) && GetMnem(ea) == "int" && get_operand_value(ea, 0) == 0x61) {
                fprintf(out, "\n=== INT 61h @ %08X ===\n%s\n", ea, GetDisasm(ea));
                start = get_func_attr(ea, FUNCATTR_START);
                if (start != BADADDR && start != last) {
                    dump_function(out, start);
                    last = start;
                }
            }
            ea = NextHead(ea, end);
        }
    }
    fclose(out);
}
