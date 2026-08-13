#!/usr/bin/env python3
"""TALK.DAT 解析、抽取與寫回。

格式見 `docs/formats/01-talk-dat.md`。摘要：

    0x000–0x7FF   1024 筆 uint16 LE 偏移表
                  [0..1021] = 每則訊息的起點
                  [1022]    = 結束哨兵（＝檔案長度 − 1）
                  [1023]    = 0，未使用
    0x800–        1022 則訊息，每則是數個 NUL 結尾的行，再以一個空行收尾

用法：

    talkdat.py dump   <TALK.DAT> <big5|shift_jis>          印出全部訊息
    talkdat.py export <TALK.DAT> <編碼> <out.json>          抽成 JSON
    talkdat.py build  <in.json>  <編碼> <out.DAT>           由 JSON 組回
    talkdat.py verify <TALK.DAT> <編碼>                     round-trip 驗證
    talkdat.py diff   <A.DAT> <encA> <B.DAT> <encB>         兩版對照
    talkdat.py correct <in.json> <corrections.json> <out.json> 套用已定案校訂

只用標準函式庫，不裝任何套件。
"""
import json
import os
import struct
import sys

TABLE_BYTES = 0x800
TABLE_ENTRIES = TABLE_BYTES // 2
MSG_COUNT = 1022                       # [1022] 是哨兵，[1023] 未使用


BIG5_FAMILY = ('big5', 'cp950', 'big5hkscs')
SJIS_FAMILY = ('shift_jis', 'cp932', 'shift_jisx0213')


def _is_lead(byte, enc):
    """判斷這個 byte 是不是雙位元組字的第一個 byte。

    這一步不能省。Big5 與 Shift-JIS 的**第二個 byte 都可能是 0x5C**（反斜線），
    直接拿 regex 找 `\\N` 會把字切一半。
    """
    if enc in BIG5_FAMILY:
        return 0x81 <= byte <= 0xfe
    if enc in SJIS_FAMILY:
        return (0x81 <= byte <= 0x9f) or (0xe0 <= byte <= 0xfc)
    raise ValueError(f'不認識的編碼 {enc}；請加進 BIG5_FAMILY 或 SJIS_FAMILY')


def parse(blob):
    """回傳 (偏移表, 1022 則訊息的原始 bytes)。

    `table[1022]` 是**最後一個 byte 的位移**，不是 one-past-end——
    所以最後一則要讀到 EOF，照 `table[1022]` 切會少一個 byte。
    """
    table = [struct.unpack('<H', blob[i:i + 2])[0]
             for i in range(0, TABLE_BYTES, 2)]
    ends = [table[i + 1] for i in range(MSG_COUNT - 1)] + [len(blob)]
    msgs = [blob[table[i]:ends[i]] for i in range(MSG_COUNT)]
    return table, msgs


def _safe_pair(pair, enc):
    """能安全來回轉換就回傳字元，否則回傳 None。

    只信「解得開**而且**編回去一模一樣」的字。Python 的 big5／shift_jis
    轉換表與 1994 年的原版並不完全一致，光靠 decode 成功不足以保證寫得回去
    —— 而寫不回去的中文化工具是不能用的。
    """
    try:
        ch = pair.decode(enc)
    except UnicodeDecodeError:
        return None
    try:
        return ch if ch.encode(enc) == pair else None
    except UnicodeEncodeError:
        return None


def decode(raw, enc):
    """把一則訊息拆成行；變數標記還原成 `{N}`，還原不了的 byte 寫成 `<xx>`。"""
    lines, cur, i = [], '', 0
    while i < len(raw):
        c = raw[i]
        if c == 0x5c and i + 1 < len(raw):            # 變數標記
            cur += '{' + chr(raw[i + 1]) + '}'
            i += 2
            continue
        if c == 0x00:                                  # 行結束（空行也是資料）
            lines.append(cur)
            cur = ''
            i += 1
            continue

        if _is_lead(c, enc) and i + 1 < len(raw):
            ch = _safe_pair(raw[i:i + 2], enc)
            if ch is not None:
                cur += ch
                i += 2
                continue
        if 0x20 <= c < 0x7f and chr(c) not in '{<':
            cur += chr(c)
        else:
            cur += f'<{c:02x}>'
        i += 1
    if cur:                                            # 沒有 NUL 收尾的殘行
        lines.append(cur)
    return lines


def encode(lines, enc):
    """decode 的反向：每一行後面一個 NUL。

    **不要自作主張補或刪結尾的空行。** 訊息結尾有幾個空行是原版的資料
    （對話框裡的留白），不是格式規定——有的訊息一個都沒有，
    有的有兩個。刪掉就寫不回去了。
    """
    out = bytearray()
    for line in lines:
        i = 0
        while i < len(line):
            if line[i] == '{' and i + 2 < len(line) and line[i + 2] == '}':
                out += b'\x5c' + line[i + 1].encode('ascii')
                i += 3
            elif line[i] == '<' and i + 3 < len(line) and line[i + 3] == '>':
                out.append(int(line[i + 1:i + 3], 16))
                i += 4
            else:
                out += line[i].encode(enc)
                i += 1
        out.append(0x00)
    return bytes(out)


