#!/usr/bin/env bash
# 建立乾淨、可交付的 dist-all 三平台發行根目錄。
#
# 此檔本身只編排 docker run；所有建置、複製、封裝、雜湊與驗收都在固定容器中。
# 每次明確重建前只清除 dist-all 的已知發行輸出；未知項目會使流程停止。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAGING_ROOT="${WOLONG_RELEASE_STAGING:-$REPO_ROOT/dist-all.staging}"
GO_IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go:latest}"
MAC_IMAGE="${WOLONG_MAC_IMAGE:-wolong-osxcross-go:20260811-event10-r4}"
APPIMAGE_IMAGE="${WOLONG_APPIMAGE_IMAGE:-u5cht/appimage:latest}"
UID_GID="$(id -u):$(id -g)"

run_repo_write() {
    docker run --rm --network none --memory 1g --cpus 1 --pids-limit 128 \
        -u "$UID_GID" -e WOLONG_DIST_ROOT=/src/$(basename "$STAGING_ROOT") \
        -v "$REPO_ROOT:/src" -w /src "$GO_IMAGE" "$@"
}

run_build() {
    docker run --rm --network none --memory 2g --cpus 2 --pids-limit 256 \
        -u "$UID_GID" -v "$REPO_ROOT:/src:ro" -v "$STAGING_ROOT:/out" \
        -v wl-gomod:/gomod:ro -w /src "$GO_IMAGE" "$@"
}

run_macos() {
    docker run --rm --network none --memory 3g --cpus 2 --pids-limit 256 \
        -u "$UID_GID" -v "$REPO_ROOT:/src:ro" -v "$STAGING_ROOT:/out" \
        -v wl-gomod:/gomod:ro -w /src "$MAC_IMAGE" "$@"
}

run_appimage() {
    docker run --rm --network none --memory 1g --cpus 1 --pids-limit 128 \
        -u "$UID_GID" -v "$STAGING_ROOT:/out" -w /out "$APPIMAGE_IMAGE" "$@"
}

run_repo_write python3 tools/release_all_fs.py rebuild

# 在容器內建立輸出根目錄，避免 Docker 為不存在的 bind mount 建出 root-owned 目錄。
run_repo_write python3 tools/release_all_fs.py prepare

run_build bash -lc '
    export PATH=/usr/local/go/bin:$PATH GOMODCACHE=/gomod GOCACHE=/tmp/gocache GOPROXY=off
    out=/out/.work/raw
    mkdir -p "$out/linux-amd64" "$out/windows-amd64" "$out/linux-arm64-tools"
    go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/linux-amd64/wlgame" ./cmd/wlgame
    go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/linux-amd64/wlview" ./cmd/wlview
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/linux-amd64/wlsim" ./cmd/wlsim
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/linux-amd64/wlshot" ./cmd/wlshot
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/windows-amd64/wlgame.exe" ./cmd/wlgame
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/windows-amd64/wlview.exe" ./cmd/wlview
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/windows-amd64/wlsim.exe" ./cmd/wlsim
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/windows-amd64/wlshot.exe" ./cmd/wlshot
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/linux-arm64-tools/wlsim" ./cmd/wlsim
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/linux-arm64-tools/wlshot" ./cmd/wlshot
'

run_macos bash -lc '
    export GOMODCACHE=/gomod GOCACHE=/tmp/gocache GOPROXY=off
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
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=1 CC="$cc" CXX="$cxx" go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/darwin-$arch/wlgame" ./cmd/wlgame
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=1 CC="$cc" CXX="$cxx" go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/darwin-$arch/wlview" ./cmd/wlview
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/darwin-$arch/wlsim" ./cmd/wlsim
        GOOS=darwin GOARCH="$arch" CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=wolong-remake-20260812" -o "$out/darwin-$arch/wlshot" ./cmd/wlshot
    done
'

run_repo_write python3 tools/release_all_fs.py stage
run_repo_write python3 tools/release_all_fs.py appdir
run_appimage bash -lc 'ARCH=x86_64 /opt/appimagetool.d/usr/bin/appimagetool --no-appstream /out/.work/appdir /out/packages/wolong-remake-linux-amd64-20260812.AppImage'
run_repo_write python3 tools/release_all_fs.py finalise

# 只有 staging 完成編譯、封裝、雜湊與 deny-list 後才交換到 dist-all。
run_repo_write python3 tools/release_all_fs.py promote

echo "完成：$REPO_ROOT/dist-all"
