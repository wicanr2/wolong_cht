#!/usr/bin/env python3
"""把散在各文件的「未解」表抽成一張總表。

    tools/py.sh tools/re_open_questions.py > docs/re/43-open-questions.md

回答的問題是「**還有什麼沒解**」。既有的三張表都不回答它：

    docs/INDEX.md   已解的斷言（欄位／常數 → 推論等級 → 出處）
    docs/re/21      函式有沒有人寫過（覆蓋率）
    docs/re/24      未讀函式大概在做什麼（角色證據）

缺口沒有總表的代價是實際發生過的：2026-08-14 重推了一次
`docs/re/11` §3.5 早就解完的戰術腳本 VM，還推錯兩處，
因為「還沒解」與「我不記得解過」在動手那一刻長得一模一樣。

## 抽取規則（只認結構，不猜語意）

1. **專門的**「未解」小節，其後的表格逐列抽出。專門的定義是標題去掉編號之後
   以未解／未定／缺口收尾且不長（≤ 12 字）——「段 1：七個圖塊已定位，UI 語意未解」
   這種**描述性標題底下擺的是已解內容**，整節抽進來會得到一堆假缺口。
2. **任何**表格裡最後一欄含「未解／假說／未定」的列。欄位表就是這個形狀
   （`| +0x16 | u16 | 恆 FFFF | 未解 |`）。
3. 收尾是「…未解／未定案／待驗／未驗」的**散句**。`docs/mechanics/` 與各文件的
   狀態行幾乎都是這個形狀而不是表——少了這一條，最該被看見的那一類會整個消失
   （第一版只認表，`規則正確性` 只抽到 1 條，而那一類的缺口最多）。
4. 列首是 `~~刪除線~~`、整列出現 ✅、那是表頭列、或句子把「未解」放在引號裡
   （談的是缺口這件事本身，不是一條缺口）的，跳過。

「擋住什麼」由**來源目錄**決定，不是猜的：`docs/mechanics/` 擋規則正確性、
`docs/formats/` 擋資料保存、`docs/re/` 擋程式碼理解。
「怎麼裁決」由關鍵字決定，只有三種值，判不出來就寫靜態。

## 為什麼結尾要印「抽不到的檔案」

只印抽得到的，會讓解析失敗長得像「那份文件沒有缺口」。
末段列出**提到未解卻一列都沒抽到**的檔案，那是這支工具自己的盲區清單。
"""
import datetime
import io
import os
import re
import sys

DOC_ROOTS = ("docs",)
SELF = "docs/re/43-open-questions.md"
SKIP = (SELF, "docs/INDEX.md")

HEADING = re.compile(r"^(#{2,4})\s+(.*)$")
OPEN_HEAD = re.compile(r"未解|未定|缺口|開放問題")
# 標題編號與裝飾，比對「是不是專門的未解小節」之前要先剝掉
DECOR = re.compile(r"^[\d.\s]*|[⭐⚠🔵✅*\s]+")
DEDICATED = re.compile(r"^.{0,10}(未解|未定案|未定|缺口|開放問題|未解範圍)$")
# 散句形式的缺口：`**未解**：…`、以及收尾是「…未解」的句子
PROSE = re.compile(r"\*\*未解\*\*[：:]\s*(.+)$|^未解[：:]\s*(.+)$")
TAIL = re.compile(r"(未解|未定案|仍未解|待驗|未驗)[。，）\)]?$")
# 談「未解」這件事本身，不是一條缺口
META = re.compile(r"「未解」|未解表|未解的表|見下方|見上方|見 §|參見")
SENT = re.compile(r"[^。\n]+[。]?")
SOLVED = re.compile(r"✅|已解|已定案")
STRUCK = re.compile(r"^~~")
SEP = re.compile(r"^\|[\s:|-]+\|$")

# 「怎麼裁決」——判不出來就是靜態。順序有意義：實測優先於兩版對照，
# 因為兩版都一樣的東西還是可能要跑起來才知道語意。
VERDICT = (
    (re.compile(r"實測|實跑|跑起來|畫出來|oracle|DOSBox|截圖"), "實測"),
    (re.compile(r"兩版|PC-98 對照"), "兩版對照"),
)
BLOCKED = re.compile(r"防拷|密碼畫面")

