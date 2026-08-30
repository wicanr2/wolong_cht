#!/usr/bin/env bash
# 建立乾淨、可交付的 dist-all 三平台發行根目錄。
#
# 此檔本身只編排 docker run；所有建置、複製、封裝、雜湊與驗收都在固定容器中。
# 每次明確重建前只清除 dist-all 的已知發行輸出；未知項目會使流程停止。
#
#   tools/release_all.sh [版本日期]        # 預設今天，格式 YYYYMMDD
#   WOLONG_BUNDLE_DATA=0 tools/release_all.sh …   # 不含原版資產的可散布批次
#
# ⚠ **預設是內含遊戲檔案的完整版**（docs/spec/72，使用者裁定 2026-08-22）：
# 四個平台的包裡都有原版資料與倚天字型，解開就能玩，**而整個 dist-all
# 不可外流**。要交付給別人的版本請明講 `WOLONG_BUNDLE_DATA=0`。
#
# 版本字串只在這裡定一次，檔名與 ldflags 都從它長出來。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAMP="${1:-$(date +%Y%m%d)}"
RELEASE_VERSION="wolong-remake-${STAMP}"
STAGING_ROOT="${WOLONG_RELEASE_STAGING:-$REPO_ROOT/dist-all.staging}"
GO_IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go:latest}"
MAC_IMAGE="${WOLONG_MAC_IMAGE:-wolong-osxcross-go:20260828}"
APPIMAGE_IMAGE="${WOLONG_APPIMAGE_IMAGE:-u5cht/appimage:latest}"
UID_GID="$(id -u):$(id -g)"
BUNDLE_DATA="${WOLONG_BUNDLE_DATA:-1}"

if [ "$BUNDLE_DATA" = 0 ]; then
    echo "── 可散布批次：不含原版資產 ──"
else
    echo "⛔ 完整版批次：內含原版資料與倚天字型，dist-all 不可外流 ──"
    [ -f "$REPO_ROOT/workplace/orig/dosv/SINARIO.DAT" ] || {
        echo "找不到 workplace/orig/dosv/SINARIO.DAT" >&2; exit 1; }
fi

run_repo_write() {
    docker run --rm --network none --memory 1g --cpus 1 --pids-limit 128 \
        -u "$UID_GID" -e WOLONG_DIST_ROOT=/src/$(basename "$STAGING_ROOT") \
        -e WOLONG_RELEASE_VERSION="$RELEASE_VERSION" \
        -e WOLONG_BUNDLE_DATA="$BUNDLE_DATA" \
        -v "$REPO_ROOT:/src" -w /src "$GO_IMAGE" "$@"
}

run_build() {
    docker run --rm --network none --memory 2g --cpus 2 --pids-limit 256 \
        -u "$UID_GID" -v "$REPO_ROOT:/src:ro" -v "$STAGING_ROOT:/out" \
        -v wl-gomod:/gomod:ro -v wl-gobuild:/gocache \
        -e RELEASE_VERSION="$RELEASE_VERSION" -w /src "$GO_IMAGE" "$@"
}

run_macos() {
    docker run --rm --network none --memory 3g --cpus 2 --pids-limit 256 \
        -u "$UID_GID" -v "$REPO_ROOT:/src:ro" -v "$STAGING_ROOT:/out" \
        -v wl-gomod:/gomod:ro -v wl-gobuild:/gocache \
        -e RELEASE_VERSION="$RELEASE_VERSION" -w /src "$MAC_IMAGE" "$@"
}

run_appimage() {
    docker run --rm --network none --memory 1g --cpus 1 --pids-limit 128 \
        -u "$UID_GID" -v "$STAGING_ROOT:/out" -w /out "$APPIMAGE_IMAGE" "$@"
}

run_repo_write python3 tools/release_all_fs.py rebuild

# 在容器內建立輸出根目錄，避免 Docker 為不存在的 bind mount 建出 root-owned 目錄。
run_repo_write python3 tools/release_all_fs.py prepare

run_build bash -lc '
    export PATH=/usr/local/go/bin:$PATH GOMODCACHE=/gomod GOCACHE=/gocache GOPROXY=off
    out=/out/.work/raw
    mkdir -p "$out/linux-amd64" "$out/windows-amd64" "$out/linux-arm64-tools"
    go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/linux-amd64/wlgame" ./cmd/wlgame
    go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/linux-amd64/wlview" ./cmd/wlview
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/linux-amd64/wlsim" ./cmd/wlsim
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/linux-amd64/wlshot" ./cmd/wlshot
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/windows-amd64/wlgame.exe" ./cmd/wlgame
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/windows-amd64/wlview.exe" ./cmd/wlview
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/windows-amd64/wlsim.exe" ./cmd/wlsim
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/windows-amd64/wlshot.exe" ./cmd/wlshot
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/linux-arm64-tools/wlsim" ./cmd/wlsim
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/linux-arm64-tools/wlshot" ./cmd/wlshot
'

run_macos bash -lc '
    export GOMODCACHE=/gomod GOCACHE=/gocache GOPROXY=off
    out=/out/.work/raw
    mkdir -p "$out/darwin-amd64" "$out/darwin-arm64"
    for arch in amd64 arm64; do
        if [ "$arch" = amd64 ]; then
            cc=/osxcross/bin/x86_64-apple-darwin24.5-clang
            cxx=/osxcross/bin/x86_64-apple-darwin24.5-clang++
        else
            cc=/osxcross/bin/aarch64-apple-darwin24.5-clang
            cxx=/osxcross/bin/aarch64-apple-darwin24.5-clang++
        fi
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=1 CC="$cc" CXX="$cxx" go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/darwin-$arch/wlgame" ./cmd/wlgame
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=1 CC="$cc" CXX="$cxx" go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/darwin-$arch/wlview" ./cmd/wlview
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/darwin-$arch/wlsim" ./cmd/wlsim
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$out/darwin-$arch/wlshot" ./cmd/wlshot
    done
'

run_repo_write python3 tools/release_all_fs.py stage
run_repo_write python3 tools/release_all_fs.py appdir
run_appimage bash -lc "ARCH=x86_64 /opt/appimagetool.d/usr/bin/appimagetool --no-appstream /out/.work/appdir /out/packages/wolong-remake-linux-amd64-${STAMP}.AppImage"
run_repo_write python3 tools/release_all_fs.py finalise

# 只有 staging 完成編譯、封裝、雜湊與 deny-list 後才交換到 dist-all。
#
# ⭐ `WOLONG_PROMOTE=0` 就停在 staging，不動 `dist-all`。
# 用途是**同時要兩個批次**：磁碟上的完整版留著，另外建一份可散布的拿去上傳。
# 沒有這個開關的話，建可散布批次會把完整版換掉——而那一批要重跑一次
# 跨平台建置才回得來。
#
#   WOLONG_BUNDLE_DATA=0 WOLONG_PROMOTE=0 \
#     WOLONG_RELEASE_STAGING=$PWD/dist-public tools/release_all.sh 20260830
#
if [ "${WOLONG_PROMOTE:-1}" = 0 ]; then
    echo "跳過 promote：批次留在 $STAGING_ROOT"
else
    run_repo_write python3 tools/release_all_fs.py promote
fi

if [ "$BUNDLE_DATA" = 0 ]; then
    echo "完成：$REPO_ROOT/dist-all（可散布）"
else
    echo "完成：$REPO_ROOT/dist-all ⛔ 內含原版資產，不可外流"
fi
