#!/usr/bin/env python3
"""過期斷言掃描：文件指到的東西**存在**，但值已經不對。

`tools/phantom_scan.py` 擋的是「指向不存在的東西」，`tools/index.py` 擋的是
「狀態行與內文矛盾」。這一支擋的是第三類，而它是本專案 2026-08-27 稽核裡
最貴的一類——**格式完全正確、連結都通、只有數字是舊的**：

- `VERIFICATION-MATRIX.md` 把 `wlgame-event3-choice.png` 的 SHA-256 記成
  `C02183…`，實際是 `CA40B8…`。照著驗的人會得出「檔案被動過」的結論。
- `docs/mobile/android-plan.md` 寫模擬器映像 `wolong-android-emulator:20260811`，
  而 `tools/android_smoke.sh` 用的是 `:20260820`。
- `docs/spec/41` 的驗證欄寫 `-open-talk`，實際旗標是 `-open-talk-index`。
- `docs/re/21` 宣稱覆蓋率 T1 100%，重跑是 738 ＋ 1 支 T2。
- `README.md` 把未解列數寫成 456 列／172 份，而 `docs/re/43`（`check.sh`
  每次重生）當時已經是 518 列／199 份。**抄過去的摘要不會自己更新。**
- `README.md` 把規格寫成 84 份／82 CONFORMED／2 READY，而 `docs/spec/`
  已經是 103 份／102／1，被點名的 `34-speed-steps` 也早就 CONFORMED 了。

六層檢查，**每一層都只在能拿到真值時才跑**（缺原版素材、缺 census 就跳過
並明講跳過了）——沉默地少跑一層，比不跑更糟。

用法：

    tools/py.sh tools/stale_scan.py

⚠ **只回報「確定不對」的。** 誤報會讓人把整個檢查關掉（同 `phantom_scan`
的原則）。判不準的一律放行，並在需要時把判準寫窄一點。
"""
from __future__ import annotations

import hashlib
import os
import re
import subprocess
import sys

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# 掃這些檔。`dist-all/` 是產物、`workplace/` 是素材，都不掃。
DOC_ROOTS = ("docs", "packaging", "android", "mobile", "docker")
DOC_FILES = ("README.md", "CLAUDE.md", "CONTEXT.md", "WORKLIST.md",
             "AGENTS.md", "REMAKE-PLAN.md", "RESEARCH-LOG.md")

# ── 第一層：雜湊 ────────────────────────────────────────────────
#
# 只認「檔名 …(80 字內，且中間出現 SHA 或「雜湊」)… 64 個 hex」這種相鄰寫法。
# 放寬到「同一行任何檔名配任何雜湊」會產生大量假陽性——一行裡常有兩三個
# 檔名與兩三個雜湊，配對本身就是猜的。
HASH_CLAIM = re.compile(r"`?\b(?P<hash>[0-9a-fA-F]{64})\b`?")

# 檔名候選：**只認白名單副檔名**。放寬到「任何 `x.y`」會把 `H.264`、
# 版本號、`44.1 kHz` 一起收進來，而它們永遠解不出檔案 → 整批變成「跳過」，
# 這一層就再次靜默失效。
HASH_FILE = re.compile(
    r"[A-Za-z0-9_\-./]+\.(?:png|jpg|gif|mp4|ogg|wav|json|md|go|py|sh|txt|csv|tsv"
    r"|exe|dat|i64|asm|bin|apk|zip|tar|gz|AppImage)\b",
    re.I,
)
# 視窗以**段落**為界（空行／標題），上限 400 字。
# 80 字太小：條列式的「- 影片：…mp4」與「- SHA-256：`…`」中間常隔一整行規格，
# 於是那三支推廣片的雜湊一條都比不到（2026-08-29 突變測試發現）。
HASH_WINDOW = 400

# ── 第二層：docker 映像標籤 ─────────────────────────────────────
IMAGE_TAG = re.compile(r"\b([a-z0-9][a-z0-9._/-]*?):(20[0-9]{6}[a-z0-9.-]*)\b")

