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
# ⭐ **目錄索引不是缺口的家。** `<目錄>/00-index.md` 的「現況」欄是**別份文件的摘要**，
# 缺口的正本在那一份自己的未解表裡；把索引也收進來，同一個缺口會被數兩次。
# 反過來，索引也不該被當成盲區——它沒有自己的缺口是正確狀態，
# 逼它寫一個「未解」小節只會產生假資料。
#
# ⚠ 代價說清楚：**索引的摘要過期時，這一份不會再提醒**（2026-09-02 就有一例，
# `spec/00-index` 對 `42` 的摘要停在「結果階段的上框未解」，而 `spec/42` §2 早就解了）。
# 那是 `tools/stale_scan.py` 與 `tools/index.py` 的守備範圍，不是缺口總表的。
INDEX_NAME = "00-index.md"

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
# ⭐ 「這一列自己承認已經解了」——只認**開頭**。
# 「已證實 A，但 B 仍未完」是合法的缺口列，整列比對會把它一起吃掉；
# 反過來，`| 段 3 0x21A0 那張空槽圖 | **已驗**：… |` 這種列是答案不是缺口，
# 而 SOLVED 只認「已解／已定案」，漏掉「已驗／已證實／已讀」這幾種寫法。
SOLVED_LEAD = re.compile(
    r"^(?:\*\*)?(?:✅|已解|已驗|已讀|已證實|已確認|已定案|已實作|已接入)")
STRUCK = re.compile(r"^~~")
SEP = re.compile(r"^\|[\s:|-]+\|$")
# 文件明講自己沒有缺口。有這一行就不算盲區——
# 「解析不到」與「真的沒有」必須分得開，否則盲區清單永遠清不完。
NO_GAPS = re.compile(r"<!--\s*缺口：無\s*-->")
# ⭐ **未解小節裡的三種東西不是缺口，而且用結構就分得出來。**
# 少了這三條，536 列裡有 14 列是工具自己造出來的——最刺眼的是
# `<!-- 缺口：無 -->`：那是文件**明講自己沒有缺口**的標記，
# 卻被當成一條缺口收進總表，**判斷完全相反**。
#
#   1. HTML 註解（`<!-- 缺口：無 -->`／`<!-- 缺口：見上表 -->`）——標記不是內容。
#   2. 小節內文只寫「（無。）」——那是「這裡沒有」的另一種寫法。
#   3. 收尾是「：」的引言句——它引出的條列會各自被收，引言本身不是缺口。
#
# ⚠ 三條都只看**結構**（註解語法、整行等於「無」、收尾標點），不猜語意。
# 想擋「這一列其實是答案」得靠 `SOLVED_LEAD` 那一組，不要往這裡加關鍵字。
EMPTY_BODY = re.compile(r"^[（(]?\s*(無|沒有|N/?A|—|-{1,3})\s*[。.，,]?\s*[）)]?$")
# ⭐ **平台層與工具鏈產物不算 remake 的缺口**
# （DOS／BIOS 介面：使用者裁定 2026-08-23；編譯器 runtime／連結進來的驅動：2026-09-02）。
# 那些是原版與 DOS 之間的介面——`INT 61h` 的音效 TSR 服務號、
# BIOS 的顯示卡暫存器、磁碟服務。remake 跑在 Go／Ebiten 上，
# **不跟 DOS TSR 講話**，知道服務號 `ah=4` 是什麼也不會改變任何一行 Go。
#
# ⚠ **用標記不用關鍵字。** 關鍵字掃過一次，`port` 中了
# `spec/64-capital-relocation-report` 的檔名——語意判斷交給寫的人，
# 工具只認明確標記。
#
# ⚠ **不是刪掉，是分開數。** 標記過的列仍然列在報告的獨立小節裡；
# 讓缺口從畫面上消失比多數幾列更糟。
PLATFORM = re.compile(r"\[DOS/BIOS\]")
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


