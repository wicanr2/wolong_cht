#!/usr/bin/env python3
"""從 docs/ 生出索引，並檢查「已經解掉的事情有沒有還被寫成未解」。

    tools/index.py check      只檢查，有問題回非 0（CI／提交前跑）
    tools/index.py generate   重生 docs/INDEX.md

## 為什麼要這個

**手寫的索引本身就是會過時的東西。** 這個專案實際踩過三次：

- `docs/playtest/04` 列了「四件事被滑鼠自動化擋住」，
  而那四件在寫下的當下就有三件已經解掉了 —— 之後好幾輪都繞著
  那道不存在的閘打轉。
- `docs/formats/07` 的狀態行寫「像素格式未解」，但同一份文件的
  §8–§10 就是像素格式的解法。
- `CONTEXT.md` 的「進行中／受阻」表列著早就完成的項目。

共同點是**沒有人會為了寫一行新結論而回頭改三份舊文件的摘要**。
所以這裡不生產「另一份要維護的文件」，而是：

1. **狀態從內文推導**，不從摘要抄；
2. **衝突自動報錯**——同一個主題在兩份文件裡等級不同就當成錯誤。

## 檢查項

| 檢查 | 為什麼 |
|---|---|
| 每份 docs 文件都要有狀態行與日期 | 沒有狀態行的文件無法判斷新舊 |
| 狀態行說「未解／受阻」但內文有對應的 confirmed／READY | 就是上面那三次 |
| A 的狀態行說某識別字未解，B 的狀態行說它 READY | `formats/05` 說 `MMAP.MAP` 未解，`formats/06` 就是它的解法 |
| 斷言表的同一個鍵在不同文件有不同等級 | 兩份文件對同一件事說法不同 |
| 文件內的相對連結指得到檔案 | 改檔名時最容易漏 |
| `[…](x.md) §N` 的小節真的存在 | 七筆引用把**行號**寫成小節號（`§590`、`§1065`…），而它們都指得到檔案，所以連結檢查一路放行 |
| `docs/spec/00-index.md` 的狀態欄與規格自己的狀態行一致 | 索引那一欄是散文，一份文件內部的檢查看不到跨文件矛盾 |
| `CONTEXT.md` 提到的 docs 路徑都存在，且有指向 `docs/INDEX.md` | 完整清單交給生成的那份，人只維護指標 |

只用標準函式庫。
"""
import os
import pathlib
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOCS = os.path.join(ROOT, "docs")
CONTEXT = os.path.join(ROOT, "CONTEXT.md")
OUT = os.path.join(DOCS, "INDEX.md")

# 推論等級，由強到弱。`docs/mechanics/00-index.md` 定義的那一組。
LEVELS = ["confirmed", "READY", "強證據", "說明書", "假說", "未解"]
LEVEL_RANK = {v: i for i, v in enumerate(LEVELS)}

STATUS_RE = re.compile(r"\*\*狀態[：:](.+?)\*\*", re.S)
DATE_RE = re.compile(r"^-\s*日期[：:]\s*(\d{4}-\d{2}-\d{2})", re.M)
TITLE_RE = re.compile(r"^#\s+(.+)$", re.M)
LINK_RE = re.compile(r"\]\((?!https?:)([^)#]+)")
# 斷言表的一列：| `+0x08` | u16 | 說明 | confirmed |
ROW_RE = re.compile(r"^\|(.+)\|\s*$", re.M)
HEAD_RE = re.compile(r"^#{2,4}\s+(.+)$", re.M)

# 狀態行裡代表「還沒好」的字眼。
OPEN_WORDS = ["未解", "受阻", "⛔", "沒解", "還沒"]


def md_files():
    for dirpath, _, names in os.walk(DOCS):
        for n in sorted(names):
            # TEMPLATE.md 是給人複製的骨架，日期／狀態欄本來就留空，
            # 拿文件的規則檢查它只會得到永遠修不掉的一條問題。
            if n.endswith(".md") and n not in ("INDEX.md", "TEMPLATE.md"):
                yield os.path.join(dirpath, n)