# ── 第三層：命令列旗標 ─────────────────────────────────────────
FLAG_CITED = re.compile(r"`-([a-z][a-z0-9-]{2,})`")
FLAG_DEF = re.compile(r'flag\.\w+\("([a-zA-Z0-9-]+)"')
# 不是本專案執行檔的旗標：emulator、ffmpeg、adb、docker、shell 工具。
# ⚠ 這份白名單只會長，不要為了讓檢查過而把真的旗標加進來。
FLAG_ALLOW = {
    "no-audio", "no-window", "no-boot-anim", "no-snapshot", "no-appstream",
    "nostdin", "encoding", "network", "memory", "cpus", "pids-limit",
    "log-opt", "rm", "appimage-extract-and-run", "appimage-extract",
    "maxdepth", "print0", "safe", "count", "run", "vet", "buildvcs",
    "trimpath", "ldflags", "selftest", "strict", "regions", "include-catalogue",
    "version", "help", "gpu", "camera-back", "camera-front", "avd",
}


def doc_paths():
    out = []
    for root in DOC_ROOTS:
        for dirpath, dirnames, filenames in os.walk(os.path.join(REPO, root)):
            dirnames[:] = [d for d in dirnames
                           if d not in (".git", "images", "build", "__pycache__")]
            out += [os.path.join(dirpath, f) for f in filenames if f.endswith(".md")]
    out += [os.path.join(REPO, f) for f in DOC_FILES
            if os.path.isfile(os.path.join(REPO, f))]
    return sorted(out)


def read(path):
    try:
        with open(path, encoding="utf-8", errors="replace") as fh:
            return fh.read()
    except OSError:
        return ""


def resolve(rel, doc):
    """把文件裡寫的路徑解成實際檔案；解不出唯一一個就回 None。"""
    for cand in (os.path.join(REPO, rel), os.path.join(os.path.dirname(doc), rel)):
        if os.path.isfile(cand):
            return cand
    return None


def hash_claims(text, doc, skipped):
    """吐出這份文件裡「某個檔案的雜湊是 X」的宣稱：(行號, 檔案路徑, 宣稱值)。

    ⭐ **自我測試與本體共用這一支。** 先前的自我測試自己重寫了一份逐行的
    比對，於是它一直在驗一個「本體早就不是這樣做」的邏輯——正對照過關，
    真正的版面（條列式，檔名與雜湊分屬不同行）卻一條都沒比到。
    """
    for m in HASH_CLAIM.finditer(text):
        start = m.start("hash")
        window = text[max(0, start - HASH_WINDOW):start]
        # ⭐ **視窗跨行，但不跨空行與標題。** 條列式的寫法本來就把檔名與
        # 雜湊拆成兩行（「- 影片：…mp4」／「- SHA-256：`…`」），逐行配對
        # 對這種版面**結構上看不到**——而那正是推廣片與發行產物的標準寫法，
        # 於是這一層長年綠燈卻一次都沒比過那幾支影片（2026-08-29 用突變
        # 測試發現：把主預告的雜湊改錯一個字元，掃描照樣通過）。
        # 空行與標題是段落邊界，跨過去就會把上一段的檔名黏到這一段的雜湊上。
        for stop in ("\n\n", "\n#"):
            if stop in window:
                window = window[window.rindex(stop) + len(stop):]
        if "SHA" not in window.upper() and "雜湊" not in window:
            continue
        # 明講是「當時」的值＝刻意保留的歷史紀錄，不是現況宣稱。
        # ⚠ 只看雜湊自己那一行與前一行——用整個段落找「當時」會讓
        # 段落裡任何一句話都能豁免整段的雜湊。
        nl = text.find("\n", start)
        near = "\n".join(window.splitlines()[-2:]) + text[start:nl if nl >= 0 else len(text)]
        if "當時" in near:
            skipped["雜湊（標明是當時的值）"] += 1
            continue
        cands = [x.group(0) for x in HASH_FILE.finditer(window)]
        if not cands:
            continue
        # ⭐ **取最近的那一個。** 同一個視窗裡常有兩個檔名
        #（「見 `docs/re/47`…原始 `KI.EXE` SHA-256 `…`」），
        # 配到遠的那個就是憑空捏一條假陽性。
        rel = cands[-1].rstrip(".，。、)）」]")
        target = resolve(rel, doc)
        if target is None:
            skipped["雜湊（檔案不在工作區）"] += 1
            continue
        # `workplace/` 底下是素材與 IDA 資料庫。**`.i64` 的雜湊本來就會漂**
        #（`idat` 一開啟就改寫它，CLAUDE.md §4.1），筆記記的是「這個結論
        # 是在哪一份資料庫上驗的」——那是出處紀錄，不是現況宣稱。
        if os.path.relpath(target, REPO).startswith("workplace" + os.sep):
            skipped["雜湊（workplace 素材／IDA 資料庫）"] += 1
            continue
        yield text.count("\n", 0, start) + 1, target, m.group("hash").lower()


