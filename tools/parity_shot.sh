#!/usr/bin/env bash
# 截一張 **640×400 原尺寸** 的 remake 畫面，給逐區對拍用（`docs/spec/90`）。
#
#   tools/parity_shot.sh out.png [wlgame 參數…]
#   tools/parity_shot.sh workplace/parity/remake-main.png \
#       -direct -scenario 0 -player 0 -seed 7 -open-window -2
#
# 與 `tools/shot.sh` 的差別：那一支是 `import -window root`，抓到的是
# **1600×900 的 Xvfb 桌面**（遊戲以 2× 放在中間）。那個尺寸給文件當插圖沒問題，
# 拿來跟原版的 640×400 比就得先縮放——而縮放會把「尺寸不對」這個
# 最重要的線索洗掉（`tools/parity_diff.py` 因此對尺寸不同直接報錯）。
#
# 這一支改用 `wlgame` 自己的 `-shot`：它寫的是 `Layout()` 回傳的
# **邏輯畫面**，也就是 640×400，一個像素都沒有經過縮放。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go}"
CACHE_VOL=wl-gomod
BUILD_VOL=wl-gobuild

OUT="${1:?用法: tools/parity_shot.sh <輸出.png> [wlgame 參數…]}"
shift || true
mkdir -p "$(dirname "$OUT")"
OUT_ABS="$(cd "$(dirname "$OUT")" && pwd)/$(basename "$OUT")"
OUT_BASE="$(basename "$OUT_ABS")"

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$REPO_ROOT:/src" \
    -v "$(dirname "$OUT_ABS"):/out" \
    -v "$CACHE_VOL:/gomod" -v "$BUILD_VOL:/gocache" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e GOCACHE=/gocache -e GOMODCACHE=/gomod \
    -e LIBGL_ALWAYS_SOFTWARE=1 \
    -w /src \
    "$IMAGE" bash -c "
set -e
export PATH=/usr/local/go/bin:\$PATH
go build -o /tmp/app ./cmd/wlgame
Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
XVFB_PID=\$!
for i in \$(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
# 指標停到左上角：Xvfb 預設把指標放在桌面中央，映進遊戲畫面正中，
# remake 的游標會畫在那裡，對拍時吃掉一塊（playtest/42 §3）。
DISPLAY=:99 xdotool mousemove 0 0 2>/dev/null || true
DISPLAY=:99 timeout 120 /tmp/app $* -shot /out/$OUT_BASE
kill -9 \$XVFB_PID 2>/dev/null || true
"
echo "$OUT"