def root_md_files():
    """repo 根目錄的 markdown（WORKLIST／README／計畫書那一批）。

    ⭐ 這些檔案**不進 `docs/` 的索引規則**（沒有狀態行、沒有日期欄），
    但「密碼頁擋住 oracle」那條斷言檢查一樣要涵蓋它們——
    WORKLIST.md 就同時留著 2026-08-12 的勘誤與三處未改的舊寫法，
    因為檢查只走 `docs/`。**會復發的斷言，檢查範圍要跟著它跑。**
    """
    for n in sorted(os.listdir(ROOT)):
        if n.endswith(".md") and os.path.isfile(os.path.join(ROOT, n)):
            yield os.path.join(ROOT, n)


def rel(path):
    return os.path.relpath(path, ROOT).replace(os.sep, "/")


def cell_level(cell):
    """一格文字裡的推論等級（取最強的那個）。沒有就回 None。"""
    best = None
    for lv in LEVELS:
        if lv in cell:
            if best is None or LEVEL_RANK[lv] < LEVEL_RANK[best]:
                best = lv
    return best


def strip_md(s):
    return s.replace("*", "").replace("`", "").strip()


class Doc:
    def __init__(self, path):
        self.path = path
        self.text = open(path, encoding="utf-8").read()
        m = TITLE_RE.search(self.text)
        self.raw_title = m.group(1) if m else os.path.basename(path)
        self.title = strip_md(self.raw_title)
        m = STATUS_RE.search(self.text)
        # ⚠ status 是給人看的（去掉 markdown），**識別字比對要用 raw_status**——
        # strip_md 會把反引號拿掉，用它去找 `MMAP.MAP` 永遠找不到。
        self.raw_status = " ".join(m.group(1).split()) if m else ""
        self.status = " ".join(strip_md(m.group(1)).split()) if m else ""
        m = DATE_RE.search(self.text)
        self.date = m.group(1) if m else ""
        self.claims = self._claims()

    def _claims(self):
        """抓出「鍵 → 等級」的斷言。只認最後一格是等級的表格列。

        鍵要**冠上最近的一個標題**：`+0x01` 在勢力表、據點表、軍團表
        是三個不同的東西，不冠標題會把它們當成同一條而互報衝突。
        """
        out = {}
        # 先記下每個標題的位置，等一下用位移找「最近的上一個標題」。
        heads = [(m.start(), strip_md(m.group(1))) for m in HEAD_RE.finditer(self.text)]

        def head_of(pos):
            name = ""
            for start, h in heads:
                if start < pos:
                    name = h
                else:
                    break
            return name

        for m in ROW_RE.finditer(self.text):
            line = m.group(1)
            cells = [c.strip() for c in line.split("|")]
            if len(cells) < 3:
                continue
            lv = cell_level(cells[-1])
            # 最後一格必須**幾乎只有**等級，否則整段散文都會被當成斷言。
            if lv is None or len(strip_md(cells[-1])) > len(lv) + 24:
                continue
            raw = strip_md(cells[0])
            if not raw or raw.startswith("---") or raw in ("偏移", "欄位", "項目"):
                continue
            h = head_of(m.start())
            key = f"{h} ▸ {raw}" if h else raw
            prev = out.get(key)
            if prev is None or LEVEL_RANK[lv] < LEVEL_RANK[prev]:
                out[key] = lv
        return out

    def open_status(self):
        return any(w in self.status for w in OPEN_WORDS)



# 敘述「當初怎麼錯的」的用語。刻意窄——寧可漏抓也不要誤報，
# 誤報會讓人把整個檢查關掉。
NARRATIVE = re.compile(r"當初|先前寫|先前記|原本寫|舊版寫|後來發現|這一節的結論|錯掉的那")

# 允許的地方：推翻紀錄本來就該集中在這裡。
#   CONTEXT.md 的「已被推翻的斷言」整節
#   CLAUDE.md 的「教訓」整節（那是規則不是敘事，但會引用具體案例）
ALLOW_SECTIONS = ("已被推翻的斷言", "從前三個 remake 專案帶過來的教訓")