def is_superseded(path):
    """狀態行自己說「已被…取代」的文件。

    ⚠ 判準刻意窄：只認**狀態行**（前 8 行）裡同時出現「已被」與「取代」。
    正文提到別的東西被取代不算——那是敘述，不是這一份的狀態。
    """
    try:
        with open(path, encoding="utf-8", errors="replace") as fh:
            head = "".join(fh.readline() for _ in range(8))
    except OSError:
        return False
    for line in head.splitlines():
        if "狀態" in line and "已被" in line and "取代" in line:
            return True
    return False


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
    in_fence = False      # 在 ``` 區塊裡
    solved_para = False   # 正處在一段「已解…」的敘述裡
    for i, line in enumerate(lines):
        # ⭐ **反組譯區塊裡的註解不是缺口。** 那些行長這樣：
        #     call sub_1E3C0                    ; ← 未解
        # 它是「這一行在做什麼還沒讀」的旁註，不是一條獨立登記的缺口；
        # 收進來會讓同一個未解被數兩次（旁註一次、未解表一次），
        # 而且沒有出處欄可以回查。2026-08-21 稽核量到 2 筆。
        if line.lstrip().startswith("```"):
            in_fence = not in_fence
            continue
        if in_fence:
            continue
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
            # ⭐ 「已解的兩條：…」那種段落，**整段**都不是缺口。
            # 只擋第一行沒有用：第二行不含「已解」，會被當成該小節的代表句
            # 抽進總表（docs/re/62 的「那一區逐像素 PASS」就是這樣進來的）。
            if not body:
                solved_para = False
            elif SOLVED.search(body):
                solved_para = True
            # `![說明](../images/x.png)` 去掉連結之後剩「!說明」，
            # 看起來像一句話。圖片行不會是缺口。
            if len(body) >= 8 and not solved_para \
                    and not body.startswith((">", "#", "```", "![", "<!--")) \
                    and not EMPTY_BODY.match(body) and not body.endswith((":", "：")) \
                    and not SOLVED.search(body) and not META.search(body) \
                    and (bullet or section not in lead_done):
                lead = SENT.search(body).group(0).strip()
                # 「（無。原先掛在這裡的兩條都收掉了：…」——第一句就說沒有，
                # 後面是為什麼沒有。**整節都不是缺口**，不能只跳過這一行，
                # 否則下一行的殘句會遞補成該小節的代表。
                if EMPTY_BODY.match(lead):
                    muted = True
                    continue
                lead_done.add(section)
                items.append((lead[:120], "（未解小節內文）", section))
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
        # ⚠ 這裡**不能**用整列比對「已解」。缺口列常常把已解的部分寫進說明
        # （「`ENDPAL` 那邊已解，開場這邊還沒做」），整列比對會把真缺口一起吃掉，
        # 而那份文件就變成盲區——看起來像「沒有缺口」。只認**開頭**。
        if "✅" in line or STRUCK.match(c[0]) \
                or any(SOLVED_LEAD.match(x) for x in c):
            continue
        # ⭐ `muted` 以前只擋散文，不擋表格列——於是「### 3.1 已解的」底下
        # 整張表照樣被收進缺口總表（`mechanics/80` 四列都是答案）。
        # 子標題說「已解」就是說「以下不是缺口」，表格也算。
        if (open_level is not None and not muted) or gap_table:
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

    rows, silent, superseded = [], [], []
    for path in files:
        rel = os.path.relpath(path, repo).replace(os.sep, "/")
        if rel in SKIP or os.path.basename(rel) == INDEX_NAME:
            continue
        # ⭐ **已被取代的文件，它的「未解」是當時的未解。**
        # 每一批發行紀錄都會再列一次「Windows／macOS 實機」「沒有音效裝置」，
        # 於是 `docs/release/` 的 33 列其實只有 12 個獨立缺口——同一個缺口
        # 被數了六次（四種寫法 ＋ 兩個變體）。**還開著的缺口一定也在最新那一份**，
        # 所以這裡只跳過舊的，不會漏。
        if is_superseded(path):
            superseded.append(rel)
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

    # ⭐ DOS／BIOS 平台層分流出去（使用者裁定 2026-08-23）。
    # 不是刪掉——另外列一節，讓「不算缺口」這個決定看得見。
    platform = [r for r in rows if PLATFORM.search(r["item"] + " " + r["status"])]
    rows = [r for r in rows if r not in platform]

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

    docs_with = len({r["file"] for r in rows})
    if docs_with:
        w("## 0. ⚠ 這個數字在量什麼\n\n")
        w("**%d 列分布在 %d 份文件，平均每份 %.1f 列。**\n\n"
          % (len(rows), docs_with, len(rows) / docs_with))
        w("⭐ **所以它比較接近「文件有多少份」，不是「原版還有多少沒解」。**\n")
        w("每寫一份新文件就帶進約三列自己的未解——而 `check.sh --strict` 還會\n")
        w("**要求**每份文件要嘛有未解小節、要嘛明講 `<!-- 缺口：無 -->`。\n")
        w("於是「解出新東西 → 寫一份文件 → 總數上升」是這個指標的常態，\n")
        w("不是退步。\n\n")
        w("⚠ 反過來也一樣：**數字變小不自動等於進度**。2026-08-21 的稽核\n")
        w("把它從 570 降到 431，而那 −139 沒有一列是靠解出新東西減掉的。\n\n")
        w("**要看進度請看別的東西**：`docs/spec/` 的 CONFORMED 份數、\n")
        w("`docs/playtest/` 的逐像素數字、`docs/re/21` 的覆蓋地圖。\n")
        w("這一份回答的是「還有什麼沒解」，**不是「還剩多少」**。\n\n")

    if platform:
        w("> ⭐ 另有 **%d 列標成 `[DOS/BIOS]`**，"
          "**不計入下面的總數**——那是原版與 DOS／BIOS 之間的介面\n"
          "> （`INT` 服務號、顯示卡暫存器、磁碟服務），"
          "而 remake 跑在 Go／Ebiten 上不跟它們講話。\n"
          "> 清單在 §9。\n\n" % len(platform))

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
    if superseded:
        w(f"⭐ **狀態行自稱「已被…取代」的 {len(superseded)} 份文件不計。**"
          "它們的「未解」是**當時**的未解，而每一批發行紀錄都會再列一次"
          "「Windows／macOS 實機」「沒有音效裝置」——一度讓 `docs/release/` 的 "
          "33 列其實只有 12 個獨立缺口。**還開著的缺口一定也在最新那一份**，"
          "所以跳過舊的不會漏。\n\n")
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

    w("\n## 9. 平台層與工具鏈產物（不計入總數）\n\n")
    w("兩類東西掛 `[DOS/BIOS]`，都**不算 remake 的缺口**：\n\n")
    w("1. **原版與 DOS／BIOS 之間的介面**：`INT` 服務號、顯示卡暫存器、磁碟服務。\n")
    w("   知道 `INT 61h` 的 `ah=4` 是什麼，不會改變任何一行 Go。\n")
    w("2. **編譯器 runtime、連結進來的驅動與程式庫**：C runtime 啟動碼、算術／字串\n")
    w("   輔助常式、被靜態連結進執行檔的驅動模組（本作的 segment `0x2000` 整段是\n")
    w("   `INT 33h` 滑鼠包裝，見 [`24`](24-unread-function-catalogue.md) §3）。\n")
    w("   **這些是 toolchain 產物，不是玩家看得到的遊戲邏輯。**\n\n")
    w("使用者裁定：第 1 類 2026-08-23、第 2 類 2026-09-02。\n\n")
    w("⚠ **排除的是「怎麼跟平台講話」，不是「原版選了什麼參數」。**\n")
    w("滑鼠驅動把游標範圍設成 640×400 是**遊戲行為**，仍然要讀；\n")
    w("`INT 33h` 本身的呼叫慣例不用。判準與流程見\n")
    w("`~/.claude/knowledge-base/retro/compiler-runtime-helper-fingerprints.md`。\n\n")
    w("⚠ **分開數不是不數。** 這些仍然是原版的未解之處，只是**不擋 remake**；\n")
    w("哪天要寫「原版怎麼跟 DOS 講話」的文件，這一節就是清單。\n\n")
    if platform:
        w("| 出處 | 缺口 | 現況 |\n|---|---|---|\n")
        for r in platform:
            w("| [`%s`](../%s) | %s | %s |\n"
              % (r["file"].replace("docs/", "", 1), r["file"].replace("docs/", "", 1),
                 r["item"], r["status"]))
    else:
        w("（沒有）\n")


if __name__ == "__main__":
    main()
