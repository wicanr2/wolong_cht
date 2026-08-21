#!/bin/sh
# 產生 60 秒 remake 推廣片。此腳本只使用 remake 自己的畫面與本專案原創合成聲，
# 應在 Docker 的 ffmpeg 容器內執行；不讀取原版影片、BGM 或 SOUND.DAT。
#
# ⭐ **大地圖與兩場戰鬥是真的在跑的遊戲**，不是截圖。素材先由
# `tools/promo_live_capture.sh` 逐幀錄出來（docs/spec/71），這裡只負責剪。
# 選單類畫面（事件視窗、校訂、存檔）本來就不動，仍用截圖。
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_dir"
out=${1:-dist/promo/wolong-remake-trailer.mp4}
font=${PROMO_FONTFILE:-/fonts/truetype/moe/MoeStandardSong.ttf}

if [ ! -f "$font" ]; then
    echo "找不到推廣片字型：$font" >&2
    exit 1
fi

mkdir -p "$(dirname -- "$out")"
tmp=$(mktemp -d /tmp/wolong-promo-XXXXXX)
cleanup() {
    rm -rf "$tmp"
}
trap cleanup EXIT INT TERM

write_text() {
    path=$1
    text=$2
    printf '%s' "$text" >"$path"
}

make_card() {
    name=$1
    duration=$2
    title=$3
    subtitle=$4
    title_file="$tmp/$name-title.txt"
    subtitle_file="$tmp/$name-subtitle.txt"
    write_text "$title_file" "$title"
    write_text "$subtitle_file" "$subtitle"
    ffmpeg -hide_banner -loglevel error -y \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$duration" \
        -vf "drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=72:x=(w-text_w)/2:y=285,drawtext=fontfile=$font:textfile=$subtitle_file:fontcolor=0xb8c8e8:fontsize=30:x=(w-text_w)/2:y=395,fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

make_image() {
    name=$1
    duration=$2
    source=$3
    label=$4
    label_file="$tmp/$name-label.txt"
    write_text "$label_file" "$label"
    ffmpeg -hide_banner -loglevel error -y \
        -loop 1 -framerate 30 -t "$duration" -i "$repo_dir/$source" \
        -vf "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2:color=0x07142e,drawtext=fontfile=$font:textfile=$label_file:fontcolor=white:fontsize=28:box=1:boxcolor=0x061128@0.82:boxborderw=12:x=30:y=28,fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

# make_frames 把一段逐幀錄下來的 PNG 接成影片。
#
# ⚠ `-start_number 0` 不能省：檔名從 f00000 開始，ffmpeg 預設從 1 找起，
# 找不到就整段落空——而落空的症狀是「少了一段」不是「報錯」。
# ⚠ `flags=neighbor`：點陣圖放大不能內插，糊掉就看不出是原版美術了。
make_frames() {
    name=$1
    duration=$2
    dir=$3
    label=$4
    label_file="$tmp/$name-label.txt"
    write_text "$label_file" "$label"
    ffmpeg -hide_banner -loglevel error -y \
        -framerate 30 -start_number 0 -i "$repo_dir/$dir/f%05d.png" -t "$duration" \
        -vf "scale=1280:720:flags=neighbor,drawtext=fontfile=$font:textfile=$label_file:fontcolor=white:fontsize=28:box=1:boxcolor=0x061128@0.82:boxborderw=12:x=30:y=28,fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

live=${WOLONG_PROMO_LIVE:-workplace/promo/live}
for seg in map battle siege; do
    if [ ! -f "$repo_dir/$live/$seg/f00000.png" ]; then
        echo "找不到動態素材 $live/$seg —— 先跑 tools/promo_live_capture.sh" >&2
        exit 1
    fi
done

make_card title 4 "臥龍傳 Remake" "DOS/V 對齊・繁中・跨平台"
make_frames natural 11 "$live/map" "大地圖：時鐘在走，兩軍在行軍"
make_image choice 2.5 docs/images/wlgame-event3-choice.png "事件 2–5 TALK 與數值選取"
make_image event3 2 docs/images/event3-44.png "事件 3 停戰金額結果"
make_image event5 2 docs/images/event5-339.png "事件 5 外交官撥款"
make_frames battle 13 "$live/battle" "戰場：45 度視角，兩軍實時交戰"
make_frames siege 9 "$live/siege" "攻城：城牆、城門與守軍"
make_image result 2.5 docs/images/wlgame-ai-battle-result.png "戰果與狀態回寫"
make_image m7 2.5 docs/images/m7-review-321.png "M7 繁中校訂與硬換行"
make_image save 2.5 docs/images/wlgame-save-ui.png "四槽存檔與重播"
make_card end 6.5 "立即開始你的三國" "Linux・Windows・macOS・Android"

concat="$tmp/concat.txt"
: >"$concat"
for name in title natural choice event3 event5 battle siege result m7 save end; do
    printf "file '%s'\n" "$tmp/$name.mp4" >>"$concat"
done

# 配樂：本專案原創合成。ffmpeg 容器裡沒有 python3，所以允許先在
# `tools/py.sh` 那邊產好再用 PROMO_SCORE 指進來。
if [ -n "${PROMO_SCORE:-}" ] && [ -f "$repo_dir/$PROMO_SCORE" ]; then
    cp "$repo_dir/$PROMO_SCORE" "$tmp/score.wav"
else
    python3 "$repo_dir/tools/promo_score.py" "$tmp/score.wav" >/dev/null
fi
ffmpeg -hide_banner -loglevel error -y \
    -f concat -safe 0 -i "$concat" -i "$tmp/score.wav" \
    -map 0:v:0 -map 1:a:0 -t 60 \
    -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
    -c:a aac -b:a 160k -ar 44100 -movflags +faststart "$out"

ffprobe -v error -show_entries format=duration:stream=codec_name,width,height,channels,sample_rate \
    -of default=noprint_wrappers=1 "$out"
echo "推廣片：$out"
