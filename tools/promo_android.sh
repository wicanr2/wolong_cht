#!/usr/bin/env bash
# 產生 Android 版推廣片。
#
#   tools/phone_capture.sh                 # 先錄 1200 張畫面
#   tools/promo_android.sh [輸出.mp4]
#
# 素材只有 **remake 自己的畫面**與本專案原創的合成配樂（`tools/promo_score.py`）；
# 不讀原版影片、BGM 或 `SOUND.DAT`。
#
# 段落由 `marks.txt` 切——腳本與剪接看同一份標記，改了時間軸不必改這裡。
#
# ⚠ **標記在主機端解析、容器一次只編一段**。把整段剪接塞進 `docker run`
# 的 `bash -c` 字串裡踩過兩個坑：巢狀引號會在改動時安靜地壞掉，
# 而 `ffmpeg` 在 `while read` 迴圈裡會把還沒讀的標記行**吃掉**
#（症狀是只切出兩段，內容還是錯的）。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRAMES="${WOLONG_PROMO_FRAMES:-$REPO_ROOT/workplace/promo/android-frames}"
OUT="${1:-$REPO_ROOT/dist-all/promo/wolong-remake-android.mp4}"
WORK="$REPO_ROOT/workplace/promo/work"
IMAGE="${WOLONG_VIDEO_IMAGE:-demonwinter-video}"
FONT=/usr/share/fonts/opentype/noto/NotoSerifCJK-Bold.ttc

[ -f "$FRAMES/marks.txt" ] || { echo "找不到 $FRAMES/marks.txt，先跑 tools/phone_capture.sh" >&2; exit 1; }
mkdir -p "$(dirname "$OUT")" "$WORK"
rm -f "$WORK"/*.mp4 "$WORK"/*.txt

# 配樂：原創合成，60 秒。影片比它短，混音時裁掉並淡出。
# ⚠ `py.sh` 把 repo 掛在 /src，路徑要用 **repo 相對**的寫法。
tools/py.sh tools/promo_score.py workplace/promo/work/score.wav >/dev/null

vid() {
    docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
        --network none --memory 4g --cpus 2 --pids-limit 256 \
        -v "$FRAMES:/frames:ro" -v "$WORK:/work" \
        -u "$(id -u):$(id -g)" -e HOME=/tmp -w /work "$IMAGE" "$@"
}

fade_out_at() { awk -v d="$1" 'BEGIN { printf "%.2f", d - 0.3 }'; }

card() {
    local name=$1 dur=$2 title=$3 sub=$4
    printf '%s' "$title" >"$WORK/$name-t.txt"
    printf '%s' "$sub" >"$WORK/$name-s.txt"
    vid ffmpeg -nostdin -hide_banner -loglevel error -y \
        -f lavfi -i "color=c=0x080a12:s=1280x720:r=30:d=$dur" \
        -vf "drawtext=fontfile=$FONT:textfile=/work/$name-t.txt:fontcolor=white:fontsize=64:x=(w-text_w)/2:y=290,drawtext=fontfile=$FONT:textfile=/work/$name-s.txt:fontcolor=0xffd666:fontsize=28:x=(w-text_w)/2:y=390,fade=t=in:st=0:d=0.3,fade=t=out:st=$(fade_out_at "$dur"):d=0.3" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p "/work/seg-$name.mp4"
}

# 每一段的標題。key 與 `mobile/wolong/demo.go` 的標記一致。
label_of() {
    case "$1" in
    map) echo "大地圖：可縮放、可拖曳" ;;
    city) echo "點一下據點：歸屬、生產力、防災、城兵" ;;
    list) echo "一覽：武將／據點／勢力／軍團" ;;
    corps) echo "軍團編成：六個位置各自挑兵種" ;;
    advise) echo "進言：玩家是軍師，指令要先過君主那一關" ;;
    battle) echo "戰場：45 度視角，六格編成與六個命令" ;;
    *) echo "" ;;
    esac
}

segment() {
    local idx=$1 key=$2 from=$3 to=$4
    local n=$((to - from))
    [ "$n" -le 5 ] && return 0
    local name
    name=$(printf 'seg-%02d-%s' "$idx" "$key")
    printf '%s' "$(label_of "$key")" >"$WORK/$key-l.txt"
    local dur
    dur=$(awk -v n="$n" 'BEGIN { printf "%.2f", n / 30.0 }')
    # `flags=neighbor`：**點陣圖放大不能內插**，糊掉就看不出是原版美術了。
    vid ffmpeg -nostdin -hide_banner -loglevel error -y \
        -framerate 30 -start_number "$from" -i /frames/f%05d.png -frames:v "$n" \
        -vf "scale=1280:720:flags=neighbor,drawtext=fontfile=$FONT:textfile=/work/$key-l.txt:fontcolor=white:fontsize=26:box=1:boxcolor=0x080a12@0.78:boxborderw=14:x=36:y=h-84,fade=t=in:st=0:d=0.25,fade=t=out:st=$(fade_out_at "$dur"):d=0.25" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p "/work/$name.mp4"
}

card 00-title 3.5 "臥龍傳 Remake" "Android・與桌面同一套規則層"

# marks.txt：`標籤 圖號`，最後一行是 `total 張數`。
keys=()
marks=()
while read -r key frame; do
    keys+=("$key")
    marks+=("$frame")
done <"$FRAMES/marks.txt"

idx=1
for i in "${!keys[@]}"; do
    [ "${keys[$i]}" = total ] && break
    segment "$idx" "${keys[$i]}" "${marks[$i]}" "${marks[$((i + 1))]}"
    idx=$((idx + 1))
done

card 99-end 5 "原版資料請自備" "不散布原版執行檔、資料、美術與字型"

: >"$WORK/concat.txt"
for f in "$WORK"/seg-*.mp4; do
    printf "file '/work/%s'\n" "$(basename "$f")" >>"$WORK/concat.txt"
done

vid ffmpeg -nostdin -hide_banner -loglevel error -y \
    -f concat -safe 0 -i /work/concat.txt -c copy /work/video.mp4
DUR=$(vid ffprobe -v error -show_entries format=duration -of csv=p=0 /work/video.mp4 | tr -d '\r')
vid ffmpeg -nostdin -hide_banner -loglevel error -y \
    -i /work/video.mp4 -i /work/score.wav -map 0:v:0 -map 1:a:0 -t "$DUR" \
    -af "afade=t=out:st=$(awk -v d="$DUR" 'BEGIN { printf "%.2f", d - 2.0 }'):d=2.0" \
    -c:v copy -c:a aac -b:a 160k -ar 44100 -movflags +faststart /work/out.mp4

cp "$WORK/out.mp4" "$OUT"
vid ffprobe -v error -show_entries format=duration:stream=codec_name,width,height \
    -of default=noprint_wrappers=1 /work/out.mp4
echo "推廣片：$OUT"
