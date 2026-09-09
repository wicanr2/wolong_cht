#!/usr/bin/env python3
"""把原版執行期建好的**道路表**從 `peek` 記錄解出來（docs/re/08 §7.1–§7.3）。

    tools/dosgolem.sh out "wait;click:320,200;click:300,151;peek:7789:0000:2048;…"
    tools/py.sh tools/orig_roadtable.py peek.log [--json 輸出.json]
    tools/py.sh tools/orig_roadtable.py --selftest

段位址從 `ipeek:19874:2` 讀（`cs:word_19874`，實測 `0x7789`）。段內版面：

| 區 | 位址 | 一筆 | 內容 |
|---|---|---|---|
| 索引表 | `0x0000`–`0x07FF` | 8 B | 每個節點兩個方向的下一個節點（`+6`／`+8`）|
| 連結記錄 | `0x0800` 起 | 16 B | `+00`／`+02` 第一筆／**最後一筆**路徑點的位址（閉區間）、`+04` 路徑點數（**不含最後那一筆**）、`+06`／`+08` 兩端的節點、`+0A`–`+0E` 外接矩形 |
| 路徑點 | 連結記錄之後 | 4 B | X（word）、Y（byte）、旗標（byte）|

⭐ **路徑點是四方向逐格的**：實測一條 leg 走
`(206,116) → (206,117) → (207,117) → … → (215,117) → (215,118) → …`，
每一步只動一個軸。remake 自己算的路徑會走對角（`straight()`），
所以同一條邊的格子序列不同——**拓樸一樣不等於路一樣**（CLAUDE.md §7 第 16 條）。

旗標的低 3 位是走訪的終止碼（`(al & 7) ≤ 1` 才繼續），
位元 6 標記 leg 的起點。
"""

import json
import re
import sys

SEG = 0x7789
INDEX_END = 0x0800
LINK_SIZE = 16
POINT_SIZE = 4

LINE = re.compile(r"([0-9A-F]{4}):([0-9A-F]{4})\s*=\s*((?:[0-9A-F]{2}\s*)+)")


def parse_peek(text):
    """把 `peek:段:偏移:長度` 的輸出接成一塊連續的 bytes（以段內偏移定址）。"""
    buf = bytearray()
    for line in text.splitlines():
        m = LINE.search(line)
        if not m:
            continue
        off = int(m.group(2), 16)
        data = bytes(int(x, 16) for x in m.group(3).split())
        if off > len(buf):
            raise ValueError(f"段內偏移 {off:#06x} 之前有洞（目前只解出 {len(buf)} byte）")
        buf[off:off + len(data)] = data
    return bytes(buf)


def u16(b, o):
    return b[o] | b[o + 1] << 8


def links(buf):
    """走連結記錄，回傳 [{'a','b','points':[(x,y,flag)…]}]。

    ⚠ **表的長度沒有欄位可查**，只能靠「記錄看起來合不合理」停下來：
    起訖位址要落在段內、路徑點數要對得上 `(end - start) / 4`。
    """
    out = []
    off = INDEX_END
    while off + LINK_SIZE <= len(buf):
        start, end = u16(buf, off), u16(buf, off + 2)
        count = u16(buf, off + 4)
        a, b = u16(buf, off + 6) // 8, u16(buf, off + 8) // 8
        if start == 0 and end == 0 and count == 0:
            break
        if not (INDEX_END <= start < end <= len(buf)) or count == 0:
            break
        # ⚠ **`+02` 指的是最後一筆，不是下一個空槽。** 那一筆的旗標一定是
        # `04`（終端城門格），而 `+04` 的計數**不含它**——`sub_1E81C` 的
        # `inc ah` 在迴圈裡，走到終端就先判斷再離開，最後一筆沒被算進去。
        # 少讀那一筆的症狀是「每條邊都剛好短一格」，而路徑前段完全相同。
        if (end - start) // POINT_SIZE != count:
            break
        pts = []
        for p in range(start, end + POINT_SIZE, POINT_SIZE):
            pts.append((u16(buf, p), buf[p + 2], buf[p + 3]))
        out.append({"a": a, "b": b, "points": pts})
        off += LINK_SIZE
    return out


