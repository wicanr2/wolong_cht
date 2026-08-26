#!/usr/bin/env python3
"""產生簡體語系包（docs/spec/84）。

    tools/langpack.sh            # 在帶網路的一次性容器裡跑本腳本
    tools/langpack.py --selftest # 驗證既有產出（不需要 opencc）

產出兩個檔，**都進版控**（是資料不是快取）：

  translations/talk-zh-hans.json   1,022 則的 OpenCC t2s 詞級機轉初稿
                                   （utf-8；{N} 變數標記逐則保留）
  translations/t2s-chars.json      字級繁→簡表，涵蓋 talk＋校訂檔出現過的
                                   全部繁體字。runtime 給 UI 詞與人名用；
                                   詞級歧義（乾／干 等）由逐句覆寫詞表處理

母本是 talk-dosv-corrected.json（60 筆校訂後的繁中）。機轉是初稿：
逐則人工校訂屬第二期（spec/84 §3）。
"""

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SRC = ROOT / "translations" / "talk-dosv-corrected.json"
OUT_TALK = ROOT / "translations" / "talk-zh-hans.json"
OUT_CHARS = ROOT / "translations" / "t2s-chars.json"

MARKER = re.compile(r"\{[1-7]\}")


def load_master():
    data = json.loads(SRC.read_text(encoding="utf-8"))
    assert len(data["messages"]) == 1022, f"母本 {len(data['messages'])} 則"
    return data


def generate():
    from opencc import OpenCC

    cc = OpenCC("t2s")
    master = load_master()

    out_msgs = []
    for i, lines in enumerate(master["messages"]):
        conv = []
        for line in lines:
            c = cc.convert(line)
            if MARKER.findall(c) != MARKER.findall(line):
                sys.exit(f"#{i} 轉換後變數標記改變：{line!r} → {c!r}")
            conv.append(c)
        out_msgs.append(conv)
    OUT_TALK.write_text(
        json.dumps({"encoding": "utf-8", "messages": out_msgs},
                   ensure_ascii=False, indent=1) + "\n",
        encoding="utf-8")

    # 字級表：涵蓋整個 Big5 常用字區（A440–C67E）＋次常用區（C940–F9FE）
    # ＋母本語料，取「單字轉出單字且不同」的對。人名與 UI 詞會用到
    # talk 語料以外的字，所以不能只掃語料。
    chars = set()
    for lines in master["messages"]:
        for line in lines:
            chars.update(line)
    for hi in range(0xA4, 0xFA):
        for lo in list(range(0x40, 0x7F)) + list(range(0xA1, 0xFF)):
            try:
                chars.add(bytes([hi, lo]).decode("cp950"))
            except UnicodeDecodeError:
                pass
    table = {}
    for ch in sorted(chars):
        if ord(ch) < 0x3000:
            continue
        s = cc.convert(ch)
        if len(s) == 1 and s != ch:
            table[ch] = s
    OUT_CHARS.write_text(
        json.dumps(table, ensure_ascii=False, indent=0, sort_keys=True) + "\n",
        encoding="utf-8")
    print(f"talk-zh-hans.json：{len(out_msgs)} 則")
    print(f"t2s-chars.json：{len(table)} 字")


def selftest():
    master = load_master()
    out = json.loads(OUT_TALK.read_text(encoding="utf-8"))
    assert out["encoding"] == "utf-8", out["encoding"]
    assert len(out["messages"]) == 1022, len(out["messages"])
    bad = 0
    for i, (a, b) in enumerate(zip(master["messages"], out["messages"])):
        if len(a) != len(b):
            sys.exit(f"#{i} 行數改變 {len(a)}→{len(b)}")
        for la, lb in zip(a, b):
            if MARKER.findall(la) != MARKER.findall(lb):
                sys.exit(f"#{i} 標記不一致：{la!r} vs {lb!r}")
    # 抽樣：常用繁體字不得殘留（機轉至少要把這些換掉）。
    joined = "\n".join(l for m in out["messages"] for l in m)
    for ch in "國軍馬騎無發歲請們":  # 全部是簡體必改字；繁簡同形字（兵 等）不能放進來
        if ch in joined:
            bad += 1
            print(f"⚠ 疑似殘留繁體字：{ch}", file=sys.stderr)
    chars = json.loads(OUT_CHARS.read_text(encoding="utf-8"))
    for k, v in chars.items():
        assert len(k) == 1 and len(v) == 1, (k, v)
    print(f"selftest ok：1022 則、標記逐則一致、字級表 {len(chars)} 字、殘留抽樣 {bad}")
    return 0 if bad == 0 else 1


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--selftest", action="store_true")
    args = ap.parse_args()
    sys.exit(selftest() if args.selftest else generate())
