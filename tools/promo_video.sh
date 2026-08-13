#!/bin/sh
# 產生 60 秒 remake 推廣片。此腳本只使用 remake 截圖與本專案原創合成聲，
# 應在 Docker 的 ffmpeg 容器內執行；不讀取原版影片、BGM 或 SOUND.DAT。
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

make_card title 4 "臥龍傳 Remake" "DOS/V 對齊・繁中・跨平台"
make_image natural 6 docs/images/wlgame-dosv-natural-remake-skeleton.png "自然策略畫面"
make_image choice 4 docs/images/wlgame-event3-choice.png "事件 2–5 TALK 與數值選取"
make_image event2 3 docs/images/event2-47.png "事件 2 協力結果"
make_image event3 3 docs/images/event3-44.png "事件 3 停戰金額結果"
make_image event4 3 docs/images/event4-284.png "事件 4 內政官撥款"
make_image event5 3 docs/images/event5-339.png "事件 5 外交官撥款"
make_image battle 5 docs/images/wlgame-ai-battle-afterpatch.png "戰術戰鬥"
make_image attack 5 docs/images/wlgame-ai-battle-attack-afterpatch.png "戰場攻擊與投射物"
make_image result 4 docs/images/wlgame-ai-battle-result.png "戰果與狀態回寫"
make_image postbattle 4 docs/images/wlgame-ai-postbattle.png "戰後回到戰略畫面"
make_image event9 3 docs/images/event9-37.png "事件 9 釋放通知"
make_image m7 3 docs/images/m7-review-321.png "M7 繁中校訂與硬換行"
make_image save 4 docs/images/wlgame-save-ui.png "四槽存檔與重播"
make_card end 6 "立即開始你的三國" "Linux・Windows・macOS 候選包"

concat="$tmp/concat.txt"
: >"$concat"
for name in title natural choice event2 event3 event4 event5 battle attack result postbattle event9 m7 save end; do
    printf "file '%s'\n" "$tmp/$name.mp4" >>"$concat"
done

python3 "$repo_dir/tools/promo_score.py" "$tmp/score.wav" >/dev/null
ffmpeg -hide_banner -loglevel error -y \
    -f concat -safe 0 -i "$concat" -i "$tmp/score.wav" \
    -map 0:v:0 -map 1:a:0 -t 60 \
    -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
    -c:a aac -b:a 160k -ar 44100 -movflags +faststart "$out"

ffprobe -v error -show_entries format=duration:stream=codec_name,width,height,channels,sample_rate \
    -of default=noprint_wrappers=1 "$out"
echo "推廣片：$out"
