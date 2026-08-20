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
# TEMPLATE.md 是骨架不是文件：它的「未解」小節是空表頭，
# 抽不到東西是正確行為，不是盲區。
SKIP = (SELF, "docs/INDEX.md", "docs/spec/TEMPLATE.md")

HEADING = re.compile(r"^(#{2,4})\s+(.*)$")
OPEN_CELL = re.compile(r"未解|未定|假說|未驗|未定位|尚未")
# 判斷「這份文件應該要有缺口」的字樣。刻意比 OPEN_CELL 窄：
# 「尚未」「未知」在正常敘述裡太常出現，放進來會讓盲區清單灌水成全庫。
SAW = re.compile(r"未解|未定案|未定位|缺口|開放問題|未驗")
# 標題編號與裝飾，比對「是不是專門的未解小節」之前要先剝掉
DECOR = re.compile(r"^[\d.\s]*|[⭐⚠🔵✅*\s]+")
DEDICATED = re.compile(r"^.{0,13}(未解|未定案|未定|缺口|開放問題|未解範圍|未完成項)$")
# 同一件事的其他寫法。標題不是只有「未解」一種講法——
# 「尚待補完的戰術分支與對拍」「還沒解的」都是缺口小節，漏認就整節看不見。
DEDICATED2 = re.compile(r"^(尚待|還沒解|待解|待補|還缺|剩下的缺口)")
# 未解小節裡的子標題：這兩組決定要不要繼續抓。
# 「已證實／強推論」底下擺的是結論不是缺口，抓進來會得到一堆假缺口。
SUB_STOP = re.compile(r"已證實|已解|已定案|強推論|已完成")
SUB_GO = re.compile(r"未知|未解|未定|缺口|待")
# 表頭就寫「缺口」的表，整張都是缺口清單
GAP_HEADER = re.compile(r"缺口|未解|待辦|下手點")
# 散句形式的缺口：`**未解**：…`、以及收尾是「…未解」的句子
PROSE = re.compile(r"\*\*未解\*\*[：:]\s*(.+)$|^未解[：:]\s*(.+)$")
TAIL = re.compile(r"(未解|未定案|仍未解|待驗|未驗|未定位)[。，）\)]?$")
# 談「未解」這件事本身，不是一條缺口
META = re.compile(r"「未解」|未解表|未解的表|見下方|見上方|見 §|參見")
SENT = re.compile(r"[^。\n]+[。]?")
# 原文裡的相對連結搬進這一份就會指錯地方（`../playtest/18` 從 docs/re/ 出發
# 才成立）。每一列本來就標了出處，連結留著只會讓連結檢查失敗——拆成純文字。
MDLINK = re.compile(r"\[([^\]]*)\]\([^)]*\)")
SOLVED = re.compile(r"✅|已解|已定案")
STRUCK = re.compile(r"^~~")
SEP = re.compile(r"^\|[\s:|-]+\|$")
# 文件明講自己沒有缺口。有這一行就不算盲區——
# 「解析不到」與「真的沒有」必須分得開，否則盲區清單永遠清不完。
NO_GAPS = re.compile(r"<!--\s*缺口：無\s*-->")
# 「…下列缺口：」之後接的條列就是缺口清單，不必另外開小節。
LEAD_IN = re.compile(r"(缺口|未解|還沒解|未完成)[：:]\s*$")

# 「怎麼裁決」——判不出來就是靜態。順序有意義：實測優先於兩版對照，
# 因為兩版都一樣的東西還是可能要跑起來才知道語意。
VERDICT = (
    (re.compile(r"實測|實跑|跑起來|畫出來|oracle|DOSBox|截圖"), "實測"),
    (re.compile(r"兩版|PC-98 對照"), "兩版對照"),
)

DOMAIN = (
    ("docs/mechanics/", "規則正確性"),
    ("docs/formats/", "資料保存"),
    ("docs/playtest/", "驗收"),
    ("docs/re/", "程式碼理解"),
    ("docs/reference/", "外部資料"),
)


