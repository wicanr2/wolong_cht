#!/usr/bin/env python3
"""產生 remake 的語系包（docs/spec/84）。

    tools/langpack.sh                 # zh-hans（需要 OpenCC，走帶網路的容器）
    tools/py.sh tools/langpack.py ja  # 日文：直接取 PC-98 原版，不是翻譯
    tools/py.sh tools/langpack.py en  # 合併 workplace/lang/en/out-*.json
    tools/py.sh tools/langpack.py selftest

產出都進版控（是資料不是快取），重跑可重生：

  translations/talk-zh-hans.json  OpenCC t2s 詞級機轉初稿
  translations/t2s-chars.json     字級繁→簡表（UI 詞與人名用）
  translations/talk-ja.json       **PC-98 日文原版的 TALK.DAT**
  translations/talk-en.json       英譯（來源見 workplace/lang/en/BRIEF.md）

母本是 `talk-dosv-corrected.json`（60 筆校訂後的繁中）。
**日文不經翻譯**：PC-98 原版與松崗版是同一個檔的兩種語言，逐則對應，
所以日文語系包是「把原版讀出來」而不是「把中文翻回去」。
"""

import argparse
import glob
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(ROOT / "tools"))

SRC = ROOT / "translations" / "talk-dosv-corrected.json"
OUT_TALK = ROOT / "translations" / "talk-zh-hans.json"
OUT_CHARS = ROOT / "translations" / "t2s-chars.json"
OUT_JA = ROOT / "translations" / "talk-ja.json"
OUT_EN = ROOT / "translations" / "talk-en.json"
PC98_TALK = ROOT / "workplace" / "orig" / "pc98" / "TALK.DAT"
EN_PARTS = ROOT / "workplace" / "lang" / "en"

MARKER = re.compile(r"\{[1-7]\}")
COUNT = 1022
# 訊息框一列 10 個全形字 ＝ 20 個半形字元、一頁 4 列（cmd/wlgame/messages.go）。
LINE_CELLS = 20
PAGE_ROWS = 4


def load_master():
    data = json.loads(SRC.read_text(encoding="utf-8"))
    assert len(data["messages"]) == COUNT, f"母本 {len(data['messages'])} 則"
    return data


def write_pack(path, messages, note=None):
    payload = {"encoding": "utf-8", "messages": messages}
    if note:
        payload["note"] = note
    path.write_text(json.dumps(payload, ensure_ascii=False, indent=1) + "\n",
                    encoding="utf-8")


def marker_set(lines):
    """回傳標記的**多重集合**（排序後的 list）。

    ⚠ 比對要用集合不是序列：日文與英文的語序與中文不同，
    `{1}的兵馬向{2}進攻` 在日文是 `{2}に{1}の軍勢が`——**順序變了不是錯**。
    """
    return sorted(MARKER.findall("".join(lines)))


# ── zh-Hans ────────────────────────────────────────────────────────
def build_zh_hans():
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
    write_pack(OUT_TALK, out_msgs)

    # 字級表：涵蓋整個 Big5 常用字區＋次常用區＋母本語料，取「單字轉出
    # 單字且不同」的對。人名與 UI 詞會用到 talk 語料以外的字，不能只掃語料。
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


# ── 日文（原版，不是翻譯）──────────────────────────────────────────
def build_ja():
    import talkdat

    if not PC98_TALK.exists():
        sys.exit(f"找不到 {PC98_TALK}（PC-98 原版素材不進版控，請自備）")
    _, messages = talkdat._load(str(PC98_TALK), "shift_jis")
    if len(messages) != COUNT:
        sys.exit(f"PC-98 TALK.DAT 解出 {len(messages)} 則，預期 {COUNT}")
    master = load_master()
    reordered = sum(
        1 for a, b in zip(master["messages"], messages)
        if marker_set(a) == marker_set(b)
        and MARKER.findall("".join(a)) != MARKER.findall("".join(b)))
    write_pack(OUT_JA, messages,
               note="PC-98 日文原版 TALK.DAT（Shift-JIS → UTF-8）。"
                    "這不是翻譯，是原版文字本身")
    print(f"talk-ja.json：{len(messages)} 則（其中 {reordered} 則的標記順序"
          f"與中文不同——日文語序，不是錯）")


