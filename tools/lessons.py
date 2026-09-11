#!/usr/bin/env python3
"""教訓的權威是 `docs/lessons.json`；`docs/lessons.md` 由這支產生。

    tools/py.sh tools/lessons.py verify     # 每條教訓的**防線**還在嗎
    tools/py.sh tools/lessons.py render     # 產生按觸發點索引的表
    tools/py.sh tools/lessons.py render --check
    tools/py.sh tools/lessons.py --selftest

⚠ **verify 的語意與 `worklist.py` 相反。** worklist 問「這一條還未完成嗎」
（為真 ＝ 仍未完成）；這裡問「**這一條的防線還在嗎**」（為真 ＝ 防線還在）。
防線不見了才是要抓的東西——規則被順手刪掉、工具被移走、測試被改名。

⭐ 為什麼要有這一份：`CLAUDE.md` §4.0 定過規矩——**會復發的斷言要同時有
規則與檢查，只有其中一個都不夠**。而教訓散在 `CONTEXT.md` §6 與
`CLAUDE.md` §7 兩處、寫成**事件敘述**，於是同一種錯誤換一個場景就認不出來。
這一份把它們去重成**模式**，每條掛三樣東西：

| 欄位 | 問的是 |
|---|---|
| `trigger` | **什麼時候要想起這一條**——按這個索引，中途也查得到 |
| `occurrences` | 犯過幾次、在哪。**≥ 2 次還只有規則就是該升級的訊號** |
| `guard` | 防線是工具、測試、還是只有文字規則 |

四種防線，由硬到軟：

| kind | 語意 |
|---|---|
| `tool` | 有工具會在 `check.sh` 每次問一次 |
| `test` | 有測試釘住 |
| `remind` | **工具會在對的時機把它印出來**——擋不住，但保證看得到 |
| `rule` | 只有文字規則，人要自己記得 |
| `none` | 連規則都還沒寫 |

⭐⭐ **`remind` 是「一定能讀到」的那一層。** 光把教訓寫進常駐文件不夠——
`CLAUDE.md` §7 的 34 條每個 session 都載入，而它們要防的動作（讀一段組語、
引用一個數字、下一個規則結論）發生在**任務中間**，那時沒有任何東西會問。
所以 `render` 會把每條教訓依 `remind_on` 攤成
`docs/lessons/<key>.txt`，各工具開頭 `cat` 它——**跑 `tools/ida.sh` 就看得到
反組譯的那幾條**，零成本、不必起 docker。
"""

import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from worklist import find_hit  # noqa: E402  掃描邏輯只留一份

ROOT = Path(__file__).resolve().parent.parent
DATA_PATH = ROOT / "docs/lessons.json"
RENDER_TARGET = ROOT / "docs/lessons.md"
BEGIN = "<!-- lessons:begin 由 tools/lessons.py render 產生，不要手改 -->"
END = "<!-- lessons:end -->"

# 復發這麼多次還只有文字規則，就該升級成工具或測試。
UPGRADE_AT = 2


def guard_alive(lesson):
    """回 (防線還在, 說明)。"""
    g = lesson["guard"]
    if g["kind"] == "none":
        return False, "還沒有防線"
    v = g.get("verify")
    if not v:
        return False, "guard 沒有 verify——說不出防線在哪"
    where = find_hit(v)
    if where:
        return True, f'防線在 {where}'
    return False, f'⚠ 防線不見了：找不到 /{v["pattern"]}/'


def needs_upgrade(lesson):
    return (lesson["guard"]["kind"] == "rule"
            and len(lesson["occurrences"]) >= UPGRADE_AT)


def load(path=DATA_PATH):
    data = json.loads(Path(path).read_text(encoding="utf-8"))
    seen = set()
    for les in data["lessons"]:
        for key in ("id", "title", "trigger", "rule", "occurrences", "guard"):
            if key not in les:
                raise SystemExit(f'教訓 {les.get("id", "?")} 缺 {key}')
        if les["id"] in seen:
            raise SystemExit(f'教訓 id 重複：{les["id"]}')
        seen.add(les["id"])
        if les["guard"]["kind"] not in data["guard_kinds"]:
            raise SystemExit(f'教訓 {les["id"]} 的 guard.kind 不認得')
        if not les["occurrences"]:
            raise SystemExit(f'教訓 {les["id"]} 沒有 occurrences——'
                             f'沒犯過的不是教訓，是規則')
    return data


