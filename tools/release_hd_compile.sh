#!/usr/bin/env bash
set -euo pipefail
[[ -f /.dockerenv ]]
version=${WOLONG_RELEASE_VERSION:?}
[[ "$version" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]]
target=${1:?linux-amd64 / windows-amd64 / macos-universal}
raw=/out/.work/raw
cd /repo
apps=(wlgame wlview wlsim wlshot wlaudio)
build() {
    local dir=$1 suffix=$2
    mkdir -p "$raw/$dir"
    for app in "${apps[@]}"; do
        go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$raw/$dir/$app$suffix" "./cmd/$app"
    done
}
case "$target" in
    linux-amd64) GOOS=linux GOARCH=amd64 CGO_ENABLED=1 build "$target" "" ;;
    windows-amd64) GOOS=windows GOARCH=amd64 CGO_ENABLED=0 build "$target" ".exe" ;;
    macos-universal)
        export PATH=/out/go/bin:$PATH
        for arch in amd64 arm64; do
            stem=x86_64
            if [[ "$arch" == arm64 ]]; then stem=aarch64; fi
            cc="/osxcross/bin/$stem-apple-darwin24.5-clang"
            GOOS=darwin GOARCH="$arch" CGO_ENABLED=1 CC="$cc" CXX="$cc++" build "darwin-$arch" ""
        done
        mkdir -p "$raw/macos-universal"
        for app in "${apps[@]}"; do
            /osxcross/bin/lipo -create "$raw/darwin-amd64/$app" "$raw/darwin-arm64/$app" -output "$raw/macos-universal/$app"
        done
        ;;
    *) exit 2 ;;
esac