def build(messages, enc):
    """由訊息串列組出完整的 TALK.DAT。"""
    body = bytearray()
    table = []
    for lines in messages:
        table.append(TABLE_BYTES + len(body))
        body += encode(lines, enc)
    table.append(TABLE_BYTES + len(body) - 1)        # [1022] 哨兵
    table.append(0)                                  # [1023] 未使用
    head = b''.join(struct.pack('<H', v) for v in table)
    assert len(head) == TABLE_BYTES, len(head)
    return head + bytes(body)


def markers(lines):
    """一則訊息用到的變數標記，排序後回傳。"""
    found = []
    for line in lines:
        i = 0
        while i < len(line):
            if line[i] == '{' and i + 2 < len(line) and line[i + 2] == '}':
                found.append(line[i + 1])
                i += 3
            else:
                i += 1
    return sorted(found)


def _load(path, enc):
    blob = open(path, 'rb').read()
    _, msgs = parse(blob)
    return blob, [decode(m, enc) for m in msgs]


def _correction_tokens(text):
    """把文字拆成可換行的 token；`{N}` 視為一個三格全形插入值。

    這不是原版渲染器的完整字寬模型。校訂工具只需要一個保守、可重現的
    JSON 版面保護：中文／標點各一格，武將／據點等 `{N}` 先按三格計。
    產出的行仍須經實機畫面抽樣，不能把這個估算當成 parity 證據。
    """
    out = []
    i = 0
    while i < len(text):
        if (text[i] == '{' and i + 2 < len(text)
                and text[i + 2] == '}' and text[i + 1].isdigit()):
            start = i
            width = 0
            # 連續的數值插入（例如 `{6}{7}`）不可被拆到兩行中間。
            while (i + 2 < len(text) and text[i] == '{'
                   and text[i + 2] == '}' and text[i + 1].isdigit()):
                width += 3
                i += 3
            out.append((text[start:i], width))
        else:
            out.append((text[i], 1))
            i += 1
    return out


def _display_width(text):
    return sum(width for _, width in _correction_tokens(text))


def _wrap_correction(text, width):
    """以原訊息的最大非空行寬度包裝校訂文字。"""
    if not text:
        return ['']
    width = max(1, width)
    lines, cur, used = [], [], 0
    for token, token_width in _correction_tokens(text):
        # 不讓中文標點孤零零落在新行開頭；這個排版保護可能讓該行多一格，
        # 呼叫端會收到畫面抽樣警告。
        if token in '，。！？；：、）】」』':
            if not cur and lines:
                lines[-1] += token
                continue
            if cur and used + token_width > width:
                cur.append(token)
                used += token_width
                continue
        if cur and used + token_width > width:
            lines.append(''.join(cur))
            cur, used = [], 0
        cur.append(token)
        used += token_width
    if cur:
        lines.append(''.join(cur))
    return lines


def _trailing_empty(lines):
    n = 0
    for line in reversed(lines):
        if line != '':
            break
        n += 1
    return n


def _correction_expected(item, by_id):
    expected = item.get('cht', '')
    if isinstance(expected, str) and expected.startswith('同 '):
        ref = int(expected.split()[1])
        if ref not in by_id or by_id[ref].get('cht', '').startswith('同 '):
            raise ValueError(f"#{item['id']} 的同文參照無效：{expected}")
        return by_id[ref]['cht']
    return expected