def check_hashes(problems, skipped):
    """文件宣稱的檔案雜湊 vs 實際。"""
    for doc in doc_paths():
        rel_doc = os.path.relpath(doc, REPO)
        # 按日期的台帳記的是「當時量到的值」，不是現況。
        if rel_doc in ("RESEARCH-LOG.md", "WORKLIST.md"):
            continue
        for i, target, claimed in hash_claims(read(doc), doc, skipped):
            with open(target, "rb") as fh:
                actual = hashlib.sha256(fh.read()).hexdigest()
            if actual != claimed:
                problems.append((
                    rel_doc, i, "雜湊過期",
                    f"{os.path.relpath(target, REPO)}："
                    f"文件說 {claimed[:16]}…，實際 {actual[:16]}…"))


def script_image_tags():
    """腳本與 Dockerfile 實際使用的映像標籤。"""
    used = {}
    for dirpath, dirnames, filenames in os.walk(REPO):
        dirnames[:] = [d for d in dirnames
                       if d not in (".git", "workplace", "dist", "dist-all",
                                    "node_modules", "build")]
        for fn in filenames:
            if not (fn.endswith((".sh", ".Dockerfile")) or fn == "Dockerfile"):
                continue
            for m in IMAGE_TAG.finditer(read(os.path.join(dirpath, fn))):
                used.setdefault(m.group(1), set()).add(m.group(2))
    return used


def check_image_tags(problems, skipped):
    """文件寫的映像標籤 vs 腳本實際用的。

    ⚠ 只在「腳本確實用到這個映像名」時才比。腳本沒提到的映像名
    （別的專案、歷史紀錄）一律放行。
    """
    used = script_image_tags()
    for doc in doc_paths():
        rel_doc = os.path.relpath(doc, REPO)
        if rel_doc in ("RESEARCH-LOG.md", "WORKLIST.md"):
            continue
        for i, line in enumerate(read(doc).splitlines(), 1):
            for m in IMAGE_TAG.finditer(line):
                name, tag = m.group(1), m.group(2)
                if name not in used or tag in used[name]:
                    continue
                problems.append((
                    rel_doc, i, "映像標籤過期",
                    f"{name}:{tag}，腳本用的是 {'、'.join(sorted(used[name]))}"))


def check_flags(problems, skipped):
    """文件引用的 `-旗標` 是否還在 `cmd/` 裡定義。"""
    defined = set()
    cmd_dir = os.path.join(REPO, "cmd")
    if not os.path.isdir(cmd_dir):
        skipped["旗標（沒有 cmd/）"] += 1
        return
    for dirpath, _, filenames in os.walk(cmd_dir):
        for fn in filenames:
            if fn.endswith(".go"):
                defined |= set(FLAG_DEF.findall(read(os.path.join(dirpath, fn))))
    for doc in doc_paths():
        rel_doc = os.path.relpath(doc, REPO)
        if rel_doc in ("RESEARCH-LOG.md", "WORKLIST.md"):
            continue
        for i, line in enumerate(read(doc).splitlines(), 1):
            for m in FLAG_CITED.finditer(line):
                flag = m.group(1)
                if flag in defined or flag in FLAG_ALLOW:
                    continue
                # 只認「明講是 wlgame／wlsim／wlview／wlshot 的旗標」那些行，
                # 否則會把 shell 與工具的選項一起抓進來。
                if not re.search(r"wl(game|sim|view|shot|audio|android)", line):
                    continue
                problems.append((rel_doc, i, "旗標不存在", f"-{flag}"))