def cmd_verify(data):
    dead, upgrade = [], []
    print(f'教訓：{len(data["lessons"])} 條')
    for les in data["lessons"]:
        alive, why = guard_alive(les)
        n = len(les["occurrences"])
        if not alive:
            dead.append((les["id"], why))
        if needs_upgrade(les):
            upgrade.append((les["id"], n))
        mark = les["guard"]["kind"]
        print(f'  {les["id"]:<34} {mark:<5} 犯過 {n} 次  {why}')
    if upgrade:
        print(f"\n⚠ {len(upgrade)} 條復發過 {UPGRADE_AT} 次以上，"
              f"而防線**只有文字規則**——該升級成工具或測試：")
        for id_, n in upgrade:
            print(f"  {id_}（{n} 次）")
    if dead:
        print(f"\n❌ {len(dead)} 條的防線不見了：")
        for id_, why in dead:
            print(f"  {id_}：{why}")
    # ⭐ 防線不見了是 gate；「該升級」只提醒——升級是人的決定。
    return 1 if dead else 0


def render_text(data):
    out = [BEGIN, "",
           f'共 **{len(data["lessons"])} 條**。權威是 '
           f'[`lessons.json`](lessons.json)，這一份由 '
           f'`tools/py.sh tools/lessons.py render` 產生。', "",
           "⭐ **按「什麼時候要想起這一條」索引**——教訓要防的動作發生在"
           "任務中間，不是開始時；按事件時間排的清單那時查不到。", ""]
    out.append("| 什麼時候想起 | 這一條 | 犯過 | 防線 |")
    out.append("|---|---|---:|---|")
    for les in data["lessons"]:
        alive, _ = guard_alive(les)
        g = les["guard"]
        badge = g["kind"]
        if not alive:
            badge += "（❌ 不見了）"
        elif needs_upgrade(les):
            badge += "（⚠ 該升級）"
        out.append(f'| {les["trigger"]} | [{les["title"]}](#{les["id"]}) '
                   f'| {len(les["occurrences"])} | {badge} |')
    out.append("")
    for les in data["lessons"]:
        out.append(f'### {les["title"]}')
        out.append("")
        out.append(f'<a id="{les["id"]}"></a>**什麼時候想起**：{les["trigger"]}')
        out.append("")
        out.append(f'**要做的**：{les["rule"]}')
        out.append("")
        out.append("| 日期 | 犯在哪 | 收據 |")
        out.append("|---|---|---|")
        for o in les["occurrences"]:
            out.append(f'| {o["date"]} | {o["what"]} | {o.get("doc", "—")} |')
        out.append("")
        g = les["guard"]
        out.append(f'**防線**：`{g["kind"]}` — {g.get("note", "")}'.rstrip())
        out.append("")
    out.append(END)
    # ⭐ `tools/index.py` 的盲區檢查會看到教訓正文裡的「未解」「還沒」，
    # 而這一份的缺口是**每條教訓自己的 `guard`**，不是一個「未解」小節。
    # 明講一次，免得每加一條教訓就多一個假警報。
    out.append("")
    out.append("<!-- 缺口：無 -->")
    return "\n".join(out)


# ⚠ `tools/index.py` 要求每份 docs/*.md 都有狀態行與日期，所以檔頭是固定的，
# `render` 只換 BEGIN..END 之間。
HEAD = """# 教訓：按觸發點索引

**狀態：由 `tools/lessons.py render` 產生，不要手改。**
權威是 [`lessons.json`](lessons.json)；`tools/ida.sh`／`tools/dosgolem.sh`
開頭會把相關的那幾條印出來。

- 日期：2026-09-10
- 規則：`~/.claude/rulebook/61-worklist-as-data-with-verify.md`（同一個做法用在教訓上）
- 相關：[`CONTEXT.md`](../CONTEXT.md) §6（事件台帳）、`CLAUDE.md` §7（歷史清單）

"""

REMIND_DIR = ROOT / "docs/lessons"
REMIND_MAX = 5


def remind_text(data, key):
    """某個工具該印的那幾條——**最多 5 條，按犯過次數排**。

    ⚠ 每次都印一大串等於沒印。寧可少而準。
    """
    picked = [l for l in data["lessons"] if key in l.get("remind_on", [])]
    picked.sort(key=lambda l: -len(l["occurrences"]))
    lines = []
    for les in picked[:REMIND_MAX]:
        n = len(les["occurrences"])
        lines.append(f'  · {les["title"]}（犯過 {n} 次）')
        lines.append(f'    {les["rule"]}')
    if not lines:
        return ""
    return ("⚠ 這一類任務犯過的（docs/lessons.md）：\n"
            + "\n".join(lines) + "\n")


def write_reminders(data, check=False):
    """把 remind 檔攤出來；check 模式只比對。"""
    keys = sorted({k for l in data["lessons"] for k in l.get("remind_on", [])})
    REMIND_DIR.mkdir(parents=True, exist_ok=True)
    stale = []
    for key in keys:
        path = REMIND_DIR / f"{key}.txt"
        want = remind_text(data, key)
        if check:
            have = path.read_text(encoding="utf-8") if path.exists() else ""
            if have != want:
                stale.append(path.name)
        else:
            path.write_text(want, encoding="utf-8")
    return keys, stale


