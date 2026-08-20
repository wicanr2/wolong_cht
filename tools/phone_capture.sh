#!/usr/bin/env bash
# 把手機版的畫面逐幀輸出成 PNG，給推廣片剪接用。
#
#   tools/phone_capture.sh [輸出目錄] [幀數]
#
# ⭐ **不錄螢幕**。X11 擷取要跟畫面搶時間，錄出來的幀率不穩還會抖；
# 逐幀輸出是確定性的——同一個 seed 跑兩次得到同一批圖。
# 操作也不送 xdotool 點擊：時間軸在 `mobile/wolong/demo.go`，
# 與畫面更新同一條時間線，不會「點下去時畫面還沒到那一格」。
#
# 輸出：`fNNNNN.png`（一張 ＝ 一個影片幀，30 fps）＋ `marks.txt`
#（`標籤 圖號`，剪接靠它切段）。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${1:-$REPO_ROOT/workplace/promo/android-frames}"
FRAMES="${2:-1200}"
mkdir -p "$OUT"
rm -f "$OUT"/f*.png "$OUT"/marks.txt

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 4g --cpus 2 --pids-limit 256 \
    -v "$REPO_ROOT:/src" -v "$OUT:/out" \
    -v wl-gomod:/gomod -v wl-gobuild:/gocache \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e GOCACHE=/gocache -e GOMODCACHE=/gomod \
    -e LIBGL_ALWAYS_SOFTWARE=1 \
    -e WOLONG_FRAMES_DIR=/out -e WOLONG_FRAMES_N="$FRAMES" \
    -e WOLONG_SCENARIO="${WOLONG_SCENARIO:-0}" -e WOLONG_PLAYER="${WOLONG_PLAYER:-0}" \
    -e WOLONG_SEED="${WOLONG_SEED:-7}" -e WOLONG_SPEED="${WOLONG_SPEED:-1}" \
    -w /src "${WOLONG_GO_IMAGE:-demonwinter-go}" bash -c '
set -e
export PATH=/usr/local/go/bin:$PATH
go build -o /tmp/app ./cmd/wlandroid
Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
for i in $(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
DISPLAY=:99 timeout 900 /tmp/app
'
echo "幀數：$(ls "$OUT"/f*.png | wc -l)"
cat "$OUT/marks.txt"
