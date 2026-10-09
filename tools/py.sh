#!/usr/bin/env bash
# 在既有 Docker Python runtime 內執行專案 Python 工具。
#
#   tools/py.sh tools/index.py generate
#   tools/py.sh tools/denylist.py --selftest
#
# Python 工具不依賴第三方套件；沿用 demonwinter-go 內的固定 Python，
# 避免把主機 Python、虛擬環境或未鎖版 library 帶進專案。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go}"
PY_CONTAINER="wolong-py-$$-$RANDOM"
PY_TIMEOUT="${WOLONG_PY_TIMEOUT:-600}"
[[ "$PY_TIMEOUT" =~ ^[1-9][0-9]*$ ]] || exit 2
ORIGINAL_MOUNT=()
if [[ -d "$REPO_ROOT/workplace/orig" ]]; then
    ORIGINAL_MOUNT=(--mount "type=bind,src=$REPO_ROOT/workplace/orig,dst=/src/workplace/orig,readonly")
fi

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "[py.sh] 找不到 $IMAGE；請先建立或修復既有工具映像" >&2
    exit 1
fi

trap 'docker rm -f "$PY_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 "$PY_TIMEOUT" docker run --rm --init --name "$PY_CONTAINER" --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 768m --cpus 1 --pids-limit 96 \
    --label wolong.project=dragon --label "wolong.task=${WOLONG_GO_TASK:-python}" \
    -v "$REPO_ROOT:/src" \
    "${ORIGINAL_MOUNT[@]}" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e TZ="${TZ:-$(cat /etc/timezone 2>/dev/null || echo UTC)}" \
    $(env | sed -n 's/^\(WOLONG_[A-Z0-9_]*\)=.*/-e \1/p') \
    -w /src \
    "$IMAGE" python3 "$@"
