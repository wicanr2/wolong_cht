#!/usr/bin/env python3
"""桌面限定封包：重用既有素材過濾與 AppDir，版本、輸入與輸出雜湊可追溯。"""
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tarfile
import zipfile

VERSION = os.environ.get("WOLONG_RELEASE_VERSION", "")
if not re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", VERSION):
    raise SystemExit("版號必須是 v.主版.次版.修訂版-YYYYMMDD")
import release_all_fs as legacy

DEST = Path("/delivery")
RAW = legacy.WORK / "raw"


def digest(path):
    with path.open("rb") as f:
        return hashlib.file_digest(f, "sha256").hexdigest()


def metadata():
    paths = [p for root in (Path("cmd"), Path("internal"), Path("translations"), Path("packaging"))
             for p in root.rglob("*") if p.is_file()]
    paths += [Path(p) for p in ("go.mod", "go.sum", "LICENSE", "tools/release_desktop.sh", "tools/release_desktop_fs.py", "tools/release_all_fs.py")]
    return {str(p): digest(p) for p in sorted(paths)}


def stage():
    for rel in ("full", "release", "promo", "verification"):
        (DEST / rel).mkdir(exist_ok=True)
    for p in (DEST, legacy.DIST):
        if p.stat().st_uid != os.getuid():
            raise SystemExit(f"輸出擁有者不符：{p}")
    # 舊 stage() 會帶 Android 與推廣片；本入口只取桌面共用的素材與 AppDir 函式。
    bundles = {
        "windows-amd64": ["windows-amd64"],
        "macos-universal": ["darwin-amd64", "darwin-arm64"],
    }
    for label, platforms in bundles.items():
        name = f"wolong-remake-{label}-{VERSION}"
        root = legacy.WORK / "stage" / name
        root.mkdir(parents=True)
        for platform in platforms:
            target = root if label.startswith("windows") else root / platform
            shutil.copytree(RAW / platform, target, dirs_exist_ok=True)
            (target / "translations").mkdir()
            shutil.copy2("translations/corrections.json", target / "translations/corrections.json")
        for src, rel in legacy.bundled_trees():
            legacy.copy_bundled_tree(src, root / rel)
        legacy.write_template("README-RELEASE.md", root / "README-RELEASE.md")
        with (root / "README-RELEASE.md").open("a") as f:
            f.write(f"\n本批版本：{VERSION}。僅桌面；原生 Windows／macOS 操作尚待實機驗收。\n")
        shutil.copy2("LICENSE", root / "LICENSE")
        legacy.write_hashes(root)
        if label.startswith("windows"):
            with zipfile.ZipFile(DEST / "full" / f"{name}.zip", "w", zipfile.ZIP_DEFLATED) as z:
                for p in sorted(root.rglob("*")):
                    if p.is_file():
                        z.write(p, str(Path(name) / p.relative_to(root)))
        else:
            with tarfile.open(DEST / "full" / f"{name}.tar.gz", "w:gz") as t:
                t.add(root, arcname=name)
    legacy.appdir()
    desktop = legacy.WORK / "appdir/wolong-remake.desktop"
    desktop.write_text(desktop.read_text().replace("X-AppImage-Version=wolong-remake-dist-all", f"X-AppImage-Version={VERSION}"))
    legacy.write_hashes(legacy.WORK / "appdir")
    info = {"version": VERSION, "source_commit": subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip(),
            "source_files": metadata(), "inputs": {str(p): digest(p) for src, _ in legacy.bundled_trees()
                for p in sorted(src.iterdir()) if p.is_file() and not p.name.startswith(".") and (src != legacy.AUDIO_SRC or p.suffix == ".ogg")},
            "binaries": {str(p.relative_to(RAW)): digest(p) for p in sorted(RAW.rglob("*")) if p.is_file()}}
    (DEST / "build-inputs.json").write_text(json.dumps(info, ensure_ascii=False, indent=2)+"\n")
    (DEST / "release/README.md").write_text(f"# {VERSION}\n\n本輪只產出本機完整版；未產生或發布公開封包。原版資產沒有公開散布授權。\n")
    (DEST / "promo/README.md").write_text(f"# {VERSION}\n\n本次桌面操作修正未重製推廣片，也不沿用舊片作本版驗收。\n")
    (DEST / "README.md").write_text(f"# 臥龍傳 {VERSION}\n\n桌面修正版，內含原版資料，只限本機使用，不可公開散布。\n\n"
        "- Linux：執行 full/ 內的 AppImage。\n- Windows：解壓 ZIP，執行 wlgame.exe。\n"
        "- macOS：解壓 tar.gz，依本機架構執行 darwin-amd64/wlgame 或 darwin-arm64/wlgame。\n\n"
        "Android 不在本次範圍。驗收狀態見 manifest.json 與 verification/。\n")


def manifest():
    packages = sorted((DEST / "full").iterdir())
    if len(packages) != 3 or any(VERSION not in p.name for p in packages):
        raise SystemExit("應恰有三個使用完整版號的桌面封包")
    info = {"version": VERSION, "distributable": False,
            "platforms": ["linux-amd64", "windows-amd64", "darwin-amd64", "darwin-arm64"],
            "android": "excluded", "validation": {
                "appimage": json.loads((DEST / "verification/normal-input-check.json").read_text())
                    if (DEST / "verification/normal-input-check.json").is_file() else "待正常操作複驗",
                "windows_macos": "封包雜湊與目標架構已查核；原生 GUI 尚未驗證"
                    if (DEST / "verification/package-check.json").is_file() else "待封包與原生 GUI 驗證"},
            "packages": [{"path": str(p.relative_to(DEST)), "bytes": p.stat().st_size, "sha256": digest(p)} for p in packages]}
    (DEST / "manifest.json").write_text(json.dumps(info, ensure_ascii=False, indent=2)+"\n")
    lines = [f"{digest(p)}  {p.relative_to(DEST)}" for p in sorted(DEST.rglob("*")) if p.is_file() and p.name != "SHA256SUMS.txt"]
    (DEST / "SHA256SUMS.txt").write_text("\n".join(lines)+"\n")


if __name__ == "__main__":
    {"stage": stage, "manifest": manifest}[sys.argv[1]]()
