#!/usr/bin/env python3
"""資產 deny-list 掃描（CLAUDE.md §10 的 [HARD] 發行閘）。

    tools/denylist.py            掃 git 追蹤的檔案
    tools/denylist.py dist/      掃指定目錄（發行包）
    tools/denylist.py --selftest 正對照：確認三層真的會動

**不得 commit 或打包任何原版 `.EXE/.COM/.DAT/.MAP/.MDL/.SCH/.MCH/.BRG/.O`，
也不得散布倚天字型。** 唯一例外是 `docs/images/` 底下的展示截圖
（那是渲染結果，不是資產本身）。

三層檢查，一層比一層難繞過：

1. **副檔名** —— 最便宜，但改個名就繞過去了。
2. **內容雜湊** —— 拿 `workplace/orig/` 全部檔案的 SHA-256 當黑名單。
   改名、換副檔名都躲不掉。原版素材不在時這一層自動跳過（會講明）。
3. **檔名族** —— 倚天字型那三個名字（`STDFONT.24` 等），
   它們不在 `workplace/orig/` 裡，雜湊那層抓不到。

第 2 層是重點。前一個 remake 專案的教訓是「gitignore 擋得住手滑，
擋不住『我先複製一份到別的地方測試』」——而 gitignore 是按路徑比對的，
複製過去就失效。**雜湊是按內容比對的，複製不會改變它。**
"""

import hashlib
import os
import subprocess
import sys

# 第 1 層：副檔名。大小寫不分。
BAD_EXT = {
    ".exe", ".com", ".dat", ".map", ".mdl", ".sch", ".mch", ".brg", ".o",
    ".fdi", ".zip", ".rar",
}

# 第 3 層：倚天中文系統的字型檔名族（`CLAUDE.md` §3.6）。
BAD_NAMES = {"stdfont.24", "stdfont.15", "ascfont.24", "ascfont.15"}

# 例外。`docs/images/` 是渲染出來的展示截圖，不是原版資產。
ALLOW_PREFIX = ("docs/images/",)

ORIG = "workplace/orig"


def tracked_files(root):
    """git 追蹤的檔案。掃 repo 時用這個而不是走檔案系統——
    我們要擋的是「進版控」，而不是「存在於工作目錄」。"""
    out = subprocess.run(["git", "-C", root, "ls-files", "-z"],
                         capture_output=True, text=True)
    if out.returncode != 0:
        return None
    return [p for p in out.stdout.split("\0") if p]


def walk_files(root):
    for dirpath, dirnames, names in os.walk(root):
        dirnames[:] = [d for d in dirnames if d != ".git"]
        for n in names:
            full = os.path.join(dirpath, n)
            yield os.path.relpath(full, root)


