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

只用標準函式庫，不裝任何套件。
"""
import json
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
          'verify': cmd_verify, 'diff': cmd_diff}.get(cmd)
    if not fn:
        sys.exit(__doc__)
    sys.exit(fn(*args) or 0)