def narrative_hits(path):
    """回傳 (行號, 內容)。允許區段內的不算。

    ⚠ `docs/playtest/` **整個目錄豁免**：那是有日期的實驗紀錄，
    文類本身就是「當時跑了什麼、結果如何、哪條路走不通」。
    對它套「只寫現況」會把紀錄的用途毀掉。
    規範的對象是**宣稱現況的文件**：規格、反組譯筆記、機制文件。

    誤報會讓人把整個檢查關掉，所以豁免要寫在這裡而不是靠人記得。
    """
    path = pathlib.Path(path)
    if "playtest" in path.parts:
        return []
    out, allowed = [], False
    for i, line in enumerate(path.read_text(encoding="utf-8").splitlines(), 1):
        if line.startswith("#"):
            allowed = any(a in line for a in ALLOW_SECTIONS)
        if not allowed and NARRATIVE.search(line):
            out.append((i, line.strip()))
    return out

# 「密碼頁擋住 oracle」的各種寫法。只認**斷言句**，不認事實描述——
# 「有防拷」「防拷畫面」本身是對的（`docs/playtest/01` §2），
# 錯的是「因此 oracle 過不去」。
BLOCKED_BY_COPY_PROTECTION = re.compile(
    r"(防拷|密碼頁|密碼畫面)[^。\n]{0,20}(擋住|擋著|阻擋|停擺|過不去|用不了)"
    r"|(被|受)(防拷|密碼頁)[^。\n]{0,6}(擋|阻)"
    # ⭐ 換個說法一樣是同一個斷言：「密碼保護使 parity 不可證實」
    # 「密碼保護造成的對拍邊界」。WORKLIST.md 三處就是這樣躲過第一版的。
    r"|(防拷|密碼保護|密碼頁)[^。\n]{0,24}"
    r"(不可證實|無法證實|不能證實|邊界|限制|受限)"
    # ⭐ 第三種寫法：把它擺成**修飾語**而不是主張。
    # 「不宣稱**密碼保護下**的原版自然長程逐像素對拍」——句子本身在講
    # 「我沒有做這件事」，但那個修飾語把「沒做」的原因暗示成防拷。
    # 前兩條都比不到它，因為斷言的動詞不在句子裡。
    r"|(防拷|密碼保護|密碼頁)(下|之下)的")
# 否定寫法（「不再是阻擋」「即可越過」）本身就是正確答案，不能誤報。
NOT_BLOCKED = re.compile(r"不再|不是阻擋|不阻擋|不擋|已可|即可|可越過|越過|不會阻擋")


def quoted(line, span):
    """命中的字串是不是包在「」裡（＝**指稱**那個斷言，不是主張它）。

    ⭐ 沒有這一條，檢查就描述不了自己：一句
    「`tools/index.py` 的『密碼頁擋住 oracle』檢查只走 docs/」
    會被自己擋下來。**規則要能寫下自己的名字。**
    """
    lo, hi = span
    return line.rfind("\u300c", 0, lo) > line.rfind("\u300d", 0, lo) \
        and line.find("\u300d", hi) != -1

# ⭐ 這兩份是**correction 本身的家**：`CLAUDE.md` §4.0 與 `CONTEXT.md` 的
# 「已被推翻的斷言」表就是用來寫下「密碼頁不擋 oracle」的，
# 裡面必然引用舊斷言的字樣。對它們套這條檢查只會得到永遠修不掉的誤報。
ASSERTION_CHECK_EXEMPT = {"CLAUDE.md", "CONTEXT.md"}


