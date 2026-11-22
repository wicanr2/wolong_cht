"""Package public engine executables and check the archive allowlist in Docker."""
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import struct
import subprocess
import tarfile
import zipfile

import denylist

ROOT = Path("/repo")
RAW = Path("/out/.work/raw")
WORK = Path("/out/public-release")
PATCH = Path("/delivery/patch")
VERSION = os.environ["WOLONG_RELEASE_VERSION"]
assert re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", VERSION)
APPS = ("wlgame", "wlview", "wlsim", "wlshot", "wlaudio")
PLATFORMS = {"linux-amd64": ("", ".tar.gz"), "windows-amd64": (".exe", ".zip"),
             "macos-universal": ("", ".tar.gz")}


def sha(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def check_arch(data, platform):
    if platform == "linux-amd64":
        assert data[:6] == b"\x7fELF\x02\x01" and data[18:20] == b"\x3e\x00"
    elif platform == "windows-amd64":
        assert data[:2] == b"MZ"
        offset, = struct.unpack_from("<I", data, 60)
        assert data[offset:offset+6] == b"PE\0\0\x64\x86"
    else:
        assert struct.unpack_from(">II", data) == (0xcafebabe, 2)
        assert {struct.unpack_from(">I", data, 8 + 20 * i)[0] for i in range(2)} == {0x1000007, 0x100000c}


def write_checksums(folder, name="SHA256SUMS.txt"):
    files = sorted(f for f in folder.rglob("*") if f.is_file() and f.name != name)
    (folder/name).write_text("".join(f"{sha(f)}  {f.relative_to(folder)}\n" for f in files))


def desktop():
    for platform, (suffix, archive) in PLATFORMS.items():
        name = f"wolong-remake-engine-{platform}-{VERSION}"
        stage = WORK / name
        stage.mkdir()
        for app in APPS:
            filename = app + suffix
            src = RAW / platform / filename
            binary = src.read_bytes()
            check_arch(binary, platform)
            if app == "wlgame":
                assert VERSION.encode() in binary
            shutil.copy2(src, stage/filename)
        shutil.copy2(ROOT/"LICENSE", stage/"LICENSE")
        (stage/"translations").mkdir()
        shutil.copy2(ROOT/"translations/corrections.json", stage/"translations/corrections.json")
        entry = "wlgame.exe" if suffix else "./wlgame"
        (stage/"README.md").write_text(
            f"# 臥龍傳 {VERSION}\n\n"
            "公開引擎版。遊戲資料與字型由玩家自備。\n\n"
            f"```text\n{entry} -orig <原版資料目錄> -font <倚天字型目錄> -save-file <可寫的存檔路徑>\n```\n\n"
            "請指定松崗 DOS/V 繁中版資料。原版音訊和 HD 圖庫可由本機路徑指定。\n"
            "引擎授權詳見 LICENSE。\n")
        bad, _, source_count = denylist.scan(str(ROOT), str(stage))
        assert not bad and source_count, bad
        write_checksums(stage)
        expected = {f"{name}/{p.relative_to(stage)}": sha(p) for p in stage.rglob("*") if p.is_file()}
        target = PATCH/(name+archive)
        if suffix:
            with zipfile.ZipFile(target,"x",zipfile.ZIP_DEFLATED,compresslevel=5) as package:
                for key in sorted(expected):
                    package.write(stage/key.split("/",1)[1],key)
            with zipfile.ZipFile(target) as package:
                assert set(package.namelist()) == set(expected)
                for key, value in expected.items():
                    assert hashlib.sha256(package.read(key)).hexdigest() == value, key
        else:
            with tarfile.open(target,"x:gz",compresslevel=5) as package:
                package.add(stage,arcname=name)
            with tarfile.open(target,"r:gz") as package:
                assert {m.name for m in package if m.isfile()} == set(expected)
                for key, value in expected.items():
                    with package.extractfile(key) as stream:
                        assert hashlib.file_digest(stream,"sha256").hexdigest() == value, key
        print(target.name, len(expected), flush=True)


def manifest():
    packages = sorted(p for p in PATCH.iterdir() if p.suffix in (".zip", ".gz", ".apk"))
    assert len(packages) == 4
    names = {f"wolong-remake-engine-{key}-{VERSION}{v[1]}" for key,v in PLATFORMS.items()}
    names.add(f"wolong-remake-engine-android-arm64-x86_64-{VERSION}.apk")
    assert {p.name for p in packages} == names
    with zipfile.ZipFile(next(p for p in packages if p.suffix == ".apk")) as apk:
        assert apk.testzip() is None
        assert not any(p.startswith("assets/") for p in apk.namelist())
        assert {p for p in apk.namelist() if p.startswith("lib/")} == {
            "lib/arm64-v8a/libgojni.so","lib/x86_64/libgojni.so"}
    data = {"version": VERSION, "source_commit": subprocess.check_output(
                ["git","-C",str(ROOT),"rev-parse","HEAD"],text=True).strip(),
            "package_scope": "公開引擎，原版資料由玩家自備",
            "packages": [{"file": p.name, "bytes": p.stat().st_size, "sha256": sha(p)} for p in packages]}
    (PATCH/"engine-manifest.json").write_text(json.dumps(data,ensure_ascii=False,indent=2)+"\n")
    write_checksums(PATCH)


if __name__ == "__main__":
    import sys
    {"desktop":desktop,"manifest":manifest}[sys.argv[1]]()
