#!/usr/bin/env bash
# 正常視窗驗收；候選可指定第三參數替換程式，正式包不帶第三參數。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${1:?請指定版號}"
CASE="${2:?請指定本輪新輸出目錄名稱}"
CANDIDATE="${3:-}"
MODE="${4:-polish}"
[[ "$MODE" == polish || "$MODE" == scroll || "$MODE" == march ]] || exit 2
[[ "$VERSION" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]] || exit 2
[[ "$CASE" =~ ^[a-z0-9-]+$ ]] || exit 2
OUT="$ROOT/workplace/desktop-polish-20260908/$CASE"
[[ ! -e "$OUT" ]] || { echo '輸出已存在，請指定新案例'; exit 2; }
mkdir -p "$OUT"
stat -c '%u:%g %n' "$OUT"
exec timeout 900 docker run --rm --name "wolong-polish-$CASE" --network none --memory 3g --cpus 3 --pids-limit 256 \
  -u "$(id -u):$(id -g)" -v "$ROOT:/src:ro" -v "$OUT:/out" \
  -e WOLONG_VERSION="$VERSION" -e WOLONG_CANDIDATE="$CANDIDATE" -e WOLONG_VERIFY_MODE="$MODE" \
  -e WOLONG_VERIFY_SCOPE="${WOLONG_VERIFY_SCOPE:-all}" \
  -e WOLONG_FOCUS_WAIT="${WOLONG_FOCUS_WAIT:-2}" \
  -e WOLONG_FOCUS_REPEATS="${WOLONG_FOCUS_REPEATS:-1}" \
  -e WOLONG_FOCUS_FREEZE="${WOLONG_FOCUS_FREEZE:-0}" \
  --entrypoint bash demonwinter-go:latest -c '
set -eu
cd /tmp
/src/dist-all/$WOLONG_VERSION/full/wolong-remake-linux-amd64-$WOLONG_VERSION.AppImage --appimage-extract >/dev/null
if [ -n "$WOLONG_CANDIDATE" ]; then cp "/src/$WOLONG_CANDIDATE" /tmp/squashfs-root/usr/bin/wlgame; fi
sha256sum /tmp/squashfs-root/usr/bin/wlgame >/out/runtime.sha256
root=root-saveb
[ "$WOLONG_VERIFY_MODE" == polish ] || root=root-noclouds
cp /src/workplace/dosgolem/$root/SAVE.DAT /out/SAVE.DAT
mkdir -p /out/home
export HOME=/out/home XDG_CONFIG_HOME=/out/home/.config DISPLAY=:99 LIBGL_ALWAYS_SOFTWARE=1
printf "pcm.!default { type null }\n" >$HOME/.asoundrc
Xvfb :99 -screen 0 1280x800x24 >/out/xvfb.log 2>&1 & xp=$!
trap '\''kill "$xp" 2>/dev/null || true'\'' EXIT
sleep 1
timeout 850 python3 /src/tools/verify_desktop_$WOLONG_VERIFY_MODE.py
'