DOMAIN = (
    ("docs/mechanics/", "規則正確性"),
    ("docs/formats/", "資料保存"),
    ("docs/playtest/", "驗收"),
    ("docs/re/", "程式碼理解"),
    ("docs/reference/", "外部資料"),
)


def domain_of(rel):
    for prefix, name in DOMAIN:
        if rel.startswith(prefix):
            return name
    return "其他"


def verdict_of(text):
    for pat, name in VERDICT:
        if pat.search(text):
            return name
    return "靜態"


def cells(line):
    if not line.startswith("|") or SEP.match(line):
        return None
    parts = [c.strip() for c in line.strip().strip("|").split("|")]
    return parts if len(parts) >= 2 else None


def collect(path, rel):
    """回傳 (items, saw_keyword)。items 是 (項目, 現況, 小節) 的清單。"""
    text = io.open(path, encoding="utf-8").read()
    lines = text.split("\n")
    items = []
    section = ""          # 目前所在的小節標題
    open_level = None     # 「未解」小節的層級；None ＝ 不在未解小節裡
    prev = ""             # 上一行非空的內文，給跨行句子補前綴
    for i, line in enumerate(lines):
        m = HEADING.match(line)
        if m:
            level, title = len(m.group(1)), m.group(2).strip()
            if open_level is not None and level <= open_level:
                open_level = None
            bare = DECOR.sub("", title).strip()
            if DEDICATED.match(bare) and not SOLVED.search(title):
                open_level, section = level, title
            elif open_level is None:
                section = title
            continue
        pm = PROSE.search(line)
        if pm and not SOLVED.search(line):
            items.append((pm.group(1) or pm.group(2), "（散句）", section))
            continue
        if not line.startswith(("|", ">", "```", "    ")):
            for sm in SENT.finditer(line):
                sent = sm.group(0).strip(" -*`")
                if len(sent) < 6 or SOLVED.search(sent) or META.search(sent):
                    continue
                if TAIL.search(sent.rstrip("*）) ")):
                    # 句子跨行時抽到的會是尾巴（「→ 假說，待驗」），
                    # 單獨看沒有資訊。補上一行當前綴。
                    if len(sent) < 16 and prev:
                        sent = prev[-60:] + " " + sent
                    items.append((sent.strip("*"), "（散句）", section))
        if line.strip():
            prev = line.strip(" -*`>")
        c = cells(line)
        if not c:
            continue
        # 表頭列：下一行是分隔線
        if i + 1 < len(lines) and SEP.match(lines[i + 1].strip()):
            continue
        if SOLVED.search(line) or STRUCK.match(c[0]):
            continue
        if open_level is not None:
            items.append((c[0], " / ".join(c[1:]), section))
        elif OPEN_HEAD.search(c[-1]) and not SOLVED.search(c[-1]):
            # 欄位表：最後一欄標未解的列
            items.append((c[0], " / ".join(c[1:]), section))
    return items, bool(OPEN_HEAD.search(text))


