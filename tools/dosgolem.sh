#!/usr/bin/env bash
# 用 dosgolem 跑原版當視覺 oracle——**不經過 DOSBox、不經過 X**。
#
#   tools/dosgolem.sh out_dir "wait;shot:z1;click:320,215;wait;shot:z2"
#
# 步驟語法與 `tools/dosv_live_capture.sh` 的時間軸同形，所以舊腳本搬得過來：
#
#   wait          畫面畫完並停住（不是「睡幾秒」）
#   steps:N       再跑 N 道指令
#   click:X,Y     點左鍵
#   move:X,Y      只移動游標
#   until:Y/M/D   跑到遊戲日期到某一天  ← 即時制的取樣點寫成日期
#   clock         印出目前的遊戲日期
#   shot:NAME     存一張 640×400 的 PNG（已經裁好，不必再跑 parity_crop.py）
#
# ⚠ **y 座標是遊戲座標，不是 DOSBox-X 的視窗座標。** 舊腳本要加 `-y dosbox`，
# 換算是 `遊戲 y ＝ 視窗 y × 399 ÷ 479`（分母 479 不是 480，
# 見 dosgolem 的 `apps/wolong/wolong.go`）。
#
# 規格：docs/spec/131-dosgolem-oracle.md
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOLEM="${WOLONG_DOSGOLEM:-$HOME/cht/dosgolem-wolong}"

if [[ ! -d "$GOLEM" ]]; then
    echo "找不到 dosgolem：$GOLEM" >&2
    echo "設 WOLONG_DOSGOLEM 指到工作副本，或 git clone https://github.com/wicanr2/dosgolem.git" >&2
    exit 1
fi

OUT="${1:?用法: tools/dosgolem.sh <輸出目錄> <時間軸> [dosbox]}"
TIMELINE="${2:?用法: tools/dosgolem.sh <輸出目錄> <時間軸> [dosbox]}"
YMODE="${3:-game}"

mkdir -p "$OUT"
OUT_ABS="$(cd "$OUT" && pwd)"

ARGS=(-exe /orig/orig/dosv/KI.EXE -root /orig/orig/dosv -dir /out -script "$TIMELINE")
[[ "$YMODE" == "dosbox" ]] && ARGS+=(-dosbox-y)

# 原版素材唯讀掛載；輸出目錄可寫。**本專案與 dosgolem 都不含原版檔案。**
DOSGOLEM_ORIG="$REPO_ROOT/workplace" \
DOSGOLEM_TIMEOUT="${WOLONG_DOSGOLEM_TIMEOUT:-30m}" \
DOSGOLEM_EXTRA_MOUNT="$OUT_ABS:/out" \
    "$GOLEM/tools/go.sh" run ./apps/wolong/cmd/shot "${ARGS[@]}"
