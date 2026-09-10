#!/usr/bin/env python3
"""未完成項的權威是 `docs/worklist.json`；`WORKLIST.md` 的那一節由這支產生。

    tools/py.sh tools/worklist.py verify     # 逐條問「這一條還成立嗎」
    tools/py.sh tools/worklist.py render     # 把未完成項寫回 WORKLIST.md
    tools/py.sh tools/worklist.py render --check  # 只比對，不寫檔（check.sh 用）
    tools/py.sh tools/worklist.py --selftest # 正反對照

⭐ **約定：`verify` 跑起來為真 ＝ 這一條仍然未完成。** 為假就是東西做好了
而條目沒跟著改——那正是要抓的過期斷言（`rulebook/61`、`rulebook/63`）。

markdown 的 `- [ ]` 清單會長出過期斷言：東西做好了，而沒有人回頭改那一條。
症狀不是報錯，是清單上留著一句自信的「還沒接」，然後有人照它去重做一遍、
或拿它當「還剩多少」的依據。它活得久是因為**沒有任何機制會問這一條還成不成立**。

四種 verify：

| kind | 語意 | 綁什麼 |
|---|---|---|
| `present` | pattern 找得到 → 仍未完成 | 程式碼或文件裡的**自承**（「remake 還沒接…」）|
| `absent` | pattern 找不到 → 仍未完成 | 東西還沒出現（型別名、函式名、檔名）|
| `json_len` | 某份 JSON 的欄位長度 ≤ max → 仍未完成 | 進度型的清單（校訂幾則、抽樣幾項）|
| `manual` | **一律回「仍未完成」並標出來** | 真的沒有機器訊號的（實機驗收、兩版並排畫面）|

⚠ **不掃測試檔。** 測試本來就會提到還沒接上的東西——為了釘住將來的行為。
把 `*_test.go` 算進來，`absent` 會因為測試裡有一行呼叫就判成「已經做了」。

⚠ **不掃 `docs/worklist.json` 自己。** 條目的 body 幾乎一定含 pattern 的字樣，
掃進去等於自己中自己，每一條都會永遠說「仍未完成」。
"""

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
DATA_PATH = ROOT / "docs/worklist.json"
RENDER_TARGET = ROOT / "WORKLIST.md"
BEGIN = "<!-- worklist:begin 由 tools/worklist.py render 產生，不要手改 -->"
END = "<!-- worklist:end -->"

SCANNED_SUFFIXES = {".go", ".py", ".sh", ".md", ".json", ".idc"}
SELF = "worklist.json"


def is_test_file(p: Path) -> bool:
    """測試檔不算——它們本來就會提到還沒接上的東西。"""
    n = p.name
    return "_test." in n or n.startswith("test_") or n.endswith("_test.py")


def scan_files(paths):
    for target in paths:
        p = ROOT / target
        if p.is_dir():
            files = sorted(p.rglob("*"))
        else:
            files = [p]
        for f in files:
            if not f.is_file() or f.suffix not in SCANNED_SUFFIXES:
                continue
            if is_test_file(f) or f.name == SELF:
                continue
            yield f


def find_hit(verify):
    """回第一個命中的檔案（相對路徑），沒有就回 None。"""
    expr = re.compile(verify["pattern"])
    for f in scan_files(verify["paths"]):
        try:
            text = f.read_text(encoding="utf-8", errors="ignore")
        except OSError:
            continue
        if expr.search(text):
            return str(f.relative_to(ROOT))
    return None


def still_open(item):
    """回 (仍未完成, 說明)。"""
    v = item["verify"]
    kind = v["kind"]
    if kind == "manual":
        return True, "要人判（沒有機器訊號）"
    if kind == "json_len":
        blob = json.loads((ROOT / v["path"]).read_text(encoding="utf-8"))
        field = blob[v["field"]] if v.get("field") else blob
        n = len(field)
        return n <= v["max"], f'{v["path"]} 有 {n} 項（門檻 {v["max"]}）'
    where = find_hit(v)
    if kind == "present":
        if where:
            return True, f"自承還在 {where}"
        return False, f'找不到 /{v["pattern"]}/'
    if kind == "absent":
        if where:
            return False, f"已經出現在 {where}"
        return True, f'還沒出現 /{v["pattern"]}/'
    raise SystemExit(f'不認得的 verify kind：{kind}（條目 {item["id"]}）')


