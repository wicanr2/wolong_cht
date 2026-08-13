#!/usr/bin/env python3
"""掃出文件裡指向不存在東西的引用（幽靈）。

    tools/py.sh tools/phantom_scan.py

規則寫著一個不存在的工具，比沒寫更糟——照做的人會先花時間去找它。
待辦掛著一個不存在的檔案，那一項永遠不可能完成，卻會佔著缺口欄位
讓每一輪都繞著它打轉。這兩種都是這支要抓的。

抓四類：

1. **路徑**：反引號或連結裡看起來像倉庫內路徑的字串，實際不存在。
2. **原版素材檔**：文件提到的 `XXX.DAT`／`.MAP`／`.EXE` 等，
   在 `workplace/orig/{dosv,pc98}/` 都找不到。
3. **Go 識別字**：`internal/...`、`cmd/...` 形式的套件路徑不存在。
4. **IDA 符號**：`sub_XXXXX` 不在函式普查的清單裡（需要 census.tsv）。

只回報「確定不存在」的。寧可漏抓也不要誤報——誤報會讓人把檢查關掉。
"""
import os
import re
import sys

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# 反引號或 markdown 連結裡的倉庫內路徑。要求含 `/` 或已知副檔名，
# 否則會把一般名詞當成路徑。
PATH_RE = re.compile(r"`([A-Za-z0-9_./-]+(?:/[A-Za-z0-9_./-]+|\.(?:go|py|sh|md|idc|json|txt))[A-Za-z0-9_./-]*)`")
LINK_RE = re.compile(r"\]\((?!https?:|#)([^)#]+)")
ASSET_RE = re.compile(r"`([A-Z0-9_]{2,12}\.(?:DAT|MAP|MDL|SCH|MCH|BRG|EXE|COM|O|BAT|SYS|CMD|\$\$\$|\d\d))`")
SYM_RE = re.compile(r"`(sub_[0-9A-Fa-f]{4,6})`")

# 這些是刻意的通配或範例，不是實際路徑。
IGNORE = re.compile(
    r"<|>|\*|\.\.\.|^/|^~|"
    r"^(?:docs|tools|internal|cmd)/?$|"
    r"YYYY|NNNN|XXXX|…"
)

# 只在這些副檔名的檔案裡掃。
SCAN_EXT = (".md",)


def repo_files():
    out = []
    for dirpath, dirnames, filenames in os.walk(REPO):
        dirnames[:] = [d for d in dirnames
                       if d not in (".git", "workplace", "node_modules",
                                    "__pycache__", "dist", "dist-all", "build")]
        for fn in filenames:
            if fn.endswith(SCAN_EXT):
                out.append(os.path.join(dirpath, fn))
    return sorted(out)


def exists(rel, base):
    """rel 可能相對於倉庫根，也可能相對於該文件所在目錄。

    本專案慣例以編號簡寫引用文件（`docs/re/22` 指 `22-strategy-command-tree.md`），
    所以最後一段是純數字時要當前綴比對。不這樣做會把幾百個正確引用當成幽靈，
    而誤報會讓人把整個檢查關掉。
    """
    for cand in (os.path.join(REPO, rel), os.path.join(base, rel)):
        if os.path.exists(cand):
            return True
        d, name = os.path.split(cand)
        if name and os.path.isdir(d):
            if any(n.startswith(name + "-") or n == name + ".md"
                   for n in os.listdir(d)):
                return True
    return False


GO_REF = re.compile(r"^((?:internal|cmd|mobile)/[A-Za-z0-9_/]+)[./]([A-Za-z_][A-Za-z0-9_.]*)$")


def go_symbol_exists(pkgdir, symbol):
    """`internal/state.AmountEdit` 這種寫法指的是套件裡的識別字，不是路徑。

    只取最後一段當識別字（`Soldier.PoseStep` → `PoseStep`），
    在該目錄的 .go 檔裡找宣告。找不到才算幽靈。
    """
    d = os.path.join(REPO, pkgdir)
    if not os.path.isdir(d):
        return False
    leaf = symbol.split(".")[-1]
    pat = re.compile(r"\b" + re.escape(leaf) + r"\b")
    for fn in os.listdir(d):
        if not fn.endswith(".go"):
            continue
        try:
            if pat.search(open(os.path.join(d, fn), encoding="utf-8",
                               errors="replace").read()):
                return True
        except OSError:
            pass
    return False