def check_coverage(problems, skipped):
    """`docs/re/21` 宣稱的覆蓋率 vs 重算。

    census.tsv 在 `workplace/`（gitignore），拿不到就跳過並明講。
    """
    census = os.path.join(REPO, "workplace", "ida", "dosv", "census", "census.tsv")
    doc = os.path.join(REPO, "docs", "re", "21-function-census.md")
    if not os.path.isfile(census) or not os.path.isfile(doc):
        skipped["覆蓋率（沒有 census.tsv）"] += 1
        return
    try:
        out = subprocess.run(
            [sys.executable, os.path.join(REPO, "tools", "re_coverage.py"), census],
            capture_output=True, text=True, cwd=REPO, timeout=300)
    except (OSError, subprocess.SubprocessError):
        skipped["覆蓋率（重算失敗）"] += 1
        return
    if out.returncode != 0:
        skipped["覆蓋率（重算失敗）"] += 1
        return
    m = re.search(r"\|\s*T1[^|]*\|\s*\*?\*?(\d+)", out.stdout)
    n = re.search(r"\|\s*T4[^|]*\|\s*\*?\*?(\d+)", out.stdout)
    if not m or not n:
        skipped["覆蓋率（讀不出重算結果）"] += 1
        return
    t1, t4 = int(m.group(1)), int(n.group(1))
    text = read(doc)
    claimed = re.search(r"\|\s*T1[^|]*\|\s*\*?\*?(\d+)", text)
    if claimed and int(claimed.group(1)) != t1:
        problems.append((
            "docs/re/21-function-census.md", 0, "覆蓋率過期",
            f"文件寫 T1 = {claimed.group(1)}，重算是 {t1}（T4 = {t4}）"))


# ── 第五層：未解列數 ─────────────────────────────────────────
# `docs/re/43` 是生成的（`check.sh` 每次重出），而別的文件會把它的總數抄過去
# 當摘要。抄過去那一份不會自己更新——README 的「456 列」就這樣掛了一個星期。
OPEN_Q_TOTAL = re.compile(r"(\d+)\s*列分布在\s*(\d+)\s*份文件")


def check_open_questions(problems, skipped):
    """別的文件抄過去的未解列數 vs `docs/re/43` 現在的數字。"""
    src = os.path.join(REPO, "docs", "re", "43-open-questions.md")
    if not os.path.isfile(src):
        skipped["未解列數（沒有 docs/re/43）"] += 1
        return
    truth = OPEN_Q_TOTAL.search(read(src))
    if not truth:
        skipped["未解列數（43 讀不出總數）"] += 1
        return
    for doc in doc_paths():
        rel = os.path.relpath(doc, REPO)
        # ⚠ 逐輪紀錄裡的數字是**當時**的數字，不是斷言（同 check_image_tags）。
        if rel in ("RESEARCH-LOG.md", "WORKLIST.md") or rel.endswith("43-open-questions.md"):
            continue
        for i, line in enumerate(read(doc).splitlines(), 1):
            m = OPEN_Q_TOTAL.search(line)
            if m and m.groups() != truth.groups():
                problems.append((
                    rel, i, "未解列數過期",
                    f"寫 {m.group(1)} 列／{m.group(2)} 份，"
                    f"`docs/re/43` 現在是 {truth.group(1)} 列／{truth.group(2)} 份"))


# ── 第六層：規格份數與狀態分布 ───────────────────────────────
# `docs/spec/` 每加一份，抄過去的摘要就舊一天，而它的格式完全正確、
# 連結也通——只有數字是舊的（本檔開頭那一段講的第三類）。
# 這一層照 `docs/spec/*.md` 自己重數，不信任何抄本。
SPEC_SUMMARY = re.compile(
    r"(\d+)\s*份\**（不含索引[^）]*）：\**(\d+)\s*CONFORMED\**／\**(\d+)\s*READY")
SPEC_STATUS = re.compile(r"狀態：([A-Za-z]+)")


def spec_truth():
    """回傳 (份數, CONFORMED, READY)；沒有 `docs/spec/` 就 None。"""
    d = os.path.join(REPO, "docs", "spec")
    if not os.path.isdir(d):
        return None
    total = conformed = ready = 0
    for fn in sorted(os.listdir(d)):
        if not fn.endswith(".md") or fn.startswith("00-index") or fn == "TEMPLATE.md":
            continue
        total += 1
        m = SPEC_STATUS.search(read(os.path.join(d, fn)))
        st = m.group(1) if m else ""
        if st == "CONFORMED":
            conformed += 1
        elif st == "READY":
            ready += 1
    return total, conformed, ready


