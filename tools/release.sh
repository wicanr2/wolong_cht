#!/usr/bin/env bash
# 發行打包（M8）。跨平台建置 → deny-list 掃描 → 壓成單一目錄。
#
#   tools/release.sh              建全部平台
#   tools/release.sh linux/amd64  只建一個
#
# [HARD] 建置一律走 docker（CLAUDE.md §10），所以這支只是 tools/go.sh 的
# 迴圈 ＋ 交叉編譯的環境變數。
#
# ⚠ **deny-list 是發行閘，不是收尾檢查。** 它擋不下來就不出包——
# 所以掃描排在打包之前，而且掃的是 dist/ 本身而不是 repo。
# repo 乾淨不代表產出乾淨：`go:embed` 可以把任何東西烤進執行檔。
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
DIST=dist

PLATFORMS=(
    linux/amd64 linux/arm64
    windows/amd64
    darwin/amd64 darwin/arm64
)
if [ $# -gt 0 ]; then PLATFORMS=("$@"); fi

# Ebiten 在 linux 上要 cgo（X11／OpenGL），交叉編譯拿不到那些標頭。
# **所以能交叉出去的只有純邏輯的工具**，遊戲本體要在目標平台上自己建。
# 這件事寫在這裡而不是留給使用者踩：
CROSS_OK=(./cmd/wlsim)          # 不依賴 Ebiten
NATIVE_ONLY=(./cmd/wlgame ./cmd/wlview ./cmd/wlshot)

rm -rf "$DIST"; mkdir -p "$DIST"

echo "── 版本 $VERSION ──"
for p in "${PLATFORMS[@]}"; do
    os="${p%%/*}"; arch="${p##*/}"
    out="$DIST/$os-$arch"; mkdir -p "$out"
    ext=""; [ "$os" = windows ] && ext=.exe

    for pkg in "${CROSS_OK[@]}"; do
        name="$(basename "$pkg")"
        echo "  $os/$arch  $name"
        GOOS=$os GOARCH=$arch CGO_ENABLED=0 \
            tools/go.sh build -trimpath \
                -ldflags "-s -w -X main.version=$VERSION" \
                -o "$out/$name$ext" "$pkg"
    done
done

# 本機平台額外建遊戲本體（需要 cgo）。
echo "── 本機平台（Ebiten 需要 cgo）──"
native="$DIST/$(go env GOOS 2>/dev/null || echo linux)-$(go env GOARCH 2>/dev/null || echo amd64)"
mkdir -p "$native"
for pkg in "${NATIVE_ONLY[@]}"; do
    name="$(basename "$pkg")"
    echo "  本機  $name"
    tools/go.sh build -trimpath \
        -ldflags "-s -w -X main.version=$VERSION" \
        -o "$native/$name" "$pkg"
done

cp README.md "$DIST/" 2>/dev/null || true
cp -r translations "$DIST/" 2>/dev/null || true

echo "── deny-list（發行閘）──"
python3 tools/denylist.py "$DIST"

echo
echo "產出在 $DIST/"
du -sh "$DIST"/* 2>/dev/null || true
