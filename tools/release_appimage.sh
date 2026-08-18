#!/usr/bin/env bash
# 單獨重打 Linux AppImage，放進 dist-all/packages。
#
#   tools/release_appimage.sh [版本日期]        # 預設今天，格式 YYYYMMDD
#
# 與 `tools/release_all.sh` 的差別：那一支重建整個 dist-all（三平台 ＋ 推廣片
# ＋ Android 附件），需要 `dist/promo/` 的四支影片與 Android APK 都在位；
# 這一支只做 Linux amd64 的 AppImage，**其餘產物一個都不動**。
#
# 邊界寫在這裡，不寫在對話裡：
#   - 只用 2 顆 CPU（這台是共用機器）
#   - 只寫 dist-all/packages 與自己的 staging 目錄
#   - 一律 --rm，不清任何 image／volume／別人的 container
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAMP="${1:-$(date +%Y%m%d)}"
VERSION="wolong-remake-${STAMP}"
NAME="wolong-remake-linux-amd64-${STAMP}.AppImage"
STAGING="$REPO_ROOT/dist-appimage.staging"
GO_IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go:latest}"
APPIMAGE_IMAGE="${WOLONG_APPIMAGE_IMAGE:-u5cht/appimage:latest}"
UID_GID="$(id -u):$(id -g)"

rm -rf "$STAGING"
mkdir -p "$STAGING/.work/raw/linux-amd64" "$STAGING/packages"

echo "[1/3] 編譯 linux-amd64（2 核）"
docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -u "$UID_GID" -v "$REPO_ROOT:/src:ro" -v "$STAGING:/out" \
    -v wl-gomod:/gomod:ro -v wl-gobuild:/gocache -w /src "$GO_IMAGE" bash -lc "
    set -e
    export PATH=/usr/local/go/bin:\$PATH GOMODCACHE=/gomod GOCACHE=/gocache GOPROXY=off
    out=/out/.work/raw/linux-amd64
    go build -trimpath -ldflags '-s -w -X main.version=${VERSION}' -o \"\$out/wlgame\" ./cmd/wlgame
    go build -trimpath -ldflags '-s -w -X main.version=${VERSION}' -o \"\$out/wlview\" ./cmd/wlview
    CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X main.version=${VERSION}' -o \"\$out/wlsim\" ./cmd/wlsim
    CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X main.version=${VERSION}' -o \"\$out/wlshot\" ./cmd/wlshot
"

echo "[2/3] 組 AppDir"
docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 1g --cpus 1 --pids-limit 128 \
    -u "$UID_GID" -e WOLONG_DIST_ROOT=/src/$(basename "$STAGING") \
    -e WOLONG_RELEASE_VERSION="$VERSION" \
    -v "$REPO_ROOT:/src" -w /src "$GO_IMAGE" python3 tools/release_all_fs.py appdir

echo "[3/3] 封裝 AppImage"
docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 1g --cpus 1 --pids-limit 128 \
    -u "$UID_GID" -v "$STAGING:/out" -w /out "$APPIMAGE_IMAGE" \
    bash -lc "ARCH=x86_64 /opt/appimagetool.d/usr/bin/appimagetool --no-appstream /out/.work/appdir /out/packages/${NAME}"

mkdir -p "$REPO_ROOT/dist-all/packages"
cp -f "$STAGING/packages/${NAME}" "$REPO_ROOT/dist-all/packages/${NAME}"
chmod +x "$REPO_ROOT/dist-all/packages/${NAME}"
rm -rf "$STAGING"

echo "完成：dist-all/packages/${NAME}"
sha256sum "$REPO_ROOT/dist-all/packages/${NAME}"