def cmd_correct(src, corrections_path, out):
    """只套用有 `fix` 的項目，保留 null 的人工裁定項。"""
    if os.path.abspath(src) == os.path.abspath(out):
        raise ValueError('correct 的輸出不可覆寫輸入 JSON')
    data = json.load(open(src, encoding='utf-8'))
    correction_data = json.load(open(corrections_path, encoding='utf-8'))
    messages = data.get('messages')
    items = correction_data.get('corrections')
    enc = data.get('encoding', 'cp950')
    if not isinstance(messages, list) or len(messages) != MSG_COUNT:
        raise ValueError(f'{src} 不是 {MSG_COUNT} 則訊息的 extract JSON')
    if not isinstance(items, list):
        raise ValueError(f'{corrections_path} 沒有 corrections 陣列')

    by_id = {int(item['id']): item for item in items}
    result = [list(lines) for lines in messages]
    applied, skipped, layout_changes, layout_warnings = [], [], [], []
    for item in items:
        ident = int(item['id'])
        fix = item.get('fix')
        if fix is None:
            skipped.append(ident)
            continue
        if ident < 0 or ident >= len(result):
            raise ValueError(f'校訂編號超出 TALK.DAT：#{ident}')
        current = ''.join(result[ident])
        expected = _correction_expected(item, by_id)
        if current != expected:
            raise ValueError(
                f'#{ident} 的現況與 corrections.json 不符：\n'
                f'  現況：{current}\n  預期：{expected}')
        jp = item.get('jp', '')
        if jp and not (isinstance(jp, str) and jp.startswith('同 ')):
            if markers([fix]) != markers([jp]):
                raise ValueError(f'#{ident} fix 與 jp 的變數集合不同')

        trailing_count = _trailing_empty(result[ident])
        body = result[ident][:-trailing_count] if trailing_count else result[ident]
        max_width = max((_display_width(line) for line in body if line), default=10)
        wrapped = _wrap_correction(fix, max_width)
        trailing = [''] * trailing_count
        try:
            encode(wrapped + trailing, enc)
        except (UnicodeEncodeError, ValueError) as exc:
            raise ValueError(f'#{ident} fix 無法以 {enc} 編碼：{exc}') from exc
        result[ident] = wrapped + trailing
        applied.append(ident)
        if len(wrapped) != len(body):
            layout_changes.append((ident, len(body), len(wrapped), max_width))
        over = [(_display_width(line), line) for line in wrapped
                if _display_width(line) > max_width]
        if over:
            layout_warnings.append((ident, max_width, over))

    output = dict(data)
    output['messages'] = result
    with open(out, 'w', encoding='utf-8') as fp:
        json.dump(output, fp, ensure_ascii=False, indent=1)
        fp.write('\n')
    print(f'✅ 套用 {len(applied)} 筆校訂 → {out}')
    print('   已套用：' + ', '.join(f'#{i}' for i in applied))
    print('   保留人工裁定：' + ', '.join(f'#{i}' for i in skipped))
    if layout_changes:
        print('   行數變更（需畫面抽樣）：' + ', '.join(
            f'#{i} {before}→{after} 行／寬 {width}'
            for i, before, after, width in layout_changes))
    if layout_warnings:
        print('   行寬警告（需畫面抽樣）：' + ', '.join(
            f'#{i} 上限 {width}、實際 ' + '/'.join(str(actual) for actual, _ in over)
            for i, width, over in layout_warnings))
    return 0


def cmd_dump(path, enc):
    _, messages = _load(path, enc)
    for i, lines in enumerate(messages):
        print(f'#{i:4d}  ' + ' ／ '.join(lines) if lines else f'#{i:4d}  （空）')


def cmd_export(path, enc, out):
    _, messages = _load(path, enc)
    with open(out, 'w', encoding='utf-8') as fp:
        json.dump({'encoding': enc, 'messages': messages},
                  fp, ensure_ascii=False, indent=1)
    print(f'{out}: {len(messages)} 則')


def cmd_build(src, enc, out):
    data = json.load(open(src, encoding='utf-8'))
    blob = build(data['messages'], enc)
    open(out, 'wb').write(blob)
    print(f'{out}: {len(blob)} bytes')


def cmd_verify(path, enc):
    """驗收標準：解出來再組回去，必須與原檔 byte-for-byte 相同。"""
    blob, messages = _load(path, enc)
    rebuilt = build(messages, enc)
    if rebuilt == blob:
        print(f'✅ round-trip 相同：{len(blob)} bytes / {len(messages)} 則')
        return 0
    print(f'❌ round-trip 不同：原 {len(blob)} B，重建 {len(rebuilt)} B')
    for i, (a, b) in enumerate(zip(blob, rebuilt)):
        if a != b:
            print(f'   第一個差異 @ 0x{i:04x}: {a:02x} != {b:02x}')
            break
    return 1


def cmd_diff(pa, ea, pb, eb):
    _, ma = _load(pa, ea)
    _, mb = _load(pb, eb)
    print(f'則數 {len(ma)} / {len(mb)}')
    bad = [i for i in range(min(len(ma), len(mb)))
           if markers(ma[i]) != markers(mb[i])]
    print(f'變數標記不一致：{len(bad)} 則 → {bad}')
    for i in bad:
        print(f'#{i}')
        print('   A ' + ' ／ '.join(ma[i]))
        print('   B ' + ' ／ '.join(mb[i]))
    return 0


if __name__ == '__main__':
    if len(sys.argv) < 2:
        sys.exit(__doc__)
    cmd, args = sys.argv[1], sys.argv[2:]
    fn = {'dump': cmd_dump, 'export': cmd_export, 'build': cmd_build,
          'verify': cmd_verify, 'diff': cmd_diff, 'correct': cmd_correct}.get(cmd)
    if not fn:
        sys.exit(__doc__)
    sys.exit(fn(*args) or 0)
