#!/usr/bin/env bash
# 在 docker 內跑 go 指令（[HARD] 本專案一律用 docker 編譯）。
#
#   tools/go.sh test ./...
#   tools/go.sh build ./cmd/wlview
#
# image 預設沿用 demonwinter-go（同一台機器上已存在、內容相同的 Go+Ebiten 環境）。
# 這是刻意的：在共用機器上多疊一份 1.9 GB 的相同 image 沒有意義。
# 要獨立的 image 就 `docker build -t wolong-go docker/go` 再設
# WOLONG_GO_IMAGE=wolong-go。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go}"
CACHE_VOL=wl-gomod
BUILD_VOL=wl-gobuild

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "[go.sh] 找不到 $IMAGE，改建 wolong-go …" >&2
    docker build -q -t wolong-go "$REPO_ROOT/docker/go" >&2
    IMAGE=wolong-go
fi

# 新建的 volume 屬 root，但下面以呼叫者 uid 執行 → 先把擁有者換過來。
if ! docker volume inspect "$CACHE_VOL" >/dev/null 2>&1 \
   || ! docker volume inspect "$BUILD_VOL" >/dev/null 2>&1; then
    docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
        -v "$CACHE_VOL:/gomod" -v "$BUILD_VOL:/gocache" "$IMAGE" \
        chown -R "$(id -u):$(id -g)" /gomod /gocache
fi

exec docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    -v "$REPO_ROOT:/src" \
    -v "$CACHE_VOL:/gomod" \
    -v "$BUILD_VOL:/gocache" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp \
    -e GOCACHE=/gocache \
    -e GOMODCACHE=/gomod \
    -w /src \
    "$IMAGE" go "$@"
