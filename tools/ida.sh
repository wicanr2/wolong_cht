#!/usr/bin/env bash
# IDA Pro 9.4 headless 包裝。image 來源 /home/anr2/ida_94_official/dist
#
#   tools/ida.sh batch <版本> <執行檔>      產 .i64 + .asm
#   tools/ida.sh raw   <版本> <idat 參數…>  直接下 idat 指令
#
# <版本> = dosv | pc98，對應 workplace/ida/<版本>/
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE=ida-pro-9.4-ver2
MODE="$1"; VER="$2"; shift 2
WORK="$ROOT/workplace/ida/$VER"
mkdir -p "$WORK"

run() {
  docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    -v "$WORK:/work" -v "$ROOT/tools:/work/tools:ro" -w /work \
    "$IMAGE" "$@"
}

case "$MODE" in
  batch)
    BIN="$1"
    cp -f "$ROOT/workplace/orig/$VER/$BIN" "$WORK/$BIN"
    chmod u+w "$WORK/$BIN"
    sha256sum "$WORK/$BIN"
    run idat -A -B "$BIN"
    ;;
  raw) run "$@" ;;
  *) echo "用法: ida.sh {batch|raw} {dosv|pc98} …" >&2; exit 2 ;;
esac
