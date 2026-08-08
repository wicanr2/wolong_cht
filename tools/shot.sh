#!/usr/bin/env bash
# 在 headless X（Xvfb）底下跑 cmd/wlview 並截一張圖。
#
#   tools/shot.sh out.png [KEYS=Right,Right,Down] [程式參數…]
#   WOLONG_SHOT_CMD=wlgame tools/shot.sh out.png …    截別支程式
#
# 用途是**驗收呈現層**。素材本身的解碼驗收用 cmd/wlshot（不需要 X），
# 這支是為了證明「Ebiten 這一層真的跑得起來、版面沒跑掉」，
# 不接受「編譯過了 / 測試綠」當作畫面正確的證據。
#
# ⚠ 不要在程式參數前加 `--`：Go 的 flag 套件遇到 `--` 就停止解析，
# 後面的旗標會被整批忽略 —— 程式照樣跑、照樣截圖，只是全用預設值。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go}"
CMD="${WOLONG_SHOT_CMD:-wlview}"
CACHE_VOL=wl-gomod
BUILD_VOL=wl-gobuild

OUT="${1:?用法: tools/shot.sh <輸出.png> [KEYS=按鍵序列] [程式參數…]}"
shift || true

KEYS=""
if [[ "${1:-}" == KEYS=* ]]; then KEYS="${1#KEYS=}"; shift; fi
if [[ "${1:-}" == "--" ]]; then shift; fi

mkdir -p "$(dirname "$OUT")"
OUT_ABS="$(cd "$(dirname "$OUT")" && pwd)/$(basename "$OUT")"

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
go build -o /tmp/app ./cmd/$CMD
Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
XVFB_PID=\$!
for i in \$(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
DISPLAY=:99 /tmp/app $* >/tmp/app.log 2>&1 &
APP_PID=\$!
sleep 3
# **一定要包 timeout。** xdotool search --sync 在視窗永遠不出現時會無限等，
# 而視窗不出現最常見的原因是程式參數打錯（flag 套件印用法就退出）。
WID=\$(timeout 20 env DISPLAY=:99 xdotool search --sync --onlyvisible --name '臥龍傳' 2>/dev/null | head -1) || true
if [ -z \"\$WID\" ]; then
    echo '!!! 視窗沒出現 —— 看下面的 app log'
    echo '--- app log ---'; cat /tmp/app.log || true
    kill -9 \$APP_PID \$XVFB_PID 2>/dev/null || true
    exit 1
fi
DISPLAY=:99 xdotool windowactivate --sync \$WID 2>/dev/null || true
sleep 1
if [ -n '$KEYS' ]; then
    for k in \$(echo '$KEYS' | tr ',' ' '); do
        DISPLAY=:99 xdotool key --window \$WID --clearmodifiers \$k
        sleep 0.25
    done
    sleep 0.6
fi
DISPLAY=:99 import -window root /out/$(basename "$OUT_ABS")
# **成功時也要把程式 log 印出來。**
# 先前只在「視窗沒出現」那條路徑印，於是啟動時的警告（載不到字型、
# 推不出道路圖…）全部被吞掉——截圖看起來正常，問題在 log 裡沒人看見。
echo '--- app log ---'; cat /tmp/app.log || true
kill -9 \$APP_PID \$XVFB_PID 2>/dev/null || true
"
echo "$OUT"