# repo 根目錄實際存在的頂層目錄；顯示名撞到它們就不能砍 `docs/`。
TOP_LEVEL_DIRS = {
    "mobile", "android", "tools", "internal", "cmd",
    "translations", "packaging", "docker", "workplace", "dist-all",
}


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
    muted = False         # 在未解小節裡，但子標題是「已證實」那一類
    prev = ""             # 上一行非空的內文，給跨行句子補前綴
    gap_table = False     # 目前這張表的表頭是不是寫著「缺口」
    lead_done = set()     # 已經收過首句的小節
    for i, line in enumerate(lines):
        m = HEADING.match(line)
        if m:
            level, title = len(m.group(1)), m.group(2).strip()
            if open_level is not None and level <= open_level:
                open_level, muted = None, False
            elif open_level is not None:
                # 未解小節裡的子標題
                if SUB_STOP.search(title):
                    muted = True
                elif SUB_GO.search(title):
                    muted = False
            bare = DECOR.sub("", title).strip()
            if (DEDICATED.match(bare) or DEDICATED2.match(bare)) \
                    and not SOLVED.search(title):
                open_level, section, muted = level, title, False
            elif open_level is None:
                section = title
            continue
        if open_level is not None and not muted and not line.startswith("|"):
            # 未解小節裡的散文與條列。docs/re/15 §5 那種「證據分級與未解範圍」
            # 底下沒有表，只認表就會整節漏掉。
            #
            # 但**不能整節逐行收**：長敘述會被拆成幾十條假缺口，
            # 總數從 231 衝到 358 而資訊量沒有增加。所以條列全收、
            # 散文只收該小節的第一句當代表，細節留在原文。
            body = line.strip(" -*`")
            bullet = bool(re.match(r"^\s*([-*]|\d+\.)\s", line))
            if len(body) >= 8 and not body.startswith((">", "#", "```")) \
                    and not SOLVED.search(body) and not META.search(body) \
                    and (bullet or section not in lead_done):
                lead_done.add(section)
                items.append((SENT.search(body).group(0).strip()[:120],
                              "（未解小節內文）", section))
                prev = body
                continue
        if LEAD_IN.search(line.rstrip()) and not SOLVED.search(line):
            open_level, muted = 9, False   # 9 ＝ 比任何標題都深，下一個標題就關掉
            section = section or "缺口條列"
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
            # 分隔線 `|---|---|` 也會讓 cells() 回 None——那時**不能**把
            # gap_table 關掉，否則表頭剛認出來就被下一行清掉，整張表一列都收不到。
            if not line.startswith("|"):
                gap_table = False
            continue
        # 表頭列：下一行是分隔線。表頭自己寫「缺口」的，整張表都收。
        if i + 1 < len(lines) and SEP.match(lines[i + 1].strip()):
            gap_table = bool(GAP_HEADER.search(line))
            continue
        if SOLVED.search(line) or STRUCK.match(c[0]):
            continue
        if open_level is not None or gap_table:
            items.append((c[0], " / ".join(c[1:]), section))
        elif OPEN_CELL.search(c[-1]) and not SOLVED.search(c[-1]):
            # 欄位表：最後一欄標未解的列
            items.append((c[0], " / ".join(c[1:]), section))
    return items, bool(SAW.search(text)) and not NO_GAPS.search(text)


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
    w(f"| **合計** | **{len(rows)}** | "
      f"{sum(1 for r in rows if r['verdict']=='靜態')} | "
      f"{sum(1 for r in rows if r['verdict']=='實測')} | "
      f"{sum(1 for r in rows if r['verdict']=='兩版對照')} |\n\n")

    # ⚠ 這個合計是**列數**，不是獨立問題數。兩件事會讓它偏高，
    # 兩件都印出來讓讀的人自己扣：
    #   1. 索引檔（`00-index.md`）的「現況」欄是**別的文件的摘要**——
    #      同一個缺口在那份文件自己的未解表裡還有一列。
    #   2. 圖例列與「這一份是缺口總表」這種**只是提到「未解」兩個字**的列。
    idx = [r for r in rows if r["file"].endswith("00-index.md")]
    by_dir = {}
    for r in rows:
        top = r["file"].split("/")[1] if r["file"].startswith("docs/") else r["file"].split("/")[0]
        by_dir[top] = by_dir.get(top, 0) + 1
    w("⚠ **這是列數，不是獨立問題數。** 索引檔的「現況」欄是別的文件的摘要，"
      f"同一個缺口在那份文件自己的未解表裡還有一列——這類共 **{len(idx)}** 列"
      "（另有少數只是提到「未解」兩個字的圖例列）。\n\n")
    w("| 來源目錄 | 列數 |\n|---|---:|\n")
    for d in sorted(by_dir, key=lambda k: -by_dir[k]):
        w(f"| `docs/{d}/` | {by_dir[d]} |\n")
    w("\n")

    for dom in order:
        sel = [r for r in rows if r["domain"] == dom]
        if not sel:
            continue
        w(f"## 2.{order.index(dom)+1} {dom}（{len(sel)} 條）\n\n")
        w("| 出處 | 缺口 | 現況 | 裁決 |\n|---|---|---|---|\n")
        for r in sorted(sel, key=lambda r: (r["file"], r["section"])):
            short = r["file"][len("docs/"):]
            link = "../" + short
            # ⚠ 顯示名去掉 `docs/` 是為了短，但**去掉之後可能撞到真的目錄**：
            # `docs/mobile/` 砍成 `mobile/…`，而 repo 根目錄真的有一個
            # `mobile/`（gomobile 綁定）。幽靈掃描會把它當成指不到的路徑，
            # 而那不是誤報——同一個字串在 repo 裡確實有兩個意思。
            # 撞到的時候就留完整路徑。
            label = short if short.split("/")[0] not in TOP_LEVEL_DIRS else r["file"]
            v = r["verdict"]
            item = MDLINK.sub(r"\1", r["item"])
            st = MDLINK.sub(r"\1", r["status"]).replace("\n", " ")
            if len(st) > 160:
                st = st[:157] + "…"
            w(f"| [`{label}`]({link}) | {item} | {st} | {v} |\n")
        w("\n")

    w("## 3. 這支工具的盲區\n\n")
    w("**目前是 0。** 每一份提到「未解／未定案／未定位／缺口／未驗」的文件，\n")
    w("要嘛抽得出至少一條，要嘛在檔尾寫了 `<!-- 缺口：無 -->` 明講自己沒有。\n")
    w("這一條由 `--strict` 把關，`check.sh` 帶著它跑——"
      "**新文件寫了「未解」卻沒有未解小節，提交會被擋下來**。\n\n")
    w("抽取只認四種結構（專門的未解小節、表格最後一欄標未解的列、\n")
    w("`**未解**：…`、收尾是「…未解」的句子）。**寫在段落中段、\n")
    w("或用別的詞說「這個還不知道」的缺口抽不到**——下列檔案提到未解\n")
    w("卻一列都沒抽出來，要嘛缺口寫成別的句式，要嘛那些字樣只是在講別的事：\n\n")
    if silent:
        for rel in silent:
            w(f"- `{rel}`\n")
    else:
        w("（沒有）\n")
    if silent and "--strict" in sys.argv:
        sys.stderr.write(
            "盲區 %d 份：這些文件提到未解卻抽不出任何一條。\n"
            "要嘛補一個「未解」小節，要嘛在檔尾加 <!-- 缺口：無 -->。\n%s\n"
            % (len(silent), "\n".join("  " + r for r in silent)))
        sys.exit(1)
    w("\n只印抽得到的部分，會讓解析失敗長得像「那份文件沒有缺口」。\n")
    w("這一節就是為了讓那個差別看得見。\n")


if __name__ == "__main__":
    main()
