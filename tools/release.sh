#!/usr/bin/env bash
# 發行打包（M8）。跨平台建置 → deny-list 掃描 → 壓成單一目錄。
#
#   tools/release.sh              建全部平台
#   tools/release.sh linux/amd64  只建一個
#
# [HARD] 建置、環境查詢與 deny-list 一律走 docker（CLAUDE.md §10），所以這支只是
# tools/go.sh／tools/py.sh 的迴圈 ＋ 交叉編譯環境變數。
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

# ── 哪些東西能交叉編譯（實測，不是推測）──
#
# **遊戲邏輯本身是 100% 純 Go**：`internal/rules/*`、`internal/state`、
# `internal/assets/*` 完全不碰 cgo，`wlsim`／`wlshot` 三個平台
# 用 `CGO_ENABLED=0` 都建得出來（Mach-O arm64／ELF aarch64／PE32+ 都驗過）。
#
# 需要 cgo 的只有**開視窗那一層**，而且**只有部分平台**：
#
#   windows  ✅ 純 Go 就能建 —— Ebiten 走 syscall／purego 動態載 DLL
#   linux    ❌ 要 cgo —— OpenGL 驅動綁 GLFW（C 函式庫），
#                CGO_ENABLED=0 會停在 `undefined: glfw.Window`
#   darwin   ❌ 同上（Cocoa／Metal）
#
# 所以矩陣是「純邏輯工具全平台 ＋ 遊戲本體只交叉到 windows」，
# linux／mac 的本體要在目標平台自己建。
CROSS_OK=(./cmd/wlsim ./cmd/wlshot)              # 不依賴 Ebiten
CROSS_WINDOWS_OK=(./cmd/wlgame ./cmd/wlview)     # 依賴 Ebiten，但 windows 純 Go
NATIVE_ONLY=(./cmd/wlgame ./cmd/wlview)

tools/py.sh tools/release_fs.py prepare

echo "── 版本 $VERSION ──"
for p in "${PLATFORMS[@]}"; do
    os="${p%%/*}"; arch="${p##*/}"
    out="$DIST/$os-$arch"
    tools/py.sh tools/release_fs.py mkdir "$out"
    ext=""; [ "$os" = windows ] && ext=.exe

    pkgs=("${CROSS_OK[@]}")
    if [ "$os" = windows ]; then pkgs+=("${CROSS_WINDOWS_OK[@]}"); fi
    for pkg in "${pkgs[@]}"; do
        name="$(basename "$pkg")"
        echo "  $os/$arch  $name"
        GOOS=$os GOARCH=$arch CGO_ENABLED=0 \
            tools/go.sh build -trimpath \
                -ldflags "-s -w -X main.version=$VERSION" \
                -o "$out/$name$ext" "$pkg"
    done
done

# 本機平台額外建遊戲本體。linux／mac 的 Ebiten 要 cgo，只能在目標平台建；
# windows 上面那圈已經建好了。
echo "── 本機平台（linux／mac 的 Ebiten 要 cgo）──"
host_os="$(tools/go.sh env GOOS)"
host_arch="$(tools/go.sh env GOARCH)"
native="$DIST/${host_os}-${host_arch}"
tools/py.sh tools/release_fs.py mkdir "$native"
for pkg in "${NATIVE_ONLY[@]}"; do
    name="$(basename "$pkg")"
    echo "  本機  $name"
    tools/go.sh build -trimpath \
        -ldflags "-s -w -X main.version=$VERSION" \
        -o "$native/$name" "$pkg"
done

tools/py.sh tools/release_fs.py finalize

echo "── deny-list（發行閘）──"
tools/py.sh tools/denylist.py "$DIST"

echo
echo "產出在 $DIST/"
tools/py.sh tools/release_fs.py report
