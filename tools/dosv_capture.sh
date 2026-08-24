#!/usr/bin/env bash
# 在一次性 DOSBox-X 容器裡跑松崗 DOS/V 原版並截圖（主機端入口）。
#
#   tools/dosv_capture.sh <輸出目錄名> "<timeline>"
#
# 例：
#   tools/dosv_capture.sh parity-main \
#     "wait:125;click:320,215;wait:3;click:300,190;wait:4;click:450,154"
#
# 輸出落在 workplace/promo-live/<輸出目錄名>/。容器內跑的是
# tools/dosv_live_capture.sh，兩支的分工：那一支在容器裡開 Xvfb 與
# DOSBox-X，這一支負責掛載、資源上限與參數。
#
# ⚠ **原版目錄唯讀掛載**：容器自己複製一份到 /tmp 再跑，
#    guest 寫的 SAVE.DAT 不會回到 workplace/orig/。
#
# ## `[HARD]` 滑鼠座標要乘 1.2
#
# DOSBox-X 的視窗是 **640×480**，遊戲畫面 640×400 放在 y 偏移 40
# （`tools/parity_crop.py`）。但 **INT 33 的對映吃的是整個視窗**：
#
#     遊戲座標 y  =  視窗座標 y × 400/480
#     視窗座標 y  =  遊戲座標 y × 1.2      ← 要點哪裡就這樣算
#
# 也就是說**黑邊也算在裡面**，不是「減掉 40」。量到的證據：
# 送 (200,165) 之後反白的是遊戲座標 y 136..151 那一列（165×5/6 ＝ 137，
# 減 40 會得到 125，落在上一列）；送 (416,56) 之後 guest 游標畫在
# 遊戲座標 y≈46（56×5/6 ＝ 46.7）。
#
# 少乘這一下的症狀是**「點了沒反應」**——按鈕之間有 8 px 的空隙，
# 座標差個十幾 px 就正好落進空隙，畫面完全不動，跟事件沒送到一模一樣。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_DOSBOXX_IMAGE:-wolong-dosboxx}"
NAME_PREFIX="${WOLONG_DOSV_CONTAINER_PREFIX:-wolong-dosv-capture}"

if [ $# -lt 1 ]; then
    echo "用法: tools/dosv_capture.sh <輸出目錄名> [timeline]" >&2
    exit 2
fi
out_name=$1
timeline=${2:-"wait:1;shot:after-confirm"}
boot_mode=${WOLONG_DOSV_BOOT_MODE:-start}

orig="$REPO_ROOT/workplace/orig/dosv"
out="$REPO_ROOT/workplace/promo-live/$out_name"
if [ ! -d "$orig" ]; then
    echo "[dosv] 找不到原版目錄 $orig" >&2
    exit 1
fi
mkdir -p "$out"

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "[dosv] 找不到映像 $IMAGE；先 docker build -f docker/dosboxx/Dockerfile" >&2
    exit 1
fi

{
    echo "source=松崗 DOS/V 唯讀實機畫面"
    echo "machine=vgaonly"
    echo "core=normal"
    echo "cputype=486"
    echo "cycles=fixed 20000"
    echo "mouse_emulation=integration"
    echo "boot_mode=$boot_mode"
    echo "timeline=$timeline"
} > "$out/capture-metadata.txt"

echo "[dosv] 輸出 $out"
echo "[dosv] timeline=$timeline"
exec docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --name "$NAME_PREFIX-$$" \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$orig:/orig:ro" \
    -v "$out:/out" \
    -v "$REPO_ROOT/tools/dosv_live_capture.sh:/capture.sh:ro" \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp \
    -e "WOLONG_DOSV_TIMELINE=$timeline" \
    -e "WOLONG_DOSV_BOOT_MODE=$boot_mode" \
    -e "WOLONG_DOSV_SEED_SAVE=${WOLONG_DOSV_SEED_SAVE:-}" \
    -w /tmp \
    "$IMAGE" bash /capture.sh