def load(path=DATA_PATH):
    data = json.loads(Path(path).read_text(encoding="utf-8"))
    seen = set()
    for item in data["items"]:
        for key in ("id", "layer", "title", "body", "acceptance", "verify"):
            if key not in item:
                raise SystemExit(f'條目 {item.get("id", "?")} 缺 {key}')
        if item["id"] in seen:
            raise SystemExit(f'條目 id 重複：{item["id"]}')
        seen.add(item["id"])
        if item["layer"] not in data["layers"]:
            raise SystemExit(f'條目 {item["id"]} 的 layer 不在 layers 裡')
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
        lines.append(f'  {item["id"]:<34} {mark:<14} {why}')
    if not quiet:
        print(f'worklist：{len(data["items"])} 條未完成項'
              f'（其中 {manual} 條要人判）')
        print("\n".join(lines))
        if stale:
            print(f"\n⚠ {len(stale)} 條的 verify 已經不成立——"
                  f"東西可能做好了而條目沒改：")
            for id_, why in stale:
                print(f"  {id_}：{why}")
    return 1 if stale else 0


def render_text(data):
    out = [BEGIN, "",
           f'共 **{len(data["items"])} 條**未完成項。權威是 '
           f'[`docs/worklist.json`](docs/worklist.json)，'
           f'每一條掛一個 verify——**跑起來為真就是這一條仍然未完成**。',
           "", "跑 `tools/py.sh tools/worklist.py verify` 逐條問一次；"
           "`check.sh` 會替你跑。", ""]
    by_layer = {}
    for item in data["items"]:
        by_layer.setdefault(item["layer"], []).append(item)
    for layer, desc in data["layers"].items():
        items = by_layer.get(layer)
        if not items:
            continue
        out.append(f"### {layer} — {desc}")
        out.append("")
        for item in items:
            open_, why = still_open(item)
            mark = "" if open_ else "（⚠ verify 已不成立，回頭看這一條）"
            out.append(f'#### {item["title"]} {mark}'.rstrip())
            out.append("")
            out.append(item["body"])
            out.append("")
            if item.get("blocked_by"):
                out.append(f'**卡在**：{item["blocked_by"]}')
                out.append("")
            out.append(f'**怎樣算做完**：{item["acceptance"]}')
            out.append("")
            v = item["verify"]
            if v["kind"] == "manual":
                out.append(f'**verify**：`manual` — {v.get("note", "要人判")}')
            elif v["kind"] == "json_len":
                out.append(f'**verify**：`json_len` `{v["path"]}` ≤ {v["max"]}')
            else:
                out.append(f'**verify**：`{v["kind"]}` `/{v["pattern"]}/` '
                           f'在 `{"`、`".join(v["paths"])}`')
            out.append("")
    out.append(END)
    return "\n".join(out)


def cmd_render(data, check=False):
    text = RENDER_TARGET.read_text(encoding="utf-8")
    block = render_text(data)
    if BEGIN not in text or END not in text:
        raise SystemExit(
            f"{RENDER_TARGET.name} 裡找不到 render 標記；"
            f"請先放一對：\n{BEGIN}\n{END}")
    head = text[:text.index(BEGIN)]
    tail = text[text.index(END) + len(END):]
    updated = head + block + tail
    if check:
        # ⚠ **markdown 與 JSON 不同步，就是過期斷言的溫床**：有人手改了那一節，
        # 或改完 JSON 忘了 render。這一關讓它當場開口。
        if updated != text:
            print(f"⚠ {RENDER_TARGET.name} 的未完成項那一節與 "
                  f"docs/worklist.json 不同步——跑 "
                  f"`tools/py.sh tools/worklist.py render`")
            return 1
        print(f'{RENDER_TARGET.name} 與 docs/worklist.json 同步'
              f'（{len(data["items"])} 條）')
        return 0
    RENDER_TARGET.write_text(updated, encoding="utf-8")
    print(f'寫回 {RENDER_TARGET.name}：{len(data["items"])} 條')
    return 0