def check_spec_counts(problems, skipped):
    """別的文件抄過去的規格份數 vs `docs/spec/` 現在的檔案。"""
    truth = spec_truth()
    if truth is None:
        skipped["規格份數（沒有 docs/spec）"] += 1
        return
    for doc in doc_paths():
        rel = os.path.relpath(doc, REPO)
        # ⚠ 逐輪紀錄裡的數字是**當時**的數字，不是斷言（同 check_open_questions）。
        if rel in ("RESEARCH-LOG.md", "WORKLIST.md"):
            continue
        for i, line in enumerate(read(doc).splitlines(), 1):
            m = SPEC_SUMMARY.search(line)
            if m and tuple(int(x) for x in m.groups()) != truth:
                problems.append((
                    rel, i, "規格份數過期",
                    f"寫 {m.group(1)} 份／{m.group(2)} CONFORMED／{m.group(3)} READY，"
                    f"`docs/spec/` 現在是 {truth[0]} 份／{truth[1]}／{truth[2]}"))


def selftest():
    """正對照：**先證明每一層抓得到，再相信它說沒事。**

    綠燈可能是「沒有問題」，也可能是「這一層根本沒在比」——兩者在輸出上
    長得一樣（`CLAUDE.md` §7 第 21 條）。所以每一層都餵一筆確定錯的資料，
    抓不到就當場失敗。
    """
    import tempfile

    ok = True

    def want(label, found, expect=True):
        nonlocal ok
        good = bool(found) == expect
        ok = ok and good
        print(f"  {'✓' if good else '✗'} {label}")

    with tempfile.TemporaryDirectory() as tmp:
        target = os.path.join(tmp, "probe.png")
        with open(target, "wb") as fh:
            fh.write(b"probe")
        real = hashlib.sha256(b"probe").hexdigest()
        doc = os.path.join(tmp, "doc.md")

        def scan(body, fn=None):
            with open(doc, "w", encoding="utf-8") as fh:
                fh.write(body)
            sink = {"雜湊（檔案不在工作區）": 0, "雜湊（標明是當時的值）": 0,
                    "雜湊（workplace 素材／IDA 資料庫）": 0}
            out = []
            for i, target, claimed in hash_claims(body, doc, sink):
                with open(target, "rb") as fh:
                    if hashlib.sha256(fh.read()).hexdigest() != claimed:
                        out.append((i, "雜湊"))
            return out

        bad = "0" * 64
        want("雜湊：擋下錯的",
             scan(f"- `probe.png` SHA-256 `{bad}`"))
        want("雜湊：對的不誤報",
             scan(f"- `probe.png` SHA-256 `{real}`"), expect=False)
        # ⭐ **條列式（檔名與雜湊分屬不同行）是推廣片與發行產物的標準寫法。**
        # 這個正對照擋的是「逐行配對」那種寫法回鍋——它會讓整批產物的雜湊
        # 靜默地不被比對，而輸出跟「沒問題」一模一樣。
        bullets = ("- 影片：[`probe.png`](probe.png)\n"
                   "- 規格：48.468 秒、1280×720、30 fps、H.264／AAC\n"
                   "- SHA-256：`{}`")
        want("雜湊：條列跨行也要比（擋下錯的）", scan(bullets.format(bad)))
        want("雜湊：條列跨行對的不誤報", scan(bullets.format(real)), expect=False)
        # 空行是段落邊界：上一段的檔名不該黏到這一段沒有檔名的雜湊上。
        want("雜湊：不跨段落配對",
             scan(f"- 工具：[`probe.png`](probe.png)\n\n產物 SHA-256：`{bad}`"),
             expect=False)
        # 標明「當時」的歷史值不算現況宣稱。
        want("雜湊：標明當時的值放行",
             scan(f"- `probe.png` SHA-256 `{bad}`（2026-08-11 當時的版本）"),
             expect=False)

    used = script_image_tags()
    want("映像標籤：腳本裡真的找得到標籤", used)
    # 隨便挑一個腳本用到的映像，把日期改掉應該要被抓到。
    hit = None
    for name, tags in used.items():
        if tags:
            hit = (name, sorted(tags)[0])
            break
    if hit:
        name, tag = hit
        fake = "20200101"
        line = f"見 `{name}:{fake}`"
        m = IMAGE_TAG.search(line)
        want("映像標籤：擋下對不上的",
             m and m.group(1) in used and m.group(2) not in used[m.group(1)])

    defined = set()
    for dirpath, _, filenames in os.walk(os.path.join(REPO, "cmd")):
        for fn in filenames:
            if fn.endswith(".go"):
                defined |= set(FLAG_DEF.findall(read(os.path.join(dirpath, fn))))
    want("旗標：讀得到 cmd/ 的旗標定義", len(defined) > 20)
    want("旗標：擋下不存在的", "open-talk" not in defined and "open-talk-index" in defined)

    truth = OPEN_Q_TOTAL.search(read(os.path.join(REPO, "docs", "re", "43-open-questions.md")))
    want("未解列數：讀得到 docs/re/43 的總數", truth)
    if truth:
        n, d = truth.groups()
        want("未解列數：擋下對不上的",
             OPEN_Q_TOTAL.search(f"共 {int(n) + 1} 列分布在 {d} 份文件").groups() != truth.groups())
        want("未解列數：對的不誤報",
             OPEN_Q_TOTAL.search(f"共 {n} 列分布在 {d} 份文件").groups() != truth.groups(),
             expect=False)

    truth = spec_truth()
    want("規格份數：數得到 docs/spec 的份數", truth and truth[0] > 20)
    if truth:
        n, c, r = truth
        want("規格份數：擋下對不上的",
             tuple(int(x) for x in SPEC_SUMMARY.search(
                 f"**{n + 1} 份**（不含索引與 X）：**{c} CONFORMED**／{r} READY"
             ).groups()) != truth)
        want("規格份數：對的不誤報",
             tuple(int(x) for x in SPEC_SUMMARY.search(
                 f"**{n} 份**（不含索引與 X）：**{c} CONFORMED**／{r} READY"
             ).groups()) != truth,
             expect=False)

    print("正對照" + ("通過" if ok else "失敗"))
    return 0 if ok else 1


