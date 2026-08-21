#!/usr/bin/env bash
# 在 headless X（Xvfb）底下跑 cmd/wlgame，把每一張畫出來的圖倒成 PNG。
#
#   tools/wlgame_frames.sh <輸出目錄> <張數> [wlgame 參數…]
#
# 推廣片的動態素材靠它（docs/spec/71）。與 `tools/shot.sh` 的差別是
# **不做 X 擷取**——程式自己寫 PNG，Xvfb 只是 Ebiten 開得起視窗的前提。
#
# ⚠ 不要在程式參數前加 `--`：Go 的 flag 套件遇到 `--` 就停止解析，
# 後面的旗標會被整批忽略——程式照樣跑、照樣錄，只是全用預設值。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go}"

OUT="${1:?用法: tools/wlgame_frames.sh <輸出目錄> <張數> [wlgame 參數…]}"
N="${2:?第二個參數是張數}"
shift 2

mkdir -p "$OUT"
OUT_ABS="$(cd "$OUT" && pwd)"
# 舊的圖要清掉：張數變少時殘留的尾巴會被 ffmpeg 一起吃進去，
# 而那看起來像「錄到了」不像「上一次的」。
rm -f "$OUT_ABS"/f*.png

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$REPO_ROOT:/src" -v "$OUT_ABS:/out" \
    -v wl-gomod:/gomod -v wl-gobuild:/gocache \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e GOCACHE=/gocache -e GOMODCACHE=/gomod \
    -e LIBGL_ALWAYS_SOFTWARE=1 -w /src \
    "$IMAGE" bash -c "
set -e
go build -o /tmp/app ./cmd/wlgame
Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
XVFB_PID=\$!
for i in \$(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
# 錄滿就自己結束；包 timeout 是因為參數打錯時 flag 套件會印用法直接退出，
# 而那種情況下等下去只是白等。
DISPLAY=:99 timeout 900 /tmp/app -frames-dir /out -frames $N $* 2>&1 | tail -20
kill -9 \$XVFB_PID 2>/dev/null || true
"

got=$(ls "$OUT_ABS"/f*.png 2>/dev/null | wc -l)
echo "錄到 $got 張 → $OUT_ABS"
[ "$got" -ge "$N" ] || { echo "⚠ 少於要求的 $N 張" >&2; exit 1; }