def sha256(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def orig_hashes(repo):
    """原版素材的 SHA-256 → 相對路徑。取不到就回 None（不是空 dict）——
    **「沒有素材可比」與「比過了沒中」是兩件事**，回報時要分得出來。"""
    base = os.path.join(repo, ORIG)
    if not os.path.isdir(base):
        return None
    table = {}
    for dirpath, _, names in os.walk(base):
        for n in names:
            full = os.path.join(dirpath, n)
            try:
                table[sha256(full)] = os.path.relpath(full, repo)
            except OSError:
                pass
    return table


def scan(repo, target):
    """回傳 (違規清單, 掃過的檔案數, 雜湊表大小或 None)。"""
    if target is None:
        files = tracked_files(repo)
        if files is None:
            print("不是 git repo，改走檔案系統", file=sys.stderr)
            files = list(walk_files(repo))
        base = repo
    else:
        base = target
        files = list(walk_files(target))

    table = orig_hashes(repo)
    bad = []
    for rel in sorted(files):
        posix = rel.replace(os.sep, "/")
        if posix.startswith(ALLOW_PREFIX):
            continue
        full = os.path.join(base, rel)
        if not os.path.isfile(full) or os.path.islink(full):
            continue

        name = os.path.basename(posix).lower()
        _, ext = os.path.splitext(name)
        if ext in BAD_EXT:
            bad.append((posix, f"副檔名 {ext}"))
            continue
        if name in BAD_NAMES:
            bad.append((posix, "倚天字型檔名"))
            continue
        if table:
            try:
                h = sha256(full)
            except OSError:
                continue
            if h in table:
                bad.append((posix, f"內容與原版素材相同（{table[h]}）"))
    return bad, len(files), (None if table is None else len(table))


def selftest(repo):
    """正對照：故意放三種違規進去，確認三層各自真的會動。

    為什麼要有這一支：deny-list 回「通過」時，**「掃過了沒中」與
    「掃的邏輯壞掉了」長得一模一樣**。發行閘沉默地失效是最糟的一種失效，
    因為它只在出事那一次才被發現。這支把正對照變成可重跑的。
    """
    import shutil
    import tempfile

    tmp = tempfile.mkdtemp(prefix="denyself-")
    try:
        os.makedirs(os.path.join(tmp, "sub"))
        cases = []

        # ① 副檔名
        with open(os.path.join(tmp, "palette.brg"), "wb") as f:
            f.write(b"not really a palette")
        cases.append(("palette.brg", "副檔名"))

        # ② 倚天字型檔名（內容是假的 → 只有檔名那層抓得到）
        with open(os.path.join(tmp, "STDFONT.24"), "wb") as f:
            f.write(b"fake")
        cases.append(("STDFONT.24", "檔名"))

        # ③ 改名 ＋ 換副檔名的真原版素材（只有內容雜湊那層抓得到）
        src = os.path.join(repo, ORIG, "dosv", "TALK.DAT")
        hashed = os.path.isfile(src)
        if hashed:
            shutil.copyfile(src, os.path.join(tmp, "sub", "strings.bin"))
            cases.append(("sub/strings.bin", "內容雜湊"))

        # ④ 乾淨檔案，不該被誤報
        with open(os.path.join(tmp, "main.go"), "wb") as f:
            f.write(b"package main\n")

        bad, _, _ = scan(repo, tmp)
        got = {p for p, _ in bad}
        ok = True
        for path, layer in cases:
            if path in got:
                print(f"  ✓ {layer}：擋下 {path}")
            else:
                print(f"  ✗ {layer}：**沒有**擋下 {path}")
                ok = False
        if "main.go" in got:
            print("  ✗ 誤報：乾淨的 main.go 被擋下")
            ok = False
        else:
            print("  ✓ 乾淨檔案沒被誤報")
        if not hashed:
            print(f"  ⚠ 找不到 {ORIG}/dosv/TALK.DAT，**內容雜湊那層沒驗到**")
        return 0 if ok else 1
    finally:
        shutil.rmtree(tmp, ignore_errors=True)


def main():
    repo = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    if len(sys.argv) > 1 and sys.argv[1] == "--selftest":
        print("deny-list 自我測試（正對照）")
        return selftest(repo)
    target = sys.argv[1] if len(sys.argv) > 1 else None
    bad, n, ntable = scan(repo, target)

    where = target or "git 追蹤的檔案"
    print(f"deny-list：掃了 {n} 個檔（{where}）")
    if ntable is None:
        # 這一行不能省。少了它，「沒跑第二層」看起來會跟「跑了沒中」一模一樣。
        print(f"⚠ 找不到 {ORIG}/，**內容雜湊那一層沒有跑**。"
              "改名過的原版素材這次抓不到。")
    else:
        print(f"內容雜湊比對了 {ntable} 個原版檔")

    if not bad:
        print("通過：沒有原版資產")
        return 0
    print(f"\n擋下 {len(bad)} 個：")
    for path, why in bad:
        print(f"  {path}  ← {why}")
    return 1


if __name__ == "__main__":
    sys.exit(main())
