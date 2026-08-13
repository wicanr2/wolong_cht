#!/bin/sh
# 產生研究用的 YouTube 原版／remake 對照片與自然畫面差異圖。
# 只使用已保存的 YouTube 代表幀，不把原版影片或原版資產放入發行包。
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
out=${1:-$repo_dir/dist/promo/wolong-remake-yt-comparison.mp4}
out_dir=$(dirname -- "$out")
review_dir=${COMPARE_REVIEW_DIR:-$repo_dir/docs/promo}
font=${COMPARE_FONTFILE:-/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf}

if [ ! -f "$font" ]; then
    echo "找不到比較片字型：$font" >&2
    exit 1
fi

mkdir -p "$out_dir" "$review_dir"
tmp=$(mktemp -d /tmp/wolong-yt-compare-XXXXXX)
cleanup() {
    rm -rf "$tmp"
}
trap cleanup EXIT INT TERM

make_panel() {
    name=$1
    original=$2
    remake=$3
    left_label=$4
    right_label=$5
    left_file="$tmp/$name-left.txt"
    right_file="$tmp/$name-right.txt"
    printf '%s' "$left_label" >"$left_file"
    printf '%s' "$right_label" >"$right_file"
    ffmpeg -hide_banner -loglevel error -y \
        -loop 1 -framerate 30 -t 4 -i "$repo_dir/$original" \
        -loop 1 -framerate 30 -t 4 -i "$repo_dir/$remake" \
        -filter_complex "[0:v]scale=640:400:force_original_aspect_ratio=decrease,pad=640:400:(ow-iw)/2:(oh-ih)/2:color=0x07142e,drawtext=fontfile=$font:textfile=$left_file:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.82:boxborderw=10:x=18:y=16[left];[1:v]scale=640:400:force_original_aspect_ratio=decrease,pad=640:400:(ow-iw)/2:(oh-ih)/2:color=0x07142e,drawtext=fontfile=$font:textfile=$right_file:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.82:boxborderw=10:x=18:y=16[right];[left][right]hstack=inputs=2,format=yuv420p" \
        -an -c:v libx264 -preset medium -crf 19 -r 30 \
        "$tmp/$name.mp4"
}

make_panel natural \
    docs/images/yt-wolong-natural-80s-640x400.png \
    docs/images/wlgame-dosv-natural-remake-skeleton.png \
    "YT ORIGINAL 80s" "REMAKE NATURAL"
make_panel choice \
    docs/images/yt-wolong-natural-160s.png \
    docs/images/wlgame-event3-choice.png \
    "YT ORIGINAL 160s" "REMAKE TALK CHOICE"
make_panel battle \
    docs/images/yt-wolong-natural-240s.png \
    docs/images/wlgame-ai-battle-afterpatch.png \
    "YT ORIGINAL 240s" "REMAKE TACTICAL"
make_panel attack \
    docs/images/yt-wolong-natural-320s.png \
    docs/images/wlgame-ai-battle-attack-afterpatch.png \
    "YT ORIGINAL 320s" "REMAKE PROJECTILE"
make_panel result \
    docs/images/yt-wolong-natural-400s.png \
    docs/images/wlgame-ai-battle-result.png \
    "YT ORIGINAL 400s" "REMAKE RESULT"
make_panel postbattle \
    docs/images/yt-wolong-natural-480s.png \
    docs/images/wlgame-ai-postbattle.png \
    "YT ORIGINAL 480s" "REMAKE POST-BATTLE"

concat="$tmp/concat.txt"
for name in natural choice battle attack result postbattle; do
    printf "file '%s'\n" "$tmp/$name.mp4" >>"$concat"
done

ffmpeg -hide_banner -loglevel error -y \
    -f concat -safe 0 -i "$concat" \
    -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
    -movflags +faststart "$out"

yt="$repo_dir/docs/images/yt-wolong-natural-80s-640x400.png"
remake="$repo_dir/docs/images/wlgame-dosv-natural-remake-skeleton.png"
convert "$yt" "$remake" -alpha off -compose difference -composite \
    "$review_dir/yt-remake-natural-difference.png"
convert "$yt" "$remake" -alpha off +append \
    "$review_dir/yt-remake-natural-side-by-side.png"

echo "比較片：$out"
echo "自然畫面 raw metrics（不同日期／鏡頭，不能當同狀態 parity）："
printf '  AE='; compare -metric AE "$yt" "$remake" null: 2>&1 || true
printf '  RMSE='; compare -metric RMSE "$yt" "$remake" null: 2>&1 || true
echo "差異圖：$review_dir/yt-remake-natural-difference.png"
echo "並排圖：$review_dir/yt-remake-natural-side-by-side.png"

ffprobe -v error \
    -show_entries format=duration:stream=codec_name,width,height,r_frame_rate \
    -of default=noprint_wrappers=1 "$out"
