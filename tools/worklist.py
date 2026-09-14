#!/usr/bin/env python3
"""以 GitHub Issues 為現行工作入口的本地輔助 verify。

    tools/py.sh tools/worklist.py verify
    tools/py.sh tools/worklist.py --selftest

docs/worklist.json 只保存本地檢查訊號與對應的 GitHub Issue，不再生成或驗證
WORKLIST.md。Issue 的狀態、內容與完成條件才是現行工作的權威；Markdown 只保存
證據、規格與歷史，不得用來登記未處理工作。

verify 跑起來為真，表示對應 Issue 仍可能未完成；跑起來為假，表示 Issue
需要回查並更新，而不是讓本地 JSON 自動關閉 Issue。
"""

import json
import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
DATA_PATH = ROOT / "docs/worklist.json"
SCANNED_SUFFIXES = {".go", ".py", ".sh", ".md", ".json", ".idc"}
SELF = "worklist.json"


def is_test_file(path):
    """測試檔不算，因為測試本來就會提到尚未接上的行為。"""
    name = path.name
    return "_test." in name or name.startswith("test_") or name.endswith("_test.py")


def scan_files(paths):
    for target in paths:
        path = ROOT / target
        files = sorted(path.rglob("*")) if path.is_dir() else [path]
        for file_path in files:
            if not file_path.is_file() or file_path.suffix not in SCANNED_SUFFIXES:
                continue
            if is_test_file(file_path) or file_path.name == SELF:
                continue
            yield file_path


def find_hit(verify):
    """回傳第一個命中的相對路徑，沒有命中則回傳 None。"""
    expression = re.compile(verify["pattern"])
    for file_path in scan_files(verify["paths"]):
        try:
            content = file_path.read_text(encoding="utf-8", errors="ignore")
        except OSError:
            continue
        if expression.search(content):
            return str(file_path.relative_to(ROOT))
    return None


def still_open(item):
    """回傳（仍未完成、說明）。"""
    verify = item["verify"]
    kind = verify["kind"]
    if kind == "manual":
        return True, "要人判（沒有機器訊號）"
    if kind == "json_len":
        blob = json.loads((ROOT / verify["path"]).read_text(encoding="utf-8"))
        field = blob[verify["field"]] if verify.get("field") else blob
        count = len(field)
        return count <= verify["max"], (
            f'{verify["path"]} 有 {count} 項（門檻 {verify["max"]}）'
        )
    where = find_hit(verify)
    if kind == "present":
        if where:
            return True, f'自承還在 {where}'
        return False, f'找不到 /{verify["pattern"]}/'
    if kind == "absent":
        if where:
            return False, f'已經出現在 {where}'
        return True, f'還沒出現 /{verify["pattern"]}/'
    raise SystemExit(f'不認得的 verify kind：{kind}（條目 {item["id"]}）')


def load(path=DATA_PATH):
    data = json.loads(Path(path).read_text(encoding="utf-8"))
    if data.get("schema") != "wolong-worklist/2":
        raise SystemExit(
            f'不支援的 worklist schema：{data.get("schema")!r}'
            "（需要 wolong-worklist/2）"
        )

    seen = set()
    for item in data["items"]:
        required = ("id", "layer", "title", "body", "acceptance", "verify", "github_issue")
        for key in required:
            if key not in item:
                raise SystemExit(f'條目 {item.get("id", "?")} 缺 {key}')
        if item["id"] in seen:
            raise SystemExit(f'條目 id 重複：{item["id"]}')
        seen.add(item["id"])
        if item["layer"] not in data["layers"]:
            raise SystemExit(f'條目 {item["id"]} 的 layer 不在 layers 裡')
        issue = item["github_issue"]
        if (
            not isinstance(issue, dict)
            or not isinstance(issue.get("number"), int)
            or issue["number"] <= 0
            or not issue.get("url")
        ):
            raise SystemExit(f'條目 {item["id"]} 的 github_issue 不完整')
    return data