# ── 英文（合併分批譯稿）────────────────────────────────────────────
def build_en():
    master = load_master()["messages"]
    merged, seen = {}, {}
    for path in sorted(glob.glob(str(EN_PARTS / "out-*.json"))):
        part = json.loads(Path(path).read_text(encoding="utf-8"))
        for k, v in part.items():
            if k in merged:
                sys.exit(f"#{k} 在 {seen[k]} 與 {path} 重複")
            merged[k], seen[k] = v, path
    missing = [i for i in range(COUNT) if str(i) not in merged]
    if missing:
        sys.exit(f"缺 {len(missing)} 則：{missing[:10]}…")
    extra = [k for k in merged if not k.isdigit() or int(k) >= COUNT]
    if extra:
        sys.exit(f"多出不存在的編號：{extra[:10]}")

    bad = []
    messages = []
    for i in range(COUNT):
        lines = merged[str(i)]
        if isinstance(lines, str):
            lines = [lines]
        if marker_set(lines) != marker_set(master[i]):
            bad.append((i, marker_set(master[i]), marker_set(lines)))
        messages.append(lines)
    if bad:
        for i, want, got in bad[:10]:
            print(f"⚠ #{i} 標記不符：母本 {want} vs 譯文 {got}", file=sys.stderr)
        sys.exit(f"{len(bad)} 則的變數標記與母本不符")
    write_pack(OUT_EN, messages,
               note="英譯（tools/langpack.py en 合併 workplace/lang/en/out-*.json）")
    print(f"talk-en.json：{COUNT} 則、{pages_report(messages)}")


def wrapped_rows(text):
    """照引擎的折行規則估算列數（半形 1 格、全形 2 格，英文以空白斷詞）。"""
    rows, width = 1, 0
    for word in text.split(" "):
        w = sum(1 if ord(c) < 0x2000 else 2 for c in word)
        if width and width + 1 + w > LINE_CELLS:
            rows += 1
            width = w
        else:
            width += (1 if width else 0) + w
    return rows


def pages_report(messages):
    over = 0
    longest = []
    for i, lines in enumerate(messages):
        rows = sum(wrapped_rows(l) for l in lines if l)
        longest.append((rows, i))
        if rows > PAGE_ROWS:
            over += 1
    longest.sort(reverse=True)
    top = "／".join(f"#{i} {r} 列" for r, i in longest[:3])
    return f"{over} 則會翻頁（最長 {top}）"


# ── selftest ───────────────────────────────────────────────────────
def selftest():
    master = load_master()["messages"]
    bad = 0
    for path, label in ((OUT_TALK, "zh-hans"), (OUT_JA, "ja"), (OUT_EN, "en")):
        if not path.exists():
            print(f"⚠ 缺 {path.name}", file=sys.stderr)
            bad += 1
            continue
        pack = json.loads(path.read_text(encoding="utf-8"))
        if pack.get("encoding") != "utf-8":
            print(f"⚠ {path.name} 的 encoding = {pack.get('encoding')!r}",
                  file=sys.stderr)
            bad += 1
        msgs = pack["messages"]
        if len(msgs) != COUNT:
            print(f"⚠ {path.name} 有 {len(msgs)} 則", file=sys.stderr)
            bad += 1
            continue
        mism = [i for i in range(COUNT)
                if marker_set(msgs[i]) != marker_set(master[i])]
        if mism:
            print(f"⚠ {path.name} 有 {len(mism)} 則標記不符：{mism[:5]}",
                  file=sys.stderr)
            bad += 1
        empty = sum(1 for m in msgs if not any(l.strip() for l in m))
        src_empty = sum(1 for m in master if not any(l.strip() for l in m))
        note = ""
        if label != "ja" and empty != src_empty:
            note = f"（空槽 {empty}，母本 {src_empty}）"
        print(f"{label}: {len(msgs)} 則、標記全對{note}")
    # 簡體：常用繁體字不得殘留（全部是簡體必改字；繁簡同形字不能放進來）
    if OUT_TALK.exists():
        joined = "\n".join(
            l for m in json.loads(OUT_TALK.read_text(encoding="utf-8"))["messages"]
            for l in m)
        for ch in "國軍馬騎無發歲請們":
            if ch in joined:
                print(f"⚠ zh-hans 疑似殘留繁體字：{ch}", file=sys.stderr)
                bad += 1
    if OUT_CHARS.exists():
        chars = json.loads(OUT_CHARS.read_text(encoding="utf-8"))
        for k, v in chars.items():
            assert len(k) == 1 and len(v) == 1, (k, v)
        print(f"t2s-chars: {len(chars)} 字")
    print(f"selftest 問題數：{bad}")
    return 1 if bad else 0


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("cmd", choices=["zh-hans", "ja", "en", "selftest"])
    args = ap.parse_args()
    sys.exit({"zh-hans": build_zh_hans, "ja": build_ja,
              "en": build_en, "selftest": selftest}[args.cmd]())
