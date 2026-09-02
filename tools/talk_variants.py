#!/usr/bin/env python3
"""把 `#406` 之後的八格變體組攤開成日中對照，並標出組內的可疑處。

    tools/py.sh tools/talk_variants.py            # 只印摘要與可疑組
    tools/py.sh tools/talk_variants.py --all      # 連同全部 77 組逐則印出
    tools/py.sh tools/talk_variants.py --group 0x1BD

## 為什麼要以組為單位

`sub_1075B` 對索引 ≥ `0x196`（406）做 ×8 展開：真正的索引是
`0x196 + (cx − 0x196) × 8 + ah`，八則是**同一情境的八種說法**，
由武將記錄 `+0x1E`（說話型）選一個（`docs/re/25` §1）。

所以「這一則譯得對不對」不是單則能回答的問題：同組譯成同一句就等於
失去變化，而 `{N}` 變數的插入位置，同組其他七則是最強的旁證
（`docs/reference/02` §14）。

## 標出來的四種可疑

| 記號 | 意思 |
|---|---|
| `重複` | **中文**組內有兩則以上相同，而日文那幾則**在中文分得出來的層次上**不同 ＝ 變化被抹平 |
| `變數` | 同一則的中日 `{N}` 標記集合不同 |
| `空` | 日文有內容而中文是空的 |
| `長度` | 中文字數不到日文的四成或超過兩倍（只是提示，不是錯誤）|

⭐ **比日文之前要先正規化。** 八格變體的差異有一大部分落在中文分不出來的
層次——第一人稱（俺／私／わし／儂）、句尾語氣（〜ぞ／〜だ／〜のだ）與
排版空白。那幾種差異譯成同一句中文是正確的，不是抹平；不先正規化，
同一批誤報每次都會再出現一次。

**這支工具不下結論**：它只把要人看的組挑出來，判讀仍然是人工的。
"""
import argparse
import json
import os
import re
import sys

VARIANT_BASE = 0x196          # 406：×8 展開的門檻
VARIANT_SPAN = 8
MSG_COUNT = 1022

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOSV = os.path.join(ROOT, 'translations', 'talk-dosv-corrected.json')
# 空格分析看的是**原文**：校訂只改字句，不會讓一則從有變空，
# 但拿校訂後的檔去數空格，等於讓工具的答案跟著校訂進度漂。
RAW_DOSV = os.path.join(ROOT, 'translations', 'extract', 'talk-dosv.json')
PC98 = os.path.join(ROOT, 'translations', 'extract', 'talk-pc98.json')
# 已經人工審過而決定不改的組（連理由一起記）。**看過的結論要留下來**，
# 否則同一批誤報每一輪都要重讀一次。
AUDIT = os.path.join(ROOT, 'translations', 'variant-audit.json')

MARKER = re.compile(r'\{(\d+)\}')

# 日文正規化：中文分不出來的三個層次。
JA_PRONOUN = re.compile(r'(俺|私|わし|儂|わたし|拙者|僕)')
JA_TAIL = re.compile(r'(のだ|んぞ|ぞ|だ|わ|ぬ|な|ね|よ|ぜ)+(?=[！？。」]|$)')
JA_SPACE = re.compile(r'[\s\u3000]+')


def normalise_ja(msg):
    """把日文正規化到「中文能分得出來」的層次。

    去掉全形空白（那是排版）、把第一人稱統一、削掉句尾語氣助詞——
    這三種差異在中文裡都會落在同一句話上。

    ⚠ **語氣助詞要逐行削**：訊息是分行存的，連成一串之後行尾的「だ」
    後面接的是下一行的第一個字，`$` 就對不上了。
    """
    out = []
    for line in msg:
        if not line:
            continue
        line = JA_SPACE.sub('', line)
        line = JA_PRONOUN.sub('私', line)
        out.append(JA_TAIL.sub('', line))
    return ''.join(out)


def load(path):
    with open(path, encoding='utf-8') as fh:
        return json.load(fh)['messages']


def reviewed():
    try:
        with open(AUDIT, encoding='utf-8') as fh:
            raw = json.load(fh).get('reviewed', {})
    except OSError:
        return {}
    return {int(k, 16): v for k, v in raw.items()}


def joined(msg):
    return ''.join(line for line in msg if line)


def markers(msg):
    return sorted(MARKER.findall(joined(msg)))


def groups():
    """回傳 [(cx, [索引…])]，cx 是呼叫端傳的值。"""
    out = []
    cx = VARIANT_BASE
    while True:
        first = VARIANT_BASE + (cx - VARIANT_BASE) * VARIANT_SPAN
        if first >= MSG_COUNT:
            break
        idx = [i for i in range(first, min(first + VARIANT_SPAN, MSG_COUNT))]
        out.append((cx, idx))
        cx += 1
    return out


