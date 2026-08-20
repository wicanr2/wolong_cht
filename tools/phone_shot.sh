#!/usr/bin/env bash
# 在桌面截一張**手機版**的畫面（960×540）。
#
#   tools/phone_shot.sh out.png [幀數]
#
# 手機 UI 的開發迴圈：同一份 `internal/ui/phone`，用桌面的 Xvfb 驗，
# 不必每次起 Android 模擬器（docs/mobile/android-plan.md §6）。
# 環境變數 WOLONG_SCENARIO／WOLONG_PLAYER／WOLONG_SEED 會傳進去。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${1:?用法: tools/phone_shot.sh <輸出.png> [幀數]}"
FRAME="${2:-60}"
mkdir -p "$(dirname "$OUT")"
OUT_ABS="$(cd "$(dirname "$OUT")" && pwd)/$(basename "$OUT")"

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$REPO_ROOT:/src" -v "$(dirname "$OUT_ABS"):/out" \
    -v wl-gomod:/gomod -v wl-gobuild:/gocache \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e GOCACHE=/gocache -e GOMODCACHE=/gomod \
    -e LIBGL_ALWAYS_SOFTWARE=1 \
    -e WOLONG_SHOT="/out/$(basename "$OUT_ABS")" -e WOLONG_SHOT_FRAME="$FRAME" \
    -e WOLONG_SCENARIO="${WOLONG_SCENARIO:-0}" -e WOLONG_PLAYER="${WOLONG_PLAYER:-0}" \
    -e WOLONG_SEED="${WOLONG_SEED:-7}" \
    -e WOLONG_SELECT="${WOLONG_SELECT:-}" -e WOLONG_ZOOM="${WOLONG_ZOOM:-}" \
    -e WOLONG_PAUSED="${WOLONG_PAUSED:-}" \
    -e WOLONG_SHEET="${WOLONG_SHEET:-}" -e WOLONG_TAB="${WOLONG_TAB:-}" \
    -e WOLONG_ADVISE="${WOLONG_ADVISE:-}" \
    -w /src "${WOLONG_GO_IMAGE:-demonwinter-go}" bash -c '
set -e
export PATH=/usr/local/go/bin:$PATH
go build -o /tmp/app ./cmd/wlandroid
Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
for i in $(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
DISPLAY=:99 timeout 120 /tmp/app
'
echo "$OUT"
