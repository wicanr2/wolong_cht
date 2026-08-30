#!/usr/bin/env bash
# 把三支現行推廣片接成**一支**發行用推廣片，配樂全片統一鋪原版曲子。
#
#   tools/promo_combined.sh [輸出.mp4]     # 預設 dist/promo/wolong-remake-promo.mp4
#
# ⚠ **輸出要寫到 `dist/promo/`，不是 `dist-all/promo/`。**
# `release_all_fs.py` 的 `promo_source()` 先找 `dist/promo/`，只寫後者的話
# 下一次打包會用舊片蓋掉新片，而 manifest 與雜湊全部正常
#（docs/promo/README.md 的同一條警告）。
#
# ## 為什麼要剪掉內部的卡
#
# 三支各自有開場卡與結尾卡。直接串接的話，主預告的結尾卡（65–72）緊接著
# DOS/V 對照片的開場卡（00–04），中段會出現十一秒的卡——看起來像三支片
# 黏在一起，而不是一支片。所以照各片**已經記錄在文件裡的分鏡**裁掉內部的卡：
#
#   主預告      0–65     去掉自己的結尾卡（65–72）；0–65 全是內容
#   DOS/V 對照  4–67     去掉開場卡（00–04）與結尾卡（67–72）
#   Android     0–48.468 完整保留：它的標題卡當第三章的章名卡，
#                        結尾卡「原版資料請自備」當全片的收尾
#
# 合計 65 ＋ 63 ＋ 48.468 ＝ **176.468 秒**。
# 剪點全部出自 docs/promo/README.md、dosv-realmachine.md 的分鏡表，
# 不是目測——目測的剪點會把內容切掉半句，而且下一次重剪要再猜一次。
#
# ## 配樂
#
# ⭐ 聲音來源是 remake 自己的 OPL3 合成（`tools/bgm2ogg.sh` 從使用者自備的
# `BGM.DAT` 算出來的，docs/spec/29）——**不是從原版錄音取樣**。
# 但它演奏的是原版的曲子，所以這條音軌是**原版衍生物**：
# 只用在推廣片，不進任何遊戲包，也不宣稱是本專案的創作。
#
# ⚠ **三支原本的音軌全部換掉**，包含 DOS/V 對照片那條 DOSBox-X 實錄的
# 原版 AdLib。使用者裁定 2026-08-30：全片統一用原版 ogg。
#
# 曲目照場景挑（docs/re/58 的曲表：2 春／3 夏／4 秋／5 冬 是大地圖，
# 9 平原野戰、10 山地林地水域 是戰場）：
#
#     0–65     沿用 original-72.wav 的前 65 秒（曲 2 春 → 曲 9 → 曲 2）
#              ——它已經對齊主預告自己的分鏡，重接反而會錯開
#     65–105   曲 5（冬・大地圖）  DOS/V 段的政略部分（開新遊戲、編成、行軍）
#     105–128  曲 10（戰場）       DOS/V 段從片內 44 秒起全是戰鬥
#     128–176  曲 3（夏・大地圖）  手機段
#
# 交界用 1 秒 acrossfade。**acrossfade 會把兩段重疊 1 秒**，所以每一段要
# 多剪 1 秒去補——不補的話總長會少 3 秒，而症狀是最後一段被切掉尾巴。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"
OUT="${1:-dist/promo/wolong-remake-promo.mp4}"
SRC="dist/promo"
AUDIO="workplace/audio"
SCORE_DIR="workplace/promo/score"
SCORE="$SCORE_DIR/combined-176.wav"
IMAGE="${WOLONG_VIDEO_IMAGE:-u5cht/video:latest}"

for f in wolong-remake-trailer.mp4 wolong-remake-dosv-realmachine.mp4 \
         wolong-remake-android.mp4; do
    [ -f "$SRC/$f" ] || { echo "找不到 $SRC/$f" >&2; exit 1; }
done
for t in 2 3 5 9 10; do
    [ -f "$AUDIO/bgm-$t.ogg" ] || {
        echo "找不到 $AUDIO/bgm-$t.ogg —— 先跑 tools/bgm2ogg.sh" >&2; exit 1; }
done
[ -f "$SCORE_DIR/original-72.wav" ] || {
    echo "找不到 $SCORE_DIR/original-72.wav —— 先跑 tools/promo_score_original.sh" >&2
    exit 1; }

mkdir -p "$SCORE_DIR" "$(dirname "$OUT")"

vrun() {
    docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
        --network none --memory 4g --cpus 2 --pids-limit 256 \
        -u "$(id -u):$(id -g)" -e HOME=/tmp \
        -v "$REPO_ROOT":/src -w /src "$IMAGE" "$@"
}

echo "[1/3] 接配樂 → $SCORE"
vrun ffmpeg -hide_banner -loglevel error -y \
    -i "$SCORE_DIR/original-72.wav" \
    -i "$AUDIO/bgm-5.ogg" -i "$AUDIO/bgm-10.ogg" -i "$AUDIO/bgm-3.ogg" \
    -filter_complex "\
[0:a]atrim=0:65,asetpts=N/SR/TB,aformat=sample_rates=44100:channel_layouts=stereo[p1];\
[1:a]atrim=0:41,asetpts=N/SR/TB,aformat=sample_rates=44100:channel_layouts=stereo[p2];\
[2:a]atrim=0:24,asetpts=N/SR/TB,aformat=sample_rates=44100:channel_layouts=stereo[p3];\
[3:a]atrim=0:50.468,asetpts=N/SR/TB,aformat=sample_rates=44100:channel_layouts=stereo[p4];\
[p1][p2]acrossfade=d=1:c1=tri:c2=tri[x1];\
[x1][p3]acrossfade=d=1:c1=tri:c2=tri[x2];\
[x2][p4]acrossfade=d=1:c1=tri:c2=tri[x3];\
[x3]atrim=0:176.468,afade=t=in:st=0:d=1,afade=t=out:st=174.468:d=2,\
loudnorm=I=-18:TP=-1.5:LRA=11[out]" \
    -map "[out]" -ar 44100 -ac 2 "$SCORE"

echo "[2/3] 接畫面並套上配樂 → $OUT"
vrun ffmpeg -hide_banner -loglevel error -y \
    -i "$SRC/wolong-remake-trailer.mp4" \
    -i "$SRC/wolong-remake-dosv-realmachine.mp4" \
    -i "$SRC/wolong-remake-android.mp4" \
    -i "$SCORE" \
    -filter_complex "\
[0:v]trim=0:65,setpts=PTS-STARTPTS,scale=1280:720,setsar=1,fps=30[v0];\
[1:v]trim=4:67,setpts=PTS-STARTPTS,scale=1280:720,setsar=1,fps=30[v1];\
[2:v]trim=0:48.468,setpts=PTS-STARTPTS,scale=1280:720,setsar=1,fps=30[v2];\
[v0][v1][v2]concat=n=3:v=1:a=0[v]" \
    -map "[v]" -map 3:a \
    -c:v libx264 -preset medium -crf 18 -pix_fmt yuv420p \
    -c:a aac -b:a 192k -ar 44100 -ac 2 -movflags +faststart -shortest "$OUT"

echo "[3/3] 驗媒體規格"
vrun ffprobe -v error -show_entries format=duration \
    -show_entries stream=codec_name,width,height,r_frame_rate,sample_rate,channels \
    -of default=nw=1 "$OUT"
sha256sum "$OUT"