def cmd_verify(data, quiet=False):
    stale = []
    manual = 0
    lines = []
    for item in data["items"]:
        open_, why = still_open(item)
        if not open_:
            stale.append((item["id"], why))
        if item["verify"]["kind"] == "manual":
            manual += 1
        mark = "仍未完成" if open_ else "⚠ 可能已完成"
        issue = item["github_issue"]
        lines.append(
            f'  #{issue["number"]:<4} {item["id"]:<34} {mark:<14} {why}'
        )
    if not quiet:
        print(
            f'GitHub Issues 對應的本地 verify：{len(data["items"])} 條工作'
            f'（其中 {manual} 條要人判）'
        )
        print("\n".join(lines))
        if stale:
            print(
                f"\n⚠ {len(stale)} 條的 verify 已經不成立——"
                "回查對應 Issue 並更新 JSON："
            )
            for item_id, why in stale:
                print(f"  {item_id}：{why}")
    return 1 if stale else 0


def selftest():
    """驗證各種訊號的正反兩面，以及現行 JSON 的 Issue 對照。"""
    import tempfile

    ok = True

    def check(name, condition):
        nonlocal ok
        print(f'  {"✓" if condition else "✗"} {name}')
        ok = ok and condition

    with tempfile.TemporaryDirectory(dir=str(ROOT / "workplace")) as temp_dir:
        directory = Path(temp_dir)
        (directory / "sub").mkdir()
        (directory / "sub" / "prod.go").write_text(
            "// 這一段 remake 還沒接\n", encoding="utf-8"
        )
        (directory / "sub" / "prod_test.go").write_text(
            "// TypeThatOnlyTestsMention\n", encoding="utf-8"
        )
        relative = str(directory.relative_to(ROOT))

        present = {
            "kind": "present",
            "paths": [relative + "/sub"],
            "pattern": "remake 還沒接",
        }
        check("present：自承還在 → 仍未完成", still_open({"verify": present})[0])
        (directory / "sub" / "prod.go").write_text("// 已經接上了\n", encoding="utf-8")
        check(
            "present：自承不見了 → 開口說可能已完成",
            not still_open({"verify": present})[0],
        )

        absent = {
            "kind": "absent",
            "paths": [relative + "/sub"],
            "pattern": "TypeThatOnlyTestsMention",
        }
        check(
            "absent：只有測試檔提到 → 仍未完成（測試檔不算）",
            still_open({"verify": absent})[0],
        )
        (directory / "sub" / "prod.go").write_text(
            "type TypeThatOnlyTestsMention int\n", encoding="utf-8"
        )
        check(
            "absent：產品碼出現了 → 開口說可能已完成",
            not still_open({"verify": absent})[0],
        )

        (directory / "n.json").write_text('{"rows": [1, 2]}', encoding="utf-8")
        json_length = {
            "kind": "json_len",
            "path": relative + "/n.json",
            "field": "rows",
            "max": 2,
        }
        check(
            "json_len：還沒超過門檻 → 仍未完成",
            still_open({"verify": json_length})[0],
        )
        (directory / "n.json").write_text('{"rows": [1, 2, 3]}', encoding="utf-8")
        check(
            "json_len：超過門檻 → 開口說可能已完成",
            not still_open({"verify": json_length})[0],
        )

        check("manual：一律回仍未完成", still_open({"verify": {"kind": "manual"}})[0])

        (directory / "worklist.json").write_text("remake 還沒接", encoding="utf-8")
        (directory / "sub" / "prod.go").write_text("// 已經接上了\n", encoding="utf-8")
        selfhit = {
            "kind": "present",
            "paths": [relative],
            "pattern": "remake 還沒接",
        }
        check(
            "不掃 worklist.json 自己（否則每一條都永遠成立）",
            not still_open({"verify": selfhit})[0],
        )

    data = load()
    check("docs/worklist.json 讀得動而且 schema 過關", bool(data["items"]))
    check(
        "每個現行工作都掛有 GitHub Issue",
        all(item.get("github_issue", {}).get("number", 0) > 0 for item in data["items"]),
    )
    print("worklist 自我測試：" + ("通過" if ok else "**失敗**"))
    return 0 if ok else 1


def main(argv):
    if "--selftest" in argv:
        return selftest()
    if len(argv) < 2 or argv[1] != "verify":
        if len(argv) >= 2 and argv[1] == "render":
            print("render 已移除：目前工作由 GitHub Issues 管理，Markdown 不再保存未處理清單。")
            return 2
        print(__doc__)
        return 2
    return cmd_verify(load())


if __name__ == "__main__":
    sys.exit(main(sys.argv))
