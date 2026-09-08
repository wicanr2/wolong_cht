#!/usr/bin/env python3
"""在 Docker 內檢查桌面封包內容、逐檔雜湊、目標架構及 AppImage 實跑身分。"""
import hashlib
import json
from pathlib import Path
import re
import struct
import sys
import tarfile
import zipfile

root = Path(sys.argv[1])
manifest = json.loads((root / "manifest.json").read_text())
version = manifest["version"]
assert re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", version)
inputs = json.loads((root / "build-inputs.json").read_text())
results = []
for pkg in manifest["packages"]:
    p = root / pkg["path"]
    assert hashlib.sha256(p.read_bytes()).hexdigest() == pkg["sha256"], p
    assert version in p.name and p.stat().st_size == pkg["bytes"], p
    if p.suffix == ".AppImage":
        data = p.read_bytes()[:64]
        assert data[:4] == b"\x7fELF" and struct.unpack_from("<H", data, 18)[0] == 62
        runtime = (root / "verification/runtime.sha256").read_text().splitlines()[0].split()[0]
        assert runtime == inputs["binaries"]["linux-amd64/wlgame"]
        assert version in (root / "verification/app.log").read_text()
        results.append({"package": p.name, "architecture": "ELF x86-64", "runtime_binary_matches": True})
        continue
    if p.suffix == ".zip":
        with zipfile.ZipFile(p) as z:
            assert z.testzip() is None
            files = {n: z.read(n) for n in z.namelist() if not n.endswith("/")}
    else:
        with tarfile.open(p) as t:
            files = {m.name: t.extractfile(m).read() for m in t.getmembers() if m.isfile()}
    prefix = next(iter(files)).split("/")[0]
    entries = {n[len(prefix)+1:]: data for n, data in files.items()}
    for line in entries["SHA256SUMS.txt"].decode().splitlines():
        expected, name = line.split("  ", 1)
        assert hashlib.sha256(entries[name]).hexdigest() == expected, name
    assert all(not n.lower().endswith(".apk") for n in entries)
    for source, expected in inputs["inputs"].items():
        source_path = Path(source)
        subdir = {"dosv": "gamedata", "eten": "fonts", "audio": "audio"}[source_path.parent.name]
        assert hashlib.sha256(entries[f"{subdir}/{source_path.name}"]).hexdigest() == expected
    binaries = {}
    for name, data in entries.items():
        if name == "wlgame.exe":
            offset = struct.unpack_from("<I", data, 0x3c)[0]
            assert data[:2] == b"MZ" and data[offset:offset+4] == b"PE\0\0"
            assert struct.unpack_from("<H", data, offset+4)[0] == 0x8664
            binaries[name] = "PE x86-64"
            expected = inputs["binaries"]["windows-amd64/wlgame.exe"]
        elif name in ("darwin-amd64/wlgame", "darwin-arm64/wlgame"):
            magic, cpu = struct.unpack_from("<II", data)
            assert magic == 0xfeedfacf
            assert cpu == (0x01000007 if "amd64" in name else 0x0100000c)
            binaries[name] = "Mach-O " + ("x86-64" if "amd64" in name else "arm64")
            expected = inputs["binaries"][name]
        else:
            continue
        assert hashlib.sha256(data).hexdigest() == expected
    assert len(binaries) == (1 if p.suffix == ".zip" else 2)
    results.append({"package": p.name, "files": len(entries), "binaries": binaries, "all_file_hashes_match": True})
receipt = {"version": version, "packages": results, "native_windows_macos": "未驗證"}
(root / "verification/package-check.json").write_text(json.dumps(receipt, ensure_ascii=False, indent=2)+"\n")
print(json.dumps(receipt, ensure_ascii=False, indent=2))
