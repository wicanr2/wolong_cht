#!/usr/bin/env bash
# 桌面限定修正版；沿用既有工具鏈，不重建 Android、不覆蓋舊包。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${1:?請指定 v.主版.次版.修訂版-YYYYMMDD}"
[[ "$VERSION" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]] || { echo '版號格式不符' >&2; exit 1; }
OUT="$ROOT/dist-all/$VERSION"
STAGE="$ROOT/workplace/desktop-package-$VERSION"
[[ ! -e "$OUT" && ! -e "$STAGE" ]] || { echo '輸出已存在；請提高修訂版號' >&2; exit 1; }
GO_IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go:latest}"
MAC_IMAGE="${WOLONG_MAC_IMAGE:-wolong-osxcross-go:20260828}"
APP_IMAGE="${WOLONG_APPIMAGE_IMAGE:-u5cht/appimage:latest}"
for image in "$GO_IMAGE" "$MAC_IMAGE" "$APP_IMAGE"; do docker image inspect "$image" >/dev/null; done
mkdir -p "$STAGE" "$OUT"
stat -c '%u:%g %n' "$STAGE" "$OUT"
run() {
    local image=$1; shift
    timeout 900 docker run --rm --network none --memory 3g --cpus 2 --pids-limit 256 \
        -u "$(id -u):$(id -g)" -v "$ROOT:/src:ro" -v "$STAGE:/out" -v "$OUT:/delivery" \
        -v wl-gomod:/gomod:ro -v wl-gobuild:/gocache \
        -e GOMODCACHE=/gomod -e GOCACHE=/gocache -e GOPROXY=off -e HOME=/tmp \
        -e WOLONG_RELEASE_VERSION="$VERSION" -e WOLONG_DIST_ROOT=/out -e WOLONG_BUNDLE_DATA=1 \
        -w /src "$image" "$@"
}
run "$GO_IMAGE" bash -c '
    set -e; export PATH=/usr/local/go/bin:$PATH
    for platform in linux-amd64 windows-amd64; do
        mkdir -p /out/.work/raw/$platform
        for cmd in wlgame wlview wlsim wlshot; do
            if [ "$platform" = windows-amd64 ]; then
                GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$WOLONG_RELEASE_VERSION" -o /out/.work/raw/$platform/$cmd.exe ./cmd/$cmd
            else
                go build -trimpath -ldflags "-s -w -X main.version=$WOLONG_RELEASE_VERSION" -o /out/.work/raw/$platform/$cmd ./cmd/$cmd
            fi
        done
    done
'
run "$MAC_IMAGE" bash -c '
    set -e
    for arch in amd64 arm64; do
        cc=/osxcross/bin/x86_64-apple-darwin24.5-clang
        [ "$arch" != arm64 ] || cc=/osxcross/bin/aarch64-apple-darwin24.5-clang
        mkdir -p /out/.work/raw/darwin-$arch
        for cmd in wlgame wlview wlsim wlshot; do
            GOOS=darwin GOARCH=$arch CGO_ENABLED=1 CC=$cc go build -trimpath -ldflags "-s -w -X main.version=$WOLONG_RELEASE_VERSION" -o /out/.work/raw/darwin-$arch/$cmd ./cmd/$cmd
        done
    done
'
run "$GO_IMAGE" python3 tools/release_desktop_fs.py stage
run "$APP_IMAGE" bash -c 'set -e; ARCH=x86_64 /opt/appimagetool.d/usr/bin/appimagetool --no-appstream /out/.work/appdir /delivery/full/wolong-remake-linux-amd64-$WOLONG_RELEASE_VERSION.AppImage'
run "$GO_IMAGE" python3 tools/release_desktop_fs.py manifest
stat -c '%u:%g %n' "$OUT" "$OUT/manifest.json"
echo "桌面候選包：$OUT；尚需對本批 AppImage 執行正常操作驗收。"
