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

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "[py.sh] 找不到 $IMAGE；請先建立或修復既有工具映像" >&2
    exit 1
fi

exec docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 768m --cpus 1 --pids-limit 96 \
    -v "$REPO_ROOT:/src" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp \
    -w /src \
    "$IMAGE" python3 "$@"
