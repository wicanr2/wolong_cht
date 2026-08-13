#!/usr/bin/env python3
"""在 Docker 內準備與整理 release 目錄。

release.sh 的建置本身已在 Docker；這個小工具把 dist/ 的刪除、建立、複製與
大小報告也放進同一個固定 Python runtime，避免包裝流程在主機執行檔案工作。
只允許操作儲存庫根目錄下的 dist/，避免誤刪其他路徑。
"""
from __future__ import annotations

import shutil
import sys
from pathlib import Path


REPO = Path.cwd().resolve()
DIST = (REPO / "dist").resolve()


def checked(path: str | None = None) -> Path:
    target = (REPO / path).resolve() if path else DIST
    if target != DIST and DIST not in target.parents:
        raise SystemExit(f"只允許操作 {DIST} 底下的 release 路徑：{target}")
    return target


def prepare() -> None:
    if DIST.exists():
        shutil.rmtree(DIST)
    DIST.mkdir(parents=True, exist_ok=True)


def make_dir(path: str) -> None:
    checked(path).mkdir(parents=True, exist_ok=True)


def finalize() -> None:
    shutil.copy2(REPO / "README.md", DIST / "README.md")
    target = DIST / "translations"
    if target.exists():
        shutil.rmtree(target)
    shutil.copytree(REPO / "translations", target)


def report() -> None:
    for child in sorted(DIST.iterdir()):
        total = sum(p.stat().st_size for p in child.rglob("*") if p.is_file()) if child.is_dir() else child.stat().st_size
        print(f"{total / 1024:.0f}K\t{child}")


def main() -> int:
    if len(sys.argv) < 2:
        raise SystemExit("用法：release_fs.py prepare|mkdir <path>|finalize|report")
    cmd = sys.argv[1]
    if cmd == "prepare":
        prepare()
    elif cmd == "mkdir" and len(sys.argv) == 3:
        make_dir(sys.argv[2])
    elif cmd == "finalize":
        finalize()
    elif cmd == "report":
        report()
    else:
        raise SystemExit("用法：release_fs.py prepare|mkdir <path>|finalize|report")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