def main():
    repo = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    files = []
    for root in DOC_ROOTS:
        for dirpath, dirnames, filenames in os.walk(os.path.join(repo, root)):
            dirnames[:] = [d for d in dirnames if d not in ("images", ".git")]
            for fn in sorted(filenames):
                if fn.endswith(".md"):
                    files.append(os.path.join(dirpath, fn))
    files.sort()

    rows, silent = [], []
    for path in files:
        rel = os.path.relpath(path, repo).replace(os.sep, "/")
        if rel in SKIP:
            continue
        items, saw = collect(path, rel)
        if not items:
            if saw:
                silent.append(rel)
            continue
        for item, status, section in items:
            rows.append({
                "file": rel,
                "item": item,
                "status": status,
                "section": section,
                "domain": domain_of(rel),
                "verdict": verdict_of(item + " " + status),
                "blocked": bool(BLOCKED.search(item + " " + status)),
            })

    out = sys.stdout
    w = out.write
    w("# 43 — 未解缺口總表（生成的）\n\n")
    w("**狀態：生成的清單，跑 `tools/py.sh tools/re_open_questions.py` 重出。\n")
    w("這一份不下結論，只把各文件的「未解」表集中到一處。**\n\n")
    # 生成日期用今天：這份是重跑就重出的東西，寫死日期等於謊報新鮮度。
    w("- 日期：%s\n" % datetime.date.today().isoformat())
    w("- 產生工具：`tools/re_open_questions.py`\n")
    w("- 來源：`docs/` 底下所有文件的未解小節、表格裡標未解的列，"
      "與收尾是「…未解」的散句\n\n")
    w("既有的三張表回答別的問題（`CLAUDE.md` §10）：`docs/INDEX.md` 是**已解**的斷言、\n")
    w("[`21`](21-function-census.md) 是函式有沒有人寫過、"
      "[`24`](24-unread-function-catalogue.md) 是未讀函式在做什麼。\n")
    w("**這一份是唯一回答「還有什麼沒解」的。**\n\n")
    w("> 「擋住什麼」由來源目錄決定，「怎麼裁決」由關鍵字決定——"
      "兩欄都是機械算出來的，\n> **不是逐條判斷過的優先序**。"
      "要排優先序請自己讀那一列指到的小節。\n\n")

    w("## 1. 總量\n\n")
    w("| 擋住什麼 | 缺口數 | 靜態可解 | 要實測 | 兩版對照 |\n")
    w("|---|---:|---:|---:|---:|\n")
    order = ["規則正確性", "資料保存", "程式碼理解", "驗收", "外部資料", "其他"]
    for dom in order:
        sel = [r for r in rows if r["domain"] == dom]
        if not sel:
            continue
        c = lambda v: sum(1 for r in sel if r["verdict"] == v)
        w(f"| {dom} | {len(sel)} | {c('靜態')} | {c('實測')} | {c('兩版對照')} |\n")
    nb = sum(1 for r in rows if r["blocked"])
    w(f"| **合計** | **{len(rows)}** | "
      f"{sum(1 for r in rows if r['verdict']=='靜態')} | "
      f"{sum(1 for r in rows if r['verdict']=='實測')} | "
      f"{sum(1 for r in rows if r['verdict']=='兩版對照')} |\n\n")
    if nb:
        w(f"其中 **{nb} 條明講被防拷擋著**——那條路沒通之前，"
          "它們不會因為多讀組語而前進。\n\n")

    for dom in order:
        sel = [r for r in rows if r["domain"] == dom]
        if not sel:
            continue
        w(f"## 2.{order.index(dom)+1} {dom}（{len(sel)} 條）\n\n")
        w("| 出處 | 缺口 | 現況 | 裁決 |\n|---|---|---|---|\n")
        for r in sorted(sel, key=lambda r: (r["file"], r["section"])):
            link = "../" + r["file"][len("docs/"):]
            v = r["verdict"] + ("・**防拷擋著**" if r["blocked"] else "")
            st = r["status"].replace("\n", " ")
            if len(st) > 160:
                st = st[:157] + "…"
            w(f"| [`{r['file'][len('docs/'):]}`]({link}) | {r['item']} | {st} | {v} |\n")
        w("\n")

    w("## 3. 這支工具的盲區\n\n")
    w("抽取只認四種結構（專門的未解小節、表格最後一欄標未解的列、\n")
    w("`**未解**：…`、收尾是「…未解」的句子）。**寫在段落中段、\n")
    w("或用別的詞說「這個還不知道」的缺口抽不到**——下列檔案提到未解\n")
    w("卻一列都沒抽出來，要嘛缺口寫成別的句式，要嘛那些字樣只是在講別的事：\n\n")
    if silent:
        for rel in silent:
            w(f"- `{rel}`\n")
    else:
        w("（沒有）\n")
    w("\n只印抽得到的部分，會讓解析失敗長得像「那份文件沒有缺口」。\n")
    w("這一節就是為了讓那個差別看得見。\n")


if __name__ == "__main__":
    main()
