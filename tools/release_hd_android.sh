#!/usr/bin/env bash
set -euo pipefail
[[ -f /.dockerenv ]]
version=${WOLONG_RELEASE_VERSION:?}
[[ "$version" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]]
[[ ${WOLONG_ANDROID_VERSION_CODE:?} =~ ^[1-9][0-9]*$ ]]
stage=/out/.work/android
bundle="/delivery/full-local/wolong-remake-android-arm64-x86_64-$version.apk"
test ! -e "$bundle"
test -s "$stage/wolong.aar"

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
export HOME="$work/home" GRADLE_USER_HOME="$work/gradle"
export ANDROID_USER_HOME="$HOME/.android" JAVA_TOOL_OPTIONS="-Duser.home=$HOME -Dfile.encoding=UTF-8"
export LANG=C.UTF-8 LC_ALL=C.UTF-8
mkdir -p "$HOME" "$GRADLE_USER_HOME/caches" "$work/src"
cp -r --reflink=auto /seed-cache/gradle-home/caches/modules-2 "$GRADLE_USER_HOME/caches/"
chmod -R u+w "$GRADLE_USER_HOME/caches"
tar --exclude=android/.gradle --exclude=android/app/build --exclude=android/app/libs \
    --exclude=android/app/src/main/assets -cf - android | tar -xf - -C "$work/src"
cd "$work/src"
mkdir -p android/app/libs
cp "$stage/wolong.aar" android/app/libs/wolong.aar
assets=android/app/src/main/assets/gamedata
mkdir -p "$assets"/{orig,eten,audio}
cp -r /out/data/gamedata/. "$assets/orig/"
cp -r /out/data/fonts/. "$assets/eten/"
cp -r /out/data/audio/. "$assets/audio/"

gradle --offline --no-daemon --console=plain --max-workers=3 -p android :app:assembleDebug
apk=android/app/build/outputs/apk/debug/app-debug.apk
apksigner verify --verbose "$apk"
zipalign -c -P 16 -v 4 "$apk" | tail -1
aapt dump badging "$apk" >"$stage/badging.txt"
grep -F "versionName='$version'" "$stage/badging.txt"
grep -F "versionCode='$WOLONG_ANDROID_VERSION_CODE'" "$stage/badging.txt"
unzip -q "$apk" 'lib/*/*.so' -d "$work/elf"

for abi in arm64-v8a x86_64; do
    lib="$work/elf/lib/$abi/libgojni.so"
    readelf -lW "$lib" | awk '/LOAD/ { if ($NF != "0x4000") exit 1; count++ } END { if (!count) exit 1 }'
    grep -aFq "$version" "$lib"
done

python3 - "$apk" "$assets" <<'PYTHON'
import hashlib
from pathlib import Path
import sys
from zipfile import ZipFile

apk, assets = map(Path, sys.argv[1:])
paths = {f"assets/gamedata/{p.relative_to(assets).as_posix()}": p for p in assets.rglob("*") if p.is_file()}
with ZipFile(apk) as package:
    assert package.testzip() is None
    assert {p for p in package.namelist() if p.startswith("assets/gamedata/")} == set(paths)
    for name, source in paths.items():
        with package.open(name) as packed, source.open("rb") as original:
            assert hashlib.file_digest(packed, "sha256").digest() == hashlib.file_digest(original, "sha256").digest(), name
print(f"Android 內嵌素材 {len(paths)} 檔驗證通過")
PYTHON
cp "$apk" "$stage/verified.apk"
mv "$stage/verified.apk" "$bundle"
sha256sum "$bundle"
