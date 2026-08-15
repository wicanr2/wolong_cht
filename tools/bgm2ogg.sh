#!/usr/bin/env bash
# 把原版的音樂渲染成 ogg。
#
#   tools/bgm2ogg.sh                       # 全部（BGM.DAT 11 首 ＋ 三個單曲檔）
#   tools/bgm2ogg.sh BGM.DAT 0             # 只做一首
#   SECONDS_PER_SONG=120 tools/bgm2ogg.sh  # 改長度
#
# 兩段路：`cmd/wlaudio`（純 Go，OPL3 合成 → WAV）＋ docker ffmpeg（WAV → ogg）。
# **Go 這邊沒有 vorbis 編碼器**，所以第二段一定要出去；這也是為什麼
# 中介的 WAV 會留著——它是「合成對不對」與「編碼對不對」的分界。
#
# ⚠ 輸出是**原版衍生物**：`workplace/` 已 gitignore，
# 不進版控、不進發行包（`CLAUDE.md` §9）。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ORIG="$ROOT/workplace/orig/dosv"
OUT="$ROOT/workplace/audio"
FFMPEG_IMAGE="${FFMPEG_IMAGE:-linuxserver/ffmpeg}"
SECONDS_PER_SONG="${SECONDS_PER_SONG:-90}"
mkdir -p "$OUT"

render() { # <資料檔> <曲號>
  local dat="$1" song="$2" base
  base="$(basename "$dat" .DAT | tr 'A-Z' 'a-z')-$song"
  "$ROOT/tools/go.sh" run ./cmd/wlaudio \
    -bgm "workplace/orig/dosv/$dat" -song "$song" \
    -seconds "$SECONDS_PER_SONG" -out "workplace/audio/$base.wav"
  docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none -v "$OUT:/out" "$FFMPEG_IMAGE" \
    -hide_banner -loglevel error -y -i "/out/$base.wav" \
    -c:a libvorbis -q:a 5 "/out/$base.ogg"
  # 容器以 root 寫檔，換回呼叫者。
  docker run --rm --network none --entrypoint chown \
    -v "$OUT:/out" "$FFMPEG_IMAGE" -R "$(id -u):$(id -g)" /out
  printf '  → %s.ogg（%s）\n' "$base" "$(du -h "$OUT/$base.ogg" | cut -f1)"
}

if [ $# -eq 2 ]; then
  render "$1" "$2"
  exit 0
fi

for n in $(seq 0 10); do render BGM.DAT "$n"; done
for f in OPENBGM.DAT ENDBGM.DAT OVERBGM.DAT; do
  [ -f "$ORIG/$f" ] && render "$f" 0
done
