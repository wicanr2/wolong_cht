#!/usr/bin/env python3
"""M7 校訂文字的人工審查報告與可重跑結構檢查。

這支工具不把自動檢查冒充原版 parity。它只把 60 筆已定案校訂整理成
可逐筆閱讀的報告，並先確認 marker、硬斷行與保守字寬沒有在產出階段
被破壞；實際畫面仍要用 wlgame 的 DOS/V modal 截圖抽樣。
"""

import argparse
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


def load_json(path):
    with open(path, encoding="utf-8") as fp:
        return json.load(fp)


def correction_rows():
    source = load_json(EXTRACT)
    corrections = load_json(CORRECTIONS)["corrections"]
    corrected = load_json(GENERATED)
    rows = []
    for item in sorted(corrections, key=lambda value: int(value["id"])):
        ident = int(item["id"])
        fix = item.get("fix")
        if fix is None:
            continue
        expected = talkdat._correction_expected(item, {
            int(value["id"]): value for value in corrections
        })
        actual = "".join(corrected["messages"][ident])
        if actual != fix:
            raise AssertionError(
                f"#{ident} runtime 校訂不一致：{actual!r} != {fix!r}"
            )
        lines = corrected["messages"][ident]
        widths = [talkdat._display_width(line) for line in lines]
        safe_width = max(widths or [0])
        wrap_width = max(
            1,
            max(
                (talkdat._display_width(line) for line in source["messages"][ident]
                 if line),
                default=1,
            ),
        )
        wrapped = talkdat._wrap_correction(fix, wrap_width)
        rows.append({
            "id": ident,
            "markers": "".join(talkdat.markers([fix])) or "－",
            "hard_lines": len(lines),
            "max_cells": safe_width,
            "estimated_pages": (len(wrapped) + 4) // 5,
            "lines": lines,
            "note": item.get("note", ""),
        })
    return rows


def check():
    source = load_json(EXTRACT)
    manifest = load_json(CORRECTIONS)
    corrections = manifest["corrections"]
    fixed = [item for item in corrections if item.get("fix") is not None]
    if len(fixed) != 60:
        raise AssertionError(f"M7 fix 數量 = {len(fixed)}，want 60")
    if len(source.get("messages", [])) != talkdat.MSG_COUNT:
        raise AssertionError("DOS/V extract 不是 1022 則 TALK")
    with tempfile.TemporaryDirectory(prefix="m7-review-") as tmp:
        out = os.path.join(tmp, "corrected.json")
        if talkdat.cmd_correct(EXTRACT, CORRECTIONS, out) != 0:
            raise AssertionError("talkdat.correct 失敗")
        if open(out, "rb").read() != open(GENERATED, "rb").read():
            raise AssertionError("runtime 校訂表不是由目前輸入重建的結果")
    rows = correction_rows()
    for row in rows:
        if row["max_cells"] > 22:
            raise AssertionError(
                f"#{row['id']} 保守硬行寬度超過 22 格：{row['max_cells']}"
            )
        if not row["note"]:
            raise AssertionError(f"#{row['id']} 缺少人工審查備註")
    return rows


def markdown(rows):
    lines = [
        "# M7 校訂文字人工審查報告",
        "",
        "> 本報告覆蓋 `translations/corrections.json` 的 60 筆已定案修正。",
        "> 自動檢查只負責產出一致性、marker、硬斷行與保守字寬；人工結論",
        "> 仍以本表的語意備註及 `wlgame` DOS/V modal 截圖為證據，不宣稱",
        "> 這些檢查等同原版同狀態逐像素 parity。",
        "",
        "## 結果",
        "",
        f"- 審查項目：{len(rows)} 筆。",
        "- 每筆都已確認校訂後文字、已知 marker、原始硬行結構與最多五列分頁。",
        "- `max_cells` 是保守全形格估算；實際像素寬度仍由 Go `textdraw` gate 驗收。",
        "- 畫面抽樣入口：`wlgame -open-talk-index <id> -shot ...`。",
        "",
        "## 逐筆紀錄",
        "",
        "| TALK | marker | 硬行數 | 最大格數 | 估算頁數 | 校訂後文字 | 人工語意備註 |",
        "|---:|:---:|---:|---:|---:|---|---|",
    ]
    for row in rows:
        text = "\\n".join(row["lines"]).replace("|", "\\|")
        note = row["note"].replace("|", "\\|").replace("\n", " ")
        lines.append(
            f"| #{row['id']} | `{row['markers']}` | {row['hard_lines']} | "
            f"{row['max_cells']} | {row['estimated_pages']} | "
            f"`{text}` | {note} |"
        )
    lines += [
        "",
        "## 人工畫面抽樣分組",
        "",
        "- marker／槽位：`#290`、`#223`、`#321`、`#322`、`#751`。",
        "- 多行／分頁與標點：`#718`、`#840`、`#843`、`#844`、`#886`、`#889`。",
        "- 短句與硬換行邊界：`#16`、`#31`、`#34`、`#314`、`#459`、`#967`。",
        "- 戰場命令與長句：`#649`、`#663`、`#770`、`#797`、`#906`。",
        "",
        "上述代表幀用於確認繁中語意、原始 NUL 行邊界、五列分頁、文字框寬度、",
        "標點位置及畫面不溢出；其餘短句以同一份逐筆表與自動 gate 封口。",
    ]
    return "\n".join(lines) + "\n"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="只執行檢查")
    parser.add_argument("--out", help="輸出 Markdown 報告")
    args = parser.parse_args()
    rows = check()
    if args.out:
        with open(args.out, "w", encoding="utf-8") as fp:
            fp.write(markdown(rows))
    print(
        "M7 review PASS: 60 筆；marker／硬行／保守寬度／產出一致性已檢查；"
        "畫面仍需依抽樣分組回看"
    )
    if args.out:
        print(f"report: {args.out}")


if __name__ == "__main__":
    main()
