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
GO_CACHE_ROOT="${WOLONG_GO_CACHE_ROOT:-}"
GO_NETWORK="${WOLONG_GO_NETWORK:-none}"
GO_TIMEOUT="${WOLONG_GO_TIMEOUT:-1200}"
GO_CONTAINER="wolong-go-$$-$RANDOM"
case "$GO_NETWORK" in none|bridge) ;; *) echo "[go.sh] network 必須是 none 或 bridge" >&2;exit 2;; esac
[[ "$GO_TIMEOUT" =~ ^[1-9][0-9]*$ ]] || exit 2

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "[go.sh] 找不到 $IMAGE，改建 wolong-go …" >&2
    docker build -q -t wolong-go "$REPO_ROOT/docker/go" >&2
    IMAGE=wolong-go
fi

# 新建的 volume 屬 root，但下面以呼叫者 uid 執行 → 先把擁有者換過來。
if [[ -z "$GO_CACHE_ROOT" ]] && { ! docker volume inspect "$CACHE_VOL" >/dev/null 2>&1 \
   || ! docker volume inspect "$BUILD_VOL" >/dev/null 2>&1; }; then
    docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
        --network none --memory 128m --cpus 0.5 --pids-limit 64 \
        -v "$CACHE_VOL:/gomod" -v "$BUILD_VOL:/gocache" "$IMAGE" \
        chown -R "$(id -u):$(id -g)" /gomod /gocache
fi
if [[ -n "$GO_CACHE_ROOT" ]]; then
    test -d "$GO_CACHE_ROOT/gomod";test -d "$GO_CACHE_ROOT/gobuild"
    test -O "$GO_CACHE_ROOT/gomod";test -O "$GO_CACHE_ROOT/gobuild"
    CACHE_MOUNTS=(--mount "type=bind,src=$GO_CACHE_ROOT/gomod,dst=/gomod" --mount "type=bind,src=$GO_CACHE_ROOT/gobuild,dst=/gocache")
else
    CACHE_MOUNTS=(-v "$CACHE_VOL:/gomod" -v "$BUILD_VOL:/gocache")
fi
ORIGINAL_MOUNT=()
if [[ -d "$REPO_ROOT/workplace/orig" ]]; then
    ORIGINAL_MOUNT=(--mount "type=bind,src=$REPO_ROOT/workplace/orig,dst=/src/workplace/orig,readonly")
fi

# ⚠ **交叉編譯的環境變數要明文傳進容器。**
# docker run 不會繼承呼叫端的環境；少了這幾行，`GOOS=windows tools/go.sh build`
# 會安靜地建出**本機平台**的執行檔——三個平台建出三個一模一樣的檔案
# （同一個 BuildID），而且每一個都 exit 0。
# 踩過：拿這個當「Ebiten 可以交叉編譯」的證據，證的其實是本機建置成功。
CROSS_ENV=()
# WOLONG_DUMP_DIR 讓測試把圖寫出來供肉眼複驗（容器內路徑）。
for v in GOOS GOARCH CGO_ENABLED GOARM GOAMD64 WOLONG_DUMP_DIR GOPROXY GOSUMDB GOTOOLCHAIN GOFLAGS; do
    if [ -n "${!v:-}" ]; then CROSS_ENV+=(-e "$v=${!v}"); fi
done

# ⚠ **Ebiten 在 init 期就要求顯示器**，所以 `cmd/wlgame` 與 `internal/ui/textdraw`
# 沒有 X 就會 panic（`glfw: The GLFW library is not initialized`）。
# 快取命中時看不出來——`go test` 直接回 `(cached)`，冷啟動才炸。
# 因此一律先起 Xvfb 再跑，不要靠快取掩蓋。
trap 'docker rm -f "$GO_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 "$GO_TIMEOUT" docker run --rm --init --name "$GO_CONTAINER" --log-opt max-size=10m --log-opt max-file=3 \
    --network "$GO_NETWORK" --memory 2g --cpus 2 --pids-limit 256 \
    --label wolong.project=dragon --label "wolong.task=${WOLONG_GO_TASK:-go}" \
    "${CROSS_ENV[@]+"${CROSS_ENV[@]}"}" \
    -v "$REPO_ROOT:/src" \
    "${CACHE_MOUNTS[@]}" \
    "${ORIGINAL_MOUNT[@]}" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp \
    -e GOCACHE=/gocache \
    -e GOMODCACHE=/gomod \
    -e DISPLAY=:99 \
    -w /src \
    "$IMAGE" sh -c 'xvfb_pid=""
                    trap '\''if [ -n "$xvfb_pid" ]; then kill "$xvfb_pid" 2>/dev/null || true;wait "$xvfb_pid" 2>/dev/null || true;fi'\'' EXIT
                    if command -v Xvfb >/dev/null 2>&1 && ! [ -e /tmp/.X11-unix/X99 ]; then
                      Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
                      xvfb_pid=$!
                      for _ in 1 2 3 4 5 6 7 8 9 10; do [ -e /tmp/.X11-unix/X99 ] && break; sleep 0.3; done
                    fi
                    go "$@"' _ "$@"
