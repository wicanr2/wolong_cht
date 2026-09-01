#!/usr/bin/env bash
# 在一次性容器裡跑 remake 並照時間軸操作、截圖（主機端入口）。
#
#   tools/wlgame_capture.sh <輸出目錄名> "<timeline>" [wlgame 參數…]
#
# 例（驗一顆已經編好的 AppImage）：
#   WOLONG_WLGAME_BIN=dist-all/packages/wolong-remake-linux-amd64-20260830.AppImage \
#   tools/wlgame_capture.sh flow-lubu "wait:3;shot:00-title;key:Return"
#
# 這一支與 `tools/dosv_capture.sh` 是一對：那一支跑原版（滑鼠時間軸），
# 這一支跑 remake（鍵盤時間軸）。容器內跑的是 tools/wlgame_live_capture.sh。
# 輸出落在 workplace/promo-live/<輸出目錄名>/，尺寸 640×400，可直接送
# tools/parity_diff.py 與原版並排。
#
# ⚠ **AppImage 要 `--appimage-extract-and-run`**：容器裡沒有 FUSE。
#    這一支會先把它解到 /tmp 再跑，所以 `WOLONG_WLGAME_BIN` 給 AppImage 或
#    裸執行檔都可以。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_GO_IMAGE:-demonwinter-go}"
NAME_PREFIX="${WOLONG_WLGAME_CONTAINER_PREFIX:-wolong-wlgame-capture}"

if [ $# -lt 1 ]; then
    echo "用法: tools/wlgame_capture.sh <輸出目錄名> [timeline] [wlgame 參數…]" >&2
    exit 2
fi
out_name=$1; shift || true
timeline=${1:-"wait:3;shot:title"}; shift || true

orig="${WOLONG_ORIG_DIR:-$REPO_ROOT/workplace/orig/dosv}"
font="${WOLONG_FONT_DIR:-$REPO_ROOT/workplace/eten}"
out="$REPO_ROOT/workplace/promo-live/$out_name"
bin="${WOLONG_WLGAME_BIN:-}"

[ -d "$orig" ] || { echo "[wlgame] 找不到原版目錄 $orig" >&2; exit 1; }
[ -d "$font" ] || { echo "[wlgame] 找不到字型目錄 $font" >&2; exit 1; }
mkdir -p "$out"

docker image inspect "$IMAGE" >/dev/null 2>&1 || {
    echo "[wlgame] 找不到映像 $IMAGE" >&2; exit 1; }

bin_mount=()
bin_env=()
if [ -n "$bin" ]; then
    bin_abs="$(cd "$(dirname "$bin")" && pwd)/$(basename "$bin")"
    [ -x "$bin_abs" ] || { echo "[wlgame] 執行檔不存在或不可執行：$bin_abs" >&2; exit 1; }
    bin_mount=(-v "$bin_abs:/pkg/app:ro")
    bin_env=(-e WOLONG_WLGAME_PREBUILT=/pkg/app)
    echo "[wlgame] 執行檔 $bin_abs"
    echo "  sha256=$(sha256sum "$bin_abs" | cut -d' ' -f1)"
else
    echo "[wlgame] 沒給 WOLONG_WLGAME_BIN，改在容器裡 go build"
fi

echo "[wlgame] 輸出 $out"
echo "[wlgame] timeline=$timeline"
exec docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --name "$NAME_PREFIX-$$" \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$REPO_ROOT:/src:ro" \
    -v "$orig:/orig:ro" \
    -v "$font:/font:ro" \
    -v "$out:/out" \
    "${bin_mount[@]}" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e LIBGL_ALWAYS_SOFTWARE=1 \
    -e "WOLONG_WLGAME_TIMELINE=$timeline" \
    -e "WOLONG_WLGAME_EXTRA_ARGS=$*" \
    "${bin_env[@]}" \
    $(env | sed -n 's/^\(WOLONG_WLGAME_[A-Z0-9_]*\)=.*/-e \1/p') \
    -w /tmp "$IMAGE" bash /src/tools/wlgame_live_capture.sh
