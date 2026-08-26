#!/usr/bin/env bash
# 把**原版配樂**接成推廣片的 72 秒音軌。
#
#   tools/promo_score_original.sh [輸出.wav]
#
# ⭐ 聲音來源是 remake 自己的 OPL3 合成（`tools/bgm2ogg.sh` 從使用者自備的
# `BGM.DAT` 算出來的，docs/spec/29）——**不是從原版錄音取樣**。
# 但它演奏的是原版的曲子，所以這條音軌是**原版衍生物**：
# 只用在推廣片，不進任何遊戲包，也不宣稱是本專案的創作。
#
# 曲目照場景挑（docs/re/58）：
#   0–18.5s  大地圖（曲 2 ＝ 春）
#   18.5–41s 戰場（曲 9 ＝ 平原野戰）
#   41–72s   回到大地圖
#
# ⚠ **剪點要對齊影片的段落**（tools/promo_video.sh 的段落表）。
# 對不上時音樂的轉折會落在畫面中間，說不出哪裡怪，只覺得鬆散。
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${1:-workplace/promo/score/original-72.wav}"
AUDIO="$REPO_ROOT/workplace/audio"
for f in bgm-2 bgm-9; do
    [ -f "$AUDIO/$f.ogg" ] || { echo "找不到 $AUDIO/$f.ogg —— 先跑 tools/bgm2ogg.sh" >&2; exit 1; }
done
mkdir -p "$REPO_ROOT/$(dirname "$OUT")"

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -u "$(id -u):$(id -g)" -e HOME=/tmp \
    -v "$AUDIO":/audio:ro -v "$REPO_ROOT/$(dirname "$OUT")":/out \
    "${WOLONG_VIDEO_IMAGE:-u5cht/video:latest}" \
    ffmpeg -hide_banner -loglevel error -y \
      -i /audio/bgm-2.ogg -i /audio/bgm-9.ogg \
      -filter_complex "\
[0:a]atrim=0:19.5,asetpts=N/SR/TB[a];\
[1:a]atrim=0:23.5,asetpts=N/SR/TB[b];\
[0:a]atrim=19:51,asetpts=N/SR/TB[c];\
[a][b]acrossfade=d=1:c1=tri:c2=tri[ab];\
[ab][c]acrossfade=d=1:c1=tri:c2=tri[abc];\
[abc]atrim=0:72,afade=t=in:st=0:d=1,afade=t=out:st=70:d=2,\
loudnorm=I=-18:TP=-1.5:LRA=11,aresample=44100[out]" \
      -map "[out]" -ac 2 -ar 44100 -c:a pcm_s16le "/out/$(basename "$OUT")"

echo "原版配樂音軌 → $OUT"