def selftest():
    """⚠ 正反對照：只驗「它印出仍未完成」證明不了機制有在看。

    每一種 kind 都要先確認**訊號在時報未完成**，再把訊號拿掉、
    確認它真的開口。
    """
    import tempfile

    ok = True

    def check(name, cond):
        nonlocal ok
        print(f'  {"✓" if cond else "✗"} {name}')
        ok = ok and cond

    with tempfile.TemporaryDirectory(dir=str(ROOT / "workplace")) as td:
        d = Path(td)
        (d / "sub").mkdir()
        (d / "sub" / "prod.go").write_text("// 這一段 remake 還沒接\n",
                                           encoding="utf-8")
        (d / "sub" / "prod_test.go").write_text("// TypeThatOnlyTestsMention\n",
                                                encoding="utf-8")
        rel = str(d.relative_to(ROOT))

        present = {"kind": "present", "paths": [rel + "/sub"],
                   "pattern": "remake 還沒接"}
        check("present：自承還在 → 仍未完成",
              still_open({"verify": present})[0])
        (d / "sub" / "prod.go").write_text("// 已經接上了\n", encoding="utf-8")
        check("present：自承不見了 → 開口說可能已完成",
              not still_open({"verify": present})[0])

        absent = {"kind": "absent", "paths": [rel + "/sub"],
                  "pattern": "TypeThatOnlyTestsMention"}
        check("absent：只有測試檔提到 → 仍未完成（測試檔不算）",
              still_open({"verify": absent})[0])
        (d / "sub" / "prod.go").write_text("type TypeThatOnlyTestsMention int\n",
                                           encoding="utf-8")
        check("absent：產品碼出現了 → 開口說可能已完成",
              not still_open({"verify": absent})[0])

        (d / "n.json").write_text('{"rows": [1, 2]}', encoding="utf-8")
        jl = {"kind": "json_len", "path": rel + "/n.json",
              "field": "rows", "max": 2}
        check("json_len：還沒超過門檻 → 仍未完成",
              still_open({"verify": jl})[0])
        (d / "n.json").write_text('{"rows": [1, 2, 3]}', encoding="utf-8")
        check("json_len：超過門檻 → 開口說可能已完成",
              not still_open({"verify": jl})[0])

        check("manual：一律回仍未完成",
              still_open({"verify": {"kind": "manual"}})[0])

        # ⚠ 自己中自己：條目的 body 幾乎一定含 pattern 的字樣。
        (d / "worklist.json").write_text("remake 還沒接", encoding="utf-8")
        (d / "sub" / "prod.go").write_text("// 已經接上了\n", encoding="utf-8")
        selfhit = {"kind": "present", "paths": [rel],
                   "pattern": "remake 還沒接"}
        check("不掃 worklist.json 自己（否則每一條都永遠成立）",
              not still_open({"verify": selfhit})[0])

    data = load()
    check("docs/worklist.json 讀得動而且 schema 過關", bool(data["items"]))
    check("render 產得出來而且含首尾標記",
          render_text(data).startswith(BEGIN) and render_text(data).endswith(END))
    print("worklist 自我測試：" + ("通過" if ok else "**失敗**"))
    return 0 if ok else 1


def main(argv):
    if "--selftest" in argv:
        return selftest()
    if len(argv) < 2 or argv[1] not in {"verify", "render"}:
        print(__doc__)
        return 2
    data = load()
    if argv[1] == "verify":
        return cmd_verify(data)
    return cmd_render(data, check="--check" in argv)


if __name__ == "__main__":
    sys.exit(main(sys.argv))