def check(docs):
    problems = []

    # ① 每份文件都要有狀態行與日期。
    for d in docs:
        if not d.status:
            problems.append(f"{rel(d.path)}：沒有「**狀態：…**」那一行")
        if not d.date:
            problems.append(f"{rel(d.path)}：沒有「- 日期：YYYY-MM-DD」")

    # ② 狀態行說未解，但內文自己有 confirmed／READY 的斷言。
    #    這就是 formats/07 那個「像素格式未解」的形狀。
    for d in docs:
        if not d.open_status():
            continue
        solid = [k for k, v in d.claims.items() if v in ("confirmed", "READY")]
        if len(solid) >= 3:
            problems.append(
                f"{rel(d.path)}：狀態行寫「{d.status[:28]}…」，"
                f"但內文有 {len(solid)} 條 confirmed／READY 斷言"
                f"（如 {'、'.join(solid[:3])}）——狀態行可能過時")

    # ③ 正文在敘述「當初怎麼錯的」。
    #
    # 規則（`~/.claude/rulebook/63`）：斷言被推翻就把正文改寫成正確答案，
    # 推翻紀錄集中到 CONTEXT.md 的「已被推翻的斷言」表，正文最多留一個指標。
    #
    # ⚠ **為什麼要做成檢查而不是再寫一遍規則**：63 的觸發條件全是稽核時
    # （review／接續 worklist／斷言完成前），而違規發生在**寫入時**——
    # 剛解出東西正要寫進文件那一刻，沒有任何觸發條件成立。
    # 這個專案裡黏住的紀律都有測試；只靠記憶的規則在「剛解出東西」
    # 那一刻最不可靠。第一次量到 18 行命中、9 個檔，那是預設寫法不是偶發。
    for d in docs:
        for n, line in narrative_hits(d.path):
            problems.append(
                f"{rel(d.path)}:{n}：正文在敘述當初怎麼錯的"
                f"（{line[:34]}…）——改寫成正確答案，"
                f"推翻紀錄放 CONTEXT.md 的「已被推翻的斷言」")

    # ③.5 已經解決的阻擋，不准再被寫成阻擋。
    #
    # 現況：**松崗 DOS/V 的密碼頁不擋 oracle**——四格留白直接按「確定」
    # 就進開場（`docs/playtest/18`，2026-08-12 三組受控實驗）。
    #
    # ⚠ 為什麼要做成檢查：這條在 `docs/playtest/01`／`13`／`17`／`18` 與
    # `CONTEXT.md` §0.1 都寫過了，**還是在四份文件裡復發**，包括
    # 同一個 session 剛寫的新筆記。原因是它長得像「已知的專案背景」，
    # 寫的時候不會想到要查——**這種斷言只有機器擋得住**。
    for path in [d.path for d in docs] + list(root_md_files()):
        if os.path.basename(path) in ASSERTION_CHECK_EXEMPT:
            continue
        try:
            text = open(path, encoding="utf-8").read()
        except OSError:
            continue
        for n, line in enumerate(text.split("\n"), 1):
            m = BLOCKED_BY_COPY_PROTECTION.search(line)
            if m and not NOT_BLOCKED.search(line) \
                    and not quoted(line, m.span()):
                problems.append(
                    f"{rel(path)}:{n}：又把密碼頁寫成 oracle 的阻擋"
                    f"（{line.strip()[:34]}…）——空白確認就會過，"
                    f"見 docs/playtest/18")

    # ④ 同一個鍵在不同文件等級不同。

    #
    # ⚠ 這一項只**提醒**不擋，因為偏移鍵天生有歧義（同一個 `+0x08`
    # 在不同的表是不同欄位），冠了標題也還是可能誤報。
    # 它的用途是「這個東西別處是不是已經解了」的線索，不是門禁。
    where = {}
    for d in docs:
        for k, v in d.claims.items():
            where.setdefault(k, []).append((v, rel(d.path)))
    warns = []
    for k, lst in sorted(where.items()):
        if len({v for v, _ in lst}) > 1:
            desc = "、".join(f"{v}（{p}）" for v, p in lst)
            warns.append(f"斷言「{k}」在不同文件等級不一致：{desc}")

    # ⑥ 跨文件：A 的狀態行說某個東西「未解」，B 的狀態行卻說它 READY。
    #
    # 抓的是狀態行裡用反引號括起來的識別字（`MMAP.MAP`、`TALK.DAT`…）。
    # 實際踩過：`docs/formats/05` 寫「`MMAP.MAP` 的編碼未解」，
    # 而隔壁的 `docs/formats/06` 整份就是它的解法且標 READY。
    ident = re.compile(r"`([A-Za-z0-9_.*/]{3,})`")
    closed = {}
    for d in docs:
        if d.status and not d.open_status():
            # **標題也要看。** `docs/formats/06` 的狀態行只有「READY。」，
            # 識別字在標題「06 — `MMAP.MAP` 的 RLE 壓縮」裡——
            # 只看狀態行的話這條檢查永遠不會觸發（做正對照才發現）。
            for t in ident.findall(d.raw_status + " " + d.raw_title):
                closed.setdefault(t, []).append(rel(d.path))
    for d in docs:
        if not d.open_status():
            continue
        for t in ident.findall(d.raw_status):
            # 只有當 t 出現在「未解」那半句附近才算數，這裡從寬：
            # 同一狀態行同時出現 READY 與未解是常態（多主題文件）。
            if t in closed and re.search(rf"{re.escape(t)}`?[^，。；]*(未解|受阻|沒解|還沒)", d.raw_status):
                problems.append(
                    f"{rel(d.path)}：狀態行說 `{t}` 未解，"
                    f"但 {'、'.join(closed[t])} 說它已經好了")

    # ⑦ `docs/spec/00-index.md` 的狀態欄不能與規格自己的狀態行相反。
    #
    # 那一欄是散文，所以 ② 的「狀態行與內文矛盾」完全看不到它——
    # ② 比的是同一份文件內部，這裡矛盾的是兩份文件。
    # 實際踩過（2026-08-21 稽核，`CONTEXT.md` §6.1）：索引把 `29-audio` 記成
    # 「**DRAFT**：播放層未做」而那份規格是 CONFORMED、音效整條接通；
    # `51-vga-dac` 記成「**READY**，尚未全面套用」而換算早就在 `toSRGB` 裡。
    #
    # ⚠ **只擋得到狀態碼那一種。** 散文形式的但書（「捲軸未解」「頭像與滑鼠未接」）
    # 沒有可比的 token，那八列還是要靠人逐列對 code。
    # **擋一半不是白做**——狀態碼相反是最誤導的一種，它會讓人以為整份規格不能用。
    spec_dir = os.path.join(DOCS, "spec")
    spec_index = os.path.join(spec_dir, "00-index.md")
    if os.path.exists(spec_index):
        own = {}
        for d in docs:
            if os.path.dirname(d.path) != spec_dir:
                continue
            m = re.match(r"([A-Z]{4,})", d.raw_status.replace("狀態：", ""))
            if m:
                own[os.path.basename(d.path)] = m.group(1)
        row = re.compile(r"\[`([0-9]+[^`]*\.md)`\]\([^)]*\)\s*\|([^|]*)\|")
        for name, cell in row.findall(open(spec_index, encoding="utf-8").read()):
            said = set(re.findall(r"(DRAFT|READY|CONFORMED)", cell))
            if said and name in own and own[name] not in said:
                problems.append(
                    f"docs/spec/00-index.md：那一列把 `{name}` 記成 "
                    f"{'／'.join(sorted(said))}，但它自己的狀態行是 {own[name]}")

    # ⑧ `[…](檔案.md) §N` 的 `§N` 要真的存在於那份文件。
    #
    # ④ 只驗檔案存不存在，驗不到小節。實際踩過（2026-08-21 稽核）：
    # 七筆引用寫的是**行號**不是小節號（`§1065`、`§590`、`§467`、`§150`、`§91`…），
    # 而它們全部指得到檔案，所以 ④ 一路放行。⭐ **一份長文件裡「§590」
    # 看起來完全合理**——沒有機械檢查就只能靠讀的人自己去翻。
    #
    # 只驗 `[…](x.md) §N` 這一種形式；`§2.1` 允許只有子節存在。
    sec_cite = re.compile(r"\]\(([^)\s]+\.md)\)\s*(?:的\s*)?§([0-9]+(?:\.[0-9]+)*)")
    heads = {}
    for d in docs:
        heads[os.path.normpath(d.path)] = {
            m.group(1) for m in re.finditer(r"^#{2,6}\s+([0-9]+(?:\.[0-9]+)*)", d.text, re.M)}
    for d in docs:
        if os.path.basename(d.path) == "43-open-questions.md":
            continue  # 生成的彙總檔，內容是別份文件的原文（引用連同上下文一起搬過來）
        for target, num in sec_cite.findall(d.text):
            full = os.path.normpath(os.path.join(os.path.dirname(d.path), target))
            own = heads.get(full)
            if own is None or num in own:
                continue
            if any(h.startswith(num + ".") for h in own):
                continue  # 父節沒編號但子節有
            problems.append(
                f"{rel(d.path)}：引用 {target} §{num}，但那份文件沒有這個小節")

    # ④ 相對連結指得到檔案。
    for d in docs:
        for target in LINK_RE.findall(d.text):
            t = target.strip()
            if not t or t.startswith("<"):
                continue
            full = os.path.normpath(os.path.join(os.path.dirname(d.path), t))
            if not os.path.exists(full):
                problems.append(f"{rel(d.path)}：連結指不到 {t}")

    # ⑤ CONTEXT.md 提到的 docs 路徑不能是死的。
    #
    # ⚠ **不檢查「有沒有列全」**。第一版要求 CONTEXT.md 的清單與 docs/
    # 完全一致，那等於強制人工複製一份清單——而那正是會爛掉的東西
    # （這個工具存在的理由）。完整清單由 docs/INDEX.md 生成，
    # CONTEXT.md 只需要指過去。
    ctx = open(CONTEXT, encoding="utf-8").read()
    # INDEX.md 不在 docs 清單裡（它是產物），但 CONTEXT 本來就該指它。
    actual = {rel(d.path) for d in docs} | {rel(OUT)}
    for extra in sorted(set(re.findall(r"`(docs/[^`]+\.md)`", ctx)) - actual):
        problems.append(f"CONTEXT.md 提到不存在的 {extra}")
    if "docs/INDEX.md" not in ctx:
        problems.append("CONTEXT.md 沒有指向 docs/INDEX.md（完整清單在那裡）")

    return problems, warns


