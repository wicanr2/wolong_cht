#!/usr/bin/env python3
"""把 dosgolem 的 `ipeek` 輸出拼成一份 SAVE.DAT，讓 remake 從**原版的那一刻**起跑。

    tools/dosgolem.sh out "…;ipeek:10CF0:2048;ipeek:114F0:2048;…"   ← 11 段
    tools/py.sh tools/orig_snapshot.py peek.log 來源SAVE.DAT 輸出SAVE.DAT [--slot N]
    tools/py.sh tools/orig_snapshot.py --selftest

⭐ **為什麼需要這個。** 對拍的地基是「兩邊從同一個局面起跑」，而原版載入
存檔之後會自己跑掉一段（實測 27 個子刻：時鐘從 10 時走到 13 時，
期間據點游標沒動、每小時世界更新跑了三次）。**那一段沒辦法在 remake 重放**
——改時鐘、補 skip、調游標都只是換一種岔法（`docs/playtest/119` §7）。
直接把原版當下的資料區塊搬過來，起點就一致了。

記憶體佈局與存檔區塊一一對應，位移 `+0x0CF0`（`docs/re/06` §0）：

    IDA linear = 區塊偏移 + 0x0CF0 + 0x10000

⚠ **只換已解的區塊內容，未解區域原封不動**——這是「改寫不是重建」
（`CLAUDE.md` §9）。來源存檔提供其餘三個槽與檔案結構。
"""

import re
import sys

BLOCK = 22208
BLOCKS = 4
FILE_SIZE = BLOCK * BLOCKS
CHUNK = 2048

# ⚠ **區塊在記憶體裡不是一整塊。** `docs/formats/08` 只保證**前 59 byte**
# 原封不動載入 `cs:0CF0h`；四張表在另一個段（`cs:word_10D52`，實測 0x2754
# ⇒ linear 0x27540），段內偏移 ＋ 0x80 才是區塊偏移（docs/spec/138 §2）。
# 把整塊當連續搬會做出一份壞存檔——軍團的勢力欄位讀成 83，remake 直接 panic。
GLOBAL_BASE, GLOBAL_LEN = 0x10CF0, 0x80      # 區塊 +0x00 起的全域欄位（ipeek，IDA linear）
# ⚠ 四張表在**動態配置的段**（`cs:word_10D52`，實測 0x2754），
# `ipeek` 的 IDA linear 讀不到它——那裡是執行檔映像，讀出來是 TALK 文字。
# 要用 `peek:段:偏移:長度` 的真實定址。段內偏移 ＋ 0x80 ＝ 區塊偏移。
TABLE_SEG, TABLE_OFF = 0x2754, 0x80
TABLE_LEN = 0x5220                           # 勢力 ＋ 據點 ＋ 軍團 ＋ 武將


def chunks():
    """回傳 (IDA linear, 區塊偏移, 長度)，兩段來源各自切片。"""
    out = []
    out.append((f"{GLOBAL_BASE:X}", 0, GLOBAL_LEN))
    off = 0
    while off < TABLE_LEN:
        n = min(CHUNK, TABLE_LEN - off)
        out.append((f"{TABLE_SEG:04X}:{off:04X}", TABLE_OFF + off, n))
        off += n
    return out


def parse(text: str) -> dict:
    """把 `ipeek` 的每一行 `ADDR = XX XX …` 依位址排好、接成連續 bytes。"""
    seen = {}
    for line in text.splitlines():
        m = re.match(r"\s*([0-9A-Fa-f:]{4,11})\s*=\s*((?:[0-9A-Fa-f]{2}\s*)+)$", line.strip())
        if not m:
            continue
        seen[m.group(1).upper()] = bytes(int(x, 16) for x in m.group(2).split())
    out = {}
    for key, dst, want in chunks():
        got = seen.get(key)
        if got is None:
            raise SystemExit(f"缺少 {key}:{want}——腳本裡的每一段都跑了嗎？")
        if len(got) < want:
            raise SystemExit(f"{key} 只有 {len(got)} byte，需要 {want}")
        out[dst] = got[:want]
    return out


def selftest() -> int:
    ok = True

    def check(label, cond):
        nonlocal ok
        print(f"  {'✓' if cond else '✗'} {label}")
        ok = ok and cond

    cs = chunks()
    check("全域那一段從 0x10CF0 起、128 B", cs[0] == ("10CF0", 0, GLOBAL_LEN))
    check("四張表走 peek:2754 且落在區塊 +0x80",
          cs[1][0].startswith("2754:") and cs[1][1] == TABLE_OFF)
    check("四張表總長 0x5220",
          sum(n for k, _, n in cs if k.startswith("2754:")) == TABLE_LEN)
    check("⚠ 不覆蓋整個區塊（未解區域保留來源）",
          sum(n for _, _, n in cs) < BLOCK)

    text = "\n".join(f"   {k} = " + " ".join(f"{(i % 256):02X}" for i in range(n))
                     for k, _, n in cs)
    got = parse(text)
    check("兩段都解出來", 0 in got and TABLE_OFF in got)
    check("全域段長度正確", len(got[0]) == GLOBAL_LEN)
    check("內容照位址順序", got[0][:4] == bytes([0, 1, 2, 3]))
    # 負對照：少一段就要報錯，不能默默補零。
    short = "\n".join(f"   {k} = " + " ".join("00" for _ in range(n))
                      for k, _, n in cs[:-1])
    try:
        parse(short)
        check("缺一段會報錯（負對照）", False)
    except SystemExit:
        check("缺一段會報錯（負對照）", True)
    return 0 if ok else 1


def main() -> int:
    args = [a for a in sys.argv[1:] if not a.startswith("--")]
    if "--selftest" in sys.argv:
        return selftest()
    if len(args) < 3:
        print(__doc__)
        return 2
    peek, src, dst = args[0], args[1], args[2]
    slot = 0
    if "--slot" in sys.argv:
        slot = int(sys.argv[sys.argv.index("--slot") + 1])
    if not 0 <= slot < BLOCKS:
        raise SystemExit("--slot 要在 0–3")

    parts = parse(open(peek, encoding="utf-8", errors="replace").read())
    data = bytearray(open(src, "rb").read())
    if len(data) != FILE_SIZE:
        raise SystemExit(f"{src} 是 {len(data)} B，預期 {FILE_SIZE}")
    base = slot * BLOCK
    n = 0
    for off, chunk in sorted(parts.items()):
        data[base + off:base + off + len(chunk)] = chunk
        n += len(chunk)
    with open(dst, "wb") as f:
        f.write(data)
    g = parts[0]
    print(f"{dst}（第 {slot + 1} 槽）：換上原版當下的 {n} byte"
          f"（其餘 {BLOCK - n} byte 保留來源，未解區域不動）")
    print(f"  時鐘 {g[6] | g[7] << 8}年{g[4]}月{g[0]}日 "
          f"{g[3]}時 子刻 {g[2]}；據點游標 {(g[0x2E] | g[0x2F] << 8) // 32}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