def check_ref(rel, base):
    """回傳問題描述，沒問題回 None。"""
    if exists(rel, base):
        return None
    m = GO_REF.match(rel)
    if m:
        if go_symbol_exists(m.group(1), m.group(2)):
            return None
        if os.path.isdir(os.path.join(REPO, m.group(1))):
            return "Go 識別字不存在"
    return "路徑不存在"


# 文件明確記載「尚未取得」的素材：提到它們是正確的，不是幽靈。
KNOWN_ABSENT = {"STDFONT.24", "STDFONT.15", "SPCFONT.15",
                "ASCFONT.24", "ASCFONT.15"}


def load_orig_names():
    names = set()
    for ver in ("dosv", "pc98", "pc98_fdi"):
        d = os.path.join(REPO, "workplace", "orig", ver)
        if os.path.isdir(d):
            names |= {n.upper() for n in os.listdir(d)}
    return names


def load_symbols():
    """兩版的符號都要收。

    `docs/re/02` 引用的 `sub_1EF24` 是 PC-98 側的符號——只比對 DOS/V 會誤報。
    版本混淆正是 CLAUDE.md §7 第 9 條要防的事，掃描器自己不能犯。
    """
    syms, found = set(), False
    for ver in ("dosv", "pc98"):
        p = os.path.join(REPO, "workplace", "ida", ver, "census", "census.tsv")
        if not os.path.exists(p):
            continue
        found = True
        with open(p, encoding="utf-8", errors="replace") as fh:
            for line in fh:
                f = line.split("\t")
                if f[0] == "FUNC" and len(f) >= 6:
                    syms.add(f[5].strip())
    return syms if found else None


def main():
    orig = load_orig_names()
    syms = load_symbols()
    problems = []

    for path in repo_files():
        rel_doc = os.path.relpath(path, REPO)
        base = os.path.dirname(path)
        try:
            text = open(path, encoding="utf-8", errors="replace").read()
        except OSError:
            continue
        lines = text.splitlines()

        for i, line in enumerate(lines, 1):
            # 一行如果本身就在說「這個東西不存在」，提到它是正確的。
            # 範圍刻意窄：只認明講的否定詞，不做語意判斷。
            if "不存在" in line or "尚未取得" in line or "從來沒有" in line:
                continue
            for m in PATH_RE.finditer(line):
                rel = m.group(1).rstrip("/")
                if IGNORE.search(rel) or not rel:
                    continue
                # 只認看起來屬於本倉庫的前綴，避免把外部路徑當幽靈。
                if not rel.split("/")[0] in ("docs", "tools", "internal", "cmd",
                                             "translations", "packaging", "mobile",
                                             "android", "docker", "workplace"):
                    continue
                verdict = check_ref(rel, base)
                if verdict:
                    problems.append((rel_doc, i, verdict, rel))

            for m in LINK_RE.finditer(line):
                rel = m.group(1).strip()
                if IGNORE.search(rel) or not rel:
                    continue
                if not exists(rel, base):
                    problems.append((rel_doc, i, "連結指不到", rel))

            if orig:
                for m in ASSET_RE.finditer(line):
                    name = m.group(1).upper()
                    if name in KNOWN_ABSENT:
                        continue
                    if name not in orig:
                        problems.append((rel_doc, i, "原版素材不存在", m.group(1)))

            if syms:
                for m in SYM_RE.finditer(line):
                    if m.group(1) not in syms:
                        problems.append((rel_doc, i, "IDA 符號不在資料庫", m.group(1)))

    if syms is None:
        print("⚠ 找不到 census.tsv，**IDA 符號那一層沒有跑**。"
              "先跑 tools/ida.sh script dosv tools/ida_function_census.idc")

    if not problems:
        print(f"幽靈掃描通過：{len(repo_files())} 份文件，沒有指向不存在東西的引用")
        return 0

    print(f"幽靈掃描：{len(problems)} 筆\n")
    by_kind = {}
    for doc, ln, kind, what in problems:
        by_kind.setdefault(kind, []).append((doc, ln, what))
    for kind in sorted(by_kind):
        print(f"## {kind}（{len(by_kind[kind])}）")
        for doc, ln, what in by_kind[kind]:
            print(f"  {doc}:{ln}  {what}")
        print()
    return 1


if __name__ == "__main__":
    sys.exit(main())
