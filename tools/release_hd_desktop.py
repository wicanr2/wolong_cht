"""組裝本機三平台完整版。用環境變數指定版號，工作目錄須在 Docker 掛載。"""
import hashlib
import json
import os
from pathlib import Path
import plistlib
import re
import shutil
import stat
import subprocess
import sys
import tarfile
import zipfile


VERSION = os.environ["WOLONG_RELEASE_VERSION"]
if not re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", VERSION):
    raise SystemExit("版號須為 v.主版.次版.修訂版-YYYYMMDD")
ROOT, WORK, OUT = Path("/repo"), Path("/out"), Path("/delivery")
RAW = WORK / ".work/raw"
PACKAGES = OUT / "full-local"
SOURCE = ROOT / "workplace/hd-redraw-local-20261104/faithful"
OLD_AUDIO = ROOT / "dist-all/packages/wolong-remake-linux-amd64-20260911.tar.gz"
APPS = ("wlgame", "wlview", "wlsim", "wlshot", "wlaudio")
HD_FLAGS = {"battle": "battle", "world": "world", "portraits": "portraits",
            "banner": "locations", "locations": "locations"}


def sha(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def write_json(path, data):
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def copy(src, dst):
    dst.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(src, dst)


def selection(src, dst, name):
    copy(src / name, dst / name)
    return json.loads((src / name).read_text())


def image(src, dst, row):
    file = row["file"]
    if Path(file).name != file or sha(src / file) != row["sha256"].lower():
        raise ValueError(f"HD 素材雜湊錯誤：{src / file}")
    copy(src / file, dst / file)


def stage_data():
    data = WORK / "data"
    data.mkdir()
    for source, folder in ((ROOT / "workplace/orig/dosv", "gamedata"), (ROOT / "workplace/eten", "fonts")):
        for path in sorted(source.iterdir()):
            if path.is_file() and not path.name.startswith("."):
                target = data / folder / path.name
                copy(path, target)
                target.chmod(0o444)
    with tarfile.open(OLD_AUDIO, "r:gz") as tar:
        audio = data / "audio"
        audio.mkdir()
        for member in tar:
            path = Path(member.name)
            if member.isfile() and path.parent.name == "audio" and path.suffix == ".ogg":
                with tar.extractfile(member) as source, (audio / path.name).open("wb") as destination:
                    shutil.copyfileobj(source, destination)
    assert {"bgm-0.ogg", "openbgm-0.ogg", "sfx-01.ogg"} <= {p.name for p in audio.iterdir()}
    for file in ("corrections.json", "talk-en.json", "talk-ja.json", "talk-zh-hans.json"):
        assert (ROOT / "translations" / file).is_file()
    for path in sorted((ROOT / "translations").rglob("*")):
        if path.is_file() and path.suffix == ".json":
            copy(path, data / path.relative_to(ROOT))
    art = data / "hd"
    battle = Path("/review/battle-hd-preview-pack")
    hd = selection(battle, art / "battle", "battle-selection.json")
    for row in hd["tiles"] + [hd[key] for key in ("water", "wood", "foliage", "slope", "roof", "ship")] + list(hd["wall"].values()) + list(hd["ground"].values()):
        image(battle, art / "battle", row)
    world = selection(SOURCE / "world", art / "world", "regions.json")
    for row in world["regions"] + world["seams"]:
        record = selection(SOURCE / "world", art / "world", row["record"])
        image(SOURCE / "world", art / "world", record)
    portraits = selection(SOURCE / "portraits", art / "portraits", "selection.json")
    for file in portraits["records"]:
        record = selection(SOURCE / "portraits", art / "portraits", file)
        image(SOURCE / "portraits", art / "portraits", record)
    locations = art / "locations"
    selection(SOURCE, locations, "originals.json")
    for sample in selection(SOURCE, locations, "manifest.json")["samples"]:
        if sample["kind"] == "scene":
            image(SOURCE / "hd-generated", locations / "hd-generated", sample["hd_record"])
    image(SOURCE, locations, selection(SOURCE, locations, "banner-selection.json"))
    print(f"封包素材 {sum(p.is_file() for p in data.rglob('*'))} 檔", flush=True)


def unix_launcher(data, binary, save_base):
    flags = " ".join(f'-hd-{mode} "$data/hd/{folder}"' for mode, folder in HD_FLAGS.items())
    return f'''#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
data={data}
save={save_base}
mkdir -p "$save"
export WOLONG_TALK_CORRECTIONS="$data/translations/corrections.json"
exec {binary} -orig "$data/gamedata" -font "$data/fonts" -audio "$data/audio" -save-file "$save/SAVE.DAT" {flags} "$@"
'''


def readme(path, command):
    path.write_text(f"# 臥龍傳 {VERSION}\n\n啟動方式：{command}\n\n完整版包含松崗 DOS/V 原版資料、音訊、字型及本機 HD 選稿。限本機使用。\n"
                    "首次以原版顯示。將滑鼠移到上緣，點原版／HD切換。戰場及大地圖可縮放 50%～200%。\n",
                    encoding="utf-8")


def stage_packages():
    data = WORK / "data"
    windows = WORK / "windows"
    shutil.copytree(data, windows)
    for app in APPS:
        copy(RAW / "windows-amd64" / f"{app}.exe", windows / f"{app}.exe")
    flags = " ".join(f'-hd-{mode} "%~dp0hd\\{folder}"' for mode, folder in HD_FLAGS.items())
    command = ('@echo off\nsetlocal\nset "SAVE=%LOCALAPPDATA%\\Wolong"\n'
               'if not exist "%SAVE%" mkdir "%SAVE%"\n'
               'set "WOLONG_TALK_CORRECTIONS=%~dp0translations\\corrections.json"\n'
               '"%~dp0wlgame.exe" -orig "%~dp0gamedata" -font "%~dp0fonts" -audio "%~dp0audio" '
               f'-save-file "%SAVE%\\SAVE.DAT" {flags} %*\n')
    (windows / "Play.cmd").write_bytes(command.replace("\n", "\r\n").encode("ascii"))
    readme(windows / "README.md", "解壓 ZIP 後執行 Play.cmd。")
    copy(ROOT / "LICENSE", windows / "LICENSE")

    app = WORK / "Wolong.app"
    shutil.copytree(data, app / "Contents/Resources")
    binaries = app / "Contents/MacOS"
    for name in APPS:
        copy(RAW / "macos-universal" / name, binaries / name)
    launcher = binaries / "launch"
    launcher.write_text(unix_launcher('"$here/../Resources"', '"$here/wlgame"', '"$HOME/Library/Application Support/Wolong"'))
    launcher.chmod(0o755)
    info = {"CFBundleExecutable":"launch", "CFBundleIdentifier":"org.wolong.remake", "CFBundleName":"Wolong",
            "CFBundleShortVersionString":VERSION[2:].split("-")[0], "CFBundleVersion":VERSION.rsplit("-",1)[1],
            "CFBundlePackageType":"APPL", "WolongReleaseVersion":VERSION, "NSHighResolutionCapable":True}
    (app / "Contents/Info.plist").write_bytes(plistlib.dumps(info))
    readme(app / "Contents/Resources/README.md", "解壓後開啟 Wolong.app。")
    copy(ROOT / "LICENSE", app / "Contents/Resources/LICENSE")

    linux = WORK / "AppDir"
    shutil.copytree(data, linux / "usr/share/wolong-remake")
    bin_dir = linux / "usr/bin"
    for name in APPS:
        copy(RAW / "linux-amd64" / name, bin_dir / name)
    entry = linux / "AppRun"
    entry.write_text(unix_launcher('"$here/usr/share/wolong-remake"', '"$here/usr/bin/wlgame"',
                                   '"${XDG_DATA_HOME:-$HOME/.local/share}/wolong"'))
    entry.chmod(0o755)
    desktop = linux / "wolong-remake.desktop"
    copy(ROOT / "packaging/appimage/wolong-remake.desktop", desktop)
    desktop.write_text(re.sub(r"(?m)^X-AppImage-Version=.*$", f"X-AppImage-Version={VERSION}", desktop.read_text()))
    copy(ROOT / "docs/images/wlgame-dosv-natural-remake-skeleton.png", linux / "wolong-remake.png")
    readme(linux / "README.md", "執行 AppImage，或用 --appimage-extract-and-run。")
    copy(ROOT / "LICENSE", linux / "LICENSE")
    name = f"wolong-remake-windows-amd64-{VERSION}"
    with zipfile.ZipFile(PACKAGES / f"{name}.zip", "x", zipfile.ZIP_DEFLATED, compresslevel=5) as bundle:
        for path in sorted(windows.rglob("*")):
            if path.is_file():
                bundle.write(path, (Path(name) / path.relative_to(windows)).as_posix())
    with tarfile.open(PACKAGES / f"wolong-remake-macos-universal-{VERSION}.tar.gz", "x:gz", compresslevel=5) as bundle:
        bundle.add(app, arcname=app.name)
    tools = ("release_hd_desktop.py", "promo_hd_capture.py", "promo_hd_render.sh")
    inputs = [ROOT / "go.mod", ROOT / "go.sum", OLD_AUDIO, *[ROOT / "tools" / name for name in tools]]
    for folder in ("cmd", "internal"):
        inputs.extend(p for p in (ROOT / folder).rglob("*.go") if p.is_file())
    write_json(OUT / "build-inputs.json", {"version":VERSION,
        "source_commit":subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip(),
        "source_files":{str(p.relative_to(ROOT)):sha(p) for p in sorted(inputs)},
        "staged_data":{str(p.relative_to(data)):sha(p) for p in sorted(data.rglob("*")) if p.is_file()},
        "binaries":{str(p.relative_to(RAW)):sha(p) for p in sorted(RAW.rglob("*")) if p.is_file()}})


if __name__ == "__main__":
    {"data": stage_data, "packages":stage_packages}[sys.argv[1]]()
