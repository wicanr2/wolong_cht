"""驗證實際封包內容、架構、成片與本機交付 SHA-256。"""
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import stat
import struct
import subprocess
import sys
import tarfile
import zipfile


VERSION = os.environ["WOLONG_RELEASE_VERSION"]
assert re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", VERSION)
OUT = Path("/delivery")
WORK = Path("/out")
APPS = ("wlgame", "wlview", "wlsim", "wlshot", "wlaudio")
INPUTS = json.loads((OUT / "build-inputs.json").read_text())
DATA = INPUTS["staged_data"]
FILM = f"wolong-remake-promo-{VERSION}.mp4"


def sha(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def arch(data, platform):
    if platform == "linux":
        assert data[:6] == b"\x7fELF\x02\x01" and data[18:20] == b"\x3e\0"
    elif platform == "windows":
        assert data[:2] == b"MZ"
        loc, = struct.unpack_from("<I", data, 60)
        assert data[loc:loc+6] == b"PE\0\0\x64\x86"
    else:
        magic, count = struct.unpack_from(">II", data)
        assert magic == 0xcafebabe and count == 2
        assert {struct.unpack_from(">I", data, 8+20*i)[0] for i in range(count)} == {0x1000007, 0x100000c}


def compare_payload(names, open_stream, prefix):
    expected = {prefix+path for path in DATA}
    assert expected <= names, expected-names
    for path, digest in DATA.items():
        with open_stream(prefix+path) as stream:
            assert hashlib.file_digest(stream, "sha256").hexdigest() == digest, path
    return len(expected)


def check_film(directory):
    film = directory / FILM
    metadata = json.loads((directory / "ffprobe.json").read_text())
    video = next(s for s in metadata["streams"] if s["codec_type"] == "video")
    audio = next(s for s in metadata["streams"] if s["codec_type"] == "audio")
    assert (video["codec_name"], video["width"], video["height"], video["r_frame_rate"]) == ("h264",1280,720,"30/1")
    assert audio["codec_name"] == "aac" and audio["channels"] == 2
    assert abs(float(video["duration"])-float(audio["duration"])) < .16
    log = (directory / "media-check.log").read_text()
    assert "black_start:" not in log
    quiet = [float(n) for n in re.findall(r"silence_duration: ([\d.]+)", log)]
    assert len(quiet) == log.count("silence_start:") and all(n < 1 for n in quiet)
    mean = float(re.search(r"mean_volume: ([-\d.]+) dB", log)[1])
    peak = float(re.search(r"max_volume: ([-\d.]+) dB", log)[1])
    assert -35 < mean < -9 and -9 < peak < 0
    result = {"path": f"promo/{FILM}", "duration_seconds": float(metadata["format"]["duration"]),
              "sha256": sha(film), "mean_db": mean, "peak_db": peak, "public_upload": False,
              "longest_quiet_seconds": max(quiet, default=0),
              "static_or_nearly_static_intervals": log.count("freeze_start:")}
    provenance = directory / "source.json"
    if provenance.is_file():
        source = json.loads(provenance.read_text())
        assert result["sha256"] == source["sha256"]
        result.update(capture_version=source["capture_version"], opening_seconds=source["opening_seconds"])
    return result


def verify_android(apk):
    receipt = json.loads((OUT / "smoke/android.json").read_text())
    digest = sha(apk)
    assert receipt["startup"] and receipt["bundled_data_boot"] and receipt["strategy_screenshot_visible"]
    assert receipt["apk_sha256"] == digest
    roots = {"gamedata": "orig", "fonts": "eten", "audio": "audio"}
    expected = {}
    for path, asset_hash in DATA.items():
        root, name = path.split("/", 1)
        if root in roots:
            expected[f"assets/gamedata/{roots[root]}/{name}"] = asset_hash
    with zipfile.ZipFile(apk) as package:
        assert {n for n in package.namelist() if n.startswith("assets/gamedata/")} == set(expected)
        assert {n for n in package.namelist() if n.startswith("lib/")} == {
            "lib/arm64-v8a/libgojni.so", "lib/x86_64/libgojni.so",
        }
        for name, asset_hash in expected.items():
            with package.open(name) as stream:
                assert hashlib.file_digest(stream, "sha256").hexdigest() == asset_hash, name
    assert receipt["packaged_assets"] == len(expected)
    return {"path": str(apk.relative_to(OUT)), "bytes": apk.stat().st_size,
            "sha256": digest, "verified_payload_files": len(expected)}


def write_manifest(manifest):
    (OUT / "manifest.json").write_text(json.dumps(manifest, ensure_ascii=False, indent=2)+"\n")
    lines = [f"{sha(file)}  {file.relative_to(OUT)}\n" for file in sorted(OUT.rglob("*"))
             if file.is_file() and file.name != "SHA256SUMS.txt"]
    (OUT / "SHA256SUMS.txt").write_text("".join(lines))
    print(json.dumps(manifest, ensure_ascii=False, indent=2))


def update_promo(candidate):
    film = check_film(candidate)
    assert abs(film["duration_seconds"] - 47) < .03
    assert (candidate / "opening-contact-sheet.png").is_file()
    subprocess.run(["sha256sum", "-c", str(candidate / "source-SHA256SUMS.txt")], check=True)
    manifest = json.loads((OUT / "manifest.json").read_text())
    assert manifest["version"] == VERSION
    for package in manifest["packages"]:
        path = OUT / package["path"]
        assert path.stat().st_size == package["bytes"] and sha(path) == package["sha256"]
    for path in sorted(candidate.iterdir()):
        assert path.is_file()
        temp = OUT / "promo" / (path.name + ".new")
        shutil.copy2(path, temp)
        temp.replace(OUT / "promo" / path.name)
    for name in ("tools/promo_hd_render.sh", "tools/promo_opening_frames.go", "tools/verify_hd_desktop.py"):
        INPUTS.setdefault("reproduction_tools", {})[name] = sha(Path("/repo") / name)
    (OUT / "build-inputs.json").write_text(json.dumps(INPUTS, ensure_ascii=False, indent=2)+"\n")
    manifest["promo"] = dict(film, opening_seconds=12, opening_source="DOS/V OPEN_S2, S3, S4, S6")
    write_manifest(manifest)


def verify():
    assert INPUTS["version"] == VERSION
    for filename, digest in INPUTS["source_files"].items():
        assert sha(Path("/repo") / filename) == digest, filename
    packages = OUT / "full-local"
    linux = packages / f"wolong-remake-linux-amd64-{VERSION}.AppImage"
    mac = packages / f"wolong-remake-macos-universal-{VERSION}.tar.gz"
    windows = packages / f"wolong-remake-windows-amd64-{VERSION}.zip"
    android = packages / f"wolong-remake-android-arm64-x86_64-{VERSION}.apk"
    expected = {linux, mac, windows, android} if android.is_file() else {linux, mac, windows}
    assert set(packages.iterdir()) == expected

    with zipfile.ZipFile(windows) as package:
        prefix = windows.stem+"/"
        count = compare_payload(set(package.namelist()), package.open, prefix)
        for app in APPS:
            name = prefix+app+".exe"
            contents = package.read(name)
            arch(contents, "windows")
            assert hashlib.sha256(contents).hexdigest() == INPUTS["binaries"]["windows-amd64/"+app+".exe"]
        command = package.read(prefix+"Play.cmd")
        assert b"-hd-battle" in command and b"%LOCALAPPDATA%" in command and b"-save-file" in command
    with tarfile.open(mac, "r:gz") as package:
        count_mac = compare_payload(set(package.getnames()), package.extractfile, "Wolong.app/Contents/Resources/")
        for app in APPS:
            entry = package.getmember("Wolong.app/Contents/MacOS/"+app)
            assert entry.mode & stat.S_IXUSR
            contents = package.extractfile(entry).read()
            arch(contents, "macos")
            assert hashlib.sha256(contents).hexdigest() == INPUTS["binaries"]["macos-universal/"+app]
            if app == "wlgame":
                assert VERSION.encode() in contents
        assert package.getmember("Wolong.app/Contents/MacOS/launch").mode & stat.S_IXUSR
        assert VERSION.encode() in package.extractfile("Wolong.app/Contents/Info.plist").read()

    extracted = WORK / "smoke/squashfs-root"
    count_linux = compare_payload({str(p.relative_to(extracted)) for p in extracted.rglob("*") if p.is_file()},
                                  lambda rel: (extracted/rel).open("rb"), "usr/share/wolong-remake/")
    for app in APPS:
        file = extracted / "usr/bin" / app
        assert file.stat().st_mode & stat.S_IXUSR
        with file.open("rb") as stream:
            arch(stream.read(24), "linux")
        assert sha(file) == INPUTS["binaries"]["linux-amd64/"+app]

    startup = {platform: json.loads((OUT / "smoke" / f"{platform}.json").read_text())
               for platform in ("linux", "windows")}
    assert all(status["startup"] for status in startup.values())
    startup["macos"] = {"native_gui": False, "verified_architectures": ["x86_64", "arm64"]}
    manifest = {
        "version": VERSION, "scope":"本機保存", "distributable":False,
        "packages":[{"path":str(file.relative_to(OUT)), "bytes":file.stat().st_size,
                     "sha256":sha(file), "verified_payload_files":num}
                    for file, num in ((linux,count_linux),(mac,count_mac),(windows,count))],
        "promo":check_film(OUT / "promo"),
        "startup": startup,
    }
    if android in expected:
        manifest["packages"].append(verify_android(android))
        manifest["startup"]["android"] = json.loads((OUT / "smoke/android.json").read_text())
    write_manifest(manifest)


if __name__ == "__main__":
    if len(sys.argv) == 3 and sys.argv[1] == "--update-promo":
        update_promo(Path(sys.argv[2]))
    else:
        verify()