def generate(docs):
    lines = [
        "# 文件索引（自動產生，不要手改）",
        "",
        "由 `tools/index.py generate` 從 `docs/**/*.md` 生出來。",
        "**狀態與日期是從各文件的內文讀的，不是另外維護的一份摘要**——",
        "手寫的索引會過時，而過時的索引比沒有索引更糟"
        "（`docs/playtest/04` 的「四件事被擋住」整張表在寫下的當下就是錯的）。",
        "",
        "提交前跑 `tools/index.py check`，它會擋下：狀態行與內文矛盾、",
        "同一斷言在兩份文件等級不同、連結壞掉、`CONTEXT.md` 漏登記。",
        "",
        "## 文件",
        "",
        "| 文件 | 標題 | 狀態 | 日期 |",
        "|---|---|---|---|",
    ]
    for d in sorted(docs, key=lambda d: rel(d.path)):
        st = d.status or "—"
        if len(st) > 60:
            st = st[:58] + "…"
        lines.append(f"| [`{rel(d.path)}`]({os.path.relpath(d.path, DOCS)}) "
                     f"| {d.title} | {st} | {d.date or '—'} |")

    # 斷言總表：一個鍵一列，按等級分組，讓「這件事解了沒」一眼可查。
    where = {}
    for d in docs:
        for k, v in d.claims.items():
            where.setdefault(k, []).append((v, rel(d.path)))
    lines += ["", "## 斷言（欄位／常數 → 推論等級 → 出處）", "",
              f"共 {len(where)} 條。**要查「這件事解了沒」先看這裡**，",
              "不要重讀整份文件，更不要重推一次。", ""]
    for lv in LEVELS:
        keys = sorted(k for k, l in where.items()
                      if min(LEVEL_RANK[x] for x, _ in l) == LEVEL_RANK[lv])
        if not keys:
            continue
        lines += [f"### {lv}（{len(keys)} 條）", "", "| 鍵 | 出處 |", "|---|---|"]
        for k in keys:
            srcs = sorted({p for v, p in where[k] if v == lv})
            lines.append(f"| {k} | {'、'.join(f'`{s}`' for s in srcs)} |")
        lines.append("")
    open(OUT, "w", encoding="utf-8").write("\n".join(lines) + "\n")
    return len(where)


def main():
    cmd = sys.argv[1] if len(sys.argv) > 1 else "check"
    docs = [Doc(p) for p in md_files()]
    if cmd == "generate":
        n = generate(docs)
        print(f"{rel(OUT)}：{len(docs)} 份文件、{n} 條斷言")
        cmd = "check"
    if cmd == "check":
        problems, warns = check(docs)
        for w in warns:
            print(f"  ? {w}")
        if warns:
            print(f"（以上 {len(warns)} 條只是提醒，不擋）")
        if not problems:
            print(f"索引檢查通過（{len(docs)} 份文件）")
            return 0
        print(f"索引檢查發現 {len(problems)} 個問題：")
        for p in problems:
            print(f"  - {p}")
        return 1
    print(__doc__)
    return 2


if __name__ == "__main__":
    sys.exit(main())
