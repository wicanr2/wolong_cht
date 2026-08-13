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
        --network none --memory 128m --cpus 0.5 --pids-limit 64 \
        -v "$CACHE_VOL:/gomod" -v "$BUILD_VOL:/gocache" "$IMAGE" \
        chown -R "$(id -u):$(id -g)" /gomod /gocache
fi

# ⚠ **交叉編譯的環境變數要明文傳進容器。**
# docker run 不會繼承呼叫端的環境；少了這幾行，`GOOS=windows tools/go.sh build`
# 會安靜地建出**本機平台**的執行檔——三個平台建出三個一模一樣的檔案
# （同一個 BuildID），而且每一個都 exit 0。
# 踩過：拿這個當「Ebiten 可以交叉編譯」的證據，證的其實是本機建置成功。
CROSS_ENV=()
for v in GOOS GOARCH CGO_ENABLED GOARM GOAMD64; do
    if [ -n "${!v:-}" ]; then CROSS_ENV+=(-e "$v=${!v}"); fi
done

exec docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    "${CROSS_ENV[@]+"${CROSS_ENV[@]}"}" \
    -v "$REPO_ROOT:/src" \
    -v "$CACHE_VOL:/gomod" \
    -v "$BUILD_VOL:/gocache" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp \
    -e GOCACHE=/gocache \
    -e GOMODCACHE=/gomod \
    -w /src \
    "$IMAGE" go "$@"
