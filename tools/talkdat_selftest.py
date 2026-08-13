#!/usr/bin/env python3
"""`talkdat.py correct` 的 Docker-only、無原版輸入自測。"""

import json
import os
import sys
import tempfile

sys.path.insert(0, os.path.dirname(__file__))
import talkdat  # noqa: E402


ROOT = os.path.dirname(os.path.dirname(__file__))
EXTRACT = os.path.join(ROOT, "translations", "extract", "talk-dosv.json")
CORRECTIONS = os.path.join(ROOT, "translations", "corrections.json")
GENERATED = os.path.join(ROOT, "translations", "talk-dosv-corrected.json")
APPLIED = {
    16, 31, 32, 34, 44, 48, 66, 114, 135, 138, 139, 149, 192, 193,
    194, 195, 211, 213, 223, 257, 258, 259, 266, 268, 290, 314, 321,
    322, 347, 351, 352, 359, 405, 426, 428, 440, 459, 470, 518, 524,
    581, 617, 649, 660, 663, 718, 720, 721, 723, 724, 751, 770, 797,
    840, 843, 844, 886, 889, 906, 967,
}
SKIPPED = set()

# `cmd/wlgame.drawMessage` uses a 384 px window, with 12 px of left/right
# content padding.  The correction width model deliberately counts ASCII and
# `{N}` formatter tokens conservatively as one / three full-width cells, so
# 22 cells (352 px) stays inside the 360 px usable content width.  This is a
# remake modal safety check, not evidence of original PC-98 hard wrapping.
MAX_DIALOG_CELLS = 22


def assert_dialog_width(messages):
    overflow = []
    for message_id, lines in enumerate(messages):
        for line_id, line in enumerate(lines):
            width = talkdat._display_width(line)
            if width > MAX_DIALOG_CELLS:
                overflow.append((message_id, line_id, width, line))
    if overflow:
        raise AssertionError(
            f"校訂後訊息超出 remake modal {MAX_DIALOG_CELLS} 格：{overflow[:4]}"
        )


def main():
    with tempfile.TemporaryDirectory(prefix="talkdat-selftest-") as tmp:
        out = os.path.join(tmp, "corrected.json")
        if talkdat.cmd_correct(EXTRACT, CORRECTIONS, out) != 0:
            raise AssertionError("correct 沒有成功回傳")
        if not os.path.isfile(GENERATED):
            raise AssertionError(f"缺少 runtime 校訂表：{GENERATED}")
        expected_bytes = open(out, "rb").read()
        generated_bytes = open(GENERATED, "rb").read()
        if generated_bytes != expected_bytes:
            raise AssertionError(
                "translations/talk-dosv-corrected.json 不是由目前 extract／corrections "
                "重建的結果；請重新執行 talkdat.py correct"
            )
        source = json.load(open(EXTRACT, encoding="utf-8"))
        corrected = json.load(open(out, encoding="utf-8"))
        corrections = json.load(open(CORRECTIONS, encoding="utf-8"))["corrections"]
        by_id = {int(item["id"]): item for item in corrections}

        changed = {
            i for i, (before, after) in enumerate(
                zip(source["messages"], corrected["messages"])
            ) if before != after
        }
        if changed != APPLIED:
            raise AssertionError(f"變更集合 {sorted(changed)} != {sorted(APPLIED)}")
        for ident in SKIPPED:
            if corrected["messages"][ident] != source["messages"][ident]:
                raise AssertionError(f"人工裁定 #{ident} 被意外套用")
        if any(by_id[i].get("fix") is None for i in APPLIED):
            raise AssertionError("APPLIED 集合含有 null fix")

        assert_dialog_width(corrected["messages"])

        raw = talkdat.build(corrected["messages"], corrected["encoding"])
        # 重新走 parse → decode → build，確保校訂 JSON 仍符合 TALK.DAT 格式。
        _, raw_messages = talkdat.parse(raw)
        decoded = [talkdat.decode(message, corrected["encoding"])
                   for message in raw_messages]
        if talkdat.build(decoded, corrected["encoding"]) != raw:
            raise AssertionError("校訂後 TALK.DAT round-trip 失敗")

        bad = os.path.join(tmp, "bad.json")
        bad_data = json.loads(json.dumps(source, ensure_ascii=False))
        bad_data["messages"][718][0] = "刻意改壞"
        with open(bad, "w", encoding="utf-8") as fp:
            json.dump(bad_data, fp, ensure_ascii=False)
        try:
            talkdat.cmd_correct(bad, CORRECTIONS, os.path.join(tmp, "bad-out.json"))
        except ValueError:
            pass
        else:
            raise AssertionError("現況不符時 correct 沒有 fail-closed")
    print("✅ talkdat correct selftest：60 筆套用、runtime 產出一致、22 格行寬 guard、round-trip、mismatch guard")


if __name__ == "__main__":
    main()