def main():
    if "--selftest" in sys.argv:
        print("過期斷言掃描自我測試（正對照）")
        return selftest()
    problems = []
    skipped = {}
    for name in ("雜湊（檔案不在工作區）", "雜湊（標明是當時的值）", "雜湊（workplace 素材／IDA 資料庫）",
                 "旗標（沒有 cmd/）",
                 "覆蓋率（沒有 census.tsv）", "覆蓋率（重算失敗）",
                 "未解列數（沒有 docs/re/43）", "未解列數（43 讀不出總數）",
                 "覆蓋率（讀不出重算結果）", "規格份數（沒有 docs/spec）"):
        skipped[name] = 0

    check_hashes(problems, skipped)
    check_image_tags(problems, skipped)
    check_flags(problems, skipped)
    check_coverage(problems, skipped)
    check_open_questions(problems, skipped)
    check_spec_counts(problems, skipped)

    # ⭐ **跳過了什麼一定要印出來。** 「沒有找到問題」與「那一層根本沒跑」
    # 在輸出上長得一樣，而後者才是危險的（CLAUDE.md §7 第 21 條）。
    noted = [f"{k} × {v}" for k, v in skipped.items() if v]
    if noted:
        print("（跳過：" + "；".join(noted) + "）")

    if not problems:
        print(f"過期斷言掃描通過：{len(doc_paths())} 份文件")
        return 0

    print(f"過期斷言掃描：{len(problems)} 筆\n")
    by_kind = {}
    for doc, ln, kind, what in problems:
        by_kind.setdefault(kind, []).append((doc, ln, what))
    for kind in sorted(by_kind):
        print(f"## {kind}（{len(by_kind[kind])}）")
        for doc, ln, what in by_kind[kind]:
            where = f"{doc}:{ln}" if ln else doc
            print(f"  {where}  {what}")
        print()
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