def audit(cx, idx, cht, jpn):
    """回傳這一組的可疑清單。"""
    flags = []
    zh = [joined(cht[i]) for i in idx]
    ja = [joined(jpn[i]) for i in idx]
    norm = [normalise_ja(jpn[i]) for i in idx]

    # 重複：中文抹平了日文的變化
    seen = {}
    for k, text in enumerate(zh):
        if not text:
            continue
        seen.setdefault(text, []).append(k)
    for text, ks in seen.items():
        if len(ks) < 2:
            continue
        if len({norm[k] for k in ks}) > 1:
            flags.append(('重複', '#%d 與 #%d 中文相同、日文不同' %
                          (idx[ks[0]], idx[ks[1]])))

    for k, i in enumerate(idx):
        if markers(cht[i]) != markers(jpn[i]):
            flags.append(('變數', '#%d 中 %s ／ 日 %s' %
                          (i, markers(cht[i]) or '—', markers(jpn[i]) or '—')))
        if ja[k] and not zh[k]:
            flags.append(('空', '#%d 日文有內容、中文是空的' % i))
        elif ja[k] and zh[k]:
            ratio = len(zh[k]) / len(ja[k])
            if ratio < 0.4 or ratio > 2.0:
                flags.append(('長度', '#%d 中 %d 字／日 %d 字' %
                              (i, len(zh[k]), len(ja[k]))))
    return flags


def show(cx, idx, cht, jpn, flags):
    print('==== cx=0x%X  #%d–#%d ====' % (cx, idx[0], idx[-1]))
    for i in idx:
        print('  #%-4d 日 %s' % (i, ' ／ '.join(x for x in jpn[i] if x) or '（空）'))
        print('        中 %s' % (' ／ '.join(x for x in cht[i] if x) or '（空）'))
    for kind, note in flags:
        print('  ⚠ %s：%s' % (kind, note))
    print()


def empty_report():
    """把「哪幾格是空的」按樣式分組印出來。

    ⭐ **空格不是資料缺損，是說話型的補集。** 八格一組由武將記錄 `+0x1E`
    選一格（0–2 主公型／3–7 臣下型，`docs/mechanics/60` §3），所以只對
    臣下說得通的情境（打了敗仗回來道歉）就把 0–2 留白，只有君主會說的
    （罵軍師財政赤字）就把 3–7 留白。

    ⚠ 兩版要一起看。**只看一版分不出「這一版漏譯」與「這一格本來就沒有」**——
    松崗版是譯本，漏譯的症狀同樣是空字串。
    """
    cht, jpn = load(RAW_DOSV), load(PC98)
    pats = {}
    for cx, idx in groups():
        empty_c = tuple(k for k, i in enumerate(idx) if not joined(cht[i]))
        empty_j = tuple(k for k, i in enumerate(idx) if not joined(jpn[i]))
        if empty_c != empty_j:
            print('⚠ 兩版空格不同 組 %s：中 %s／日 %s' % (hex(cx), list(empty_c), list(empty_j)))
        if empty_c:
            pats.setdefault(empty_c, []).append((cx, idx))
    total = sum(len(g) * len(k) for k, g in pats.items())
    print('空訊息 %d 則，落在 %d 組，%d 種樣式\n'
          % (total, sum(len(g) for g in pats.values()), len(pats)))
    for empty, gs in sorted(pats.items(), key=lambda kv: -len(kv[1])):
        filled = [k for k in range(VARIANT_SPAN) if k not in empty]
        print('── 空格 %s（有字的是 %s）：%d 組 %s'
              % (list(empty), filled, len(gs), [hex(cx) for cx, _ in gs]))
        cx, idx = gs[0]
        print('   例 %s 格%d：%s' % (hex(cx), filled[0], joined(cht[idx[filled[0]]])[:40]))
    return 0


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--all', action='store_true', help='連沒有可疑處的組也印')
    ap.add_argument('--group', help='只看某一組（十六進位的 cx）')
    ap.add_argument('--empty', action='store_true',
                    help='只報空格分佈（哪幾組的哪幾格沒有字，兩版一起比）')
    args = ap.parse_args()

    if args.empty:
        return empty_report()

    cht, jpn = load(DOSV), load(PC98)
    if len(cht) != len(jpn):
        print('兩版則數不同：%d / %d' % (len(cht), len(jpn)))
        return 1

    only = int(args.group, 16) if args.group else None
    seen_before = reviewed()
    total, flagged, closed = 0, 0, 0
    counts = {}
    for cx, idx in groups():
        if only is not None and cx != only:
            continue
        total += 1
        flags = audit(cx, idx, cht, jpn)
        note = seen_before.get(cx)
        if flags and note and not args.all and only is None:
            closed += 1
            continue
        for kind, _ in flags:
            counts[kind] = counts.get(kind, 0) + 1
        if flags:
            flagged += 1
        if flags or args.all or only is not None:
            show(cx, idx, cht, jpn, flags)
            if note:
                print('  ✓ 已審（%s）：%s' % (note.get('日期', ''), note.get('理由', '')))
                print()
    print('組數 %d，待審 %d 組，已審過而不改 %d 組' % (total, flagged, closed))
    for kind in sorted(counts):
        print('  %s：%d 筆' % (kind, counts[kind]))
    return 0


if __name__ == '__main__':
    sys.exit(main())