def four_way(points):
    """回傳「相鄰兩點只動一個軸」的比例——用來確認原版是不是四方向。"""
    ok = tot = 0
    for (x0, y0, _), (x1, y1, _) in zip(points, points[1:]):
        tot += 1
        if (x0 == x1) != (y0 == y1):
            ok += 1
    return ok, tot


def selftest():
    fails = []

    def check(label, cond):
        print(("  ok  " if cond else "  FAIL") + "  " + label)
        if not cond:
            fails.append(label)

    # 一筆連結記錄 ＋ 三個路徑點，手工排出來。
    buf = bytearray(0x820 + 3 * POINT_SIZE)
    start, end = 0x0810, 0x0818
    buf[0x800:0x802] = start.to_bytes(2, "little")
    buf[0x802:0x804] = end.to_bytes(2, "little")
    buf[0x804:0x806] = (2).to_bytes(2, "little")  # 不含最後那一筆
    buf[0x806:0x808] = (5 * 8).to_bytes(2, "little")
    buf[0x808:0x80A] = (9 * 8).to_bytes(2, "little")
    for i, (x, y, f) in enumerate([(206, 116, 0x44), (206, 117, 0), (207, 117, 0x04)]):
        p = start + i * POINT_SIZE
        buf[p:p + 2] = x.to_bytes(2, "little")
        buf[p + 2] = y
        buf[p + 3] = f
    got = links(bytes(buf))
    check("解出一條連結", len(got) == 1)
    check("最後一筆（旗標 04）也要收進來",
          got and got[0]["points"][-1][2] == 0x04)
    check("起訖節點 5 → 9", got and got[0]["a"] == 5 and got[0]["b"] == 9)
    check("三個路徑點且第一個是 (206,116)",
          got and len(got[0]["points"]) == 3 and got[0]["points"][0][:2] == (206, 116))
    check("四方向判定：三點兩段都只動一個軸", four_way(got[0]["points"]) == (2, 2))

    # 負對照一：路徑點數對不上就要停，不能默默照收。
    bad = bytearray(buf)
    bad[0x804:0x806] = (99).to_bytes(2, "little")
    check("負對照：路徑點數對不上就停", links(bytes(bad)) == [])

    # 負對照二：起訖位址落在索引表裡是壞資料。
    bad = bytearray(buf)
    bad[0x800:0x802] = (0x0004).to_bytes(2, "little")
    check("負對照：起訖位址落在索引表就停", links(bytes(bad)) == [])

    # 負對照三：對角線的點要被判成「不是四方向」。
    check("負對照：對角線不算四方向",
          four_way([(1, 1, 0), (2, 2, 0)]) == (0, 1))

    # 正對照：peek 版面解得回來。
    text = "   7789:0000 = 01 02 03 04\n   7789:0004 = 05 06\n"
    check("peek 版面解得回來", parse_peek(text) == bytes([1, 2, 3, 4, 5, 6]))
    return 1 if fails else 0


def main():
    args = sys.argv[1:]
    if args[:1] == ["--selftest"]:
        return selftest()
    if not args:
        print(__doc__)
        return 2
    out_json = None
    if "--json" in args:
        i = args.index("--json")
        out_json = args[i + 1]
        del args[i:i + 2]
    buf = parse_peek(open(args[0], encoding="utf-8", errors="replace").read())
    print(f"解出段內 {len(buf)} byte")
    ls = links(buf)
    pts = sum(len(l["points"]) for l in ls)
    print(f"連結記錄 {len(ls)} 筆，路徑點 {pts} 個")
    ok = tot = 0
    for l in ls:
        a, b = four_way(l["points"])
        ok += a
        tot += b
    if tot:
        print(f"相鄰兩點只動一個軸：{ok}/{tot}（{ok * 100.0 / tot:.1f}%）")
    if out_json:
        json.dump([{"a": l["a"], "b": l["b"],
                    "points": [[p[0], p[1], p[2]] for p in l["points"]]} for l in ls],
                  open(out_json, "w"), separators=(",", ":"))
        print(f"寫到 {out_json}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
