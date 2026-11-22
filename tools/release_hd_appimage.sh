#!/usr/bin/env bash
set -euo pipefail
[[ -f /.dockerenv ]]
version=${WOLONG_RELEASE_VERSION:?}
[[ "$version" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]]
runtime=/opt/runtime-x86_64
[[ "$(od -An -tx1 -j 8 -N 3 "$runtime" | xargs)" == "41 49 02" ]]
dst=/delivery/full-local/wolong-remake-linux-amd64-$version.AppImage
[[ ! -e "$dst" && -x /out/AppDir/AppRun ]]
scratch=$(mktemp -d /out/appimage-XXXXXX)
trap 'rm -rf "$scratch"' EXIT
mksquashfs /out/AppDir "$scratch/data.sqfs" -noappend -all-root -comp gzip -processors 3 -no-progress
cat "$runtime" "$scratch/data.sqfs" > "$scratch/output"
chmod 0755 "$scratch/output"
unsquashfs -s -o "$(stat -c %s "$runtime")" "$scratch/output"
mv "$scratch/output" "$dst"