def cmd_render(data, check=False):
    if not RENDER_TARGET.exists():
        RENDER_TARGET.write_text(HEAD + f"{BEGIN}\n{END}\n", encoding="utf-8")
    text = RENDER_TARGET.read_text(encoding="utf-8")
    block = render_text(data)
    if BEGIN not in text or END not in text:
        raise SystemExit(f"{RENDER_TARGET.name} 裡找不到 render 標記")
    updated = (text[:text.index(BEGIN)] + block
               + text[text.index(END) + len(END):])
    keys, stale = write_reminders(data, check)
    if check:
        if updated != text or stale:
            what = [] if updated == text else [RENDER_TARGET.name]
            print(f'⚠ {"、".join(what + stale)} 與 docs/lessons.json 不同步'
                  f'——跑 `tools/py.sh tools/lessons.py render`')
            return 1
        print(f'{RENDER_TARGET.name} 與 {len(keys)} 份提示檔都同步'
              f'（{len(data["lessons"])} 條）')
        return 0
    RENDER_TARGET.write_text(updated, encoding="utf-8")
    print(f'寫回 {RENDER_TARGET.name} 與 {len(keys)} 份提示檔'
          f'（{"、".join(keys)}）：{len(data["lessons"])} 條')
    return 0


def selftest():
    """⚠ 正反對照：只驗「它說防線還在」證明不了機制有在看。"""
    import tempfile

    ok = True

    def check(name, cond):
        nonlocal ok
        print(f'  {"✓" if cond else "✗"} {name}')
        ok = ok and cond

    with tempfile.TemporaryDirectory(dir=str(ROOT / "workplace")) as td:
        d = Path(td)
        (d / "r.md").write_text("先懷疑自我修改碼\n", encoding="utf-8")
        rel = str(d.relative_to(ROOT))
        les = {"guard": {"kind": "rule", "verify": {
            "kind": "present", "paths": [rel], "pattern": "先懷疑自我修改碼"}},
            "occurrences": [{"date": "x", "what": "y"}]}
        check("規則文字還在 → 防線還在", guard_alive(les)[0])
        (d / "r.md").write_text("（被順手刪掉了）\n", encoding="utf-8")
        check("規則文字被刪 → 防線不見了（開口）", not guard_alive(les)[0])

        check("guard.kind ＝ none → 一律當成沒有防線",
              not guard_alive({"guard": {"kind": "none"},
                               "occurrences": []})[0])
        check("guard 沒有 verify → 說不出防線在哪",
              not guard_alive({"guard": {"kind": "tool"},
                               "occurrences": []})[0])

        les["occurrences"] = [{"date": "a", "what": "1"}]
        check(f"只有規則、犯過 1 次 → 還不用升級", not needs_upgrade(les))
        les["occurrences"].append({"date": "b", "what": "2"})
        check(f"只有規則、犯過 {UPGRADE_AT} 次 → 該升級", needs_upgrade(les))
        les["guard"]["kind"] = "tool"
        check("已經是工具了 → 不再提醒升級", not needs_upgrade(les))

    data = load()
    check("docs/lessons.json 讀得動而且 schema 過關", bool(data["lessons"]))
    check("每條都有 occurrences（沒犯過的不是教訓）",
          all(l["occurrences"] for l in data["lessons"]))
    # ⚠ 判準不能用 `endswith(END)`：本體會在 END 之後再加一行缺口註腳，
    #   於是這一格從某一輪起恆為假——而 `check.sh` 在它之前的幽靈掃描
    #   就停住了，所以沒有人看到它紅。改成問「標記成對且順序正確」。
    text = render_text(data)
    check("render 產得出來而且首尾標記成對",
          text.startswith(BEGIN) and text.count(BEGIN) == 1
          and text.count(END) == 1 and text.index(END) > text.index(BEGIN))
    # ⚠ 提示檔要真的挑得出東西——空字串與「沒有這個 key」長得一樣。
    keys = sorted({k for l in data["lessons"] for k in l.get("remind_on", [])})
    check("每條教訓至少掛一個 remind_on（否則工具印不到它）",
          all(l.get("remind_on") for l in data["lessons"]))
    check("每個 remind key 都挑得出內容",
          bool(keys) and all(remind_text(data, k) for k in keys))
    check("不存在的 key 回空字串", remind_text(data, "沒有這個key") == "")
    print("lessons 自我測試：" + ("通過" if ok else "**失敗**"))
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
