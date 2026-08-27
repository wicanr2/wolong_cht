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

四層檢查，**每一層都只在能拿到真值時才跑**（缺原版素材、缺 census 就跳過
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
HASH_CLAIM = re.compile(
    r"(?P<file>[A-Za-z0-9_\-./]+\.(?:png|mp4|ogg|wav|json|md|go|py|sh))"
    r"(?P<mid>.{0,80}?)`?(?P<hash>[0-9a-fA-F]{64})`?",
    re.S,
)

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


def check_hashes(problems, skipped):
    """文件宣稱的檔案雜湊 vs 實際。"""
    for doc in doc_paths():
        rel_doc = os.path.relpath(doc, REPO)
        # 按日期的台帳記的是「當時量到的值」，不是現況。
        if rel_doc in ("RESEARCH-LOG.md", "WORKLIST.md"):
            continue
        for i, line in enumerate(read(doc).splitlines(), 1):
            for m in HASH_CLAIM.finditer(line):
                mid = m.group("mid")
                if "SHA" not in mid.upper() and "雜湊" not in mid:
                    continue
                target = resolve(m.group("file"), doc)
                if target is None:
                    skipped["雜湊（檔案不在工作區）"] += 1
                    continue
                with open(target, "rb") as fh:
                    actual = hashlib.sha256(fh.read()).hexdigest()
                if actual != m.group("hash").lower():
                    problems.append((
                        rel_doc, i, "雜湊過期",
                        f"{os.path.relpath(target, REPO)}："
                        f"文件說 {m.group('hash')[:16]}…，實際 {actual[:16]}…"))


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

        def scan(body, fn):
            with open(doc, "w", encoding="utf-8") as fh:
                fh.write(body)
            found = []
            for i, line in enumerate(body.splitlines(), 1):
                for m in HASH_CLAIM.finditer(line):
                    if "SHA" not in m.group("mid").upper() and "雜湊" not in m.group("mid"):
                        continue
                    t = resolve(m.group("file"), doc)
                    if t is None:
                        continue
                    with open(t, "rb") as fh:
                        if hashlib.sha256(fh.read()).hexdigest() != m.group("hash").lower():
                            found.append((i, "雜湊"))
            return found

        bad = "0" * 64
        want("雜湊：擋下錯的",
             scan(f"- `probe.png` SHA-256 `{bad}`", None))
        want("雜湊：對的不誤報",
             scan(f"- `probe.png` SHA-256 `{real}`", None), expect=False)

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

    print("正對照" + ("通過" if ok else "失敗"))
    return 0 if ok else 1


def main():
    if "--selftest" in sys.argv:
        print("過期斷言掃描自我測試（正對照）")
        return selftest()
    problems = []
    skipped = {}
    for name in ("雜湊（檔案不在工作區）", "旗標（沒有 cmd/）",
                 "覆蓋率（沒有 census.tsv）", "覆蓋率（重算失敗）",
                 "覆蓋率（讀不出重算結果）"):
        skipped[name] = 0

    check_hashes(problems, skipped)
    check_image_tags(problems, skipped)
    check_flags(problems, skipped)
    check_coverage(problems, skipped)

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
