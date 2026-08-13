// 非破壞性匯出指定位址的原始位元組；只讀 .i64，不改名、不改資料庫。
#include <idc.idc>

static main() {
    auto start, count, out, i;
    Wait();
    start = xtol(ARGV[1]);
    count = xtol(ARGV[2]);
    out = fopen("/work/bytes-dump.txt", "w");
    if (out == 0) { return; }
    fprintf(out, "start=0x%X count=%d\n", start, count);
    for (i = 0; i < count; i = i + 1) {
        fprintf(out, "0x%X %02X\n", start + i, Byte(start + i));
    }
    fclose(out);
}
